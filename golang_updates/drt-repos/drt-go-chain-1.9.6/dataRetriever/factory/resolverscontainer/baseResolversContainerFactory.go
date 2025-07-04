package resolverscontainer

import (
	"fmt"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data/typeConverters"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	logger "github.com/TerraDharitri/drt-go-chain-logger"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever/resolvers"
	topicsender "github.com/TerraDharitri/drt-go-chain/dataRetriever/topicSender"
	"github.com/TerraDharitri/drt-go-chain/p2p"
	"github.com/TerraDharitri/drt-go-chain/process/factory"
	"github.com/TerraDharitri/drt-go-chain/sharding"
)

// EmptyExcludePeersOnTopic is an empty topic
const EmptyExcludePeersOnTopic = ""

var log = logger.GetOrCreate("dataRetriever/factory/resolverscontainer")

type baseResolversContainerFactory struct {
	container                       dataRetriever.ResolversContainer
	shardCoordinator                sharding.Coordinator
	mainMessenger                   p2p.Messenger
	fullArchiveMessenger            p2p.Messenger
	store                           dataRetriever.StorageService
	marshalizer                     marshal.Marshalizer
	dataPools                       dataRetriever.PoolsHolder
	uint64ByteSliceConverter        typeConverters.Uint64ByteSliceConverter
	dataPacker                      dataRetriever.DataPacker
	triesContainer                  common.TriesHolder
	inputAntifloodHandler           dataRetriever.P2PAntifloodHandler
	outputAntifloodHandler          dataRetriever.P2PAntifloodHandler
	throttler                       dataRetriever.ResolverThrottler
	trieNodesThrottler              dataRetriever.ResolverThrottler
	intraShardTopic                 string
	isFullHistoryNode               bool
	mainPreferredPeersHolder        dataRetriever.PreferredPeersHolderHandler
	fullArchivePreferredPeersHolder dataRetriever.PreferredPeersHolderHandler
	payloadValidator                dataRetriever.PeerAuthenticationPayloadValidator
}

func (brcf *baseResolversContainerFactory) checkParams() error {
	if check.IfNil(brcf.shardCoordinator) {
		return dataRetriever.ErrNilShardCoordinator
	}
	if check.IfNil(brcf.mainMessenger) {
		return fmt.Errorf("%w for main network", dataRetriever.ErrNilMessenger)
	}
	if check.IfNil(brcf.fullArchiveMessenger) {
		return fmt.Errorf("%w for full archive network", dataRetriever.ErrNilMessenger)
	}
	if check.IfNil(brcf.store) {
		return dataRetriever.ErrNilStore
	}
	if check.IfNil(brcf.marshalizer) {
		return dataRetriever.ErrNilMarshalizer
	}
	if check.IfNil(brcf.dataPools) {
		return dataRetriever.ErrNilDataPoolHolder
	}
	if check.IfNil(brcf.uint64ByteSliceConverter) {
		return dataRetriever.ErrNilUint64ByteSliceConverter
	}
	if check.IfNil(brcf.dataPacker) {
		return dataRetriever.ErrNilDataPacker
	}
	if check.IfNil(brcf.triesContainer) {
		return dataRetriever.ErrNilTrieDataGetter
	}
	if check.IfNil(brcf.inputAntifloodHandler) {
		return fmt.Errorf("%w for input", dataRetriever.ErrNilAntifloodHandler)
	}
	if check.IfNil(brcf.outputAntifloodHandler) {
		return fmt.Errorf("%w for output", dataRetriever.ErrNilAntifloodHandler)
	}
	if check.IfNil(brcf.throttler) {
		return fmt.Errorf("%w for the main throttler", dataRetriever.ErrNilThrottler)
	}
	if check.IfNil(brcf.trieNodesThrottler) {
		return fmt.Errorf("%w for the trie nodes throttler", dataRetriever.ErrNilThrottler)
	}
	if check.IfNil(brcf.mainPreferredPeersHolder) {
		return fmt.Errorf("%w for main network", dataRetriever.ErrNilPreferredPeersHolder)
	}
	if check.IfNil(brcf.fullArchivePreferredPeersHolder) {
		return fmt.Errorf("%w for full archive network", dataRetriever.ErrNilPreferredPeersHolder)
	}

	return nil
}

