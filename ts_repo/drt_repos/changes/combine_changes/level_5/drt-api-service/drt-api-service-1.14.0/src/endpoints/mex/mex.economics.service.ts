import { CacheService } from "@terradharitri/sdk-nestjs-cache";
import { BadRequestException, Injectable } from "@nestjs/common";
import { gql } from "graphql-request";
import { CacheInfo } from "src/utils/cache.info";
import { GraphQlService } from "src/common/graphql/graphql.service";
import { MoaSettingsService } from "./moa.settings.service";
import { MoaEconomics } from "./entities/moa.economics";

@Injectable()
export class MoaEconomicsService {
  constructor(
    private readonly moaSettingService: MoaSettingsService,
    private readonly cachingService: CacheService,
    private readonly graphQlService: GraphQlService
  ) { }

  async refreshMoaEconomics() {
    const economics = await this.getMoaEconomicsRaw();
    await this.cachingService.setRemote(CacheInfo.MoaEconomics.key, economics, CacheInfo.MoaEconomics.ttl);
  }

  async getMoaEconomics(): Promise<MoaEconomics> {
    return await this.cachingService.getOrSet(
      CacheInfo.MoaEconomics.key,
      async () => await this.getMoaEconomicsRaw(),
      CacheInfo.MoaEconomics.ttl,
    );
  }

  async getMoaEconomicsRaw(): Promise<MoaEconomics> {
    const settings = await this.moaSettingService.getSettings();
    if (!settings) {
      throw new BadRequestException('Could not fetch MOA settings');
    }

    const variables = {
      "moaID": settings.moaId,
      "days": 7,
    };

    const query = gql`
      query ($days: Int!, $moaID: String!) {
        totalAggregatedRewards(days: $days)
        moaPriceUSD: getTokenPriceUSD(tokenID: $moaID)
        moaSupply: totalTokenSupply(tokenID: $moaID)
        factory {
          totalVolumeUSD24h
          __typename
        }
      }
    `;

    const response: any = await this.graphQlService.getExchangeServiceData(query, variables);
    if (!response) {
      throw new BadRequestException('Could not fetch MOA economics data from MOA microservice');
    }

    const moaEconomics = MoaEconomics.fromQueryResponse(response, settings);
    return moaEconomics;
  }
}
