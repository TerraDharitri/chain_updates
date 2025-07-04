import { registerEnumType } from '@nestjs/graphql';

export enum DcdtSubType {
  NonFungibleDCDTv2 = 'NonFungibleDCDTv2',
  DynamicNonFungibleDCDT = 'DynamicNonFungibleDCDT',
  DynamicSemiFungibleDCDT = 'DynamicSemiFungibleDCDT',
  DynamicMetaDCDT = 'DynamicMetaDCDT',
}

registerEnumType(DcdtSubType, {
  name: 'DcdtSubType',
  description: 'Dcdt sub type enum.',
  valuesMap: {
    NonFungibleDCDTv2: {
      description: 'Non-fungible DCDT v2 sub type.',
    },
    DynamicNonFungibleDCDT: {
      description: 'Dynamic non-fungible sub type.',
    },
    DynamicSemiFungibleDCDT: {
      description: 'Dynamic semi-fungible sub type.',
    },
    DynamicMetaDCDT: {
      description: 'Dynamic meta DCDT sub type.',
    },
  },
});
