package redundancy

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
)

// P2PMessenger defines a subset of the p2p.Messenger interface
type P2PMessenger interface {
	ID() core.PeerID
	IsInterfaceNil() bool
}
