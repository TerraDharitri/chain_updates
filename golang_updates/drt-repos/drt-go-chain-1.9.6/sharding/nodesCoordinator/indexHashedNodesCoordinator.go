package nodesCoordinator

import (
	"bytes"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	atomicFlags "github.com/TerraDharitri/drt-go-chain-core/core/atomic"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain-core/data/endProcess"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	logger "github.com/TerraDharitri/drt-go-chain-logger"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/epochStart"
	"github.com/TerraDharitri/drt-go-chain/state"
	"github.com/TerraDharitri/drt-go-chain/storage"
)

var _ NodesCoordinator = (*indexHashedNodesCoordinator)(nil)
var _ PublicKeysSelector = (*indexHashedNodesCoordinator)(nil)

const (
	keyFormat               = "%s_%v_%v_%v"
	defaultSelectionChances = uint32(1)
	minEpochsToWait         = uint32(1)
	leaderSelectionSize     = 1
)

// TODO: move this to config parameters
const nodesCoordinatorStoredEpochs = 4

type validatorWithShardID struct {
	validator Validator
	shardID   uint32
}

// savedConsensusGroup holds the leader and consensus group for a specific selection
type savedConsensusGroup struct {
	leader         Validator
	consensusGroup []Validator
}

type validatorList []Validator

// Len will return the length of the validatorList
func (v validatorList) Len() int { return len(v) }

// Swap will interchange the objects on input indexes
func (v validatorList) Swap(i, j int) { v[i], v[j] = v[j], v[i] }

// Less will return true if object on index i should appear before object in index j
// Sorting of validators should be by index and public key
func (v validatorList) Less(i, j int) bool {
	if v[i].Index() == v[j].Index() {
		return bytes.Compare(v[i].PubKey(), v[j].PubKey()) < 0
	}
	return v[i].Index() < v[j].Index()
}

// TODO: add a parameter for shardID  when acting as observer
type epochNodesConfig struct {
	nbShards       uint32
	shardID        uint32
	eligibleMap    map[uint32][]Validator
	waitingMap     map[uint32][]Validator
	selectors      map[uint32]RandomSelector
	leavingMap     map[uint32][]Validator
	shuffledOutMap map[uint32][]Validator
	newList        []Validator
	auctionList    []Validator
	mutNodesMaps   sync.RWMutex
	lowWaitingList bool
}

type indexHashedNodesCoordinator struct {
	shardIDAsObserver               uint32
	currentEpoch                    uint32
	chainParametersHandler          ChainParametersHandler
	numTotalEligible                uint64
	selfPubKey                      []byte
	savedStateKey                   []byte
	marshalizer                     marshal.Marshalizer
	hasher                          hashing.Hasher
	shuffler                        NodesShuffler
	epochStartRegistrationHandler   EpochStartEventNotifier
	bootStorer                      storage.Storer
	nodesConfig                     map[uint32]*epochNodesConfig
	mutNodesConfig                  sync.RWMutex
	mutSavedStateKey                sync.RWMutex
	nodesCoordinatorHelper          NodesCoordinatorHelper
	consensusGroupCacher            Cacher
	loadingFromDisk                 atomic.Value
	shuffledOutHandler              ShuffledOutHandler
	startEpoch                      uint32
	publicKeyToValidatorMap         map[string]*validatorWithShardID
	isFullArchive                   bool
	chanStopNode                    chan endProcess.ArgEndProcess
	nodeTypeProvider                NodeTypeProviderHandler
	enableEpochsHandler             common.EnableEpochsHandler
	validatorInfoCacher             epochStart.ValidatorInfoCacher
	genesisNodesSetupHandler        GenesisNodesSetupHandler
	flagStakingV4Step2              atomicFlags.Flag
	nodesCoordinatorRegistryFactory NodesCoordinatorRegistryFactory
	flagStakingV4Started            atomicFlags.Flag
}

// NewIndexHashedNodesCoordinator creates a new index hashed group selector
func NewIndexHashedNodesCoordinator(arguments ArgNodesCoordinator) (*indexHashedNodesCoordinator, error) {
	err := checkArguments(arguments)
	if err != nil {
		return nil, err
	}

	nodesConfig := make(map[uint32]*epochNodesConfig, nodesCoordinatorStoredEpochs)

	nodesConfig[arguments.Epoch] = &epochNodesConfig{
		nbShards:       arguments.NbShards,
		shardID:        arguments.ShardIDAsObserver,
		eligibleMap:    make(map[uint32][]Validator),
		waitingMap:     make(map[uint32][]Validator),
		selectors:      make(map[uint32]RandomSelector),
		leavingMap:     make(map[uint32][]Validator),
		shuffledOutMap: make(map[uint32][]Validator),
		newList:        make([]Validator, 0),
		auctionList:    make([]Validator, 0),
		lowWaitingList: false,
	}

	// todo: if not genesis, use previous randomness from start of epoch meta block
	savedKey := arguments.Hasher.Compute(string(arguments.SelfPublicKey))

	ihnc := &indexHashedNodesCoordinator{
		marshalizer:                     arguments.Marshalizer,
		hasher:                          arguments.Hasher,
		shuffler:                        arguments.Shuffler,
		epochStartRegistrationHandler:   arguments.EpochStartNotifier,
		bootStorer:                      arguments.BootStorer,
		selfPubKey:                      arguments.SelfPublicKey,
		nodesConfig:                     nodesConfig,
		currentEpoch:                    arguments.Epoch,
		savedStateKey:                   savedKey,
		chainParametersHandler:          arguments.ChainParametersHandler,
		consensusGroupCacher:            arguments.ConsensusGroupCache,
		shardIDAsObserver:               arguments.ShardIDAsObserver,
		shuffledOutHandler:              arguments.ShuffledOutHandler,
		startEpoch:                      arguments.StartEpoch,
		publicKeyToValidatorMap:         make(map[string]*validatorWithShardID),
		chanStopNode:                    arguments.ChanStopNode,
		nodeTypeProvider:                arguments.NodeTypeProvider,
		isFullArchive:                   arguments.IsFullArchive,
		enableEpochsHandler:             arguments.EnableEpochsHandler,
		validatorInfoCacher:             arguments.ValidatorInfoCacher,
		genesisNodesSetupHandler:        arguments.GenesisNodesSetupHandler,
		nodesCoordinatorRegistryFactory: arguments.NodesCoordinatorRegistryFactory,
	}

	ihnc.loadingFromDisk.Store(false)

	ihnc.nodesCoordinatorHelper = ihnc
	err = ihnc.setNodesPerShards(arguments.EligibleNodes, arguments.WaitingNodes, nil, nil, arguments.Epoch, false)
	if err != nil {
		return nil, err
	}

	ihnc.fillPublicKeyToValidatorMap()
	err = ihnc.saveState(ihnc.savedStateKey, arguments.Epoch)
	if err != nil {
		log.Error("saving initial nodes coordinator config failed",
			"error", err.Error())
	}
	log.Info("new nodes config is set for epoch", "epoch", arguments.Epoch)
	currentNodesConfig := ihnc.nodesConfig[arguments.Epoch]
	if currentNodesConfig == nil {
		return nil, fmt.Errorf("%w epoch=%v", ErrEpochNodesConfigDoesNotExist, arguments.Epoch)
	}

	currentConfig := nodesConfig[arguments.Epoch]
	if currentConfig == nil {
		return nil, fmt.Errorf("%w epoch=%v", ErrEpochNodesConfigDoesNotExist, arguments.Epoch)
	}

	displayNodesConfiguration(
		currentConfig.eligibleMap,
		currentConfig.waitingMap,
		currentConfig.leavingMap,
		make(map[uint32][]Validator),
		currentConfig.shuffledOutMap,
		currentConfig.nbShards)

	ihnc.epochStartRegistrationHandler.RegisterHandler(ihnc)

	return ihnc, nil
}

