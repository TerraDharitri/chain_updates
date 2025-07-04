package broadcast_test

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/consensus"
	testscommonConsensus "github.com/TerraDharitri/drt-go-chain/testscommon/consensus"
	"github.com/TerraDharitri/drt-go-chain/testscommon/pool"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/atomic"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/stretchr/testify/assert"

	"github.com/TerraDharitri/drt-go-chain/consensus/broadcast"
	"github.com/TerraDharitri/drt-go-chain/consensus/broadcast/shared"
	"github.com/TerraDharitri/drt-go-chain/consensus/mock"
	"github.com/TerraDharitri/drt-go-chain/consensus/spos"
	"github.com/TerraDharitri/drt-go-chain/p2p"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/factory"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/hashingMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/p2pmocks"
)

var expectedErr = errors.New("expected error")

func createDelayData(prefix string) ([]byte, *block.Header, map[uint32][]byte, map[string][][]byte) {
	miniblocks := make(map[uint32][]byte)
	receiverShardID := uint32(1)
	miniblocks[receiverShardID] = []byte(prefix + "miniblock data")

	transactions := make(map[string][][]byte)
	topic := "txBlockBodies_0_1"
	transactions[topic] = [][]byte{
		[]byte(prefix + "tx0"),
		[]byte(prefix + "tx1"),
	}
	headerHash := []byte(prefix + "header hash")
	header := &block.Header{
		Round:        0,
		PrevRandSeed: []byte(prefix),
	}

	return headerHash, header, miniblocks, transactions
}

func createInterceptorContainer() process.InterceptorsContainer {
	return &testscommon.InterceptorsContainerStub{
		GetCalled: func(topic string) (process.Interceptor, error) {
			return &testscommon.InterceptorStub{
				ProcessReceivedMessageCalled: func(message p2p.MessageP2P) ([]byte, error) {
					return nil, nil
				},
			}, nil
		},
	}
}

func createDefaultShardChainArgs() broadcast.ShardChainMessengerArgs {
	marshalizerMock := &mock.MarshalizerMock{}
	hasher := &hashingMocks.HasherMock{}
	messengerMock := &p2pmocks.MessengerStub{}
	shardCoordinatorMock := &mock.ShardCoordinatorMock{}
	singleSignerMock := &mock.SingleSignerMock{}
	headersSubscriber := &pool.HeadersPoolStub{}
	interceptorsContainer := createInterceptorContainer()
	peerSigHandler := &mock.PeerSignatureHandler{
		Signer: singleSignerMock,
	}
	alarmScheduler := &testscommon.AlarmSchedulerStub{}
	delayedBroadcaster := &testscommonConsensus.DelayedBroadcasterMock{}

	return broadcast.ShardChainMessengerArgs{
		CommonMessengerArgs: broadcast.CommonMessengerArgs{
			Marshalizer:                marshalizerMock,
			Hasher:                     hasher,
			Messenger:                  messengerMock,
			ShardCoordinator:           shardCoordinatorMock,
			PeerSignatureHandler:       peerSigHandler,
			HeadersSubscriber:          headersSubscriber,
			InterceptorsContainer:      interceptorsContainer,
			MaxDelayCacheSize:          1,
			MaxValidatorDelayCacheSize: 1,
			AlarmScheduler:             alarmScheduler,
			KeysHandler:                &testscommon.KeysHandlerStub{},
			DelayedBroadcaster:         delayedBroadcaster,
		},
	}
}

func newBlockWithEmptyMiniblock() *block.Body {
	return &block.Body{
		MiniBlocks: []*block.MiniBlock{
			{
				TxHashes:        [][]byte{},
				ReceiverShardID: 0,
				SenderShardID:   0,
				Type:            0,
			},
		},
	}
}

func TestShardChainMessenger_NewShardChainMessengerNilMarshalizerShouldFail(t *testing.T) {
	args := createDefaultShardChainArgs()
	args.Marshalizer = nil

	scm, err := broadcast.NewShardChainMessenger(args)

	assert.Nil(t, scm)
	assert.Equal(t, spos.ErrNilMarshalizer, err)
}

