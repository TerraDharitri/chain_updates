package process

import (
	"github.com/TerraDharitri/drt-go-chain/node/external"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/vm"
)

type genesisProcessors struct {
	txCoordinator       process.TransactionCoordinator
	systemSCs           vm.SystemSCContainer
	txProcessor         process.TransactionProcessor
	scProcessor         process.SmartContractProcessor
	scrProcessor        process.SmartContractResultProcessor
	rwdProcessor        process.RewardTransactionProcessor
	blockchainHook      process.BlockChainHookHandler
	queryService        external.SCQueryService
	vmContainersFactory process.VirtualMachinesContainerFactory
	vmContainer         process.VirtualMachinesContainer
}
