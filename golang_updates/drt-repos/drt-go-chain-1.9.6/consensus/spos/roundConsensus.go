package spos

import (
	"sync"

	"github.com/TerraDharitri/drt-go-chain-core/core/check"

	"github.com/TerraDharitri/drt-go-chain/consensus"
)

// roundConsensus defines the data needed by spos to do the consensus in each round
type roundConsensus struct {
	eligibleNodes        map[string]struct{}
	mutEligible          sync.RWMutex
	consensusGroup       []string
	leader               string
	consensusGroupSize   int
	selfPubKey           string
	validatorRoundStates map[string]*roundState
	mut                  sync.RWMutex
	keysHandler          consensus.KeysHandler
}

// NewRoundConsensus creates a new roundConsensus object
func NewRoundConsensus(
	eligibleNodes map[string]struct{},
	consensusGroupSize int,
	selfId string,
	keysHandler consensus.KeysHandler,
) (*roundConsensus, error) {
	if check.IfNil(keysHandler) {
		return nil, ErrNilKeysHandler
	}

	return &roundConsensus{
		eligibleNodes:        eligibleNodes,
		consensusGroupSize:   consensusGroupSize,
		selfPubKey:           selfId,
		mutEligible:          sync.RWMutex{},
		validatorRoundStates: make(map[string]*roundState),
		keysHandler:          keysHandler,
	}, nil
}

// ConsensusGroupIndex returns the index of given public key in the current consensus group
func (rcns *roundConsensus) ConsensusGroupIndex(pubKey string) (int, error) {
	for i, pk := range rcns.consensusGroup {
		if pk == pubKey {
			return i, nil
		}
	}
	return 0, ErrNotFoundInConsensus
}

// SelfConsensusGroupIndex returns the index of self public key in current consensus group
func (rcns *roundConsensus) SelfConsensusGroupIndex() (int, error) {
	return rcns.ConsensusGroupIndex(rcns.selfPubKey)
}

// SetEligibleList sets the eligible list ID's
func (rcns *roundConsensus) SetEligibleList(eligibleList map[string]struct{}) {
	rcns.mutEligible.Lock()
	rcns.eligibleNodes = eligibleList
	rcns.mutEligible.Unlock()
}

// ConsensusGroup returns the consensus group ID's
func (rcns *roundConsensus) ConsensusGroup() []string {
	rcns.mut.RLock()
	defer rcns.mut.RUnlock()

	return rcns.consensusGroup
}

// SetConsensusGroup sets the consensus group ID's
func (rcns *roundConsensus) SetConsensusGroup(consensusGroup []string) {
	rcns.mut.Lock()

	rcns.consensusGroup = consensusGroup

	rcns.validatorRoundStates = make(map[string]*roundState)

	for i := 0; i < len(consensusGroup); i++ {
		rcns.validatorRoundStates[rcns.consensusGroup[i]] = NewRoundState()
	}

	rcns.mut.Unlock()
}

// Leader returns the leader for the current consensus
func (rcns *roundConsensus) Leader() string {
	rcns.mut.RLock()
	defer rcns.mut.RUnlock()

	return rcns.leader
}

// SetLeader sets the leader for the current consensus
func (rcns *roundConsensus) SetLeader(leader string) {
	rcns.mut.Lock()
	defer rcns.mut.Unlock()

	rcns.leader = leader
}

// ConsensusGroupSize returns the consensus group size
func (rcns *roundConsensus) ConsensusGroupSize() int {
	return rcns.consensusGroupSize
}

// SetConsensusGroupSize sets the consensus group size
func (rcns *roundConsensus) SetConsensusGroupSize(consensusGroupSize int) {
	rcns.consensusGroupSize = consensusGroupSize
}

// SelfPubKey returns selfPubKey ID
func (rcns *roundConsensus) SelfPubKey() string {
	return rcns.selfPubKey
}

// SetSelfPubKey sets selfPubKey ID
func (rcns *roundConsensus) SetSelfPubKey(selfPubKey string) {
	rcns.selfPubKey = selfPubKey
}

