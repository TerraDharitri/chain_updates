import type { CustomElementsDefineOptions } from '@terradharitri/sdk-dapp-ui/dist/loader';
export type { DrtLedgerConnect } from '@terradharitri/sdk-dapp-ui/dist/web-components/drt-ledger-connect';
export type { DrtSignTransactionsPanel } from '@terradharitri/sdk-dapp-ui/dist/web-components/drt-sign-transactions-panel';
export type { DrtWalletConnect } from '@terradharitri/sdk-dapp-ui/dist/web-components/drt-wallet-connect';
export type { DrtPendingTransactionsPanel } from '@terradharitri/sdk-dapp-ui/dist/web-components/drt-pending-transactions-panel';
export type { DrtNotificationsFeed } from '@terradharitri/sdk-dapp-ui/dist/web-components/drt-notifications-feed';
export type { DrtToastList } from '@terradharitri/sdk-dapp-ui/dist/web-components/drt-toast-list';
export type { DrtUnlockPanel } from '@terradharitri/sdk-dapp-ui/dist/web-components/drt-unlock-panel';
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
