import { useEffect, useMemo, useRef, useState } from 'react';
import { useAuthorizationContext } from 'components/SwapAuthorizationProvider';
import {
  FilteredTokensQueryType,
  GET_FILTERED_TOKENS,
  GET_FILTERED_TOKENS_AND_BALANCE
} from 'queries';
import {
  DcdtType,
  FactoryType,
  TokensPaginationType,
  UserDcdtType
} from 'types';
import { getSortedTokensByUsdValue, mergeTokens } from 'utils';
import { useFetchTokenPrices } from './useFetchTokenPrices';
import { useIntersectionObserver } from './useIntersectionObserver';
import { useLazyQueryWrapper } from './useLazyQueryWrapper';

const DEFAULT_OFFSET = 0;
const DEFAULT_LIMIT = 500;
const DEFAULT_ENABLED_SWAPS = true;
const DEFAULT_PRICE_POLLING = false;
const DEFAULT_IDENTIFIERS: string[] = [];
const DEFAULT_PAGINATION: TokensPaginationType = { first: 20, after: '' };
const DEFAULT_SEARCH_INPUT = '';

interface GetTokensType {
  limit?: number;
  offset?: number;
  identifiers?: string[];
  enabledSwaps?: boolean;
  pagination?: TokensPaginationType;
  searchInput?: string;
}

interface UseTokensType {
  pricePolling?: boolean;
  observerId?: string;
  searchInput?: string;
  identifiers?: string[];
}

export const useFilteredTokens = (options?: UseTokensType) => {
  const isInitialLoad = useRef(true);
  const { client, isAuthenticated } = useAuthorizationContext();

  if (!client) {
    throw new Error('Swap GraphQL client not initialized');
  }

  const pricePolling = options?.pricePolling ?? DEFAULT_PRICE_POLLING;
  const searchInput = options?.searchInput;

  const [pagination, setPagination] = useState<TokensPaginationType>({
    first: 20,
    after: ''
  });
  const [hasMore, setHasMore] = useState(true);
  const [currentCursor, setCurrentCursor] = useState<string>();
  const [loadedCursors, setLoadedCursors] = useState<Set<string>>(new Set());

  const [tokens, setTokens] = useState<UserDcdtType[]>([]);
  const [wrappedRewa, setWrappedRewa] = useState<DcdtType>();
  const [swapConfig, setSwapConfig] = useState<FactoryType>();
  const [tokensCount, setTokensCount] = useState<number>();
  let ignoreNextHasMore = false;

  const { tokenPrices } = useFetchTokenPrices({
    isPollingEnabled: pricePolling
  });

  const handleOnCompleted = (data?: FilteredTokensQueryType | null) => {
    if (!data) return;

    const { wrappingInfo, userTokens, factory, filteredTokens } = data;
    const { edges, pageInfo, pageData } = filteredTokens;

    setTokensCount(pageData?.count);
    if (factory) setSwapConfig(factory);

    const newWrappedRewa =
      wrappingInfo && wrappingInfo.length
        ? wrappingInfo[0].wrappedToken
        : undefined;
    setWrappedRewa(newWrappedRewa);

    if (!edges) return;

    setCurrentCursor(edges[edges.length - 1]?.cursor);
    const tokensWithBalance: UserDcdtType[] = edges.map((token) => ({
      ...token.node,
      balance: '0',
      valueUSD: '0'
    }));

    const mergedTokens = mergeTokens(tokensWithBalance, userTokens);
    const sortedTokensWithBalance = getSortedTokensByUsdValue({
      tokens: mergedTokens,
      wrappedRewa: newWrappedRewa
    });

    setTokens((prevTokens) => mergeTokens(prevTokens, sortedTokensWithBalance));

    if (!pageInfo?.hasNextPage && !ignoreNextHasMore) {
      setHasMore(false);
    } else {
      ignoreNextHasMore = false;
    }
  };

  const {
    isError,
    isLoading,
    execute: getTokensTrigger
  } = useLazyQueryWrapper<FilteredTokensQueryType>({
    query: isAuthenticated
      ? GET_FILTERED_TOKENS_AND_BALANCE
      : GET_FILTERED_TOKENS,
    queryOptions: {
      client,
      onCompleted: handleOnCompleted
    }
  });

  const getTokens = (options?: GetTokensType, continueFetching?: boolean) => {
    const variables = {
      limit: options?.limit ?? DEFAULT_LIMIT,
      offset: options?.offset ?? DEFAULT_OFFSET,
      identifiers: options?.identifiers ?? DEFAULT_IDENTIFIERS,
      enabledSwaps: options?.enabledSwaps ?? DEFAULT_ENABLED_SWAPS,
      pagination: options?.pagination ?? DEFAULT_PAGINATION,
      searchInput: options?.searchInput ?? DEFAULT_SEARCH_INPUT
    };

    getTokensTrigger({
      variables
    });

    if (continueFetching) {
      ignoreNextHasMore = true;
    }
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

  useEffect(() => {
    if (isInitialLoad.current) {
      isInitialLoad.current = false;
      return;
    }
    setPagination({ first: 20, after: '' });
    setLoadedCursors(new Set());
    setHasMore(true);
    getTokens({ pagination: { first: 20, after: '' }, searchInput });
  }, [searchInput]);

  useEffect(() => {
    if (pagination.after) {
      setLoadedCursors((prev) => new Set(prev).add(pagination.after as string));
      getTokens({ pagination, searchInput });
    }
  }, [pagination]);

  useIntersectionObserver({
    tokens,
    hasMore,
    isLoading: isLoading ?? false,
    observerId: options?.observerId ?? '',
    loadedCursors,
    currentCursor: currentCursor ?? '',
    setPagination
  });

  return {
    swapConfig,
    wrappedRewa,
    isTokensError: isError,
    isTokensLoading: isLoading,
    tokens: tokensWithUpdatedPrice,
    totalTokensCount: tokensCount,
    getTokens,
    refetch: getTokensTrigger
  };
};
