import { responseTypeMap } from '@multiversx/sdk-web-wallet-cross-window-provider/out/constants/windowProviderConstants';
import { WindowProviderResponseEnums } from '@multiversx/sdk-web-wallet-cross-window-provider/out/enums';
import { WindowProviderRequestEnums } from '@multiversx/sdk-web-wallet-cross-window-provider/out/enums/windowProviderEnums';
import {
  PostMessageParamsType,
  PostMessageReturnType
} from '@multiversx/sdk-web-wallet-cross-window-provider/out/types';
import { WindowManager } from '@multiversx/sdk-web-wallet-cross-window-provider/out/WindowManager';
import {
  IframeLoginTypes,
  iframeWindowReadyEvent,
  safeDocument,
  safeWindow
} from '../constants';
import { IframeProviderEventDataType } from '../IframeProvider';
import {
  ExtendedIframeLoginType,
  IframeProviderContentWindowModel,
  LoginBrandingType
} from './IframeManager.types';

export class IframeManager extends WindowManager {
  private iframeWalletComponent: IframeProviderContentWindowModel | null = null;
  private readonly iframeId = 'mx-iframe-wallet';
  private loginType: ExtendedIframeLoginType = IframeLoginTypes.metamask;
  private hasHandshake: boolean;

  constructor(props?: { onDisconnect?: () => Promise<boolean> }) {
    super();
    this.registerToChildResponse({
      onDisconnect: props?.onDisconnect
    });

    this.hasHandshake = false;
  }

  public get iframeWallet() {
    return this.iframeWalletComponent;
  }

  public override async postMessage<T extends WindowProviderRequestEnums>({
    type,
    payload
  }: PostMessageParamsType<T>): Promise<PostMessageReturnType<T>> {
    this.hasHandshake = await this.handshake(type);

    if (!this.hasHandshake) {
      throw new Error('Cannot establish handshake');
    }

    this.walletWindow?.postMessage(
      {
        type,
        payload
      },
      this.walletUrl
    );

    return await this.listenOnce(responseTypeMap[type]);
  }

  public override async closeConnection(): Promise<boolean> {
    const result = await super.closeConnection();
    this.iframeWalletComponent?.remove();
    this.walletWindow = null;
    return result;
  }

  public async setLoginType(loginType: ExtendedIframeLoginType) {
    this.loginType = loginType;
  }

  public override isWalletOpened(): boolean {
    return Boolean(this.walletWindow) && this.hasHandshake;
  }

  public override closeWalletWindow(): void {
    if (!this.walletWindow) {
      return;
    }
    this.iframeWallet?.setWalletVisible(false);
  }

  public override async setWalletWindow(): Promise<void> {
    if (this.walletWindow) {
      this.iframeWallet?.setWalletVisible(true);
      return;
    }

    const anchor = safeDocument.getElementById?.('root');

    const module = await import('./IframeProviderContentWindow');
    const IframeProviderContentWindow = module.IframeProviderContentWindow;

    this.iframeWalletComponent = new IframeProviderContentWindow({
      id: this.iframeId,
      anchor,
      url: this.walletUrl,
      loginType: this.loginType
    });
    this.iframeWalletComponent.walletAddress = this.walletUrl;

    const iframe = await new Promise(
      (resolve: (value?: HTMLIFrameElement) => void) => {
        this.iframeWalletComponent?.addEventListener(
          iframeWindowReadyEvent,
          (event: Event & { detail?: HTMLIFrameElement }) => {
            resolve(event.detail);
          }
        );
      }
    );

    if (!iframe) {
      throw new Error('Cannot initialize iframe window');
    }

    this.walletWindow = iframe.contentWindow;
    this.setWalletVisible(true);
  }

  public setWalletVisible(visible: boolean): void {
    this.iframeWalletComponent?.setWalletVisible(visible);
  }

  public setLoginBranding(loginBranding: LoginBrandingType): void {
    this.iframeWalletComponent?.setLoginBranding(loginBranding);
  }

  private registerToChildResponse = <
    T extends WindowProviderResponseEnums
  >(props?: {
    onDisconnect?: () => Promise<boolean>;
  }) => {
    safeWindow.addEventListener?.(
      'message',
      async (event: MessageEvent<IframeProviderEventDataType<T>>) => {
        const { data } = event;

        const type = data.type;

        if (event.origin !== this.walletUrl) {
          return;
        }

        if (type === WindowProviderResponseEnums.disconnectResponse) {
          await props?.onDisconnect?.();
          sessionStorage.clear();
          localStorage.clear();
          window.location.reload();
          return;
        }
      }
    );
  };
}
