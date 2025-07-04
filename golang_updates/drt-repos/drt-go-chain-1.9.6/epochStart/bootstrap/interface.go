package bootstrap

import (
	"context"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
)

// StartOfEpochNodesConfigHandler defines the methods to process nodesConfig from epoch start metablocks
type StartOfEpochNodesConfigHandler interface {
	NodesConfigFromMetaBlock(currMetaBlock data.HeaderHandler, prevMetaBlock data.HeaderHandler) (nodesCoordinator.NodesCoordinatorRegistryHandler, uint32, []*block.MiniBlock, error)
	IsInterfaceNil() bool
}

// EpochStartMetaBlockInterceptorProcessor defines the methods to sync an epoch start metablock
type EpochStartMetaBlockInterceptorProcessor interface {
	process.InterceptorProcessor
	GetEpochStartMetaBlock(ctx context.Context) (data.MetaHeaderHandler, error)
}

// StartInEpochNodesCoordinator defines the methods to process and save nodesCoordinator information to storage
type StartInEpochNodesCoordinator interface {
	EpochStartPrepare(metaHdr data.HeaderHandler, body data.BodyHandler)
	NodesCoordinatorToRegistry(epoch uint32) nodesCoordinator.NodesCoordinatorRegistryHandler
	ShardIdForEpoch(epoch uint32) (uint32, error)
	IsInterfaceNil() bool
}

// Messenger defines which methods a p2p messenger should implement
type Messenger interface {
	dataRetriever.MessageHandler
	dataRetriever.TopicHandler
	UnregisterMessageProcessor(topic string, identifier string) error
	UnregisterAllMessageProcessors() error
	UnJoinAllTopics() error
	ConnectedPeers() []core.PeerID
	Verify(payload []byte, pid core.PeerID, signature []byte) error
	Broadcast(topic string, buff []byte)
	BroadcastUsingPrivateKey(topic string, buff []byte, pid core.PeerID, skBytes []byte)
	Sign(payload []byte) ([]byte, error)
	SignUsingPrivateKey(skBytes []byte, payload []byte) ([]byte, error)
}

// RequestHandler defines which methods a request handler should implement
type RequestHandler interface {
	RequestStartOfEpochMetaBlock(epoch uint32)
	RequestMetaHeaderByNonce(nonce uint64)
	SetNumPeersToQuery(topic string, intra int, cross int) error
	GetNumPeersToQuery(topic string) (int, int, error)
	RequestEquivalentProofByNonce(headerShard uint32, headerNonce uint64)
	RequestEquivalentProofByHash(headerShard uint32, headerHash []byte)
	SetEpoch(epoch uint32)
	IsInterfaceNil() bool
}

// NodeTypeProviderHandler defines the actions needed for a component that can handle the node type
type NodeTypeProviderHandler interface {
	SetType(nodeType core.NodeType)
	GetType() core.NodeType
	IsInterfaceNil() bool
}

// ProofsPool defines the behaviour of a proofs pool components
type ProofsPool interface {
	RegisterHandler(handler func(headerProof data.HeaderProofHandler))
	GetProof(shardID uint32, headerHash []byte) (data.HeaderProofHandler, error)
	GetProofByNonce(headerNonce uint64, shardID uint32) (data.HeaderProofHandler, error)
	HasProof(shardID uint32, headerHash []byte) bool
	IsInterfaceNil() bool
}
