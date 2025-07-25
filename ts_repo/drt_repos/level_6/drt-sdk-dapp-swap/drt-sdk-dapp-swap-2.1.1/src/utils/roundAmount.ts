import { stringIsFloat } from '@terradharitri/sdk-dapp/utils/validation/stringIsFloat';
import BigNumber from 'bignumber.js';
import { DIGITS } from 'constants/index';

export const roundAmount = (amount: string, digits = DIGITS) => {
  if (!stringIsFloat(amount)) {
    return '0';
  }

  const bNamount = new BigNumber(amount);

  if (bNamount.isZero()) {
    return '0';
  }

  let formattedAmount = new BigNumber(amount).toFormat(digits);

  formattedAmount =
    parseFloat(formattedAmount) > 0
      ? formattedAmount
      : new BigNumber(amount).toFormat(digits + 4);

  return parseFloat(formattedAmount) > 0
    ? formattedAmount
    : new BigNumber(amount).toFormat(digits + 9);
};
