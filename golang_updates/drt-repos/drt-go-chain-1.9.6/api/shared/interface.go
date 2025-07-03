package shared

import (
	"math/big"

	"github.com/gin-gonic/gin"
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data/alteredAccount"
	"github.com/TerraDharitri/drt-go-chain-core/data/api"
	"github.com/TerraDharitri/drt-go-chain-core/data/dcdt"
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	"github.com/TerraDharitri/drt-go-chain-core/data/validator"
	"github.com/TerraDharitri/drt-go-chain-core/data/vm"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/debug"
	"github.com/TerraDharitri/drt-go-chain/heartbeat/data"
	"github.com/TerraDharitri/drt-go-chain/node/external"
	"github.com/TerraDharitri/drt-go-chain/process"
	txSimData "github.com/TerraDharitri/drt-go-chain/process/transactionEvaluator/data"
	"github.com/TerraDharitri/drt-go-chain/state"
)

// HttpServerCloser defines the basic actions of starting and closing that a web server should be able to do
type HttpServerCloser interface {
	Start()
	Close() error
	IsInterfaceNil() bool
}

// MiddlewareProcessor defines a processor used internally by the web server when processing requests
type MiddlewareProcessor interface {
	MiddlewareHandlerFunc() gin.HandlerFunc
	IsInterfaceNil() bool
}

// ApiFacadeHandler interface defines methods that can be used by the web server to interact with the node
type ApiFacadeHandler interface {
	RestApiInterface() string
	RestAPIServerDebugMode() bool
	PprofEnabled() bool
	IsInterfaceNil() bool
}

// UpgradeableHttpServerHandler defines the actions that an upgradeable http server need to do
type UpgradeableHttpServerHandler interface {
	StartHttpServer() error
	UpdateFacade(facade FacadeHandler) error
	Close() error
	IsInterfaceNil() bool
}

// GroupHandler defines the actions needed to be performed by a gin API group
type GroupHandler interface {
	UpdateFacade(newFacade interface{}) error
	RegisterRoutes(
		ws *gin.RouterGroup,
		apiConfig config.ApiRoutesConfig,
	)
	IsInterfaceNil() bool
}

