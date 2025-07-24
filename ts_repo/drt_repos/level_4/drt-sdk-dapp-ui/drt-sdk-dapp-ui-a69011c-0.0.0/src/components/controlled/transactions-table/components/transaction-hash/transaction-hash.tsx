import { Component, h, Prop } from '@stencil/core';
import classNames from 'classnames';
import { DataTestIdsEnum } from 'constants/dataTestIds.enum';

import type { ITransactionsTableRow } from '../../transactions-table.type';

@Component({
  tag: 'drt-transaction-hash',
  styleUrl: 'transaction-hash.css',
})
export class TransactionHash {
  @Prop() class?: string;
  @Prop() transaction: ITransactionsTableRow;

  render() {
    if (!this.transaction) {
      return null;
    }

    return (
      <div class={classNames(this.class, 'transaction-hash')}>
        <drt-transaction-icon iconInfo={this.transaction.iconInfo}></drt-transaction-icon>
        <drt-explorer-link dataTestId={DataTestIdsEnum.transactionLink} link={this.transaction.link} text={this.transaction.txHash}></drt-explorer-link>
      </div>
    );
  }
}
