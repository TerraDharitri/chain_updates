package storagerequesterscontainer

import (
	"github.com/TerraDharitri/drt-go-chain-core/data/endProcess"
	"github.com/TerraDharitri/drt-go-chain-core/data/typeConverters"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/p2p"
	"github.com/TerraDharitri/drt-go-chain/sharding"
)

// FactoryArgs will hold the arguments for RequestersContainerFactory for both shard and meta
type FactoryArgs struct {
	GeneralConfig            config.Config
	ShardIDForTries          uint32
	ChainID                  string
	WorkingDirectory         string
	Hasher                   hashing.Hasher
	ShardCoordinator         sharding.Coordinator
	Messenger                p2p.Messenger
	Store                    dataRetriever.StorageService
	Marshalizer              marshal.Marshalizer
	Uint64ByteSliceConverter typeConverters.Uint64ByteSliceConverter
	DataPacker               dataRetriever.DataPacker
	ManualEpochStartNotifier dataRetriever.ManualEpochStartNotifier
	ChanGracefullyClose      chan endProcess.ArgEndProcess
	EnableEpochsHandler      common.EnableEpochsHandler
	StateStatsHandler        common.StateStatisticsHandler
}
