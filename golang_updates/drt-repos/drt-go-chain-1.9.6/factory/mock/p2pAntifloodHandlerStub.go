package mock

import (
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain/p2p"
	"github.com/TerraDharitri/drt-go-chain/process"
)

// P2PAntifloodHandlerStub -
type P2PAntifloodHandlerStub struct {
	CanProcessMessageCalled            func(message p2p.MessageP2P, fromConnectedPeer core.PeerID) error
	CanProcessMessagesOnTopicCalled    func(peer core.PeerID, topic string, numMessages uint32, totalSize uint64, sequence []byte) error
	ApplyConsensusSizeCalled           func(size int)
	SetDebuggerCalled                  func(debugger process.AntifloodDebugger) error
	BlacklistPeerCalled                func(peer core.PeerID, reason string, duration time.Duration)
	IsOriginatorEligibleForTopicCalled func(pid core.PeerID, topic string) error
	SetConsensusSizeNotifierCalled     func(subscriber process.ChainParametersSubscriber, shardID uint32)
}

// CanProcessMessage -
func (p2pahs *P2PAntifloodHandlerStub) CanProcessMessage(message p2p.MessageP2P, fromConnectedPeer core.PeerID) error {
	if p2pahs.CanProcessMessageCalled == nil {
		return nil
	}
	return p2pahs.CanProcessMessageCalled(message, fromConnectedPeer)
}

// IsOriginatorEligibleForTopic -
func (p2pahs *P2PAntifloodHandlerStub) IsOriginatorEligibleForTopic(pid core.PeerID, topic string) error {
	if p2pahs.IsOriginatorEligibleForTopicCalled != nil {
		return p2pahs.IsOriginatorEligibleForTopicCalled(pid, topic)
	}
	return nil
}

// CanProcessMessagesOnTopic -
func (p2pahs *P2PAntifloodHandlerStub) CanProcessMessagesOnTopic(peer core.PeerID, topic string, numMessages uint32, totalSize uint64, sequence []byte) error {
	if p2pahs.CanProcessMessagesOnTopicCalled == nil {
		return nil
	}
	return p2pahs.CanProcessMessagesOnTopicCalled(peer, topic, numMessages, totalSize, sequence)
}

// SetConsensusSizeNotifier -
func (p2pahs *P2PAntifloodHandlerStub) SetConsensusSizeNotifier(subscriber process.ChainParametersSubscriber, shardID uint32) {
	if p2pahs.SetConsensusSizeNotifierCalled != nil {
		p2pahs.SetConsensusSizeNotifierCalled(subscriber, shardID)
	}
}

// SetDebugger -
func (p2pahs *P2PAntifloodHandlerStub) SetDebugger(debugger process.AntifloodDebugger) error {
	if p2pahs.SetDebuggerCalled != nil {
		return p2pahs.SetDebuggerCalled(debugger)
	}

	return nil
}

// BlacklistPeer -
func (p2pahs *P2PAntifloodHandlerStub) BlacklistPeer(peer core.PeerID, reason string, duration time.Duration) {
	if p2pahs.BlacklistPeerCalled != nil {
		p2pahs.BlacklistPeerCalled(peer, reason, duration)
	}
}

// ResetForTopic -
func (p2pahs *P2PAntifloodHandlerStub) ResetForTopic(_ string) {

}

// SetMaxMessagesForTopic -
func (p2pahs *P2PAntifloodHandlerStub) SetMaxMessagesForTopic(_ string, _ uint32) {

}

// SetPeerValidatorMapper -
func (p2pahs *P2PAntifloodHandlerStub) SetPeerValidatorMapper(_ process.PeerValidatorMapper) error {
	return nil
}

// SetTopicsForAll -
func (p2pahs *P2PAntifloodHandlerStub) SetTopicsForAll(_ ...string) {

}

// Close -
func (p2pahs *P2PAntifloodHandlerStub) Close() error {
	return nil
}

// IsInterfaceNil -
func (p2pahs *P2PAntifloodHandlerStub) IsInterfaceNil() bool {
	return p2pahs == nil
}
