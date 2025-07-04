import { Inject, Injectable, Logger, forwardRef } from '@nestjs/common';
import { drtConfig } from 'src/config';
import { AuctionEntity } from 'src/db/auctions';
import { KroganSwapAuctionEventEnum } from 'src/modules/assets/models';
import { AuctionsGetterService, AuctionsSetterService } from 'src/modules/auctions';
import { MarketplacesService } from 'src/modules/marketplaces/marketplaces.service';
import { MarketplaceTypeEnum } from 'src/modules/marketplaces/models/MarketplaceType.enum';
import { UsdPriceService } from 'src/modules/usdPrice/usd-price.service';
import { BigNumberUtils } from 'src/utils/bigNumber-utils';
import { NumbatSwapUpdateEvent } from '../../entities/auction/numbatnftswap/numbatswap-updateAuction.event';

@Injectable()
export class SwapUpdateEventHandler {
  private readonly logger = new Logger(SwapUpdateEventHandler.name);
  constructor(
    private auctionsGetterService: AuctionsGetterService,
    private auctionsService: AuctionsSetterService,
    @Inject(forwardRef(() => MarketplacesService))
    private readonly marketplaceService: MarketplacesService,
    private usdPriceService: UsdPriceService,
  ) {}

  async handle(event: any, hash: string, marketplaceType: MarketplaceTypeEnum) {
    const updateEvent = new NumbatSwapUpdateEvent(event);
    const topics = updateEvent.getTopics();
    const marketplace = await this.marketplaceService.getMarketplaceByType(updateEvent.getAddress(), marketplaceType);
    this.logger.log(`${updateEvent.getIdentifier()}  auction event detected for hash '${hash}' and marketplace '${marketplace?.name}'`);
    let auction = await this.auctionsGetterService.getAuctionByIdAndMarketplace(parseInt(topics.auctionId, 16), marketplace.key);

    if (auction) {
      await this.updateAuctionPrice(auction, topics, hash);
      this.auctionsService.updateAuction(auction, KroganSwapAuctionEventEnum.NftSwapUpdate);
    }
  }

  private async updateAuctionPrice(
    auction: AuctionEntity,
    topics: {
      seller: string;
      collection: string;
      nonce: string;
      auctionId: string;
      nrAuctionTokens: number;
      price: string;
      deadline: number;
    },
    hash: string,
  ) {
    const paymentToken = await this.usdPriceService.getToken(auction.paymentToken);
    const decimals = paymentToken?.decimals ?? drtConfig.decimals;
    auction.minBid = topics.price;
    auction.minBidDenominated = BigNumberUtils.denominateAmount(topics.price, decimals);
    auction.maxBid = auction.minBid;
    auction.maxBidDenominated = auction.minBidDenominated;
    auction.endDate = topics.deadline;
    auction.nrAuctionedTokens = topics.nrAuctionTokens;
    auction.blockHash = hash;
  }
}
