package nodesCoordinator

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain-core/data/endProcess"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/hashing/sha256"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever/dataPool"
	"github.com/TerraDharitri/drt-go-chain/epochStart"
	"github.com/TerraDharitri/drt-go-chain/sharding/mock"
	"github.com/TerraDharitri/drt-go-chain/state"
	"github.com/TerraDharitri/drt-go-chain/storage/cache"
	"github.com/TerraDharitri/drt-go-chain/testscommon/chainParameters"
	"github.com/TerraDharitri/drt-go-chain/testscommon/enableEpochsHandlerMock"
	testscommonConsensus "github.com/TerraDharitri/drt-go-chain/testscommon/epochstartmock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/genericMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/hashingMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/nodeTypeProviderMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/shardingMocks/nodesCoordinatorMocks"
	vic "github.com/TerraDharitri/drt-go-chain/testscommon/validatorInfoCacher"
)

const stakingV4Epoch = 444

func createDummyNodesList(nbNodes uint32, suffix string) []Validator {
	list := make([]Validator, 0)
	hasher := sha256.NewSha256()

	for j := uint32(0); j < nbNodes; j++ {
		pk := hasher.Compute(fmt.Sprintf("pk%s_%d", suffix, j))
		list = append(list, newValidatorMock(pk, 1, defaultSelectionChances))
	}

	return list
}

func createDummyNodesMap(nodesPerShard uint32, nbShards uint32, suffix string) map[uint32][]Validator {
	nodesMap := make(map[uint32][]Validator)

	var shard uint32

	for i := uint32(0); i <= nbShards; i++ {
		shard = i
		if i == nbShards {
			shard = core.MetachainShardId
		}
		list := createDummyNodesList(nodesPerShard, suffix+"_i")
		nodesMap[shard] = list
	}

	return nodesMap
}

func isStringSubgroup(a []string, b []string) bool {
	var found bool
	for _, va := range a {
		found = false
		for _, vb := range b {
			if va == vb {
				found = true
				break
			}
		}
		if !found {
			return found
		}
	}

	return found
}

func createNodesCoordinatorRegistryFactory() NodesCoordinatorRegistryFactory {
	ncf, _ := NewNodesCoordinatorRegistryFactory(
		&marshal.GogoProtoMarshalizer{},
		stakingV4Epoch,
	)
	return ncf
}

func createArguments() ArgNodesCoordinator {
	nbShards := uint32(1)
	eligibleMap := createDummyNodesMap(10, nbShards, "eligible")
	waitingMap := createDummyNodesMap(3, nbShards, "waiting")
	shufflerArgs := &NodesShufflerArgs{
		ShuffleBetweenShards: shuffleBetweenShards,
		EnableEpochsHandler:  &mock.EnableEpochsHandlerMock{},
	}
	nodeShuffler, _ := NewHashValidatorsShuffler(shufflerArgs)

	epochStartSubscriber := &testscommonConsensus.EpochStartNotifierStub{}
	bootStorer := genericMocks.NewStorerMock()

	arguments := ArgNodesCoordinator{
		ChainParametersHandler: &chainParameters.ChainParametersHandlerStub{
			CurrentChainParametersCalled: func() config.ChainParametersByEpochConfig {
				return config.ChainParametersByEpochConfig{
					ShardConsensusGroupSize:     1,
					MetachainConsensusGroupSize: 1,
				}
			},
			ChainParametersForEpochCalled: func(_ uint32) (config.ChainParametersByEpochConfig, error) {
				return config.ChainParametersByEpochConfig{
					ShardConsensusGroupSize:     1,
					ShardMinNumNodes:            10,
					MetachainConsensusGroupSize: 1,
					MetachainMinNumNodes:        10,
				}, nil
			},
		},
		Marshalizer:         &mock.MarshalizerMock{},
		Hasher:              &hashingMocks.HasherMock{},
		Shuffler:            nodeShuffler,
		EpochStartNotifier:  epochStartSubscriber,
		BootStorer:          bootStorer,
		NbShards:            nbShards,
		EligibleNodes:       eligibleMap,
		WaitingNodes:        waitingMap,
		SelfPublicKey:       []byte("test"),
		ConsensusGroupCache: &mock.NodesCoordinatorCacheMock{},
		ShuffledOutHandler:  &mock.ShuffledOutHandlerStub{},
		IsFullArchive:       false,
		ChanStopNode:        make(chan endProcess.ArgEndProcess),
		NodeTypeProvider:    &nodeTypeProviderMock.NodeTypeProviderStub{},
		EnableEpochsHandler: &mock.EnableEpochsHandlerMock{
			IsRefactorPeersMiniBlocksFlagEnabledField: true,
		},
		GenesisNodesSetupHandler:        &mock.NodesSetupMock{},
		ValidatorInfoCacher:             &vic.ValidatorInfoCacherStub{},
		NodesCoordinatorRegistryFactory: createNodesCoordinatorRegistryFactory(),
	}
	return arguments
}

func validatorsPubKeys(validators []Validator) []string {
	pKeys := make([]string, len(validators))
	for _, v := range validators {
		pKeys = append(pKeys, string(v.PubKey()))
	}

	return pKeys
}

// ------- NewIndexHashedNodesCoordinator

func TestNewIndexHashedNodesCoordinator_NilHasherShouldErr(t *testing.T) {
	t.Parallel()

	arguments := createArguments()
	arguments.Hasher = nil
	ihnc, err := NewIndexHashedNodesCoordinator(arguments)

	require.Equal(t, ErrNilHasher, err)
	require.Nil(t, ihnc)
}

func TestNewIndexHashedNodesCoordinator_NilChainParametersHandleShouldErr(t *testing.T) {
	t.Parallel()

	arguments := createArguments()
	arguments.ChainParametersHandler = nil
	ihnc, err := NewIndexHashedNodesCoordinator(arguments)

	require.Equal(t, ErrNilChainParametersHandler, err)
	require.Nil(t, ihnc)
}

func TestNewIndexHashedNodesCoordinator_ZeroNbShardsShouldErr(t *testing.T) {
	t.Parallel()

	arguments := createArguments()
	arguments.NbShards = 0
	ihnc, err := NewIndexHashedNodesCoordinator(arguments)

	require.Equal(t, ErrInvalidNumberOfShards, err)
	require.Nil(t, ihnc)
}

func TestNewIndexHashedNodesCoordinator_InvalidShardIdShouldErr(t *testing.T) {
	t.Parallel()

	arguments := createArguments()
	arguments.ShardIDAsObserver = 10
	ihnc, err := NewIndexHashedNodesCoordinator(arguments)

	require.Equal(t, ErrInvalidShardId, err)
	require.Nil(t, ihnc)
}

func TestNewIndexHashedNodesCoordinator_NilSelfPublicKeyShouldErr(t *testing.T) {
	t.Parallel()

	arguments := createArguments()
	arguments.SelfPublicKey = nil
	ihnc, err := NewIndexHashedNodesCoordinator(arguments)

	require.Equal(t, ErrNilPubKey, err)
	require.Nil(t, ihnc)
}

func TestNewIndexHashedNodesCoordinator_NilCacherShouldErr(t *testing.T) {
	arguments := createArguments()
	arguments.ConsensusGroupCache = nil
	ihnc, err := NewIndexHashedNodesCoordinator(arguments)

	require.Equal(t, ErrNilCacher, err)
	require.Nil(t, ihnc)
}

func TestNewIndexHashedNodesCoordinator_NilEnableEpochsHandlerShouldErr(t *testing.T) {
	arguments := createArguments()
	arguments.EnableEpochsHandler = nil
	ihnc, err := NewIndexHashedNodesCoordinator(arguments)

	require.Equal(t, ErrNilEnableEpochsHandler, err)
	require.Nil(t, ihnc)
}

func TestNewIndexHashedNodesCoordinator_InvalidEnableEpochsHandlerShouldErr(t *testing.T) {
	arguments := createArguments()
	arguments.EnableEpochsHandler = enableEpochsHandlerMock.NewEnableEpochsHandlerStubWithNoFlagsDefined()
	ihnc, err := NewIndexHashedNodesCoordinator(arguments)

	require.True(t, errors.Is(err, core.ErrInvalidEnableEpochsHandler))
	require.Nil(t, ihnc)
}

func TestNewIndexHashedNodesCoordinator_NilGenesisNodesSetupHandlerShouldErr(t *testing.T) {
	arguments := createArguments()
	arguments.GenesisNodesSetupHandler = nil
	ihnc, err := NewIndexHashedNodesCoordinator(arguments)
	require.Equal(t, ErrNilGenesisNodesSetupHandler, err)
	require.Nil(t, ihnc)
}

func TestNewIndexHashedGroupSelector_OkValsShouldWork(t *testing.T) {
	t.Parallel()

	arguments := createArguments()
	ihnc, err := NewIndexHashedNodesCoordinator(arguments)

	require.NotNil(t, ihnc)
	require.Nil(t, err)
}

// ------- LoadEligibleList

func TestIndexHashedNodesCoordinator_SetNilEligibleMapShouldErr(t *testing.T) {
	t.Parallel()

	waitingMap := createDummyNodesMap(3, 3, "waiting")
	arguments := createArguments()

	ihnc, _ := NewIndexHashedNodesCoordinator(arguments)
	require.Equal(t, ErrNilInputNodesMap, ihnc.setNodesPerShards(nil, waitingMap, nil, nil, 0, false))
}

func TestIndexHashedNodesCoordinator_SetNilWaitingMapShouldErr(t *testing.T) {
	t.Parallel()

	eligibleMap := createDummyNodesMap(10, 3, "eligible")
	arguments := createArguments()

	ihnc, _ := NewIndexHashedNodesCoordinator(arguments)
	require.Equal(t, ErrNilInputNodesMap, ihnc.setNodesPerShards(eligibleMap, nil, nil, nil, 0, false))
}

func TestIndexHashedNodesCoordinator_OkValShouldWork(t *testing.T) {
	t.Parallel()

	eligibleMap := createDummyNodesMap(10, 3, "eligible")
	waitingMap := createDummyNodesMap(3, 3, "waiting")

	shufflerArgs := &NodesShufflerArgs{
		ShuffleBetweenShards: shuffleBetweenShards,
		MaxNodesEnableConfig: nil,
		EnableEpochsHandler:  &mock.EnableEpochsHandlerMock{},
	}
	nodeShuffler, err := NewHashValidatorsShuffler(shufflerArgs)
	require.Nil(t, err)

	epochStartSubscriber := &testscommonConsensus.EpochStartNotifierStub{}
	bootStorer := genericMocks.NewStorerMock()

	arguments := ArgNodesCoordinator{
		ChainParametersHandler: &chainParameters.ChainParametersHandlerStub{
			CurrentChainParametersCalled: func() config.ChainParametersByEpochConfig {
				return config.ChainParametersByEpochConfig{
					ShardConsensusGroupSize:     2,
					MetachainConsensusGroupSize: 1,
				}
			},
		},
		Marshalizer:                     &mock.MarshalizerMock{},
		Hasher:                          &hashingMocks.HasherMock{},
		Shuffler:                        nodeShuffler,
		EpochStartNotifier:              epochStartSubscriber,
		BootStorer:                      bootStorer,
		NbShards:                        1,
		EligibleNodes:                   eligibleMap,
		WaitingNodes:                    waitingMap,
		SelfPublicKey:                   []byte("key"),
		ConsensusGroupCache:             &mock.NodesCoordinatorCacheMock{},
		ShuffledOutHandler:              &mock.ShuffledOutHandlerStub{},
		ChanStopNode:                    make(chan endProcess.ArgEndProcess),
		NodeTypeProvider:                &nodeTypeProviderMock.NodeTypeProviderStub{},
		EnableEpochsHandler:             &mock.EnableEpochsHandlerMock{},
		ValidatorInfoCacher:             &vic.ValidatorInfoCacherStub{},
		GenesisNodesSetupHandler:        &mock.NodesSetupMock{},
		NodesCoordinatorRegistryFactory: createNodesCoordinatorRegistryFactory(),
	}

	ihnc, err := NewIndexHashedNodesCoordinator(arguments)
	require.Nil(t, err)

	readEligible := ihnc.nodesConfig[arguments.Epoch].eligibleMap[0]
	require.Equal(t, eligibleMap[0], readEligible)
}

// ------- ComputeValidatorsGroup

func TestIndexHashedNodesCoordinator_NewCoordinatorTooFewNodesShouldErr(t *testing.T) {
	t.Parallel()

	eligibleMap := createDummyNodesMap(5, 3, "eligible")
	waitingMap := createDummyNodesMap(3, 3, "waiting")
	shufflerArgs := &NodesShufflerArgs{
		ShuffleBetweenShards: shuffleBetweenShards,
		MaxNodesEnableConfig: nil,
		EnableEpochsHandler:  &mock.EnableEpochsHandlerMock{},
	}
	nodeShuffler, err := NewHashValidatorsShuffler(shufflerArgs)
	require.Nil(t, err)

	epochStartSubscriber := &testscommonConsensus.EpochStartNotifierStub{}
	bootStorer := genericMocks.NewStorerMock()

	arguments := ArgNodesCoordinator{
		ChainParametersHandler: &chainParameters.ChainParametersHandlerStub{
			CurrentChainParametersCalled: func() config.ChainParametersByEpochConfig {
				return config.ChainParametersByEpochConfig{
					ShardConsensusGroupSize:     10,
					MetachainConsensusGroupSize: 1,
				}
			},
			ChainParametersForEpochCalled: func(_ uint32) (config.ChainParametersByEpochConfig, error) {
				return config.ChainParametersByEpochConfig{
					ShardConsensusGroupSize:     10,
					MetachainConsensusGroupSize: 1,
				}, nil
			},
		},
		Marshalizer:                     &mock.MarshalizerMock{},
		Hasher:                          &hashingMocks.HasherMock{},
		Shuffler:                        nodeShuffler,
		EpochStartNotifier:              epochStartSubscriber,
		BootStorer:                      bootStorer,
		NbShards:                        1,
		EligibleNodes:                   eligibleMap,
		WaitingNodes:                    waitingMap,
		SelfPublicKey:                   []byte("key"),
		ConsensusGroupCache:             &mock.NodesCoordinatorCacheMock{},
		ShuffledOutHandler:              &mock.ShuffledOutHandlerStub{},
		ChanStopNode:                    make(chan endProcess.ArgEndProcess),
		NodeTypeProvider:                &nodeTypeProviderMock.NodeTypeProviderStub{},
		EnableEpochsHandler:             &mock.EnableEpochsHandlerMock{},
		ValidatorInfoCacher:             &vic.ValidatorInfoCacherStub{},
		GenesisNodesSetupHandler:        &mock.NodesSetupMock{},
		NodesCoordinatorRegistryFactory: createNodesCoordinatorRegistryFactory(),
	}
	ihnc, err := NewIndexHashedNodesCoordinator(arguments)

	require.Equal(t, ErrSmallShardEligibleListSize, err)
	require.Nil(t, ihnc)
}

func TestIndexHashedNodesCoordinator_ComputeValidatorsGroupNilRandomnessShouldErr(t *testing.T) {
	t.Parallel()

	arguments := createArguments()
	ihnc, _ := NewIndexHashedNodesCoordinator(arguments)
	leader, list2, err := ihnc.ComputeConsensusGroup(nil, 0, 0, 0)

	require.Equal(t, ErrNilRandomness, err)
	require.Nil(t, list2)
	require.Nil(t, leader)
}

func TestIndexHashedNodesCoordinator_ComputeValidatorsGroupInvalidShardIdShouldErr(t *testing.T) {
	t.Parallel()

	arguments := createArguments()
	ihnc, _ := NewIndexHashedNodesCoordinator(arguments)
	leader, list2, err := ihnc.ComputeConsensusGroup([]byte("radomness"), 0, 5, 0)

	require.Equal(t, ErrInvalidShardId, err)
	require.Nil(t, list2)
	require.Nil(t, leader)
}

// ------- functionality tests

