import { useMemo, useState } from 'react';
import { useAuthorizationContext } from 'components/SwapAuthorizationProvider';
import { GET_TOKENS, GET_TOKENS_AND_BALANCE, TokensType } from 'queries';
import { DcdtType, FactoryType, UserDcdtType } from 'types';
import { getSortedTokensByUsdValue } from 'utils';
import { useFetchTokenPrices } from './useFetchTokenPrices';
import { useLazyQueryWrapper } from './useLazyQueryWrapper';

const DEFAULT_OFFSET = 0;
const DEFAULT_LIMIT = 500;
const DEFAULT_ENABLED_SWAPS = true;
const DEFAULT_PRICE_POLLING = false;
const DEFAULT_IDENTIFIERS: string[] = [];

interface GetTokensType {
  limit?: number;
  offset?: number;
  identifiers?: string[];
  enabledSwaps?: boolean;
}

interface UseTokensType {
  pricePolling?: boolean;
}

export const useTokens = (options?: UseTokensType) => {
  const { client, isAuthenticated } = useAuthorizationContext();

  if (!client) {
    throw new Error('Swap GraphQL client not initialized');
  }

  const pricePolling = options?.pricePolling ?? DEFAULT_PRICE_POLLING;

  const [tokens, setTokens] = useState<UserDcdtType[]>([]);
  const [wrappedRewa, setWrappedRewa] = useState<DcdtType>();
  const [swapConfig, setSwapConfig] = useState<FactoryType>();

  const { tokenPrices } = useFetchTokenPrices({
    isPollingEnabled: pricePolling
  });

  const handleOnCompleted = (data?: TokensType | null) => {
    if (!data) {
      return;
    }

    const { tokens: swapTokens, wrappingInfo, userTokens, factory } = data;

    if (factory) {
      setSwapConfig(factory);
    }

    const newWrappedRewa =
      wrappingInfo && wrappingInfo.length
        ? wrappingInfo[0].wrappedToken
        : undefined;

    setWrappedRewa(newWrappedRewa);

    if (!swapTokens) {
      setTokens([]);
      return;
    }

    const tokensWithBalance: UserDcdtType[] = swapTokens.map((token) => {
      const tokenFound = userTokens?.find(
        ({ identifier }) => identifier === token.identifier
      );

      return {
        ...token,
        balance: tokenFound?.balance ?? '0',
        valueUSD: tokenFound?.valueUSD ?? '0'
      };
    });

    const sortedTokensWithBalance = getSortedTokensByUsdValue({
      tokens: tokensWithBalance,
      wrappedRewa: newWrappedRewa
    });

    setTokens(sortedTokensWithBalance);
  };

  const {
    isError,
    isLoading,
    execute: getTokensTrigger
  } = useLazyQueryWrapper<TokensType>({
    query: isAuthenticated ? GET_TOKENS_AND_BALANCE : GET_TOKENS,
    queryOptions: {
      client,
      onCompleted: handleOnCompleted
    }
  });

  const getTokens = (options?: GetTokensType) => {
    const variables = {
      limit: options?.limit ?? DEFAULT_LIMIT,
      offset: options?.offset ?? DEFAULT_OFFSET,
      identifiers: options?.identifiers ?? DEFAULT_IDENTIFIERS,
      enabledSwaps: options?.enabledSwaps ?? DEFAULT_ENABLED_SWAPS
    };

    getTokensTrigger({
      variables
    });
  };

  const tokensWithUpdatedPrice = useMemo(
    () =>
      tokens.map((token) => {
        const tokenPrice = tokenPrices?.find(
          ({ identifier }) => identifier === token.identifier
        )?.price;

        return {
          ...token,
          price: tokenPrice ?? token.price
        };
      }),
    [tokens, tokenPrices]
  );

  return {
    swapConfig,
    wrappedRewa,
    isTokensError: isError,
    isTokensLoading: isLoading,
    tokens: tokensWithUpdatedPrice,
    getTokens,
    refetch: getTokensTrigger
  };
};
