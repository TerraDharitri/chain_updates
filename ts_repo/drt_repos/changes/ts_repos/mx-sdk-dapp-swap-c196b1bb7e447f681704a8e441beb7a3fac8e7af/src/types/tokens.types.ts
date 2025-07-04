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

export interface EsdtType extends TokenBaseType {
  balance: string | null;
  identifier: string;
  owner: string;
  price?: string;
  type?: 'FungibleESDT-LP';
  previous24hPrice?: string;
  previous7dPrice?: string;
}

export interface NftCollectionType extends TokenBaseType {
  collection: string;
}

export interface UserEsdtType extends EsdtType {
  valueUSD: string;
  usdPrice?: string;
}

export interface WrappingInfoType {
  wrappedToken: EsdtType;
}

export interface TokensPaginationType {
  before?: string;
  after?: string;
  first?: number;
  last?: number;
}
