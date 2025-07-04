import { Component, OnInit, OnDestroy, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { faArrowUp, faArrowDown } from '@fortawesome/free-solid-svg-icons';
import { Subject, combineLatest } from 'rxjs';
import { takeUntil, map } from 'rxjs/operators';

import { OutputContainerComponent } from '../../components/output-container/output-container.component';
import { ButtonComponent } from '../../components/button/button.component';
import {
  PingPongOutputComponent,
  SignedTransactionType,
} from '../../components/ping-pong-output/ping-pong-output.component';
import { contractAddress } from '../../../config';
import { PingPongService } from '../../services/ping-pong.service';
import {
  calculatePingPongState,
  formatTimeRemaining,
} from '../../helpers/countdown.helpers';
import BigNumber from 'bignumber.js';

interface ViewState {
  timeToPong: number | null;
  timeRemaining: string;
  canPing: boolean;
  canPong: boolean;
  isLoading: boolean;
  error: string | null;
}

@Component({
  selector: 'app-ping-pong-raw',
  standalone: true,
  imports: [
    CommonModule,
    FontAwesomeModule,
    OutputContainerComponent,
    ButtonComponent,
    PingPongOutputComponent,
  ],
  templateUrl: './ping-pong-raw.component.html',
  styleUrls: ['./ping-pong-raw.component.css'],
})
export class PingPongRawComponent implements OnInit, OnDestroy {
  // Icons
  readonly faArrowUp = faArrowUp;
  readonly faArrowDown = faArrowDown;
  readonly contractAddress = contractAddress;

  // State
  viewState: ViewState = {
    timeToPong: null,
    timeRemaining: '00:00',
    canPing: true,
    canPong: false,
    isLoading: false,
    error: null,
  };

  public pendingTransactions: SignedTransactionType[] = [];

  // Lifecycle management
  private readonly destroy$ = new Subject<void>();

  // Dependency injection (modern Angular pattern)
  private readonly pingPongService = inject(PingPongService);

  ngOnInit(): void {
    this.initializeComponent();
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  private initializeComponent(): void {
    combineLatest([
      this.pingPongService.timeToPong$,
      this.pingPongService.loading$,
      this.pingPongService.error$,
    ])
      .pipe(
        map(([timeToPong, isLoading, error]) => {
          const pingPongState = calculatePingPongState(timeToPong);
          const timeRemaining = formatTimeRemaining(timeToPong || 0);

          return {
            timeToPong,
            timeRemaining,
            canPing: pingPongState.canPing,
            canPong: pingPongState.canPong,
            isLoading,
            error,
          };
        }),
        takeUntil(this.destroy$)
      )
      .subscribe(viewState => {
        this.viewState = viewState;
      });
  }

  formatAmount(amount: string): string {
    // Convert from wei to EGLD (simple formatting)
    const amountBN = new BigNumber(amount);
    const egldAmount = amountBN.dividedBy(new BigNumber(10).pow(18));
    return egldAmount.toFixed(4) + ' EGLD';
  }

  async onSendPingTransaction(): Promise<void> {
    try {
      await this.pingPongService.sendPingTransaction();
      this.pendingTransactions = this.pingPongService.getPendingTransactions();
    } catch (error) {
      console.error('Error sending ping transaction:', error);
    }
  }

  async onSendPongTransaction(): Promise<void> {
    try {
      await this.pingPongService.sendPongTransaction();
      this.pendingTransactions = this.pingPongService.getPendingTransactions();
    } catch (error) {
      console.error('Error sending pong transaction:', error);
    }
  }

  clearError(): void {
    this.pingPongService.clearError();
  }
}
