import { EnvironmentsEnum, NetworkType } from 'types';

export const fallbackNetworkConfigurations: Record<
  keyof typeof EnvironmentsEnum,
  NetworkType
> = {
  devnet: {
    id: 'devnet',
    chainId: 'D',
    name: 'Devnet',
    rewaLabel: 'xREWA',
    decimals: '18',
    digits: '4',
    gasPerDataByte: '1500',
    walletConnectDeepLink:
      'https://durian.page.link/?apn=com.numbat.durian.wallet&isi=1519405832&ibi=com.numbat.durian.wallet&link=https://xportal.com/',
    walletConnectBridgeAddresses: ['https://bridge.walletconnect.org'],
    walletConnectV2RelayAddresses: ['wss://relay.walletconnect.com'],
    walletAddress: 'https://devnet-wallet.dharitri.org',
    xAliasAddress: 'https://devnet.xalias.com',
    apiAddress: 'https://devnet-api.dharitri.org',
    extrasApiAddress: 'https://devnet-extras-api.dharitri.org',
    explorerAddress: 'http://devnet-explorer.dharitri.org',
    metamaskSnapWalletAddress: 'https://devnet-iframe-wallet.dharitri.org',
    apiTimeout: '4000',
    roundDuration: 6000
  },
  testnet: {
    id: 'testnet',
    chainId: 'T',
    name: 'Testnet',
    rewaLabel: 'xREWA',
    decimals: '18',
    digits: '4',
    gasPerDataByte: '1500',
    walletConnectDeepLink:
      'https://durian.page.link/?apn=com.numbat.durian.wallet&isi=1519405832&ibi=com.numbat.durian.wallet&link=https://xportal.com/',
    walletConnectBridgeAddresses: ['https://bridge.walletconnect.org'],
    walletConnectV2RelayAddresses: ['wss://relay.walletconnect.com'],
    walletAddress: 'https://testnet-wallet.dharitri.org',
    xAliasAddress: 'https://testnet.xalias.com',
    apiAddress: 'https://testnet-api.dharitri.org',
    extrasApiAddress: 'https://testnet-extras-api.dharitri.org',
    explorerAddress: 'http://testnet-explorer.dharitri.org',
    metamaskSnapWalletAddress: 'https://testnet-iframe-wallet.dharitri.org',
    apiTimeout: '4000',
    roundDuration: 6000
  },
  mainnet: {
    id: 'mainnet',
    chainId: '1',
    name: 'Mainnet',
    rewaLabel: 'REWA',
    decimals: '18',
    digits: '4',
    gasPerDataByte: '1500',
    walletConnectDeepLink:
      'https://durian.page.link/?apn=com.numbat.durian.wallet&isi=1519405832&ibi=com.numbat.durian.wallet&link=https://xportal.com/',
    walletConnectBridgeAddresses: ['https://bridge.walletconnect.org'],
    walletConnectV2RelayAddresses: ['wss://relay.walletconnect.com'],
    walletAddress: 'https://wallet.dharitri.org',
    xAliasAddress: 'https://xalias.com',
    apiAddress: 'https://api.dharitri.org',
    extrasApiAddress: 'https://extras-api.dharitri.org',
    explorerAddress: 'https://explorer.dharitri.org',
    metamaskSnapWalletAddress: 'https://iframe-wallet.dharitri.org',
    apiTimeout: '4000',
    roundDuration: 6000
  }
};

const { chainId: DEVNET_CHAIN_ID, rewaLabel: DEVNET_REWA_LABEL } =
  fallbackNetworkConfigurations[EnvironmentsEnum.devnet];
const { chainId: TESTNET_CHAIN_ID, rewaLabel: TESTNET_REWA_LABEL } =
  fallbackNetworkConfigurations[EnvironmentsEnum.testnet];
const { chainId: MAINNET_CHAIN_ID, rewaLabel: MAINNET_REWA_LABEL } =
  fallbackNetworkConfigurations[EnvironmentsEnum.mainnet];

export {
  DEVNET_CHAIN_ID,
  TESTNET_CHAIN_ID,
  MAINNET_CHAIN_ID,
  DEVNET_REWA_LABEL,
  TESTNET_REWA_LABEL,
  MAINNET_REWA_LABEL
};

export const chainIdByEnvironment: Record<EnvironmentsEnum, string> = {
  [EnvironmentsEnum.devnet]: DEVNET_CHAIN_ID,
  [EnvironmentsEnum.testnet]: TESTNET_CHAIN_ID,
  [EnvironmentsEnum.mainnet]: MAINNET_CHAIN_ID
};

export const chainIdToEnvironment: Record<string, EnvironmentsEnum> = {
  [DEVNET_CHAIN_ID]: EnvironmentsEnum.devnet,
  [TESTNET_CHAIN_ID]: EnvironmentsEnum.testnet,
  [MAINNET_CHAIN_ID]: EnvironmentsEnum.mainnet
};

export const DHARITRI_ORIGIN = 'dharitri://';
export const TIMEOUT = 3000;

/**
 * PPU value that is used when the user doesn't want a higher gas price for the transaction
 */
export const EMPTY_PPU = 0;
