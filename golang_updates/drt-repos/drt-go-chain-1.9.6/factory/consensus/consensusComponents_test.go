package consensus_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	crypto "github.com/TerraDharitri/drt-go-chain-crypto"
	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/consensus"
	retriever "github.com/TerraDharitri/drt-go-chain/dataRetriever"
	errorsDrt "github.com/TerraDharitri/drt-go-chain/errors"
	consensusComp "github.com/TerraDharitri/drt-go-chain/factory/consensus"
	"github.com/TerraDharitri/drt-go-chain/factory/mock"
	testsMocks "github.com/TerraDharitri/drt-go-chain/integrationTests/mock"
	"github.com/TerraDharitri/drt-go-chain/p2p"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
	"github.com/TerraDharitri/drt-go-chain/storage"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/cache"
	consensusMocks "github.com/TerraDharitri/drt-go-chain/testscommon/consensus"
	"github.com/TerraDharitri/drt-go-chain/testscommon/cryptoMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/dataRetriever"
	dataRetrieverMocks "github.com/TerraDharitri/drt-go-chain/testscommon/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/testscommon/dblookupext"
	"github.com/TerraDharitri/drt-go-chain/testscommon/enableEpochsHandlerMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/epochNotifier"
	factoryMocks "github.com/TerraDharitri/drt-go-chain/testscommon/factory"
	"github.com/TerraDharitri/drt-go-chain/testscommon/genericMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/genesisMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/marshallerMock"
	outportMocks "github.com/TerraDharitri/drt-go-chain/testscommon/outport"
	"github.com/TerraDharitri/drt-go-chain/testscommon/p2pmocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/shardingMocks"
	stateMocks "github.com/TerraDharitri/drt-go-chain/testscommon/state"
	"github.com/TerraDharitri/drt-go-chain/testscommon/statusHandler"
	"github.com/TerraDharitri/drt-go-chain/testscommon/storageManager"
	"github.com/TerraDharitri/drt-go-chain/update"
)

