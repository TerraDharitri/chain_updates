package factory

import (
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/heartbeat"
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
)

type interceptedHeartbeatDataFactory struct {
	marshalizer marshal.Marshalizer
	peerID      core.PeerID
}

// NewInterceptedHeartbeatDataFactory creates an instance of interceptedHeartbeatDataFactory
func NewInterceptedHeartbeatDataFactory(arg ArgInterceptedDataFactory) (*interceptedHeartbeatDataFactory, error) {
	if check.IfNil(arg.CoreComponents.InternalMarshalizer()) {
		return nil, process.ErrNilMarshalizer
	}
	if len(arg.PeerID) == 0 {
		return nil, process.ErrEmptyPeerID
	}

	return &interceptedHeartbeatDataFactory{
		marshalizer: arg.CoreComponents.InternalMarshalizer(),
		peerID:      arg.PeerID,
	}, nil
}

// Create creates instances of InterceptedData by unmarshalling provided buffer
func (ihdf *interceptedHeartbeatDataFactory) Create(buff []byte, _ core.PeerID) (process.InterceptedData, error) {
	arg := heartbeat.ArgBaseInterceptedHeartbeat{
		DataBuff:   buff,
		Marshaller: ihdf.marshalizer,
	}

	return heartbeat.NewInterceptedHeartbeat(arg)
}

// IsInterfaceNil returns true if there is no value under the interface
func (ihdf *interceptedHeartbeatDataFactory) IsInterfaceNil() bool {
	return ihdf == nil
}
