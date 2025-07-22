export interface WalletWindowMockType {
  close: jest.Func;
  postMessage: jest.Func;
}

export const getWalletWindowMock = (): WalletWindowMockType => ({
  close: jest.fn(),
  postMessage: jest.fn()
});
