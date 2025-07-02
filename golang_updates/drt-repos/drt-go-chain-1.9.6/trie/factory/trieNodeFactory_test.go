package factory

import (
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain/trie"
	"github.com/stretchr/testify/assert"
)

func TestNewTrieNodeFactory(t *testing.T) {
	t.Parallel()

	tnf := NewTrieNodeFactory()
	assert.False(t, check.IfNil(tnf))
}

func TestTrieNodeFactory_CreateEmpty(t *testing.T) {
	t.Parallel()

	tnf := NewTrieNodeFactory()

	emptyInterceptedNode := tnf.CreateEmpty()
	n, ok := emptyInterceptedNode.(*trie.InterceptedTrieNode)
	assert.True(t, ok)
	assert.NotNil(t, n)
}
