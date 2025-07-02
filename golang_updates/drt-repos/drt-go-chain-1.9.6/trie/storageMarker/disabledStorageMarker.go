package storageMarker

import "github.com/TerraDharitri/drt-go-chain/common"

type disabledStorageMarker struct {
}

// NewDisabledStorageMarker creates a new instance of disabledStorageMarker
func NewDisabledStorageMarker() *disabledStorageMarker {
	return &disabledStorageMarker{}
}

// MarkStorerAsSyncedAndActive does nothing for this implementation
func (dsm *disabledStorageMarker) MarkStorerAsSyncedAndActive(_ common.StorageManager) {
}

// IsInterfaceNil returns true if there is no value under the interface
func (dsm *disabledStorageMarker) IsInterfaceNil() bool {
	return dsm == nil
}