func createMockConsensusComponentsFactoryArgs() consensusComp.ConsensusComponentsFactoryArgs {
	return consensusComp.ConsensusComponentsFactoryArgs{
		Config:              testscommon.GetGeneralConfig(),
		BootstrapRoundIndex: 0,
		CoreComponents: &mock.CoreComponentsMock{
			IntMarsh: &marshallerMock.MarshalizerStub{},
			Hash: &testscommon.HasherStub{
				SizeCalled: func() int {
					return 1
				},
			},
			UInt64ByteSliceConv: &testsMocks.Uint64ByteSliceConverterMock{},
			AddrPubKeyConv:      &testscommon.PubkeyConverterStub{},
			WatchdogTimer:       &testscommon.WatchdogMock{},
			AlarmSch:            &testscommon.AlarmSchedulerStub{},
			NtpSyncTimer:        &testscommon.SyncTimerStub{},
			GenesisBlockTime:    time.Time{},
			NodesConfig: &genesisMocks.NodesSetupStub{
				GetShardConsensusGroupSizeCalled: func() uint32 {
					return 2
				},
				GetMetaConsensusGroupSizeCalled: func() uint32 {
					return 2
				},
			},
			EpochChangeNotifier:      &epochNotifier.EpochNotifierStub{},
			StartTime:                time.Time{},
			EnableEpochsHandlerField: &enableEpochsHandlerMock.EnableEpochsHandlerStub{},
		},
		NetworkComponents: &testsMocks.NetworkComponentsStub{
			Messenger:      &p2pmocks.MessengerStub{},
			InputAntiFlood: &testsMocks.P2PAntifloodHandlerStub{},
			PeerHonesty:    &testscommon.PeerHonestyHandlerStub{},
		},
		CryptoComponents: &testsMocks.CryptoComponentsStub{
			PrivKey:         &cryptoMocks.PrivateKeyStub{},
			PubKey:          &cryptoMocks.PublicKeyStub{},
			PubKeyString:    "pub key string",
			PeerSignHandler: &testsMocks.PeerSignatureHandler{},
			MultiSigContainer: &cryptoMocks.MultiSignerContainerMock{
				MultiSigner: &cryptoMocks.MultisignerMock{},
			},
			BlKeyGen:         &cryptoMocks.KeyGenStub{},
			BlockSig:         &cryptoMocks.SingleSignerStub{},
			KeysHandlerField: &testscommon.KeysHandlerStub{},
			SigHandler:       &consensusMocks.SigningHandlerStub{},
		},
		DataComponents: &testsMocks.DataComponentsStub{
			DataPool: &dataRetriever.PoolsHolderStub{
				MiniBlocksCalled: func() storage.Cacher {
					return &cache.CacherStub{}
				},
				TrieNodesCalled: func() storage.Cacher {
					return &cache.CacherStub{}
				},
				HeadersCalled: func() retriever.HeadersPool {
					return &testsMocks.HeadersCacherStub{}
				},
				ProofsCalled: func() retriever.ProofsPool {
					return &dataRetrieverMocks.ProofsPoolMock{}
				},
			},
			BlockChain: &testscommon.ChainHandlerStub{
				GetGenesisHeaderHashCalled: func() []byte {
					return []byte("genesis hash")
				},
				GetGenesisHeaderCalled: func() data.HeaderHandler {
					return &testscommon.HeaderHandlerStub{}
				},
			},
			MbProvider: &testsMocks.MiniBlocksProviderStub{},
			Store:      &genericMocks.ChainStorerMock{},
		},
		ProcessComponents: &testsMocks.ProcessComponentsStub{
			EpochTrigger:                  &testsMocks.EpochStartTriggerStub{},
			EpochNotifier:                 &testsMocks.EpochStartNotifierStub{},
			NodesCoord:                    &shardingMocks.NodesCoordinatorMock{},
			NodeRedundancyHandlerInternal: &testsMocks.RedundancyHandlerStub{},
			HardforkTriggerField:          &testscommon.HardforkTriggerStub{},
			ReqHandler:                    &testscommon.RequestHandlerStub{},
			MainPeerMapper:                &testsMocks.PeerShardMapperStub{},
			FullArchivePeerMapper:         &testsMocks.PeerShardMapperStub{},
			ShardCoord:                    testscommon.NewMultiShardsCoordinatorMock(2),
			RoundHandlerField: &testscommon.RoundHandlerMock{
				TimeDurationCalled: func() time.Duration {
					return time.Second
				},
			},
			BootSore:                             &mock.BootstrapStorerMock{},
			ForkDetect:                           &mock.ForkDetectorMock{},
			BlockProcess:                         &testscommon.BlockProcessorStub{},
			BlockTrack:                           &mock.BlockTrackerStub{},
			ScheduledTxsExecutionHandlerInternal: &testscommon.ScheduledTxsExecutionStub{},
			ProcessedMiniBlocksTrackerInternal:   &testscommon.ProcessedMiniBlocksTrackerStub{},
			PendingMiniBlocksHdl:                 &mock.PendingMiniBlocksHandlerStub{},
			BlackListHdl:                         &testscommon.TimeCacheStub{},
			CurrentEpochProviderInternal:         &testsMocks.CurrentNetworkEpochProviderStub{},
			HistoryRepositoryInternal:            &dblookupext.HistoryRepositoryStub{},
			IntContainer:                         &testscommon.InterceptorsContainerStub{},
			HeaderSigVerif:                       &consensusMocks.HeaderSigVerifierMock{},
			HeaderIntegrVerif:                    &mock.HeaderIntegrityVerifierStub{},
			FallbackHdrValidator:                 &testscommon.FallBackHeaderValidatorStub{},
			SentSignaturesTrackerInternal:        &testscommon.SentSignatureTrackerStub{},
		},
		StateComponents: &factoryMocks.StateComponentsMock{
			StorageManagers: map[string]common.StorageManager{
				retriever.UserAccountsUnit.String(): &storageManager.StorageManagerStub{},
				retriever.PeerAccountsUnit.String(): &storageManager.StorageManagerStub{},
			},
			Accounts:             &stateMocks.AccountsStub{},
			PeersAcc:             &stateMocks.AccountsStub{},
			MissingNodesNotifier: &testscommon.MissingTrieNodesNotifierStub{},
		},
		StatusComponents: &testsMocks.StatusComponentsStub{
			Outport: &outportMocks.OutportStub{},
		},
		StatusCoreComponents: &factoryMocks.StatusCoreComponentsStub{
			AppStatusHandlerField: &statusHandler.AppStatusHandlerStub{},
		},
		ScheduledProcessor:    &consensusMocks.ScheduledProcessorStub{},
		IsInImportMode:        false,
		ShouldDisableWatchdog: false,
	}
}

