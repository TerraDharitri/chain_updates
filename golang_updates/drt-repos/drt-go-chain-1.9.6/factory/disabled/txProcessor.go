package disabled

import (
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	"github.com/TerraDharitri/drt-go-chain/state"
)

// TxProcessor implements the TransactionProcessor interface but does nothing as it is disabled
type TxProcessor struct {
}

// ProcessTransaction does nothing as it is disabled
func (txProc *TxProcessor) ProcessTransaction(_ *transaction.Transaction) (vmcommon.ReturnCode, error) {
	return 0, nil
}

// VerifyTransaction does nothing as it is disabled
func (txProc *TxProcessor) VerifyTransaction(_ *transaction.Transaction) error {
	return nil
}

// VerifyGuardian does nothing as it is disabled
func (txProc *TxProcessor) VerifyGuardian(_ *transaction.Transaction, _ state.UserAccountHandler) error {
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (txProc *TxProcessor) IsInterfaceNil() bool {
	return txProc == nil
}
