package mock

import (
	"github.com/TerraDharitri/drt-go-chain/p2p"
	"github.com/TerraDharitri/drt-go-chain-core/core"
)

// ResolverStub -
type ResolverStub struct {
	ProcessReceivedMessageCalled func(message p2p.MessageP2P) error
}

// ProcessReceivedMessage -
func (rs *ResolverStub) ProcessReceivedMessage(message p2p.MessageP2P, _ core.PeerID) error {
	return rs.ProcessReceivedMessageCalled(message)
}

// IsInterfaceNil returns true if there is no value under the interface
func (rs *ResolverStub) IsInterfaceNil() bool {
	return rs == nil
}
