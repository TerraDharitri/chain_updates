import { CONFIG_ENDPOINT } from 'apiCalls/endpoints';
import { NetworkType } from 'types';
import { useApiFetch } from '../utils';

export const useGetServerConfigurationFromApi = () => {
  const { data, isLoading, error } = useApiFetch({
    apiEndpoint: CONFIG_ENDPOINT
  });

  if (isLoading) {
    return null;
  }

  if (error) {
    console.error('Error fetching server configuration:', error);
  }

  const rewaDenomination = 'rewaDenomination';

  if (!(rewaDenomination in data)) {
    return data;
  }

  const {
    [rewaDenomination]: decimals,
    decimals: digits,
    ...rest
  } = data as NetworkType & {
    [rewaDenomination]: string;
  };

  const networkConfig: NetworkType = {
    ...rest,
    decimals,
    digits
  };

  return networkConfig;
};
