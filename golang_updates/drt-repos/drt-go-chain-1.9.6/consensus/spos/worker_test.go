package spos_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	crypto "github.com/TerraDharitri/drt-go-chain-crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/consensus"
	"github.com/TerraDharitri/drt-go-chain/consensus/mock"
	"github.com/TerraDharitri/drt-go-chain/consensus/spos"
	"github.com/TerraDharitri/drt-go-chain/consensus/spos/bls"
	"github.com/TerraDharitri/drt-go-chain/p2p"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/bootstrapperStubs"
	"github.com/TerraDharitri/drt-go-chain/testscommon/cache"
	consensusMocks "github.com/TerraDharitri/drt-go-chain/testscommon/consensus"
	"github.com/TerraDharitri/drt-go-chain/testscommon/enableEpochsHandlerMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/hashingMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/p2pmocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/processMocks"
	statusHandlerMock "github.com/TerraDharitri/drt-go-chain/testscommon/statusHandler"
)

const roundTimeDuration = 100 * time.Millisecond

var fromConnectedPeerId = core.PeerID("connected peer id")

const HashSize = 32
const SignatureSize = 48
const PublicKeySize = 96

var blockHeaderHash = make([]byte, HashSize)
var invalidBlockHeaderHash = make([]byte, HashSize+1)
var signature = make([]byte, SignatureSize)
var invalidSignature = make([]byte, SignatureSize+1)
var publicKey = make([]byte, PublicKeySize)

func createDefaultWorkerArgs(appStatusHandler core.AppStatusHandler) *spos.WorkerArgs {
	blockchainMock := &testscommon.ChainHandlerStub{}
	blockProcessor := &testscommon.BlockProcessorStub{
		DecodeBlockHeaderCalled: func(dta []byte) data.HeaderHandler {
			return nil
		},
		RevertCurrentBlockCalled: func() {
		},
		DecodeBlockBodyCalled: func(dta []byte) data.BodyHandler {
			return nil
		},
	}
	bootstrapperMock := &bootstrapperStubs.BootstrapperStub{}
	broadcastMessengerMock := &consensusMocks.BroadcastMessengerMock{}
	consensusState := initConsensusState()
	forkDetectorMock := &processMocks.ForkDetectorStub{}
	forkDetectorMock.AddHeaderCalled = func(header data.HeaderHandler, hash []byte, state process.BlockHeaderState, selfNotarizedHeaders []data.HeaderHandler, selfNotarizedHeadersHashes [][]byte) error {
		return nil
	}
	keyGeneratorMock, _, _ := consensusMocks.InitKeys()
	marshalizerMock := mock.MarshalizerMock{}
	roundHandlerMock := initRoundHandlerMock()
	shardCoordinatorMock := mock.ShardCoordinatorMock{}
	singleSignerMock := &mock.SingleSignerMock{
		SignStub: func(private crypto.PrivateKey, msg []byte) ([]byte, error) {
			return []byte("signed"), nil
		},
		VerifyStub: func(public crypto.PublicKey, msg []byte, sig []byte) error {
			return nil
		},
	}
	syncTimerMock := &consensusMocks.SyncTimerMock{}
	hasher := &hashingMocks.HasherMock{}
	blsService, _ := bls.NewConsensusService()
	poolAdder := cache.NewCacherMock()

	scheduledProcessorArgs := spos.ScheduledProcessorWrapperArgs{
		SyncTimer:                syncTimerMock,
		Processor:                blockProcessor,
		RoundTimeDurationHandler: roundHandlerMock,
	}
	scheduledProcessor, _ := spos.NewScheduledProcessorWrapper(scheduledProcessorArgs)

	peerSigHandler := &mock.PeerSignatureHandler{Signer: singleSignerMock, KeyGen: keyGeneratorMock}

	workerArgs := &spos.WorkerArgs{
		ConsensusService:         blsService,
		BlockChain:               blockchainMock,
		BlockProcessor:           blockProcessor,
		ScheduledProcessor:       scheduledProcessor,
		Bootstrapper:             bootstrapperMock,
		BroadcastMessenger:       broadcastMessengerMock,
		ConsensusState:           consensusState,
		ForkDetector:             forkDetectorMock,
		Marshalizer:              marshalizerMock,
		Hasher:                   hasher,
		RoundHandler:             roundHandlerMock,
		ShardCoordinator:         shardCoordinatorMock,
		PeerSignatureHandler:     peerSigHandler,
		SyncTimer:                syncTimerMock,
		HeaderSigVerifier:        &consensusMocks.HeaderSigVerifierMock{},
		HeaderIntegrityVerifier:  &testscommon.HeaderVersionHandlerStub{},
		ChainID:                  chainID,
		NetworkShardingCollector: &p2pmocks.NetworkShardingCollectorStub{},
		AntifloodHandler:         createMockP2PAntifloodHandler(),
		PoolAdder:                poolAdder,
		SignatureSize:            SignatureSize,
		PublicKeySize:            PublicKeySize,
		AppStatusHandler:         appStatusHandler,
		NodeRedundancyHandler:    &mock.NodeRedundancyHandlerStub{},
		PeerBlacklistHandler:     &mock.PeerBlacklistHandlerStub{},
		EnableEpochsHandler:      &enableEpochsHandlerMock.EnableEpochsHandlerStub{},
		InvalidSignersCache:      &consensusMocks.InvalidSignersCacheMock{},
	}

	return workerArgs
}

func createMockP2PAntifloodHandler() *mock.P2PAntifloodHandlerStub {
	return &mock.P2PAntifloodHandlerStub{
		CanProcessMessageCalled: func(message p2p.MessageP2P, fromConnectedPeer core.PeerID) error {
			return nil
		},
		CanProcessMessagesOnTopicCalled: func(peer core.PeerID, topic string, numMessages uint32, totalSize uint64, sequence []byte) error {
			return nil
		},
	}
}

func initWorker(appStatusHandler core.AppStatusHandler) *spos.Worker {
	workerArgs := createDefaultWorkerArgs(appStatusHandler)
	sposWorker, _ := spos.NewWorker(workerArgs)

	sposWorker.ConsensusState().SetHeader(&block.HeaderV2{})

	return sposWorker
}

func initRoundHandlerMock() *consensusMocks.RoundHandlerMock {
	return &consensusMocks.RoundHandlerMock{
		RoundIndex: 0,
		TimeStampCalled: func() time.Time {
			return time.Unix(0, 0)
		},
		TimeDurationCalled: func() time.Duration {
			return roundTimeDuration
		},
	}
}

func TestWorker_NewWorkerConsensusServiceNilShouldFail(t *testing.T) {
	t.Parallel()

	workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
	workerArgs.ConsensusService = nil
	wrk, err := spos.NewWorker(workerArgs)

	assert.Nil(t, wrk)
	assert.Equal(t, spos.ErrNilConsensusService, err)
}

func TestWorker_NewWorkerBlockChainNilShouldFail(t *testing.T) {
	t.Parallel()
	workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
	workerArgs.BlockChain = nil
	wrk, err := spos.NewWorker(workerArgs)

	assert.Nil(t, wrk)
	assert.Equal(t, spos.ErrNilBlockChain, err)
}

func TestWorker_NewWorkerBlockProcessorNilShouldFail(t *testing.T) {
	t.Parallel()

	workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
	workerArgs.BlockProcessor = nil
	wrk, err := spos.NewWorker(workerArgs)

	assert.Nil(t, wrk)
	assert.Equal(t, spos.ErrNilBlockProcessor, err)
}

func TestWorker_NewWorkerBootstrapperNilShouldFail(t *testing.T) {
	t.Parallel()

	workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
	workerArgs.Bootstrapper = nil
	wrk, err := spos.NewWorker(workerArgs)

	assert.Nil(t, wrk)
	assert.Equal(t, spos.ErrNilBootstrapper, err)
}

func TestWorker_NewWorkerBroadcastMessengerNilShouldFail(t *testing.T) {
	t.Parallel()

	workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
	workerArgs.BroadcastMessenger = nil
	wrk, err := spos.NewWorker(workerArgs)

	assert.Nil(t, wrk)
	assert.Equal(t, spos.ErrNilBroadcastMessenger, err)
}

func TestWorker_NewWorkerConsensusStateNilShouldFail(t *testing.T) {
	t.Parallel()

	workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
	workerArgs.ConsensusState = nil
	wrk, err := spos.NewWorker(workerArgs)

	assert.Nil(t, wrk)
	assert.Equal(t, spos.ErrNilConsensusState, err)
}

func TestWorker_NewWorkerForkDetectorNilShouldFail(t *testing.T) {
	t.Parallel()

	workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
	workerArgs.ForkDetector = nil
	wrk, err := spos.NewWorker(workerArgs)

	assert.Nil(t, wrk)
	assert.Equal(t, spos.ErrNilForkDetector, err)
}

func TestWorker_NewWorkerMarshalizerNilShouldFail(t *testing.T) {
	t.Parallel()

	workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
	workerArgs.Marshalizer = nil
	wrk, err := spos.NewWorker(workerArgs)

	assert.Nil(t, wrk)
	assert.Equal(t, spos.ErrNilMarshalizer, err)
}

func TestWorker_NewWorkerHasherNilShouldFail(t *testing.T) {
	t.Parallel()

	workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
	workerArgs.Hasher = nil
	wrk, err := spos.NewWorker(workerArgs)

	assert.Nil(t, wrk)
	assert.Equal(t, spos.ErrNilHasher, err)
}

func TestWorker_NewWorkerRoundHandlerNilShouldFail(t *testing.T) {
	t.Parallel()

	workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
	workerArgs.RoundHandler = nil
	wrk, err := spos.NewWorker(workerArgs)

	assert.Nil(t, wrk)
	assert.Equal(t, spos.ErrNilRoundHandler, err)
}

func TestWorker_NewWorkerShardCoordinatorNilShouldFail(t *testing.T) {
	t.Parallel()

	workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
	workerArgs.ShardCoordinator = nil
	wrk, err := spos.NewWorker(workerArgs)

	assert.Nil(t, wrk)
	assert.Equal(t, spos.ErrNilShardCoordinator, err)
}

