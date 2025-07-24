import { useEffect, useState } from 'react';
import { Transaction } from '@terradharitri/sdk-core';
import {
  MultiSignTransactionType,
  TransactionDataTokenType,
  TransactionsDataTokensType
} from 'types';
import { parseMultiDcdtTransferDataForMultipleTransactions } from 'utils/transactions/parseMultiDcdtTransferDataForMultipleTransactions';

const defaultTransactionInfo: TransactionDataTokenType = {
  tokenId: '',
  amount: '',
  type: '',
  multiTxData: '',
  receiver: ''
};

interface UseParseMultiDcdtTransferDataPropsType {
  transactions?: Transaction[];
}

interface UseParseMultiDcdtTransferDataReturnType {
  parsedTransactionsByDataField: TransactionsDataTokensType;
  getTxInfoByDataField: (
    data: string,
    multiTransactionData?: string
  ) => TransactionDataTokenType;
  allTransactions: MultiSignTransactionType[];
}

export function useParseMultiDcdtTransferData({
  transactions
}: UseParseMultiDcdtTransferDataPropsType): UseParseMultiDcdtTransferDataReturnType {
  const [parsedTransactionsByDataField, setParsedTransactions] =
    useState<TransactionsDataTokensType>({});
  const [allTransactions, setAllTransactions] = useState<
    MultiSignTransactionType[]
  >([]);

  function getTxInfoByDataField(
    data: string,
    multiTransactionData?: string
  ): TransactionDataTokenType {
    if (parsedTransactionsByDataField == null) {
      return defaultTransactionInfo;
    }

    if (data in parsedTransactionsByDataField) {
      return parsedTransactionsByDataField[data];
    }

    if (
      multiTransactionData != null &&
      String(multiTransactionData) in parsedTransactionsByDataField
    ) {
      return parsedTransactionsByDataField[multiTransactionData];
    }

    return defaultTransactionInfo;
  }

  function extractTransactionDCDTData() {
    const { allTransactions, parsedTransactionsByDataField } =
      parseMultiDcdtTransferDataForMultipleTransactions({
        transactions
      });

    setAllTransactions(allTransactions);
    setParsedTransactions(parsedTransactionsByDataField);
  }

  useEffect(() => {
    extractTransactionDCDTData();
  }, [transactions?.length]);

  return {
    parsedTransactionsByDataField,
    getTxInfoByDataField,
    allTransactions
  };
}