func checkArguments(arguments ArgNodesCoordinator) error {
	if check.IfNil(arguments.ChainParametersHandler) {
		return ErrNilChainParametersHandler
	}
	if arguments.NbShards < 1 {
		return ErrInvalidNumberOfShards
	}
	if arguments.ShardIDAsObserver >= arguments.NbShards && arguments.ShardIDAsObserver != core.MetachainShardId {
		return ErrInvalidShardId
	}
	if check.IfNil(arguments.Hasher) {
		return ErrNilHasher
	}
	if len(arguments.SelfPublicKey) == 0 {
		return ErrNilPubKey
	}
	if check.IfNil(arguments.Shuffler) {
		return ErrNilShuffler
	}
	if check.IfNil(arguments.BootStorer) {
		return ErrNilBootStorer
	}
	if check.IfNilReflect(arguments.ConsensusGroupCache) {
		return ErrNilCacher
	}
	if check.IfNil(arguments.Marshalizer) {
		return ErrNilMarshalizer
	}
	if check.IfNil(arguments.ShuffledOutHandler) {
		return ErrNilShuffledOutHandler
	}
	if check.IfNil(arguments.NodeTypeProvider) {
		return ErrNilNodeTypeProvider
	}
	if check.IfNil(arguments.NodesCoordinatorRegistryFactory) {
		return ErrNilNodesCoordinatorRegistryFactory
	}
	if nil == arguments.ChanStopNode {
		return ErrNilNodeStopChannel
	}
	if check.IfNil(arguments.EnableEpochsHandler) {
		return ErrNilEnableEpochsHandler
	}
	err := core.CheckHandlerCompatibility(arguments.EnableEpochsHandler, []core.EnableEpochFlag{
		common.RefactorPeersMiniBlocksFlag,
	})
	if err != nil {
		return err
	}
	if check.IfNil(arguments.ValidatorInfoCacher) {
		return ErrNilValidatorInfoCacher
	}
	if check.IfNil(arguments.GenesisNodesSetupHandler) {
		return ErrNilGenesisNodesSetupHandler
	}

	return nil
}

// setNodesPerShards loads the distribution of nodes per shard into the nodes management component
func (ihnc *indexHashedNodesCoordinator) setNodesPerShards(
	eligible map[uint32][]Validator,
	waiting map[uint32][]Validator,
	leaving map[uint32][]Validator,
	shuffledOut map[uint32][]Validator,
	epoch uint32,
	lowWaitingList bool,
) error {
	ihnc.mutNodesConfig.Lock()
	defer ihnc.mutNodesConfig.Unlock()

	nodesConfig, ok := ihnc.nodesConfig[epoch]
	if !ok {
		log.Debug("Did not find nodesConfig", "epoch", epoch)
		nodesConfig = &epochNodesConfig{}
	}

	nodesConfig.mutNodesMaps.Lock()
	defer nodesConfig.mutNodesMaps.Unlock()

	if eligible == nil || waiting == nil {
		return ErrNilInputNodesMap
	}

	currentChainParameters, err := ihnc.chainParametersHandler.ChainParametersForEpoch(epoch)
	if err != nil {
		return err
	}

	nodesList := eligible[core.MetachainShardId]
	if len(nodesList) < int(currentChainParameters.MetachainConsensusGroupSize) {
		return ErrSmallMetachainEligibleListSize
	}

	numTotalEligible := uint64(len(nodesList))
	for shardId := uint32(0); shardId < uint32(len(eligible)-1); shardId++ {
		nbNodesShard := len(eligible[shardId])
		if nbNodesShard < int(currentChainParameters.ShardConsensusGroupSize) {
			return ErrSmallShardEligibleListSize
		}
		numTotalEligible += uint64(nbNodesShard)
	}

	var isCurrentNodeValidator bool
	// nbShards holds number of shards without meta
	nodesConfig.nbShards = uint32(len(eligible) - 1)
	nodesConfig.eligibleMap = eligible
	nodesConfig.waitingMap = waiting
	nodesConfig.leavingMap = leaving
	nodesConfig.shuffledOutMap = shuffledOut
	nodesConfig.lowWaitingList = lowWaitingList
	nodesConfig.shardID, isCurrentNodeValidator = ihnc.computeShardForSelfPublicKey(nodesConfig)
	nodesConfig.selectors, err = ihnc.createSelectors(nodesConfig)
	if err != nil {
		return err
	}

	ihnc.nodesConfig[epoch] = nodesConfig
	ihnc.numTotalEligible = numTotalEligible
	ihnc.setNodeType(isCurrentNodeValidator)

	if ihnc.isFullArchive && isCurrentNodeValidator {
		ihnc.chanStopNode <- endProcess.ArgEndProcess{
			Reason:      common.WrongConfiguration,
			Description: ErrValidatorCannotBeFullArchive.Error(),
		}

		return nil
	}

	return nil
}

func (ihnc *indexHashedNodesCoordinator) setNodeType(isValidator bool) {
	if isValidator {
		ihnc.nodeTypeProvider.SetType(core.NodeTypeValidator)
		return
	}

	ihnc.nodeTypeProvider.SetType(core.NodeTypeObserver)
}

// ComputeAdditionalLeaving - computes extra leaving validators based on computation at the start of epoch
func (ihnc *indexHashedNodesCoordinator) ComputeAdditionalLeaving(_ []*state.ShardValidatorInfo) (map[uint32][]Validator, error) {
	return make(map[uint32][]Validator), nil
}

