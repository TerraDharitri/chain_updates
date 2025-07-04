import {
  ProviderTypeEnum,
  ProviderType
} from 'providers/types/providerFactory.types';

export const providerLabels: Record<ProviderType, string> = {
  [ProviderTypeEnum.crossWindow]: 'DharitrI Web Wallet',
  [ProviderTypeEnum.extension]: 'DharitrI Wallet Extension',
  [ProviderTypeEnum.walletConnect]: 'xPortal App',
  [ProviderTypeEnum.ledger]: 'Ledger',
  [ProviderTypeEnum.metamask]: 'MetaMask Snap',
  [ProviderTypeEnum.passkey]: 'Passkey',
  [ProviderTypeEnum.webview]: 'Webview',
  [ProviderTypeEnum.none]: ''
};
