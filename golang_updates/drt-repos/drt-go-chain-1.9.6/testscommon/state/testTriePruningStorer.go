package state

import (
	"sync"

	"github.com/TerraDharitri/drt-go-chain/common/statistics/disabled"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain/storage"
	"github.com/TerraDharitri/drt-go-chain/storage/database"
	storageMock "github.com/TerraDharitri/drt-go-chain/storage/mock"
	"github.com/TerraDharitri/drt-go-chain/storage/pruning"
	"github.com/TerraDharitri/drt-go-chain/storage/storageunit"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
)

// CreateTestingTriePruningStorer creates a new trie pruning storer that is used for testing
func CreateTestingTriePruningStorer(coordinator sharding.Coordinator, notifier pruning.EpochStartNotifier) (storage.Storer, *persisterMap, error) {
	cacheConf := storageunit.CacheConfig{
		Capacity: 10,
		Type:     "LRU",
		Shards:   3,
	}
	dbConf := storageunit.DBConfig{
		FilePath:          "path/Epoch_0/Shard_1",
		Type:              "LvlDBSerial",
		BatchDelaySeconds: 500,
		MaxBatchSize:      1,
		MaxOpenFiles:      1000,
	}

	persistersMap := NewPersistersMap()
	persisterFactory := &storageMock.PersisterFactoryStub{
		CreateCalled: func(path string) (storage.Persister, error) {
			return persistersMap.GetPersister(path), nil
		},
	}
	epochsData := pruning.EpochArgs{
		NumOfEpochsToKeep:     4,
		NumOfActivePersisters: 4,
	}
	args := pruning.StorerArgs{
		PruningEnabled:         true,
		Identifier:             "id",
		ShardCoordinator:       coordinator,
		PathManager:            &testscommon.PathManagerStub{},
		CacheConf:              cacheConf,
		DbPath:                 dbConf.FilePath,
		PersisterFactory:       persisterFactory,
		EpochsData:             epochsData,
		Notifier:               notifier,
		OldDataCleanerProvider: &testscommon.OldDataCleanerProviderStub{},
		CustomDatabaseRemover:  &testscommon.CustomDatabaseRemoverStub{},
		MaxBatchSize:           10,
		PersistersTracker:      pruning.NewPersistersTracker(epochsData),
		StateStatsHandler:      disabled.NewStateStatistics(),
	}

	tps, err := pruning.NewTriePruningStorer(args)
	return tps, persistersMap, err
}

type persisterMap struct {
	persisters map[string]storage.Persister
	mutex      sync.Mutex
}

// NewPersistersMap returns a new persisterMap
func NewPersistersMap() *persisterMap {
	return &persisterMap{
		persisters: make(map[string]storage.Persister),
		mutex:      sync.Mutex{},
	}
}

// GetPersister returns the persister for the given path, or creates a new persister if it does not exist
func (pm *persisterMap) GetPersister(path string) storage.Persister {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	persister, exists := pm.persisters[path]
	if !exists {
		persister = database.NewMemDB()
		pm.persisters[path] = persister
	}

	return persister
}
