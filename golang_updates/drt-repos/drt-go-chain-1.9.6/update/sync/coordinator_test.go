package sync

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	dataTransaction "github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/state"
	"github.com/TerraDharitri/drt-go-chain/storage"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/cache"
	stateMock "github.com/TerraDharitri/drt-go-chain/testscommon/state"
	storageStubs "github.com/TerraDharitri/drt-go-chain/testscommon/storage"
	"github.com/TerraDharitri/drt-go-chain/testscommon/syncer"
	trieMock "github.com/TerraDharitri/drt-go-chain/testscommon/trie"
	"github.com/TerraDharitri/drt-go-chain/update"
	"github.com/TerraDharitri/drt-go-chain/update/mock"
)

func createHeaderSyncHandler(retErr bool) update.HeaderSyncHandler {
	meta := &block.MetaBlock{
		Nonce: 1, Epoch: 1, RootHash: []byte("metaRootHash"),
		EpochStart: block.EpochStart{
			LastFinalizedHeaders: []block.EpochStartShardData{
				{
					ShardID:                 0,
					RootHash:                []byte("shardDataRootHash"),
					PendingMiniBlockHeaders: []block.MiniBlockHeader{{Hash: []byte("hash")}},
					FirstPendingMetaBlock:   []byte("firstPending"),
				},
			},
		},
	}
	args := createMockHeadersSyncHandlerArgs()
	args.StorageService = &storageStubs.ChainStorerStub{GetStorerCalled: func(unitType dataRetriever.UnitType) (storage.Storer, error) {
		return &storageStubs.StorerStub{
			GetCalled: func(key []byte) (bytes []byte, err error) {
				if retErr {
					return nil, errors.New("err")
				}

				return json.Marshal(meta)
			},
		}, nil
	}}

	if !retErr {
		args.StorageService = initStore()
		byteArray := args.Uint64Converter.ToByteSlice(meta.Nonce)
		_ = args.StorageService.Put(dataRetriever.MetaHdrNonceHashDataUnit, byteArray, []byte("firstPending"))
		marshaledData, _ := json.Marshal(meta)
		_ = args.StorageService.Put(dataRetriever.MetaBlockUnit, []byte("firstPending"), marshaledData)

		_ = args.StorageService.Put(dataRetriever.MetaBlockUnit, []byte(core.EpochStartIdentifier(meta.Epoch)), marshaledData)
	}

	headersSyncHandler, _ := NewHeadersSyncHandler(args)
	return headersSyncHandler
}

func createPendingMiniBlocksSyncHandler() update.EpochStartPendingMiniBlocksSyncHandler {
	txHash := []byte("txHash")
	mb := &block.MiniBlock{TxHashes: [][]byte{txHash}}
	args := ArgsNewPendingMiniBlocksSyncer{
		Storage: &storageStubs.StorerStub{},
		Cache: &cache.CacherStub{
			RegisterHandlerCalled: func(f func(key []byte, val interface{})) {},
			PeekCalled: func(key []byte) (value interface{}, ok bool) {
				return mb, true
			},
		},
		Marshalizer:    &mock.MarshalizerFake{},
		RequestHandler: &testscommon.RequestHandlerStub{},
	}

	pendingMiniBlocksSyncer, _ := NewPendingMiniBlocksSyncer(args)
	return pendingMiniBlocksSyncer
}

func createPendingTxSyncHandler() update.TransactionsSyncHandler {
	args := createMockArgs()
	args.Storages = &storageStubs.ChainStorerStub{
		GetStorerCalled: func(unitType dataRetriever.UnitType) (storage.Storer, error) {
			return &storageStubs.StorerStub{
				GetCalled: func(key []byte) (bytes []byte, err error) {
					tx := &dataTransaction.Transaction{
						Nonce: 1, Value: big.NewInt(10), SndAddr: []byte("snd"), RcvAddr: []byte("rcv"),
					}
					return json.Marshal(tx)
				},
			}, nil
		},
	}

	pendingTxsSyncer, _ := NewTransactionsSyncer(args)
	return pendingTxsSyncer
}

func createSyncTrieState(retErr bool) update.EpochStartTriesSyncHandler {
	args := ArgsNewSyncAccountsDBsHandler{
		AccountsDBsSyncers: &mock.AccountsDBSyncersStub{
			GetCalled: func(key string) (syncer update.AccountsDBSyncer, err error) {
				return &mock.AccountsDBSyncerStub{
					SyncAccountsCalled: func(rootHash []byte, _ common.StorageMarker) error {
						if retErr {
							return errors.New("err")
						}
						return nil
					},
				}, nil
			},
		},
		ActiveAccountsDBs: make(map[state.AccountsDbIdentifier]state.AccountsAdapter),
	}

	args.ActiveAccountsDBs[state.UserAccountsState] = &stateMock.AccountsStub{
		RecreateAllTriesCalled: func(rootHash []byte) (map[string]common.Trie, error) {
			tries := make(map[string]common.Trie)
			tries[string(rootHash)] = &trieMock.TrieStub{
				CommitCalled: func() error {
					if retErr {
						return errors.New("err")
					}
					return nil
				},
			}
			return tries, nil
		},
	}

	args.ActiveAccountsDBs[state.PeerAccountsState] = &stateMock.AccountsStub{
		RecreateAllTriesCalled: func(rootHash []byte) (map[string]common.Trie, error) {
			tries := make(map[string]common.Trie)
			tries[string(rootHash)] = &trieMock.TrieStub{
				CommitCalled: func() error {
					if retErr {
						return errors.New("err")
					}
					return nil
				},
			}
			return tries, nil
		},
	}

	triesSyncHandler, _ := NewSyncAccountsDBsHandler(args)
	return triesSyncHandler
}

