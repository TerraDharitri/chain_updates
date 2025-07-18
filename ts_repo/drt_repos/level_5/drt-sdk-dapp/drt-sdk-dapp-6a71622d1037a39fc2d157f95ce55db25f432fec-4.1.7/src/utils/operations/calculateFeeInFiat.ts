import { DIGITS, DECIMALS } from 'constants/index';

import { formatAmount } from './formatAmount';
import { getUsdValue } from './getUsdValue';

export interface CalculateFeeInFiatType {
  feeLimit: string;
  rewaPriceInUsd: number;
  hideEqualSign?: boolean;
}

export const calculateFeeInFiat = ({
  feeLimit,
  rewaPriceInUsd,
  hideEqualSign
}: CalculateFeeInFiatType) => {
  const amount = formatAmount({
    input: feeLimit,
    decimals: DECIMALS,
    digits: DIGITS,
    showLastNonZeroDecimal: true
  });

  const feeAsUsdValue = getUsdValue({
    amount,
    usd: rewaPriceInUsd,
    decimals: DIGITS
  });

  if (hideEqualSign) {
    return feeAsUsdValue;
  }

  return `≈ ${feeAsUsdValue}`;
};
