import { Injectable, Logger } from '@nestjs/common';
import { DrtApiService, DrtStatsService } from 'src/common';
import { Asset } from '../assets/models';

@Injectable()
export class TrendingService {
  constructor(private apiService: DrtApiService, private statsService: DrtStatsService, private readonly logger: Logger) {}

  async getTrendingAssets(): Promise<[Asset[], number]> {
    try {
      const trendingNfts = await this.statsService.getTrending('identifier');
      const nfts = await this.apiService.getNftsByIdentifiers(trendingNfts);
      return [nfts?.map((nft) => Asset.fromNft(nft)), trendingNfts.length];
    } catch (err) {
      this.logger.error('An error occurred while loading trending assets.', {
        path: 'FeaturedNftsService.getTrendingAssets',
        exception: err,
      });
    }
  }
}
