import { getWalletWindowMock, WalletWindowMockType } from '../../test-utils';
import { CrossWindowProvider } from '../CrossWindowProvider';
import { Message } from '@terradharitri/sdk-core/out';

describe('CrossWindowProvider Login', () => {
  let crossWindowProvider: CrossWindowProvider;
  let walletWindowMock: WalletWindowMockType;
  let windowOpenSpy: jest.SpyInstance;

  beforeEach(() => {
    walletWindowMock = getWalletWindowMock();
    CrossWindowProvider.getInstance().getWindowManager().postMessage = jest
      .fn()
      .mockImplementation(() => ({
        payload: {
          data: {
            status: 'signed',
            signature: Buffer.from('testSignature').toString('hex')
          }
        }
      }));

    crossWindowProvider = CrossWindowProvider.getInstance();
    windowOpenSpy = jest.spyOn(window, 'open');
    windowOpenSpy.mockImplementation(() => walletWindowMock);
  });

  it('should sign a message correctly', async () => {
    const mockMessage = new Message({ data: Buffer.from('test') });
    crossWindowProvider.setAddress(
      'drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf'
    );
    await crossWindowProvider.init();
    const result = await crossWindowProvider.signMessage(mockMessage);

    expect(result.signature?.toString()).toBe('testSignature');
  });
});
