package sync

import (
	"bytes"
	"math"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/consensus"
	"github.com/TerraDharitri/drt-go-chain/process"
)

var _ process.ForkDetector = (*shardForkDetector)(nil)

// shardForkDetector implements the shard fork detector mechanism
type shardForkDetector struct {
	*baseForkDetector
}

// NewShardForkDetector method creates a new shardForkDetector object
func NewShardForkDetector(
	roundHandler consensus.RoundHandler,
	blackListHandler process.TimeCacher,
	blockTracker process.BlockTracker,
	genesisTime int64,
	enableEpochsHandler common.EnableEpochsHandler,
	proofsPool process.ProofsPool,
) (*shardForkDetector, error) {

	if check.IfNil(roundHandler) {
		return nil, process.ErrNilRoundHandler
	}
	if check.IfNil(blackListHandler) {
		return nil, process.ErrNilBlackListCacher
	}
	if check.IfNil(blockTracker) {
		return nil, process.ErrNilBlockTracker
	}
	if check.IfNil(enableEpochsHandler) {
		return nil, process.ErrNilEnableEpochsHandler
	}
	if check.IfNil(proofsPool) {
		return nil, process.ErrNilProofsPool
	}

	genesisHdr, _, err := blockTracker.GetSelfNotarizedHeader(core.MetachainShardId, 0)
	if err != nil {
		return nil, err
	}

	bfd := &baseForkDetector{
		roundHandler:        roundHandler,
		blackListHandler:    blackListHandler,
		genesisTime:         genesisTime,
		blockTracker:        blockTracker,
		genesisNonce:        genesisHdr.GetNonce(),
		genesisRound:        genesisHdr.GetRound(),
		genesisEpoch:        genesisHdr.GetEpoch(),
		enableEpochsHandler: enableEpochsHandler,
		proofsPool:          proofsPool,
	}

	bfd.headers = make(map[uint64][]*headerInfo)
	bfd.fork.checkpoint = make([]*checkpointInfo, 0)
	checkpoint := &checkpointInfo{
		nonce: bfd.genesisNonce,
		round: bfd.genesisRound,
	}
	bfd.setFinalCheckpoint(checkpoint)
	bfd.addCheckpoint(checkpoint)
	bfd.fork.rollBackNonce = math.MaxUint64
	bfd.fork.probableHighestNonce = bfd.genesisNonce
	bfd.fork.highestNonceReceived = bfd.genesisNonce

	sfd := shardForkDetector{
		baseForkDetector: bfd,
	}

	sfd.blockTracker.RegisterSelfNotarizedFromCrossHeadersHandler(sfd.ReceivedSelfNotarizedFromCrossHeaders)

	bfd.forkDetector = &sfd

	return &sfd, nil
}

// AddHeader method adds a new header to headers map
func (sfd *shardForkDetector) AddHeader(
	header data.HeaderHandler,
	headerHash []byte,
	state process.BlockHeaderState,
	selfNotarizedHeaders []data.HeaderHandler,
	selfNotarizedHeadersHashes [][]byte,
) error {
	return sfd.addHeader(
		header,
		headerHash,
		state,
		selfNotarizedHeaders,
		selfNotarizedHeadersHashes,
		sfd.doJobOnBHProcessed,
	)
}

func (sfd *shardForkDetector) doJobOnBHProcessed(
	header data.HeaderHandler,
	headerHash []byte,
	selfNotarizedHeaders []data.HeaderHandler,
	selfNotarizedHeadersHashes [][]byte,
) {
	_ = sfd.appendSelfNotarizedHeaders(selfNotarizedHeaders, selfNotarizedHeadersHashes, core.MetachainShardId)
	sfd.computeFinalCheckpoint()
	newCheckpoint := &checkpointInfo{nonce: header.GetNonce(), round: header.GetRound(), hash: headerHash}
	sfd.addCheckpoint(newCheckpoint)
	// first shard block with proof does not have increased consensus
	// so instant finality will only be set after the first block with increased consensus
	if common.IsFlagEnabledAfterEpochsStartBlock(header, sfd.enableEpochsHandler, common.AndromedaFlag) {
		sfd.setFinalCheckpoint(newCheckpoint)
	}
	sfd.removePastOrInvalidRecords()
}

