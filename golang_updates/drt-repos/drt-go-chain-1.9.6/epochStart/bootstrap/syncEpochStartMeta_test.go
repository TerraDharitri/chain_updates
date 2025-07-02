package bootstrap

import (
	"errors"
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/common/graceperiod"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/epochStart"
	"github.com/TerraDharitri/drt-go-chain/epochStart/mock"
	"github.com/TerraDharitri/drt-go-chain/p2p"
	processMock "github.com/TerraDharitri/drt-go-chain/process/mock"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/cryptoMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/testscommon/economicsmocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/hashingMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/p2pmocks"
)

func TestNewEpochStartMetaSyncer_NilsShouldError(t *testing.T) {
	t.Parallel()

	args := getEpochStartSyncerArgs()
	args.CoreComponentsHolder = nil
	ess, err := NewEpochStartMetaSyncer(args)
	assert.True(t, check.IfNil(ess))
	assert.Equal(t, epochStart.ErrNilCoreComponentsHolder, err)

	args = getEpochStartSyncerArgs()
	args.CryptoComponentsHolder = nil
	ess, err = NewEpochStartMetaSyncer(args)
	assert.True(t, check.IfNil(ess))
	assert.Equal(t, epochStart.ErrNilCryptoComponentsHolder, err)

	args = getEpochStartSyncerArgs()
	args.HeaderIntegrityVerifier = nil
	ess, err = NewEpochStartMetaSyncer(args)
	assert.True(t, check.IfNil(ess))
	assert.Equal(t, epochStart.ErrNilHeaderIntegrityVerifier, err)

	args = getEpochStartSyncerArgs()
	args.MetaBlockProcessor = nil
	ess, err = NewEpochStartMetaSyncer(args)
	assert.True(t, check.IfNil(ess))
	assert.Equal(t, epochStart.ErrNilMetablockProcessor, err)

	args = getEpochStartSyncerArgs()
	args.InterceptedDataVerifierFactory = nil
	ess, err = NewEpochStartMetaSyncer(args)
	assert.True(t, check.IfNil(ess))
	assert.Equal(t, epochStart.ErrNilInterceptedDataVerifierFactory, err)
}

func TestNewEpochStartMetaSyncer_ShouldWork(t *testing.T) {
	t.Parallel()

	args := getEpochStartSyncerArgs()
	ess, err := NewEpochStartMetaSyncer(args)
	require.NoError(t, err)
	require.False(t, check.IfNil(ess))
}

func TestEpochStartMetaSyncer_SyncEpochStartMetaRegisterMessengerProcessorFailsShouldErr(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("expected error")

	args := getEpochStartSyncerArgs()
	messenger := &p2pmocks.MessengerStub{
		RegisterMessageProcessorCalled: func(_ string, _ string, _ p2p.MessageProcessor) error {
			return expectedErr
		},
	}
	args.Messenger = messenger
	ess, err := NewEpochStartMetaSyncer(args)
	require.NoError(t, err)

	mb, err := ess.SyncEpochStartMeta(time.Second)
	require.Equal(t, expectedErr, err)
	require.Nil(t, mb)
}

func TestEpochStartMetaSyncer_SyncEpochStartMetaProcessorFailsShouldErr(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("expected error")

	args := getEpochStartSyncerArgs()
	messenger := &p2pmocks.MessengerStub{
		ConnectedPeersCalled: func() []core.PeerID {
			return []core.PeerID{"peer_0", "peer_1", "peer_2", "peer_3", "peer_4", "peer_5"}
		},
	}
	args.Messenger = messenger
	ess, _ := NewEpochStartMetaSyncer(args)

	mbIntercProc := &mock.MetaBlockInterceptorProcessorStub{
		GetEpochStartMetaBlockCalled: func() (data.MetaHeaderHandler, error) {
			return nil, expectedErr
		},
	}
	ess.SetEpochStartMetaBlockInterceptorProcessor(mbIntercProc)

	mb, err := ess.SyncEpochStartMeta(time.Second)
	require.Equal(t, expectedErr, err)
	require.Nil(t, mb)
}

func TestEpochStartMetaSyncer_SyncEpochStartMetaShouldWork(t *testing.T) {
	t.Parallel()

	expectedMb := &block.MetaBlock{Nonce: 37}

	args := getEpochStartSyncerArgs()
	messenger := &p2pmocks.MessengerStub{
		ConnectedPeersCalled: func() []core.PeerID {
			return []core.PeerID{"peer_0", "peer_1", "peer_2", "peer_3", "peer_4", "peer_5"}
		},
	}
	args.Messenger = messenger
	ess, _ := NewEpochStartMetaSyncer(args)

	mbIntercProc := &mock.MetaBlockInterceptorProcessorStub{
		GetEpochStartMetaBlockCalled: func() (data.MetaHeaderHandler, error) {
			return expectedMb, nil
		},
	}
	ess.SetEpochStartMetaBlockInterceptorProcessor(mbIntercProc)

	mb, err := ess.SyncEpochStartMeta(time.Second)
	require.NoError(t, err)
	require.Equal(t, expectedMb, mb)
}

func getEpochStartSyncerArgs() ArgsNewEpochStartMetaSyncer {
	gracePeriod, _ := graceperiod.NewEpochChangeGracePeriod([]config.EpochChangeGracePeriodByEpoch{{EnableEpoch: 0, GracePeriodInRounds: 1}})
	return ArgsNewEpochStartMetaSyncer{
		CoreComponentsHolder: &mock.CoreComponentsMock{
			IntMarsh:            &mock.MarshalizerMock{},
			Marsh:               &mock.MarshalizerMock{},
			Hash:                &hashingMocks.HasherMock{},
			UInt64ByteSliceConv: &mock.Uint64ByteSliceConverterMock{},
			AddrPubKeyConv:      testscommon.NewPubkeyConverterMock(32),
			PathHdl:             &testscommon.PathManagerStub{},
			ChainIdCalled: func() string {
				return "chain-ID"
			},
			EpochChangeGracePeriodHandlerField: gracePeriod,
		},
		CryptoComponentsHolder: &mock.CryptoComponentsMock{
			PubKey:   &cryptoMocks.PublicKeyStub{},
			BlockSig: &cryptoMocks.SignerStub{},
			TxSig:    &cryptoMocks.SignerStub{},
			BlKeyGen: &cryptoMocks.KeyGenStub{},
			TxKeyGen: &cryptoMocks.KeyGenStub{},
		},
		RequestHandler:   &testscommon.RequestHandlerStub{},
		Messenger:        &p2pmocks.MessengerStub{},
		ShardCoordinator: mock.NewMultiShardsCoordinatorMock(2),
		EconomicsData:    &economicsmocks.EconomicsHandlerMock{},
		WhitelistHandler: &testscommon.WhiteListHandlerStub{},
		StartInEpochConfig: config.EpochStartConfig{
			MinNumConnectedPeersToStart:       2,
			MinNumOfPeersToConsiderBlockValid: 2,
		},
		HeaderIntegrityVerifier:        &mock.HeaderIntegrityVerifierStub{},
		MetaBlockProcessor:             &mock.EpochStartMetaBlockProcessorStub{},
		InterceptedDataVerifierFactory: &processMock.InterceptedDataVerifierFactoryMock{},
		ProofsPool:                     &dataRetriever.ProofsPoolMock{},
		ProofsInterceptorProcessor:     &processMock.InterceptorProcessorStub{},
	}
}
