package sync

import (
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/typeConverters"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/consensus"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/dblookupext"
	"github.com/TerraDharitri/drt-go-chain/outport"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain/state"
)

// ArgBaseBootstrapper holds all dependencies required by the bootstrap data factory in order to create
// new instances
type ArgBaseBootstrapper struct {
	HistoryRepo                  dblookupext.HistoryRepository
	PoolsHolder                  dataRetriever.PoolsHolder
	Store                        dataRetriever.StorageService
	ChainHandler                 data.ChainHandler
	RoundHandler                 consensus.RoundHandler
	BlockProcessor               process.BlockProcessor
	WaitTime                     time.Duration
	Hasher                       hashing.Hasher
	Marshalizer                  marshal.Marshalizer
	ForkDetector                 process.ForkDetector
	RequestHandler               process.RequestHandler
	ShardCoordinator             sharding.Coordinator
	Accounts                     state.AccountsAdapter
	BlackListHandler             process.TimeCacher
	NetworkWatcher               process.NetworkConnectionWatcher
	BootStorer                   process.BootStorer
	StorageBootstrapper          process.BootstrapperFromStorage
	EpochHandler                 dataRetriever.EpochHandler
	MiniblocksProvider           process.MiniBlockProvider
	Uint64Converter              typeConverters.Uint64ByteSliceConverter
	AppStatusHandler             core.AppStatusHandler
	OutportHandler               outport.OutportHandler
	AccountsDBSyncer             process.AccountsDBSyncer
	CurrentEpochProvider         process.CurrentNetworkEpochProviderHandler
	IsInImportMode               bool
	ScheduledTxsExecutionHandler process.ScheduledTxsExecutionHandler
	ProcessWaitTime              time.Duration
	RepopulateTokensSupplies     bool
	EnableEpochsHandler          common.EnableEpochsHandler
}

// ArgShardBootstrapper holds all dependencies required by the bootstrap data factory in order to create
// new instances of shard bootstrapper
type ArgShardBootstrapper struct {
	ArgBaseBootstrapper
}

// ArgMetaBootstrapper holds all dependencies required by the bootstrap data factory in order to create
// new instances of meta bootstrapper
type ArgMetaBootstrapper struct {
	ArgBaseBootstrapper
	EpochBootstrapper           process.EpochBootstrapper
	ValidatorStatisticsDBSyncer process.AccountsDBSyncer
	ValidatorAccountsDB         state.AccountsAdapter
}
