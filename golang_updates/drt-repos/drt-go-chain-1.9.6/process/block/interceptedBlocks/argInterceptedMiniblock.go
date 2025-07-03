package interceptedBlocks

import (
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
)

// ArgInterceptedMiniblock is the argument for the intercepted miniblock
type ArgInterceptedMiniblock struct {
	MiniblockBuff    []byte
	Marshalizer      marshal.Marshalizer
	Hasher           hashing.Hasher
	ShardCoordinator sharding.Coordinator
}
