package mock

import crypto "github.com/TerraDharitri/drt-go-chain-core"

type PrivateKeysLoaderStub struct {
	PrivateKeysByShardCalled func() (map[uint32][]crypto.PrivateKey, error)
}

func (pkls *PrivateKeysLoaderStub) PrivateKeysByShard() (map[uint32][]crypto.PrivateKey, error) {
	return pkls.PrivateKeysByShardCalled()
}
