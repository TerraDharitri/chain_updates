import { DECIMALS } from 'lib/sdkDappUtils';
import { ServerTransactionType } from 'types/serverTransactions.types';

import { getAssetAmount } from './getAssetAmount';
import { getAssetPrice } from './getAssetPrice';

interface IGetTransactionAssetsParams {
  transaction: ServerTransactionType;
  userIsReceiver: boolean;
  rewaLabel?: string;
}

export interface ITransactionAsset {
  assetPrefix: string;
  assetTicker: string;
  assetAmount: string;
  assetImage?: string;
  assetPrice?: string;
  type: string;
}

export const getTransactionAssets = ({
  transaction,
  userIsReceiver,
  rewaLabel
}: IGetTransactionAssetsParams): ITransactionAsset[] => {
  const transactionAction = transaction.action;
  const transactionArguments = transactionAction && transactionAction.arguments;

  const transactionTransfers: Record<string, string>[] =
    transactionArguments && Array.isArray(transactionArguments.transfers)
      ? transactionArguments.transfers
      : [];

  const isRewaTransfer = transactionTransfers.length === 0;
  const processedRewaLabel = rewaLabel ?? 'REWA';
  const assetPrefix = userIsReceiver ? '+' : '-';

  const rewaTransferAsset: ITransactionAsset = {
    assetPrefix,
    type: processedRewaLabel,
    assetTicker: processedRewaLabel,
    assetAmount: getAssetAmount({
      value: transaction.value,
      decimals: String(DECIMALS)
    })
  };

  const transfersAssets = transactionTransfers.map(
    (transactionTransfer): ITransactionAsset => ({
      assetPrefix,
      type: transactionTransfer.type,
      assetTicker: transactionTransfer.ticker,
      assetImage: transactionTransfer.svgUrl,
      assetAmount: getAssetAmount(transactionTransfer),
      assetPrice: getAssetPrice(transactionTransfer)
    })
  );

  if (isRewaTransfer) {
    return [rewaTransferAsset];
  }

  return transfersAssets;
};