// ComputeConsensusGroup will generate a list of validators based on the eligible list
// and each eligible validator weight/chance
func (ihnc *indexHashedNodesCoordinator) ComputeConsensusGroup(
	randomness []byte,
	round uint64,
	shardID uint32,
	epoch uint32,
) (leader Validator, validatorsGroup []Validator, err error) {
	var selector RandomSelector
	var eligibleList []Validator

	log.Trace("computing consensus group for",
		"epoch", epoch,
		"shardID", shardID,
		"randomness", randomness,
		"round", round)

	if len(randomness) == 0 {
		return nil, nil, ErrNilRandomness
	}

	ihnc.mutNodesConfig.RLock()
	nodesConfig, ok := ihnc.nodesConfig[epoch]
	if ok {
		if shardID >= nodesConfig.nbShards && shardID != core.MetachainShardId {
			log.Warn("shardID is not ok", "shardID", shardID, "nbShards", nodesConfig.nbShards)
			ihnc.mutNodesConfig.RUnlock()
			return nil, nil, ErrInvalidShardId
		}
		selector = nodesConfig.selectors[shardID]
		eligibleList = nodesConfig.eligibleMap[shardID]
	}
	ihnc.mutNodesConfig.RUnlock()

	if !ok {
		return nil, nil, fmt.Errorf("%w epoch=%v", ErrEpochNodesConfigDoesNotExist, epoch)
	}

	key := []byte(fmt.Sprintf(keyFormat, string(randomness), round, shardID, epoch))
	savedCG := ihnc.searchConsensusForKey(key)
	if savedCG != nil {
		return savedCG.leader, savedCG.consensusGroup, nil
	}

	consensusSize := ihnc.ConsensusGroupSizeForShardAndEpoch(shardID, epoch)
	randomness = []byte(fmt.Sprintf("%d-%s", round, randomness))

	log.Debug("computeValidatorsGroup",
		"randomness", randomness,
		"consensus size", consensusSize,
		"eligible list length", len(eligibleList),
		"epoch", epoch,
		"round", round,
		"shardID", shardID)

	leader, validatorsGroup, err = ihnc.selectLeaderAndConsensusGroup(selector, randomness, eligibleList, consensusSize, epoch)
	if err != nil {
		return nil, nil, err
	}

	ihnc.cacheConsensusGroup(key, validatorsGroup, leader)

	return leader, validatorsGroup, nil
}

func (ihnc *indexHashedNodesCoordinator) cacheConsensusGroup(key []byte, consensusGroup []Validator, leader Validator) {
	size := leader.Size() * len(consensusGroup)
	savedCG := &savedConsensusGroup{
		leader:         leader,
		consensusGroup: consensusGroup,
	}
	ihnc.consensusGroupCacher.Put(key, savedCG, size)
}

func (ihnc *indexHashedNodesCoordinator) selectLeaderAndConsensusGroup(
	selector RandomSelector,
	randomness []byte,
	eligibleList []Validator,
	consensusSize int,
	epoch uint32,
) (Validator, []Validator, error) {
	leaderPositionInSelection := 0
	if !ihnc.enableEpochsHandler.IsFlagEnabledInEpoch(common.AndromedaFlag, epoch) {
		tempList, err := selectValidators(selector, randomness, uint32(consensusSize), eligibleList)
		if err != nil {
			return nil, nil, err
		}

		if len(tempList) == 0 {
			return nil, nil, ErrEmptyValidatorsList
		}

		return tempList[leaderPositionInSelection], tempList, nil
	}

	selectedValidators, err := selectValidators(selector, randomness, leaderSelectionSize, eligibleList)
	if err != nil {
		return nil, nil, err
	}
	return selectedValidators[leaderPositionInSelection], eligibleList, nil
}

func (ihnc *indexHashedNodesCoordinator) searchConsensusForKey(key []byte) *savedConsensusGroup {
	value, ok := ihnc.consensusGroupCacher.Get(key)
	if ok {
		savedCG, typeOk := value.(*savedConsensusGroup)
		if typeOk {
			return savedCG
		}
	}
	return nil
}

// GetValidatorWithPublicKey gets the validator with the given public key
func (ihnc *indexHashedNodesCoordinator) GetValidatorWithPublicKey(publicKey []byte) (Validator, uint32, error) {
	if len(publicKey) == 0 {
		return nil, 0, ErrNilPubKey
	}
	ihnc.mutNodesConfig.RLock()
	v, ok := ihnc.publicKeyToValidatorMap[string(publicKey)]
	ihnc.mutNodesConfig.RUnlock()
	if ok {
		return v.validator, v.shardID, nil
	}

	return nil, 0, ErrValidatorNotFound
}

// GetConsensusValidatorsPublicKeys calculates the validators consensus group for a specific shard, randomness and round number,
// returning their public keys
func (ihnc *indexHashedNodesCoordinator) GetConsensusValidatorsPublicKeys(
	randomness []byte,
	round uint64,
	shardID uint32,
	epoch uint32,
) (string, []string, error) {
	leader, consensusNodes, err := ihnc.ComputeConsensusGroup(randomness, round, shardID, epoch)
	if err != nil {
		return "", nil, err
	}

	pubKeys := make([]string, 0)

	for _, v := range consensusNodes {
		pubKeys = append(pubKeys, string(v.PubKey()))
	}

	return string(leader.PubKey()), pubKeys, nil
}

// GetAllEligibleValidatorsPublicKeysForShard will return all validators public keys for the provided shard
func (ihnc *indexHashedNodesCoordinator) GetAllEligibleValidatorsPublicKeysForShard(epoch uint32, shardID uint32) ([]string, error) {
	ihnc.mutNodesConfig.RLock()
	nodesConfig, ok := ihnc.nodesConfig[epoch]
	ihnc.mutNodesConfig.RUnlock()

	if !ok {
		return nil, fmt.Errorf("%w epoch=%v", ErrEpochNodesConfigDoesNotExist, epoch)
	}

	nodesConfig.mutNodesMaps.RLock()
	defer nodesConfig.mutNodesMaps.RUnlock()

	shardEligible := nodesConfig.eligibleMap[shardID]
	validatorsPubKeys := make([]string, 0, len(shardEligible))
	for i := 0; i < len(shardEligible); i++ {
		validatorsPubKeys = append(validatorsPubKeys, string(shardEligible[i].PubKey()))
	}

	return validatorsPubKeys, nil
}

// GetAllEligibleValidatorsPublicKeys will return all validators public keys for all shards
func (ihnc *indexHashedNodesCoordinator) GetAllEligibleValidatorsPublicKeys(epoch uint32) (map[uint32][][]byte, error) {
	validatorsPubKeys := make(map[uint32][][]byte)

	ihnc.mutNodesConfig.RLock()
	nodesConfig, ok := ihnc.nodesConfig[epoch]
	ihnc.mutNodesConfig.RUnlock()

	if !ok {
		return nil, fmt.Errorf("%w epoch=%v", ErrEpochNodesConfigDoesNotExist, epoch)
	}

	nodesConfig.mutNodesMaps.RLock()
	defer nodesConfig.mutNodesMaps.RUnlock()

	for shardID, shardEligible := range nodesConfig.eligibleMap {
		for i := 0; i < len(shardEligible); i++ {
			validatorsPubKeys[shardID] = append(validatorsPubKeys[shardID], shardEligible[i].PubKey())
		}
	}

	return validatorsPubKeys, nil
}

