import { gql } from '@apollo/client';
import { useQueryWrapper } from 'hooks';
import { DcdtType } from 'types';

export const useFetchTokenPrices = ({
  skip = false,
  isPollingEnabled = false
}: {
  skip?: boolean;
  isPollingEnabled?: boolean;
}) => {
  const tokensQuery = gql`
    query swapPackageTokenPrices {
      tokens(enabledSwaps: true) {
        price
        identifier
      }
    }
  `;

  const { data, isError, isLoading } = useQueryWrapper<{ tokens: DcdtType[] }>({
    isPollingEnabled,
    query: tokensQuery,
    queryOptions: {
      skip
    }
  });

  return {
    tokenPrices: data?.tokens,
    isTokenPricesError: isError,
    isTokenPricesLoading: isLoading
  };
};
