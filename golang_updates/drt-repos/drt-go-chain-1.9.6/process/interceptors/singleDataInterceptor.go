package interceptors

import (
	"errors"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/debug/handler"
	"github.com/TerraDharitri/drt-go-chain/p2p"
	"github.com/TerraDharitri/drt-go-chain/process"
)

// ArgSingleDataInterceptor is the argument for the single-data interceptor
type ArgSingleDataInterceptor struct {
	Topic                   string
	DataFactory             process.InterceptedDataFactory
	Processor               process.InterceptorProcessor
	Throttler               process.InterceptorThrottler
	AntifloodHandler        process.P2PAntifloodHandler
	WhiteListRequest        process.WhiteListHandler
	PreferredPeersHolder    process.PreferredPeersHolderHandler
	CurrentPeerId           core.PeerID
	InterceptedDataVerifier process.InterceptedDataVerifier
}

// SingleDataInterceptor is used for intercepting packed multi data
type SingleDataInterceptor struct {
	*baseDataInterceptor
	factory          process.InterceptedDataFactory
	whiteListRequest process.WhiteListHandler
}

// NewSingleDataInterceptor hooks a new interceptor for single data
func NewSingleDataInterceptor(arg ArgSingleDataInterceptor) (*SingleDataInterceptor, error) {
	if len(arg.Topic) == 0 {
		return nil, process.ErrEmptyTopic
	}
	if check.IfNil(arg.DataFactory) {
		return nil, process.ErrNilInterceptedDataFactory
	}
	if check.IfNil(arg.Processor) {
		return nil, process.ErrNilInterceptedDataProcessor
	}
	if check.IfNil(arg.Throttler) {
		return nil, process.ErrNilInterceptorThrottler
	}
	if check.IfNil(arg.AntifloodHandler) {
		return nil, process.ErrNilAntifloodHandler
	}
	if check.IfNil(arg.WhiteListRequest) {
		return nil, process.ErrNilWhiteListHandler
	}
	if check.IfNil(arg.PreferredPeersHolder) {
		return nil, process.ErrNilPreferredPeersHolder
	}
	if check.IfNil(arg.InterceptedDataVerifier) {
		return nil, process.ErrNilInterceptedDataVerifier
	}
	if len(arg.CurrentPeerId) == 0 {
		return nil, process.ErrEmptyPeerID
	}

	singleDataIntercept := &SingleDataInterceptor{
		baseDataInterceptor: &baseDataInterceptor{
			throttler:               arg.Throttler,
			antifloodHandler:        arg.AntifloodHandler,
			topic:                   arg.Topic,
			currentPeerId:           arg.CurrentPeerId,
			processor:               arg.Processor,
			preferredPeersHolder:    arg.PreferredPeersHolder,
			debugHandler:            handler.NewDisabledInterceptorDebugHandler(),
			interceptedDataVerifier: arg.InterceptedDataVerifier,
		},
		factory:          arg.DataFactory,
		whiteListRequest: arg.WhiteListRequest,
	}

	return singleDataIntercept, nil
}

// ProcessReceivedMessage is the callback func from the p2p.Messenger and will be called each time a new message was received
// (for the topic this validator was registered to)
func (sdi *SingleDataInterceptor) ProcessReceivedMessage(message p2p.MessageP2P, fromConnectedPeer core.PeerID, _ p2p.MessageHandler) ([]byte, error) {
	err := sdi.preProcessMesage(message, fromConnectedPeer)
	if err != nil {
		return nil, err
	}

	interceptedData, err := sdi.factory.Create(message.Data(), message.Peer())
	if err != nil {
		sdi.throttler.EndProcessing()

		// this situation is so severe that we need to black list the peers
		reason := "can not create object from received bytes, topic " + sdi.topic + ", error " + err.Error()
		sdi.antifloodHandler.BlacklistPeer(message.Peer(), reason, common.InvalidMessageBlacklistDuration)
		sdi.antifloodHandler.BlacklistPeer(fromConnectedPeer, reason, common.InvalidMessageBlacklistDuration)

		return nil, err
	}

	sdi.receivedDebugInterceptedData(interceptedData)
	err = sdi.interceptedDataVerifier.Verify(interceptedData)
	if err != nil {
		sdi.throttler.EndProcessing()
		sdi.processDebugInterceptedData(interceptedData, err)

		isWrongVersion := errors.Is(err, process.ErrInvalidTransactionVersion) || errors.Is(err, process.ErrInvalidChainID)
		if isWrongVersion {
			// this situation is so severe that we need to black list de peers
			reason := "wrong version of received intercepted data, topic " + sdi.topic + ", error " + err.Error()
			sdi.antifloodHandler.BlacklistPeer(message.Peer(), reason, common.InvalidMessageBlacklistDuration)
			sdi.antifloodHandler.BlacklistPeer(fromConnectedPeer, reason, common.InvalidMessageBlacklistDuration)
		}

		return nil, err
	}

	errOriginator := sdi.antifloodHandler.IsOriginatorEligibleForTopic(message.Peer(), sdi.topic)
	isWhiteListed := sdi.whiteListRequest.IsWhiteListed(interceptedData)
	if !isWhiteListed && errOriginator != nil {
		log.Trace("got message from peer on topic only for validators",
			"originator", p2p.PeerIdToShortString(message.Peer()), "topic",
			sdi.topic, "err", errOriginator)
		sdi.throttler.EndProcessing()
		return nil, errOriginator
	}

	messageID := interceptedData.Hash()
	isForCurrentShard := interceptedData.IsForCurrentShard()
	shouldProcess := isForCurrentShard || isWhiteListed
	if !shouldProcess {
		sdi.throttler.EndProcessing()
		log.Trace("intercepted data is for other shards",
			"pid", p2p.MessageOriginatorPid(message),
			"seq no", p2p.MessageOriginatorSeq(message),
			"topic", message.Topic(),
			"hash", interceptedData.Hash(),
			"is for current shard", isForCurrentShard,
			"is white listed", isWhiteListed,
		)

		return messageID, nil
	}

	go func() {
		sdi.processInterceptedData(interceptedData, message)
		sdi.throttler.EndProcessing()
	}()

	return messageID, nil
}

// RegisterHandler registers a callback function to be notified on received data
func (sdi *SingleDataInterceptor) RegisterHandler(handler func(topic string, hash []byte, data interface{})) {
	sdi.processor.RegisterHandler(handler)
}

// Close returns nil
func (sdi *SingleDataInterceptor) Close() error {
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (sdi *SingleDataInterceptor) IsInterfaceNil() bool {
	return sdi == nil
}