func TestWorker_NewWorkerPeerSignatureHandlerNilShouldFail(t *testing.T) {
	t.Parallel()

	workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
	workerArgs.PeerSignatureHandler = nil
	wrk, err := spos.NewWorker(workerArgs)

	assert.Nil(t, wrk)
	assert.Equal(t, spos.ErrNilPeerSignatureHandler, err)
}

func TestWorker_NewWorkerSyncTimerNilShouldFail(t *testing.T) {
	t.Parallel()

	workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
	workerArgs.SyncTimer = nil
	wrk, err := spos.NewWorker(workerArgs)

	assert.Nil(t, wrk)
	assert.Equal(t, spos.ErrNilSyncTimer, err)
}

func TestWorker_NewWorkerHeaderSigVerifierNilShouldFail(t *testing.T) {
	t.Parallel()

	workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
	workerArgs.HeaderSigVerifier = nil
	wrk, err := spos.NewWorker(workerArgs)

	assert.Nil(t, wrk)
	assert.Equal(t, spos.ErrNilHeaderSigVerifier, err)
}

func TestWorker_NewWorkerHeaderIntegrityVerifierShouldFail(t *testing.T) {
	t.Parallel()

	workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
	workerArgs.HeaderIntegrityVerifier = nil
	wrk, err := spos.NewWorker(workerArgs)

	assert.Nil(t, wrk)
	assert.Equal(t, spos.ErrNilHeaderIntegrityVerifier, err)
}

func TestWorker_NewWorkerEmptyChainIDShouldFail(t *testing.T) {
	t.Parallel()

	workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
	workerArgs.ChainID = nil
	wrk, err := spos.NewWorker(workerArgs)

	assert.Nil(t, wrk)
	assert.Equal(t, spos.ErrInvalidChainID, err)
}

func TestWorker_NewWorkerNilNetworkShardingCollectorShouldFail(t *testing.T) {
	t.Parallel()

	workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
	workerArgs.NetworkShardingCollector = nil
	wrk, err := spos.NewWorker(workerArgs)

	assert.Nil(t, wrk)
	assert.Equal(t, spos.ErrNilNetworkShardingCollector, err)
}

func TestWorker_NewWorkerNilAntifloodHandlerShouldFail(t *testing.T) {
	t.Parallel()

	workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
	workerArgs.AntifloodHandler = nil
	wrk, err := spos.NewWorker(workerArgs)

	assert.Nil(t, wrk)
	assert.Equal(t, spos.ErrNilAntifloodHandler, err)
}

func TestWorker_NewWorkerPoolAdderNilShouldFail(t *testing.T) {
	t.Parallel()

	workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
	workerArgs.PoolAdder = nil
	wrk, err := spos.NewWorker(workerArgs)

	assert.Nil(t, wrk)
	assert.Equal(t, spos.ErrNilPoolAdder, err)
}

func TestWorker_NewWorkerNodeRedundancyHandlerShouldFail(t *testing.T) {
	t.Parallel()

	workerArgs := createDefaultWorkerArgs(statusHandlerMock.NewAppStatusHandlerMock())
	workerArgs.NodeRedundancyHandler = nil
	wrk, err := spos.NewWorker(workerArgs)

	assert.Nil(t, wrk)
	assert.Equal(t, spos.ErrNilNodeRedundancyHandler, err)
}

func TestWorker_NewWorkerPoolEnableEpochsHandlerNilShouldFail(t *testing.T) {
	t.Parallel()

	workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
	workerArgs.EnableEpochsHandler = nil
	wrk, err := spos.NewWorker(workerArgs)

	assert.Nil(t, wrk)
	assert.Equal(t, spos.ErrNilEnableEpochsHandler, err)
}

func TestWorker_NewWorkerPoolInvalidSignersCacheNilShouldFail(t *testing.T) {
	t.Parallel()

	workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
	workerArgs.InvalidSignersCache = nil
	wrk, err := spos.NewWorker(workerArgs)

	assert.Nil(t, wrk)
	assert.Equal(t, spos.ErrNilInvalidSignersCache, err)
}

func TestWorker_NewWorkerShouldWork(t *testing.T) {
	t.Parallel()

	workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
	wrk, err := spos.NewWorker(workerArgs)

	assert.Nil(t, err)
	assert.False(t, check.IfNil(wrk))
}

func TestWorker_ProcessReceivedMessageShouldErrIfFloodIsDetectedOnTopic(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("flood detected")
	workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
	antifloodHandler := &mock.P2PAntifloodHandlerStub{
		CanProcessMessageCalled: func(message p2p.MessageP2P, fromConnectedPeer core.PeerID) error {
			return nil
		},
		CanProcessMessagesOnTopicCalled: func(peer core.PeerID, topic string, numMessages uint32, totalSize uint64, sequence []byte) error {
			return expectedErr
		},
	}

	workerArgs.AntifloodHandler = antifloodHandler
	wrk, _ := spos.NewWorker(workerArgs)

	msg := &p2pmocks.P2PMessageMock{
		DataField:      []byte("aaa"),
		TopicField:     "topic1",
		SignatureField: []byte("signature"),
	}
	msgID, err := wrk.ProcessReceivedMessage(msg, "peer", &p2pmocks.MessengerStub{})
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, msgID)
}

func TestWorker_ReceivedSyncStateShouldNotSendOnChannelWhenInputIsFalse(t *testing.T) {
	t.Parallel()
	wrk := initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	wrk.ReceivedSyncState(false)
	rcv := false
	select {
	case rcv = <-wrk.ConsensusStateChangedChannel():
	case <-time.After(100 * time.Millisecond):
	}

	assert.False(t, rcv)
}

func TestWorker_ReceivedSyncStateShouldNotSendOnChannelWhenChannelIsBusy(t *testing.T) {
	t.Parallel()
	wrk := initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	wrk.ConsensusStateChangedChannel() <- false
	wrk.ReceivedSyncState(true)
	rcv := false
	select {
	case rcv = <-wrk.ConsensusStateChangedChannel():
	case <-time.After(100 * time.Millisecond):
	}

	assert.False(t, rcv)
}

func TestWorker_ReceivedSyncStateShouldSendOnChannel(t *testing.T) {
	t.Parallel()
	wrk := initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	wrk.ReceivedSyncState(true)
	rcv := false
	select {
	case rcv = <-wrk.ConsensusStateChangedChannel():
	case <-time.After(100 * time.Millisecond):
	}

	assert.True(t, rcv)
}

func TestWorker_InitReceivedMessagesShouldInitMap(t *testing.T) {
	t.Parallel()
	wrk := initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	wrk.NilReceivedMessages()
	wrk.InitReceivedMessages()

	assert.NotNil(t, wrk.ReceivedMessages()[bls.MtBlockBody])
}

func TestWorker_AddReceivedMessageCallShouldWork(t *testing.T) {
	t.Parallel()
	wrk := initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	receivedMessageCall := func(context.Context, *consensus.Message) bool {
		return true
	}
	wrk.AddReceivedMessageCall(bls.MtBlockBody, receivedMessageCall)
	receivedMessageCalls := wrk.ReceivedMessagesCalls()

	assert.Equal(t, 1, len(receivedMessageCalls))
	assert.NotNil(t, receivedMessageCalls[bls.MtBlockBody])
	assert.True(t, receivedMessageCalls[bls.MtBlockBody][0](context.Background(), nil))
}

func TestWorker_RemoveAllReceivedMessageCallsShouldWork(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	receivedMessageCall := func(context.Context, *consensus.Message) bool {
		return true
	}
	wrk.AddReceivedMessageCall(bls.MtBlockBody, receivedMessageCall)
	receivedMessageCalls := wrk.ReceivedMessagesCalls()

	assert.Equal(t, 1, len(receivedMessageCalls))
	assert.NotNil(t, receivedMessageCalls[bls.MtBlockBody])
	assert.True(t, receivedMessageCalls[bls.MtBlockBody][0](context.Background(), nil))

	wrk.RemoveAllReceivedMessagesCalls()
	receivedMessageCalls = wrk.ReceivedMessagesCalls()

	assert.Equal(t, 0, len(receivedMessageCalls))
	assert.Nil(t, receivedMessageCalls[bls.MtBlockBody])
}

func TestWorker_ProcessReceivedMessageTxBlockBodyShouldRetNil(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	blk := &block.Body{}
	blkStr, _ := mock.MarshalizerMock{}.Marshal(blk)
	cnsMsg := consensus.NewConsensusMessage(
		nil,
		nil,
		blkStr,
		nil,
		[]byte(wrk.ConsensusState().ConsensusGroup()[0]),
		signature,
		int(bls.MtBlockBody),
		0,
		chainID,
		nil,
		nil,
		nil,
		currentPid,
		nil,
	)
	buff, _ := wrk.Marshalizer().Marshal(cnsMsg)
	time.Sleep(time.Second)
	msg := &p2pmocks.P2PMessageMock{
		DataField:      buff,
		PeerField:      currentPid,
		SignatureField: []byte("signature"),
	}
	msgID, err := wrk.ProcessReceivedMessage(msg, fromConnectedPeerId, &p2pmocks.MessengerStub{})
	assert.Nil(t, err)
	assert.Len(t, msgID, 0)
}

func TestWorker_ProcessReceivedMessageNilMessageShouldErr(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	msgID, err := wrk.ProcessReceivedMessage(nil, fromConnectedPeerId, &p2pmocks.MessengerStub{})
	time.Sleep(time.Second)

	assert.Equal(t, 0, len(wrk.ReceivedMessages()[bls.MtBlockBody]))
	assert.Equal(t, spos.ErrNilMessage, err)
	assert.Nil(t, msgID)
}

func TestWorker_ProcessReceivedMessageNilMessageDataFieldShouldErr(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	msgID, err := wrk.ProcessReceivedMessage(&p2pmocks.P2PMessageMock{}, fromConnectedPeerId, &p2pmocks.MessengerStub{})
	time.Sleep(time.Second)

	assert.Equal(t, 0, len(wrk.ReceivedMessages()[bls.MtBlockBody]))
	assert.Equal(t, spos.ErrNilDataToProcess, err)
	assert.Nil(t, msgID)
}

