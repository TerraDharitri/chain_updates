export enum TokenTypesEnum {
  core = 'Core',
  ecosystem = 'Ecosystem',
  community = 'Community',
  experimental = 'Experimental',
  unlisted = 'Unlisted'
}

export interface AssetsType {
  website?: string;
  description?: string;
  status?: string;
  pngUrl?: string;
  svgUrl?: string;
  social?: any;
}

export interface TokenBaseType {
  name: string;
  decimals: number;
  ticker: string;
  assets?: AssetsType;
}

export interface DcdtType extends TokenBaseType {
  balance: string | null;
  identifier: string;
  owner: string;
  price?: string;
  type?: TokenTypesEnum | 'FungibleDCDT-LP';
  previous24hPrice?: string;
  previous7dPrice?: string;
}

export interface NftCollectionType extends TokenBaseType {
  collection: string;
}

export interface UserDcdtType extends DcdtType {
  valueUSD: string;
  usdPrice?: string;
}

export interface WrappingInfoType {
  wrappedToken: DcdtType;
}
