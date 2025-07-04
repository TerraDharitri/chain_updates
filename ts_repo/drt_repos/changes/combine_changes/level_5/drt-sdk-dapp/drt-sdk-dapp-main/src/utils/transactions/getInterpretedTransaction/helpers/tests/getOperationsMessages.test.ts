import {
  OperationType,
  TransactionOperationActionTypeEnum,
  VisibleTransactionOperationType
} from 'types/serverTransactions.types';
import { getOperationsMessages } from '../getOperationsMessages';
import { baseTransactionMock } from './base-transaction-mock';

describe('getOperationsMessages', () => {
  it('receive empty array if no operations messages exists on the transaction', () => {
    const result = getOperationsMessages(baseTransactionMock);

    expect(result).toEqual([]);
  });

  it('receive an array with all operations messages', () => {
    const baseTransactionOperation = {
      name: 'send',
      type: VisibleTransactionOperationType.rewa,
      action: TransactionOperationActionTypeEnum.transfer,
      dcdtType: 'FungibleDCDT',
      receiver:
        'drt1qqqqqqqqqqqqqpgq4gdcg0k83u7lpv4s4532w3au9y9h0vm70eqq88sr45',
      sender: 'drt1axhx4kenjlae6sknq7zjg2g4fvzavv979r2fg425p62wkl84avtqd4f00j',
      value: '1',
      decimals: 4,
      identifier: 'REWA'
    };

    const transaction = {
      ...baseTransactionMock,
      operations: [
        {
          ...baseTransactionOperation,
          message: 'message 0'
        },
        {
          ...baseTransactionOperation,
          message: 'message 1'
        }
      ] as OperationType[]
    };

    const result = getOperationsMessages(transaction);

    expect(result).toEqual(['message 0', 'message 1']);
  });
});
