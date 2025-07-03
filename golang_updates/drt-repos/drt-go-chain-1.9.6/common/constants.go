package common

import (
	"math"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
)

// NodeOperation defines the p2p node operation
type NodeOperation string

// NormalOperation defines the normal mode operation: either seeder, observer or validator
const NormalOperation NodeOperation = "normal operation"

// FullArchiveMode defines the node operation as a full archive mode
const FullArchiveMode NodeOperation = "full archive mode"

// PeerType represents the type of peer
type PeerType string

// EligibleList represents the list of peers who participate in consensus inside a shard
const EligibleList PeerType = "eligible"

// WaitingList represents the list of peers who don't participate in consensus but will join the next epoch
const WaitingList PeerType = "waiting"

// LeavingList represents the list of peers who were taken out of eligible and waiting because of rating
const LeavingList PeerType = "leaving"

// InactiveList represents the list of peers who were taken out because they were leaving
const InactiveList PeerType = "inactive"

// JailedList represents the list of peers who have stake but are in jail
const JailedList PeerType = "jailed"

// ObserverList represents the list of peers who don't participate in consensus but will join the next epoch
const ObserverList PeerType = "observer"

// NewList represents the list of peers who have stake and are pending to become eligible
const NewList PeerType = "new"

// MetachainTopicIdentifier is the identifier used in topics to define the metachain shard ID
const MetachainTopicIdentifier = "META" // TODO - move this to drt-go-chain-core and change wherever we use the string value

// AuctionList represents the list of peers which don't participate in consensus yet, but will be selected
// based on their top up stake
const AuctionList PeerType = "auction"

// SelectedFromAuctionList represents the list of peers which have been selected from AuctionList based on
// their top up to be distributed on the WaitingList in the next epoch
const SelectedFromAuctionList PeerType = "selectedFromAuction"

// CombinedPeerType - represents the combination of two peerTypes
const CombinedPeerType = "%s (%s)"

// UnVersionedAppString represents the default app version that indicate that the binary wasn't build by setting
// the appVersion flag
const UnVersionedAppString = "undefined"

// DisabledShardIDAsObserver defines the uint32 identifier which tells that the node hasn't configured any preferred
// shard to start in as observer
const DisabledShardIDAsObserver = uint32(0xFFFFFFFF) - 7

// MaxTxNonceDeltaAllowed specifies the maximum difference between an account's nonce and a received transaction's nonce
// in order to mark the transaction as valid.
const MaxTxNonceDeltaAllowed = 100

// MaxBulkTransactionSize specifies the maximum size of one bulk with txs which can be send over the network
// TODO convert this const into a var and read it from config when this code moves to another binary
const MaxBulkTransactionSize = 1 << 18 // 256KB bulks

// MaxTxsToRequest specifies the maximum number of txs to request
const MaxTxsToRequest = 1000

// NodesSetupJsonFileName specifies the name of the json file which contains the setup of the nodes
const NodesSetupJsonFileName = "nodesSetup.json"

// ConsensusTopic is the topic used in consensus algorithm
const ConsensusTopic = "consensus"

// GenesisTxSignatureString is the string used to generate genesis transaction signature as 128 hex characters
const GenesisTxSignatureString = "GENESISGENESISGENESISGENESISGENESISGENESISGENESISGENESISGENESISG"

// HeartbeatV2Topic is the topic used for heartbeatV2 signaling
const HeartbeatV2Topic = "heartbeatV2"

// PeerAuthenticationTopic is the topic used for peer authentication signaling
const PeerAuthenticationTopic = "peerAuthentication"

// ConnectionTopic represents the topic used when sending the new connection message data
const ConnectionTopic = "connection"

// ValidatorInfoTopic is the topic used for validatorInfo signaling
const ValidatorInfoTopic = "validatorInfo"

// EquivalentProofsTopic is the topic used for equivalent proofs
const EquivalentProofsTopic = "equivalentProofs"

// MetricCurrentRound is the metric for monitoring the current round of a node
const MetricCurrentRound = "drt_current_round"

// MetricNonce is the metric for monitoring the nonce of a node
const MetricNonce = "drt_nonce"

// MetricBlockTimestamp is the metric for monitoring the timestamp of the last synchronized block
const MetricBlockTimestamp = "drt_block_timestamp"

// MetricProbableHighestNonce is the metric for monitoring the max speculative nonce received by the node by listening on the network
const MetricProbableHighestNonce = "drt_probable_highest_nonce"

// MetricNumConnectedPeers is the metric for monitoring the number of connected peers
const MetricNumConnectedPeers = "drt_num_connected_peers"

// MetricNumConnectedPeersClassification is the metric for monitoring the number of connected peers split on the connection type
const MetricNumConnectedPeersClassification = "drt_num_connected_peers_classification"

// MetricSynchronizedRound is the metric for monitoring the synchronized round of a node
const MetricSynchronizedRound = "drt_synchronized_round"

// MetricIsSyncing is the metric for monitoring if a node is syncing
const MetricIsSyncing = "drt_is_syncing"

// MetricPublicKeyBlockSign is the metric for monitoring public key of a node used in block signing
const MetricPublicKeyBlockSign = "drt_public_key_block_sign"

// MetricShardId is the metric for monitoring shard id of a node
const MetricShardId = "drt_shard_id"

// MetricNumShardsWithoutMetachain is the metric for monitoring the number of shards (excluding meta)
const MetricNumShardsWithoutMetachain = "drt_num_shards_without_meta"

// MetricTxPoolLoad is the metric for monitoring number of transactions from pool of a node
const MetricTxPoolLoad = "drt_tx_pool_load"

// MetricCountLeader is the metric for monitoring number of rounds when a node was leader
const MetricCountLeader = "drt_count_leader"

// MetricCountConsensus is the metric for monitoring number of rounds when a node was in consensus group
const MetricCountConsensus = "drt_count_consensus"

// MetricCountAcceptedBlocks is the metric for monitoring number of blocks that was accepted proposed by a node
const MetricCountAcceptedBlocks = "drt_count_accepted_blocks"

// MetricNodeType is the metric for monitoring the type of the node
const MetricNodeType = "drt_node_type"

// MetricLiveValidatorNodes is the metric for the number of live validators on the network
const MetricLiveValidatorNodes = "drt_live_validator_nodes"

// MetricConnectedNodes is the metric for monitoring total connected nodes on the network
const MetricConnectedNodes = "drt_connected_nodes"

// MetricNumIntraShardValidatorNodes is the metric for the number of intra-shard validators
const MetricNumIntraShardValidatorNodes = "drt_intra_shard_validator_nodes"

// MetricCpuLoadPercent is the metric for monitoring CPU load [%]
const MetricCpuLoadPercent = "drt_cpu_load_percent"

// MetricMemLoadPercent is the metric for monitoring memory load [%]
const MetricMemLoadPercent = "drt_mem_load_percent"

// MetricMemTotal is the metric for monitoring total memory bytes
const MetricMemTotal = "drt_mem_total"

// MetricMemUsedGolang is a metric for monitoring the memory ("total")
const MetricMemUsedGolang = "drt_mem_used_golang"

// MetricMemUsedSystem is a metric for monitoring the memory ("sys mem")
const MetricMemUsedSystem = "drt_mem_used_sys"

// MetricMemHeapInUse is a metric for monitoring the memory ("heap in use")
const MetricMemHeapInUse = "drt_mem_heap_inuse"

// MetricMemStackInUse is a metric for monitoring the memory ("stack in use")
const MetricMemStackInUse = "drt_mem_stack_inuse"

// MetricNetworkRecvPercent is the metric for monitoring network receive load [%]
const MetricNetworkRecvPercent = "drt_network_recv_percent"

// MetricNetworkRecvBps is the metric for monitoring network received bytes per second
const MetricNetworkRecvBps = "drt_network_recv_bps"

// MetricNetworkRecvBpsPeak is the metric for monitoring network received peak bytes per second
const MetricNetworkRecvBpsPeak = "drt_network_recv_bps_peak"

// MetricNetworkRecvBytesInCurrentEpochPerHost is the metric for monitoring network received bytes in current epoch per host
const MetricNetworkRecvBytesInCurrentEpochPerHost = "drt_network_recv_bytes_in_epoch_per_host"

// MetricNetworkSendBytesInCurrentEpochPerHost is the metric for monitoring network send bytes in current epoch per host
const MetricNetworkSendBytesInCurrentEpochPerHost = "drt_network_sent_bytes_in_epoch_per_host"

// MetricNetworkSentPercent is the metric for monitoring network sent load [%]
const MetricNetworkSentPercent = "drt_network_sent_percent"

// MetricNetworkSentBps is the metric for monitoring network sent bytes per second
const MetricNetworkSentBps = "drt_network_sent_bps"

// MetricNetworkSentBpsPeak is the metric for monitoring network sent peak bytes per second
const MetricNetworkSentBpsPeak = "drt_network_sent_bps_peak"

// MetricRoundTime is the metric for round time in seconds
const MetricRoundTime = "drt_round_time"

// MetricEpochNumber is the metric for the number of epoch
const MetricEpochNumber = "drt_epoch_number"

// MetricAppVersion is the metric for the current app version
const MetricAppVersion = "drt_app_version"

// MetricNumTxInBlock is the metric for the number of transactions in the proposed block
const MetricNumTxInBlock = "drt_num_tx_block"

// MetricConsensusState is the metric for consensus state of node proposer,participant or not consensus group
const MetricConsensusState = "drt_consensus_state"

// MetricNumMiniBlocks is the metric for number of miniblocks in a block
const MetricNumMiniBlocks = "drt_num_mini_blocks"

// MetricConsensusRoundState is the metric for consensus round state for a block
const MetricConsensusRoundState = "drt_consensus_round_state"

// MetricCrossCheckBlockHeight is the metric that store cross block height
const MetricCrossCheckBlockHeight = "drt_cross_check_block_height"

// MetricCrossCheckBlockHeightMeta is the metric that store metachain cross block height
const MetricCrossCheckBlockHeightMeta = "drt_cross_check_block_height_meta"

