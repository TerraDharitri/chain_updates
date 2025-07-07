import { NftType, NftTypeEnum } from 'types';

export const hasNftOverview = ({ type, metadata, rarities }: NftType) => {
  return Boolean(
    type &&
      type !== NftTypeEnum.MetaDCDT &&
      (metadata?.attributes || (rarities && Object.keys(rarities).length > 0))
  );
};
