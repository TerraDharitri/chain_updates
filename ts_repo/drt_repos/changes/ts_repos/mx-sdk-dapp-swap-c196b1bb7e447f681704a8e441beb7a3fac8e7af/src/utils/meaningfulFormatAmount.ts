import { DIGITS } from 'constants/index';
import { formatAmount, stringIsInteger } from 'lib';

export const meaningfulFormatAmount = ({
  amount = '0',
  decimals = 0,
  digits = DIGITS,
  addCommas = false,
  showLastNonZeroDecimal = false
}: {
  digits?: number;
  decimals?: number;
  addCommas?: boolean;
  amount?: string | null;
  showLastNonZeroDecimal?: boolean;
}) => {
  if (amount == null || !stringIsInteger(amount)) {
    return '0';
  }

  const tryFormatAmount = (tryDigits: number) =>
    formatAmount({
      decimals,
      addCommas,
      input: amount,
      digits: tryDigits,
      showLastNonZeroDecimal
    });

  let formattedAmount = tryFormatAmount(digits);

  formattedAmount =
    parseFloat(formattedAmount) > 0
      ? formattedAmount
      : tryFormatAmount(DIGITS + 4);

  formattedAmount =
    parseFloat(formattedAmount) > 0
      ? formattedAmount
      : tryFormatAmount(DIGITS + 9);

  return parseFloat(formattedAmount) > 0
    ? formattedAmount
    : tryFormatAmount(DIGITS + 14);
};