func TestIndexHashedNodesCoordinator_ComputeValidatorsGroup1ValidatorShouldReturnSame(t *testing.T) {
	t.Parallel()

	list := []Validator{
		newValidatorMock([]byte("pk0"), 1, defaultSelectionChances),
	}
	tmp := createDummyNodesMap(2, 1, "meta")
	nodesMap := make(map[uint32][]Validator)
	nodesMap[0] = list
	nodesMap[core.MetachainShardId] = tmp[core.MetachainShardId]
	shufflerArgs := &NodesShufflerArgs{
		ShuffleBetweenShards: shuffleBetweenShards,
		MaxNodesEnableConfig: nil,
		EnableEpochsHandler:  &mock.EnableEpochsHandlerMock{},
	}
	nodeShuffler, err := NewHashValidatorsShuffler(shufflerArgs)
	require.Nil(t, err)

	epochStartSubscriber := &testscommonConsensus.EpochStartNotifierStub{}
	bootStorer := genericMocks.NewStorerMock()

	arguments := ArgNodesCoordinator{
		ChainParametersHandler: &chainParameters.ChainParametersHandlerStub{
			CurrentChainParametersCalled: func() config.ChainParametersByEpochConfig {
				return config.ChainParametersByEpochConfig{
					ShardConsensusGroupSize:     1,
					MetachainConsensusGroupSize: 1,
				}
			},
			ChainParametersForEpochCalled: func(_ uint32) (config.ChainParametersByEpochConfig, error) {
				return config.ChainParametersByEpochConfig{
					ShardConsensusGroupSize:     1,
					MetachainConsensusGroupSize: 1,
				}, nil
			},
		},
		Marshalizer:                     &mock.MarshalizerMock{},
		Hasher:                          &hashingMocks.HasherMock{},
		Shuffler:                        nodeShuffler,
		EpochStartNotifier:              epochStartSubscriber,
		BootStorer:                      bootStorer,
		NbShards:                        1,
		EligibleNodes:                   nodesMap,
		WaitingNodes:                    make(map[uint32][]Validator),
		SelfPublicKey:                   []byte("key"),
		ConsensusGroupCache:             &mock.NodesCoordinatorCacheMock{},
		ShuffledOutHandler:              &mock.ShuffledOutHandlerStub{},
		ChanStopNode:                    make(chan endProcess.ArgEndProcess),
		NodeTypeProvider:                &nodeTypeProviderMock.NodeTypeProviderStub{},
		EnableEpochsHandler:             &mock.EnableEpochsHandlerMock{},
		ValidatorInfoCacher:             &vic.ValidatorInfoCacherStub{},
		GenesisNodesSetupHandler:        &mock.NodesSetupMock{},
		NodesCoordinatorRegistryFactory: createNodesCoordinatorRegistryFactory(),
	}
	ihnc, _ := NewIndexHashedNodesCoordinator(arguments)
	leader, list2, err := ihnc.ComputeConsensusGroup([]byte("randomness"), 0, 0, 0)

	require.Nil(t, err)
	require.Equal(t, list, list2)
	require.Equal(t, list[0], leader)
}

func TestIndexHashedNodesCoordinator_ComputeValidatorsGroup400of400For10locksNoMemoization(t *testing.T) {
	consensusGroupSize := 400
	nodesPerShard := uint32(400)
	waitingMap := make(map[uint32][]Validator)
	eligibleMap := createDummyNodesMap(nodesPerShard, 1, "eligible")
	shufflerArgs := &NodesShufflerArgs{
		ShuffleBetweenShards: shuffleBetweenShards,
		MaxNodesEnableConfig: nil,
		EnableEpochsHandler:  &mock.EnableEpochsHandlerMock{},
	}
	nodeShuffler, err := NewHashValidatorsShuffler(shufflerArgs)
	require.Nil(t, err)

	epochStartSubscriber := &testscommonConsensus.EpochStartNotifierStub{}
	bootStorer := genericMocks.NewStorerMock()

	getCounter := int32(0)
	putCounter := int32(0)

	lruCache := &mock.NodesCoordinatorCacheMock{
		PutCalled: func(key []byte, value interface{}, sizeInBytes int) (evicted bool) {
			atomic.AddInt32(&putCounter, 1)
			return false
		},
		GetCalled: func(key []byte) (value interface{}, ok bool) {
			atomic.AddInt32(&getCounter, 1)
			return nil, false
		},
	}

	arguments := ArgNodesCoordinator{
		ChainParametersHandler: &chainParameters.ChainParametersHandlerStub{
			CurrentChainParametersCalled: func() config.ChainParametersByEpochConfig {
				return config.ChainParametersByEpochConfig{
					ShardConsensusGroupSize:     uint32(consensusGroupSize),
					MetachainConsensusGroupSize: 1,
				}
			},
			ChainParametersForEpochCalled: func(_ uint32) (config.ChainParametersByEpochConfig, error) {
				return config.ChainParametersByEpochConfig{
					ShardConsensusGroupSize:     uint32(consensusGroupSize),
					MetachainConsensusGroupSize: 1,
				}, nil
			},
		},
		Marshalizer:                     &mock.MarshalizerMock{},
		Hasher:                          &hashingMocks.HasherMock{},
		Shuffler:                        nodeShuffler,
		EpochStartNotifier:              epochStartSubscriber,
		BootStorer:                      bootStorer,
		NbShards:                        1,
		EligibleNodes:                   eligibleMap,
		WaitingNodes:                    waitingMap,
		SelfPublicKey:                   []byte("key"),
		ConsensusGroupCache:             lruCache,
		ShuffledOutHandler:              &mock.ShuffledOutHandlerStub{},
		ChanStopNode:                    make(chan endProcess.ArgEndProcess),
		NodeTypeProvider:                &nodeTypeProviderMock.NodeTypeProviderStub{},
		EnableEpochsHandler:             &mock.EnableEpochsHandlerMock{},
		ValidatorInfoCacher:             &vic.ValidatorInfoCacherStub{},
		GenesisNodesSetupHandler:        &mock.NodesSetupMock{},
		NodesCoordinatorRegistryFactory: createNodesCoordinatorRegistryFactory(),
	}

	ihnc, err := NewIndexHashedNodesCoordinator(arguments)
	require.Nil(t, err)

	miniBlocks := 10

	var list2 []Validator
	var leader Validator
	for i := 0; i < miniBlocks; i++ {
		for j := 0; j <= i; j++ {
			randomness := strconv.Itoa(j)
			leader, list2, err = ihnc.ComputeConsensusGroup([]byte(randomness), uint64(j), 0, 0)
			require.Nil(t, err)
			require.Equal(t, consensusGroupSize, len(list2))
			require.NotNil(t, leader)
		}
	}

	computationNr := miniBlocks * (miniBlocks + 1) / 2

	require.Equal(t, int32(computationNr), getCounter)
	require.Equal(t, int32(computationNr), putCounter)
}

func TestIndexHashedNodesCoordinator_ComputeValidatorsGroup400of400For10BlocksMemoization(t *testing.T) {
	consensusGroupSize := 400
	nodesPerShard := uint32(400)
	waitingMap := make(map[uint32][]Validator)
	eligibleMap := createDummyNodesMap(nodesPerShard, 1, "eligible")
	shufflerArgs := &NodesShufflerArgs{
		ShuffleBetweenShards: false,
		MaxNodesEnableConfig: nil,
		EnableEpochsHandler:  &mock.EnableEpochsHandlerMock{},
	}
	nodeShuffler, err := NewHashValidatorsShuffler(shufflerArgs)
	require.Nil(t, err)

	epochStartSubscriber := &testscommonConsensus.EpochStartNotifierStub{}
	bootStorer := genericMocks.NewStorerMock()

	getCounter := 0
	putCounter := 0

	mut := sync.Mutex{}

	// consensusGroup := list[0:21]
	cacheMap := make(map[string]interface{})
	lruCache := &mock.NodesCoordinatorCacheMock{
		PutCalled: func(key []byte, value interface{}, sizeInBytes int) (evicted bool) {
			mut.Lock()
			defer mut.Unlock()
			putCounter++
			cacheMap[string(key)] = value
			return false
		},
		GetCalled: func(key []byte) (value interface{}, ok bool) {
			mut.Lock()
			defer mut.Unlock()
			getCounter++
			val, ok := cacheMap[string(key)]
			if ok {
				return val, true
			}
			return nil, false
		},
	}

	arguments := ArgNodesCoordinator{
		ChainParametersHandler: &chainParameters.ChainParametersHandlerStub{
			CurrentChainParametersCalled: func() config.ChainParametersByEpochConfig {
				return config.ChainParametersByEpochConfig{
					ShardConsensusGroupSize:     uint32(consensusGroupSize),
					MetachainConsensusGroupSize: 1,
				}
			},
			ChainParametersForEpochCalled: func(_ uint32) (config.ChainParametersByEpochConfig, error) {
				return config.ChainParametersByEpochConfig{
					ShardConsensusGroupSize:     uint32(consensusGroupSize),
					MetachainConsensusGroupSize: 1,
				}, nil
			},
		},
		Marshalizer:                     &mock.MarshalizerMock{},
		Hasher:                          &hashingMocks.HasherMock{},
		Shuffler:                        nodeShuffler,
		EpochStartNotifier:              epochStartSubscriber,
		BootStorer:                      bootStorer,
		NbShards:                        1,
		EligibleNodes:                   eligibleMap,
		WaitingNodes:                    waitingMap,
		SelfPublicKey:                   []byte("key"),
		ConsensusGroupCache:             lruCache,
		ShuffledOutHandler:              &mock.ShuffledOutHandlerStub{},
		ChanStopNode:                    make(chan endProcess.ArgEndProcess),
		NodeTypeProvider:                &nodeTypeProviderMock.NodeTypeProviderStub{},
		EnableEpochsHandler:             &mock.EnableEpochsHandlerMock{},
		ValidatorInfoCacher:             &vic.ValidatorInfoCacherStub{},
		GenesisNodesSetupHandler:        &mock.NodesSetupMock{},
		NodesCoordinatorRegistryFactory: createNodesCoordinatorRegistryFactory(),
	}

	ihnc, err := NewIndexHashedNodesCoordinator(arguments)
	require.Nil(t, err)

	miniBlocks := 10

	var list2 []Validator
	var leader Validator
	for i := 0; i < miniBlocks; i++ {
		for j := 0; j <= i; j++ {
			randomness := strconv.Itoa(j)
			leader, list2, err = ihnc.ComputeConsensusGroup([]byte(randomness), uint64(j), 0, 0)
			require.Nil(t, err)
			require.Equal(t, consensusGroupSize, len(list2))
			require.NotNil(t, leader)
		}
	}

	computationNr := miniBlocks * (miniBlocks + 1) / 2

	require.Equal(t, computationNr, getCounter)
	require.Equal(t, miniBlocks, putCounter)
}

func TestIndexHashedNodesCoordinator_ComputeValidatorsGroup63of400TestEqualSameParams(t *testing.T) {
	t.Skip("testing consistency - to be run manually")
	lruCache := &mock.NodesCoordinatorCacheMock{
		GetCalled: func(key []byte) (value interface{}, ok bool) {
			return nil, false
		},
		PutCalled: func(key []byte, value interface{}, sizeInBytes int) (evicted bool) {
			return false
		},
	}

	consensusGroupSize := 63
	nodesPerShard := uint32(400)
	waitingMap := make(map[uint32][]Validator)
	eligibleMap := createDummyNodesMap(nodesPerShard, 1, "eligible")

	shufflerArgs := &NodesShufflerArgs{
		ShuffleBetweenShards: shuffleBetweenShards,
		MaxNodesEnableConfig: nil,
		EnableEpochsHandler:  &mock.EnableEpochsHandlerMock{},
	}
	nodeShuffler, err := NewHashValidatorsShuffler(shufflerArgs)
	require.Nil(t, err)

	epochStartSubscriber := &testscommonConsensus.EpochStartNotifierStub{}
	bootStorer := genericMocks.NewStorerMock()

	arguments := ArgNodesCoordinator{
		ChainParametersHandler: &chainParameters.ChainParametersHandlerStub{
			CurrentChainParametersCalled: func() config.ChainParametersByEpochConfig {
				return config.ChainParametersByEpochConfig{
					ShardConsensusGroupSize:     uint32(consensusGroupSize),
					MetachainConsensusGroupSize: 1,
				}
			},
		},
		Marshalizer:              &mock.MarshalizerMock{},
		Hasher:                   &hashingMocks.HasherMock{},
		Shuffler:                 nodeShuffler,
		EpochStartNotifier:       epochStartSubscriber,
		BootStorer:               bootStorer,
		NbShards:                 1,
		EligibleNodes:            eligibleMap,
		WaitingNodes:             waitingMap,
		SelfPublicKey:            []byte("key"),
		ConsensusGroupCache:      lruCache,
		ChanStopNode:             make(chan endProcess.ArgEndProcess),
		NodeTypeProvider:         &nodeTypeProviderMock.NodeTypeProviderStub{},
		EnableEpochsHandler:      &mock.EnableEpochsHandlerMock{},
		ValidatorInfoCacher:      &vic.ValidatorInfoCacherStub{},
		GenesisNodesSetupHandler: &mock.NodesSetupMock{},
	}

	ihnc, err := NewIndexHashedNodesCoordinator(arguments)
	require.Nil(t, err)

	nbDifferentSamplings := 1000
	repeatPerSampling := 100

	list := make([][]Validator, repeatPerSampling)
	var leader Validator
	for i := 0; i < nbDifferentSamplings; i++ {
		randomness := arguments.Hasher.Compute(strconv.Itoa(i))
		fmt.Printf("starting selection with randomness: %s\n", hex.EncodeToString(randomness))
		for j := 0; j < repeatPerSampling; j++ {
			leader, list[j], err = ihnc.ComputeConsensusGroup(randomness, 0, 0, 0)
			require.Nil(t, err)
			require.Equal(t, consensusGroupSize, len(list[j]))
			require.NotNil(t, leader)
		}

		for j := 1; j < repeatPerSampling; j++ {
			require.Equal(t, list[0], list[j])
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func BenchmarkIndexHashedGroupSelector_ComputeValidatorsGroup21of400(b *testing.B) {
	consensusGroupSize := 21
	nodesPerShard := uint32(400)
	waitingMap := make(map[uint32][]Validator)
	eligibleMap := createDummyNodesMap(nodesPerShard, 1, "eligible")
	shufflerArgs := &NodesShufflerArgs{
		ShuffleBetweenShards: shuffleBetweenShards,
		MaxNodesEnableConfig: nil,
		EnableEpochsHandler:  &mock.EnableEpochsHandlerMock{},
	}
	nodeShuffler, err := NewHashValidatorsShuffler(shufflerArgs)
	require.Nil(b, err)

	epochStartSubscriber := &testscommonConsensus.EpochStartNotifierStub{}
	bootStorer := genericMocks.NewStorerMock()

	arguments := ArgNodesCoordinator{
		ChainParametersHandler: &chainParameters.ChainParametersHandlerStub{
			CurrentChainParametersCalled: func() config.ChainParametersByEpochConfig {
				return config.ChainParametersByEpochConfig{
					ShardConsensusGroupSize:     uint32(consensusGroupSize),
					MetachainConsensusGroupSize: 1,
				}
			},
		},
		Marshalizer:              &mock.MarshalizerMock{},
		Hasher:                   &hashingMocks.HasherMock{},
		Shuffler:                 nodeShuffler,
		EpochStartNotifier:       epochStartSubscriber,
		BootStorer:               bootStorer,
		NbShards:                 1,
		EligibleNodes:            eligibleMap,
		WaitingNodes:             waitingMap,
		SelfPublicKey:            []byte("key"),
		ConsensusGroupCache:      &mock.NodesCoordinatorCacheMock{},
		ShuffledOutHandler:       &mock.ShuffledOutHandlerStub{},
		ChanStopNode:             make(chan endProcess.ArgEndProcess),
		NodeTypeProvider:         &nodeTypeProviderMock.NodeTypeProviderStub{},
		EnableEpochsHandler:      &mock.EnableEpochsHandlerMock{},
		ValidatorInfoCacher:      &vic.ValidatorInfoCacherStub{},
		GenesisNodesSetupHandler: &mock.NodesSetupMock{},
	}
	ihnc, _ := NewIndexHashedNodesCoordinator(arguments)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		randomness := strconv.Itoa(i)
		leader, list2, _ := ihnc.ComputeConsensusGroup([]byte(randomness), 0, 0, 0)

		require.Equal(b, consensusGroupSize, len(list2))
		require.NotNil(b, leader)
	}
}

func BenchmarkIndexHashedNodesCoordinator_CopyMaps(b *testing.B) {
	previousConfig := &epochNodesConfig{}

	eligibleMap := generateValidatorMap(400, 3)
	waitingMap := generateValidatorMap(400, 3)

	previousConfig.eligibleMap = eligibleMap
	previousConfig.waitingMap = waitingMap

	testMutex := sync.RWMutex{}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		testMutex.RLock()

		copiedPrevious := &epochNodesConfig{}
		copiedPrevious.eligibleMap = copyValidatorMap(previousConfig.eligibleMap)
		copiedPrevious.waitingMap = copyValidatorMap(previousConfig.waitingMap)
		copiedPrevious.nbShards = previousConfig.nbShards

		testMutex.RUnlock()
	}
}

func runBenchmark(consensusGroupCache Cacher, consensusGroupSize int, nodesMap map[uint32][]Validator, b *testing.B) {
	waitingMap := make(map[uint32][]Validator)
	shufflerArgs := &NodesShufflerArgs{
		ShuffleBetweenShards: shuffleBetweenShards,
		MaxNodesEnableConfig: nil,
		EnableEpochsHandler:  &mock.EnableEpochsHandlerMock{},
	}
	nodeShuffler, err := NewHashValidatorsShuffler(shufflerArgs)
	require.Nil(b, err)

	epochStartSubscriber := &testscommonConsensus.EpochStartNotifierStub{}
	bootStorer := genericMocks.NewStorerMock()

	arguments := ArgNodesCoordinator{
		ChainParametersHandler: &chainParameters.ChainParametersHandlerStub{
			CurrentChainParametersCalled: func() config.ChainParametersByEpochConfig {
				return config.ChainParametersByEpochConfig{
					ShardConsensusGroupSize:     uint32(consensusGroupSize),
					MetachainConsensusGroupSize: 1,
				}
			},
		},
		Marshalizer:              &mock.MarshalizerMock{},
		Hasher:                   &hashingMocks.HasherMock{},
		EpochStartNotifier:       epochStartSubscriber,
		Shuffler:                 nodeShuffler,
		BootStorer:               bootStorer,
		NbShards:                 1,
		EligibleNodes:            nodesMap,
		WaitingNodes:             waitingMap,
		SelfPublicKey:            []byte("key"),
		ConsensusGroupCache:      consensusGroupCache,
		ShuffledOutHandler:       &mock.ShuffledOutHandlerStub{},
		ChanStopNode:             make(chan endProcess.ArgEndProcess),
		NodeTypeProvider:         &nodeTypeProviderMock.NodeTypeProviderStub{},
		EnableEpochsHandler:      &mock.EnableEpochsHandlerMock{},
		ValidatorInfoCacher:      &vic.ValidatorInfoCacherStub{},
		GenesisNodesSetupHandler: &mock.NodesSetupMock{},
	}
	ihnc, _ := NewIndexHashedNodesCoordinator(arguments)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		missedBlocks := 1000
		for j := 0; j < missedBlocks; j++ {
			randomness := strconv.Itoa(j)
			leader, list2, _ := ihnc.ComputeConsensusGroup([]byte(randomness), uint64(j), 0, 0)
			require.Equal(b, consensusGroupSize, len(list2))
			require.NotNil(b, leader)
		}
	}
}

