import { Constants } from '@terradharitri/sdk-nestjs-common';
import { CacheService } from '@terradharitri/sdk-nestjs-cache';
import { BadRequestException, Injectable } from '@nestjs/common';
import { CacheInfo } from 'src/utils/cache.info';
import { GraphQlService } from 'src/common/graphql/graphql.service';
import { MoaPair } from './entities/moa.pair';
import { MoaPairState } from './entities/moa.pair.state';
import { MoaPairType } from './entities/moa.pair.type';
import { MoaSettingsService } from './moa.settings.service';
import { OriginLogger } from '@terradharitri/sdk-nestjs-common';
import { ApiConfigService } from 'src/common/api-config/api.config.service';
import { MoaPairExchange } from './entities/moa.pair.exchange';
import { MoaPairsFilter } from './entities/moa.pairs..filter';
import { MoaPairStatus } from './entities/moa.pair.status';
import { filteredPairsQuery } from './graphql/filtered.pairs.query';

@Injectable()
export class MoaPairService {
  private readonly logger = new OriginLogger(MoaPairService.name);

  constructor(
    private readonly cachingService: CacheService,
    private readonly moaSettingService: MoaSettingsService,
    private readonly graphQlService: GraphQlService,
    private readonly apiConfigService: ApiConfigService,
  ) { }

  async refreshMoaPairs(): Promise<void> {
    const pairs = await this.getAllMoaPairsRaw(false);
    await this.cachingService.setRemote(CacheInfo.MoaPairs.key, pairs, CacheInfo.MoaPairs.ttl);
    this.cachingService.setLocal(CacheInfo.MoaPairs.key, pairs, Constants.oneSecond() * 30);
  }

  async getMoaPairs(from: number, size: number, filter?: MoaPairsFilter): Promise<any> {
    let allMoaPairs = await this.getAllMoaPairs(filter?.includeFarms ?? false);
    allMoaPairs = this.applyFilters(allMoaPairs, filter);

    return allMoaPairs.slice(from, from + size);
  }

  async getMoaPair(baseId: string, quoteId: string, includeFarms: boolean = false): Promise<MoaPair | undefined> {
    const allMoaPairs = await this.getAllMoaPairs(includeFarms);
    return allMoaPairs.find(pair => pair.baseId === baseId && pair.quoteId === quoteId);
  }

  async getAllMoaPairs(includeFarms: boolean = false): Promise<MoaPair[]> {
    if (!this.apiConfigService.isExchangeEnabled()) {
      return [];
    }

    const cacheKey = includeFarms ? CacheInfo.MoaPairsWithFarms.key : CacheInfo.MoaPairs.key;
    const ttl = includeFarms ? CacheInfo.MoaPairsWithFarms.ttl : CacheInfo.MoaPairs.ttl;

    return await this.cachingService.getOrSet(
      cacheKey,
      async () => await this.getAllMoaPairsRaw(includeFarms),
      ttl,
      Constants.oneSecond() * 30,
    );
  }

  async getMoaPairsCount(filter?: MoaPairsFilter): Promise<number> {
    const drtPairs = await this.getAllMoaPairs(filter?.includeFarms ?? false);
    const filteredPairs = this.applyFilters(drtPairs, filter);

    return filteredPairs.length;
  }

  async getAllMoaPairsRaw(includeFarms: boolean = false): Promise<MoaPair[]> {
    try {
      const settings = await this.moaSettingService.getSettings();
      if (!settings) {
        throw new BadRequestException('Could not fetch MOA settings');
      }

      const allPairs: MoaPair[] = [];
      let cursor: string | null = null;
      let hasNextPage = true;

      while (hasNextPage) {
        const variables = {
          pagination: { first: 25, after: cursor },
          filters: { state: [MoaPairStatus.active] },
        };

        const query = filteredPairsQuery(includeFarms);
        const result: any = await this.graphQlService.getExchangeServiceData(query, variables);

        if (!result) {
          break;
        }

        const pairs = result.filteredPairs.edges.map((edge: any) => this.getPairInfo(edge.node, includeFarms));
        allPairs.push(...pairs.filter((pair: MoaPair | undefined) => pair !== undefined));

        hasNextPage = result.filteredPairs.pageInfo.hasNextPage;
        cursor = result.filteredPairs.edges.length > 0 ? result.filteredPairs.edges[result.filteredPairs.edges.length - 1].cursor : null;
      }

      return allPairs;
    } catch (error) {
      this.logger.error('An error occurred while getting all moa pairs from the exchange');
      this.logger.error(error);
      return [];
    }
  }