// GetAllWaitingValidatorsPublicKeys will return all validators public keys for all shards
func (ihnc *indexHashedNodesCoordinator) GetAllWaitingValidatorsPublicKeys(epoch uint32) (map[uint32][][]byte, error) {
	validatorsPubKeys := make(map[uint32][][]byte)

	ihnc.mutNodesConfig.RLock()
	nodesConfig, ok := ihnc.nodesConfig[epoch]
	ihnc.mutNodesConfig.RUnlock()

	if !ok {
		return nil, fmt.Errorf("%w epoch=%v", ErrEpochNodesConfigDoesNotExist, epoch)
	}

	nodesConfig.mutNodesMaps.RLock()
	defer nodesConfig.mutNodesMaps.RUnlock()

	for shardID, shardWaiting := range nodesConfig.waitingMap {
		for i := 0; i < len(shardWaiting); i++ {
			validatorsPubKeys[shardID] = append(validatorsPubKeys[shardID], shardWaiting[i].PubKey())
		}
	}

	return validatorsPubKeys, nil
}

// GetAllLeavingValidatorsPublicKeys will return all leaving validators public keys for all shards
func (ihnc *indexHashedNodesCoordinator) GetAllLeavingValidatorsPublicKeys(epoch uint32) (map[uint32][][]byte, error) {
	validatorsPubKeys := make(map[uint32][][]byte)

	ihnc.mutNodesConfig.RLock()
	nodesConfig, ok := ihnc.nodesConfig[epoch]
	ihnc.mutNodesConfig.RUnlock()

	if !ok {
		return nil, fmt.Errorf("%w epoch=%v", ErrEpochNodesConfigDoesNotExist, epoch)
	}

	nodesConfig.mutNodesMaps.RLock()
	defer nodesConfig.mutNodesMaps.RUnlock()

	for shardID, shardLeaving := range nodesConfig.leavingMap {
		for i := 0; i < len(shardLeaving); i++ {
			validatorsPubKeys[shardID] = append(validatorsPubKeys[shardID], shardLeaving[i].PubKey())
		}
	}

	return validatorsPubKeys, nil
}

// GetAllShuffledOutValidatorsPublicKeys will return all shuffled out validator public keys from all shards
func (ihnc *indexHashedNodesCoordinator) GetAllShuffledOutValidatorsPublicKeys(epoch uint32) (map[uint32][][]byte, error) {
	validatorsPubKeys := make(map[uint32][][]byte)

	ihnc.mutNodesConfig.RLock()
	nodesConfig, ok := ihnc.nodesConfig[epoch]
	ihnc.mutNodesConfig.RUnlock()

	if !ok {
		return nil, fmt.Errorf("%w epoch=%v", ErrEpochNodesConfigDoesNotExist, epoch)
	}

	nodesConfig.mutNodesMaps.RLock()
	defer nodesConfig.mutNodesMaps.RUnlock()

	for shardID, shuffledOutList := range nodesConfig.shuffledOutMap {
		for _, shuffledOutValidator := range shuffledOutList {
			validatorsPubKeys[shardID] = append(validatorsPubKeys[shardID], shuffledOutValidator.PubKey())
		}
	}

	return validatorsPubKeys, nil
}

// GetShuffledOutToAuctionValidatorsPublicKeys will return shuffled out to auction validators public keys
func (ihnc *indexHashedNodesCoordinator) GetShuffledOutToAuctionValidatorsPublicKeys(epoch uint32) (map[uint32][][]byte, error) {
	validatorsPubKeys := make(map[uint32][][]byte)

	ihnc.mutNodesConfig.RLock()
	nodesConfig, ok := ihnc.nodesConfig[epoch]
	ihnc.mutNodesConfig.RUnlock()

	if !ok {
		return nil, fmt.Errorf("%w epoch=%v", ErrEpochNodesConfigDoesNotExist, epoch)
	}

	if nodesConfig.lowWaitingList {
		// in case of low waiting list the nodes do not go through auction but directly to waiting
		return validatorsPubKeys, nil
	}

	return ihnc.GetAllShuffledOutValidatorsPublicKeys(epoch)
}

// GetValidatorsIndexes will return validators indexes for a block
func (ihnc *indexHashedNodesCoordinator) GetValidatorsIndexes(
	publicKeys []string,
	epoch uint32,
) ([]uint64, error) {
	signersIndexes := make([]uint64, 0)

	validatorsPubKeys, err := ihnc.GetAllEligibleValidatorsPublicKeys(epoch)
	if err != nil {
		return nil, err
	}

	ihnc.mutNodesConfig.RLock()
	nodesConfig := ihnc.nodesConfig[epoch]
	ihnc.mutNodesConfig.RUnlock()

	for _, pubKey := range publicKeys {
		for index, value := range validatorsPubKeys[nodesConfig.shardID] {
			if bytes.Equal([]byte(pubKey), value) {
				signersIndexes = append(signersIndexes, uint64(index))
			}
		}
	}

	if len(publicKeys) != len(signersIndexes) {
		strHaving := "having the following keys: \n"
		for index, value := range validatorsPubKeys[nodesConfig.shardID] {
			strHaving += fmt.Sprintf(" index %d  key %s\n", index, logger.DisplayByteSlice(value))
		}

		strNeeded := "needed the following keys: \n"
		for _, pubKey := range publicKeys {
			strNeeded += fmt.Sprintf(" key %s\n", logger.DisplayByteSlice([]byte(pubKey)))
		}

		log.Error("public keys not found\n"+strHaving+"\n"+strNeeded+"\n",
			"len pubKeys", len(publicKeys),
			"len signers", len(signersIndexes),
		)

		return nil, ErrInvalidNumberPubKeys
	}

	return signersIndexes, nil
}

// GetCachedEpochs returns all epochs cached
func (ihnc *indexHashedNodesCoordinator) GetCachedEpochs() map[uint32]struct{} {
	cachedEpochs := make(map[uint32]struct{}, nodesCoordinatorStoredEpochs)

	ihnc.mutNodesConfig.RLock()
	for epoch := range ihnc.nodesConfig {
		cachedEpochs[epoch] = struct{}{}
	}
	ihnc.mutNodesConfig.RUnlock()

	return cachedEpochs
}

