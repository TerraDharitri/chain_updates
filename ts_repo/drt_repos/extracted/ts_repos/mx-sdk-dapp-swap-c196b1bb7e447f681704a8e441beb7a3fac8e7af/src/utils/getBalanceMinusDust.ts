import BigNumber from 'bignumber.js';
import { MIN_EGLD_DUST } from 'constants/general';

export const getBalanceMinusDust = (balance?: string) => {
  const bNbalance = new BigNumber(balance ?? '0');
  const bNminDust = new BigNumber(MIN_EGLD_DUST);
  const balanceSubtractDust = bNbalance.minus(bNminDust);

  return balanceSubtractDust.isGreaterThan(0)
    ? balanceSubtractDust.toString(10)
    : '0';
};
