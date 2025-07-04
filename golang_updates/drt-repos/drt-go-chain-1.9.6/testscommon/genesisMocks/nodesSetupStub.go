package genesisMocks

import (
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
)

// NodesSetupStub -
type NodesSetupStub struct {
	InitialNodesPubKeysCalled                 func() map[uint32][]string
	InitialEligibleNodesPubKeysForShardCalled func(shardId uint32) ([]string, error)
	GetShardIDForPubKeyCalled                 func(pubKey []byte) (uint32, error)
	NumberOfShardsCalled                      func() uint32
	GetShardConsensusGroupSizeCalled          func() uint32
	GetMetaConsensusGroupSizeCalled           func() uint32
	GetRoundDurationCalled                    func() uint64
	MinNumberOfMetaNodesCalled                func() uint32
	MinNumberOfShardNodesCalled               func() uint32
	GetHysteresisCalled                       func() float32
	GetAdaptivityCalled                       func() bool
	InitialNodesInfoForShardCalled            func(shardId uint32) ([]nodesCoordinator.GenesisNodeInfoHandler, []nodesCoordinator.GenesisNodeInfoHandler, error)
	InitialNodesInfoCalled                    func() (map[uint32][]nodesCoordinator.GenesisNodeInfoHandler, map[uint32][]nodesCoordinator.GenesisNodeInfoHandler)
	GetStartTimeCalled                        func() int64
	MinNumberOfNodesCalled                    func() uint32
	AllInitialNodesCalled                     func() []nodesCoordinator.GenesisNodeInfoHandler
	MinNumberOfNodesWithHysteresisCalled      func() uint32
	MinShardHysteresisNodesCalled             func() uint32
	MinMetaHysteresisNodesCalled              func() uint32
	GetChainIdCalled                          func() string
	GetMinTransactionVersionCalled            func() uint32
	ExportNodesConfigCalled                   func() config.NodesConfig
}

// InitialNodesPubKeys -
func (n *NodesSetupStub) InitialNodesPubKeys() map[uint32][]string {
	if n.InitialNodesPubKeysCalled != nil {
		return n.InitialNodesPubKeysCalled()
	}

	return map[uint32][]string{0: {"val1", "val2"}}
}

// InitialEligibleNodesPubKeysForShard -
func (n *NodesSetupStub) InitialEligibleNodesPubKeysForShard(shardId uint32) ([]string, error) {
	if n.InitialEligibleNodesPubKeysForShardCalled != nil {
		return n.InitialEligibleNodesPubKeysForShardCalled(shardId)
	}

	return []string{"val1", "val2"}, nil
}

// NumberOfShards -
func (n *NodesSetupStub) NumberOfShards() uint32 {
	if n.NumberOfShardsCalled != nil {
		return n.NumberOfShardsCalled()
	}
	return 1
}

// GetShardIDForPubKey -
func (n *NodesSetupStub) GetShardIDForPubKey(pubkey []byte) (uint32, error) {
	if n.GetShardIDForPubKeyCalled != nil {
		return n.GetShardIDForPubKeyCalled(pubkey)
	}
	return 0, nil
}

// GetShardConsensusGroupSize -
func (n *NodesSetupStub) GetShardConsensusGroupSize() uint32 {
	if n.GetShardConsensusGroupSizeCalled != nil {
		return n.GetShardConsensusGroupSizeCalled()
	}
	return 1
}

// GetMetaConsensusGroupSize -
func (n *NodesSetupStub) GetMetaConsensusGroupSize() uint32 {
	if n.GetMetaConsensusGroupSizeCalled != nil {
		return n.GetMetaConsensusGroupSizeCalled()
	}
	return 1
}

// GetRoundDuration -
func (n *NodesSetupStub) GetRoundDuration() uint64 {
	if n.GetRoundDurationCalled != nil {
		return n.GetRoundDurationCalled()
	}
	return 4000
}

