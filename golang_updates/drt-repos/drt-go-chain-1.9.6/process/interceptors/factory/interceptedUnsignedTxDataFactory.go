package factory

import (
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/unsigned"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
)

var _ process.InterceptedDataFactory = (*interceptedUnsignedTxDataFactory)(nil)

type interceptedUnsignedTxDataFactory struct {
	protoMarshalizer marshal.Marshalizer
	hasher           hashing.Hasher
	pubkeyConverter  core.PubkeyConverter
	shardCoordinator sharding.Coordinator
}

// NewInterceptedUnsignedTxDataFactory creates an instance of interceptedUnsignedTxDataFactory
func NewInterceptedUnsignedTxDataFactory(argument *ArgInterceptedDataFactory) (*interceptedUnsignedTxDataFactory, error) {
	if argument == nil {
		return nil, process.ErrNilArgumentStruct
	}
	if check.IfNil(argument.CoreComponents) {
		return nil, process.ErrNilCoreComponentsHolder
	}
	if check.IfNil(argument.CoreComponents.InternalMarshalizer()) {
		return nil, process.ErrNilMarshalizer
	}
	if check.IfNil(argument.CoreComponents.TxMarshalizer()) {
		return nil, process.ErrNilMarshalizer
	}
	if check.IfNil(argument.CoreComponents.Hasher()) {
		return nil, process.ErrNilHasher
	}
	if check.IfNil(argument.CoreComponents.AddressPubKeyConverter()) {
		return nil, process.ErrNilPubkeyConverter
	}
	if check.IfNil(argument.ShardCoordinator) {
		return nil, process.ErrNilShardCoordinator
	}

	return &interceptedUnsignedTxDataFactory{
		protoMarshalizer: argument.CoreComponents.InternalMarshalizer(),
		hasher:           argument.CoreComponents.Hasher(),
		pubkeyConverter:  argument.CoreComponents.AddressPubKeyConverter(),
		shardCoordinator: argument.ShardCoordinator,
	}, nil
}

// Create creates instances of InterceptedData by unmarshalling provided buffer
func (iutdf *interceptedUnsignedTxDataFactory) Create(buff []byte, _ core.PeerID) (process.InterceptedData, error) {
	return unsigned.NewInterceptedUnsignedTransaction(
		buff,
		iutdf.protoMarshalizer,
		iutdf.hasher,
		iutdf.pubkeyConverter,
		iutdf.shardCoordinator,
	)
}

// IsInterfaceNil returns true if there is no value under the interface
func (iutdf *interceptedUnsignedTxDataFactory) IsInterfaceNil() bool {
	return iutdf == nil
}
