import { h } from '@stencil/core';
import type { IProviderBase } from 'types/provider.types';
import { ProviderTypeEnum } from 'types/provider.types';

export const getProviderButtonIcon = (providerType: IProviderBase['type']): HTMLElement => {
  switch (providerType) {
    case ProviderTypeEnum.extension:
      return <drt-extension-provider-icon />;
    case ProviderTypeEnum.metamask:
      return <drt-metamask-provider-icon />;
    case ProviderTypeEnum.passkey:
      return <drt-passkey-provider-icon />;
    case ProviderTypeEnum.walletConnect:
      return <drt-dharitri-logo-icon />;
    case ProviderTypeEnum.ledger:
      return <drt-ledger-provider-icon />;
    case ProviderTypeEnum.crossWindow:
      return <drt-wallet-provider-icon />;

    default:
      return <drt-dharitri-logo-icon />;
  }
};
