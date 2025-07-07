import { IPlainTransactionObject } from 'lib';

export interface WrappingQueryResponseType {
  wrapEgld?: IPlainTransactionObject;
  unwrapEgld?: IPlainTransactionObject;
}
