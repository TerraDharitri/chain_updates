import { NetworkType } from 'redux/slices';
import { sharedNetorks } from './sharedNetworks';

export * from './sharedConfig';

// This config is used for puppeteer tests

const sovereignNetwork = sharedNetorks.find(
  (network) => network.id === 'sovereign'
);

if (!sovereignNetwork) {
  throw new Error('Sovereign network not found');
}

export const networks: NetworkType[] = [
  ...sharedNetorks.filter((network) => network.id !== 'sovereign'),
  {
    ...sovereignNetwork,
    default: true,
    apiAddress: 'https://api-sovereign-test.numbat.ro',
    gatewayUrl: '',
    extrasApi: 'https://extras-api-sovereign-test.numbat.ro',
    sampleAuthenticatedDomains: ['https://api-sovereign-test.numbat.ro'],
    sovereignContractAddress: '',
    walletAddress: 'https://wallet.voyager1.dev'
  }
];
