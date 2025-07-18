import axios from 'axios';
import { CONFIG_ENDPOINT } from 'apiCalls/endpoints';
import { NetworkType } from 'types';

export async function getServerConfiguration(apiAddress: string) {
  const configUrl = `${apiAddress}/${CONFIG_ENDPOINT}`;

  try {
    const { data } = await axios.get<NetworkType>(configUrl);
    if (data != null) {
      // TODO: rewaDenomination will be removed from API when dapp-core v1 will be discontinued
      const rewaDenomination = 'rewaDenomination';
      if (rewaDenomination in data) {
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
      }
      return data;
    }
  } catch (err) {
    console.error('error fetching configuration for ', configUrl);
  }
  return null;
}
