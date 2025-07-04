package peer

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"sync"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	logger "github.com/TerraDharitri/drt-go-chain-logger"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/common/errChan"
	"github.com/TerraDharitri/drt-go-chain/common/holders"
	"github.com/TerraDharitri/drt-go-chain/common/validatorInfo"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
	"github.com/TerraDharitri/drt-go-chain/state"
	"github.com/TerraDharitri/drt-go-chain/state/accounts"
	"github.com/TerraDharitri/drt-go-chain/state/parsers"
)

var log = logger.GetOrCreate("process/peer")

var _ process.ValidatorStatisticsProcessor = (*validatorStatistics)(nil)

type validatorActionType uint8

const (
	unknownAction             validatorActionType = 0
	leaderSuccess             validatorActionType = 1
	leaderFail                validatorActionType = 2
	validatorSuccess          validatorActionType = 3
	validatorIgnoredSignature validatorActionType = 4
)

// ArgValidatorStatisticsProcessor holds all dependencies for the validatorStatistics
type ArgValidatorStatisticsProcessor struct {
	Marshalizer                          marshal.Marshalizer
	NodesCoordinator                     nodesCoordinator.NodesCoordinator
	ShardCoordinator                     sharding.Coordinator
	DataPool                             DataPool
	StorageService                       dataRetriever.StorageService
	PubkeyConv                           core.PubkeyConverter
	PeerAdapter                          state.AccountsAdapter
	Rater                                sharding.PeerAccountListAndRatingHandler
	RewardsHandler                       process.RewardsHandler
	MaxComputableRounds                  uint64
	MaxConsecutiveRoundsOfRatingDecrease uint64
	NodesSetup                           sharding.GenesisNodesSetupHandler
	GenesisNonce                         uint64
	RatingEnableEpoch                    uint32
	EnableEpochsHandler                  common.EnableEpochsHandler
}

type validatorStatistics struct {
	marshalizer                          marshal.Marshalizer
	dataPool                             DataPool
	storageService                       dataRetriever.StorageService
	nodesCoordinator                     nodesCoordinator.NodesCoordinator
	shardCoordinator                     sharding.Coordinator
	pubkeyConv                           core.PubkeyConverter
	peerAdapter                          state.AccountsAdapter
	rater                                sharding.PeerAccountListAndRatingHandler
	rewardsHandler                       process.RewardsHandler
	maxComputableRounds                  uint64
	maxConsecutiveRoundsOfRatingDecrease uint64
	missedBlocksCounters                 validatorRoundCounters
	mutValidatorStatistics               sync.RWMutex
	genesisNonce                         uint64
	ratingEnableEpoch                    uint32
	lastFinalizedRootHash                []byte
	enableEpochsHandler                  common.EnableEpochsHandler
}

// NewValidatorStatisticsProcessor instantiates a new validatorStatistics structure responsible for keeping account of
//
//	each validator actions in the consensus process
func NewValidatorStatisticsProcessor(arguments ArgValidatorStatisticsProcessor) (*validatorStatistics, error) {
	if check.IfNil(arguments.PeerAdapter) {
		return nil, process.ErrNilPeerAccountsAdapter
	}
	if check.IfNil(arguments.PubkeyConv) {
		return nil, process.ErrNilPubkeyConverter
	}
	if check.IfNil(arguments.DataPool) {
		return nil, process.ErrNilDataPoolHolder
	}
	if check.IfNil(arguments.StorageService) {
		return nil, process.ErrNilStorage
	}
	if check.IfNil(arguments.NodesCoordinator) {
		return nil, process.ErrNilNodesCoordinator
	}
	if check.IfNil(arguments.ShardCoordinator) {
		return nil, process.ErrNilShardCoordinator
	}
	if check.IfNil(arguments.Marshalizer) {
		return nil, process.ErrNilMarshalizer
	}
	if arguments.MaxComputableRounds == 0 {
		return nil, process.ErrZeroMaxComputableRounds
	}
	if arguments.MaxConsecutiveRoundsOfRatingDecrease == 0 {
		return nil, process.ErrZeroMaxConsecutiveRoundsOfRatingDecrease
	}
	if check.IfNil(arguments.Rater) {
		return nil, process.ErrNilRater
	}
	if check.IfNil(arguments.RewardsHandler) {
		return nil, process.ErrNilRewardsHandler
	}
	if check.IfNil(arguments.NodesSetup) {
		return nil, process.ErrNilNodesSetup
	}
	if check.IfNil(arguments.EnableEpochsHandler) {
		return nil, process.ErrNilEnableEpochsHandler
	}
	err := core.CheckHandlerCompatibility(arguments.EnableEpochsHandler, []core.EnableEpochFlag{
		common.StopDecreasingValidatorRatingWhenStuckFlag,
		common.SwitchJailWaitingFlag,
		common.StakingV2FlagAfterEpoch,
		common.BelowSignedThresholdFlag,
	})
	if err != nil {
		return nil, err
	}

	vs := &validatorStatistics{
		peerAdapter:                          arguments.PeerAdapter,
		pubkeyConv:                           arguments.PubkeyConv,
		nodesCoordinator:                     arguments.NodesCoordinator,
		shardCoordinator:                     arguments.ShardCoordinator,
		dataPool:                             arguments.DataPool,
		storageService:                       arguments.StorageService,
		marshalizer:                          arguments.Marshalizer,
		missedBlocksCounters:                 make(validatorRoundCounters),
		rater:                                arguments.Rater,
		rewardsHandler:                       arguments.RewardsHandler,
		maxComputableRounds:                  arguments.MaxComputableRounds,
		maxConsecutiveRoundsOfRatingDecrease: arguments.MaxConsecutiveRoundsOfRatingDecrease,
		genesisNonce:                         arguments.GenesisNonce,
		enableEpochsHandler:                  arguments.EnableEpochsHandler,
	}

	err = vs.saveInitialState(arguments.NodesSetup)
	if err != nil {
		return nil, err
	}

	return vs, nil
}

