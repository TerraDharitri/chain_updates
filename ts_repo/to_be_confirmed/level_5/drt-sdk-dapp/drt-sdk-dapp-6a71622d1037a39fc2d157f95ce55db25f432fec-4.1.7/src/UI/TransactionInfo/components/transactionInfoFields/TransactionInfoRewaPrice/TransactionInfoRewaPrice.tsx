import React from 'react';
import { N_A } from 'constants/index';
import { withStyles, WithStylesImportType } from 'hocs/withStyles';
import { getUsdValue } from 'utils/operations/getUsdValue';
import {
  WithClassnameType,
  WithTransactionType
} from '../../../../../UI/types';
import { DetailItem } from '../../DetailItem';

const TransactionInfoRewaPriceComponent = ({
  className,
  transaction,
  styles
}: WithTransactionType & WithClassnameType & WithStylesImportType) => (
  <DetailItem className={className} title='REWA Price'>
    {transaction.price != null ? (
      getUsdValue({ amount: '1', usd: transaction.price })
    ) : (
      <span className={styles?.price}>{N_A}</span>
    )}
  </DetailItem>
);

export const TransactionInfoRewaPrice = withStyles(
  TransactionInfoRewaPriceComponent,
  {
    ssrStyles: () =>
      import(
        'UI/TransactionInfo/components/transactionInfoFields/TransactionInfoRewaPrice/styles.scss'
      ),
    clientStyles: () =>
      require('UI/TransactionInfo/components/transactionInfoFields/TransactionInfoRewaPrice/styles.scss')
        .default
  }
);