func TestShardChainMessenger_NewShardChainMessengerNilMessengerShouldFail(t *testing.T) {
	args := createDefaultShardChainArgs()
	args.Messenger = nil
	scm, err := broadcast.NewShardChainMessenger(args)

	assert.Nil(t, scm)
	assert.Equal(t, spos.ErrNilMessenger, err)
}

func TestShardChainMessenger_NewShardChainMessengerNilShardCoordinatorShouldFail(t *testing.T) {
	args := createDefaultShardChainArgs()
	args.ShardCoordinator = nil
	scm, err := broadcast.NewShardChainMessenger(args)

	assert.Nil(t, scm)
	assert.Equal(t, spos.ErrNilShardCoordinator, err)
}

func TestShardChainMessenger_NewShardChainMessengerNilPeerSignatureHandlerShouldFail(t *testing.T) {
	args := createDefaultShardChainArgs()
	args.PeerSignatureHandler = nil
	scm, err := broadcast.NewShardChainMessenger(args)

	assert.Nil(t, scm)
	assert.Equal(t, spos.ErrNilPeerSignatureHandler, err)
}

func TestShardChainMessenger_NewShardChainMessengerNilInterceptorsContainerShouldFail(t *testing.T) {
	args := createDefaultShardChainArgs()
	args.InterceptorsContainer = nil
	scm, err := broadcast.NewShardChainMessenger(args)

	assert.Nil(t, scm)
	assert.Equal(t, spos.ErrNilInterceptorsContainer, err)
}

func TestShardChainMessenger_NewShardChainMessengerNilHeadersSubscriberShouldFail(t *testing.T) {
	args := createDefaultShardChainArgs()
	args.HeadersSubscriber = nil
	scm, err := broadcast.NewShardChainMessenger(args)

	assert.Nil(t, scm)
	assert.Equal(t, spos.ErrNilHeadersSubscriber, err)
}

func TestShardChainMessenger_NilDelayedBroadcasterShouldError(t *testing.T) {
	args := createDefaultShardChainArgs()
	args.DelayedBroadcaster = nil
	scm, err := broadcast.NewShardChainMessenger(args)

	assert.Nil(t, scm)
	assert.Equal(t, broadcast.ErrNilDelayedBroadcaster, err)
}

func TestShardChainMessenger_NilKeysHandlerShouldError(t *testing.T) {
	args := createDefaultShardChainArgs()
	args.KeysHandler = nil
	scm, err := broadcast.NewShardChainMessenger(args)

	assert.Nil(t, scm)
	assert.Equal(t, broadcast.ErrNilKeysHandler, err)
}

func TestShardChainMessenger_NewShardChainMessengerShouldWork(t *testing.T) {
	args := createDefaultShardChainArgs()
	scm, err := broadcast.NewShardChainMessenger(args)

	assert.NotNil(t, scm)
	assert.Equal(t, nil, err)
	assert.False(t, scm.IsInterfaceNil())
}

func TestShardChainMessenger_NewShardChainMessengerShouldErr(t *testing.T) {

	args := createDefaultShardChainArgs()
	args.DelayedBroadcaster = &testscommonConsensus.DelayedBroadcasterMock{
		SetBroadcastHandlersCalled: func(
			mbBroadcast func(mbData map[uint32][]byte, pkBytes []byte) error,
			txBroadcast func(txData map[string][][]byte, pkBytes []byte) error,
			headerBroadcast func(header data.HeaderHandler, pkBytes []byte) error,
			consensusMessageBroadcast func(message *consensus.Message) error,
		) error {
			return expectedErr
		}}

	_, err := broadcast.NewShardChainMessenger(args)

	assert.Equal(t, expectedErr, err)

}

func TestShardChainMessenger_BroadcastBlockShouldErrNilBody(t *testing.T) {
	args := createDefaultShardChainArgs()
	scm, _ := broadcast.NewShardChainMessenger(args)

	err := scm.BroadcastBlock(nil, &block.Header{})
	assert.Equal(t, spos.ErrNilBody, err)
}

func TestShardChainMessenger_BroadcastBlockShouldErrNilHeader(t *testing.T) {
	args := createDefaultShardChainArgs()
	scm, _ := broadcast.NewShardChainMessenger(args)

	err := scm.BroadcastBlock(newTestBlockBody(), nil)
	assert.Equal(t, spos.ErrNilHeader, err)
}