// MinNumberOfMetaNodes -
func (n *NodesSetupStub) MinNumberOfMetaNodes() uint32 {
	if n.MinNumberOfMetaNodesCalled != nil {
		return n.MinNumberOfMetaNodesCalled()
	}
	return 1
}

// MinNumberOfShardNodes -
func (n *NodesSetupStub) MinNumberOfShardNodes() uint32 {
	if n.MinNumberOfShardNodesCalled != nil {
		return n.MinNumberOfShardNodesCalled()
	}
	return 1
}

// GetHysteresis -
func (n *NodesSetupStub) GetHysteresis() float32 {
	if n.GetHysteresisCalled != nil {
		return n.GetHysteresisCalled()
	}
	return 0
}

// GetAdaptivity -
func (n *NodesSetupStub) GetAdaptivity() bool {
	if n.GetAdaptivityCalled != nil {
		return n.GetAdaptivityCalled()
	}
	return false
}

// InitialNodesInfoForShard -
func (n *NodesSetupStub) InitialNodesInfoForShard(
	shardId uint32,
) ([]nodesCoordinator.GenesisNodeInfoHandler, []nodesCoordinator.GenesisNodeInfoHandler, error) {
	if n.InitialNodesInfoForShardCalled != nil {
		return n.InitialNodesInfoForShardCalled(shardId)
	}

	return nil, nil, nil
}

// InitialNodesInfo -
func (n *NodesSetupStub) InitialNodesInfo() (map[uint32][]nodesCoordinator.GenesisNodeInfoHandler, map[uint32][]nodesCoordinator.GenesisNodeInfoHandler) {
	if n.InitialNodesInfoCalled != nil {
		return n.InitialNodesInfoCalled()
	}

	return nil, nil
}

// GetStartTime -
func (n *NodesSetupStub) GetStartTime() int64 {
	if n.GetStartTimeCalled != nil {
		return n.GetStartTimeCalled()
	}
	return 0
}

// MinNumberOfNodes -
func (n *NodesSetupStub) MinNumberOfNodes() uint32 {
	if n.MinNumberOfNodesCalled != nil {
		return n.MinNumberOfNodesCalled()
	}
	return 1
}

// MinNumberOfNodesWithHysteresis -
func (n *NodesSetupStub) MinNumberOfNodesWithHysteresis() uint32 {
	if n.MinNumberOfNodesWithHysteresisCalled != nil {
		return n.MinNumberOfNodesWithHysteresisCalled()
	}
	return n.MinNumberOfNodes()
}

// AllInitialNodes -
func (n *NodesSetupStub) AllInitialNodes() []nodesCoordinator.GenesisNodeInfoHandler {
	if n.AllInitialNodesCalled != nil {
		return n.AllInitialNodesCalled()
	}
	return nil
}

// GetChainId -
func (n *NodesSetupStub) GetChainId() string {
	if n.GetChainIdCalled != nil {
		return n.GetChainIdCalled()
	}
	return "chainID"
}

// GetMinTransactionVersion -
func (n *NodesSetupStub) GetMinTransactionVersion() uint32 {
	if n.GetMinTransactionVersionCalled != nil {
		return n.GetMinTransactionVersionCalled()
	}
	return 1
}

// MinShardHysteresisNodes -
func (n *NodesSetupStub) MinShardHysteresisNodes() uint32 {
	if n.MinShardHysteresisNodesCalled != nil {
		return n.MinShardHysteresisNodesCalled()
	}
	return 1
}

// MinMetaHysteresisNodes -
func (n *NodesSetupStub) MinMetaHysteresisNodes() uint32 {
	if n.MinMetaHysteresisNodesCalled != nil {
		return n.MinMetaHysteresisNodesCalled()
	}
	return 1
}

// ExportNodesConfig -
func (n *NodesSetupStub) ExportNodesConfig() config.NodesConfig {
	if n.ExportNodesConfigCalled != nil {
		return n.ExportNodesConfigCalled()
	}

	return config.NodesConfig{}
}

// IsInterfaceNil -
func (n *NodesSetupStub) IsInterfaceNil() bool {
	return n == nil
}