// MetricNumProcessedTxs is the metric that stores the number of transactions processed
const MetricNumProcessedTxs = "drt_num_transactions_processed"

// MetricCurrentBlockHash is the metric that stores the current block hash
const MetricCurrentBlockHash = "drt_current_block_hash"

// MetricCurrentRoundTimestamp is the metric that stores current round timestamp
const MetricCurrentRoundTimestamp = "drt_current_round_timestamp"

// MetricHeaderSize is the metric that stores the current block size
const MetricHeaderSize = "drt_current_block_size"

// MetricMiniBlocksSize is the metric that stores the current block size
const MetricMiniBlocksSize = "drt_mini_blocks_size"

// MetricNumShardHeadersFromPool is the metric that stores number of shard header from pool
const MetricNumShardHeadersFromPool = "drt_num_shard_headers_from_pool"

// MetricNumShardHeadersProcessed is the metric that stores number of shard header processed
const MetricNumShardHeadersProcessed = "drt_num_shard_headers_processed"

// MetricNumTimesInForkChoice is the metric that counts how many times a node was in fork choice
const MetricNumTimesInForkChoice = "drt_fork_choice_count"

// MetricHighestFinalBlock is the metric for the nonce of the highest final block
const MetricHighestFinalBlock = "drt_highest_final_nonce"

// MetricLatestTagSoftwareVersion is the metric that stores the latest tag software version
const MetricLatestTagSoftwareVersion = "drt_latest_tag_software_version"

// MetricCountConsensusAcceptedBlocks is the metric for monitoring number of blocks accepted when the node was in consensus group
const MetricCountConsensusAcceptedBlocks = "drt_count_consensus_accepted_blocks"

// MetricNodeDisplayName is the metric that stores the name of the node
const MetricNodeDisplayName = "drt_node_display_name"

// MetricConsensusGroupSize is the metric for consensus group size for the current shard/meta
const MetricConsensusGroupSize = "drt_consensus_group_size"

// MetricShardConsensusGroupSize is the metric for the shard consensus group size
const MetricShardConsensusGroupSize = "drt_shard_consensus_group_size"

// MetricMetaConsensusGroupSize is the metric for the metachain consensus group size
const MetricMetaConsensusGroupSize = "drt_meta_consensus_group_size"

// MetricNumNodesPerShard is the metric which holds the number of nodes in a shard
const MetricNumNodesPerShard = "drt_num_nodes_in_shard"

// MetricNumMetachainNodes is the metric which holds the number of nodes in metachain
const MetricNumMetachainNodes = "drt_num_metachain_nodes"

// MetricNumValidators is the metric for the number of validators
const MetricNumValidators = "drt_num_validators"

// MetricPeerType is the metric which tells the peer's type (in eligible list, in waiting list, or observer)
const MetricPeerType = "drt_peer_type"

// MetricPeerSubType is the metric which tells the peer's subtype (regular observer or full history observer)
const MetricPeerSubType = "drt_peer_subtype"

// MetricLeaderPercentage is the metric for leader rewards percentage
const MetricLeaderPercentage = "drt_leader_percentage"

// MetricDenomination is the metric for exposing the denomination
const MetricDenomination = "drt_denomination"

// MetricRoundAtEpochStart is the metric for storing the first round of the current epoch
const MetricRoundAtEpochStart = "drt_round_at_epoch_start"

// MetricNonceAtEpochStart is the metric for storing the first nonce of the current epoch
const MetricNonceAtEpochStart = "drt_nonce_at_epoch_start"

// MetricRoundsPerEpoch is the metric that tells the number of rounds in an epoch
const MetricRoundsPerEpoch = "drt_rounds_per_epoch"

// MetricRoundsPassedInCurrentEpoch is the metric that tells the number of rounds passed in current epoch
const MetricRoundsPassedInCurrentEpoch = "drt_rounds_passed_in_current_epoch"

// MetricNoncesPassedInCurrentEpoch is the metric that tells the number of nonces passed in current epoch
const MetricNoncesPassedInCurrentEpoch = "drt_nonces_passed_in_current_epoch"

// MetricReceivedProposedBlock is the metric that specifies the moment in the round when the received block has reached the
// current node. The value is provided in percent (0 meaning it has been received just after the round started and
// 100 meaning that the block has been received in the last moment of the round)
const MetricReceivedProposedBlock = "drt_consensus_received_proposed_block"

// MetricCreatedProposedBlock is the metric that specifies the percent of the block subround used for header and body
// creation (0 meaning that the block was created in no-time and 100 meaning that the block creation used all the
// subround spare duration)
const MetricCreatedProposedBlock = "drt_consensus_created_proposed_block"

// MetricRedundancyLevel is the metric that specifies the redundancy level of the current node
const MetricRedundancyLevel = "drt_redundancy_level"

// MetricRedundancyIsMainActive is the metric that specifies data about the redundancy main machine
const MetricRedundancyIsMainActive = "drt_redundancy_is_main_active"

// MetricRedundancyStepInReason is the metric that specifies why the back-up machine stepped in
const MetricRedundancyStepInReason = "drt_redundancy_step_in_reason"

// MetricValueNA represents the value to be used when a metric is not available/applicable
const MetricValueNA = "N/A"

// MetricProcessedProposedBlock is the metric that specify the percent of the block subround used for header and body
// processing (0 meaning that the block was processed in no-time and 100 meaning that the block processing used all the
// subround spare duration)
const MetricProcessedProposedBlock = "drt_consensus_processed_proposed_block"

// MetricMinGasPrice is the metric that specifies min gas price
const MetricMinGasPrice = "drt_min_gas_price"

// MetricMinGasLimit is the metric that specifies the minimum gas limit
const MetricMinGasLimit = "drt_min_gas_limit"

// MetricExtraGasLimitGuardedTx specifies the extra gas limit required for guarded transactions
const MetricExtraGasLimitGuardedTx = "drt_extra_gas_limit_guarded_tx"

// MetricExtraGasLimitRelayedTx specifies the extra gas limit required for relayed v3 transactions
const MetricExtraGasLimitRelayedTx = "drt_extra_gas_limit_relayed_tx"

// MetricRewardsTopUpGradientPoint is the metric that specifies the rewards top up gradient point
const MetricRewardsTopUpGradientPoint = "drt_rewards_top_up_gradient_point"

// MetricGasPriceModifier is the metric that specifies the gas price modifier
const MetricGasPriceModifier = "drt_gas_price_modifier"

// MetricTopUpFactor is the metric that specifies the top-up factor
const MetricTopUpFactor = "drt_top_up_factor"

// MetricMinTransactionVersion is the metric that specifies the minimum transaction version
const MetricMinTransactionVersion = "drt_min_transaction_version"

// MetricGatewayMetricsEndpoint is the metric that specifies gateway endpoint
const MetricGatewayMetricsEndpoint = "drt_gateway_metrics_endpoint"

// MetricGasPerDataByte is the metric that specifies the required gas for a data byte
const MetricGasPerDataByte = "drt_gas_per_data_byte"

// MetricMaxGasPerTransaction is the metric that specifies the maximum gas limit for a transaction
const MetricMaxGasPerTransaction = "drt_max_gas_per_transaction"

// MetricChainId is the metric that specifies current chain id
const MetricChainId = "drt_chain_id"

// MetricStartTime is the metric that specifies the genesis start time
const MetricStartTime = "drt_start_time"

// MetricRoundDuration is the metric that specifies the round duration in milliseconds
const MetricRoundDuration = "drt_round_duration"

// MetricTotalSupply holds the total supply value for the last epoch
const MetricTotalSupply = "drt_total_supply"

// MetricTotalBaseStakedValue holds the total base staked value
const MetricTotalBaseStakedValue = "drt_total_base_staked_value"

// MetricTopUpValue holds the total top up value
const MetricTopUpValue = "drt_total_top_up_value"

// MetricInflation holds the inflation value for the last epoch
const MetricInflation = "drt_inflation"

// MetricDevRewardsInEpoch holds the developers' rewards value for the last epoch
const MetricDevRewardsInEpoch = "drt_dev_rewards"

// MetricTotalFees holds the total fees value for the last epoch
const MetricTotalFees = "drt_total_fees"

// MetricEpochForEconomicsData holds the epoch for which economics data are computed
const MetricEpochForEconomicsData = "drt_epoch_for_economics_data"

// MetachainShardId will be used to identify a shard ID as metachain
const MetachainShardId = uint32(0xFFFFFFFF)

// BaseOperationCost represents the field name for base operation costs
const BaseOperationCost = "BaseOperationCost"

// BuiltInCost represents the field name for built-in operation costs
const BuiltInCost = "BuiltInCost"

// MetaChainSystemSCsCost represents the field name for metachain system smart contract operation costs
const MetaChainSystemSCsCost = "MetaChainSystemSCsCost"

// BaseOpsAPICost represents the field name of the SC API (EEI) gas costs
const BaseOpsAPICost = "BaseOpsAPICost"

// MaxPerTransaction represents the field name of max counts per transaction in block chain hook
const MaxPerTransaction = "MaxPerTransaction"

// AsyncCallStepField is the field name for the gas cost for any of the two steps required to execute an async call
const AsyncCallStepField = "AsyncCallStep"

// AsyncCallbackGasLockField is the field name for the gas amount to be locked
// before executing the destination async call, to be put aside for the async callback
const AsyncCallbackGasLockField = "AsyncCallbackGasLock"

