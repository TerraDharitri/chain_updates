package dataRetriever

import (
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"

	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever/dataPool"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever/dataPool/headersCache"
	proofscache "github.com/TerraDharitri/drt-go-chain/dataRetriever/dataPool/proofsCache"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever/shardedData"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever/txpool"
	"github.com/TerraDharitri/drt-go-chain/storage"
	"github.com/TerraDharitri/drt-go-chain/storage/cache"
	"github.com/TerraDharitri/drt-go-chain/storage/storageunit"
	"github.com/TerraDharitri/drt-go-chain/testscommon/txcachemocks"
)

// PoolsHolderMock -
type PoolsHolderMock struct {
	transactions           dataRetriever.ShardedDataCacherNotifier
	unsignedTransactions   dataRetriever.ShardedDataCacherNotifier
	rewardTransactions     dataRetriever.ShardedDataCacherNotifier
	headers                dataRetriever.HeadersPool
	miniBlocks             storage.Cacher
	peerChangesBlocks      storage.Cacher
	trieNodes              storage.Cacher
	trieNodesChunks        storage.Cacher
	smartContracts         storage.Cacher
	currBlockTxs           dataRetriever.TransactionCacher
	currEpochValidatorInfo dataRetriever.ValidatorInfoCacher
	peerAuthentications    storage.Cacher
	heartbeats             storage.Cacher
	validatorsInfo         dataRetriever.ShardedDataCacherNotifier
	proofs                 dataRetriever.ProofsPool
}

// NewPoolsHolderMock -
func NewPoolsHolderMock() *PoolsHolderMock {
	var err error
	holder := &PoolsHolderMock{}

	holder.transactions, err = txpool.NewShardedTxPool(
		txpool.ArgShardedTxPool{
			Config: storageunit.CacheConfig{
				Capacity:             100000,
				SizePerSender:        1000,
				SizeInBytes:          1000000000,
				SizeInBytesPerSender: 10000000,
				Shards:               16,
			},
			TxGasHandler:   txcachemocks.NewTxGasHandlerMock(),
			Marshalizer:    &marshal.GogoProtoMarshalizer{},
			NumberOfShards: 1,
		},
	)
	panicIfError("NewPoolsHolderMock", err)

	holder.unsignedTransactions, err = shardedData.NewShardedData("unsignedTxPool", storageunit.CacheConfig{
		Capacity:    10000,
		SizeInBytes: 1000000000,
		Shards:      1,
	})
	panicIfError("NewPoolsHolderMock", err)

	holder.rewardTransactions, err = shardedData.NewShardedData("rewardsTxPool", storageunit.CacheConfig{
		Capacity:    100,
		SizeInBytes: 100000,
		Shards:      1,
	})
	panicIfError("NewPoolsHolderMock", err)

	holder.headers, err = headersCache.NewHeadersPool(config.HeadersPoolConfig{MaxHeadersPerShard: 1000, NumElementsToRemoveOnEviction: 100})
	panicIfError("NewPoolsHolderMock", err)

	holder.miniBlocks, err = storageunit.NewCache(storageunit.CacheConfig{Type: storageunit.LRUCache, Capacity: 10000, Shards: 1, SizeInBytes: 0})
	panicIfError("NewPoolsHolderMock", err)

	holder.peerChangesBlocks, err = storageunit.NewCache(storageunit.CacheConfig{Type: storageunit.LRUCache, Capacity: 10000, Shards: 1, SizeInBytes: 0})
	panicIfError("NewPoolsHolderMock", err)

	holder.currBlockTxs = dataPool.NewCurrentBlockTransactionsPool()
	holder.currEpochValidatorInfo = dataPool.NewCurrentEpochValidatorInfoPool()

	holder.trieNodes, err = storageunit.NewCache(storageunit.CacheConfig{Type: storageunit.SizeLRUCache, Capacity: 900000, Shards: 1, SizeInBytes: 314572800})
	panicIfError("NewPoolsHolderMock", err)

	holder.trieNodesChunks, err = storageunit.NewCache(storageunit.CacheConfig{Type: storageunit.SizeLRUCache, Capacity: 900000, Shards: 1, SizeInBytes: 314572800})
	panicIfError("NewPoolsHolderMock", err)

	holder.smartContracts, err = storageunit.NewCache(storageunit.CacheConfig{Type: storageunit.LRUCache, Capacity: 10000, Shards: 1, SizeInBytes: 0})
	panicIfError("NewPoolsHolderMock", err)

	holder.peerAuthentications, err = cache.NewTimeCacher(cache.ArgTimeCacher{
		DefaultSpan: 10 * time.Second,
		CacheExpiry: 10 * time.Second,
	})
	panicIfError("NewPoolsHolderMock", err)

	holder.heartbeats, err = storageunit.NewCache(storageunit.CacheConfig{Type: storageunit.LRUCache, Capacity: 10000, Shards: 1, SizeInBytes: 0})
	panicIfError("NewPoolsHolderMock", err)

	holder.validatorsInfo, err = shardedData.NewShardedData("validatorsInfoPool", storageunit.CacheConfig{
		Capacity:    100,
		SizeInBytes: 100000,
		Shards:      1,
	})
	panicIfError("NewPoolsHolderMock", err)

	holder.proofs = proofscache.NewProofsPool(3, 100)

	return holder
}

