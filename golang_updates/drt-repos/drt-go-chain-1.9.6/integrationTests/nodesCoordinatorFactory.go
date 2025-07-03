package integrationTests

import (
	"fmt"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/integrationTests/mock"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
	"github.com/TerraDharitri/drt-go-chain/storage"
	"github.com/TerraDharitri/drt-go-chain/testscommon/chainParameters"
	"github.com/TerraDharitri/drt-go-chain/testscommon/enableEpochsHandlerMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/genesisMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/nodeTypeProviderMock"
	vic "github.com/TerraDharitri/drt-go-chain/testscommon/validatorInfoCacher"
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data/endProcess"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
)

// ArgIndexHashedNodesCoordinatorFactory -
type ArgIndexHashedNodesCoordinatorFactory struct {
	nodesPerShard           int
	nbMetaNodes             int
	shardConsensusGroupSize int
	metaConsensusGroupSize  int
	shardId                 uint32
	nbShards                int
	validatorsMap           map[uint32][]nodesCoordinator.Validator
	waitingMap              map[uint32][]nodesCoordinator.Validator
	keyIndex                int
	cp                      *CryptoParams
	epochStartSubscriber    nodesCoordinator.EpochStartEventNotifier
	hasher                  hashing.Hasher
	consensusGroupCache     nodesCoordinator.Cacher
	bootStorer              storage.Storer
}

// IndexHashedNodesCoordinatorFactory -
type IndexHashedNodesCoordinatorFactory struct {
}

// CreateNodesCoordinator -
func (tpn *IndexHashedNodesCoordinatorFactory) CreateNodesCoordinator(arg ArgIndexHashedNodesCoordinatorFactory) nodesCoordinator.NodesCoordinator {

	keys := arg.cp.NodesKeys[arg.shardId][arg.keyIndex]
	pubKeyBytes, _ := keys.MainKey.Pk.ToByteArray()

	nodeShufflerArgs := &nodesCoordinator.NodesShufflerArgs{
		ShuffleBetweenShards: shuffleBetweenShards,
		MaxNodesEnableConfig: nil,
		EnableEpochsHandler:  &enableEpochsHandlerMock.EnableEpochsHandlerStub{},
	}

	nodeShuffler, _ := nodesCoordinator.NewHashValidatorsShuffler(nodeShufflerArgs)
	nodesCoordinatorRegistryFactory, _ := nodesCoordinator.NewNodesCoordinatorRegistryFactory(
		TestMarshalizer,
		StakingV4Step2EnableEpoch,
	)
	argumentsNodesCoordinator := nodesCoordinator.ArgNodesCoordinator{
		ChainParametersHandler: &chainParameters.ChainParametersHandlerStub{
			ChainParametersForEpochCalled: func(_ uint32) (config.ChainParametersByEpochConfig, error) {
				return config.ChainParametersByEpochConfig{
					ShardMinNumNodes:            uint32(arg.nodesPerShard),
					MetachainMinNumNodes:        uint32(arg.nbMetaNodes),
					Hysteresis:                  hysteresis,
					Adaptivity:                  adaptivity,
					ShardConsensusGroupSize:     uint32(arg.shardConsensusGroupSize),
					MetachainConsensusGroupSize: uint32(arg.metaConsensusGroupSize),
				}, nil
			},
		},
		Marshalizer:         TestMarshalizer,
		Hasher:              arg.hasher,
		Shuffler:            nodeShuffler,
		EpochStartNotifier:  arg.epochStartSubscriber,
		ShardIDAsObserver:   arg.shardId,
		NbShards:            uint32(arg.nbShards),
		EligibleNodes:       arg.validatorsMap,
		WaitingNodes:        arg.waitingMap,
		SelfPublicKey:       pubKeyBytes,
		ConsensusGroupCache: arg.consensusGroupCache,
		BootStorer:          arg.bootStorer,
		ShuffledOutHandler:  &mock.ShuffledOutHandlerStub{},
		ChanStopNode:        endProcess.GetDummyEndProcessChannel(),
		NodeTypeProvider:    &nodeTypeProviderMock.NodeTypeProviderStub{},
		IsFullArchive:       false,
		EnableEpochsHandler: &enableEpochsHandlerMock.EnableEpochsHandlerStub{
			GetActivationEpochCalled: func(flag core.EnableEpochFlag) uint32 {
				if flag == common.RefactorPeersMiniBlocksFlag || flag == common.StakingV4Step2Flag {
					return UnreachableEpoch
				}
				return 0
			},
		},
		ValidatorInfoCacher:             &vic.ValidatorInfoCacherStub{},
		GenesisNodesSetupHandler:        &genesisMocks.NodesSetupStub{},
		NodesCoordinatorRegistryFactory: nodesCoordinatorRegistryFactory,
	}
	nodesCoord, err := nodesCoordinator.NewIndexHashedNodesCoordinator(argumentsNodesCoordinator)
	if err != nil {
		fmt.Println("Error creating node coordinator")
	}

	return nodesCoord
}