func computeMemoryRequirements(consensusGroupCache Cacher, consensusGroupSize int, nodesMap map[uint32][]Validator, b *testing.B) {
	waitingMap := make(map[uint32][]Validator)
	shufflerArgs := &NodesShufflerArgs{
		ShuffleBetweenShards: shuffleBetweenShards,
		MaxNodesEnableConfig: nil,
		EnableEpochsHandler:  &mock.EnableEpochsHandlerMock{},
	}
	nodeShuffler, err := NewHashValidatorsShuffler(shufflerArgs)
	require.Nil(b, err)

	epochStartSubscriber := &testscommonConsensus.EpochStartNotifierStub{}
	bootStorer := genericMocks.NewStorerMock()

	arguments := ArgNodesCoordinator{
		ChainParametersHandler: &chainParameters.ChainParametersHandlerStub{
			CurrentChainParametersCalled: func() config.ChainParametersByEpochConfig {
				return config.ChainParametersByEpochConfig{
					ShardConsensusGroupSize:     uint32(consensusGroupSize),
					MetachainConsensusGroupSize: 1,
				}
			},
		},
		Marshalizer:              &mock.MarshalizerMock{},
		Hasher:                   &hashingMocks.HasherMock{},
		EpochStartNotifier:       epochStartSubscriber,
		Shuffler:                 nodeShuffler,
		BootStorer:               bootStorer,
		NbShards:                 1,
		EligibleNodes:            nodesMap,
		WaitingNodes:             waitingMap,
		SelfPublicKey:            []byte("key"),
		ConsensusGroupCache:      consensusGroupCache,
		ShuffledOutHandler:       &mock.ShuffledOutHandlerStub{},
		ChanStopNode:             make(chan endProcess.ArgEndProcess),
		NodeTypeProvider:         &nodeTypeProviderMock.NodeTypeProviderStub{},
		EnableEpochsHandler:      &mock.EnableEpochsHandlerMock{},
		ValidatorInfoCacher:      &vic.ValidatorInfoCacherStub{},
		GenesisNodesSetupHandler: &mock.NodesSetupMock{},
	}
	ihnc, err := NewIndexHashedNodesCoordinator(arguments)
	require.Nil(b, err)

	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)

	missedBlocks := 1000
	for i := 0; i < missedBlocks; i++ {
		randomness := strconv.Itoa(i)
		leader, list2, _ := ihnc.ComputeConsensusGroup([]byte(randomness), uint64(i), 0, 0)
		require.Equal(b, consensusGroupSize, len(list2))
		require.NotNil(b, leader)
	}

	m2 := runtime.MemStats{}
	runtime.ReadMemStats(&m2)

	fmt.Printf("Used %d MB\n", (m2.HeapAlloc-m.HeapAlloc)/1024/1024)
}

func BenchmarkIndexHashedNodesCoordinator_ComputeValidatorsGroup63of400RecomputeEveryGroup(b *testing.B) {
	consensusGroupSize := 63
	nodesPerShard := uint32(400)
	eligibleMap := createDummyNodesMap(nodesPerShard, 1, "eligible")

	consensusGroupCache, _ := cache.NewLRUCache(1)
	computeMemoryRequirements(consensusGroupCache, consensusGroupSize, eligibleMap, b)
	consensusGroupCache, _ = cache.NewLRUCache(1)
	runBenchmark(consensusGroupCache, consensusGroupSize, eligibleMap, b)
}

func BenchmarkIndexHashedNodesCoordinator_ComputeValidatorsGroup400of400RecomputeEveryGroup(b *testing.B) {
	consensusGroupSize := 400
	nodesPerShard := uint32(400)
	eligibleMap := createDummyNodesMap(nodesPerShard, 1, "eligible")

	consensusGroupCache, _ := cache.NewLRUCache(1)
	computeMemoryRequirements(consensusGroupCache, consensusGroupSize, eligibleMap, b)
	consensusGroupCache, _ = cache.NewLRUCache(1)
	runBenchmark(consensusGroupCache, consensusGroupSize, eligibleMap, b)
}

func BenchmarkIndexHashedNodesCoordinator_ComputeValidatorsGroup63of400Memoization(b *testing.B) {
	consensusGroupSize := 63
	nodesPerShard := uint32(400)
	eligibleMap := createDummyNodesMap(nodesPerShard, 1, "eligible")

	consensusGroupCache, _ := cache.NewLRUCache(10000)
	computeMemoryRequirements(consensusGroupCache, consensusGroupSize, eligibleMap, b)
	consensusGroupCache, _ = cache.NewLRUCache(10000)
	runBenchmark(consensusGroupCache, consensusGroupSize, eligibleMap, b)
}

func BenchmarkIndexHashedNodesCoordinator_ComputeValidatorsGroup400of400Memoization(b *testing.B) {
	consensusGroupSize := 400
	nodesPerShard := uint32(400)
	eligibleMap := createDummyNodesMap(nodesPerShard, 1, "eligible")

	consensusGroupCache, _ := cache.NewLRUCache(1000)
	computeMemoryRequirements(consensusGroupCache, consensusGroupSize, eligibleMap, b)
	consensusGroupCache, _ = cache.NewLRUCache(1000)
	runBenchmark(consensusGroupCache, consensusGroupSize, eligibleMap, b)
}

func TestIndexHashedNodesCoordinator_GetValidatorWithPublicKeyShouldReturnErrNilPubKey(t *testing.T) {
	t.Parallel()

	arguments := createArguments()
	ihnc, _ := NewIndexHashedNodesCoordinator(arguments)

	_, _, err := ihnc.GetValidatorWithPublicKey(nil)
	require.Equal(t, ErrNilPubKey, err)
}

func TestIndexHashedNodesCoordinator_GetValidatorWithPublicKeyShouldReturnErrValidatorNotFound(t *testing.T) {
	t.Parallel()

	arguments := createArguments()
	ihnc, _ := NewIndexHashedNodesCoordinator(arguments)

	_, _, err := ihnc.GetValidatorWithPublicKey([]byte("pk1"))
	require.Equal(t, ErrValidatorNotFound, err)
}

func TestIndexHashedNodesCoordinator_GetValidatorWithPublicKeyShouldWork(t *testing.T) {
	t.Parallel()

	listMeta := []Validator{
		newValidatorMock([]byte("pk0_meta"), 1, defaultSelectionChances),
		newValidatorMock([]byte("pk1_meta"), 1, defaultSelectionChances),
		newValidatorMock([]byte("pk2_meta"), 1, defaultSelectionChances),
	}
	listShard0 := []Validator{
		newValidatorMock([]byte("pk0_shard0"), 1, defaultSelectionChances),
		newValidatorMock([]byte("pk1_shard0"), 1, defaultSelectionChances),
		newValidatorMock([]byte("pk2_shard0"), 1, defaultSelectionChances),
	}
	listShard1 := []Validator{
		newValidatorMock([]byte("pk0_shard1"), 1, defaultSelectionChances),
		newValidatorMock([]byte("pk1_shard1"), 1, defaultSelectionChances),
		newValidatorMock([]byte("pk2_shard1"), 1, defaultSelectionChances),
	}

	eligibleMap := make(map[uint32][]Validator)
	eligibleMap[core.MetachainShardId] = listMeta
	eligibleMap[0] = listShard0
	eligibleMap[1] = listShard1
	shufflerArgs := &NodesShufflerArgs{
		ShuffleBetweenShards: shuffleBetweenShards,
		MaxNodesEnableConfig: nil,
		EnableEpochsHandler:  &mock.EnableEpochsHandlerMock{},
	}
	nodeShuffler, err := NewHashValidatorsShuffler(shufflerArgs)
	require.Nil(t, err)

	epochStartSubscriber := &testscommonConsensus.EpochStartNotifierStub{}
	bootStorer := genericMocks.NewStorerMock()

	arguments := ArgNodesCoordinator{
		ChainParametersHandler: &chainParameters.ChainParametersHandlerStub{
			CurrentChainParametersCalled: func() config.ChainParametersByEpochConfig {
				return config.ChainParametersByEpochConfig{
					ShardConsensusGroupSize:     1,
					MetachainConsensusGroupSize: 1,
				}
			},
		},
		Marshalizer:                     &mock.MarshalizerMock{},
		Hasher:                          &hashingMocks.HasherMock{},
		Shuffler:                        nodeShuffler,
		EpochStartNotifier:              epochStartSubscriber,
		BootStorer:                      bootStorer,
		NbShards:                        2,
		EligibleNodes:                   eligibleMap,
		WaitingNodes:                    make(map[uint32][]Validator),
		SelfPublicKey:                   []byte("key"),
		ConsensusGroupCache:             &mock.NodesCoordinatorCacheMock{},
		ShuffledOutHandler:              &mock.ShuffledOutHandlerStub{},
		ChanStopNode:                    make(chan endProcess.ArgEndProcess),
		NodeTypeProvider:                &nodeTypeProviderMock.NodeTypeProviderStub{},
		EnableEpochsHandler:             &mock.EnableEpochsHandlerMock{},
		ValidatorInfoCacher:             &vic.ValidatorInfoCacherStub{},
		GenesisNodesSetupHandler:        &mock.NodesSetupMock{},
		NodesCoordinatorRegistryFactory: createNodesCoordinatorRegistryFactory(),
	}
	ihnc, _ := NewIndexHashedNodesCoordinator(arguments)

	v, shardId, err := ihnc.GetValidatorWithPublicKey([]byte("pk0_meta"))
	require.Nil(t, err)
	require.Equal(t, core.MetachainShardId, shardId)
	require.Equal(t, []byte("pk0_meta"), v.PubKey())

	v, shardId, err = ihnc.GetValidatorWithPublicKey([]byte("pk1_shard0"))
	require.Nil(t, err)
	require.Equal(t, uint32(0), shardId)
	require.Equal(t, []byte("pk1_shard0"), v.PubKey())

	v, shardId, err = ihnc.GetValidatorWithPublicKey([]byte("pk2_shard1"))
	require.Nil(t, err)
	require.Equal(t, uint32(1), shardId)
	require.Equal(t, []byte("pk2_shard1"), v.PubKey())
}

func TestIndexHashedGroupSelector_GetAllEligibleValidatorsPublicKeys(t *testing.T) {
	t.Parallel()

	shardZeroId := uint32(0)
	shardOneId := uint32(1)
	expectedValidatorsPubKeys := map[uint32][][]byte{
		shardZeroId:           {[]byte("pk0_shard0"), []byte("pk1_shard0"), []byte("pk2_shard0")},
		shardOneId:            {[]byte("pk0_shard1"), []byte("pk1_shard1"), []byte("pk2_shard1")},
		core.MetachainShardId: {[]byte("pk0_meta"), []byte("pk1_meta"), []byte("pk2_meta")},
	}

	listMeta := []Validator{
		newValidatorMock(expectedValidatorsPubKeys[core.MetachainShardId][0], 1, defaultSelectionChances),
		newValidatorMock(expectedValidatorsPubKeys[core.MetachainShardId][1], 1, defaultSelectionChances),
		newValidatorMock(expectedValidatorsPubKeys[core.MetachainShardId][2], 1, defaultSelectionChances),
	}
	listShard0 := []Validator{
		newValidatorMock(expectedValidatorsPubKeys[shardZeroId][0], 1, defaultSelectionChances),
		newValidatorMock(expectedValidatorsPubKeys[shardZeroId][1], 1, defaultSelectionChances),
		newValidatorMock(expectedValidatorsPubKeys[shardZeroId][2], 1, defaultSelectionChances),
	}
	listShard1 := []Validator{
		newValidatorMock(expectedValidatorsPubKeys[shardOneId][0], 1, defaultSelectionChances),
		newValidatorMock(expectedValidatorsPubKeys[shardOneId][1], 1, defaultSelectionChances),
		newValidatorMock(expectedValidatorsPubKeys[shardOneId][2], 1, defaultSelectionChances),
	}

	eligibleMap := make(map[uint32][]Validator)
	eligibleMap[core.MetachainShardId] = listMeta
	eligibleMap[shardZeroId] = listShard0
	eligibleMap[shardOneId] = listShard1
	shufflerArgs := &NodesShufflerArgs{
		ShuffleBetweenShards: shuffleBetweenShards,
		MaxNodesEnableConfig: nil,
		EnableEpochsHandler:  &mock.EnableEpochsHandlerMock{},
	}
	nodeShuffler, err := NewHashValidatorsShuffler(shufflerArgs)
	require.Nil(t, err)

	epochStartSubscriber := &testscommonConsensus.EpochStartNotifierStub{}
	bootStorer := genericMocks.NewStorerMock()

	arguments := ArgNodesCoordinator{
		ChainParametersHandler: &chainParameters.ChainParametersHandlerStub{
			CurrentChainParametersCalled: func() config.ChainParametersByEpochConfig {
				return config.ChainParametersByEpochConfig{
					ShardConsensusGroupSize:     1,
					MetachainConsensusGroupSize: 1,
				}
			},
		},
		Marshalizer:                     &mock.MarshalizerMock{},
		Hasher:                          &hashingMocks.HasherMock{},
		Shuffler:                        nodeShuffler,
		EpochStartNotifier:              epochStartSubscriber,
		BootStorer:                      bootStorer,
		ShardIDAsObserver:               shardZeroId,
		NbShards:                        2,
		EligibleNodes:                   eligibleMap,
		WaitingNodes:                    make(map[uint32][]Validator),
		SelfPublicKey:                   []byte("key"),
		ConsensusGroupCache:             &mock.NodesCoordinatorCacheMock{},
		ShuffledOutHandler:              &mock.ShuffledOutHandlerStub{},
		ChanStopNode:                    make(chan endProcess.ArgEndProcess),
		NodeTypeProvider:                &nodeTypeProviderMock.NodeTypeProviderStub{},
		EnableEpochsHandler:             &mock.EnableEpochsHandlerMock{},
		ValidatorInfoCacher:             &vic.ValidatorInfoCacherStub{},
		GenesisNodesSetupHandler:        &mock.NodesSetupMock{},
		NodesCoordinatorRegistryFactory: createNodesCoordinatorRegistryFactory(),
	}

	ihnc, _ := NewIndexHashedNodesCoordinator(arguments)

	allValidatorsPublicKeys, err := ihnc.GetAllEligibleValidatorsPublicKeys(0)
	require.Equal(t, expectedValidatorsPubKeys, allValidatorsPublicKeys)
	require.Nil(t, err)
}