// FacadeHandler defines all the methods that a facade should implement
type FacadeHandler interface {
	GetBalance(address string, options api.AccountQueryOptions) (*big.Int, api.BlockInfo, error)
	GetUsername(address string, options api.AccountQueryOptions) (string, api.BlockInfo, error)
	GetCodeHash(address string, options api.AccountQueryOptions) ([]byte, api.BlockInfo, error)
	GetValueForKey(address string, key string, options api.AccountQueryOptions) (string, api.BlockInfo, error)
	GetAccount(address string, options api.AccountQueryOptions) (api.AccountResponse, api.BlockInfo, error)
	GetAccounts(addresses []string, options api.AccountQueryOptions) (map[string]*api.AccountResponse, api.BlockInfo, error)
	GetDCDTData(address string, key string, nonce uint64, options api.AccountQueryOptions) (*dcdt.DCDigitalToken, api.BlockInfo, error)
	GetDCDTsRoles(address string, options api.AccountQueryOptions) (map[string][]string, api.BlockInfo, error)
	GetNFTTokenIDsRegisteredByAddress(address string, options api.AccountQueryOptions) ([]string, api.BlockInfo, error)
	GetDCDTsWithRole(address string, role string, options api.AccountQueryOptions) ([]string, api.BlockInfo, error)
	GetAllDCDTTokens(address string, options api.AccountQueryOptions) (map[string]*dcdt.DCDigitalToken, api.BlockInfo, error)
	GetKeyValuePairs(address string, options api.AccountQueryOptions) (map[string]string, api.BlockInfo, error)
	GetGuardianData(address string, options api.AccountQueryOptions) (api.GuardianData, api.BlockInfo, error)
	GetBlockByHash(hash string, options api.BlockQueryOptions) (*api.Block, error)
	GetBlockByNonce(nonce uint64, options api.BlockQueryOptions) (*api.Block, error)
	GetBlockByRound(round uint64, options api.BlockQueryOptions) (*api.Block, error)
	GetAlteredAccountsForBlock(options api.GetAlteredAccountsForBlockOptions) ([]*alteredAccount.AlteredAccount, error)
	GetInternalShardBlockByNonce(format common.ApiOutputFormat, nonce uint64) (interface{}, error)
	GetInternalShardBlockByHash(format common.ApiOutputFormat, hash string) (interface{}, error)
	GetInternalShardBlockByRound(format common.ApiOutputFormat, round uint64) (interface{}, error)
	GetInternalMetaBlockByNonce(format common.ApiOutputFormat, nonce uint64) (interface{}, error)
	GetInternalMetaBlockByHash(format common.ApiOutputFormat, hash string) (interface{}, error)
	GetInternalMetaBlockByRound(format common.ApiOutputFormat, round uint64) (interface{}, error)
	GetInternalStartOfEpochMetaBlock(format common.ApiOutputFormat, epoch uint32) (interface{}, error)
	GetInternalStartOfEpochValidatorsInfo(epoch uint32) ([]*state.ShardValidatorInfo, error)
	GetInternalMiniBlockByHash(format common.ApiOutputFormat, hash string, epoch uint32) (interface{}, error)
	Trigger(epoch uint32, withEarlyEndOfEpoch bool) error
	IsSelfTrigger() bool
	GetTotalStakedValue() (*api.StakeValues, error)
	GetDirectStakedList() ([]*api.DirectStakedValue, error)
	GetDelegatorsList() ([]*api.Delegator, error)
	StatusMetrics() external.StatusMetricsHandler
	GetTokenSupply(token string) (*api.DCDTSupply, error)
	GetAllIssuedDCDTs(tokenType string) ([]string, error)
	GetHeartbeats() ([]data.PubKeyHeartbeat, error)
	GetQueryHandler(name string) (debug.QueryHandler, error)
	GetEpochStartDataAPI(epoch uint32) (*common.EpochStartDataAPI, error)
	GetPeerInfo(pid string) ([]core.QueryP2PPeerInfo, error)
	GetConnectedPeersRatingsOnMainNetwork() (string, error)
	GetProof(rootHash string, address string) (*common.GetProofResponse, error)
	GetProofDataTrie(rootHash string, address string, key string) (*common.GetProofResponse, *common.GetProofResponse, error)
	GetProofCurrentRootHash(address string) (*common.GetProofResponse, error)
	VerifyProof(rootHash string, address string, proof [][]byte) (bool, error)
	GetThrottlerForEndpoint(endpoint string) (core.Throttler, bool)
	CreateTransaction(txArgs *external.ArgsCreateTransaction) (*transaction.Transaction, []byte, error)
	ValidateTransaction(tx *transaction.Transaction) error
	ValidateTransactionForSimulation(tx *transaction.Transaction, checkSignature bool) error
	SendBulkTransactions([]*transaction.Transaction) (uint64, error)
	SimulateTransactionExecution(tx *transaction.Transaction) (*txSimData.SimulationResultsWithVMOutput, error)
	GetTransaction(hash string, withResults bool) (*transaction.ApiTransactionResult, error)
	ComputeTransactionGasLimit(tx *transaction.Transaction) (*transaction.CostResponse, error)
	EncodeAddressPubkey(pk []byte) (string, error)
	ValidatorStatisticsApi() (map[string]*validator.ValidatorStatistics, error)
	AuctionListApi() ([]*common.AuctionListValidatorAPIResponse, error)
	ExecuteSCQuery(*process.SCQuery) (*vm.VMOutputApi, api.BlockInfo, error)
	DecodeAddressPubkey(pk string) ([]byte, error)
	RestApiInterface() string
	RestAPIServerDebugMode() bool
	PprofEnabled() bool
	GetGenesisNodesPubKeys() (map[uint32][]string, map[uint32][]string, error)
	GetGenesisBalances() ([]*common.InitialAccountAPI, error)
	GetGasConfigs() (map[string]map[string]uint64, error)
	GetTransactionsPool(fields string) (*common.TransactionsPoolAPIResponse, error)
	GetTransactionsPoolForSender(sender, fields string) (*common.TransactionsPoolForSenderApiResponse, error)
	GetLastPoolNonceForSender(sender string) (uint64, error)
	GetTransactionsPoolNonceGapsForSender(sender string) (*common.TransactionsPoolNonceGapsForSenderApiResponse, error)
	IsDataTrieMigrated(address string, options api.AccountQueryOptions) (bool, error)
	GetManagedKeysCount() int
	GetManagedKeys() []string
	GetLoadedKeys() []string
	GetEligibleManagedKeys() ([]string, error)
	GetWaitingManagedKeys() ([]string, error)
	GetWaitingEpochsLeftForPublicKey(publicKey string) (uint32, error)
	GetSCRsByTxHash(txHash string, scrHash string) ([]*transaction.ApiSmartContractResult, error)
	P2PPrometheusMetricsEnabled() bool
	IsInterfaceNil() bool
}
