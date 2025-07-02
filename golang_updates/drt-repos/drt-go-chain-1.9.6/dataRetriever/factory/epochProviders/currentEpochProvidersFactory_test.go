package epochProviders

import (
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever/resolvers/epochproviders"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever/resolvers/epochproviders/disabled"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateCurrentEpochProvider_NilCurrentEpochProvider(t *testing.T) {
	t.Parallel()

	cnep, err := CreateCurrentEpochProvider(
		config.Config{},
		0,
		0,
		false,
	)

	assert.Nil(t, err)
	assert.IsType(t, disabled.NewEpochProvider(), cnep)
}

func TestCreateCurrentEpochProvider_ArithmeticEpochProvider(t *testing.T) {
	t.Parallel()

	cnep, err := CreateCurrentEpochProvider(
		config.Config{
			EpochStartConfig: config.EpochStartConfig{
				RoundsPerEpoch: 1,
			},
		},
		1,
		1,
		true,
	)
	require.Nil(t, err)

	aep, _ := epochproviders.NewArithmeticEpochProvider(
		epochproviders.ArgArithmeticEpochProvider{
			RoundsPerEpoch:          1,
			RoundTimeInMilliseconds: 1,
			StartTime:               1,
		},
	)
	require.False(t, check.IfNil(aep))
	assert.IsType(t, aep, cnep)
}
