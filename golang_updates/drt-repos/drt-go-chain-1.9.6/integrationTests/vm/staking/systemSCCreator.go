package staking

import (
	"bytes"
	"strconv"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	vmcommonMock "github.com/TerraDharitri/drt-go-chain-vm-common/mock"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/epochStart"
	"github.com/TerraDharitri/drt-go-chain/epochStart/metachain"
	epochStartMock "github.com/TerraDharitri/drt-go-chain/epochStart/mock"
	"github.com/TerraDharitri/drt-go-chain/epochStart/notifier"
	"github.com/TerraDharitri/drt-go-chain/factory"
	"github.com/TerraDharitri/drt-go-chain/genesis/process/disabled"
	"github.com/TerraDharitri/drt-go-chain/process"
	metaProcess "github.com/TerraDharitri/drt-go-chain/process/factory/metachain"
	"github.com/TerraDharitri/drt-go-chain/process/peer"
	"github.com/TerraDharitri/drt-go-chain/process/smartContract/builtInFunctions"
	"github.com/TerraDharitri/drt-go-chain/process/smartContract/hooks"
	"github.com/TerraDharitri/drt-go-chain/process/smartContract/hooks/counters"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
	"github.com/TerraDharitri/drt-go-chain/state"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/cryptoMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/genesisMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/guardianMocks"
	"github.com/TerraDharitri/drt-go-chain/vm"
)

func createSystemSCProcessor(
	nc nodesCoordinator.NodesCoordinator,
	coreComponents factory.CoreComponentsHolder,
	stateComponents factory.StateComponentsHandler,
	shardCoordinator sharding.Coordinator,
	maxNodesConfig []config.MaxNodesChangeConfig,
	validatorStatisticsProcessor process.ValidatorStatisticsProcessor,
	systemVM vmcommon.VMExecutionHandler,
	stakingDataProvider epochStart.StakingDataProvider,
) process.EpochStartSystemSCProcessor {
	maxNodesChangeConfigProvider, _ := notifier.NewNodesConfigProvider(
		coreComponents.EpochNotifier(),
		maxNodesConfig,
	)

	auctionCfg := config.SoftAuctionConfig{
		TopUpStep:             "10",
		MinTopUp:              "1",
		MaxTopUp:              "32000000",
		MaxNumberOfIterations: 100000,
	}
	ald, _ := metachain.NewAuctionListDisplayer(metachain.ArgsAuctionListDisplayer{
		TableDisplayHandler:      metachain.NewTableDisplayer(),
		ValidatorPubKeyConverter: &testscommon.PubkeyConverterMock{},
		AddressPubKeyConverter:   &testscommon.PubkeyConverterMock{},
		AuctionConfig:            auctionCfg,
	})

	argsAuctionListSelector := metachain.AuctionListSelectorArgs{
		ShardCoordinator:             shardCoordinator,
		StakingDataProvider:          stakingDataProvider,
		MaxNodesChangeConfigProvider: maxNodesChangeConfigProvider,
		AuctionListDisplayHandler:    ald,
		SoftAuctionConfig:            auctionCfg,
	}
	auctionListSelector, _ := metachain.NewAuctionListSelector(argsAuctionListSelector)

	args := metachain.ArgsNewEpochStartSystemSCProcessing{
		SystemVM:                     systemVM,
		UserAccountsDB:               stateComponents.AccountsAdapter(),
		PeerAccountsDB:               stateComponents.PeerAccounts(),
		Marshalizer:                  coreComponents.InternalMarshalizer(),
		StartRating:                  initialRating,
		ValidatorInfoCreator:         validatorStatisticsProcessor,
		EndOfEpochCallerAddress:      vm.EndOfEpochAddress,
		StakingSCAddress:             vm.StakingSCAddress,
		ChanceComputer:               &epochStartMock.ChanceComputerStub{},
		EpochNotifier:                coreComponents.EpochNotifier(),
		GenesisNodesConfig:           &genesisMocks.NodesSetupStub{},
		StakingDataProvider:          stakingDataProvider,
		NodesConfigProvider:          nc,
		ShardCoordinator:             shardCoordinator,
		DCDTOwnerAddressBytes:        bytes.Repeat([]byte{1}, 32),
		EnableEpochsHandler:          coreComponents.EnableEpochsHandler(),
		MaxNodesChangeConfigProvider: maxNodesChangeConfigProvider,
		AuctionListSelector:          auctionListSelector,
	}

	systemSCProcessor, _ := metachain.NewSystemSCProcessor(args)
	return systemSCProcessor
}