func TestNewConsensusComponentsFactory(t *testing.T) {
	t.Parallel()

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		ccf, err := consensusComp.NewConsensusComponentsFactory(createMockConsensusComponentsFactoryArgs())

		require.NotNil(t, ccf)
		require.Nil(t, err)
	})
	t.Run("nil CoreComponents should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		args.CoreComponents = nil
		ccf, err := consensusComp.NewConsensusComponentsFactory(args)

		require.Nil(t, ccf)
		require.Equal(t, errorsDrt.ErrNilCoreComponentsHolder, err)
	})
	t.Run("nil GenesisNodesSetup should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		args.CoreComponents = &mock.CoreComponentsMock{
			NodesConfig: nil,
		}
		ccf, err := consensusComp.NewConsensusComponentsFactory(args)

		require.Nil(t, ccf)
		require.Equal(t, errorsDrt.ErrNilGenesisNodesSetupHandler, err)
	})
	t.Run("nil DataComponents should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		args.DataComponents = nil
		ccf, err := consensusComp.NewConsensusComponentsFactory(args)

		require.Nil(t, ccf)
		require.Equal(t, errorsDrt.ErrNilDataComponentsHolder, err)
	})
	t.Run("nil Datapool should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		args.DataComponents = &testsMocks.DataComponentsStub{
			DataPool: nil,
		}
		ccf, err := consensusComp.NewConsensusComponentsFactory(args)

		require.Nil(t, ccf)
		require.Equal(t, errorsDrt.ErrNilDataPoolsHolder, err)
	})
	t.Run("nil BlockChain should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		args.DataComponents = &testsMocks.DataComponentsStub{
			DataPool:   &dataRetriever.PoolsHolderStub{},
			BlockChain: nil,
		}
		ccf, err := consensusComp.NewConsensusComponentsFactory(args)

		require.Nil(t, ccf)
		require.Equal(t, errorsDrt.ErrNilBlockChainHandler, err)
	})
	t.Run("nil CryptoComponents should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		args.CryptoComponents = nil
		ccf, err := consensusComp.NewConsensusComponentsFactory(args)

		require.Nil(t, ccf)
		require.Equal(t, errorsDrt.ErrNilCryptoComponentsHolder, err)
	})
	t.Run("nil PublicKey should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		args.CryptoComponents = &testsMocks.CryptoComponentsStub{
			PubKey: nil,
		}
		ccf, err := consensusComp.NewConsensusComponentsFactory(args)

		require.Nil(t, ccf)
		require.Equal(t, errorsDrt.ErrNilPublicKey, err)
	})
	t.Run("nil PrivateKey should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		args.CryptoComponents = &testsMocks.CryptoComponentsStub{
			PubKey:  &cryptoMocks.PublicKeyStub{},
			PrivKey: nil,
		}
		ccf, err := consensusComp.NewConsensusComponentsFactory(args)

		require.Nil(t, ccf)
		require.Equal(t, errorsDrt.ErrNilPrivateKey, err)
	})
	t.Run("nil NetworkComponents should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		args.NetworkComponents = nil
		ccf, err := consensusComp.NewConsensusComponentsFactory(args)

		require.Nil(t, ccf)
		require.Equal(t, errorsDrt.ErrNilNetworkComponentsHolder, err)
	})
	t.Run("nil Messenger should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		args.NetworkComponents = &testsMocks.NetworkComponentsStub{
			Messenger: nil,
		}
		ccf, err := consensusComp.NewConsensusComponentsFactory(args)

		require.Nil(t, ccf)
		require.Equal(t, errorsDrt.ErrNilMessenger, err)
	})
	t.Run("nil ProcessComponents should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		args.ProcessComponents = nil
		ccf, err := consensusComp.NewConsensusComponentsFactory(args)

		require.Nil(t, ccf)
		require.Equal(t, errorsDrt.ErrNilProcessComponentsHolder, err)
	})
	t.Run("nil NodesCoordinator should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		args.ProcessComponents = &testsMocks.ProcessComponentsStub{
			NodesCoord: nil,
		}
		ccf, err := consensusComp.NewConsensusComponentsFactory(args)

		require.Nil(t, ccf)
		require.Equal(t, errorsDrt.ErrNilNodesCoordinator, err)
	})
	t.Run("nil ShardCoordinator should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		args.ProcessComponents = &testsMocks.ProcessComponentsStub{
			NodesCoord: &shardingMocks.NodesCoordinatorMock{},
			ShardCoord: nil,
		}
		ccf, err := consensusComp.NewConsensusComponentsFactory(args)

		require.Nil(t, ccf)
		require.Equal(t, errorsDrt.ErrNilShardCoordinator, err)
	})
	t.Run("nil RoundHandler should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		args.ProcessComponents = &testsMocks.ProcessComponentsStub{
			NodesCoord:        &shardingMocks.NodesCoordinatorMock{},
			ShardCoord:        &testscommon.ShardsCoordinatorMock{},
			RoundHandlerField: nil,
		}
		ccf, err := consensusComp.NewConsensusComponentsFactory(args)

		require.Nil(t, ccf)
		require.Equal(t, errorsDrt.ErrNilRoundHandler, err)
	})
	t.Run("nil HardforkTrigger should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		args.ProcessComponents = &testsMocks.ProcessComponentsStub{
			NodesCoord:           &shardingMocks.NodesCoordinatorMock{},
			ShardCoord:           &testscommon.ShardsCoordinatorMock{},
			RoundHandlerField:    &testscommon.RoundHandlerMock{},
			HardforkTriggerField: nil,
		}
		ccf, err := consensusComp.NewConsensusComponentsFactory(args)

		require.Nil(t, ccf)
		require.Equal(t, errorsDrt.ErrNilHardforkTrigger, err)
	})
	t.Run("nil StateComponents should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		args.StateComponents = nil
		ccf, err := consensusComp.NewConsensusComponentsFactory(args)

		require.Nil(t, ccf)
		require.Equal(t, errorsDrt.ErrNilStateComponentsHolder, err)
	})
	t.Run("nil StatusComponents should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		args.StatusComponents = nil
		ccf, err := consensusComp.NewConsensusComponentsFactory(args)

		require.Nil(t, ccf)
		require.Equal(t, errorsDrt.ErrNilStatusComponentsHolder, err)
	})
	t.Run("nil OutportHandler should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		args.StatusComponents = &testsMocks.StatusComponentsStub{
			Outport: nil,
		}
		ccf, err := consensusComp.NewConsensusComponentsFactory(args)

		require.Nil(t, ccf)
		require.Equal(t, errorsDrt.ErrNilOutportHandler, err)
	})
	t.Run("nil ScheduledProcessor should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		args.ScheduledProcessor = nil
		ccf, err := consensusComp.NewConsensusComponentsFactory(args)

		require.Nil(t, ccf)
		require.Equal(t, errorsDrt.ErrNilScheduledProcessor, err)
	})
	t.Run("nil StatusCoreComponents should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		args.StatusCoreComponents = nil
		ccf, err := consensusComp.NewConsensusComponentsFactory(args)

		require.Nil(t, ccf)
		require.Equal(t, errorsDrt.ErrNilStatusCoreComponents, err)
	})
}

