package v2

import (
	"context"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	outportcore "github.com/TerraDharitri/drt-go-chain-core/data/outport"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/consensus/spos"
	"github.com/TerraDharitri/drt-go-chain/outport"
	"github.com/TerraDharitri/drt-go-chain/outport/disabled"
)

// subroundStartRound defines the data needed by the subround StartRound
type subroundStartRound struct {
	*spos.Subround
	processingThresholdPercentage int

	sentSignatureTracker spos.SentSignaturesTracker
	worker               spos.WorkerHandler
	outportHandler       outport.OutportHandler
	outportMutex         sync.RWMutex
}

// NewSubroundStartRound creates a subroundStartRound object
func NewSubroundStartRound(
	baseSubround *spos.Subround,
	processingThresholdPercentage int,
	sentSignatureTracker spos.SentSignaturesTracker,
	worker spos.WorkerHandler,
) (*subroundStartRound, error) {
	err := checkNewSubroundStartRoundParams(
		baseSubround,
	)
	if err != nil {
		return nil, err
	}
	if check.IfNil(sentSignatureTracker) {
		return nil, ErrNilSentSignatureTracker
	}
	if check.IfNil(worker) {
		return nil, spos.ErrNilWorker
	}

	srStartRound := subroundStartRound{
		Subround:                      baseSubround,
		processingThresholdPercentage: processingThresholdPercentage,
		sentSignatureTracker:          sentSignatureTracker,
		worker:                        worker,
		outportHandler:                disabled.NewDisabledOutport(),
		outportMutex:                  sync.RWMutex{},
	}
	srStartRound.Job = srStartRound.doStartRoundJob
	srStartRound.Check = srStartRound.doStartRoundConsensusCheck
	srStartRound.Extend = worker.Extend
	baseSubround.EpochStartRegistrationHandler().RegisterHandler(&srStartRound)

	return &srStartRound, nil
}

func checkNewSubroundStartRoundParams(
	baseSubround *spos.Subround,
) error {
	if baseSubround == nil {
		return spos.ErrNilSubround
	}
	if check.IfNil(baseSubround.ConsensusStateHandler) {
		return spos.ErrNilConsensusState
	}

	err := spos.ValidateConsensusCore(baseSubround.ConsensusCoreHandler)

	return err
}

// SetOutportHandler method sets outport handler
func (sr *subroundStartRound) SetOutportHandler(outportHandler outport.OutportHandler) error {
	if check.IfNil(outportHandler) {
		return outport.ErrNilDriver
	}

	sr.outportMutex.Lock()
	sr.outportHandler = outportHandler
	sr.outportMutex.Unlock()

	return nil
}

// doStartRoundJob method does the job of the subround StartRound
func (sr *subroundStartRound) doStartRoundJob(_ context.Context) bool {
	sr.ResetConsensusState()
	sr.SetRoundIndex(sr.RoundHandler().Index())
	sr.SetRoundTimeStamp(sr.RoundHandler().TimeStamp())
	topic := spos.GetConsensusTopicID(sr.ShardCoordinator())
	sr.GetAntiFloodHandler().ResetForTopic(topic)
	sr.worker.ResetConsensusMessages()
	sr.worker.ResetInvalidSignersCache()

	return true
}

// doStartRoundConsensusCheck method checks if the consensus is achieved in the subround StartRound
func (sr *subroundStartRound) doStartRoundConsensusCheck() bool {
	if sr.GetRoundCanceled() {
		return false
	}

	if sr.IsSubroundFinished(sr.Current()) {
		return true
	}

	if sr.initCurrentRound() {
		return true
	}

	return false
}

