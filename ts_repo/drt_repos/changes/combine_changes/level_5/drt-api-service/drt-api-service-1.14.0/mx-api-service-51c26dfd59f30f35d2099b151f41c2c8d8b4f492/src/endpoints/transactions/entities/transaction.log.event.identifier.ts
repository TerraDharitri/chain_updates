export enum TransactionLogEventIdentifier {
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
  DCDTFreeze = 'DCDTFreeze',
  transferValueOnly = 'transferValueOnly',
  writeLog = 'writeLog',
  signalError = 'signalError'
}