func (brcf *baseResolversContainerFactory) generateTxResolvers(
	topic string,
	unit dataRetriever.UnitType,
	dataPool dataRetriever.ShardedDataCacherNotifier,
) error {

	shardC := brcf.shardCoordinator
	noOfShards := shardC.NumberOfShards()

	keys := make([]string, noOfShards+1)
	resolverSlice := make([]dataRetriever.Resolver, noOfShards+1)

	for idx := uint32(0); idx < noOfShards; idx++ {
		identifierTx := topic + shardC.CommunicationIdentifier(idx)
		excludePeersFromTopic := topic + shardC.CommunicationIdentifier(shardC.SelfId())

		resolver, err := brcf.createTxResolver(identifierTx, excludePeersFromTopic, unit, dataPool, idx)
		if err != nil {
			return err
		}

		resolverSlice[idx] = resolver
		keys[idx] = identifierTx
	}

	identifierTx := topic + shardC.CommunicationIdentifier(core.MetachainShardId)
	excludePeersFromTopic := topic + shardC.CommunicationIdentifier(shardC.SelfId())

	resolver, err := brcf.createTxResolver(identifierTx, excludePeersFromTopic, unit, dataPool, core.MetachainShardId)
	if err != nil {
		return err
	}

	resolverSlice[noOfShards] = resolver
	keys[noOfShards] = identifierTx

	return brcf.container.AddMultiple(keys, resolverSlice)
}

func (brcf *baseResolversContainerFactory) createTxResolver(
	topic string,
	excludedTopic string,
	unit dataRetriever.UnitType,
	dataPool dataRetriever.ShardedDataCacherNotifier,
	targetShardID uint32,
) (dataRetriever.Resolver, error) {

	txStorer, err := brcf.store.GetStorer(unit)
	if err != nil {
		return nil, err
	}

	resolverSender, err := brcf.createOneResolverSenderWithSpecifiedNumRequests(topic, excludedTopic, targetShardID)
	if err != nil {
		return nil, err
	}

	arg := resolvers.ArgTxResolver{
		ArgBaseResolver: resolvers.ArgBaseResolver{
			SenderResolver:   resolverSender,
			Marshaller:       brcf.marshalizer,
			AntifloodHandler: brcf.inputAntifloodHandler,
			Throttler:        brcf.throttler,
		},
		TxPool:            dataPool,
		TxStorage:         txStorer,
		DataPacker:        brcf.dataPacker,
		IsFullHistoryNode: brcf.isFullHistoryNode,
	}
	resolver, err := resolvers.NewTxResolver(arg)
	if err != nil {
		return nil, err
	}

	err = brcf.mainMessenger.RegisterMessageProcessor(resolver.RequestTopic(), common.DefaultResolversIdentifier, resolver)
	if err != nil {
		return nil, err
	}

	err = brcf.fullArchiveMessenger.RegisterMessageProcessor(resolver.RequestTopic(), common.DefaultResolversIdentifier, resolver)
	if err != nil {
		return nil, err
	}

	return resolver, nil
}

func (brcf *baseResolversContainerFactory) generateMiniBlocksResolvers() error {
	shardC := brcf.shardCoordinator
	noOfShards := shardC.NumberOfShards()
	keys := make([]string, noOfShards+2)
	resolverSlice := make([]dataRetriever.Resolver, noOfShards+2)

	for idx := uint32(0); idx < noOfShards; idx++ {
		identifierMiniBlocks := factory.MiniBlocksTopic + shardC.CommunicationIdentifier(idx)
		excludePeersFromTopic := factory.MiniBlocksTopic + shardC.CommunicationIdentifier(shardC.SelfId())

		resolver, err := brcf.createMiniBlocksResolver(identifierMiniBlocks, excludePeersFromTopic, idx)
		if err != nil {
			return err
		}

		resolverSlice[idx] = resolver
		keys[idx] = identifierMiniBlocks
	}

	identifierMiniBlocks := factory.MiniBlocksTopic + shardC.CommunicationIdentifier(core.MetachainShardId)
	excludePeersFromTopic := factory.MiniBlocksTopic + shardC.CommunicationIdentifier(shardC.SelfId())

	resolver, err := brcf.createMiniBlocksResolver(identifierMiniBlocks, excludePeersFromTopic, core.MetachainShardId)
	if err != nil {
		return err
	}

	resolverSlice[noOfShards] = resolver
	keys[noOfShards] = identifierMiniBlocks

	identifierAllShardMiniBlocks := factory.MiniBlocksTopic + shardC.CommunicationIdentifier(core.AllShardId)
	allShardMiniblocksResolver, err := brcf.createMiniBlocksResolver(identifierAllShardMiniBlocks, EmptyExcludePeersOnTopic, brcf.shardCoordinator.SelfId())
	if err != nil {
		return err
	}

	resolverSlice[noOfShards+1] = allShardMiniblocksResolver
	keys[noOfShards+1] = identifierAllShardMiniBlocks

	return brcf.container.AddMultiple(keys, resolverSlice)
}

