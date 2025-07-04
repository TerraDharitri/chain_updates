package staking

import (
	"math/big"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/TerraDharitri/drt-go-chain-storage/lrucache"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever/dataPool"
	"github.com/TerraDharitri/drt-go-chain/factory"
	integrationMocks "github.com/TerraDharitri/drt-go-chain/integrationTests/mock"
	"github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
	"github.com/TerraDharitri/drt-go-chain/state/accounts"
	"github.com/TerraDharitri/drt-go-chain/storage"
	"github.com/TerraDharitri/drt-go-chain/testscommon/chainParameters"
	nodesSetupMock "github.com/TerraDharitri/drt-go-chain/testscommon/genesisMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/stakingcommon"
)

const (
	shuffleBetweenShards = false
	adaptivity           = false
	hysteresis           = float32(0.2)
	initialRating        = 5
)

func createNodesCoordinator(
	eligibleMap map[uint32][]nodesCoordinator.Validator,
	waitingMap map[uint32][]nodesCoordinator.Validator,
	numOfMetaNodes uint32,
	numOfShards uint32,
	numOfEligibleNodesPerShard uint32,
	shardConsensusGroupSize int,
	metaConsensusGroupSize int,
	coreComponents factory.CoreComponentsHolder,
	bootStorer storage.Storer,
	nodesCoordinatorRegistryFactory nodesCoordinator.NodesCoordinatorRegistryFactory,
	maxNodesConfig []config.MaxNodesChangeConfig,
) nodesCoordinator.NodesCoordinator {
	shufflerArgs := &nodesCoordinator.NodesShufflerArgs{
		ShuffleBetweenShards: shuffleBetweenShards,
		MaxNodesEnableConfig: maxNodesConfig,
		EnableEpochs: config.EnableEpochs{
			StakingV4Step2EnableEpoch: stakingV4Step2EnableEpoch,
			StakingV4Step3EnableEpoch: stakingV4Step3EnableEpoch,
		},
		EnableEpochsHandler: coreComponents.EnableEpochsHandler(),
	}

	nodeShuffler, _ := nodesCoordinator.NewHashValidatorsShuffler(shufflerArgs)
	cache, _ := lrucache.NewCache(10000)
	argumentsNodesCoordinator := nodesCoordinator.ArgNodesCoordinator{
		ChainParametersHandler: &chainParameters.ChainParametersHandlerStub{
			ChainParametersForEpochCalled: func(epoch uint32) (config.ChainParametersByEpochConfig, error) {
				return config.ChainParametersByEpochConfig{
					RoundDuration:               0,
					Hysteresis:                  hysteresis,
					ShardConsensusGroupSize:     uint32(shardConsensusGroupSize),
					ShardMinNumNodes:            numOfEligibleNodesPerShard,
					MetachainConsensusGroupSize: uint32(metaConsensusGroupSize),
					MetachainMinNumNodes:        numOfMetaNodes,
					Adaptivity:                  adaptivity,
				}, nil
			},
		},
		Marshalizer:                     coreComponents.InternalMarshalizer(),
		Hasher:                          coreComponents.Hasher(),
		ShardIDAsObserver:               core.MetachainShardId,
		NbShards:                        numOfShards,
		EligibleNodes:                   eligibleMap,
		WaitingNodes:                    waitingMap,
		SelfPublicKey:                   eligibleMap[core.MetachainShardId][0].PubKey(),
		ConsensusGroupCache:             cache,
		ShuffledOutHandler:              &integrationMocks.ShuffledOutHandlerStub{},
		ChanStopNode:                    coreComponents.ChanStopNodeProcess(),
		IsFullArchive:                   false,
		Shuffler:                        nodeShuffler,
		BootStorer:                      bootStorer,
		EpochStartNotifier:              coreComponents.EpochStartNotifierWithConfirm(),
		NodesCoordinatorRegistryFactory: nodesCoordinatorRegistryFactory,
		NodeTypeProvider:                coreComponents.NodeTypeProvider(),
		EnableEpochsHandler:             coreComponents.EnableEpochsHandler(),
		ValidatorInfoCacher:             dataPool.NewCurrentEpochValidatorInfoPool(),
		GenesisNodesSetupHandler:        &nodesSetupMock.NodesSetupStub{},
	}

	baseNodesCoordinator, _ := nodesCoordinator.NewIndexHashedNodesCoordinator(argumentsNodesCoordinator)
	nodesCoord, _ := nodesCoordinator.NewIndexHashedNodesCoordinatorWithRater(baseNodesCoordinator, coreComponents.Rater())
	return nodesCoord
}

