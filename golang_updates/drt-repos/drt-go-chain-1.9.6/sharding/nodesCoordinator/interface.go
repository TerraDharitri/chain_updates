package nodesCoordinator

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"

	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/epochStart"
	"github.com/TerraDharitri/drt-go-chain/state"
)

// Validator defines a node that can be allocated to a shard for participation in a consensus group as validator
// or block proposer
type Validator interface {
	PubKey() []byte
	Chances() uint32
	Index() uint32
	Size() int
}

// NodesCoordinator defines the behaviour of a struct able to do validator group selection
type NodesCoordinator interface {
	NodesCoordinatorHelper
	PublicKeysSelector
	ComputeConsensusGroup(randomness []byte, round uint64, shardId uint32, epoch uint32) (leader Validator, validatorsGroup []Validator, err error)
	GetValidatorWithPublicKey(publicKey []byte) (validator Validator, shardId uint32, err error)
	LoadState(key []byte) error
	GetSavedStateKey() []byte
	ShardIdForEpoch(epoch uint32) (uint32, error)
	ShuffleOutForEpoch(_ uint32)
	GetConsensusWhitelistedNodes(epoch uint32) (map[string]struct{}, error)
	ConsensusGroupSizeForShardAndEpoch(uint32, uint32) int
	GetNumTotalEligible() uint64
	GetWaitingEpochsLeftForPublicKey(publicKey []byte) (uint32, error)
	GetCachedEpochs() map[uint32]struct{}
	IsInterfaceNil() bool
}

// EpochStartEventNotifier provides Register and Unregister functionality for the end of epoch events
type EpochStartEventNotifier interface {
	RegisterHandler(handler epochStart.ActionHandler)
	UnregisterHandler(handler epochStart.ActionHandler)
	IsInterfaceNil() bool
}

// PublicKeysSelector allows retrieval of eligible validators public keys
type PublicKeysSelector interface {
	GetValidatorsIndexes(publicKeys []string, epoch uint32) ([]uint64, error)
	GetAllEligibleValidatorsPublicKeys(epoch uint32) (map[uint32][][]byte, error)
	GetAllEligibleValidatorsPublicKeysForShard(epoch uint32, shardID uint32) ([]string, error)
	GetAllWaitingValidatorsPublicKeys(epoch uint32) (map[uint32][][]byte, error)
	GetAllLeavingValidatorsPublicKeys(epoch uint32) (map[uint32][][]byte, error)
	GetAllShuffledOutValidatorsPublicKeys(epoch uint32) (map[uint32][][]byte, error)
	GetShuffledOutToAuctionValidatorsPublicKeys(epoch uint32) (map[uint32][][]byte, error)
	GetConsensusValidatorsPublicKeys(randomness []byte, round uint64, shardId uint32, epoch uint32) (string, []string, error)
	GetOwnPublicKey() []byte
}

// NodesShuffler provides shuffling functionality for nodes
type NodesShuffler interface {
	UpdateNodeLists(args ArgsUpdateNodes) (*ResUpdateNodes, error)
	IsInterfaceNil() bool
}

// NodesCoordinatorHelper provides polymorphism functionality for nodesCoordinator
type NodesCoordinatorHelper interface {
	ValidatorsWeights(validators []Validator) ([]uint32, error)
	ComputeAdditionalLeaving(allValidators []*state.ShardValidatorInfo) (map[uint32][]Validator, error)
	GetChance(uint32) uint32
}

// ChanceComputer provides chance computation capabilities based on a rating
type ChanceComputer interface {
	// GetChance returns the chances for the rating
	GetChance(uint32) uint32
	// IsInterfaceNil verifies if the interface is nil
	IsInterfaceNil() bool
}

// Cacher provides the capabilities needed to store and retrieve information needed in the NodesCoordinator
type Cacher interface {
	// Clear is used to completely clear the cache.
	Clear()
	// Put adds a value to the cache.  Returns true if an eviction occurred.
	Put(key []byte, value interface{}, sizeInBytes int) (evicted bool)
	// Get looks up a key's value from the cache.
	Get(key []byte) (value interface{}, ok bool)
}

