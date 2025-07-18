import { registerEnumType } from '@nestjs/graphql';

export enum TokenType {
  FungibleDCDT = 'FungibleDCDT',
  MetaDCDT = 'MetaDCDT',
}

registerEnumType(TokenType, {
  name: 'TokenType',
  description: 'Token type enum.',
  valuesMap: {
    FungibleDCDT: {
      description: 'Fungible DCDT token type.',
    },
    MetaDCDT: {
      description: 'Meta DCDT token type.',
    },
  },
});
