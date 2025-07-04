import BigNumber from 'bignumber.js';
import { PairType } from 'types';

const empty = {
  burn: undefined,
  feesAPR: undefined,
  totalFee: undefined,
  lpHolders: undefined
};

export const getPairFeeDetails = (pair?: PairType) => {
  if (!pair) {
    return empty;
  }

  const { totalFeePercent, specialFeePercent, feesAPR: pairFeesAPR } = pair;

  const totalFee = totalFeePercent * 100;
  const burn = specialFeePercent * 100;
  const lpHolderPercentage = totalFeePercent - specialFeePercent;

  const lpHolders = lpHolderPercentage * 100;

  const feesAPR = new BigNumber(pairFeesAPR).multipliedBy(100).toNumber();

  return { totalFee, burn, lpHolders, feesAPR };
};