func TestWorker_ProcessReceivedMessageEmptySignatureFieldShouldErr(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	msgID, err := wrk.ProcessReceivedMessage(
		&p2pmocks.P2PMessageMock{
			DataField: []byte("data field"),
		},
		fromConnectedPeerId,
		&p2pmocks.MessengerStub{},
	)
	time.Sleep(time.Second)

	assert.Equal(t, 0, len(wrk.ReceivedMessages()[bls.MtBlockBody]))
	assert.Equal(t, spos.ErrNilSignatureOnP2PMessage, err)
	assert.Nil(t, msgID)
}

func TestWorker_ProcessReceivedMessageRedundancyNodeShouldResetInactivityIfNeeded(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	var wasCalled bool
	nodeRedundancyMock := &mock.NodeRedundancyHandlerStub{
		IsRedundancyNodeCalled: func() bool {
			return true
		},
		ResetInactivityIfNeededCalled: func(selfPubKey string, consensusMsgPubKey string, consensusMsgPeerID core.PeerID) {
			wasCalled = true
		},
	}
	wrk.SetNodeRedundancyHandler(nodeRedundancyMock)
	buff, _ := wrk.Marshalizer().Marshal(&consensus.Message{})
	_, _ = wrk.ProcessReceivedMessage(
		&p2pmocks.P2PMessageMock{
			DataField:      buff,
			SignatureField: []byte("signature"),
		},
		fromConnectedPeerId,
		&p2pmocks.MessengerStub{},
	)

	assert.True(t, wasCalled)
}

func TestWorker_ProcessReceivedMessageNodeNotInEligibleListShouldErr(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	blk := &block.Body{}
	blkStr, _ := mock.MarshalizerMock{}.Marshal(blk)
	cnsMsg := consensus.NewConsensusMessage(
		nil,
		nil,
		blkStr,
		nil,
		publicKey,
		signature,
		int(bls.MtBlockBody),
		0,
		chainID,
		nil,
		nil,
		nil,
		currentPid,
		nil,
	)
	buff, _ := wrk.Marshalizer().Marshal(cnsMsg)
	msgID, err := wrk.ProcessReceivedMessage(
		&p2pmocks.P2PMessageMock{
			DataField:      buff,
			SignatureField: []byte("signature"),
		},
		fromConnectedPeerId,
		&p2pmocks.MessengerStub{},
	)
	time.Sleep(time.Second)

	assert.Equal(t, 0, len(wrk.ReceivedMessages()[bls.MtBlockBody]))
	assert.True(t, errors.Is(err, spos.ErrNodeIsNotInEligibleList))
	assert.Nil(t, msgID)
}

func TestWorker_ProcessReceivedMessageComputeReceivedProposedBlockMetric(t *testing.T) {
	t.Parallel()

	t.Run("normal operation", func(t *testing.T) {
		t.Parallel()

		roundDuration := time.Millisecond * 1000
		delay := time.Millisecond * 430
		roundStartTimeStamp := time.Now()

		receivedValue, redundancyReason, redundancyStatus := testWorkerProcessReceivedMessageComputeReceivedProposedBlockMetric(
			t,
			roundStartTimeStamp,
			delay,
			roundDuration,
			&mock.NodeRedundancyHandlerStub{},
			&testscommon.KeysHandlerStub{})

		minimumExpectedValue := uint64(delay * 100 / roundDuration)
		assert.True(t,
			receivedValue >= minimumExpectedValue,
			fmt.Sprintf("minimum expected was %d, got %d", minimumExpectedValue, receivedValue),
		)
		assert.Empty(t, redundancyReason)
		assert.True(t, redundancyStatus)
	})
	t.Run("time.Since returns negative value", func(t *testing.T) {
		// test the edgecase when the returned NTP time stored in the round handler is
		// slightly advanced when comparing with time.Now.
		t.Parallel()

		roundDuration := time.Millisecond * 1000
		delay := time.Millisecond * 430
		roundStartTimeStamp := time.Now().Add(time.Minute)

		receivedValue, redundancyReason, redundancyStatus := testWorkerProcessReceivedMessageComputeReceivedProposedBlockMetric(
			t,
			roundStartTimeStamp,
			delay,
			roundDuration,
			&mock.NodeRedundancyHandlerStub{},
			&testscommon.KeysHandlerStub{})

		assert.Zero(t, receivedValue)
		assert.Empty(t, redundancyReason)
		assert.True(t, redundancyStatus)
	})
	t.Run("normal operation as a single-key redundancy node", func(t *testing.T) {
		t.Parallel()

		roundDuration := time.Millisecond * 1000
		delay := time.Millisecond * 430
		roundStartTimeStamp := time.Now()

		receivedValue, redundancyReason, redundancyStatus := testWorkerProcessReceivedMessageComputeReceivedProposedBlockMetric(
			t,
			roundStartTimeStamp,
			delay,
			roundDuration,
			&mock.NodeRedundancyHandlerStub{
				IsMainMachineActiveCalled: func() bool {
					return false
				},
			},
			&testscommon.KeysHandlerStub{})

		minimumExpectedValue := uint64(delay * 100 / roundDuration)
		assert.True(t,
			receivedValue >= minimumExpectedValue,
			fmt.Sprintf("minimum expected was %d, got %d", minimumExpectedValue, receivedValue),
		)
		assert.Equal(t, spos.RedundancySingleKeySteppedIn, redundancyReason)
		assert.False(t, redundancyStatus)
	})
	t.Run("normal operation as a multikey-key redundancy node", func(t *testing.T) {
		t.Parallel()

		roundDuration := time.Millisecond * 1000
		delay := time.Millisecond * 430
		roundStartTimeStamp := time.Now()

		multikeyReason := "multikey step in reason"
		receivedValue, redundancyReason, redundancyStatus := testWorkerProcessReceivedMessageComputeReceivedProposedBlockMetric(
			t,
			roundStartTimeStamp,
			delay,
			roundDuration,
			&mock.NodeRedundancyHandlerStub{},
			&testscommon.KeysHandlerStub{
				GetRedundancyStepInReasonCalled: func() string {
					return multikeyReason
				},
			})

		minimumExpectedValue := uint64(delay * 100 / roundDuration)
		assert.True(t,
			receivedValue >= minimumExpectedValue,
			fmt.Sprintf("minimum expected was %d, got %d", minimumExpectedValue, receivedValue),
		)
		assert.Equal(t, multikeyReason, redundancyReason)
		assert.False(t, redundancyStatus)
	})
}

func testWorkerProcessReceivedMessageComputeReceivedProposedBlockMetric(
	t *testing.T,
	roundStartTimeStamp time.Time,
	delay time.Duration,
	roundDuration time.Duration,
	redundancyHandler consensus.NodeRedundancyHandler,
	keysHandler consensus.KeysHandler,
) (uint64, string, bool) {
	marshaller := mock.MarshalizerMock{}
	receivedValue := uint64(0)
	redundancyReason := ""
	redundancyStatus := false
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{
		SetUInt64ValueHandler: func(key string, value uint64) {
			receivedValue = value
		},
		SetStringValueHandler: func(key string, value string) {
			if key == common.MetricRedundancyIsMainActive {
				var err error
				redundancyStatus, err = strconv.ParseBool(value)
				assert.Nil(t, err)
			}
			if key == common.MetricRedundancyStepInReason {
				redundancyReason = value
			}
		},
	})
	wrk.SetBlockProcessor(&testscommon.BlockProcessorStub{
		DecodeBlockHeaderCalled: func(dta []byte) data.HeaderHandler {
			header := &block.Header{}
			_ = marshaller.Unmarshal(header, dta)

			return header
		},
		RevertCurrentBlockCalled: func() {
		},
		DecodeBlockBodyCalled: func(dta []byte) data.BodyHandler {
			return nil
		},
	})

	wrk.SetRoundHandler(&consensusMocks.RoundHandlerMock{
		RoundIndex: 0,
		TimeDurationCalled: func() time.Duration {
			return roundDuration
		},
		TimeStampCalled: func() time.Time {
			return roundStartTimeStamp
		},
	})
	wrk.SetRedundancyHandler(redundancyHandler)
	wrk.SetKeysHandler(keysHandler)
	hdr := &block.Header{
		ChainID:         chainID,
		PrevHash:        []byte("prev hash"),
		PrevRandSeed:    []byte("prev rand seed"),
		RandSeed:        []byte("rand seed"),
		RootHash:        []byte("roothash"),
		SoftwareVersion: []byte("software version"),
		AccumulatedFees: big.NewInt(0),
		DeveloperFees:   big.NewInt(0),
	}
	hdrHash, _ := core.CalculateHash(mock.MarshalizerMock{}, &hashingMocks.HasherMock{}, hdr)
	hdrStr, _ := marshaller.Marshal(hdr)
	cnsMsg := consensus.NewConsensusMessage(
		hdrHash,
		nil,
		nil,
		hdrStr,
		[]byte(wrk.ConsensusState().Leader()),
		signature,
		int(bls.MtBlockHeader),
		0,
		chainID,
		nil,
		nil,
		nil,
		currentPid,
		nil,
	)

	time.Sleep(delay)

	buff, _ := wrk.Marshalizer().Marshal(cnsMsg)
	msg := &p2pmocks.P2PMessageMock{
		DataField:      buff,
		PeerField:      currentPid,
		SignatureField: []byte("signature"),
	}
	_, _ = wrk.ProcessReceivedMessage(msg, "", &p2pmocks.MessengerStub{})

	return receivedValue, redundancyReason, redundancyStatus
}

func TestWorker_ProcessReceivedMessageInconsistentChainIDInConsensusMessageShouldErr(t *testing.T) {
	t.Parallel()

	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	blk := &block.Body{}
	blkStr, _ := mock.MarshalizerMock{}.Marshal(blk)
	cnsMsg := consensus.NewConsensusMessage(
		blockHeaderHash,
		nil,
		blkStr,
		nil,
		[]byte(wrk.ConsensusState().ConsensusGroup()[0]),
		[]byte("sig"),
		int(bls.MtBlockBody),
		1,
		[]byte("inconsistent chain ID"),
		nil,
		nil,
		nil,
		currentPid,
		nil,
	)
	buff, _ := wrk.Marshalizer().Marshal(cnsMsg)
	msgID, err := wrk.ProcessReceivedMessage(
		&p2pmocks.P2PMessageMock{
			DataField:      buff,
			SignatureField: []byte("signature"),
		},
		fromConnectedPeerId,
		&p2pmocks.MessengerStub{},
	)

	assert.True(t, errors.Is(err, spos.ErrInvalidChainID))
	assert.Nil(t, msgID)
}

