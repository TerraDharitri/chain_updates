package chainSimulator

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	chainData "github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain/api/shared"
	"github.com/TerraDharitri/drt-go-chain/consensus"
	"github.com/TerraDharitri/drt-go-chain/factory"
	"github.com/TerraDharitri/drt-go-chain/node/chainSimulator/dtos"
	"github.com/TerraDharitri/drt-go-chain/sharding"
)

// NodeHandlerMock -
type NodeHandlerMock struct {
	GetProcessComponentsCalled    func() factory.ProcessComponentsHolder
	GetChainHandlerCalled         func() chainData.ChainHandler
	GetBroadcastMessengerCalled   func() consensus.BroadcastMessenger
	GetShardCoordinatorCalled     func() sharding.Coordinator
	GetCryptoComponentsCalled     func() factory.CryptoComponentsHolder
	GetCoreComponentsCalled       func() factory.CoreComponentsHolder
	GetDataComponentsCalled       func() factory.DataComponentsHandler
	GetStateComponentsCalled      func() factory.StateComponentsHolder
	GetFacadeHandlerCalled        func() shared.FacadeHandler
	GetStatusCoreComponentsCalled func() factory.StatusCoreComponentsHolder
	GetNetworkComponentsCalled    func() factory.NetworkComponentsHolder
	SetKeyValueForAddressCalled   func(addressBytes []byte, state map[string]string) error
	SetStateForAddressCalled      func(address []byte, state *dtos.AddressState) error
	RemoveAccountCalled           func(address []byte) error
	GetBasePeersCalled            func() map[uint32]core.PeerID
	SetBasePeersCalled            func(basePeers map[uint32]core.PeerID)
	CloseCalled                   func() error
}

// ForceChangeOfEpoch -
func (mock *NodeHandlerMock) ForceChangeOfEpoch() error {
	return nil
}

// GetProcessComponents -
func (mock *NodeHandlerMock) GetProcessComponents() factory.ProcessComponentsHolder {
	if mock.GetProcessComponentsCalled != nil {
		return mock.GetProcessComponentsCalled()
	}
	return nil
}

// GetChainHandler -
func (mock *NodeHandlerMock) GetChainHandler() chainData.ChainHandler {
	if mock.GetChainHandlerCalled != nil {
		return mock.GetChainHandlerCalled()
	}
	return nil
}

// GetBroadcastMessenger -
func (mock *NodeHandlerMock) GetBroadcastMessenger() consensus.BroadcastMessenger {
	if mock.GetBroadcastMessengerCalled != nil {
		return mock.GetBroadcastMessengerCalled()
	}
	return nil
}

// GetShardCoordinator -
func (mock *NodeHandlerMock) GetShardCoordinator() sharding.Coordinator {
	if mock.GetShardCoordinatorCalled != nil {
		return mock.GetShardCoordinatorCalled()
	}
	return nil
}

// GetCryptoComponents -
func (mock *NodeHandlerMock) GetCryptoComponents() factory.CryptoComponentsHolder {
	if mock.GetCryptoComponentsCalled != nil {
		return mock.GetCryptoComponentsCalled()
	}
	return nil
}

// GetCoreComponents -
func (mock *NodeHandlerMock) GetCoreComponents() factory.CoreComponentsHolder {
	if mock.GetCoreComponentsCalled != nil {
		return mock.GetCoreComponentsCalled()
	}
	return nil
}

// GetDataComponents -
func (mock *NodeHandlerMock) GetDataComponents() factory.DataComponentsHolder {
	if mock.GetDataComponentsCalled != nil {
		return mock.GetDataComponentsCalled()
	}
	return nil
}

// GetStateComponents -
func (mock *NodeHandlerMock) GetStateComponents() factory.StateComponentsHolder {
	if mock.GetStateComponentsCalled != nil {
		return mock.GetStateComponentsCalled()
	}
	return nil
}

// GetFacadeHandler -
func (mock *NodeHandlerMock) GetFacadeHandler() shared.FacadeHandler {
	if mock.GetFacadeHandlerCalled != nil {
		return mock.GetFacadeHandlerCalled()
	}
	return nil
}

// GetStatusCoreComponents -
func (mock *NodeHandlerMock) GetStatusCoreComponents() factory.StatusCoreComponentsHolder {
	if mock.GetStatusCoreComponentsCalled != nil {
		return mock.GetStatusCoreComponentsCalled()
	}
	return nil
}

// GetNetworkComponents -
func (mock *NodeHandlerMock) GetNetworkComponents() factory.NetworkComponentsHolder {
	if mock.GetNetworkComponentsCalled != nil {
		return mock.GetNetworkComponentsCalled()
	}
	return nil
}

// SetKeyValueForAddress -
func (mock *NodeHandlerMock) SetKeyValueForAddress(addressBytes []byte, state map[string]string) error {
	if mock.SetKeyValueForAddressCalled != nil {
		return mock.SetKeyValueForAddressCalled(addressBytes, state)
	}
	return nil
}

// SetStateForAddress -
func (mock *NodeHandlerMock) SetStateForAddress(address []byte, state *dtos.AddressState) error {
	if mock.SetStateForAddressCalled != nil {
		return mock.SetStateForAddressCalled(address, state)
	}
	return nil
}

// RemoveAccount -
func (mock *NodeHandlerMock) RemoveAccount(address []byte) error {
	if mock.RemoveAccountCalled != nil {
		return mock.RemoveAccountCalled(address)
	}

	return nil
}

// GetBasePeers -
func (mock *NodeHandlerMock) GetBasePeers() map[uint32]core.PeerID {
	if mock.GetBasePeersCalled != nil {
		return mock.GetBasePeersCalled()
	}

	return nil
}

// SetBasePeers -
func (mock *NodeHandlerMock) SetBasePeers(basePeers map[uint32]core.PeerID) {
	if mock.SetBasePeersCalled != nil {
		mock.SetBasePeersCalled(basePeers)
	}
}

// Close -
func (mock *NodeHandlerMock) Close() error {
	if mock.CloseCalled != nil {
		return mock.CloseCalled()
	}
	return nil
}

// IsInterfaceNil -
func (mock *NodeHandlerMock) IsInterfaceNil() bool {
	return mock == nil
}
