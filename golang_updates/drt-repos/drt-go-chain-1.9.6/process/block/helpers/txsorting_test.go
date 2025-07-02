package helpers

import (
	"testing"

	"github.com/TerraDharitri/drt-go-chain/testscommon/enableEpochsHandlerMock"
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/stretchr/testify/require"
)

func TestComputeRandomnessForTxSorting(t *testing.T) {
	t.Parallel()

	header := &block.Header{
		RandSeed:     []byte{0x01},
		PrevRandSeed: []byte{0x02},
	}

	t.Run("flag not active should return previous randomness", func(t *testing.T) {
		t.Parallel()

		enableEpochsHandler := &enableEpochsHandlerMock.EnableEpochsHandlerStub{
			IsFlagEnabledCalled: func(flag core.EnableEpochFlag) bool {
				return false
			},
		}
		require.Equal(t, header.PrevRandSeed, ComputeRandomnessForTxSorting(header, enableEpochsHandler))
	})
	t.Run("flag active should return current randomness", func(t *testing.T) {
		t.Parallel()

		enableEpochsHandler := &enableEpochsHandlerMock.EnableEpochsHandlerStub{
			IsFlagEnabledCalled: func(flag core.EnableEpochFlag) bool {
				return true
			},
		}
		require.Equal(t, header.RandSeed, ComputeRandomnessForTxSorting(header, enableEpochsHandler))
	})
}
