package integrationTests

import (
	"fmt"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data"

	"github.com/TerraDharitri/drt-go-chain/common/enablers"
	"github.com/TerraDharitri/drt-go-chain/common/forking"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever/provider"
	"github.com/TerraDharitri/drt-go-chain/integrationTests/mock"
	"github.com/TerraDharitri/drt-go-chain/outport/disabled"
	"github.com/TerraDharitri/drt-go-chain/process/block"
	"github.com/TerraDharitri/drt-go-chain/process/block/bootstrapStorage"
	"github.com/TerraDharitri/drt-go-chain/process/sync"
	"github.com/TerraDharitri/drt-go-chain/state"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/dblookupext"
	"github.com/TerraDharitri/drt-go-chain/testscommon/enableEpochsHandlerMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/factory"
	"github.com/TerraDharitri/drt-go-chain/testscommon/outport"
	statusHandlerMock "github.com/TerraDharitri/drt-go-chain/testscommon/statusHandler"
)

func (tpn *TestProcessorNode) addGenesisBlocksIntoStorage() {
	for shardId, header := range tpn.GenesisBlocks {
		buffHeader, _ := TestMarshalizer.Marshal(header)
		headerHash := TestHasher.Compute(string(buffHeader))

		if shardId == core.MetachainShardId {
			metablockStorer, _ := tpn.Storage.GetStorer(dataRetriever.MetaBlockUnit)
			_ = metablockStorer.Put(headerHash, buffHeader)
		} else {
			shardblockStorer, _ := tpn.Storage.GetStorer(dataRetriever.BlockHeaderUnit)
			_ = shardblockStorer.Put(headerHash, buffHeader)
		}
	}
}

