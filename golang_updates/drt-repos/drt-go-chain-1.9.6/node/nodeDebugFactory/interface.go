package nodeDebugFactory

import "github.com/TerraDharitri/drt-go-chain/debug"

// NodeWrapper is the interface that defines the behavior of a Node that can work with debug handlers
type NodeWrapper interface {
	AddQueryHandler(name string, handler debug.QueryHandler) error
	IsInterfaceNil() bool
}
