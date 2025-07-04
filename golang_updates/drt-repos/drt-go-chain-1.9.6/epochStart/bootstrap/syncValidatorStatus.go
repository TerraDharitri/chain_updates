package bootstrap

import (
	"context"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain-core/data/endProcess"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/epochStart"
	"github.com/TerraDharitri/drt-go-chain/epochStart/bootstrap/disabled"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
	"github.com/TerraDharitri/drt-go-chain/storage"
	"github.com/TerraDharitri/drt-go-chain/storage/cache"
	"github.com/TerraDharitri/drt-go-chain/update"
	"github.com/TerraDharitri/drt-go-chain/update/sync"
)

const consensusGroupCacheSize = 50

type syncValidatorStatus struct {
	miniBlocksSyncer    epochStart.PendingMiniBlocksSyncHandler
	transactionsSyncer  update.TransactionsSyncHandler
	dataPool            dataRetriever.PoolsHolder
	marshalizer         marshal.Marshalizer
	requestHandler      process.RequestHandler
	nodeCoordinator     StartInEpochNodesCoordinator
	genesisNodesConfig  sharding.GenesisNodesSetupHandler
	memDB               storage.Storer
	enableEpochsHandler common.EnableEpochsHandler
}

// ArgsNewSyncValidatorStatus holds the arguments needed for creating a new validator status process component
type ArgsNewSyncValidatorStatus struct {
	DataPool                        dataRetriever.PoolsHolder
	Marshalizer                     marshal.Marshalizer
	Hasher                          hashing.Hasher
	RequestHandler                  process.RequestHandler
	ChanceComputer                  nodesCoordinator.ChanceComputer
	GenesisNodesConfig              sharding.GenesisNodesSetupHandler
	ChainParametersHandler          process.ChainParametersHandler
	NodeShuffler                    nodesCoordinator.NodesShuffler
	PubKey                          []byte
	ShardIdAsObserver               uint32
	ChanNodeStop                    chan endProcess.ArgEndProcess
	NodeTypeProvider                NodeTypeProviderHandler
	IsFullArchive                   bool
	EnableEpochsHandler             common.EnableEpochsHandler
	NodesCoordinatorRegistryFactory nodesCoordinator.NodesCoordinatorRegistryFactory
}

// NewSyncValidatorStatus creates a new validator status process component
func NewSyncValidatorStatus(args ArgsNewSyncValidatorStatus) (*syncValidatorStatus, error) {
	if args.ChanNodeStop == nil {
		return nil, nodesCoordinator.ErrNilNodeStopChannel
	}

	s := &syncValidatorStatus{
		dataPool:            args.DataPool,
		marshalizer:         args.Marshalizer,
		requestHandler:      args.RequestHandler,
		genesisNodesConfig:  args.GenesisNodesConfig,
		enableEpochsHandler: args.EnableEpochsHandler,
	}

	var err error

	syncMiniBlocksArgs := sync.ArgsNewPendingMiniBlocksSyncer{
		Storage:        disabled.CreateMemUnit(),
		Cache:          s.dataPool.MiniBlocks(),
		Marshalizer:    s.marshalizer,
		RequestHandler: s.requestHandler,
	}
	s.miniBlocksSyncer, err = sync.NewPendingMiniBlocksSyncer(syncMiniBlocksArgs)
	if err != nil {
		return nil, err
	}

	syncTxsArgs := sync.ArgsNewTransactionsSyncer{
		DataPools:      s.dataPool,
		Storages:       disabled.NewChainStorer(),
		Marshaller:     s.marshalizer,
		RequestHandler: s.requestHandler,
	}
	s.transactionsSyncer, err = sync.NewTransactionsSyncer(syncTxsArgs)
	if err != nil {
		return nil, err
	}

	eligibleNodesInfo, waitingNodesInfo := args.GenesisNodesConfig.InitialNodesInfo()

	eligibleValidators, err := nodesCoordinator.NodesInfoToValidators(eligibleNodesInfo)
	if err != nil {
		return nil, err
	}

	waitingValidators, err := nodesCoordinator.NodesInfoToValidators(waitingNodesInfo)
	if err != nil {
		return nil, err
	}

	consensusGroupCache, err := cache.NewLRUCache(consensusGroupCacheSize)
	if err != nil {
		return nil, err
	}

	s.memDB = disabled.CreateMemUnit()

	argsNodesCoordinator := nodesCoordinator.ArgNodesCoordinator{
		ChainParametersHandler:          args.ChainParametersHandler,
		Marshalizer:                     args.Marshalizer,
		Hasher:                          args.Hasher,
		Shuffler:                        args.NodeShuffler,
		EpochStartNotifier:              &disabled.EpochStartNotifier{},
		BootStorer:                      s.memDB,
		ShardIDAsObserver:               args.ShardIdAsObserver,
		NbShards:                        args.GenesisNodesConfig.NumberOfShards(),
		EligibleNodes:                   eligibleValidators,
		WaitingNodes:                    waitingValidators,
		SelfPublicKey:                   args.PubKey,
		ConsensusGroupCache:             consensusGroupCache,
		ShuffledOutHandler:              disabled.NewShuffledOutHandler(),
		ChanStopNode:                    args.ChanNodeStop,
		NodeTypeProvider:                args.NodeTypeProvider,
		IsFullArchive:                   args.IsFullArchive,
		EnableEpochsHandler:             args.EnableEpochsHandler,
		ValidatorInfoCacher:             s.dataPool.CurrentEpochValidatorInfo(),
		GenesisNodesSetupHandler:        s.genesisNodesConfig,
		NodesCoordinatorRegistryFactory: args.NodesCoordinatorRegistryFactory,
	}
	baseNodesCoordinator, err := nodesCoordinator.NewIndexHashedNodesCoordinator(argsNodesCoordinator)
	if err != nil {
		return nil, err
	}

	nodesCoord, err := nodesCoordinator.NewIndexHashedNodesCoordinatorWithRater(baseNodesCoordinator, args.ChanceComputer)
	if err != nil {
		return nil, err
	}

	s.nodeCoordinator = nodesCoord

	return s, nil
}

