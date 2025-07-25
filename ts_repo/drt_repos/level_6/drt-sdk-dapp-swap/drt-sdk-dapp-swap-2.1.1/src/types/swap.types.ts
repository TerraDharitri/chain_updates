import { RawTransactionType } from '@terradharitri/sdk-dapp/types/transactions.types';
import { PairType } from './pairs.types';

export interface SwapRouteType {
  amountIn: string;
  tokenInID: string;
  tokenInPriceUSD: string;
  tokenInExchangeRateDenom: string;

  amountOut: string;
  tokenOutID: string;
  tokenOutPriceUSD: string;
  tokenOutExchangeRateDenom: string;

  fees: string[];
  swapType: number; // 0 fixedInput, 1 fixedOutput
  tokenRoute: string[];
  pricesImpact: string[];

  maxPriceDeviationPercent: number;
  tokensPriceDeviationPercent: number | null;

  intermediaryAmounts: string[];
  pairs: PairType[];
  transactions?: RawTransactionType[];
}

export enum SwapActionTypesEnum {
  wrap = 'wrap',
  unwrap = 'unwrap',
  swap = 'swap'
}

export interface SwapFeeDetailsType {
  totalFee?: number;
  burn?: number;
  lpHolders?: number;
}
