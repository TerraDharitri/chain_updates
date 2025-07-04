package resolvers_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data/batch"
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	"github.com/stretchr/testify/assert"

	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever/mock"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever/resolvers"
	"github.com/TerraDharitri/drt-go-chain/p2p"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/p2pmocks"
	storageStubs "github.com/TerraDharitri/drt-go-chain/testscommon/storage"
)

var connectedPeerId = core.PeerID("connected peer id")

func createMockArgTxResolver() resolvers.ArgTxResolver {
	return resolvers.ArgTxResolver{
		ArgBaseResolver: createMockArgBaseResolver(),
		TxPool:          testscommon.NewShardedDataStub(),
		TxStorage:       &storageStubs.StorerStub{},
		DataPacker:      &mock.DataPackerStub{},
	}
}

func TestNewTxResolver_NilResolverShouldErr(t *testing.T) {
	t.Parallel()

	arg := createMockArgTxResolver()
	arg.SenderResolver = nil
	txRes, err := resolvers.NewTxResolver(arg)

	assert.Equal(t, dataRetriever.ErrNilResolverSender, err)
	assert.Nil(t, txRes)
}

func TestNewTxResolver_NilTxPoolShouldErr(t *testing.T) {
	t.Parallel()

	arg := createMockArgTxResolver()
	arg.TxPool = nil
	txRes, err := resolvers.NewTxResolver(arg)

	assert.Equal(t, dataRetriever.ErrNilTxDataPool, err)
	assert.Nil(t, txRes)
}

func TestNewTxResolver_NilTxStorageShouldErr(t *testing.T) {
	t.Parallel()

	arg := createMockArgTxResolver()
	arg.TxStorage = nil
	txRes, err := resolvers.NewTxResolver(arg)

	assert.Equal(t, dataRetriever.ErrNilTxStorage, err)
	assert.Nil(t, txRes)
}

func TestNewTxResolver_NilMarshalizerShouldErr(t *testing.T) {
	t.Parallel()

	arg := createMockArgTxResolver()
	arg.Marshaller = nil
	txRes, err := resolvers.NewTxResolver(arg)

	assert.Equal(t, dataRetriever.ErrNilMarshalizer, err)
	assert.Nil(t, txRes)
}

func TestNewTxResolver_NilDataPackerShouldErr(t *testing.T) {
	t.Parallel()

	arg := createMockArgTxResolver()
	arg.DataPacker = nil
	txRes, err := resolvers.NewTxResolver(arg)

	assert.Equal(t, dataRetriever.ErrNilDataPacker, err)
	assert.Nil(t, txRes)
}

func TestNewTxResolver_NilAntifloodHandlerShouldErr(t *testing.T) {
	t.Parallel()

	arg := createMockArgTxResolver()
	arg.AntifloodHandler = nil
	txRes, err := resolvers.NewTxResolver(arg)

	assert.Equal(t, dataRetriever.ErrNilAntifloodHandler, err)
	assert.Nil(t, txRes)
}

func TestNewTxResolver_NilThrottlerShouldErr(t *testing.T) {
	t.Parallel()

	arg := createMockArgTxResolver()
	arg.Throttler = nil
	txRes, err := resolvers.NewTxResolver(arg)

	assert.Equal(t, dataRetriever.ErrNilThrottler, err)
	assert.Nil(t, txRes)
}

func TestNewTxResolver_OkValsShouldWork(t *testing.T) {
	t.Parallel()

	arg := createMockArgTxResolver()
	txRes, err := resolvers.NewTxResolver(arg)

	assert.Nil(t, err)
	assert.False(t, check.IfNil(txRes))
}

