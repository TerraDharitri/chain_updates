import { formatAmount } from '@terradharitri/sdk-dapp/utils';
import { parseAmount } from '@terradharitri/sdk-dapp/utils/operations/parseAmount';
import BigNumber from 'bignumber.js';
import { SwapActionTypesEnum } from 'types';

export const calculateMinimumReceived = ({
  tolerance,
  secondAmount,
  isFixedOutput,
  swapActionType,
  secondTokenDecimals
}: {
  tolerance: string;
  secondAmount?: string;
  isFixedOutput: boolean;
  swapActionType?: SwapActionTypesEnum;
  secondTokenDecimals?: number;
}) => {
  if (!secondAmount) {
    return;
  }

  if (isFixedOutput || swapActionType !== SwapActionTypesEnum.swap) {
    return secondAmount;
  }

  const bnSecondAmount = new BigNumber(
    parseAmount(secondAmount, secondTokenDecimals)
  );

  const minAmount = new BigNumber(1)
    .dividedBy(new BigNumber(tolerance).dividedBy(100).plus(1))
    .times(bnSecondAmount)
    .toFixed(0);

  return formatAmount({ input: minAmount, decimals: secondTokenDecimals });
};