// saveNodesCoordinatorUpdates is called at the first block after start in epoch to update the state trie according to
// the shuffling and changes done by the nodesCoordinator at the end of the epoch
func (vs *validatorStatistics) saveNodesCoordinatorUpdates(epoch uint32) (bool, error) {
	log.Debug("save nodes coordinator updates ", "epoch", epoch)

	nodesMap, err := vs.nodesCoordinator.GetAllEligibleValidatorsPublicKeys(epoch)
	if err != nil {
		return false, err
	}

	var tmpNodeForcedToRemain, nodeForcedToRemain bool
	tmpNodeForcedToRemain, err = vs.saveUpdatesForNodesMap(nodesMap, common.EligibleList)
	if err != nil {
		return false, err
	}
	nodeForcedToRemain = nodeForcedToRemain || tmpNodeForcedToRemain

	nodesMap, err = vs.nodesCoordinator.GetAllWaitingValidatorsPublicKeys(epoch)
	if err != nil {
		return false, err
	}

	tmpNodeForcedToRemain, err = vs.saveUpdatesForNodesMap(nodesMap, common.WaitingList)
	if err != nil {
		return false, err
	}
	nodeForcedToRemain = nodeForcedToRemain || tmpNodeForcedToRemain

	nodesMap, err = vs.nodesCoordinator.GetAllLeavingValidatorsPublicKeys(epoch)
	if err != nil {
		return false, err
	}

	tmpNodeForcedToRemain, err = vs.saveUpdatesForNodesMap(nodesMap, common.InactiveList)
	if err != nil {
		return false, err
	}
	nodeForcedToRemain = nodeForcedToRemain || tmpNodeForcedToRemain

	if vs.enableEpochsHandler.IsFlagEnabled(common.StakingV4Step2Flag) {
		nodesMap, err = vs.nodesCoordinator.GetShuffledOutToAuctionValidatorsPublicKeys(epoch)
		if err != nil {
			return false, err
		}

		_, err = vs.saveUpdatesForNodesMap(nodesMap, common.AuctionList)
		if err != nil {
			return false, err
		}
	}

	return nodeForcedToRemain, nil
}

func (vs *validatorStatistics) saveUpdatesForNodesMap(
	nodesMap map[uint32][][]byte,
	peerType common.PeerType,
) (bool, error) {
	nodeForcedToRemain := false
	for shardID := uint32(0); shardID < vs.shardCoordinator.NumberOfShards(); shardID++ {
		tmpNodeForcedToRemain, err := vs.saveUpdatesForList(nodesMap[shardID], shardID, peerType)
		if err != nil {
			return false, err
		}

		nodeForcedToRemain = nodeForcedToRemain || tmpNodeForcedToRemain
	}

	tmpNodeForcedToRemain, err := vs.saveUpdatesForList(nodesMap[core.MetachainShardId], core.MetachainShardId, peerType)
	if err != nil {
		return false, err
	}

	nodeForcedToRemain = nodeForcedToRemain || tmpNodeForcedToRemain
	return nodeForcedToRemain, nil
}

func (vs *validatorStatistics) saveUpdatesForList(
	pks [][]byte,
	shardID uint32,
	peerType common.PeerType,
) (bool, error) {
	nodeForcedToStay := false
	for index, pubKey := range pks {
		peerAcc, err := vs.loadPeerAccount(pubKey)
		if err != nil {
			log.Debug("error getting peer account", "error", err, "key", pubKey)
			return false, err
		}

		isNodeLeaving := (peerType == common.WaitingList || peerType == common.EligibleList) && peerAcc.GetList() == string(common.LeavingList)
		isNodeWithLowRating := vs.isValidatorWithLowRating(peerAcc)
		isNodeJailed := vs.enableEpochsHandler.IsFlagEnabled(common.SwitchJailWaitingFlag) && peerType == common.InactiveList && isNodeWithLowRating
		isStakingV4Started := vs.enableEpochsHandler.IsFlagEnabled(common.StakingV4StartedFlag)
		if isNodeJailed {
			peerAcc.SetListAndIndex(shardID, string(common.JailedList), uint32(index), isStakingV4Started)
		} else if isNodeLeaving {
			peerAcc.SetListAndIndex(shardID, string(common.LeavingList), uint32(index), isStakingV4Started)
			if isStakingV4Started {
				peerAcc.SetPreviousList(string(peerType))
			}
		} else {
			peerAcc.SetListAndIndex(shardID, string(peerType), uint32(index), isStakingV4Started)
		}

		err = vs.peerAdapter.SaveAccount(peerAcc)
		if err != nil {
			return false, err
		}

		nodeForcedToStay = nodeForcedToStay || isNodeLeaving
	}

	return nodeForcedToStay, nil
}

// saveInitialState takes an initial peer list, validates it and sets up the initial state for each of the peers
func (vs *validatorStatistics) saveInitialState(nodesConfig sharding.GenesisNodesSetupHandler) error {
	eligibleNodesInfo, waitingNodesInfo := nodesConfig.InitialNodesInfo()
	err := vs.saveInitialValueForMap(eligibleNodesInfo, common.EligibleList)
	if err != nil {
		return err
	}

	err = vs.saveInitialValueForMap(waitingNodesInfo, common.WaitingList)
	if err != nil {
		return err
	}

	hash, err := vs.peerAdapter.Commit()
	if err != nil {
		return err
	}

	log.Trace("committed peer adapter", "root hash", hex.EncodeToString(hash))

	return nil
}

func (vs *validatorStatistics) saveInitialValueForMap(
	nodesInfo map[uint32][]nodesCoordinator.GenesisNodeInfoHandler,
	peerType common.PeerType,
) error {
	if len(nodesInfo) == 0 {
		return nil
	}

	for shardID := uint32(0); shardID < vs.shardCoordinator.NumberOfShards(); shardID++ {
		nodeInfoList := nodesInfo[shardID]
		for index, nodeInfo := range nodeInfoList {
			err := vs.initializeNode(nodeInfo, shardID, peerType, uint32(index))
			if err != nil {
				return err
			}
		}
	}

	shardID := core.MetachainShardId
	nodeInfoList := nodesInfo[shardID]
	for index, nodeInfo := range nodeInfoList {
		err := vs.initializeNode(nodeInfo, shardID, peerType, uint32(index))
		if err != nil {
			return err
		}
	}

	return nil
}