func TestWorker_ProcessReceivedMessageTypeInvalidShouldErr(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	blk := &block.Body{}
	blkStr, _ := mock.MarshalizerMock{}.Marshal(blk)
	cnsMsg := consensus.NewConsensusMessage(
		blockHeaderHash,
		nil,
		blkStr,
		nil,
		[]byte(wrk.ConsensusState().ConsensusGroup()[0]),
		[]byte("sig"),
		666,
		0,
		chainID,
		nil,
		nil,
		nil,
		currentPid,
		nil,
	)
	buff, _ := wrk.Marshalizer().Marshal(cnsMsg)
	msgID, err := wrk.ProcessReceivedMessage(
		&p2pmocks.P2PMessageMock{
			DataField:      buff,
			SignatureField: []byte("signature"),
		},
		fromConnectedPeerId,
		&p2pmocks.MessengerStub{},
	)
	time.Sleep(time.Second)

	assert.Equal(t, 0, len(wrk.ReceivedMessages()[666]))
	assert.True(t, errors.Is(err, spos.ErrInvalidMessageType), err)
	assert.Nil(t, msgID)
}

func TestWorker_ProcessReceivedHeaderHashSizeInvalidShouldErr(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	blk := &block.Body{}
	blkStr, _ := mock.MarshalizerMock{}.Marshal(blk)
	cnsMsg := consensus.NewConsensusMessage(
		invalidBlockHeaderHash,
		nil,
		blkStr,
		nil,
		[]byte(wrk.ConsensusState().ConsensusGroup()[0]),
		[]byte("sig"),
		int(bls.MtBlockBody),
		0,
		chainID,
		nil,
		nil,
		nil,
		currentPid,
		nil,
	)
	buff, _ := wrk.Marshalizer().Marshal(cnsMsg)
	msgID, err := wrk.ProcessReceivedMessage(
		&p2pmocks.P2PMessageMock{
			DataField:      buff,
			SignatureField: []byte("signature"),
		},
		fromConnectedPeerId,
		&p2pmocks.MessengerStub{},
	)
	time.Sleep(time.Second)

	assert.Equal(t, 0, len(wrk.ReceivedMessages()[bls.MtBlockBody]))
	assert.True(t, errors.Is(err, spos.ErrInvalidHeaderHashSize), err)
	assert.Nil(t, msgID)
}

func TestWorker_ProcessReceivedMessageForFutureRoundShouldErr(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	blk := &block.Body{}
	blkStr, _ := mock.MarshalizerMock{}.Marshal(blk)
	cnsMsg := consensus.NewConsensusMessage(
		nil,
		nil,
		blkStr,
		nil,
		[]byte(wrk.ConsensusState().ConsensusGroup()[0]),
		signature,
		int(bls.MtBlockBody),
		2,
		chainID,
		nil,
		nil,
		nil,
		currentPid,
		nil,
	)
	buff, _ := wrk.Marshalizer().Marshal(cnsMsg)
	msgID, err := wrk.ProcessReceivedMessage(
		&p2pmocks.P2PMessageMock{
			DataField:      buff,
			SignatureField: []byte("signature"),
		},
		fromConnectedPeerId,
		&p2pmocks.MessengerStub{},
	)
	time.Sleep(time.Second)

	assert.Equal(t, 0, len(wrk.ReceivedMessages()[bls.MtBlockBody]))
	assert.True(t, errors.Is(err, spos.ErrMessageForFutureRound))
	assert.Nil(t, msgID)
}

func TestWorker_ProcessReceivedMessageForPastRoundShouldErr(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	blk := &block.Body{}
	blkStr, _ := mock.MarshalizerMock{}.Marshal(blk)
	cnsMsg := consensus.NewConsensusMessage(
		nil,
		nil,
		blkStr,
		nil,
		[]byte(wrk.ConsensusState().ConsensusGroup()[0]),
		signature,
		int(bls.MtBlockBody),
		-1,
		chainID,
		nil,
		nil,
		nil,
		currentPid,
		nil,
	)
	buff, _ := wrk.Marshalizer().Marshal(cnsMsg)
	msgID, err := wrk.ProcessReceivedMessage(
		&p2pmocks.P2PMessageMock{
			DataField:      buff,
			SignatureField: []byte("signature"),
		},
		fromConnectedPeerId,
		&p2pmocks.MessengerStub{},
	)
	time.Sleep(time.Second)

	assert.Equal(t, 0, len(wrk.ReceivedMessages()[bls.MtBlockBody]))
	assert.True(t, errors.Is(err, spos.ErrMessageForPastRound))
	assert.Nil(t, msgID)
}

func TestWorker_ProcessReceivedMessageTypeLimitReachedShouldErr(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	blk := &block.Body{}
	blkStr, _ := mock.MarshalizerMock{}.Marshal(blk)
	cnsMsg := consensus.NewConsensusMessage(
		nil,
		nil,
		blkStr,
		nil,
		[]byte(wrk.ConsensusState().ConsensusGroup()[0]),
		signature,
		int(bls.MtBlockBody),
		0,
		chainID,
		nil,
		nil,
		nil,
		currentPid,
		nil,
	)
	buff, _ := wrk.Marshalizer().Marshal(cnsMsg)
	msg := &p2pmocks.P2PMessageMock{
		DataField:      buff,
		PeerField:      currentPid,
		SignatureField: []byte("signature"),
	}

	msgID, err := wrk.ProcessReceivedMessage(msg, fromConnectedPeerId, &p2pmocks.MessengerStub{})
	time.Sleep(time.Second)
	assert.Equal(t, 1, len(wrk.ReceivedMessages()[bls.MtBlockBody]))
	assert.Nil(t, err)
	assert.Len(t, msgID, 0)

	msgID, err = wrk.ProcessReceivedMessage(msg, fromConnectedPeerId, &p2pmocks.MessengerStub{})
	time.Sleep(time.Second)
	assert.Equal(t, 1, len(wrk.ReceivedMessages()[bls.MtBlockBody]))
	assert.True(t, errors.Is(err, spos.ErrMessageTypeLimitReached))
	assert.Len(t, msgID, 0)

	msgID, err = wrk.ProcessReceivedMessage(msg, fromConnectedPeerId, &p2pmocks.MessengerStub{})
	time.Sleep(time.Second)
	assert.Equal(t, 1, len(wrk.ReceivedMessages()[bls.MtBlockBody]))
	assert.True(t, errors.Is(err, spos.ErrMessageTypeLimitReached))
	assert.Len(t, msgID, 0)
}

func TestWorker_ProcessReceivedMessageInvalidSignatureShouldErr(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	blk := &block.Body{}
	blkStr, _ := mock.MarshalizerMock{}.Marshal(blk)
	cnsMsg := consensus.NewConsensusMessage(
		nil,
		nil,
		blkStr,
		nil,
		[]byte(wrk.ConsensusState().ConsensusGroup()[0]),
		invalidSignature,
		int(bls.MtBlockBody),
		0,
		chainID,
		nil,
		nil,
		nil,
		currentPid,
		nil,
	)
	buff, _ := wrk.Marshalizer().Marshal(cnsMsg)
	msgID, err := wrk.ProcessReceivedMessage(
		&p2pmocks.P2PMessageMock{
			DataField:      buff,
			SignatureField: []byte("signature"),
		},
		fromConnectedPeerId,
		&p2pmocks.MessengerStub{},
	)
	time.Sleep(time.Second)

	assert.Equal(t, 0, len(wrk.ReceivedMessages()[bls.MtBlockBody]))
	assert.True(t, errors.Is(err, spos.ErrInvalidSignatureSize))
	assert.Nil(t, msgID)
}

func TestWorker_ProcessReceivedMessageReceivedMessageIsFromSelfShouldRetNilAndNotProcess(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	blk := &block.Body{}
	blkStr, _ := mock.MarshalizerMock{}.Marshal(blk)
	cnsMsg := consensus.NewConsensusMessage(
		nil,
		nil,
		blkStr,
		nil,
		[]byte(wrk.ConsensusState().SelfPubKey()),
		signature,
		int(bls.MtBlockBody),
		0,
		chainID,
		nil,
		nil,
		nil,
		currentPid,
		nil,
	)
	buff, _ := wrk.Marshalizer().Marshal(cnsMsg)
	msg := &p2pmocks.P2PMessageMock{
		DataField:      buff,
		PeerField:      currentPid,
		SignatureField: []byte("signature"),
	}
	msgID, err := wrk.ProcessReceivedMessage(msg, fromConnectedPeerId, &p2pmocks.MessengerStub{})
	time.Sleep(time.Second)

	assert.Equal(t, 0, len(wrk.ReceivedMessages()[bls.MtBlockBody]))
	assert.Nil(t, err)
	assert.Len(t, msgID, 0)
}

func TestWorker_ProcessReceivedMessageWhenRoundIsCanceledShouldRetNilAndNotProcess(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	wrk.ConsensusState().RoundCanceled = true
	blk := &block.Body{}
	blkStr, _ := mock.MarshalizerMock{}.Marshal(blk)
	cnsMsg := consensus.NewConsensusMessage(
		nil,
		nil,
		blkStr,
		nil,
		[]byte(wrk.ConsensusState().ConsensusGroup()[0]),
		signature,
		int(bls.MtBlockBody),
		0,
		chainID,
		nil,
		nil,
		nil,
		currentPid,
		nil,
	)
	buff, _ := wrk.Marshalizer().Marshal(cnsMsg)
	msg := &p2pmocks.P2PMessageMock{
		DataField:      buff,
		PeerField:      currentPid,
		SignatureField: []byte("signature"),
	}
	msgID, err := wrk.ProcessReceivedMessage(msg, fromConnectedPeerId, &p2pmocks.MessengerStub{})
	time.Sleep(time.Second)

	assert.Equal(t, 0, len(wrk.ReceivedMessages()[bls.MtBlockBody]))
	assert.Nil(t, err)
	assert.Len(t, msgID, 0)
}

