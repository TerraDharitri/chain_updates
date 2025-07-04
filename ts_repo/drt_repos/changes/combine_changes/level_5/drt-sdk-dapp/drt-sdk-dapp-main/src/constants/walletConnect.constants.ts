import { WalletConnectConfig } from 'providers/strategies/WalletConnectProviderStrategy/types/walletConnect.types';

export const fallbackWalletConnectConfigurations: WalletConnectConfig = {
  walletConnectV2ProjectId: '',
  walletConnectDeepLink:
    'https://durian.page.link/?apn=com.numbat.durian.wallet&isi=1519405832&ibi=com.numbat.durian.wallet&link=https://xportal.com/',
  walletConnectV2RelayAddress: 'wss://relay.walletconnect.com'
};
