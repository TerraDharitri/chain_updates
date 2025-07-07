import { registerEnumType } from "@nestjs/graphql";

export enum NftType {
  NonFungibleDCDT = 'NonFungibleDCDT',
  SemiFungibleDCDT = 'SemiFungibleDCDT',
  MetaDCDT = 'MetaDCDT',
}

registerEnumType(NftType, {
  name: 'NftType',
  description: 'NFT type.',
  valuesMap: {
    NonFungibleDCDT: {
      description: 'Non-fungible NFT type.',
    },
    SemiFungibleDCDT: {
      description: 'Semi-fungible NFT type.',
    },
    MetaDCDT: {
      description: 'Meta DCDT NFT type.',
    },
  },
});
