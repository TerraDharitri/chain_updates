import { Message, Transaction } from '@multiversx/sdk-core';
import {
  CrossWindowProvider,
  IProviderAccount
  // !!! IMPORTANT !!! It is necessary to import explicitly from the file because the module exports (can export) some classes
  // that are using window API and will break the build process on the SSR environments (e.g. PopupConsent)
} from '@multiversx/sdk-web-wallet-cross-window-provider/out/CrossWindowProvider/CrossWindowProvider';
import {
  WindowProviderRequestEnums,
  WindowProviderResponseEnums
} from '@multiversx/sdk-web-wallet-cross-window-provider/out/enums';
import {
  ErrCouldNotLogin,
  ErrCouldNotSignTransactions,
  ErrProviderNotInitialized,
  ErrTransactionCancelled
} from '@multiversx/sdk-web-wallet-cross-window-provider/out/errors';
import { ReplyWithPostMessagePayloadType } from '@multiversx/sdk-web-wallet-cross-window-provider/out/types';
import { IframeLoginTypes } from '../constants';
import { IframeManager } from '../IframeManager/IframeManager';
import {
  ExtendedIframeLoginType,
  LoginBrandingType
} from '../IframeManager/IframeManager.types';

export type IframeProviderEventDataType<T extends WindowProviderResponseEnums> =
  {
    type: T;
    payload: ReplyWithPostMessagePayloadType<T>;
  };

export class IframeProvider extends CrossWindowProvider {
  protected static _instance: IframeProvider | null = null;
  protected readonly windowManager: IframeManager;
  private loginType: ExtendedIframeLoginType = IframeLoginTypes.metamask;

  public constructor() {
    super();
    this.windowManager = new IframeManager({
      onDisconnect: this.logout.bind(this)
    });
  }

  public static getInstance(): IframeProvider {
    if (!IframeProvider._instance) {
      IframeProvider._instance = new IframeProvider();
      return <IframeProvider>IframeProvider._instance;
    }

    return <IframeProvider>IframeProvider._instance;
  }

  public override async init(): Promise<boolean> {
    const initialized = await super.init();
    return initialized;
  }

  public setLoginType(loginType: ExtendedIframeLoginType): void {
    this.loginType = loginType;
    this.windowManager.setLoginType(loginType);
  }

  public setLoginBranding(loginBranding: LoginBrandingType): void {
    this.windowManager.setLoginBranding(loginBranding);
  }

  public override setWalletUrl(url: string): CrossWindowProvider {
    const newUrl = `${url}/?iframeProviderLoginType=${this.loginType}`;
    return super.setWalletUrl(newUrl);
  }

  public override async login(
    options: {
      token?: string;
    } = {}
  ): Promise<IProviderAccount> {
    await this.windowManager.setWalletWindow();

    const account = await super.login(options);

    if (!account.address) {
      this.windowManager.iframeWallet?.remove();
      this.windowManager.walletWindow = null;
      throw new ErrCouldNotLogin();
    }

    this.windowManager.closeWalletWindow();

    return account;
  }

  public override async dispose(): Promise<boolean> {
    return super.dispose();
  }

  public override async logout(): Promise<boolean> {
    if (!this.initialized) {
      throw new ErrProviderNotInitialized();
    }

    try {
      this.ensureConnected();
      await this.windowManager.closeConnection();
    } catch (e) {
      console.error(e);
    }

    this.initialized = false;
    this.disconnect();

    return true;
  }

  public override async signTransaction(
    transaction: Transaction
  ): Promise<Transaction> {
    await this.windowManager.setWalletWindow();
    const data = await super.signTransaction(transaction);
    this.windowManager.closeWalletWindow();
    return data;
  }

  public override async signTransactions(
    transactions: Transaction[]
  ): Promise<Transaction[]> {
    await this.windowManager.setWalletWindow();
    this.ensureConnected();

    const {
      type,
      payload: { data: signedPlainTransactions, error }
    } = await this.windowManager.postMessage({
      type: WindowProviderRequestEnums.signTransactionsRequest,
      payload: transactions.map((tx) => tx.toPlainObject())
    });

    if (error || !signedPlainTransactions) {
      this.windowManager.closeWalletWindow();
      throw new ErrCouldNotSignTransactions();
    }

    if (type === WindowProviderResponseEnums.cancelResponse) {
      this.windowManager.closeWalletWindow();
      throw new ErrTransactionCancelled();
    }

    const hasTransactions = signedPlainTransactions?.length > 0;

    if (!hasTransactions) {
      throw new ErrCouldNotSignTransactions();
    }

    const data = signedPlainTransactions.map((tx) =>
      Transaction.newFromPlainObject(tx)
    );

    this.windowManager.closeWalletWindow();
    return data;
  }

  public override async guardTransactions(
    transactions: Transaction[]
  ): Promise<Transaction[]> {
    await this.windowManager.setWalletWindow();
    const data = await super.guardTransactions(transactions);
    this.windowManager.closeWalletWindow();
    return data;
  }

  public override async signMessage(messageToSign: Message): Promise<Message> {
    await this.windowManager.setWalletWindow();
    messageToSign.signer = this.loginType;
    const data = await super.signMessage(messageToSign);
    this.windowManager.closeWalletWindow();
    return data;
  }

  public override async openPopupConsent(): Promise<boolean> {
    return true;
  }
}
