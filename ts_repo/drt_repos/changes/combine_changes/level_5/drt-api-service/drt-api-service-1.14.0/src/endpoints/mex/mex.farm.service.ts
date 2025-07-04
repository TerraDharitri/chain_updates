import { Constants } from "@terradharitri/sdk-nestjs-common";
import { CacheService } from "@terradharitri/sdk-nestjs-cache";
import { forwardRef, Inject, Injectable } from "@nestjs/common";
import { QueryPagination } from "src/common/entities/query.pagination";
import { CacheInfo } from "src/utils/cache.info";
import { GraphQlService } from "src/common/graphql/graphql.service";
import { MoaFarm } from "./entities/moa.farm";
import { MoaTokenService } from "./moa.token.service";
import { MoaStakingProxy } from "./entities/moa.staking.proxy";
import { ApiConfigService } from "src/common/api-config/api.config.service";
import { farmsQuery } from "./graphql/farms.query";
import { stakingProxyQuery } from "./graphql/staking.proxy.query";

@Injectable()
export class MoaFarmService {
  constructor(
    private readonly cachingService: CacheService,
    private readonly graphQlService: GraphQlService,
    @Inject(forwardRef(() => MoaTokenService))
    private readonly moaTokenService: MoaTokenService,
    private readonly apiConfigService: ApiConfigService,
  ) { }

  async refreshMoaFarms(): Promise<void> {
    const farms = await this.getAllMoaFarmsRaw();
    await this.cachingService.setRemote(CacheInfo.MoaFarms.key, farms, CacheInfo.MoaFarms.ttl);
    this.cachingService.setLocal(CacheInfo.MoaFarms.key, farms, Constants.oneSecond() * 30);
  }

  async getMoaFarms(pagination: QueryPagination): Promise<MoaFarm[]> {
    const moaFarms = await this.getAllMoaFarms();
    const { from, size } = pagination;

    return moaFarms.slice(from, from + size);
  }

  async getAllMoaFarms(): Promise<MoaFarm[]> {
    if (!this.apiConfigService.isExchangeEnabled()) {
      return [];
    }

    return await this.cachingService.getOrSet(
      CacheInfo.MoaFarms.key,
      async () => await this.getAllMoaFarmsRaw(),
      CacheInfo.MoaFarms.ttl,
      Constants.oneSecond() * 30,
    );
  }

  async getMoaFarmsCount(): Promise<number> {
    const moaFarms = await this.getAllMoaFarms();

    return moaFarms.length;
  }

  private async getAllMoaFarmsRaw(): Promise<MoaFarm[]> {
    const response: any = await this.graphQlService.getExchangeServiceData(farmsQuery, {});
    if (!response) {
      return [];
    }

    const pairs = await this.moaTokenService.getIndexedMoaTokens();

    const farms = response.farms.map((farmResponse: any) => MoaFarm.fromFarmQueryResponse(farmResponse));

    const stakingFarms = response.stakingFarms.map((stakingFarm: any) => MoaFarm.fromStakingFarmResponse(stakingFarm, pairs));

    return [...farms, ...stakingFarms];
  }

  async getAllStakingProxies(): Promise<MoaStakingProxy[]> {
    if (!this.apiConfigService.isExchangeEnabled()) {
      return [];
    }

    return await this.cachingService.getOrSet(
      CacheInfo.StakingProxies.key,
      async () => await this.getAllStakingProxiesRaw(),
      CacheInfo.StakingProxies.ttl,
      Constants.oneSecond() * 30,
    );
  }

  private async getAllStakingProxiesRaw(): Promise<MoaStakingProxy[]> {
    const response: any = await this.graphQlService.getExchangeServiceData(stakingProxyQuery, {});
    if (!response) {
      return [];
    }

    const stakingProxies = response.stakingProxies.map((stakingProxyRaw: any) => MoaStakingProxy.fromQueryResponse(stakingProxyRaw));
    return stakingProxies;
  }
}
