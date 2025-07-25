import { parseAmount } from '@terradharitri/sdk-dapp/utils';
import { FIXED_INPUT } from 'constants/general';
import { SwapRouteType, UserDcdtType } from 'types';
import { meaningfulFormatAmount } from './meaningfulFormatAmount';

export const getCorrectAmountsOnTokenChange = ({
  newToken,
  firstToken,
  secondToken,
  activeRoute,
  needsParsing
}: {
  needsParsing: boolean;
  newToken?: UserDcdtType;
  firstToken?: UserDcdtType;
  secondToken?: UserDcdtType;
  activeRoute?: SwapRouteType;
}) => {
  if (!activeRoute || !newToken)
    return {
      amountIn: undefined,
      amountOut: undefined
    };

  const isFixedInput = activeRoute?.swapType === FIXED_INPUT;
  const inputParsedAmount = isFixedInput
    ? activeRoute?.amountIn
    : activeRoute?.amountOut;

  const inputToken = isFixedInput
    ? activeRoute?.tokenInID
    : activeRoute?.tokenOutID;

  const inputTokenDecimals =
    inputToken === firstToken?.identifier
      ? firstToken?.decimals
      : secondToken?.decimals;

  // format input amount with the input decimals
  // and preserve precision in case it was a max button amount
  const formattedInputParsedAmount = meaningfulFormatAmount({
    amount: inputParsedAmount,
    showLastNonZeroDecimal: true,
    decimals: inputTokenDecimals
  });

  // parse the precise old amount with the new decimals
  const parsedNewAmount = needsParsing
    ? parseAmount(formattedInputParsedAmount, newToken?.decimals)
    : inputParsedAmount;

  const amountIn = isFixedInput ? parsedNewAmount : undefined;
  const amountOut = isFixedInput ? undefined : parsedNewAmount;

  return {
    amountIn,
    amountOut
  };
};
