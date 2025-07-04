import { useAuthorizationContext } from 'components/SwapAuthorizationProvider';
import { wrapRewaQuery } from 'queries';
import { WrappingQueryResponseType } from 'types';
import { useLazyQueryWrapper } from './useLazyQueryWrapper';

export const useWrapRewa = () => {
  const { client } = useAuthorizationContext();

  if (!client) {
    throw new Error('Swap GraphQL client not initialized');
  }

  const { execute: wrapRewa, isLoading } =
    useLazyQueryWrapper<WrappingQueryResponseType>({
      query: wrapRewaQuery,
      queryOptions: {
        client
      }
    });

  return {
    wrapRewa,
    isWrapRewaLoading: isLoading
  };
};