func (brcf *baseResolversContainerFactory) createMiniBlocksResolver(
	topic string,
	excludedTopic string,
	targetShardID uint32,
) (dataRetriever.Resolver, error) {
	miniBlocksStorer, err := brcf.store.GetStorer(dataRetriever.MiniBlockUnit)
	if err != nil {
		return nil, err
	}

	resolverSender, err := brcf.createOneResolverSenderWithSpecifiedNumRequests(topic, excludedTopic, targetShardID)
	if err != nil {
		return nil, err
	}

	arg := resolvers.ArgMiniblockResolver{
		ArgBaseResolver: resolvers.ArgBaseResolver{
			SenderResolver:   resolverSender,
			Marshaller:       brcf.marshalizer,
			AntifloodHandler: brcf.inputAntifloodHandler,
			Throttler:        brcf.throttler,
		},
		MiniBlockPool:     brcf.dataPools.MiniBlocks(),
		MiniBlockStorage:  miniBlocksStorer,
		DataPacker:        brcf.dataPacker,
		IsFullHistoryNode: brcf.isFullHistoryNode,
	}
	txBlkResolver, err := resolvers.NewMiniblockResolver(arg)
	if err != nil {
		return nil, err
	}

	err = brcf.mainMessenger.RegisterMessageProcessor(txBlkResolver.RequestTopic(), common.DefaultResolversIdentifier, txBlkResolver)
	if err != nil {
		return nil, err
	}

	err = brcf.fullArchiveMessenger.RegisterMessageProcessor(txBlkResolver.RequestTopic(), common.DefaultResolversIdentifier, txBlkResolver)
	if err != nil {
		return nil, err
	}

	return txBlkResolver, nil
}

func (brcf *baseResolversContainerFactory) generatePeerAuthenticationResolver() error {
	identifierPeerAuth := common.PeerAuthenticationTopic
	shardC := brcf.shardCoordinator
	resolverSender, err := brcf.createOneResolverSenderWithSpecifiedNumRequests(identifierPeerAuth, EmptyExcludePeersOnTopic, shardC.SelfId())
	if err != nil {
		return err
	}

	arg := resolvers.ArgPeerAuthenticationResolver{
		ArgBaseResolver: resolvers.ArgBaseResolver{
			SenderResolver:   resolverSender,
			Marshaller:       brcf.marshalizer,
			AntifloodHandler: brcf.inputAntifloodHandler,
			Throttler:        brcf.throttler,
		},
		PeerAuthenticationPool: brcf.dataPools.PeerAuthentications(),
		DataPacker:             brcf.dataPacker,
		PayloadValidator:       brcf.payloadValidator,
	}
	peerAuthResolver, err := resolvers.NewPeerAuthenticationResolver(arg)
	if err != nil {
		return err
	}

	err = brcf.mainMessenger.RegisterMessageProcessor(peerAuthResolver.RequestTopic(), common.DefaultResolversIdentifier, peerAuthResolver)
	if err != nil {
		return err
	}

	err = brcf.fullArchiveMessenger.RegisterMessageProcessor(peerAuthResolver.RequestTopic(), common.DefaultResolversIdentifier, peerAuthResolver)
	if err != nil {
		return err
	}

	return brcf.container.Add(identifierPeerAuth, peerAuthResolver)
}

func (brcf *baseResolversContainerFactory) createOneResolverSenderWithSpecifiedNumRequests(
	topic string,
	excludedTopic string,
	targetShardId uint32,
) (dataRetriever.TopicResolverSender, error) {

	log.Trace("baseResolversContainerFactory.createOneResolverSenderWithSpecifiedNumRequests",
		"topic", topic, "intraShardTopic", brcf.intraShardTopic, "excludedTopic", excludedTopic)

	arg := topicsender.ArgTopicResolverSender{
		ArgBaseTopicSender: topicsender.ArgBaseTopicSender{
			MainMessenger:                   brcf.mainMessenger,
			FullArchiveMessenger:            brcf.fullArchiveMessenger,
			TopicName:                       topic,
			OutputAntiflooder:               brcf.outputAntifloodHandler,
			MainPreferredPeersHolder:        brcf.mainPreferredPeersHolder,
			FullArchivePreferredPeersHolder: brcf.fullArchivePreferredPeersHolder,
			TargetShardId:                   targetShardId,
		},
	}
	// TODO instantiate topic sender resolver with the shard IDs for which this resolver is supposed to serve the data
	// this will improve the serving of transactions as the searching will be done only on 2 sharded data units
	resolverSender, err := topicsender.NewTopicResolverSender(arg)
	if err != nil {
		return nil, err
	}

	return resolverSender, nil
}

