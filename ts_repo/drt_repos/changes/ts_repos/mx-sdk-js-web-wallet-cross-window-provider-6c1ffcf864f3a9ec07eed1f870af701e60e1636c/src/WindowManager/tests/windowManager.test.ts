import {
  WindowProviderRequestEnums,
  WindowProviderResponseEnums
} from '../../enums';
import { getWalletWindowMock } from '../../test-utils';
import { WindowManager } from '../WindowManager';

describe('WindowManager', () => {
  let windowOpenSpy: jest.SpyInstance;
  let windowCloseSpy: jest.SpyInstance;
  let windowAddListenerSpy: jest.SpyInstance;
  let windowRemoveListenerSpy: jest.SpyInstance;
  let walletWindowMock: { close: jest.Func; postMessage: jest.Func };

  beforeEach(() => {
    walletWindowMock = getWalletWindowMock();

    windowAddListenerSpy = jest.spyOn(window, 'addEventListener');
    windowAddListenerSpy.mockImplementation(jest.fn());

    windowRemoveListenerSpy = jest.spyOn(window, 'removeEventListener');
    windowRemoveListenerSpy.mockImplementation(jest.fn());

    windowCloseSpy = jest.spyOn(window, 'close');
    windowCloseSpy.mockImplementation(jest.fn());

    windowOpenSpy = jest.spyOn(window, 'open');
    windowOpenSpy.mockImplementation(() => walletWindowMock);
  });

  afterEach(() => {
    jest.resetAllMocks();
  });

  it('should set wallet URL and initialize', async () => {
    const windowManager = new WindowManager();
    windowManager.setWalletUrl('https://wallet.example.com');
    await windowManager.init();
    expect(windowManager.isInitialized()).toBeTruthy();
  });

  it('should call init successfully correctly', async () => {
    const windowManager = new WindowManager();
    await windowManager.init();
    expect(windowManager.isInitialized()).toBeTruthy();
  });

  it('should call handshake successfully', async () => {
    const windowManager = new WindowManager();
    windowManager.handshake(WindowProviderRequestEnums.loginRequest);
    expect(windowAddListenerSpy).toHaveBeenCalledTimes(1);
    expect(windowOpenSpy).toHaveBeenCalledTimes(1);
    expect(windowCloseSpy).toHaveBeenCalledTimes(0);
  });

  it('should close connections correctly', async () => {
    const windowManager = new WindowManager();
    windowManager.init();
    windowManager.setWalletUrl('https://wallet.example.com');
    windowManager.closeConnection();

    expect(windowAddListenerSpy).toHaveBeenCalledTimes(1);
    expect(windowOpenSpy).toHaveBeenCalledTimes(1);
    expect(windowCloseSpy).toHaveBeenCalledTimes(0);
  });

  it('should call add event listener successfully', async () => {
    const windowManager = new WindowManager();
    await windowManager.init();
    windowManager.handshake(WindowProviderRequestEnums.loginRequest);
    windowManager.listenOnce(WindowProviderResponseEnums.loginResponse);

    expect(windowAddListenerSpy).toHaveBeenCalledTimes(2);
  });

  it('should call postMessage successfully', async () => {
    const windowManager = new WindowManager();
    windowManager.postMessage({
      type: WindowProviderRequestEnums.loginRequest,
      payload: {
        token: 'token'
      }
    });

    expect(windowAddListenerSpy).toHaveBeenCalledTimes(1);
    expect(windowOpenSpy).toHaveBeenCalledTimes(1);
    expect(windowCloseSpy).toHaveBeenCalledTimes(0);
  });
});