func TestTxResolver_ProcessReceivedMessageCanProcessMessageErrorsShouldErr(t *testing.T) {
	t.Parallel()

	arg := createMockArgTxResolver()
	arg.AntifloodHandler = &mock.P2PAntifloodHandlerStub{
		CanProcessMessageCalled: func(message p2p.MessageP2P, fromConnectedPeer core.PeerID) error {
			return expectedErr
		},
		CanProcessMessagesOnTopicCalled: func(peer core.PeerID, topic string, numMessages uint32, totalSize uint64, sequence []byte) error {
			return nil
		},
	}
	txRes, _ := resolvers.NewTxResolver(arg)

	msgID, err := txRes.ProcessReceivedMessage(&p2pmocks.P2PMessageMock{}, connectedPeerId, &p2pmocks.MessengerStub{})

	assert.True(t, errors.Is(err, expectedErr))
	assert.False(t, arg.Throttler.(*mock.ThrottlerStub).StartWasCalled())
	assert.False(t, arg.Throttler.(*mock.ThrottlerStub).EndWasCalled())
	assert.Nil(t, msgID)
}

func TestTxResolver_ProcessReceivedMessageNilMessageShouldErr(t *testing.T) {
	t.Parallel()

	arg := createMockArgTxResolver()
	txRes, _ := resolvers.NewTxResolver(arg)

	msgID, err := txRes.ProcessReceivedMessage(nil, connectedPeerId, &p2pmocks.MessengerStub{})

	assert.Equal(t, dataRetriever.ErrNilMessage, err)
	assert.False(t, arg.Throttler.(*mock.ThrottlerStub).StartWasCalled())
	assert.False(t, arg.Throttler.(*mock.ThrottlerStub).EndWasCalled())
	assert.Nil(t, msgID)
}

func TestTxResolver_ProcessReceivedMessageWrongTypeShouldErr(t *testing.T) {
	t.Parallel()

	arg := createMockArgTxResolver()
	txRes, _ := resolvers.NewTxResolver(arg)

	data, _ := arg.Marshaller.Marshal(&dataRetriever.RequestData{Type: dataRetriever.NonceType, Value: []byte("aaa")})

	msg := &p2pmocks.P2PMessageMock{DataField: data}

	msgID, err := txRes.ProcessReceivedMessage(msg, connectedPeerId, &p2pmocks.MessengerStub{})

	assert.True(t, errors.Is(err, dataRetriever.ErrRequestTypeNotImplemented))
	assert.True(t, arg.Throttler.(*mock.ThrottlerStub).StartWasCalled())
	assert.True(t, arg.Throttler.(*mock.ThrottlerStub).EndWasCalled())
	assert.Nil(t, msgID)
}

func TestTxResolver_ProcessReceivedMessageNilValueShouldErr(t *testing.T) {
	t.Parallel()

	arg := createMockArgTxResolver()
	txRes, _ := resolvers.NewTxResolver(arg)

	data, _ := arg.Marshaller.Marshal(&dataRetriever.RequestData{Type: dataRetriever.HashType, Value: nil})

	msg := &p2pmocks.P2PMessageMock{DataField: data}

	msgID, err := txRes.ProcessReceivedMessage(msg, connectedPeerId, &p2pmocks.MessengerStub{})

	assert.Equal(t, dataRetriever.ErrNilValue, err)
	assert.True(t, arg.Throttler.(*mock.ThrottlerStub).StartWasCalled())
	assert.True(t, arg.Throttler.(*mock.ThrottlerStub).EndWasCalled())
	assert.Nil(t, msgID)
}

func TestTxResolver_ProcessReceivedMessageFoundInTxPoolShouldSearchAndSend(t *testing.T) {
	t.Parallel()

	marshalizer := &mock.MarshalizerMock{}
	searchWasCalled := false
	sendWasCalled := false
	txReturned := &transaction.Transaction{
		Nonce: 10,
	}
	txPool := testscommon.NewShardedDataStub()
	txPool.SearchFirstDataCalled = func(key []byte) (value interface{}, ok bool) {
		if bytes.Equal([]byte("aaa"), key) {
			searchWasCalled = true
			return txReturned, true
		}

		return nil, false
	}

	arg := createMockArgTxResolver()
	arg.SenderResolver = &mock.TopicResolverSenderStub{
		SendCalled: func(buff []byte, peer core.PeerID, source p2p.MessageHandler) error {
			sendWasCalled = true
			return nil
		},
	}
	arg.TxPool = txPool
	txRes, _ := resolvers.NewTxResolver(arg)

	data, _ := marshalizer.Marshal(&dataRetriever.RequestData{Type: dataRetriever.HashType, Value: []byte("aaa")})

	msg := &p2pmocks.P2PMessageMock{DataField: data}

	msgID, err := txRes.ProcessReceivedMessage(msg, connectedPeerId, &p2pmocks.MessengerStub{})

	assert.Nil(t, err)
	assert.True(t, searchWasCalled)
	assert.True(t, sendWasCalled)
	assert.True(t, arg.Throttler.(*mock.ThrottlerStub).StartWasCalled())
	assert.True(t, arg.Throttler.(*mock.ThrottlerStub).EndWasCalled())
	assert.Len(t, msgID, 0)
}