func TestConsensusComponentsFactory_Create(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("expected error")
	t.Run("invalid shard id should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		processCompStub, ok := args.ProcessComponents.(*testsMocks.ProcessComponentsStub)
		require.True(t, ok)
		processCompStub.ShardCoord = &testscommon.ShardsCoordinatorMock{
			SelfIDCalled: func() uint32 {
				return 5
			},
			NoShards: 2,
		}
		ccf, _ := consensusComp.NewConsensusComponentsFactory(args)
		require.NotNil(t, ccf)

		cc, err := ccf.Create()
		require.Equal(t, sharding.ErrShardIdOutOfRange, err)
		require.Nil(t, cc)
	})
	t.Run("genesis block not initialized should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		processCompStub, ok := args.ProcessComponents.(*testsMocks.ProcessComponentsStub)
		require.True(t, ok)
		processCompStub.ShardCoord = &testscommon.ShardsCoordinatorMock{
			SelfIDCalled: func() uint32 {
				return core.MetachainShardId // coverage
			},
			NoShards: 2,
		}

		dataCompStub, ok := args.DataComponents.(*testsMocks.DataComponentsStub)
		require.True(t, ok)
		dataCompStub.BlockChain = &testscommon.ChainHandlerStub{
			GetGenesisHeaderHashCalled: func() []byte {
				return []byte("")
			},
		}
		ccf, _ := consensusComp.NewConsensusComponentsFactory(args)
		require.NotNil(t, ccf)

		cc, err := ccf.Create()
		require.Equal(t, errorsDrt.ErrGenesisBlockNotInitialized, err)
		require.Nil(t, cc)
	})
	t.Run("createChronology fails should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		processCompStub, ok := args.ProcessComponents.(*testsMocks.ProcessComponentsStub)
		require.True(t, ok)
		cnt := 0
		processCompStub.RoundHandlerCalled = func() consensus.RoundHandler {
			cnt++
			if cnt > 1 {
				return nil
			}
			return &testscommon.RoundHandlerMock{}
		}

		args.IsInImportMode = true        // coverage
		args.ShouldDisableWatchdog = true // coverage
		statusCompStub, ok := args.StatusComponents.(*testsMocks.StatusComponentsStub)
		require.True(t, ok)
		statusCompStub.Outport = &outportMocks.OutportStub{
			HasDriversCalled: func() bool {
				return true // coverage
			},
		}
		ccf, _ := consensusComp.NewConsensusComponentsFactory(args)
		require.NotNil(t, ccf)

		cc, err := ccf.Create()
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), "roundHandler"))
		require.Nil(t, cc)
	})
	t.Run("createBootstrapper fails due to nil shard coordinator should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		processCompStub, ok := args.ProcessComponents.(*testsMocks.ProcessComponentsStub)
		require.True(t, ok)
		cnt := 0
		processCompStub.ShardCoordinatorCalled = func() sharding.Coordinator {
			cnt++
			if cnt > 1 {
				return nil // createBootstrapper fails
			}
			return testscommon.NewMultiShardsCoordinatorMock(2)
		}
		ccf, _ := consensusComp.NewConsensusComponentsFactory(args)
		require.NotNil(t, ccf)

		cc, err := ccf.Create()
		require.Equal(t, errorsDrt.ErrNilShardCoordinator, err)
		require.Nil(t, cc)
	})
	t.Run("createBootstrapper fails due to invalid shard coordinator should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		processCompStub, ok := args.ProcessComponents.(*testsMocks.ProcessComponentsStub)
		require.True(t, ok)
		cnt := 0
		shardC := testscommon.NewMultiShardsCoordinatorMock(2)
		processCompStub.ShardCoordinatorCalled = func() sharding.Coordinator {
			cnt++
			if cnt > 1 {
				shardC.SelfIDCalled = func() uint32 {
					return shardC.NoShards + 1 // createBootstrapper returns ErrShardIdOutOfRange
				}
				return shardC
			}
			return shardC
		}
		ccf, _ := consensusComp.NewConsensusComponentsFactory(args)
		require.NotNil(t, ccf)

		cc, err := ccf.Create()
		require.Equal(t, sharding.ErrShardIdOutOfRange, err)
		require.Nil(t, cc)
	})
	t.Run("createUserAccountsSyncer fails due to missing UserAccountTrie should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		stateCompStub, ok := args.StateComponents.(*factoryMocks.StateComponentsMock)
		require.True(t, ok)
		stateCompStub.StorageManagers = make(map[string]common.StorageManager) // missing UserAccountTrie
		ccf, _ := consensusComp.NewConsensusComponentsFactory(args)
		require.NotNil(t, ccf)

		cc, err := ccf.Create()
		require.Equal(t, errorsDrt.ErrNilTrieStorageManager, err)
		require.Nil(t, cc)
	})
	t.Run("createUserAccountsSyncer fails due to invalid NumConcurrentTrieSyncers should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		args.Config.TrieSync.NumConcurrentTrieSyncers = 0
		ccf, _ := consensusComp.NewConsensusComponentsFactory(args)
		require.NotNil(t, ccf)

		cc, err := ccf.Create()
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), "value is not positive"))
		require.Nil(t, cc)
	})
	t.Run("createUserAccountsSyncer fails due to missing UserAccountTrie should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		stateCompStub, ok := args.StateComponents.(*factoryMocks.StateComponentsMock)
		require.True(t, ok)
		stateCompStub.StorageManagers = make(map[string]common.StorageManager) // missing UserAccountTrie
		processCompStub, ok := args.ProcessComponents.(*testsMocks.ProcessComponentsStub)
		require.True(t, ok)
		processCompStub.ShardCoordinatorCalled = func() sharding.Coordinator {
			shardC := testscommon.NewMultiShardsCoordinatorMock(2)
			shardC.CurrentShard = core.MetachainShardId
			return shardC
		}
		ccf, _ := consensusComp.NewConsensusComponentsFactory(args)
		require.NotNil(t, ccf)

		cc, err := ccf.Create()
		require.Equal(t, errorsDrt.ErrNilTrieStorageManager, err)
		require.Nil(t, cc)
	})
	t.Run("createValidatorAccountsSyncer fails due to missing PeerAccountTrie should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		stateCompStub, ok := args.StateComponents.(*factoryMocks.StateComponentsMock)
		require.True(t, ok)
		stateCompStub.StorageManagers = map[string]common.StorageManager{
			retriever.UserAccountsUnit.String(): &storageManager.StorageManagerStub{},
		} // missing PeerAccountTrie
		processCompStub, ok := args.ProcessComponents.(*testsMocks.ProcessComponentsStub)
		require.True(t, ok)
		processCompStub.ShardCoordinatorCalled = func() sharding.Coordinator {
			shardC := testscommon.NewMultiShardsCoordinatorMock(2)
			shardC.CurrentShard = core.MetachainShardId
			return shardC
		}
		ccf, _ := consensusComp.NewConsensusComponentsFactory(args)
		require.NotNil(t, ccf)

		cc, err := ccf.Create()
		require.Equal(t, errorsDrt.ErrNilTrieStorageManager, err)
		require.Nil(t, cc)
	})
	t.Run("createConsensusState fails due to nil public key should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		cryptoCompStub, ok := args.CryptoComponents.(*testsMocks.CryptoComponentsStub)
		require.True(t, ok)
		cnt := 0
		cryptoCompStub.PublicKeyCalled = func() crypto.PublicKey {
			cnt++
			if cnt > 1 {
				return nil
			}
			return &cryptoMocks.PublicKeyStub{}
		}
		processCompStub, ok := args.ProcessComponents.(*testsMocks.ProcessComponentsStub)
		require.True(t, ok)
		processCompStub.ShardCoordinatorCalled = func() sharding.Coordinator {
			shardC := testscommon.NewMultiShardsCoordinatorMock(2)
			shardC.CurrentShard = core.MetachainShardId // coverage
			return shardC
		}
		ccf, _ := consensusComp.NewConsensusComponentsFactory(args)
		require.NotNil(t, ccf)

		cc, err := ccf.Create()
		require.Equal(t, errorsDrt.ErrNilPublicKey, err)
		require.Nil(t, cc)
	})
	t.Run("createConsensusState fails due to ToByteArray failure should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		cryptoCompStub, ok := args.CryptoComponents.(*testsMocks.CryptoComponentsStub)
		require.True(t, ok)
		cryptoCompStub.PubKey = &cryptoMocks.PublicKeyStub{
			ToByteArrayStub: func() ([]byte, error) {
				return nil, expectedErr
			},
		}
		ccf, _ := consensusComp.NewConsensusComponentsFactory(args)
		require.NotNil(t, ccf)

		cc, err := ccf.Create()
		require.Equal(t, expectedErr, err)
		require.Nil(t, cc)
	})
	t.Run("createConsensusState fails due to GetConsensusWhitelistedNodes failure should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		processCompStub, ok := args.ProcessComponents.(*testsMocks.ProcessComponentsStub)
		require.True(t, ok)
		processCompStub.NodesCoordinatorCalled = func() nodesCoordinator.NodesCoordinator {
			return &shardingMocks.NodesCoordinatorMock{
				GetConsensusWhitelistedNodesCalled: func(epoch uint32) (map[string]struct{}, error) {
					return nil, expectedErr
				},
			}
		}
		ccf, _ := consensusComp.NewConsensusComponentsFactory(args)
		require.NotNil(t, ccf)

		cc, err := ccf.Create()
		require.Equal(t, expectedErr, err)
		require.Nil(t, cc)
	})
	t.Run("GetConsensusCoreFactory failure should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		args.Config.Consensus.Type = "invalid" // GetConsensusCoreFactory fails
		ccf, _ := consensusComp.NewConsensusComponentsFactory(args)
		require.NotNil(t, ccf)

		cc, err := ccf.Create()
		require.Error(t, err)
		require.Nil(t, cc)
	})
	t.Run("GetBroadcastMessenger failure should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		processCompStub, ok := args.ProcessComponents.(*testsMocks.ProcessComponentsStub)
		require.True(t, ok)
		cnt := 0
		processCompStub.ShardCoordinatorCalled = func() sharding.Coordinator {
			cnt++
			if cnt > 6 {
				return nil // GetBroadcastMessenger fails
			}
			return testscommon.NewMultiShardsCoordinatorMock(2)
		}
		dataCompStub, ok := args.DataComponents.(*testsMocks.DataComponentsStub)
		require.True(t, ok)
		dataCompStub.BlockChain = &testscommon.ChainHandlerStub{
			GetGenesisHeaderHashCalled: func() []byte {
				return []byte("genesis hash")
			},
			GetGenesisHeaderCalled: func() data.HeaderHandler {
				return &testscommon.HeaderHandlerStub{}
			},
			GetCurrentBlockHeaderCalled: func() data.HeaderHandler {
				return &testscommon.HeaderHandlerStub{} // coverage
			},
		}
		ccf, _ := consensusComp.NewConsensusComponentsFactory(args)
		require.NotNil(t, ccf)

		cc, err := ccf.Create()
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), "shard coordinator"))
		require.Nil(t, cc)
	})
	t.Run("NewWorker failure should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		args.Config.Marshalizer.SizeCheckDelta = 1 // coverage
		processCompStub, ok := args.ProcessComponents.(*testsMocks.ProcessComponentsStub)
		require.True(t, ok)
		processCompStub.HeaderIntegrVerif = nil
		ccf, _ := consensusComp.NewConsensusComponentsFactory(args)
		require.NotNil(t, ccf)

		cc, err := ccf.Create()
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), "header integrity verifier"))
		require.Nil(t, cc)
	})
	t.Run("createConsensusTopic fails due nil shard coordinator should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		processCompStub, ok := args.ProcessComponents.(*testsMocks.ProcessComponentsStub)
		require.True(t, ok)
		cnt := 0
		processCompStub.ShardCoordinatorCalled = func() sharding.Coordinator {
			cnt++
			if cnt >= 10 {
				return nil // createConsensusTopic fails
			}
			return testscommon.NewMultiShardsCoordinatorMock(2)
		}
		ccf, _ := consensusComp.NewConsensusComponentsFactory(args)
		require.NotNil(t, ccf)

		cc, err := ccf.Create()
		require.Equal(t, errorsDrt.ErrNilShardCoordinator, err)
		require.Nil(t, cc)
	})
	t.Run("createConsensusTopic fails due nil messenger should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		netwCompStub, ok := args.NetworkComponents.(*testsMocks.NetworkComponentsStub)
		require.True(t, ok)
		cnt := 0
		netwCompStub.MessengerCalled = func() p2p.Messenger {
			cnt++
			if cnt > 4 {
				return nil
			}
			return &p2pmocks.MessengerStub{}
		}
		ccf, _ := consensusComp.NewConsensusComponentsFactory(args)
		require.NotNil(t, ccf)

		cc, err := ccf.Create()
		require.Equal(t, errorsDrt.ErrNilMessenger, err)
		require.Nil(t, cc)
	})
	t.Run("createConsensusTopic fails due CreateTopic failure should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		netwCompStub, ok := args.NetworkComponents.(*testsMocks.NetworkComponentsStub)
		require.True(t, ok)
		netwCompStub.Messenger = &p2pmocks.MessengerStub{
			HasTopicCalled: func(name string) bool {
				return false
			},
			CreateTopicCalled: func(name string, createChannelForTopic bool) error {
				return expectedErr
			},
		}
		ccf, _ := consensusComp.NewConsensusComponentsFactory(args)
		require.NotNil(t, ccf)

		cc, err := ccf.Create()
		require.Equal(t, expectedErr, err)
		require.Nil(t, cc)
	})
	t.Run("createConsensusState fails due to nil KeysHandler should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		cryptoCompStub, ok := args.CryptoComponents.(*testsMocks.CryptoComponentsStub)
		require.True(t, ok)
		cnt := 0
		cryptoCompStub.KeysHandlerCalled = func() consensus.KeysHandler {
			cnt++
			if cnt > 0 {
				return nil
			}
			return &testscommon.KeysHandlerStub{}
		}
		ccf, _ := consensusComp.NewConsensusComponentsFactory(args)
		require.NotNil(t, ccf)

		cc, err := ccf.Create()
		require.Error(t, err)
		require.Contains(t, err.Error(), "keys handler")
		require.Nil(t, cc)
	})
	t.Run("NewConsensusCore failure should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		cryptoCompStub, ok := args.CryptoComponents.(*testsMocks.CryptoComponentsStub)
		require.True(t, ok)
		cryptoCompStub.SigHandler = nil
		ccf, _ := consensusComp.NewConsensusComponentsFactory(args)
		require.NotNil(t, ccf)

		cc, err := ccf.Create()
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), "signing handler"))
		require.Nil(t, cc)
	})
	t.Run("addCloserInstances failure should error", func(t *testing.T) {
		t.Parallel()

		args := createMockConsensusComponentsFactoryArgs()
		processCompStub, ok := args.ProcessComponents.(*testsMocks.ProcessComponentsStub)
		require.True(t, ok)
		processCompStub.HardforkTriggerField = &testscommon.HardforkTriggerStub{
			AddCloserCalled: func(closer update.Closer) error {
				return expectedErr
			},
		}
		ccf, _ := consensusComp.NewConsensusComponentsFactory(args)
		require.NotNil(t, ccf)

		cc, err := ccf.Create()
		require.Equal(t, expectedErr, err)
		require.Nil(t, cc)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		ccf, _ := consensusComp.NewConsensusComponentsFactory(createMockConsensusComponentsFactoryArgs())
		require.NotNil(t, ccf)

		cc, err := ccf.Create()
		require.NoError(t, err)
		require.NotNil(t, cc)

		require.Nil(t, cc.Close())
	})
}
