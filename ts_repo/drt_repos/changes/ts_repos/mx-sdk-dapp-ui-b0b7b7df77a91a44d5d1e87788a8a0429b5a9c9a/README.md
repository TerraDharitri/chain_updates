# DharitrI UI library for Front-End DApps

DharitrI Front-End Library for JavaScript and TypeScript (written in TypeScript).

## Introduction

`sdk-dapp-ui` is a library that holds components to display user information from the DharitrI blockchain.

Since the library is built using [Stencil](https://stenciljs.com/), it can be used in any front-end framework, such as [React](https://github.com/TerraDharitri/mx-template-dapp), Angular, or [Solid.js](https://github.com/TerraDharitri/mx-solidjs-template-dapp), but also in back-end frameworks like [Next.js](https://github.com/TerraDharitri/mx-template-dapp-nextjs).

## GitHub project

The GitHub repository can be found here: [https://github.com/TerraDharitri/mx-sdk-dapp-ui](https://github.com/TerraDharitri/mx-sdk-dapp-ui)

## Live demo: template-dapp

See [Template dApp](https://template-dapp.dharitri.org/) for live demo or checkout usage in the [Github repo](https://github.com/TerraDharitri/mx-template-dapp)

## Requirements

- Node.js version 20.13.1+
- Npm version 10.5.2+

## Distribution

[npm](https://www.npmjs.com/package/@terradharitri/sdk-dapp-ui)

## Installation

The library can be installed via npm or yarn.

```bash
npm install @terradharitri/sdk-dapp-ui
```

or

```bash
yarn add @terradharitri/sdk-dapp-ui
```

## Usage

`sdk-dapp-ui` library is primarily designed to work with [@terradharitri/sdk-dapp](https://www.npmjs.com/package/@terradharitri/sdk-dapp), since components are designed to display data and emit user events, but do not hold any business logic.

The basic usage of the components would be importing the component and its corresponding interface and creating a wrapper for it in your application.

It outputs both web components (working out of the box) and React components (you need to create a wrapper for them).

### React Components import

```tsx
export {
  MvxCopyButton,
  MvxExplorerLink,
  MvxFormatAmount,
  MvxTransactionsTable,
  MvxUnlockButton,
} from '@terradharitri/sdk-dapp-ui/react';
```

The library is divided into three main categories of components:

1. The ones that only display data (visual)
2. The ones that display data provided by a controller (controlled)
3. The ones that are designed for user interaction (functional)

### 1. Visual Components

Visual components are the most basic building blocks that handle pure presentation. They are controlled through props and don't contain any business logic. These components focus on consistent styling and user interface elements.

Components:

- **Preloader** (`mvx-preloader`): A loading indicator for asynchronous operations
- **Font Awesome Icon** (`mvx-fa-icon`): Icon component with Font Awesome integration
- **Side Panel** (`mvx-side-panel`): Sliding panel with header and content sections
- **Tooltip** (`mvx-tooltip`): Contextual information display with hover/click activation
- **Transaction List Item**: Structured display of transaction information
- **Pagination** (`mvx-pagination`): Navigation controls for paginated content

### Visual Component Example

```tsx
export { getStore } from '@terradharitri/sdk-dapp/out/store/store';
export type { ExplorerLink as ExplorerLinkSDKPropsType } from '@terradharitri/sdk-dapp-ui/dist/types/components/visual/explorer-link/explorer-link.d.ts';
export { networkSelector } from '@terradharitri/sdk-dapp/out/store/selectors/networkSelectors';
import { IPropsWithClass, IPropsWithChildren } from 'types';

interface ExplorerLinkPropsType extends Partial<ExplorerLinkSDKPropsType> {
  'page': string;
  'class'?: string;
  'data-testid'?: string;
  'children'?: JSX.Element;
}

export const ExplorerLink = ({
  children,
  page,
  'class': className,
  'data-testid': dataTestId,
  ...rest
}: ExplorerLinkPropsType) => {
  const store = getStore();
  const network = networkSelector(store.getState());
  return (
    <mvx-explorer-link link={`${network.explorerAddress}${page}`} class={className} data-testid={dataTestId} {...rest}>
      {children ? <div>{children}</div> : null}
    </mvx-explorer-link>
  );
};
```

### 2. Controlled Components

Controlled components are designed to display data that is processed by a controller. They receive formatted data through props and focus on consistent data presentation. These components are typically used in data-heavy sections of the application.

Components:

- **Format Amount** (`mvx-format-amount`): Numerical amount formatting with validation
- **Transactions Table** (`mvx-transactions-table`): Structured display of transaction data

### Controlled Component Example

```tsx
import { TransactionsTableController } from '@terradharitri/sdk-dapp/out/controllers/TransactionsTableController';
import { accountSelector } from '@terradharitri/sdk-dapp/out/store/selectors/accountSelectors';
import { networkSelector } from '@terradharitri/sdk-dapp/out/store/selectors/networkSelectors';
import { getStore } from '@terradharitri/sdk-dapp/out/store/store';

export const TransactionsTable = () => {
  const store = getStore();
  const network = networkSelector(store.getState());
  const account = accountSelector(store.getState());

  const data = await TransactionsTableController.processTransactions({
    address: account().address,
    egldLabel: network().egldLabel,
    explorerAddress: network().explorerAddress,
    transactions: props.transactions || [],
  });

  return <mvx-transactions-table transactions={data} />;
};
```

### 3. Functional Components

Functional components handle specific application functionality and business logic. They integrate with the application's event system and manage user interactions. These components are typically used in complex workflows like authentication and transaction signing.

Components:

- **Sign Transactions Panel** (`mvx-sign-transactions-panel`): Transaction signing workflow
- **Notifications Feed** (`mvx-notifications-feed`): Transaction notifications and history
- **Wallet Connect** (`mvx-wallet-connect`): Wallet connection flow
- **Unlock Panel** (`mvx-unlock-panel`): Wallet authentication
- **Toast List** (`mvx-toast-list`): Notification management
- **Ledger Connect** (`mvx-ledger-connect`): Hardware wallet connection

You can check out the way these components are used in `@terradharitri/sdk-dapp` [here](https://github.com/TerraDharitri/mx-sdk-dapp/blob/main/src/managers/UnlockPanelManager/UnlockPanelManager.ts).

## Debugging your dApp

The recommended way to debug your application is by using [lerna](https://lerna.js.org/). Make sure you have the same package version in sdk-daap-core's package.json and in your project's package.json.

If you preffer to use [npm link](https://docs.npmjs.com/cli/v11/commands/npm-link), make sure to use the `preserveSymlinks` option in the server configuration:

```js
  resolve: {
    preserveSymlinks: true, // 👈
    alias: {
      src: "/src",
    },
  },
```

To build the library, run:

```bash
npm run build
```

To run the unit tests, run:

```bash
npm test
```

To run a specific test file in Stencil, run:

```bash
npx stencil test src/components/visual/transaction-list-item/tests/transaction-list-item.spec.tsx --spec
```

To run an individual test from a specific test file in Stencil, run:

```bash
npx stencil test src/components/visual/transaction-list-item/tests/transaction-list-item.spec.tsx --spec -t 'renders with asset icon'
```
