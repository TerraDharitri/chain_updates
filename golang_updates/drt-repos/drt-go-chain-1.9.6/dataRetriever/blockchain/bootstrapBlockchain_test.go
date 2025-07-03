package blockchain

import (
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/stretchr/testify/assert"
)

func TestNewBootstrapBlockchain(t *testing.T) {
	t.Parallel()

	blockchain := NewBootstrapBlockchain()
	assert.False(t, check.IfNil(blockchain))
	providedHeaderHandler := &testscommon.HeaderHandlerStub{}
	assert.Nil(t, blockchain.SetCurrentBlockHeaderAndRootHash(providedHeaderHandler, nil))
	assert.Equal(t, providedHeaderHandler, blockchain.GetCurrentBlockHeader())
}