const (
	// MetricScDeployEnableEpoch represents the epoch when the deployment of smart contracts is enabled
	MetricScDeployEnableEpoch = "drt_smart_contract_deploy_enable_epoch"

	// MetricBuiltInFunctionsEnableEpoch represents the epoch when the built-in functions is enabled
	MetricBuiltInFunctionsEnableEpoch = "drt_built_in_functions_enable_epoch"

	// MetricRelayedTransactionsEnableEpoch represents the epoch when the relayed transactions is enabled
	MetricRelayedTransactionsEnableEpoch = "drt_relayed_transactions_enable_epoch"

	// MetricPenalizedTooMuchGasEnableEpoch represents the epoch when the penalization for using too much gas is enabled
	MetricPenalizedTooMuchGasEnableEpoch = "drt_penalized_too_much_gas_enable_epoch"

	// MetricSwitchJailWaitingEnableEpoch represents the epoch when the system smart contract processing at end of epoch is enabled
	MetricSwitchJailWaitingEnableEpoch = "drt_switch_jail_waiting_enable_epoch"

	// MetricSwitchHysteresisForMinNodesEnableEpoch represents the epoch when the system smart contract changes its config to consider
	// also (minimum) hysteresis nodes for the minimum number of nodes
	MetricSwitchHysteresisForMinNodesEnableEpoch = "drt_switch_hysteresis_for_min_nodes_enable_epoch"

	// MetricBelowSignedThresholdEnableEpoch represents the epoch when the change for computing rating for validators below signed rating is enabled
	MetricBelowSignedThresholdEnableEpoch = "drt_below_signed_threshold_enable_epoch"

	// MetricTransactionSignedWithTxHashEnableEpoch represents the epoch when the node will also accept transactions that are
	// signed with the hash of transaction
	MetricTransactionSignedWithTxHashEnableEpoch = "drt_transaction_signed_with_txhash_enable_epoch"

	// MetricMetaProtectionEnableEpoch represents the epoch when the transactions to the metachain are checked to have enough gas
	MetricMetaProtectionEnableEpoch = "drt_meta_protection_enable_epoch"

	// MetricAheadOfTimeGasUsageEnableEpoch represents the epoch when the cost of smart contract prepare changes from compiler
	// per byte to ahead of time prepare per byte
	MetricAheadOfTimeGasUsageEnableEpoch = "drt_ahead_of_time_gas_usage_enable_epoch"

	// MetricGasPriceModifierEnableEpoch represents the epoch when the gas price modifier in fee computation is enabled
	MetricGasPriceModifierEnableEpoch = "drt_gas_price_modifier_enable_epoch"

	// MetricRepairCallbackEnableEpoch represents the epoch when the callback repair is activated for smart contract results
	MetricRepairCallbackEnableEpoch = "drt_repair_callback_enable_epoch"

	// MetricBlockGasAndFreeRecheckEnableEpoch represents the epoch when gas and fees used in each created or processed block are re-checked
	MetricBlockGasAndFreeRecheckEnableEpoch = "drt_block_gas_and_fee_recheck_enable_epoch"

	// MetricStakingV2EnableEpoch represents the epoch when staking v2 is enabled
	MetricStakingV2EnableEpoch = "drt_staking_v2_enable_epoch"

	// MetricStakeEnableEpoch represents the epoch when staking is enabled
	MetricStakeEnableEpoch = "drt_stake_enable_epoch"

	// MetricDoubleKeyProtectionEnableEpoch represents the epoch when double key protection is enabled
	MetricDoubleKeyProtectionEnableEpoch = "drt_double_key_protection_enable_epoch"

	// MetricDcdtEnableEpoch represents the epoch when DCDT is enabled
	MetricDcdtEnableEpoch = "drt_dcdt_enable_epoch"

	// MetricGovernanceEnableEpoch  represents the epoch when governance is enabled
	MetricGovernanceEnableEpoch = "drt_governance_enable_epoch"

	// MetricDelegationManagerEnableEpoch represents the epoch when the delegation manager is enabled
	MetricDelegationManagerEnableEpoch = "drt_delegation_manager_enable_epoch"

	// MetricDelegationSmartContractEnableEpoch represents the epoch when delegation smart contract is enabled
	MetricDelegationSmartContractEnableEpoch = "drt_delegation_smart_contract_enable_epoch"

	// MetricCorrectLastUnjailedEnableEpoch represents the epoch when the correction on the last unjailed node is applied
	MetricCorrectLastUnjailedEnableEpoch = "drt_correct_last_unjailed_enable_epoch"

	// MetricBalanceWaitingListsEnableEpoch represents the epoch when the balance waiting lists on shards fix is applied
	MetricBalanceWaitingListsEnableEpoch = "drt_balance_waiting_lists_enable_epoch"

	// MetricReturnDataToLastTransferEnableEpoch represents the epoch when the return data to last transfer is applied
	MetricReturnDataToLastTransferEnableEpoch = "drt_return_data_to_last_transfer_enable_epoch"

	// MetricSenderInOutTransferEnableEpoch represents the epoch when the sender in out transfer is applied
	MetricSenderInOutTransferEnableEpoch = "drt_sender_in_out_transfer_enable_epoch"

	// MetricRelayedTransactionsV2EnableEpoch represents the epoch when the relayed transactions v2 is enabled
	MetricRelayedTransactionsV2EnableEpoch = "drt_relayed_transactions_v2_enable_epoch"

	// MetricFixRelayedBaseCostEnableEpoch represents the epoch when the fix for relayed base cost is enabled
	MetricFixRelayedBaseCostEnableEpoch = "drt_fix_relayed_base_cost_enable_epoch"

	// MetricUnbondTokensV2EnableEpoch represents the epoch when the unbond tokens v2 is applied
	MetricUnbondTokensV2EnableEpoch = "drt_unbond_tokens_v2_enable_epoch"

	// MetricSaveJailedAlwaysEnableEpoch represents the epoch the save jailed fix is applied
	MetricSaveJailedAlwaysEnableEpoch = "drt_save_jailed_always_enable_epoch"

	// MetricValidatorToDelegationEnableEpoch represents the epoch when the validator to delegation feature (staking v3.5) is enabled
	MetricValidatorToDelegationEnableEpoch = "drt_validator_to_delegation_enable_epoch"

	// MetricReDelegateBelowMinCheckEnableEpoch represents the epoch when the re-delegation below minimum value fix is applied
	MetricReDelegateBelowMinCheckEnableEpoch = "drt_redelegate_below_min_check_enable_epoch"

	// MetricIncrementSCRNonceInMultiTransferEnableEpoch represents the epoch when the fix for multi transfer SCR is enabled
	MetricIncrementSCRNonceInMultiTransferEnableEpoch = "drt_increment_scr_nonce_in_multi_transfer_enable_epoch"

	// MetricScheduledMiniBlocksEnableEpoch represents the epoch when the scheduled miniblocks feature is enabled
	MetricScheduledMiniBlocksEnableEpoch = "drt_scheduled_miniblocks_enable_epoch"

	// MetricDCDTMultiTransferEnableEpoch represents the epoch when the DCDT multi transfer feature is enabled
	MetricDCDTMultiTransferEnableEpoch = "drt_dcdt_multi_transfer_enable_epoch"

	// MetricGlobalMintBurnDisableEpoch represents the epoch when the global mint and burn feature is disabled
	MetricGlobalMintBurnDisableEpoch = "drt_global_mint_burn_disable_epoch"

	// MetricDCDTTransferRoleEnableEpoch represents the epoch when the DCDT transfer role feature is enabled
	MetricDCDTTransferRoleEnableEpoch = "drt_dcdt_transfer_role_enable_epoch"

	// MetricComputeRewardCheckpointEnableEpoch represents the epoch when compute reward checkpoint feature is enabled
	MetricComputeRewardCheckpointEnableEpoch = "drt_compute_reward_checkpoint_enable_epoch"

	// MetricSCRSizeInvariantCheckEnableEpoch represents the epoch when scr size invariant check is enabled
	MetricSCRSizeInvariantCheckEnableEpoch = "drt_scr_size_invariant_check_enable_epoch"

	// MetricBackwardCompSaveKeyValueEnableEpoch represents the epoch when backward compatibility save key valu is enabled
	MetricBackwardCompSaveKeyValueEnableEpoch = "drt_backward_comp_save_keyvalue_enable_epoch"

	// MetricDCDTNFTCreateOnMultiShardEnableEpoch represents the epoch when dcdt nft create on multi shard is enabled
	MetricDCDTNFTCreateOnMultiShardEnableEpoch = "drt_dcdt_nft_create_on_multi_shard_enable_epoch"

	// MetricMetaDCDTSetEnableEpoch represents the epoch when meta dcdt set is enabled
	MetricMetaDCDTSetEnableEpoch = "drt_meta_dcdt_set_enable_epoch"

	// MetricAddTokensToDelegationEnableEpoch represents the epoch when add tokens to delegation
	MetricAddTokensToDelegationEnableEpoch = "drt_add_tokens_to_delegation_enable_epoch"

	// MetricMultiDCDTTransferFixOnCallBackOnEnableEpoch represents the epoch when multi dcdt transfer fix on callback on is enabled
	MetricMultiDCDTTransferFixOnCallBackOnEnableEpoch = "drt_multi_dcdt_transfer_fix_on_callback_enable_epoch"

	// MetricOptimizeGasUsedInCrossMiniBlocksEnableEpoch represents the epoch when optimize gas used in cross miniblocks is enabled
	MetricOptimizeGasUsedInCrossMiniBlocksEnableEpoch = "drt_optimize_gas_used_in_cross_miniblocks_enable_epoch"
	// MetricCorrectFirstQueuedEpoch represents the epoch when correct first queued fix is enabled
	MetricCorrectFirstQueuedEpoch = "drt_correct_first_queued_enable_epoch"

	// MetricCorrectJailedNotUnstakedEmptyQueueEpoch represents the epoch when correct jailed not unstacked meptry queue fix is enabled
	MetricCorrectJailedNotUnstakedEmptyQueueEpoch = "drt_correct_jailed_not_unstaked_empty_queue_enable_epoch"

	// MetricFixOOGReturnCodeEnableEpoch represents the epoch when OOG return code fix is enabled
	MetricFixOOGReturnCodeEnableEpoch = "drt_fix_oog_return_code_enable_epoch"

	// MetricRemoveNonUpdatedStorageEnableEpoch represents the epoch when remove non updated storage fix is enabled
	MetricRemoveNonUpdatedStorageEnableEpoch = "drt_remove_non_updated_storage_enable_epoch"

	// MetricDeleteDelegatorAfterClaimRewardsEnableEpoch represents the epoch when delete delegator after claim rewards fix is enabled
	MetricDeleteDelegatorAfterClaimRewardsEnableEpoch = "drt_delete_delegator_after_claim_rewards_enable_epoch"

	// MetricOptimizeNFTStoreEnableEpoch represents the epoch when optimize nft store feature is enabled
	MetricOptimizeNFTStoreEnableEpoch = "drt_optimize_nft_store_enable_epoch"

	// MetricCreateNFTThroughExecByCallerEnableEpoch represents the epoch when create nft through exec by caller functionality is enabled
	MetricCreateNFTThroughExecByCallerEnableEpoch = "drt_create_nft_through_exec_by_caller_enable_epoch"

	// MetricStopDecreasingValidatorRatingWhenStuckEnableEpoch represents the epoch when stop decreaing validator rating when stuck functionality is enabled
	MetricStopDecreasingValidatorRatingWhenStuckEnableEpoch = "drt_stop_decreasing_validator_rating_when_stuck_enable_epoch"

	// MetricFrontRunningProtectionEnableEpoch represents the epoch when front running protection feature is enabled
	MetricFrontRunningProtectionEnableEpoch = "drt_front_running_protection_enable_epoch"

	// MetricIsPayableBySCEnableEpoch represents the epoch when is payable by SC feature is enabled
	MetricIsPayableBySCEnableEpoch = "drt_is_payable_by_sc_enable_epoch"

	// MetricCleanUpInformativeSCRsEnableEpoch represents the epoch when cleanup informative scrs functionality is enabled
	MetricCleanUpInformativeSCRsEnableEpoch = "drt_cleanup_informative_scrs_enable_epoch"

	// MetricStorageAPICostOptimizationEnableEpoch represents the epoch when storage api cost optimization feature is enabled
	MetricStorageAPICostOptimizationEnableEpoch = "drt_storage_api_cost_optimization_enable_epoch"

	// MetricTransformToMultiShardCreateEnableEpoch represents the epoch when transform to multi shard create functionality is enabled
	MetricTransformToMultiShardCreateEnableEpoch = "drt_transform_to_multi_shard_create_enable_epoch"

	// MetricDCDTRegisterAndSetAllRolesEnableEpoch represents the epoch when dcdt register and set all roles functionality is enabled
	MetricDCDTRegisterAndSetAllRolesEnableEpoch = "drt_dcdt_register_and_set_all_roles_enable_epoch"

	// MetricDoNotReturnOldBlockInBlockchainHookEnableEpoch represents the epoch when do not return old block in blockchain hook fix is enabled
	MetricDoNotReturnOldBlockInBlockchainHookEnableEpoch = "drt_do_not_returns_old_block_in_blockchain_hook_enable_epoch"

	// MetricAddFailedRelayedTxToInvalidMBsDisableEpoch represents the epoch when add failed relayed tx to invalid miniblocks functionality is enabled
	MetricAddFailedRelayedTxToInvalidMBsDisableEpoch = "drt_add_failed_relayed_tx_to_invalid_mbs_enable_epoch"

	// MetricSCRSizeInvariantOnBuiltInResultEnableEpoch represents the epoch when scr size invariant on builtin result functionality is enabled
	MetricSCRSizeInvariantOnBuiltInResultEnableEpoch = "drt_scr_size_invariant_on_builtin_result_enable_epoch"

	// MetricCheckCorrectTokenIDForTransferRoleEnableEpoch represents the epoch when check correct tokenID for transfer role fix is enabled
	MetricCheckCorrectTokenIDForTransferRoleEnableEpoch = "drt_check_correct_tokenid_for_transfer_role_enable_epoch"

	// MetricDisableExecByCallerEnableEpoch represents the epoch when disable exec by caller functionality is enabled
	MetricDisableExecByCallerEnableEpoch = "drt_disable_exec_by_caller_enable_epoch"

	// MetricFailExecutionOnEveryAPIErrorEnableEpoch represents the epoch when fail execution on every api error functionality is enabled
	MetricFailExecutionOnEveryAPIErrorEnableEpoch = "drt_fail_execution_on_every_api_error_enable_epoch"

	// MetricManagedCryptoAPIsEnableEpoch represents the epoch when managed cypto apis functionality is enabled
	MetricManagedCryptoAPIsEnableEpoch = "drt_managed_crypto_apis_enable_epoch"

	// MetricRefactorContextEnableEpoch represents the epoch when refactor context functionality is enabled
	MetricRefactorContextEnableEpoch = "drt_refactor_context_enable_epoch"

	// MetricCheckFunctionArgumentEnableEpoch represents the epoch when check function argument functionality is enabled
	MetricCheckFunctionArgumentEnableEpoch = "drt_check_function_argument_enable_epoch"

	// MetricCheckExecuteOnReadOnlyEnableEpoch represents the epoch when check execute on read only fix is enabled
	MetricCheckExecuteOnReadOnlyEnableEpoch = "drt_check_execute_on_readonly_enable_epoch"

	// MetricMiniBlockPartialExecutionEnableEpoch represents the epoch when miniblock partial execution feature is enabled
	MetricMiniBlockPartialExecutionEnableEpoch = "drt_miniblock_partial_execution_enable_epoch"

	// MetricDCDTMetadataContinuousCleanupEnableEpoch represents the epoch when dcdt metadata continuous clenaup functionality is enabled
	MetricDCDTMetadataContinuousCleanupEnableEpoch = "drt_dcdt_metadata_continuous_cleanup_enable_epoch"

	// MetricFixAsyncCallBackArgsListEnableEpoch represents the epoch when fix async callback args list is enabled
	MetricFixAsyncCallBackArgsListEnableEpoch = "drt_fix_async_callback_args_list_enable_epoch"

	// MetricFixOldTokenLiquidityEnableEpoch represents the epoch when fix old token liquidity is enabled
	MetricFixOldTokenLiquidityEnableEpoch = "drt_fix_old_token_liquidity_enable_epoch"

	// MetricRuntimeMemStoreLimitEnableEpoch represents the epoch when runtime mem store limit functionality is enabled
	MetricRuntimeMemStoreLimitEnableEpoch = "drt_runtime_mem_store_limit_enable_epoch"

	// MetricRuntimeCodeSizeFixEnableEpoch represents the epoch when runtime code size fix is enabled
	MetricRuntimeCodeSizeFixEnableEpoch = "drt_runtime_code_size_fix_enable_epoch"

	// MetricSetSenderInEeiOutputTransferEnableEpoch represents the epoch when set sender in eei output transfer functionality is enabled
	MetricSetSenderInEeiOutputTransferEnableEpoch = "drt_set_sender_in_eei_output_transfer_enable_epoch"

	// MetricRefactorPeersMiniBlocksEnableEpoch represents the epoch when refactor peers miniblock feature is enabled
	MetricRefactorPeersMiniBlocksEnableEpoch = "drt_refactor_peers_miniblocks_enable_epoch"

	// MetricSCProcessorV2EnableEpoch represents the epoch when SC processor V2 feature is enabled
	MetricSCProcessorV2EnableEpoch = "drt_sc_processorv2_enable_epoch"

	// MetricMaxBlockchainHookCountersEnableEpoch represents the epoch when max blockchain hook counters functionality is enabled
	MetricMaxBlockchainHookCountersEnableEpoch = "drt_max_blockchain_hook_counters_enable_epoch"

	// MetricWipeSingleNFTLiquidityDecreaseEnableEpoch represents the epoch when wipe single NFT liquidity decrease functionality is enabled
	MetricWipeSingleNFTLiquidityDecreaseEnableEpoch = "drt_wipe_single_nft_liquidity_decrease_enable_epoch"

	// MetricAlwaysSaveTokenMetaDataEnableEpoch represents the epoch when always save token metadata functionality is enabled
	MetricAlwaysSaveTokenMetaDataEnableEpoch = "drt_always_save_token_metadata_enable_epoch"

	// MetricSetGuardianEnableEpoch represents the epoch when the guardian feature is enabled
	MetricSetGuardianEnableEpoch = "drt_set_guardian_feature_enable_epoch"

	// MetricSetScToScLogEventEnableEpoch represents the epoch when the sc to sc log event feature is enabled
	MetricSetScToScLogEventEnableEpoch = "drt_set_sc_to_sc_log_event_enable_epoch"

	// MetricRelayedNonceFixEnableEpoch represents the epoch when relayed nonce fix is enabled
	MetricRelayedNonceFixEnableEpoch = "drt_relayed_nonce_fix_enable_epoch"

	// MetricDeterministicSortOnValidatorsInfoEnableEpoch represents the epoch when deterministic sort on validators info functionality is enabled
	MetricDeterministicSortOnValidatorsInfoEnableEpoch = "drt_deterministic_sort_on_validators_info_enable_epoch"

	// MetricKeepExecOrderOnCreatedSCRsEnableEpoch represents the epoch when keep exec order on created scs fix is enabled
	MetricKeepExecOrderOnCreatedSCRsEnableEpoch = "drt_keep_exec_order_on_created_scrs_enable_epoch"

	// MetricMultiClaimOnDelegationEnableEpoch represents the epoch when multi claim on delegation functionality is enabled
	MetricMultiClaimOnDelegationEnableEpoch = "drt_multi_claim_on_delegation_enable_epoch"

	// MetricChangeUsernameEnableEpoch represents the epoch when change username functionality is enabled
	MetricChangeUsernameEnableEpoch = "drt_change_username_enable_epoch"

	// MetricAutoBalanceDataTriesEnableEpoch represents the epoch when auto balance data tries feature is enabled
	MetricAutoBalanceDataTriesEnableEpoch = "drt_auto_balance_data_tries_enable_epoch"

	// MetricMigrateDataTrieEnableEpoch represents the epoch when migrate data trie feature is enabled
	MetricMigrateDataTrieEnableEpoch = "drt_migrate_datatrie_enable_epoch"

	// MetricConsistentTokensValuesLengthCheckEnableEpoch represents the epoch when consistent tokens values length check is enabled
	MetricConsistentTokensValuesLengthCheckEnableEpoch = "drt_consistent_tokens_values_length_check_enable_epoch"

	// MetricFixDelegationChangeOwnerOnAccountEnableEpoch represents the epoch when fix delegation change owner on account is enabled
	MetricFixDelegationChangeOwnerOnAccountEnableEpoch = "drt_fix_delegation_change_owner_on_account_enable_epoch"

	// MetricDynamicGasCostForDataTrieStorageLoadEnableEpoch represents the epoch when dynamic gas cost for data tries storage load functionality is enabled
	MetricDynamicGasCostForDataTrieStorageLoadEnableEpoch = "drt_dynamic_gas_cost_for_datatrie_storage_load_enable_epoch"

	// MetricNFTStopCreateEnableEpoch represents the epoch when NFT stop create functionality is enabled
	MetricNFTStopCreateEnableEpoch = "drt_nft_stop_create_enable_epoch"

	// MetricChangeOwnerAddressCrossShardThroughSCEnableEpoch represents the epoch when change owner address cross shard through SC functionality is enabled
	MetricChangeOwnerAddressCrossShardThroughSCEnableEpoch = "drt_change_owner_address_cross_shard_through_sc_enable_epoch"

	// MetricFixGasRemainingForSaveKeyValueBuiltinFunctionEnableEpoch represents the epoch when fix gas remaining for save key value builin function is enabled
	MetricFixGasRemainingForSaveKeyValueBuiltinFunctionEnableEpoch = "drt_fix_gas_remainig_for_save_keyvalue_builtin_function_enable_epoch"

	// MetricCurrentRandomnessOnSortingEnableEpoch represents the epoch when current randomness on sorting functionality is enabled
	MetricCurrentRandomnessOnSortingEnableEpoch = "drt_current_randomness_on_sorting_enable_epoch"

	// MetricStakeLimitsEnableEpoch represents the epoch when stake limits functionality is enabled
	MetricStakeLimitsEnableEpoch = "drt_stake_limits_enable_epoch"

	// MetricStakingV4Step1EnableEpoch represents the epoch when staking v4 step 1 feature is enabled
	MetricStakingV4Step1EnableEpoch = "drt_staking_v4_step1_enable_epoch"

	// MetricStakingV4Step2EnableEpoch represents the epoch when staking v4 step 2 feature is enabled
	MetricStakingV4Step2EnableEpoch = "drt_staking_v4_step2_enable_epoch"

	// MetricStakingV4Step3EnableEpoch represents the epoch when staking v4 step 3 feature is enabled
	MetricStakingV4Step3EnableEpoch = "drt_staking_v4_step3_enable_epoch"

	// MetricCleanupAuctionOnLowWaitingListEnableEpoch represents the epoch when cleanup auction on low waiting list fix is enabled
	MetricCleanupAuctionOnLowWaitingListEnableEpoch = "drt_cleanup_auction_on_low_waiting_list_enable_epoch"

	// MetricAlwaysMergeContextsInEEIEnableEpoch represents the epoch when always merge contexts in EEI fix is enabled
	MetricAlwaysMergeContextsInEEIEnableEpoch = "drt_always_merge_contexts_in_eei_enable_epoch"

	// MetricDynamicDCDTEnableEpoch represents the epoch when dynamic DCDT feature is enabled
	MetricDynamicDCDTEnableEpoch = "drt_dynamic_dcdt_enable_epoch"

	// MetricREWAInMultiTransferEnableEpoch represents the epoch when REWA in multi transfer feature is enabled
	MetricREWAInMultiTransferEnableEpoch = "drt_rewa_in_multi_transfer_enable_epoch"

	// MetricCheckBuiltInCallOnTransferValueAndFailEnableRound represents the round when check builtincall on transfer value and fail is enabled
	MetricCheckBuiltInCallOnTransferValueAndFailEnableRound = "drt_checkbuiltincall_ontransfervalueandfail_enable_round"

	// MetricMultiDCDTNFTTransferAndExecuteByUserEnableEpoch represents the epoch when enshrined sovereign opcodes are enabled
	MetricMultiDCDTNFTTransferAndExecuteByUserEnableEpoch = "drt_multi_dcdt_transfer_execute_by_user_enable_epoch"

	// MetricFixRelayedMoveBalanceToNonPayableSCEnableEpoch represents the epoch when the fix for relayed move balance to non-payable sc is enabled
	MetricFixRelayedMoveBalanceToNonPayableSCEnableEpoch = "drt_fix_relayed_move_balance_to_non_payable_sc_enable_epoch"

	// MetricRelayedTransactionsV3EnableEpoch represents the epoch when the relayed transactions v3 are enabled
	MetricRelayedTransactionsV3EnableEpoch = "drt_relayed_transactions_v3_enable_epoch"

	// MetricRelayedTransactionsV3FixDCDTTransferEnableEpoch represents the epoch when the fix for relayed transactions v3 with dcdt transfer are enabled
	MetricRelayedTransactionsV3FixDCDTTransferEnableEpoch = "drt_relayed_transactions_v3_fix_dcdt_transfer_enable_epoch"

	// MetricMaxNodesChangeEnableEpoch holds configuration for changing the maximum number of nodes and the enabling epoch
	MetricMaxNodesChangeEnableEpoch = "drt_max_nodes_change_enable_epoch"

	// MetricCryptoOpcodesV2EnableEpoch represents the epoch when crypto opcodes v2 feature is enabled
	MetricCryptoOpcodesV2EnableEpoch = "drt_crypto_opcodes_v2_enable_epoch"

	// MetricEpochEnable represents the epoch when the max nodes change configuration is applied
	MetricEpochEnable = "drt_epoch_enable"

	// EpochEnableSuffix represents the suffix for EpochEnable item in MaxNodesChangeEnableEpoch list
	EpochEnableSuffix = "_epoch_enable"

	// MetricMaxNumNodes represents the maximum number of nodes than can be enabled in a max nodes change configuration setup
	MetricMaxNumNodes = "drt_max_num_nodes"

	// MaxNumNodesSuffix represents the suffix for MaxNumNodes item in MaxNodesChangeEnableEpoch list
	MaxNumNodesSuffix = "_max_num_nodes"

	// MetricNodesToShufflePerShard represents the nodes to be shuffled per shard
	MetricNodesToShufflePerShard = "drt_nodes_to_shuffle_per_shard"

	// NodesToShufflePerShardSuffix represents the suffix for NodesToShufflePerShard item in MaxNodesChangeEnableEpoch list
	NodesToShufflePerShardSuffix = "_nodes_to_shuffle_per_shard"
)

