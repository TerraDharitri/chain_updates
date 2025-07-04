import { DcdtEnumType, NftEnumType } from 'types/tokens.types';

export const getTokenType = (type?: NftEnumType) => {
  switch (type) {
    case NftEnumType.NonFungibleDCDT:
    case NftEnumType.SemiFungibleDCDT:
    case NftEnumType.MetaDCDT:
      return type;
    default:
      return DcdtEnumType.FungibleDCDT;
  }
};
