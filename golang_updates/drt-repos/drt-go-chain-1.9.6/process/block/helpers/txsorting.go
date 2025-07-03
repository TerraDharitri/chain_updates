package helpers

import (
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain-core/data"
)

// ComputeRandomnessForTxSorting returns the randomness for transactions sorting
func ComputeRandomnessForTxSorting(header data.HeaderHandler, enableEpochsHandler common.EnableEpochsHandler) []byte {
	if enableEpochsHandler.IsFlagEnabled(common.CurrentRandomnessOnSortingFlag) {
		return header.GetRandSeed()
	}

	return header.GetPrevRandSeed()
}