// SaveNodesCoordinatorUpdates saves the results from the nodes coordinator changes after end of epoch
func (vs *validatorStatistics) SaveNodesCoordinatorUpdates(epoch uint32) (bool, error) {
	nodeForcedToRemain, err := vs.saveNodesCoordinatorUpdates(epoch)
	if err != nil {
		log.Error("could not update info from nodesCoordinator")
		return nodeForcedToRemain, err
	}

	return nodeForcedToRemain, nil
}

// UpdatePeerState takes a header, updates the peer state for all of the
// consensus members and returns the new root hash
func (vs *validatorStatistics) UpdatePeerState(header data.MetaHeaderHandler, cache map[string]data.HeaderHandler) ([]byte, error) {
	if header.GetNonce() == vs.genesisNonce {
		return vs.peerAdapter.RootHash()
	}

	vs.mutValidatorStatistics.Lock()
	vs.missedBlocksCounters.reset()
	vs.mutValidatorStatistics.Unlock()

	previousHeader, ok := cache[string(header.GetPrevHash())]
	if !ok {
		return nil, fmt.Errorf("%w - updatePeerState get header from cache - hash: %s, round: %v, nonce: %v",
			process.ErrMissingHeader,
			hex.EncodeToString(header.GetPrevHash()),
			header.GetRound(),
			header.GetNonce())
	}

	epoch := computeEpoch(header)
	err := vs.checkForMissedBlocks(
		header.GetRound(),
		previousHeader.GetRound(),
		previousHeader.GetRandSeed(),
		previousHeader.GetShardID(),
		epoch,
	)
	if err != nil {
		return nil, err
	}

	log.Debug("UpdatePeerState - registering shard leader fees", "metaNonce", header.GetNonce())
	err = vs.updateShardDataPeerState(header, cache)
	if err != nil {
		return nil, err
	}

	err = vs.updateMissedBlocksCounters()
	if err != nil {
		return nil, err
	}

	if header.GetNonce() == vs.genesisNonce+1 {
		return vs.peerAdapter.RootHash()
	}
	log.Trace("Increasing", "round", previousHeader.GetRound(), "prevRandSeed", previousHeader.GetPrevRandSeed())

	consensusGroupEpoch := computeEpoch(previousHeader)
	leader, consensusGroup, err := vs.nodesCoordinator.ComputeConsensusGroup(
		previousHeader.GetPrevRandSeed(),
		previousHeader.GetRound(),
		previousHeader.GetShardID(),
		consensusGroupEpoch)
	if err != nil {
		return nil, err
	}

	encodedLeaderPk := vs.pubkeyConv.SilentEncode(leader.PubKey(), log)

	leaderPK := core.GetTrimmedPk(encodedLeaderPk)
	log.Trace("Increasing for leader", "leader", leaderPK, "round", previousHeader.GetRound())

	log.Debug("UpdatePeerState - registering meta previous leader fees", "metaNonce", previousHeader.GetNonce())

	bitmap := vs.getBitmapForHeader(previousHeader)
	err = vs.updateValidatorInfoOnSuccessfulBlock(
		leader,
		consensusGroup,
		bitmap,
		big.NewInt(0).Sub(previousHeader.GetAccumulatedFees(), previousHeader.GetDeveloperFees()),
		previousHeader.GetShardID(),
		previousHeader.GetEpoch(),
	)
	if err != nil {
		return nil, err
	}

	rootHash, err := vs.peerAdapter.RootHash()
	if err != nil {
		return nil, err
	}

	log.Trace("after updating validator stats", "rootHash", rootHash, "round", header.GetRound(), "selfId", vs.shardCoordinator.SelfId())

	return rootHash, nil
}

func (vs *validatorStatistics) getBitmapForHeader(header data.HeaderHandler) []byte {
	bitmap := header.GetPubKeysBitmap()
	if vs.enableEpochsHandler.IsFlagEnabledInEpoch(common.AndromedaFlag, header.GetEpoch()) {
		bitmap = vs.getBitmapForFullConsensus(header.GetShardID(), header.GetEpoch())
	}
	return bitmap
}

func computeEpoch(header data.HeaderHandler) uint32 {
	// TODO: change if start of epoch block needs to be validated by the new epoch nodes
	// previous block was proposed by the consensus group of the previous epoch
	epoch := header.GetEpoch()
	if header.IsStartOfEpochBlock() && epoch > 0 {
		epoch = epoch - 1
	}

	return epoch
}

// DisplayRatings will print the ratings
func (vs *validatorStatistics) DisplayRatings(epoch uint32) {
	validatorPKs, err := vs.nodesCoordinator.GetAllEligibleValidatorsPublicKeys(epoch)
	if err != nil {
		log.Warn("could not get ValidatorPublicKeys", "epoch", epoch)
		return
	}
	log.Trace("started printing tempRatings")
	for shardID, list := range validatorPKs {
		for _, pk := range list {
			log.Trace("tempRating", "PK", pk, "tempRating", vs.getTempRating(string(pk)), "ShardID", shardID)
		}
	}
	log.Trace("finished printing tempRatings")
}

// Commit commits the validator statistics trie and returns the root hash
func (vs *validatorStatistics) Commit() ([]byte, error) {
	return vs.peerAdapter.Commit()
}

// RootHash returns the root hash of the validator statistics trie
func (vs *validatorStatistics) RootHash() ([]byte, error) {
	return vs.peerAdapter.RootHash()
}

func (vs *validatorStatistics) getValidatorDataFromLeaves(
	leavesChannels *common.TrieIteratorChannels,
) (state.ShardValidatorsInfoMapHandler, error) {
	validators := state.NewShardValidatorsInfoMap()
	for pa := range leavesChannels.LeavesChan {
		peerAccount, err := vs.unmarshalPeer(pa)
		if err != nil {
			return nil, err
		}

		validatorInfoData := vs.PeerAccountToValidatorInfo(peerAccount)
		err = validators.Add(validatorInfoData)
		if err != nil {
			return nil, err
		}
	}

	err := leavesChannels.ErrChan.ReadFromChanNonBlocking()
	if err != nil {
		return nil, err
	}

	return validators, nil
}

func getActualList(peerAccount state.PeerAccountHandler) string {
	savedList := peerAccount.GetList()
	if peerAccount.GetUnStakedEpoch() == common.DefaultUnstakedEpoch {
		if savedList == string(common.InactiveList) {
			return string(common.JailedList)
		}
		return savedList
	}
	if savedList == string(common.InactiveList) {
		return savedList
	}

	return string(common.LeavingList)
}

