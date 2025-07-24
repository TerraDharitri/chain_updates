import React from 'react';
import { faCoins } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import classNames from 'classnames';

import { DataTestIdsEnum, MULTI_TRANSFER_REWA_TOKEN } from 'constants/index';
import { withStyles } from 'hocs/withStyles';
import { TokenOptionType, useGetNetworkConfig } from 'hooks';
import { ActiveLedgerTransactionType } from 'types';
import { Balance } from 'UI/Balance';
import { UsdValue } from 'UI/UsdValue';
import { formatAmount } from 'utils';

import { WithStylesImportType } from '../../../../../../../../hocs/useStyles';

export interface ConfirmAmountDataPropsType extends WithStylesImportType {
  isRewa: boolean;
  tokenPrice?: number;
  isNftOrSft: boolean;
  amount: string;
  handleReference: (element: HTMLElement | null) => void;
  currentTransaction: ActiveLedgerTransactionType;
  tokenDetails: TokenOptionType;
}

const ConfirmAmountDataComponent = ({
  isRewa,
  styles,
  tokenPrice,
  isNftOrSft,
  handleReference,
  currentTransaction,
  amount,
  tokenDetails
}: ConfirmAmountDataPropsType) => {
  const { network } = useGetNetworkConfig();
  const { tokenAvatar, tokenDecimals, identifier } = tokenDetails;

  const getFormattedAmount = ({ addCommas }: { addCommas: boolean }) =>
    formatAmount({
      input: isRewa
        ? currentTransaction.transaction.getValue().toString()
        : amount,
      decimals: isRewa ? Number(network.decimals) : tokenDecimals,
      digits: Number(network.digits),
      showLastNonZeroDecimal: false,
      addCommas
    });

  const isRewaToken =
    isRewa || tokenDetails.identifier === MULTI_TRANSFER_REWA_TOKEN;
  const formattedAmount = getFormattedAmount({ addCommas: true });
  const rawAmount = getFormattedAmount({ addCommas: false });

  return (
    <div className={styles?.confirmAmountData}>
      <div className={styles?.confirmAmountDataWrapper}>
        {!isRewaToken && tokenAvatar && (
          <img src={tokenAvatar} className={styles?.confirmAmountDataIcon} />
        )}

        {!isRewaToken && !tokenAvatar && (
          <div
            className={classNames(
              styles?.confirmAmountDataIcon,
              styles?.confirmAmountDataIconDefault
            )}
          >
            <FontAwesomeIcon
              icon={faCoins}
              className={styles?.confirmAmountDataIconDefaultIcon}
            />
          </div>
        )}

        <div
          className={styles?.confirmAmountDataBalanceWrapper}
          ref={handleReference}
        >
          <Balance
            amount={formattedAmount}
            rewaIcon={isRewaToken}
            data-testid={DataTestIdsEnum.confirmAmount}
            showTokenLabel
            showTokenLabelSup
            tokenLabel={isRewaToken ? network.rewaLabel : identifier}
            className={styles?.confirmAmountDataBalance}
          />
        </div>
      </div>

      {!isNftOrSft && tokenPrice && (
        <UsdValue
          amount={rawAmount}
          usd={tokenPrice}
          data-testid={DataTestIdsEnum.confirmUsdValue}
          className={styles?.confirmAmountDataPrice}
        />
      )}

      {!isNftOrSft && !tokenPrice && (
        <div className={styles?.confirmAmountDataPrice}>Price Unknown</div>
      )}
    </div>
  );
};

export const ConfirmAmountData = withStyles(ConfirmAmountDataComponent, {
  ssrStyles: () =>
    import(
      'UI/SignTransactionsModals/SignWithDeviceModal/components/components/ConfirmAmount/components/ConfirmAmountData/confirmAmountDataStyles.scss'
    ),
  clientStyles: () =>
    require('UI/SignTransactionsModals/SignWithDeviceModal/components/components/ConfirmAmount/components/ConfirmAmountData/confirmAmountDataStyles.scss')
      .default
});
