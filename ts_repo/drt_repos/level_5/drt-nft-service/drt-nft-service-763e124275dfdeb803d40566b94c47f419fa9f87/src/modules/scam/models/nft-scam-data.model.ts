import { MxApiAbout } from 'src/common/services/drt-communication/models/drt-api-about.model';
import { Asset } from 'src/modules/assets/models';
import { NftScamInfoModel } from './nft-scam-info.model';

export class NftScamRelatedData {
  drtApiAbout?: MxApiAbout;
  nftFromApi?: Asset;
  nftFromElastic?: any;
  nftFromDb?: NftScamInfoModel;

  constructor(init?: Partial<NftScamRelatedData>) {
    Object.assign(this, init);
  }
}