func TestIndexHashedGroupSelector_GetAllWaitingValidatorsPublicKeys(t *testing.T) {
	t.Parallel()

	shardZeroId := uint32(0)
	shardOneId := uint32(1)
	expectedValidatorsPubKeys := map[uint32][][]byte{
		shardZeroId:           {[]byte("pk0_shard0"), []byte("pk1_shard0"), []byte("pk2_shard0")},
		shardOneId:            {[]byte("pk0_shard1"), []byte("pk1_shard1"), []byte("pk2_shard1")},
		core.MetachainShardId: {[]byte("pk0_meta"), []byte("pk1_meta"), []byte("pk2_meta")},
	}

	listMeta := []Validator{
		newValidatorMock(expectedValidatorsPubKeys[core.MetachainShardId][0], 1, defaultSelectionChances),
		newValidatorMock(expectedValidatorsPubKeys[core.MetachainShardId][1], 1, defaultSelectionChances),
		newValidatorMock(expectedValidatorsPubKeys[core.MetachainShardId][2], 1, defaultSelectionChances),
	}
	listShard0 := []Validator{
		newValidatorMock(expectedValidatorsPubKeys[shardZeroId][0], 1, defaultSelectionChances),
		newValidatorMock(expectedValidatorsPubKeys[shardZeroId][1], 1, defaultSelectionChances),
		newValidatorMock(expectedValidatorsPubKeys[shardZeroId][2], 1, defaultSelectionChances),
	}
	listShard1 := []Validator{
		newValidatorMock(expectedValidatorsPubKeys[shardOneId][0], 1, defaultSelectionChances),
		newValidatorMock(expectedValidatorsPubKeys[shardOneId][1], 1, defaultSelectionChances),
		newValidatorMock(expectedValidatorsPubKeys[shardOneId][2], 1, defaultSelectionChances),
	}

	waitingMap := make(map[uint32][]Validator)
	waitingMap[core.MetachainShardId] = listMeta
	waitingMap[shardZeroId] = listShard0
	waitingMap[shardOneId] = listShard1

	shufflerArgs := &NodesShufflerArgs{
		ShuffleBetweenShards: shuffleBetweenShards,
		MaxNodesEnableConfig: nil,
		EnableEpochsHandler:  &mock.EnableEpochsHandlerMock{},
	}
	nodeShuffler, err := NewHashValidatorsShuffler(shufflerArgs)
	require.Nil(t, err)

	epochStartSubscriber := &testscommonConsensus.EpochStartNotifierStub{}
	bootStorer := genericMocks.NewStorerMock()

	eligibleMap := make(map[uint32][]Validator)
	eligibleMap[core.MetachainShardId] = []Validator{&validator{}}
	eligibleMap[shardZeroId] = []Validator{&validator{}}

	arguments := ArgNodesCoordinator{
		ChainParametersHandler: &chainParameters.ChainParametersHandlerStub{
			CurrentChainParametersCalled: func() config.ChainParametersByEpochConfig {
				return config.ChainParametersByEpochConfig{
					ShardConsensusGroupSize:     1,
					MetachainConsensusGroupSize: 1,
				}
			},
		},
		Marshalizer:                     &mock.MarshalizerMock{},
		Hasher:                          &hashingMocks.HasherMock{},
		Shuffler:                        nodeShuffler,
		EpochStartNotifier:              epochStartSubscriber,
		BootStorer:                      bootStorer,
		ShardIDAsObserver:               shardZeroId,
		NbShards:                        2,
		EligibleNodes:                   eligibleMap,
		WaitingNodes:                    waitingMap,
		SelfPublicKey:                   []byte("key"),
		ConsensusGroupCache:             &mock.NodesCoordinatorCacheMock{},
		ShuffledOutHandler:              &mock.ShuffledOutHandlerStub{},
		ChanStopNode:                    make(chan endProcess.ArgEndProcess),
		NodeTypeProvider:                &nodeTypeProviderMock.NodeTypeProviderStub{},
		EnableEpochsHandler:             &mock.EnableEpochsHandlerMock{},
		ValidatorInfoCacher:             &vic.ValidatorInfoCacherStub{},
		GenesisNodesSetupHandler:        &mock.NodesSetupMock{},
		NodesCoordinatorRegistryFactory: createNodesCoordinatorRegistryFactory(),
	}

	ihnc, _ := NewIndexHashedNodesCoordinator(arguments)

	allValidatorsPublicKeys, err := ihnc.GetAllWaitingValidatorsPublicKeys(0)
	require.Equal(t, expectedValidatorsPubKeys, allValidatorsPublicKeys)
	require.Nil(t, err)
}

func createBlockBodyFromNodesCoordinator(ihnc *indexHashedNodesCoordinator, epoch uint32, validatorInfoCacher epochStart.ValidatorInfoCacher) *block.Body {
	body := &block.Body{MiniBlocks: make([]*block.MiniBlock, 0)}

	mbs := createMiniBlocksForNodesMap(ihnc.nodesConfig[epoch].eligibleMap, string(common.EligibleList), ihnc.marshalizer, ihnc.hasher, validatorInfoCacher)
	body.MiniBlocks = append(body.MiniBlocks, mbs...)

	mbs = createMiniBlocksForNodesMap(ihnc.nodesConfig[epoch].waitingMap, string(common.WaitingList), ihnc.marshalizer, ihnc.hasher, validatorInfoCacher)
	body.MiniBlocks = append(body.MiniBlocks, mbs...)

	mbs = createMiniBlocksForNodesMap(ihnc.nodesConfig[epoch].leavingMap, string(common.LeavingList), ihnc.marshalizer, ihnc.hasher, validatorInfoCacher)
	body.MiniBlocks = append(body.MiniBlocks, mbs...)

	return body
}

func createMiniBlocksForNodesMap(
	nodesMap map[uint32][]Validator,
	list string,
	marshaller marshal.Marshalizer,
	hasher hashing.Hasher,
	validatorInfoCacher epochStart.ValidatorInfoCacher,
) []*block.MiniBlock {

	miniBlocks := make([]*block.MiniBlock, 0)
	for shId, eligibleList := range nodesMap {
		miniBlock := &block.MiniBlock{Type: block.PeerBlock}
		for index, eligible := range eligibleList {
			shardValidatorInfo := &state.ShardValidatorInfo{
				PublicKey:  eligible.PubKey(),
				ShardId:    shId,
				List:       list,
				Index:      uint32(index),
				TempRating: 10,
			}

			shardValidatorInfoHash, _ := core.CalculateHash(marshaller, hasher, shardValidatorInfo)

			miniBlock.TxHashes = append(miniBlock.TxHashes, shardValidatorInfoHash)
			validatorInfoCacher.AddValidatorInfo(shardValidatorInfoHash, shardValidatorInfo)
		}
		miniBlocks = append(miniBlocks, miniBlock)
	}
	return miniBlocks
}

func TestIndexHashedNodesCoordinator_EpochStart(t *testing.T) {
	t.Parallel()

	arguments := createArguments()
	arguments.ValidatorInfoCacher = dataPool.NewCurrentEpochValidatorInfoPool()
	ihnc, err := NewIndexHashedNodesCoordinator(arguments)
	require.Nil(t, err)
	epoch := uint32(1)

	header := &block.MetaBlock{
		PrevRandSeed: []byte("rand seed"),
		EpochStart:   block.EpochStart{LastFinalizedHeaders: []block.EpochStartShardData{{}}},
		Epoch:        epoch,
	}

	ihnc.nodesConfig[epoch] = ihnc.nodesConfig[0]

	body := createBlockBodyFromNodesCoordinator(ihnc, epoch, ihnc.validatorInfoCacher)
	ihnc.EpochStartPrepare(header, body)
	ihnc.EpochStartAction(header)

	validators, err := ihnc.GetAllEligibleValidatorsPublicKeys(epoch)
	require.Nil(t, err)
	require.NotNil(t, validators)

	computedShardId, isValidator := ihnc.computeShardForSelfPublicKey(ihnc.nodesConfig[0])
	// should remain in same shard with intra shard shuffling
	require.Equal(t, arguments.ShardIDAsObserver, computedShardId)
	require.False(t, isValidator)
}

func TestIndexHashedNodesCoordinator_setNodesPerShardsShouldTriggerWrongConfiguration(t *testing.T) {
	t.Parallel()

	chanStopNode := make(chan endProcess.ArgEndProcess, 1)
	arguments := createArguments()
	arguments.ChanStopNode = chanStopNode
	arguments.IsFullArchive = true

	pk := []byte("pk")
	arguments.SelfPublicKey = pk
	ihnc, err := NewIndexHashedNodesCoordinator(arguments)
	require.Nil(t, err)

	eligibleMap := map[uint32][]Validator{
		core.MetachainShardId: {
			newValidatorMock(pk, 1, 1),
		},
	}

	err = ihnc.setNodesPerShards(eligibleMap, map[uint32][]Validator{}, map[uint32][]Validator{}, map[uint32][]Validator{}, 2, false)
	require.NoError(t, err)

	value := <-chanStopNode
	require.Equal(t, common.WrongConfiguration, value.Reason)
}

func TestIndexHashedNodesCoordinator_setNodesPerShardsShouldNotTriggerWrongConfiguration(t *testing.T) {
	t.Parallel()

	chanStopNode := make(chan endProcess.ArgEndProcess, 1)
	arguments := createArguments()
	arguments.ChanStopNode = chanStopNode
	arguments.IsFullArchive = false

	pk := []byte("pk")
	arguments.SelfPublicKey = pk
	ihnc, err := NewIndexHashedNodesCoordinator(arguments)
	require.Nil(t, err)

	eligibleMap := map[uint32][]Validator{
		core.MetachainShardId: {
			newValidatorMock(pk, 1, 1),
		},
	}

	err = ihnc.setNodesPerShards(eligibleMap, map[uint32][]Validator{}, map[uint32][]Validator{}, map[uint32][]Validator{}, 2, false)
	require.NoError(t, err)

	require.Empty(t, chanStopNode)
}

func TestIndexHashedNodesCoordinator_setNodesPerShardsShouldSetNodeTypeValidator(t *testing.T) {
	t.Parallel()

	arguments := createArguments()
	arguments.IsFullArchive = false

	var nodeTypeResult core.NodeType
	var setTypeWasCalled bool
	arguments.NodeTypeProvider = &nodeTypeProviderMock.NodeTypeProviderStub{
		SetTypeCalled: func(nodeType core.NodeType) {
			nodeTypeResult = nodeType
			setTypeWasCalled = true
		},
	}

	pk := []byte("pk")
	arguments.SelfPublicKey = pk
	ihnc, err := NewIndexHashedNodesCoordinator(arguments)
	require.Nil(t, err)

	eligibleMap := map[uint32][]Validator{
		core.MetachainShardId: {
			newValidatorMock(pk, 1, 1),
		},
	}

	err = ihnc.setNodesPerShards(eligibleMap, map[uint32][]Validator{}, map[uint32][]Validator{}, map[uint32][]Validator{}, 2, false)
	require.NoError(t, err)
	require.True(t, setTypeWasCalled)
	require.Equal(t, core.NodeTypeValidator, nodeTypeResult)
}

func TestIndexHashedNodesCoordinator_setNodesPerShardsShouldSetNodeTypeObserver(t *testing.T) {
	t.Parallel()

	arguments := createArguments()
	arguments.IsFullArchive = false

	var nodeTypeResult core.NodeType
	var setTypeWasCalled bool
	arguments.NodeTypeProvider = &nodeTypeProviderMock.NodeTypeProviderStub{
		SetTypeCalled: func(nodeType core.NodeType) {
			nodeTypeResult = nodeType
			setTypeWasCalled = true
		},
	}

	pk := []byte("observer pk")
	arguments.SelfPublicKey = pk
	ihnc, err := NewIndexHashedNodesCoordinator(arguments)
	require.Nil(t, err)

	eligibleMap := map[uint32][]Validator{
		core.MetachainShardId: {
			newValidatorMock([]byte("validator pk"), 1, 1),
		},
	}

	err = ihnc.setNodesPerShards(eligibleMap, map[uint32][]Validator{}, map[uint32][]Validator{}, map[uint32][]Validator{}, 2, false)
	require.NoError(t, err)
	require.True(t, setTypeWasCalled)
	require.Equal(t, core.NodeTypeObserver, nodeTypeResult)
}

func TestIndexHashedNodesCoordinator_EpochStartInEligible(t *testing.T) {
	t.Parallel()

	arguments := createArguments()
	arguments.ValidatorInfoCacher = dataPool.NewCurrentEpochValidatorInfoPool()
	pk := []byte("pk")
	arguments.SelfPublicKey = pk
	ihnc, err := NewIndexHashedNodesCoordinator(arguments)
	require.Nil(t, err)
	epoch := uint32(2)

	header := &block.MetaBlock{
		PrevRandSeed: []byte("rand seed"),
		EpochStart:   block.EpochStart{LastFinalizedHeaders: []block.EpochStartShardData{{}}},
		Epoch:        epoch,
	}

	validatorShard := core.MetachainShardId
	ihnc.nodesConfig = map[uint32]*epochNodesConfig{
		epoch: {
			shardID: validatorShard,
			eligibleMap: map[uint32][]Validator{
				validatorShard: {newValidatorMock(pk, 1, 1)},
			},
		},
	}
	body := createBlockBodyFromNodesCoordinator(ihnc, epoch, ihnc.validatorInfoCacher)
	ihnc.EpochStartPrepare(header, body)
	ihnc.EpochStartAction(header)

	computedShardId, isValidator := ihnc.computeShardForSelfPublicKey(ihnc.nodesConfig[epoch])

	require.Equal(t, validatorShard, computedShardId)
	require.True(t, isValidator)
}

func TestIndexHashedNodesCoordinator_computeShardForSelfPublicKeyWithStakingV4(t *testing.T) {
	t.Parallel()

	arguments := createArguments()
	pk := []byte("pk")
	arguments.SelfPublicKey = pk
	nc, _ := NewIndexHashedNodesCoordinator(arguments)
	epoch := uint32(2)

	metaShard := core.MetachainShardId
	nc.nodesConfig = map[uint32]*epochNodesConfig{
		epoch: {
			shardID: metaShard,
			shuffledOutMap: map[uint32][]Validator{
				metaShard: {newValidatorMock(pk, 1, 1)},
			},
		},
	}

	computedShardId, isValidator := nc.computeShardForSelfPublicKey(nc.nodesConfig[epoch])
	require.Equal(t, nc.shardIDAsObserver, computedShardId)
	require.False(t, isValidator)

	nc.flagStakingV4Step2.SetValue(true)

	computedShardId, isValidator = nc.computeShardForSelfPublicKey(nc.nodesConfig[epoch])
	require.Equal(t, metaShard, computedShardId)
	require.True(t, isValidator)
}

func TestIndexHashedNodesCoordinator_EpochStartInWaiting(t *testing.T) {
	t.Parallel()

	arguments := createArguments()
	arguments.ValidatorInfoCacher = dataPool.NewCurrentEpochValidatorInfoPool()
	pk := []byte("pk")
	arguments.SelfPublicKey = pk
	ihnc, err := NewIndexHashedNodesCoordinator(arguments)
	require.Nil(t, err)

	epoch := uint32(2)
	header := &block.MetaBlock{
		PrevRandSeed: []byte("rand seed"),
		EpochStart:   block.EpochStart{LastFinalizedHeaders: []block.EpochStartShardData{{}}},
		Epoch:        epoch,
	}

	validatorShard := core.MetachainShardId
	ihnc.nodesConfig = map[uint32]*epochNodesConfig{
		epoch: {
			shardID: validatorShard,
			waitingMap: map[uint32][]Validator{
				validatorShard: {newValidatorMock(pk, 1, 1)},
			},
		},
	}
	body := createBlockBodyFromNodesCoordinator(ihnc, epoch, ihnc.validatorInfoCacher)
	ihnc.EpochStartPrepare(header, body)
	ihnc.EpochStartAction(header)

	computedShardId, isValidator := ihnc.computeShardForSelfPublicKey(ihnc.nodesConfig[epoch])
	require.Equal(t, validatorShard, computedShardId)
	require.True(t, isValidator)
}

