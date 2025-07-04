package statistics

import (
	"fmt"
	"sync"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/stretchr/testify/assert"
)

func TestNewStateStatistics_ShouldWork(t *testing.T) {
	t.Parallel()

	ss := NewStateStatistics()

	assert.False(t, check.IfNil(ss))
}

func TestStateStatistics_Processing(t *testing.T) {
	t.Parallel()

	t.Run("trie operations", func(t *testing.T) {
		t.Parallel()

		ss := NewStateStatistics()

		assert.Equal(t, uint64(0), ss.Trie())

		ss.IncrTrie()
		ss.IncrTrie()
		assert.Equal(t, uint64(2), ss.Trie())

		ss.IncrTrie()
		assert.Equal(t, uint64(3), ss.Trie())

		ss.Reset()
		assert.Equal(t, uint64(0), ss.Trie())
	})

	t.Run("persister operations", func(t *testing.T) {
		t.Parallel()

		ss := NewStateStatistics()

		epoch := uint32(1)

		assert.Equal(t, uint64(0), ss.Persister(epoch))

		ss.IncrPersister(epoch)
		ss.IncrPersister(epoch)
		assert.Equal(t, uint64(2), ss.Persister(epoch))

		ss.IncrPersister(epoch)
		assert.Equal(t, uint64(3), ss.Persister(epoch))

		ss.IncrWritePersister(epoch)
		ss.IncrWritePersister(epoch)
		ss.IncrWritePersister(epoch)
		assert.Equal(t, uint64(3), ss.WritePersister(epoch))

		ss.Reset()
		assert.Equal(t, uint64(0), ss.Persister(epoch))
		assert.Equal(t, uint64(0), ss.WritePersister(epoch))
	})

	t.Run("cache operations", func(t *testing.T) {
		t.Parallel()

		ss := NewStateStatistics()

		assert.Equal(t, uint64(0), ss.Cache())

		ss.IncrCache()
		ss.IncrCache()
		assert.Equal(t, uint64(2), ss.Cache())

		ss.IncrCache()
		assert.Equal(t, uint64(3), ss.Cache())

		ss.Reset()
		assert.Equal(t, uint64(0), ss.Cache())
	})
}

func TestStateStatistics_Snapshot(t *testing.T) {
	t.Parallel()

	t.Run("persister operations", func(t *testing.T) {
		t.Parallel()

		ss := NewStateStatistics()

		epoch := uint32(1)

		assert.Equal(t, uint64(0), ss.SnapshotPersister(epoch))

		ss.IncrSnapshotPersister(epoch)
		ss.IncrSnapshotPersister(epoch)
		assert.Equal(t, uint64(2), ss.SnapshotPersister(epoch))

		ss.IncrSnapshotPersister(epoch)
		assert.Equal(t, uint64(3), ss.SnapshotPersister(epoch))

		ss.ResetSnapshot()
		assert.Equal(t, uint64(0), ss.SnapshotPersister(epoch))
	})

	t.Run("cache operations", func(t *testing.T) {
		t.Parallel()

		ss := NewStateStatistics()

		assert.Equal(t, uint64(0), ss.Cache())

		ss.IncrSnapshotCache()
		ss.IncrSnapshotCache()
		assert.Equal(t, uint64(2), ss.SnapshotCache())

		ss.IncrSnapshotCache()
		assert.Equal(t, uint64(3), ss.SnapshotCache())

		ss.ResetSnapshot()
		assert.Equal(t, uint64(0), ss.SnapshotCache())
	})
}

func TestStateStatistics_ConcurrenyOperations(t *testing.T) {
	t.Parallel()

	defer func() {
		r := recover()
		if r != nil {
			assert.Fail(t, fmt.Sprintf("should have not panicked: %v", r))
		}
	}()

	numIterations := 10000

	epoch := uint32(1)

	ss := NewStateStatistics()

	wg := sync.WaitGroup{}
	wg.Add(numIterations)

	for i := 0; i < numIterations; i++ {
		go func(idx int) {
			switch idx % 13 {
			case 0:
				ss.Reset()
			case 1:
				ss.IncrCache()
			case 2:
				ss.IncrPersister(epoch)
			case 3:
				ss.IncrTrie()
			case 7:
				_ = ss.Cache()
			case 8:
				_ = ss.Persister(epoch)
			case 9:
				_ = ss.Trie()
			case 10:
				_ = ss.ProcessingStats()
			case 11:
				ss.IncrWritePersister(epoch)
			case 12:
				_ = ss.WritePersister(epoch)
			}

			wg.Done()
		}(i)
	}

	wg.Wait()
}
