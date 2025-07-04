import { InterpretedTransactionType } from 'types/serverTransactions.types';
import { getRewaValueData } from './getRewaValueData';
import { getVisibleOperations } from './getVisibleOperations';

let warningLogged = false;

const logError = (hash: string) => {
  if (!warningLogged) {
    console.error(
      `Operations field missing for txHash: ${hash}.
        Unable to compute value field.`
    );
    warningLogged = true;
  }
};

export function getValueFromOperations(
  transaction: InterpretedTransactionType
) {
  try {
    if (transaction.operations) {
      const [operation] = getVisibleOperations(transaction);
      return getRewaValueData(operation?.value);
    } else {
      logError(transaction.txHash);
    }
  } catch {
    logError(transaction.txHash);
  }
  return getRewaValueData(transaction.value);
}
