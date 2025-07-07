import { useAuthorizationContext } from 'components/SwapAuthorizationProvider';
import { unwrapRewaQuery } from 'queries';
import { WrappingQueryResponseType } from 'types';
import { useLazyQueryWrapper } from './useLazyQueryWrapper';

export const useUnwrapRewa = () => {
  const { client } = useAuthorizationContext();

  if (!client) {
    throw new Error('Swap GraphQL client not initialized');
  }

  const { execute: unwrapRewa, isLoading } =
    useLazyQueryWrapper<WrappingQueryResponseType>({
      query: unwrapRewaQuery,
      queryOptions: {
        client
      }
    });

  return {
    unwrapRewa,
    isUnwrapRewaLoading: isLoading
  };
};
