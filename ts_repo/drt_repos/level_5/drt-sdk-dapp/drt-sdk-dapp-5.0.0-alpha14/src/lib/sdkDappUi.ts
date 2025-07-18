import type { CustomElementsDefineOptions } from '@terradharitri/sdk-dapp-ui/dist/loader';
export type { MvxLedgerConnect } from '@terradharitri/sdk-dapp-ui/dist/web-components/drt-ledger-connect';
export type { MvxSignTransactionsPanel } from '@terradharitri/sdk-dapp-ui/dist/web-components/drt-sign-transactions-panel';
export type { MvxWalletConnect } from '@terradharitri/sdk-dapp-ui/dist/web-components/drt-wallet-connect';
export type { MvxPendingTransactionsPanel } from '@terradharitri/sdk-dapp-ui/dist/web-components/drt-pending-transactions-panel';
export type { MvxNotificationsFeed } from '@terradharitri/sdk-dapp-ui/dist/web-components/drt-notifications-feed';
export type { MvxToastList } from '@terradharitri/sdk-dapp-ui/dist/web-components/drt-toast-list';
export type { MvxUnlockPanel } from '@terradharitri/sdk-dapp-ui/dist/web-components/drt-unlock-panel';
export type { IEventBus } from '@terradharitri/sdk-dapp-ui/dist/types/utils/EventBus';
export type {
  ITransactionListItem,
  ITransactionListItemAsset,
  ITransactionListItemAction
} from '@terradharitri/sdk-dapp-ui/dist/types/components/visual/transaction-list-item/transaction-list-item.types.d.ts';

export async function defineCustomElements(
  win?: Window,
  opts?: CustomElementsDefineOptions
): Promise<void> {
  try {
    const loader = await import('@terradharitri/sdk-dapp-ui/dist/loader');
    loader.defineCustomElements(win, opts);
  } catch (err) {
    throw new Error('@terradharitri/sdk-dapp-ui not found' + err);
  }
}
