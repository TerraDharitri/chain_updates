import type { CustomElementsDefineOptions } from '@terradharitri/sdk-dapp-core-ui/dist/loader';
export type { MvxLedgerFlow } from '@terradharitri/sdk-dapp-core-ui/dist/web-components/drt-ledger-flow';
export type { MvxSignTransactionsPanel } from '@terradharitri/sdk-dapp-core-ui/dist/web-components/drt-sign-transactions-panel';
export type { MvxWalletConnectProvider } from '@terradharitri/sdk-dapp-core-ui/dist/web-components/drt-wallet-connect-provider';
export type { MvxPendingTransactionsPanel } from '@terradharitri/sdk-dapp-core-ui/dist/web-components/drt-pending-transactions-panel';
export type { IPendingTransactionsPanelData } from '@terradharitri/sdk-dapp-core-ui/dist/types/components/functional/pending-transactions-panel/pending-transactions-panel.types.d';
export type { MvxNotificationsFeed } from '@terradharitri/sdk-dapp-core-ui/dist/web-components/drt-notifications-feed';
export type { MvxToastList } from '@terradharitri/sdk-dapp-core-ui/dist/web-components/drt-toast-list';
export type { MvxUnlockPanel } from '@terradharitri/sdk-dapp-core-ui/dist/web-components/drt-unlock-panel';
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