func TestTxResolver_ProcessReceivedMessageFoundInTxPoolMarshalizerFailShouldRetNilAndErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("MarshalizerMock generic error")

	marshalizerMock := &mock.MarshalizerMock{}
	marshalizerStub := &mock.MarshalizerStub{
		MarshalCalled: func(obj interface{}) (i []byte, e error) {
			return nil, errExpected
		},
		UnmarshalCalled: func(obj interface{}, buff []byte) error {
			return marshalizerMock.Unmarshal(obj, buff)
		},
	}
	txReturned := &transaction.Transaction{
		Nonce: 10,
	}
	txPool := testscommon.NewShardedDataStub()
	txPool.SearchFirstDataCalled = func(key []byte) (value interface{}, ok bool) {
		if bytes.Equal([]byte("aaa"), key) {
			return txReturned, true
		}

		return nil, false
	}

	arg := createMockArgTxResolver()
	arg.TxPool = txPool
	arg.Marshaller = marshalizerStub
	txRes, _ := resolvers.NewTxResolver(arg)

	data, _ := marshalizerMock.Marshal(&dataRetriever.RequestData{Type: dataRetriever.HashType, Value: []byte("aaa")})

	msg := &p2pmocks.P2PMessageMock{DataField: data}

	msgID, err := txRes.ProcessReceivedMessage(msg, connectedPeerId, &p2pmocks.MessengerStub{})

	assert.True(t, errors.Is(err, errExpected))
	assert.True(t, arg.Throttler.(*mock.ThrottlerStub).StartWasCalled())
	assert.True(t, arg.Throttler.(*mock.ThrottlerStub).EndWasCalled())
	assert.Nil(t, msgID)
}

func TestTxResolver_ProcessReceivedMessageBatchMarshalFailShouldRetNilAndErr(t *testing.T) {
	t.Parallel()

	marshalizerMock := &mock.MarshalizerMock{}
	cnt := 0
	marshalizerStub := &mock.MarshalizerStub{
		MarshalCalled: func(obj interface{}) (i []byte, e error) {
			cnt++
			if cnt > 1 {
				return nil, expectedErr
			}
			return marshalizerMock.Marshal(obj)
		},
		UnmarshalCalled: func(obj interface{}, buff []byte) error {
			return marshalizerMock.Unmarshal(obj, buff)
		},
	}
	txReturned := &transaction.Transaction{
		Nonce: 10,
	}
	txPool := testscommon.NewShardedDataStub()
	txPool.SearchFirstDataCalled = func(key []byte) (value interface{}, ok bool) {
		if bytes.Equal([]byte("aaa"), key) {
			return txReturned, true
		}

		return nil, false
	}

	arg := createMockArgTxResolver()
	arg.TxPool = txPool
	arg.Marshaller = marshalizerStub
	txRes, _ := resolvers.NewTxResolver(arg)

	data, _ := marshalizerMock.Marshal(&dataRetriever.RequestData{Type: dataRetriever.HashType, Value: []byte("aaa")})

	msg := &p2pmocks.P2PMessageMock{DataField: data}

	msgID, err := txRes.ProcessReceivedMessage(msg, connectedPeerId, &p2pmocks.MessengerStub{})

	assert.True(t, errors.Is(err, expectedErr))
	assert.True(t, arg.Throttler.(*mock.ThrottlerStub).StartWasCalled())
	assert.True(t, arg.Throttler.(*mock.ThrottlerStub).EndWasCalled())
	assert.Nil(t, msgID)
}

