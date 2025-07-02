package factory

import (
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/trie"
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
)

var _ process.InterceptedDataFactory = (*interceptedTrieNodeDataFactory)(nil)

type interceptedTrieNodeDataFactory struct {
	hasher hashing.Hasher
}

// NewInterceptedTrieNodeDataFactory creates an instance of interceptedTrieNodeDataFactory
func NewInterceptedTrieNodeDataFactory(
	argument *ArgInterceptedDataFactory,
) (*interceptedTrieNodeDataFactory, error) {

	if argument == nil {
		return nil, process.ErrNilArgumentStruct
	}
	if check.IfNil(argument.CoreComponents) {
		return nil, process.ErrNilCoreComponentsHolder
	}
	if check.IfNil(argument.CoreComponents.Hasher()) {
		return nil, process.ErrNilHasher
	}

	return &interceptedTrieNodeDataFactory{
		hasher: argument.CoreComponents.Hasher(),
	}, nil
}

// Create creates instances of InterceptedData by unmarshalling provided buffer
func (sidf *interceptedTrieNodeDataFactory) Create(buff []byte, _ core.PeerID) (process.InterceptedData, error) {
	return trie.NewInterceptedTrieNode(buff, sidf.hasher)
}

// IsInterfaceNil returns true if there is no value under the interface
func (sidf *interceptedTrieNodeDataFactory) IsInterfaceNil() bool {
	return sidf == nil
}
