package requesterscontainer

import (
	"fmt"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data/typeConverters"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	logger "github.com/TerraDharitri/drt-go-chain-logger"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever/requestHandlers/requesters"
	topicsender "github.com/TerraDharitri/drt-go-chain/dataRetriever/topicSender"
	"github.com/TerraDharitri/drt-go-chain/p2p"
	"github.com/TerraDharitri/drt-go-chain/process/factory"
	"github.com/TerraDharitri/drt-go-chain/sharding"
)

// EmptyExcludePeersOnTopic is an empty topic
const EmptyExcludePeersOnTopic = ""

var log = logger.GetOrCreate("dataRetriever/factory/requesterscontainer")

type baseRequestersContainerFactory struct {
	container                       dataRetriever.RequestersContainer
	shardCoordinator                sharding.Coordinator
	mainMessenger                   p2p.Messenger
	fullArchiveMessenger            p2p.Messenger
	marshaller                      marshal.Marshalizer
	uint64ByteSliceConverter        typeConverters.Uint64ByteSliceConverter
	intRandomizer                   dataRetriever.IntRandomizer
	outputAntifloodHandler          dataRetriever.P2PAntifloodHandler
	intraShardTopic                 string
	currentNetworkEpochProvider     dataRetriever.CurrentNetworkEpochProviderHandler
	mainPreferredPeersHolder        dataRetriever.PreferredPeersHolderHandler
	fullArchivePreferredPeersHolder dataRetriever.PreferredPeersHolderHandler
	peersRatingHandler              dataRetriever.PeersRatingHandler
	enableEpochsHandler             common.EnableEpochsHandler
	numCrossShardPeers              int
	numIntraShardPeers              int
	numTotalPeers                   int
	numFullHistoryPeers             int
}

func (brcf *baseRequestersContainerFactory) checkParams() error {
	if check.IfNil(brcf.shardCoordinator) {
		return dataRetriever.ErrNilShardCoordinator
	}
	if check.IfNil(brcf.mainMessenger) {
		return fmt.Errorf("%w on main network", dataRetriever.ErrNilMessenger)
	}
	if check.IfNil(brcf.fullArchiveMessenger) {
		return fmt.Errorf("%w on full archive network", dataRetriever.ErrNilMessenger)
	}
	if check.IfNil(brcf.marshaller) {
		return dataRetriever.ErrNilMarshalizer
	}
	if check.IfNil(brcf.uint64ByteSliceConverter) {
		return dataRetriever.ErrNilUint64ByteSliceConverter
	}
	if check.IfNil(brcf.outputAntifloodHandler) {
		return fmt.Errorf("%w for output", dataRetriever.ErrNilAntifloodHandler)
	}
	if check.IfNil(brcf.currentNetworkEpochProvider) {
		return dataRetriever.ErrNilCurrentNetworkEpochProvider
	}
	if check.IfNil(brcf.mainPreferredPeersHolder) {
		return fmt.Errorf("%w on main network", dataRetriever.ErrNilPreferredPeersHolder)
	}
	if check.IfNil(brcf.fullArchivePreferredPeersHolder) {
		return fmt.Errorf("%w on full archive network", dataRetriever.ErrNilPreferredPeersHolder)
	}
	if check.IfNil(brcf.peersRatingHandler) {
		return dataRetriever.ErrNilPeersRatingHandler
	}
	if brcf.numCrossShardPeers <= 0 {
		return fmt.Errorf("%w for numCrossShardPeers", dataRetriever.ErrInvalidValue)
	}
	if brcf.numTotalPeers <= brcf.numCrossShardPeers {
		return fmt.Errorf("%w for numTotalPeers", dataRetriever.ErrInvalidValue)
	}
	if brcf.numFullHistoryPeers <= 0 {
		return fmt.Errorf("%w for numFullHistoryPeers", dataRetriever.ErrInvalidValue)
	}

	return nil
}

func (brcf *baseRequestersContainerFactory) generateCommonRequesters() error {
	err := brcf.generateTxRequesters(factory.TransactionTopic)
	if err != nil {
		return err
	}

	err = brcf.generateTxRequesters(factory.UnsignedTransactionTopic)
	if err != nil {
		return err
	}

	err = brcf.generateMiniBlocksRequesters()
	if err != nil {
		return err
	}

	err = brcf.generatePeerAuthenticationRequester()
	if err != nil {
		return err
	}

	err = brcf.generateValidatorInfoRequester()
	if err != nil {
		return err
	}

	return nil
}