func TestShardChainMessenger_BroadcastBlockShouldErrMiniBlockEmpty(t *testing.T) {
	args := createDefaultShardChainArgs()
	scm, _ := broadcast.NewShardChainMessenger(args)

	err := scm.BroadcastBlock(newBlockWithEmptyMiniblock(), &block.Header{})
	assert.Equal(t, data.ErrMiniBlockEmpty, err)
}

func TestShardChainMessenger_BroadcastBlockShouldErrMockMarshalizer(t *testing.T) {
	marshalizer := mock.MarshalizerMock{
		Fail: true,
	}
	args := createDefaultShardChainArgs()
	args.Marshalizer = marshalizer
	scm, _ := broadcast.NewShardChainMessenger(args)

	err := scm.BroadcastBlock(newTestBlockBody(), &block.Header{})
	assert.Equal(t, mock.ErrMockMarshalizer, err)
}

func TestShardChainMessenger_BroadcastBlockShouldWork(t *testing.T) {
	messenger := &p2pmocks.MessengerStub{
		BroadcastCalled: func(topic string, buff []byte) {
		},
	}
	args := createDefaultShardChainArgs()
	args.Messenger = messenger
	scm, _ := broadcast.NewShardChainMessenger(args)

	err := scm.BroadcastBlock(newTestBlockBody(), &block.Header{})
	assert.Nil(t, err)
}

func TestShardChainMessenger_BroadcastMiniBlocksShouldBeDone(t *testing.T) {
	channelBroadcastCalled := make(chan bool, 100)
	channelBroadcastUsingPrivateKeyCalled := make(chan bool, 100)

	messenger := &p2pmocks.MessengerStub{
		BroadcastCalled: func(topic string, buff []byte) {
			channelBroadcastCalled <- true
		},
		BroadcastUsingPrivateKeyCalled: func(topic string, buff []byte, pid core.PeerID, skBytes []byte) {
			channelBroadcastUsingPrivateKeyCalled <- true
		},
	}
	args := createDefaultShardChainArgs()
	args.Messenger = messenger
	args.KeysHandler = &testscommon.KeysHandlerStub{
		IsOriginalPublicKeyOfTheNodeCalled: func(pkBytes []byte) bool {
			return bytes.Equal(nodePkBytes, pkBytes)
		},
	}
	scm, _ := broadcast.NewShardChainMessenger(args)

	miniBlocks := make(map[uint32][]byte)
	miniBlocks[0] = make([]byte, 0)
	miniBlocks[1] = make([]byte, 0)
	miniBlocks[2] = make([]byte, 0)
	miniBlocks[3] = make([]byte, 0)

	t.Run("original public key of the node", func(t *testing.T) {
		err := scm.BroadcastMiniBlocks(miniBlocks, nodePkBytes)

		called := 0
		for i := 0; i < 4; i++ {
			select {
			case <-channelBroadcastCalled:
				called++
			case <-time.After(time.Millisecond * 100):
				break
			}
		}

		assert.Nil(t, err)
		assert.Equal(t, 4, called)
	})
	t.Run("managed key", func(t *testing.T) {
		err := scm.BroadcastMiniBlocks(miniBlocks, []byte("managed key"))

		called := 0
		for i := 0; i < 4; i++ {
			select {
			case <-channelBroadcastUsingPrivateKeyCalled:
				called++
			case <-time.After(time.Millisecond * 100):
				break
			}
		}

		assert.Nil(t, err)
		assert.Equal(t, 4, called)
	})
}

func TestShardChainMessenger_BroadcastTransactionsShouldNotBeCalled(t *testing.T) {
	channelCalled := make(chan bool, 1)

	messenger := &p2pmocks.MessengerStub{
		BroadcastCalled: func(topic string, buff []byte) {
			channelCalled <- true
		},
	}
	args := createDefaultShardChainArgs()
	args.Messenger = messenger
	scm, _ := broadcast.NewShardChainMessenger(args)

	transactions := make(map[string][][]byte)
	err := scm.BroadcastTransactions(transactions, []byte("pk bytes"))

	wasCalled := false
	select {
	case <-channelCalled:
		wasCalled = true
	case <-time.After(time.Millisecond * 100):
	}

	assert.Nil(t, err)
	assert.False(t, wasCalled)

	transactions[factory.TransactionTopic] = make([][]byte, 0)
	err = scm.BroadcastTransactions(transactions, []byte("pk bytes"))

	wasCalled = false
	select {
	case <-channelCalled:
		wasCalled = true
	case <-time.After(time.Millisecond * 100):
	}

	assert.Nil(t, err)
	assert.False(t, wasCalled)
}

