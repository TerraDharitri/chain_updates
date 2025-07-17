import { Injectable, OnDestroy } from '@angular/core';
import {
  BehaviorSubject,
  Observable,
  timer,
  EMPTY,
  Subject,
  of,
  firstValueFrom,
} from 'rxjs';
import {
  switchMap,
  map,
  catchError,
  takeUntil,
  shareReplay,
  distinctUntilChanged,
} from 'rxjs/operators';
import { HttpClient } from '@angular/common/http';
import { contractAddress } from '../../config';
import { Address, Transaction } from '@terradharitri/sdk-core/out';
import { GAS_PRICE } from '@terradharitri/sdk-dapp/out/constants/mvx.constants';
import BigNumber from 'bignumber.js';
import { getAccount } from '@terradharitri/sdk-dapp/out/methods/account/getAccount';
import { getNetworkConfig } from '@terradharitri/sdk-dapp/out/methods/network/getNetworkConfig';
import { getAccountProvider } from '@terradharitri/sdk-dapp/out/providers/helpers/accountProvider';
import { TransactionManager } from '@terradharitri/sdk-dapp/out/managers/TransactionManager';
import { SignedTransactionType } from '../components/ping-pong-output/ping-pong-output.component';
import { getPendingTransactions } from '@terradharitri/sdk-dapp/out/methods/transactions/getPendingTransactions';

export interface PingPongResponseType {
  data: {
    data: {
      returnData: string[];
    };
  };
}

export interface TransactionsDisplayInfoType {
  processingMessage: string;
  errorMessage: string;
  successMessage: string;
}

const PING_TRANSACTION_INFO: TransactionsDisplayInfoType = {
  processingMessage: 'Processing Ping transaction',
  errorMessage: 'An error has occurred during Ping',
  successMessage: 'Ping transaction successful',
};

const PONG_TRANSACTION_INFO: TransactionsDisplayInfoType = {
  processingMessage: 'Processing Pong transaction',
  errorMessage: 'An error has occurred during Pong',
  successMessage: 'Pong transaction successful',
};

@Injectable({
  providedIn: 'root',
})
export class PingPongService implements OnDestroy {
  private readonly destroy$ = new Subject<void>();

  // State subjects
  private readonly pingAmountSubject = new BehaviorSubject<string>('0');
  private readonly timeToPongSubject = new BehaviorSubject<number | null>(null);
  private readonly loadingSubject = new BehaviorSubject<boolean>(false);
  private readonly errorSubject = new BehaviorSubject<string | null>(null);

  // Public observables
  public readonly pingAmount$ = this.pingAmountSubject.asObservable();
  public readonly timeToPong$ = this.timeToPongSubject.asObservable();
  public readonly loading$ = this.loadingSubject.asObservable();
  public readonly error$ = this.errorSubject.asObservable();

  public readonly timeToPongRefresh$ = timer(0, 1000).pipe(
    takeUntil(this.destroy$),
    switchMap(() => this.fetchTimeToPong()),
    distinctUntilChanged(),
    shareReplay(1)
  );

