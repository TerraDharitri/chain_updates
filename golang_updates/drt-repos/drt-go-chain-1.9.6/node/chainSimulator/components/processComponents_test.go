package components

import (
	"math/big"
	"sync"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	coreData "github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/endProcess"
	"github.com/TerraDharitri/drt-go-chain-core/hashing/blake2b"
	"github.com/TerraDharitri/drt-go-chain-core/hashing/keccak"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/stretchr/testify/require"

	commonFactory "github.com/TerraDharitri/drt-go-chain/common/factory"
	"github.com/TerraDharitri/drt-go-chain/common/graceperiod"
	disabledStatistics "github.com/TerraDharitri/drt-go-chain/common/statistics/disabled"
	"github.com/TerraDharitri/drt-go-chain/config"
	retriever "github.com/TerraDharitri/drt-go-chain/dataRetriever"
	mockFactory "github.com/TerraDharitri/drt-go-chain/factory/mock"
	"github.com/TerraDharitri/drt-go-chain/integrationTests/mock"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	chainStorage "github.com/TerraDharitri/drt-go-chain/storage"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/bootstrapMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/chainParameters"
	"github.com/TerraDharitri/drt-go-chain/testscommon/components"
	"github.com/TerraDharitri/drt-go-chain/testscommon/cryptoMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/testscommon/economicsmocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/enableEpochsHandlerMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/epochNotifier"
	"github.com/TerraDharitri/drt-go-chain/testscommon/factory"
	"github.com/TerraDharitri/drt-go-chain/testscommon/genericMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/guardianMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/mainFactoryMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/nodeTypeProviderMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/outport"
	"github.com/TerraDharitri/drt-go-chain/testscommon/p2pmocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/shardingMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/statusHandler"
	"github.com/TerraDharitri/drt-go-chain/testscommon/storage"
	updateMocks "github.com/TerraDharitri/drt-go-chain/update/mock"
)

const testingProtocolSustainabilityAddress = "drt1932eft30w753xyvme8d49qejgkjc09n5e49w4mwdjtm0neld797spn6u9l"

var (
	addrPubKeyConv, _ = commonFactory.NewPubkeyConverter(config.PubkeyConfig{
		Length:          32,
		Type:            "bech32",
		SignatureLength: 0,
		Hrp:             "drt",
	})
	valPubKeyConv, _ = commonFactory.NewPubkeyConverter(config.PubkeyConfig{
		Length:          96,
		Type:            "hex",
		SignatureLength: 48,
	})
)

