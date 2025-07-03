package factory

import (
	"crypto"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-proxy/common"
	"github.com/TerraDharitri/drt-go-chain-proxy/data"
	"github.com/TerraDharitri/drt-go-chain-proxy/observer"
)

// Processor defines what a processor should be able to do
type Processor interface {
	ComputeShardId(addressBuff []byte) (uint32, error)
	CallGetRestEndPoint(address string, path string, value interface{}) (int, error)
	CallPostRestEndPoint(address string, path string, data interface{}, response interface{}) (int, error)
	GetObserversOnePerShard(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error)
	GetShardIDs() []uint32
	GetFullHistoryNodesOnePerShard(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error)
	GetObservers(shardID uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error)
	GetAllObservers(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error)
	GetFullHistoryNodes(shardID uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error)
	GetAllFullHistoryNodes(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error)
	GetShardCoordinator() common.Coordinator
	GetPubKeyConverter() core.PubkeyConverter
	GetObserverProvider() observer.NodesProviderHandler
	GetFullHistoryNodesProvider() observer.NodesProviderHandler
	IsInterfaceNil() bool
}

// PrivateKeysLoaderHandler defines what a component which handles loading of the private keys file should do
type PrivateKeysLoaderHandler interface {
	PrivateKeysByShard() (map[uint32][]crypto.PrivateKey, error)
}
