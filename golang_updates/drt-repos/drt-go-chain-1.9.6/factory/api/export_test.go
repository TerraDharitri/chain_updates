package api

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/factory"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/vm"
)

// SCQueryElementArgs -
type SCQueryElementArgs struct {
	GeneralConfig         *config.Config
	EpochConfig           *config.EpochConfig
	CoreComponents        factory.CoreComponentsHolder
	StateComponents       factory.StateComponentsHolder
	StatusCoreComponents  factory.StatusCoreComponentsHolder
	DataComponents        factory.DataComponentsHolder
	ProcessComponents     factory.ProcessComponentsHolder
	GasScheduleNotifier   core.GasScheduleNotifier
	MessageSigVerifier    vm.MessageSignVerifier
	SystemSCConfig        *config.SystemSmartContractsConfig
	Bootstrapper          process.Bootstrapper
	AllowVMQueriesChan    chan struct{}
	WorkingDir            string
	Index                 int
	GuardedAccountHandler process.GuardedAccountHandler
}

// CreateScQueryElement -
func CreateScQueryElement(args SCQueryElementArgs) (process.SCQueryService, common.StorageManager, error) {
	return createScQueryElement(scQueryElementArgs{
		generalConfig:         args.GeneralConfig,
		epochConfig:           args.EpochConfig,
		coreComponents:        args.CoreComponents,
		stateComponents:       args.StateComponents,
		dataComponents:        args.DataComponents,
		processComponents:     args.ProcessComponents,
		statusCoreComponents:  args.StatusCoreComponents,
		gasScheduleNotifier:   args.GasScheduleNotifier,
		messageSigVerifier:    args.MessageSigVerifier,
		systemSCConfig:        args.SystemSCConfig,
		bootstrapper:          args.Bootstrapper,
		allowVMQueriesChan:    args.AllowVMQueriesChan,
		workingDir:            args.WorkingDir,
		index:                 args.Index,
		guardedAccountHandler: args.GuardedAccountHandler,
	})
}

// CreateBlockchainForScQuery -
func CreateBlockchainForScQuery(selfShardID uint32) (data.ChainHandler, error) {
	return createBlockchainForScQuery(selfShardID)
}
