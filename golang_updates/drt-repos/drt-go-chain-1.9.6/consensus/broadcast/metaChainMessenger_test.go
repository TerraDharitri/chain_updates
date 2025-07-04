package broadcast_test

import (
	"bytes"
	"sync"
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/consensus"
	"github.com/TerraDharitri/drt-go-chain/consensus/broadcast"
	"github.com/TerraDharitri/drt-go-chain/consensus/broadcast/shared"
	"github.com/TerraDharitri/drt-go-chain/consensus/mock"
	"github.com/TerraDharitri/drt-go-chain/consensus/spos"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	consensusMock "github.com/TerraDharitri/drt-go-chain/testscommon/consensus"
	"github.com/TerraDharitri/drt-go-chain/testscommon/hashingMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/p2pmocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/pool"
)

var nodePkBytes = []byte("node public key bytes")

func createDefaultMetaChainArgs() broadcast.MetaChainMessengerArgs {
	marshalizerMock := &mock.MarshalizerMock{}
	messengerMock := &p2pmocks.MessengerStub{}
	shardCoordinatorMock := &mock.ShardCoordinatorMock{}
	singleSignerMock := &mock.SingleSignerMock{}
	hasher := &hashingMocks.HasherMock{}
	headersSubscriber := &pool.HeadersPoolStub{}
	interceptorsContainer := createInterceptorContainer()
	peerSigHandler := &mock.PeerSignatureHandler{Signer: singleSignerMock}
	alarmScheduler := &testscommon.AlarmSchedulerStub{}
	delayedBroadcaster := &consensusMock.DelayedBroadcasterMock{}

	return broadcast.MetaChainMessengerArgs{
		CommonMessengerArgs: broadcast.CommonMessengerArgs{
			Marshalizer:                marshalizerMock,
			Hasher:                     hasher,
			Messenger:                  messengerMock,
			ShardCoordinator:           shardCoordinatorMock,
			PeerSignatureHandler:       peerSigHandler,
			HeadersSubscriber:          headersSubscriber,
			InterceptorsContainer:      interceptorsContainer,
			MaxValidatorDelayCacheSize: 2,
			MaxDelayCacheSize:          2,
			AlarmScheduler:             alarmScheduler,
			KeysHandler:                &testscommon.KeysHandlerStub{},
			DelayedBroadcaster:         delayedBroadcaster,
		},
	}
}

func TestMetaChainMessenger_NewMetaChainMessengerNilMarshalizerShouldFail(t *testing.T) {
	args := createDefaultMetaChainArgs()
	args.Marshalizer = nil
	mcm, err := broadcast.NewMetaChainMessenger(args)

	assert.Nil(t, mcm)
	assert.Equal(t, spos.ErrNilMarshalizer, err)
}

func TestMetaChainMessenger_NewMetaChainMessengerNilMessengerShouldFail(t *testing.T) {
	args := createDefaultMetaChainArgs()
	args.Messenger = nil
	mcm, err := broadcast.NewMetaChainMessenger(args)

	assert.Nil(t, mcm)
	assert.Equal(t, spos.ErrNilMessenger, err)
}

func TestMetaChainMessenger_NewMetaChainMessengerNilShardCoordinatorShouldFail(t *testing.T) {
	args := createDefaultMetaChainArgs()
	args.ShardCoordinator = nil
	mcm, err := broadcast.NewMetaChainMessenger(args)

	assert.Nil(t, mcm)
	assert.Equal(t, spos.ErrNilShardCoordinator, err)
}

func TestMetaChainMessenger_NewMetaChainMessengerNilPeerSignatureHandlerShouldFail(t *testing.T) {
	args := createDefaultMetaChainArgs()
	args.PeerSignatureHandler = nil
	mcm, err := broadcast.NewMetaChainMessenger(args)

	assert.Nil(t, mcm)
	assert.Equal(t, spos.ErrNilPeerSignatureHandler, err)
}

func TestMetaChainMessenger_NilKeysHandlerShouldError(t *testing.T) {
	args := createDefaultMetaChainArgs()
	args.KeysHandler = nil
	mcm, err := broadcast.NewMetaChainMessenger(args)

	assert.Nil(t, mcm)
	assert.Equal(t, broadcast.ErrNilKeysHandler, err)
}

