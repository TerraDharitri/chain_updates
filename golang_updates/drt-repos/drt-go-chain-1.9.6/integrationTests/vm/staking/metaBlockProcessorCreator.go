package staking

import (
	"math/big"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"

	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/epochStart"
	"github.com/TerraDharitri/drt-go-chain/epochStart/metachain"
	"github.com/TerraDharitri/drt-go-chain/factory"
	integrationMocks "github.com/TerraDharitri/drt-go-chain/integrationTests/mock"
	"github.com/TerraDharitri/drt-go-chain/process"
	blproc "github.com/TerraDharitri/drt-go-chain/process/block"
	"github.com/TerraDharitri/drt-go-chain/process/block/bootstrapStorage"
	"github.com/TerraDharitri/drt-go-chain/process/block/postprocess"
	"github.com/TerraDharitri/drt-go-chain/process/block/processedMb"
	"github.com/TerraDharitri/drt-go-chain/process/mock"
	"github.com/TerraDharitri/drt-go-chain/process/scToProtocol"
	"github.com/TerraDharitri/drt-go-chain/process/smartContract"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
	"github.com/TerraDharitri/drt-go-chain/state"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/dblookupext"
	factory2 "github.com/TerraDharitri/drt-go-chain/testscommon/factory"
	"github.com/TerraDharitri/drt-go-chain/testscommon/integrationtests"
	"github.com/TerraDharitri/drt-go-chain/testscommon/outport"
	statusHandlerMock "github.com/TerraDharitri/drt-go-chain/testscommon/statusHandler"
)

