import { UpdateMarketplaceArgs } from '../UpdateMarketplaceArgs';
import { WhitelistMarketplaceArgs } from '../WhitelistMarketplaceArgs';
import { drtConfig } from 'src/config';

export class UpdateMarketplaceRequest {
  marketplaceName: string;
  marketplaceKey: string;
  marketplaceUrl: string;
  marketplaceScAddress: string = drtConfig.nftMarketplaceAddress;

  constructor(init?: Partial<UpdateMarketplaceRequest>) {
    Object.assign(this, init);
  }

  static fromArgs(args: UpdateMarketplaceArgs) {
    return new UpdateMarketplaceRequest({
      marketplaceName: args.marketplaceName,
      marketplaceKey: args.marketplaceKey,
      marketplaceUrl: args.marketplaceUrl,
      marketplaceScAddress: args.marketplaceScAddress,
    });
  }
}
