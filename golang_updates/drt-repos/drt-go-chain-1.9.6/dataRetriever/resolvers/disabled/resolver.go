package disabled

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"

	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/p2p"
)

type resolver struct {
}

// NewDisabledResolver returns a new instance of disabled resolver
func NewDisabledResolver() *resolver {
	return &resolver{}
}

// ProcessReceivedMessage returns nil as it is disabled
func (r *resolver) ProcessReceivedMessage(_ p2p.MessageP2P, _ core.PeerID, _ p2p.MessageHandler) ([]byte, error) {
	return []byte{}, nil
}

// SetDebugHandler returns nil as it is disabled
func (r *resolver) SetDebugHandler(_ dataRetriever.DebugHandler) error {
	return nil
}

// SetEpochHandler does nothing and returns nil
func (r *resolver) SetEpochHandler(_ dataRetriever.EpochHandler) error {
	return nil
}

// Close returns nil as it is disabled
func (r *resolver) Close() error {
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (r *resolver) IsInterfaceNil() bool {
	return r == nil
}