const (
	// MetricRatingsGeneralStartRating represents the starting rating used by the rater
	MetricRatingsGeneralStartRating = "drt_ratings_general_start_rating"

	// MetricRatingsGeneralMaxRating represents the maximum rating limit
	MetricRatingsGeneralMaxRating = "drt_ratings_general_max_rating"

	// MetricRatingsGeneralMinRating represents the minimum rating limit
	MetricRatingsGeneralMinRating = "drt_ratings_general_min_rating"

	// MetricRatingsGeneralSignedBlocksThreshold represents the signed blocks threshold
	MetricRatingsGeneralSignedBlocksThreshold = "drt_ratings_general_signed_blocks_threshold"

	// MetricRatingsGeneralSelectionChances represents the selection chances thresholds
	MetricRatingsGeneralSelectionChances = "drt_ratings_general_selection_chances"

	// MetricSelectionChancesMaxThreshold represents the max threshold for a selection chances item
	MetricSelectionChancesMaxThreshold = "drt_max_threshold"

	// SelectionChancesMaxThresholdSuffix represents the SelectionChances suffix for MaxThreshold
	SelectionChancesMaxThresholdSuffix = "_max_threshold"

	// MetricSelectionChancesChancePercent represents the chance percentage for a selection chances metric
	MetricSelectionChancesChancePercent = "drt_chance_percent"

	// SelectionChancesChancePercentSuffix represents the SelectionChances suffix for ChancePercent
	SelectionChancesChancePercentSuffix = "_chance_percent"

	// MetricRatingsShardChainHoursToMaxRatingFromStartRating represents the hours to max rating from start rating
	MetricRatingsShardChainHoursToMaxRatingFromStartRating = "drt_ratings_shardchain_hours_to_max_rating_from_start_rating"

	// MetricRatingsShardChainProposerValidatorImportance represents the proposer validator importance index
	MetricRatingsShardChainProposerValidatorImportance = "drt_ratings_shardchain_proposer_validator_importance"

	// MetricRatingsShardChainProposerDecreaseFactor represents the proposer decrease factor
	MetricRatingsShardChainProposerDecreaseFactor = "drt_ratings_shardchain_proposer_decrease_factor"

	// MetricRatingsShardChainValidatorDecreaseFactor represents the validator decrease factor
	MetricRatingsShardChainValidatorDecreaseFactor = "drt_ratings_shardchain_validator_decrease_factor"

	// MetricRatingsShardChainConsecutiveMissedBlocksPenalty represents the consecutive missed block penalty
	MetricRatingsShardChainConsecutiveMissedBlocksPenalty = "drt_ratings_shardchain_consecutive_missed_blocks_penalty"

	// MetricRatingsMetaChainHoursToMaxRatingFromStartRating represents the hours to max rating from start rating
	MetricRatingsMetaChainHoursToMaxRatingFromStartRating = "drt_ratings_metachain_hours_to_max_rating_from_start_rating"

	// MetricRatingsMetaChainProposerValidatorImportance represents the proposer validator importance index
	MetricRatingsMetaChainProposerValidatorImportance = "drt_ratings_metachain_proposer_validator_importance"

	// MetricRatingsMetaChainProposerDecreaseFactor represents the proposer decrease factor
	MetricRatingsMetaChainProposerDecreaseFactor = "drt_ratings_metachain_proposer_decrease_factor"

	// MetricRatingsMetaChainValidatorDecreaseFactor represents the validator decrease factor
	MetricRatingsMetaChainValidatorDecreaseFactor = "drt_ratings_metachain_validator_decrease_factor"

	// MetricRatingsMetaChainConsecutiveMissedBlocksPenalty represents the consecutive missed blocks penalty
	MetricRatingsMetaChainConsecutiveMissedBlocksPenalty = "drt_ratings_metachain_consecutive_missed_blocks_penalty"

	// MetricRatingsPeerHonestyDecayCoefficient represents the peer honesty decay coefficient
	MetricRatingsPeerHonestyDecayCoefficient = "drt_ratings_peerhonesty_decay_coefficient"

	// MetricRatingsPeerHonestyDecayUpdateIntervalInSeconds represents the decat update interval in seconds
	MetricRatingsPeerHonestyDecayUpdateIntervalInSeconds = "drt_ratings_peerhonesty_decay_update_interval_inseconds"

	// MetricRatingsPeerHonestyMaxScore represents the peer honesty max score allowed
	MetricRatingsPeerHonestyMaxScore = "drt_ratings_peerhonesty_max_score"

	// MetricRatingsPeerHonestyMinScore represents the peer honesty min score
	MetricRatingsPeerHonestyMinScore = "drt_ratings_peerhonesty_min_score"

	// MetricRatingsPeerHonestyBadPeerThreshold represents the peer honesty bad peer threshold
	MetricRatingsPeerHonestyBadPeerThreshold = "drt_ratings_peerhonesty_bad_peer_threshold"

	// MetricRatingsPeerHonestyUnitValue represents the peer honesty unit value
	MetricRatingsPeerHonestyUnitValue = "drt_ratings_peerhonesty_unit_value"

	// MetricHysteresis represents the hysteresis threshold
	MetricHysteresis = "drt_hysteresis"

	// MetricAdaptivity represents a boolean to determine if adaptivity will be enabled or not
	MetricAdaptivity = "drt_adaptivity"
)

