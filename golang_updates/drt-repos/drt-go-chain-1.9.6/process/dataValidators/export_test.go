package dataValidators

import (
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	"github.com/TerraDharitri/drt-go-chain/process"
)

// CheckAccount -
func (txv *txValidator) CheckAccount(
	interceptedTx process.InterceptedTransactionHandler,
	accountHandler vmcommon.AccountHandler,
) error {
	return txv.checkAccount(interceptedTx, accountHandler)
}

// GetTxData -
func GetTxData(interceptedTx process.InterceptedTransactionHandler) ([]byte, error) {
	return getTxData(interceptedTx)
}