func createArgsProcessComponentsHolder() ArgsProcessComponentsHolder {
	var nodesConfig config.NodesConfig
	_ = core.LoadJsonFile(&nodesConfig, "../../../integrationTests/factory/testdata/nodesSetup.json")
	nodesSetup, _ := sharding.NewNodesSetup(nodesConfig, &chainParameters.ChainParametersHolderMock{}, addrPubKeyConv, valPubKeyConv, 3)
	gracePeriod, _ := graceperiod.NewEpochChangeGracePeriod([]config.EpochChangeGracePeriodByEpoch{{EnableEpoch: 0, GracePeriodInRounds: 1}})
	args := ArgsProcessComponentsHolder{
		Config: testscommon.GetGeneralConfig(),
		EpochConfig: config.EpochConfig{
			GasSchedule: config.GasScheduleConfig{
				GasScheduleByEpochs: []config.GasScheduleByEpochs{
					{
						StartEpoch: 0,
						FileName:   "../../../cmd/node/config/gasSchedules/gasScheduleV8.toml",
					},
				},
			},
		},
		RoundConfig:    testscommon.GetDefaultRoundsConfig(),
		PrefsConfig:    config.Preferences{},
		ImportDBConfig: config.ImportDbConfig{},
		FlagsConfig: config.ContextFlagsConfig{
			Version: "v1.0.0",
		},
		NodesCoordinator: &shardingMocks.NodesCoordinatorStub{},
		SystemSCConfig: config.SystemSmartContractsConfig{
			DCDTSystemSCConfig: config.DCDTSystemSCConfig{
				BaseIssuingCost: "1000",
				OwnerAddress:    "drt1fpkcgel4gcmh8zqqdt043yfcn5tyx8373kg6q2qmkxzu4dqamc0snh8ehx",
			},
			GovernanceSystemSCConfig: config.GovernanceSystemSCConfig{
				V1: config.GovernanceSystemSCConfigV1{
					ProposalCost:     "500",
					NumNodes:         100,
					MinQuorum:        50,
					MinPassThreshold: 50,
					MinVetoThreshold: 50,
				},
				Active: config.GovernanceSystemSCConfigActive{
					ProposalCost:     "500",
					MinQuorum:        0.5,
					MinPassThreshold: 0.5,
					MinVetoThreshold: 0.5,
				},
				OwnerAddress: "drt1vxy22x0fj4zv6hktmydg8vpfh6euv02cz4yg0aaws6rrad5a5awq4up8y3",
			},
			StakingSystemSCConfig: config.StakingSystemSCConfig{
				GenesisNodePrice:                     "2500000000000000000000",
				MinStakeValue:                        "1",
				UnJailValue:                          "1",
				MinStepValue:                         "1",
				UnBondPeriod:                         0,
				NumRoundsWithoutBleed:                0,
				MaximumPercentageToBleed:             0,
				BleedPercentagePerRound:              0,
				MaxNumberOfNodesForStake:             10,
				ActivateBLSPubKeyMessageVerification: false,
				MinUnstakeTokensValue:                "1",
				NodeLimitPercentage:                  0.1,
				StakeLimitPercentage:                 1,
				UnBondPeriodInEpochs:                 10,
			},
			DelegationManagerSystemSCConfig: config.DelegationManagerSystemSCConfig{
				MinCreationDeposit:  "100",
				MinStakeAmount:      "100",
				ConfigChangeAddress: "drt1vxy22x0fj4zv6hktmydg8vpfh6euv02cz4yg0aaws6rrad5a5awq4up8y3",
			},
			DelegationSystemSCConfig: config.DelegationSystemSCConfig{
				MinServiceFee: 0,
				MaxServiceFee: 100,
			},
		},
		DataComponents: &mock.DataComponentsStub{
			DataPool: dataRetriever.NewPoolsHolderMock(),
			BlockChain: &testscommon.ChainHandlerStub{
				GetGenesisHeaderHashCalled: func() []byte {
					return []byte("genesis hash")
				},
				GetGenesisHeaderCalled: func() coreData.HeaderHandler {
					return &testscommon.HeaderHandlerStub{}
				},
			},
			MbProvider: &mock.MiniBlocksProviderStub{},
			Store:      genericMocks.NewChainStorerMock(0),
		},
		CoreComponents: &mockFactory.CoreComponentsMock{
			IntMarsh:            &marshal.GogoProtoMarshalizer{},
			TxMarsh:             &marshal.JsonMarshalizer{},
			UInt64ByteSliceConv: &mock.Uint64ByteSliceConverterMock{},
			AddrPubKeyConv:      addrPubKeyConv,
			ValPubKeyConv:       valPubKeyConv,
			NodesConfig:         nodesSetup,
			EpochChangeNotifier: &epochNotifier.EpochNotifierStub{},
			EconomicsHandler: &economicsmocks.EconomicsHandlerMock{
				ProtocolSustainabilityAddressCalled: func() string {
					return testingProtocolSustainabilityAddress
				},
				GenesisTotalSupplyCalled: func() *big.Int {
					return big.NewInt(0).Mul(big.NewInt(1000000000000000000), big.NewInt(20000000))
				},
			},
			Hash:                               blake2b.NewBlake2b(),
			TxVersionCheckHandler:              &testscommon.TxVersionCheckerStub{},
			RatingHandler:                      &testscommon.RaterMock{},
			EnableEpochsHandlerField:           &enableEpochsHandlerMock.EnableEpochsHandlerStub{},
			EnableRoundsHandlerField:           &testscommon.EnableRoundsHandlerStub{},
			EpochNotifierWithConfirm:           &updateMocks.EpochStartNotifierStub{},
			RoundHandlerField:                  &testscommon.RoundHandlerMock{},
			RoundChangeNotifier:                &epochNotifier.RoundNotifierStub{},
			ChanStopProcess:                    make(chan endProcess.ArgEndProcess, 1),
			TxSignHasherField:                  keccak.NewKeccak(),
			HardforkTriggerPubKeyField:         []byte("hardfork pub key"),
			WasmVMChangeLockerInternal:         &sync.RWMutex{},
			NodeTypeProviderField:              &nodeTypeProviderMock.NodeTypeProviderStub{},
			RatingsConfig:                      &testscommon.RatingsInfoMock{},
			PathHdl:                            &testscommon.PathManagerStub{},
			ProcessStatusHandlerInternal:       &testscommon.ProcessStatusHandlerStub{},
			EpochChangeGracePeriodHandlerField: gracePeriod,
		},
		CryptoComponents: &mock.CryptoComponentsStub{
			BlKeyGen: &cryptoMocks.KeyGenStub{},
			BlockSig: &cryptoMocks.SingleSignerStub{},
			MultiSigContainer: &cryptoMocks.MultiSignerContainerMock{
				MultiSigner: &cryptoMocks.MultisignerMock{},
			},
			PrivKey:                 &cryptoMocks.PrivateKeyStub{},
			PubKey:                  &cryptoMocks.PublicKeyStub{},
			PubKeyString:            "pub key string",
			PubKeyBytes:             []byte("pub key bytes"),
			TxKeyGen:                &cryptoMocks.KeyGenStub{},
			TxSig:                   &cryptoMocks.SingleSignerStub{},
			PeerSignHandler:         &cryptoMocks.PeerSignatureHandlerStub{},
			MsgSigVerifier:          &testscommon.MessageSignVerifierMock{},
			ManagedPeersHolderField: &testscommon.ManagedPeersHolderStub{},
			KeysHandlerField:        &testscommon.KeysHandlerStub{},
		},
		NetworkComponents: &mock.NetworkComponentsStub{
			Messenger:                        &p2pmocks.MessengerStub{},
			FullArchiveNetworkMessengerField: &p2pmocks.MessengerStub{},
			InputAntiFlood:                   &mock.P2PAntifloodHandlerStub{},
			OutputAntiFlood:                  &mock.P2PAntifloodHandlerStub{},
			PreferredPeersHolder:             &p2pmocks.PeersHolderStub{},
			PeersRatingHandlerField:          &p2pmocks.PeersRatingHandlerStub{},
			FullArchivePreferredPeersHolder:  &p2pmocks.PeersHolderStub{},
		},
		BootstrapComponents: &mainFactoryMocks.BootstrapComponentsStub{
			ShCoordinator:              mock.NewMultiShardsCoordinatorMock(2),
			BootstrapParams:            &bootstrapMocks.BootstrapParamsHandlerMock{},
			HdrIntegrityVerifier:       &mock.HeaderIntegrityVerifierStub{},
			GuardedAccountHandlerField: &guardianMocks.GuardedAccountHandlerStub{},
			VersionedHdrFactory:        &testscommon.VersionedHeaderFactoryStub{},
		},
		StatusComponents: &mock.StatusComponentsStub{
			Outport: &outport.OutportStub{},
		},
		StatusCoreComponents: &factory.StatusCoreComponentsStub{
			AppStatusHandlerField:  &statusHandler.AppStatusHandlerStub{},
			StateStatsHandlerField: disabledStatistics.NewStateStatistics(),
		},
		EconomicsConfig: config.EconomicsConfig{
			GlobalSettings: config.GlobalSettings{
				GenesisTotalSupply:          "20000000000000000000000000",
				MinimumInflation:            0,
				GenesisMintingSenderAddress: "drt17rc0pu8s7rc0pu8s7rc0pu8s7rc0pu8s7rc0pu8s7rc0pu8s7rcqa2qg80",
				YearSettings: []*config.YearSetting{
					{
						Year:             0,
						MaximumInflation: 0.01,
					},
				},
			},
		},
		ConfigurationPathsHolder: config.ConfigurationPathsHolder{
			Genesis:        "../../../integrationTests/factory/testdata/genesis.json",
			SmartContracts: "../../../integrationTests/factory/testdata/genesisSmartContracts.json",
			Nodes:          "../../../integrationTests/factory/testdata/genesis.json",
		},
	}

	args.StateComponents = components.GetStateComponents(args.CoreComponents, args.StatusCoreComponents)
	return args
}

