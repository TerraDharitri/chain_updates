import { DcdtType, NftCollectionType } from './tokens.types';

export enum PairStatesEnum {
  active = 'Active',
  inactive = 'Inactive',
  partialActive = 'PartialActive'
}

export interface SimpleLockType {
  address: string;

  lockedToken: NftCollectionType; // LKDCDT
  lpProxyToken: NftCollectionType; // LKLP
  farmProxyToken: NftCollectionType; // LKFARM

  intermediatedPairs: string[];
  intermediatedFarms: string[];
}

export interface LockedTokensInfoType {
  unlockEpoch: number;
  lockingSC: SimpleLockType;
}

export interface PairType {
  address: string;
  firstToken: DcdtType;
  secondToken: DcdtType;
  feesAPR: string;
  totalFeePercent: number;
  specialFeePercent: number;
  lockedTokensInfo?: LockedTokensInfoType;
}
