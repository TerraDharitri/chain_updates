package receipts

import (
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
)

// ArgsNewReceiptsRepository holds arguments for creating a receiptsRepository
type ArgsNewReceiptsRepository struct {
	Marshaller marshal.Marshalizer
	Hasher     hashing.Hasher
	Store      dataRetriever.StorageService
}

func (args *ArgsNewReceiptsRepository) check() error {
	if check.IfNil(args.Marshaller) {
		return core.ErrNilMarshalizer
	}
	if check.IfNil(args.Hasher) {
		return core.ErrNilHasher
	}
	if check.IfNil(args.Store) {
		return core.ErrNilStore
	}

	return nil
}