func TestIndexHashedNodesCoordinator_EpochStartInLeaving(t *testing.T) {
	t.Parallel()

	arguments := createArguments()
	arguments.ValidatorInfoCacher = dataPool.NewCurrentEpochValidatorInfoPool()
	pk := []byte("pk")
	arguments.SelfPublicKey = pk
	ihnc, err := NewIndexHashedNodesCoordinator(arguments)
	require.Nil(t, err)

	epoch := uint32(2)
	header := &block.MetaBlock{
		PrevRandSeed: []byte("rand seed"),
		EpochStart:   block.EpochStart{LastFinalizedHeaders: []block.EpochStartShardData{{}}},
		Epoch:        epoch,
	}

	validatorShard := core.MetachainShardId
	ihnc.nodesConfig = map[uint32]*epochNodesConfig{
		epoch: {
			shardID: validatorShard,
			eligibleMap: map[uint32][]Validator{
				validatorShard: {
					newValidatorMock([]byte("eligiblePk"), 1, 1),
				},
			},
			leavingMap: map[uint32][]Validator{
				validatorShard: {newValidatorMock(pk, 1, 1)},
			},
		},
	}
	body := createBlockBodyFromNodesCoordinator(ihnc, epoch, ihnc.validatorInfoCacher)
	ihnc.EpochStartPrepare(header, body)
	ihnc.EpochStartAction(header)

	computedShardId, isValidator := ihnc.computeShardForSelfPublicKey(ihnc.nodesConfig[epoch])
	require.Equal(t, validatorShard, computedShardId)
	require.True(t, isValidator)
}

func TestIndexHashedNodesCoordinator_EpochStart_EligibleSortedAscendingByIndex(t *testing.T) {
	t.Parallel()

	nbShards := uint32(1)
	eligibleMap := make(map[uint32][]Validator)

	pk1 := []byte{2}
	pk2 := []byte{1}

	list := []Validator{
		newValidatorMock(pk1, 1, 1),
		newValidatorMock(pk2, 1, 1),
	}
	eligibleMap[core.MetachainShardId] = list

	shufflerArgs := &NodesShufflerArgs{
		ShuffleBetweenShards: shuffleBetweenShards,
		MaxNodesEnableConfig: nil,
		EnableEpochsHandler:  &mock.EnableEpochsHandlerMock{},
	}
	nodeShuffler, err := NewHashValidatorsShuffler(shufflerArgs)
	require.Nil(t, err)

	epochStartSubscriber := &testscommonConsensus.EpochStartNotifierStub{}
	bootStorer := genericMocks.NewStorerMock()

	arguments := ArgNodesCoordinator{
		ChainParametersHandler: &chainParameters.ChainParametersHandlerStub{
			CurrentChainParametersCalled: func() config.ChainParametersByEpochConfig {
				return config.ChainParametersByEpochConfig{
					ShardConsensusGroupSize:     1,
					MetachainConsensusGroupSize: 1,
				}
			},
			ChainParametersForEpochCalled: func(_ uint32) (config.ChainParametersByEpochConfig, error) {
				return config.ChainParametersByEpochConfig{
					ShardConsensusGroupSize:     1,
					MetachainConsensusGroupSize: 1,
					ShardMinNumNodes:            2,
					MetachainMinNumNodes:        2,
				}, nil
			},
		},
		Marshalizer:         &mock.MarshalizerMock{},
		Hasher:              &hashingMocks.HasherMock{},
		Shuffler:            nodeShuffler,
		EpochStartNotifier:  epochStartSubscriber,
		BootStorer:          bootStorer,
		NbShards:            nbShards,
		EligibleNodes:       eligibleMap,
		WaitingNodes:        map[uint32][]Validator{},
		SelfPublicKey:       []byte("test"),
		ConsensusGroupCache: &mock.NodesCoordinatorCacheMock{},
		ShuffledOutHandler:  &mock.ShuffledOutHandlerStub{},
		ChanStopNode:        make(chan endProcess.ArgEndProcess),
		NodeTypeProvider:    &nodeTypeProviderMock.NodeTypeProviderStub{},
		EnableEpochsHandler: &mock.EnableEpochsHandlerMock{
			IsRefactorPeersMiniBlocksFlagEnabledField: true,
		},
		ValidatorInfoCacher:             dataPool.NewCurrentEpochValidatorInfoPool(),
		GenesisNodesSetupHandler:        &mock.NodesSetupMock{},
		NodesCoordinatorRegistryFactory: createNodesCoordinatorRegistryFactory(),
	}

	ihnc, err := NewIndexHashedNodesCoordinator(arguments)
	require.Nil(t, err)
	epoch := uint32(1)

	header := &block.MetaBlock{
		PrevRandSeed: []byte("rand seed"),
		EpochStart:   block.EpochStart{LastFinalizedHeaders: []block.EpochStartShardData{{}}},
		Epoch:        epoch,
	}

	ihnc.nodesConfig[epoch] = ihnc.nodesConfig[0]

	body := createBlockBodyFromNodesCoordinator(ihnc, epoch, ihnc.validatorInfoCacher)
	ihnc.EpochStartPrepare(header, body)

	newNodesConfig := ihnc.nodesConfig[1]

	firstEligible := newNodesConfig.eligibleMap[core.MetachainShardId][0]
	secondEligible := newNodesConfig.eligibleMap[core.MetachainShardId][1]
	assert.True(t, firstEligible.Index() < secondEligible.Index())
}

func TestIndexHashedNodesCoordinator_GetConsensusValidatorsPublicKeysNotExistingEpoch(t *testing.T) {
	t.Parallel()

	args := createArguments()
	ihnc, err := NewIndexHashedNodesCoordinator(args)
	require.Nil(t, err)

	var pKeys []string
	var leader string
	randomness := []byte("randomness")
	leader, pKeys, err = ihnc.GetConsensusValidatorsPublicKeys(randomness, 0, 0, 1)
	require.True(t, errors.Is(err, ErrEpochNodesConfigDoesNotExist))
	require.Nil(t, pKeys)
	require.Empty(t, leader)
}

func TestIndexHashedNodesCoordinator_GetConsensusValidatorsPublicKeysExistingEpoch(t *testing.T) {
	t.Parallel()

	args := createArguments()
	ihnc, err := NewIndexHashedNodesCoordinator(args)
	require.Nil(t, err)

	shard0PubKeys := validatorsPubKeys(args.EligibleNodes[0])

	var pKeys []string
	var leader string
	randomness := []byte("randomness")
	leader, pKeys, err = ihnc.GetConsensusValidatorsPublicKeys(randomness, 0, 0, 0)
	require.Nil(t, err)
	require.True(t, len(pKeys) > 0)
	require.True(t, isStringSubgroup(pKeys, shard0PubKeys))
	require.NotEmpty(t, leader)
}

func TestIndexHashedNodesCoordinator_GetValidatorsIndexes(t *testing.T) {
	t.Parallel()

	args := createArguments()
	ihnc, err := NewIndexHashedNodesCoordinator(args)
	require.Nil(t, err)
	randomness := []byte("randomness")

	var pKeys []string
	var leader string
	leader, pKeys, err = ihnc.GetConsensusValidatorsPublicKeys(randomness, 0, 0, 0)
	require.Nil(t, err)

	var indexes []uint64
	indexes, err = ihnc.GetValidatorsIndexes(pKeys, 0)
	require.Nil(t, err)
	require.Equal(t, len(pKeys), len(indexes))
	require.NotEmpty(t, leader)
}

func TestIndexHashedNodesCoordinator_GetValidatorsIndexesInvalidPubKey(t *testing.T) {
	t.Parallel()

	args := createArguments()
	ihnc, err := NewIndexHashedNodesCoordinator(args)
	require.Nil(t, err)
	randomness := []byte("randomness")

	var pKeys []string
	var leader string
	leader, pKeys, err = ihnc.GetConsensusValidatorsPublicKeys(randomness, 0, 0, 0)
	require.Nil(t, err)
	require.NotEmpty(t, leader)

	var indexes []uint64
	pKeys[0] = "dummy"
	indexes, err = ihnc.GetValidatorsIndexes(pKeys, 0)
	require.Equal(t, ErrInvalidNumberPubKeys, err)
	require.Nil(t, indexes)
}

func TestIndexHashedNodesCoordinator_GetSavedStateKey(t *testing.T) {
	t.Parallel()

	args := createArguments()
	args.ValidatorInfoCacher = dataPool.NewCurrentEpochValidatorInfoPool()
	ihnc, err := NewIndexHashedNodesCoordinator(args)
	require.Nil(t, err)

	header := &block.MetaBlock{
		PrevRandSeed: []byte("rand seed"),
		EpochStart:   block.EpochStart{LastFinalizedHeaders: []block.EpochStartShardData{{}}},
		Epoch:        1,
	}

	body := createBlockBodyFromNodesCoordinator(ihnc, 0, ihnc.validatorInfoCacher)
	ihnc.EpochStartPrepare(header, body)
	ihnc.EpochStartAction(header)

	key := ihnc.GetSavedStateKey()
	require.Equal(t, []byte("rand seed"), key)
}

func TestIndexHashedNodesCoordinator_GetSavedStateKeyEpoch0(t *testing.T) {
	t.Parallel()

	args := createArguments()
	ihnc, err := NewIndexHashedNodesCoordinator(args)
	require.Nil(t, err)

	expectedKey := args.Hasher.Compute(string(args.SelfPublicKey))
	key := ihnc.GetSavedStateKey()
	require.Equal(t, expectedKey, key)
}

func TestIndexHashedNodesCoordinator_ShardIdForEpochInvalidEpoch(t *testing.T) {
	t.Parallel()

	args := createArguments()
	ihnc, err := NewIndexHashedNodesCoordinator(args)
	require.Nil(t, err)

	shardId, err := ihnc.ShardIdForEpoch(1)
	require.True(t, errors.Is(err, ErrEpochNodesConfigDoesNotExist))
	require.Equal(t, uint32(0), shardId)
}

func TestIndexHashedNodesCoordinator_ShardIdForEpochValidEpoch(t *testing.T) {
	t.Parallel()

	args := createArguments()
	ihnc, err := NewIndexHashedNodesCoordinator(args)
	require.Nil(t, err)

	shardId, err := ihnc.ShardIdForEpoch(0)
	require.Nil(t, err)
	require.Equal(t, uint32(0), shardId)
}

func TestIndexHashedNodesCoordinator_GetConsensusWhitelistedNodesEpoch0(t *testing.T) {
	t.Parallel()

	args := createArguments()
	ihnc, err := NewIndexHashedNodesCoordinator(args)
	require.Nil(t, err)

	nodesCurrentEpoch, err := ihnc.GetAllEligibleValidatorsPublicKeys(0)
	require.Nil(t, err)

	allNodesList := make([]string, 0)
	for _, nodesList := range nodesCurrentEpoch {
		for _, nodeKey := range nodesList {
			allNodesList = append(allNodesList, string(nodeKey))
		}
	}

	whitelistedNodes, err := ihnc.GetConsensusWhitelistedNodes(0)
	require.Nil(t, err)
	require.Greater(t, len(whitelistedNodes), 0)

	for key := range whitelistedNodes {
		require.True(t, isStringSubgroup([]string{key}, allNodesList))
	}
}

func TestIndexHashedNodesCoordinator_GetConsensusWhitelistedNodesEpoch1(t *testing.T) {
	t.Parallel()

	arguments := createArguments()
	arguments.ValidatorInfoCacher = dataPool.NewCurrentEpochValidatorInfoPool()
	ihnc, err := NewIndexHashedNodesCoordinator(arguments)
	require.Nil(t, err)

	header := &block.MetaBlock{
		PrevRandSeed: []byte("rand seed"),
		EpochStart:   block.EpochStart{LastFinalizedHeaders: []block.EpochStartShardData{{}}},
		Epoch:        1,
	}

	body := createBlockBodyFromNodesCoordinator(ihnc, 0, ihnc.validatorInfoCacher)
	ihnc.EpochStartPrepare(header, body)
	ihnc.EpochStartAction(header)

	nodesPrevEpoch, err := ihnc.GetAllEligibleValidatorsPublicKeys(0)
	require.Nil(t, err)
	nodesCurrentEpoch, err := ihnc.GetAllEligibleValidatorsPublicKeys(1)
	require.Nil(t, err)

	allNodesList := make([]string, 0)
	for shardId := range nodesPrevEpoch {
		for _, nodeKey := range nodesPrevEpoch[shardId] {
			allNodesList = append(allNodesList, string(nodeKey))
		}
		for _, nodeKey := range nodesCurrentEpoch[shardId] {
			allNodesList = append(allNodesList, string(nodeKey))
		}
	}

	whitelistedNodes, err := ihnc.GetConsensusWhitelistedNodes(1)
	require.Nil(t, err)
	require.Greater(t, len(whitelistedNodes), 0)

	for key := range whitelistedNodes {
		require.True(t, isStringSubgroup([]string{key}, allNodesList))
	}
}

func TestIndexHashedNodesCoordinator_GetAllEligibleValidatorsPublicKeysForShard(t *testing.T) {
	t.Parallel()

	t.Run("missing nodes config should error", func(t *testing.T) {
		t.Parallel()

		arguments := createArguments()
		arguments.ValidatorInfoCacher = dataPool.NewCurrentEpochValidatorInfoPool()
		ihnc, err := NewIndexHashedNodesCoordinator(arguments)
		require.Nil(t, err)

		validators, err := ihnc.GetAllEligibleValidatorsPublicKeysForShard(100, 0)
		require.True(t, errors.Is(err, ErrEpochNodesConfigDoesNotExist))
		require.Nil(t, validators)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		arguments := createArguments()
		arguments.ValidatorInfoCacher = dataPool.NewCurrentEpochValidatorInfoPool()
		ihnc, err := NewIndexHashedNodesCoordinator(arguments)
		require.Nil(t, err)

		expectedValidators := make([]string, 0, len(arguments.EligibleNodes[0]))
		for _, val := range arguments.EligibleNodes[0] {
			expectedValidators = append(expectedValidators, string(val.PubKey()))
		}
		validators, err := ihnc.GetAllEligibleValidatorsPublicKeysForShard(0, 0)
		require.NoError(t, err)
		require.Equal(t, expectedValidators, validators)
	})
}

func TestIndexHashedNodesCoordinator_GetConsensusWhitelistedNodesAfterRevertToEpoch(t *testing.T) {
	t.Parallel()

	arguments := createArguments()
	arguments.ValidatorInfoCacher = dataPool.NewCurrentEpochValidatorInfoPool()
	ihnc, err := NewIndexHashedNodesCoordinator(arguments)
	require.Nil(t, err)

	header := &block.MetaBlock{
		PrevRandSeed: []byte("rand seed"),
		EpochStart:   block.EpochStart{LastFinalizedHeaders: []block.EpochStartShardData{{}}},
		Epoch:        1,
	}

	body := createBlockBodyFromNodesCoordinator(ihnc, 0, ihnc.validatorInfoCacher)
	ihnc.EpochStartPrepare(header, body)
	ihnc.EpochStartAction(header)

	body = createBlockBodyFromNodesCoordinator(ihnc, 1, ihnc.validatorInfoCacher)
	header = &block.MetaBlock{
		PrevRandSeed: []byte("rand seed"),
		EpochStart:   block.EpochStart{LastFinalizedHeaders: []block.EpochStartShardData{{}}},
		Epoch:        2,
	}
	ihnc.EpochStartPrepare(header, body)
	ihnc.EpochStartAction(header)

	body = createBlockBodyFromNodesCoordinator(ihnc, 2, ihnc.validatorInfoCacher)
	header = &block.MetaBlock{
		PrevRandSeed: []byte("rand seed"),
		EpochStart:   block.EpochStart{LastFinalizedHeaders: []block.EpochStartShardData{{}}},
		Epoch:        3,
	}
	ihnc.EpochStartPrepare(header, body)
	ihnc.EpochStartAction(header)

	body = createBlockBodyFromNodesCoordinator(ihnc, 3, ihnc.validatorInfoCacher)
	header = &block.MetaBlock{
		PrevRandSeed: []byte("rand seed"),
		EpochStart:   block.EpochStart{LastFinalizedHeaders: []block.EpochStartShardData{{}}},
		Epoch:        4,
	}
	ihnc.EpochStartPrepare(header, body)
	ihnc.EpochStartAction(header)

	nodesEpoch1, err := ihnc.GetAllEligibleValidatorsPublicKeys(1)
	require.Nil(t, err)

	allNodesList := make([]string, 0)
	for _, nodesList := range nodesEpoch1 {
		for _, nodeKey := range nodesList {
			allNodesList = append(allNodesList, string(nodeKey))
		}
	}

	whitelistedNodes, err := ihnc.GetConsensusWhitelistedNodes(1)
	require.Nil(t, err)
	require.Greater(t, len(whitelistedNodes), 0)

	for key := range whitelistedNodes {
		require.True(t, isStringSubgroup([]string{key}, allNodesList))
	}
}

