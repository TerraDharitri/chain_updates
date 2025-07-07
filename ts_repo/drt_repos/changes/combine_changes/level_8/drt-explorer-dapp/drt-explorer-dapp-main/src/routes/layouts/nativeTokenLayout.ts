import { networks } from 'config';

import { NativeTokenLayout } from 'layouts/NativeTokenLayout';
import { NativeTokenAccounts } from 'pages/NativeToken/NativeTokenAccounts';
import { NativeTokenTransactions } from 'pages/NativeToken/NativeTokenTransactions';

import { TitledRouteObject } from '../routes';

export const nativeTokenLayout: TitledRouteObject[] = [];

networks.forEach((network) => {
  if (!network.rewaLabel) {
    return;
  }

  const networkPath = `/${network.rewaLabel?.toLowerCase()}`;

  const routeExists = nativeTokenLayout.find(
    (route) => route.path === networkPath
  );

  if (routeExists) {
    return;
  }

  nativeTokenLayout.push({
    path: networkPath,
    preventScroll: true,
    Component: NativeTokenLayout,
    children: [
      {
        path: networkPath,
        title: network.rewaLabel,
        preventScroll: true,
        Component: NativeTokenTransactions
      },
      {
        path: `${networkPath}/accounts`,
        title: 'Holders',
        preventScroll: true,
        Component: NativeTokenAccounts
      }
    ]
  });
});
