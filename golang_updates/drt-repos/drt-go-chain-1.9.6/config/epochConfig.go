package config

// EpochConfig will hold epoch configurations
type EpochConfig struct {
	EnableEpochs EnableEpochs
	GasSchedule  GasScheduleConfig
}

// GasScheduleConfig represents the versioning config area for the gas schedule toml
type GasScheduleConfig struct {
	GasScheduleByEpochs []GasScheduleByEpochs
}

// EnableEpochs will hold the configuration for activation epochs
type EnableEpochs struct {
	SCDeployEnableEpoch                                      uint32
	BuiltInFunctionsEnableEpoch                              uint32
	RelayedTransactionsEnableEpoch                           uint32
	PenalizedTooMuchGasEnableEpoch                           uint32
	SwitchJailWaitingEnableEpoch                             uint32
	SwitchHysteresisForMinNodesEnableEpoch                   uint32
	BelowSignedThresholdEnableEpoch                          uint32
	TransactionSignedWithTxHashEnableEpoch                   uint32
	MetaProtectionEnableEpoch                                uint32
	AheadOfTimeGasUsageEnableEpoch                           uint32
	GasPriceModifierEnableEpoch                              uint32
	RepairCallbackEnableEpoch                                uint32
	MaxNodesChangeEnableEpoch                                []MaxNodesChangeConfig
	BlockGasAndFeesReCheckEnableEpoch                        uint32
	StakingV2EnableEpoch                                     uint32
	StakeEnableEpoch                                         uint32
	DoubleKeyProtectionEnableEpoch                           uint32
	DCDTEnableEpoch                                          uint32
	GovernanceEnableEpoch                                    uint32
	DelegationManagerEnableEpoch                             uint32
	DelegationSmartContractEnableEpoch                       uint32
	CorrectLastUnjailedEnableEpoch                           uint32
	BalanceWaitingListsEnableEpoch                           uint32
	ReturnDataToLastTransferEnableEpoch                      uint32
	SenderInOutTransferEnableEpoch                           uint32
	RelayedTransactionsV2EnableEpoch                         uint32
	UnbondTokensV2EnableEpoch                                uint32
	SaveJailedAlwaysEnableEpoch                              uint32
	ValidatorToDelegationEnableEpoch                         uint32
	ReDelegateBelowMinCheckEnableEpoch                       uint32
	IncrementSCRNonceInMultiTransferEnableEpoch              uint32
	ScheduledMiniBlocksEnableEpoch                           uint32
	DCDTMultiTransferEnableEpoch                             uint32
	GlobalMintBurnDisableEpoch                               uint32
	DCDTTransferRoleEnableEpoch                              uint32
	ComputeRewardCheckpointEnableEpoch                       uint32
	SCRSizeInvariantCheckEnableEpoch                         uint32
	BackwardCompSaveKeyValueEnableEpoch                      uint32
	DCDTNFTCreateOnMultiShardEnableEpoch                     uint32
	MetaDCDTSetEnableEpoch                                   uint32
	AddTokensToDelegationEnableEpoch                         uint32
	MultiDCDTTransferFixOnCallBackOnEnableEpoch              uint32
	OptimizeGasUsedInCrossMiniBlocksEnableEpoch              uint32
	CorrectFirstQueuedEpoch                                  uint32
	CorrectJailedNotUnstakedEmptyQueueEpoch                  uint32
	FixOOGReturnCodeEnableEpoch                              uint32
	RemoveNonUpdatedStorageEnableEpoch                       uint32
	DeleteDelegatorAfterClaimRewardsEnableEpoch              uint32
	OptimizeNFTStoreEnableEpoch                              uint32
	CreateNFTThroughExecByCallerEnableEpoch                  uint32
	StopDecreasingValidatorRatingWhenStuckEnableEpoch        uint32
	FrontRunningProtectionEnableEpoch                        uint32
	IsPayableBySCEnableEpoch                                 uint32
	CleanUpInformativeSCRsEnableEpoch                        uint32
	StorageAPICostOptimizationEnableEpoch                    uint32
	TransformToMultiShardCreateEnableEpoch                   uint32
	DCDTRegisterAndSetAllRolesEnableEpoch                    uint32
	DoNotReturnOldBlockInBlockchainHookEnableEpoch           uint32
	AddFailedRelayedTxToInvalidMBsDisableEpoch               uint32
	SCRSizeInvariantOnBuiltInResultEnableEpoch               uint32
	CheckCorrectTokenIDForTransferRoleEnableEpoch            uint32
	DisableExecByCallerEnableEpoch                           uint32
	FailExecutionOnEveryAPIErrorEnableEpoch                  uint32
	ManagedCryptoAPIsEnableEpoch                             uint32
	RefactorContextEnableEpoch                               uint32
	CheckFunctionArgumentEnableEpoch                         uint32
	CheckExecuteOnReadOnlyEnableEpoch                        uint32
	MiniBlockPartialExecutionEnableEpoch                     uint32
	DCDTMetadataContinuousCleanupEnableEpoch                 uint32
	FixAsyncCallBackArgsListEnableEpoch                      uint32
	FixOldTokenLiquidityEnableEpoch                          uint32
	RuntimeMemStoreLimitEnableEpoch                          uint32
	RuntimeCodeSizeFixEnableEpoch                            uint32
	SetSenderInEeiOutputTransferEnableEpoch                  uint32
	RefactorPeersMiniBlocksEnableEpoch                       uint32
	SCProcessorV2EnableEpoch                                 uint32
	MaxBlockchainHookCountersEnableEpoch                     uint32
	WipeSingleNFTLiquidityDecreaseEnableEpoch                uint32
	AlwaysSaveTokenMetaDataEnableEpoch                       uint32
	SetGuardianEnableEpoch                                   uint32
	ScToScLogEventEnableEpoch                                uint32
	RelayedNonceFixEnableEpoch                               uint32
	DeterministicSortOnValidatorsInfoEnableEpoch             uint32
	KeepExecOrderOnCreatedSCRsEnableEpoch                    uint32
	MultiClaimOnDelegationEnableEpoch                        uint32
	ChangeUsernameEnableEpoch                                uint32
	AutoBalanceDataTriesEnableEpoch                          uint32
	MigrateDataTrieEnableEpoch                               uint32
	ConsistentTokensValuesLengthCheckEnableEpoch             uint32
	FixDelegationChangeOwnerOnAccountEnableEpoch             uint32
	DynamicGasCostForDataTrieStorageLoadEnableEpoch          uint32
	NFTStopCreateEnableEpoch                                 uint32
	ChangeOwnerAddressCrossShardThroughSCEnableEpoch         uint32
	FixGasRemainingForSaveKeyValueBuiltinFunctionEnableEpoch uint32
	CurrentRandomnessOnSortingEnableEpoch                    uint32
	StakeLimitsEnableEpoch                                   uint32
	StakingV4Step1EnableEpoch                                uint32
	StakingV4Step2EnableEpoch                                uint32
	StakingV4Step3EnableEpoch                                uint32
	CleanupAuctionOnLowWaitingListEnableEpoch                uint32
	AlwaysMergeContextsInEEIEnableEpoch                      uint32
	UseGasBoundedShouldFailExecutionEnableEpoch              uint32
	DynamicDCDTEnableEpoch                                   uint32
	REWAInMultiTransferEnableEpoch                           uint32
	CryptoOpcodesV2EnableEpoch                               uint32
	UnJailCleanupEnableEpoch                                 uint32
	FixRelayedBaseCostEnableEpoch                            uint32
	MultiDCDTNFTTransferAndExecuteByUserEnableEpoch          uint32
	FixRelayedMoveBalanceToNonPayableSCEnableEpoch           uint32
	RelayedTransactionsV3EnableEpoch                         uint32
	RelayedTransactionsV3FixDCDTTransferEnableEpoch          uint32
	AndromedaEnableEpoch                                     uint32
	CheckBuiltInCallOnTransferValueAndFailEnableRound        uint32
	BLSMultiSignerEnableEpoch                                []MultiSignerConfig
}

// GasScheduleByEpochs represents a gas schedule toml entry that will be applied from the provided epoch
type GasScheduleByEpochs struct {
	StartEpoch uint32
	FileName   string
}
