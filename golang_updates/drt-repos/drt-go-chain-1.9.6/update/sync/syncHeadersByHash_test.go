package sync

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	dataRetrieverMocks "github.com/TerraDharitri/drt-go-chain/testscommon/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/testscommon/enableEpochsHandlerMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/genericMocks"
	storageStubs "github.com/TerraDharitri/drt-go-chain/testscommon/storage"
	"github.com/TerraDharitri/drt-go-chain/update"
	"github.com/TerraDharitri/drt-go-chain/update/mock"
	"github.com/stretchr/testify/require"
)

func TestNewMissingheadersByHashSyncer_NilParamsShouldErr(t *testing.T) {
	t.Parallel()

	okArgs := getMisingHeadersByHashSyncerArgs()

	testInput := make(map[ArgsNewMissingHeadersByHashSyncer]error)

	nilStorerArgs := okArgs
	nilStorerArgs.Storage = nil
	testInput[nilStorerArgs] = dataRetriever.ErrNilHeadersStorage

	nilCacheArgs := okArgs
	nilCacheArgs.Cache = nil
	testInput[nilCacheArgs] = update.ErrNilCacher

	nilMarshalizerArgs := okArgs
	nilMarshalizerArgs.Marshalizer = nil
	testInput[nilMarshalizerArgs] = dataRetriever.ErrNilMarshalizer

	nilRequestHandlerArgs := okArgs
	nilRequestHandlerArgs.RequestHandler = nil
	testInput[nilRequestHandlerArgs] = update.ErrNilRequestHandler

	nilProofsPoolArgs := okArgs
	nilProofsPoolArgs.ProofsPool = nil
	testInput[nilProofsPoolArgs] = dataRetriever.ErrNilProofsPool

	nilEnableEpochsHandlerArgs := okArgs
	nilEnableEpochsHandlerArgs.EnableEpochsHandler = nil
	testInput[nilEnableEpochsHandlerArgs] = process.ErrNilEnableEpochsHandler

	for args, expectedErr := range testInput {
		mhhs, err := NewMissingheadersByHashSyncer(args)
		require.True(t, check.IfNil(mhhs))
		require.Equal(t, expectedErr, err)
	}
}

func TestNewMissingheadersByHashSyncer_OkValsShouldWork(t *testing.T) {
	t.Parallel()

	args := getMisingHeadersByHashSyncerArgs()
	mhhs, err := NewMissingheadersByHashSyncer(args)
	require.NoError(t, err)
	require.NotNil(t, mhhs)
}

func TestSyncHeadersByHash_SyncMissingHeadersByHashHeaderFoundInCacheShouldWork(t *testing.T) {
	t.Parallel()

	args := getMisingHeadersByHashSyncerArgs()
	args.Cache = &mock.HeadersCacherStub{
		GetHeaderByHashCalled: func(_ []byte) (data.HeaderHandler, error) {
			return &block.MetaBlock{Nonce: 37}, nil
		},
	}
	mhhs, _ := NewMissingheadersByHashSyncer(args)

	err := mhhs.SyncMissingHeadersByHash([]uint32{0, 1}, [][]byte{[]byte("hash234")}, context.Background())
	require.NoError(t, err)
}

func TestSyncHeadersByHash_SyncMissingHeadersByHashHeaderFoundInStorageShouldWork(t *testing.T) {
	t.Parallel()

	args := getMisingHeadersByHashSyncerArgs()
	args.Cache = &mock.HeadersCacherStub{
		GetHeaderByHashCalled: func(_ []byte) (data.HeaderHandler, error) {
			return nil, errors.New("not found")
		},
	}
	args.Storage = &storageStubs.StorerStub{
		GetCalled: func(key []byte) ([]byte, error) {
			mb := &block.MetaBlock{Nonce: 37}
			mbBytes, _ := args.Marshalizer.Marshal(mb)
			return mbBytes, nil
		},
	}
	mhhs, _ := NewMissingheadersByHashSyncer(args)

	err := mhhs.SyncMissingHeadersByHash([]uint32{0, 1}, [][]byte{[]byte("hash234")}, context.Background())
	require.NoError(t, err)
}