func createStakingDataProvider(
	enableEpochsHandler common.EnableEpochsHandler,
	systemVM vmcommon.VMExecutionHandler,
) epochStart.StakingDataProvider {
	argsStakingDataProvider := metachain.StakingDataProviderArgs{
		EnableEpochsHandler: enableEpochsHandler,
		SystemVM:            systemVM,
		MinNodePrice:        strconv.Itoa(nodePrice),
	}
	stakingSCProvider, _ := metachain.NewStakingDataProvider(argsStakingDataProvider)

	return stakingSCProvider
}

func createValidatorStatisticsProcessor(
	dataComponents factory.DataComponentsHolder,
	coreComponents factory.CoreComponentsHolder,
	nc nodesCoordinator.NodesCoordinator,
	shardCoordinator sharding.Coordinator,
	peerAccounts state.AccountsAdapter,
) process.ValidatorStatisticsProcessor {
	argsValidatorsProcessor := peer.ArgValidatorStatisticsProcessor{
		Marshalizer:                          coreComponents.InternalMarshalizer(),
		NodesCoordinator:                     nc,
		ShardCoordinator:                     shardCoordinator,
		DataPool:                             dataComponents.Datapool(),
		StorageService:                       dataComponents.StorageService(),
		PubkeyConv:                           coreComponents.AddressPubKeyConverter(),
		PeerAdapter:                          peerAccounts,
		Rater:                                coreComponents.Rater(),
		RewardsHandler:                       &epochStartMock.RewardsHandlerStub{},
		NodesSetup:                           &genesisMocks.NodesSetupStub{},
		MaxComputableRounds:                  1,
		MaxConsecutiveRoundsOfRatingDecrease: 2000,
		EnableEpochsHandler:                  coreComponents.EnableEpochsHandler(),
	}
	validatorStatisticsProcessor, _ := peer.NewValidatorStatisticsProcessor(argsValidatorsProcessor)
	return validatorStatisticsProcessor
}

func createBlockChainHook(
	dataComponents factory.DataComponentsHolder,
	coreComponents factory.CoreComponentsHolder,
	accountsAdapter state.AccountsAdapter,
	shardCoordinator sharding.Coordinator,
	gasScheduleNotifier core.GasScheduleNotifier,
) (hooks.ArgBlockChainHook, process.BlockChainHookWithAccountsAdapter) {
	argsBuiltIn := builtInFunctions.ArgsCreateBuiltInFunctionContainer{
		GasSchedule:               gasScheduleNotifier,
		MapDNSAddresses:           make(map[string]struct{}),
		Marshalizer:               coreComponents.InternalMarshalizer(),
		Accounts:                  accountsAdapter,
		ShardCoordinator:          shardCoordinator,
		EpochNotifier:             coreComponents.EpochNotifier(),
		EnableEpochsHandler:       coreComponents.EnableEpochsHandler(),
		AutomaticCrawlerAddresses: [][]byte{core.SystemAccountAddress},
		MaxNumNodesInTransferRole: 1,
		GuardedAccountHandler:     &guardianMocks.GuardedAccountHandlerStub{},
		MapDNSV2Addresses:         make(map[string]struct{}),
	}

	builtInFunctionsContainer, _ := builtInFunctions.CreateBuiltInFunctionsFactory(argsBuiltIn)
	_ = builtInFunctionsContainer.CreateBuiltInFunctionContainer()
	builtInFunctionsContainer.BuiltInFunctionContainer()

	argsHook := hooks.ArgBlockChainHook{
		Accounts:                 accountsAdapter,
		PubkeyConv:               coreComponents.AddressPubKeyConverter(),
		StorageService:           dataComponents.StorageService(),
		BlockChain:               dataComponents.Blockchain(),
		ShardCoordinator:         shardCoordinator,
		Marshalizer:              coreComponents.InternalMarshalizer(),
		Uint64Converter:          coreComponents.Uint64ByteSliceConverter(),
		NFTStorageHandler:        &testscommon.SimpleNFTStorageHandlerStub{},
		BuiltInFunctions:         builtInFunctionsContainer.BuiltInFunctionContainer(),
		DataPool:                 dataComponents.Datapool(),
		CompiledSCPool:           dataComponents.Datapool().SmartContracts(),
		EpochNotifier:            coreComponents.EpochNotifier(),
		GlobalSettingsHandler:    &vmcommonMock.GlobalSettingsHandlerStub{},
		NilCompiledSCStore:       true,
		EnableEpochsHandler:      coreComponents.EnableEpochsHandler(),
		GasSchedule:              gasScheduleNotifier,
		Counter:                  counters.NewDisabledCounter(),
		MissingTrieNodesNotifier: &testscommon.MissingTrieNodesNotifierStub{},
	}

	blockChainHook, _ := hooks.NewBlockChainHookImpl(argsHook)
	return argsHook, blockChainHook
}