  private getPairInfo(pair: any, includeFarms: boolean = false): MoaPair | undefined {
    const firstTokenSymbol = pair.firstToken.identifier.split('-')[0];
    const secondTokenSymbol = pair.secondToken.identifier.split('-')[0];
    const state = this.getPairState(pair.state);
    const type = this.getPairType(pair.type);

    if (!type || [MoaPairType.unlisted].includes(type)) {
      return undefined;
    }

    const dharitrixTypes = [
      MoaPairType.core,
      MoaPairType.community,
      MoaPairType.experimental,
      MoaPairType.ecosystem,
    ];

    let exchange: MoaPairExchange;

    if (dharitrixTypes.includes(type)) {
      exchange = MoaPairExchange.dharitrix;
    } else {
      exchange = MoaPairExchange.unknown;
    }

    const baseInfo = {
      address: pair.address,
      id: pair.liquidityPoolToken.identifier,
      symbol: pair.liquidityPoolToken.identifier.split('-')[0],
      name: pair.liquidityPoolToken.name,
      price: Number(pair.liquidityPoolTokenPriceUSD),
      totalValue: Number(pair.lockedValueUSD),
      volume24h: Number(pair.volumeUSD24h),
      tradesCount: Number(pair.tradesCount),
      tradesCount24h: Number(pair.tradesCount24h),
      deployedAt: Number(pair.deployedAt),
      state,
      type,
      exchange,
      ...(includeFarms && {
        hasFarms: pair.hasFarms ?? false,
        hasDualFarms: pair.hasDualFarms ?? false,
      }),
    };

    if ((firstTokenSymbol === 'WREWA' && secondTokenSymbol === 'USDC') || secondTokenSymbol === 'WREWA') {
      return {
        ...baseInfo,
        basePrevious24hPrice: Number(pair.firstToken.previous24hPrice),
        quotePrevious24hPrice: Number(pair.secondToken.previous24hPrice),
        baseId: pair.firstToken.identifier,
        basePrice: Number(pair.firstTokenPriceUSD),
        baseSymbol: firstTokenSymbol,
        baseName: pair.firstToken.name,
        quoteId: pair.secondToken.identifier,
        quotePrice: Number(pair.secondTokenPriceUSD),
        quoteSymbol: secondTokenSymbol,
        quoteName: pair.secondToken.name,
      };
    }

    return {
      ...baseInfo,
      basePrevious24hPrice: Number(pair.secondToken.previous24hPrice),
      quotePrevious24hPrice: Number(pair.firstToken.previous24hPrice),
      baseId: pair.secondToken.identifier,
      basePrice: Number(pair.secondTokenPriceUSD),
      baseSymbol: secondTokenSymbol,
      baseName: pair.secondToken.name,
      quoteId: pair.firstToken.identifier,
      quotePrice: Number(pair.firstTokenPriceUSD),
      quoteSymbol: firstTokenSymbol,
      quoteName: pair.firstToken.name,
    };
  }

  private getPairState(state: string): MoaPairState {
    switch (state) {
      case 'Active':
        return MoaPairState.active;
      case 'Inactive':
        return MoaPairState.inactive;
      case 'ActiveNoSwaps':
        return MoaPairState.paused;
      case 'PartialActive':
        return MoaPairState.partial;
      default:
        throw new Error(`Unsupported pair state '${state}'`);
    }
  }

  private getPairType(type: string): MoaPairType | undefined {
    switch (type) {
      case 'Core':
        return MoaPairType.core;
      case 'Community':
        return MoaPairType.community;
      case 'Ecosystem':
        return MoaPairType.ecosystem;
      case 'Experimental':
        return MoaPairType.experimental;
      case 'Unlisted':
        return MoaPairType.unlisted;
      default:
        this.logger.error(`Unsupported pair type '${type}'`);
        return undefined;
    }
  }

  private applyFilters(drtPairs: MoaPair[], filter?: MoaPairsFilter): MoaPair[] {
    if (!filter) {
      return drtPairs;
    }

    let filteredPairs = drtPairs;

    if (filter.exchange) {
      filteredPairs = filteredPairs.filter(pair => pair.exchange === filter.exchange);
    }

    return filteredPairs;
  }
}