// EpochStartPrepare is called when an epoch start event is observed, but not yet confirmed/committed.
// Some components may need to do some initialisation on this event
func (ihnc *indexHashedNodesCoordinator) EpochStartPrepare(metaHdr data.HeaderHandler, body data.BodyHandler) {
	if !metaHdr.IsStartOfEpochBlock() {
		log.Error("could not process EpochStartPrepare on nodesCoordinator - not epoch start block")
		return
	}

	_, castOk := metaHdr.(*block.MetaBlock)
	if !castOk {
		log.Error("could not process EpochStartPrepare on nodesCoordinator - not metaBlock")
		return
	}

	randomness := metaHdr.GetPrevRandSeed()
	newEpoch := metaHdr.GetEpoch()

	if check.IfNil(body) && newEpoch == ihnc.currentEpoch {
		log.Debug("nil body provided for epoch start prepare, it is normal in case of revertStateToBlock")
		return
	}

	ihnc.updateEpochFlags(newEpoch)

	allValidatorInfo, err := ihnc.createValidatorInfoFromBody(body, ihnc.numTotalEligible, newEpoch)
	if err != nil {
		log.Error("could not create validator info from body - do nothing on nodesCoordinator epochStartPrepare", "error", err.Error())
		return
	}

	// TODO: compare with previous nodesConfig if exists
	newNodesConfig, err := ihnc.computeNodesConfigFromList(allValidatorInfo)
	if err != nil {
		log.Error("could not compute nodes config from list - do nothing on nodesCoordinator epochStartPrepare")
		return
	}

	additionalLeavingMap, err := ihnc.nodesCoordinatorHelper.ComputeAdditionalLeaving(allValidatorInfo)
	if err != nil {
		log.Error("could not compute additionalLeaving Nodes  - do nothing on nodesCoordinator epochStartPrepare")
		return
	}

	unStakeLeavingList := ihnc.createSortedListFromMap(newNodesConfig.leavingMap)
	additionalLeavingList := ihnc.createSortedListFromMap(additionalLeavingMap)

	chainParamsForEpoch, err := ihnc.chainParametersHandler.ChainParametersForEpoch(newEpoch)
	if err != nil {
		log.Warn("indexHashedNodesCoordinator.EpochStartPrepare: could not compute chain params for epoch. "+
			"Will use the current chain parameters", "epoch", newEpoch, "error", err)
		chainParamsForEpoch = ihnc.chainParametersHandler.CurrentChainParameters()
	}
	shufflerArgs := ArgsUpdateNodes{
		ChainParameters:   chainParamsForEpoch,
		Eligible:          newNodesConfig.eligibleMap,
		Waiting:           newNodesConfig.waitingMap,
		NewNodes:          newNodesConfig.newList,
		Auction:           newNodesConfig.auctionList,
		UnStakeLeaving:    unStakeLeavingList,
		AdditionalLeaving: additionalLeavingList,
		Rand:              randomness,
		NbShards:          newNodesConfig.nbShards,
		Epoch:             newEpoch,
	}

	resUpdateNodes, err := ihnc.shuffler.UpdateNodeLists(shufflerArgs)
	if err != nil {
		log.Error("could not compute UpdateNodeLists - do nothing on nodesCoordinator epochStartPrepare", "err", err.Error())
		return
	}

	leavingNodesMap, stillRemainingNodesMap := createActuallyLeavingPerShards(
		newNodesConfig.leavingMap,
		additionalLeavingMap,
		resUpdateNodes.Leaving,
	)

	err = ihnc.setNodesPerShards(resUpdateNodes.Eligible, resUpdateNodes.Waiting, leavingNodesMap, resUpdateNodes.ShuffledOut, newEpoch, resUpdateNodes.LowWaitingList)
	if err != nil {
		log.Error("set nodes per shard failed", "error", err.Error())
	}

	ihnc.fillPublicKeyToValidatorMap()
	err = ihnc.saveState(randomness, newEpoch)
	ihnc.handleErrorLog(err, "saving nodes coordinator config failed")

	displayNodesConfiguration(
		resUpdateNodes.Eligible,
		resUpdateNodes.Waiting,
		leavingNodesMap,
		stillRemainingNodesMap,
		resUpdateNodes.ShuffledOut,
		newNodesConfig.nbShards)

	ihnc.mutSavedStateKey.Lock()
	ihnc.savedStateKey = randomness
	ihnc.mutSavedStateKey.Unlock()

	ihnc.consensusGroupCacher.Clear()
}

func (ihnc *indexHashedNodesCoordinator) fillPublicKeyToValidatorMap() {
	ihnc.mutNodesConfig.Lock()
	defer ihnc.mutNodesConfig.Unlock()

	index := 0
	epochList := make([]uint32, len(ihnc.nodesConfig))
	mapAllValidators := make(map[uint32]map[string]*validatorWithShardID)
	for epoch, epochConfig := range ihnc.nodesConfig {
		epochConfig.mutNodesMaps.RLock()
		mapAllValidators[epoch] = ihnc.createPublicKeyToValidatorMap(epochConfig.eligibleMap, epochConfig.waitingMap)
		epochConfig.mutNodesMaps.RUnlock()

		epochList[index] = epoch
		index++
	}

	sort.Slice(epochList, func(i, j int) bool {
		return epochList[i] < epochList[j]
	})

	ihnc.publicKeyToValidatorMap = make(map[string]*validatorWithShardID)
	for _, epoch := range epochList {
		validatorsForEpoch := mapAllValidators[epoch]
		for pubKey, vInfo := range validatorsForEpoch {
			ihnc.publicKeyToValidatorMap[pubKey] = vInfo
		}
	}
}

func (ihnc *indexHashedNodesCoordinator) createSortedListFromMap(validatorsMap map[uint32][]Validator) []Validator {
	sortedList := make([]Validator, 0)
	for _, validators := range validatorsMap {
		sortedList = append(sortedList, validators...)
	}
	sort.Sort(validatorList(sortedList))
	return sortedList
}

// GetChance will return default chance
func (ihnc *indexHashedNodesCoordinator) GetChance(_ uint32) uint32 {
	return defaultSelectionChances
}

func (ihnc *indexHashedNodesCoordinator) computeNodesConfigFromList(
	validatorInfos []*state.ShardValidatorInfo,
) (*epochNodesConfig, error) {
	eligibleMap := make(map[uint32][]Validator)
	waitingMap := make(map[uint32][]Validator)
	leavingMap := make(map[uint32][]Validator)
	newNodesList := make([]Validator, 0)
	auctionList := make([]Validator, 0)
	if len(validatorInfos) == 0 {
		log.Warn("computeNodesConfigFromList - validatorInfos len is 0")
	}

	for _, validatorInfo := range validatorInfos {
		chance := ihnc.nodesCoordinatorHelper.GetChance(validatorInfo.TempRating)
		currentValidator, err := NewValidator(validatorInfo.PublicKey, chance, validatorInfo.Index)
		if err != nil {
			return nil, err
		}

		switch validatorInfo.List {
		case string(common.WaitingList):
			waitingMap[validatorInfo.ShardId] = append(waitingMap[validatorInfo.ShardId], currentValidator)
		case string(common.EligibleList):
			eligibleMap[validatorInfo.ShardId] = append(eligibleMap[validatorInfo.ShardId], currentValidator)
		case string(common.LeavingList):
			log.Debug("leaving node validatorInfo",
				"pk", validatorInfo.PublicKey,
				"previous list", validatorInfo.PreviousList,
				"current index", validatorInfo.Index,
				"previous index", validatorInfo.PreviousIndex,
				"shardId", validatorInfo.ShardId)
			leavingMap[validatorInfo.ShardId] = append(leavingMap[validatorInfo.ShardId], currentValidator)
			ihnc.addValidatorToPreviousMap(
				eligibleMap,
				waitingMap,
				currentValidator,
				validatorInfo,
			)
		case string(common.NewList):
			if ihnc.flagStakingV4Step2.IsSet() {
				return nil, epochStart.ErrReceivedNewListNodeInStakingV4
			}
			log.Debug("new node registered", "pk", validatorInfo.PublicKey)
			newNodesList = append(newNodesList, currentValidator)
		case string(common.InactiveList):
			log.Debug("inactive validator", "pk", validatorInfo.PublicKey)
		case string(common.JailedList):
			log.Debug("jailed validator", "pk", validatorInfo.PublicKey)
		case string(common.SelectedFromAuctionList):
			log.Debug("selected node from auction", "pk", validatorInfo.PublicKey)
			if ihnc.flagStakingV4Step2.IsSet() {
				auctionList = append(auctionList, currentValidator)
			} else {
				return nil, ErrReceivedAuctionValidatorsBeforeStakingV4
			}
		}
	}

	sort.Sort(validatorList(newNodesList))
	sort.Sort(validatorList(auctionList))
	for _, eligibleList := range eligibleMap {
		sort.Sort(validatorList(eligibleList))
	}
	for _, waitingList := range waitingMap {
		sort.Sort(validatorList(waitingList))
	}
	for _, leavingList := range leavingMap {
		sort.Sort(validatorList(leavingList))
	}

	if len(eligibleMap) == 0 {
		return nil, fmt.Errorf("%w eligible map size is zero. No validators found", ErrMapSizeZero)
	}

	nbShards := len(eligibleMap) - 1

	newNodesConfig := &epochNodesConfig{
		eligibleMap: eligibleMap,
		waitingMap:  waitingMap,
		leavingMap:  leavingMap,
		newList:     newNodesList,
		auctionList: auctionList,
		nbShards:    uint32(nbShards),
	}

	return newNodesConfig, nil
}

