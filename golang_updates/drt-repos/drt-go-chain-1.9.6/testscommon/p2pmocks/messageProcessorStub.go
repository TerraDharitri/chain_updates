package p2pmocks

import (
	"github.com/TerraDharitri/drt-go-chain-communication/p2p"
	"github.com/TerraDharitri/drt-go-chain-core/core"
)

// MessageProcessorStub -
type MessageProcessorStub struct {
	ProcessReceivedMessageCalled func(message p2p.MessageP2P, fromConnectedPeer core.PeerID, source p2p.MessageHandler) ([]byte, error)
}

// ProcessReceivedMessage -
func (stub *MessageProcessorStub) ProcessReceivedMessage(message p2p.MessageP2P, fromConnectedPeer core.PeerID, source p2p.MessageHandler) ([]byte, error) {
	if stub.ProcessReceivedMessageCalled != nil {
		return stub.ProcessReceivedMessageCalled(message, fromConnectedPeer, source)
	}

	return nil, nil
}

// IsInterfaceNil -
func (stub *MessageProcessorStub) IsInterfaceNil() bool {
	return stub == nil
}
