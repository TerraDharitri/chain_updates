package chronology_test

import (
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/consensus"
	"github.com/TerraDharitri/drt-go-chain/consensus/chronology"
	"github.com/TerraDharitri/drt-go-chain/consensus/mock"
	consensusMocks "github.com/TerraDharitri/drt-go-chain/testscommon/consensus"
	statusHandlerMock "github.com/TerraDharitri/drt-go-chain/testscommon/statusHandler"
)

func initSubroundHandlerMock() *mock.SubroundHandlerMock {
	srm := &mock.SubroundHandlerMock{}
	srm.CurrentCalled = func() int {
		return 0
	}
	srm.NextCalled = func() int {
		return 1
	}
	srm.DoWorkCalled = func(roundHandler consensus.RoundHandler) bool {
		return false
	}
	srm.NameCalled = func() string {
		return "(TEST)"
	}
	return srm
}

func TestChronology_NewChronologyNilRoundHandlerShouldFail(t *testing.T) {
	t.Parallel()

	arg := getDefaultChronologyArg()
	arg.RoundHandler = nil
	chr, err := chronology.NewChronology(arg)

	assert.Nil(t, chr)
	assert.Equal(t, err, chronology.ErrNilRoundHandler)
}

func TestChronology_NewChronologyNilSyncerShouldFail(t *testing.T) {
	t.Parallel()

	arg := getDefaultChronologyArg()
	arg.SyncTimer = nil
	chr, err := chronology.NewChronology(arg)

	assert.Nil(t, chr)
	assert.Equal(t, err, chronology.ErrNilSyncTimer)
}

func TestChronology_NewChronologyNilWatchdogShouldFail(t *testing.T) {
	t.Parallel()

	arg := getDefaultChronologyArg()
	arg.Watchdog = nil
	chr, err := chronology.NewChronology(arg)

	assert.Nil(t, chr)
	assert.Equal(t, err, chronology.ErrNilWatchdog)
}

func TestChronology_NewChronologyNilAppStatusHandlerShouldFail(t *testing.T) {
	t.Parallel()

	arg := getDefaultChronologyArg()
	arg.AppStatusHandler = nil
	chr, err := chronology.NewChronology(arg)

	assert.Nil(t, chr)
	assert.Equal(t, err, chronology.ErrNilAppStatusHandler)
}

func TestChronology_NewChronologyShouldWork(t *testing.T) {
	t.Parallel()

	arg := getDefaultChronologyArg()
	chr, err := chronology.NewChronology(arg)

	assert.Nil(t, err)
	assert.False(t, check.IfNil(chr))
}

func TestChronology_AddSubroundShouldWork(t *testing.T) {
	t.Parallel()

	arg := getDefaultChronologyArg()
	chr, _ := chronology.NewChronology(arg)

	chr.AddSubround(initSubroundHandlerMock())
	chr.AddSubround(initSubroundHandlerMock())
	chr.AddSubround(initSubroundHandlerMock())

	assert.Equal(t, 3, len(chr.SubroundHandlers()))
}

func TestChronology_RemoveAllSubroundsShouldReturnEmptySubroundHandlersArray(t *testing.T) {
	t.Parallel()

	arg := getDefaultChronologyArg()
	chr, _ := chronology.NewChronology(arg)

	chr.AddSubround(initSubroundHandlerMock())
	chr.AddSubround(initSubroundHandlerMock())
	chr.AddSubround(initSubroundHandlerMock())

	assert.Equal(t, 3, len(chr.SubroundHandlers()))
	chr.RemoveAllSubrounds()
	assert.Equal(t, 0, len(chr.SubroundHandlers()))
}

func TestChronology_StartRoundShouldReturnWhenRoundIndexIsNegative(t *testing.T) {
	t.Parallel()

	arg := getDefaultChronologyArg()
	roundHandlerMock := &consensusMocks.RoundHandlerMock{}
	roundHandlerMock.IndexCalled = func() int64 {
		return -1
	}
	roundHandlerMock.BeforeGenesisCalled = func() bool {
		return true
	}
	arg.RoundHandler = roundHandlerMock
	chr, _ := chronology.NewChronology(arg)

	srm := initSubroundHandlerMock()
	chr.AddSubround(srm)
	chr.SetSubroundId(0)
	chr.StartRound()

	assert.Equal(t, srm.Current(), chr.SubroundId())
}

func TestChronology_StartRoundShouldReturnWhenLoadSubroundHandlerReturnsNil(t *testing.T) {
	t.Parallel()

	arg := getDefaultChronologyArg()
	chr, _ := chronology.NewChronology(arg)

	initSubroundHandlerMock()
	chr.StartRound()

	assert.Equal(t, -1, chr.SubroundId())
}

