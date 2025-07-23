import { Component, h, Prop } from '@stencil/core';
import { DataTestIdsEnum } from 'constants/dataTestIds.enum';

import type { TransactionRowType } from '../../transactions-table.type';

const transactionHashClasses: Record<string, string> = {
  explorerLink: 'drt:text-blue-link!',
};

@Component({
  tag: 'drt-transaction-hash',
})
export class TransactionHash {
  @Prop() class?: string;
  @Prop() transaction: TransactionRowType;

  render() {
    if (!this.transaction) {
      return null;
    }

    return (
      <div class={{ 'transaction-hash': true, [this.class]: Boolean(this.class) }}>
        <drt-transaction-icon iconInfo={this.transaction.iconInfo} />

        <drt-explorer-link
          dataTestId={DataTestIdsEnum.transactionLink}
          link={this.transaction.link}
          class={transactionHashClasses.explorerLink}
        >
          <span>{this.transaction.txHash}</span>
        </drt-explorer-link>
      </div>
    );
  }
}
