import { IframeLoginTypes } from '@terradharitri/sdk-web-wallet-iframe-provider/out/constants';

export type IframeProviderType = {
  type: IframeLoginTypes;
  address?: string;
  walletUrl?: string;
};
