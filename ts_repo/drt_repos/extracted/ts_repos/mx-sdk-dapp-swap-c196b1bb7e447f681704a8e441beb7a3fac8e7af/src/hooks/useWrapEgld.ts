import { useAuthorizationContext } from 'components/SwapAuthorizationProvider';
import { wrapEgldQuery } from 'queries';
import { WrappingQueryResponseType } from 'types';
import { useLazyQueryWrapper } from './useLazyQueryWrapper';

export const useWrapEgld = () => {
  const { client } = useAuthorizationContext();

  if (!client) {
    throw new Error('Swap GraphQL client not initialized');
  }

  const { execute: wrapEgld, isLoading } =
    useLazyQueryWrapper<WrappingQueryResponseType>({
      query: wrapEgldQuery,
      queryOptions: {
        client
      }
    });

  return {
    wrapEgld,
    isWrapEgldLoading: isLoading
  };
};