func TestChronology_StartRoundShouldReturnWhenDoWorkReturnsFalse(t *testing.T) {
	t.Parallel()

	arg := getDefaultChronologyArg()
	roundHandlerMock := &consensusMocks.RoundHandlerMock{}
	roundHandlerMock.UpdateRound(roundHandlerMock.TimeStamp(), roundHandlerMock.TimeStamp().Add(roundHandlerMock.TimeDuration()))
	arg.RoundHandler = roundHandlerMock
	chr, _ := chronology.NewChronology(arg)

	srm := initSubroundHandlerMock()
	chr.AddSubround(srm)
	chr.SetSubroundId(0)
	chr.StartRound()

	assert.Equal(t, -1, chr.SubroundId())
}

func TestChronology_StartRoundShouldWork(t *testing.T) {
	t.Parallel()

	arg := getDefaultChronologyArg()
	roundHandlerMock := &consensusMocks.RoundHandlerMock{}
	roundHandlerMock.UpdateRound(roundHandlerMock.TimeStamp(), roundHandlerMock.TimeStamp().Add(roundHandlerMock.TimeDuration()))
	arg.RoundHandler = roundHandlerMock
	chr, _ := chronology.NewChronology(arg)

	srm := initSubroundHandlerMock()
	srm.DoWorkCalled = func(roundHandler consensus.RoundHandler) bool {
		return true
	}
	chr.AddSubround(srm)
	chr.SetSubroundId(0)
	chr.StartRound()

	assert.Equal(t, srm.Next(), chr.SubroundId())
}

func TestChronology_UpdateRoundShouldInitRound(t *testing.T) {
	t.Parallel()

	arg := getDefaultChronologyArg()
	chr, _ := chronology.NewChronology(arg)

	srm := initSubroundHandlerMock()
	chr.AddSubround(srm)
	chr.UpdateRound()

	assert.Equal(t, srm.Current(), chr.SubroundId())
}

func TestChronology_LoadSubroundHandlerShouldReturnNilWhenSubroundHandlerNotExists(t *testing.T) {
	t.Parallel()

	arg := getDefaultChronologyArg()
	chr, _ := chronology.NewChronology(arg)

	assert.Nil(t, chr.LoadSubroundHandler(0))
}

func TestChronology_LoadSubroundHandlerShouldReturnNilWhenIndexIsOutOfBound(t *testing.T) {
	t.Parallel()
	arg := getDefaultChronologyArg()
	chr, _ := chronology.NewChronology(arg)

	chr.AddSubround(initSubroundHandlerMock())
	chr.SetSubroundHandlers(make([]consensus.SubroundHandler, 0))

	assert.Nil(t, chr.LoadSubroundHandler(0))
}

func TestChronology_InitRoundShouldNotSetSubroundWhenRoundIndexIsNegative(t *testing.T) {
	t.Parallel()

	arg := getDefaultChronologyArg()
	roundHandlerMock := &consensusMocks.RoundHandlerMock{}
	arg.RoundHandler = roundHandlerMock
	arg.GenesisTime = arg.SyncTimer.CurrentTime()
	chr, _ := chronology.NewChronology(arg)

	chr.AddSubround(initSubroundHandlerMock())
	roundHandlerMock.IndexCalled = func() int64 {
		return -1
	}
	roundHandlerMock.BeforeGenesisCalled = func() bool {
		return true
	}
	chr.InitRound()

	assert.Equal(t, -1, chr.SubroundId())
}

func TestChronology_InitRoundShouldSetSubroundWhenRoundIndexIsPositive(t *testing.T) {
	t.Parallel()

	arg := getDefaultChronologyArg()
	roundHandlerMock := &consensusMocks.RoundHandlerMock{}
	roundHandlerMock.UpdateRound(roundHandlerMock.TimeStamp(), roundHandlerMock.TimeStamp().Add(roundHandlerMock.TimeDuration()))
	arg.RoundHandler = roundHandlerMock
	arg.GenesisTime = arg.SyncTimer.CurrentTime()
	chr, _ := chronology.NewChronology(arg)

	sr := initSubroundHandlerMock()
	chr.AddSubround(sr)
	chr.InitRound()

	assert.Equal(t, sr.Current(), chr.SubroundId())
}

func TestChronology_StartRoundShouldNotUpdateRoundWhenCurrentRoundIsNotFinished(t *testing.T) {
	t.Parallel()

	arg := getDefaultChronologyArg()
	roundHandlerMock := &consensusMocks.RoundHandlerMock{}
	arg.RoundHandler = roundHandlerMock
	arg.GenesisTime = arg.SyncTimer.CurrentTime()
	chr, _ := chronology.NewChronology(arg)

	chr.SetSubroundId(0)
	chr.StartRound()

	assert.Equal(t, int64(0), roundHandlerMock.Index())
}