const (
	// StorerOrder defines the order of storers to be notified of a start of epoch event
	StorerOrder = iota
	// ChainParametersOrder defines the order in which ChainParameters is notified of a start of epoch event
	ChainParametersOrder
	// NodesCoordinatorOrder defines the order in which NodesCoordinator is notified of a start of epoch event
	NodesCoordinatorOrder
	// ConsensusHandlerOrder defines the order in which ConsensusHandler is notified of a start of epoch event
	ConsensusHandlerOrder
	// ConsensusStartRoundOrder defines the order in which Consensus StartRound subround is notified of a start of epoch event
	ConsensusStartRoundOrder
	// NetworkShardingOrder defines the order in which the network sharding subsystem is notified of a start of epoch event
	NetworkShardingOrder
	// IndexerOrder defines the order in which indexer is notified of a start of epoch event
	IndexerOrder
	// NetStatisticsOrder defines the order in which netStatistic component is notified of a start of epoch event
	NetStatisticsOrder
	// OldDatabaseCleanOrder defines the order in which oldDatabaseCleaner component is notified of a start of epoch event
	OldDatabaseCleanOrder
)

// NodeState specifies what type of state a node could have
type NodeState int

const (
	// NsSynchronized defines ID of a state of synchronized
	NsSynchronized NodeState = iota
	// NsNotSynchronized defines ID of a state of not synchronized
	NsNotSynchronized
	// NsNotCalculated defines ID of a state which is not calculated
	NsNotCalculated
)