func TestShardChainMessenger_BroadcastTransactionsShouldBeCalled(t *testing.T) {
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

	args := createDefaultShardChainArgs()
	args.Messenger = messenger
	args.KeysHandler = &testscommon.KeysHandlerStub{
		IsOriginalPublicKeyOfTheNodeCalled: func(pkBytes []byte) bool {
			return bytes.Equal(pkBytes, nodePkBytes)
		},
	}
	scm, _ := broadcast.NewShardChainMessenger(args)

	transactions := make(map[string][][]byte)
	txs := make([][]byte, 0)
	txs = append(txs, []byte(""))
	transactions[factory.TransactionTopic] = txs
	t.Run("original public key of the node", func(t *testing.T) {
		err := scm.BroadcastTransactions(transactions, nodePkBytes)

		wasCalled := false
		for i := 0; i < 4; i++ {
			select {
			case <-channelBroadcastCalled:
				wasCalled = true
			case <-time.After(time.Millisecond * 100):
				break
			}
		}

		assert.Nil(t, err)
		assert.True(t, wasCalled)
	})
	t.Run("managed key", func(t *testing.T) {
		err := scm.BroadcastTransactions(transactions, []byte("managed key"))

		wasCalled := false
		for i := 0; i < 4; i++ {
			select {
			case <-channelBroadcastUsingPrivateKeyCalled:
				wasCalled = true
			case <-time.After(time.Millisecond * 100):
				break
			}
		}

		assert.Nil(t, err)
		assert.True(t, wasCalled)
	})
}

func TestShardChainMessenger_BroadcastHeaderNilHeaderShouldErr(t *testing.T) {
	args := createDefaultShardChainArgs()
	scm, _ := broadcast.NewShardChainMessenger(args)

	err := scm.BroadcastHeader(nil, []byte("pk bytes"))
	assert.Equal(t, spos.ErrNilHeader, err)
}

func TestShardChainMessenger_BroadcastHeaderShouldErr(t *testing.T) {
	marshalizer := mock.MarshalizerMock{
		Fail: true,
	}

	args := createDefaultShardChainArgs()
	args.Marshalizer = marshalizer
	scm, _ := broadcast.NewShardChainMessenger(args)

	err := scm.BroadcastHeader(&block.MetaBlock{Nonce: 10}, []byte("pk bytes"))
	assert.Equal(t, mock.ErrMockMarshalizer, err)
}

func TestShardChainMessenger_BroadcastHeaderShouldWork(t *testing.T) {
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
	args := createDefaultShardChainArgs()
	args.KeysHandler = &testscommon.KeysHandlerStub{
		IsOriginalPublicKeyOfTheNodeCalled: func(pkBytes []byte) bool {
			return bytes.Equal(pkBytes, nodePkBytes)
		},
	}
	args.Messenger = messenger
	scm, _ := broadcast.NewShardChainMessenger(args)

	hdr := &block.MetaBlock{Nonce: 10}
	t.Run("original public key of the node", func(t *testing.T) {
		err := scm.BroadcastHeader(hdr, nodePkBytes)

		wasCalled := false
		for i := 0; i < 4; i++ {
			select {
			case <-channelBroadcastCalled:
				wasCalled = true
			case <-time.After(time.Millisecond * 100):
				break
			}
		}

		assert.Nil(t, err)
		assert.True(t, wasCalled)
	})
	t.Run("managed key", func(t *testing.T) {
		err := scm.BroadcastHeader(hdr, []byte("managed key"))

		wasCalled := false
		for i := 0; i < 4; i++ {
			select {
			case <-channelBroadcastUsingPrivateKeyCalled:
				wasCalled = true
			case <-time.After(time.Millisecond * 100):
				break
			}
		}

		assert.Nil(t, err)
		assert.True(t, wasCalled)
	})
}