// JobDone returns the state of the action done, by the node represented by the key parameter,
// in subround given by the subroundId parameter
func (rcns *roundConsensus) JobDone(key string, subroundId int) (bool, error) {
	rcns.mut.RLock()
	currentRoundState := rcns.validatorRoundStates[key]

	if currentRoundState == nil {
		rcns.mut.RUnlock()
		return false, ErrInvalidKey
	}

	retcode := currentRoundState.JobDone(subroundId)
	rcns.mut.RUnlock()

	return retcode, nil
}

// SetJobDone set the state of the action done, by the node represented by the key parameter,
// in subround given by the subroundId parameter
func (rcns *roundConsensus) SetJobDone(key string, subroundId int, value bool) error {
	rcns.mut.Lock()

	currentRoundState := rcns.validatorRoundStates[key]

	if currentRoundState == nil {
		rcns.mut.Unlock()
		return ErrInvalidKey
	}

	currentRoundState.SetJobDone(subroundId, value)
	rcns.mut.Unlock()

	return nil
}

// SelfJobDone returns the self state of the action done in subround given by the subroundId parameter
func (rcns *roundConsensus) SelfJobDone(subroundId int) (bool, error) {
	return rcns.JobDone(rcns.selfPubKey, subroundId)
}

// IsNodeInConsensusGroup method checks if the node is part of consensus group of the current round
func (rcns *roundConsensus) IsNodeInConsensusGroup(node string) bool {
	rcns.mut.RLock()
	defer rcns.mut.RUnlock()

	for i := 0; i < len(rcns.consensusGroup); i++ {
		if rcns.consensusGroup[i] == node {
			return true
		}
	}

	return false
}

// IsNodeInEligibleList method checks if the node is part of the eligible list
func (rcns *roundConsensus) IsNodeInEligibleList(node string) bool {
	rcns.mutEligible.RLock()
	_, ok := rcns.eligibleNodes[node]
	rcns.mutEligible.RUnlock()

	return ok
}

// ComputeSize method returns the number of messages received from the nodes belonging to the current jobDone group
// related to this subround
func (rcns *roundConsensus) ComputeSize(subroundId int) int {
	n := 0

	for i := 0; i < len(rcns.consensusGroup); i++ {
		isJobDone, err := rcns.JobDone(rcns.consensusGroup[i], subroundId)
		if err != nil {
			log.Debug("JobDone", "error", err.Error())
			continue
		}

		if isJobDone {
			n++
		}
	}

	return n
}

// ResetRoundState method resets the state of each node from the current jobDone group, regarding the
// consensus validatorRoundStates
func (rcns *roundConsensus) ResetRoundState() {
	rcns.mut.Lock()

	var currentRoundState *roundState
	for i := 0; i < len(rcns.consensusGroup); i++ {
		currentRoundState = rcns.validatorRoundStates[rcns.consensusGroup[i]]
		if currentRoundState == nil {
			log.Debug("validatorRoundStates", "error", ErrNilRoundState.Error())
			continue
		}

		currentRoundState.ResetJobsDone()
	}

	rcns.mut.Unlock()
}

// IsMultiKeyInConsensusGroup method checks if one of the nodes which are controlled by this instance
// is in consensus group in the current round
func (rcns *roundConsensus) IsMultiKeyInConsensusGroup() bool {
	for i := 0; i < len(rcns.consensusGroup); i++ {
		if rcns.IsKeyManagedBySelf([]byte(rcns.consensusGroup[i])) {
			return true
		}
	}

	return false
}

// IsKeyManagedBySelf returns true if the key is managed by the current node
func (rcns *roundConsensus) IsKeyManagedBySelf(pkBytes []byte) bool {
	return rcns.keysHandler.IsKeyManagedByCurrentNode(pkBytes)
}

// IncrementRoundsWithoutReceivedMessages increments the number of rounds without received messages on a provided public key
func (rcns *roundConsensus) IncrementRoundsWithoutReceivedMessages(pkBytes []byte) {
	rcns.keysHandler.IncrementRoundsWithoutReceivedMessages(pkBytes)
}

// GetKeysHandler returns the keysHandler instance
func (rcns *roundConsensus) GetKeysHandler() consensus.KeysHandler {
	return rcns.keysHandler
}
