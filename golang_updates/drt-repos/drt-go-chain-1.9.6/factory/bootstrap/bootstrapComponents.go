package bootstrap

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	logger "github.com/TerraDharitri/drt-go-chain-logger"
	interceptorFactory "github.com/TerraDharitri/drt-go-chain/process/interceptors/factory"

	nodeFactory "github.com/TerraDharitri/drt-go-chain/cmd/node/factory"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/epochStart/bootstrap"
	"github.com/TerraDharitri/drt-go-chain/errors"
	"github.com/TerraDharitri/drt-go-chain/factory"
	"github.com/TerraDharitri/drt-go-chain/factory/block"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/guardian"
	"github.com/TerraDharitri/drt-go-chain/process/headerCheck"
	"github.com/TerraDharitri/drt-go-chain/process/smartContract"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
	"github.com/TerraDharitri/drt-go-chain/storage"
	"github.com/TerraDharitri/drt-go-chain/storage/directoryhandler"
	storageFactory "github.com/TerraDharitri/drt-go-chain/storage/factory"
	"github.com/TerraDharitri/drt-go-chain/storage/latestData"
	"github.com/TerraDharitri/drt-go-chain/storage/storageunit"
)

var log = logger.GetOrCreate("factory")

// BootstrapComponentsFactoryArgs holds the arguments needed to create a bootstrap components factory
type BootstrapComponentsFactoryArgs struct {
	Config                         config.Config
	RoundConfig                    config.RoundConfig
	PrefConfig                     config.Preferences
	ImportDbConfig                 config.ImportDbConfig
	FlagsConfig                    config.ContextFlagsConfig
	WorkingDir                     string
	CoreComponents                 factory.CoreComponentsHolder
	CryptoComponents               factory.CryptoComponentsHolder
	NetworkComponents              factory.NetworkComponentsHolder
	StatusCoreComponents           factory.StatusCoreComponentsHolder
	InterceptedDataVerifierFactory process.InterceptedDataVerifierFactory
}

type bootstrapComponentsFactory struct {
	config               config.Config
	prefConfig           config.Preferences
	importDbConfig       config.ImportDbConfig
	flagsConfig          config.ContextFlagsConfig
	workingDir           string
	coreComponents       factory.CoreComponentsHolder
	cryptoComponents     factory.CryptoComponentsHolder
	networkComponents    factory.NetworkComponentsHolder
	statusCoreComponents factory.StatusCoreComponentsHolder
}

type bootstrapComponents struct {
	epochStartBootstrapper          factory.EpochStartBootstrapper
	bootstrapParamsHolder           factory.BootstrapParamsHolder
	nodeType                        core.NodeType
	shardCoordinator                sharding.Coordinator
	headerVersionHandler            nodeFactory.HeaderVersionHandler
	versionedHeaderFactory          nodeFactory.VersionedHeaderFactory
	headerIntegrityVerifier         nodeFactory.HeaderIntegrityVerifierHandler
	guardedAccountHandler           process.GuardedAccountHandler
	nodesCoordinatorRegistryFactory nodesCoordinator.NodesCoordinatorRegistryFactory
}

// NewBootstrapComponentsFactory creates an instance of bootstrapComponentsFactory
func NewBootstrapComponentsFactory(args BootstrapComponentsFactoryArgs) (*bootstrapComponentsFactory, error) {
	if check.IfNil(args.CoreComponents) {
		return nil, errors.ErrNilCoreComponentsHolder
	}
	if check.IfNil(args.CoreComponents.EnableEpochsHandler()) {
		return nil, errors.ErrNilEnableEpochsHandler
	}
	if check.IfNil(args.CryptoComponents) {
		return nil, errors.ErrNilCryptoComponentsHolder
	}
	if check.IfNil(args.NetworkComponents) {
		return nil, errors.ErrNilNetworkComponentsHolder
	}
	if check.IfNil(args.StatusCoreComponents) {
		return nil, errors.ErrNilStatusCoreComponents
	}
	if check.IfNil(args.StatusCoreComponents.TrieSyncStatistics()) {
		return nil, errors.ErrNilTrieSyncStatistics
	}
	if args.WorkingDir == "" {
		return nil, errors.ErrInvalidWorkingDir
	}
	if check.IfNil(args.StatusCoreComponents.AppStatusHandler()) {
		return nil, errors.ErrNilAppStatusHandler
	}

	return &bootstrapComponentsFactory{
		config:               args.Config,
		prefConfig:           args.PrefConfig,
		importDbConfig:       args.ImportDbConfig,
		flagsConfig:          args.FlagsConfig,
		workingDir:           args.WorkingDir,
		coreComponents:       args.CoreComponents,
		cryptoComponents:     args.CryptoComponents,
		networkComponents:    args.NetworkComponents,
		statusCoreComponents: args.StatusCoreComponents,
	}, nil
}

