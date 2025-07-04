package mock

import "github.com/TerraDharitri/drt-go-chain-core/core"

// FloodPreventerStub -
type FloodPreventerStub struct {
	IncreaseLoadCalled       func(pid core.PeerID, size uint64) error
	ApplyConsensusSizeCalled func(size int)
	ResetCalled              func()
}

// IncreaseLoad -
func (fps *FloodPreventerStub) IncreaseLoad(pid core.PeerID, size uint64) error {
	if fps.IncreaseLoadCalled != nil {
		return fps.IncreaseLoadCalled(pid, size)
	}

	return nil
}

// ApplyConsensusSize -
func (fps *FloodPreventerStub) ApplyConsensusSize(size int) {
	if fps.ApplyConsensusSizeCalled != nil {
		fps.ApplyConsensusSizeCalled(size)
	}
}

// Reset -
func (fps *FloodPreventerStub) Reset() {
	fps.ResetCalled()
}

// IsInterfaceNil -
func (fps *FloodPreventerStub) IsInterfaceNil() bool {
	return fps == nil
}