func (ihnc *indexHashedNodesCoordinator) addValidatorToPreviousMap(
	eligibleMap map[uint32][]Validator,
	waitingMap map[uint32][]Validator,
	currentValidator *validator,
	validatorInfo *state.ShardValidatorInfo,
) {
	shardId := validatorInfo.ShardId
	previousList := validatorInfo.PreviousList

	log.Debug("checking leaving node",
		"current list", validatorInfo.List,
		"previous list", previousList,
		"current index", validatorInfo.Index,
		"previous index", validatorInfo.PreviousIndex,
		"pk", currentValidator.PubKey(),
		"shardId", shardId)

	if !ihnc.flagStakingV4Started.IsSet() || len(previousList) == 0 {
		log.Debug("leaving node before staking v4 or with not previous list set node found in",
			"list", "eligible", "shardId", shardId, "previous list", previousList)
		eligibleMap[shardId] = append(eligibleMap[shardId], currentValidator)
		return
	}

	if previousList == string(common.EligibleList) {
		log.Debug("leaving node found in", "list", "eligible", "shardId", shardId)
		currentValidator.index = validatorInfo.PreviousIndex
		eligibleMap[shardId] = append(eligibleMap[shardId], currentValidator)
		return
	}

	if previousList == string(common.WaitingList) {
		log.Debug("leaving node found in", "list", "waiting", "shardId", shardId)
		currentValidator.index = validatorInfo.PreviousIndex
		waitingMap[shardId] = append(waitingMap[shardId], currentValidator)
		return
	}

	log.Debug("leaving node not found in eligible or waiting",
		"previous list", previousList,
		"current index", validatorInfo.Index,
		"previous index", validatorInfo.PreviousIndex,
		"pk", currentValidator.PubKey(),
		"shardId", shardId)
}

func (ihnc *indexHashedNodesCoordinator) handleErrorLog(err error, message string) {
	if err == nil {
		return
	}

	logLevel := logger.LogError
	if core.IsClosingError(err) {
		logLevel = logger.LogDebug
	}

	log.Log(logLevel, message, "error", err.Error())
}

// EpochStartAction is called upon a start of epoch event.
// NodeCoordinator has to get the nodes assignment to shards using the shuffler.
func (ihnc *indexHashedNodesCoordinator) EpochStartAction(hdr data.HeaderHandler) {
	newEpoch := hdr.GetEpoch()
	epochToRemove := int32(newEpoch) - nodesCoordinatorStoredEpochs
	needToRemove := epochToRemove >= 0
	ihnc.currentEpoch = newEpoch

	err := ihnc.saveState(ihnc.savedStateKey, newEpoch)
	ihnc.handleErrorLog(err, "saving nodes coordinator config failed")

	ihnc.mutNodesConfig.Lock()
	if needToRemove {
		for epoch := range ihnc.nodesConfig {
			if epoch <= uint32(epochToRemove) {
				delete(ihnc.nodesConfig, epoch)
			}
		}
	}
	ihnc.mutNodesConfig.Unlock()
}

// NotifyOrder returns the notification order for a start of epoch event
func (ihnc *indexHashedNodesCoordinator) NotifyOrder() uint32 {
	return common.NodesCoordinatorOrder
}

// GetSavedStateKey returns the key for the last nodes coordinator saved state
func (ihnc *indexHashedNodesCoordinator) GetSavedStateKey() []byte {
	ihnc.mutSavedStateKey.RLock()
	key := ihnc.savedStateKey
	ihnc.mutSavedStateKey.RUnlock()

	return key
}

// ShardIdForEpoch returns the nodesCoordinator configured ShardId for specified epoch if epoch configuration exists,
// otherwise error
func (ihnc *indexHashedNodesCoordinator) ShardIdForEpoch(epoch uint32) (uint32, error) {
	ihnc.mutNodesConfig.RLock()
	nodesConfig, ok := ihnc.nodesConfig[epoch]
	ihnc.mutNodesConfig.RUnlock()

	if !ok {
		return 0, fmt.Errorf("%w epoch=%v", ErrEpochNodesConfigDoesNotExist, epoch)
	}

	return nodesConfig.shardID, nil
}

// ShuffleOutForEpoch verifies if the shards changed in the new epoch and calls the shuffleOutHandler
func (ihnc *indexHashedNodesCoordinator) ShuffleOutForEpoch(epoch uint32) {
	log.Debug("shuffle out called for", "epoch", epoch)

	ihnc.mutNodesConfig.Lock()
	nodesConfig := ihnc.nodesConfig[epoch]
	ihnc.mutNodesConfig.Unlock()

	if nodesConfig == nil {
		log.Warn("shuffleOutForEpoch failed",
			"epoch", epoch,
			"error", ErrEpochNodesConfigDoesNotExist)
		return
	}

	if isValidator(nodesConfig, ihnc.selfPubKey) {
		err := ihnc.shuffledOutHandler.Process(nodesConfig.shardID)
		if err != nil {
			log.Warn("shuffle out process failed", "err", err)
		}
	}
}

func isValidator(config *epochNodesConfig, pk []byte) bool {
	if config == nil {
		return false
	}

	config.mutNodesMaps.RLock()
	defer config.mutNodesMaps.RUnlock()

	found := false
	found, _ = searchInMap(config.eligibleMap, pk)
	if found {
		return true
	}

	found, _ = searchInMap(config.waitingMap, pk)
	return found
}

func searchInMap(validatorMap map[uint32][]Validator, pk []byte) (bool, uint32) {
	for shardId, validatorsInShard := range validatorMap {
		for _, val := range validatorsInShard {
			if bytes.Equal(val.PubKey(), pk) {
				return true, shardId
			}
		}
	}
	return false, 0
}

