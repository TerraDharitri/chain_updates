package nodesCoordinator

import (
	"github.com/TerraDharitri/drt-go-chain-core/data/endProcess"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/epochStart"
	"github.com/TerraDharitri/drt-go-chain/storage"
)

// ArgNodesCoordinator holds all dependencies required by the nodes coordinator in order to create new instances
type ArgNodesCoordinator struct {
	ChainParametersHandler          ChainParametersHandler
	Marshalizer                     marshal.Marshalizer
	Hasher                          hashing.Hasher
	Shuffler                        NodesShuffler
	EpochStartNotifier              EpochStartEventNotifier
	BootStorer                      storage.Storer
	ShardIDAsObserver               uint32
	NbShards                        uint32
	EligibleNodes                   map[uint32][]Validator
	WaitingNodes                    map[uint32][]Validator
	SelfPublicKey                   []byte
	Epoch                           uint32
	StartEpoch                      uint32
	ConsensusGroupCache             Cacher
	ShuffledOutHandler              ShuffledOutHandler
	ChanStopNode                    chan endProcess.ArgEndProcess
	NodeTypeProvider                NodeTypeProviderHandler
	IsFullArchive                   bool
	EnableEpochsHandler             common.EnableEpochsHandler
	ValidatorInfoCacher             epochStart.ValidatorInfoCacher
	GenesisNodesSetupHandler        GenesisNodesSetupHandler
	NodesCoordinatorRegistryFactory NodesCoordinatorRegistryFactory
}