func TestIndexHashedNodesCoordinator_ConsensusGroupSize(t *testing.T) {
	t.Parallel()

	testEpoch := uint32(37)
	shardConsensusGroupSize, metaConsensusGroupSize := 1, 1
	arguments := createArguments()
	arguments.Epoch = testEpoch - 1
	numTimesChainParametersForEpochWasCalled := 0
	arguments.ChainParametersHandler = &chainParameters.ChainParametersHandlerStub{
		CurrentChainParametersCalled: func() config.ChainParametersByEpochConfig {
			return config.ChainParametersByEpochConfig{
				ShardConsensusGroupSize:     uint32(shardConsensusGroupSize),
				MetachainConsensusGroupSize: uint32(metaConsensusGroupSize),
			}
		},
		ChainParametersForEpochCalled: func(epoch uint32) (config.ChainParametersByEpochConfig, error) {
			if numTimesChainParametersForEpochWasCalled == 0 {
				require.Equal(t, testEpoch-1, epoch)
			} else {
				require.Equal(t, testEpoch, epoch)
			}
			numTimesChainParametersForEpochWasCalled++

			return config.ChainParametersByEpochConfig{
				ShardConsensusGroupSize:     1,
				MetachainConsensusGroupSize: 1,
			}, nil
		},
	}
	ihnc, err := NewIndexHashedNodesCoordinator(arguments)
	require.Nil(t, err)

	consensusSizeShard := ihnc.ConsensusGroupSizeForShardAndEpoch(0, testEpoch)
	consensusSizeMeta := ihnc.ConsensusGroupSizeForShardAndEpoch(core.MetachainShardId, testEpoch)

	require.Equal(t, shardConsensusGroupSize, consensusSizeShard)
	require.Equal(t, metaConsensusGroupSize, consensusSizeMeta)

	// consensus group size from chain parameters should have been called once from the constructor, once for shard and once for meta
	require.Equal(t, 3, numTimesChainParametersForEpochWasCalled)
}

func TestIndexHashedNodesCoordinator_GetNumTotalEligible(t *testing.T) {
	t.Parallel()

	arguments := createArguments()
	ihnc, err := NewIndexHashedNodesCoordinator(arguments)
	require.Nil(t, err)

	expectedNbNodes := uint64(0)
	for _, nodesList := range arguments.EligibleNodes {
		expectedNbNodes += uint64(len(nodesList))
	}

	nbNodes := ihnc.GetNumTotalEligible()
	require.Equal(t, expectedNbNodes, nbNodes)
}

func TestIndexHashedNodesCoordinator_GetOwnPublicKey(t *testing.T) {
	t.Parallel()

	arguments := createArguments()
	ihnc, err := NewIndexHashedNodesCoordinator(arguments)
	require.Nil(t, err)

	ownPubKey := ihnc.GetOwnPublicKey()
	require.Equal(t, arguments.SelfPublicKey, ownPubKey)
}

func TestIndexHashedNodesCoordinator_ShuffleOutWithEligible(t *testing.T) {
	t.Parallel()

	processCalled := false
	newShard := uint32(0)

	arguments := createArguments()
	arguments.ShuffledOutHandler = &mock.ShuffledOutHandlerStub{
		ProcessCalled: func(newShardID uint32) error {
			processCalled = true
			newShard = newShardID
			return nil
		},
	}
	pk := []byte("pk")
	arguments.SelfPublicKey = pk
	ihnc, err := NewIndexHashedNodesCoordinator(arguments)
	require.Nil(t, err)

	epoch := uint32(2)
	validatorShard := uint32(7)
	ihnc.nodesConfig = map[uint32]*epochNodesConfig{
		epoch: {
			shardID: validatorShard,
			eligibleMap: map[uint32][]Validator{
				validatorShard: {newValidatorMock(pk, 1, 1)},
			},
		},
	}

	ihnc.ShuffleOutForEpoch(epoch)
	require.True(t, processCalled)
	require.Equal(t, validatorShard, newShard)
}

func TestIndexHashedNodesCoordinator_ShuffleOutWithWaiting(t *testing.T) {
	t.Parallel()

	processCalled := false
	newShard := uint32(0)

	arguments := createArguments()
	arguments.ShuffledOutHandler = &mock.ShuffledOutHandlerStub{
		ProcessCalled: func(newShardID uint32) error {
			processCalled = true
			newShard = newShardID
			return nil
		},
	}
	pk := []byte("pk")
	arguments.SelfPublicKey = pk
	ihnc, err := NewIndexHashedNodesCoordinator(arguments)
	require.Nil(t, err)

	epoch := uint32(2)
	validatorShard := uint32(7)
	ihnc.nodesConfig = map[uint32]*epochNodesConfig{
		epoch: {
			shardID: validatorShard,
			waitingMap: map[uint32][]Validator{
				validatorShard: {newValidatorMock(pk, 1, 1)},
			},
		},
	}

	ihnc.ShuffleOutForEpoch(epoch)
	require.True(t, processCalled)
	require.Equal(t, validatorShard, newShard)
}

func TestIndexHashedNodesCoordinator_ShuffleOutWithObserver(t *testing.T) {
	t.Parallel()

	processCalled := false
	newShard := uint32(0)

	arguments := createArguments()
	arguments.ShuffledOutHandler = &mock.ShuffledOutHandlerStub{
		ProcessCalled: func(newShardID uint32) error {
			processCalled = true
			newShard = newShardID
			return nil
		},
	}
	pk := []byte("pk")
	arguments.SelfPublicKey = pk
	ihnc, err := NewIndexHashedNodesCoordinator(arguments)
	require.Nil(t, err)

	epoch := uint32(2)
	validatorShard := uint32(7)
	ihnc.nodesConfig = map[uint32]*epochNodesConfig{
		epoch: {
			shardID: validatorShard,
			eligibleMap: map[uint32][]Validator{
				validatorShard: {newValidatorMock([]byte("eligibleKey"), 1, 1)},
			},
			waitingMap: map[uint32][]Validator{
				validatorShard: {newValidatorMock([]byte("waitingKey"), 1, 1)},
			},
			leavingMap: map[uint32][]Validator{
				validatorShard: {newValidatorMock(pk, 1, 1)}},
		},
	}

	ihnc.ShuffleOutForEpoch(epoch)
	require.False(t, processCalled)
	expectedShardForLeaving := uint32(0)
	require.Equal(t, expectedShardForLeaving, newShard)
}

func TestIndexHashedNodesCoordinator_ShuffleOutNotFound(t *testing.T) {
	t.Parallel()

	processCalled := false
	newShard := uint32(0)

	arguments := createArguments()
	arguments.ShuffledOutHandler = &mock.ShuffledOutHandlerStub{
		ProcessCalled: func(newShardID uint32) error {
			processCalled = true
			newShard = newShardID
			return nil
		},
	}
	pk := []byte("pk")
	arguments.SelfPublicKey = pk
	ihnc, err := NewIndexHashedNodesCoordinator(arguments)
	require.Nil(t, err)

	epoch := uint32(2)
	validatorShard := uint32(7)
	ihnc.nodesConfig = map[uint32]*epochNodesConfig{
		epoch: {
			shardID: validatorShard,
			eligibleMap: map[uint32][]Validator{
				validatorShard: {newValidatorMock([]byte("eligibleKey"), 1, 1)},
			},
			waitingMap: map[uint32][]Validator{
				validatorShard: {newValidatorMock([]byte("waitingKey"), 1, 1)},
			},
			leavingMap: map[uint32][]Validator{
				validatorShard: {newValidatorMock([]byte("observerKey"), 1, 1)},
			},
		},
	}

	ihnc.ShuffleOutForEpoch(epoch)
	require.False(t, processCalled)
	expectedShardForNotFound := uint32(0)
	require.Equal(t, expectedShardForNotFound, newShard)
}

func TestIndexHashedNodesCoordinator_ShuffleOutNilConfig(t *testing.T) {
	t.Parallel()

	processCalled := false
	newShard := uint32(0)

	arguments := createArguments()
	arguments.ShuffledOutHandler = &mock.ShuffledOutHandlerStub{
		ProcessCalled: func(newShardID uint32) error {
			processCalled = true
			newShard = newShardID
			return nil
		},
	}
	pk := []byte("pk")
	arguments.SelfPublicKey = pk
	ihnc, err := NewIndexHashedNodesCoordinator(arguments)
	require.Nil(t, err)

	epoch := uint32(2)
	ihnc.nodesConfig = map[uint32]*epochNodesConfig{
		epoch: nil,
	}

	ihnc.ShuffleOutForEpoch(epoch)
	require.False(t, processCalled)
	expectedShardForNotFound := uint32(0)
	require.Equal(t, expectedShardForNotFound, newShard)
}

func TestIndexHashedNodesCoordinator_computeNodesConfigFromListNoValidators(t *testing.T) {
	t.Parallel()

	arguments := createArguments()
	pk := []byte("pk")
	arguments.SelfPublicKey = pk
	ihnc, _ := NewIndexHashedNodesCoordinator(arguments)

	validatorInfos := make([]*state.ShardValidatorInfo, 0)
	newNodesConfig, err := ihnc.computeNodesConfigFromList(validatorInfos)

	assert.Nil(t, newNodesConfig)
	assert.True(t, errors.Is(err, ErrMapSizeZero))

	newNodesConfig, err = ihnc.computeNodesConfigFromList(nil)

	assert.Nil(t, newNodesConfig)
	assert.True(t, errors.Is(err, ErrMapSizeZero))
}

func TestIndexHashedNodesCoordinator_computeNodesConfigFromListNilPk(t *testing.T) {
	t.Parallel()

	arguments := createArguments()
	pk := []byte("pk")
	arguments.SelfPublicKey = pk
	ihnc, _ := NewIndexHashedNodesCoordinator(arguments)

	validatorInfos :=
		[]*state.ShardValidatorInfo{
			{
				PublicKey:  pk,
				ShardId:    0,
				List:       "test1",
				Index:      0,
				TempRating: 0,
			},
			{
				PublicKey:  nil,
				ShardId:    0,
				List:       "test",
				Index:      0,
				TempRating: 0,
			},
		}

	newNodesConfig, err := ihnc.computeNodesConfigFromList(validatorInfos)

	assert.Nil(t, newNodesConfig)
	assert.NotNil(t, err)
	assert.Equal(t, ErrNilPubKey, err)
}

func TestIndexHashedNodesCoordinator_computeNodesConfigFromListWithStakingV4(t *testing.T) {
	t.Parallel()
	arguments := createArguments()
	nc, _ := NewIndexHashedNodesCoordinator(arguments)

	shard0Eligible := &state.ShardValidatorInfo{
		PublicKey:  []byte("pk0"),
		List:       string(common.EligibleList),
		Index:      1,
		TempRating: 2,
		ShardId:    0,
	}
	shard0Auction := &state.ShardValidatorInfo{
		PublicKey:  []byte("pk1"),
		List:       string(common.SelectedFromAuctionList),
		Index:      3,
		TempRating: 2,
		ShardId:    0,
	}
	shard1Auction := &state.ShardValidatorInfo{
		PublicKey:  []byte("pk2"),
		List:       string(common.SelectedFromAuctionList),
		Index:      2,
		TempRating: 2,
		ShardId:    1,
	}
	validatorInfos := []*state.ShardValidatorInfo{shard0Eligible, shard0Auction, shard1Auction}

	newNodesConfig, err := nc.computeNodesConfigFromList(validatorInfos)
	require.Equal(t, ErrReceivedAuctionValidatorsBeforeStakingV4, err)
	require.Nil(t, newNodesConfig)

	nc.updateEpochFlags(stakingV4Epoch)

	newNodesConfig, err = nc.computeNodesConfigFromList(validatorInfos)
	require.Nil(t, err)
	v1, _ := NewValidator([]byte("pk2"), 1, 2)
	v2, _ := NewValidator([]byte("pk1"), 1, 3)
	require.Equal(t, []Validator{v1, v2}, newNodesConfig.auctionList)

	validatorInfos = append(validatorInfos, &state.ShardValidatorInfo{
		PublicKey: []byte("pk3"),
		List:      string(common.NewList),
	})
	newNodesConfig, err = nc.computeNodesConfigFromList(validatorInfos)
	require.Equal(t, epochStart.ErrReceivedNewListNodeInStakingV4, err)
	require.Nil(t, newNodesConfig)
}

func TestIndexHashedNodesCoordinator_computeNodesConfigFromListValidatorsWithFix(t *testing.T) {
	t.Parallel()

	arguments := createArguments()
	pk := []byte("pk")
	arguments.SelfPublicKey = pk
	ihnc, _ := NewIndexHashedNodesCoordinator(arguments)
	_ = ihnc.flagStakingV4Started.SetReturningPrevious()

	shard0Eligible0 := &state.ShardValidatorInfo{
		PublicKey:  []byte("pk0"),
		List:       string(common.EligibleList),
		Index:      1,
		TempRating: 2,
		ShardId:    0,
	}
	shard0Eligible1 := &state.ShardValidatorInfo{
		PublicKey:  []byte("pk1"),
		List:       string(common.EligibleList),
		Index:      2,
		TempRating: 2,
		ShardId:    0,
	}
	shardmetaEligible0 := &state.ShardValidatorInfo{
		PublicKey:  []byte("pk2"),
		ShardId:    core.MetachainShardId,
		List:       string(common.EligibleList),
		Index:      1,
		TempRating: 4,
	}
	shard0Waiting0 := &state.ShardValidatorInfo{
		PublicKey: []byte("pk3"),
		List:      string(common.WaitingList),
		Index:     14,
		ShardId:   0,
	}
	shardmetaWaiting0 := &state.ShardValidatorInfo{
		PublicKey: []byte("pk4"),
		ShardId:   core.MetachainShardId,
		List:      string(common.WaitingList),
		Index:     15,
	}
	shard0New0 := &state.ShardValidatorInfo{
		PublicKey: []byte("pk5"),
		List:      string(common.NewList), Index: 3,
		ShardId: 0,
	}
	shard0Leaving0 := &state.ShardValidatorInfo{
		PublicKey:    []byte("pk6"),
		List:         string(common.LeavingList),
		PreviousList: string(common.EligibleList),
		ShardId:      0,
	}
	shardMetaLeaving1 := &state.ShardValidatorInfo{
		PublicKey:     []byte("pk7"),
		List:          string(common.LeavingList),
		PreviousList:  string(common.WaitingList),
		Index:         1,
		PreviousIndex: 1,
		ShardId:       core.MetachainShardId,
	}

	validatorInfos :=
		[]*state.ShardValidatorInfo{
			shard0Eligible0,
			shard0Eligible1,
			shardmetaEligible0,
			shard0Waiting0,
			shardmetaWaiting0,
			shard0New0,
			shard0Leaving0,
			shardMetaLeaving1,
		}

	newNodesConfig, err := ihnc.computeNodesConfigFromList(validatorInfos)
	assert.Nil(t, err)

	assert.Equal(t, uint32(1), newNodesConfig.nbShards)

	verifySizes(t, newNodesConfig)
	verifyLeavingNodesInEligibleOrWaiting(t, newNodesConfig)

	// maps have the correct validators inside
	eligibleListShardZero := createValidatorList(ihnc,
		[]*state.ShardValidatorInfo{shard0Eligible0, shard0Eligible1, shard0Leaving0})
	assert.Equal(t, eligibleListShardZero, newNodesConfig.eligibleMap[0])
	eligibleListMeta := createValidatorList(ihnc,
		[]*state.ShardValidatorInfo{shardmetaEligible0})
	assert.Equal(t, eligibleListMeta, newNodesConfig.eligibleMap[core.MetachainShardId])

	waitingListShardZero := createValidatorList(ihnc,
		[]*state.ShardValidatorInfo{shard0Waiting0})
	assert.Equal(t, waitingListShardZero, newNodesConfig.waitingMap[0])
	waitingListMeta := createValidatorList(ihnc,
		[]*state.ShardValidatorInfo{shardmetaWaiting0, shardMetaLeaving1})
	assert.Equal(t, waitingListMeta, newNodesConfig.waitingMap[core.MetachainShardId])

	leavingListShardZero := createValidatorList(ihnc,
		[]*state.ShardValidatorInfo{shard0Leaving0})
	assert.Equal(t, leavingListShardZero, newNodesConfig.leavingMap[0])

	leavingListMeta := createValidatorList(ihnc,
		[]*state.ShardValidatorInfo{shardMetaLeaving1})
	assert.Equal(t, leavingListMeta, newNodesConfig.leavingMap[core.MetachainShardId])

	newListShardZero := createValidatorList(ihnc,
		[]*state.ShardValidatorInfo{shard0New0})
	assert.Equal(t, newListShardZero, newNodesConfig.newList)
}

