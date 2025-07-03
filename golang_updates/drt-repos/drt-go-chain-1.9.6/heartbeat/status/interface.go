package status

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain/heartbeat/data"
)

// HeartbeatMonitor defines the operations that a monitor should implement
type HeartbeatMonitor interface {
	GetHeartbeats() []data.PubKeyHeartbeat
	IsInterfaceNil() bool
}

// HeartbeatSenderInfoProvider is able to provide correct information about the current sender
type HeartbeatSenderInfoProvider interface {
	GetCurrentNodeType() (string, core.P2PPeerSubType, error)
	IsInterfaceNil() bool
}
