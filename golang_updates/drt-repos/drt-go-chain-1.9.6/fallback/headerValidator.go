package fallback

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	logger "github.com/TerraDharitri/drt-go-chain-logger"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/process"
)

var log = logger.GetOrCreate("fallback")

type fallbackHeaderValidator struct {
	headersPool    dataRetriever.HeadersPool
	marshalizer    marshal.Marshalizer
	storageService dataRetriever.StorageService
}

// NewFallbackHeaderValidator creates a new fallbackHeaderValidator object
func NewFallbackHeaderValidator(
	headersPool dataRetriever.HeadersPool,
	marshalizer marshal.Marshalizer,
	storageService dataRetriever.StorageService,
) (*fallbackHeaderValidator, error) {

	if check.IfNil(headersPool) {
		return nil, process.ErrNilHeadersDataPool
	}
	if check.IfNil(marshalizer) {
		return nil, process.ErrNilMarshalizer
	}
	if check.IfNil(storageService) {
		return nil, process.ErrNilStorage
	}

	hv := &fallbackHeaderValidator{
		headersPool:    headersPool,
		marshalizer:    marshalizer,
		storageService: storageService,
	}

	return hv, nil
}

// ShouldApplyFallbackValidationForHeaderWith returns if for the given header data fallback validation could be applied or not
func (fhv *fallbackHeaderValidator) ShouldApplyFallbackValidationForHeaderWith(shardID uint32, startOfEpochBlock bool, round uint64, prevHeaderHash []byte) bool {
	if shardID != core.MetachainShardId {
		return false
	}
	if !startOfEpochBlock {
		return false
	}

	previousHeader, err := process.GetMetaHeader(prevHeaderHash, fhv.headersPool, fhv.marshalizer, fhv.storageService)
	if err != nil {
		log.Debug("ShouldApplyFallbackValidation", "GetMetaHeader", err.Error())
		return false
	}

	isRoundTooOld := int64(round)-int64(previousHeader.GetRound()) >= common.MaxRoundsWithoutCommittedStartInEpochBlock
	return isRoundTooOld
}

// ShouldApplyFallbackValidation returns if for the given header could be applied fallback validation or not
func (fhv *fallbackHeaderValidator) ShouldApplyFallbackValidation(headerHandler data.HeaderHandler) bool {
	if check.IfNil(headerHandler) {
		return false
	}

	return fhv.ShouldApplyFallbackValidationForHeaderWith(headerHandler.GetShardID(), headerHandler.IsStartOfEpochBlock(), headerHandler.GetRound(), headerHandler.GetPrevHash())
}

// IsInterfaceNil returns true if there is no value under the interface
func (fhv *fallbackHeaderValidator) IsInterfaceNil() bool {
	return fhv == nil
}