func TestIndexHashedNodesCoordinator_computeNodesConfigFromListValidatorsNoFix(t *testing.T) {
	t.Parallel()

	arguments := createArguments()
	pk := []byte("pk")
	arguments.SelfPublicKey = pk
	ihnc, _ := NewIndexHashedNodesCoordinator(arguments)

	shard0Eligible0 := &state.ShardValidatorInfo{
		PublicKey:  []byte("pk0"),
		List:       string(common.EligibleList),
		Index:      1,
		TempRating: 2,
		ShardId:    0,
	}
	shard0Eligible1 := &state.ShardValidatorInfo{
		PublicKey:  []byte("pk1"),
		List:       string(common.EligibleList),
		Index:      2,
		TempRating: 2,
		ShardId:    0,
	}
	shardmetaEligible0 := &state.ShardValidatorInfo{
		PublicKey:  []byte("pk2"),
		ShardId:    core.MetachainShardId,
		List:       string(common.EligibleList),
		Index:      1,
		TempRating: 4,
	}
	shard0Waiting0 := &state.ShardValidatorInfo{
		PublicKey: []byte("pk3"),
		List:      string(common.WaitingList),
		Index:     14,
		ShardId:   0,
	}
	shardmetaWaiting0 := &state.ShardValidatorInfo{
		PublicKey: []byte("pk4"),
		ShardId:   core.MetachainShardId,
		List:      string(common.WaitingList),
		Index:     15,
	}
	shard0New0 := &state.ShardValidatorInfo{
		PublicKey: []byte("pk5"),
		List:      string(common.NewList), Index: 3,
		ShardId: 0,
	}
	shard0Leaving0 := &state.ShardValidatorInfo{
		PublicKey: []byte("pk6"),
		List:      string(common.LeavingList),
		ShardId:   0,
	}
	shardMetaLeaving1 := &state.ShardValidatorInfo{
		PublicKey: []byte("pk7"),
		List:      string(common.LeavingList),
		Index:     1,
		ShardId:   core.MetachainShardId,
	}

	validatorInfos :=
		[]*state.ShardValidatorInfo{
			shard0Eligible0,
			shard0Eligible1,
			shardmetaEligible0,
			shard0Waiting0,
			shardmetaWaiting0,
			shard0New0,
			shard0Leaving0,
			shardMetaLeaving1,
		}

	ihnc.flagStakingV4Started.Reset()
	newNodesConfig, err := ihnc.computeNodesConfigFromList(validatorInfos)
	assert.Nil(t, err)

	assert.Equal(t, uint32(1), newNodesConfig.nbShards)

	verifySizes(t, newNodesConfig)
	verifyLeavingNodesInEligible(t, newNodesConfig)

	// maps have the correct validators inside
	eligibleListShardZero := createValidatorList(ihnc,
		[]*state.ShardValidatorInfo{shard0Eligible0, shard0Eligible1, shard0Leaving0})
	assert.Equal(t, eligibleListShardZero, newNodesConfig.eligibleMap[0])
	eligibleListMeta := createValidatorList(ihnc,
		[]*state.ShardValidatorInfo{shardmetaEligible0, shardMetaLeaving1})
	assert.Equal(t, eligibleListMeta, newNodesConfig.eligibleMap[core.MetachainShardId])

	waitingListShardZero := createValidatorList(ihnc,
		[]*state.ShardValidatorInfo{shard0Waiting0})
	assert.Equal(t, waitingListShardZero, newNodesConfig.waitingMap[0])
	waitingListMeta := createValidatorList(ihnc,
		[]*state.ShardValidatorInfo{shardmetaWaiting0})
	assert.Equal(t, waitingListMeta, newNodesConfig.waitingMap[core.MetachainShardId])

	leavingListShardZero := createValidatorList(ihnc,
		[]*state.ShardValidatorInfo{shard0Leaving0})
	assert.Equal(t, leavingListShardZero, newNodesConfig.leavingMap[0])

	leavingListMeta := createValidatorList(ihnc,
		[]*state.ShardValidatorInfo{shardMetaLeaving1})
	assert.Equal(t, leavingListMeta, newNodesConfig.leavingMap[core.MetachainShardId])

	newListShardZero := createValidatorList(ihnc,
		[]*state.ShardValidatorInfo{shard0New0})
	assert.Equal(t, newListShardZero, newNodesConfig.newList)
}

func createValidatorList(ihnc *indexHashedNodesCoordinator, shardValidators []*state.ShardValidatorInfo) []Validator {
	validators := make([]Validator, len(shardValidators))
	for i, v := range shardValidators {
		shardValidator, _ := NewValidator(
			v.PublicKey,
			ihnc.GetChance(v.TempRating),
			v.Index)
		validators[i] = shardValidator
	}
	sort.Sort(validatorList(validators))
	return validators
}

func verifyLeavingNodesInEligible(t *testing.T, newNodesConfig *epochNodesConfig) {
	for leavingShardId, leavingValidators := range newNodesConfig.leavingMap {
		for _, leavingValidator := range leavingValidators {
			found, shardId := searchInMap(newNodesConfig.eligibleMap, leavingValidator.PubKey())
			assert.True(t, found)
			assert.Equal(t, leavingShardId, shardId)
		}
	}
}

func verifyLeavingNodesInEligibleOrWaiting(t *testing.T, newNodesConfig *epochNodesConfig) {
	for leavingShardId, leavingValidators := range newNodesConfig.leavingMap {
		for _, leavingValidator := range leavingValidators {
			found, shardId := searchInMap(newNodesConfig.eligibleMap, leavingValidator.PubKey())
			if !found {
				found, shardId = searchInMap(newNodesConfig.waitingMap, leavingValidator.PubKey())
			}
			assert.True(t, found)
			assert.Equal(t, leavingShardId, shardId)
		}
	}
}

func verifySizes(t *testing.T, newNodesConfig *epochNodesConfig) {
	expectedEligibleSize := 2
	expectedWaitingSize := 2
	expectedNewSize := 1
	expectedLeavingSize := 2

	assert.NotNil(t, newNodesConfig)
	assert.Equal(t, uint32(expectedEligibleSize-1), newNodesConfig.nbShards)
	assert.Equal(t, expectedEligibleSize, len(newNodesConfig.eligibleMap))
	assert.Equal(t, expectedWaitingSize, len(newNodesConfig.waitingMap))
	assert.Equal(t, expectedNewSize, len(newNodesConfig.newList))
	assert.Equal(t, expectedLeavingSize, len(newNodesConfig.leavingMap))
}

func TestIndexHashedNodesCoordinator_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	var ihnc NodesCoordinator
	require.True(t, check.IfNil(ihnc))

	var ihnc2 *indexHashedNodesCoordinator
	require.True(t, check.IfNil(ihnc2))

	arguments := createArguments()
	ihnc3, err := NewIndexHashedNodesCoordinator(arguments)
	require.Nil(t, err)
	require.False(t, check.IfNil(ihnc3))
}

func TestIndexHashedNodesCoordinator_GetShardValidatorInfoData(t *testing.T) {
	t.Parallel()

	t.Run("get shard validator info data before refactor peers mini block activation flag is set", func(t *testing.T) {
		t.Parallel()

		txHash := []byte("txHash")
		svi := &state.ShardValidatorInfo{PublicKey: []byte("x")}

		arguments := createArguments()
		arguments.EnableEpochsHandler = &enableEpochsHandlerMock.EnableEpochsHandlerStub{
			IsFlagEnabledInEpochCalled: func(flag core.EnableEpochFlag, epoch uint32) bool {
				if flag == common.RefactorPeersMiniBlocksFlag {
					return epoch >= 1
				}
				return false
			},
		}
		arguments.ValidatorInfoCacher = &vic.ValidatorInfoCacherStub{
			GetValidatorInfoCalled: func(validatorInfoHash []byte) (*state.ShardValidatorInfo, error) {
				if bytes.Equal(validatorInfoHash, txHash) {
					return svi, nil
				}
				return nil, errors.New("error")
			},
		}
		ihnc, _ := NewIndexHashedNodesCoordinator(arguments)

		marshalledSVI, _ := arguments.Marshalizer.Marshal(svi)
		shardValidatorInfo, _ := ihnc.getShardValidatorInfoData(marshalledSVI, 0)
		require.Equal(t, svi, shardValidatorInfo)
	})

	t.Run("get shard validator info data after refactor peers mini block activation flag is set", func(t *testing.T) {
		t.Parallel()

		txHash := []byte("txHash")
		svi := &state.ShardValidatorInfo{PublicKey: []byte("x")}

		arguments := createArguments()
		arguments.ValidatorInfoCacher = &vic.ValidatorInfoCacherStub{
			GetValidatorInfoCalled: func(validatorInfoHash []byte) (*state.ShardValidatorInfo, error) {
				if bytes.Equal(validatorInfoHash, txHash) {
					return svi, nil
				}
				return nil, errors.New("error")
			},
		}
		ihnc, _ := NewIndexHashedNodesCoordinator(arguments)

		shardValidatorInfo, _ := ihnc.getShardValidatorInfoData(txHash, 0)
		require.Equal(t, svi, shardValidatorInfo)
	})
}

func TestIndexHashedGroupSelector_GetWaitingEpochsLeftForPublicKey(t *testing.T) {
	t.Parallel()

	t.Run("missing nodes config for current epoch should error ", func(t *testing.T) {
		t.Parallel()

		epochStartSubscriber := &testscommonConsensus.EpochStartNotifierStub{}
		bootStorer := genericMocks.NewStorerMock()

		shufflerArgs := &NodesShufflerArgs{
			ShuffleBetweenShards: shuffleBetweenShards,
			MaxNodesEnableConfig: nil,
			EnableEpochsHandler:  &mock.EnableEpochsHandlerMock{},
		}
		nodeShuffler, err := NewHashValidatorsShuffler(shufflerArgs)
		require.Nil(t, err)

		arguments := ArgNodesCoordinator{
			ChainParametersHandler: &chainParameters.ChainParametersHandlerStub{
				CurrentChainParametersCalled: func() config.ChainParametersByEpochConfig {
					return config.ChainParametersByEpochConfig{
						ShardConsensusGroupSize:     1,
						MetachainConsensusGroupSize: 1,
					}
				},
			},
			Marshalizer:        &mock.MarshalizerMock{},
			Hasher:             &hashingMocks.HasherMock{},
			Shuffler:           nodeShuffler,
			EpochStartNotifier: epochStartSubscriber,
			BootStorer:         bootStorer,
			ShardIDAsObserver:  0,
			NbShards:           2,
			EligibleNodes: map[uint32][]Validator{
				core.MetachainShardId: {newValidatorMock([]byte("pk"), 1, 0)},
			},
			WaitingNodes:        make(map[uint32][]Validator),
			SelfPublicKey:       []byte("key"),
			ConsensusGroupCache: &mock.NodesCoordinatorCacheMock{},
			ShuffledOutHandler:  &mock.ShuffledOutHandlerStub{},
			ChanStopNode:        make(chan endProcess.ArgEndProcess),
			NodeTypeProvider:    &nodeTypeProviderMock.NodeTypeProviderStub{},
			EnableEpochsHandler: &mock.EnableEpochsHandlerMock{
				CurrentEpoch: 1,
			},
			ValidatorInfoCacher:             &vic.ValidatorInfoCacherStub{},
			GenesisNodesSetupHandler:        &mock.NodesSetupMock{},
			NodesCoordinatorRegistryFactory: createNodesCoordinatorRegistryFactory(),
		}

		ihnc, _ := NewIndexHashedNodesCoordinator(arguments)

		epochsLeft, err := ihnc.GetWaitingEpochsLeftForPublicKey([]byte("pk"))
		require.True(t, errors.Is(err, ErrEpochNodesConfigDoesNotExist))
		require.Equal(t, uint32(0), epochsLeft)
	})
	t.Run("min hysteresis nodes returns 0 should work", func(t *testing.T) {
		t.Parallel()

		shardZeroId := uint32(0)
		expectedValidatorsPubKeys := map[uint32][][]byte{
			shardZeroId:           {[]byte("pk0_shard0")},
			core.MetachainShardId: {[]byte("pk0_meta")},
		}

		listMeta := []Validator{
			newValidatorMock(expectedValidatorsPubKeys[core.MetachainShardId][0], 1, defaultSelectionChances),
		}
		listShard0 := []Validator{
			newValidatorMock(expectedValidatorsPubKeys[shardZeroId][0], 1, defaultSelectionChances),
		}

		waitingMap := make(map[uint32][]Validator)
		waitingMap[core.MetachainShardId] = listMeta
		waitingMap[shardZeroId] = listShard0

		epochStartSubscriber := &testscommonConsensus.EpochStartNotifierStub{}
		bootStorer := genericMocks.NewStorerMock()

		eligibleMap := make(map[uint32][]Validator)
		eligibleMap[core.MetachainShardId] = []Validator{&validator{}}
		eligibleMap[shardZeroId] = []Validator{&validator{}}

		shufflerArgs := &NodesShufflerArgs{
			ShuffleBetweenShards: shuffleBetweenShards,
			MaxNodesEnableConfig: nil,
			EnableEpochsHandler:  &mock.EnableEpochsHandlerMock{},
		}
		nodeShuffler, err := NewHashValidatorsShuffler(shufflerArgs)
		require.Nil(t, err)

		arguments := ArgNodesCoordinator{
			ChainParametersHandler: &chainParameters.ChainParametersHandlerStub{
				CurrentChainParametersCalled: func() config.ChainParametersByEpochConfig {
					return config.ChainParametersByEpochConfig{
						ShardConsensusGroupSize:     1,
						MetachainConsensusGroupSize: 1,
					}
				},
			},
			Marshalizer:         &mock.MarshalizerMock{},
			Hasher:              &hashingMocks.HasherMock{},
			Shuffler:            nodeShuffler,
			EpochStartNotifier:  epochStartSubscriber,
			BootStorer:          bootStorer,
			ShardIDAsObserver:   shardZeroId,
			NbShards:            2,
			EligibleNodes:       eligibleMap,
			WaitingNodes:        waitingMap,
			SelfPublicKey:       []byte("key"),
			ConsensusGroupCache: &mock.NodesCoordinatorCacheMock{},
			ShuffledOutHandler:  &mock.ShuffledOutHandlerStub{},
			ChanStopNode:        make(chan endProcess.ArgEndProcess),
			NodeTypeProvider:    &nodeTypeProviderMock.NodeTypeProviderStub{},
			EnableEpochsHandler: &mock.EnableEpochsHandlerMock{},
			ValidatorInfoCacher: &vic.ValidatorInfoCacherStub{},
			GenesisNodesSetupHandler: &mock.NodesSetupMock{
				MinShardHysteresisNodesCalled: func() uint32 {
					return 0
				},
				MinMetaHysteresisNodesCalled: func() uint32 {
					return 0
				},
			},
			NodesCoordinatorRegistryFactory: createNodesCoordinatorRegistryFactory(),
		}

		ihnc, _ := NewIndexHashedNodesCoordinator(arguments)

		epochsLeft, err := ihnc.GetWaitingEpochsLeftForPublicKey([]byte("pk0_shard0"))
		require.NoError(t, err)
		require.Equal(t, uint32(1), epochsLeft)

		epochsLeft, err = ihnc.GetWaitingEpochsLeftForPublicKey([]byte("pk0_meta"))
		require.NoError(t, err)
		require.Equal(t, uint32(1), epochsLeft)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		shardZeroId := uint32(0)
		expectedValidatorsPubKeys := map[uint32][][]byte{
			shardZeroId:           {[]byte("pk0_shard0"), []byte("pk1_shard0"), []byte("pk2_shard0")},
			core.MetachainShardId: {[]byte("pk0_meta"), []byte("pk1_meta"), []byte("pk2_meta"), []byte("pk3_meta"), []byte("pk4_meta")},
		}

		listMeta := []Validator{
			newValidatorMock(expectedValidatorsPubKeys[core.MetachainShardId][0], 1, defaultSelectionChances),
			newValidatorMock(expectedValidatorsPubKeys[core.MetachainShardId][1], 1, defaultSelectionChances),
			newValidatorMock(expectedValidatorsPubKeys[core.MetachainShardId][2], 1, defaultSelectionChances),
			newValidatorMock(expectedValidatorsPubKeys[core.MetachainShardId][3], 1, defaultSelectionChances),
			newValidatorMock(expectedValidatorsPubKeys[core.MetachainShardId][4], 1, defaultSelectionChances),
		}
		listShard0 := []Validator{
			newValidatorMock(expectedValidatorsPubKeys[shardZeroId][0], 1, defaultSelectionChances),
			newValidatorMock(expectedValidatorsPubKeys[shardZeroId][1], 1, defaultSelectionChances),
			newValidatorMock(expectedValidatorsPubKeys[shardZeroId][2], 1, defaultSelectionChances),
		}

		waitingMap := make(map[uint32][]Validator)
		waitingMap[core.MetachainShardId] = listMeta
		waitingMap[shardZeroId] = listShard0

		epochStartSubscriber := &testscommonConsensus.EpochStartNotifierStub{}
		bootStorer := genericMocks.NewStorerMock()

		eligibleMap := make(map[uint32][]Validator)
		eligibleMap[core.MetachainShardId] = []Validator{&validator{}}
		eligibleMap[shardZeroId] = []Validator{&validator{}}

		shufflerArgs := &NodesShufflerArgs{
			ShuffleBetweenShards: shuffleBetweenShards,
			MaxNodesEnableConfig: nil,
			EnableEpochsHandler:  &mock.EnableEpochsHandlerMock{},
		}
		nodeShuffler, err := NewHashValidatorsShuffler(shufflerArgs)
		require.Nil(t, err)

		arguments := ArgNodesCoordinator{
			ChainParametersHandler: &chainParameters.ChainParametersHandlerStub{
				CurrentChainParametersCalled: func() config.ChainParametersByEpochConfig {
					return config.ChainParametersByEpochConfig{
						ShardConsensusGroupSize:     1,
						MetachainConsensusGroupSize: 1,
					}
				},
			},
			Marshalizer:         &mock.MarshalizerMock{},
			Hasher:              &hashingMocks.HasherMock{},
			Shuffler:            nodeShuffler,
			EpochStartNotifier:  epochStartSubscriber,
			BootStorer:          bootStorer,
			ShardIDAsObserver:   shardZeroId,
			NbShards:            2,
			EligibleNodes:       eligibleMap,
			WaitingNodes:        waitingMap,
			SelfPublicKey:       []byte("key"),
			ConsensusGroupCache: &mock.NodesCoordinatorCacheMock{},
			ShuffledOutHandler:  &mock.ShuffledOutHandlerStub{},
			ChanStopNode:        make(chan endProcess.ArgEndProcess),
			NodeTypeProvider:    &nodeTypeProviderMock.NodeTypeProviderStub{},
			EnableEpochsHandler: &mock.EnableEpochsHandlerMock{},
			ValidatorInfoCacher: &vic.ValidatorInfoCacherStub{},
			GenesisNodesSetupHandler: &mock.NodesSetupMock{
				MinShardHysteresisNodesCalled: func() uint32 {
					return 2
				},
				MinMetaHysteresisNodesCalled: func() uint32 {
					return 2
				},
			},
			NodesCoordinatorRegistryFactory: createNodesCoordinatorRegistryFactory(),
		}

		ihnc, _ := NewIndexHashedNodesCoordinator(arguments)

		epochsLeft, err := ihnc.GetWaitingEpochsLeftForPublicKey(nil)
		require.Equal(t, ErrNilPubKey, err)
		require.Zero(t, epochsLeft)

		epochsLeft, err = ihnc.GetWaitingEpochsLeftForPublicKey([]byte("missing_pk"))
		require.Equal(t, ErrKeyNotFoundInWaitingList, err)
		require.Zero(t, epochsLeft)

		epochsLeft, err = ihnc.GetWaitingEpochsLeftForPublicKey([]byte("pk0_shard0"))
		require.NoError(t, err)
		require.Equal(t, uint32(1), epochsLeft)

		epochsLeft, err = ihnc.GetWaitingEpochsLeftForPublicKey([]byte("pk1_shard0"))
		require.NoError(t, err)
		require.Equal(t, uint32(1), epochsLeft)

		epochsLeft, err = ihnc.GetWaitingEpochsLeftForPublicKey([]byte("pk2_shard0"))
		require.NoError(t, err)
		require.Equal(t, uint32(2), epochsLeft)

		epochsLeft, err = ihnc.GetWaitingEpochsLeftForPublicKey([]byte("pk0_meta"))
		require.NoError(t, err)
		require.Equal(t, uint32(1), epochsLeft)

		epochsLeft, err = ihnc.GetWaitingEpochsLeftForPublicKey([]byte("pk1_meta"))
		require.NoError(t, err)
		require.Equal(t, uint32(1), epochsLeft)

		epochsLeft, err = ihnc.GetWaitingEpochsLeftForPublicKey([]byte("pk2_meta"))
		require.NoError(t, err)
		require.Equal(t, uint32(2), epochsLeft)

		epochsLeft, err = ihnc.GetWaitingEpochsLeftForPublicKey([]byte("pk3_meta"))
		require.NoError(t, err)
		require.Equal(t, uint32(2), epochsLeft)

		epochsLeft, err = ihnc.GetWaitingEpochsLeftForPublicKey([]byte("pk4_meta"))
		require.NoError(t, err)
		require.Equal(t, uint32(3), epochsLeft)
	})
}

