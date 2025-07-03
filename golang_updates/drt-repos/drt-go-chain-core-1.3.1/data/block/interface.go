package block

import "github.com/TerraDharitri/drt-go-chain-core/data"

// EmptyBlockCreator is able to create empty block instances
type EmptyBlockCreator interface {
	CreateNewHeader() data.HeaderHandler
	IsInterfaceNil() bool
}
