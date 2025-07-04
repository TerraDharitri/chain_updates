import { BadRequestException, forwardRef, Inject, Injectable } from "@nestjs/common";
import { CacheInfo } from "src/utils/cache.info";
import { MoaToken } from "./entities/moa.token";
import { MoaPairService } from "./moa.pair.service";
import { MoaPair } from "./entities/moa.pair";
import { ApiConfigService } from "src/common/api-config/api.config.service";
import { MoaFarmService } from "./moa.farm.service";
import { MoaSettingsService } from "./moa.settings.service";
import { Constants } from "@terradharitri/sdk-nestjs-common";
import { CacheService } from "@terradharitri/sdk-nestjs-cache";
import { OriginLogger } from "@terradharitri/sdk-nestjs-common";
import { QueryPagination } from "src/common/entities/query.pagination";
import { MoaTokenType } from "./entities/moa.token.type";
import { GraphQlService } from "src/common/graphql/graphql.service";
import { tokensQuery } from "./graphql/tokens.query";

@Injectable()
export class MoaTokenService {
  private readonly logger = new OriginLogger(MoaTokenService.name);

  constructor(
    private readonly cachingService: CacheService,
    private readonly apiConfigService: ApiConfigService,
    private readonly drtPairService: MoaPairService,
    @Inject(forwardRef(() => MoaFarmService))
    private readonly moaFarmService: MoaFarmService,
    private readonly moaSettingsService: MoaSettingsService,
    private readonly graphQlService: GraphQlService,
  ) { }

  async refreshMoaTokens(): Promise<void> {
    const tokens = await this.getAllMoaTokensRaw();
    await this.cachingService.setRemote(CacheInfo.MoaTokens.key, tokens, CacheInfo.MoaTokens.ttl);
    this.cachingService.setLocal(CacheInfo.MoaTokens.key, tokens, Constants.oneSecond() * 30);

    const tokenTypes = await this.getAllMoaTokenTypesRaw();
    await this.cachingService.setRemote(CacheInfo.MoaTokenTypes.key, tokenTypes, CacheInfo.MoaTokenTypes.ttl);
    this.cachingService.setLocal(CacheInfo.MoaTokenTypes.key, tokenTypes, Constants.oneSecond() * 30);

    const indexedTokens = await this.getIndexedMoaTokensRaw();
    await this.cachingService.setRemote(CacheInfo.MoaTokensIndexed.key, indexedTokens, CacheInfo.MoaTokensIndexed.ttl);
    this.cachingService.setLocal(CacheInfo.MoaTokensIndexed.key, indexedTokens, Constants.oneSecond() * 30);

    const indexedPrices = await this.getMoaPricesRaw();
    await this.cachingService.setRemote(CacheInfo.MoaPrices.key, indexedPrices, CacheInfo.MoaPrices.ttl);
    this.cachingService.setLocal(CacheInfo.MoaPrices.key, indexedPrices, Constants.oneSecond() * 30);
  }

  async getMoaTokens(queryPagination: QueryPagination): Promise<MoaToken[]> {
    const { from, size } = queryPagination;
    let allMoaTokens = await this.getAllMoaTokens();
    allMoaTokens = JSON.parse(JSON.stringify(allMoaTokens));

    return allMoaTokens.slice(from, from + size);
  }

  async getMoaTokenByIdentifier(identifier: string): Promise<MoaToken | undefined> {
    const moaTokens = await this.getAllMoaTokens();
    return moaTokens.find(x => x.id === identifier);
  }

  async getMoaPrices(): Promise<Record<string, { price: number, isToken: boolean }>> {
    return await this.cachingService.getOrSet(
      CacheInfo.MoaPrices.key,
      async () => await this.getMoaPricesRaw(),
      CacheInfo.MoaPrices.ttl,
      Constants.oneSecond() * 30
    );
  }

  async getMoaPricesRaw(): Promise<Record<string, { price: number, isToken: boolean }>> {
    try {
      const result: Record<string, { price: number, isToken: boolean }> = {};

      const tokens = await this.getAllMoaTokens();
      for (const token of tokens) {
        result[token.id] = {
          price: token.price,
          isToken: true,
        };
      }

      const pairs = await this.drtPairService.getAllMoaPairs();
      for (const pair of pairs) {
        result[pair.id] = {
          price: pair.price,
          isToken: false,
        };
      }

      const farms = await this.moaFarmService.getAllMoaFarms();
      for (const farm of farms) {
        result[farm.id] = {
          price: farm.price,
          isToken: false,
        };
      }

      const settings = await this.moaSettingsService.getSettings();
      if (settings) {
        const moaToken = tokens.find(x => x.symbol === 'MOA');
        if (moaToken) {
          const lkmoaIdentifier = settings.lockedAssetIdentifier;
          if (lkmoaIdentifier) {
            result[lkmoaIdentifier] = {
              price: moaToken.price,
              isToken: false,
            };
          }

          const xmoaIdentifier = settings.lockedAssetIdentifierV2;
          if (xmoaIdentifier) {
            result[xmoaIdentifier] = {
              price: moaToken.price,
              isToken: false,
            };
          }
        }
      }

      return result;
    } catch (error) {
      this.logger.error('An error occurred while fetching moa prices');
      this.logger.error(error);
      return {};
    }
  }

