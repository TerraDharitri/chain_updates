package statusCore

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	logger "github.com/TerraDharitri/drt-go-chain-logger"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/common/statistics"
	"github.com/TerraDharitri/drt-go-chain/common/statistics/disabled"
	"github.com/TerraDharitri/drt-go-chain/common/statistics/machine"
	"github.com/TerraDharitri/drt-go-chain/config"
	errorsDrt "github.com/TerraDharitri/drt-go-chain/errors"
	"github.com/TerraDharitri/drt-go-chain/factory"
	"github.com/TerraDharitri/drt-go-chain/node/external"
	"github.com/TerraDharitri/drt-go-chain/node/metrics"
	"github.com/TerraDharitri/drt-go-chain/statusHandler"
	"github.com/TerraDharitri/drt-go-chain/statusHandler/persister"
	trieStatistics "github.com/TerraDharitri/drt-go-chain/trie/statistics"
)

var log = logger.GetOrCreate("factory")

// StatusCoreComponentsFactoryArgs holds the arguments needed for creating a status core components factory
type StatusCoreComponentsFactoryArgs struct {
	Config          config.Config
	EpochConfig     config.EpochConfig
	RoundConfig     config.RoundConfig
	RatingsConfig   config.RatingsConfig
	EconomicsConfig config.EconomicsConfig
	CoreComp        factory.CoreComponentsHolder
}

// statusCoreComponentsFactory is responsible for creating the status core components
type statusCoreComponentsFactory struct {
	config          config.Config
	epochConfig     config.EpochConfig
	roundConfig     config.RoundConfig
	ratingsConfig   config.RatingsConfig
	economicsConfig config.EconomicsConfig
	coreComp        factory.CoreComponentsHolder
}

// statusCoreComponents is the DTO used for core components
type statusCoreComponents struct {
	resourceMonitor    factory.ResourceMonitor
	networkStatistics  factory.NetworkStatisticsProvider
	trieSyncStatistics factory.TrieSyncStatisticsProvider
	appStatusHandler   core.AppStatusHandler
	statusMetrics      external.StatusMetricsHandler
	persistentHandler  factory.PersistentStatusHandler
	stateStatsHandler  common.StateStatisticsHandler
}

// NewStatusCoreComponentsFactory initializes the factory which is responsible to creating status core components
func NewStatusCoreComponentsFactory(args StatusCoreComponentsFactoryArgs) (*statusCoreComponentsFactory, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	return &statusCoreComponentsFactory{
		config:          args.Config,
		epochConfig:     args.EpochConfig,
		roundConfig:     args.RoundConfig,
		ratingsConfig:   args.RatingsConfig,
		economicsConfig: args.EconomicsConfig,
		coreComp:        args.CoreComp,
	}, nil
}

func checkArgs(args StatusCoreComponentsFactoryArgs) error {
	if check.IfNil(args.CoreComp) {
		return errorsDrt.ErrNilCoreComponents
	}
	if check.IfNil(args.CoreComp.EconomicsData()) {
		return errorsDrt.ErrNilEconomicsData
	}

	return nil
}

// Create creates the status core components
func (sccf *statusCoreComponentsFactory) Create() (*statusCoreComponents, error) {
	netStats := machine.NewNetStatistics()

	resourceMonitor, err := statistics.NewResourceMonitor(sccf.config, netStats)
	if err != nil {
		return nil, err
	}

	if sccf.config.ResourceStats.Enabled {
		resourceMonitor.StartMonitoring()
	}

	appStatusHandler, statusMetrics, persistentStatusHandler, err := sccf.createStatusHandler()
	if err != nil {
		return nil, err
	}

	stateStatsHandler := sccf.createStateStatsHandler()

	ssc := &statusCoreComponents{
		resourceMonitor:    resourceMonitor,
		networkStatistics:  netStats,
		trieSyncStatistics: trieStatistics.NewTrieSyncStatistics(),
		appStatusHandler:   appStatusHandler,
		statusMetrics:      statusMetrics,
		persistentHandler:  persistentStatusHandler,
		stateStatsHandler:  stateStatsHandler,
	}

	return ssc, nil
}

func (sccf *statusCoreComponentsFactory) createStateStatsHandler() common.StateStatisticsHandler {
	if sccf.config.StateTriesConfig.StateStatisticsEnabled {
		return statistics.NewStateStatistics()
	}

	return disabled.NewStateStatistics()
}

func (sccf *statusCoreComponentsFactory) createStatusHandler() (core.AppStatusHandler, external.StatusMetricsHandler, factory.PersistentStatusHandler, error) {
	var appStatusHandlers []core.AppStatusHandler
	var handler core.AppStatusHandler
	statusMetrics := statusHandler.NewStatusMetrics()
	appStatusHandlers = append(appStatusHandlers, statusMetrics)

	persistentHandler, err := persister.NewPersistentStatusHandler(sccf.coreComp.InternalMarshalizer(), sccf.coreComp.Uint64ByteSliceConverter())
	if err != nil {
		return nil, nil, nil, err
	}
	appStatusHandlers = append(appStatusHandlers, persistentHandler)
	if len(appStatusHandlers) > 0 {
		handler, err = statusHandler.NewAppStatusFacadeWithHandlers(appStatusHandlers...)
		if err != nil {
			log.Warn("cannot init AppStatusFacade, will start with NilStatusHandler", "error", err)
			handler = statusHandler.NewNilStatusHandler()
		}
	} else {
		handler = statusHandler.NewNilStatusHandler()
		log.Debug("no AppStatusHandler used: will start with NilStatusHandler")
	}

	err = metrics.InitBaseMetrics(handler)
	if err != nil {
		return nil, nil, nil, err
	}

	err = metrics.InitConfigMetrics(handler, sccf.epochConfig, sccf.economicsConfig, sccf.coreComp.GenesisNodesSetup(), sccf.config.GatewayMetricsConfig)
	if err != nil {
		return nil, nil, nil, err
	}

	err = metrics.InitRatingsMetrics(handler, sccf.ratingsConfig)
	if err != nil {
		return nil, nil, nil, err
	}

	err = sccf.coreComp.EconomicsData().SetStatusHandler(handler)
	if err != nil {
		log.Debug("cannot set status handler to economicsData", "error", err)
		return nil, nil, nil, err
	}

	err = sccf.coreComp.RatingsData().SetStatusHandler(handler)
	if err != nil {
		log.Debug("cannot set status handler to ratingsData", "error", err)
		return nil, nil, nil, err
	}

	return handler, statusMetrics, persistentHandler, nil
}

// Close closes all underlying components
func (scc *statusCoreComponents) Close() error {
	var errNetStats error
	var errResourceMonitor error
	if !check.IfNil(scc.networkStatistics) {
		errNetStats = scc.networkStatistics.Close()
	}
	if !check.IfNil(scc.resourceMonitor) {
		errResourceMonitor = scc.resourceMonitor.Close()
	}
	if !check.IfNil(scc.appStatusHandler) {
		scc.appStatusHandler.Close()
	}

	if errNetStats != nil {
		return errNetStats
	}
	return errResourceMonitor
}