func createVMContainerFactory(
	coreComponents factory.CoreComponentsHolder,
	gasScheduleNotifier core.GasScheduleNotifier,
	blockChainHook process.BlockChainHookWithAccountsAdapter,
	argsBlockChainHook hooks.ArgBlockChainHook,
	stateComponents factory.StateComponentsHandler,
	shardCoordinator sharding.Coordinator,
	nc nodesCoordinator.NodesCoordinator,
	maxNumNodes uint32,
) process.VirtualMachinesContainerFactory {
	signVerifer, _ := disabled.NewMessageSignVerifier(&cryptoMocks.KeyGenStub{})

	argsNewVMContainerFactory := metaProcess.ArgsNewVMContainerFactory{
		BlockChainHook:      blockChainHook,
		PubkeyConv:          coreComponents.AddressPubKeyConverter(),
		Economics:           coreComponents.EconomicsData(),
		MessageSignVerifier: signVerifer,
		GasSchedule:         gasScheduleNotifier,
		NodesConfigProvider: &genesisMocks.NodesSetupStub{},
		Hasher:              coreComponents.Hasher(),
		Marshalizer:         coreComponents.InternalMarshalizer(),
		SystemSCConfig: &config.SystemSmartContractsConfig{
			DCDTSystemSCConfig: config.DCDTSystemSCConfig{
				BaseIssuingCost: "1000",
				OwnerAddress:    "aaaaaa",
			},
			GovernanceSystemSCConfig: config.GovernanceSystemSCConfig{
				V1: config.GovernanceSystemSCConfigV1{
					NumNodes:         2000,
					ProposalCost:     "500",
					MinQuorum:        50,
					MinPassThreshold: 10,
					MinVetoThreshold: 10,
				},
				OwnerAddress: "3132333435363738393031323334353637383930313233343536373839303234",
			},
			StakingSystemSCConfig: config.StakingSystemSCConfig{
				GenesisNodePrice:                     strconv.Itoa(nodePrice),
				UnJailValue:                          "10",
				MinStepValue:                         "10",
				MinStakeValue:                        "1",
				UnBondPeriod:                         1,
				NumRoundsWithoutBleed:                1,
				MaximumPercentageToBleed:             1,
				BleedPercentagePerRound:              1,
				MaxNumberOfNodesForStake:             uint64(maxNumNodes),
				ActivateBLSPubKeyMessageVerification: false,
				MinUnstakeTokensValue:                "1",
				StakeLimitPercentage:                 100.0,
				NodeLimitPercentage:                  100.0,
			},
			DelegationManagerSystemSCConfig: config.DelegationManagerSystemSCConfig{
				MinCreationDeposit:  "100",
				MinStakeAmount:      "100",
				ConfigChangeAddress: "3132333435363738393031323334353637383930313233343536373839303234",
			},
			DelegationSystemSCConfig: config.DelegationSystemSCConfig{
				MinServiceFee: 0,
				MaxServiceFee: 100,
			},
			SoftAuctionConfig: config.SoftAuctionConfig{
				TopUpStep:             "10",
				MinTopUp:              "1",
				MaxTopUp:              "32000000",
				MaxNumberOfIterations: 100000,
			},
		},
		ValidatorAccountsDB: stateComponents.PeerAccounts(),
		ChanceComputer:      coreComponents.Rater(),
		EnableEpochsHandler: coreComponents.EnableEpochsHandler(),
		ShardCoordinator:    shardCoordinator,
		NodesCoordinator:    nc,
		UserAccountsDB:      stateComponents.AccountsAdapter(),
		ArgBlockChainHook:   argsBlockChainHook,
	}

	metaVmFactory, _ := metaProcess.NewVMContainerFactory(argsNewVMContainerFactory)
	return metaVmFactory
}
