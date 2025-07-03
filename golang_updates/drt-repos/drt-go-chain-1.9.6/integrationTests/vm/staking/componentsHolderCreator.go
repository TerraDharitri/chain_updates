package staking

import (
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/nodetype"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain-core/data/endProcess"
	"github.com/TerraDharitri/drt-go-chain-core/data/typeConverters/uint64ByteSlice"
	"github.com/TerraDharitri/drt-go-chain-core/hashing/sha256"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/common/enablers"
	"github.com/TerraDharitri/drt-go-chain/common/forking"
	"github.com/TerraDharitri/drt-go-chain/common/graceperiod"
	"github.com/TerraDharitri/drt-go-chain/common/statistics/disabled"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever/blockchain"
	"github.com/TerraDharitri/drt-go-chain/epochStart/notifier"
	"github.com/TerraDharitri/drt-go-chain/factory"
	"github.com/TerraDharitri/drt-go-chain/integrationTests"
	integrationMocks "github.com/TerraDharitri/drt-go-chain/integrationTests/mock"
	mockFactory "github.com/TerraDharitri/drt-go-chain/node/mock/factory"
	"github.com/TerraDharitri/drt-go-chain/process/mock"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
	"github.com/TerraDharitri/drt-go-chain/state"
	stateFactory "github.com/TerraDharitri/drt-go-chain/state/factory"
	"github.com/TerraDharitri/drt-go-chain/state/storagePruningManager"
	"github.com/TerraDharitri/drt-go-chain/state/storagePruningManager/evictionWaitingList"
	"github.com/TerraDharitri/drt-go-chain/statusHandler"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	dataRetrieverMock "github.com/TerraDharitri/drt-go-chain/testscommon/dataRetriever"
	notifierMocks "github.com/TerraDharitri/drt-go-chain/testscommon/epochNotifier"
	factoryTests "github.com/TerraDharitri/drt-go-chain/testscommon/factory"
	"github.com/TerraDharitri/drt-go-chain/testscommon/mainFactoryMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/outport"
	"github.com/TerraDharitri/drt-go-chain/testscommon/stakingcommon"
	stateTests "github.com/TerraDharitri/drt-go-chain/testscommon/state"
	statusHandlerMock "github.com/TerraDharitri/drt-go-chain/testscommon/statusHandler"
	"github.com/TerraDharitri/drt-go-chain/trie"
)

const hashSize = 32

func createComponentHolders(numOfShards uint32) (
	factory.CoreComponentsHolder,
	factory.DataComponentsHolder,
	factory.BootstrapComponentsHolder,
	factory.StatusComponentsHolder,
	factory.StateComponentsHandler,
) {
	coreComponents := createCoreComponents()
	statusComponents := createStatusComponents()
	stateComponents := createStateComponents(coreComponents)
	dataComponents := createDataComponents(coreComponents, numOfShards)
	bootstrapComponents := createBootstrapComponents(coreComponents.InternalMarshalizer(), numOfShards)

	return coreComponents, dataComponents, bootstrapComponents, statusComponents, stateComponents
}