func (brcf *baseRequestersContainerFactory) generateTxRequesters(topic string) error {

	shardC := brcf.shardCoordinator
	noOfShards := shardC.NumberOfShards()

	keys := make([]string, noOfShards+1)
	requestersSlice := make([]dataRetriever.Requester, noOfShards+1)

	for idx := uint32(0); idx < noOfShards; idx++ {
		identifierTx := topic + shardC.CommunicationIdentifier(idx)
		excludePeersFromTopic := topic + shardC.CommunicationIdentifier(shardC.SelfId())

		requester, err := brcf.createTxRequester(identifierTx, excludePeersFromTopic, idx, brcf.numCrossShardPeers, brcf.numIntraShardPeers)
		if err != nil {
			return err
		}

		requestersSlice[idx] = requester
		keys[idx] = identifierTx
	}

	identifierTx := topic + shardC.CommunicationIdentifier(core.MetachainShardId)
	excludePeersFromTopic := topic + shardC.CommunicationIdentifier(shardC.SelfId())

	requester, err := brcf.createTxRequester(identifierTx, excludePeersFromTopic, core.MetachainShardId, brcf.numCrossShardPeers, brcf.numIntraShardPeers)
	if err != nil {
		return err
	}

	requestersSlice[noOfShards] = requester
	keys[noOfShards] = identifierTx

	return brcf.container.AddMultiple(keys, requestersSlice)
}

func (brcf *baseRequestersContainerFactory) createTxRequester(
	topic string,
	excludedTopic string,
	targetShardID uint32,
	numCrossShardPeers int,
	numIntraShardPeers int,
) (dataRetriever.Requester, error) {

	requestSender, err := brcf.createOneRequestSenderWithSpecifiedNumRequests(topic, excludedTopic, targetShardID, numCrossShardPeers, numIntraShardPeers)
	if err != nil {
		return nil, err
	}

	arg := requesters.ArgTransactionRequester{
		ArgBaseRequester: requesters.ArgBaseRequester{
			RequestSender: requestSender,
			Marshaller:    brcf.marshaller,
		},
	}
	return requesters.NewTransactionRequester(arg)
}

func (brcf *baseRequestersContainerFactory) generateMiniBlocksRequesters() error {
	shardC := brcf.shardCoordinator
	noOfShards := shardC.NumberOfShards()
	keys := make([]string, noOfShards+2)
	requestersSlice := make([]dataRetriever.Requester, noOfShards+2)

	for idx := uint32(0); idx < noOfShards; idx++ {
		identifierMiniBlocks := factory.MiniBlocksTopic + shardC.CommunicationIdentifier(idx)
		excludePeersFromTopic := factory.MiniBlocksTopic + shardC.CommunicationIdentifier(shardC.SelfId())

		requester, err := brcf.createMiniBlocksRequester(identifierMiniBlocks, excludePeersFromTopic, idx, brcf.numCrossShardPeers, brcf.numIntraShardPeers)
		if err != nil {
			return err
		}

		requestersSlice[idx] = requester
		keys[idx] = identifierMiniBlocks
	}

	identifierMiniBlocks := factory.MiniBlocksTopic + shardC.CommunicationIdentifier(core.MetachainShardId)
	excludePeersFromTopic := factory.MiniBlocksTopic + shardC.CommunicationIdentifier(shardC.SelfId())

	requester, err := brcf.createMiniBlocksRequester(identifierMiniBlocks, excludePeersFromTopic, core.MetachainShardId, brcf.numCrossShardPeers, brcf.numIntraShardPeers)
	if err != nil {
		return err
	}

	requestersSlice[noOfShards] = requester
	keys[noOfShards] = identifierMiniBlocks

	identifierAllShardMiniBlocks := factory.MiniBlocksTopic + shardC.CommunicationIdentifier(core.AllShardId)
	allShardMiniblocksResolver, err := brcf.createMiniBlocksRequester(identifierAllShardMiniBlocks, EmptyExcludePeersOnTopic, brcf.shardCoordinator.SelfId(), brcf.numCrossShardPeers, brcf.numIntraShardPeers)
	if err != nil {
		return err
	}

	requestersSlice[noOfShards+1] = allShardMiniblocksResolver
	keys[noOfShards+1] = identifierAllShardMiniBlocks

	return brcf.container.AddMultiple(keys, requestersSlice)
}

func (brcf *baseRequestersContainerFactory) createMiniBlocksRequester(
	topic string,
	excludedTopic string,
	targetShardID uint32,
	numCrossShardPeers int,
	numIntraShardPeers int,
) (dataRetriever.Requester, error) {
	requestSender, err := brcf.createOneRequestSenderWithSpecifiedNumRequests(topic, excludedTopic, targetShardID, numCrossShardPeers, numIntraShardPeers)
	if err != nil {
		return nil, err
	}

	arg := requesters.ArgMiniblockRequester{
		ArgBaseRequester: requesters.ArgBaseRequester{
			RequestSender: requestSender,
			Marshaller:    brcf.marshaller,
		},
	}
	return requesters.NewMiniblockRequester(arg)
}

func (brcf *baseRequestersContainerFactory) generatePeerAuthenticationRequester() error {
	identifierPeerAuth := common.PeerAuthenticationTopic
	shardC := brcf.shardCoordinator
	requestSender, err := brcf.createOneRequestSenderWithSpecifiedNumRequests(identifierPeerAuth, EmptyExcludePeersOnTopic, shardC.SelfId(), brcf.numCrossShardPeers, brcf.numIntraShardPeers)
	if err != nil {
		return err
	}

	arg := requesters.ArgPeerAuthenticationRequester{
		ArgBaseRequester: requesters.ArgBaseRequester{
			RequestSender: requestSender,
			Marshaller:    brcf.marshaller,
		},
	}
	requester, err := requesters.NewPeerAuthenticationRequester(arg)
	if err != nil {
		return err
	}

	return brcf.container.Add(identifierPeerAuth, requester)
}

