package sync

import (
	"context"
	"fmt"
	"math"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/storage"
)

// ShardBootstrap implements the bootstrap mechanism
type ShardBootstrap struct {
	*baseBootstrap
}

// NewShardBootstrap creates a new Bootstrap object
func NewShardBootstrap(arguments ArgShardBootstrapper) (*ShardBootstrap, error) {
	if check.IfNil(arguments.PoolsHolder) {
		return nil, process.ErrNilPoolsHolder
	}
	if check.IfNil(arguments.PoolsHolder.Headers()) {
		return nil, process.ErrNilHeadersDataPool
	}
	if check.IfNil(arguments.PoolsHolder.Proofs()) {
		return nil, process.ErrNilProofsPool
	}
	if check.IfNil(arguments.PoolsHolder.MiniBlocks()) {
		return nil, process.ErrNilTxBlockBody
	}

	err := checkBaseBootstrapParameters(arguments.ArgBaseBootstrapper)
	if err != nil {
		return nil, err
	}

	base := &baseBootstrap{
		chainHandler:                 arguments.ChainHandler,
		blockProcessor:               arguments.BlockProcessor,
		store:                        arguments.Store,
		headers:                      arguments.PoolsHolder.Headers(),
		proofs:                       arguments.PoolsHolder.Proofs(),
		roundHandler:                 arguments.RoundHandler,
		waitTime:                     arguments.WaitTime,
		hasher:                       arguments.Hasher,
		marshalizer:                  arguments.Marshalizer,
		forkDetector:                 arguments.ForkDetector,
		requestHandler:               arguments.RequestHandler,
		shardCoordinator:             arguments.ShardCoordinator,
		accounts:                     arguments.Accounts,
		blackListHandler:             arguments.BlackListHandler,
		networkWatcher:               arguments.NetworkWatcher,
		bootStorer:                   arguments.BootStorer,
		storageBootstrapper:          arguments.StorageBootstrapper,
		epochHandler:                 arguments.EpochHandler,
		miniBlocksProvider:           arguments.MiniblocksProvider,
		uint64Converter:              arguments.Uint64Converter,
		poolsHolder:                  arguments.PoolsHolder,
		statusHandler:                arguments.AppStatusHandler,
		outportHandler:               arguments.OutportHandler,
		accountsDBSyncer:             arguments.AccountsDBSyncer,
		currentEpochProvider:         arguments.CurrentEpochProvider,
		isInImportMode:               arguments.IsInImportMode,
		historyRepo:                  arguments.HistoryRepo,
		scheduledTxsExecutionHandler: arguments.ScheduledTxsExecutionHandler,
		processWaitTime:              arguments.ProcessWaitTime,
		repopulateTokensSupplies:     arguments.RepopulateTokensSupplies,
		enableEpochsHandler:          arguments.EnableEpochsHandler,
	}

	if base.isInImportMode {
		log.Warn("using always-not-synced status because the node is running in import-db")
	}

	boot := ShardBootstrap{
		baseBootstrap: base,
	}

	base.blockBootstrapper = &boot
	base.syncStarter = &boot
	base.requestMiniBlocks = boot.requestMiniBlocksFromHeaderWithNonceIfMissing

	// placed in struct fields for performance reasons
	base.headerStore, err = boot.store.GetStorer(dataRetriever.BlockHeaderUnit)
	if err != nil {
		return nil, err
	}

	hdrNonceHashDataUnit := dataRetriever.GetHdrNonceHashDataUnit(boot.shardCoordinator.SelfId())
	base.headerNonceHashStore, err = boot.store.GetStorer(hdrNonceHashDataUnit)
	if err != nil {
		return nil, err
	}

	base.init()

	return &boot, nil
}

func (boot *ShardBootstrap) getBlockBody(headerHandler data.HeaderHandler) (data.BodyHandler, error) {
	header, ok := headerHandler.(data.ShardHeaderHandler)
	if !ok {
		return nil, process.ErrWrongTypeAssertion
	}

	hashes := make([][]byte, len(header.GetMiniBlockHeaderHandlers()))
	for i := 0; i < len(header.GetMiniBlockHeaderHandlers()); i++ {
		hashes[i] = header.GetMiniBlockHeaderHandlers()[i].GetHash()
	}

	miniBlocksAndHashes, missingMiniBlocksHashes := boot.miniBlocksProvider.GetMiniBlocks(hashes)
	if len(missingMiniBlocksHashes) > 0 {
		return nil, process.ErrMissingBody
	}

	miniBlocks := make([]*block.MiniBlock, len(miniBlocksAndHashes))
	for index, miniBlockAndHash := range miniBlocksAndHashes {
		miniBlocks[index] = miniBlockAndHash.Miniblock
	}

	return &block.Body{MiniBlocks: miniBlocks}, nil
}

// StartSyncingBlocks method will start syncing blocks as a go routine
func (boot *ShardBootstrap) StartSyncingBlocks() error {
	errNotCritical := boot.storageBootstrapper.LoadFromStorage()
	if errNotCritical != nil {
		log.Debug("boot.syncFromStorer",
			"error", errNotCritical.Error(),
		)
	}

	var ctx context.Context
	ctx, boot.cancelFunc = context.WithCancel(context.Background())

	err := boot.handleAccountsTrieIteration()
	if err != nil {
		return fmt.Errorf("%w while handling accounts trie iteration", err)
	}

	go boot.syncBlocks(ctx)

	return nil
}

