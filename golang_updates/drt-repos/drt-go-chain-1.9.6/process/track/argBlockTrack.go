package track

import (
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/sharding"
)

// ArgBaseTracker holds all dependencies required by the process data factory in order to create
// new instances of shard/meta block tracker
type ArgBaseTracker struct {
	Hasher                        hashing.Hasher
	HeaderValidator               process.HeaderConstructionValidator
	Marshalizer                   marshal.Marshalizer
	RequestHandler                process.RequestHandler
	RoundHandler                  process.RoundHandler
	ShardCoordinator              sharding.Coordinator
	Store                         dataRetriever.StorageService
	StartHeaders                  map[uint32]data.HeaderHandler
	PoolsHolder                   dataRetriever.PoolsHolder
	WhitelistHandler              process.WhiteListHandler
	FeeHandler                    process.FeeHandler
	EnableEpochsHandler           common.EnableEpochsHandler
	ProofsPool                    process.ProofsPool
	EpochChangeGracePeriodHandler common.EpochChangeGracePeriodHandler
	IsImportDBMode                bool
}

// ArgShardTracker holds all dependencies required by the process data factory in order to create
// new instances of shard block tracker
type ArgShardTracker struct {
	ArgBaseTracker
}

// ArgMetaTracker holds all dependencies required by the process data factory in order to create
// new instances of meta block tracker
type ArgMetaTracker struct {
	ArgBaseTracker
}