// MetricP2PPeerInfo is the metric for the node's p2p info
const MetricP2PPeerInfo = "drt_p2p_peer_info"

// MetricP2PIntraShardValidators is the metric that outputs the intra-shard connected validators
const MetricP2PIntraShardValidators = "drt_p2p_intra_shard_validators"

// MetricP2PCrossShardValidators is the metric that outputs the cross-shard connected validators
const MetricP2PCrossShardValidators = "drt_p2p_cross_shard_validators"

// MetricP2PIntraShardObservers is the metric that outputs the intra-shard connected observers
const MetricP2PIntraShardObservers = "drt_p2p_intra_shard_observers"

// MetricP2PCrossShardObservers is the metric that outputs the cross-shard connected observers
const MetricP2PCrossShardObservers = "drt_p2p_cross_shard_observers"

// MetricP2PUnknownPeers is the metric that outputs the unknown-shard connected peers
const MetricP2PUnknownPeers = "drt_p2p_unknown_shard_peers"

// MetricP2PNumConnectedPeersClassification is the metric for monitoring the number of connected peers split on the connection type
const MetricP2PNumConnectedPeersClassification = "drt_p2p_num_connected_peers_classification"

// MetricAreVMQueriesReady will hold the string representation of the boolean that indicated if the node is ready
// to process VM queries
const MetricAreVMQueriesReady = "drt_are_vm_queries_ready"

// HighestRoundFromBootStorage is the key for the highest round that is saved in storage
const HighestRoundFromBootStorage = "highestRoundFromBootStorage"

// TriggerRegistryKeyPrefix is the key prefix to save epoch start registry to storage
const TriggerRegistryKeyPrefix = "epochStartTrigger_"

// TriggerRegistryInitialKeyPrefix is the key prefix to save initial data to storage
const TriggerRegistryInitialKeyPrefix = "initial_value_epoch_"

// NodesCoordinatorRegistryKeyPrefix is the key prefix to save epoch start registry to storage
const NodesCoordinatorRegistryKeyPrefix = "indexHashed_"

// ShuffledOut signals that a restart is pending because the node was shuffled out
const ShuffledOut = "shuffledOut"

// WrongConfiguration signals that the node has a malformed configuration and cannot continue processing
const WrongConfiguration = "wrongConfiguration"

// ImportComplete signals that a node restart will be done because the import did complete
const ImportComplete = "importComplete"

// DefaultStatsPath is the default path where the node stats are logged
const DefaultStatsPath = "stats"

// DefaultDBPath is the default path for nodes databases
const DefaultDBPath = "db"

// MetachainShardName is the string identifier of the metachain shard
const MetachainShardName = "metachain"

// TemporaryPath is the default temporary path directory
const TemporaryPath = "temp"

// TimeToWaitForP2PBootstrap is the wait time for the P2P to bootstrap
const TimeToWaitForP2PBootstrap = 20 * time.Second

// MaxSoftwareVersionLengthInBytes represents the maximum length for the software version to be saved in block header
const MaxSoftwareVersionLengthInBytes = 10

// ExtraDelayForBroadcastBlockInfo represents the number of seconds to wait since a block has been broadcast and the
// moment when its components, like mini blocks and transactions, would be broadcast too
const ExtraDelayForBroadcastBlockInfo = 1 * time.Second

// ExtraDelayBetweenBroadcastMbsAndTxs represents the number of seconds to wait since miniblocks have been broadcast
// and the moment when theirs transactions would be broadcast too
const ExtraDelayBetweenBroadcastMbsAndTxs = 1 * time.Second

// ExtraDelayForRequestBlockInfo represents the number of seconds to wait since a block has been received and the
// moment when its components, like mini blocks and transactions, would be requested too if they are still missing
const ExtraDelayForRequestBlockInfo = ExtraDelayForBroadcastBlockInfo + ExtraDelayBetweenBroadcastMbsAndTxs + time.Second

// CommitMaxTime represents max time accepted for a commit action, after which a warn message is displayed
const CommitMaxTime = 3 * time.Second

// PutInStorerMaxTime represents max time accepted for a put action, after which a warn message is displayed
const PutInStorerMaxTime = time.Second

// DefaultUnstakedEpoch represents the default epoch that is set for a validator that has not unstaked yet
const DefaultUnstakedEpoch = math.MaxUint32

// InvalidMessageBlacklistDuration represents the time to keep a peer in the black list if it sends a message that
// does not follow the protocol: example not using the same marshaler as the other peers
const InvalidMessageBlacklistDuration = time.Second * 3600

// PublicKeyBlacklistDuration represents the time to keep a public key in the black list if it will degrade its
// rating to a minimum threshold due to improper messages
const PublicKeyBlacklistDuration = time.Second * 7200

// InvalidSigningBlacklistDuration defines the time to keep a peer id in blacklist if it signs a message with invalid signature
const InvalidSigningBlacklistDuration = time.Second * 7200

// MaxWaitingTimeToReceiveRequestedItem represents the maximum waiting time in seconds needed to receive the requested items
const MaxWaitingTimeToReceiveRequestedItem = 5 * time.Second

// DefaultLogProfileIdentifier represents the default log profile used when the logviewer/termui applications do not
// need to change the current logging profile
const DefaultLogProfileIdentifier = "[default log profile]"

// NotSetDestinationShardID represents the shardIdString when the destinationShardId is not set in the prefs
const NotSetDestinationShardID = "disabled"

// AdditionalScrForEachScCallOrSpecialTx specifies the additional number of smart contract results which should be
// considered by a node, when it includes sc calls or special txs in a miniblock.
// Ex.: normal txs -> aprox. 27000, sc calls or special txs -> aprox. 6250 = 27000 / (AdditionalScrForEachScCallOrSpecialTx + 1),
// considering that constant below is set to 3
const AdditionalScrForEachScCallOrSpecialTx = 3

// MaxRoundsWithoutCommittedStartInEpochBlock defines the maximum rounds to wait for start in epoch block to be committed,
// before a special action to be applied
const MaxRoundsWithoutCommittedStartInEpochBlock = 50

// DefaultResolversIdentifier represents the identifier that is used in conjunction with regular resolvers
// (that makes the node run properly)
const DefaultResolversIdentifier = "default resolver"

// DefaultInterceptorsIdentifier represents the identifier that is used in conjunction with regular interceptors
// (that makes the node run properly)
const DefaultInterceptorsIdentifier = "default interceptor"

// HardforkInterceptorsIdentifier represents the identifier that is used in the hardfork process
const HardforkInterceptorsIdentifier = "hardfork interceptor"

// HardforkResolversIdentifier represents the resolver that is used in the hardfork process
const HardforkResolversIdentifier = "hardfork resolver"

// EpochStartInterceptorsIdentifier represents the identifier that is used in the start-in-epoch process
const EpochStartInterceptorsIdentifier = "epoch start interceptor"

// TimeoutGettingTrieNodes defines the timeout in trie sync operation if no node is received
const TimeoutGettingTrieNodes = 2 * time.Minute // to consider syncing a very large trie node of 64MB at ~1MB/s

// TimeoutGettingTrieNodesInHardfork represents the maximum time allowed between 2 nodes fetches (and commits)
// during the hardfork process
const TimeoutGettingTrieNodesInHardfork = time.Minute * 10