func TestShardChainMessenger_BroadcastBlockDataLeaderNilHeaderShouldErr(t *testing.T) {
	args := createDefaultShardChainArgs()
	scm, _ := broadcast.NewShardChainMessenger(args)

	_, _, miniblocks, transactions := createDelayData("1")

	err := scm.BroadcastBlockDataLeader(nil, miniblocks, transactions, []byte("pk bytes"))
	assert.Equal(t, spos.ErrNilHeader, err)
}

func TestShardChainMessenger_BroadcastBlockDataLeaderNilMiniblocksShouldReturnNil(t *testing.T) {
	args := createDefaultShardChainArgs()
	scm, _ := broadcast.NewShardChainMessenger(args)

	_, header, _, transactions := createDelayData("1")

	err := scm.BroadcastBlockDataLeader(header, nil, transactions, []byte("pk bytes"))
	assert.Nil(t, err)
}

func TestShardChainMessenger_BroadcastBlockDataLeaderShouldErr(t *testing.T) {
	marshalizer := mock.MarshalizerMock{
		Fail: true,
	}

	args := createDefaultShardChainArgs()
	args.Marshalizer = marshalizer

	scm, _ := broadcast.NewShardChainMessenger(args)

	_, header, miniblocks, transactions := createDelayData("1")

	err := scm.BroadcastBlockDataLeader(header, miniblocks, transactions, []byte("pk bytes"))
	assert.Equal(t, mock.ErrMockMarshalizer, err)
}

func TestShardChainMessenger_BroadcastBlockDataLeaderShouldErrDelayedBroadcaster(t *testing.T) {

	args := createDefaultShardChainArgs()

	args.DelayedBroadcaster = &testscommonConsensus.DelayedBroadcasterMock{
		SetLeaderDataCalled: func(data *shared.DelayedBroadcastData) error {
			return expectedErr
		}}

	scm, _ := broadcast.NewShardChainMessenger(args)
	require.NotNil(t, scm)

	_, header, miniblocks, transactions := createDelayData("1")

	err := scm.BroadcastBlockDataLeader(header, miniblocks, transactions, []byte("pk bytes"))

	assert.Equal(t, expectedErr, err)
}

func TestShardChainMessenger_BroadcastBlockDataLeaderShouldTriggerWaitingDelayedMessage(t *testing.T) {
	broadcastWasCalled := atomic.Flag{}
	broadcastUsingPrivateKeyWasCalled := atomic.Flag{}
	messenger := &p2pmocks.MessengerStub{
		BroadcastCalled: func(topic string, buff []byte) {
			broadcastWasCalled.SetValue(true)
		},
		BroadcastUsingPrivateKeyCalled: func(topic string, buff []byte, pid core.PeerID, skBytes []byte) {
			broadcastUsingPrivateKeyWasCalled.SetValue(true)
		},
	}
	args := createDefaultShardChainArgs()
	args.Messenger = messenger
	args.KeysHandler = &testscommon.KeysHandlerStub{
		IsOriginalPublicKeyOfTheNodeCalled: func(pkBytes []byte) bool {
			return bytes.Equal(pkBytes, nodePkBytes)
		},
	}
	argsDelayedBroadcaster := broadcast.ArgsDelayedBlockBroadcaster{
		InterceptorsContainer: args.InterceptorsContainer,
		HeadersSubscriber:     args.HeadersSubscriber,
		ShardCoordinator:      args.ShardCoordinator,
		LeaderCacheSize:       args.MaxDelayCacheSize,
		ValidatorCacheSize:    args.MaxDelayCacheSize,
		AlarmScheduler:        args.AlarmScheduler,
	}

	// Using real component in order to properly simulate the expected behavior
	args.DelayedBroadcaster, _ = broadcast.NewDelayedBlockBroadcaster(&argsDelayedBroadcaster)

	scm, _ := broadcast.NewShardChainMessenger(args)

	t.Run("original public key of the node", func(t *testing.T) {
		_, header, miniBlocksMarshalled, transactions := createDelayData("1")
		err := scm.BroadcastBlockDataLeader(header, miniBlocksMarshalled, transactions, nodePkBytes)
		time.Sleep(10 * time.Millisecond)
		assert.Nil(t, err)
		assert.False(t, broadcastWasCalled.IsSet())

		broadcastWasCalled.Reset()
		_, header2, miniBlocksMarshalled2, transactions2 := createDelayData("2")
		err = scm.BroadcastBlockDataLeader(header2, miniBlocksMarshalled2, transactions2, nodePkBytes)
		time.Sleep(10 * time.Millisecond)
		assert.Nil(t, err)
		assert.True(t, broadcastWasCalled.IsSet())
	})
	t.Run("managed key", func(t *testing.T) {
		_, header, miniBlocksMarshalled, transactions := createDelayData("1")
		err := scm.BroadcastBlockDataLeader(header, miniBlocksMarshalled, transactions, []byte("managed key"))
		time.Sleep(10 * time.Millisecond)
		assert.Nil(t, err)
		assert.False(t, broadcastUsingPrivateKeyWasCalled.IsSet())

		broadcastWasCalled.Reset()
		_, header2, miniBlocksMarshalled2, transactions2 := createDelayData("2")
		err = scm.BroadcastBlockDataLeader(header2, miniBlocksMarshalled2, transactions2, []byte("managed key"))
		time.Sleep(10 * time.Millisecond)
		assert.Nil(t, err)
		assert.True(t, broadcastUsingPrivateKeyWasCalled.IsSet())
	})
}