  async getIndexedMoaTokens(): Promise<Record<string, MoaToken>> {
    if (!this.apiConfigService.getExchangeServiceUrl()) {
      return {};
    }

    return await this.cachingService.getOrSet(
      CacheInfo.MoaTokensIndexed.key,
      async () => await this.getIndexedMoaTokensRaw(),
      CacheInfo.MoaTokensIndexed.ttl,
      Constants.oneSecond() * 30
    );
  }

  async getIndexedMoaTokensRaw(): Promise<Record<string, MoaToken>> {
    const result: Record<string, MoaToken> = {};

    const tokens = await this.getAllMoaTokens();
    for (const token of tokens) {
      result[token.id] = token;
    }

    return result;
  }

  async getMoaTokensCount(): Promise<number> {
    const moaTokens = await this.getAllMoaTokens();

    return moaTokens.length;
  }

  private async getAllMoaTokens(): Promise<MoaToken[]> {
    if (!this.apiConfigService.getExchangeServiceUrl()) {
      return [];
    }

    return await this.cachingService.getOrSet(
      CacheInfo.MoaTokens.key,
      async () => await this.getAllMoaTokensRaw(),
      CacheInfo.MoaTokens.ttl,
      Constants.oneSecond() * 30
    );
  }

  private async getAllMoaTokensRaw(): Promise<MoaToken[]> {
    const pairs = await this.drtPairService.getAllMoaPairs();

    const moaTokens: MoaToken[] = [];
    for (const pair of pairs) {
      if (pair.baseSymbol === 'WREWA' && pair.quoteSymbol === "USDC") {
        const wrewaToken = new MoaToken();
        wrewaToken.id = pair.baseId;
        wrewaToken.symbol = pair.baseSymbol;
        wrewaToken.name = pair.baseName;
        wrewaToken.price = pair.basePrice;
        wrewaToken.previous24hPrice = pair.basePrevious24hPrice;
        wrewaToken.previous24hVolume = pair.volume24h;
        wrewaToken.tradesCount = this.computeTradesCountForMoaToken(wrewaToken, pairs);
        moaTokens.push(wrewaToken);
      }

      const moaToken = this.getMoaToken(pair);
      if (!moaToken) {
        continue;
      }

      moaToken.tradesCount = this.computeTradesCountForMoaToken(moaToken, pairs);

      moaTokens.push(moaToken);
    }

    return moaTokens.distinct(x => x.id);
  }

  private getMoaToken(pair: MoaPair): MoaToken | null {
    if (pair.baseSymbol === 'WREWA' && pair.quoteSymbol === "USDC") {
      return {
        id: pair.quoteId,
        symbol: pair.quoteSymbol,
        name: pair.quoteName,
        price: pair.quotePrice,
        previous24hPrice: pair.quotePrevious24hPrice,
        previous24hVolume: pair.volume24h,
        tradesCount: 0,
      };
    }

    if (['WREWA', 'USDC'].includes(pair.quoteSymbol)) {
      return {
        id: pair.baseId,
        symbol: pair.baseSymbol,
        name: pair.baseName,
        price: pair.basePrice,
        previous24hPrice: pair.basePrevious24hPrice,
        previous24hVolume: pair.volume24h,
        tradesCount: 0,
      };
    }

    if (['WREWA', 'USDC'].includes(pair.baseSymbol)) {
      return {
        id: pair.quoteId,
        symbol: pair.quoteSymbol,
        name: pair.quoteName,
        price: pair.quotePrice,
        previous24hPrice: pair.quotePrevious24hPrice,
        previous24hVolume: pair.volume24h,
        tradesCount: 0,
      };
    }

    return null;
  }

  async getAllMoaTokenTypes(): Promise<MoaTokenType[]> {
    if (!this.apiConfigService.getExchangeServiceUrl()) {
      return [];
    }

    return await this.cachingService.getOrSet(
      CacheInfo.MoaTokenTypes.key,
      async () => await this.getAllMoaTokenTypesRaw(),
      CacheInfo.MoaTokenTypes.ttl,
      Constants.oneSecond() * 30
    );
  }

  private async getAllMoaTokenTypesRaw(): Promise<MoaTokenType[]> {
    try {
      const settings = await this.moaSettingsService.getSettings();
      if (!settings) {
        throw new BadRequestException('Could not fetch MOA tokens');
      }

      const result: any = await this.graphQlService.getExchangeServiceData(tokensQuery);
      if (!result || !result.tokens) {
        return [];
      }

      return result.tokens.map((token: MoaTokenType) => ({
        identifier: token.identifier,
        type: token.type.toLowerCase(),
      }));
    } catch (error) {
      this.logger.error('An error occurred while fetching all moa token types');
      this.logger.error(error);
      return [];
    }
  }

  private computeTradesCountForMoaToken(moaToken: MoaToken, filteredPairs: MoaPair[]): number {
    const pairs = filteredPairs.filter(x => x.baseId === moaToken.id || x.quoteId === moaToken.id);
    const computeResult = pairs.sum(pair => pair.tradesCount ?? 0);
    return computeResult;
  }
}