func (brcf *baseResolversContainerFactory) createTrieNodesResolver(
	topic string,
	trieId string,
	targetShardID uint32,
) (dataRetriever.Resolver, error) {
	resolverSender, err := brcf.createOneResolverSenderWithSpecifiedNumRequests(
		topic,
		EmptyExcludePeersOnTopic,
		targetShardID,
	)
	if err != nil {
		return nil, err
	}

	trie := brcf.triesContainer.Get([]byte(trieId))
	argTrie := resolvers.ArgTrieNodeResolver{
		ArgBaseResolver: resolvers.ArgBaseResolver{
			SenderResolver:   resolverSender,
			Marshaller:       brcf.marshalizer,
			AntifloodHandler: brcf.inputAntifloodHandler,
			Throttler:        brcf.trieNodesThrottler,
		},
		TrieDataGetter: trie,
	}
	resolver, err := resolvers.NewTrieNodeResolver(argTrie)
	if err != nil {
		return nil, err
	}

	err = brcf.mainMessenger.RegisterMessageProcessor(resolver.RequestTopic(), common.DefaultResolversIdentifier, resolver)
	if err != nil {
		return nil, err
	}

	err = brcf.fullArchiveMessenger.RegisterMessageProcessor(resolver.RequestTopic(), common.DefaultResolversIdentifier, resolver)
	if err != nil {
		return nil, err
	}

	return resolver, nil
}

func (brcf *baseResolversContainerFactory) generateValidatorInfoResolver() error {
	identifierValidatorInfo := common.ValidatorInfoTopic
	shardC := brcf.shardCoordinator
	resolverSender, err := brcf.createOneResolverSenderWithSpecifiedNumRequests(identifierValidatorInfo, EmptyExcludePeersOnTopic, shardC.SelfId())
	if err != nil {
		return err
	}

	validatorInfoStorage, err := brcf.store.GetStorer(dataRetriever.UnsignedTransactionUnit)
	if err != nil {
		return err
	}

	arg := resolvers.ArgValidatorInfoResolver{
		SenderResolver:       resolverSender,
		Marshaller:           brcf.marshalizer,
		AntifloodHandler:     brcf.inputAntifloodHandler,
		Throttler:            brcf.throttler,
		ValidatorInfoPool:    brcf.dataPools.ValidatorsInfo(),
		ValidatorInfoStorage: validatorInfoStorage,
		DataPacker:           brcf.dataPacker,
		IsFullHistoryNode:    brcf.isFullHistoryNode,
	}
	validatorInfoResolver, err := resolvers.NewValidatorInfoResolver(arg)
	if err != nil {
		return err
	}

	err = brcf.mainMessenger.RegisterMessageProcessor(validatorInfoResolver.RequestTopic(), common.DefaultResolversIdentifier, validatorInfoResolver)
	if err != nil {
		return err
	}

	err = brcf.fullArchiveMessenger.RegisterMessageProcessor(validatorInfoResolver.RequestTopic(), common.DefaultResolversIdentifier, validatorInfoResolver)
	if err != nil {
		return err
	}

	return brcf.container.Add(identifierValidatorInfo, validatorInfoResolver)
}

func (brcf *baseResolversContainerFactory) createEquivalentProofsResolver(
	topic string,
	targetShardID uint32,
) (dataRetriever.Resolver, error) {
	resolverSender, err := brcf.createOneResolverSenderWithSpecifiedNumRequests(
		topic,
		EmptyExcludePeersOnTopic,
		targetShardID,
	)
	if err != nil {
		return nil, err
	}

	arg := resolvers.ArgEquivalentProofsResolver{
		ArgBaseResolver: resolvers.ArgBaseResolver{
			SenderResolver:   resolverSender,
			Marshaller:       brcf.marshalizer,
			AntifloodHandler: brcf.inputAntifloodHandler,
			Throttler:        brcf.trieNodesThrottler,
		},
		DataPacker:           brcf.dataPacker,
		Storage:              brcf.store,
		EquivalentProofsPool: brcf.dataPools.Proofs(),
		NonceConverter:       brcf.uint64ByteSliceConverter,
		IsFullHistoryNode:    brcf.isFullHistoryNode,
	}
	resolver, err := resolvers.NewEquivalentProofsResolver(arg)
	if err != nil {
		return nil, err
	}

	err = brcf.mainMessenger.RegisterMessageProcessor(resolver.RequestTopic(), common.DefaultResolversIdentifier, resolver)
	if err != nil {
		return nil, err
	}

	err = brcf.fullArchiveMessenger.RegisterMessageProcessor(resolver.RequestTopic(), common.DefaultResolversIdentifier, resolver)
	if err != nil {
		return nil, err
	}

	return resolver, nil
}