// PeerAccountToValidatorInfo creates a validator info from the given peer account
func (vs *validatorStatistics) PeerAccountToValidatorInfo(peerAccount state.PeerAccountHandler) *state.ValidatorInfo {
	chance := vs.rater.GetChance(peerAccount.GetRating())
	startRatingChance := vs.rater.GetChance(vs.rater.GetStartRating())
	ratingModifier := float32(chance) / float32(startRatingChance)

	list := ""
	if vs.enableEpochsHandler.IsFlagEnabled(common.SwitchJailWaitingFlag) {
		list = peerAccount.GetList()
	} else {
		list = getActualList(peerAccount)
	}

	return &state.ValidatorInfo{
		PublicKey:                       peerAccount.AddressBytes(),
		ShardId:                         peerAccount.GetShardId(),
		List:                            list,
		PreviousList:                    peerAccount.GetPreviousList(),
		Index:                           peerAccount.GetIndexInList(),
		PreviousIndex:                   peerAccount.GetPreviousIndexInList(),
		TempRating:                      peerAccount.GetTempRating(),
		Rating:                          peerAccount.GetRating(),
		RatingModifier:                  ratingModifier,
		RewardAddress:                   peerAccount.GetRewardAddress(),
		LeaderSuccess:                   peerAccount.GetLeaderSuccessRate().GetNumSuccess(),
		LeaderFailure:                   peerAccount.GetLeaderSuccessRate().GetNumFailure(),
		ValidatorSuccess:                peerAccount.GetValidatorSuccessRate().GetNumSuccess(),
		ValidatorFailure:                peerAccount.GetValidatorSuccessRate().GetNumFailure(),
		ValidatorIgnoredSignatures:      peerAccount.GetValidatorIgnoredSignaturesRate(),
		TotalLeaderSuccess:              peerAccount.GetTotalLeaderSuccessRate().GetNumSuccess(),
		TotalLeaderFailure:              peerAccount.GetTotalLeaderSuccessRate().GetNumFailure(),
		TotalValidatorSuccess:           peerAccount.GetTotalValidatorSuccessRate().GetNumSuccess(),
		TotalValidatorFailure:           peerAccount.GetTotalValidatorSuccessRate().GetNumFailure(),
		TotalValidatorIgnoredSignatures: peerAccount.GetTotalValidatorIgnoredSignaturesRate(),
		NumSelectedInSuccessBlocks:      peerAccount.GetNumSelectedInSuccessBlocks(),
		AccumulatedFees:                 big.NewInt(0).Set(peerAccount.GetAccumulatedFees()),
	}
}

// IsLowRating returns true if temp rating is under 0 chance value
func (vs *validatorStatistics) IsLowRating(blsKey []byte) bool {
	acc, err := vs.peerAdapter.GetExistingAccount(blsKey)
	if err != nil {
		return false
	}

	validatorAccount, ok := acc.(state.PeerAccountHandler)
	if !ok {
		return false
	}

	return vs.isValidatorWithLowRating(validatorAccount)
}

func (vs *validatorStatistics) isValidatorWithLowRating(validatorAccount state.PeerAccountHandler) bool {
	minChance := vs.rater.GetChance(0)
	return vs.rater.GetChance(validatorAccount.GetTempRating()) < minChance
}

func (vs *validatorStatistics) jailValidatorIfBadRatingAndInactive(validatorAccount state.PeerAccountHandler) {
	if !vs.enableEpochsHandler.IsFlagEnabled(common.SwitchJailWaitingFlag) {
		return
	}

	if validatorAccount.GetList() != string(common.InactiveList) {
		return
	}
	if !vs.isValidatorWithLowRating(validatorAccount) {
		return
	}

	validatorAccount.SetListAndIndex(validatorAccount.GetShardId(), string(common.JailedList), validatorAccount.GetIndexInList(), vs.enableEpochsHandler.IsFlagEnabled(common.StakingV4StartedFlag))
}

func (vs *validatorStatistics) unmarshalPeer(peerAccountData core.KeyValueHolder) (state.PeerAccountHandler, error) {
	peerAccount, err := accounts.NewPeerAccount(peerAccountData.Key())
	if err != nil {
		return nil, err
	}
	err = vs.marshalizer.Unmarshal(peerAccount, peerAccountData.Value())
	if err != nil {
		return nil, err
	}
	return peerAccount, nil
}

// GetValidatorInfoForRootHash returns all the peer accounts from the trie with the given rootHash
func (vs *validatorStatistics) GetValidatorInfoForRootHash(rootHash []byte) (state.ShardValidatorsInfoMapHandler, error) {
	sw := core.NewStopWatch()
	sw.Start("GetValidatorInfoForRootHash")
	defer func() {
		sw.Stop("GetValidatorInfoForRootHash")
		log.Debug("GetValidatorInfoForRootHash", sw.GetMeasurements()...)
	}()

	leavesChannels := &common.TrieIteratorChannels{
		LeavesChan: make(chan core.KeyValueHolder, common.TrieLeavesChannelDefaultCapacity),
		ErrChan:    errChan.NewErrChanWrapper(),
	}
	err := vs.peerAdapter.GetAllLeaves(leavesChannels, context.Background(), rootHash, parsers.NewMainTrieLeafParser())
	if err != nil {
		return nil, err
	}

	vInfos, err := vs.getValidatorDataFromLeaves(leavesChannels)
	if err != nil {
		return nil, err
	}

	return vInfos, err
}

