import {
  GAS_LIMIT,
  GAS_PRICE,
  VERSION,
  isStringBase64,
  getStore,
  accountSelector,
  Address,
  Transaction,
  TransactionOptions,
  TransactionVersion,
  IPlainTransactionObject
} from 'lib';

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

  const { address } = accountSelector(getStore().getState());

  const dataPayload = data
    ? isStringBase64(data)
      ? Buffer.from(data, 'base64')
      : Buffer.from(data.trim())
    : undefined;

  const transaction = new Transaction({
    value: BigInt(value),
    data: dataPayload,
    receiver: new Address(receiver),
    sender: new Address(sender && sender !== '' ? sender : address),
    gasLimit: BigInt(gasLimit.valueOf() ?? GAS_LIMIT),
    gasPrice: BigInt(gasPrice.valueOf() ?? GAS_PRICE),
    chainID: chainID.valueOf(),
    version: new TransactionVersion(version ?? VERSION).valueOf()
  });

  if (options) {
    transaction.options = new TransactionOptions(options).valueOf();
  }

  return transaction;
};
