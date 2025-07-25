import type { Transaction } from '@terradharitri/sdk-core';

import {
  MultiSignTransactionType,
  TransactionsDataTokensType
} from 'types/transactions.types';
import { getTokenFromData } from './getTokenFromData';
import { parseMultiDcdtTransferData } from './parseMultiDcdtTransferData';

export function parseMultiDcdtTransferDataForMultipleTransactions({
  transactions
}: {
  transactions?: Transaction[];
}) {
  const allTransactions: MultiSignTransactionType[] = [];
  const parsedTransactionsByDataField: TransactionsDataTokensType = {};

  if (!transactions || transactions.length === 0) {
    return {
      allTransactions,
      parsedTransactionsByDataField
    };
  }

  transactions.forEach((transaction, transactionIndex) => {
    const txData = transaction.getData().toString();
    const multiTxs = parseMultiDcdtTransferData(txData);

    if (multiTxs.length > 0) {
      multiTxs.forEach((trx, idx) => {
        const isLastView = idx === multiTxs.length - 1;

        const newTx: MultiSignTransactionType = {
          transaction,
          multiTxData: trx.data,
          transactionIndex: idx,
          needsSigning: isLastView
        };

        parsedTransactionsByDataField[trx.data] = {
          tokenId: trx.token ? trx.token : '',
          amount: trx.amount ? trx.amount : '',
          type: trx.type,
          nonce: trx.nonce ? trx.nonce : '',
          multiTxData: trx.data,
          receiver: trx.receiver
        };

        allTransactions.push(newTx);
      });
    } else {
      const transactionData = transaction.getData().toString();

      const { tokenId, amount } = getTokenFromData(transactionData);

      if (tokenId) {
        parsedTransactionsByDataField[transactionData] = {
          tokenId,
          amount,
          receiver: transaction.getReceiver().bech32()
        };
      }
      allTransactions.push({
        transaction,
        transactionIndex,
        multiTxData: transactionData,
        needsSigning: true
      });
    }
  });

  return {
    allTransactions,
    parsedTransactionsByDataField
  };
}
