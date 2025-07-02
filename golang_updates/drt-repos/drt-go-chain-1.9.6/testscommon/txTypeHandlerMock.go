package testscommon

import (
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain/process"
)

// TxTypeHandlerMock -
type TxTypeHandlerMock struct {
	ComputeTransactionTypeCalled func(tx data.TransactionHandler) (process.TransactionType, process.TransactionType, bool)
}

// ComputeTransactionType -
func (th *TxTypeHandlerMock) ComputeTransactionType(tx data.TransactionHandler) (process.TransactionType, process.TransactionType, bool) {
	if th.ComputeTransactionTypeCalled == nil {
		return process.MoveBalance, process.MoveBalance, false
	}

	return th.ComputeTransactionTypeCalled(tx)
}

// IsInterfaceNil -
func (th *TxTypeHandlerMock) IsInterfaceNil() bool {
	return th == nil
}
