import { useMemo } from 'react';
import { formatAmount } from '@terradharitri/sdk-dapp/utils';
import BigNumber from 'bignumber.js';
import { SwapRouteType, SelectOptionType } from 'types';
import { roundAmount } from 'utils';

// input amount can be truncated
// this hooks returns the usd value for the complete input amount
export const useInputAmountUsdValue = ({
  swapRoute,
  firstToken,
  firstAmount,
  secondToken,
  secondAmount
}: {
  firstAmount: string;
  secondAmount: string;
  swapRoute?: SwapRouteType;
  firstToken?: SelectOptionType;
  secondToken?: SelectOptionType;
}) => {
  return useMemo(() => {
    if (!swapRoute) {
      return {
        firstAmountUsdValue: undefined,
        secondAmountUsdValue: undefined
      };
    }

    const {
      swapType,
      amountIn,
      amountOut,
      tokenInID,
      tokenOutID,
      tokenInPriceUSD,
      tokenOutPriceUSD
    } = swapRoute;

    const isFixedInput = swapType === 0;
    const isFixedOutput = swapType === 1;

    const shouldComputeFirstUsdValue = Boolean(
      isFixedOutput && secondAmount && firstToken?.value === tokenInID
    );
    const shouldComputeSecondUsdValue = Boolean(
      isFixedInput && firstAmount && secondToken?.value === tokenOutID
    );

    if (!shouldComputeFirstUsdValue && !shouldComputeSecondUsdValue) {
      return {
        firstAmountUsdValue: undefined,
        secondAmountUsdValue: undefined
      };
    }

    const input = shouldComputeFirstUsdValue ? amountIn : amountOut;
    const decimals = shouldComputeFirstUsdValue
      ? firstToken?.token.decimals
      : secondToken?.token.decimals;

    const usd = shouldComputeFirstUsdValue ? tokenInPriceUSD : tokenOutPriceUSD;
    const amount = formatAmount({
      input,
      decimals,
      showLastNonZeroDecimal: true
    });
    const amountWithoutCommas = amount.replace(/,/g, '');

    const computedUsdValue = new BigNumber(amountWithoutCommas)
      .times(new BigNumber(usd))
      .toString(10);
    const usdValue = roundAmount(computedUsdValue, 2);

    return {
      firstAmountUsdValue: shouldComputeFirstUsdValue ? usdValue : undefined,
      secondAmountUsdValue: shouldComputeSecondUsdValue ? usdValue : undefined
    };
  }, [swapRoute, secondAmount, firstAmount, firstToken, secondToken]);
};
