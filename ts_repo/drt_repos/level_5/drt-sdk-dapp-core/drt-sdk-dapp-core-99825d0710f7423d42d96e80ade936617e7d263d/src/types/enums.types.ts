export enum EnvironmentsEnum {
  testnet = 'testnet',
  devnet = 'devnet',
  mainnet = 'mainnet'
}

export enum TransactionServerStatusesEnum {
  pending = 'pending',
  fail = 'fail',
  invalid = 'invalid',
  success = 'success',
  executed = 'executed',
  notExecuted = 'not executed',
  rewardReverted = 'reward-reverted'
}

export enum TransactionTypesEnum {
  MultiDCDTNFTTransfer = 'MultiDCDTNFTTransfer',
  DCDTTransfer = 'DCDTTransfer',
  DCDTNFTBurn = 'DCDTNFTBurn',
  DCDTNFTTransfer = 'DCDTNFTTransfer',
  dcdtTransaction = 'dcdtTransaction',
  nftTransaction = 'nftTransaction',
  scCall = 'scCall'
}

export enum TransactionBatchStatusesEnum {
  signed = 'signed',
  cancelled = 'cancelled',
  success = 'success',
  sent = 'sent',
  fail = 'fail',
  timedOut = 'timedOut',
  invalid = 'invalid'
}

export enum TypesOfSmartContractCallsEnum {
  MultiDCDTNFTTransfer = 'MultiDCDTNFTTransfer',
  DCDTNFTTransfer = 'DCDTNFTTransfer'
}

export enum DCDTTransferTypesEnum {
  DCDTNFTTransfer = 'DCDTNFTTransfer',
  DCDTNFTBurn = 'DCDTNFTBurn',
  DCDTNFTAddQuantity = 'DCDTNFTAddQuantity',
  DCDTNFTCreate = 'DCDTNFTCreate',
  MultiDCDTNFTTransfer = 'MultiDCDTNFTTransfer',
  DCDTTransfer = 'DCDTTransfer',
  DCDTBurn = 'DCDTBurn',
  DCDTLocalMint = 'DCDTLocalMint',
  DCDTLocalBurn = 'DCDTLocalBurn',
  DCDTWipe = 'DCDTWipe',
  DCDTFreeze = 'DCDTFreeze'
}

export enum GuardianActionsEnum {
  SetGuardian = 'SetGuardian',
  GuardAccount = 'GuardAccount',
  UnGuardAccount = 'UnGuardAccount'
}

export enum SigningErrorsEnum {
  notInitialized = 'Provider not initialized',
  errorSigning = 'Error when signing',
  errorSigningTx = 'Error signing transaction',
  missingProviderMessage = 'You need a signer/valid signer to send a transaction, use either WalletProvider, LedgerProvider or WalletConnect',
  defaultTransactionStatusMessage = 'Undefined transaction status',
  secondLoginAttemptError = 'Action not allowed. User is logged in. Call logout() first',
  senderDifferentThanLoggedInAddress = 'You cannot sign transactions from a different account'
}
export enum SigningWarningsEnum {
  cancelled = 'Signing canceled',
  transactionCancelled = 'Transaction canceled'
}
