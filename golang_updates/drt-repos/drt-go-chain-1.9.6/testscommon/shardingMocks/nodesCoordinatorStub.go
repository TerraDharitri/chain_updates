package shardingMocks

import (
	"github.com/TerraDharitri/drt-go-chain-core/data"

	"github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
	"github.com/TerraDharitri/drt-go-chain/state"
)

// NodesCoordinatorStub -
type NodesCoordinatorStub struct {
	GetValidatorsPublicKeysCalled                    func(randomness []byte, round uint64, shardId uint32, epoch uint32) (string, []string, error)
	GetValidatorsRewardsAddressesCalled              func(randomness []byte, round uint64, shardId uint32, epoch uint32) ([]string, error)
	GetValidatorWithPublicKeyCalled                  func(publicKey []byte) (validator nodesCoordinator.Validator, shardId uint32, err error)
	GetAllValidatorsPublicKeysCalled                 func() (map[uint32][][]byte, error)
	GetAllWaitingValidatorsPublicKeysCalled          func(_ uint32) (map[uint32][][]byte, error)
	GetAllEligibleValidatorsPublicKeysCalled         func(epoch uint32) (map[uint32][][]byte, error)
	GetAllEligibleValidatorsPublicKeysForShardCalled func(epoch uint32, shardID uint32) ([]string, error)
	GetValidatorsIndexesCalled                       func(pubKeys []string, epoch uint32) ([]uint64, error)
	ConsensusGroupSizeCalled                         func(shardID uint32, epoch uint32) int
	ComputeConsensusGroupCalled                      func(randomness []byte, round uint64, shardId uint32, epoch uint32) (leader nodesCoordinator.Validator, validatorsGroup []nodesCoordinator.Validator, err error)
	EpochStartPrepareCalled                          func(metaHdr data.HeaderHandler, body data.BodyHandler)
	GetConsensusWhitelistedNodesCalled               func(epoch uint32) (map[string]struct{}, error)
	GetOwnPublicKeyCalled                            func() []byte
	GetWaitingEpochsLeftForPublicKeyCalled           func(publicKey []byte) (uint32, error)
	GetNumTotalEligibleCalled                        func() uint64
	ShardIdForEpochCalled                            func(epoch uint32) (uint32, error)
	GetCachedEpochsCalled                            func() map[uint32]struct{}
}

// NodesCoordinatorToRegistry -
func (ncm *NodesCoordinatorStub) NodesCoordinatorToRegistry(uint32) nodesCoordinator.NodesCoordinatorRegistryHandler {
	return nil
}

// EpochStartPrepare -
func (ncm *NodesCoordinatorStub) EpochStartPrepare(metaHdr data.HeaderHandler, body data.BodyHandler) {
	if ncm.EpochStartPrepareCalled != nil {
		ncm.EpochStartPrepareCalled(metaHdr, body)
	}
}

// GetChance -
func (ncm *NodesCoordinatorStub) GetChance(uint32) uint32 {
	return 1
}

// ValidatorsWeights -
func (ncm *NodesCoordinatorStub) ValidatorsWeights(_ []nodesCoordinator.Validator) ([]uint32, error) {
	return nil, nil
}

// GetAllLeavingValidatorsPublicKeys -
func (ncm *NodesCoordinatorStub) GetAllLeavingValidatorsPublicKeys(_ uint32) (map[uint32][][]byte, error) {
	return nil, nil
}

// SetConfig -
func (ncm *NodesCoordinatorStub) SetConfig(_ nodesCoordinator.NodesCoordinatorRegistryHandler) error {
	return nil
}

// ComputeAdditionalLeaving -
func (ncm *NodesCoordinatorStub) ComputeAdditionalLeaving(_ []*state.ShardValidatorInfo) (map[uint32][]nodesCoordinator.Validator, error) {
	return nil, nil
}

// GetAllEligibleValidatorsPublicKeys -
func (ncm *NodesCoordinatorStub) GetAllEligibleValidatorsPublicKeys(epoch uint32) (map[uint32][][]byte, error) {
	if ncm.GetAllEligibleValidatorsPublicKeysCalled != nil {
		return ncm.GetAllEligibleValidatorsPublicKeysCalled(epoch)
	}
	return nil, nil
}

// GetAllEligibleValidatorsPublicKeysForShard -
func (ncm *NodesCoordinatorStub) GetAllEligibleValidatorsPublicKeysForShard(epoch uint32, shardID uint32) ([]string, error) {
	if ncm.GetAllEligibleValidatorsPublicKeysForShardCalled != nil {
		return ncm.GetAllEligibleValidatorsPublicKeysForShardCalled(epoch, shardID)
	}
	return nil, nil
}

// GetAllWaitingValidatorsPublicKeys -
func (ncm *NodesCoordinatorStub) GetAllWaitingValidatorsPublicKeys(epoch uint32) (map[uint32][][]byte, error) {
	if ncm.GetAllWaitingValidatorsPublicKeysCalled != nil {
		return ncm.GetAllWaitingValidatorsPublicKeysCalled(epoch)
	}

	return nil, nil
}

// GetAllShuffledOutValidatorsPublicKeys -
func (ncm *NodesCoordinatorStub) GetAllShuffledOutValidatorsPublicKeys(_ uint32) (map[uint32][][]byte, error) {
	return nil, nil
}

// GetShuffledOutToAuctionValidatorsPublicKeys -
func (ncm *NodesCoordinatorStub) GetShuffledOutToAuctionValidatorsPublicKeys(_ uint32) (map[uint32][][]byte, error) {
	return nil, nil
}

