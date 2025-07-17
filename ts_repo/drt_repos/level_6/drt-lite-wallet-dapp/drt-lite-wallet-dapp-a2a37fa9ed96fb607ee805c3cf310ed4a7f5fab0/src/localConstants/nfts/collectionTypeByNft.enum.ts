import { NftEnumType } from 'types/sdkDapp.types';

export const CollectionTypeByNftEnum = {
  [NftEnumType.SemiFungibleDCDT]: 'sft',
  [NftEnumType.NonFungibleDCDT]: 'nft',
  [NftEnumType.MetaDCDT]: 'meta'
};
