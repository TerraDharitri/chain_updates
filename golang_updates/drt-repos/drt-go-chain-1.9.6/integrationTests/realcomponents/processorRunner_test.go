package realcomponents

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/testscommon"
)

func TestNewProcessorRunnerAndClose(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	cfg, err := testscommon.CreateTestConfigs(t.TempDir(), "../../cmd/node/config")
	require.Nil(t, err)

	pr := NewProcessorRunner(t, *cfg)
	pr.Close(t)
}
