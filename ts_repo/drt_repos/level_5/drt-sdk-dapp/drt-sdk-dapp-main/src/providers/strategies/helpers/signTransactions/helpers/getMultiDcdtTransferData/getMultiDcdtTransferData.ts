import { Transaction } from '@terradharitri/sdk-core';
import {
  MultiSignTransactionType,
  TransactionDataTokenType
} from 'types/transactions.types';
import { parseMultiDcdtTransferDataForMultipleTransactions } from './helpers/parseMultiDcdtTransferDataForMultipleTransactions';

export type MultiDcdtTransferDataReturnType = ReturnType<
  typeof getMultiDcdtTransferData
>;

export function getMultiDcdtTransferData(transactions?: Transaction[]): {
  parsedTransactionsByDataField: Record<string, TransactionDataTokenType>;
  allTransactions: MultiSignTransactionType[];
} {
  const { allTransactions, parsedTransactionsByDataField } =
    parseMultiDcdtTransferDataForMultipleTransactions({ transactions });

  return {
    parsedTransactionsByDataField,
    allTransactions
  };
}
