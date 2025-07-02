package mock

import (
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain-core/core"
)

// InterceptedDataFactoryStub -
type InterceptedDataFactoryStub struct {
	CreateCalled func(buff []byte) (process.InterceptedData, error)
}

// Create -
func (idfs *InterceptedDataFactoryStub) Create(buff []byte, _ core.PeerID) (process.InterceptedData, error) {
	return idfs.CreateCalled(buff)
}

// IsInterfaceNil -
func (idfs *InterceptedDataFactoryStub) IsInterfaceNil() bool {
	return idfs == nil
}