// ProcessRatingsEndOfEpoch makes end of epoch process on the rating
func (vs *validatorStatistics) ProcessRatingsEndOfEpoch(
	validatorInfos state.ShardValidatorsInfoMapHandler,
	epoch uint32,
) error {
	if validatorInfos == nil || len(validatorInfos.GetAllValidatorsInfo()) == 0 {
		return process.ErrNilValidatorInfos
	}

	if epoch > 0 {
		epoch = epoch - 1
	}

	signedThreshold := vs.rater.GetSignedBlocksThreshold()
	for shardId, validators := range validatorInfos.GetShardValidatorsInfoMap() {
		for _, validator := range validators {
			if !vs.enableEpochsHandler.IsFlagEnabled(common.StakingV2FlagAfterEpoch) {
				if validator.GetList() != string(common.EligibleList) {
					continue
				}
			} else {
				if validator.GetList() != string(common.EligibleList) && !validatorInfo.WasLeavingEligibleInCurrentEpoch(validator) {
					continue
				}
			}

			err := vs.verifySignaturesBelowSignedThreshold(validator, signedThreshold, shardId, epoch)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (vs *validatorStatistics) verifySignaturesBelowSignedThreshold(
	validator state.ValidatorInfoHandler,
	signedThreshold float32,
	shardId uint32,
	epoch uint32,
) error {
	if epoch < vs.ratingEnableEpoch {
		return nil
	}

	if vs.enableEpochsHandler.IsFlagEnabledInEpoch(common.AndromedaFlag, epoch) {
		return nil
	}

	validatorOccurrences := core.MaxUint32(1, validator.GetValidatorSuccess()+validator.GetValidatorFailure()+validator.GetValidatorIgnoredSignatures())
	computedThreshold := float32(validator.GetValidatorSuccess()) / float32(validatorOccurrences)

	if computedThreshold <= signedThreshold {
		increasedRatingTimes := uint32(0)
		if !vs.enableEpochsHandler.IsFlagEnabled(common.BelowSignedThresholdFlag) {
			increasedRatingTimes = validator.GetValidatorFailure()
		} else {
			increasedRatingTimes = validator.GetValidatorSuccess() + validator.GetValidatorIgnoredSignatures()
		}

		newTempRating := vs.rater.RevertIncreaseValidator(shardId, validator.GetTempRating(), increasedRatingTimes)
		pa, err := vs.loadPeerAccount(validator.GetPublicKey())
		if err != nil {
			return err
		}

		pa.SetTempRating(newTempRating)
		vs.jailValidatorIfBadRatingAndInactive(pa)
		err = vs.peerAdapter.SaveAccount(pa)
		if err != nil {
			return err
		}

		log.Debug("below signed blocks threshold",
			"pk", validator.GetPublicKey(),
			"signed %", computedThreshold,
			"validatorSuccess", validator.GetValidatorSuccess(),
			"validatorFailure", validator.GetValidatorFailure(),
			"validatorIgnored", validator.GetValidatorIgnoredSignatures(),
			"new tempRating", newTempRating,
			"old tempRating", validator.GetTempRating(),
		)

		validator.SetTempRating(newTempRating)
	}

	return nil
}

// ResetValidatorStatisticsAtNewEpoch resets the validator info at the start of a new epoch
func (vs *validatorStatistics) ResetValidatorStatisticsAtNewEpoch(vInfos state.ShardValidatorsInfoMapHandler) error {
	sw := core.NewStopWatch()
	sw.Start("ResetValidatorStatisticsAtNewEpoch")
	defer func() {
		sw.Stop("ResetValidatorStatisticsAtNewEpoch")
		log.Debug("ResetValidatorStatisticsAtNewEpoch", sw.GetMeasurements()...)
	}()

	for _, validator := range vInfos.GetAllValidatorsInfo() {
		account, err := vs.peerAdapter.LoadAccount(validator.GetPublicKey())
		if err != nil {
			return err
		}

		peerAccount, ok := account.(state.PeerAccountHandler)
		if !ok {
			return process.ErrWrongTypeAssertion
		}
		peerAccount.ResetAtNewEpoch()
		vs.setToJailedIfNeeded(peerAccount, validator)

		err = vs.peerAdapter.SaveAccount(peerAccount)
		if err != nil {
			return err
		}
	}

	return nil
}

func (vs *validatorStatistics) setToJailedIfNeeded(
	peerAccount state.PeerAccountHandler,
	validator state.ValidatorInfoHandler,
) {
	if !vs.enableEpochsHandler.IsFlagEnabled(common.SwitchJailWaitingFlag) {
		return
	}

	if validator.GetList() == string(common.WaitingList) || validator.GetList() == string(common.EligibleList) {
		return
	}

	if validator.GetList() == string(common.JailedList) && peerAccount.GetList() != string(common.JailedList) {
		peerAccount.SetListAndIndex(validator.GetShardId(), string(common.JailedList), validator.GetIndex(), vs.enableEpochsHandler.IsFlagEnabled(common.StakingV4StartedFlag))
		return
	}

	if vs.isValidatorWithLowRating(peerAccount) {
		peerAccount.SetListAndIndex(validator.GetShardId(), string(common.JailedList), validator.GetIndex(), vs.enableEpochsHandler.IsFlagEnabled(common.StakingV4StartedFlag))
	}
}

func (vs *validatorStatistics) checkForMissedBlocks(
	currentHeaderRound,
	previousHeaderRound uint64,
	prevRandSeed []byte,
	shardID uint32,
	epoch uint32,
) error {
	missedRounds := currentHeaderRound - previousHeaderRound
	if missedRounds <= 1 {
		return nil
	}
	if vs.enableEpochsHandler.IsFlagEnabled(common.StopDecreasingValidatorRatingWhenStuckFlag) {
		if missedRounds > vs.maxConsecutiveRoundsOfRatingDecrease {
			return nil
		}
	}

	tooManyComputations := missedRounds > vs.maxComputableRounds
	if !tooManyComputations {
		return vs.computeDecrease(previousHeaderRound, currentHeaderRound, prevRandSeed, shardID, epoch)
	}

	return vs.decreaseAll(shardID, missedRounds-1, epoch)
}

func (vs *validatorStatistics) computeDecrease(
	previousHeaderRound uint64,
	currentHeaderRound uint64,
	prevRandSeed []byte,
	shardID uint32,
	epoch uint32,
) error {
	if epoch < vs.ratingEnableEpoch {
		return nil
	}

	sw := core.NewStopWatch()
	sw.Start("checkForMissedBlocks")
	defer func() {
		sw.Stop("checkForMissedBlocks")
		log.Trace("measurements checkForMissedBlocks", sw.GetMeasurements()...)
	}()

	for i := previousHeaderRound + 1; i < currentHeaderRound; i++ {
		swInner := core.NewStopWatch()

		swInner.Start("ComputeValidatorsGroup")
		log.Debug("decreasing", "round", i, "prevRandSeed", prevRandSeed, "shardId", shardID)
		leader, consensusGroup, err := vs.nodesCoordinator.ComputeConsensusGroup(prevRandSeed, i, shardID, epoch)
		swInner.Stop("ComputeValidatorsGroup")
		if err != nil {
			return err
		}

		swInner.Start("loadPeerAccount")
		leaderPeerAcc, err := vs.loadPeerAccount(leader.PubKey())

		encodedLeaderPk := vs.pubkeyConv.SilentEncode(leader.PubKey(), log)
		leaderPK := core.GetTrimmedPk(encodedLeaderPk)
		swInner.Stop("loadPeerAccount")
		if err != nil {
			return err
		}

		vs.mutValidatorStatistics.Lock()
		vs.missedBlocksCounters.decreaseLeader(leader.PubKey())
		vs.mutValidatorStatistics.Unlock()

		swInner.Start("ComputeDecreaseProposer")
		newRating := vs.rater.ComputeDecreaseProposer(
			shardID,
			leaderPeerAcc.GetTempRating(),
			leaderPeerAcc.GetConsecutiveProposerMisses())
		swInner.Stop("ComputeDecreaseProposer")

		swInner.Start("SetConsecutiveProposerMisses")
		leaderPeerAcc.SetConsecutiveProposerMisses(leaderPeerAcc.GetConsecutiveProposerMisses() + 1)
		swInner.Stop("SetConsecutiveProposerMisses")

		swInner.Start("SetTempRating")
		leaderPeerAcc.SetTempRating(newRating)
		log.Debug("decreasing for leader",
			"leader", leaderPK,
			"round", i,
			"temp rating", newRating,
			"consecutive misses", leaderPeerAcc.GetConsecutiveProposerMisses())
		vs.jailValidatorIfBadRatingAndInactive(leaderPeerAcc)

		err = vs.peerAdapter.SaveAccount(leaderPeerAcc)
		swInner.Stop("SetTempRating")
		if err != nil {
			return err
		}

		swInner.Start("ComputeDecreaseAllValidators")
		err = vs.decreaseForConsensusValidators(consensusGroup, shardID, epoch)
		swInner.Stop("ComputeDecreaseAllValidators")
		if err != nil {
			return err
		}
		sw.Add(swInner)
	}
	return nil
}

func (vs *validatorStatistics) decreaseForConsensusValidators(
	consensusGroup []nodesCoordinator.Validator,
	shardId uint32,
	epoch uint32,
) error {
	if epoch < vs.ratingEnableEpoch {
		return nil
	}

	vs.mutValidatorStatistics.Lock()
	defer vs.mutValidatorStatistics.Unlock()

	for j := 1; j < len(consensusGroup); j++ {
		validatorPeerAccount, verr := vs.loadPeerAccount(consensusGroup[j].PubKey())
		if verr != nil {
			return verr
		}
		vs.missedBlocksCounters.decreaseValidator(consensusGroup[j].PubKey())

		newRating := vs.rater.ComputeDecreaseValidator(shardId, validatorPeerAccount.GetTempRating())
		validatorPeerAccount.SetTempRating(newRating)
		vs.jailValidatorIfBadRatingAndInactive(validatorPeerAccount)
		err := vs.peerAdapter.SaveAccount(validatorPeerAccount)
		if err != nil {
			return err
		}
	}

	return nil
}

// RevertPeerState takes the current and previous headers and undos the peer state
// for all of the consensus members
func (vs *validatorStatistics) RevertPeerState(header data.MetaHeaderHandler) error {
	rootHashHolder := holders.NewDefaultRootHashesHolder(header.GetValidatorStatsRootHash())
	return vs.peerAdapter.RecreateTrie(rootHashHolder)
}

// TODO: check if this can be taken from somewhere else
func (vs *validatorStatistics) getBitmapForFullConsensus(shardID uint32, epoch uint32) []byte {
	consensusSize := vs.nodesCoordinator.ConsensusGroupSizeForShardAndEpoch(shardID, epoch)
	bitmap := make([]byte, consensusSize/8+1)
	for i := 0; i < consensusSize; i++ {
		bitmap[i/8] |= 1 << (uint16(i) % 8)
	}

	return bitmap
}

func (vs *validatorStatistics) updateShardDataPeerState(
	header data.HeaderHandler,
	cacheMap map[string]data.HeaderHandler,
) error {
	metaHeader, ok := header.(*block.MetaBlock)
	if !ok {
		return process.ErrInvalidMetaHeader
	}

	var currentHeader data.HeaderHandler
	for _, h := range metaHeader.ShardInfo {
		if h.Nonce == vs.genesisNonce {
			continue
		}

		currentHeader, ok = cacheMap[string(h.HeaderHash)]
		if !ok {
			return fmt.Errorf("%w - updateShardDataPeerState header from cache - hash: %s, round: %v, nonce: %v",
				process.ErrMissingHeader,
				hex.EncodeToString(h.HeaderHash),
				h.GetRound(),
				h.GetNonce())
		}

		epoch := computeEpoch(currentHeader)

		leader, shardConsensus, shardInfoErr := vs.nodesCoordinator.ComputeConsensusGroup(h.PrevRandSeed, h.Round, h.ShardID, epoch)
		if shardInfoErr != nil {
			return shardInfoErr
		}

		log.Debug("updateShardDataPeerState - registering shard leader fees", "shard headerHash", h.HeaderHash, "accumulatedFees", h.AccumulatedFees.String(), "developerFees", h.DeveloperFees.String())
		bitmap := h.PubKeysBitmap
		if vs.enableEpochsHandler.IsFlagEnabledInEpoch(common.AndromedaFlag, h.Epoch) {
			bitmap = vs.getBitmapForFullConsensus(h.ShardID, h.Epoch)
		}
		shardInfoErr = vs.updateValidatorInfoOnSuccessfulBlock(
			leader,
			shardConsensus,
			bitmap,
			big.NewInt(0).Sub(h.AccumulatedFees, h.DeveloperFees),
			h.ShardID,
			currentHeader.GetEpoch(),
		)
		if shardInfoErr != nil {
			return shardInfoErr
		}

		if h.Nonce == vs.genesisNonce+1 {
			continue
		}

		prevShardData, shardInfoErr := vs.searchInMap(h.PrevHash, cacheMap)
		if shardInfoErr != nil {
			return shardInfoErr
		}

		shardInfoErr = vs.checkForMissedBlocks(
			h.Round,
			prevShardData.GetRound(),
			prevShardData.GetRandSeed(),
			h.ShardID,
			epoch,
		)
		if shardInfoErr != nil {
			return shardInfoErr
		}
	}

	return nil
}

func (vs *validatorStatistics) searchInMap(hash []byte, cacheMap map[string]data.HeaderHandler) (data.HeaderHandler, error) {
	blkHandler := cacheMap[string(hash)]
	if check.IfNil(blkHandler) {
		return nil, fmt.Errorf("%w : searchInMap hash = %s",
			process.ErrMissingHeader, logger.DisplayByteSlice(hash))
	}

	blk, ok := blkHandler.(data.ShardHeaderHandler)
	if !ok {
		return nil, process.ErrWrongTypeAssertion
	}

	return blk, nil
}

func (vs *validatorStatistics) initializeNode(
	node nodesCoordinator.GenesisNodeInfoHandler,
	shardID uint32,
	peerType common.PeerType,
	index uint32,
) error {
	peerAccount, err := vs.loadPeerAccount(node.PubKeyBytes())
	if err != nil {
		return err
	}

	return vs.savePeerAccountData(peerAccount, node, node.GetInitialRating(), shardID, peerType, index)
}

func (vs *validatorStatistics) savePeerAccountData(
	peerAccount state.PeerAccountHandler,
	node nodesCoordinator.GenesisNodeInfoHandler,
	startRating uint32,
	shardID uint32,
	peerType common.PeerType,
	index uint32,
) error {
	log.Trace("validatorStatistics - savePeerAccountData",
		"pubkey", node.PubKeyBytes(),
		"reward address", node.AddressBytes(),
		"initial rating", node.GetInitialRating())
	err := peerAccount.SetRewardAddress(node.AddressBytes())
	if err != nil {
		return err
	}

	peerAccount.SetRating(startRating)
	peerAccount.SetTempRating(startRating)
	peerAccount.SetListAndIndex(shardID, string(peerType), index, vs.enableEpochsHandler.IsFlagEnabled(common.StakingV4StartedFlag))

	return vs.peerAdapter.SaveAccount(peerAccount)
}

func (vs *validatorStatistics) updateValidatorInfoOnSuccessfulBlock(
	leader nodesCoordinator.Validator,
	validatorList []nodesCoordinator.Validator,
	signingBitmap []byte,
	accumulatedFees *big.Int,
	shardId uint32,
	epoch uint32,
) error {

	if len(signingBitmap) == 0 {
		return process.ErrNilPubKeysBitmap
	}
	lenValidators := len(validatorList)
	for i := 0; i < lenValidators; i++ {
		peerAcc, err := vs.loadPeerAccount(validatorList[i].PubKey())
		if err != nil {
			return err
		}

		peerAcc.IncreaseNumSelectedInSuccessBlocks()

		newRating := peerAcc.GetRating()
		isLeader := bytes.Equal(leader.PubKey(), validatorList[i].PubKey())
		validatorSigned := (signingBitmap[i/8] & (1 << (uint16(i) % 8))) != 0
		actionType := vs.computeValidatorActionType(isLeader, validatorSigned)

		switch actionType {
		case leaderSuccess:
			peerAcc.IncreaseLeaderSuccessRate(1)
			peerAcc.SetConsecutiveProposerMisses(0)
			newRating = vs.rater.ComputeIncreaseProposer(shardId, peerAcc.GetTempRating())
			var leaderAccumulatedFees *big.Int
			if vs.enableEpochsHandler.IsFlagEnabled(common.StakingV2FlagAfterEpoch) {
				leaderAccumulatedFees = core.GetIntTrimmedPercentageOfValue(accumulatedFees, vs.rewardsHandler.LeaderPercentageInEpoch(epoch))
			} else {
				leaderAccumulatedFees = core.GetApproximatePercentageOfValue(accumulatedFees, vs.rewardsHandler.LeaderPercentageInEpoch(epoch))
			}

			peerAcc.AddToAccumulatedFees(leaderAccumulatedFees)
			log.Debug("updateValidatorInfoOnSuccessfulBlock",
				"leaderAccumulatedFees in current block", leaderAccumulatedFees.String(),
				"leader fees in Epoch", peerAcc.GetAccumulatedFees().String(),
				"leader", core.GetTrimmedPk(string(peerAcc.AddressBytes())))
		case validatorSuccess:
			peerAcc.IncreaseValidatorSuccessRate(1)
			newRating = vs.rater.ComputeIncreaseValidator(shardId, peerAcc.GetTempRating())
		case validatorIgnoredSignature:
			peerAcc.IncreaseValidatorIgnoredSignaturesRate(1)
			newRating = vs.rater.ComputeIncreaseValidator(shardId, peerAcc.GetTempRating())
		}

		peerAcc.SetTempRating(newRating)

		err = vs.peerAdapter.SaveAccount(peerAcc)
		if err != nil {
			return err
		}
	}

	return nil
}

func (vs *validatorStatistics) loadPeerAccount(address []byte) (state.PeerAccountHandler, error) {
	account, err := vs.peerAdapter.LoadAccount(address)
	if err != nil {
		return nil, err
	}

	peerAccount, ok := account.(state.PeerAccountHandler)
	if !ok {
		return nil, process.ErrInvalidPeerAccount
	}

	return peerAccount, nil
}

func (vs *validatorStatistics) getMatchingPrevShardData(currentShardData block.ShardData, shardInfo []block.ShardData) *block.ShardData {
	for _, prevShardData := range shardInfo {
		if currentShardData.ShardID != prevShardData.ShardID {
			continue
		}
		if currentShardData.Nonce == prevShardData.Nonce+1 {
			return &prevShardData
		}
	}

	return nil
}

func (vs *validatorStatistics) updateMissedBlocksCounters() error {
	vs.mutValidatorStatistics.Lock()
	defer func() {
		vs.missedBlocksCounters.reset()
		vs.mutValidatorStatistics.Unlock()
	}()

	for pubKey, roundCounters := range vs.missedBlocksCounters {
		peerAccount, err := vs.loadPeerAccount([]byte(pubKey))
		if err != nil {
			return err
		}

		if roundCounters.leaderDecreaseCount > 0 {
			peerAccount.DecreaseLeaderSuccessRate(roundCounters.leaderDecreaseCount)
		}

		if roundCounters.validatorDecreaseCount > 0 {
			peerAccount.DecreaseValidatorSuccessRate(roundCounters.validatorDecreaseCount)
		}

		err = vs.peerAdapter.SaveAccount(peerAccount)
		if err != nil {
			return err
		}
	}

	return nil
}

func (vs *validatorStatistics) computeValidatorActionType(isLeader, validatorSigned bool) validatorActionType {
	if isLeader && validatorSigned {
		return leaderSuccess
	}
	if isLeader && !validatorSigned {
		return leaderFail
	}
	if !isLeader && validatorSigned {
		return validatorSuccess
	}
	if !isLeader && !validatorSigned {
		return validatorIgnoredSignature
	}

	return unknownAction
}

// IsInterfaceNil returns true if there is no value under the interface
func (vs *validatorStatistics) IsInterfaceNil() bool {
	return vs == nil
}

func (vs *validatorStatistics) getTempRating(s string) uint32 {
	peer, err := vs.loadPeerAccount([]byte(s))

	if err != nil {
		log.Debug("Error getting peer account", "error", err)
		return vs.rater.GetStartRating()
	}

	return peer.GetTempRating()
}

func (vs *validatorStatistics) display(validatorKey string) {
	if log.GetLevel() != logger.LogTrace {
		// do not need to load peer account if not log level trace
		return
	}

	peerAcc, err := vs.loadPeerAccount([]byte(validatorKey))
	if err != nil {
		log.Trace("display peer acc", "error", err)
		return
	}

	log.Trace("validator statistics",
		"pk", core.GetTrimmedPk(hex.EncodeToString(peerAcc.AddressBytes())),
		"leader fail", peerAcc.GetLeaderSuccessRate().GetNumFailure(),
		"leader success", peerAcc.GetLeaderSuccessRate().GetNumSuccess(),
		"val success", peerAcc.GetValidatorSuccessRate().GetNumSuccess(),
		"val ignored sigs", peerAcc.GetValidatorIgnoredSignaturesRate(),
		"val fail", peerAcc.GetValidatorSuccessRate().GetNumFailure(),
		"temp rating", peerAcc.GetTempRating(),
		"rating", peerAcc.GetRating(),
	)
}

func (vs *validatorStatistics) decreaseAll(
	shardID uint32,
	missedRounds uint64,
	epoch uint32,
) error {
	if epoch < vs.ratingEnableEpoch {
		return nil
	}

	log.Debug("ValidatorStatistics decreasing all", "shardID", shardID, "missedRounds", missedRounds)
	consensusGroupSize := vs.nodesCoordinator.ConsensusGroupSizeForShardAndEpoch(shardID, epoch)
	validators, err := vs.nodesCoordinator.GetAllEligibleValidatorsPublicKeys(epoch)
	if err != nil {
		return err
	}

	shardValidators := validators[shardID]
	validatorsCount := len(shardValidators)
	percentageRoundMissedFromTotalValidators := float64(missedRounds) / float64(validatorsCount)
	leaderAppearances := uint32(percentageRoundMissedFromTotalValidators + 1 - math.SmallestNonzeroFloat64)
	consensusGroupAppearances := uint32(float64(consensusGroupSize)*percentageRoundMissedFromTotalValidators +
		1 - math.SmallestNonzeroFloat64)
	ratingDifference := uint32(0)

	for i, validator := range shardValidators {
		validatorPeerAccount, errLoad := vs.loadPeerAccount(validator)
		if errLoad != nil {
			return errLoad
		}
		validatorPeerAccount.DecreaseLeaderSuccessRate(leaderAppearances)
		validatorPeerAccount.DecreaseValidatorSuccessRate(consensusGroupAppearances)

		currentTempRating := validatorPeerAccount.GetTempRating()
		for ct := uint32(0); ct < leaderAppearances; ct++ {
			currentTempRating = vs.rater.ComputeDecreaseProposer(shardID, currentTempRating, 0)
		}

		for ct := uint32(0); ct < consensusGroupAppearances; ct++ {
			currentTempRating = vs.rater.ComputeDecreaseValidator(shardID, currentTempRating)
		}

		if i == 0 {
			ratingDifference = validatorPeerAccount.GetTempRating() - currentTempRating
		}

		validatorPeerAccount.SetTempRating(currentTempRating)
		vs.jailValidatorIfBadRatingAndInactive(validatorPeerAccount)
		err = vs.peerAdapter.SaveAccount(validatorPeerAccount)
		if err != nil {
			return err
		}

		vs.display(string(validator))
	}

	log.Trace(fmt.Sprintf("Decrease leader: %v, decrease validator: %v, ratingDifference: %v", leaderAppearances, consensusGroupAppearances, ratingDifference))

	return nil
}

// Process - processes a validatorInfo and updates fields
func (vs *validatorStatistics) Process(svi data.ShardValidatorInfoHandler) error {
	log.Trace("ValidatorInfoData", "pk", svi.GetPublicKey(), "tempRating", svi.GetTempRating())

	pa, err := vs.loadPeerAccount(svi.GetPublicKey())
	if err != nil {
		return err
	}

	pa.SetRating(svi.GetTempRating())
	return vs.peerAdapter.SaveAccount(pa)
}

// SetLastFinalizedRootHash - sets the last finalized root hash needed for correct validatorStatistics computations
func (vs *validatorStatistics) SetLastFinalizedRootHash(lastFinalizedRootHash []byte) {
	if len(lastFinalizedRootHash) == 0 {
		return
	}

	vs.mutValidatorStatistics.Lock()
	vs.lastFinalizedRootHash = lastFinalizedRootHash
	vs.mutValidatorStatistics.Unlock()
}

// LastFinalizedRootHash returns the root hash of the validator statistics trie that was last finalized
func (vs *validatorStatistics) LastFinalizedRootHash() []byte {
	vs.mutValidatorStatistics.RLock()
	defer vs.mutValidatorStatistics.RUnlock()
	return vs.lastFinalizedRootHash
}
