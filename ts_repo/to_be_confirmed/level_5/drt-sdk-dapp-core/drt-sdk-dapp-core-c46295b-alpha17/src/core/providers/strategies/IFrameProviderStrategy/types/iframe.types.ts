import { IframeLoginTypes } from '@terradharitri/sdk-web-wallet-iframe-provider/out/constants';

export type IFrameProviderType = {
  type: IframeLoginTypes;
  address?: string;
};
