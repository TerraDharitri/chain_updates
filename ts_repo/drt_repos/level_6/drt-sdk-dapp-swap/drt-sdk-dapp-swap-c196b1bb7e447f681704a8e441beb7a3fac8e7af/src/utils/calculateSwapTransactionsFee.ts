import BigNumber from 'bignumber.js';
import { IPlainTransactionObject } from 'lib';
import { getTransactionFee } from './getTransactionFee';

export const calculateSwapTransactionsFee = (
  transactions?: IPlainTransactionObject[]
): string | undefined =>
  transactions
    ?.reduce((acc, tx) => {
      const bnTransactionFee = new BigNumber(getTransactionFee(tx) ?? '0');
      return acc.plus(bnTransactionFee);
    }, new BigNumber('0'))
    .toString(10);
