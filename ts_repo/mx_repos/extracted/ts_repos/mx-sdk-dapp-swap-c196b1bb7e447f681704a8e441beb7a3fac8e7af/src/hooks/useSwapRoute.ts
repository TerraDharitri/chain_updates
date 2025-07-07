import { useEffect, useMemo, useRef, useState } from 'react';
import { OperationVariables } from '@apollo/client';
import { useAuthorizationContext } from 'components/SwapAuthorizationProvider';
import { FIXED_INPUT, FIXED_OUTPUT } from 'constants/general';
import { IPlainTransactionObject } from 'lib';
import {
  swapQuery,
  wrapEgldQuery,
  unwrapEgldQuery,
  SwapRouteQueryResponseType,
  swapWithoutTransactionsQuery
} from 'queries';
import {
  EsdtType,
  SwapRouteType,
  SwapActionTypesEnum,
  WrappingQueryResponseType
} from 'types';
import { translateSwapError, getSwapActionType } from 'utils';
import { useQueryWrapper } from './useQueryWrapper';

export interface GetSwapRouteType {
  amountIn?: string;
  tokenInID: string;
  amountOut?: string;
  tokenOutID: string;
  tolerancePercentage?: number;
  wrappingAmount?: string; // used only by wrapping queries
}

export interface GetSwapRouteVariablesType
  extends Omit<GetSwapRouteType, 'tolerancePercentage'> {
  tolerance?: number;
}

export interface UseSwapRouteType {
  swapRouteError?: string;
  swapRoute?: SwapRouteType;
  isSwapRouteError?: boolean;
  isAmountInLoading: boolean;
  isSwapRouteLoading: boolean;
  isAmountOutLoading: boolean;
  transactions?: IPlainTransactionObject[];
  swapActionType?: SwapActionTypesEnum;
  previousFetchVariables: React.MutableRefObject<
    GetSwapRouteVariablesType | undefined
  >;
  getSwapRoute: (props: GetSwapRouteType) => void;
  refetch: (variables?: Partial<OperationVariables>) => void;
}

export const useSwapRoute = ({
  wrappedEgld,
  isPollingEnabled = false
}: {
  wrappedEgld?: EsdtType;
  isPollingEnabled?: boolean;
}): UseSwapRouteType => {
  const { client, isAuthenticated } = useAuthorizationContext();

  if (!client) {
    throw new Error('Swap GraphQL client not initialized');
  }

  const [variables, setVariables] = useState<GetSwapRouteType>();
  const [swapRoute, setSwapRoute] = useState<SwapRouteType>();
  const [swapRouteError, setSwapRouteError] = useState<string>();

  const previousFetchVariablesRef = useRef<GetSwapRouteVariablesType>();

  const swapActionType = useMemo(
    () =>
      getSwapActionType({
        firstTokenId: variables?.tokenInID,
        secondTokenId: variables?.tokenOutID,
        wrappedEgld
      }),
    [variables, wrappedEgld]
  );

  const query = useMemo(() => {
    if (swapActionType === SwapActionTypesEnum.wrap) return wrapEgldQuery;
    if (swapActionType === SwapActionTypesEnum.unwrap) return unwrapEgldQuery;

    return isAuthenticated ? swapQuery : swapWithoutTransactionsQuery;
  }, [isAuthenticated, swapActionType]);

  const skip = useMemo(() => {
    if (!variables) return true;

    const { amountIn, amountOut } = variables;
    const hasAmount = Boolean(amountIn ?? amountOut);

    if (!hasAmount) return true;

    return false;
  }, [variables]);

  const { data, error, refetch, isRefetching, isLoading, isError } =
    useQueryWrapper<SwapRouteQueryResponseType | WrappingQueryResponseType>({
      query,
      queryOptions: {
        skip,
        client,
        variables
      },
      isPollingEnabled
    });

  const handleOnCompleted = () => {
    if (!variables) {
      setSwapRoute(undefined);
      setSwapRouteError(translateSwapError(error?.message));
      return;
    }

    switch (swapActionType) {
      case SwapActionTypesEnum.wrap:
      case SwapActionTypesEnum.unwrap:
        const wrapTx = (data as WrappingQueryResponseType)?.wrapEgld;
        const unwrapTx = (data as WrappingQueryResponseType)?.unwrapEgld;
        const tx = wrapTx ?? unwrapTx;

        const swapType = variables.amountIn ? FIXED_INPUT : FIXED_OUTPUT;
        const amount = variables.amountIn ?? variables.amountOut;

        const wrappingSwapRoute: SwapRouteType = {
          amountIn: amount ?? '0',
          tokenInID: variables.tokenInID,
          tokenInPriceUSD: wrappedEgld?.price ?? '0',
          tokenInExchangeRateDenom: '1',
          maxPriceDeviationPercent: 0,
          tokensPriceDeviationPercent: 0,

          amountOut: amount ?? '0',
          tokenOutID: variables.tokenOutID,
          tokenOutPriceUSD: wrappedEgld?.price ?? '0',
          tokenOutExchangeRateDenom: '1',

          fees: [],
          swapType,
          tokenRoute: [],
          pricesImpact: [],

          intermediaryAmounts: [],
          pairs: [],
          transactions: tx ? [tx] : []
        };

        // ignore errors for wrapping because the endpoints don't work when not authenticated
        setSwapRouteError(undefined);
        setSwapRoute(wrappingSwapRoute);
        break;
      default:
        const swap = (data as SwapRouteQueryResponseType)?.swap;
        const translatedError = translateSwapError(error?.message);

        // we do not set partial routes
        setSwapRouteError(translatedError);
        setSwapRoute(translatedError ? undefined : swap);
    }
  };

  const getSwapRoute = ({
    amountIn,
    amountOut,
    tokenInID,
    tokenOutID,
    tolerancePercentage = 1
  }: GetSwapRouteType) => {
    if (!amountIn && !amountOut) {
      setVariables(undefined);
      return;
    }

    const guardedTolerancePercentage =
      tolerancePercentage < 0 || tolerancePercentage > 100
        ? 1
        : tolerancePercentage;

    const variables: GetSwapRouteVariablesType = {
      wrappingAmount: amountIn ?? amountOut,
      amountIn,
      amountOut,
      tokenInID,
      tokenOutID,
      tolerance: guardedTolerancePercentage / 100
    };

    setVariables(variables);
    previousFetchVariablesRef.current = variables;
  };

  useEffect(handleOnCompleted, [data, error, variables]);

  const isAmountOutLoading = Boolean(
    (isLoading || isRefetching) && previousFetchVariablesRef.current?.amountIn
  );

  const isAmountInLoading = Boolean(
    (isLoading || isRefetching) && previousFetchVariablesRef.current?.amountOut
  );

  return {
    refetch,
    getSwapRoute,
    swapRoute,
    swapActionType,
    swapRouteError,
    isAmountInLoading,
    isAmountOutLoading,
    isSwapRouteError: isError,
    isSwapRouteLoading: isLoading,
    transactions: swapRoute?.transactions,
    previousFetchVariables: previousFetchVariablesRef
  };
};
