package process

import (
	"math/big"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/typeConverters"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	crypto "github.com/TerraDharitri/drt-go-chain-crypto"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/dblookupext"
	"github.com/TerraDharitri/drt-go-chain/genesis"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain/state"
	"github.com/TerraDharitri/drt-go-chain/update"
)

type coreComponentsHandler interface {
	InternalMarshalizer() marshal.Marshalizer
	TxMarshalizer() marshal.Marshalizer
	Hasher() hashing.Hasher
	AddressPubKeyConverter() core.PubkeyConverter
	Uint64ByteSliceConverter() typeConverters.Uint64ByteSliceConverter
	TxVersionChecker() process.TxVersionCheckerHandler
	ChainID() string
	EnableEpochsHandler() common.EnableEpochsHandler
	EconomicsData() process.EconomicsDataHandler
	IsInterfaceNil() bool
}

type dataComponentsHandler interface {
	StorageService() dataRetriever.StorageService
	Blockchain() data.ChainHandler
	Datapool() dataRetriever.PoolsHolder
	SetBlockchain(chain data.ChainHandler) error
	Clone() interface{}
	IsInterfaceNil() bool
}

// ArgsGenesisBlockCreator holds the arguments which are needed to create a genesis block
type ArgsGenesisBlockCreator struct {
	GenesisTime             uint64
	GenesisNonce            uint64
	GenesisRound            uint64
	StartEpochNum           uint32
	GenesisEpoch            uint32
	Data                    dataComponentsHandler
	Core                    coreComponentsHandler
	Accounts                state.AccountsAdapter
	ValidatorAccounts       state.AccountsAdapter
	InitialNodesSetup       genesis.InitialNodesHandler
	Economics               process.EconomicsDataHandler
	ShardCoordinator        sharding.Coordinator
	AccountsParser          genesis.AccountsParser
	SmartContractParser     genesis.InitialSmartContractParser
	GasSchedule             core.GasScheduleNotifier
	TxLogsProcessor         process.TransactionLogProcessor
	VirtualMachineConfig    config.VirtualMachineConfig
	HardForkConfig          config.HardforkConfig
	TrieStorageManagers     map[string]common.StorageManager
	SystemSCConfig          config.SystemSmartContractsConfig
	RoundConfig             config.RoundConfig
	EpochConfig             config.EpochConfig
	HeaderVersionConfigs    config.VersionsConfig
	WorkingDir              string
	BlockSignKeyGen         crypto.KeyGenerator
	HistoryRepository       dblookupext.HistoryRepository
	TxExecutionOrderHandler common.TxExecutionOrderHandler

	GenesisNodePrice *big.Int
	GenesisString    string

	// created components
	importHandler          update.ImportHandler
	versionedHeaderFactory genesis.VersionedHeaderFactory
}
