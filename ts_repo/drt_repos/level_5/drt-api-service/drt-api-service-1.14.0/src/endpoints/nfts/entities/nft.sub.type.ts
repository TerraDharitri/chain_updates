import { registerEnumType } from "@nestjs/graphql";

export enum NftSubType {
  NonFungibleDCDT = 'NonFungibleDCDT',
  SemiFungibleDCDT = 'SemiFungibleDCDT',
  MetaDCDT = 'MetaDCDT',
  NonFungibleDCDTv2 = 'NonFungibleDCDTv2',
  DynamicNonFungibleDCDT = 'DynamicNonFungibleDCDT',
  DynamicSemiFungibleDCDT = 'DynamicSemiFungibleDCDT',
  DynamicMetaDCDT = 'DynamicMetaDCDT',
  None = '',
}

registerEnumType(NftSubType, {
  name: 'NftSubType',
  description: 'NFT subtype.',
  valuesMap: {
    NonFungibleDCDT: {
      description: 'Non-fungible DCDT NFT type.',
    },
    SemiFungibleDCDT: {
      description: 'Semi-fungible DCDT NFT type.',
    },
    MetaDCDT: {
      description: 'Meta DCDT NFT type.',
    },
    NonFungibleDCDTv2: {
      description: 'Non-fungible DCDT v2 NFT type.',
    },
    DynamicNonFungibleDCDT: {
      description: 'Dynamic non-fungible NFT type.',
    },
    DynamicSemiFungibleDCDT: {
      description: 'Dynamic semi-fungible NFT type.',
    },
    DynamicMetaDCDT: {
      description: 'Dynamic meta DCDT NFT type.',
    },
    None: {
      description: '',
    },
  },
});