func TestMetaChainMessenger_NilDelayedBroadcasterShouldError(t *testing.T) {
	args := createDefaultMetaChainArgs()
	args.DelayedBroadcaster = nil
	scm, err := broadcast.NewMetaChainMessenger(args)

	assert.Nil(t, scm)
	assert.Equal(t, broadcast.ErrNilDelayedBroadcaster, err)
}
func TestMetaChainMessenger_NewMetaChainMessengerShouldWork(t *testing.T) {
	args := createDefaultMetaChainArgs()
	mcm, err := broadcast.NewMetaChainMessenger(args)

	assert.NotNil(t, mcm)
	assert.Equal(t, nil, err)
	assert.False(t, mcm.IsInterfaceNil())
}

func TestMetaChainMessenger_BroadcastBlockShouldErrNilMetaHeader(t *testing.T) {
	args := createDefaultMetaChainArgs()
	mcm, _ := broadcast.NewMetaChainMessenger(args)

	err := mcm.BroadcastBlock(newTestBlockBody(), nil)
	assert.Equal(t, spos.ErrNilMetaHeader, err)
}

func TestMetaChainMessenger_BroadcastBlockShouldErrMockMarshalizer(t *testing.T) {
	marshalizer := &mock.MarshalizerMock{
		Fail: true,
	}
	args := createDefaultMetaChainArgs()
	args.Marshalizer = marshalizer
	mcm, _ := broadcast.NewMetaChainMessenger(args)

	err := mcm.BroadcastBlock(newTestBlockBody(), &block.MetaBlock{})
	assert.Equal(t, mock.ErrMockMarshalizer, err)
}

func TestMetaChainMessenger_BroadcastBlockShouldWork(t *testing.T) {
	messenger := &p2pmocks.MessengerStub{
		BroadcastCalled: func(topic string, buff []byte) {
		},
	}
	args := createDefaultMetaChainArgs()
	args.Messenger = messenger
	mcm, _ := broadcast.NewMetaChainMessenger(args)

	err := mcm.BroadcastBlock(newTestBlockBody(), &block.MetaBlock{})
	assert.Nil(t, err)
}

func TestMetaChainMessenger_BroadcastMiniBlocksShouldWork(t *testing.T) {
	args := createDefaultMetaChainArgs()
	mcm, _ := broadcast.NewMetaChainMessenger(args)

	err := mcm.BroadcastMiniBlocks(nil, []byte("pk bytes"))
	assert.Nil(t, err)
}

func TestMetaChainMessenger_BroadcastTransactionsShouldWork(t *testing.T) {
	args := createDefaultMetaChainArgs()
	mcm, _ := broadcast.NewMetaChainMessenger(args)

	err := mcm.BroadcastTransactions(nil, []byte("pk bytes"))
	assert.Nil(t, err)
}

func TestMetaChainMessenger_BroadcastHeaderNilHeaderShouldErr(t *testing.T) {
	args := createDefaultMetaChainArgs()
	mcm, _ := broadcast.NewMetaChainMessenger(args)

	err := mcm.BroadcastHeader(nil, []byte("pk bytes"))
	assert.Equal(t, spos.ErrNilHeader, err)
}

func TestMetaChainMessenger_BroadcastHeaderOkHeaderShouldWork(t *testing.T) {
	channelBroadcastCalled := make(chan bool, 1)
	channelBroadcastUsingPrivateKeyCalled := make(chan bool, 1)

	messenger := &p2pmocks.MessengerStub{
		BroadcastCalled: func(topic string, buff []byte) {
			channelBroadcastCalled <- true
		},
		BroadcastUsingPrivateKeyCalled: func(topic string, buff []byte, pid core.PeerID, skBytes []byte) {
			channelBroadcastUsingPrivateKeyCalled <- true
		},
	}
	args := createDefaultMetaChainArgs()
	args.KeysHandler = &testscommon.KeysHandlerStub{
		IsOriginalPublicKeyOfTheNodeCalled: func(pkBytes []byte) bool {
			return bytes.Equal(pkBytes, nodePkBytes)
		},
	}
	args.Messenger = messenger
	mcm, _ := broadcast.NewMetaChainMessenger(args)

	hdr := block.Header{
		Nonce: 10,
	}

	t.Run("original public key of the node", func(t *testing.T) {
		err := mcm.BroadcastHeader(&hdr, nodePkBytes)
		assert.Nil(t, err)

		wasCalled := false
		select {
		case <-channelBroadcastCalled:
			wasCalled = true
		case <-time.After(time.Millisecond * 100):
		}

		assert.Nil(t, err)
		assert.True(t, wasCalled)
	})
	t.Run("managed key", func(t *testing.T) {
		err := mcm.BroadcastHeader(&hdr, []byte("managed key"))
		assert.Nil(t, err)

		wasCalled := false
		select {
		case <-channelBroadcastUsingPrivateKeyCalled:
			wasCalled = true
		case <-time.After(time.Millisecond * 100):
		}

		assert.Nil(t, err)
		assert.True(t, wasCalled)
	})

}

