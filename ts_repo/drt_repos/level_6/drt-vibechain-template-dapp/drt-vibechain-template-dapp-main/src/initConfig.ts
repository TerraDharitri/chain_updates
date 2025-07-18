import {
  EnvironmentsEnum,
  ICustomProvider,
  InitAppType,
  ProviderTypeEnum,
} from "./lib";
import { walletConnectV2ProjectId } from "./config";

const ADDITIONAL_PROVIDERS = {} as const;

export const ExtendedProviders = {
  ...ProviderTypeEnum,
  ...ADDITIONAL_PROVIDERS,
} as const;

const DEFAULT_TOAST_LIEFTIME = 5000;

const providers: ICustomProvider<ProviderTypeEnum>[] = [];

(window as any).dharitri = {};
// Option 1: Add providers using the `window.providers` array
(window as any).dharitri.providers = providers;

export const config: InitAppType = {
  storage: { getStorageCallback: () => sessionStorage },
  dAppConfig: {
    nativeAuth: true,
    environment: EnvironmentsEnum.devnet,
    network: {
      walletAddress: "https://devnet-wallet.dharitri.org",
    },
    providers: {
      walletConnect: {
        walletConnectV2ProjectId,
      },
    },
    successfulToastLifetime: DEFAULT_TOAST_LIEFTIME,
  },

  // Option 2: Add providers using the config `customProviders` array
  // customProviders: [customWalletProvider]
};
