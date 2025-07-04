import { useAuthorizationContext } from 'components/SwapAuthorizationProvider';
import { unwrapEgldQuery } from 'queries';
import { WrappingQueryResponseType } from 'types';
import { useLazyQueryWrapper } from './useLazyQueryWrapper';

export const useUnwrapEgld = () => {
  const { client } = useAuthorizationContext();

  if (!client) {
    throw new Error('Swap GraphQL client not initialized');
  }

  const { execute: unwrapEgld, isLoading } =
    useLazyQueryWrapper<WrappingQueryResponseType>({
      query: unwrapEgldQuery,
      queryOptions: {
        client
      }
    });

  return {
    unwrapEgld,
    isUnwrapEgldLoading: isLoading
  };
};