func TestShardChainMessenger_PrepareBroadcastHeaderValidatorShouldFailHeaderNil(t *testing.T) {

	pkBytes := make([]byte, 32)
	args := createDefaultShardChainArgs()

	args.DelayedBroadcaster = &testscommonConsensus.DelayedBroadcasterMock{
		SetHeaderForValidatorCalled: func(vData *shared.ValidatorHeaderBroadcastData) error {
			require.Fail(t, "SetHeaderForValidator should not be called")
			return nil
		}}

	scm, _ := broadcast.NewShardChainMessenger(args)
	require.NotNil(t, scm)

	scm.PrepareBroadcastHeaderValidator(nil, nil, nil, 1, pkBytes)
}

func TestShardChainMessenger_PrepareBroadcastHeaderValidatorShouldFailCalculateHashErr(t *testing.T) {

	pkBytes := make([]byte, 32)
	headerMock := &testscommon.HeaderHandlerStub{}

	args := createDefaultShardChainArgs()

	args.DelayedBroadcaster = &testscommonConsensus.DelayedBroadcasterMock{
		SetHeaderForValidatorCalled: func(vData *shared.ValidatorHeaderBroadcastData) error {
			require.Fail(t, "SetHeaderForValidator should not be called")
			return nil
		}}

	args.Marshalizer = &testscommon.MarshallerStub{MarshalCalled: func(obj interface{}) ([]byte, error) {
		return nil, expectedErr
	}}

	scm, _ := broadcast.NewShardChainMessenger(args)
	require.NotNil(t, scm)

	scm.PrepareBroadcastHeaderValidator(headerMock, nil, nil, 1, pkBytes)
}

func TestShardChainMessenger_PrepareBroadcastHeaderValidatorShouldWork(t *testing.T) {

	pkBytes := make([]byte, 32)
	headerMock := &testscommon.HeaderHandlerStub{}

	args := createDefaultShardChainArgs()

	varSetHeaderForValidatorCalled := false

	args.DelayedBroadcaster = &testscommonConsensus.DelayedBroadcasterMock{
		SetHeaderForValidatorCalled: func(vData *shared.ValidatorHeaderBroadcastData) error {
			varSetHeaderForValidatorCalled = true
			return nil
		}}

	args.Marshalizer = &testscommon.MarshallerStub{MarshalCalled: func(obj interface{}) ([]byte, error) {
		return nil, nil
	}}
	args.Hasher = &testscommon.HasherStub{ComputeCalled: func(s string) []byte {
		return nil
	}}

	scm, _ := broadcast.NewShardChainMessenger(args)
	require.NotNil(t, scm)

	scm.PrepareBroadcastHeaderValidator(headerMock, nil, nil, 1, pkBytes)

	assert.True(t, varSetHeaderForValidatorCalled)
}

