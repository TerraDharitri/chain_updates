package round_test

import (
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/consensus/round"
	consensusMocks "github.com/TerraDharitri/drt-go-chain/testscommon/consensus"

	"github.com/stretchr/testify/assert"
)

const roundTimeDuration = 10 * time.Millisecond

func TestRound_NewRoundShouldErrNilSyncTimer(t *testing.T) {
	t.Parallel()

	genesisTime := time.Now()

	rnd, err := round.NewRound(genesisTime, genesisTime, roundTimeDuration, nil, 0)

	assert.Nil(t, rnd)
	assert.Equal(t, round.ErrNilSyncTimer, err)
}

func TestRound_NewRoundShouldWork(t *testing.T) {
	t.Parallel()

	genesisTime := time.Now()

	syncTimerMock := &consensusMocks.SyncTimerMock{}

	rnd, err := round.NewRound(genesisTime, genesisTime, roundTimeDuration, syncTimerMock, 0)

	assert.Nil(t, err)
	assert.False(t, check.IfNil(rnd))
}

func TestRound_UpdateRoundShouldNotChangeAnything(t *testing.T) {
	t.Parallel()

	genesisTime := time.Now()

	syncTimerMock := &consensusMocks.SyncTimerMock{}

	rnd, _ := round.NewRound(genesisTime, genesisTime, roundTimeDuration, syncTimerMock, 0)
	oldIndex := rnd.Index()
	oldTimeStamp := rnd.TimeStamp()

	rnd.UpdateRound(genesisTime, genesisTime)

	newIndex := rnd.Index()
	newTimeStamp := rnd.TimeStamp()

	assert.Equal(t, oldIndex, newIndex)
	assert.Equal(t, oldTimeStamp, newTimeStamp)
}

func TestRound_UpdateRoundShouldAdvanceOneRound(t *testing.T) {
	t.Parallel()

	genesisTime := time.Now()

	syncTimerMock := &consensusMocks.SyncTimerMock{}

	rnd, _ := round.NewRound(genesisTime, genesisTime, roundTimeDuration, syncTimerMock, 0)
	oldIndex := rnd.Index()
	rnd.UpdateRound(genesisTime, genesisTime.Add(roundTimeDuration))
	newIndex := rnd.Index()

	assert.Equal(t, oldIndex, newIndex-1)
}

func TestRound_IndexShouldReturnFirstIndex(t *testing.T) {
	t.Parallel()

	genesisTime := time.Now()

	syncTimerMock := &consensusMocks.SyncTimerMock{}

	rnd, _ := round.NewRound(genesisTime, genesisTime, roundTimeDuration, syncTimerMock, 0)
	rnd.UpdateRound(genesisTime, genesisTime.Add(roundTimeDuration/2))
	index := rnd.Index()

	assert.Equal(t, int64(0), index)
}

func TestRound_TimeStampShouldReturnTimeStampOfTheNextRound(t *testing.T) {
	t.Parallel()

	genesisTime := time.Now()

	syncTimerMock := &consensusMocks.SyncTimerMock{}

	rnd, _ := round.NewRound(genesisTime, genesisTime, roundTimeDuration, syncTimerMock, 0)
	rnd.UpdateRound(genesisTime, genesisTime.Add(roundTimeDuration+roundTimeDuration/2))
	timeStamp := rnd.TimeStamp()

	assert.Equal(t, genesisTime.Add(roundTimeDuration), timeStamp)
}

func TestRound_TimeDurationShouldReturnTheDurationOfOneRound(t *testing.T) {
	t.Parallel()

	genesisTime := time.Now()

	syncTimerMock := &consensusMocks.SyncTimerMock{}

	rnd, _ := round.NewRound(genesisTime, genesisTime, roundTimeDuration, syncTimerMock, 0)
	timeDuration := rnd.TimeDuration()

	assert.Equal(t, roundTimeDuration, timeDuration)
}

func TestRound_RemainingTimeInCurrentRoundShouldReturnPositiveValue(t *testing.T) {
	t.Parallel()

	genesisTime := time.Unix(0, 0)

	syncTimerMock := &consensusMocks.SyncTimerMock{}

	timeElapsed := int64(roundTimeDuration - 1)

	syncTimerMock.CurrentTimeCalled = func() time.Time {
		return time.Unix(0, timeElapsed)
	}

	rnd, _ := round.NewRound(genesisTime, genesisTime, roundTimeDuration, syncTimerMock, 0)

	remainingTime := rnd.RemainingTime(rnd.TimeStamp(), roundTimeDuration)

	assert.Equal(t, time.Duration(int64(rnd.TimeDuration())-timeElapsed), remainingTime)
	assert.True(t, remainingTime > 0)
}

func TestRound_RemainingTimeInCurrentRoundShouldReturnNegativeValue(t *testing.T) {
	t.Parallel()

	genesisTime := time.Unix(0, 0)

	syncTimerMock := &consensusMocks.SyncTimerMock{}

	timeElapsed := int64(roundTimeDuration + 1)

	syncTimerMock.CurrentTimeCalled = func() time.Time {
		return time.Unix(0, timeElapsed)
	}

	rnd, _ := round.NewRound(genesisTime, genesisTime, roundTimeDuration, syncTimerMock, 0)

	remainingTime := rnd.RemainingTime(rnd.TimeStamp(), roundTimeDuration)

	assert.Equal(t, time.Duration(int64(rnd.TimeDuration())-timeElapsed), remainingTime)
	assert.True(t, remainingTime < 0)
}

func TestRound_RevertOneRound(t *testing.T) {
	t.Parallel()

	genesisTime := time.Now()

	syncTimerMock := &consensusMocks.SyncTimerMock{}

	startRound := int64(10)
	rnd, _ := round.NewRound(genesisTime, genesisTime, roundTimeDuration, syncTimerMock, startRound)
	index := rnd.Index()
	require.Equal(t, startRound, index)

	rnd.RevertOneRound()
	index = rnd.Index()
	require.Equal(t, startRound-1, index)
}

func TestRound_BeforeGenesis(t *testing.T) {
	t.Parallel()

	genesisTime := time.Now()

	syncTimerMock := &consensusMocks.SyncTimerMock{}

	startRound := int64(-1)
	rnd, _ := round.NewRound(genesisTime, genesisTime, roundTimeDuration, syncTimerMock, startRound)
	require.True(t, rnd.BeforeGenesis())

	time.Sleep(roundTimeDuration * 2)
	currentTime := time.Now()

	rnd.UpdateRound(genesisTime, currentTime)
	require.False(t, rnd.BeforeGenesis())
}
