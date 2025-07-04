// types here need to be synced with the types in sdk-dapp providerFactory.constants.ts
import { ProviderTypeEnum } from 'types/provider.types';

export const providerLabels: Record<string, string> = {
  [ProviderTypeEnum.crossWindow]: 'DharitrI Web Wallet',
  [ProviderTypeEnum.extension]: 'DharitrI Wallet Extension',
  [ProviderTypeEnum.walletConnect]: 'xPortal App',
  [ProviderTypeEnum.ledger]: 'Ledger',
  [ProviderTypeEnum.metamask]: 'MetaMask Snap',
  [ProviderTypeEnum.passkey]: 'Passkey',
};