func TestChronology_StartRoundShouldUpdateRoundWhenCurrentRoundIsFinished(t *testing.T) {
	t.Parallel()
	arg := getDefaultChronologyArg()
	roundHandlerMock := &consensusMocks.RoundHandlerMock{}
	arg.RoundHandler = roundHandlerMock
	arg.GenesisTime = arg.SyncTimer.CurrentTime()
	chr, _ := chronology.NewChronology(arg)

	chr.SetSubroundId(-1)
	chr.StartRound()

	assert.Equal(t, int64(1), roundHandlerMock.Index())
}

func TestChronology_CheckIfStatusHandlerWorks(t *testing.T) {
	t.Parallel()

	chanDone := make(chan bool, 2)
	arg := getDefaultChronologyArg()
	arg.GenesisTime = arg.SyncTimer.CurrentTime()
	arg.AppStatusHandler = &statusHandlerMock.AppStatusHandlerStub{
		SetUInt64ValueHandler: func(key string, value uint64) {
			chanDone <- true
		},
	}
	chr, err := chronology.NewChronology(arg)

	assert.Nil(t, err)

	srm := initSubroundHandlerMock()
	srm.DoWorkCalled = func(roundHandler consensus.RoundHandler) bool {
		return true
	}

	chr.AddSubround(srm)
	chr.StartRound()

	select {
	case <-chanDone:
	case <-time.After(1 * time.Second):
		assert.Fail(t, "AppStatusHandler not working")
	}
}

func getDefaultChronologyArg() chronology.ArgChronology {
	return chronology.ArgChronology{
		GenesisTime:      time.Now(),
		RoundHandler:     &consensusMocks.RoundHandlerMock{},
		SyncTimer:        &consensusMocks.SyncTimerMock{},
		AppStatusHandler: statusHandlerMock.NewAppStatusHandlerMock(),
		Watchdog:         &mock.WatchdogMock{},
	}
}

func TestChronology_CloseWatchDogStop(t *testing.T) {
	t.Parallel()

	arg := getDefaultChronologyArg()
	stopCalled := false
	arg.Watchdog = &mock.WatchdogMock{
		StopCalled: func(alarmID string) {
			stopCalled = true
		},
	}

	chr, err := chronology.NewChronology(arg)
	require.Nil(t, err)
	chr.SetCancelFunc(nil)

	err = chr.Close()
	assert.Nil(t, err)
	assert.True(t, stopCalled)
}

func TestChronology_Close(t *testing.T) {
	t.Parallel()

	arg := getDefaultChronologyArg()
	stopCalled := false
	arg.Watchdog = &mock.WatchdogMock{
		StopCalled: func(alarmID string) {
			stopCalled = true
		},
	}

	chr, err := chronology.NewChronology(arg)
	require.Nil(t, err)

	cancelCalled := false
	chr.SetCancelFunc(func() {
		cancelCalled = true
	})

	err = chr.Close()
	assert.Nil(t, err)
	assert.True(t, stopCalled)
	assert.True(t, cancelCalled)
}

func TestChronology_StartRounds(t *testing.T) {
	t.Parallel()

	arg := getDefaultChronologyArg()

	chr, err := chronology.NewChronology(arg)
	require.Nil(t, err)
	doneFuncCalled := false

	ctx := &mock.ContextMock{
		DoneFunc: func() <-chan struct{} {
			done := make(chan struct{})
			close(done)
			doneFuncCalled = true
			return done
		},
	}
	chr.StartRoundsTest(ctx)
	assert.True(t, doneFuncCalled)
}

func TestChronology_StartRoundsShouldWork(t *testing.T) {
	t.Parallel()

	arg := getDefaultChronologyArg()
	roundHandlerMock := &consensusMocks.RoundHandlerMock{}
	roundHandlerMock.UpdateRound(roundHandlerMock.TimeStamp(), roundHandlerMock.TimeStamp().Add(roundHandlerMock.TimeDuration()))
	arg.RoundHandler = roundHandlerMock
	chr, _ := chronology.NewChronology(arg)

	srm := initSubroundHandlerMock()
	srm.DoWorkCalled = func(roundHandler consensus.RoundHandler) bool {
		return true
	}
	chr.AddSubround(srm)
	chr.SetSubroundId(1)
	chr.StartRounds()
	defer chr.Close()

	assert.Equal(t, srm.Next(), chr.SubroundId())
	time.Sleep(time.Millisecond * 10)
}