func createMetaBlockProcessor(
	nc nodesCoordinator.NodesCoordinator,
	systemSCProcessor process.EpochStartSystemSCProcessor,
	coreComponents factory.CoreComponentsHolder,
	dataComponents factory.DataComponentsHolder,
	bootstrapComponents factory.BootstrapComponentsHolder,
	statusComponents factory.StatusComponentsHolder,
	stateComponents factory.StateComponentsHandler,
	validatorsInfoCreator process.ValidatorStatisticsProcessor,
	blockChainHook process.BlockChainHookHandler,
	metaVMFactory process.VirtualMachinesContainerFactory,
	epochStartHandler process.EpochStartTriggerHandler,
	vmContainer process.VirtualMachinesContainer,
	txCoordinator process.TransactionCoordinator,
) process.BlockProcessor {
	blockTracker := createBlockTracker(
		dataComponents.Blockchain().GetGenesisHeader(),
		bootstrapComponents.ShardCoordinator(),
	)
	epochStartDataCreator := createEpochStartDataCreator(
		coreComponents,
		dataComponents,
		bootstrapComponents.ShardCoordinator(),
		epochStartHandler,
		blockTracker,
	)

	accountsDb := make(map[state.AccountsDbIdentifier]state.AccountsAdapter)
	accountsDb[state.UserAccountsState] = stateComponents.AccountsAdapter()
	accountsDb[state.PeerAccountsState] = stateComponents.PeerAccounts()

	bootStrapStorer, _ := dataComponents.StorageService().GetStorer(dataRetriever.BootstrapUnit)
	bootStorer, _ := bootstrapStorage.NewBootstrapStorer(
		coreComponents.InternalMarshalizer(),
		bootStrapStorer,
	)

	headerValidator := createHeaderValidator(coreComponents)
	valInfoCreator := createValidatorInfoCreator(coreComponents, dataComponents, bootstrapComponents.ShardCoordinator())
	stakingToPeer := createSCToProtocol(coreComponents, stateComponents, dataComponents.Datapool().CurrentBlockTxs())

	args := blproc.ArgMetaProcessor{
		ArgBaseProcessor: blproc.ArgBaseProcessor{
			CoreComponents:      coreComponents,
			DataComponents:      dataComponents,
			BootstrapComponents: bootstrapComponents,
			StatusComponents:    statusComponents,
			StatusCoreComponents: &factory2.StatusCoreComponentsStub{
				AppStatusHandlerField: &statusHandlerMock.AppStatusHandlerStub{},
			},
			AccountsDB:                     accountsDb,
			ForkDetector:                   &integrationMocks.ForkDetectorStub{},
			NodesCoordinator:               nc,
			FeeHandler:                     postprocess.NewFeeAccumulator(),
			RequestHandler:                 &testscommon.RequestHandlerStub{},
			BlockChainHook:                 blockChainHook,
			TxCoordinator:                  txCoordinator,
			EpochStartTrigger:              epochStartHandler,
			HeaderValidator:                headerValidator,
			BootStorer:                     bootStorer,
			BlockTracker:                   blockTracker,
			BlockSizeThrottler:             &mock.BlockSizeThrottlerStub{},
			HistoryRepository:              &dblookupext.HistoryRepositoryStub{},
			VMContainersFactory:            metaVMFactory,
			VmContainer:                    vmContainer,
			GasHandler:                     &mock.GasHandlerMock{},
			ScheduledTxsExecutionHandler:   &testscommon.ScheduledTxsExecutionStub{},
			ScheduledMiniBlocksEnableEpoch: 10000,
			ProcessedMiniBlocksTracker:     processedMb.NewProcessedMiniBlocksTracker(),
			OutportDataProvider:            &outport.OutportDataProviderStub{},
			ReceiptsRepository:             &testscommon.ReceiptsRepositoryStub{},
			ManagedPeersHolder:             &testscommon.ManagedPeersHolderStub{},
			BlockProcessingCutoffHandler:   &testscommon.BlockProcessingCutoffStub{},
			SentSignaturesTracker:          &testscommon.SentSignatureTrackerStub{},
		},
		SCToProtocol:             stakingToPeer,
		PendingMiniBlocksHandler: &mock.PendingMiniBlocksHandlerStub{},
		EpochStartDataCreator:    epochStartDataCreator,
		EpochEconomics:           &mock.EpochEconomicsStub{},
		EpochRewardsCreator: &testscommon.RewardsCreatorStub{
			GetLocalTxCacheCalled: func() epochStart.TransactionCacher {
				return dataComponents.Datapool().CurrentBlockTxs()
			},
		},
		EpochValidatorInfoCreator:    valInfoCreator,
		ValidatorStatisticsProcessor: validatorsInfoCreator,
		EpochSystemSCProcessor:       systemSCProcessor,
	}

	metaProc, _ := blproc.NewMetaProcessor(args)
	return metaProc
}

func createValidatorInfoCreator(
	coreComponents factory.CoreComponentsHolder,
	dataComponents factory.DataComponentsHolder,
	shardCoordinator sharding.Coordinator,
) process.EpochStartValidatorInfoCreator {
	mbStorer, _ := dataComponents.StorageService().GetStorer(dataRetriever.MiniBlockUnit)

	args := metachain.ArgsNewValidatorInfoCreator{
		ShardCoordinator:     shardCoordinator,
		MiniBlockStorage:     mbStorer,
		Hasher:               coreComponents.Hasher(),
		Marshalizer:          coreComponents.InternalMarshalizer(),
		DataPool:             dataComponents.Datapool(),
		EnableEpochsHandler:  coreComponents.EnableEpochsHandler(),
		ValidatorInfoStorage: integrationtests.CreateMemUnit(),
	}

	valInfoCreator, _ := metachain.NewValidatorInfoCreator(args)
	return valInfoCreator
}

