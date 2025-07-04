import { NftTypeEnum } from 'types';

export const getNftText = (type: NftTypeEnum) => {
  switch (type) {
    case NftTypeEnum.SemiFungibleDCDT:
      return 'SFT';
    case NftTypeEnum.NonFungibleDCDT:
      return 'NFT';
    case NftTypeEnum.MetaDCDT:
      return 'Meta-DCDT';
    default:
      return '';
  }
};