func TestCreateProcessComponents(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	t.Run("should work", func(t *testing.T) {
		comp, err := CreateProcessComponents(createArgsProcessComponentsHolder())
		require.NoError(t, err)
		require.NotNil(t, comp)

		require.Nil(t, comp.Create())
		require.Nil(t, comp.Close())
	})
	t.Run("NewImportStartHandler failure should error", func(t *testing.T) {
		t.Parallel()

		args := createArgsProcessComponentsHolder()
		args.FlagsConfig.Version = ""
		comp, err := CreateProcessComponents(args)
		require.Error(t, err)
		require.Nil(t, comp)
	})
	t.Run("total supply conversion failure should error", func(t *testing.T) {
		t.Parallel()

		args := createArgsProcessComponentsHolder()
		args.EconomicsConfig.GlobalSettings.GenesisTotalSupply = "invalid number"
		comp, err := CreateProcessComponents(args)
		require.Error(t, err)
		require.Nil(t, comp)
	})
	t.Run("NewAccountsParser failure should error", func(t *testing.T) {
		t.Parallel()

		args := createArgsProcessComponentsHolder()
		args.ConfigurationPathsHolder.Genesis = ""
		comp, err := CreateProcessComponents(args)
		require.Error(t, err)
		require.Nil(t, comp)
	})
	t.Run("NewSmartContractsParser failure should error", func(t *testing.T) {
		t.Parallel()

		args := createArgsProcessComponentsHolder()
		args.ConfigurationPathsHolder.SmartContracts = ""
		comp, err := CreateProcessComponents(args)
		require.Error(t, err)
		require.Nil(t, comp)
	})
	t.Run("NewHistoryRepositoryFactory failure should error", func(t *testing.T) {
		t.Parallel()

		args := createArgsProcessComponentsHolder()
		dataMock, ok := args.DataComponents.(*mock.DataComponentsStub)
		require.True(t, ok)
		dataMock.Store = nil
		comp, err := CreateProcessComponents(args)
		require.Error(t, err)
		require.Nil(t, comp)
	})
	t.Run("historyRepositoryFactory.Create failure should error", func(t *testing.T) {
		t.Parallel()

		args := createArgsProcessComponentsHolder()
		args.Config.DbLookupExtensions.Enabled = true
		dataMock, ok := args.DataComponents.(*mock.DataComponentsStub)
		require.True(t, ok)
		dataMock.Store = &storage.ChainStorerStub{
			GetStorerCalled: func(unitType retriever.UnitType) (chainStorage.Storer, error) {
				if unitType == retriever.DCDTSuppliesUnit {
					return nil, expectedErr
				}
				return &storage.StorerStub{}, nil
			},
		}
		comp, err := CreateProcessComponents(args)
		require.Error(t, err)
		require.Nil(t, comp)
	})
	t.Run("NewGasScheduleNotifier failure should error", func(t *testing.T) {
		t.Parallel()

		args := createArgsProcessComponentsHolder()
		args.EpochConfig.GasSchedule = config.GasScheduleConfig{}
		comp, err := CreateProcessComponents(args)
		require.Error(t, err)
		require.Nil(t, comp)
	})
	t.Run("NewProcessComponentsFactory failure should error", func(t *testing.T) {
		t.Parallel()

		args := createArgsProcessComponentsHolder()
		dataMock, ok := args.DataComponents.(*mock.DataComponentsStub)
		require.True(t, ok)
		dataMock.BlockChain = nil
		comp, err := CreateProcessComponents(args)
		require.Error(t, err)
		require.Nil(t, comp)
	})
	t.Run("managedProcessComponents.Create failure should error", func(t *testing.T) {
		t.Parallel()

		args := createArgsProcessComponentsHolder()
		args.NodesCoordinator = nil
		comp, err := CreateProcessComponents(args)
		require.Error(t, err)
		require.Nil(t, comp)
	})
}