func TestTxResolver_ProcessReceivedMessageFoundInTxStorageShouldRetValAndSend(t *testing.T) {
	t.Parallel()

	marshalizer := &mock.MarshalizerMock{}

	txPool := testscommon.NewShardedDataStub()
	txPool.SearchFirstDataCalled = func(key []byte) (value interface{}, ok bool) {
		// not found in txPool
		return nil, false
	}
	searchWasCalled := false
	sendWasCalled := false
	txReturned := &transaction.Transaction{
		Nonce: 10,
	}
	txReturnedAsBuffer, _ := marshalizer.Marshal(txReturned)
	txStorage := &storageStubs.StorerStub{}
	txStorage.SearchFirstCalled = func(key []byte) (i []byte, e error) {
		if bytes.Equal([]byte("aaa"), key) {
			searchWasCalled = true
			return txReturnedAsBuffer, nil
		}

		return nil, nil
	}

	arg := createMockArgTxResolver()
	arg.SenderResolver = &mock.TopicResolverSenderStub{
		SendCalled: func(buff []byte, peer core.PeerID, source p2p.MessageHandler) error {
			sendWasCalled = true
			return nil
		},
	}
	arg.TxPool = txPool
	arg.TxStorage = txStorage
	txRes, _ := resolvers.NewTxResolver(arg)

	data, _ := marshalizer.Marshal(&dataRetriever.RequestData{Type: dataRetriever.HashType, Value: []byte("aaa")})

	msg := &p2pmocks.P2PMessageMock{DataField: data}

	msgID, err := txRes.ProcessReceivedMessage(msg, connectedPeerId, &p2pmocks.MessengerStub{})

	assert.Nil(t, err)
	assert.True(t, searchWasCalled)
	assert.True(t, sendWasCalled)
	assert.True(t, arg.Throttler.(*mock.ThrottlerStub).StartWasCalled())
	assert.True(t, arg.Throttler.(*mock.ThrottlerStub).EndWasCalled())
	assert.Len(t, msgID, 0)
}

func TestTxResolver_ProcessReceivedMessageFoundInTxStorageCheckRetError(t *testing.T) {
	t.Parallel()

	marshalizer := &mock.MarshalizerMock{}

	txPool := testscommon.NewShardedDataStub()
	txPool.SearchFirstDataCalled = func(key []byte) (value interface{}, ok bool) {
		// not found in txPool
		return nil, false
	}

	errExpected := errors.New("expected error")

	txStorage := &storageStubs.StorerStub{}
	txStorage.SearchFirstCalled = func(key []byte) (i []byte, e error) {
		if bytes.Equal([]byte("aaa"), key) {
			return nil, errExpected
		}

		return nil, nil
	}

	arg := createMockArgTxResolver()
	arg.TxPool = txPool
	arg.TxStorage = txStorage
	txRes, _ := resolvers.NewTxResolver(arg)

	data, _ := marshalizer.Marshal(&dataRetriever.RequestData{Type: dataRetriever.HashType, Value: []byte("aaa")})

	msg := &p2pmocks.P2PMessageMock{DataField: data}

	msgID, err := txRes.ProcessReceivedMessage(msg, connectedPeerId, &p2pmocks.MessengerStub{})

	assert.True(t, errors.Is(err, errExpected))
	assert.True(t, arg.Throttler.(*mock.ThrottlerStub).StartWasCalled())
	assert.True(t, arg.Throttler.(*mock.ThrottlerStub).EndWasCalled())
	assert.Nil(t, msgID)
}