func createGenesisNodes(
	numOfMetaNodes uint32,
	numOfShards uint32,
	numOfNodesPerShard uint32,
	numOfWaitingNodesPerShard uint32,
	marshaller marshal.Marshalizer,
	stateComponents factory.StateComponentsHandler,
) (map[uint32][]nodesCoordinator.Validator, map[uint32][]nodesCoordinator.Validator) {
	addressStartIdx := uint32(0)
	eligibleGenesisNodes := generateGenesisNodeInfoMap(numOfMetaNodes, numOfShards, numOfNodesPerShard, addressStartIdx)
	eligibleValidators, _ := nodesCoordinator.NodesInfoToValidators(eligibleGenesisNodes)

	addressStartIdx = numOfMetaNodes + numOfShards*numOfNodesPerShard
	waitingGenesisNodes := generateGenesisNodeInfoMap(numOfWaitingNodesPerShard, numOfShards, numOfWaitingNodesPerShard, addressStartIdx)
	waitingValidators, _ := nodesCoordinator.NodesInfoToValidators(waitingGenesisNodes)

	registerValidators(eligibleValidators, stateComponents, marshaller, common.EligibleList)
	registerValidators(waitingValidators, stateComponents, marshaller, common.WaitingList)

	return eligibleValidators, waitingValidators
}

func createGenesisNodesWithCustomConfig(
	owners map[string]*OwnerStats,
	marshaller marshal.Marshalizer,
	stateComponents factory.StateComponentsHandler,
) (map[uint32][]nodesCoordinator.Validator, map[uint32][]nodesCoordinator.Validator) {
	eligibleGenesis := make(map[uint32][]nodesCoordinator.GenesisNodeInfoHandler)
	waitingGenesis := make(map[uint32][]nodesCoordinator.GenesisNodeInfoHandler)

	for owner, ownerStats := range owners {
		registerOwnerKeys(
			[]byte(owner),
			ownerStats.EligibleBlsKeys,
			ownerStats.TotalStake,
			stateComponents,
			marshaller,
			common.EligibleList,
			eligibleGenesis,
		)

		registerOwnerKeys(
			[]byte(owner),
			ownerStats.WaitingBlsKeys,
			ownerStats.TotalStake,
			stateComponents,
			marshaller,
			common.WaitingList,
			waitingGenesis,
		)
	}

	eligible, _ := nodesCoordinator.NodesInfoToValidators(eligibleGenesis)
	waiting, _ := nodesCoordinator.NodesInfoToValidators(waitingGenesis)

	return eligible, waiting
}

func generateGenesisNodeInfoMap(
	numOfMetaNodes uint32,
	numOfShards uint32,
	numOfNodesPerShard uint32,
	addressStartIdx uint32,
) map[uint32][]nodesCoordinator.GenesisNodeInfoHandler {
	validatorsMap := make(map[uint32][]nodesCoordinator.GenesisNodeInfoHandler)
	id := addressStartIdx
	for shardId := uint32(0); shardId < numOfShards; shardId++ {
		for n := uint32(0); n < numOfNodesPerShard; n++ {
			addr := generateAddress(id)
			validator := integrationMocks.NewNodeInfo(addr, addr, shardId, initialRating)
			validatorsMap[shardId] = append(validatorsMap[shardId], validator)
			id++
		}
	}

	for n := uint32(0); n < numOfMetaNodes; n++ {
		addr := generateAddress(id)
		validator := integrationMocks.NewNodeInfo(addr, addr, core.MetachainShardId, initialRating)
		validatorsMap[core.MetachainShardId] = append(validatorsMap[core.MetachainShardId], validator)
		id++
	}

	return validatorsMap
}

func registerOwnerKeys(
	owner []byte,
	ownerPubKeys map[uint32][][]byte,
	totalStake *big.Int,
	stateComponents factory.StateComponentsHolder,
	marshaller marshal.Marshalizer,
	list common.PeerType,
	allNodes map[uint32][]nodesCoordinator.GenesisNodeInfoHandler,
) {
	for shardID, pubKeysInShard := range ownerPubKeys {
		for _, pubKey := range pubKeysInShard {
			validator := integrationMocks.NewNodeInfo(pubKey, pubKey, shardID, initialRating)
			allNodes[shardID] = append(allNodes[shardID], validator)

			savePeerAcc(stateComponents, pubKey, shardID, list)
		}
		stakingcommon.RegisterValidatorKeys(
			stateComponents.AccountsAdapter(),
			owner,
			owner,
			pubKeysInShard,
			totalStake,
			marshaller,
		)
	}
}

func registerValidators(
	validators map[uint32][]nodesCoordinator.Validator,
	stateComponents factory.StateComponentsHolder,
	marshaller marshal.Marshalizer,
	list common.PeerType,
) {
	for shardID, validatorsInShard := range validators {
		for idx, val := range validatorsInShard {
			pubKey := val.PubKey()
			savePeerAcc(stateComponents, pubKey, shardID, list)

			stakingcommon.RegisterValidatorKeys(
				stateComponents.AccountsAdapter(),
				pubKey,
				pubKey,
				[][]byte{pubKey},
				big.NewInt(nodePrice+int64(idx)),
				marshaller,
			)
		}
	}
}

func savePeerAcc(
	stateComponents factory.StateComponentsHolder,
	pubKey []byte,
	shardID uint32,
	list common.PeerType,
) {
	peerAccount, _ := accounts.NewPeerAccount(pubKey)
	peerAccount.SetTempRating(initialRating)
	peerAccount.ShardId = shardID
	peerAccount.BLSPublicKey = pubKey
	peerAccount.List = string(list)
	_ = stateComponents.PeerAccounts().SaveAccount(peerAccount)
}
