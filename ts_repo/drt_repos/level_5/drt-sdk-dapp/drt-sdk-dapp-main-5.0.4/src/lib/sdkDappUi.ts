export { MvxLedgerConnect } from '@terradharitri/sdk-dapp-ui/web-components/drt-ledger-connect';
export type { MvxSignTransactionsPanel } from '@terradharitri/sdk-dapp-ui/web-components/drt-sign-transactions-panel';
export type { MvxWalletConnect } from '@terradharitri/sdk-dapp-ui/web-components/drt-wallet-connect';
export type { MvxPendingTransactionsPanel } from '@terradharitri/sdk-dapp-ui/web-components/drt-pending-transactions-panel';
export type { MvxNotificationsFeed } from '@terradharitri/sdk-dapp-ui/web-components/drt-notifications-feed';
export type { MvxToastList } from '@terradharitri/sdk-dapp-ui/web-components/drt-toast-list';
export type { MvxUnlockPanel } from '@terradharitri/sdk-dapp-ui/web-components/drt-unlock-panel';

export type { IEventBus } from '@terradharitri/sdk-dapp-ui/types/utils/EventBus';
export type {
  ITransactionListItem,
  ITransactionListItemAsset,
  ITransactionListItemAction
} from '@terradharitri/sdk-dapp-ui/types/components/visual/transaction-list-item/transaction-list-item.types';

export async function defineCustomElements(opts?: any): Promise<void> {
  try {
    const loader = await import('@terradharitri/sdk-dapp-ui');
    loader.defineCustomElements(opts);
  } catch (err) {
    throw new Error('@terradharitri/sdk-dapp-ui not found' + err);
  }
}
