import { getWalletWindowMock, WalletWindowMockType } from '../../test-utils';
import { CrossWindowProvider } from '../CrossWindowProvider';

describe('CrossWindowProvider', () => {
  let crossWindowProvider: CrossWindowProvider;
  let walletWindowMock: WalletWindowMockType;
  let windowOpenSpy: jest.SpyInstance;

  beforeEach(() => {
    walletWindowMock = getWalletWindowMock();
    CrossWindowProvider.getInstance().getWindowManager().postMessage = jest
      .fn()
      .mockImplementation(() => true);

    CrossWindowProvider.getInstance().getWindowManager().closeConnection = jest
      .fn()
      .mockImplementation(() => true);

    crossWindowProvider = CrossWindowProvider.getInstance();
    windowOpenSpy = jest.spyOn(window, 'open');
    windowOpenSpy.mockImplementation(() => walletWindowMock);
  });

  it('should handle logout correctly', async () => {
    await crossWindowProvider.init();
    crossWindowProvider.setAddress('testAddress');
    const result = await crossWindowProvider.logout();
    expect(result).toBe(true);
  });
});
