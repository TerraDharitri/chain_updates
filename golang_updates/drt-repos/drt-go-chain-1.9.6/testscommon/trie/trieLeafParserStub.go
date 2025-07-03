package trie

import "github.com/TerraDharitri/drt-go-chain-core/core"

// TrieLeafParserStub -
type TrieLeafParserStub struct {
	ParseLeafCalled func(key []byte, val []byte, version core.TrieNodeVersion) (core.KeyValueHolder, error)
}

// ParseLeaf -
func (tlps *TrieLeafParserStub) ParseLeaf(key []byte, val []byte, version core.TrieNodeVersion) (core.KeyValueHolder, error) {
	if tlps.ParseLeafCalled != nil {
		return tlps.ParseLeafCalled(key, val, version)
	}

	return nil, nil
}

// IsInterfaceNil -
func (tlps *TrieLeafParserStub) IsInterfaceNil() bool {
	return tlps == nil
}
