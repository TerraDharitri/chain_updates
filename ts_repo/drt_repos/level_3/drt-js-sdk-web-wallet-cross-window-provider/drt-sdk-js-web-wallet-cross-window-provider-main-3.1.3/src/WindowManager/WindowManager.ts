import { safeWindow, responseTypeMap } from '../constants';
import { confirmationDialogTag } from '../CrossWindowProvider/PopupConsent/constants';
import { PopupConsentModel } from '../CrossWindowProvider/PopupConsent/PopupConsent.model';
import {
  WindowProviderRequestEnums,
  WindowProviderResponseEnums
} from '../enums';
import {
  ErrCannotEstablishHandshake,
  ErrCouldNotLogin,
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
  public isPopupBlocked = false;
  private _isLoggingOut = false;

  private activeListeners: Map<string, (event: MessageEvent) => void> =
    new Map();
  private static readonly SESSION_STORAGE_KEY = 'mx-wallet-session-id';

  constructor() {
    safeWindow.addEventListener?.('beforeunload', () => {
      this.closeWalletWindow();
    });

    safeWindow.name = safeWindow.location?.origin;
    this._loadSessionFromStorage();
  }

  public get walletUrl(): string {
    return this._walletUrl;
  }

  public setWalletUrl(url: string): WindowManager {
    this._walletUrl = url;
    return this;
  }

  private _loadSessionFromStorage(): void {
    try {
      const storedSession = safeWindow.sessionStorage?.getItem(
        WindowManager.SESSION_STORAGE_KEY
      );
      if (storedSession) {
        this._session = storedSession;
      }
    } catch (error) {
      console.warn('Failed to load session from sessionStorage:', error);
    }
  }

  private _saveSessionToStorage(sessionId: string): void {
    try {
      safeWindow.sessionStorage?.setItem(
        WindowManager.SESSION_STORAGE_KEY,
        sessionId
      );
    } catch (error) {
      console.warn('Failed to save session to sessionStorage:', error);
    }
  }

  private _clearSessionFromStorage(): void {
    try {
      safeWindow.sessionStorage?.removeItem(WindowManager.SESSION_STORAGE_KEY);
    } catch (error) {
      console.warn('Failed to clear session from sessionStorage:', error);
    }
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

  public async openPopupConsent(): Promise<boolean> {
    if (
      typeof document === 'undefined' ||
      typeof window === 'undefined' ||
      !this.isPopupBlocked
    ) {
      return true;
    }

    const module = await import(
      '../CrossWindowProvider/PopupConsent/PopupConsent'
    );
    const PopupConsent = module.PopupConsent;

    const customElements = safeWindow.customElements;
    if (customElements && !customElements.get(confirmationDialogTag)) {
      customElements.define(confirmationDialogTag, PopupConsent);
    }

    const popup = document.createElement(
      confirmationDialogTag
    ) as PopupConsentModel & HTMLElement;

    popup.walletUrl = this.walletUrl;

    document.body.appendChild(popup);

    const popupConsentResponse: boolean = await new Promise<boolean>(
      (resolve) => {
        popup.onConfirm = () => {
          resolve(true);
          document.body.removeChild(popup);
        };
        popup.onCancel = () => {
          resolve(false);
          document.body.removeChild(popup);
        };
      }
    );

    return popupConsentResponse;
  }

  async handshake(type: WindowProviderRequestEnums): Promise<boolean> {
    const isOpened = this.isWalletOpened(type);

    if (isOpened) {
      return true;
    }

    this.closeWalletWindow();
    const isWindowInitialized = await this.setWalletWindow();

    if (!isWindowInitialized) {
      const popupConsentResponse = await this.openPopupConsent();

      if (!popupConsentResponse) {
        throw new ErrCouldNotLogin();
      }
    }

    const { payload } = await this.listenOnce(
      WindowProviderResponseEnums.handshakeResponse
    );

    if (!payload) {
      throw new ErrCannotEstablishHandshake();
    }

    const receivedSession =
      typeof payload === 'string' ? payload : payload.data;

    this._session = this._session || receivedSession || Date.now().toString();
    // Save the current session in sessionStorage to preserve it across page reloads
    this._saveSessionToStorage(this._session);

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
              this._saveSessionToStorage(this._session);
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
    // If the wallet window was previously blocked (e.g., by a popup blocker),
    // reinitialize the connection by opening the wallet window again.
    if (!this.walletWindow) {
      const isWindowInitialized = await this.setWalletWindow();

      if (!isWindowInitialized) {
        throw new ErrProviderNotInstantiated();
      }
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
    if (this._isLoggingOut) {
      console.warn('Logout already in progress');
      return false;
    }

    if (!this.initialized) {
      throw new ErrProviderNotInitialized();
    }

    this._isLoggingOut = true;

    try {
      await this.postMessage({
        type: WindowProviderRequestEnums.logoutRequest,
        payload: undefined
      });

      this._session = Date.now().toString();
      this._clearSessionFromStorage();
      this.initialized = false;

      return true;
    } finally {
      this._isLoggingOut = false;
    }
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

  public async setWalletWindow(): Promise<boolean> {
    this.walletWindow =
      safeWindow.open?.(this.walletUrl, this.walletUrl) ?? null;

    const isWindowInitialized = Boolean(this.walletWindow);
    this.isPopupBlocked = !isWindowInitialized;

    return isWindowInitialized;
  }
}
