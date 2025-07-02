package mock

import (
	"github.com/TerraDharitri/drt-go-chain/outport"
)

// StatusComponentsMock -
type StatusComponentsMock struct {
	Outport outport.OutportHandler
}

// OutportHandler -
func (scm *StatusComponentsMock) OutportHandler() outport.OutportHandler {
	return scm.Outport
}

// IsInterfaceNil -
func (scm *StatusComponentsMock) IsInterfaceNil() bool {
	return scm == nil
}
