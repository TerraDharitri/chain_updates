import {
  Address,
  Transaction,
  TransactionOptions,
  TransactionPayload,
  TransactionVersion,
  IPlainTransactionObject
} from '@terradharitri/sdk-core';
import { GAS_LIMIT, GAS_PRICE, VERSION } from '@terradharitri/sdk-dapp/constants';
import { accountSelector } from '@terradharitri/sdk-dapp/reduxStore/selectors';
import { store } from '@terradharitri/sdk-dapp/reduxStore/store';
import { isStringBase64 } from '@terradharitri/sdk-dapp/utils/decoders/base64Utils';

export const createTransactionFromRaw = (
  rawTransaction: IPlainTransactionObject
) => {
  const {
    data,
    value,
    receiver,
    sender,
    gasLimit,
    gasPrice,
    chainID,
    version,
    options
  } = rawTransaction;

  const { address } = accountSelector(store.getState());

  const dataPayload = isStringBase64(data ?? '')
    ? TransactionPayload.fromEncoded(data)
    : new TransactionPayload(data);

  return new Transaction({
    value: value.valueOf(),
    data: dataPayload,
    receiver: new Address(receiver),
    sender: new Address(sender && sender !== '' ? sender : address),
    gasLimit: gasLimit.valueOf() ?? GAS_LIMIT,
    gasPrice: gasPrice.valueOf() ?? GAS_PRICE,
    chainID: chainID.valueOf(),
    version: new TransactionVersion(version ?? VERSION),
    ...(options ? { options: new TransactionOptions(options) } : {})
  });
};