// ReceivedSelfNotarizedFromCrossHeaders is a registered call handler through which fork detector is notified
// when metachain notarized new headers from self shard
func (sfd *shardForkDetector) ReceivedSelfNotarizedFromCrossHeaders(
	shardID uint32,
	selfNotarizedHeaders []data.HeaderHandler,
	selfNotarizedHeadersHashes [][]byte,
) {
	// accept only self notarized headers by meta
	if shardID != core.MetachainShardId {
		return
	}

	appended := sfd.appendSelfNotarizedHeaders(selfNotarizedHeaders, selfNotarizedHeadersHashes, shardID)
	if appended {
		sfd.computeFinalCheckpoint()
	}
}

func (sfd *shardForkDetector) appendSelfNotarizedHeaders(
	selfNotarizedHeaders []data.HeaderHandler,
	selfNotarizedHeadersHashes [][]byte,
	shardID uint32,
) bool {

	selfNotarizedHeaderAdded := false
	finalNonce := sfd.finalCheckpoint().nonce

	for i := 0; i < len(selfNotarizedHeaders); i++ {
		if selfNotarizedHeaders[i].GetNonce() <= finalNonce {
			continue
		}

		hasProof := sfd.proofsPool.HasProof(selfNotarizedHeaders[i].GetShardID(), selfNotarizedHeadersHashes[i])
		appended := sfd.append(&headerInfo{
			nonce:    selfNotarizedHeaders[i].GetNonce(),
			round:    selfNotarizedHeaders[i].GetRound(),
			hash:     selfNotarizedHeadersHashes[i],
			state:    process.BHNotarized,
			hasProof: hasProof,
		})
		if appended {
			log.Debug("added self notarized header in fork detector",
				"notarized by shard", shardID,
				"round", selfNotarizedHeaders[i].GetRound(),
				"nonce", selfNotarizedHeaders[i].GetNonce(),
				"hash", selfNotarizedHeadersHashes[i])

			selfNotarizedHeaderAdded = true
		}
	}

	return selfNotarizedHeaderAdded
}

func (sfd *shardForkDetector) computeFinalCheckpoint() {
	finalCheckpoint := &checkpointInfo{}
	finalCheckpointWasSet := false

	sfd.mutHeaders.RLock()
	for nonce, headersInfo := range sfd.headers {
		if finalCheckpoint.nonce >= nonce {
			continue
		}

		indexBHProcessed, indexBHNotarized := sfd.getProcessedAndNotarizedIndexes(headersInfo)
		isProcessedBlockAlreadyNotarized := indexBHProcessed != -1 && indexBHNotarized != -1
		if !isProcessedBlockAlreadyNotarized {
			continue
		}

		sameHash := bytes.Equal(headersInfo[indexBHNotarized].hash, headersInfo[indexBHProcessed].hash)
		if !sameHash {
			continue
		}

		finalCheckpoint = &checkpointInfo{
			nonce: nonce,
			round: headersInfo[indexBHNotarized].round,
			hash:  headersInfo[indexBHNotarized].hash,
		}

		finalCheckpointWasSet = true
	}
	sfd.mutHeaders.RUnlock()

	if finalCheckpointWasSet {
		sfd.setFinalCheckpoint(finalCheckpoint)
	}
}

func (sfd *shardForkDetector) getProcessedAndNotarizedIndexes(headersInfo []*headerInfo) (int, int) {
	indexBHProcessed := -1
	indexBHNotarized := -1

	for index, hdrInfo := range headersInfo {
		switch hdrInfo.state {
		case process.BHProcessed:
			indexBHProcessed = index
		case process.BHNotarized:
			indexBHNotarized = index
		}
	}

	return indexBHProcessed, indexBHNotarized
}

// IsInterfaceNil returns true if there is no value under the interface
func (sfd *shardForkDetector) IsInterfaceNil() bool {
	return sfd == nil
}