// Create creates the bootstrap components
func (bcf *bootstrapComponentsFactory) Create() (*bootstrapComponents, error) {
	destShardIdAsObserver, err := common.ProcessDestinationShardAsObserver(bcf.prefConfig.Preferences.DestinationShardAsObserver)
	if err != nil {
		return nil, err
	}

	versionsCache, err := storageunit.NewCache(storageFactory.GetCacherFromConfig(bcf.config.Versions.Cache))
	if err != nil {
		return nil, err
	}

	headerVersionHandler, err := block.NewHeaderVersionHandler(
		bcf.config.Versions.VersionsByEpochs,
		bcf.config.Versions.DefaultVersion,
		versionsCache,
	)
	if err != nil {
		return nil, err
	}

	headerIntegrityVerifier, err := headerCheck.NewHeaderIntegrityVerifier(
		[]byte(bcf.coreComponents.ChainID()),
		headerVersionHandler,
	)
	if err != nil {
		return nil, err
	}

	genesisShardCoordinator, nodeType, err := CreateShardCoordinator(
		bcf.coreComponents.GenesisNodesSetup(),
		bcf.cryptoComponents.PublicKey(),
		bcf.prefConfig.Preferences,
		log,
	)
	if err != nil {
		return nil, err
	}

	bootstrapDataProvider, err := storageFactory.NewBootstrapDataProvider(bcf.coreComponents.InternalMarshalizer())
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrNewBootstrapDataProvider, err)
	}

	parentDir := filepath.Join(
		bcf.workingDir,
		common.DefaultDBPath,
		bcf.coreComponents.ChainID())

	latestStorageDataProvider, err := createLatestStorageDataProvider(
		bootstrapDataProvider,
		bcf.config,
		parentDir,
		storage.DefaultEpochString,
		storage.DefaultShardString,
	)
	if err != nil {
		return nil, err
	}

	unitOpener, err := createUnitOpener(
		bootstrapDataProvider,
		latestStorageDataProvider,
		storage.DefaultEpochString,
		storage.DefaultShardString,
	)
	if err != nil {
		return nil, err
	}

	dataSyncerFactory := bootstrap.NewScheduledDataSyncerFactory()

	// increment num received to make sure that first heartbeat message
	// will have value 1, thus explorer will display status in progress
	tss := bcf.statusCoreComponents.TrieSyncStatistics()
	tss.AddNumProcessed(1)

	setGuardianEpochsDelay := bcf.config.GeneralSettings.SetGuardianEpochsDelay
	guardedAccountHandler, err := guardian.NewGuardedAccount(bcf.coreComponents.InternalMarshalizer(), bcf.coreComponents.EpochNotifier(), setGuardianEpochsDelay)
	if err != nil {
		return nil, err
	}

	nodesCoordinatorRegistryFactory, err := nodesCoordinator.NewNodesCoordinatorRegistryFactory(
		bcf.coreComponents.InternalMarshalizer(),
		bcf.coreComponents.EnableEpochsHandler().GetActivationEpoch(common.StakingV4Step2Flag),
	)
	if err != nil {
		return nil, err
	}

	// create a new instance of interceptedDataVerifier which will be used for bootstrap only
	interceptedDataVerifierFactory := interceptorFactory.NewInterceptedDataVerifierFactory(interceptorFactory.InterceptedDataVerifierFactoryArgs{
		CacheSpan:   time.Duration(bcf.config.InterceptedDataVerifier.CacheSpanInSec) * time.Second,
		CacheExpiry: time.Duration(bcf.config.InterceptedDataVerifier.CacheExpiryInSec) * time.Second,
	})

	epochStartBootstrapArgs := bootstrap.ArgsEpochStartBootstrap{
		CoreComponentsHolder:            bcf.coreComponents,
		CryptoComponentsHolder:          bcf.cryptoComponents,
		MainMessenger:                   bcf.networkComponents.NetworkMessenger(),
		FullArchiveMessenger:            bcf.networkComponents.FullArchiveNetworkMessenger(),
		GeneralConfig:                   bcf.config,
		PrefsConfig:                     bcf.prefConfig.Preferences,
		FlagsConfig:                     bcf.flagsConfig,
		EconomicsData:                   bcf.coreComponents.EconomicsData(),
		GenesisNodesConfig:              bcf.coreComponents.GenesisNodesSetup(),
		GenesisShardCoordinator:         genesisShardCoordinator,
		StorageUnitOpener:               unitOpener,
		Rater:                           bcf.coreComponents.Rater(),
		DestinationShardAsObserver:      destShardIdAsObserver,
		NodeShuffler:                    bcf.coreComponents.NodesShuffler(),
		RoundHandler:                    bcf.coreComponents.RoundHandler(),
		LatestStorageDataProvider:       latestStorageDataProvider,
		ArgumentsParser:                 smartContract.NewArgumentParser(),
		StatusHandler:                   bcf.statusCoreComponents.AppStatusHandler(),
		HeaderIntegrityVerifier:         headerIntegrityVerifier,
		DataSyncerCreator:               dataSyncerFactory,
		ScheduledSCRsStorer:             nil, // will be updated after sync from network
		TrieSyncStatisticsProvider:      tss,
		NodeProcessingMode:              common.GetNodeProcessingMode(&bcf.importDbConfig),
		StateStatsHandler:               bcf.statusCoreComponents.StateStatsHandler(),
		NodesCoordinatorRegistryFactory: nodesCoordinatorRegistryFactory,
		EnableEpochsHandler:             bcf.coreComponents.EnableEpochsHandler(),
		InterceptedDataVerifierFactory:  interceptedDataVerifierFactory,
	}

	var epochStartBootstrapper factory.EpochStartBootstrapper
	if bcf.importDbConfig.IsImportDBMode {
		storageArg := bootstrap.ArgsStorageEpochStartBootstrap{
			ArgsEpochStartBootstrap:    epochStartBootstrapArgs,
			ImportDbConfig:             bcf.importDbConfig,
			ChanGracefullyClose:        bcf.coreComponents.ChanStopNodeProcess(),
			TimeToWaitForRequestedData: bootstrap.DefaultTimeToWaitForRequestedData,
		}

		epochStartBootstrapper, err = bootstrap.NewStorageEpochStartBootstrap(storageArg)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", errors.ErrNewStorageEpochStartBootstrap, err)
		}
	} else {
		epochStartBootstrapper, err = bootstrap.NewEpochStartBootstrap(epochStartBootstrapArgs)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", errors.ErrNewEpochStartBootstrap, err)
		}
	}

	bootstrapParameters, err := epochStartBootstrapper.Bootstrap()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrBootstrap, err)
	}

	log.Info("bootstrap parameters",
		"shardId", bootstrapParameters.SelfShardId,
		"epoch", bootstrapParameters.Epoch,
		"numShards", bootstrapParameters.NumOfShards,
	)

	shardCoordinator, err := sharding.NewMultiShardCoordinator(
		bootstrapParameters.NumOfShards,
		bootstrapParameters.SelfShardId)
	if err != nil {
		return nil, err
	}

	versionedHeaderFactory, err := bcf.createHeaderFactory(headerVersionHandler, bootstrapParameters.SelfShardId)
	if err != nil {
		return nil, err
	}

	return &bootstrapComponents{
		epochStartBootstrapper: epochStartBootstrapper,
		bootstrapParamsHolder: &bootstrapParams{
			bootstrapParams: bootstrapParameters,
		},
		nodeType:                        nodeType,
		shardCoordinator:                shardCoordinator,
		headerVersionHandler:            headerVersionHandler,
		headerIntegrityVerifier:         headerIntegrityVerifier,
		versionedHeaderFactory:          versionedHeaderFactory,
		guardedAccountHandler:           guardedAccountHandler,
		nodesCoordinatorRegistryFactory: nodesCoordinatorRegistryFactory,
	}, nil
}