// CurrentBlockTxs -
func (holder *PoolsHolderMock) CurrentBlockTxs() dataRetriever.TransactionCacher {
	return holder.currBlockTxs
}

// CurrentEpochValidatorInfo -
func (holder *PoolsHolderMock) CurrentEpochValidatorInfo() dataRetriever.ValidatorInfoCacher {
	return holder.currEpochValidatorInfo
}

// Transactions -
func (holder *PoolsHolderMock) Transactions() dataRetriever.ShardedDataCacherNotifier {
	return holder.transactions
}

// UnsignedTransactions -
func (holder *PoolsHolderMock) UnsignedTransactions() dataRetriever.ShardedDataCacherNotifier {
	return holder.unsignedTransactions
}

// RewardTransactions -
func (holder *PoolsHolderMock) RewardTransactions() dataRetriever.ShardedDataCacherNotifier {
	return holder.rewardTransactions
}

// Headers -
func (holder *PoolsHolderMock) Headers() dataRetriever.HeadersPool {
	return holder.headers
}

// SetHeadersPool -
func (holder *PoolsHolderMock) SetHeadersPool(headersPool dataRetriever.HeadersPool) {
	holder.headers = headersPool
}

// MiniBlocks -
func (holder *PoolsHolderMock) MiniBlocks() storage.Cacher {
	return holder.miniBlocks
}

// PeerChangesBlocks -
func (holder *PoolsHolderMock) PeerChangesBlocks() storage.Cacher {
	return holder.peerChangesBlocks
}

// SetTransactions -
func (holder *PoolsHolderMock) SetTransactions(pool dataRetriever.ShardedDataCacherNotifier) {
	holder.transactions = pool
}

// SetUnsignedTransactions -
func (holder *PoolsHolderMock) SetUnsignedTransactions(pool dataRetriever.ShardedDataCacherNotifier) {
	holder.unsignedTransactions = pool
}

// TrieNodes -
func (holder *PoolsHolderMock) TrieNodes() storage.Cacher {
	return holder.trieNodes
}

// TrieNodesChunks -
func (holder *PoolsHolderMock) TrieNodesChunks() storage.Cacher {
	return holder.trieNodesChunks
}

// SmartContracts -
func (holder *PoolsHolderMock) SmartContracts() storage.Cacher {
	return holder.smartContracts
}

// PeerAuthentications -
func (holder *PoolsHolderMock) PeerAuthentications() storage.Cacher {
	return holder.peerAuthentications
}

// Heartbeats -
func (holder *PoolsHolderMock) Heartbeats() storage.Cacher {
	return holder.heartbeats
}

// ValidatorsInfo -
func (holder *PoolsHolderMock) ValidatorsInfo() dataRetriever.ShardedDataCacherNotifier {
	return holder.validatorsInfo
}

// Proofs -
func (holder *PoolsHolderMock) Proofs() dataRetriever.ProofsPool {
	return holder.proofs
}

// Close -
func (holder *PoolsHolderMock) Close() error {
	var lastError error
	if !check.IfNil(holder.trieNodes) {
		err := holder.trieNodes.Close()
		if err != nil {
			lastError = err
		}
	}

	if !check.IfNil(holder.peerAuthentications) {
		err := holder.peerAuthentications.Close()
		if err != nil {
			lastError = err
		}
	}

	return lastError
}

// IsInterfaceNil returns true if there is no value under the interface
func (holder *PoolsHolderMock) IsInterfaceNil() bool {
	return holder == nil
}
