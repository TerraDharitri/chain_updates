package facade

import (
	"github.com/TerraDharitri/drt-go-chain/ntp"
)

// GetSyncer returns the current syncer
func (nf *nodeFacade) GetSyncer() ntp.SyncTimer {
	return nf.syncer
}