func TestMetaChainMessenger_BroadcastBlockDataLeader(t *testing.T) {
	countersBroadcast := make(map[string]int)
	mutCounters := &sync.Mutex{}

	messengerMock := &p2pmocks.MessengerStub{
		BroadcastCalled: func(topic string, buff []byte) {
			mutCounters.Lock()
			countersBroadcast[broadcastMethodPrefix+topic]++
			mutCounters.Unlock()
		},
		BroadcastUsingPrivateKeyCalled: func(topic string, buff []byte, pid core.PeerID, skBytes []byte) {
			mutCounters.Lock()
			countersBroadcast[broadcastUsingPrivateKeyCalledMethodPrefix+topic]++
			mutCounters.Unlock()
		},
	}

	args := createDefaultMetaChainArgs()
	args.KeysHandler = &testscommon.KeysHandlerStub{
		IsOriginalPublicKeyOfTheNodeCalled: func(pkBytes []byte) bool {
			return bytes.Equal(pkBytes, nodePkBytes)
		},
	}
	args.Messenger = messengerMock
	mcm, _ := broadcast.NewMetaChainMessenger(args)

	miniBlocks := map[uint32][]byte{0: []byte("mbs data1"), 1: []byte("mbs data2")}
	transactions := map[string][][]byte{"topic1": {[]byte("txdata1"), []byte("txdata2")}, "topic2": {[]byte("txdata3")}}

	t.Run("original public key of the node", func(t *testing.T) {
		mutCounters.Lock()
		countersBroadcast = make(map[string]int)
		mutCounters.Unlock()

		err := mcm.BroadcastBlockDataLeader(nil, miniBlocks, transactions, nodePkBytes)
		require.Nil(t, err)
		sleepTime := common.ExtraDelayBetweenBroadcastMbsAndTxs +
			common.ExtraDelayForBroadcastBlockInfo +
			time.Millisecond*100
		time.Sleep(sleepTime)

		mutCounters.Lock()
		defer mutCounters.Unlock()

		numBroadcast := countersBroadcast[broadcastMethodPrefix+"txBlockBodies_0"]
		numBroadcast += countersBroadcast[broadcastMethodPrefix+"txBlockBodies_0_1"]
		assert.Equal(t, len(miniBlocks), numBroadcast)

		numBroadcast = countersBroadcast[broadcastMethodPrefix+"topic1"]
		numBroadcast += countersBroadcast[broadcastMethodPrefix+"topic2"]
		assert.Equal(t, len(transactions), numBroadcast)
	})
	t.Run("managed key", func(t *testing.T) {
		mutCounters.Lock()
		countersBroadcast = make(map[string]int)
		mutCounters.Unlock()

		err := mcm.BroadcastBlockDataLeader(nil, miniBlocks, transactions, []byte("pk bytes"))
		require.Nil(t, err)
		sleepTime := common.ExtraDelayBetweenBroadcastMbsAndTxs +
			common.ExtraDelayForBroadcastBlockInfo +
			time.Millisecond*100
		time.Sleep(sleepTime)

		mutCounters.Lock()
		defer mutCounters.Unlock()

		numBroadcast := countersBroadcast[broadcastUsingPrivateKeyCalledMethodPrefix+"txBlockBodies_0"]
		numBroadcast += countersBroadcast[broadcastUsingPrivateKeyCalledMethodPrefix+"txBlockBodies_0_1"]
		assert.Equal(t, len(miniBlocks), numBroadcast)

		numBroadcast = countersBroadcast[broadcastUsingPrivateKeyCalledMethodPrefix+"topic1"]
		numBroadcast += countersBroadcast[broadcastUsingPrivateKeyCalledMethodPrefix+"topic2"]
		assert.Equal(t, len(transactions), numBroadcast)
	})
}

