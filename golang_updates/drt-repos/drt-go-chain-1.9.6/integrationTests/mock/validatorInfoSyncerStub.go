package mock

import (
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
)

// ValidatorInfoSyncerStub -
type ValidatorInfoSyncerStub struct {
}

// SyncMiniBlocks -
func (vip *ValidatorInfoSyncerStub) SyncMiniBlocks(*block.MetaBlock) ([][]byte, data.BodyHandler, error) {
	return nil, nil, nil
}

// IsInterfaceNil -
func (vip *ValidatorInfoSyncerStub) IsInterfaceNil() bool {
	return vip == nil
}
