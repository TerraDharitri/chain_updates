import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { LabelComponent } from '../label/label.component';
import { FormatAmountComponent } from '../format-amount/format-amount.component';
import { MxLinkComponent } from '../mx-link/mx-link.component';
import { SignedTransactionType } from '../ping-pong-output/ping-pong-output.component';

@Component({
  selector: 'app-transaction-output',
  standalone: true,
  imports: [
    CommonModule,
    LabelComponent,
    FormatAmountComponent,
    MxLinkComponent,
  ],
  templateUrl: './transaction-output.component.html',
})
export class TransactionOutputComponent {
  @Input() transaction?: SignedTransactionType;

  get decodedData(): string {
    if (!this.transaction?.data) {
      return 'N/A';
    }

    try {
      return Buffer.from(this.transaction.data, 'base64').toString('ascii');
    } catch {
      return 'N/A';
    }
  }
}