// SyncBlock method actually does the synchronization. It requests the next block header from the pool
// and if it is not found there it will be requested from the network. After the header is received,
// it requests the block body in the same way(pool and then, if it is not found in the pool, from network).
// If either header and body are received the ProcessBlock and CommitBlock method will be called successively.
// These methods will execute the block and its transactions. Finally, if everything works, the block will be committed
// in the blockchain, and all this mechanism will be reiterated for the next block.
func (boot *ShardBootstrap) SyncBlock(ctx context.Context) error {
	err := boot.syncBlock()
	if core.IsGetNodeFromDBError(err) {
		getNodeErr := core.UnwrapGetNodeFromDBErr(err)
		if getNodeErr == nil {
			return err
		}

		errSync := boot.syncUserAccountsState(getNodeErr.GetKey())
		boot.handleTrieSyncError(errSync, ctx)
	}

	return err
}

// Close closes the synchronization loop
func (boot *ShardBootstrap) Close() error {
	if check.IfNil(boot.baseBootstrap) {
		return nil
	}

	return boot.baseBootstrap.Close()
}

func (boot *ShardBootstrap) getPrevHeader(
	header data.HeaderHandler,
	headerStore storage.Storer,
) (data.HeaderHandler, error) {

	prevHash := header.GetPrevHash()
	buffHeader, err := headerStore.Get(prevHash)
	if err != nil {
		return nil, err
	}

	prevHeader, err := process.UnmarshalShardHeader(boot.marshalizer, buffHeader)
	if err != nil {
		return nil, err
	}

	return prevHeader, nil
}

func (boot *ShardBootstrap) getCurrHeader() (data.HeaderHandler, error) {
	blockHeader := boot.chainHandler.GetCurrentBlockHeader()
	if check.IfNil(blockHeader) {
		return nil, process.ErrNilBlockHeader
	}

	header, ok := blockHeader.(data.ShardHeaderHandler)
	if !ok {
		return nil, process.ErrWrongTypeAssertion
	}

	return header, nil
}

func (boot *ShardBootstrap) requestMiniBlocksFromHeaderWithNonceIfMissing(headerHandler data.HeaderHandler) {
	nextBlockNonce := boot.getNonceForNextBlock()
	maxNonce := core.MinUint64(nextBlockNonce+process.MaxHeadersToRequestInAdvance-1, boot.forkDetector.ProbableHighestNonce())
	if headerHandler.GetNonce() < nextBlockNonce || headerHandler.GetNonce() > maxNonce {
		return
	}

	header, ok := headerHandler.(data.ShardHeaderHandler)
	if !ok {
		log.Warn("cannot convert headerHandler in block.Header")
		return
	}

	hashes := make([][]byte, 0, len(header.GetMiniBlockHeaderHandlers()))
	for i := 0; i < len(header.GetMiniBlockHeaderHandlers()); i++ {
		hashes = append(hashes, header.GetMiniBlockHeaderHandlers()[i].GetHash())
	}

	_, missingMiniBlocksHashes := boot.miniBlocksProvider.GetMiniBlocksFromPool(hashes)
	if len(missingMiniBlocksHashes) > 0 {
		log.Trace("requesting in advance mini blocks",
			"num miniblocks", len(missingMiniBlocksHashes),
			"header nonce", header.GetNonce(),
		)
		boot.requestHandler.RequestMiniBlocks(boot.shardCoordinator.SelfId(), missingMiniBlocksHashes)
	}
}

func (boot *ShardBootstrap) getBlockBodyRequestingIfMissing(headerHandler data.HeaderHandler) (data.BodyHandler, error) {
	header, ok := headerHandler.(data.ShardHeaderHandler)
	if !ok {
		return nil, process.ErrWrongTypeAssertion
	}

	hashes := make([][]byte, len(header.GetMiniBlockHeaderHandlers()))
	for i := 0; i < len(header.GetMiniBlockHeaderHandlers()); i++ {
		hashes[i] = header.GetMiniBlockHeaderHandlers()[i].GetHash()
	}

	boot.setRequestedMiniBlocks(nil)

	miniBlockSlice, err := boot.getMiniBlocksRequestingIfMissing(hashes)
	if err != nil {
		return nil, err
	}

	blockBody := &block.Body{MiniBlocks: miniBlockSlice}

	return blockBody, nil
}

func (boot *ShardBootstrap) isForkTriggeredByMeta() bool {
	return boot.forkInfo.IsDetected &&
		boot.forkInfo.Nonce != math.MaxUint64 &&
		boot.forkInfo.Round == process.MinForkRound &&
		boot.forkInfo.Hash != nil
}

func (boot *ShardBootstrap) requestHeaderByNonce(nonce uint64) {
	boot.requestHandler.RequestShardHeaderByNonce(boot.shardCoordinator.SelfId(), nonce)
}

func (boot *ShardBootstrap) requestProofByNonce(nonce uint64) {
	boot.requestHandler.RequestEquivalentProofByNonce(boot.shardCoordinator.SelfId(), nonce)
}

// IsInterfaceNil returns true if there is no value under the interface
func (boot *ShardBootstrap) IsInterfaceNil() bool {
	return boot == nil
}