// GetNumTotalEligible -
func (ncm *NodesCoordinatorStub) GetNumTotalEligible() uint64 {
	if ncm.GetNumTotalEligibleCalled != nil {
		return ncm.GetNumTotalEligibleCalled()
	}
	return 1
}

// GetAllValidatorsPublicKeys -
func (ncm *NodesCoordinatorStub) GetAllValidatorsPublicKeys(_ uint32) (map[uint32][][]byte, error) {
	if ncm.GetAllValidatorsPublicKeysCalled != nil {
		return ncm.GetAllValidatorsPublicKeysCalled()
	}

	return nil, nil
}

// GetValidatorsIndexes -
func (ncm *NodesCoordinatorStub) GetValidatorsIndexes(pubkeys []string, epoch uint32) ([]uint64, error) {
	if ncm.GetValidatorsIndexesCalled != nil {
		return ncm.GetValidatorsIndexesCalled(pubkeys, epoch)
	}
	return nil, nil
}

// ComputeConsensusGroup -
func (ncm *NodesCoordinatorStub) ComputeConsensusGroup(
	randomness []byte,
	round uint64,
	shardId uint32,
	epoch uint32,
) (leader nodesCoordinator.Validator, validatorsGroup []nodesCoordinator.Validator, err error) {
	if ncm.ComputeConsensusGroupCalled != nil {
		return ncm.ComputeConsensusGroupCalled(randomness, round, shardId, epoch)
	}

	return nil, nil, nil
}

// ConsensusGroupSizeForShardAndEpoch -
func (ncm *NodesCoordinatorStub) ConsensusGroupSizeForShardAndEpoch(shardID uint32, epoch uint32) int {
	if ncm.ConsensusGroupSizeCalled != nil {
		return ncm.ConsensusGroupSizeCalled(shardID, epoch)
	}
	return 1
}

// GetConsensusValidatorsPublicKeys -
func (ncm *NodesCoordinatorStub) GetConsensusValidatorsPublicKeys(
	randomness []byte,
	round uint64,
	shardId uint32,
	epoch uint32,
) (string, []string, error) {
	if ncm.GetValidatorsPublicKeysCalled != nil {
		return ncm.GetValidatorsPublicKeysCalled(randomness, round, shardId, epoch)
	}

	return "", nil, nil
}

// SetNodesPerShards -
func (ncm *NodesCoordinatorStub) SetNodesPerShards(_ map[uint32][]nodesCoordinator.Validator, _ map[uint32][]nodesCoordinator.Validator, _ []nodesCoordinator.Validator, _ uint32) error {
	return nil
}

// LoadState -
func (ncm *NodesCoordinatorStub) LoadState(_ []byte) error {
	return nil
}

// GetSavedStateKey -
func (ncm *NodesCoordinatorStub) GetSavedStateKey() []byte {
	return []byte("key")
}

// ShardIdForEpoch returns the nodesCoordinator configured ShardId for specified epoch if epoch configuration exists,
// otherwise error
func (ncm *NodesCoordinatorStub) ShardIdForEpoch(epoch uint32) (uint32, error) {

	if ncm.ShardIdForEpochCalled != nil {
		return ncm.ShardIdForEpochCalled(epoch)
	}
	return 0, nil
}

// ShuffleOutForEpoch verifies if the shards changed in the new epoch and calls the shuffleOutHandler
func (ncm *NodesCoordinatorStub) ShuffleOutForEpoch(_ uint32) {
}

// GetConsensusWhitelistedNodes return the whitelisted nodes allowed to send consensus messages, for each of the shards
func (ncm *NodesCoordinatorStub) GetConsensusWhitelistedNodes(epoch uint32) (map[string]struct{}, error) {
	if ncm.GetConsensusWhitelistedNodesCalled != nil {
		return ncm.GetConsensusWhitelistedNodesCalled(epoch)
	}
	return nil, nil
}

// GetSelectedPublicKeys -
func (ncm *NodesCoordinatorStub) GetSelectedPublicKeys(_ []byte, _ uint32, _ uint32) ([]string, error) {
	panic("implement me")
}

// GetValidatorWithPublicKey -
func (ncm *NodesCoordinatorStub) GetValidatorWithPublicKey(publicKey []byte) (nodesCoordinator.Validator, uint32, error) {
	if ncm.GetValidatorWithPublicKeyCalled != nil {
		return ncm.GetValidatorWithPublicKeyCalled(publicKey)
	}
	return nil, 0, nil
}

// GetOwnPublicKey -
func (ncm *NodesCoordinatorStub) GetOwnPublicKey() []byte {
	if ncm.GetOwnPublicKeyCalled != nil {
		return ncm.GetOwnPublicKeyCalled()
	}
	return []byte("key")
}

// GetWaitingEpochsLeftForPublicKey -
func (ncm *NodesCoordinatorStub) GetWaitingEpochsLeftForPublicKey(publicKey []byte) (uint32, error) {
	if ncm.GetWaitingEpochsLeftForPublicKeyCalled != nil {
		return ncm.GetWaitingEpochsLeftForPublicKeyCalled(publicKey)
	}
	return 0, nil
}

// GetCachedEpochs -
func (ncm *NodesCoordinatorStub) GetCachedEpochs() map[uint32]struct{} {
	if ncm.GetCachedEpochsCalled != nil {
		return ncm.GetCachedEpochsCalled()
	}
	return make(map[uint32]struct{})
}

// IsInterfaceNil returns true if there is no value under the interface
func (ncm *NodesCoordinatorStub) IsInterfaceNil() bool {
	return ncm == nil
}
