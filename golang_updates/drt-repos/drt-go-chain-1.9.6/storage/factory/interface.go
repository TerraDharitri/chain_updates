package factory

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/block/bootstrapStorage"
	"github.com/TerraDharitri/drt-go-chain/storage"
)

// BootstrapDataProviderHandler defines which actions should be done for loading bootstrap data from the boot storer
type BootstrapDataProviderHandler interface {
	LoadForPath(persisterFactory storage.PersisterFactory, path string) (*bootstrapStorage.BootstrapData, storage.Storer, error)
	GetStorer(storer storage.Storer) (process.BootStorer, error)
	IsInterfaceNil() bool
}

// NodeTypeProviderHandler defines the actions needed for a component that can handle the node type
type NodeTypeProviderHandler interface {
	SetType(nodeType core.NodeType)
	GetType() core.NodeType
	IsInterfaceNil() bool
}