func TestWorker_ProcessReceivedMessageWrongChainIDInProposedBlockShouldError(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	wrk.SetBlockProcessor(
		&testscommon.BlockProcessorStub{
			DecodeBlockHeaderCalled: func(dta []byte) data.HeaderHandler {
				return &testscommon.HeaderHandlerStub{
					CheckChainIDCalled: func(reference []byte) error {
						return spos.ErrInvalidChainID
					},
					GetPrevHashCalled: func() []byte {
						return make([]byte, 0)
					},
				}
			},
			RevertCurrentBlockCalled: func() {
			},
		},
	)

	hdr := &block.Header{ChainID: wrongChainID}
	hdrHash, _ := core.CalculateHash(mock.MarshalizerMock{}, &hashingMocks.HasherMock{}, hdr)
	cnsMsg := consensus.NewConsensusMessage(
		hdrHash,
		nil,
		nil,
		nil,
		[]byte(wrk.ConsensusState().ConsensusGroup()[0]),
		[]byte("sig"),
		int(bls.MtBlockHeader),
		0,
		wrongChainID,
		nil,
		nil,
		nil,
		currentPid,
		nil,
	)
	buff, _ := wrk.Marshalizer().Marshal(cnsMsg)
	msgID, err := wrk.ProcessReceivedMessage(
		&p2pmocks.P2PMessageMock{
			DataField:      buff,
			SignatureField: []byte("signature"),
		},
		fromConnectedPeerId,
		&p2pmocks.MessengerStub{},
	)
	time.Sleep(time.Second)

	assert.True(t, errors.Is(err, spos.ErrInvalidChainID))
	assert.Nil(t, msgID)
}

func TestWorker_ProcessReceivedMessageWithABadOriginatorShouldErr(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	wrk.SetBlockProcessor(
		&testscommon.BlockProcessorStub{
			DecodeBlockHeaderCalled: func(dta []byte) data.HeaderHandler {
				return &testscommon.HeaderHandlerStub{
					CheckChainIDCalled: func(reference []byte) error {
						return nil
					},
					GetPrevHashCalled: func() []byte {
						return make([]byte, 0)
					},
				}
			},
			RevertCurrentBlockCalled: func() {
			},
			DecodeBlockBodyCalled: func(dta []byte) data.BodyHandler {
				return nil
			},
		},
	)

	hdr := &block.Header{ChainID: chainID}
	hdrHash, _ := core.CalculateHash(mock.MarshalizerMock{}, &hashingMocks.HasherMock{}, hdr)
	hdrStr, _ := mock.MarshalizerMock{}.Marshal(hdr)
	cnsMsg := consensus.NewConsensusMessage(
		hdrHash,
		nil,
		nil,
		hdrStr,
		[]byte(wrk.ConsensusState().ConsensusGroup()[0]),
		signature,
		int(bls.MtBlockHeader),
		0,
		chainID,
		nil,
		nil,
		nil,
		currentPid,
		nil,
	)
	buff, _ := wrk.Marshalizer().Marshal(cnsMsg)
	msg := &p2pmocks.P2PMessageMock{
		DataField:      buff,
		PeerField:      "other originator",
		SignatureField: []byte("signature"),
	}
	msgID, err := wrk.ProcessReceivedMessage(msg, fromConnectedPeerId, &p2pmocks.MessengerStub{})
	time.Sleep(time.Second)

	assert.Equal(t, 0, len(wrk.ReceivedMessages()[bls.MtBlockHeader]))
	assert.True(t, errors.Is(err, spos.ErrOriginatorMismatch))
	assert.Nil(t, msgID)
}

func TestWorker_ProcessReceivedMessageWithHeaderAndWrongHash(t *testing.T) {
	t.Parallel()

	workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
	wrk, _ := spos.NewWorker(workerArgs)
	wrk.ConsensusState().SetHeader(&block.HeaderV2{})

	wrk.SetBlockProcessor(
		&testscommon.BlockProcessorStub{
			DecodeBlockHeaderCalled: func(dta []byte) data.HeaderHandler {
				return &testscommon.HeaderHandlerStub{
					CheckChainIDCalled: func(reference []byte) error {
						return nil
					},
					GetPrevHashCalled: func() []byte {
						return make([]byte, 0)
					},
				}
			},
			RevertCurrentBlockCalled: func() {
			},
			DecodeBlockBodyCalled: func(dta []byte) data.BodyHandler {
				return nil
			},
		},
	)

	hdr := &block.Header{ChainID: chainID}
	hdrHash := make([]byte, 32) // wrong hash
	hdrStr, _ := mock.MarshalizerMock{}.Marshal(hdr)
	cnsMsg := consensus.NewConsensusMessage(
		hdrHash,
		nil,
		nil,
		hdrStr,
		[]byte(wrk.ConsensusState().ConsensusGroup()[0]),
		signature,
		int(bls.MtBlockHeader),
		0,
		chainID,
		nil,
		nil,
		nil,
		currentPid,
		nil,
	)
	buff, _ := wrk.Marshalizer().Marshal(cnsMsg)
	msg := &p2pmocks.P2PMessageMock{
		DataField:      buff,
		PeerField:      currentPid,
		SignatureField: []byte("signature"),
	}
	msgID, err := wrk.ProcessReceivedMessage(msg, fromConnectedPeerId, &p2pmocks.MessengerStub{})
	time.Sleep(time.Second)

	assert.Equal(t, 0, len(wrk.ReceivedMessages()[bls.MtBlockHeader]))
	assert.ErrorIs(t, err, spos.ErrWrongHashForHeader)
	assert.Nil(t, msgID)
}

func TestWorker_ProcessReceivedMessageOkValsShouldWork(t *testing.T) {
	t.Parallel()

	workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
	expectedShardID := workerArgs.ShardCoordinator.SelfId()
	expectedPK := []byte(workerArgs.ConsensusState.ConsensusGroup()[0])
	wasUpdatePeerIDInfoCalled := false
	workerArgs.NetworkShardingCollector = &p2pmocks.NetworkShardingCollectorStub{
		UpdatePeerIDInfoCalled: func(pid core.PeerID, pk []byte, shardID uint32) {
			assert.Equal(t, currentPid, pid)
			assert.Equal(t, expectedPK, pk)
			assert.Equal(t, expectedShardID, shardID)
			wasUpdatePeerIDInfoCalled = true
		},
	}
	wrk, _ := spos.NewWorker(workerArgs)
	wrk.ConsensusState().SetHeader(&block.HeaderV2{})

	wrk.SetBlockProcessor(
		&testscommon.BlockProcessorStub{
			DecodeBlockHeaderCalled: func(dta []byte) data.HeaderHandler {
				return &testscommon.HeaderHandlerStub{
					CheckChainIDCalled: func(reference []byte) error {
						return nil
					},
					GetPrevHashCalled: func() []byte {
						return make([]byte, 0)
					},
				}
			},
			RevertCurrentBlockCalled: func() {
			},
			DecodeBlockBodyCalled: func(dta []byte) data.BodyHandler {
				return nil
			},
		},
	)

	hdr := &block.Header{ChainID: chainID}
	hdrHash, _ := core.CalculateHash(mock.MarshalizerMock{}, &hashingMocks.HasherMock{}, hdr)
	hdrStr, _ := mock.MarshalizerMock{}.Marshal(hdr)
	cnsMsg := consensus.NewConsensusMessage(
		hdrHash,
		nil,
		nil,
		hdrStr,
		[]byte(wrk.ConsensusState().ConsensusGroup()[0]),
		signature,
		int(bls.MtBlockHeader),
		0,
		chainID,
		nil,
		nil,
		nil,
		currentPid,
		nil,
	)
	buff, _ := wrk.Marshalizer().Marshal(cnsMsg)
	msg := &p2pmocks.P2PMessageMock{
		DataField:      buff,
		PeerField:      currentPid,
		SignatureField: []byte("signature"),
	}
	msgID, err := wrk.ProcessReceivedMessage(msg, fromConnectedPeerId, &p2pmocks.MessengerStub{})
	time.Sleep(time.Second)

	assert.Equal(t, 1, len(wrk.ReceivedMessages()[bls.MtBlockHeader]))
	assert.Nil(t, err)
	assert.True(t, wasUpdatePeerIDInfoCalled)
	assert.Len(t, msgID, 0)
}

func TestWorker_CheckSelfStateShouldErrMessageFromItself(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	cnsMsg := consensus.NewConsensusMessage(
		nil,
		nil,
		nil,
		nil,
		[]byte(wrk.ConsensusState().SelfPubKey()),
		nil,
		0,
		0,
		chainID,
		nil,
		nil,
		nil,
		currentPid,
		nil,
	)
	err := wrk.CheckSelfState(cnsMsg)
	assert.Equal(t, spos.ErrMessageFromItself, err)
}

func TestWorker_CheckSelfStateShouldErrRoundCanceled(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	wrk.ConsensusState().RoundCanceled = true
	cnsMsg := consensus.NewConsensusMessage(
		nil,
		nil,
		nil,
		nil,
		[]byte(wrk.ConsensusState().ConsensusGroup()[0]),
		nil,
		0,
		0,
		chainID,
		nil,
		nil,
		nil,
		currentPid,
		nil,
	)
	err := wrk.CheckSelfState(cnsMsg)
	assert.Equal(t, spos.ErrRoundCanceled, err)
}

func TestWorker_CheckSelfStateShouldNotErr(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	cnsMsg := consensus.NewConsensusMessage(
		nil,
		nil,
		nil,
		nil,
		[]byte(wrk.ConsensusState().ConsensusGroup()[0]),
		nil,
		0,
		0,
		chainID,
		nil,
		nil,
		nil,
		currentPid,
		nil,
	)
	err := wrk.CheckSelfState(cnsMsg)
	assert.Nil(t, err)
}

func TestWorker_CheckSignatureShouldReturnNilErr(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	blk := &block.Body{}
	blkStr, _ := mock.MarshalizerMock{}.Marshal(blk)
	cnsMsg := consensus.NewConsensusMessage(
		nil,
		nil,
		blkStr,
		nil,
		[]byte(wrk.ConsensusState().ConsensusGroup()[0]),
		[]byte("sig"),
		int(bls.MtBlockBody),
		0,
		chainID,
		nil,
		nil,
		nil,
		currentPid,
		nil,
	)
	err := wrk.CheckSignature(cnsMsg)

	assert.Nil(t, err)
}

