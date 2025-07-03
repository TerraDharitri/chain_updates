package processing_test

import (
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain/factory/processing"
	"github.com/TerraDharitri/drt-go-chain/process/mock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/components"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManagedProcessComponents_createAPITransactionEvaluator(t *testing.T) {
	t.Parallel()

	shardCoordinatorForShardID2 := mock.NewMultiShardsCoordinatorMock(3)
	shardCoordinatorForShardID2.CurrentShard = 2

	shardCoordinatorForMetachain := mock.NewMultiShardsCoordinatorMock(3)
	shardCoordinatorForMetachain.CurrentShard = core.MetachainShardId

	// no further t.Parallel as these tests are quite heavy (they open netMessengers and other components that start a lot of goroutines)
	t.Run("invalid VMOutputCacher config should error", func(t *testing.T) {
		processArgs := components.GetProcessComponentsFactoryArgs(shardCoordinatorForShardID2)
		processArgs.Config.VMOutputCacher.Type = "invalid"
		pcf, err := processing.NewProcessComponentsFactory(processArgs)
		require.Nil(t, err)

		apiTransactionEvaluator, vmContainerFactory, err := pcf.CreateAPITransactionEvaluator()
		require.NotNil(t, err)
		assert.True(t, check.IfNil(apiTransactionEvaluator))
		assert.True(t, check.IfNil(vmContainerFactory))
		assert.Contains(t, err.Error(), "not supported cache type")
	})
	t.Run("should work for shard", func(t *testing.T) {
		processArgs := components.GetProcessComponentsFactoryArgs(shardCoordinatorForShardID2)
		pcf, err := processing.NewProcessComponentsFactory(processArgs)
		require.Nil(t, err)

		apiTransactionEvaluator, vmContainerFactory, err := pcf.CreateAPITransactionEvaluator()
		require.Nil(t, err)
		assert.False(t, check.IfNil(apiTransactionEvaluator))
		assert.False(t, check.IfNil(vmContainerFactory))
		require.NoError(t, vmContainerFactory.Close())
	})
	t.Run("should work for metachain", func(t *testing.T) {
		processArgs := components.GetProcessComponentsFactoryArgs(shardCoordinatorForMetachain)
		pcf, err := processing.NewProcessComponentsFactory(processArgs)
		require.Nil(t, err)

		apiTransactionEvaluator, vmContainerFactory, err := pcf.CreateAPITransactionEvaluator()
		require.Nil(t, err)
		assert.False(t, check.IfNil(apiTransactionEvaluator))
		assert.False(t, check.IfNil(vmContainerFactory))
		require.NoError(t, vmContainerFactory.Close())
	})
}
