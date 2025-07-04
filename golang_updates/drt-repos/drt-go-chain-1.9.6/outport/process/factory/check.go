package factory

import (
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain/outport/process"
	"github.com/TerraDharitri/drt-go-chain/outport/process/alteredaccounts"
	"github.com/TerraDharitri/drt-go-chain/outport/process/transactionsfee"
)

func checkArgOutportDataProviderFactory(arg ArgOutportDataProviderFactory) error {
	if check.IfNil(arg.AddressConverter) {
		return alteredaccounts.ErrNilPubKeyConverter
	}
	if check.IfNil(arg.AccountsDB) {
		return alteredaccounts.ErrNilAccountsDB
	}
	if check.IfNil(arg.Marshaller) {
		return transactionsfee.ErrNilMarshaller
	}
	if check.IfNil(arg.DcdtDataStorageHandler) {
		return alteredaccounts.ErrNilDCDTDataStorageHandler
	}
	if check.IfNil(arg.TransactionsStorer) {
		return transactionsfee.ErrNilStorage
	}
	if check.IfNil(arg.EconomicsData) {
		return transactionsfee.ErrNilTransactionFeeCalculator
	}
	if check.IfNil(arg.ShardCoordinator) {
		return transactionsfee.ErrNilShardCoordinator
	}
	if check.IfNil(arg.TxCoordinator) {
		return process.ErrNilTransactionCoordinator
	}
	if check.IfNil(arg.NodesCoordinator) {
		return process.ErrNilNodesCoordinator
	}
	if check.IfNil(arg.GasConsumedProvider) {
		return process.ErrNilGasConsumedProvider
	}
	if check.IfNil(arg.Hasher) {
		return process.ErrNilHasher
	}
	if check.IfNil(arg.MbsStorer) {
		return process.ErrNilStorer
	}
	if check.IfNil(arg.EnableEpochsHandler) {
		return process.ErrNilEnableEpochsHandler
	}
	if check.IfNil(arg.ExecutionOrderGetter) {
		return process.ErrNilExecutionOrderGetter
	}
	if check.IfNil(arg.ProofsPool) {
		return process.ErrNilProofsPool
	}

	return nil
}