func (sr *subroundStartRound) initCurrentRound() bool {
	nodeState := sr.BootStrapper().GetNodeState()
	if nodeState != common.NsSynchronized { // if node is not synchronized yet, it has to continue the bootstrapping mechanism
		return false
	}

	sr.AppStatusHandler().SetStringValue(common.MetricConsensusRoundState, "")

	err := sr.generateNextConsensusGroup(sr.RoundHandler().Index())
	if err != nil {
		log.Debug("initCurrentRound.generateNextConsensusGroup",
			"round index", sr.RoundHandler().Index(),
			"error", err.Error())

		sr.SetRoundCanceled(true)

		return false
	}

	if sr.NodeRedundancyHandler().IsRedundancyNode() {
		sr.NodeRedundancyHandler().AdjustInactivityIfNeeded(
			sr.SelfPubKey(),
			sr.ConsensusGroup(),
			sr.RoundHandler().Index(),
		)
		// we should not return here, the multikey redundancy system relies on it
		// the NodeRedundancyHandler "thinks" it is in redundancy mode even if we use the multikey redundancy system
	}

	leader, err := sr.GetLeader()
	if err != nil {
		log.Debug("initCurrentRound.GetLeader", "error", err.Error())

		sr.SetRoundCanceled(true)

		return false
	}

	msg := sr.GetLeaderStartRoundMessage()
	if len(msg) != 0 {
		sr.AppStatusHandler().Increment(common.MetricCountLeader)
		sr.AppStatusHandler().SetStringValue(common.MetricConsensusRoundState, "proposed")
		sr.AppStatusHandler().SetStringValue(common.MetricConsensusState, "proposer")
	}

	log.Debug("step 0: preparing the round",
		"leader", core.GetTrimmedPk(hex.EncodeToString([]byte(leader))),
		"messsage", msg)
	sr.sentSignatureTracker.StartRound()

	pubKeys := sr.ConsensusGroup()
	numMultiKeysInConsensusGroup := sr.computeNumManagedKeysInConsensusGroup(pubKeys)
	if numMultiKeysInConsensusGroup > 0 {
		log.Debug("in consensus group with multi keys identities", "num", numMultiKeysInConsensusGroup)
	}

	sr.indexRoundIfNeeded(pubKeys)

	if !sr.IsSelfInConsensusGroup() {
		log.Debug("not in consensus group")
		sr.AppStatusHandler().SetStringValue(common.MetricConsensusState, "not in consensus group")
	} else {
		if !sr.IsSelfLeader() {
			sr.AppStatusHandler().Increment(common.MetricCountConsensus)
			sr.AppStatusHandler().SetStringValue(common.MetricConsensusState, "participant")
		}
	}

	err = sr.SigningHandler().Reset(pubKeys)
	if err != nil {
		log.Debug("initCurrentRound.Reset", "error", err.Error())

		sr.SetRoundCanceled(true)

		return false
	}

	startTime := sr.GetRoundTimeStamp()
	maxTime := sr.RoundHandler().TimeDuration() * time.Duration(sr.processingThresholdPercentage) / 100
	if sr.RoundHandler().RemainingTime(startTime, maxTime) < 0 {
		log.Debug("canceled round, time is out",
			"round", sr.SyncTimer().FormattedCurrentTime(), sr.RoundHandler().Index(),
			"subround", sr.Name())

		sr.SetRoundCanceled(true)

		return false
	}

	sr.SetStatus(sr.Current(), spos.SsFinished)

	// execute stored messages which were received in this new round but before this initialisation
	go sr.worker.ExecuteStoredMessages()

	return true
}

func (sr *subroundStartRound) computeNumManagedKeysInConsensusGroup(pubKeys []string) int {
	numMultiKeysInConsensusGroup := 0
	for _, pk := range pubKeys {
		pkBytes := []byte(pk)
		if sr.IsKeyManagedBySelf(pkBytes) {
			numMultiKeysInConsensusGroup++
			log.Trace("in consensus group with multi key",
				"pk", core.GetTrimmedPk(hex.EncodeToString(pkBytes)))
		}
		sr.IncrementRoundsWithoutReceivedMessages(pkBytes)
	}

	return numMultiKeysInConsensusGroup
}

