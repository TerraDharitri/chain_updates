import { Config } from '@stencil/core';
import { sass } from '@stencil/sass';
import { reactOutputTarget } from '@stencil/react-output-target';
import nodePolyfills from 'rollup-plugin-node-polyfills';
import tailwind from 'stencil-tailwind-plugin';

const excludeComponents = [
  'drt-sign-transactions-panel',
  'drt-transaction-fee-component',
  'drt-pending-transactions-panel',
  'drt-ledger-connect-panel',
  'drt-ledger-flow',
  'drt-ledger-account-screen',
  'drt-ledger-connect-screen',
  'drt-ledger-confirm-screen',
  'drt-toast-list',
  'drt-generic-toast',
  'drt-custom-toast',
  'drt-simple-toast',
  'drt-transaction-toast-details-body',
  'drt-transaction-toast-details',
  'drt-transaction-toast-content',
  'drt-transaction-toast',
  'drt-transaction-toast-wrapper',
  'drt-sign-transaction-component',
  'drt-wallet-connect-provider',
  'drt-wallet-connect-panel',
  'drt-transaction-toast-progress',
  'drt-token-component',
  'drt-fungible-component',
  'drt-balance-component',
  'drt-unlock-panel',
];

export const config: Config = {
  namespace: 'sdk-dapp-core-ui',
  globalScript: './src/global/scripts/fonts-loader.ts',
  plugins: [
    sass(),
    tailwind({
      tailwindCssPath: './src/global/tailwind.css',
    }),
  ],
  outputTargets: [
    reactOutputTarget({
      outDir: './dist/react',
      stencilPackageName: '../../dist/types',
      customElementsDir: '../web-components',
      excludeComponents,
    }),
    {
      type: 'dist-custom-elements',
      externalRuntime: false,
      generateTypeDeclarations: true,
      dir: './dist/web-components',
    },
    {
      type: 'dist',
      copy: [{ src: 'assets', dest: 'assets' }],
      esmLoaderPath: './loader',
    },
    // this is only for testing purposes
    // {
    //   type: 'www',
    //   serviceWorker: null,
    // },
  ],
  rollupPlugins: {
    before: [nodePolyfills()],
  },
  extras: {
    enableImportInjection: true,
  },
};