func (tpn *TestProcessorNode) initBlockProcessorWithSync() {
	var err error

	accountsDb := make(map[state.AccountsDbIdentifier]state.AccountsAdapter)
	accountsDb[state.UserAccountsState] = tpn.AccntState
	accountsDb[state.PeerAccountsState] = tpn.PeerState

	if tpn.EpochNotifier == nil {
		tpn.EpochNotifier = forking.NewGenericEpochNotifier()
	}
	if tpn.EnableEpochsHandler == nil {
		tpn.EnableEpochsHandler, _ = enablers.NewEnableEpochsHandler(CreateEnableEpochsConfig(), tpn.EpochNotifier)
	}
	coreComponents := GetDefaultCoreComponents(tpn.EnableEpochsHandler, tpn.EpochNotifier)
	coreComponents.InternalMarshalizerField = TestMarshalizer
	coreComponents.HasherField = TestHasher
	coreComponents.Uint64ByteSliceConverterField = TestUint64Converter
	coreComponents.EpochNotifierField = tpn.EpochNotifier
	coreComponents.RoundNotifierField = tpn.RoundNotifier

	dataComponents := GetDefaultDataComponents()
	dataComponents.Store = tpn.Storage
	dataComponents.DataPool = tpn.DataPool
	dataComponents.BlockChain = tpn.BlockChain

	bootstrapComponents := getDefaultBootstrapComponents(tpn.ShardCoordinator, tpn.EnableEpochsHandler)
	bootstrapComponents.HdrIntegrityVerifier = tpn.HeaderIntegrityVerifier

	statusComponents := GetDefaultStatusComponents()

	statusCoreComponents := &factory.StatusCoreComponentsStub{
		AppStatusHandlerField: &statusHandlerMock.AppStatusHandlerStub{},
	}

	argumentsBase := block.ArgBaseProcessor{
		CoreComponents:       coreComponents,
		DataComponents:       dataComponents,
		BootstrapComponents:  bootstrapComponents,
		StatusComponents:     statusComponents,
		StatusCoreComponents: statusCoreComponents,
		Config:               config.Config{},
		AccountsDB:           accountsDb,
		ForkDetector:         nil,
		NodesCoordinator:     tpn.NodesCoordinator,
		FeeHandler:           tpn.FeeAccumulator,
		RequestHandler:       tpn.RequestHandler,
		BlockChainHook:       &testscommon.BlockChainHookStub{},
		EpochStartTrigger:    &mock.EpochStartTriggerStub{},
		HeaderValidator:      tpn.HeaderValidator,
		BootStorer: &mock.BoostrapStorerMock{
			PutCalled: func(round int64, bootData bootstrapStorage.BootstrapData) error {
				return nil
			},
		},
		BlockTracker:                 tpn.BlockTracker,
		BlockSizeThrottler:           TestBlockSizeThrottler,
		HistoryRepository:            tpn.HistoryRepository,
		GasHandler:                   tpn.GasHandler,
		ScheduledTxsExecutionHandler: &testscommon.ScheduledTxsExecutionStub{},
		ProcessedMiniBlocksTracker:   &testscommon.ProcessedMiniBlocksTrackerStub{},
		ReceiptsRepository:           &testscommon.ReceiptsRepositoryStub{},
		OutportDataProvider:          &outport.OutportDataProviderStub{},
		BlockProcessingCutoffHandler: &testscommon.BlockProcessingCutoffStub{},
		ManagedPeersHolder:           &testscommon.ManagedPeersHolderStub{},
		SentSignaturesTracker:        &testscommon.SentSignatureTrackerStub{},
	}

	if tpn.ShardCoordinator.SelfId() == core.MetachainShardId {
		tpn.ForkDetector, _ = sync.NewMetaForkDetector(
			tpn.RoundHandler,
			tpn.BlockBlackListHandler,
			tpn.BlockTracker,
			0,
			tpn.EnableEpochsHandler,
			tpn.DataPool.Proofs())
		argumentsBase.ForkDetector = tpn.ForkDetector
		argumentsBase.TxCoordinator = &mock.TransactionCoordinatorMock{}
		arguments := block.ArgMetaProcessor{
			ArgBaseProcessor:          argumentsBase,
			SCToProtocol:              &mock.SCToProtocolStub{},
			PendingMiniBlocksHandler:  &mock.PendingMiniBlocksHandlerStub{},
			EpochStartDataCreator:     &mock.EpochStartDataCreatorStub{},
			EpochEconomics:            &mock.EpochEconomicsStub{},
			EpochRewardsCreator:       &testscommon.RewardsCreatorStub{},
			EpochValidatorInfoCreator: &testscommon.EpochValidatorInfoCreatorStub{},
			ValidatorStatisticsProcessor: &testscommon.ValidatorStatisticsProcessorStub{
				UpdatePeerStateCalled: func(header data.MetaHeaderHandler) ([]byte, error) {
					return []byte("validator stats root hash"), nil
				},
			},
			EpochSystemSCProcessor: &testscommon.EpochStartSystemSCStub{},
		}

		tpn.BlockProcessor, err = block.NewMetaProcessor(arguments)
	} else {
		tpn.ForkDetector, _ = sync.NewShardForkDetector(
			tpn.RoundHandler,
			tpn.BlockBlackListHandler,
			tpn.BlockTracker,
			0,
			tpn.EnableEpochsHandler,
			tpn.DataPool.Proofs())
		argumentsBase.ForkDetector = tpn.ForkDetector
		argumentsBase.BlockChainHook = tpn.BlockchainHook
		argumentsBase.TxCoordinator = tpn.TxCoordinator
		argumentsBase.ScheduledTxsExecutionHandler = &testscommon.ScheduledTxsExecutionStub{}
		arguments := block.ArgShardProcessor{
			ArgBaseProcessor: argumentsBase,
		}

		tpn.BlockProcessor, err = block.NewShardProcessor(arguments)
	}

	if err != nil {
		panic(fmt.Sprintf("Error creating blockprocessor: %s", err.Error()))
	}
}

func (tpn *TestProcessorNode) createShardBootstrapper() (TestBootstrapper, error) {
	argsBaseBootstrapper := sync.ArgBaseBootstrapper{
		PoolsHolder:                  tpn.DataPool,
		Store:                        tpn.Storage,
		ChainHandler:                 tpn.BlockChain,
		RoundHandler:                 tpn.RoundHandler,
		BlockProcessor:               tpn.BlockProcessor,
		WaitTime:                     tpn.RoundHandler.TimeDuration(),
		Hasher:                       TestHasher,
		Marshalizer:                  TestMarshalizer,
		ForkDetector:                 tpn.ForkDetector,
		RequestHandler:               tpn.RequestHandler,
		ShardCoordinator:             tpn.ShardCoordinator,
		Accounts:                     tpn.AccntState,
		BlackListHandler:             tpn.BlockBlackListHandler,
		NetworkWatcher:               tpn.MainMessenger,
		BootStorer:                   tpn.BootstrapStorer,
		StorageBootstrapper:          tpn.StorageBootstrapper,
		EpochHandler:                 tpn.EpochStartTrigger,
		MiniblocksProvider:           tpn.MiniblocksProvider,
		Uint64Converter:              TestUint64Converter,
		AppStatusHandler:             TestAppStatusHandler,
		OutportHandler:               disabled.NewDisabledOutport(),
		AccountsDBSyncer:             &mock.AccountsDBSyncerStub{},
		CurrentEpochProvider:         &testscommon.CurrentEpochProviderStub{},
		IsInImportMode:               false,
		HistoryRepo:                  &dblookupext.HistoryRepositoryStub{},
		ScheduledTxsExecutionHandler: &testscommon.ScheduledTxsExecutionStub{},
		ProcessWaitTime:              tpn.RoundHandler.TimeDuration(),
		RepopulateTokensSupplies:     false,
		EnableEpochsHandler:          tpn.EnableEpochsHandler,
	}

	argsShardBootstrapper := sync.ArgShardBootstrapper{
		ArgBaseBootstrapper: argsBaseBootstrapper,
	}

	bootstrap, err := sync.NewShardBootstrap(argsShardBootstrapper)
	if err != nil {
		return nil, err
	}

	return &sync.TestShardBootstrap{
		ShardBootstrap: bootstrap,
	}, nil
}