func TestShardChainMessenger_PrepareBroadcastBlockDataValidatorShouldFailHeaderNil(t *testing.T) {

	pkBytes := make([]byte, 32)
	args := createDefaultShardChainArgs()

	args.DelayedBroadcaster = &testscommonConsensus.DelayedBroadcasterMock{
		SetValidatorDataCalled: func(data *shared.DelayedBroadcastData) error {
			require.Fail(t, "SetValidatorData should not be called")
			return nil
		}}

	scm, _ := broadcast.NewShardChainMessenger(args)
	require.NotNil(t, scm)

	scm.PrepareBroadcastBlockDataValidator(nil, nil, nil, 1, pkBytes)
}

func TestShardChainMessenger_PrepareBroadcastBlockDataValidatorShouldFailMiniBlocksLenZero(t *testing.T) {

	pkBytes := make([]byte, 32)
	miniBlocks := make(map[uint32][]byte)
	headerMock := &testscommon.HeaderHandlerStub{}

	args := createDefaultShardChainArgs()

	args.DelayedBroadcaster = &testscommonConsensus.DelayedBroadcasterMock{
		SetValidatorDataCalled: func(data *shared.DelayedBroadcastData) error {
			require.Fail(t, "SetValidatorData should not be called")
			return nil
		}}

	scm, _ := broadcast.NewShardChainMessenger(args)
	require.NotNil(t, scm)

	scm.PrepareBroadcastBlockDataValidator(headerMock, miniBlocks, nil, 1, pkBytes)
}

func TestShardChainMessenger_PrepareBroadcastBlockDataValidatorShouldFailCalculateHashErr(t *testing.T) {

	pkBytes := make([]byte, 32)
	miniBlocks := map[uint32][]byte{1: {}}
	headerMock := &testscommon.HeaderHandlerStub{}

	args := createDefaultShardChainArgs()

	args.DelayedBroadcaster = &testscommonConsensus.DelayedBroadcasterMock{
		SetValidatorDataCalled: func(data *shared.DelayedBroadcastData) error {
			require.Fail(t, "SetValidatorData should not be called")
			return nil
		}}

	args.Marshalizer = &testscommon.MarshallerStub{
		MarshalCalled: func(obj interface{}) ([]byte, error) {
			return nil, expectedErr
		},
	}

	scm, _ := broadcast.NewShardChainMessenger(args)
	require.NotNil(t, scm)

	scm.PrepareBroadcastBlockDataValidator(headerMock, miniBlocks, nil, 1, pkBytes)
}

func TestShardChainMessenger_PrepareBroadcastBlockDataValidatorShouldWork(t *testing.T) {

	pkBytes := make([]byte, 32)
	miniBlocks := map[uint32][]byte{1: {}}
	headerMock := &testscommon.HeaderHandlerStub{}

	args := createDefaultShardChainArgs()

	varSetValidatorDataCalled := false
	args.DelayedBroadcaster = &testscommonConsensus.DelayedBroadcasterMock{
		SetValidatorDataCalled: func(data *shared.DelayedBroadcastData) error {
			varSetValidatorDataCalled = true
			return nil
		}}

	args.Marshalizer = &testscommon.MarshallerStub{
		MarshalCalled: func(obj interface{}) ([]byte, error) {
			return nil, nil
		},
	}

	args.Hasher = &testscommon.HasherStub{
		ComputeCalled: func(s string) []byte {
			return nil
		},
	}

	scm, _ := broadcast.NewShardChainMessenger(args)
	require.NotNil(t, scm)

	scm.PrepareBroadcastBlockDataValidator(headerMock, miniBlocks, nil, 1, pkBytes)

	assert.True(t, varSetValidatorDataCalled)
}

func TestShardChainMessenger_CloseShouldWork(t *testing.T) {

	args := createDefaultShardChainArgs()

	varCloseCalled := false
	args.DelayedBroadcaster = &testscommonConsensus.DelayedBroadcasterMock{
		CloseCalled: func() {
			varCloseCalled = true
		},
	}

	scm, _ := broadcast.NewShardChainMessenger(args)
	require.NotNil(t, scm)

	scm.Close()
	assert.True(t, varCloseCalled)

}
