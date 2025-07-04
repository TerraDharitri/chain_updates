import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TransactionOutputComponent } from '../transaction-output/transaction-output.component';
import { SignedTransactionType } from '../ping-pong-output/ping-pong-output.component';

@Component({
  selector: 'app-transactions-output',
  standalone: true,
  imports: [CommonModule, TransactionOutputComponent],
  templateUrl: './transactions-output.component.html',
})
export class TransactionsOutputComponent {
  @Input() transactions: SignedTransactionType[] = [];

  trackByHash(index: number, transaction: SignedTransactionType): string {
    return transaction.hash;
  }
}
