import {
  ExplorerLink,
  FormatAmount,
  TransactionsTable
} from '@terradharitri/sdk-dapp-core-ui/dist/types/components';

declare module 'solid-js' {
  namespace JSX {
    interface IntrinsicElements {
      'drt-format-amount': FormatAmount;
      'drt-explorer-link': ExplorerLink;
      'drt-transactions-table': TransactionsTable;
    }
  }
}
