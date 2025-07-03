package factory

import (
	"path"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/TerraDharitri/drt-go-chain/common"
	commonDisabled "github.com/TerraDharitri/drt-go-chain/common/disabled"
	"github.com/TerraDharitri/drt-go-chain/common/statistics"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain/state"
	"github.com/TerraDharitri/drt-go-chain/storage/factory"
	storageFactory "github.com/TerraDharitri/drt-go-chain/storage/factory"
	"github.com/TerraDharitri/drt-go-chain/storage/storageunit"
	"github.com/TerraDharitri/drt-go-chain/trie"
	"github.com/TerraDharitri/drt-go-chain/update"
	"github.com/TerraDharitri/drt-go-chain/update/genesis"
)

// ArgsNewDataTrieFactory is the argument structure for the new data trie factory
type ArgsNewDataTrieFactory struct {
	StorageConfig        config.StorageConfig
	SyncFolder           string
	Marshalizer          marshal.Marshalizer
	Hasher               hashing.Hasher
	ShardCoordinator     sharding.Coordinator
	EnableEpochsHandler  common.EnableEpochsHandler
	StateStatsCollector  common.StateStatisticsHandler
	MaxTrieLevelInMemory uint
}

type dataTrieFactory struct {
	shardCoordinator     sharding.Coordinator
	trieStorage          common.StorageManager
	marshalizer          marshal.Marshalizer
	hasher               hashing.Hasher
	enableEpochsHandler  common.EnableEpochsHandler
	maxTrieLevelInMemory uint
}

// NewDataTrieFactory creates a data trie factory
func NewDataTrieFactory(args ArgsNewDataTrieFactory) (*dataTrieFactory, error) {
	if len(args.SyncFolder) < 2 {
		return nil, update.ErrInvalidFolderName
	}
	if check.IfNil(args.ShardCoordinator) {
		return nil, update.ErrNilShardCoordinator
	}
	if check.IfNil(args.Marshalizer) {
		return nil, update.ErrNilMarshalizer
	}
	if check.IfNil(args.Hasher) {
		return nil, update.ErrNilHasher
	}
	if check.IfNil(args.EnableEpochsHandler) {
		return nil, update.ErrNilEnableEpochsHandler
	}
	if check.IfNil(args.StateStatsCollector) {
		return nil, statistics.ErrNilStateStatsHandler
	}

	dbConfig := storageFactory.GetDBFromConfig(args.StorageConfig.DB)
	dbConfig.FilePath = path.Join(args.SyncFolder, args.StorageConfig.DB.FilePath)

	persisterFactory, err := factory.NewPersisterFactory(args.StorageConfig.DB)
	if err != nil {
		return nil, err
	}

	accountsTrieStorage, err := storageunit.NewStorageUnitFromConf(
		storageFactory.GetCacherFromConfig(args.StorageConfig.Cache),
		dbConfig,
		persisterFactory,
	)
	if err != nil {
		return nil, err
	}
	tsmArgs := trie.NewTrieStorageManagerArgs{
		MainStorer:  accountsTrieStorage,
		Marshalizer: args.Marshalizer,
		Hasher:      args.Hasher,
		GeneralConfig: config.TrieStorageManagerConfig{
			SnapshotsGoroutineNum: 2,
		},
		IdleProvider:   commonDisabled.NewProcessStatusHandler(),
		Identifier:     dataRetriever.UserAccountsUnit.String(),
		StatsCollector: args.StateStatsCollector,
	}
	options := trie.StorageManagerOptions{
		PruningEnabled:   false,
		SnapshotsEnabled: false,
	}
	trieStorage, err := trie.CreateTrieStorageManager(tsmArgs, options)
	if err != nil {
		return nil, err
	}

	d := &dataTrieFactory{
		shardCoordinator:     args.ShardCoordinator,
		trieStorage:          trieStorage,
		marshalizer:          args.Marshalizer,
		hasher:               args.Hasher,
		maxTrieLevelInMemory: args.MaxTrieLevelInMemory,
		enableEpochsHandler:  args.EnableEpochsHandler,
	}

	return d, nil
}

// TrieStorageManager returns trie storage manager
func (d *dataTrieFactory) TrieStorageManager() common.StorageManager {
	return d.trieStorage
}

// Create creates a TriesHolder container to hold all the states
func (d *dataTrieFactory) Create() (common.TriesHolder, error) {
	container := state.NewDataTriesHolder()

	for i := uint32(0); i < d.shardCoordinator.NumberOfShards(); i++ {
		err := d.createAndAddOneTrie(i, genesis.UserAccount, container)
		if err != nil {
			return nil, err
		}
	}

	err := d.createAndAddOneTrie(core.MetachainShardId, genesis.UserAccount, container)
	if err != nil {
		return nil, err
	}

	err = d.createAndAddOneTrie(core.MetachainShardId, genesis.ValidatorAccount, container)
	if err != nil {
		return nil, err
	}

	return container, nil
}

func (d *dataTrieFactory) createAndAddOneTrie(shId uint32, accType genesis.Type, container common.TriesHolder) error {
	dataTrie, err := trie.NewTrie(d.trieStorage, d.marshalizer, d.hasher, d.enableEpochsHandler, d.maxTrieLevelInMemory)
	if err != nil {
		return err
	}

	trieId := genesis.CreateTrieIdentifier(shId, accType)
	container.Put([]byte(trieId), dataTrie)

	return nil
}

// IsInterfaceNil returns true if underlying object is nil
func (d *dataTrieFactory) IsInterfaceNil() bool {
	return d == nil
}
