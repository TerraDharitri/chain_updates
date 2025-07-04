package bootstrap

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/common/statistics/disabled"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/epochStart/mock"
	"github.com/TerraDharitri/drt-go-chain/process/block/bootstrapStorage"
	"github.com/TerraDharitri/drt-go-chain/storage"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	dataRetrieverMocks "github.com/TerraDharitri/drt-go-chain/testscommon/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/testscommon/enableEpochsHandlerMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/hashingMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/nodeTypeProviderMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/shardingMocks"
	storageStubs "github.com/TerraDharitri/drt-go-chain/testscommon/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createStorageHandlerArgs() StorageHandlerArgs {
	return StorageHandlerArgs{
		GeneralConfig:                   testscommon.GetGeneralConfig(),
		PreferencesConfig:               config.PreferencesConfig{},
		ShardCoordinator:                &mock.ShardCoordinatorStub{},
		PathManagerHandler:              &testscommon.PathManagerStub{},
		Marshaller:                      &mock.MarshalizerMock{},
		Hasher:                          &hashingMocks.HasherMock{},
		CurrentEpoch:                    0,
		Uint64Converter:                 &mock.Uint64ByteSliceConverterMock{},
		NodeTypeProvider:                &nodeTypeProviderMock.NodeTypeProviderStub{},
		NodesCoordinatorRegistryFactory: &shardingMocks.NodesCoordinatorRegistryFactoryMock{},
		ManagedPeersHolder:              &testscommon.ManagedPeersHolderStub{},
		SnapshotsEnabled:                false,
		NodeProcessingMode:              common.Normal,
		StateStatsHandler:               disabled.NewStateStatistics(),
		RepopulateTokensSupplies:        false,
		ProofsPool:                      &dataRetrieverMocks.ProofsPoolMock{},
		EnableEpochsHandler:             &enableEpochsHandlerMock.EnableEpochsHandlerStub{},
	}
}

func TestNewMetaStorageHandler_InvalidConfigErr(t *testing.T) {
	args := createStorageHandlerArgs()
	args.GeneralConfig = config.Config{}

	mtStrHandler, err := NewMetaStorageHandler(args)
	assert.True(t, check.IfNil(mtStrHandler))
	assert.NotNil(t, err)
}

func TestNewMetaStorageHandler_CreateForMetaErr(t *testing.T) {
	defer func() {
		_ = os.RemoveAll("./Epoch_0")
	}()

	args := createStorageHandlerArgs()
	mtStrHandler, err := NewMetaStorageHandler(args)
	assert.False(t, check.IfNil(mtStrHandler))
	assert.Nil(t, err)
}

func TestMetaStorageHandler_saveLastHeader(t *testing.T) {
	defer func() {
		_ = os.RemoveAll("./Epoch_0")
	}()

	args := createStorageHandlerArgs()
	mtStrHandler, _ := NewMetaStorageHandler(args)
	header := &block.MetaBlock{Nonce: 0}

	headerHash, _ := core.CalculateHash(args.Marshaller, args.Hasher, header)
	expectedBootInfo := bootstrapStorage.BootstrapHeaderInfo{
		ShardId: core.MetachainShardId, Hash: headerHash,
	}

	bootHeaderInfo, err := mtStrHandler.saveLastHeader(header)
	assert.Nil(t, err)
	assert.Equal(t, expectedBootInfo, bootHeaderInfo)
}

func TestMetaStorageHandler_saveLastCrossNotarizedHeaders(t *testing.T) {
	defer func() {
		_ = os.RemoveAll("./Epoch_0")
	}()

	args := createStorageHandlerArgs()
	mtStrHandler, _ := NewMetaStorageHandler(args)

	hdr1 := &block.Header{Nonce: 1}
	hdr2 := &block.Header{Nonce: 2}
	hdrHash1, _ := core.CalculateHash(args.Marshaller, args.Hasher, hdr1)
	hdrHash2, _ := core.CalculateHash(args.Marshaller, args.Hasher, hdr2)

	hdr3 := &block.MetaBlock{
		Nonce: 3,
		EpochStart: block.EpochStart{LastFinalizedHeaders: []block.EpochStartShardData{
			{HeaderHash: hdrHash1}, {HeaderHash: hdrHash2},
		}},
	}

	hdrs := map[string]data.HeaderHandler{string(hdrHash1): hdr1, string(hdrHash2): hdr2}
	crossNotarizedHdrs, err := mtStrHandler.saveLastCrossNotarizedHeaders(hdr3, hdrs)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(crossNotarizedHdrs))
}

