import { TokenTransfer } from '@terradharitri/sdk-core';
import { DECIMALS } from '../constants';

export function parseAmount(amount: string, numDecimals: number = DECIMALS) {
  return TokenTransfer.fungibleFromAmount('', amount, numDecimals).toString();
}