func TestWorker_ExecuteMessagesShouldNotExecuteWhenConsensusDataIsNil(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	blk := &block.Body{}
	blkStr, _ := mock.MarshalizerMock{}.Marshal(blk)
	wrk.InitReceivedMessages()
	cnsMsg := consensus.NewConsensusMessage(
		nil,
		nil,
		blkStr,
		nil,
		[]byte(wrk.ConsensusState().ConsensusGroup()[0]),
		[]byte("sig"),
		int(bls.MtBlockBody),
		0,
		chainID,
		nil,
		nil,
		nil,
		currentPid,
		nil,
	)
	msgType := consensus.MessageType(cnsMsg.MsgType)
	cnsDataList := wrk.ReceivedMessages()[msgType]
	cnsDataList = append(cnsDataList, nil)
	wrk.SetReceivedMessages(msgType, cnsDataList)
	wrk.ExecuteMessage(cnsDataList)

	assert.Nil(t, wrk.ReceivedMessages()[msgType][0])
}

func TestWorker_ExecuteMessagesShouldNotExecuteWhenMessageIsForOtherRound(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	blk := &block.Body{}
	blkStr, _ := mock.MarshalizerMock{}.Marshal(blk)
	wrk.InitReceivedMessages()
	cnsMsg := consensus.NewConsensusMessage(
		nil,
		nil,
		blkStr,
		nil,
		[]byte(wrk.ConsensusState().ConsensusGroup()[0]),
		[]byte("sig"),
		int(bls.MtBlockBody),
		-1,
		chainID,
		nil,
		nil,
		nil,
		currentPid,
		nil,
	)
	msgType := consensus.MessageType(cnsMsg.MsgType)
	cnsDataList := wrk.ReceivedMessages()[msgType]
	cnsDataList = append(cnsDataList, cnsMsg)
	wrk.SetReceivedMessages(msgType, cnsDataList)
	wrk.ExecuteMessage(cnsDataList)

	assert.NotNil(t, wrk.ReceivedMessages()[msgType][0])
}

func TestWorker_ExecuteBlockBodyMessagesShouldNotExecuteWhenStartRoundIsNotFinished(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	blk := &block.Body{}
	blkStr, _ := mock.MarshalizerMock{}.Marshal(blk)
	wrk.InitReceivedMessages()
	cnsMsg := consensus.NewConsensusMessage(
		nil,
		nil,
		blkStr,
		nil,
		[]byte(wrk.ConsensusState().ConsensusGroup()[0]),
		[]byte("sig"),
		int(bls.MtBlockBody),
		0,
		chainID,
		nil,
		nil,
		nil,
		currentPid,
		nil,
	)
	msgType := consensus.MessageType(cnsMsg.MsgType)
	cnsDataList := wrk.ReceivedMessages()[msgType]
	cnsDataList = append(cnsDataList, cnsMsg)
	wrk.SetReceivedMessages(msgType, cnsDataList)
	wrk.ExecuteMessage(cnsDataList)

	assert.NotNil(t, wrk.ReceivedMessages()[msgType][0])
}

func TestWorker_ExecuteBlockHeaderMessagesShouldNotExecuteWhenStartRoundIsNotFinished(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	blk := &block.Body{}
	blkStr, _ := mock.MarshalizerMock{}.Marshal(blk)
	wrk.InitReceivedMessages()
	cnsMsg := consensus.NewConsensusMessage(
		nil,
		nil,
		blkStr,
		nil,
		[]byte(wrk.ConsensusState().ConsensusGroup()[0]),
		[]byte("sig"),
		int(bls.MtBlockHeader),
		0,
		chainID,
		nil,
		nil,
		nil,
		currentPid,
		nil,
	)
	msgType := consensus.MessageType(cnsMsg.MsgType)
	cnsDataList := wrk.ReceivedMessages()[msgType]
	cnsDataList = append(cnsDataList, cnsMsg)
	wrk.SetReceivedMessages(msgType, cnsDataList)
	wrk.ExecuteMessage(cnsDataList)

	assert.NotNil(t, wrk.ReceivedMessages()[msgType][0])
}

func TestWorker_ExecuteSignatureMessagesShouldNotExecuteWhenBlockIsNotFinished(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	blk := &block.Body{}
	blkStr, _ := mock.MarshalizerMock{}.Marshal(blk)
	wrk.InitReceivedMessages()
	cnsMsg := consensus.NewConsensusMessage(
		nil,
		nil,
		blkStr,
		nil,
		[]byte(wrk.ConsensusState().ConsensusGroup()[0]),
		[]byte("sig"),
		int(bls.MtSignature),
		0,
		chainID,
		nil,
		nil,
		nil,
		currentPid,
		nil,
	)
	msgType := consensus.MessageType(cnsMsg.MsgType)
	cnsDataList := wrk.ReceivedMessages()[msgType]
	cnsDataList = append(cnsDataList, cnsMsg)
	wrk.SetReceivedMessages(msgType, cnsDataList)
	wrk.ExecuteMessage(cnsDataList)

	assert.NotNil(t, wrk.ReceivedMessages()[msgType][0])
}

func TestWorker_ExecuteMessagesShouldExecute(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	wrk.StartWorking()
	blk := &block.Body{}
	blkStr, _ := mock.MarshalizerMock{}.Marshal(blk)
	wrk.InitReceivedMessages()
	cnsMsg := consensus.NewConsensusMessage(
		nil,
		nil,
		blkStr,
		nil,
		[]byte(wrk.ConsensusState().ConsensusGroup()[0]),
		[]byte("sig"),
		int(bls.MtBlockBody),
		0,
		chainID,
		nil,
		nil,
		nil,
		currentPid,
		nil,
	)
	msgType := consensus.MessageType(cnsMsg.MsgType)
	cnsDataList := wrk.ReceivedMessages()[msgType]
	cnsDataList = append(cnsDataList, cnsMsg)
	wrk.SetReceivedMessages(msgType, cnsDataList)
	wrk.ConsensusState().SetStatus(bls.SrStartRound, spos.SsFinished)
	wrk.ExecuteMessage(cnsDataList)

	assert.Nil(t, wrk.ReceivedMessages()[msgType][0])

	_ = wrk.Close()
}

func TestWorker_CheckChannelsShouldWork(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	wrk.StartWorking()
	wrk.AppendReceivedMessagesCalls(bls.MtBlockHeader, func(ctx context.Context, cnsMsg *consensus.Message) bool {
		_ = wrk.ConsensusState().SetJobDone(wrk.ConsensusState().ConsensusGroup()[0], bls.SrBlock, true)
		return true
	})
	rnd := wrk.RoundHandler()
	roundDuration := rnd.TimeDuration()
	rnd.UpdateRound(time.Now(), time.Now().Add(roundDuration))
	cnsGroup := wrk.ConsensusState().ConsensusGroup()
	hdr := &block.Header{}
	hdr.Nonce = 1
	hdr.TimeStamp = uint64(wrk.RoundHandler().TimeStamp().Unix())
	hdrStr, _ := mock.MarshalizerMock{}.Marshal(hdr)
	cnsMsg := consensus.NewConsensusMessage(
		nil,
		nil,
		nil,
		hdrStr,
		[]byte(cnsGroup[0]),
		[]byte("sig"),
		int(bls.MtBlockHeader),
		1,
		chainID,
		nil,
		nil,
		nil,
		currentPid,
		nil,
	)
	wrk.ExecuteMessageChannel() <- cnsMsg
	time.Sleep(1000 * time.Millisecond)
	isBlockJobDone, err := wrk.ConsensusState().JobDone(cnsGroup[0], bls.SrBlock)

	assert.Nil(t, err)
	assert.True(t, isBlockJobDone)

	_ = wrk.Close()
}

func TestWorker_ConvertHeaderToConsensusMessage(t *testing.T) {
	t.Parallel()

	t.Run("nil header should error", func(t *testing.T) {
		wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
		_, err := wrk.ConvertHeaderToConsensusMessage(nil)
		require.Equal(t, spos.ErrInvalidHeader, err)
	})
	t.Run("valid header v2 should not error", func(t *testing.T) {
		wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
		marshaller := wrk.Marshalizer()
		hdr := &block.HeaderV2{
			Header: &block.Header{
				Round: 100,
			},
		}

		hdrStr, _ := marshaller.Marshal(hdr)
		expectedConsensusMsg := &consensus.Message{
			Header:     hdrStr,
			MsgType:    int64(bls.MtBlockHeader),
			RoundIndex: 100,
		}

		message, err := wrk.ConvertHeaderToConsensusMessage(hdr)
		require.Nil(t, err)
		require.Equal(t, expectedConsensusMsg, message)
	})
	t.Run("valid header metaHeader should not error", func(t *testing.T) {
		wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
		marshaller := wrk.Marshalizer()
		hdr := &block.MetaBlock{
			Round: 100,
		}

		hdrStr, _ := marshaller.Marshal(hdr)
		expectedConsensusMsg := &consensus.Message{
			Header:     hdrStr,
			MsgType:    int64(bls.MtBlockHeader),
			RoundIndex: 100,
		}

		message, err := wrk.ConvertHeaderToConsensusMessage(hdr)
		require.Nil(t, err)
		require.Equal(t, expectedConsensusMsg, message)
	})
}

