import { CollectionLayout } from 'layouts/CollectionLayout';
import { NftLayout } from 'layouts/NftLayout';
import { TokenLayout } from 'layouts/TokenLayout';
import { CollectionRoles } from 'pages/CollectionDetails/CollectionRoles';
import { CollectionTransactions } from 'pages/CollectionDetails/CollectionTransactions';
import { NftAccounts } from 'pages/NftDetails/NftAccounts';
import { NftTransactions } from 'pages/NftDetails/NftTransactions';
import { TokenDetailsAccounts } from 'pages/TokenDetails/TokenAccounts';
import { TokenDetailsLockedAccounts } from 'pages/TokenDetails/TokenLockedAccounts';
import { TokenDetailsRoles } from 'pages/TokenDetails/TokenRoles';
import { TokenTransactions } from 'pages/TokenDetails/TokenTransactions';
import { Tokens } from 'pages/Tokens';
import { TokensMeta } from 'pages/TokensMeta';

import { TitledRouteObject } from '../routes';

export const tokensRoutes = {
  tokens: '/tokens',
  tokensMeta: '/meta-tokens',
  tokensMetaDcdt: '/meta-dcdt',
  tokensMetaDcdtDetails: '/meta-dcdt/:hash',
  tokensMetaDcdtDetailsRoles: '/meta-dcdt/:hash/roles',
  tokensProofDetails: '/proofs/:hash',
  tokensProofDetailsAccounts: '/proofs/:hash/accounts',
  tokenDetails: '/tokens/:hash',
  tokenDetailsAccounts: '/tokens/:hash/accounts',
  tokenDetailsLockedAccounts: '/tokens/:hash/locked-accounts',
  tokenDetailsRoles: '/tokens/:hash/roles'
};

export const tokenLayout: TitledRouteObject[] = [
  {
    path: tokensRoutes.tokensMeta,
    title: 'Meta-DCDT Tokens',
    preventScroll: true,
    Component: TokensMeta
  },
  {
    path: tokensRoutes.tokensMetaDcdt,
    title: 'Meta-DCDT Tokens',
    preventScroll: true,
    Component: TokensMeta
  },
  {
    path: tokensRoutes.tokensMetaDcdtDetails,
    preventScroll: true,
    Component: CollectionLayout,
    children: [
      {
        path: tokensRoutes.tokensMetaDcdtDetails,
        title: 'Meta-DCDT Transactions',
        preventScroll: true,
        Component: CollectionTransactions
      },
      {
        path: tokensRoutes.tokensMetaDcdtDetailsRoles,
        title: 'Meta-DCDT Roles',
        preventScroll: true,
        Component: CollectionRoles
      }
    ]
  },
  {
    path: tokensRoutes.tokensProofDetails,
    preventScroll: true,
    Component: NftLayout,
    children: [
      {
        path: tokensRoutes.tokensProofDetails,
        title: 'Proof Transactions',
        preventScroll: true,
        Component: NftTransactions
      },
      {
        path: tokensRoutes.tokensProofDetailsAccounts,
        title: 'Proof Holders',
        preventScroll: true,
        Component: NftAccounts
      }
    ]
  },
  {
    path: tokensRoutes.tokens,
    title: 'Tokens',
    preventScroll: true,
    Component: Tokens
  },
  {
    path: tokensRoutes.tokenDetails,
    preventScroll: true,
    Component: TokenLayout,
    children: [
      {
        path: tokensRoutes.tokenDetails,
        title: 'Token Details',
        preventScroll: true,
        Component: TokenTransactions
      },
      {
        path: tokensRoutes.tokenDetailsAccounts,
        title: 'Token Holders',
        preventScroll: true,
        Component: TokenDetailsAccounts
      },
      {
        path: tokensRoutes.tokenDetailsLockedAccounts,
        title: 'Locked Token Accounts',
        preventScroll: true,
        Component: TokenDetailsLockedAccounts
      },
      {
        path: tokensRoutes.tokenDetailsRoles,
        title: 'Token Roles',
        preventScroll: true,
        Component: TokenDetailsRoles
      }
    ]
  }
];
