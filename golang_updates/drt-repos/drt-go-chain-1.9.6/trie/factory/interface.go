package factory

import (
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/storage"
)

type coreComponentsHandler interface {
	InternalMarshalizer() marshal.Marshalizer
	Hasher() hashing.Hasher
	PathHandler() storage.PathManagerHandler
	ProcessStatusHandler() common.ProcessStatusHandler
	EnableEpochsHandler() common.EnableEpochsHandler
}