func TestMetaStorageHandler_saveTriggerRegistry(t *testing.T) {
	defer func() {
		_ = os.RemoveAll("./Epoch_0")
	}()

	args := createStorageHandlerArgs()
	mtStrHandler, _ := NewMetaStorageHandler(args)

	components := &ComponentsNeededForBootstrap{
		EpochStartMetaBlock: &block.MetaBlock{Nonce: 3},
		PreviousEpochStart:  &block.MetaBlock{Nonce: 2},
	}

	_, err := mtStrHandler.saveTriggerRegistry(components)
	assert.Nil(t, err)
}

func TestMetaStorageHandler_saveDataToStorage(t *testing.T) {
	defer func() {
		_ = os.RemoveAll("./Epoch_0")
	}()

	args := createStorageHandlerArgs()
	mtStrHandler, _ := NewMetaStorageHandler(args)

	components := &ComponentsNeededForBootstrap{
		EpochStartMetaBlock: &block.MetaBlock{Nonce: 3},
		PreviousEpochStart:  &block.MetaBlock{Nonce: 2},
	}

	err := mtStrHandler.SaveDataToStorage(components)
	assert.Nil(t, err)
}

func TestMetaStorageHandler_SaveDataToStorageMissingStorer(t *testing.T) {
	t.Parallel()

	t.Run("missing BootstrapUnit", testMetaWithMissingStorer(dataRetriever.BootstrapUnit, 1))
	t.Run("missing MetaBlockUnit", testMetaWithMissingStorer(dataRetriever.MetaBlockUnit, 1))
	t.Run("missing MetaHdrNonceHashDataUnit", testMetaWithMissingStorer(dataRetriever.MetaHdrNonceHashDataUnit, 1))
	t.Run("missing MetaBlockUnit", testMetaWithMissingStorer(dataRetriever.MetaBlockUnit, 2))                       // saveMetaHdrForEpochTrigger(components.EpochStartMetaBlock)
	t.Run("missing BootstrapUnit", testMetaWithMissingStorer(dataRetriever.BootstrapUnit, 2))                       // saveMetaHdrForEpochTrigger(components.EpochStartMetaBlock)
	t.Run("missing MetaBlockUnit", testMetaWithMissingStorer(dataRetriever.MetaBlockUnit, 3))                       // saveMetaHdrForEpochTrigger(components.PreviousEpochStart)
	t.Run("missing BootstrapUnit", testMetaWithMissingStorer(dataRetriever.BootstrapUnit, 3))                       // saveMetaHdrForEpochTrigger(components.PreviousEpochStart)
	t.Run("missing MetaBlockUnit", testMetaWithMissingStorer(dataRetriever.MetaBlockUnit, 4))                       // saveMetaHdrToStorage(components.PreviousEpochStart)
	t.Run("missing MetaHdrNonceHashDataUnit", testMetaWithMissingStorer(dataRetriever.MetaHdrNonceHashDataUnit, 2)) // saveMetaHdrToStorage(components.PreviousEpochStart)
}

func testMetaWithMissingStorer(missingUnit dataRetriever.UnitType, atCallNumber int) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		defer func() {
			_ = os.RemoveAll("./Epoch_0")
		}()

		args := createStorageHandlerArgs()
		mtStrHandler, _ := NewMetaStorageHandler(args)
		counter := 0
		mtStrHandler.storageService = &storageStubs.ChainStorerStub{
			GetStorerCalled: func(unitType dataRetriever.UnitType) (storage.Storer, error) {
				counter++
				if counter < atCallNumber {
					return &storageStubs.StorerStub{}, nil
				}

				if unitType == missingUnit ||
					strings.Contains(unitType.String(), missingUnit.String()) {
					return nil, fmt.Errorf("%w for %s", storage.ErrKeyNotFound, missingUnit.String())
				}

				return &storageStubs.StorerStub{}, nil
			},
		}
		components := &ComponentsNeededForBootstrap{
			EpochStartMetaBlock: &block.MetaBlock{Nonce: 3},
			PreviousEpochStart:  &block.MetaBlock{Nonce: 2},
		}

		err := mtStrHandler.SaveDataToStorage(components)
		require.NotNil(t, err)
		require.True(t, strings.Contains(err.Error(), storage.ErrKeyNotFound.Error()))
		require.True(t, strings.Contains(err.Error(), missingUnit.String()))
	}
}
