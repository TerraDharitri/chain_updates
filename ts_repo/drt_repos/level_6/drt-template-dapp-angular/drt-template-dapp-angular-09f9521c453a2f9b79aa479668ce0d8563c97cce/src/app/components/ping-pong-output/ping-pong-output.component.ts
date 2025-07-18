import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TransactionsOutputComponent } from '../transactions-output/transactions-output.component';
import { contractAddress } from '../../../config';

export interface SignedTransactionType {
  hash: string;
  receiver: string;
  value: string;
  gasPrice: string;
  gasLimit: string;
  data?: string;
}

@Component({
  selector: 'app-ping-pong-output',
  standalone: true,
  imports: [CommonModule, TransactionsOutputComponent],
  templateUrl: './ping-pong-output.component.html',
})
export class PingPongOutputComponent {
  @Input() transactions?: SignedTransactionType[] | null = null;

  contractAddress = contractAddress;
}
