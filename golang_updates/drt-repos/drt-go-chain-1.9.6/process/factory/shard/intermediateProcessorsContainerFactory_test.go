package shard_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/factory/shard"
	"github.com/TerraDharitri/drt-go-chain/process/mock"
	"github.com/TerraDharitri/drt-go-chain/storage"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/cache"
	txExecOrderStub "github.com/TerraDharitri/drt-go-chain/testscommon/common"
	dataRetrieverMock "github.com/TerraDharitri/drt-go-chain/testscommon/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/testscommon/economicsmocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/enableEpochsHandlerMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/hashingMocks"
	storageStubs "github.com/TerraDharitri/drt-go-chain/testscommon/storage"
)

func createDataPools() dataRetriever.PoolsHolder {
	pools := dataRetrieverMock.NewPoolsHolderStub()
	pools.TransactionsCalled = func() dataRetriever.ShardedDataCacherNotifier {
		return testscommon.NewShardedDataStub()
	}
	pools.HeadersCalled = func() dataRetriever.HeadersPool {
		return &mock.HeadersCacherStub{}
	}
	pools.MiniBlocksCalled = func() storage.Cacher {
		return cache.NewCacherStub()
	}
	pools.PeerChangesBlocksCalled = func() storage.Cacher {
		return cache.NewCacherStub()
	}
	pools.MetaBlocksCalled = func() storage.Cacher {
		return cache.NewCacherStub()
	}
	pools.UnsignedTransactionsCalled = func() dataRetriever.ShardedDataCacherNotifier {
		return testscommon.NewShardedDataStub()
	}
	pools.RewardTransactionsCalled = func() dataRetriever.ShardedDataCacherNotifier {
		return testscommon.NewShardedDataStub()
	}
	pools.TrieNodesCalled = func() storage.Cacher {
		return cache.NewCacherStub()
	}
	pools.CurrBlockTxsCalled = func() dataRetriever.TransactionCacher {
		return &mock.TxForCurrentBlockStub{}
	}
	return pools
}

func createMockPubkeyConverter() *testscommon.PubkeyConverterMock {
	return testscommon.NewPubkeyConverterMock(32)
}

func createMockArgsNewIntermediateProcessorsFactory() shard.ArgsNewIntermediateProcessorsContainerFactory {
	args := shard.ArgsNewIntermediateProcessorsContainerFactory{
		Hasher:                  &hashingMocks.HasherMock{},
		Marshalizer:             &mock.MarshalizerMock{},
		ShardCoordinator:        mock.NewMultiShardsCoordinatorMock(5),
		PubkeyConverter:         createMockPubkeyConverter(),
		Store:                   &storageStubs.ChainStorerStub{},
		PoolsHolder:             createDataPools(),
		EconomicsFee:            &economicsmocks.EconomicsHandlerMock{},
		EnableEpochsHandler:     enableEpochsHandlerMock.NewEnableEpochsHandlerStub(common.KeepExecOrderOnCreatedSCRsFlag),
		TxExecutionOrderHandler: &txExecOrderStub.TxExecutionOrderHandlerStub{},
	}
	return args
}

func TestNewIntermediateProcessorsContainerFactory_NilShardCoord(t *testing.T) {
	t.Parallel()

	args := createMockArgsNewIntermediateProcessorsFactory()
	args.ShardCoordinator = nil
	ipcf, err := shard.NewIntermediateProcessorsContainerFactory(args)

	assert.Nil(t, ipcf)
	assert.Equal(t, process.ErrNilShardCoordinator, err)
}

func TestNewIntermediateProcessorsContainerFactory_NilMarshalizer(t *testing.T) {
	t.Parallel()

	args := createMockArgsNewIntermediateProcessorsFactory()
	args.Marshalizer = nil
	ipcf, err := shard.NewIntermediateProcessorsContainerFactory(args)

	assert.Nil(t, ipcf)
	assert.Equal(t, process.ErrNilMarshalizer, err)
}

func TestNewIntermediateProcessorsContainerFactory_NilHasher(t *testing.T) {
	t.Parallel()

	args := createMockArgsNewIntermediateProcessorsFactory()
	args.Hasher = nil
	ipcf, err := shard.NewIntermediateProcessorsContainerFactory(args)

	assert.Nil(t, ipcf)
	assert.Equal(t, process.ErrNilHasher, err)
}

func TestNewIntermediateProcessorsContainerFactory_NilAdrConv(t *testing.T) {
	t.Parallel()

	args := createMockArgsNewIntermediateProcessorsFactory()
	args.PubkeyConverter = nil
	ipcf, err := shard.NewIntermediateProcessorsContainerFactory(args)

	assert.Nil(t, ipcf)
	assert.Equal(t, process.ErrNilPubkeyConverter, err)
}

func TestNewIntermediateProcessorsContainerFactory_NilStorer(t *testing.T) {
	t.Parallel()

	args := createMockArgsNewIntermediateProcessorsFactory()
	args.Store = nil
	ipcf, err := shard.NewIntermediateProcessorsContainerFactory(args)

	assert.Nil(t, ipcf)
	assert.Equal(t, process.ErrNilStorage, err)
}

func TestNewIntermediateProcessorsContainerFactory_NilPoolsHolder(t *testing.T) {
	t.Parallel()

	args := createMockArgsNewIntermediateProcessorsFactory()
	args.PoolsHolder = nil
	ipcf, err := shard.NewIntermediateProcessorsContainerFactory(args)

	assert.Nil(t, ipcf)
	assert.Equal(t, process.ErrNilPoolsHolder, err)
}

func TestNewIntermediateProcessorsContainerFactory_NilEconomicsFeeHandler(t *testing.T) {
	t.Parallel()

	args := createMockArgsNewIntermediateProcessorsFactory()
	args.EconomicsFee = nil
	ipcf, err := shard.NewIntermediateProcessorsContainerFactory(args)

	assert.Nil(t, ipcf)
	assert.Equal(t, process.ErrNilEconomicsFeeHandler, err)
}

func TestNewIntermediateProcessorsContainerFactory_NilEnableEpochsHandler(t *testing.T) {
	t.Parallel()

	args := createMockArgsNewIntermediateProcessorsFactory()
	args.EnableEpochsHandler = nil
	ipcf, err := shard.NewIntermediateProcessorsContainerFactory(args)

	assert.Nil(t, ipcf)
	assert.Equal(t, process.ErrNilEnableEpochsHandler, err)
}

func TestNewIntermediateProcessorsContainerFactory(t *testing.T) {
	t.Parallel()

	args := createMockArgsNewIntermediateProcessorsFactory()
	ipcf, err := shard.NewIntermediateProcessorsContainerFactory(args)

	assert.Nil(t, err)
	assert.NotNil(t, ipcf)
	assert.False(t, ipcf.IsInterfaceNil())
}

func TestIntermediateProcessorsContainerFactory_Create(t *testing.T) {
	t.Parallel()

	args := createMockArgsNewIntermediateProcessorsFactory()
	ipcf, err := shard.NewIntermediateProcessorsContainerFactory(args)
	assert.Nil(t, err)
	assert.NotNil(t, ipcf)

	container, err := ipcf.Create()
	assert.Nil(t, err)
	assert.Equal(t, 3, container.Len())
}
