package testscommon

import "github.com/TerraDharitri/drt-go-chain-core/data"

// FieldsSizeCheckerMock -
type FieldsSizeCheckerMock struct {
	IsProofSizeValidCalled func(proof data.HeaderProofHandler) bool
}

// IsProofSizeValid -
func (mock *FieldsSizeCheckerMock) IsProofSizeValid(proof data.HeaderProofHandler) bool {
	if mock.IsProofSizeValidCalled != nil {
		return mock.IsProofSizeValidCalled(proof)
	}

	return true
}

// IsInterfaceNil -
func (mock *FieldsSizeCheckerMock) IsInterfaceNil() bool {
	return mock == nil
}