func TestNewSyncState_Ok(t *testing.T) {
	t.Parallel()

	args := ArgsNewSyncState{
		Headers:      createHeaderSyncHandler(false),
		Tries:        createSyncTrieState(false),
		MiniBlocks:   createPendingMiniBlocksSyncHandler(),
		Transactions: createPendingTxSyncHandler(),
	}

	ss, err := NewSyncState(args)
	require.Nil(t, err)
	require.False(t, ss.IsInterfaceNil())

	err = ss.SyncAllState(1)
	require.Nil(t, err)
}

func TestNewSyncState_CannotSyncHeaderErr(t *testing.T) {
	t.Parallel()

	args := ArgsNewSyncState{
		Headers:      createHeaderSyncHandler(true),
		Tries:        createSyncTrieState(false),
		MiniBlocks:   createPendingMiniBlocksSyncHandler(),
		Transactions: createPendingTxSyncHandler(),
	}

	ss, err := NewSyncState(args)
	require.Nil(t, err)

	err = ss.SyncAllState(1)
	require.NotNil(t, err)
}

func TestNewSyncState_CannotSyncTriesErr(t *testing.T) {
	t.Parallel()

	args := ArgsNewSyncState{
		Headers:      createHeaderSyncHandler(false),
		Tries:        createSyncTrieState(true),
		MiniBlocks:   createPendingMiniBlocksSyncHandler(),
		Transactions: createPendingTxSyncHandler(),
	}

	ss, err := NewSyncState(args)
	require.Nil(t, err)

	err = ss.SyncAllState(1)
	require.NotNil(t, err)
}

func TestSyncState_SyncAllStatePendingMiniBlocksErr(t *testing.T) {
	t.Parallel()

	localErr := errors.New("err")
	args := ArgsNewSyncState{
		Headers: &mock.HeaderSyncHandlerStub{
			SyncUnFinishedMetaHeadersCalled: func(epoch uint32) error {
				return nil
			},
			GetEpochStartMetaBlockCalled: func() (metaBlock data.MetaHeaderHandler, err error) {
				return &block.MetaBlock{}, nil
			},
		},
		Tries: &mock.EpochStartTriesSyncHandlerMock{},
		MiniBlocks: &mock.EpochStartPendingMiniBlocksSyncHandlerMock{
			SyncPendingMiniBlocksFromMetaCalled: func(meta data.MetaHeaderHandler, unFinished map[string]data.MetaHeaderHandler, ctx context.Context) error {
				return localErr
			},
		},
		Transactions: &syncer.TransactionsSyncHandlerMock{},
	}

	ss, err := NewSyncState(args)
	require.Nil(t, err)

	err = ss.SyncAllState(0)
	require.True(t, errors.Is(err, localErr))
}

func TestSyncState_SyncAllStateGetMiniBlocksErr(t *testing.T) {
	t.Parallel()

	localErr := errors.New("err")
	args := ArgsNewSyncState{
		Headers: &mock.HeaderSyncHandlerStub{
			SyncUnFinishedMetaHeadersCalled: func(epoch uint32) error {
				return nil
			},
			GetEpochStartMetaBlockCalled: func() (metaBlock data.MetaHeaderHandler, err error) {
				return &block.MetaBlock{}, nil
			},
		},
		Tries: &mock.EpochStartTriesSyncHandlerMock{},
		MiniBlocks: &mock.EpochStartPendingMiniBlocksSyncHandlerMock{
			GetMiniBlocksCalled: func() (m map[string]*block.MiniBlock, err error) {
				return nil, localErr
			},
		},
		Transactions: &syncer.TransactionsSyncHandlerMock{},
	}

	ss, err := NewSyncState(args)
	require.Nil(t, err)

	err = ss.SyncAllState(0)
	require.True(t, errors.Is(err, localErr))
}

func TestSyncState_SyncAllStateSyncTxsErr(t *testing.T) {
	t.Parallel()

	localErr := errors.New("err")
	args := ArgsNewSyncState{
		Headers: &mock.HeaderSyncHandlerStub{
			SyncUnFinishedMetaHeadersCalled: func(epoch uint32) error {
				return nil
			},
			GetEpochStartMetaBlockCalled: func() (metaBlock data.MetaHeaderHandler, err error) {
				return &block.MetaBlock{}, nil
			},
		},
		Tries:      &mock.EpochStartTriesSyncHandlerMock{},
		MiniBlocks: &mock.EpochStartPendingMiniBlocksSyncHandlerMock{},
		Transactions: &syncer.TransactionsSyncHandlerMock{
			SyncTransactionsForCalled: func(miniBlocks map[string]*block.MiniBlock, epoch uint32, ctx context.Context) error {
				return localErr
			},
		},
	}

	ss, err := NewSyncState(args)
	require.Nil(t, err)

	err = ss.SyncAllState(0)
	require.True(t, errors.Is(err, localErr))
}