func TestMetaChainMessenger_Close(t *testing.T) {
	t.Parallel()

	args := createDefaultMetaChainArgs()
	closeCalled := false
	delayedBroadcaster := &consensusMock.DelayedBroadcasterMock{
		CloseCalled: func() {
			closeCalled = true
		},
	}
	args.DelayedBroadcaster = delayedBroadcaster

	mcm, _ := broadcast.NewMetaChainMessenger(args)
	require.NotNil(t, mcm)
	mcm.Close()
	assert.True(t, closeCalled)
}

func TestMetaChainMessenger_PrepareBroadcastHeaderValidator(t *testing.T) {
	t.Parallel()

	t.Run("Nil header", func(t *testing.T) {
		t.Parallel()

		args := createDefaultMetaChainArgs()
		delayedBroadcaster := &consensusMock.DelayedBroadcasterMock{
			SetHeaderForValidatorCalled: func(vData *shared.ValidatorHeaderBroadcastData) error {
				require.Fail(t, "SetHeaderForValidator should not be called")
				return nil
			},
		}
		args.DelayedBroadcaster = delayedBroadcaster

		mcm, _ := broadcast.NewMetaChainMessenger(args)
		require.NotNil(t, mcm)
		mcm.PrepareBroadcastHeaderValidator(nil, make(map[uint32][]byte), make(map[string][][]byte), 0, make([]byte, 0))
	})
	t.Run("Err on core.CalculateHash", func(t *testing.T) {
		t.Parallel()

		args := createDefaultMetaChainArgs()
		delayedBroadcaster := &consensusMock.DelayedBroadcasterMock{
			SetHeaderForValidatorCalled: func(vData *shared.ValidatorHeaderBroadcastData) error {
				require.Fail(t, "SetHeaderForValidator should not be called")
				return nil
			},
		}
		args.DelayedBroadcaster = delayedBroadcaster

		header := &block.Header{}
		mcm, _ := broadcast.NewMetaChainMessenger(args)
		require.NotNil(t, mcm)
		mcm.SetMarshalizerMeta(nil)
		mcm.PrepareBroadcastHeaderValidator(header, make(map[uint32][]byte), make(map[string][][]byte), 0, make([]byte, 0))
	})
	t.Run("Err on SetHeaderForValidator", func(t *testing.T) {
		t.Parallel()

		args := createDefaultMetaChainArgs()
		checkVarModified := false
		delayedBroadcaster := &consensusMock.DelayedBroadcasterMock{
			SetHeaderForValidatorCalled: func(vData *shared.ValidatorHeaderBroadcastData) error {
				checkVarModified = true
				return expectedErr
			},
		}
		args.DelayedBroadcaster = delayedBroadcaster

		mcm, _ := broadcast.NewMetaChainMessenger(args)
		require.NotNil(t, mcm)
		header := &block.Header{}
		mcm.PrepareBroadcastHeaderValidator(header, make(map[uint32][]byte), make(map[string][][]byte), 0, make([]byte, 0))
		assert.True(t, checkVarModified)
	})
}

func TestMetaChainMessenger_BroadcastBlock(t *testing.T) {
	t.Parallel()

	t.Run("Err nil blockData", func(t *testing.T) {
		args := createDefaultMetaChainArgs()
		mcm, _ := broadcast.NewMetaChainMessenger(args)
		require.NotNil(t, mcm)
		err := mcm.BroadcastBlock(nil, nil)
		assert.NotNil(t, err)
	})
}

func TestMetaChainMessenger_NewMetaChainMessengerFailSetBroadcast(t *testing.T) {
	t.Parallel()

	args := createDefaultMetaChainArgs()
	varModified := false
	delayedBroadcaster := &consensusMock.DelayedBroadcasterMock{
		SetBroadcastHandlersCalled: func(
			mbBroadcast func(mbData map[uint32][]byte, pkBytes []byte) error,
			txBroadcast func(txData map[string][][]byte, pkBytes []byte) error,
			headerBroadcast func(header data.HeaderHandler, pkBytes []byte) error,
			consensusMessageBroadcast func(message *consensus.Message) error) error {
			varModified = true
			return expectedErr
		},
	}
	args.DelayedBroadcaster = delayedBroadcaster

	mcm, err := broadcast.NewMetaChainMessenger(args)
	assert.Nil(t, mcm)
	assert.NotNil(t, err)
	assert.True(t, varModified)
}
