package mock

import "github.com/TerraDharitri/drt-go-chain-proxy/data"

// DCDTSuppliesProcessorStub -
type DCDTSuppliesProcessorStub struct {
	GetDCDTSupplyCalled func(token string) (*data.DCDTSupplyResponse, error)
}

// GetDCDTSupply -
func (e *DCDTSuppliesProcessorStub) GetDCDTSupply(token string) (*data.DCDTSupplyResponse, error) {
	if e.GetDCDTSupplyCalled != nil {
		return e.GetDCDTSupplyCalled(token)
	}

	return nil, nil
}
