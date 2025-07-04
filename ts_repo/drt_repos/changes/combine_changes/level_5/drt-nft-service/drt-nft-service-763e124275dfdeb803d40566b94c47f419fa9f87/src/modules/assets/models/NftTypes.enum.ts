import { registerEnumType } from '@nestjs/graphql';

export enum NftTypeEnum {
  NonFungibleDCDT = 'NonFungibleDCDT',
  SemiFungibleDCDT = 'SemiFungibleDCDT',
}

registerEnumType(NftTypeEnum, {
  name: 'NftTypeEnum',
});

export enum ScamInfoTypeEnum {
  none = 'none',
  potentialScam = 'potentialScam',
  scam = 'scam',
}

registerEnumType(ScamInfoTypeEnum, {
  name: 'ScamInfoTypeEnum',
});
