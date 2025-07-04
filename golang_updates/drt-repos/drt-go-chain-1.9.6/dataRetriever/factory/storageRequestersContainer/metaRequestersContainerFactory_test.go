package storagerequesterscontainer_test

import (
	"strings"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data/endProcess"
	"github.com/TerraDharitri/drt-go-chain/common/statistics"
	"github.com/TerraDharitri/drt-go-chain/common/statistics/disabled"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	storagerequesterscontainer "github.com/TerraDharitri/drt-go-chain/dataRetriever/factory/storageRequestersContainer"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever/mock"
	"github.com/TerraDharitri/drt-go-chain/p2p"
	"github.com/TerraDharitri/drt-go-chain/storage"
	"github.com/TerraDharitri/drt-go-chain/testscommon/enableEpochsHandlerMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/hashingMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/p2pmocks"
	storageStubs "github.com/TerraDharitri/drt-go-chain/testscommon/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createMessengerStubForMeta(matchStrToErrOnCreate string, matchStrToErrOnRegister string) p2p.Messenger {
	stub := &p2pmocks.MessengerStub{}

	stub.CreateTopicCalled = func(name string, createChannelForTopic bool) error {
		if matchStrToErrOnCreate == "" {
			return nil
		}
		if strings.Contains(name, matchStrToErrOnCreate) {
			return errExpected
		}

		return nil
	}

	stub.RegisterMessageProcessorCalled = func(topic string, identifier string, handler p2p.MessageProcessor) error {
		if matchStrToErrOnRegister == "" {
			return nil
		}
		if strings.Contains(topic, matchStrToErrOnRegister) {
			return errExpected
		}

		return nil
	}

	return stub
}

func createStoreForMeta() dataRetriever.StorageService {
	return &storageStubs.ChainStorerStub{
		GetStorerCalled: func(unitType dataRetriever.UnitType) (storage.Storer, error) {
			return &storageStubs.StorerStub{}, nil
		},
	}
}

func TestNewMetaRequestersContainerFactory_NilShardCoordinatorShouldErr(t *testing.T) {
	t.Parallel()

	args := getArgumentsMeta()
	args.ShardCoordinator = nil
	rcf, err := storagerequesterscontainer.NewMetaRequestersContainerFactory(args)

	assert.Nil(t, rcf)
	assert.Equal(t, dataRetriever.ErrNilShardCoordinator, err)
}

func TestNewMetaRequestersContainerFactory_NilMessengerShouldErr(t *testing.T) {
	t.Parallel()

	args := getArgumentsMeta()
	args.Messenger = nil
	rcf, err := storagerequesterscontainer.NewMetaRequestersContainerFactory(args)

	assert.Nil(t, rcf)
	assert.Equal(t, dataRetriever.ErrNilMessenger, err)
}

func TestNewMetaRequestersContainerFactory_NilStoreShouldErr(t *testing.T) {
	t.Parallel()

	args := getArgumentsMeta()
	args.Store = nil
	rcf, err := storagerequesterscontainer.NewMetaRequestersContainerFactory(args)

	assert.Nil(t, rcf)
	assert.Equal(t, dataRetriever.ErrNilStore, err)
}

func TestNewMetaRequestersContainerFactory_NilMarshalizerShouldErr(t *testing.T) {
	t.Parallel()

	args := getArgumentsMeta()
	args.Marshalizer = nil
	rcf, err := storagerequesterscontainer.NewMetaRequestersContainerFactory(args)

	assert.Nil(t, rcf)
	assert.Equal(t, dataRetriever.ErrNilMarshalizer, err)
}

func TestNewMetaRequestersContainerFactory_NilUint64SliceConverterShouldErr(t *testing.T) {
	t.Parallel()

	args := getArgumentsMeta()
	args.Uint64ByteSliceConverter = nil
	rcf, err := storagerequesterscontainer.NewMetaRequestersContainerFactory(args)

	assert.Nil(t, rcf)
	assert.Equal(t, dataRetriever.ErrNilUint64ByteSliceConverter, err)
}

func TestNewMetaRequestersContainerFactory_NilDataPackerShouldErr(t *testing.T) {
	t.Parallel()

	args := getArgumentsMeta()
	args.DataPacker = nil
	rcf, err := storagerequesterscontainer.NewMetaRequestersContainerFactory(args)

	assert.Nil(t, rcf)
	assert.Equal(t, dataRetriever.ErrNilDataPacker, err)
}

