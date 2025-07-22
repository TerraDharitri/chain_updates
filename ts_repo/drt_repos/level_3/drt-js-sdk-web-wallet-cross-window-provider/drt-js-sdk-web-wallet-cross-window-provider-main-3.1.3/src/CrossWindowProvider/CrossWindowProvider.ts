import { Address, Message, Transaction } from '@terradharitri/sdk-core';
import {
  WindowProviderRequestEnums,
  WindowProviderResponseEnums,
  SignMessageStatusEnum
} from '../enums';
import {
  ErrAccountNotConnected,
  ErrCannotSignSingleTransaction,
  ErrCouldNotGuardTransactions,
  ErrCouldNotLogin,
  ErrCouldNotSignMessage,
  ErrCouldNotSignTransactions,
  ErrProviderNotInitialized,
  ErrTransactionCancelled
} from '../errors';
import { WindowManager } from '../WindowManager';

export interface IProviderAccount {
  address: string;
  signature?: string;
  multisig?: string;
  impersonate?: string;
}

export class CrossWindowProvider {
  private account: IProviderAccount = { address: '' };
  private accessToken: string | undefined = undefined;

  protected initialized = false;
  protected windowManager: WindowManager;
  protected static _instance: CrossWindowProvider | null = null;

  public constructor() {
    this.windowManager = new WindowManager();
  }

  protected ensureConnected() {
    if (!this.account.address) {
      throw new ErrAccountNotConnected();
    }
  }

  protected disconnect() {
    this.account = { address: '' };
  }

  public static getInstance(): CrossWindowProvider {
    if (!CrossWindowProvider._instance) {
      CrossWindowProvider._instance = new CrossWindowProvider();
      return CrossWindowProvider._instance;
    }

    return CrossWindowProvider._instance;
  }

  public getWindowManager(): WindowManager {
    return this.windowManager;
  }

  public setAddress(address: string): CrossWindowProvider {
    this.account.address = address;
    return this;
  }

  public setWalletUrl(url: string): CrossWindowProvider {
    this.windowManager.setWalletUrl(url);
    return this;
  }

  public setWalletWindow(): Promise<boolean> {
    return this.windowManager.setWalletWindow();
  }

  async init(): Promise<boolean> {
    this.initialized = await this.windowManager.init();
    return this.initialized;
  }

  onDestroy(): boolean {
    CrossWindowProvider._instance = null;
    this.initialized = this.windowManager.onDestroy();
    return this.initialized;
  }

  async login(
    options: {
      token?: string;
    } = {}
  ): Promise<IProviderAccount> {
    if (!this.initialized) {
      throw new ErrProviderNotInitialized();
    }

    const isRelogin = this.isConnected();

    if (isRelogin) {
      const { address, signature, multisig, impersonate } = this.account;
      return {
        address,
        signature,
        multisig,
        impersonate
      };
    }

    this.accessToken = options.token;

    const {
      payload: { data, error }
    } = await this.windowManager.postMessage({
      type: WindowProviderRequestEnums.loginRequest,
      payload: {
        token: this.accessToken
      }
    });

    if (error || !data) {
      throw new ErrCouldNotLogin();
    }

    this.account.address = data.address;
    this.account.signature = data.signature;
    this.account.multisig = data.multisig;
    this.account.impersonate = data.impersonate;

    return {
      address: this.account.address,
      signature: this.account.signature,
      multisig: this.account.multisig,
      impersonate: this.account.impersonate
    };
  }

  async dispose(): Promise<boolean> {
    const connectionClosed = await this.windowManager.closeConnection();
    this.initialized = !connectionClosed;
    CrossWindowProvider._instance = null;
    return connectionClosed;
  }

  async logout(): Promise<boolean> {
    if (!this.initialized) {
      throw new ErrProviderNotInitialized();
    }

    this.ensureConnected();
    const connectionClosed = await this.dispose();
    this.disconnect();

    return connectionClosed;
  }

  async getAddress(): Promise<string> {
    if (!this.initialized) {
      throw new ErrProviderNotInitialized();
    }

    return this.account?.address ?? '';
  }

  isInitialized(): boolean {
    return this.initialized;
  }

  isConnected(): boolean {
    return Boolean(this.account.address);
  }

  getAccount(): IProviderAccount | null {
    return this.account;
  }

  setAccount(account: IProviderAccount): void {
    this.account = account;
  }

  async signTransaction(transaction: Transaction): Promise<Transaction> {
    this.ensureConnected();

    const signedTransactions = await this.signTransactions([transaction]);

    if (signedTransactions.length != 1) {
      throw new ErrCannotSignSingleTransaction();
    }

    return signedTransactions[0];
  }

  async signTransactions(transactions: Transaction[]): Promise<Transaction[]> {
    this.ensureConnected();

    const {
      type,
      payload: { data: signedPlainTransactions, error }
    } = await this.windowManager.postMessage({
      type: WindowProviderRequestEnums.signTransactionsRequest,
      payload: transactions.map((tx) => tx.toPlainObject())
    });

    if (error || !signedPlainTransactions) {
      throw new ErrCouldNotSignTransactions();
    }

    if (type === WindowProviderResponseEnums.cancelResponse) {
      throw new ErrTransactionCancelled();
    }

    const hasTransactions = signedPlainTransactions?.length > 0;

    if (!hasTransactions) {
      throw new ErrCouldNotSignTransactions();
    }

    return signedPlainTransactions.map((tx) => Transaction.fromPlainObject(tx));
  }

  async guardTransactions(transactions: Transaction[]): Promise<Transaction[]> {
    this.ensureConnected();

    const {
      type,
      payload: { data: signedPlainTransactions, error }
    } = await this.windowManager.postMessage({
      type: WindowProviderRequestEnums.guardTransactionsRequest,
      payload: transactions.map((tx) => tx.toPlainObject())
    });

    if (error || !signedPlainTransactions) {
      throw new ErrCouldNotSignTransactions();
    }

    if (type === WindowProviderResponseEnums.cancelResponse) {
      throw new ErrTransactionCancelled();
    }

    const hasTransactions = signedPlainTransactions?.length > 0;

    if (!hasTransactions) {
      throw new ErrCouldNotGuardTransactions();
    }

    return signedPlainTransactions.map((tx) => Transaction.fromPlainObject(tx));
  }

  async signMessage(messageToSign: Message): Promise<Message> {
    this.ensureConnected();

    const {
      payload: { data, error }
    } = await this.windowManager.postMessage({
      type: WindowProviderRequestEnums.signMessageRequest,
      payload: {
        message: Buffer.from(messageToSign.data).toString()
      }
    });

    if (error || !data) {
      throw new ErrCouldNotSignMessage();
    }

    const { status, signature } = data;

    if (status !== SignMessageStatusEnum.signed) {
      throw new ErrCouldNotSignMessage();
    }

    return new Message({
      data: Buffer.from(messageToSign.data),
      address:
        messageToSign.address ?? Address.fromBech32(this.account.address),
      signer: messageToSign.signer || 'wallet-cross-window',
      version: messageToSign.version,
      signature: Buffer.from(String(signature), 'hex')
    });
  }

  cancelAction() {
    const isWalletOpened = this.windowManager?.isWalletOpened(
      WindowProviderRequestEnums.cancelAction
    );

    if (!isWalletOpened) {
      return;
    }

    return this.windowManager?.postMessage({
      type: WindowProviderRequestEnums.cancelAction,
      payload: undefined
    });
  }
}
