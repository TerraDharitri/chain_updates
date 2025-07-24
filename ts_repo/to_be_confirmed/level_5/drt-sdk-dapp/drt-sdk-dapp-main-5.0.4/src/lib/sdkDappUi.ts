export { DrtLedgerConnect } from '@terradharitri/sdk-dapp-ui/web-components/drt-ledger-connect';
export type { DrtSignTransactionsPanel } from '@terradharitri/sdk-dapp-ui/web-components/drt-sign-transactions-panel';
export type { DrtWalletConnect } from '@terradharitri/sdk-dapp-ui/web-components/drt-wallet-connect';
export type { DrtPendingTransactionsPanel } from '@terradharitri/sdk-dapp-ui/web-components/drt-pending-transactions-panel';
export type { DrtNotificationsFeed } from '@terradharitri/sdk-dapp-ui/web-components/drt-notifications-feed';
export type { DrtToastList } from '@terradharitri/sdk-dapp-ui/web-components/drt-toast-list';
export type { DrtUnlockPanel } from '@terradharitri/sdk-dapp-ui/web-components/drt-unlock-panel';

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