// NodesConfigFromMetaBlock synces and creates registry from epoch start metablock
func (s *syncValidatorStatus) NodesConfigFromMetaBlock(
	currMetaBlock data.HeaderHandler,
	prevMetaBlock data.HeaderHandler,
) (nodesCoordinator.NodesCoordinatorRegistryHandler, uint32, []*block.MiniBlock, error) {
	if currMetaBlock.GetNonce() > 1 && !currMetaBlock.IsStartOfEpochBlock() {
		return nil, 0, nil, epochStart.ErrNotEpochStartBlock
	}
	if prevMetaBlock.GetNonce() > 1 && !prevMetaBlock.IsStartOfEpochBlock() {
		return nil, 0, nil, epochStart.ErrNotEpochStartBlock
	}

	allMiniblocks := make([]*block.MiniBlock, 0)
	prevMiniBlocks, err := s.processValidatorChangesFor(prevMetaBlock)
	if err != nil {
		return nil, 0, nil, err
	}
	allMiniblocks = append(allMiniblocks, prevMiniBlocks...)

	currentMiniBlocks, err := s.processValidatorChangesFor(currMetaBlock)
	if err != nil {
		return nil, 0, nil, err
	}
	allMiniblocks = append(allMiniblocks, currentMiniBlocks...)

	selfShardId, err := s.nodeCoordinator.ShardIdForEpoch(currMetaBlock.GetEpoch())
	if err != nil {
		return nil, 0, nil, err
	}

	nodesConfig := s.nodeCoordinator.NodesCoordinatorToRegistry(currMetaBlock.GetEpoch())
	nodesConfig.SetCurrentEpoch(currMetaBlock.GetEpoch())
	return nodesConfig, selfShardId, allMiniblocks, nil
}

func (s *syncValidatorStatus) processValidatorChangesFor(metaBlock data.HeaderHandler) ([]*block.MiniBlock, error) {
	if metaBlock.GetEpoch() == 0 {
		// no need to process for genesis - already created
		return make([]*block.MiniBlock, 0), nil
	}

	blockBody, miniBlocks, err := s.getPeerBlockBodyForMeta(metaBlock)
	if err != nil {
		return nil, err
	}
	s.nodeCoordinator.EpochStartPrepare(metaBlock, blockBody)

	return miniBlocks, nil
}

func findPeerMiniBlockHeaders(metaBlock data.HeaderHandler) []data.MiniBlockHeaderHandler {
	shardMBHeaderHandlers := make([]data.MiniBlockHeaderHandler, 0)
	mbHeaderHandlers := metaBlock.GetMiniBlockHeaderHandlers()
	for i, mbHeader := range mbHeaderHandlers {
		if mbHeader.GetTypeInt32() != int32(block.PeerBlock) {
			continue
		}

		shardMBHeaderHandlers = append(shardMBHeaderHandlers, mbHeaderHandlers[i])
	}
	return shardMBHeaderHandlers
}

func (s *syncValidatorStatus) getPeerBlockBodyForMeta(
	metaBlock data.HeaderHandler,
) (data.BodyHandler, []*block.MiniBlock, error) {
	shardMBHeaders := findPeerMiniBlockHeaders(metaBlock)

	s.miniBlocksSyncer.ClearFields()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	err := s.miniBlocksSyncer.SyncPendingMiniBlocks(shardMBHeaders, ctx)
	cancel()
	if err != nil {
		return nil, nil, err
	}

	peerMiniBlocks, err := s.miniBlocksSyncer.GetMiniBlocks()
	if err != nil {
		return nil, nil, err
	}

	if s.enableEpochsHandler.IsFlagEnabledInEpoch(common.RefactorPeersMiniBlocksFlag, metaBlock.GetEpoch()) {
		s.transactionsSyncer.ClearFields()
		ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
		err = s.transactionsSyncer.SyncTransactionsFor(peerMiniBlocks, metaBlock.GetEpoch(), ctx)
		cancel()
		if err != nil {
			return nil, nil, err
		}

		validatorsInfo, err := s.transactionsSyncer.GetValidatorsInfo()
		if err != nil {
			return nil, nil, err
		}

		currentEpochValidatorInfoPool := s.dataPool.CurrentEpochValidatorInfo()
		for validatorInfoHash, validatorInfo := range validatorsInfo {
			currentEpochValidatorInfoPool.AddValidatorInfo([]byte(validatorInfoHash), validatorInfo)
		}
	}

	blockBody := &block.Body{MiniBlocks: make([]*block.MiniBlock, 0, len(peerMiniBlocks))}
	for _, mbHeader := range shardMBHeaders {
		blockBody.MiniBlocks = append(blockBody.MiniBlocks, peerMiniBlocks[string(mbHeader.GetHash())])
	}

	return blockBody, blockBody.MiniBlocks, nil
}

// IsInterfaceNil returns true if underlying object is nil
func (s *syncValidatorStatus) IsInterfaceNil() bool {
	return s == nil
}