func createEpochStartDataCreator(
	coreComponents factory.CoreComponentsHolder,
	dataComponents factory.DataComponentsHolder,
	shardCoordinator sharding.Coordinator,
	epochStartTrigger process.EpochStartTriggerHandler,
	blockTracker process.BlockTracker,
) process.EpochStartDataCreator {
	argsEpochStartDataCreator := metachain.ArgsNewEpochStartData{
		Marshalizer:         coreComponents.InternalMarshalizer(),
		Hasher:              coreComponents.Hasher(),
		Store:               dataComponents.StorageService(),
		DataPool:            dataComponents.Datapool(),
		BlockTracker:        blockTracker,
		ShardCoordinator:    shardCoordinator,
		EpochStartTrigger:   epochStartTrigger,
		RequestHandler:      &testscommon.RequestHandlerStub{},
		GenesisEpoch:        0,
		EnableEpochsHandler: coreComponents.EnableEpochsHandler(),
	}
	epochStartDataCreator, _ := metachain.NewEpochStartData(argsEpochStartDataCreator)
	return epochStartDataCreator
}

func createBlockTracker(
	genesisMetaHeader data.HeaderHandler,
	shardCoordinator sharding.Coordinator,
) process.BlockTracker {
	genesisBlocks := make(map[uint32]data.HeaderHandler)
	for ShardID := uint32(0); ShardID < shardCoordinator.NumberOfShards(); ShardID++ {
		genesisBlocks[ShardID] = createGenesisBlock(ShardID)
	}

	genesisBlocks[core.MetachainShardId] = genesisMetaHeader
	return mock.NewBlockTrackerMock(shardCoordinator, genesisBlocks)
}

func createGenesisBlock(shardID uint32) *block.Header {
	rootHash := []byte("roothash")
	return &block.Header{
		Nonce:           0,
		Round:           0,
		Signature:       rootHash,
		RandSeed:        rootHash,
		PrevRandSeed:    rootHash,
		ShardID:         shardID,
		PubKeysBitmap:   rootHash,
		RootHash:        rootHash,
		PrevHash:        rootHash,
		AccumulatedFees: big.NewInt(0),
		DeveloperFees:   big.NewInt(0),
	}
}

func createGenesisMetaBlock() *block.MetaBlock {
	rootHash := []byte("roothash")
	return &block.MetaBlock{
		Nonce:                  0,
		Round:                  0,
		Signature:              rootHash,
		RandSeed:               rootHash,
		PrevRandSeed:           rootHash,
		PubKeysBitmap:          rootHash,
		RootHash:               rootHash,
		PrevHash:               rootHash,
		AccumulatedFees:        big.NewInt(0),
		DeveloperFees:          big.NewInt(0),
		AccumulatedFeesInEpoch: big.NewInt(0),
		DevFeesInEpoch:         big.NewInt(0),
	}
}

func createHeaderValidator(coreComponents factory.CoreComponentsHolder) epochStart.HeaderValidator {
	argsHeaderValidator := blproc.ArgsHeaderValidator{
		Hasher:              coreComponents.Hasher(),
		Marshalizer:         coreComponents.InternalMarshalizer(),
		EnableEpochsHandler: coreComponents.EnableEpochsHandler(),
	}
	headerValidator, _ := blproc.NewHeaderValidator(argsHeaderValidator)
	return headerValidator
}

func createSCToProtocol(
	coreComponents factory.CoreComponentsHolder,
	stateComponents factory.StateComponentsHandler,
	txCacher dataRetriever.TransactionCacher,
) process.SmartContractToProtocolHandler {
	args := scToProtocol.ArgStakingToPeer{
		PubkeyConv:          coreComponents.AddressPubKeyConverter(),
		Hasher:              coreComponents.Hasher(),
		Marshalizer:         coreComponents.InternalMarshalizer(),
		PeerState:           stateComponents.PeerAccounts(),
		BaseState:           stateComponents.AccountsAdapter(),
		ArgParser:           smartContract.NewArgumentParser(),
		CurrTxs:             txCacher,
		RatingsData:         &mock.RatingsInfoMock{},
		EnableEpochsHandler: coreComponents.EnableEpochsHandler(),
	}
	stakingToPeer, _ := scToProtocol.NewStakingToPeer(args)
	return stakingToPeer
}
