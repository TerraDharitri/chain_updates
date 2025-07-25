import type { CustomElementsDefineOptions } from '@terradharitri/sdk-dapp-core-ui/dist/loader';
export type { LedgerConnectPanel } from '@terradharitri/sdk-dapp-core-ui/dist/web-components/ledger-connect-panel';
export type { LedgerConnect } from '@terradharitri/sdk-dapp-core-ui/dist/web-components/ledger-connect';
export type {
  IAccountScreenData,
  IConfirmScreenData,
  IConnectScreenData,
  ILedgerAccount,
  ILedgerConnectPanelData
} from '@terradharitri/sdk-dapp-core-ui/dist/types/components/functional/ledger-connect-components/ledger-connect.types.d';
export type { SignTransactionsPanel } from '@terradharitri/sdk-dapp-core-ui/dist/web-components/sign-transactions-panel';
export type { WalletConnectPanel } from '@terradharitri/sdk-dapp-core-ui/dist/web-components/wallet-connect-panel';
export type { WalletConnect } from '@terradharitri/sdk-dapp-core-ui/dist/web-components/wallet-connect';
export type { IWalletConnectPanelData } from '@terradharitri/sdk-dapp-core-ui/dist/types/components/functional/wallet-connect-components/wallet-connect-panel.types.d';
export type { PendingTransactionsPanel } from '@terradharitri/sdk-dapp-core-ui/dist/web-components/pending-transactions-panel';
export { type IPendingTransactionsPanelData } from '@terradharitri/sdk-dapp-core-ui/dist/types/components/functional/pending-transactions-panel/pending-transactions-panel.types.d';
export type { NotificationsFeed } from '@terradharitri/sdk-dapp-core-ui/dist/web-components/notifications-feed';
export type { ToastList } from '@terradharitri/sdk-dapp-core-ui/dist/web-components/toast-list';
export type { IEventBus } from '@terradharitri/sdk-dapp-core-ui/dist/types/utils/EventBus';
export type {
  ITransactionListItem,
  ITransactionListItemAsset,
  ITransactionListItemAction
} from '@terradharitri/sdk-dapp-core-ui/dist/types/components/visual/transaction-list-item/transaction-list-item.types.d.ts';

export async function defineCustomElements(
  win?: Window,
  opts?: CustomElementsDefineOptions
): Promise<void> {
  try {
    const loader = await import('@terradharitri/sdk-dapp-core-ui/dist/loader');
    loader.defineCustomElements(win, opts);
  } catch (err) {
    throw new Error('@terradharitri/sdk-dapp-core-ui not found' + err);
  }
}