func (bcf *bootstrapComponentsFactory) createHeaderFactory(handler nodeFactory.HeaderVersionHandler, shardID uint32) (nodeFactory.VersionedHeaderFactory, error) {
	if shardID == core.MetachainShardId {
		return block.NewMetaHeaderFactory(handler)
	}
	return block.NewShardHeaderFactory(handler)
}

// Close closes the bootstrap components, closing at the same time any running goroutines
func (bc *bootstrapComponents) Close() error {
	// TODO: close all components
	if !check.IfNil(bc.epochStartBootstrapper) {
		return bc.epochStartBootstrapper.Close()
	}

	return nil
}

// NodeType returns the node type
func (bc *bootstrapComponents) NodeType() core.NodeType {
	return bc.nodeType
}

// ShardCoordinator returns the shard coordinator
func (bc *bootstrapComponents) ShardCoordinator() sharding.Coordinator {
	return bc.shardCoordinator
}

// HeaderVersionHandler returns the header version handler
func (bc *bootstrapComponents) HeaderVersionHandler() nodeFactory.HeaderVersionHandler {
	return bc.headerVersionHandler
}

// VersionedHeaderFactory returns the versioned header factory
func (bc *bootstrapComponents) VersionedHeaderFactory() nodeFactory.VersionedHeaderFactory {
	return bc.versionedHeaderFactory
}

