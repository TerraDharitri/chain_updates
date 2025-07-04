import { useEffect, useState } from 'react';
import { SwapRouteType } from 'types';

export const useRateCalculator = ({
  activeRoute
}: {
  activeRoute?: SwapRouteType;
}) => {
  const [direction, setDirection] = useState<'normal' | 'reverse'>('normal');

  const [exchangeRate, setExchangeRate] = useState<string>();
  const [tokenInId, setTokenInId] = useState<string>();
  const [tokenInPriceUsd, setTokenInPriceUsd] = useState<string>();
  const [tokenOutId, setTokenOutId] = useState<string>();

  const [reverseExchangeRate, setReverseExchangeRate] = useState<string>();
  const [reverseTokenInId, setReverseTokenInId] = useState<string>();
  const [reverseTokenInPriceUsd, setReverseTokenInPriceUsd] =
    useState<string>();
  const [reverseTokenOutId, setReverseTokenOutId] = useState<string>();

  const calcRate = () => {
    if (!activeRoute) {
      return;
    }

    const {
      tokenInID,
      tokenOutID,
      tokenInExchangeRateDenom,
      tokenInPriceUSD,
      tokenOutExchangeRateDenom,
      tokenOutPriceUSD
    } = activeRoute;

    setTokenInId(tokenOutID);
    setTokenOutId(tokenInID);
    setExchangeRate(tokenOutExchangeRateDenom);
    setTokenInPriceUsd(tokenOutPriceUSD);

    setReverseTokenInId(tokenInID);
    setReverseTokenOutId(tokenOutID);
    setReverseExchangeRate(tokenInExchangeRateDenom);
    setReverseTokenInPriceUsd(tokenInPriceUSD);
  };

  const switchTokensDirection = () => {
    setDirection((existing) => (existing === 'normal' ? 'reverse' : 'normal'));
  };

  useEffect(calcRate, [activeRoute, direction]);

  return {
    tokenInId: direction === 'normal' ? tokenInId : reverseTokenInId,
    tokenInIdPriceUsd:
      direction === 'normal' ? tokenInPriceUsd : reverseTokenInPriceUsd,
    tokenOutId: direction === 'normal' ? tokenOutId : reverseTokenOutId,
    exchangeRate: direction === 'normal' ? exchangeRate : reverseExchangeRate,
    switchTokensDirection
  };
};
