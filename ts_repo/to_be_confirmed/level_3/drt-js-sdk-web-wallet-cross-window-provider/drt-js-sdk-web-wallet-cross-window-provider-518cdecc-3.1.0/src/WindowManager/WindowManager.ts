import { safeWindow, responseTypeMap } from '../constants';
import {
  WindowProviderRequestEnums,
  WindowProviderResponseEnums
} from '../enums';
import {
  ErrCannotEstablishHandshake,
  ErrProviderNotInitialized,
  ErrProviderNotInstantiated
} from '../errors';
import {
  PostMessageParamsType,
  PostMessageReturnType,
  ReplyWithPostMessageEventType,
  ReplyWithPostMessagePayloadType
} from '../types';

export class WindowManager {
  private _session = '';
  private _walletUrl = '';
  protected initialized = false;
  public walletWindow: Window | null = null;
  private activeListeners: Map<string, (event: MessageEvent) => void> =
    new Map();

  constructor() {
    safeWindow.addEventListener?.('beforeunload', () => {
      this.closeWalletWindow();
    });

    safeWindow.name = safeWindow.location?.origin;
  }

  public get walletUrl(): string {
    return this._walletUrl;
  }

  public setWalletUrl(url: string): WindowManager {
    this._walletUrl = url;
    return this;
  }

  async init(): Promise<boolean> {
    this.initialized = typeof window !== 'undefined';
    return this.initialized;
  }

  onDestroy(): boolean {
    this.initialized = false;
    return this.initialized;
  }

  public isWalletOpened(type?: WindowProviderRequestEnums) {
    return (
      type === WindowProviderRequestEnums.cancelAction &&
      Boolean(this.walletWindow)
    );
  }

  async handshake(type: WindowProviderRequestEnums): Promise<boolean> {
    const isOpened = this.isWalletOpened(type);

    if (isOpened) {
      return true;
    }

    this.closeWalletWindow();
    await this.setWalletWindow();

    const { payload } = await this.listenOnce(
      WindowProviderResponseEnums.handshakeResponse
    );

    if (!payload) {
      throw new ErrCannotEstablishHandshake();
    }

    this._session = this._session || payload.data || Date.now().toString();
    this.walletWindow?.postMessage(
      {
        type: WindowProviderRequestEnums.finalizeHandshakeRequest,
        payload: this._session
      },
      this.walletUrl
    );

    this.addHandshakeChangeListener();
    return true;
  }

  private async addHandshakeChangeListener() {
    const walletUrl = this.walletUrl;

    const eventHandler = (
      event: MessageEvent<ReplyWithPostMessageEventType>
    ) => {
      try {
        const { type, payload } = event.data;
        const isWalletEvent = event.origin === new URL(walletUrl).origin;

        if (!isWalletEvent) {
          return;
        }

        switch (type) {
          case WindowProviderResponseEnums.handshakeResponse:
            if (payload === '') {
              this.walletWindow?.close();
              this.walletWindow = null;
              safeWindow.removeEventListener?.('message', eventHandler);
              break;
            }

            // Save the current session and send it in later handshake requests
            if (typeof payload === 'string') {
              this._session = payload;
            }

            break;
        }
      } catch (e) {
        console.error('Handshake response error', e);
      }
    };

    safeWindow.addEventListener?.('message', eventHandler);
  }

  async listenOnce<T extends WindowProviderResponseEnums>(
    action: T
  ): Promise<{
    type: T;
    payload: ReplyWithPostMessagePayloadType<T>;
  }> {
    if (!this.walletWindow) {
      throw new ErrProviderNotInstantiated();
    }

    return new Promise((resolve) => {
      const walletUrl = this.walletUrl;

      // Prevent duplicate listeners for the same action
      if (this.activeListeners.has(action)) {
        const existingHandler = this.activeListeners.get(action);
        if (existingHandler) {
          safeWindow.removeEventListener('message', existingHandler);
        }
      }

      const eventHandler = (
        event: MessageEvent<{
          type: T;
          payload: ReplyWithPostMessagePayloadType<T>;
        }>
      ) => {
        const { type, payload } = event.data;
        const isWalletEvent = event.origin === new URL(walletUrl).origin;

        const isCurrentAction =
          action === type ||
          type === WindowProviderResponseEnums.cancelResponse;

        if (!isCurrentAction || !isWalletEvent) {
          return;
        }

        safeWindow.removeEventListener('message', eventHandler);
        this.activeListeners.delete(action);

        resolve({ type, payload });
      };

      safeWindow.addEventListener('message', eventHandler);
      this.activeListeners.set(action, eventHandler);
    });
  }

  async closeConnection(): Promise<boolean> {
    if (!this.initialized) {
      throw new ErrProviderNotInitialized();
    }

    await this.postMessage({
      type: WindowProviderRequestEnums.logoutRequest,
      payload: undefined
    });

    // Reset the session on logout
    this._session = Date.now().toString();
    this.initialized = false;

    return true;
  }

  isInitialized(): boolean {
    return this.initialized;
  }

  async postMessage<T extends WindowProviderRequestEnums>({
    type,
    payload
  }: PostMessageParamsType<T>): Promise<PostMessageReturnType<T>> {
    await this.handshake(type);

    this.walletWindow?.postMessage(
      {
        type,
        payload
      },
      this.walletUrl
    );

    const data = await this.listenOnce(responseTypeMap[type]);
    this.closeWalletWindow();

    return data;
  }

  public closeWalletWindow() {
    this.walletWindow?.close();
  }

  public async setWalletWindow(): Promise<void> {
    this.walletWindow =
      safeWindow.open?.(this.walletUrl, this.walletUrl) ?? null;
  }
}
