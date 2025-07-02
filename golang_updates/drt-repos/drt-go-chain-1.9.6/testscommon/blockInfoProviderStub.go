package testscommon

import (
	"github.com/TerraDharitri/drt-go-chain/common"
)

// BlockInfoProviderStub -
type BlockInfoProviderStub struct {
	GetBlockInfoCalled func() common.BlockInfo
}

// GetBlockInfo -
func (stub *BlockInfoProviderStub) GetBlockInfo() common.BlockInfo {
	if stub.GetBlockInfoCalled != nil {
		return stub.GetBlockInfoCalled()
	}

	return nil
}

// IsInterfaceNil -
func (stub *BlockInfoProviderStub) IsInterfaceNil() bool {
	return stub == nil
}
