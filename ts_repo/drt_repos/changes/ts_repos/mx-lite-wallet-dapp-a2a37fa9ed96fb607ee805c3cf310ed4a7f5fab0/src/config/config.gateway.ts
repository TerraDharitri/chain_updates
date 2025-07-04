import { EnvironmentsEnum } from 'types';
import { sharedNetorks } from './sharedNetworks';
import { NetworkType } from '../redux/slices';

export * from './sharedConfig';

export const networks: NetworkType[] = [
  ...sharedNetorks,
  {
    default: true,
    id: EnvironmentsEnum.devnet,
    name: 'Gateway',
    apiAddress: '',
    gatewayUrl: 'https://devnet-gateway.dharitri.org',
    extrasApi: 'https://devnet-extras-api.dharitri.org',
    sampleAuthenticatedDomains: [''],
    sovereignContractAddress: '',
    walletAddress: 'https://devnet-wallet.dharitri.org',
    WEGLDid: ''
  }
];
