import React from 'react';
import { render } from '@testing-library/react';
import { getMockTransactionByIndex } from 'UI/TransactionsTable/components/tests/getMockTransactionByIndex';
import { TransactionInfoTo } from '../TransactionInfoTo';

describe('TransactionInfoTo tests', () => {
  it('should display contract', () => {
    const methods = render(
      <TransactionInfoTo transaction={getMockTransactionByIndex(0)} />
    );

    const transactionToContract = methods.getByTestId('transactionToContract');

    expect(transactionToContract.textContent).toBe('Contract');

    const transactionToExplorerLink = methods.getByTestId(
      'transactionToExplorerLink'
    );

    expect(transactionToExplorerLink.getAttribute('href')).toBe(
      'http://devnet-explorer.dharitri.org/accounts/drt1yvesqqqqqqqqqqqqqqqqqqqqqqqqyvesqqqqqqqqqqqqqz8llllsan65ch'
    );

    const transactionToAccount = methods.getByTestId('transactionToAccount');

    expect(transactionToAccount.textContent).toBe('Staking: ARC Stake');

    expect(transactionToAccount.getAttribute('title')).toBe(
      'Staking: ARC Stake (drt1yvesqqqqqqqqqqqqqqqqqqqqqqqqyvesqqqqqqqqqqqqqz8llllsan65ch)'
    );

    const transactionToShard = methods.getByTestId('transactionToShard');

    expect(transactionToShard.textContent).toBe('Metachain');
  });
});