// GetConsensusWhitelistedNodes return the whitelisted nodes allowed to send consensus messages, for each of the shards
func (ihnc *indexHashedNodesCoordinator) GetConsensusWhitelistedNodes(
	epoch uint32,
) (map[string]struct{}, error) {
	var err error
	shardEligible := make(map[string]struct{})
	publicKeysPrevEpoch := make(map[uint32][][]byte)
	prevEpochConfigExists := false

	if epoch > ihnc.startEpoch {
		publicKeysPrevEpoch, err = ihnc.GetAllEligibleValidatorsPublicKeys(epoch - 1)
		if err == nil {
			prevEpochConfigExists = true
		} else {
			log.Warn("get consensus whitelisted nodes", "error", err.Error())
		}
	}

	var prevEpochShardId uint32
	if prevEpochConfigExists {
		prevEpochShardId, err = ihnc.ShardIdForEpoch(epoch - 1)
		if err == nil {
			for _, pubKey := range publicKeysPrevEpoch[prevEpochShardId] {
				shardEligible[string(pubKey)] = struct{}{}
			}
		} else {
			log.Trace("not critical error getting shardID for epoch", "epoch", epoch-1, "error", err)
		}
	}

	publicKeysNewEpoch, errGetEligible := ihnc.GetAllEligibleValidatorsPublicKeys(epoch)
	if errGetEligible != nil {
		return nil, errGetEligible
	}

	epochShardId, errShardIdForEpoch := ihnc.ShardIdForEpoch(epoch)
	if errShardIdForEpoch != nil {
		return nil, errShardIdForEpoch
	}

	for _, pubKey := range publicKeysNewEpoch[epochShardId] {
		shardEligible[string(pubKey)] = struct{}{}
	}

	return shardEligible, nil
}

func (ihnc *indexHashedNodesCoordinator) createPublicKeyToValidatorMap(
	eligible map[uint32][]Validator,
	waiting map[uint32][]Validator,
) map[string]*validatorWithShardID {
	publicKeyToValidatorMap := make(map[string]*validatorWithShardID)
	for shardId, shardEligible := range eligible {
		for i := 0; i < len(shardEligible); i++ {
			publicKeyToValidatorMap[string(shardEligible[i].PubKey())] = &validatorWithShardID{
				validator: shardEligible[i],
				shardID:   shardId,
			}
		}
	}
	for shardId, shardWaiting := range waiting {
		for i := 0; i < len(shardWaiting); i++ {
			publicKeyToValidatorMap[string(shardWaiting[i].PubKey())] = &validatorWithShardID{
				validator: shardWaiting[i],
				shardID:   shardId,
			}
		}
	}

	return publicKeyToValidatorMap
}

func (ihnc *indexHashedNodesCoordinator) computeShardForSelfPublicKey(nodesConfig *epochNodesConfig) (uint32, bool) {
	pubKey := ihnc.selfPubKey
	selfShard := ihnc.shardIDAsObserver
	epNodesConfig, ok := ihnc.nodesConfig[ihnc.currentEpoch]
	if ok {
		log.Trace("computeShardForSelfPublicKey found existing config",
			"shard", epNodesConfig.shardID,
		)
		selfShard = epNodesConfig.shardID
	}

	found, shardId := searchInMap(nodesConfig.eligibleMap, pubKey)
	if found {
		log.Trace("computeShardForSelfPublicKey found validator in eligible",
			"epoch", ihnc.currentEpoch,
			"shard", shardId,
			"validator PK", pubKey,
		)
		return shardId, true
	}

	found, shardId = searchInMap(nodesConfig.waitingMap, pubKey)
	if found {
		log.Trace("computeShardForSelfPublicKey found validator in waiting",
			"epoch", ihnc.currentEpoch,
			"shard", shardId,
			"validator PK", pubKey,
		)
		return shardId, true
	}

	found, shardId = searchInMap(nodesConfig.leavingMap, pubKey)
	if found {
		log.Trace("computeShardForSelfPublicKey found validator in leaving",
			"epoch", ihnc.currentEpoch,
			"shard", shardId,
			"validator PK", pubKey,
		)
		return shardId, true
	}

	if ihnc.flagStakingV4Step2.IsSet() {
		found, shardId = searchInMap(nodesConfig.shuffledOutMap, pubKey)
		if found {
			log.Trace("computeShardForSelfPublicKey found validator in shuffled out",
				"epoch", ihnc.currentEpoch,
				"shard", shardId,
				"validator PK", pubKey,
			)
			return shardId, true
		}
	}

	log.Trace("computeShardForSelfPublicKey returned default",
		"shard", selfShard,
	)
	return selfShard, false
}

// ConsensusGroupSizeForShardAndEpoch returns the consensus group size for a specific shard in a given epoch
func (ihnc *indexHashedNodesCoordinator) ConsensusGroupSizeForShardAndEpoch(
	shardID uint32,
	epoch uint32,
) int {
	return common.ConsensusGroupSizeForShardAndEpoch(log, ihnc.chainParametersHandler, shardID, epoch)
}

// GetNumTotalEligible returns the number of total eligible accross all shards from current setup
func (ihnc *indexHashedNodesCoordinator) GetNumTotalEligible() uint64 {
	return ihnc.numTotalEligible
}

// GetOwnPublicKey will return current node public key  for block sign
func (ihnc *indexHashedNodesCoordinator) GetOwnPublicKey() []byte {
	return ihnc.selfPubKey
}

// IsInterfaceNil returns true if there is no value under the interface
func (ihnc *indexHashedNodesCoordinator) IsInterfaceNil() bool {
	return ihnc == nil
}

// createSelectors creates the consensus group selectors for each shard
// Not concurrent safe, needs to be called under mutex
func (ihnc *indexHashedNodesCoordinator) createSelectors(
	nodesConfig *epochNodesConfig,
) (map[uint32]RandomSelector, error) {
	var err error
	var weights []uint32

	selectors := make(map[uint32]RandomSelector)
	// weights for validators are computed according to each validator rating
	for shard, vList := range nodesConfig.eligibleMap {
		log.Debug("create selectors", "shard", shard)
		weights, err = ihnc.nodesCoordinatorHelper.ValidatorsWeights(vList)
		if err != nil {
			return nil, err
		}

		selectors[shard], err = NewSelectorExpandedList(weights, ihnc.hasher)
		if err != nil {
			return nil, err
		}
	}

	return selectors, nil
}

// ValidatorsWeights returns the weights/chances for each of the given validators
func (ihnc *indexHashedNodesCoordinator) ValidatorsWeights(validators []Validator) ([]uint32, error) {
	weights := make([]uint32, len(validators))
	for i := range validators {
		weights[i] = defaultSelectionChances
	}

	return weights, nil
}