func TestTxResolver_ProcessReceivedMessageRequestedTwoSmallTransactionsShouldCallSliceSplitter(t *testing.T) {
	t.Parallel()

	txHash1 := []byte("txHash1")
	txHash2 := []byte("txHash2")

	tx1 := &transaction.Transaction{
		Nonce: 10,
	}
	tx2 := &transaction.Transaction{
		Nonce: 20,
	}

	marshalizer := &mock.MarshalizerMock{}
	txPool := testscommon.NewShardedDataStub()
	txPool.SearchFirstDataCalled = func(key []byte) (value interface{}, ok bool) {
		if bytes.Equal(txHash1, key) {
			return tx1, true
		}
		if bytes.Equal(txHash2, key) {
			return tx2, true
		}

		return nil, false
	}

	splitSliceWasCalled := false
	sendWasCalled := false
	arg := createMockArgTxResolver()
	arg.SenderResolver = &mock.TopicResolverSenderStub{
		SendCalled: func(buff []byte, peer core.PeerID, source p2p.MessageHandler) error {
			sendWasCalled = true
			return nil
		},
	}
	arg.TxPool = txPool
	arg.DataPacker = &mock.DataPackerStub{
		PackDataInChunksCalled: func(data [][]byte, limit int) ([][]byte, error) {
			if len(data) != 2 {
				return nil, errors.New("should have been 2 data pieces")
			}

			splitSliceWasCalled = true
			return make([][]byte, 1), nil
		},
	}
	txRes, _ := resolvers.NewTxResolver(arg)

	buff, _ := marshalizer.Marshal(&batch.Batch{Data: [][]byte{txHash1, txHash2}})
	data, _ := marshalizer.Marshal(&dataRetriever.RequestData{Type: dataRetriever.HashArrayType, Value: buff})

	msg := &p2pmocks.P2PMessageMock{DataField: data}

	msgID, err := txRes.ProcessReceivedMessage(msg, connectedPeerId, &p2pmocks.MessengerStub{})

	assert.Nil(t, err)
	assert.True(t, splitSliceWasCalled)
	assert.True(t, sendWasCalled)
	assert.True(t, arg.Throttler.(*mock.ThrottlerStub).StartWasCalled())
	assert.True(t, arg.Throttler.(*mock.ThrottlerStub).EndWasCalled())
	assert.Len(t, msgID, 0)
}

func TestTxResolver_ProcessReceivedMessageRequestedTwoSmallTransactionsFoundOnlyOneShouldWork(t *testing.T) {
	t.Parallel()

	txHash1 := []byte("txHash1")
	txHash2 := []byte("txHash2")

	tx1 := &transaction.Transaction{
		Nonce: 10,
	}

	marshalizer := &mock.MarshalizerMock{}
	txPool := testscommon.NewShardedDataStub()
	txPool.SearchFirstDataCalled = func(key []byte) (value interface{}, ok bool) {
		if bytes.Equal(txHash1, key) {
			return tx1, true
		}

		return nil, false
	}

	splitSliceWasCalled := false
	sendWasCalled := false
	arg := createMockArgTxResolver()
	arg.SenderResolver = &mock.TopicResolverSenderStub{
		SendCalled: func(buff []byte, peer core.PeerID, source p2p.MessageHandler) error {
			sendWasCalled = true
			return nil
		},
	}
	arg.TxStorage = &storageStubs.StorerStub{
		SearchFirstCalled: func(key []byte) (i []byte, err error) {
			return nil, errors.New("not found")
		},
	}
	arg.TxPool = txPool
	arg.DataPacker = &mock.DataPackerStub{
		PackDataInChunksCalled: func(data [][]byte, limit int) ([][]byte, error) {
			if len(data) != 1 {
				return nil, errors.New("should have been 1 data piece")
			}

			splitSliceWasCalled = true
			return make([][]byte, 1), nil
		},
	}
	txRes, _ := resolvers.NewTxResolver(arg)

	buff, _ := marshalizer.Marshal(&batch.Batch{Data: [][]byte{txHash1, txHash2}})
	data, _ := marshalizer.Marshal(&dataRetriever.RequestData{Type: dataRetriever.HashArrayType, Value: buff})

	msg := &p2pmocks.P2PMessageMock{DataField: data}

	msgID, err := txRes.ProcessReceivedMessage(msg, connectedPeerId, &p2pmocks.MessengerStub{})

	assert.NotNil(t, err)
	assert.True(t, splitSliceWasCalled)
	assert.True(t, sendWasCalled)
	assert.True(t, arg.Throttler.(*mock.ThrottlerStub).StartWasCalled())
	assert.True(t, arg.Throttler.(*mock.ThrottlerStub).EndWasCalled())
	assert.Nil(t, msgID)
}