// RetrialIntervalForOutportDriver is the interval in which the outport driver should try to call the driver again
const RetrialIntervalForOutportDriver = time.Second * 10

// NodeProcessingMode represents the processing mode in which the node was started
type NodeProcessingMode int

const (
	// Normal means that the node has started in the normal processing mode
	Normal NodeProcessingMode = iota

	// ImportDb means that the node has started in the import-db mode
	ImportDb
)

const (
	// ActiveDBKey is the key at which ActiveDBVal will be saved
	ActiveDBKey = "activeDB"

	// ActiveDBVal is the value that will be saved at ActiveDBKey
	ActiveDBVal = "yes"

	// TrieSyncedKey is the key at which TrieSyncedVal will be saved
	TrieSyncedKey = "synced"

	// TrieSyncedVal is the value that will be saved at TrieSyncedKey
	TrieSyncedVal = "yes"

	// TrieLeavesChannelDefaultCapacity represents the default value to be used as capacity for getting all trie leaves on
	// a channel
	TrieLeavesChannelDefaultCapacity = 100

	// TrieLeavesChannelSyncCapacity represents the value to be used as capacity for getting main trie
	// leaf nodes for trie sync
	TrieLeavesChannelSyncCapacity = 1000
)

// ApiOutputFormat represents the format type returned by api
type ApiOutputFormat uint8

const (
	// ApiOutputFormatJSON outport format returns struct directly, will be serialized into JSON by gin
	ApiOutputFormatJSON ApiOutputFormat = 0

	// ApiOutputFormatProto outport format returns the bytes of the proto object
	ApiOutputFormatProto ApiOutputFormat = 1
)

// BlockProcessingCutoffMode represents the type to be used to identify the mode of the block processing cutoff
type BlockProcessingCutoffMode string

const (
	// BlockProcessingCutoffModePause represents the mode where the node will pause the processing at the given coordinates
	BlockProcessingCutoffModePause = "pause"
	// BlockProcessingCutoffModeProcessError represents the mode where the node will reprocess with error the block at the given coordinates
	BlockProcessingCutoffModeProcessError = "process-error"
)

// BlockProcessingCutoffTrigger represents the trigger of the cutoff potentially used in block processing
type BlockProcessingCutoffTrigger string

const (
	// BlockProcessingCutoffByNonce represents the cutoff by nonce
	BlockProcessingCutoffByNonce BlockProcessingCutoffTrigger = "nonce"
	// BlockProcessingCutoffByRound represents the cutoff by round
	BlockProcessingCutoffByRound BlockProcessingCutoffTrigger = "round"
	// BlockProcessingCutoffByEpoch represents the cutoff by epoch
	BlockProcessingCutoffByEpoch BlockProcessingCutoffTrigger = "epoch"
)

// MaxIndexOfTxInMiniBlock defines the maximum index of a tx inside one mini block
const MaxIndexOfTxInMiniBlock = int32(29999)

// MetricAccountsSnapshotInProgress is the metric that outputs the status of the accounts' snapshot, if it's in progress or not
const MetricAccountsSnapshotInProgress = "drt_accounts_snapshot_in_progress"

// MetricLastAccountsSnapshotDurationSec is the metric that outputs the duration in seconds of the last accounts db snapshot. If snapshot is in progress it will be set to 0
const MetricLastAccountsSnapshotDurationSec = "drt_accounts_snapshot_last_duration_in_seconds"

// MetricPeersSnapshotInProgress is the metric that outputs the status of the peers' snapshot, if it's in progress or not
const MetricPeersSnapshotInProgress = "drt_peers_snapshot_in_progress"

// MetricLastPeersSnapshotDurationSec is the metric that outputs the duration in seconds of the last peers db snapshot. If snapshot is in progress it will be set to 0
const MetricLastPeersSnapshotDurationSec = "drt_peers_snapshot_last_duration_in_seconds"

// GenesisStorageSuffix defines the storage suffix used for genesis altered data
const GenesisStorageSuffix = "_genesis"

// MetricAccountsSnapshotNumNodes is the metric that outputs the number of trie nodes written for accounts after snapshot
const MetricAccountsSnapshotNumNodes = "drt_accounts_snapshot_num_nodes"

// MetricTrieSyncNumReceivedBytes is the metric that outputs the number of bytes received for accounts during trie sync
const MetricTrieSyncNumReceivedBytes = "drt_trie_sync_num_bytes_received"

// MetricTrieSyncNumProcessedNodes is the metric that outputs the number of trie nodes processed for accounts during trie sync
const MetricTrieSyncNumProcessedNodes = "drt_trie_sync_num_nodes_processed"

// FullArchiveMetricSuffix is the suffix added to metrics specific for full archive network
const FullArchiveMetricSuffix = "_full_archive"

