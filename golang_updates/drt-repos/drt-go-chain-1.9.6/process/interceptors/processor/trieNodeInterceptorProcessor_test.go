package processor_test

import (
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core/check"

	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/interceptors/processor"
	"github.com/TerraDharitri/drt-go-chain/process/mock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/cache"

	"github.com/stretchr/testify/assert"
)

func TestNewTrieNodesInterceptorProcessor_NilCacherShouldErr(t *testing.T) {
	t.Parallel()

	tnip, err := processor.NewTrieNodesInterceptorProcessor(nil)
	assert.Nil(t, tnip)
	assert.Equal(t, process.ErrNilCacher, err)
}

func TestNewTrieNodesInterceptorProcessor_OkValsShouldWork(t *testing.T) {
	t.Parallel()

	tnip, err := processor.NewTrieNodesInterceptorProcessor(cache.NewCacherMock())
	assert.Nil(t, err)
	assert.NotNil(t, tnip)
}

// ------- Validate

func TestTrieNodesInterceptorProcessor_ValidateShouldWork(t *testing.T) {
	t.Parallel()

	tnip, _ := processor.NewTrieNodesInterceptorProcessor(cache.NewCacherMock())

	assert.Nil(t, tnip.Validate(nil, ""))
}

// ------- Save

func TestTrieNodesInterceptorProcessor_SaveWrongTypeAssertion(t *testing.T) {
	t.Parallel()

	tnip, _ := processor.NewTrieNodesInterceptorProcessor(cache.NewCacherMock())

	err := tnip.Save(nil, "", "")
	assert.Equal(t, process.ErrWrongTypeAssertion, err)
}

func TestTrieNodesInterceptorProcessor_SaveShouldPutInCacher(t *testing.T) {
	t.Parallel()

	nodeSize := 2242
	nodeHash := []byte("hash")
	interceptedTrieNode := &mock.InterceptedTrieNodeStub{
		HashField: nodeHash,
		SizeInBytesCalled: func() int {
			return nodeSize
		},
	}

	putCalled := false
	cacher := &cache.CacherStub{
		PutCalled: func(key []byte, value interface{}, sizeInBytes int) (evicted bool) {
			putCalled = true
			assert.Equal(t, len(nodeHash)+nodeSize, sizeInBytes)
			return false
		},
	}
	tnip, _ := processor.NewTrieNodesInterceptorProcessor(cacher)

	err := tnip.Save(interceptedTrieNode, "", "")
	assert.Nil(t, err)
	assert.True(t, putCalled)
}

// ------- IsInterfaceNil

func TestTrieNodesInterceptorProcessor_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	var tnip *processor.TrieNodeInterceptorProcessor
	assert.True(t, check.IfNil(tnip))
}
