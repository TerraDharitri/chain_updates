package integrationTests

import (
	"math/big"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/debug"
	"github.com/TerraDharitri/drt-go-chain/epochStart"
	"github.com/TerraDharitri/drt-go-chain/heartbeat/data"
	"github.com/TerraDharitri/drt-go-chain/node/external"
	"github.com/TerraDharitri/drt-go-chain/process"
	txSimData "github.com/TerraDharitri/drt-go-chain/process/transactionEvaluator/data"
	"github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data/alteredAccount"
	"github.com/TerraDharitri/drt-go-chain-core/data/api"
	dataApi "github.com/TerraDharitri/drt-go-chain-core/data/api"
	"github.com/TerraDharitri/drt-go-chain-core/data/dcdt"
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	"github.com/TerraDharitri/drt-go-chain-core/data/validator"
	"github.com/TerraDharitri/drt-go-chain-core/data/vm"
)

// TestBootstrapper extends the Bootstrapper interface with some functions intended to be used only in tests
// as it simplifies the reproduction of edge cases
type TestBootstrapper interface {
	process.Bootstrapper
	RollBack(revertUsingForkNonce bool) error
	SetProbableHighestNonce(nonce uint64)
}

// TestEpochStartTrigger extends the epochStart trigger interface with some functions intended to by used only
// in tests as it simplifies the reproduction of test scenarios
type TestEpochStartTrigger interface {
	epochStart.TriggerHandler
	GetRoundsPerEpoch() uint64
	SetTrigger(triggerHandler epochStart.TriggerHandler)
	SetRoundsPerEpoch(roundsPerEpoch uint64)
	SetMinRoundsBetweenEpochs(minRoundsPerEpoch uint64)
	SetEpoch(epoch uint32)
}

// NodesCoordinatorFactory is used for creating a nodesCoordinator in the integration tests
type NodesCoordinatorFactory interface {
	CreateNodesCoordinator(arg ArgIndexHashedNodesCoordinatorFactory) nodesCoordinator.NodesCoordinator
}

// NetworkShardingUpdater defines the updating methods used by the network sharding component
type NetworkShardingUpdater interface {
	GetPeerInfo(pid core.PeerID) core.P2PPeerInfo
	UpdatePeerIDPublicKeyPair(pid core.PeerID, pk []byte)
	PutPeerIdShardId(pid core.PeerID, shardID uint32)
	UpdatePeerIDInfo(pid core.PeerID, pk []byte, shardID uint32)
	PutPeerIdSubType(pid core.PeerID, peerSubType core.P2PPeerSubType)
	IsInterfaceNil() bool
}

