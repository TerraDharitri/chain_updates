package mock

import (
	"github.com/TerraDharitri/drt-go-chain-proxy/common"
	"github.com/TerraDharitri/drt-go-chain-proxy/data"
)

// BlocksProcessorStub -
type BlocksProcessorStub struct {
	GetBlocksByRoundCalled func(round uint64, options common.BlockQueryOptions) (*data.BlocksApiResponse, error)
}

// GetBlocksByRound -
func (bps *BlocksProcessorStub) GetBlocksByRound(round uint64, options common.BlockQueryOptions) (*data.BlocksApiResponse, error) {
	if bps.GetBlocksByRoundCalled != nil {
		return bps.GetBlocksByRoundCalled(round, options)
	}
	return nil, nil
}
