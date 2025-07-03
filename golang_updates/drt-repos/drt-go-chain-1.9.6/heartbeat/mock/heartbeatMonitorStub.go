package mock

import "github.com/TerraDharitri/drt-go-chain/heartbeat/data"

// HeartbeatMonitorStub -
type HeartbeatMonitorStub struct {
	GetHeartbeatsCalled func() []data.PubKeyHeartbeat
}

// GetHeartbeats -
func (stub *HeartbeatMonitorStub) GetHeartbeats() []data.PubKeyHeartbeat {
	if stub.GetHeartbeatsCalled != nil {
		return stub.GetHeartbeatsCalled()
	}

	return nil
}

// IsInterfaceNil -
func (stub *HeartbeatMonitorStub) IsInterfaceNil() bool {
	return stub == nil
}
