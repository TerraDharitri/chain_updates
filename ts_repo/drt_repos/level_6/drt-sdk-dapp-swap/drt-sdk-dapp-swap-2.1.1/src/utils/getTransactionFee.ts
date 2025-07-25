import { IPlainTransactionObject, NetworkConfig } from '@terradharitri/sdk-core';
import {
  GAS_LIMIT,
  GAS_PER_DATA_BYTE,
  GAS_PRICE_MODIFIER
} from '@terradharitri/sdk-dapp/constants';
import { createTransactionFromRaw } from './createTransactionFromRaw';

export const getTransactionFee = (rawTransaction: IPlainTransactionObject) => {
  const transaction = createTransactionFromRaw(rawTransaction);

  const networkConfig = new NetworkConfig();
  networkConfig.MinGasLimit = GAS_LIMIT;
  networkConfig.GasPerDataByte = GAS_PER_DATA_BYTE;
  networkConfig.GasPriceModifier = GAS_PRICE_MODIFIER;

  try {
    return transaction.computeFee(networkConfig).toString(10);
  } catch (err) {
    return 0;
  }
};
