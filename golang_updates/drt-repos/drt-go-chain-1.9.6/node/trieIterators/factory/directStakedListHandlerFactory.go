package factory

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain/node/external"
	"github.com/TerraDharitri/drt-go-chain/node/trieIterators"
	"github.com/TerraDharitri/drt-go-chain/node/trieIterators/disabled"
)

// CreateDirectStakedListHandler will create a new instance of DirectStakedListHandler
func CreateDirectStakedListHandler(args trieIterators.ArgTrieIteratorProcessor) (external.DirectStakedListHandler, error) {
	//TODO add unit tests
	if args.ShardID != core.MetachainShardId {
		return disabled.NewDisabledDirectStakedListProcessor(), nil
	}

	return trieIterators.NewDirectStakedListProcessor(args)
}
