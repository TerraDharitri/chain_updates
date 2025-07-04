import type { SessionTransactionType } from 'types/transactions.types';

export enum TransactionIconTypeEnum {
  FAILED = 'failed',
  MULTIPLE_ASSETS = 'multiple_assets',
  ASSET = 'asset',
  NFT = 'nft',
  SFT = 'sft',
  CONTRACT = 'contract',
  SYMBOL = 'symbol'
}
export interface IBaseTransactionParams {
  address: string;
  explorerAddress: string;
  rewaLabel: string;
}

export interface IGetHistoricalTransactionsParams
  extends IBaseTransactionParams {
  transactionsSessions: Record<string, SessionTransactionType>;
}

export enum TransactionActionMethodTypeEnum {
  transfer = 'Transfer',
  delegate = 'Delegate',
  stake = 'Stake',
  unDelegate = 'Undelegate',
  stakeClaimRewards = 'Stake Claim Rewards',
  reDelegateRewards = 'Redelegate Rewards',
  withdraw = 'Withdraw',
  claimLockedAssets = 'Claim Locked Assets',
  swapTokensFixedInput = 'Swap',
  swapTokensFixedOutput = 'Swap',
  swap = 'Swap',
  multiPairSwap = 'Multiple Pair Swap',
  aggregateRewa = 'Aggregate REWA',
  addLiquidity = 'Add Liquidity',
  addLiquidityProxy = 'Add Liquidity Proxy',
  removeLiquidity = 'Remove Liquidity',
  removeLiquidityProxy = 'Remove Liquidity Proxy',
  enterFarm = 'Enter Farm',
  enterFarmProxy = 'Enter Farm Proxy',
  enterFarmAndLockRewards = 'Enter Farm & Lock Rewards',
  enterFarmAndLockRewardsProxy = 'Enter Farm & Lock Rewards Proxy',
  exitFarm = 'Exit Farm',
  exitFarmProxy = 'Exit Farm Proxy',
  claimRewards = 'Claim Rewards',
  claimRewardsProxy = 'Claim Rewards Proxy',
  compoundRewards = 'Reinvest Rewards',
  compoundRewardsProxy = 'Reinvest Rewards Proxy',
  createNftMinter = 'Create NFT Minter',
  scDeploy = 'Smart Contract Deploy',
  wrapRewa = 'Wrap REWA',
  unwrapRewa = 'Unwrap REWA',
  lockAssets = 'Lock Assets',
  unlockAssets = 'Unlock Assets',
  mergeLockedAssetTokens = 'Merge Locked Tokens',
  stakeFarm = 'Stake Farm',
  stakeFarmProxy = 'Stake Farm Proxy',
  stakeFarmTokens = 'Stake Farm Tokens',
  stakeFarmTokensProxy = 'Stake Farm Tokens Proxy',
  unstakeFarm = 'Unstake Farm',
  unstakeFarmProxy = 'Unstake Farm Proxy Proxy',
  unstakeFarmTokens = 'Unstake Farm Tokens',
  unstakeFarmTokensProxy = 'Unstake Farm Proxy',
  claimDualYield = 'Claim Dual Yield',
  claimDualYieldProxy = 'Claim Dual Yield Proxy',
  unbondFarm = 'Unbond Farm',
  ClaimDeveloperRewards = 'Claim Developer Rewards',
  ChangeOwnerAddress = 'Change Owner Address',
  SetUserName = 'Set Username',
  SaveKeyValue = 'Save Key Value',
  DCDTTransfer = 'Transfer DCDT',
  DCDTBurn = 'Burn DCDT',
  DCDTFreeze = 'Freeze DCDT',
  DCDTUnFreeze = 'Unfreeze DCDT',
  DCDTWipe = 'Wipe DCDT',
  DCDTPause = 'Pause DCDT',
  DCDTUnPause = 'Unpause DCDTBurn',
  DCDTSetRole = 'Set DCDT Role',
  DCDTUnSetRole = 'Unset DCDT Role',
  DCDTSetLimitedTransfer = 'Set DCDT Limited Transfer',
  DCDTUnSetLimitedTransfer = 'Unset DCDT Limited Transfer',
  DCDTLocalBurn = 'Local DCDT Burn',
  DCDTLocalMint = 'Local DCDT Mint',
  DCDTNFTTransfer = 'Transfer DCDT/NFT',
  DCDTNFTCreate = 'Create DCDT/NFT',
  DCDTNFTAddQuantity = 'Add DCDT/NFT Quantity',
  DCDTNFTBurn = 'Burn DCDT/NFT',
  DCDTNFTAddURI = 'Add DCDT/NFT URI',
  DCDTNFTUpdateAttributes = 'Update DCDT/NFT Attributes',
  MultiDCDTNFTTransfer = 'Multiple DCDT/NFT Transfer',
  SetGuardian = 'Register Guardian',
  GuardAccount = 'Activate Guardian',
  UnGuardAccount = 'Deactivate Guardian'
}
