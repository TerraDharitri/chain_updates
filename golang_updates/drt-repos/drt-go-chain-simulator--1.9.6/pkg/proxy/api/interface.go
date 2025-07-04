package api

import (
	dtosc "github.com/TerraDharitri/drt-go-chain-simulator/pkg/dtos"
	"github.com/TerraDharitri/drt-go-chain/node/chainSimulator/dtos"
)

// SimulatorFacade defines what a simulator facade should be able to do
type SimulatorFacade interface {
	GenerateBlocks(numOfBlocks int) error
	GetInitialWalletKeys() *dtos.InitialWalletKeys
	SetKeyValueForAddress(address string, keyValueMap map[string]string) error
	SetStateMultiple(stateSlice []*dtos.AddressState, noGenerate bool) error
	SetStateMultipleOverwrite(stateSlice []*dtos.AddressState, noGenerate bool) error
	AddValidatorKeys(validators *dtosc.ValidatorKeys) error
	GenerateBlocksUntilEpochIsReached(targetEpoch int32) error
	GenerateBlocksUntilTransactionIsProcessed(txHash string, maxNumOfBlocksToGenerate int) error
	ForceUpdateValidatorStatistics() error
	GetObserversInfo() (map[uint32]*dtosc.ObserverInfo, error)
	ForceChangeOfEpoch(targetEpoch uint32) error
	IsInterfaceNil() bool
}