func (sr *subroundStartRound) indexRoundIfNeeded(pubKeys []string) {
	sr.outportMutex.RLock()
	defer sr.outportMutex.RUnlock()

	if !sr.outportHandler.HasDrivers() {
		return
	}

	currentHeader := sr.Blockchain().GetCurrentBlockHeader()
	if check.IfNil(currentHeader) {
		currentHeader = sr.Blockchain().GetGenesisHeader()
	}

	epoch := currentHeader.GetEpoch()
	shardId := sr.ShardCoordinator().SelfId()
	nodesCoordinatorShardID, err := sr.NodesCoordinator().ShardIdForEpoch(epoch)
	if err != nil {
		log.Debug("initCurrentRound.ShardIdForEpoch",
			"epoch", epoch,
			"error", err.Error())
		return
	}

	if shardId != nodesCoordinatorShardID {
		log.Debug("initCurrentRound.ShardIdForEpoch",
			"epoch", epoch,
			"shardCoordinator.ShardID", shardId,
			"nodesCoordinator.ShardID", nodesCoordinatorShardID)
		return
	}

	round := sr.RoundHandler().Index()

	roundInfo := &outportcore.RoundInfo{
		Round:            uint64(round),
		SignersIndexes:   make([]uint64, 0),
		BlockWasProposed: false,
		ShardId:          shardId,
		Epoch:            epoch,
		Timestamp:        uint64(sr.GetRoundTimeStamp().Unix()),
	}
	roundsInfo := &outportcore.RoundsInfo{
		ShardID:    shardId,
		RoundsInfo: []*outportcore.RoundInfo{roundInfo},
	}
	sr.outportHandler.SaveRoundsInfo(roundsInfo)
}

func (sr *subroundStartRound) generateNextConsensusGroup(roundIndex int64) error {
	currentHeader := sr.Blockchain().GetCurrentBlockHeader()
	if check.IfNil(currentHeader) {
		currentHeader = sr.Blockchain().GetGenesisHeader()
		if check.IfNil(currentHeader) {
			return spos.ErrNilHeader
		}
	}

	randomSeed := currentHeader.GetRandSeed()

	log.Debug("random source for the next consensus group",
		"rand", randomSeed)

	shardId := sr.ShardCoordinator().SelfId()

	leader, nextConsensusGroup, err := sr.GetNextConsensusGroup(
		randomSeed,
		uint64(sr.GetRoundIndex()),
		shardId,
		sr.NodesCoordinator(),
		currentHeader.GetEpoch(),
	)
	if err != nil {
		return err
	}

	log.Trace("consensus group is formed by next validators:",
		"round", roundIndex)

	for i := 0; i < len(nextConsensusGroup); i++ {
		log.Trace(core.GetTrimmedPk(hex.EncodeToString([]byte(nextConsensusGroup[i]))))
	}

	sr.SetConsensusGroup(nextConsensusGroup)
	sr.SetLeader(leader)

	consensusGroupSizeForEpoch := sr.NodesCoordinator().ConsensusGroupSizeForShardAndEpoch(shardId, currentHeader.GetEpoch())
	sr.SetConsensusGroupSize(consensusGroupSizeForEpoch)

	return nil
}

// EpochStartPrepare wis called when an epoch start event is observed, but not yet confirmed/committed.
// Some components may need to do initialisation on this event
func (sr *subroundStartRound) EpochStartPrepare(metaHdr data.HeaderHandler, _ data.BodyHandler) {
	log.Trace(fmt.Sprintf("epoch %d start prepare in consensus", metaHdr.GetEpoch()))
}

// EpochStartAction is called upon a start of epoch event.
func (sr *subroundStartRound) EpochStartAction(hdr data.HeaderHandler) {
	log.Trace(fmt.Sprintf("epoch %d start action in consensus", hdr.GetEpoch()))

	sr.changeEpoch(hdr.GetEpoch())
}

func (sr *subroundStartRound) changeEpoch(currentEpoch uint32) {
	epochNodes, err := sr.NodesCoordinator().GetConsensusWhitelistedNodes(currentEpoch)
	if err != nil {
		panic(fmt.Sprintf("consensus changing epoch failed with error %s", err.Error()))
	}

	sr.SetEligibleList(epochNodes)
}

// NotifyOrder returns the notification order for a start of epoch event
func (sr *subroundStartRound) NotifyOrder() uint32 {
	return common.ConsensusStartRoundOrder
}