func createCoreComponents() factory.CoreComponentsHolder {
	epochNotifier := forking.NewGenericEpochNotifier()
	configEnableEpochs := config.EnableEpochs{
		StakingV4Step1EnableEpoch:          stakingV4Step1EnableEpoch,
		StakingV4Step2EnableEpoch:          stakingV4Step2EnableEpoch,
		StakingV4Step3EnableEpoch:          stakingV4Step3EnableEpoch,
		GovernanceEnableEpoch:              integrationTests.UnreachableEpoch,
		RefactorPeersMiniBlocksEnableEpoch: integrationTests.UnreachableEpoch,
		AndromedaEnableEpoch:               integrationTests.UnreachableEpoch,
	}

	enableEpochsHandler, _ := enablers.NewEnableEpochsHandler(configEnableEpochs, epochNotifier)
	gracePeriod, _ := graceperiod.NewEpochChangeGracePeriod([]config.EpochChangeGracePeriodByEpoch{{EnableEpoch: 0, GracePeriodInRounds: 1}})
	return &integrationMocks.CoreComponentsStub{
		InternalMarshalizerField:           &marshal.GogoProtoMarshalizer{},
		HasherField:                        sha256.NewSha256(),
		Uint64ByteSliceConverterField:      uint64ByteSlice.NewBigEndianConverter(),
		StatusHandlerField:                 statusHandler.NewStatusMetrics(),
		RoundHandlerField:                  &mock.RoundHandlerMock{RoundTimeDuration: time.Second},
		EpochStartNotifierWithConfirmField: notifier.NewEpochStartSubscriptionHandler(),
		EpochNotifierField:                 epochNotifier,
		RaterField:                         &testscommon.RaterMock{Chance: 5},
		AddressPubKeyConverterField:        testscommon.NewPubkeyConverterMock(addressLength),
		EconomicsDataField:                 stakingcommon.CreateEconomicsData(),
		ChanStopNodeProcessField:           endProcess.GetDummyEndProcessChannel(),
		NodeTypeProviderField:              nodetype.NewNodeTypeProvider(core.NodeTypeValidator),
		ProcessStatusHandlerInternal:       statusHandler.NewProcessStatusHandler(),
		EnableEpochsHandlerField:           enableEpochsHandler,
		EnableRoundsHandlerField:           &testscommon.EnableRoundsHandlerStub{},
		RoundNotifierField:                 &notifierMocks.RoundNotifierStub{},
		EpochChangeGracePeriodHandlerField: gracePeriod,
	}
}

func createDataComponents(coreComponents factory.CoreComponentsHolder, numOfShards uint32) factory.DataComponentsHolder {
	genesisBlock := createGenesisMetaBlock()
	genesisBlockHash, _ := coreComponents.InternalMarshalizer().Marshal(genesisBlock)
	genesisBlockHash = coreComponents.Hasher().Compute(string(genesisBlockHash))

	blockChain, _ := blockchain.NewMetaChain(&statusHandlerMock.AppStatusHandlerStub{})
	_ = blockChain.SetGenesisHeader(createGenesisMetaBlock())
	blockChain.SetGenesisHeaderHash(genesisBlockHash)

	chainStorer := dataRetriever.NewChainStorer()
	chainStorer.AddStorer(dataRetriever.BootstrapUnit, integrationTests.CreateMemUnit())
	chainStorer.AddStorer(dataRetriever.MetaHdrNonceHashDataUnit, integrationTests.CreateMemUnit())
	chainStorer.AddStorer(dataRetriever.MetaBlockUnit, integrationTests.CreateMemUnit())
	chainStorer.AddStorer(dataRetriever.MiniBlockUnit, integrationTests.CreateMemUnit())
	chainStorer.AddStorer(dataRetriever.BlockHeaderUnit, integrationTests.CreateMemUnit())
	for i := uint32(0); i < numOfShards; i++ {
		unit := dataRetriever.ShardHdrNonceHashDataUnit + dataRetriever.UnitType(i)
		chainStorer.AddStorer(unit, integrationTests.CreateMemUnit())
	}

	return &mockFactory.DataComponentsMock{
		Store:         chainStorer,
		DataPool:      dataRetrieverMock.NewPoolsHolderMock(),
		BlockChain:    blockChain,
		EconomicsData: coreComponents.EconomicsData(),
	}
}

func createBootstrapComponents(
	marshaller marshal.Marshalizer,
	numOfShards uint32,
) factory.BootstrapComponentsHolder {
	shardCoordinator, _ := sharding.NewMultiShardCoordinator(numOfShards, core.MetachainShardId)
	ncr, _ := nodesCoordinator.NewNodesCoordinatorRegistryFactory(
		marshaller,
		stakingV4Step2EnableEpoch,
	)

	return &mainFactoryMocks.BootstrapComponentsStub{
		ShCoordinator:        shardCoordinator,
		HdrIntegrityVerifier: &mock.HeaderIntegrityVerifierStub{},
		VersionedHdrFactory: &testscommon.VersionedHeaderFactoryStub{
			CreateCalled: func(epoch uint32) data.HeaderHandler {
				return &block.MetaBlock{Epoch: epoch}
			},
		},
		NodesCoordinatorRegistryFactoryField: ncr,
	}
}

