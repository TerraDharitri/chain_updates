import { registerEnumType } from '@nestjs/graphql';

export enum DcdtType {
  FungibleDCDT = 'FungibleDCDT',
  NonFungibleDCDT = 'NonFungibleDCDT',
  SemiFungibleDCDT = 'SemiFungibleDCDT',
  MetaDCDT = 'MetaDCDT',
}

registerEnumType(DcdtType, {
  name: 'DcdtType',
  description: 'Dcdt type enum.',
  valuesMap: {
    FungibleDCDT: {
      description: 'Fungible DCDT token type.',
    },
    NonFungibleDCDT: {
      description: 'Non-fungible DCDT token type.',
    },
    SemiFungibleDCDT: {
      description: 'Semi-fungible DCDT token type.',
    },
    MetaDCDT: {
      description: 'Meta DCDT token type.',
    },
  },
});