  constructor(private http: HttpClient) {
    this.initializeAutoRefresh();
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  // Public getters for immediate values
  get currentPingAmount(): string {
    return this.pingAmountSubject.value;
  }

  get currentTimeToPong(): number | null {
    return this.timeToPongSubject.value;
  }

  get isLoading(): boolean {
    return this.loadingSubject.value;
  }

  // Private methods
  private initializeAutoRefresh(): void {
    this.timeToPongRefresh$.subscribe({
      next: time => this.timeToPongSubject.next(time),
      error: error => this.handleError('Failed to fetch time to pong', error),
    });
  }

  private async fetchPingAmount(): Promise<string> {
    try {
      const data = await firstValueFrom(this.makeVmQuery('getPingAmount', []));

      if (!data) {
        return '0';
      }

      return this.decodeAmount(data);
    } catch (error) {
      console.error('Unable to call getPingAmount', error);
      return '0';
    }
  }

  private fetchTimeToPong(): Observable<number | null> {
    const account = getAccount();

    try {
      const args = new Address(account.address).toHex();

      return this.makeVmQuery('getTimeToPong', [args]).pipe(
        map(data => this.decodeTime(data)),
        catchError(error => {
          console.error('Unable to call getTimeToPong', error);
          return of(0);
        })
      );
    } catch (error) {
      console.error('Invalid address format', error);
      return of(0);
    }
  }

  private makeVmQuery(
    funcName: string,
    args: string[]
  ): Observable<PingPongResponseType> {
    const networkConfig = getNetworkConfig();

    const body = {
      scAddress: contractAddress,
      funcName,
      args,
    };

    return this.http.post<PingPongResponseType>(
      `${networkConfig.network.apiAddress}/vm-values/query`,
      body
    );
  }

  private decodeAmount(data: PingPongResponseType): string {
    if (!data.data.data.returnData) {
      return '0';
    }

    const returnValue = data.data.data.returnData[0];
    if (!returnValue) return '0';

    const decodedString = Buffer.from(returnValue, 'base64').toString('hex');
    return new BigNumber(decodedString, 16).toString(10);
  }

  private decodeTime(data: PingPongResponseType): number | null {
    if (!data.data.data.returnData) {
      return null;
    }

    const returnValue = data.data.data.returnData[0];
    if (returnValue === '' || !returnValue) {
      return 0;
    }

    const decodedString = Buffer.from(returnValue, 'base64').toString('hex');
    return new BigNumber(decodedString, 16).toNumber();
  }

  private handleError(message: string, error: any): void {
    console.error(message, error);
    this.errorSubject.next(message);
    this.loadingSubject.next(false);
  }

  public refreshTimeToPong(): Observable<number | null> {
    return this.fetchTimeToPong();
  }

  // Transaction methods (simplified - in real app would use proper DharitrI services)
  public async sendPingTransaction(): Promise<string> {
    this.loadingSubject.next(true);
    this.errorSubject.next(null);

    try {
      const pingAmount = await this.fetchPingAmount();
      const account = getAccount();
      const networkConfig = getNetworkConfig();

      const pingTransaction = new Transaction({
        value: BigInt(pingAmount),
        data: Buffer.from('ping'),
        receiver: new Address(contractAddress),
        gasLimit: BigInt(6000000),
        gasPrice: BigInt(GAS_PRICE),
        chainID: networkConfig.network.chainId,
        sender: new Address(account.address),
        version: 1,
      });

      const sessionId = await this.signAndSendTransactions(
        [pingTransaction],
        PING_TRANSACTION_INFO
      );

      this.loadingSubject.next(false);
      return sessionId;
    } catch (error) {
      this.handleError('Failed to send ping transaction', error);
      throw error;
    }
  }

  public async sendPongTransaction(): Promise<string> {
    this.loadingSubject.next(true);
    this.errorSubject.next(null);

    try {
      const account = getAccount();
      const networkConfig = getNetworkConfig();

      const pongTransaction = new Transaction({
        value: BigInt(0),
        data: Buffer.from('pong'),
        receiver: new Address(contractAddress),
        gasLimit: BigInt(6000000),
        gasPrice: BigInt(GAS_PRICE),
        chainID: networkConfig.network.chainId,
        sender: new Address(account.address),
        version: 1,
      });

      const sessionId = await this.signAndSendTransactions(
        [pongTransaction],
        PONG_TRANSACTION_INFO
      );

      this.loadingSubject.next(false);
      return sessionId;
    } catch (error) {
      this.handleError('Failed to send pong transaction', error);
      throw error;
    }
  }

  public async signAndSendTransactions(
    transactions: Transaction[],
    transactionsDisplayInfo?: TransactionsDisplayInfoType
  ): Promise<string> {
    const provider = getAccountProvider();
    const txManager = TransactionManager.getInstance();

    const signedTransactions = await provider.signTransactions(transactions);
    const sentTransactions = await txManager.send(signedTransactions);
    const sessionId = await txManager.track(sentTransactions, {
      transactionsDisplayInfo,
    });

    return sessionId;
  }

  public getPendingTransactions(): SignedTransactionType[] {
    const pendingTransactions = getPendingTransactions();

    return pendingTransactions.map(tx => ({
      hash: tx.hash,
      receiver: tx.receiver,
      value: tx.value?.toString() || '0',
      gasPrice: tx.gasPrice?.toString() || '0',
      gasLimit: tx.gasLimit?.toString() || '0',
    }));
  }

  public clearError(): void {
    this.errorSubject.next(null);
  }
}
