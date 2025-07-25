import { useCallback, useEffect } from 'react';
import { useQuery, DocumentNode, QueryHookOptions } from '@apollo/client';
import { useAuthorizationContext } from 'components/SwapAuthorizationProvider';
import { POLLING_INTERVAL } from 'constants/index';
import { useIsPageVisible } from 'hooks';

export const useQueryWrapper = <TData>({
  query,
  queryOptions,
  refetchTrigger,
  isPollingEnabled = false,
  isRefetchEnabled = false
}: {
  query: DocumentNode;
  refetchTrigger?: number;
  isPollingEnabled?: boolean;
  isRefetchEnabled?: boolean;
  queryOptions?: QueryHookOptions<TData>;
}) => {
  const isPageVisible = useIsPageVisible();
  const { client } = useAuthorizationContext();

  const {
    error,
    loading,
    previousData,
    data = previousData,
    refetch,
    stopPolling,
    startPolling,
    ...rest
  } = useQuery<TData>(query, {
    client,
    notifyOnNetworkStatusChange: true,
    fetchPolicy: 'no-cache', // used for first run
    nextFetchPolicy: 'no-cache', // used for subsequent runs
    ...queryOptions
  });

  // listening on queryOptions resets the polling interval -> posibile race condition fix
  const startPollingCallback = useCallback(() => {
    if (isPageVisible && isPollingEnabled && !error) {
      startPolling(POLLING_INTERVAL);
    } else {
      stopPolling();
    }
  }, [isPageVisible, isPollingEnabled, error, queryOptions]);

  // mount and unmount
  useEffect(() => {
    startPollingCallback();

    return () => {
      stopPolling();
    };
  }, [startPollingCallback]);

  useEffect(() => {
    const isInitialRenderRefetch = refetchTrigger === 0;

    if (!isRefetchEnabled || isInitialRenderRefetch) {
      return;
    }

    refetch();
    startPollingCallback();
  }, [refetchTrigger, isRefetchEnabled]);

  const isLoading = data == null && loading;
  const isError = Boolean(error);
  const isRefetching = loading;

  return {
    data: isError ? undefined : data,
    error,
    isError,
    isLoading,
    isRefetching,
    refetch,
    ...rest
  };
};