func createActuallyLeavingPerShards(
	unstakeLeaving map[uint32][]Validator,
	additionalLeaving map[uint32][]Validator,
	leaving []Validator,
) (map[uint32][]Validator, map[uint32][]Validator) {
	actuallyLeaving := make(map[uint32][]Validator)
	actuallyRemaining := make(map[uint32][]Validator)
	processedValidatorsMap := make(map[string]bool)

	computeActuallyLeaving(unstakeLeaving, leaving, actuallyLeaving, actuallyRemaining, processedValidatorsMap)
	computeActuallyLeaving(additionalLeaving, leaving, actuallyLeaving, actuallyRemaining, processedValidatorsMap)

	return actuallyLeaving, actuallyRemaining
}

func computeActuallyLeaving(
	unstakeLeaving map[uint32][]Validator,
	leaving []Validator,
	actuallyLeaving map[uint32][]Validator,
	actuallyRemaining map[uint32][]Validator,
	processedValidatorsMap map[string]bool,
) {
	sortedShardIds := sortKeys(unstakeLeaving)
	for _, shardId := range sortedShardIds {
		leavingValidatorsPerShard := unstakeLeaving[shardId]
		for _, v := range leavingValidatorsPerShard {
			if processedValidatorsMap[string(v.PubKey())] {
				continue
			}
			processedValidatorsMap[string(v.PubKey())] = true
			found := false
			for _, leavingValidator := range leaving {
				if bytes.Equal(v.PubKey(), leavingValidator.PubKey()) {
					found = true
					break
				}
			}
			if found {
				actuallyLeaving[shardId] = append(actuallyLeaving[shardId], v)
			} else {
				actuallyRemaining[shardId] = append(actuallyRemaining[shardId], v)
			}
		}
	}
}

func selectValidators(
	selector RandomSelector,
	randomness []byte,
	selectionSize uint32,
	eligibleList []Validator,
) ([]Validator, error) {
	if check.IfNil(selector) {
		return nil, ErrNilRandomSelector
	}
	if len(randomness) == 0 {
		return nil, ErrNilRandomness
	}

	// todo: checks for indexes
	selectedIndexes, err := selector.Select(randomness, selectionSize)
	if err != nil {
		return nil, err
	}

	selectedValidators := make([]Validator, selectionSize)
	for i := range selectedValidators {
		selectedValidators[i] = eligibleList[selectedIndexes[i]]
	}

	displayValidatorsForRandomness(selectedValidators, randomness)

	return selectedValidators, nil
}

// createValidatorInfoFromBody unmarshalls body data to create validator info
func (ihnc *indexHashedNodesCoordinator) createValidatorInfoFromBody(
	body data.BodyHandler,
	previousTotal uint64,
	epoch uint32,
) ([]*state.ShardValidatorInfo, error) {
	if check.IfNil(body) {
		return nil, ErrNilBlockBody
	}

	blockBody, ok := body.(*block.Body)
	if !ok {
		return nil, ErrWrongTypeAssertion
	}

	allValidatorInfo := make([]*state.ShardValidatorInfo, 0, previousTotal)
	for _, peerMiniBlock := range blockBody.MiniBlocks {
		if peerMiniBlock.Type != block.PeerBlock {
			continue
		}

		for _, txHash := range peerMiniBlock.TxHashes {
			shardValidatorInfo, err := ihnc.getShardValidatorInfoData(txHash, epoch)
			if err != nil {
				return nil, err
			}

			allValidatorInfo = append(allValidatorInfo, shardValidatorInfo)
		}
	}

	return allValidatorInfo, nil
}

func (ihnc *indexHashedNodesCoordinator) getShardValidatorInfoData(txHash []byte, epoch uint32) (*state.ShardValidatorInfo, error) {
	if ihnc.enableEpochsHandler.IsFlagEnabledInEpoch(common.RefactorPeersMiniBlocksFlag, epoch) {
		shardValidatorInfo, err := ihnc.validatorInfoCacher.GetValidatorInfo(txHash)
		if err != nil {
			return nil, err
		}

		return shardValidatorInfo, nil
	}

	shardValidatorInfo := &state.ShardValidatorInfo{}
	err := ihnc.marshalizer.Unmarshal(shardValidatorInfo, txHash)
	if err != nil {
		return nil, err
	}

	return shardValidatorInfo, nil
}

func (ihnc *indexHashedNodesCoordinator) updateEpochFlags(epoch uint32) {
	ihnc.flagStakingV4Started.SetValue(epoch >= ihnc.enableEpochsHandler.GetActivationEpoch(common.StakingV4Step1Flag))
	log.Debug("indexHashedNodesCoordinator: flagStakingV4Started", "enabled", ihnc.flagStakingV4Started.IsSet())

	ihnc.flagStakingV4Step2.SetValue(epoch >= ihnc.enableEpochsHandler.GetActivationEpoch(common.StakingV4Step2Flag))
	log.Debug("indexHashedNodesCoordinator: flagStakingV4Step2", "enabled", ihnc.flagStakingV4Step2.IsSet())
}

// GetWaitingEpochsLeftForPublicKey returns the number of epochs left for the public key until it becomes eligible
func (ihnc *indexHashedNodesCoordinator) GetWaitingEpochsLeftForPublicKey(publicKey []byte) (uint32, error) {
	if len(publicKey) == 0 {
		return 0, ErrNilPubKey
	}

	currentEpoch := ihnc.enableEpochsHandler.GetCurrentEpoch()

	ihnc.mutNodesConfig.RLock()
	nodesConfig, ok := ihnc.nodesConfig[currentEpoch]
	ihnc.mutNodesConfig.RUnlock()

	if !ok {
		return 0, fmt.Errorf("%w epoch=%v", ErrEpochNodesConfigDoesNotExist, currentEpoch)
	}

	nodesConfig.mutNodesMaps.RLock()
	defer nodesConfig.mutNodesMaps.RUnlock()

	for shardId, shardWaiting := range nodesConfig.waitingMap {
		epochsLeft, err := ihnc.searchWaitingEpochsLeftForPublicKeyInShard(publicKey, shardId, shardWaiting)
		if err != nil {
			continue
		}

		return epochsLeft, err
	}

	return 0, ErrKeyNotFoundInWaitingList
}

func (ihnc *indexHashedNodesCoordinator) searchWaitingEpochsLeftForPublicKeyInShard(publicKey []byte, shardId uint32, shardWaiting []Validator) (uint32, error) {
	for idx, val := range shardWaiting {
		if !bytes.Equal(val.PubKey(), publicKey) {
			continue
		}

		minHysteresisNodes := ihnc.getMinHysteresisNodes(shardId)
		if minHysteresisNodes == 0 {
			return minEpochsToWait, nil
		}

		return uint32(idx)/minHysteresisNodes + minEpochsToWait, nil
	}

	return 0, ErrKeyNotFoundInWaitingList
}

func (ihnc *indexHashedNodesCoordinator) getMinHysteresisNodes(shardId uint32) uint32 {
	if shardId == common.MetachainShardId {
		return ihnc.genesisNodesSetupHandler.MinMetaHysteresisNodes()
	}

	return ihnc.genesisNodesSetupHandler.MinShardHysteresisNodes()
}