func TestProcessComponentsHolder_IsInterfaceNil(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	var comp *processComponentsHolder
	require.True(t, comp.IsInterfaceNil())

	comp, _ = CreateProcessComponents(createArgsProcessComponentsHolder())
	require.False(t, comp.IsInterfaceNil())
	require.Nil(t, comp.Close())
}

func TestProcessComponentsHolder_Getters(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	comp, err := CreateProcessComponents(createArgsProcessComponentsHolder())
	require.NoError(t, err)

	require.NotNil(t, comp.SentSignaturesTracker())
	require.NotNil(t, comp.NodesCoordinator())
	require.NotNil(t, comp.ShardCoordinator())
	require.NotNil(t, comp.InterceptorsContainer())
	require.NotNil(t, comp.FullArchiveInterceptorsContainer())
	require.NotNil(t, comp.ResolversContainer())
	require.NotNil(t, comp.RequestersFinder())
	require.NotNil(t, comp.RoundHandler())
	require.NotNil(t, comp.EpochStartTrigger())
	require.NotNil(t, comp.EpochStartNotifier())
	require.NotNil(t, comp.ForkDetector())
	require.NotNil(t, comp.BlockProcessor())
	require.NotNil(t, comp.BlackListHandler())
	require.NotNil(t, comp.BootStorer())
	require.NotNil(t, comp.HeaderSigVerifier())
	require.NotNil(t, comp.HeaderIntegrityVerifier())
	require.NotNil(t, comp.ValidatorsStatistics())
	require.NotNil(t, comp.ValidatorsProvider())
	require.NotNil(t, comp.BlockTracker())
	require.NotNil(t, comp.PendingMiniBlocksHandler())
	require.NotNil(t, comp.RequestHandler())
	require.NotNil(t, comp.TxLogsProcessor())
	require.NotNil(t, comp.HeaderConstructionValidator())
	require.NotNil(t, comp.PeerShardMapper())
	require.NotNil(t, comp.FullArchivePeerShardMapper())
	require.NotNil(t, comp.FallbackHeaderValidator())
	require.NotNil(t, comp.APITransactionEvaluator())
	require.NotNil(t, comp.WhiteListHandler())
	require.NotNil(t, comp.WhiteListerVerifiedTxs())
	require.NotNil(t, comp.HistoryRepository())
	require.NotNil(t, comp.ImportStartHandler())
	require.NotNil(t, comp.RequestedItemsHandler())
	require.NotNil(t, comp.NodeRedundancyHandler())
	require.NotNil(t, comp.CurrentEpochProvider())
	require.NotNil(t, comp.ScheduledTxsExecutionHandler())
	require.NotNil(t, comp.TxsSenderHandler())
	require.NotNil(t, comp.HardforkTrigger())
	require.NotNil(t, comp.ProcessedMiniBlocksTracker())
	require.NotNil(t, comp.DCDTDataStorageHandlerForAPI())
	require.NotNil(t, comp.AccountsParser())
	require.NotNil(t, comp.ReceiptsRepository())
	require.NotNil(t, comp.EpochSystemSCProcessor())
	require.Nil(t, comp.CheckSubcomponents())
	require.Empty(t, comp.String())

	require.Nil(t, comp.Close())
}
