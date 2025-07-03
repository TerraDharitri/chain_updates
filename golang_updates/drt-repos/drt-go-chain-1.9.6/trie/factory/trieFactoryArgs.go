package factory

import (
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/storage"
)

// TrieFactoryArgs holds the arguments for creating a trie factory
type TrieFactoryArgs struct {
	Marshalizer              marshal.Marshalizer
	Hasher                   hashing.Hasher
	PathManager              storage.PathManagerHandler
	TrieStorageManagerConfig config.TrieStorageManagerConfig
	StateStatsHandler        common.StateStatisticsHandler
}
