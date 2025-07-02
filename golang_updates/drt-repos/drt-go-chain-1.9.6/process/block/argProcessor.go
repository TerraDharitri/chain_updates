package block

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/typeConverters"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"

	nodeFactory "github.com/TerraDharitri/drt-go-chain/cmd/node/factory"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/consensus"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/dblookupext"
	"github.com/TerraDharitri/drt-go-chain/outport"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/block/cutoff"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
	"github.com/TerraDharitri/drt-go-chain/state"
)

type coreComponentsHolder interface {
	Hasher() hashing.Hasher
	InternalMarshalizer() marshal.Marshalizer
	Uint64ByteSliceConverter() typeConverters.Uint64ByteSliceConverter
	EpochNotifier() process.EpochNotifier
	EnableEpochsHandler() common.EnableEpochsHandler
	RoundNotifier() process.RoundNotifier
	EnableRoundsHandler() process.EnableRoundsHandler
	EpochChangeGracePeriodHandler() common.EpochChangeGracePeriodHandler
	RoundHandler() consensus.RoundHandler
	EconomicsData() process.EconomicsDataHandler
	ProcessStatusHandler() common.ProcessStatusHandler
	IsInterfaceNil() bool
}

type dataComponentsHolder interface {
	StorageService() dataRetriever.StorageService
	Datapool() dataRetriever.PoolsHolder
	Blockchain() data.ChainHandler
	IsInterfaceNil() bool
}

type bootstrapComponentsHolder interface {
	ShardCoordinator() sharding.Coordinator
	VersionedHeaderFactory() nodeFactory.VersionedHeaderFactory
	HeaderIntegrityVerifier() nodeFactory.HeaderIntegrityVerifierHandler
	IsInterfaceNil() bool
}

type statusComponentsHolder interface {
	OutportHandler() outport.OutportHandler
	IsInterfaceNil() bool
}

type statusCoreComponentsHolder interface {
	AppStatusHandler() core.AppStatusHandler
	IsInterfaceNil() bool
}

// ArgBaseProcessor holds all dependencies required by the process data factory in order to create
// new instances
type ArgBaseProcessor struct {
	CoreComponents       coreComponentsHolder
	DataComponents       dataComponentsHolder
	BootstrapComponents  bootstrapComponentsHolder
	StatusComponents     statusComponentsHolder
	StatusCoreComponents statusCoreComponentsHolder

	Config                         config.Config
	PrefsConfig                    config.Preferences
	AccountsDB                     map[state.AccountsDbIdentifier]state.AccountsAdapter
	ForkDetector                   process.ForkDetector
	NodesCoordinator               nodesCoordinator.NodesCoordinator
	FeeHandler                     process.TransactionFeeHandler
	RequestHandler                 process.RequestHandler
	BlockChainHook                 process.BlockChainHookHandler
	TxCoordinator                  process.TransactionCoordinator
	EpochStartTrigger              process.EpochStartTriggerHandler
	HeaderValidator                process.HeaderConstructionValidator
	BootStorer                     process.BootStorer
	BlockTracker                   process.BlockTracker
	BlockSizeThrottler             process.BlockSizeThrottler
	Version                        string
	HistoryRepository              dblookupext.HistoryRepository
	VMContainersFactory            process.VirtualMachinesContainerFactory
	VmContainer                    process.VirtualMachinesContainer
	GasHandler                     gasConsumedProvider
	OutportDataProvider            outport.DataProviderOutport
	ScheduledTxsExecutionHandler   process.ScheduledTxsExecutionHandler
	ScheduledMiniBlocksEnableEpoch uint32
	ProcessedMiniBlocksTracker     process.ProcessedMiniBlocksTracker
	ReceiptsRepository             receiptsRepository
	BlockProcessingCutoffHandler   cutoff.BlockProcessingCutoffHandler
	ManagedPeersHolder             common.ManagedPeersHolder
	SentSignaturesTracker          process.SentSignaturesTracker
}

// ArgShardProcessor holds all dependencies required by the process data factory in order to create
// new instances of shard processor
type ArgShardProcessor struct {
	ArgBaseProcessor
}

// ArgMetaProcessor holds all dependencies required by the process data factory in order to create
// new instances of meta processor
type ArgMetaProcessor struct {
	ArgBaseProcessor
	PendingMiniBlocksHandler     process.PendingMiniBlocksHandler
	SCToProtocol                 process.SmartContractToProtocolHandler
	EpochStartDataCreator        process.EpochStartDataCreator
	EpochEconomics               process.EndOfEpochEconomics
	EpochRewardsCreator          process.RewardsCreator
	EpochValidatorInfoCreator    process.EpochStartValidatorInfoCreator
	EpochSystemSCProcessor       process.EpochStartSystemSCProcessor
	ValidatorStatisticsProcessor process.ValidatorStatisticsProcessor
}
