import BigNumber from 'bignumber.js';
import { InterpretedTransactionType } from 'types/serverTransactions.types';
import { getRewaValueData } from './getRewaValueData';

let warningLogged = false;

export function getValueFromActions(transaction: InterpretedTransactionType) {
  const value = new BigNumber(transaction.action?.arguments?.value);

  if (!value.isNaN()) {
    return getRewaValueData(transaction.action?.arguments?.value);
  }

  if (!warningLogged) {
    console.error(
      `Unable to interpret ${transaction.action?.name} data for txHash: ${transaction.txHash}`
    );
    warningLogged = true;
  }

  // fallback on transaction value
  return getRewaValueData(transaction.value);
}
