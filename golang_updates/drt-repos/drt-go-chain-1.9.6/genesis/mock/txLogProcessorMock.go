package mock

import (
	"fmt"

	"github.com/TerraDharitri/drt-go-chain-core/data"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
)

// TxLogProcessorMock -
type TxLogProcessorMock struct {
}

// GetLog -
func (tlpm *TxLogProcessorMock) GetLog(_ []byte) (data.LogHandler, error) {
	return nil, fmt.Errorf("log not found for provided tx hash")
}

// SaveLog -
func (tlpm *TxLogProcessorMock) SaveLog(_ []byte, _ data.TransactionHandler, _ []*vmcommon.LogEntry) error {
	return nil
}

// Clean -
func (tlpm *TxLogProcessorMock) Clean() {
}

// IsInterfaceNil -
func (tlpm *TxLogProcessorMock) IsInterfaceNil() bool {
	return tlpm == nil
}

// GetAllCurrentLogs -
func (tlpm *TxLogProcessorMock) GetAllCurrentLogs() []*data.LogData {
	return nil
}
