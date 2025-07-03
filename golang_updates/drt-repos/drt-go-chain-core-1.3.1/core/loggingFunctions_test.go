package core_test

import (
	"fmt"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/mock"
	"github.com/stretchr/testify/require"
)

func TestDumpGoRoutinesToLogShouldNotPanic(t *testing.T) {
	t.Parallel()

	defer func() {
		r := recover()
		if r != nil {
			require.Fail(t, fmt.Sprintf("should have not paniced %v", r))
		}
	}()

	core.DumpGoRoutinesToLog(0, &mock.LoggerMock{})
}

func TestGetRunningGoRoutines(t *testing.T) {
	t.Parallel()

	res := core.GetRunningGoRoutines(&mock.LoggerMock{})
	require.NotNil(t, res)
}
