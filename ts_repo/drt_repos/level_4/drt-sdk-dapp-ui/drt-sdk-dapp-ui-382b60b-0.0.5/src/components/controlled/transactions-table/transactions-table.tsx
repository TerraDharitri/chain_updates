import { Component, h, Prop } from '@stencil/core';
import classNames from 'classnames';

import { DataTestIdsEnum } from '../../../constants/dataTestIds.enum';
import type { ITransactionsTableRow } from './transactions-table.type';

const COLUMNS = ['TxHash', 'Age', 'Shard', 'From', 'To', 'Method', 'Value'];

@Component({
  tag: 'drt-transactions-table',
  styleUrl: 'transactions-table.css',
})
export class TransactionsTable {
  @Prop() class?: string;
  @Prop() transactions: ITransactionsTableRow[];

  render() {
    return (
      <table class={classNames(this.class, 'transactions-table')}>
        <thead class="transactions-table-header">
          <tr>
            {COLUMNS.map(column => (
              <th key={column} scope="col" class="transactions-table-header-cell">
                {column}
              </th>
            ))}
          </tr>
        </thead>
        <tbody class="transactions-table-body">
          {this.transactions.map(transaction => (
            <tr class="transactions-table-body-row">
              <td class="transactions-table-body-cell">
                <drt-transaction-hash transaction={transaction}></drt-transaction-hash>
              </td>
              <td class="transactions-table-body-cell">
                <drt-transaction-age age={transaction.age.timeAgo} tooltip={transaction.age.tooltip}></drt-transaction-age>
              </td>
              <td class="transactions-table-body-cell">
                <drt-transaction-shards transaction={transaction}></drt-transaction-shards>
              </td>
              <td class="transactions-table-body-cell">
                <drt-transaction-account
                  account={transaction.sender}
                  dataTestId={DataTestIdsEnum.transactionSender}
                  scope="sender"
                  showLockedAccounts={true}
                ></drt-transaction-account>
              </td>
              <td class="transactions-table-body-cell">
                <drt-transaction-account
                  account={transaction.receiver}
                  dataTestId={DataTestIdsEnum.transactionReceiver}
                  scope="receiver"
                  showLockedAccounts={true}
                ></drt-transaction-account>
              </td>
              <td class="transactions-table-body-cell">
                <drt-transaction-method method={transaction.method.name} actionDescription={transaction.method.actionDescription}></drt-transaction-method>
              </td>
              <td class="transactions-table-body-cell">
                <drt-transaction-value value={transaction.value}></drt-transaction-value>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    );
  }
}
