import { registerEnumType } from "@nestjs/graphql";

export interface TokenAccount {
  identifier: string;
  address: string;
  balance: string;
  balanceNum: number;
  token: string;
  timestamp: number;
  type: TokenType;
  data: any;
  nft_scamInfoType: string;
  nft_scamInfoDescription: string;
}

export enum TokenType {
  FungibleDCDT = 'FungibleDCDT',
  NonFungibleDCDT = 'NonFungibleDCDT',
  SemiFungibleDCDT = 'SemiFungibleDCDT',
  MetaDCDT = 'MetaDCDT',
}

registerEnumType(TokenType, {
  name: 'TokenType',
  description: 'Token Type object.',
  valuesMap: {
    FungibleDCDT: {
      description: 'FungibleDCDT.',
    },
    NonFungibleDCDT: {
      description: 'NonFungibleDCDT.',
    },
    SemiFungibleDCDT: {
      description: 'SemiFungibleDCDT.',
    },
    MetaDCDT: {
      description: 'MetaDCDT.',
    },
  },
});
