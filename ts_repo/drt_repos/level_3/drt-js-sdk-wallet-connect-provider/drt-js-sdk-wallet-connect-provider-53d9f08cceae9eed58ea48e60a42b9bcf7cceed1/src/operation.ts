export enum Operation {
  SIGN_TRANSACTION = "drt_signTransaction",
  SIGN_TRANSACTIONS = "drt_signTransactions",
  SIGN_MESSAGE = "drt_signMessage",
}

export enum OptionalOperation {
  SIGN_LOGIN_TOKEN = "drt_signLoginToken",
  SIGN_NATIVE_AUTH_TOKEN = "drt_signNativeAuthToken",
  CANCEL_ACTION = "drt_cancelAction",
}
