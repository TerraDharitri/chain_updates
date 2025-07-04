package processMocks

import (
	"github.com/TerraDharitri/drt-go-chain-core/data"

	"github.com/TerraDharitri/drt-go-chain/process"
)

// ForkDetectorStub -
type ForkDetectorStub struct {
	AddHeaderCalled                 func(header data.HeaderHandler, hash []byte, state process.BlockHeaderState, selfNotarizedHeaders []data.HeaderHandler, selfNotarizedHeadersHashes [][]byte) error
	RemoveHeaderCalled              func(nonce uint64, hash []byte)
	CheckForkCalled                 func() *process.ForkInfo
	GetHighestFinalBlockNonceCalled func() uint64
	GetHighestFinalBlockHashCalled  func() []byte
	ProbableHighestNonceCalled      func() uint64
	ResetForkCalled                 func()
	GetNotarizedHeaderHashCalled    func(nonce uint64) []byte
	SetRollBackNonceCalled          func(nonce uint64)
	RestoreToGenesisCalled          func()
	ResetProbableHighestNonceCalled func()
	SetFinalToLastCheckpointCalled  func()
	ReceivedProofCalled             func(proof data.HeaderProofHandler)
	AddCheckpointCalled             func(nonce uint64, round uint64, hash []byte)
}

// RestoreToGenesis -
func (fdm *ForkDetectorStub) RestoreToGenesis() {
	fdm.RestoreToGenesisCalled()
}

// AddHeader -
func (fdm *ForkDetectorStub) AddHeader(header data.HeaderHandler, hash []byte, state process.BlockHeaderState, selfNotarizedHeaders []data.HeaderHandler, selfNotarizedHeadersHashes [][]byte) error {
	if fdm.AddHeaderCalled != nil {
		return fdm.AddHeaderCalled(header, hash, state, selfNotarizedHeaders, selfNotarizedHeadersHashes)
	}
	return nil
}

// RemoveHeader -
func (fdm *ForkDetectorStub) RemoveHeader(nonce uint64, hash []byte) {
	fdm.RemoveHeaderCalled(nonce, hash)
}

// CheckFork -
func (fdm *ForkDetectorStub) CheckFork() *process.ForkInfo {
	return fdm.CheckForkCalled()
}

// GetHighestFinalBlockNonce -
func (fdm *ForkDetectorStub) GetHighestFinalBlockNonce() uint64 {
	if fdm.GetHighestFinalBlockNonceCalled != nil {
		return fdm.GetHighestFinalBlockNonceCalled()
	}
	return 0
}

// GetHighestFinalBlockHash -
func (fdm *ForkDetectorStub) GetHighestFinalBlockHash() []byte {
	return fdm.GetHighestFinalBlockHashCalled()
}

// ProbableHighestNonce -
func (fdm *ForkDetectorStub) ProbableHighestNonce() uint64 {
	return fdm.ProbableHighestNonceCalled()
}

// SetRollBackNonce -
func (fdm *ForkDetectorStub) SetRollBackNonce(nonce uint64) {
	if fdm.SetRollBackNonceCalled != nil {
		fdm.SetRollBackNonceCalled(nonce)
	}
}

// ResetFork -
func (fdm *ForkDetectorStub) ResetFork() {
	fdm.ResetForkCalled()
}

// GetNotarizedHeaderHash -
func (fdm *ForkDetectorStub) GetNotarizedHeaderHash(nonce uint64) []byte {
	return fdm.GetNotarizedHeaderHashCalled(nonce)
}

// ResetProbableHighestNonce -
func (fdm *ForkDetectorStub) ResetProbableHighestNonce() {
	if fdm.ResetProbableHighestNonceCalled != nil {
		fdm.ResetProbableHighestNonceCalled()
	}
}

// SetFinalToLastCheckpoint -
func (fdm *ForkDetectorStub) SetFinalToLastCheckpoint() {
	if fdm.SetFinalToLastCheckpointCalled != nil {
		fdm.SetFinalToLastCheckpointCalled()
	}
}

// ReceivedProof -
func (fdm *ForkDetectorStub) ReceivedProof(proof data.HeaderProofHandler) {
	if fdm.ReceivedProofCalled != nil {
		fdm.ReceivedProofCalled(proof)
	}
}

// AddCheckpoint -
func (fdm *ForkDetectorStub) AddCheckpoint(nonce uint64, round uint64, hash []byte) {
	if fdm.AddCheckpointCalled != nil {
		fdm.AddCheckpointCalled(nonce, round, hash)
	}
}

// IsInterfaceNil returns true if there is no value under the interface
func (fdm *ForkDetectorStub) IsInterfaceNil() bool {
	return fdm == nil
}