func TestNewMetaRequestersContainerFactory_NilStateStatsShouldErr(t *testing.T) {
	t.Parallel()

	args := getArgumentsMeta()
	args.StateStatsHandler = nil
	rcf, err := storagerequesterscontainer.NewMetaRequestersContainerFactory(args)

	assert.Nil(t, rcf)
	assert.Equal(t, statistics.ErrNilStateStatsHandler, err)
}

func TestNewMetaRequestersContainerFactory_ShouldWork(t *testing.T) {
	t.Parallel()

	args := getArgumentsMeta()
	rcf, err := storagerequesterscontainer.NewMetaRequestersContainerFactory(args)

	assert.Nil(t, err)
	assert.False(t, check.IfNil(rcf))
}

func TestMetaRequestersContainerFactory_CreateShouldWork(t *testing.T) {
	t.Parallel()

	args := getArgumentsMeta()
	rcf, _ := storagerequesterscontainer.NewMetaRequestersContainerFactory(args)

	container, err := rcf.Create()

	assert.NotNil(t, container)
	assert.Nil(t, err)
}

func TestMetaRequestersContainerFactory_With4ShardsShouldWork(t *testing.T) {
	t.Parallel()

	noOfShards := 4
	shardCoordinator := mock.NewMultipleShardsCoordinatorMock()
	shardCoordinator.SetNoShards(uint32(noOfShards))
	shardCoordinator.CurrentShard = 1

	args := getArgumentsMeta()
	args.ShardCoordinator = shardCoordinator
	rcf, err := storagerequesterscontainer.NewMetaRequestersContainerFactory(args)
	require.Nil(t, err)

	container, err := rcf.Create()
	require.Nil(t, err)

	numRequestersShardHeadersForMetachain := noOfShards
	numRequesterMetablocks := 1
	numRequestersMiniBlocks := noOfShards + 2
	numRequestersUnsigned := noOfShards + 1
	numRequestersRewards := noOfShards
	numRequestersTxs := noOfShards + 1
	numPeerAuthentication := 1
	numValidatorInfo := 1
	numEquivalentProofs := noOfShards + 1
	totalRequesters := numRequestersShardHeadersForMetachain + numRequesterMetablocks + numRequestersMiniBlocks +
		numRequestersUnsigned + numRequestersTxs + numRequestersRewards + numPeerAuthentication +
		numValidatorInfo + numEquivalentProofs

	assert.Equal(t, totalRequesters, container.Len())
	assert.Equal(t, totalRequesters, container.Len())
}

func getMockStorageConfig() config.StorageConfig {
	return config.StorageConfig{
		Cache: config.CacheConfig{
			Name:     "mock",
			Type:     "LRU",
			Capacity: 1000,
			Shards:   1,
		},
		DB: config.DBConfig{
			FilePath:          "",
			Type:              "MemoryDB",
			BatchDelaySeconds: 1,
			MaxBatchSize:      1,
			MaxOpenFiles:      10,
		},
	}
}

func getArgumentsMeta() storagerequesterscontainer.FactoryArgs {
	return storagerequesterscontainer.FactoryArgs{
		GeneralConfig: config.Config{
			AccountsTrieStorage:     getMockStorageConfig(),
			PeerAccountsTrieStorage: getMockStorageConfig(),
			TrieStorageManagerConfig: config.TrieStorageManagerConfig{
				PruningBufferLen:      255,
				SnapshotsBufferLen:    255,
				SnapshotsGoroutineNum: 2,
			},
			StateTriesConfig: config.StateTriesConfig{
				AccountsStatePruningEnabled: false,
				PeerStatePruningEnabled:     false,
				MaxStateTrieLevelInMemory:   5,
				MaxPeerTrieLevelInMemory:    5,
			},
		},
		ShardIDForTries:          0,
		ChainID:                  "T",
		WorkingDirectory:         "",
		Hasher:                   &hashingMocks.HasherMock{},
		ShardCoordinator:         mock.NewOneShardCoordinatorMock(),
		Messenger:                createMessengerStubForMeta("", ""),
		Store:                    createStoreForMeta(),
		Marshalizer:              &mock.MarshalizerMock{},
		Uint64ByteSliceConverter: &mock.Uint64ByteSliceConverterMock{},
		DataPacker:               &mock.DataPackerStub{},
		ManualEpochStartNotifier: &mock.ManualEpochStartNotifierStub{},
		ChanGracefullyClose:      make(chan endProcess.ArgEndProcess),
		EnableEpochsHandler:      &enableEpochsHandlerMock.EnableEpochsHandlerStub{},
		StateStatsHandler:        disabled.NewStateStatistics(),
	}
}
