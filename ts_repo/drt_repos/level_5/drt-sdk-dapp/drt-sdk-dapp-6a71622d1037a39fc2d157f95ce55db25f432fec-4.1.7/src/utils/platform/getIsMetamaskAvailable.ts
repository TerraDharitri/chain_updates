import { MetamaskProvider } from '@terradharitri/sdk-metamask-provider/out/metamaskProvider';

export const getIsMetamaskAvailable = () => {
  return MetamaskProvider.isMetamaskInstalled();
};