func (tpn *TestProcessorNode) createMetaChainBootstrapper() (TestBootstrapper, error) {
	argsBaseBootstrapper := sync.ArgBaseBootstrapper{
		PoolsHolder:                  tpn.DataPool,
		Store:                        tpn.Storage,
		ChainHandler:                 tpn.BlockChain,
		RoundHandler:                 tpn.RoundHandler,
		BlockProcessor:               tpn.BlockProcessor,
		WaitTime:                     tpn.RoundHandler.TimeDuration(),
		Hasher:                       TestHasher,
		Marshalizer:                  TestMarshalizer,
		ForkDetector:                 tpn.ForkDetector,
		RequestHandler:               tpn.RequestHandler,
		ShardCoordinator:             tpn.ShardCoordinator,
		Accounts:                     tpn.AccntState,
		BlackListHandler:             tpn.BlockBlackListHandler,
		NetworkWatcher:               tpn.MainMessenger,
		BootStorer:                   tpn.BootstrapStorer,
		StorageBootstrapper:          tpn.StorageBootstrapper,
		EpochHandler:                 tpn.EpochStartTrigger,
		MiniblocksProvider:           tpn.MiniblocksProvider,
		Uint64Converter:              TestUint64Converter,
		AppStatusHandler:             TestAppStatusHandler,
		OutportHandler:               disabled.NewDisabledOutport(),
		AccountsDBSyncer:             &mock.AccountsDBSyncerStub{},
		CurrentEpochProvider:         &testscommon.CurrentEpochProviderStub{},
		IsInImportMode:               false,
		HistoryRepo:                  &dblookupext.HistoryRepositoryStub{},
		ScheduledTxsExecutionHandler: &testscommon.ScheduledTxsExecutionStub{},
		ProcessWaitTime:              tpn.RoundHandler.TimeDuration(),
		RepopulateTokensSupplies:     false,
		EnableEpochsHandler:          &enableEpochsHandlerMock.EnableEpochsHandlerStub{},
	}

	argsMetaBootstrapper := sync.ArgMetaBootstrapper{
		ArgBaseBootstrapper:         argsBaseBootstrapper,
		EpochBootstrapper:           tpn.EpochStartTrigger,
		ValidatorAccountsDB:         tpn.PeerState,
		ValidatorStatisticsDBSyncer: &mock.AccountsDBSyncerStub{},
	}

	bootstrap, err := sync.NewMetaBootstrap(argsMetaBootstrapper)
	if err != nil {
		return nil, err
	}

	return &sync.TestMetaBootstrap{
		MetaBootstrap: bootstrap,
	}, nil
}

func (tpn *TestProcessorNode) initBootstrapper() {
	tpn.createMiniblocksProvider()

	if tpn.ShardCoordinator.SelfId() < tpn.ShardCoordinator.NumberOfShards() {
		tpn.Bootstrapper, _ = tpn.createShardBootstrapper()
	} else {
		tpn.Bootstrapper, _ = tpn.createMetaChainBootstrapper()
	}
}

func (tpn *TestProcessorNode) createMiniblocksProvider() {
	storer, _ := tpn.Storage.GetStorer(dataRetriever.MiniBlockUnit)
	arg := provider.ArgMiniBlockProvider{
		MiniBlockPool:    tpn.DataPool.MiniBlocks(),
		MiniBlockStorage: storer,
		Marshalizer:      TestMarshalizer,
	}

	miniblockGetter, err := provider.NewMiniBlockProvider(arg)
	log.LogIfError(err)

	tpn.MiniblocksProvider = miniblockGetter
}
