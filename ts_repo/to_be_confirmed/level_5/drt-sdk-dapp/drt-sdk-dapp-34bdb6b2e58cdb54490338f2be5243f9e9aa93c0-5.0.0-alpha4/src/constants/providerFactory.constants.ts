import { ProviderTypeEnum } from 'providers/types/providerFactory.types';

export const providerLabels: Record<string, string> = {
  [ProviderTypeEnum.crossWindow]: 'Dharitri Web Wallet',
  [ProviderTypeEnum.extension]: 'Dharitri Wallet Extension',
  [ProviderTypeEnum.walletConnect]: 'xPortal App',
  [ProviderTypeEnum.ledger]: 'Ledger',
  [ProviderTypeEnum.metamask]: 'MetaMask Snap',
  [ProviderTypeEnum.passkey]: 'Passkey',
  [ProviderTypeEnum.none]: ''
};
