package requesterscontainer

import (
	"github.com/TerraDharitri/drt-go-chain-core/data/typeConverters"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/p2p"
	"github.com/TerraDharitri/drt-go-chain/sharding"
)

// FactoryArgs will hold the arguments for RequestersContainerFactory for both shard and meta
type FactoryArgs struct {
	RequesterConfig                 config.RequesterConfig
	ShardCoordinator                sharding.Coordinator
	MainMessenger                   p2p.Messenger
	FullArchiveMessenger            p2p.Messenger
	Marshaller                      marshal.Marshalizer
	Uint64ByteSliceConverter        typeConverters.Uint64ByteSliceConverter
	OutputAntifloodHandler          dataRetriever.P2PAntifloodHandler
	CurrentNetworkEpochProvider     dataRetriever.CurrentNetworkEpochProviderHandler
	MainPreferredPeersHolder        p2p.PreferredPeersHolderHandler
	FullArchivePreferredPeersHolder p2p.PreferredPeersHolderHandler
	PeersRatingHandler              dataRetriever.PeersRatingHandler
	SizeCheckDelta                  uint32
	EnableEpochsHandler             common.EnableEpochsHandler
}
