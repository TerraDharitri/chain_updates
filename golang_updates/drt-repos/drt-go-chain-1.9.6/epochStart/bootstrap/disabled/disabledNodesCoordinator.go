package disabled

import (
	nodesCoord "github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
	"github.com/TerraDharitri/drt-go-chain/state"
)

// nodesCoordinator -
type nodesCoordinator struct {
}

// NewNodesCoordinator returns a new instance of nodesCoordinator
func NewNodesCoordinator() *nodesCoordinator {
	return &nodesCoordinator{}
}

// GetChance -
func (n *nodesCoordinator) GetChance(uint32) uint32 {
	return 1
}

// GetAllLeavingValidatorsPublicKeys -
func (n *nodesCoordinator) GetAllLeavingValidatorsPublicKeys(_ uint32) (map[uint32][][]byte, error) {
	return nil, nil
}

// ValidatorsWeights -
func (n *nodesCoordinator) ValidatorsWeights(validators []nodesCoord.Validator) ([]uint32, error) {
	return make([]uint32, len(validators)), nil
}

// ComputeAdditionalLeaving -
func (n *nodesCoordinator) ComputeAdditionalLeaving(_ []*state.ShardValidatorInfo) (map[uint32][]nodesCoord.Validator, error) {
	return nil, nil
}

// GetValidatorsIndexes -
func (n *nodesCoordinator) GetValidatorsIndexes(_ []string, _ uint32) ([]uint64, error) {
	return nil, nil
}

// GetAllEligibleValidatorsPublicKeys -
func (n *nodesCoordinator) GetAllEligibleValidatorsPublicKeys(_ uint32) (map[uint32][][]byte, error) {
	return nil, nil
}

// GetAllEligibleValidatorsPublicKeysForShard -
func (n *nodesCoordinator) GetAllEligibleValidatorsPublicKeysForShard(_ uint32, _ uint32) ([]string, error) {
	return nil, nil
}

// GetAllWaitingValidatorsPublicKeys -
func (n *nodesCoordinator) GetAllWaitingValidatorsPublicKeys(_ uint32) (map[uint32][][]byte, error) {
	return nil, nil
}

// GetAllShuffledOutValidatorsPublicKeys -
func (n *nodesCoordinator) GetAllShuffledOutValidatorsPublicKeys(_ uint32) (map[uint32][][]byte, error) {
	return nil, nil
}

// GetShuffledOutToAuctionValidatorsPublicKeys -
func (n *nodesCoordinator) GetShuffledOutToAuctionValidatorsPublicKeys(_ uint32) (map[uint32][][]byte, error) {
	return nil, nil
}

// GetConsensusValidatorsPublicKeys -
func (n *nodesCoordinator) GetConsensusValidatorsPublicKeys(_ []byte, _ uint64, _ uint32, _ uint32) (string, []string, error) {
	return "", nil, nil
}

// GetOwnPublicKey -
func (n *nodesCoordinator) GetOwnPublicKey() []byte {
	return nil
}

// ComputeConsensusGroup -
func (n *nodesCoordinator) ComputeConsensusGroup(_ []byte, _ uint64, _ uint32, _ uint32) (leader nodesCoord.Validator, validatorsGroup []nodesCoord.Validator, err error) {
	return nil, nil, nil
}

// GetValidatorWithPublicKey -
func (n *nodesCoordinator) GetValidatorWithPublicKey(_ []byte) (validator nodesCoord.Validator, shardId uint32, err error) {
	return nil, 0, nil
}

// LoadState -
func (n *nodesCoordinator) LoadState(_ []byte) error {
	return nil
}

// GetSavedStateKey -
func (n *nodesCoordinator) GetSavedStateKey() []byte {
	return nil
}

// ShardIdForEpoch -
func (n *nodesCoordinator) ShardIdForEpoch(_ uint32) (uint32, error) {
	return 0, nil
}

// ShuffleOutForEpoch verifies if the shards changed in the new epoch and calls the shuffleOutHandler
func (n *nodesCoordinator) ShuffleOutForEpoch(_ uint32) {
}

// GetConsensusWhitelistedNodes -
func (n *nodesCoordinator) GetConsensusWhitelistedNodes(_ uint32) (map[string]struct{}, error) {
	return nil, nil
}

// ConsensusGroupSizeForShardAndEpoch -
func (n *nodesCoordinator) ConsensusGroupSizeForShardAndEpoch(uint32, uint32) int {
	return 0
}

// GetNumTotalEligible -
func (n *nodesCoordinator) GetNumTotalEligible() uint64 {
	return 0
}

// GetWaitingEpochsLeftForPublicKey returns 0
func (n *nodesCoordinator) GetWaitingEpochsLeftForPublicKey(_ []byte) (uint32, error) {
	return 0, nil
}

// GetCachedEpochs returns an empty map
func (n *nodesCoordinator) GetCachedEpochs() map[uint32]struct{} {
	return make(map[uint32]struct{})
}

// IsInterfaceNil -
func (n *nodesCoordinator) IsInterfaceNil() bool {
	return n == nil
}
