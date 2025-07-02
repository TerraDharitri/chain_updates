package process

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	chainData "github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain/api/shared"
	"github.com/TerraDharitri/drt-go-chain/consensus"
	"github.com/TerraDharitri/drt-go-chain/factory"
	"github.com/TerraDharitri/drt-go-chain/node/chainSimulator/dtos"
	"github.com/TerraDharitri/drt-go-chain/sharding"
)

// NodeHandler defines what a node handler should be able to do
type NodeHandler interface {
	GetProcessComponents() factory.ProcessComponentsHolder
	GetChainHandler() chainData.ChainHandler
	GetBroadcastMessenger() consensus.BroadcastMessenger
	GetShardCoordinator() sharding.Coordinator
	GetCryptoComponents() factory.CryptoComponentsHolder
	GetCoreComponents() factory.CoreComponentsHolder
	GetDataComponents() factory.DataComponentsHolder
	GetStateComponents() factory.StateComponentsHolder
	GetFacadeHandler() shared.FacadeHandler
	GetStatusCoreComponents() factory.StatusCoreComponentsHolder
	GetNetworkComponents() factory.NetworkComponentsHolder
	SetKeyValueForAddress(addressBytes []byte, state map[string]string) error
	SetStateForAddress(address []byte, state *dtos.AddressState) error
	RemoveAccount(address []byte) error
	ForceChangeOfEpoch() error
	GetBasePeers() map[uint32]core.PeerID
	SetBasePeers(basePeers map[uint32]core.PeerID)
	Close() error
	IsInterfaceNil() bool
}
