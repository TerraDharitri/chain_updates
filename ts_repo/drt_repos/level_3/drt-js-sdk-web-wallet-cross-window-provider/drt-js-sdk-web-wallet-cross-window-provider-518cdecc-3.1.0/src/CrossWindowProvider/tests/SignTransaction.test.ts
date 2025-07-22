import { Transaction } from '@terradharitri/sdk-core';
import {
  createMockTransaction,
  getWalletWindowMock,
  WalletWindowMockType
} from '../../test-utils';
import { CrossWindowProvider } from '../CrossWindowProvider';

describe('CrossWindowProvider Login', () => {
  let crossWindowProvider: CrossWindowProvider;
  let walletWindowMock: WalletWindowMockType;
  let windowOpenSpy: jest.SpyInstance;

  let mockTransaction: Transaction;

  beforeEach(() => {
    walletWindowMock = getWalletWindowMock();
    mockTransaction = createMockTransaction({
      data: 'data',
      receiverUsername: 'receiver',
      senderUsername: 'sender'
    });

    CrossWindowProvider.getInstance().getWindowManager().postMessage = jest
      .fn()
      .mockImplementation(() => ({
        payload: { data: [mockTransaction.toPlainObject()] }
      }));

    crossWindowProvider = CrossWindowProvider.getInstance();
    windowOpenSpy = jest.spyOn(window, 'open');
    windowOpenSpy.mockImplementation(() => walletWindowMock);
  });

  it('should sign a transaction correctly', async () => {
    crossWindowProvider.setAddress('testAddress');
    await crossWindowProvider.init();

    const result = await crossWindowProvider.signTransaction(mockTransaction);
    expect(result).toStrictEqual(mockTransaction);
  });
});
