package testscommon

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain/process"
)

// RatingsInfoMock -
type RatingsInfoMock struct {
	StartRatingProperty           uint32
	MaxRatingProperty             uint32
	MinRatingProperty             uint32
	SignedBlocksThresholdProperty float32
	MetaRatingsStepDataProperty   process.RatingsStepHandler
	ShardRatingsStepDataProperty  process.RatingsStepHandler
	SelectionChancesProperty      []process.SelectionChance
	SetStatusHandlerCalled        func(handler core.AppStatusHandler) error
}

// StartRating -
func (rd *RatingsInfoMock) StartRating() uint32 {
	return rd.StartRatingProperty
}

// MaxRating -
func (rd *RatingsInfoMock) MaxRating() uint32 {
	return rd.MaxRatingProperty
}

// MinRating -
func (rd *RatingsInfoMock) MinRating() uint32 {
	return rd.MinRatingProperty
}

// SignedBlocksThreshold -
func (rd *RatingsInfoMock) SignedBlocksThreshold() float32 {
	return rd.SignedBlocksThresholdProperty
}

// SelectionChances -
func (rd *RatingsInfoMock) SelectionChances() []process.SelectionChance {
	return rd.SelectionChancesProperty
}

// MetaChainRatingsStepHandler -
func (rd *RatingsInfoMock) MetaChainRatingsStepHandler() process.RatingsStepHandler {
	return rd.MetaRatingsStepDataProperty
}

// ShardChainRatingsStepHandler -
func (rd *RatingsInfoMock) ShardChainRatingsStepHandler() process.RatingsStepHandler {
	return rd.ShardRatingsStepDataProperty
}

// SetStatusHandler -
func (rd *RatingsInfoMock) SetStatusHandler(handler core.AppStatusHandler) error {
	if rd.SetStatusHandlerCalled != nil {
		return rd.SetStatusHandlerCalled(handler)
	}
	return nil
}

// IsInterfaceNil -
func (rd *RatingsInfoMock) IsInterfaceNil() bool {
	return rd == nil
}