// IndexHashedNodesCoordinatorWithRaterFactory -
type IndexHashedNodesCoordinatorWithRaterFactory struct {
	sharding.PeerAccountListAndRatingHandler
}

// CreateNodesCoordinator is used for creating a nodes coordinator in the integration tests
// based on the provided parameters
func (ihncrf *IndexHashedNodesCoordinatorWithRaterFactory) CreateNodesCoordinator(
	arg ArgIndexHashedNodesCoordinatorFactory,
) nodesCoordinator.NodesCoordinator {
	keys := arg.cp.NodesKeys[arg.shardId][arg.keyIndex]
	pubKeyBytes, _ := keys.MainKey.Pk.ToByteArray()

	shufflerArgs := &nodesCoordinator.NodesShufflerArgs{
		ShuffleBetweenShards: shuffleBetweenShards,
		MaxNodesEnableConfig: nil,
		EnableEpochsHandler:  &enableEpochsHandlerMock.EnableEpochsHandlerStub{},
	}

	nodeShuffler, _ := nodesCoordinator.NewHashValidatorsShuffler(shufflerArgs)
	nodesCoordinatorRegistryFactory, _ := nodesCoordinator.NewNodesCoordinatorRegistryFactory(
		TestMarshalizer,
		StakingV4Step2EnableEpoch,
	)
	argumentsNodesCoordinator := nodesCoordinator.ArgNodesCoordinator{
		ChainParametersHandler: &chainParameters.ChainParametersHandlerStub{
			ChainParametersForEpochCalled: func(_ uint32) (config.ChainParametersByEpochConfig, error) {
				return config.ChainParametersByEpochConfig{
					ShardMinNumNodes:            uint32(arg.nodesPerShard),
					MetachainMinNumNodes:        uint32(arg.nbMetaNodes),
					Hysteresis:                  hysteresis,
					Adaptivity:                  adaptivity,
					ShardConsensusGroupSize:     uint32(arg.shardConsensusGroupSize),
					MetachainConsensusGroupSize: uint32(arg.metaConsensusGroupSize),
				}, nil
			},
		},
		Marshalizer:         TestMarshalizer,
		Hasher:              arg.hasher,
		Shuffler:            nodeShuffler,
		EpochStartNotifier:  arg.epochStartSubscriber,
		ShardIDAsObserver:   arg.shardId,
		NbShards:            uint32(arg.nbShards),
		EligibleNodes:       arg.validatorsMap,
		WaitingNodes:        arg.waitingMap,
		SelfPublicKey:       pubKeyBytes,
		ConsensusGroupCache: arg.consensusGroupCache,
		BootStorer:          arg.bootStorer,
		ShuffledOutHandler:  &mock.ShuffledOutHandlerStub{},
		ChanStopNode:        endProcess.GetDummyEndProcessChannel(),
		NodeTypeProvider:    &nodeTypeProviderMock.NodeTypeProviderStub{},
		IsFullArchive:       false,
		EnableEpochsHandler: &enableEpochsHandlerMock.EnableEpochsHandlerStub{
			GetActivationEpochCalled: func(flag core.EnableEpochFlag) uint32 {
				if flag == common.RefactorPeersMiniBlocksFlag {
					return UnreachableEpoch
				}
				return 0
			},
		},
		ValidatorInfoCacher:             &vic.ValidatorInfoCacherStub{},
		GenesisNodesSetupHandler:        &genesisMocks.NodesSetupStub{},
		NodesCoordinatorRegistryFactory: nodesCoordinatorRegistryFactory,
	}

	baseCoordinator, err := nodesCoordinator.NewIndexHashedNodesCoordinator(argumentsNodesCoordinator)
	if err != nil {
		log.Debug("Error creating node coordinator")
	}

	nodesCoord, err := nodesCoordinator.NewIndexHashedNodesCoordinatorWithRater(baseCoordinator, ihncrf.PeerAccountListAndRatingHandler)
	if err != nil {
		log.Debug("Error creating node coordinator")
	}

	return &NodesWithRater{
		NodesCoordinator: nodesCoord,
		rater:            ihncrf.PeerAccountListAndRatingHandler,
	}
}

// NodesWithRater -
type NodesWithRater struct {
	nodesCoordinator.NodesCoordinator
	rater sharding.PeerAccountListAndRatingHandler
}

// IsInterfaceNil -
func (nwr *NodesWithRater) IsInterfaceNil() bool {
	return nwr == nil
}