func TestWorker_StoredHeadersExecution(t *testing.T) {
	t.Parallel()

	hdr := &block.HeaderV2{
		Header: &block.Header{
			Round: 100,
		},
	}

	t.Run("Test stored headers before current round advances to same round should not finalize round", func(t *testing.T) {
		wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
		wrk.StartWorking()
		wrk.AddReceivedHeaderHandler(func(handler data.HeaderHandler) {
			_ = wrk.ConsensusState().SetJobDone(wrk.ConsensusState().ConsensusGroup()[0], bls.SrBlock, true)
		})

		roundIndex := &atomic.Int64{}
		roundIndex.Store(99)
		roundHandler := &consensusMocks.RoundHandlerMock{
			IndexCalled: func() int64 {
				return roundIndex.Load()
			},
		}
		wrk.SetRoundHandler(roundHandler)
		wrk.ConsensusState().SetRoundIndex(99)
		cnsGroup := wrk.ConsensusState().ConsensusGroup()

		wrk.BlockProcessor().(*testscommon.BlockProcessorStub).DecodeBlockHeaderCalled = func(dta []byte) data.HeaderHandler {
			return hdr
		}
		wrk.SetEnableEpochsHandler(&enableEpochsHandlerMock.EnableEpochsHandlerStub{
			IsFlagEnabledInEpochCalled: func(flag core.EnableEpochFlag, epoch uint32) bool {
				return true
			},
			IsFlagEnabledCalled: func(flag core.EnableEpochFlag) bool {
				return true
			},
		})

		wrk.ConsensusState().SetStatus(bls.SrStartRound, spos.SsFinished)
		wrk.AddFutureHeaderToProcessIfNeeded(hdr)
		time.Sleep(200 * time.Millisecond)
		wrk.ExecuteStoredMessages()
		time.Sleep(200 * time.Millisecond)

		isBlockJobDone, err := wrk.ConsensusState().JobDone(cnsGroup[0], bls.SrBlock)

		assert.Nil(t, err)
		assert.False(t, isBlockJobDone)

		_ = wrk.Close()
	})
	t.Run("Test stored headers should finalize round after roundIndex advances", func(t *testing.T) {
		wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
		wrk.StartWorking()
		wrk.AddReceivedHeaderHandler(func(handler data.HeaderHandler) {
			_ = wrk.ConsensusState().SetJobDone(wrk.ConsensusState().ConsensusGroup()[0], bls.SrBlock, true)
		})

		roundIndex := &atomic.Int64{}
		roundIndex.Store(99)
		roundHandler := &consensusMocks.RoundHandlerMock{
			IndexCalled: func() int64 {
				return roundIndex.Load()
			},
		}
		wrk.SetRoundHandler(roundHandler)

		wrk.ConsensusState().SetRoundIndex(99)
		cnsGroup := wrk.ConsensusState().ConsensusGroup()

		wrk.BlockProcessor().(*testscommon.BlockProcessorStub).DecodeBlockHeaderCalled = func(dta []byte) data.HeaderHandler {
			return hdr
		}
		wrk.SetEnableEpochsHandler(&enableEpochsHandlerMock.EnableEpochsHandlerStub{
			IsFlagEnabledInEpochCalled: func(flag core.EnableEpochFlag, epoch uint32) bool {
				return true
			},
			IsFlagEnabledCalled: func(flag core.EnableEpochFlag) bool {
				return true
			},
		})

		wrk.ConsensusState().SetStatus(bls.SrStartRound, spos.SsFinished)
		wrk.AddFutureHeaderToProcessIfNeeded(hdr)
		time.Sleep(200 * time.Millisecond)
		roundIndex.Store(100)
		wrk.ConsensusState().SetRoundIndex(100)
		wrk.ExecuteStoredMessages()
		time.Sleep(200 * time.Millisecond)

		isBlockJobDone, err := wrk.ConsensusState().JobDone(cnsGroup[0], bls.SrBlock)

		assert.Nil(t, err)
		assert.True(t, isBlockJobDone)

		_ = wrk.Close()
	})
	t.Run("Test stored meta headers should finalize round after roundIndex advances", func(t *testing.T) {
		hdr := &block.MetaBlock{
			Round: 100,
		}

		wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
		wrk.StartWorking()
		wrk.AddReceivedHeaderHandler(func(handler data.HeaderHandler) {
			_ = wrk.ConsensusState().SetJobDone(wrk.ConsensusState().ConsensusGroup()[0], bls.SrBlock, true)
		})

		roundIndex := &atomic.Int64{}
		roundIndex.Store(99)
		roundHandler := &consensusMocks.RoundHandlerMock{
			IndexCalled: func() int64 {
				return roundIndex.Load()
			},
		}
		wrk.SetRoundHandler(roundHandler)

		wrk.ConsensusState().SetRoundIndex(99)
		cnsGroup := wrk.ConsensusState().ConsensusGroup()

		wrk.BlockProcessor().(*testscommon.BlockProcessorStub).DecodeBlockHeaderCalled = func(dta []byte) data.HeaderHandler {
			return hdr
		}
		wrk.SetEnableEpochsHandler(&enableEpochsHandlerMock.EnableEpochsHandlerStub{
			IsFlagEnabledInEpochCalled: func(flag core.EnableEpochFlag, epoch uint32) bool {
				return true
			},
			IsFlagEnabledCalled: func(flag core.EnableEpochFlag) bool {
				return true
			},
		})

		wrk.ConsensusState().SetStatus(bls.SrStartRound, spos.SsFinished)
		wrk.AddFutureHeaderToProcessIfNeeded(hdr)
		time.Sleep(200 * time.Millisecond)
		roundIndex.Store(100)
		wrk.ConsensusState().SetRoundIndex(100)
		wrk.ExecuteStoredMessages()
		time.Sleep(200 * time.Millisecond)

		isBlockJobDone, err := wrk.ConsensusState().JobDone(cnsGroup[0], bls.SrBlock)

		assert.Nil(t, err)
		assert.True(t, isBlockJobDone)

		_ = wrk.Close()
	})
}

func TestWorker_ExtendShouldReturnWhenRoundIsCanceled(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	executed := false
	bootstrapperMock := &bootstrapperStubs.BootstrapperStub{
		GetNodeStateCalled: func() common.NodeState {
			return common.NsNotSynchronized
		},
		CreateAndCommitEmptyBlockCalled: func(shardForCurrentNode uint32) (data.BodyHandler, data.HeaderHandler, error) {
			executed = true
			return nil, nil, errors.New("error")
		},
	}
	wrk.SetBootstrapper(bootstrapperMock)
	wrk.ConsensusState().RoundCanceled = true
	wrk.Extend(0)

	assert.False(t, executed)
}

func TestWorker_ExtendShouldReturnWhenGetNodeStateNotReturnSynchronized(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	executed := false
	bootstrapperMock := &bootstrapperStubs.BootstrapperStub{
		GetNodeStateCalled: func() common.NodeState {
			return common.NsNotSynchronized
		},
		CreateAndCommitEmptyBlockCalled: func(shardForCurrentNode uint32) (data.BodyHandler, data.HeaderHandler, error) {
			executed = true
			return nil, nil, errors.New("error")
		},
	}
	wrk.SetBootstrapper(bootstrapperMock)
	wrk.Extend(0)

	assert.False(t, executed)
}

func TestWorker_ExtendShouldReturnWhenCreateEmptyBlockFail(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	executed := false
	bmm := &consensusMocks.BroadcastMessengerMock{
		BroadcastBlockCalled: func(handler data.BodyHandler, handler2 data.HeaderHandler) error {
			executed = true
			return nil
		},
	}
	wrk.SetBroadcastMessenger(bmm)
	bootstrapperMock := &bootstrapperStubs.BootstrapperStub{
		CreateAndCommitEmptyBlockCalled: func(shardForCurrentNode uint32) (data.BodyHandler, data.HeaderHandler, error) {
			return nil, nil, errors.New("error")
		}}
	wrk.SetBootstrapper(bootstrapperMock)
	wrk.Extend(0)

	assert.False(t, executed)
}

func TestWorker_ExtendShouldWorkAfterAWhile(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	executed := int32(0)
	blockProcessor := &testscommon.BlockProcessorStub{
		RevertCurrentBlockCalled: func() {
			atomic.AddInt32(&executed, 1)
		},
	}
	wrk.SetBlockProcessor(blockProcessor)
	wrk.ConsensusState().SetProcessingBlock(true)
	n := 10
	go func() {
		for n > 0 {
			time.Sleep(100 * time.Millisecond)
			n--
		}
		wrk.ConsensusState().SetProcessingBlock(false)
	}()
	wrk.Extend(1)

	assert.Equal(t, int32(1), atomic.LoadInt32(&executed))
	assert.Equal(t, 0, n)
}

func TestWorker_ExtendShouldWork(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	executed := int32(0)
	blockProcessor := &testscommon.BlockProcessorStub{
		RevertCurrentBlockCalled: func() {
			atomic.AddInt32(&executed, 1)
		},
	}
	wrk.SetBlockProcessor(blockProcessor)
	wrk.Extend(1)
	time.Sleep(1000 * time.Millisecond)

	assert.Equal(t, int32(1), atomic.LoadInt32(&executed))
}

func TestWorker_ExecuteStoredMessagesShouldWork(t *testing.T) {
	t.Parallel()
	wrk := *initWorker(&statusHandlerMock.AppStatusHandlerStub{})
	wrk.StartWorking()
	blk := &block.Body{}
	blkStr, _ := mock.MarshalizerMock{}.Marshal(blk)
	wrk.InitReceivedMessages()
	cnsMsg := consensus.NewConsensusMessage(
		nil,
		nil,
		blkStr,
		nil,
		[]byte(wrk.ConsensusState().ConsensusGroup()[0]),
		[]byte("sig"),
		int(bls.MtBlockBody),
		0,
		chainID,
		nil,
		nil,
		nil,
		currentPid,
		nil,
	)
	msgType := consensus.MessageType(cnsMsg.MsgType)
	cnsDataList := wrk.ReceivedMessages()[msgType]
	cnsDataList = append(cnsDataList, cnsMsg)
	wrk.SetReceivedMessages(msgType, cnsDataList)
	wrk.ConsensusState().SetStatus(bls.SrStartRound, spos.SsFinished)

	rcvMsg := wrk.ReceivedMessages()
	assert.Equal(t, 1, len(rcvMsg[msgType]))

	wrk.ExecuteStoredMessages()

	rcvMsg = wrk.ReceivedMessages()
	assert.Equal(t, 0, len(rcvMsg[msgType]))

	_ = wrk.Close()
}

func TestWorker_AppStatusHandlerNilShouldErr(t *testing.T) {
	t.Parallel()

	workerArgs := createDefaultWorkerArgs(nil)
	_, err := spos.NewWorker(workerArgs)

	assert.Equal(t, spos.ErrNilAppStatusHandler, err)
}

