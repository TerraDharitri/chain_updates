import { useRef, useState } from 'react';
import { EGLD_IDENTIFIER, FIXED_INPUT, FIXED_OUTPUT } from 'constants/general';
import { stringIsFloat, parseAmount } from 'lib';
import { SwapRouteType, SelectOptionType } from 'types';
import {
  getTokenDecimals,
  getBalanceMinusDust,
  meaningfulFormatAmount,
  getCorrectAmountsOnTokenChange
} from 'utils';

import { useInputAmountUsdValue } from './useInputAmountUsdValue';
import { usePrevious } from './usePrevious';
import { GetSwapRouteType } from './useSwapRoute';

type SwapFormStateType = {
  firstToken?: SelectOptionType;
  secondToken?: SelectOptionType;
  firstAmount: string;
  secondAmount: string;
  resetSwapping?: boolean;
};

enum InputTouchedEnum {
  first = 'first',
  second = 'second'
}

export const useSwapFormHandlers = ({
  getSwapRoute,
  tolerancePercentage
}: {
  getSwapRoute: ({}: GetSwapRouteType) => void;
  tolerancePercentage?: number;
}) => {
  const [formState, setFormState] = useState<SwapFormStateType>({
    firstAmount: '',
    secondAmount: ''
  });
  const [activeRoute, setActiveRoute] = useState<SwapRouteType>();
  const isSwitching = useRef<boolean>();
  const lastInputTouched = useRef<InputTouchedEnum>();
  const prevSecondToken = usePrevious(formState.secondToken);

  const inputAmountsUsdValue = useInputAmountUsdValue({
    swapRoute: activeRoute,
    firstToken: formState.firstToken,
    firstAmount: formState.firstAmount,
    secondToken: formState.secondToken,
    secondAmount: formState.secondAmount
  });

  const handleOnChangeSwapRoute = (swapRoute?: SwapRouteType) => {
    if (!swapRoute) {
      setActiveRoute(undefined);
      return;
    }

    const { swapType, tokenInID, tokenOutID, amountIn, amountOut, pairs } =
      swapRoute;

    setActiveRoute(swapRoute);

    const isFixedInput = swapType === FIXED_INPUT;
    const amount = isFixedInput ? amountOut : amountIn;
    const identifier = isFixedInput ? tokenOutID : tokenInID;
    const decimals = getTokenDecimals({ identifier, pairs });

    const amountToDisplayDenom = meaningfulFormatAmount({ amount, decimals });

    const shouldUpdateSecondAmount =
      isFixedInput && Boolean(formState.firstAmount);
    const shouldUpdateFirstAmount =
      !isFixedInput && Boolean(formState.secondAmount);

    if (isSwitching.current) {
      const firstAmountDecimals = getTokenDecimals({
        identifier: swapRoute.tokenInID,
        pairs
      });
      const secondAmountDecimals = getTokenDecimals({
        identifier: swapRoute.tokenOutID,
        pairs
      });

      const firstAmountToDisplayDenom = meaningfulFormatAmount({
        amount: swapRoute.amountIn,
        decimals: firstAmountDecimals
      });
      const secondAmountToDisplayDenom = meaningfulFormatAmount({
        amount: swapRoute.amountOut,
        decimals: secondAmountDecimals
      });

      const newFirstAmount = firstAmountToDisplayDenom;
      const newSecondAmount = secondAmountToDisplayDenom;

      setFormState((current) => ({
        ...current,
        firstAmount: newFirstAmount,
        secondAmount: newSecondAmount
      }));

      isSwitching.current = false;
    } else {
      setFormState((current) => ({
        ...current,
        firstAmount: shouldUpdateFirstAmount
          ? amountToDisplayDenom
          : current.firstAmount,
        secondAmount: shouldUpdateSecondAmount
          ? amountToDisplayDenom
          : current.secondAmount
      }));
    }
  };

  const handleOnChangeFirstAmount = (amount: string) => {
    if (isSwitching.current) {
      return;
    }

    setFormState((current) => {
      const tokenInID = current.firstToken?.value;
      const tokenOutID = current.secondToken?.value;
      const hasBothTokens = tokenInID != null && tokenOutID != null;

      if (hasBothTokens) {
        lastInputTouched.current = InputTouchedEnum.first;

        const decimals = current.firstToken?.token.decimals;

        const amountIn = stringIsFloat(amount)
          ? parseAmount(amount, decimals)
          : '';

        getSwapRoute({
          amountIn,
          tokenInID,
          tokenOutID,
          tolerancePercentage
        });
      }

      return {
        ...current,
        firstAmount: amount,
        secondAmount: Boolean(amount) ? current.secondAmount : ''
      };
    });
  };

  const handleOnChangeSecondAmount = (amount: string) => {
    if (isSwitching.current) {
      return;
    }

    setFormState((current) => {
      const tokenInID = current.firstToken?.value;
      const tokenOutID = current.secondToken?.value;
      const hasBothTokens = tokenInID != null && tokenOutID != null;

      if (hasBothTokens) {
        lastInputTouched.current = InputTouchedEnum.second;

        const decimals = current.secondToken?.token.decimals;
        const amountOut = stringIsFloat(amount)
          ? parseAmount(amount, decimals)
          : '';

        getSwapRoute({
          amountOut,
          tokenInID,
          tokenOutID,
          tolerancePercentage
        });
      }

      return {
        ...current,
        secondAmount: amount,
        firstAmount: Boolean(amount) ? current.firstAmount : ''
      };
    });
  };

  const handleOnChangeFirstSelect = (option?: SelectOptionType) => {
    if (isSwitching.current) {
      return;
    }

    setFormState((current) => {
      const tokenInID = option?.value;
      const tokenOutID = current.secondToken?.value;
      const hasBothTokens = tokenInID != null && tokenOutID != null;

      if (hasBothTokens) {
        const isFixedInput = activeRoute?.swapType === FIXED_INPUT;

        const { amountIn, amountOut } = getCorrectAmountsOnTokenChange({
          activeRoute,
          newToken: option?.token,
          needsParsing: isFixedInput,
          firstToken: current.firstToken?.token,
          secondToken: current.secondToken?.token
        });

        getSwapRoute({
          amountIn,
          amountOut,
          tokenInID,
          tokenOutID,
          tolerancePercentage
        });
      }

      return {
        ...current,
        firstToken: option
      };
    });
  };

  const handleOnChangeSecondSelect = (option?: SelectOptionType) => {
    if (isSwitching.current) {
      return;
    }

    setFormState((current) => {
      const tokenOutID = option?.value;
      const tokenInID = current.firstToken?.value;
      const hasBothTokens = tokenInID != null && tokenOutID != null;

      if (hasBothTokens) {
        const isFixedOutput = activeRoute?.swapType === FIXED_OUTPUT;

        const { amountIn, amountOut } = getCorrectAmountsOnTokenChange({
          activeRoute,
          newToken: option?.token,
          needsParsing: isFixedOutput,
          firstToken: current.firstToken?.token,
          secondToken: current.secondToken?.token
        });

        getSwapRoute({
          amountIn,
          amountOut,
          tokenInID,
          tokenOutID,
          tolerancePercentage
        });
      }

      return {
        ...current,
        secondToken: option
      };
    });

    if (!prevSecondToken && formState.firstAmount) {
      handleOnChangeFirstAmount(formState.firstAmount);
    }
  };

  const handleOnFirstMaxBtnChange = () => {
    setFormState((current) => {
      const tokenInID = current.firstToken?.value;
      const tokenOutID = current.secondToken?.value;

      const balance = current.firstToken?.token.balance ?? '0';
      const amountIn =
        tokenInID === EGLD_IDENTIFIER ? getBalanceMinusDust(balance) : balance;

      const hasBothTokens = tokenInID != null && tokenOutID != null;

      if (!hasBothTokens) {
        return current;
      }

      const balanceDenom = meaningfulFormatAmount({
        amount: amountIn,
        decimals: current.firstToken?.token.decimals
      });

      getSwapRoute({
        amountIn,
        tokenInID,
        tokenOutID,
        tolerancePercentage
      });

      return {
        ...current,
        firstAmount: balanceDenom
      };
    });
  };

  const handleSwitchTokens = () => {
    setFormState((current) => {
      const newFirstAmount = current.secondAmount;
      const newSecondAmount = current.firstAmount;
      const newFirstToken = current.secondToken;
      const newSecondToken = current.firstToken;

      const tokenInID = newFirstToken?.value;
      const tokenOutID = newSecondToken?.value;
      const hasBothTokens = tokenInID != null && tokenOutID != null;

      if (!hasBothTokens) {
        return current;
      }

      // Enable switching only if we have input amounts otherwise, we enable only tokens switching (selects)
      isSwitching.current =
        lastInputTouched.current !== undefined &&
        newFirstAmount !== '' &&
        newSecondAmount !== '';

      const shouldGetAmountId =
        lastInputTouched.current === 'first' && newSecondAmount !== '';

      const shouldGetAmountOut =
        lastInputTouched.current === 'second' && newFirstAmount !== '';

      if (shouldGetAmountId) {
        lastInputTouched.current = InputTouchedEnum.second;

        const decimals = newSecondToken?.token.decimals;
        const amountOut = parseAmount(newSecondAmount, decimals);

        getSwapRoute({
          amountOut,
          tokenInID,
          tokenOutID,
          tolerancePercentage
        });
      }

      if (shouldGetAmountOut) {
        lastInputTouched.current = InputTouchedEnum.first;

        const decimals = newFirstToken?.token.decimals;
        const amountIn = parseAmount(newFirstAmount, decimals);

        getSwapRoute({
          amountIn,
          tokenInID,
          tokenOutID,
          tolerancePercentage
        });
      }

      return {
        firstAmount: newFirstAmount,
        secondAmount: newSecondAmount,
        firstToken: newFirstToken,
        secondToken: newSecondToken
      };
    });
  };

  return {
    firstToken: formState.firstToken,
    firstAmount: formState.firstAmount,
    secondToken: formState.secondToken,
    secondAmount: formState.secondAmount,
    activeRoute,
    inputAmountsUsdValue,
    handleOnChangeFirstAmount,
    handleOnChangeSecondAmount,
    handleOnChangeFirstSelect,
    handleOnChangeSecondSelect,
    handleOnChangeSwapRoute,
    handleOnFirstMaxBtnChange,
    handleSwitchTokens
  };
};
