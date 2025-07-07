import { getEnvironmentForChainId } from '@terradharitri/sdk-dapp/apiCalls/configuration/getEnvironmentForChainId';
import { fallbackNetworkConfigurations } from '@terradharitri/sdk-dapp/constants/index';
import { EnvironmentsEnum } from '@terradharitri/sdk-dapp/types/enums.types';
import { useSelector } from 'react-redux';
import { activeNetworkSelector } from 'redux/selectors';

export const useGetEnvironment = () => {
  const { chainId, id: networkId } = useSelector(activeNetworkSelector);
  if (chainId) {
    const environmentForChainId = getEnvironmentForChainId(chainId);
    if (environmentForChainId) {
      return environmentForChainId;
    }
  }

  if (networkId) {
    const hasFallbackConfig = fallbackNetworkConfigurations?.[networkId];
    if (hasFallbackConfig) {
      return networkId as EnvironmentsEnum;
    }
  }

  return;
};