func TestNodesCoordinator_CustomConsensusGroupSize(t *testing.T) {
	arguments := createArguments()
	numShards := uint32(2)
	nodesPerShard := uint32(3)
	eligibleMap := createDummyNodesMap(nodesPerShard, numShards, "eligible")
	waitingMap := createDummyNodesMap(0, numShards, "waiting")
	arguments.EligibleNodes = eligibleMap
	arguments.WaitingNodes = waitingMap
	arguments.ValidatorInfoCacher = dataPool.NewCurrentEpochValidatorInfoPool()

	consensusParams := []struct {
		enableEpoch   uint32
		shardCnsSize  uint32
		metaCnsSize   uint32
		shardMinNodes uint32
		metaMinNodes  uint32
	}{
		{
			enableEpoch:   9,
			shardCnsSize:  3,
			shardMinNodes: 3,
			metaCnsSize:   3,
			metaMinNodes:  3,
		},
		{
			enableEpoch:   6,
			shardCnsSize:  3,
			shardMinNodes: 3,
			metaCnsSize:   2,
			metaMinNodes:  2,
		},
		{
			enableEpoch:   3,
			shardCnsSize:  3,
			shardMinNodes: 3,
			metaCnsSize:   3,
			metaMinNodes:  3,
		},
		{
			enableEpoch:   0,
			shardCnsSize:  2,
			shardMinNodes: 2,
			metaCnsSize:   3,
			metaMinNodes:  3,
		},
	}
	arguments.ChainParametersHandler = &chainParameters.ChainParametersHandlerStub{
		ChainParametersForEpochCalled: func(epoch uint32) (config.ChainParametersByEpochConfig, error) {
			for _, cfg := range consensusParams {
				if epoch >= cfg.enableEpoch {
					return config.ChainParametersByEpochConfig{
						ShardConsensusGroupSize:     cfg.shardCnsSize,
						ShardMinNumNodes:            cfg.shardMinNodes,
						MetachainConsensusGroupSize: cfg.metaCnsSize,
						MetachainMinNumNodes:        cfg.metaMinNodes,
					}, nil
				}
			}

			return config.ChainParametersByEpochConfig{}, errors.New("wrong test setup")
		},
	}

	shufflerArgs := &NodesShufflerArgs{
		ShuffleBetweenShards: shuffleBetweenShards,
		EnableEpochsHandler:  &mock.EnableEpochsHandlerMock{},
		MaxNodesEnableConfig: []config.MaxNodesChangeConfig{
			{EpochEnable: 0, MaxNumNodes: nodesPerShard * (numShards + 1), NodesToShufflePerShard: 2},
			{EpochEnable: 3, MaxNumNodes: nodesPerShard * (numShards + 1), NodesToShufflePerShard: 3},
		},
	}
	arguments.Shuffler, _ = NewHashValidatorsShuffler(shufflerArgs)

	ihnc, _ := NewIndexHashedNodesCoordinator(arguments)
	require.NotNil(t, ihnc)

	numEpochsToCheck := uint32(100)
	checksCounter := 0
	for ep := uint32(0); ep < numEpochsToCheck; ep++ {
		for _, cfg := range consensusParams {
			if ep >= cfg.enableEpoch {
				changeEpochAndTestNewConsensusSizes(&consensusSizeChangeTestArgs{
					t:                     t,
					ihnc:                  ihnc,
					epoch:                 ep,
					expectedShardMinNodes: cfg.shardMinNodes,
					expectedMetaMinNodes:  cfg.metaMinNodes,
				})
				checksCounter++
				break
			}
		}
	}
	require.Equal(t, numEpochsToCheck, uint32(checksCounter))
}

func TestIndexHashedNodesCoordinator_cacheConsensusGroup(t *testing.T) {
	t.Parallel()

	maxNumValuesCache := 3
	key := []byte("key")

	leader := &validator{
		pubKey:  []byte("leader"),
		chances: 10,
		index:   20,
	}
	validator1 := &validator{
		pubKey:  []byte("validator1"),
		chances: 10,
		index:   20,
	}

	t.Run("adding a key should work", func(t *testing.T) {
		t.Parallel()

		arguments := createArguments()

		arguments.ConsensusGroupCache, _ = cache.NewLRUCache(maxNumValuesCache)
		nodesCoordinator, err := NewIndexHashedNodesCoordinator(arguments)
		require.Nil(t, err)

		consensusGroup := []Validator{leader, validator1}
		expectedData := &savedConsensusGroup{
			leader:         leader,
			consensusGroup: consensusGroup,
		}

		nodesCoordinator.cacheConsensusGroup(key, consensusGroup, leader)
		value := nodesCoordinator.searchConsensusForKey(key)

		require.NotNil(t, value)
		require.Equal(t, expectedData, value)
	})

	t.Run("adding a key twice should overwrite the value", func(t *testing.T) {
		t.Parallel()

		arguments := createArguments()

		arguments.ConsensusGroupCache, _ = cache.NewLRUCache(maxNumValuesCache)
		nodesCoordinator, err := NewIndexHashedNodesCoordinator(arguments)
		require.Nil(t, err)

		cg1 := []Validator{leader, validator1}
		cg2 := []Validator{leader}
		expectedData := &savedConsensusGroup{
			leader:         leader,
			consensusGroup: cg2,
		}

		nodesCoordinator.cacheConsensusGroup(key, cg1, leader)
		nodesCoordinator.cacheConsensusGroup(key, cg2, leader)
		value := nodesCoordinator.searchConsensusForKey(key)
		require.NotNil(t, value)
		require.Equal(t, expectedData, value)
	})

	t.Run("adding more keys than the cache size should remove the oldest key", func(t *testing.T) {
		t.Parallel()

		arguments := createArguments()

		key1 := []byte("key1")
		key2 := []byte("key2")
		key3 := []byte("key3")
		key4 := []byte("key4")

		cg1 := []Validator{leader, validator1}
		cg2 := []Validator{leader}
		cg3 := []Validator{validator1}
		cg4 := []Validator{leader, validator1, validator1}

		arguments.ConsensusGroupCache, _ = cache.NewLRUCache(maxNumValuesCache)
		nodesCoordinator, err := NewIndexHashedNodesCoordinator(arguments)
		require.Nil(t, err)

		nodesCoordinator.cacheConsensusGroup(key1, cg1, leader)
		nodesCoordinator.cacheConsensusGroup(key2, cg2, leader)
		nodesCoordinator.cacheConsensusGroup(key3, cg3, leader)
		nodesCoordinator.cacheConsensusGroup(key4, cg4, leader)

		value := nodesCoordinator.searchConsensusForKey(key1)
		require.Nil(t, value)

		value = nodesCoordinator.searchConsensusForKey(key2)
		require.Equal(t, cg2, value.consensusGroup)

		value = nodesCoordinator.searchConsensusForKey(key3)
		require.Equal(t, cg3, value.consensusGroup)

		value = nodesCoordinator.searchConsensusForKey(key4)
		require.Equal(t, cg4, value.consensusGroup)
	})
}

func TestIndexHashedNodesCoordinator_selectLeaderAndConsensusGroup(t *testing.T) {
	t.Parallel()

	validator1 := &validator{pubKey: []byte("validator1")}
	validator2 := &validator{pubKey: []byte("validator2")}
	validator3 := &validator{pubKey: []byte("validator3")}
	validator4 := &validator{pubKey: []byte("validator4")}

	randomness := []byte("randomness")
	epoch := uint32(1)

	eligibleList := []Validator{validator1, validator2, validator3, validator4}
	consensusSize := len(eligibleList)
	expectedError := errors.New("expected error")
	selectFunc := func(randSeed []byte, sampleSize uint32) ([]uint32, error) {
		if len(eligibleList) < int(sampleSize) {
			return nil, expectedError
		}

		result := make([]uint32, sampleSize)
		for i := 0; i < int(sampleSize); i++ {
			// reverse order from eligible list
			result[i] = uint32(len(eligibleList) - 1 - i)
		}

		return result, nil
	}
	expectedConsensusFixedOrder := []Validator{validator1, validator2, validator3, validator4}
	expectedConsensusNotFixedOrder := []Validator{validator4, validator3, validator2, validator1}
	expectedLeader := validator4

	t.Run("with fixed ordering enabled, data not cached", func(t *testing.T) {
		t.Parallel()

		arguments := createArguments()
		arguments.EnableEpochsHandler = &enableEpochsHandlerMock.EnableEpochsHandlerStub{
			IsFlagEnabledInEpochCalled: func(flag core.EnableEpochFlag, epoch uint32) bool {
				return true
			},
		}

		ihnc, err := NewIndexHashedNodesCoordinator(arguments)
		require.Nil(t, err)

		selector := &nodesCoordinatorMocks.RandomSelectorMock{
			SelectCalled: selectFunc,
		}

		leader, cg, err := ihnc.selectLeaderAndConsensusGroup(selector, randomness, eligibleList, consensusSize, epoch)
		require.Nil(t, err)
		require.Equal(t, validator4, leader)
		require.Equal(t, expectedLeader, leader)
		require.Equal(t, expectedConsensusFixedOrder, cg)
	})
	t.Run("with fixed ordering disabled, data not cached", func(t *testing.T) {
		t.Parallel()
		arguments := createArguments()
		arguments.EnableEpochsHandler = &enableEpochsHandlerMock.EnableEpochsHandlerStub{
			IsFlagEnabledInEpochCalled: func(flag core.EnableEpochFlag, epoch uint32) bool {
				return false
			},
		}

		ihnc, err := NewIndexHashedNodesCoordinator(arguments)
		require.Nil(t, err)

		selector := &nodesCoordinatorMocks.RandomSelectorMock{
			SelectCalled: selectFunc,
		}

		leader, cg, err := ihnc.selectLeaderAndConsensusGroup(selector, randomness, eligibleList, consensusSize, epoch)
		require.Nil(t, err)
		require.Equal(t, expectedLeader, leader)
		require.Equal(t, expectedConsensusNotFixedOrder, cg)
	})
}

func TestIndexHashedNodesCoordinator_GetCachedEpochs(t *testing.T) {
	t.Parallel()

	arguments := createArguments()
	ihnc, err := NewIndexHashedNodesCoordinator(arguments)
	require.Nil(t, err)

	cachedEpochs := ihnc.GetCachedEpochs()
	require.Equal(t, 1, len(cachedEpochs))

	// add new epoch
	ihnc.AddDummyEpoch(1)
	cachedEpochs = ihnc.GetCachedEpochs()
	require.Equal(t, 2, len(cachedEpochs))

	// add new epoch
	ihnc.AddDummyEpoch(2)
	cachedEpochs = ihnc.GetCachedEpochs()
	require.Equal(t, 3, len(cachedEpochs))
}

type consensusSizeChangeTestArgs struct {
	t                     *testing.T
	ihnc                  *indexHashedNodesCoordinator
	epoch                 uint32
	expectedShardMinNodes uint32
	expectedMetaMinNodes  uint32
}

func changeEpochAndTestNewConsensusSizes(args *consensusSizeChangeTestArgs) {
	header := &block.MetaBlock{
		PrevRandSeed: []byte("rand seed"),
		EpochStart:   block.EpochStart{LastFinalizedHeaders: []block.EpochStartShardData{{}}},
	}

	header.Epoch = args.epoch
	epochForPrevConfig := uint32(0)
	if args.epoch > 0 {
		epochForPrevConfig = args.epoch - 1
	}
	args.ihnc.nodesConfig[args.epoch] = args.ihnc.nodesConfig[epochForPrevConfig]
	body := createBlockBodyFromNodesCoordinator(args.ihnc, args.epoch, args.ihnc.validatorInfoCacher)
	args.ihnc.EpochStartPrepare(header, body)
	args.ihnc.EpochStartAction(header)
	require.Len(args.t, args.ihnc.nodesConfig[args.epoch].eligibleMap[0], int(args.expectedShardMinNodes))
	require.Len(args.t, args.ihnc.nodesConfig[args.epoch].eligibleMap[common.MetachainShardId], int(args.expectedMetaMinNodes))
}