func TestWorker_ProcessReceivedMessageWrongHeaderShouldErr(t *testing.T) {
	t.Parallel()

	workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
	headerSigVerifier := &consensusMocks.HeaderSigVerifierMock{}
	headerSigVerifier.VerifyRandSeedCalled = func(header data.HeaderHandler) error {
		return process.ErrRandSeedDoesNotMatch
	}

	workerArgs.HeaderSigVerifier = headerSigVerifier
	wrk, _ := spos.NewWorker(workerArgs)
	wrk.ConsensusState().SetHeader(&block.HeaderV2{})

	hdr := &block.Header{}
	hdr.Nonce = 1
	hdr.TimeStamp = uint64(wrk.RoundHandler().TimeStamp().Unix())
	hdrStr, _ := mock.MarshalizerMock{}.Marshal(hdr)
	hdrHash := (&hashingMocks.HasherMock{}).Compute(string(hdrStr))
	cnsMsg := consensus.NewConsensusMessage(
		hdrHash,
		nil,
		nil,
		hdrStr,
		[]byte(wrk.ConsensusState().ConsensusGroup()[0]),
		signature,
		int(bls.MtBlockHeader),
		0,
		chainID,
		nil,
		nil,
		nil,
		currentPid,
		nil,
	)
	buff, _ := wrk.Marshalizer().Marshal(cnsMsg)
	time.Sleep(time.Second)
	msg := &p2pmocks.P2PMessageMock{
		DataField:      buff,
		PeerField:      currentPid,
		SignatureField: []byte("signature"),
	}
	msgID, err := wrk.ProcessReceivedMessage(msg, "", &p2pmocks.MessengerStub{})
	assert.True(t, errors.Is(err, spos.ErrInvalidHeader))
	assert.Nil(t, msgID)
}

func TestWorker_ProcessReceivedMessageWithSignature(t *testing.T) {
	t.Parallel()

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
		wrk, _ := spos.NewWorker(workerArgs)
		wrk.ConsensusState().SetHeader(&block.HeaderV2{})

		hdr := &block.Header{}
		hdr.Nonce = 1
		hdr.TimeStamp = uint64(wrk.RoundHandler().TimeStamp().Unix())
		hdrStr, _ := mock.MarshalizerMock{}.Marshal(hdr)
		hdrHash := (&hashingMocks.HasherMock{}).Compute(string(hdrStr))
		pubKey := []byte(wrk.ConsensusState().ConsensusGroup()[0])

		cnsMsg := consensus.NewConsensusMessage(
			hdrHash,
			bytes.Repeat([]byte("a"), SignatureSize),
			nil,
			nil,
			pubKey,
			bytes.Repeat([]byte("a"), SignatureSize),
			int(bls.MtSignature),
			0,
			chainID,
			nil,
			nil,
			nil,
			currentPid,
			nil,
		)
		buff, err := wrk.Marshalizer().Marshal(cnsMsg)
		require.Nil(t, err)

		time.Sleep(time.Second)
		msg := &p2pmocks.P2PMessageMock{
			DataField:      buff,
			PeerField:      currentPid,
			SignatureField: []byte("signature"),
		}
		msgID, err := wrk.ProcessReceivedMessage(msg, "", &p2pmocks.MessengerStub{})
		assert.Nil(t, err)
		assert.Len(t, msgID, 0)

		p2pMsgWithSignature, ok := wrk.ConsensusState().GetMessageWithSignature(string(pubKey))
		require.True(t, ok)
		require.Equal(t, msg, p2pMsgWithSignature)
	})
}

func TestWorker_ProcessReceivedMessageWithInvalidSigners(t *testing.T) {
	t.Parallel()

	workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
	cntCheckKnownInvalidSignersCalled := 0
	workerArgs.InvalidSignersCache = &consensusMocks.InvalidSignersCacheMock{
		CheckKnownInvalidSignersCalled: func(headerHash []byte, invalidSigners []byte) bool {
			cntCheckKnownInvalidSignersCalled++
			return cntCheckKnownInvalidSignersCalled > 1
		},
	}
	workerArgs.AntifloodHandler = &mock.P2PAntifloodHandlerStub{
		CanProcessMessageCalled: func(message p2p.MessageP2P, fromConnectedPeer core.PeerID) error {
			return nil
		},
		CanProcessMessagesOnTopicCalled: func(peer core.PeerID, topic string, numMessages uint32, totalSize uint64, sequence []byte) error {
			return nil
		},
		BlacklistPeerCalled: func(peer core.PeerID, reason string, duration time.Duration) {
			require.Fail(t, "should have not been called")
		},
	}
	workerArgs.EnableEpochsHandler = &enableEpochsHandlerMock.EnableEpochsHandlerStub{
		IsFlagEnabledCalled: func(flag core.EnableEpochFlag) bool {
			return true
		},
	}
	wrk, _ := spos.NewWorker(workerArgs)
	wrk.ConsensusState().SetHeader(&block.HeaderV2{})

	hdr := &block.Header{}
	hdr.Nonce = 1
	hdr.TimeStamp = uint64(wrk.RoundHandler().TimeStamp().Unix())
	hdrStr, _ := mock.MarshalizerMock{}.Marshal(hdr)
	hdrHash := (&hashingMocks.HasherMock{}).Compute(string(hdrStr))
	pubKey := []byte(wrk.ConsensusState().ConsensusGroup()[0])

	invalidSigners := []byte("invalid signers")
	cnsMsg := consensus.NewConsensusMessage(
		hdrHash,
		nil,
		nil,
		nil,
		pubKey,
		bytes.Repeat([]byte("a"), SignatureSize),
		int(bls.MtInvalidSigners),
		0,
		chainID,
		nil,
		nil,
		nil,
		currentPid,
		invalidSigners,
	)
	buff, err := wrk.Marshalizer().Marshal(cnsMsg)
	require.Nil(t, err)

	msg := &p2pmocks.P2PMessageMock{
		DataField:      buff,
		PeerField:      currentPid,
		SignatureField: []byte("signature"),
	}

	// first call should be ok
	msgID, err := wrk.ProcessReceivedMessage(msg, "", &p2pmocks.MessengerStub{})
	require.Nil(t, err)
	require.Len(t, msgID, 0)

	// reset the received messages to allow a second one of the same type
	wrk.ResetConsensusMessages()

	// second call should see this message as already received and return error
	msgID, err = wrk.ProcessReceivedMessage(msg, "", &p2pmocks.MessengerStub{})
	require.Equal(t, spos.ErrInvalidSignersAlreadyReceived, err)
	require.Nil(t, msgID)
}

func TestWorker_ReceivedHeader(t *testing.T) {
	t.Parallel()

	t.Run("nil header should early exit", func(t *testing.T) {
		t.Parallel()

		workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
		wrk, _ := spos.NewWorker(workerArgs)
		wrk.ConsensusState().SetHeader(&block.HeaderV2{})

		rcvHeaderHandler := func(header data.HeaderHandler) {
			require.Fail(t, "should have not been called")
		}
		wrk.AddReceivedHeaderHandler(rcvHeaderHandler)
		wrk.ReceivedHeader(nil, nil)
	})
	t.Run("unprocessable header should early exit", func(t *testing.T) {
		t.Parallel()

		workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
		wrk, _ := spos.NewWorker(workerArgs)
		wrk.ConsensusState().SetHeader(&block.HeaderV2{})

		rcvHeaderHandler := func(header data.HeaderHandler) {
			require.Fail(t, "should have not been called")
		}
		wrk.AddReceivedHeaderHandler(rcvHeaderHandler)
		wrk.ReceivedHeader(&block.Header{
			ShardID: workerArgs.ShardCoordinator.SelfId(),
			Round:   uint64(workerArgs.RoundHandler.Index() + 1), // should not process this one
		}, nil)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		wasSetUInt64ValueCalled := false
		setStringValueCnt := 0
		appStatusHandler := &statusHandlerMock.AppStatusHandlerStub{
			SetUInt64ValueHandler: func(key string, value uint64) {
				require.Equal(t, common.MetricReceivedProposedBlock, key)
				wasSetUInt64ValueCalled = true
			},
			SetStringValueHandler: func(key string, value string) {
				setStringValueCnt++
				if key != common.MetricRedundancyIsMainActive &&
					key != common.MetricRedundancyStepInReason {
					require.Fail(t, "unexpected key for SetStringValue")
				}
			},
		}
		workerArgs := createDefaultWorkerArgs(appStatusHandler)
		workerArgs.EnableEpochsHandler = &enableEpochsHandlerMock.EnableEpochsHandlerStub{
			IsFlagEnabledInEpochCalled: func(flag core.EnableEpochFlag, epoch uint32) bool {
				return flag == common.AndromedaFlag
			},
		}
		wrk, _ := spos.NewWorker(workerArgs)
		wrk.ConsensusState().SetHeader(&block.HeaderV2{})

		wasHandlerCalled := false
		rcvHeaderHandler := func(header data.HeaderHandler) {
			wasHandlerCalled = true
		}
		wrk.AddReceivedHeaderHandler(rcvHeaderHandler)
		wrk.ReceivedHeader(&block.Header{
			ShardID: workerArgs.ShardCoordinator.SelfId(),
			Round:   uint64(workerArgs.RoundHandler.Index()),
		}, nil)
		require.True(t, wasHandlerCalled)

		wrk.RemoveAllReceivedHeaderHandlers() // coverage only
		require.True(t, wasSetUInt64ValueCalled)
		require.Equal(t, 2, setStringValueCnt)
	})
}

func TestWorker_ReceivedProof(t *testing.T) {
	t.Parallel()

	t.Run("nil proof should early exit", func(t *testing.T) {
		t.Parallel()

		workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
		wrk, _ := spos.NewWorker(workerArgs)
		wrk.ConsensusState().SetHeader(&block.HeaderV2{})

		rcvProofHandler := func(proof consensus.ProofHandler) {
			require.Fail(t, "should have not been called")
		}
		wrk.AddReceivedProofHandler(rcvProofHandler)
		wrk.ReceivedProof(nil)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		workerArgs := createDefaultWorkerArgs(&statusHandlerMock.AppStatusHandlerStub{})
		wrk, _ := spos.NewWorker(workerArgs)
		wrk.ConsensusState().SetHeader(&block.HeaderV2{})

		wasHandlerCalled := false
		rcvProofHandler := func(proof consensus.ProofHandler) {
			wasHandlerCalled = true
		}
		wrk.AddReceivedProofHandler(rcvProofHandler)
		wrk.ReceivedProof(&block.HeaderProof{})
		require.True(t, wasHandlerCalled)
	})
}