func TestSyncHeadersByHash_SyncMissingHeadersByHashHeaderNotFoundShouldTimeout(t *testing.T) {
	t.Parallel()

	var errNotFound = errors.New("not found")
	args := getMisingHeadersByHashSyncerArgs()
	args.Cache = &mock.HeadersCacherStub{
		GetHeaderByHashCalled: func(_ []byte) (data.HeaderHandler, error) {
			return nil, errNotFound
		},
	}
	args.Storage = &storageStubs.StorerStub{
		GetCalled: func(_ []byte) ([]byte, error) {
			return nil, errNotFound
		},
	}
	mhhs, _ := NewMissingheadersByHashSyncer(args)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	err := mhhs.SyncMissingHeadersByHash([]uint32{0, 1}, [][]byte{[]byte("hash234")}, ctx)
	cancel()

	require.Equal(t, update.ErrTimeIsOut, err)
}

func TestSyncHeadersByHash_GetHeadersNotSyncedShouldErr(t *testing.T) {
	t.Parallel()

	args := getMisingHeadersByHashSyncerArgs()
	mhhs, _ := NewMissingheadersByHashSyncer(args)
	require.NotNil(t, mhhs)

	res, err := mhhs.GetHeaders()
	require.Nil(t, res)
	require.Equal(t, update.ErrNotSynced, err)
}

func TestSyncHeadersByHash_GetHeadersShouldReceiveAndReturnOkMb(t *testing.T) {
	t.Parallel()

	var handlerToNotify func(header data.HeaderHandler, shardHeaderHash []byte)
	var errNotFound = errors.New("not found")
	args := getMisingHeadersByHashSyncerArgs()
	args.Storage = &storageStubs.StorerStub{
		GetCalled: func(_ []byte) ([]byte, error) {
			return nil, errNotFound
		},
	}
	args.Cache = &mock.HeadersCacherStub{
		GetHeaderByHashCalled: func(_ []byte) (data.HeaderHandler, error) {
			return nil, errNotFound
		},
		RegisterHandlerCalled: func(handler func(header data.HeaderHandler, shardHeaderHash []byte)) {
			handlerToNotify = handler
		},
	}

	var wg sync.WaitGroup
	expectedHash := []byte("hash")
	expectedMB := &block.MetaBlock{Nonce: 37}
	args.RequestHandler = &testscommon.RequestHandlerStub{
		RequestShardHeaderCalled: func(shardID uint32, hash []byte) {
			wg.Add(1)
			go func() {
				handlerToNotify(expectedMB, expectedHash)
				wg.Done()
			}()
		},
	}

	mhhs, _ := NewMissingheadersByHashSyncer(args)
	require.NotNil(t, mhhs)

	err := mhhs.SyncMissingHeadersByHash([]uint32{0}, [][]byte{[]byte("hash")}, context.Background())
	require.NoError(t, err)
	wg.Wait()

	res, err := mhhs.GetHeaders()
	require.NoError(t, err)
	require.NotNil(t, res)

	actualMb, ok := res[string(expectedHash)]
	require.True(t, ok)
	require.Equal(t, expectedMB, actualMb)
}

func getMisingHeadersByHashSyncerArgs() ArgsNewMissingHeadersByHashSyncer {
	return ArgsNewMissingHeadersByHashSyncer{
		Storage:             genericMocks.NewStorerMock(),
		Cache:               &mock.HeadersCacherStub{},
		ProofsPool:          &dataRetrieverMocks.ProofsPoolMock{},
		Marshalizer:         &mock.MarshalizerMock{},
		RequestHandler:      &testscommon.RequestHandlerStub{},
		EnableEpochsHandler: &enableEpochsHandlerMock.EnableEpochsHandlerStub{},
	}
}
