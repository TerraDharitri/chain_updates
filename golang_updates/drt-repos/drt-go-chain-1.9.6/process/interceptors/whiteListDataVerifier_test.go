package interceptors

import (
	"bytes"
	"errors"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core/check"

	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/cache"

	"github.com/stretchr/testify/assert"
)

func TestNewWhiteListDataVerifier_NilCacherShouldErr(t *testing.T) {
	t.Parallel()

	wldv, err := NewWhiteListDataVerifier(nil)

	assert.True(t, check.IfNil(wldv))
	assert.True(t, errors.Is(err, process.ErrNilCacher))
}

func TestNewWhiteListDataVerifier_ShouldWork(t *testing.T) {
	t.Parallel()

	wldv, err := NewWhiteListDataVerifier(cache.NewCacherStub())

	assert.False(t, check.IfNil(wldv))
	assert.Nil(t, err)
}

func TestWhiteListDataVerifier_Add(t *testing.T) {
	t.Parallel()

	keys := [][]byte{[]byte("key1"), []byte("key2")}
	added := map[string]struct{}{}
	cacher := &cache.CacherStub{
		PutCalled: func(key []byte, value interface{}, sizeInBytes int) (evicted bool) {
			added[string(key)] = struct{}{}
			return false
		},
	}
	wldv, _ := NewWhiteListDataVerifier(cacher)

	wldv.Add(keys)

	for _, key := range keys {
		_, ok := added[string(key)]
		assert.True(t, ok)
	}
}

func TestWhiteListDataVerifier_Remove(t *testing.T) {
	t.Parallel()

	keys := [][]byte{[]byte("key1"), []byte("key2")}
	removed := map[string]struct{}{}
	cacher := &cache.CacherStub{
		RemoveCalled: func(key []byte) {
			removed[string(key)] = struct{}{}
		},
	}
	wldv, _ := NewWhiteListDataVerifier(cacher)

	wldv.Remove(keys)

	for _, key := range keys {
		_, ok := removed[string(key)]
		assert.True(t, ok)
	}
}

func TestWhiteListDataVerifier_IsWhiteListedNilInterceptedDataShouldRetFalse(t *testing.T) {
	t.Parallel()

	wldv, _ := NewWhiteListDataVerifier(cache.NewCacherStub())

	assert.False(t, wldv.IsWhiteListed(nil))
}

func TestWhiteListDataVerifier_IsWhiteListedNotFoundShouldRetFalse(t *testing.T) {
	t.Parallel()

	keyCheck := []byte("key")
	wldv, _ := NewWhiteListDataVerifier(
		&cache.CacherStub{
			HasCalled: func(key []byte) bool {
				return !bytes.Equal(key, keyCheck)
			},
		},
	)

	ids := &testscommon.InterceptedDataStub{
		HashCalled: func() []byte {
			return keyCheck
		},
	}

	assert.False(t, wldv.IsWhiteListed(ids))
}

func TestWhiteListDataVerifier_IsWhiteListedFoundShouldRetTrue(t *testing.T) {
	t.Parallel()

	keyCheck := []byte("key")
	wldv, _ := NewWhiteListDataVerifier(
		&cache.CacherStub{
			HasCalled: func(key []byte) bool {
				return bytes.Equal(key, keyCheck)
			},
		},
	)

	ids := &testscommon.InterceptedDataStub{
		IdentifiersCalled: func() [][]byte {
			return [][]byte{keyCheck}
		},
	}

	assert.True(t, wldv.IsWhiteListed(ids))
}
