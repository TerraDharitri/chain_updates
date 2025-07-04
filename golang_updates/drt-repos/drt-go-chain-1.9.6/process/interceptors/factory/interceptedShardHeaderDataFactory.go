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

var _ process.InterceptedDataFactory = (*interceptedShardHeaderDataFactory)(nil)

type interceptedShardHeaderDataFactory struct {
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

// NewInterceptedShardHeaderDataFactory creates an instance of interceptedShardHeaderDataFactory
func NewInterceptedShardHeaderDataFactory(argument *ArgInterceptedDataFactory) (*interceptedShardHeaderDataFactory, error) {
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

	return &interceptedShardHeaderDataFactory{
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
func (ishdf *interceptedShardHeaderDataFactory) Create(buff []byte, _ core.PeerID) (process.InterceptedData, error) {
	arg := &interceptedBlocks.ArgInterceptedBlockHeader{
		HdrBuff:                       buff,
		Marshalizer:                   ishdf.marshalizer,
		Hasher:                        ishdf.hasher,
		ShardCoordinator:              ishdf.shardCoordinator,
		HeaderSigVerifier:             ishdf.headerSigVerifier,
		HeaderIntegrityVerifier:       ishdf.headerIntegrityVerifier,
		ValidityAttester:              ishdf.validityAttester,
		EpochStartTrigger:             ishdf.epochStartTrigger,
		EnableEpochsHandler:           ishdf.enableEpochsHandler,
		EpochChangeGracePeriodHandler: ishdf.epochChangeGracePeriodHandler,
	}

	return interceptedBlocks.NewInterceptedHeader(arg)
}

// IsInterfaceNil returns true if there is no value under the interface
func (ishdf *interceptedShardHeaderDataFactory) IsInterfaceNil() bool {
	return ishdf == nil
}