// ShuffledOutHandler defines the methods needed for the computation of a shuffled out event
type ShuffledOutHandler interface {
	Process(newShardID uint32) error
	RegisterHandler(handler func(newShardID uint32))
	CurrentShardID() uint32
	IsInterfaceNil() bool
}

// RandomSelector selects randomly a subset of elements from a set of data
type RandomSelector interface {
	Select(randSeed []byte, sampleSize uint32) ([]uint32, error)
	IsInterfaceNil() bool
}

// EpochStartActionHandler defines the action taken on epoch start event
type EpochStartActionHandler interface {
	EpochStartAction(hdr data.HeaderHandler)
	EpochStartPrepare(metaHdr data.HeaderHandler, body data.BodyHandler)
	NotifyOrder() uint32
}

// NodeTypeProviderHandler defines the actions needed for a component that can handle the node type
type NodeTypeProviderHandler interface {
	SetType(nodeType core.NodeType)
	GetType() core.NodeType
	IsInterfaceNil() bool
}

// GenesisNodeInfoHandler defines the public methods for the genesis nodes info
type GenesisNodeInfoHandler interface {
	AssignedShard() uint32
	AddressBytes() []byte
	PubKeyBytes() []byte
	GetInitialRating() uint32
	IsInterfaceNil() bool
}

// ValidatorsDistributor distributes validators across shards
type ValidatorsDistributor interface {
	DistributeValidators(destination map[uint32][]Validator, source map[uint32][]Validator, rand []byte, balanced bool) error
	IsInterfaceNil() bool
}

// EpochsConfigUpdateHandler specifies the behaviour needed to update nodes config epochs
type EpochsConfigUpdateHandler interface {
	NodesCoordinator
	SetNodesConfigFromValidatorsInfo(epoch uint32, randomness []byte, validatorsInfo []*state.ShardValidatorInfo) error
	IsEpochInConfig(epoch uint32) bool
}

// GenesisNodesSetupHandler defines a component able to provide the genesis nodes info
type GenesisNodesSetupHandler interface {
	MinShardHysteresisNodes() uint32
	MinMetaHysteresisNodes() uint32
	IsInterfaceNil() bool
}

// EpochValidatorsHandler defines what one epoch configuration for a nodes coordinator should hold
type EpochValidatorsHandler interface {
	GetEligibleValidators() map[string][]*SerializableValidator
	GetWaitingValidators() map[string][]*SerializableValidator
	GetLeavingValidators() map[string][]*SerializableValidator
}

// EpochValidatorsHandlerWithAuction defines what one epoch configuration for a nodes coordinator should hold + shuffled out validators
type EpochValidatorsHandlerWithAuction interface {
	EpochValidatorsHandler
	GetShuffledOutValidators() map[string][]*SerializableValidator
	GetLowWaitingList() bool
}

// NodesCoordinatorRegistryHandler defines what is used to initialize nodes coordinator
type NodesCoordinatorRegistryHandler interface {
	GetEpochsConfig() map[string]EpochValidatorsHandler
	GetCurrentEpoch() uint32
	SetCurrentEpoch(epoch uint32)
}

// NodesCoordinatorRegistryFactory handles NodesCoordinatorRegistryHandler marshall/unmarshall
type NodesCoordinatorRegistryFactory interface {
	CreateNodesCoordinatorRegistry(buff []byte) (NodesCoordinatorRegistryHandler, error)
	GetRegistryData(registry NodesCoordinatorRegistryHandler, epoch uint32) ([]byte, error)
	IsInterfaceNil() bool
}

// EpochNotifier can notify upon an epoch change and provide the current epoch
type EpochNotifier interface {
	RegisterNotifyHandler(handler vmcommon.EpochSubscriberHandler)
	CurrentEpoch() uint32
	CheckEpoch(header data.HeaderHandler)
	IsInterfaceNil() bool
}

// ChainParametersHandler defines the actions that need to be done by a component that can handle chain parameters
type ChainParametersHandler interface {
	CurrentChainParameters() config.ChainParametersByEpochConfig
	AllChainParameters() []config.ChainParametersByEpochConfig
	ChainParametersForEpoch(epoch uint32) (config.ChainParametersByEpochConfig, error)
	IsInterfaceNil() bool
}