// Enable epoch flags definitions
const (
	SCDeployFlag                                        core.EnableEpochFlag = "SCDeployFlag"
	BuiltInFunctionsFlag                                core.EnableEpochFlag = "BuiltInFunctionsFlag"
	RelayedTransactionsFlag                             core.EnableEpochFlag = "RelayedTransactionsFlag"
	PenalizedTooMuchGasFlag                             core.EnableEpochFlag = "PenalizedTooMuchGasFlag"
	SwitchJailWaitingFlag                               core.EnableEpochFlag = "SwitchJailWaitingFlag"
	BelowSignedThresholdFlag                            core.EnableEpochFlag = "BelowSignedThresholdFlag"
	SwitchHysteresisForMinNodesFlagInSpecificEpochOnly  core.EnableEpochFlag = "SwitchHysteresisForMinNodesFlagInSpecificEpochOnly"
	TransactionSignedWithTxHashFlag                     core.EnableEpochFlag = "TransactionSignedWithTxHashFlag"
	MetaProtectionFlag                                  core.EnableEpochFlag = "MetaProtectionFlag"
	AheadOfTimeGasUsageFlag                             core.EnableEpochFlag = "AheadOfTimeGasUsageFlag"
	GasPriceModifierFlag                                core.EnableEpochFlag = "GasPriceModifierFlag"
	RepairCallbackFlag                                  core.EnableEpochFlag = "RepairCallbackFlag"
	ReturnDataToLastTransferFlagAfterEpoch              core.EnableEpochFlag = "ReturnDataToLastTransferFlagAfterEpoch"
	SenderInOutTransferFlag                             core.EnableEpochFlag = "SenderInOutTransferFlag"
	StakeFlag                                           core.EnableEpochFlag = "StakeFlag"
	StakingV2Flag                                       core.EnableEpochFlag = "StakingV2Flag"
	StakingV2OwnerFlagInSpecificEpochOnly               core.EnableEpochFlag = "StakingV2OwnerFlagInSpecificEpochOnly"
	StakingV2FlagAfterEpoch                             core.EnableEpochFlag = "StakingV2FlagAfterEpoch"
	DoubleKeyProtectionFlag                             core.EnableEpochFlag = "DoubleKeyProtectionFlag"
	DCDTFlag                                            core.EnableEpochFlag = "DCDTFlag"
	DCDTFlagInSpecificEpochOnly                         core.EnableEpochFlag = "DCDTFlagInSpecificEpochOnly"
	GovernanceFlag                                      core.EnableEpochFlag = "GovernanceFlag"
	GovernanceFlagInSpecificEpochOnly                   core.EnableEpochFlag = "GovernanceFlagInSpecificEpochOnly"
	DelegationManagerFlag                               core.EnableEpochFlag = "DelegationManagerFlag"
	DelegationSmartContractFlag                         core.EnableEpochFlag = "DelegationSmartContractFlag"
	DelegationSmartContractFlagInSpecificEpochOnly      core.EnableEpochFlag = "DelegationSmartContractFlagInSpecificEpochOnly"
	CorrectLastUnJailedFlag                             core.EnableEpochFlag = "CorrectLastUnJailedFlag"
	CorrectLastUnJailedFlagInSpecificEpochOnly          core.EnableEpochFlag = "CorrectLastUnJailedFlagInSpecificEpochOnly"
	RelayedTransactionsV2Flag                           core.EnableEpochFlag = "RelayedTransactionsV2Flag"
	UnBondTokensV2Flag                                  core.EnableEpochFlag = "UnBondTokensV2Flag"
	SaveJailedAlwaysFlag                                core.EnableEpochFlag = "SaveJailedAlwaysFlag"
	ReDelegateBelowMinCheckFlag                         core.EnableEpochFlag = "ReDelegateBelowMinCheckFlag"
	ValidatorToDelegationFlag                           core.EnableEpochFlag = "ValidatorToDelegationFlag"
	IncrementSCRNonceInMultiTransferFlag                core.EnableEpochFlag = "IncrementSCRNonceInMultiTransferFlag"
	DCDTMultiTransferFlag                               core.EnableEpochFlag = "DCDTMultiTransferFlag"
	GlobalMintBurnFlag                                  core.EnableEpochFlag = "GlobalMintBurnFlag"
	DCDTTransferRoleFlag                                core.EnableEpochFlag = "DCDTTransferRoleFlag"
	ComputeRewardCheckpointFlag                         core.EnableEpochFlag = "ComputeRewardCheckpointFlag"
	SCRSizeInvariantCheckFlag                           core.EnableEpochFlag = "SCRSizeInvariantCheckFlag"
	BackwardCompSaveKeyValueFlag                        core.EnableEpochFlag = "BackwardCompSaveKeyValueFlag"
	DCDTNFTCreateOnMultiShardFlag                       core.EnableEpochFlag = "DCDTNFTCreateOnMultiShardFlag"
	MetaDCDTSetFlag                                     core.EnableEpochFlag = "MetaDCDTSetFlag"
	AddTokensToDelegationFlag                           core.EnableEpochFlag = "AddTokensToDelegationFlag"
	MultiDCDTTransferFixOnCallBackFlag                  core.EnableEpochFlag = "MultiDCDTTransferFixOnCallBackFlag"
	OptimizeGasUsedInCrossMiniBlocksFlag                core.EnableEpochFlag = "OptimizeGasUsedInCrossMiniBlocksFlag"
	CorrectFirstQueuedFlag                              core.EnableEpochFlag = "CorrectFirstQueuedFlag"
	DeleteDelegatorAfterClaimRewardsFlag                core.EnableEpochFlag = "DeleteDelegatorAfterClaimRewardsFlag"
	RemoveNonUpdatedStorageFlag                         core.EnableEpochFlag = "RemoveNonUpdatedStorageFlag"
	OptimizeNFTStoreFlag                                core.EnableEpochFlag = "OptimizeNFTStoreFlag"
	CreateNFTThroughExecByCallerFlag                    core.EnableEpochFlag = "CreateNFTThroughExecByCallerFlag"
	StopDecreasingValidatorRatingWhenStuckFlag          core.EnableEpochFlag = "StopDecreasingValidatorRatingWhenStuckFlag"
	FrontRunningProtectionFlag                          core.EnableEpochFlag = "FrontRunningProtectionFlag"
	PayableBySCFlag                                     core.EnableEpochFlag = "PayableBySCFlag"
	CleanUpInformativeSCRsFlag                          core.EnableEpochFlag = "CleanUpInformativeSCRsFlag"
	StorageAPICostOptimizationFlag                      core.EnableEpochFlag = "StorageAPICostOptimizationFlag"
	DCDTRegisterAndSetAllRolesFlag                      core.EnableEpochFlag = "DCDTRegisterAndSetAllRolesFlag"
	ScheduledMiniBlocksFlag                             core.EnableEpochFlag = "ScheduledMiniBlocksFlag"
	CorrectJailedNotUnStakedEmptyQueueFlag              core.EnableEpochFlag = "CorrectJailedNotUnStakedEmptyQueueFlag"
	DoNotReturnOldBlockInBlockchainHookFlag             core.EnableEpochFlag = "DoNotReturnOldBlockInBlockchainHookFlag"
	AddFailedRelayedTxToInvalidMBsFlag                  core.EnableEpochFlag = "AddFailedRelayedTxToInvalidMBsFlag"
	SCRSizeInvariantOnBuiltInResultFlag                 core.EnableEpochFlag = "SCRSizeInvariantOnBuiltInResultFlag"
	CheckCorrectTokenIDForTransferRoleFlag              core.EnableEpochFlag = "CheckCorrectTokenIDForTransferRoleFlag"
	FailExecutionOnEveryAPIErrorFlag                    core.EnableEpochFlag = "FailExecutionOnEveryAPIErrorFlag"
	MiniBlockPartialExecutionFlag                       core.EnableEpochFlag = "MiniBlockPartialExecutionFlag"
	ManagedCryptoAPIsFlag                               core.EnableEpochFlag = "ManagedCryptoAPIsFlag"
	DCDTMetadataContinuousCleanupFlag                   core.EnableEpochFlag = "DCDTMetadataContinuousCleanupFlag"
	DisableExecByCallerFlag                             core.EnableEpochFlag = "DisableExecByCallerFlag"
	RefactorContextFlag                                 core.EnableEpochFlag = "RefactorContextFlag"
	CheckFunctionArgumentFlag                           core.EnableEpochFlag = "CheckFunctionArgumentFlag"
	CheckExecuteOnReadOnlyFlag                          core.EnableEpochFlag = "CheckExecuteOnReadOnlyFlag"
	SetSenderInEeiOutputTransferFlag                    core.EnableEpochFlag = "SetSenderInEeiOutputTransferFlag"
	FixAsyncCallbackCheckFlag                           core.EnableEpochFlag = "FixAsyncCallbackCheckFlag"
	SaveToSystemAccountFlag                             core.EnableEpochFlag = "SaveToSystemAccountFlag"
	CheckFrozenCollectionFlag                           core.EnableEpochFlag = "CheckFrozenCollectionFlag"
	SendAlwaysFlag                                      core.EnableEpochFlag = "SendAlwaysFlag"
	ValueLengthCheckFlag                                core.EnableEpochFlag = "ValueLengthCheckFlag"
	CheckTransferFlag                                   core.EnableEpochFlag = "CheckTransferFlag"
	DCDTNFTImprovementV1Flag                            core.EnableEpochFlag = "DCDTNFTImprovementV1Flag"
	ChangeDelegationOwnerFlag                           core.EnableEpochFlag = "ChangeDelegationOwnerFlag"
	RefactorPeersMiniBlocksFlag                         core.EnableEpochFlag = "RefactorPeersMiniBlocksFlag"
	SCProcessorV2Flag                                   core.EnableEpochFlag = "SCProcessorV2Flag"
	FixAsyncCallBackArgsListFlag                        core.EnableEpochFlag = "FixAsyncCallBackArgsListFlag"
	FixOldTokenLiquidityFlag                            core.EnableEpochFlag = "FixOldTokenLiquidityFlag"
	RuntimeMemStoreLimitFlag                            core.EnableEpochFlag = "RuntimeMemStoreLimitFlag"
	RuntimeCodeSizeFixFlag                              core.EnableEpochFlag = "RuntimeCodeSizeFixFlag"
	MaxBlockchainHookCountersFlag                       core.EnableEpochFlag = "MaxBlockchainHookCountersFlag"
	WipeSingleNFTLiquidityDecreaseFlag                  core.EnableEpochFlag = "WipeSingleNFTLiquidityDecreaseFlag"
	AlwaysSaveTokenMetaDataFlag                         core.EnableEpochFlag = "AlwaysSaveTokenMetaDataFlag"
	SetGuardianFlag                                     core.EnableEpochFlag = "SetGuardianFlag"
	RelayedNonceFixFlag                                 core.EnableEpochFlag = "RelayedNonceFixFlag"
	ConsistentTokensValuesLengthCheckFlag               core.EnableEpochFlag = "ConsistentTokensValuesLengthCheckFlag"
	KeepExecOrderOnCreatedSCRsFlag                      core.EnableEpochFlag = "KeepExecOrderOnCreatedSCRsFlag"
	MultiClaimOnDelegationFlag                          core.EnableEpochFlag = "MultiClaimOnDelegationFlag"
	ChangeUsernameFlag                                  core.EnableEpochFlag = "ChangeUsernameFlag"
	AutoBalanceDataTriesFlag                            core.EnableEpochFlag = "AutoBalanceDataTriesFlag"
	MigrateDataTrieFlag                                 core.EnableEpochFlag = "MigrateDataTrieFlag"
	FixDelegationChangeOwnerOnAccountFlag               core.EnableEpochFlag = "FixDelegationChangeOwnerOnAccountFlag"
	FixOOGReturnCodeFlag                                core.EnableEpochFlag = "FixOOGReturnCodeFlag"
	DeterministicSortOnValidatorsInfoFixFlag            core.EnableEpochFlag = "DeterministicSortOnValidatorsInfoFixFlag"
	DynamicGasCostForDataTrieStorageLoadFlag            core.EnableEpochFlag = "DynamicGasCostForDataTrieStorageLoadFlag"
	ScToScLogEventFlag                                  core.EnableEpochFlag = "ScToScLogEventFlag"
	BlockGasAndFeesReCheckFlag                          core.EnableEpochFlag = "BlockGasAndFeesReCheckFlag"
	BalanceWaitingListsFlag                             core.EnableEpochFlag = "BalanceWaitingListsFlag"
	NFTStopCreateFlag                                   core.EnableEpochFlag = "NFTStopCreateFlag"
	FixGasRemainingForSaveKeyValueFlag                  core.EnableEpochFlag = "FixGasRemainingForSaveKeyValueFlag"
	IsChangeOwnerAddressCrossShardThroughSCFlag         core.EnableEpochFlag = "IsChangeOwnerAddressCrossShardThroughSCFlag"
	CurrentRandomnessOnSortingFlag                      core.EnableEpochFlag = "CurrentRandomnessOnSortingFlag"
	StakeLimitsFlag                                     core.EnableEpochFlag = "StakeLimitsFlag"
	StakingV4Step1Flag                                  core.EnableEpochFlag = "StakingV4Step1Flag"
	StakingV4Step2Flag                                  core.EnableEpochFlag = "StakingV4Step2Flag"
	StakingV4Step3Flag                                  core.EnableEpochFlag = "StakingV4Step3Flag"
	CleanupAuctionOnLowWaitingListFlag                  core.EnableEpochFlag = "CleanupAuctionOnLowWaitingListFlag"
	StakingV4StartedFlag                                core.EnableEpochFlag = "StakingV4StartedFlag"
	AlwaysMergeContextsInEEIFlag                        core.EnableEpochFlag = "AlwaysMergeContextsInEEIFlag"
	UseGasBoundedShouldFailExecutionFlag                core.EnableEpochFlag = "UseGasBoundedShouldFailExecutionFlag"
	DynamicDCDTFlag                                     core.EnableEpochFlag = "DynamicDcdtFlag"
	REWAInDCDTMultiTransferFlag                         core.EnableEpochFlag = "REWAInDCDTMultiTransferFlag"
	CryptoOpcodesV2Flag                                 core.EnableEpochFlag = "CryptoOpcodesV2Flag"
	UnJailCleanupFlag                                   core.EnableEpochFlag = "UnJailCleanupFlag"
	FixRelayedBaseCostFlag                              core.EnableEpochFlag = "FixRelayedBaseCostFlag"
	MultiDCDTNFTTransferAndExecuteByUserFlag            core.EnableEpochFlag = "MultiDCDTNFTTransferAndExecuteByUserFlag"
	FixRelayedMoveBalanceToNonPayableSCFlag             core.EnableEpochFlag = "FixRelayedMoveBalanceToNonPayableSCFlag"
	RelayedTransactionsV3Flag                           core.EnableEpochFlag = "RelayedTransactionsV3Flag"
	RelayedTransactionsV3FixDCDTTransferFlag            core.EnableEpochFlag = "RelayedTransactionsV3FixDCDTTransferFlag"
	AndromedaFlag                                       core.EnableEpochFlag = "AndromedaFlag"
	CheckBuiltInCallOnTransferValueAndFailExecutionFlag core.EnableEpochFlag = "CheckBuiltInCallOnTransferValueAndFailExecutionFlag"
	// all new flags must be added to createAllFlagsMap method, as part of enableEpochsHandler allFlagsDefined
)