// Facade is the node facade used to decouple the node implementation with the web server. Used in integration tests
type Facade interface {
	GetBalance(address string, options api.AccountQueryOptions) (*big.Int, api.BlockInfo, error)
	GetUsername(address string, options api.AccountQueryOptions) (string, api.BlockInfo, error)
	GetCodeHash(address string, options api.AccountQueryOptions) ([]byte, api.BlockInfo, error)
	GetValueForKey(address string, key string, options api.AccountQueryOptions) (string, api.BlockInfo, error)
	GetAccount(address string, options api.AccountQueryOptions) (dataApi.AccountResponse, api.BlockInfo, error)
	GetAccounts(addresses []string, options api.AccountQueryOptions) (map[string]*api.AccountResponse, api.BlockInfo, error)
	GetDCDTData(address string, key string, nonce uint64, options api.AccountQueryOptions) (*dcdt.DCDigitalToken, api.BlockInfo, error)
	GetNFTTokenIDsRegisteredByAddress(address string, options api.AccountQueryOptions) ([]string, api.BlockInfo, error)
	GetDCDTsWithRole(address string, role string, options api.AccountQueryOptions) ([]string, api.BlockInfo, error)
	GetAllDCDTTokens(address string, options api.AccountQueryOptions) (map[string]*dcdt.DCDigitalToken, api.BlockInfo, error)
	GetDCDTsRoles(address string, options api.AccountQueryOptions) (map[string][]string, api.BlockInfo, error)
	GetKeyValuePairs(address string, options api.AccountQueryOptions) (map[string]string, api.BlockInfo, error)
	GetGuardianData(address string, options api.AccountQueryOptions) (api.GuardianData, api.BlockInfo, error)
	GetBlockByHash(hash string, options api.BlockQueryOptions) (*dataApi.Block, error)
	GetBlockByNonce(nonce uint64, options api.BlockQueryOptions) (*dataApi.Block, error)
	GetBlockByRound(round uint64, options api.BlockQueryOptions) (*dataApi.Block, error)
	Trigger(epoch uint32, withEarlyEndOfEpoch bool) error
	IsSelfTrigger() bool
	GetTotalStakedValue() (*dataApi.StakeValues, error)
	GetDirectStakedList() ([]*dataApi.DirectStakedValue, error)
	GetDelegatorsList() ([]*dataApi.Delegator, error)
	GetAllIssuedDCDTs(tokenType string) ([]string, error)
	GetTokenSupply(token string) (*dataApi.DCDTSupply, error)
	GetHeartbeats() ([]data.PubKeyHeartbeat, error)
	StatusMetrics() external.StatusMetricsHandler
	GetQueryHandler(name string) (debug.QueryHandler, error)
	GetEpochStartDataAPI(epoch uint32) (*common.EpochStartDataAPI, error)
	GetPeerInfo(pid string) ([]core.QueryP2PPeerInfo, error)
	GetConnectedPeersRatingsOnMainNetwork() (string, error)
	CreateTransaction(txArgs *external.ArgsCreateTransaction) (*transaction.Transaction, []byte, error)
	ValidateTransaction(tx *transaction.Transaction) error
	ValidateTransactionForSimulation(tx *transaction.Transaction, bypassSignature bool) error
	SendBulkTransactions([]*transaction.Transaction) (uint64, error)
	SimulateTransactionExecution(tx *transaction.Transaction) (*txSimData.SimulationResultsWithVMOutput, error)
	GetTransaction(hash string, withResults bool) (*transaction.ApiTransactionResult, error)
	ComputeTransactionGasLimit(tx *transaction.Transaction) (*transaction.CostResponse, error)
	EncodeAddressPubkey(pk []byte) (string, error)
	GetThrottlerForEndpoint(endpoint string) (core.Throttler, bool)
	ValidatorStatisticsApi() (map[string]*validator.ValidatorStatistics, error)
	AuctionListApi() ([]*common.AuctionListValidatorAPIResponse, error)
	ExecuteSCQuery(*process.SCQuery) (*vm.VMOutputApi, api.BlockInfo, error)
	DecodeAddressPubkey(pk string) ([]byte, error)
	GetProof(rootHash string, address string) (*common.GetProofResponse, error)
	GetProofDataTrie(rootHash string, address string, key string) (*common.GetProofResponse, *common.GetProofResponse, error)
	GetProofCurrentRootHash(address string) (*common.GetProofResponse, error)
	VerifyProof(rootHash string, address string, proof [][]byte) (bool, error)
	GetGenesisNodesPubKeys() (map[uint32][]string, map[uint32][]string, error)
	GetGenesisBalances() ([]*common.InitialAccountAPI, error)
	GetGasConfigs() (map[string]map[string]uint64, error)
	GetTransactionsPool(fields string) (*common.TransactionsPoolAPIResponse, error)
	GetTransactionsPoolForSender(sender, fields string) (*common.TransactionsPoolForSenderApiResponse, error)
	GetLastPoolNonceForSender(sender string) (uint64, error)
	GetTransactionsPoolNonceGapsForSender(sender string) (*common.TransactionsPoolNonceGapsForSenderApiResponse, error)
	GetAlteredAccountsForBlock(options dataApi.GetAlteredAccountsForBlockOptions) ([]*alteredAccount.AlteredAccount, error)
	IsDataTrieMigrated(address string, options api.AccountQueryOptions) (bool, error)
	GetManagedKeysCount() int
	GetManagedKeys() []string
	GetLoadedKeys() []string
	GetEligibleManagedKeys() ([]string, error)
	GetWaitingManagedKeys() ([]string, error)
	GetWaitingEpochsLeftForPublicKey(publicKey string) (uint32, error)
	GetSCRsByTxHash(txHash string, scrHash string) ([]*transaction.ApiSmartContractResult, error)
	IsInterfaceNil() bool
}
