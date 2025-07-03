package mock

import (
	"github.com/TerraDharitri/drt-go-chain-core/data/vm"
	"github.com/TerraDharitri/drt-go-chain-proxy/data"
)

// SCQueryServiceStub -
type SCQueryServiceStub struct {
	ExecuteQueryCalled func(*data.SCQuery) (*vm.VMOutputApi, data.BlockInfo, error)
}

// ExecuteQuery -
func (serviceStub *SCQueryServiceStub) ExecuteQuery(query *data.SCQuery) (*vm.VMOutputApi, data.BlockInfo, error) {
	return serviceStub.ExecuteQueryCalled(query)
}