func createStatusComponents() factory.StatusComponentsHolder {
	return &integrationMocks.StatusComponentsStub{
		Outport:                  &outport.OutportStub{},
		SoftwareVersionCheck:     &integrationMocks.SoftwareVersionCheckerMock{},
		ManagedPeersMonitorField: &testscommon.ManagedPeersMonitorStub{},
	}
}

func createStateComponents(coreComponents factory.CoreComponentsHolder) factory.StateComponentsHandler {
	tsmArgs := getNewTrieStorageManagerArgs(coreComponents)
	tsm, _ := trie.CreateTrieStorageManager(tsmArgs, trie.StorageManagerOptions{})
	trieFactoryManager, _ := trie.NewTrieStorageManagerWithoutPruning(tsm)

	argsAccCreator := stateFactory.ArgsAccountCreator{
		Hasher:              coreComponents.Hasher(),
		Marshaller:          coreComponents.InternalMarshalizer(),
		EnableEpochsHandler: coreComponents.EnableEpochsHandler(),
	}

	accCreator, _ := stateFactory.NewAccountCreator(argsAccCreator)

	userAccountsDB := createAccountsDB(coreComponents, accCreator, trieFactoryManager)
	peerAccountsDB := createAccountsDB(coreComponents, stateFactory.NewPeerAccountCreator(), trieFactoryManager)

	_ = userAccountsDB.SetSyncer(&mock.AccountsDBSyncerStub{})
	_ = peerAccountsDB.SetSyncer(&mock.AccountsDBSyncerStub{})

	return &factoryTests.StateComponentsMock{
		PeersAcc: peerAccountsDB,
		Accounts: userAccountsDB,
	}
}

func getNewTrieStorageManagerArgs(coreComponents factory.CoreComponentsHolder) trie.NewTrieStorageManagerArgs {
	return trie.NewTrieStorageManagerArgs{
		MainStorer:     testscommon.CreateMemUnit(),
		Marshalizer:    coreComponents.InternalMarshalizer(),
		Hasher:         coreComponents.Hasher(),
		GeneralConfig:  config.TrieStorageManagerConfig{SnapshotsGoroutineNum: 1},
		IdleProvider:   &testscommon.ProcessStatusHandlerStub{},
		Identifier:     "id",
		StatsCollector: disabled.NewStateStatistics(),
	}
}

func createAccountsDB(
	coreComponents factory.CoreComponentsHolder,
	accountFactory state.AccountFactory,
	trieStorageManager common.StorageManager,
) *state.AccountsDB {
	tr, _ := trie.NewTrie(
		trieStorageManager,
		coreComponents.InternalMarshalizer(),
		coreComponents.Hasher(),
		coreComponents.EnableEpochsHandler(),
		5,
	)

	argsEvictionWaitingList := evictionWaitingList.MemoryEvictionWaitingListArgs{
		RootHashesSize: 10,
		HashesSize:     hashSize,
	}
	ewl, _ := evictionWaitingList.NewMemoryEvictionWaitingList(argsEvictionWaitingList)
	spm, _ := storagePruningManager.NewStoragePruningManager(ewl, 10)
	argsAccountsDb := state.ArgsAccountsDB{
		Trie:                  tr,
		Hasher:                coreComponents.Hasher(),
		Marshaller:            coreComponents.InternalMarshalizer(),
		AccountFactory:        accountFactory,
		StoragePruningManager: spm,
		AddressConverter:      coreComponents.AddressPubKeyConverter(),
		SnapshotsManager:      &stateTests.SnapshotsManagerStub{},
	}
	adb, _ := state.NewAccountsDB(argsAccountsDb)
	return adb
}
