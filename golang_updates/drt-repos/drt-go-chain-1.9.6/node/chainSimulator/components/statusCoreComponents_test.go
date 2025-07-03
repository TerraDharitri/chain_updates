package components

import (
	"testing"

	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/factory"
	"github.com/TerraDharitri/drt-go-chain/factory/mock"
	mockTests "github.com/TerraDharitri/drt-go-chain/integrationTests/mock"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/components"
	"github.com/TerraDharitri/drt-go-chain/testscommon/economicsmocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/genesisMocks"
	"github.com/stretchr/testify/require"
)

func createArgs() (config.Configs, factory.CoreComponentsHolder) {
	generalCfg := testscommon.GetGeneralConfig()
	ratingsCfg := components.CreateDummyRatingsConfig()
	economicsCfg := components.CreateDummyEconomicsConfig()
	cfg := config.Configs{
		GeneralConfig: &generalCfg,
		EpochConfig: &config.EpochConfig{
			GasSchedule: config.GasScheduleConfig{
				GasScheduleByEpochs: []config.GasScheduleByEpochs{
					{
						StartEpoch: 0,
						FileName:   "gasScheduleV1.toml",
					},
				},
			},
		},
		RoundConfig: &config.RoundConfig{
			RoundActivations: map[string]config.ActivationRoundByName{
				"Example": {
					Round: "18446744073709551615",
				},
			},
		},
		RatingsConfig:   &ratingsCfg,
		EconomicsConfig: &economicsCfg,
	}

	return cfg, &mock.CoreComponentsMock{
		EconomicsHandler:    &economicsmocks.EconomicsHandlerMock{},
		IntMarsh:            &testscommon.MarshallerStub{},
		UInt64ByteSliceConv: &mockTests.Uint64ByteSliceConverterMock{},
		NodesConfig:         &genesisMocks.NodesSetupStub{},
		RatingsConfig:       &testscommon.RatingsInfoMock{},
	}
}

func TestCreateStatusCoreComponents(t *testing.T) {
	t.Parallel()

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		cfg, coreComp := createArgs()
		comp, err := CreateStatusCoreComponents(cfg, coreComp)
		require.NoError(t, err)
		require.NotNil(t, comp)

		require.Nil(t, comp.Create())
		require.Nil(t, comp.Close())
	})
	t.Run("NewStatusCoreComponentsFactory failure should error", func(t *testing.T) {
		t.Parallel()

		cfg, _ := createArgs()
		comp, err := CreateStatusCoreComponents(cfg, nil)
		require.Error(t, err)
		require.Nil(t, comp)
	})
	t.Run("managedStatusCoreComponents.Create failure should error", func(t *testing.T) {
		t.Parallel()

		cfg, coreComp := createArgs()
		cfg.GeneralConfig.ResourceStats.RefreshIntervalInSec = 0
		comp, err := CreateStatusCoreComponents(cfg, coreComp)
		require.Error(t, err)
		require.Nil(t, comp)
	})
}

func TestStatusCoreComponentsHolder_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	var comp *statusCoreComponentsHolder
	require.True(t, comp.IsInterfaceNil())

	cfg, coreComp := createArgs()
	comp, _ = CreateStatusCoreComponents(cfg, coreComp)
	require.False(t, comp.IsInterfaceNil())
	require.Nil(t, comp.Close())
}

func TestStatusCoreComponentsHolder_Getters(t *testing.T) {
	t.Parallel()

	cfg, coreComp := createArgs()
	comp, err := CreateStatusCoreComponents(cfg, coreComp)
	require.NoError(t, err)

	require.NotNil(t, comp.ResourceMonitor())
	require.NotNil(t, comp.NetworkStatistics())
	require.NotNil(t, comp.TrieSyncStatistics())
	require.NotNil(t, comp.AppStatusHandler())
	require.NotNil(t, comp.StatusMetrics())
	require.NotNil(t, comp.PersistentStatusHandler())
	require.NotNil(t, comp.StateStatsHandler())
	require.Nil(t, comp.CheckSubcomponents())
	require.Empty(t, comp.String())
}
