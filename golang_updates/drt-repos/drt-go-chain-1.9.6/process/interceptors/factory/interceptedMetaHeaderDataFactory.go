package factory

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/block/interceptedBlocks"
	"github.com/TerraDharitri/drt-go-chain/sharding"
)

var _ process.InterceptedDataFactory = (*interceptedMetaHeaderDataFactory)(nil)

// ArgInterceptedMetaHeaderFactory is the DTO used to create a new instance of meta header factory
type ArgInterceptedMetaHeaderFactory struct {
	ArgInterceptedDataFactory
}

type interceptedMetaHeaderDataFactory struct {
	marshalizer                   marshal.Marshalizer
	hasher                        hashing.Hasher
	shardCoordinator              sharding.Coordinator
	headerSigVerifier             process.InterceptedHeaderSigVerifier
	headerIntegrityVerifier       process.HeaderIntegrityVerifier
	validityAttester              process.ValidityAttester
	epochStartTrigger             process.EpochStartTriggerHandler
	enableEpochsHandler           common.EnableEpochsHandler
	epochChangeGracePeriodHandler common.EpochChangeGracePeriodHandler
}

// NewInterceptedMetaHeaderDataFactory creates an instance of interceptedMetaHeaderDataFactory
func NewInterceptedMetaHeaderDataFactory(argument *ArgInterceptedMetaHeaderFactory) (*interceptedMetaHeaderDataFactory, error) {
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
	if check.IfNil(argument.ShardCoordinator) {
		return nil, process.ErrNilShardCoordinator
	}
	if check.IfNil(argument.HeaderSigVerifier) {
		return nil, process.ErrNilHeaderSigVerifier
	}
	if check.IfNil(argument.HeaderIntegrityVerifier) {
		return nil, process.ErrNilHeaderIntegrityVerifier
	}
	if check.IfNil(argument.EpochStartTrigger) {
		return nil, process.ErrNilEpochStartTrigger
	}
	if len(argument.CoreComponents.ChainID()) == 0 {
		return nil, process.ErrInvalidChainID
	}
	if check.IfNil(argument.ValidityAttester) {
		return nil, process.ErrNilValidityAttester
	}

	return &interceptedMetaHeaderDataFactory{
		marshalizer:                   argument.CoreComponents.InternalMarshalizer(),
		hasher:                        argument.CoreComponents.Hasher(),
		shardCoordinator:              argument.ShardCoordinator,
		headerSigVerifier:             argument.HeaderSigVerifier,
		headerIntegrityVerifier:       argument.HeaderIntegrityVerifier,
		validityAttester:              argument.ValidityAttester,
		epochStartTrigger:             argument.EpochStartTrigger,
		enableEpochsHandler:           argument.CoreComponents.EnableEpochsHandler(),
		epochChangeGracePeriodHandler: argument.CoreComponents.EpochChangeGracePeriodHandler(),
	}, nil
}

// Create creates instances of InterceptedData by unmarshalling provided buffer
func (imhdf *interceptedMetaHeaderDataFactory) Create(buff []byte, _ core.PeerID) (process.InterceptedData, error) {
	arg := &interceptedBlocks.ArgInterceptedBlockHeader{
		HdrBuff:                       buff,
		Marshalizer:                   imhdf.marshalizer,
		Hasher:                        imhdf.hasher,
		ShardCoordinator:              imhdf.shardCoordinator,
		HeaderSigVerifier:             imhdf.headerSigVerifier,
		HeaderIntegrityVerifier:       imhdf.headerIntegrityVerifier,
		ValidityAttester:              imhdf.validityAttester,
		EpochStartTrigger:             imhdf.epochStartTrigger,
		EnableEpochsHandler:           imhdf.enableEpochsHandler,
		EpochChangeGracePeriodHandler: imhdf.epochChangeGracePeriodHandler,
	}

	return interceptedBlocks.NewInterceptedMetaHeader(arg)
}

// IsInterfaceNil returns true if there is no value under the interface
func (imhdf *interceptedMetaHeaderDataFactory) IsInterfaceNil() bool {
	return imhdf == nil
}