func (brcf *baseRequestersContainerFactory) createOneRequestSenderWithSpecifiedNumRequests(
	topic string,
	excludedTopic string,
	targetShardId uint32,
	numCrossShardPeers int,
	numIntraShardPeers int,
) (dataRetriever.TopicRequestSender, error) {

	log.Trace("baseRequestersContainerFactory.createOneRequestSenderWithSpecifiedNumRequests",
		"topic", topic, "intraShardTopic", brcf.intraShardTopic, "excludedTopic", excludedTopic,
		"numCrossShardPeers", numCrossShardPeers, "numIntraShardPeers", numIntraShardPeers)

	peerListCreator, err := topicsender.NewDiffPeerListCreator(brcf.mainMessenger, topic, brcf.intraShardTopic, excludedTopic)
	if err != nil {
		return nil, err
	}

	arg := topicsender.ArgTopicRequestSender{
		ArgBaseTopicSender: topicsender.ArgBaseTopicSender{
			MainMessenger:                   brcf.mainMessenger,
			FullArchiveMessenger:            brcf.fullArchiveMessenger,
			TopicName:                       topic,
			OutputAntiflooder:               brcf.outputAntifloodHandler,
			MainPreferredPeersHolder:        brcf.mainPreferredPeersHolder,
			FullArchivePreferredPeersHolder: brcf.fullArchivePreferredPeersHolder,
			TargetShardId:                   targetShardId,
		},
		Marshaller:                  brcf.marshaller,
		Randomizer:                  brcf.intRandomizer,
		PeerListCreator:             peerListCreator,
		NumIntraShardPeers:          numIntraShardPeers,
		NumCrossShardPeers:          numCrossShardPeers,
		NumFullHistoryPeers:         brcf.numFullHistoryPeers,
		CurrentNetworkEpochProvider: brcf.currentNetworkEpochProvider,
		SelfShardIdProvider:         brcf.shardCoordinator,
		PeersRatingHandler:          brcf.peersRatingHandler,
	}
	return topicsender.NewTopicRequestSender(arg)
}

func (brcf *baseRequestersContainerFactory) createTrieNodesRequester(
	topic string,
	numCrossShardPeers int,
	numIntraShardPeers int,
	targetShardID uint32,
) (dataRetriever.Requester, error) {
	requestSender, err := brcf.createOneRequestSenderWithSpecifiedNumRequests(
		topic,
		EmptyExcludePeersOnTopic,
		targetShardID,
		numCrossShardPeers,
		numIntraShardPeers,
	)
	if err != nil {
		return nil, err
	}

	arg := requesters.ArgTrieNodeRequester{
		ArgBaseRequester: requesters.ArgBaseRequester{
			RequestSender: requestSender,
			Marshaller:    brcf.marshaller,
		},
	}
	return requesters.NewTrieNodeRequester(arg)
}

func (brcf *baseRequestersContainerFactory) generateValidatorInfoRequester() error {
	identifierValidatorInfo := common.ValidatorInfoTopic
	shardC := brcf.shardCoordinator
	requestSender, err := brcf.createOneRequestSenderWithSpecifiedNumRequests(identifierValidatorInfo, EmptyExcludePeersOnTopic, shardC.SelfId(), brcf.numCrossShardPeers, brcf.numIntraShardPeers)
	if err != nil {
		return err
	}

	arg := requesters.ArgValidatorInfoRequester{
		ArgBaseRequester: requesters.ArgBaseRequester{
			RequestSender: requestSender,
			Marshaller:    brcf.marshaller,
		},
	}
	requester, err := requesters.NewValidatorInfoRequester(arg)
	if err != nil {
		return err
	}

	return brcf.container.Add(identifierValidatorInfo, requester)
}

func (brcf *baseRequestersContainerFactory) createEquivalentProofsRequester(
	topic string,
	numCrossShardPeers int,
	numIntraShardPeers int,
	targetShardID uint32,
) (dataRetriever.Requester, error) {
	requestSender, err := brcf.createOneRequestSenderWithSpecifiedNumRequests(
		topic,
		EmptyExcludePeersOnTopic,
		targetShardID,
		numCrossShardPeers,
		numIntraShardPeers,
	)
	if err != nil {
		return nil, err
	}

	arg := requesters.ArgEquivalentProofsRequester{
		ArgBaseRequester: requesters.ArgBaseRequester{
			RequestSender: requestSender,
			Marshaller:    brcf.marshaller,
		},
		EnableEpochsHandler: brcf.enableEpochsHandler,
	}
	return requesters.NewEquivalentProofsRequester(arg)
}
