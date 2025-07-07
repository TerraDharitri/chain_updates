import { useEffect, useMemo, useState } from 'react';
import BigNumber from 'bignumber.js';
import {
  SelectOptionType,
  SwapRouteType,
  SwapFeeDetailsType,
  SwapActionTypesEnum
} from 'types';
import {
  getPairFeeDetails,
  calculateMinimumReceived,
  removeCommas,
  calculateSwapTransactionsFee
} from 'utils';
import { useRateCalculator } from './useRateCalculator';

interface UseSwapInfoProps {
  tolerance: string;
  firstToken?: SelectOptionType;
  firstAmount?: string;
  secondToken?: SelectOptionType;
  secondAmount?: string;
  activeRoute?: SwapRouteType;
  swapActionType?: SwapActionTypesEnum;
}

export const useSwapInfo = ({
  tolerance,
  activeRoute,
  firstToken,
  firstAmount,
  secondToken,
  secondAmount,
  swapActionType
}: UseSwapInfoProps) => {
  const [minimumAmountReceived, setMinimumAmountReceived] = useState<string>();
  const [feeAmounts, setFeeAmounts] = useState<string[]>();
  const [feePercentages, setFeePercentages] = useState<string[]>();
  const [pricesImpact, setPricesImpact] = useState<string[]>();
  const [feeDetails, setFeeDetails] = useState<SwapFeeDetailsType>();

  const { tokenInId, tokenOutId, exchangeRate, switchTokensDirection } =
    useRateCalculator({
      activeRoute
    });

  const totalTransactionsFee = useMemo(
    () => calculateSwapTransactionsFee(activeRoute?.transactions),
    [activeRoute?.transactions]
  );

  const cleanExchangeRate = removeCommas(exchangeRate || '');
  const cleanFirstTokenPrice = removeCommas(firstToken?.token?.price ?? '');
  const cleanSecondTokenPrice = removeCommas(secondToken?.token?.price ?? '');

  const tokenPrice =
    tokenInId === firstToken?.value
      ? cleanSecondTokenPrice
      : cleanFirstTokenPrice;

  const exchangeRateUsdValue = new BigNumber(cleanExchangeRate || '0')
    .multipliedBy(tokenPrice || '0')
    .toString(10);

  const handleUpdateStats = () => {
    if (!activeRoute) {
      return;
    }

    const { swapType, pairs } = activeRoute;

    const minimumReceived = calculateMinimumReceived({
      tolerance,
      secondAmount,
      secondTokenDecimals: secondToken?.token.decimals,
      isFixedOutput: swapType === 1,
      swapActionType
    });

    const newFeePercentages: string[] = [];

    pairs.forEach((pair) => {
      const { totalFee, burn, lpHolders } = getPairFeeDetails(pair);
      const isDirectSwap = pairs.length === 1;
      const shouldPushFeePercentage = !isDirectSwap && totalFee;

      if (shouldPushFeePercentage) {
        newFeePercentages.push(`${totalFee}%`);
      }

      const newFeeTooltip = {
        totalFee,
        burn,
        lpHolders
      };

      setFeeDetails(newFeeTooltip);
    });

    const newPricesImpacts = firstAmount ? activeRoute.pricesImpact : undefined;
    const newFeeAmounts = firstAmount ? activeRoute.fees : undefined;

    setMinimumAmountReceived(minimumReceived);
    setFeeAmounts(newFeeAmounts);
    setPricesImpact(newPricesImpacts);
    setFeePercentages(newFeePercentages);
  };

  useEffect(handleUpdateStats, [
    tolerance,
    activeRoute,
    firstToken,
    firstAmount,
    secondToken,
    secondAmount,
    swapActionType
  ]);

  return {
    exchangeRate,
    feeDetails,
    feeAmounts,
    pricesImpact,
    feePercentages,
    tokenInId,
    tokenOutId,
    switchTokensDirection,
    totalTransactionsFee,
    minimumAmountReceived,
    exchangeRateUsdValue
  };
};
