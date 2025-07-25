import { DECIMALS } from '@terradharitri/sdk-dapp/constants';

export const canParseAmount = (amount: string, paramDecimals?: number) => {
  const isFloat = amount.indexOf('.') >= 0;

  if (!isFloat) {
    return true;
  }

  const maxDecimals = paramDecimals ?? DECIMALS;
  const decimalPart = amount.split('.').pop() ?? '';
  const lengthDoesNotExceedMaxDecimals = decimalPart.length <= maxDecimals;

  return lengthDoesNotExceedMaxDecimals;
};