// HeaderIntegrityVerifier returns the header integrity verifier
func (bc *bootstrapComponents) HeaderIntegrityVerifier() nodeFactory.HeaderIntegrityVerifierHandler {
	return bc.headerIntegrityVerifier
}

// createLatestStorageDataProvider will create the latest storage data provider handler
func createLatestStorageDataProvider(
	bootstrapDataProvider storageFactory.BootstrapDataProviderHandler,
	generalConfig config.Config,
	parentDir string,
	defaultEpochString string,
	defaultShardString string,
) (storage.LatestStorageDataProviderHandler, error) {
	directoryReader := directoryhandler.NewDirectoryReader()

	latestStorageDataArgs := latestData.ArgsLatestDataProvider{
		GeneralConfig:         generalConfig,
		BootstrapDataProvider: bootstrapDataProvider,
		DirectoryReader:       directoryReader,
		ParentDir:             parentDir,
		DefaultEpochString:    defaultEpochString,
		DefaultShardString:    defaultShardString,
	}

	return latestData.NewLatestDataProvider(latestStorageDataArgs)
}

// createUnitOpener will create a new unit opener handler
func createUnitOpener(
	bootstrapDataProvider storageFactory.BootstrapDataProviderHandler,
	latestDataFromStorageProvider storage.LatestStorageDataProviderHandler,
	defaultEpochString string,
	defaultShardString string,
) (storage.UnitOpenerHandler, error) {
	argsStorageUnitOpener := storageFactory.ArgsNewOpenStorageUnits{
		BootstrapDataProvider:     bootstrapDataProvider,
		LatestStorageDataProvider: latestDataFromStorageProvider,
		DefaultEpochString:        defaultEpochString,
		DefaultShardString:        defaultShardString,
	}

	return storageFactory.NewStorageUnitOpenHandler(argsStorageUnitOpener)
}
