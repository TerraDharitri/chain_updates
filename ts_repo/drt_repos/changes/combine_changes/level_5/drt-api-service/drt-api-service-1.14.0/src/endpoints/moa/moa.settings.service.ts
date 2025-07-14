import { Constants } from "@terradharitri/sdk-nestjs-common";
import { CacheService } from "@terradharitri/sdk-nestjs-cache";
import { Injectable } from "@nestjs/common";
import { CacheInfo } from "src/utils/cache.info";
import { GraphQlService } from "src/common/graphql/graphql.service";
import { TransactionMetadata } from "../transactions/transaction-action/entities/transaction.metadata";
import { TransactionMetadataTransfer } from "../transactions/transaction-action/entities/transaction.metadata.transfer";
import { MoaSettings } from "./entities/moa.settings";
import { ApiConfigService } from "src/common/api-config/api.config.service";
import { settingsQuery } from "./graphql/settings.query";
import { pairCountQuery } from "./graphql/pairs.count.query";

@Injectable()
export class MoaSettingsService {
  private wrewaId: string | undefined;

  constructor(
    private readonly cachingService: CacheService,
    private readonly graphQlService: GraphQlService,
    private readonly apiConfigService: ApiConfigService,
  ) { }

  getTransfers(metadata: TransactionMetadata): TransactionMetadataTransfer[] | undefined {
    const transfers = metadata.transfers;
    if (!transfers || transfers.length === 0) {
      return undefined;
    }

    return transfers;
  }

  async isMoaInteraction(metadata: TransactionMetadata): Promise<boolean> {
    const moaContracts = await this.getMoaContracts();
    return moaContracts.has(metadata.receiver);
  }

  async refreshSettings(): Promise<void> {
    const settings = await this.getSettingsRaw();
    await this.cachingService.setRemote(CacheInfo.MoaSettings.key, settings, CacheInfo.MoaSettings.ttl);
    this.cachingService.setLocal(CacheInfo.MoaSettings.key, settings, Constants.oneMinute() * 10);

    const contracts = await this.getMoaContractsRaw();
    await this.cachingService.setRemote(CacheInfo.MoaContracts.key, contracts, CacheInfo.MoaContracts.ttl);
    this.cachingService.setLocal(CacheInfo.MoaContracts.key, contracts, Constants.oneMinute() * 10);
  }

  async getSettings(): Promise<MoaSettings | null> {
    if (!this.apiConfigService.isExchangeEnabled()) {
      return null;
    }

    const settings = await this.cachingService.getOrSet(
      CacheInfo.MoaSettings.key,
      async () => await this.getSettingsRaw(),
      CacheInfo.MoaSettings.ttl,
      Constants.oneMinute() * 10,
    );

    this.wrewaId = settings?.wrewaId;

    return settings;
  }

  async getMoaContracts(): Promise<Set<string>> {
    let contracts = await this.cachingService.getLocal<Set<string>>(CacheInfo.MoaContracts.key);
    if (!contracts) {
      contracts = await this.getMoaContractsRaw();
      this.cachingService.setLocal(CacheInfo.MoaContracts.key, contracts, Constants.oneMinute() * 10);
    }

    return contracts;
  }

  async getMoaContractsRaw(): Promise<Set<string>> {
    const settings = await this.getSettings();
    if (!settings) {
      return new Set<string>();
    }

    return new Set<string>([
      settings.distributionContract,
      settings.lockedAssetContract,
      settings.routerFactoryContract,
      ...settings.farmContracts,
      ...settings.pairContracts,
      ...settings.wrapContracts,
    ]);
  }

  public async getSettingsRaw(): Promise<MoaSettings | null> {
    const pairLimitCount = await this.getPairLimitCount();
    const response = await this.graphQlService.getExchangeServiceData(settingsQuery(pairLimitCount));
    if (!response) {
      return null;
    }

    const transformedResponse = {
      ...response,
      pairs: response.filteredPairs.edges.map((edge: { node: { address: string } }) => ({
        address: edge.node.address,
      })),
    };

    const settings = MoaSettings.fromQueryResponse(transformedResponse);
    return settings;
  }

  getWrewaId(): string | undefined {
    return this.wrewaId;
  }

  private async getPairLimitCount(): Promise<number> {
    const response = await this.graphQlService.getExchangeServiceData(pairCountQuery);
    if (!response) {
      return 500;
    }

    return response.factory.pairCount;
  }
}
