package data

import (
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
)

// SimulationResultsWithVMOutput is the data transfer object which will hold results for simulation a transaction's execution
type SimulationResultsWithVMOutput struct {
	transaction.SimulationResults
	VMOutput *vmcommon.VMOutput `json:"-"`
}
