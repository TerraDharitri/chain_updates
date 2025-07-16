import { Injectable } from "@nestjs/common";
import { GraphQlService } from "src/common/graphql/graphql.service";
import { OriginLogger } from "@terradharitri/sdk-nestjs-common";
import { gql } from 'graphql-request';
import { MoaTokenChart } from "./entities/moa.token.chart";
import { MoaTokenService } from "./moa.token.service";
import { CacheService } from "@terradharitri/sdk-nestjs-cache";
import { CacheInfo } from "src/utils/cache.info";
import { tokenPricesHourResolutionQuery } from "./graphql/token.prices.hour.resolution.query";

@Injectable()
export class MoaTokenChartsService {
  private readonly logger = new OriginLogger(MoaTokenChartsService.name);

  constructor(
    private readonly graphQlService: GraphQlService,
    private readonly moaTokenService: MoaTokenService,
    private readonly cachingService: CacheService,
  ) { }

  async getTokenPricesHourResolution(tokenIdentifier: string): Promise<MoaTokenChart[] | undefined> {
    return await this.cachingService.getOrSet(
      CacheInfo.TokenHourChart(tokenIdentifier).key,
      async () => await this.getTokenPricesHourResolutionRaw(tokenIdentifier),
      CacheInfo.TokenHourChart(tokenIdentifier).ttl,
    );
  }

  async getTokenPricesHourResolutionRaw(tokenIdentifier: string): Promise<MoaTokenChart[] | undefined> {
    const isMoaToken = await this.isMoaToken(tokenIdentifier);
    if (!isMoaToken) {
      return undefined;
    }

    try {
      const query = tokenPricesHourResolutionQuery(tokenIdentifier);
      const data = await this.graphQlService.getExchangeServiceData(query);
      return this.convertToMoaTokenChart(data?.values24h) || [];
    } catch (error) {
      this.logger.error(`An error occurred while fetching hourly token prices for ${tokenIdentifier}`, error);
      return [];
    }
  }

  async getTokenPricesDayResolution(tokenIdentifier: string): Promise<MoaTokenChart[] | undefined> {
    return await this.cachingService.getOrSet(
      CacheInfo.TokenDailyChart(tokenIdentifier).key,
      async () => await this.getTokenPricesDayResolutionRaw(tokenIdentifier),
      CacheInfo.TokenDailyChart(tokenIdentifier).ttl,
    );
  }

  async getTokenPricesDayResolutionRaw(tokenIdentifier: string): Promise<MoaTokenChart[] | undefined> {
    const isMoaToken = await this.isMoaToken(tokenIdentifier);
    if (!isMoaToken) {
      return undefined;
    }

    const query = gql`
      query tokenPriceDayResolution {
        latestCompleteValues(
          series: "${tokenIdentifier}",
          metric: "priceUSD",
        ) {
          timestamp
          value
        }
      }
    `;

    try {
      const data = await this.graphQlService.getExchangeServiceData(query);
      return this.convertToMoaTokenChart(data?.latestCompleteValues) || [];
    } catch (error) {
      this.logger.error(`An error occurred while fetching daily token prices for ${tokenIdentifier}`, error);
      return [];
    }
  }

  private convertToMoaTokenChart(data: { timestamp: string; value: string }[]): MoaTokenChart[] {
    return data?.map(item => new MoaTokenChart({
      timestamp: Math.floor(new Date(item.timestamp).getTime() / 1000),
      value: Number(item.value),
    })) || [];
  }

  private async isMoaToken(tokenIdentifier: string): Promise<boolean> {
    const token = await this.moaTokenService.getMoaTokenByIdentifier(tokenIdentifier);
    return token !== undefined;
  }
}