func TestTxResolver_ProcessReceivedMessageHashArrayUnmarshalFails(t *testing.T) {
	t.Parallel()

	arg := createMockArgTxResolver()
	marshalizer := arg.Marshaller
	cnt := 0
	arg.Marshaller = &mock.MarshalizerStub{
		UnmarshalCalled: func(obj interface{}, buff []byte) error {
			cnt++
			if cnt > 1 {
				return expectedErr
			}
			return marshalizer.Unmarshal(obj, buff)
		},
	}
	txRes, _ := resolvers.NewTxResolver(arg)

	data, _ := marshalizer.Marshal(&dataRetriever.RequestData{Type: dataRetriever.HashArrayType, Value: []byte("buff")})
	msg := &p2pmocks.P2PMessageMock{DataField: data}

	msgID, err := txRes.ProcessReceivedMessage(msg, connectedPeerId, &p2pmocks.MessengerStub{})

	assert.True(t, errors.Is(err, expectedErr))
	assert.True(t, arg.Throttler.(*mock.ThrottlerStub).StartWasCalled())
	assert.True(t, arg.Throttler.(*mock.ThrottlerStub).EndWasCalled())
	assert.Nil(t, msgID)
}

func TestTxResolver_ProcessReceivedMessageHashArrayPackDataInChunksFails(t *testing.T) {
	t.Parallel()

	txHash1 := []byte("txHash1")
	txHash2 := []byte("txHash2")

	arg := createMockArgTxResolver()
	arg.DataPacker = &mock.DataPackerStub{
		PackDataInChunksCalled: func(data [][]byte, limit int) ([][]byte, error) {
			return nil, expectedErr
		},
	}
	txRes, _ := resolvers.NewTxResolver(arg)

	buff, _ := arg.Marshaller.Marshal(&batch.Batch{Data: [][]byte{txHash1, txHash2}})
	data, _ := arg.Marshaller.Marshal(&dataRetriever.RequestData{Type: dataRetriever.HashArrayType, Value: buff})
	msg := &p2pmocks.P2PMessageMock{DataField: data}

	msgID, err := txRes.ProcessReceivedMessage(msg, connectedPeerId, &p2pmocks.MessengerStub{})

	assert.True(t, errors.Is(err, expectedErr))
	assert.True(t, arg.Throttler.(*mock.ThrottlerStub).StartWasCalled())
	assert.True(t, arg.Throttler.(*mock.ThrottlerStub).EndWasCalled())
	assert.Nil(t, msgID)
}

func TestTxResolver_ProcessReceivedMessageHashArraySendFails(t *testing.T) {
	t.Parallel()

	txHash1 := []byte("txHash1")
	txHash2 := []byte("txHash2")

	arg := createMockArgTxResolver()
	arg.SenderResolver = &mock.TopicResolverSenderStub{
		SendCalled: func(buff []byte, peer core.PeerID, source p2p.MessageHandler) error {
			return expectedErr
		},
	}
	txRes, _ := resolvers.NewTxResolver(arg)

	buff, _ := arg.Marshaller.Marshal(&batch.Batch{Data: [][]byte{txHash1, txHash2}})
	data, _ := arg.Marshaller.Marshal(&dataRetriever.RequestData{Type: dataRetriever.HashArrayType, Value: buff})
	msg := &p2pmocks.P2PMessageMock{DataField: data}

	msgID, err := txRes.ProcessReceivedMessage(msg, connectedPeerId, &p2pmocks.MessengerStub{})

	assert.True(t, errors.Is(err, expectedErr))
	assert.True(t, arg.Throttler.(*mock.ThrottlerStub).StartWasCalled())
	assert.True(t, arg.Throttler.(*mock.ThrottlerStub).EndWasCalled())
	assert.Nil(t, msgID)
}

func TestTxResolver_Close(t *testing.T) {
	t.Parallel()

	arg := createMockArgTxResolver()
	txRes, _ := resolvers.NewTxResolver(arg)

	assert.Nil(t, txRes.Close())
}
