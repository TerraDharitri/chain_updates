# sdk-dapp-swap

> A library to hold the main logic for swapping between tokens on the DharitrI blockchain

[![NPM](https://img.shields.io/npm/v/@terradharitri/sdk-dapp-swap.svg)](https://www.npmjs.com/package/@terradharitri/sdk-dapp-swap) [![JavaScript Style Guide](https://img.shields.io/badge/code_style-standard-brightgreen.svg)](https://standardjs.com)

# Installation

The library can be installed via npm or yarn.

```bash
npm install @terradharitri/sdk-dapp-swap
```

or

```bash
yarn add @terradharitri/sdk-dapp-swap
```

# Usage

sdk-dapp-swap proposes to simplify the process of implementing the swap functionality on DharitrI, providing easy to use utility functions and hooks for a specific application.

## Prerequisites

There are a couple of requirements that need to be met for the application to work properly.

**_If you experience bugs, please make sure that you read these, before opening an issue_**

<details>
  <summary>
      React
  </summary>

### React

This library was built for applications that use React, it might not be suitable for usage with other libraries or frameworks.

  </details>

<details>
  <summary>
    SwapAuthorizationProvider
 </summary>

### SwapAuthorizationProvider

You need to wrap your application with the **SwapAuthorizationProvider** component, which is exported by the library, as we need to create a global Context to be able to manipulate the data.

- import the Provider:

```typescript
import { SwapAuthorizationProvider } from '@terradharitri/sdk-dapp-swap/components/SwapAuthorizationProvider/SwapAuthorizationProvider';
or;
import { SwapAuthorizationProvider } from '@terradharitri/sdk-dapp-swap/components';
```

- Wrap your application with this Provider.

```jsx
<SwapAuthorizationProvider
    graphQLAddress={network.graphQlAddress}
    getAccessToken={getAccessToken}
    getAuthorizationHeaders={getAuthorizationHeaders} // optional
>
```

`graphQLAddress` is a required key that is needed to tell the Apollo Client on what graphQL address it should send the requests. Accepted values are any url string as well as the xExchange GraphQL service environment based endpoint (provided in `SwapGraphQLAddressEnum`).

`getAccessToken` is an async function that returns the accessToken mandatory for authorizing the requests.

SwapAuthorizationProvider also accepts an optional `getAuthorizationHeaders` async function that returns the headers object that will be provided on each request.
This allows the user to specify extra rules for the authorization of the requests.
`getAuthorizationHeaders` accepts the `requestParams` parameter which is of type `AuthorizationHeadersRequestParamsType` (also exported from the package).

  </details>

# Hooks

This area covers the hooks required to implement a basic swap functionality.

These hooks are exposed as named exports, which can be imported from sdk-dapp-swap.

```typescript
import { useTokens } from '@terradharitri/sdk-dapp-swap/hooks';
or;
import { useTokens } from '@terradharitri/sdk-dapp-swap/hooks/useTokens';

const {
  tokens,
  wrappedEgld,
  swapConfig,
  getTokens,
  isTokensLoading,
  isTokensError
} = useTokens(options);
```

The useTokens hook returns the state of the tokens based on whether you are authenticated or not. This hook accepts the `options` object in case the user wants to provide extra variables information for the query (see GetTokensType type).

```typescript
import { useSwapRoute } from '@terradharitri/sdk-dapp-swap/hooks';
or;
import { useSwapRoute } from '@terradharitri/sdk-dapp-swap/hooks/useSwapRoute';

const {
  getSwapRoute,
  swapRoute,
  isSwapRouteLoading,
  isAmountInLoading,
  isAmountOutLoading,
  swapRouteError,
  swapActionType
} = useSwapRoute({
  wrappedEgld,
  isPollingEnabled: true
});
```

The useSwapRoute hook returns the swapRoute required to compute the informations regarding both tokens provided in the fields. This hook accepts 2 optional parameters: `wrappedEgld` EsdtType object and a boolean parameter called `isPollingEnabled` to decide if the query should have polling or not.

```typescript
import { useSwapFormHandlers } from '@terradharitri/sdk-dapp-swap/hooks';
or;
import { useSwapFormHandlers } from '@terradharitri/sdk-dapp-swap/hooks/useSwapFormHandlers';

const {
  firstToken,
  firstAmount,
  secondToken,
  secondAmount,
  activeRoute,
  handleOnChangeFirstAmount,
  handleOnChangeSecondAmount,
  handleOnChangeFirstSelect,
  handleOnChangeSecondSelect,
  handleOnChangeSwapRoute,
  handleOnFirstMaxBtnChange,
  handleSwitchTokens
} = useSwapFormHandlers({
  getSwapRoute,
  tolerancePercentage
});
```

The useSwapFormHandlers hook returns the state of both the tokens and amounts, as well as the callbacks required for the inputs, selects, the maximum amount button and the button that triggers both tokens to switch. This hook accepts a mandatory parameter called `getSwapRoute` which is a function that comes from useSwapRoute (see useSwapRoute for reference) and an optional `tolerancePercentage`, which is the tolerance of the user account.

```typescript
import { useSwapInfo } from '@terradharitri/sdk-dapp-swap/hooks';
or;
import { useSwapInfo } from '@terradharitri/sdk-dapp-swap/hooks/useSwapInfo';

const {
  exchangeRate,
  feeDetails,
  feeAmounts,
  pricesImpact,
  feePercentages,
  tokenInId,
  tokenOutId,
  totalTransactionsFee,
  minimumAmountReceived,
  exchangeRateUsdValue,
  switchTokensDirection
} = useSwapInfo({
  tolerance,
  activeRoute,
  firstToken,
  firstAmount,
  secondToken,
  secondAmount,
  swapActionType
});
```

The useSwapInfo hook provides all the informations regarding the swap route, like exchange rate, details about fees, minimum amount received etc. This hook accepts a couple of parameters that implies both tokens and amounts, as well as the active swap route and account tolerance.

```typescript
import { useSwapValidationSchema } from '@terradharitri/sdk-dapp-swap/validation/hooks';
or;
import { useSwapValidationSchema } from '@terradharitri/sdk-dapp-swap/validation/hooks/useSwapValidationSchema';

const validationSchema = useSwapValidationSchema({
  firstToken,
  secondToken,
  minAcceptedAmount: minSwapAmount
});
```

The useSwapInfo hook provides the validation schema in order to facilitate a proper swap. This schema provides rules like insufficient balance, token / minimum token required etc. This hook accepts 2 mandatory `firstToken` and `secondToken` parameters, as well as an optional number called `minAcceptedAmount` which is the minimum amount required for the validation schema.

# sdk-dapp-swap exports

We recommend to import everything from its own file, in this way you won't be stressed about importing unused code.

You can either go into their specific folder in the module for extra trimming, or import everything together.

for example, these 2 imports are both valid:

```typescript
import {
  useSwapRoute,
  useSwapFormHandlers
} from '@terradharitri/sdk-dapp-swap/hooks';
```

and

```typescript
import { useSwapRoute } from '@terradharitri/sdk-dapp-swap/hooks/useSwapRoute';
import { useSwapFormHandlers } from '@terradharitri/sdk-dapp-swap/hooks/useSwapFormHandlers';
```

## constants exports

```typescript
import {
  FIXED_INPUT,
  FIXED_OUTPUT,
  EGLD_IDENTIFIER,
  DIGITS,
  MIN_EGLD_DUST,
  POLLING_INTERVAL
} from '@terradharitri/sdk-dapp-swap/constants/general';
```

## hooks exports

```typescript
import {
  useTokens,
  useWrapEgld,
  useSwapInfo,
  useSwapRoute,
  useUnwrapEgld,
  useQueryWrapper,
  useIsPageVisible,
  useRateCalculator,
  useLazyQueryWrapper,
  useSwapFormHandlers,
  useInputAmountUsdValue,
  useFetchMaintenanceFlag
} from '@terradharitri/sdk-dapp-swap/hooks';
```

## utils exports

```typescript
import {
  roundAmount,
  removeCommas,
  canParseAmount,
  getTokenDecimals,
  getPairFeeDetails,
  getTransactionFee,
  translateSwapError,
  getSwapActionType,
  getBalanceMinusDust,
  meaningfulFormatAmount,
  createTransactionFromRaw,
  calculateMinimumReceived,
  getSortedTokensByUsdValue
} from '@terradharitri/sdk-dapp-swap/utils';
```

### components exports

```typescript
import {
  SwapForm,
  SelectOpion,
  TokenSelect,
  SwapAuthorizationProvider
} from '@terradharitri/sdk-dapp-swap/components';
```

## Roadmap

See the [open issues](https://github.com/TerraDharitri/mx-sdk-dapp-swap/issues) for a list of proposed features (and known issues).

## Contributing

Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any contributions you make are **greatly appreciated**.

One can contribute by creating _pull requests_, or by opening _issues_ for discovered bugs or desired features.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

GPL-3.0-or-later
