package logs

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
)

// ArgsNewLogsFacade holds the arguments for constructing a logsFacade
type ArgsNewLogsFacade struct {
	StorageService  dataRetriever.StorageService
	Marshaller      marshal.Marshalizer
	PubKeyConverter core.PubkeyConverter
}

func (args *ArgsNewLogsFacade) check() error {
	if check.IfNil(args.StorageService) {
		return core.ErrNilStore
	}
	if check.IfNil(args.Marshaller) {
		return core.ErrNilMarshalizer
	}
	if check.IfNil(args.PubKeyConverter) {
		return core.ErrNilPubkeyConverter
	}

	return nil
}
