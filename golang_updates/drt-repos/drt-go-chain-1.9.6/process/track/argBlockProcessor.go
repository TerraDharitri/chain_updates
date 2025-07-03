package track

import (
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/sharding"
)

// ArgBlockProcessor holds all dependencies required to process tracked blocks in order to create new instances of
// block processor
type ArgBlockProcessor struct {
	HeaderValidator                       process.HeaderConstructionValidator
	RequestHandler                        process.RequestHandler
	ShardCoordinator                      sharding.Coordinator
	BlockTracker                          blockTrackerHandler
	CrossNotarizer                        blockNotarizerHandler
	SelfNotarizer                         blockNotarizerHandler
	CrossNotarizedHeadersNotifier         blockNotifierHandler
	SelfNotarizedFromCrossHeadersNotifier blockNotifierHandler
	SelfNotarizedHeadersNotifier          blockNotifierHandler
	FinalMetachainHeadersNotifier         blockNotifierHandler
	RoundHandler                          process.RoundHandler
	EnableEpochsHandler                   common.EnableEpochsHandler
	ProofsPool                            process.ProofsPool
	Marshaller                            marshal.Marshalizer
	Hasher                                hashing.Hasher
	HeadersPool                           dataRetriever.HeadersPool
	IsImportDBMode                        bool
}
