package components

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain/api/gin"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/common/forking"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/facade"
	apiComp "github.com/TerraDharitri/drt-go-chain/factory/api"
	nodePack "github.com/TerraDharitri/drt-go-chain/node"
	"github.com/TerraDharitri/drt-go-chain/node/metrics"
	"github.com/TerraDharitri/drt-go-chain/process/mock"
)

func (node *testOnlyProcessingNode) createFacade(configs config.Configs, apiInterface APIConfigurator, vmQueryDelayAfterStartInMs uint64) error {
	log.Debug("creating api resolver structure")

	err := node.createMetrics(configs)
	if err != nil {
		return err
	}

	argsGasScheduleNotifier := forking.ArgsNewGasScheduleNotifier{
		GasScheduleConfig:  configs.EpochConfig.GasSchedule,
		ConfigDir:          configs.ConfigurationPathsHolder.GasScheduleDirectoryName,
		EpochNotifier:      node.CoreComponentsHolder.EpochNotifier(),
		WasmVMChangeLocker: node.CoreComponentsHolder.WasmVMChangeLocker(),
	}
	gasScheduleNotifier, err := forking.NewGasScheduleNotifier(argsGasScheduleNotifier)
	if err != nil {
		return err
	}

	allowVMQueriesChan := make(chan struct{})
	go func() {
		time.Sleep(time.Duration(vmQueryDelayAfterStartInMs) * time.Millisecond)
		close(allowVMQueriesChan)
		node.StatusCoreComponents.AppStatusHandler().SetStringValue(common.MetricAreVMQueriesReady, strconv.FormatBool(true))
	}()

	apiResolverArgs := &apiComp.ApiResolverArgs{
		Configs:              &configs,
		CoreComponents:       node.CoreComponentsHolder,
		DataComponents:       node.DataComponentsHolder,
		StateComponents:      node.StateComponentsHolder,
		BootstrapComponents:  node.BootstrapComponentsHolder,
		CryptoComponents:     node.CryptoComponentsHolder,
		ProcessComponents:    node.ProcessComponentsHolder,
		StatusCoreComponents: node.StatusCoreComponents,
		GasScheduleNotifier:  gasScheduleNotifier,
		Bootstrapper: &mock.BootstrapperStub{
			GetNodeStateCalled: func() common.NodeState {
				return common.NsSynchronized
			},
		},
		AllowVMQueriesChan: allowVMQueriesChan,
		StatusComponents:   node.StatusComponentsHolder,
		ProcessingMode:     common.GetNodeProcessingMode(configs.ImportDbConfig),
	}

	apiResolver, err := apiComp.CreateApiResolver(apiResolverArgs)
	if err != nil {
		return err
	}

	log.Debug("creating dharitri node facade")

	flagsConfig := configs.FlagsConfig

	nd, err := nodePack.NewNode(
		nodePack.WithStatusCoreComponents(node.StatusCoreComponents),
		nodePack.WithCoreComponents(node.CoreComponentsHolder),
		nodePack.WithCryptoComponents(node.CryptoComponentsHolder),
		nodePack.WithBootstrapComponents(node.BootstrapComponentsHolder),
		nodePack.WithStateComponents(node.StateComponentsHolder),
		nodePack.WithDataComponents(node.DataComponentsHolder),
		nodePack.WithStatusComponents(node.StatusComponentsHolder),
		nodePack.WithProcessComponents(node.ProcessComponentsHolder),
		nodePack.WithNetworkComponents(node.NetworkComponentsHolder),
		nodePack.WithInitialNodesPubKeys(node.CoreComponentsHolder.GenesisNodesSetup().InitialNodesPubKeys()),
		nodePack.WithRoundDuration(node.CoreComponentsHolder.GenesisNodesSetup().GetRoundDuration()),
		nodePack.WithGenesisTime(node.CoreComponentsHolder.GenesisTime()),
		nodePack.WithConsensusType(configs.GeneralConfig.Consensus.Type),
		nodePack.WithRequestedItemsHandler(node.ProcessComponentsHolder.RequestedItemsHandler()),
		nodePack.WithAddressSignatureSize(configs.GeneralConfig.AddressPubkeyConverter.SignatureLength),
		nodePack.WithValidatorSignatureSize(configs.GeneralConfig.ValidatorPubkeyConverter.SignatureLength),
		nodePack.WithPublicKeySize(configs.GeneralConfig.ValidatorPubkeyConverter.Length),
		nodePack.WithNodeStopChannel(node.CoreComponentsHolder.ChanStopNodeProcess()),
		nodePack.WithImportMode(configs.ImportDbConfig.IsImportDBMode),
		nodePack.WithDCDTNFTStorageHandler(node.ProcessComponentsHolder.DCDTDataStorageHandlerForAPI()),
	)
	if err != nil {
		return errors.New("error creating node: " + err.Error())
	}

	shardID := node.GetShardCoordinator().SelfId()
	restApiInterface := apiInterface.RestApiInterface(shardID)

	argNodeFacade := facade.ArgNodeFacade{
		Node:                   nd,
		ApiResolver:            apiResolver,
		RestAPIServerDebugMode: flagsConfig.EnableRestAPIServerDebugMode,
		WsAntifloodConfig:      configs.GeneralConfig.WebServerAntiflood,
		FacadeConfig: config.FacadeConfig{
			RestApiInterface: restApiInterface,
			PprofEnabled:     flagsConfig.EnablePprof,
		},
		ApiRoutesConfig: *configs.ApiRoutesConfig,
		AccountsState:   node.StateComponentsHolder.AccountsAdapter(),
		PeerState:       node.StateComponentsHolder.PeerAccounts(),
		Blockchain:      node.DataComponentsHolder.Blockchain(),
	}

	ef, err := facade.NewNodeFacade(argNodeFacade)
	if err != nil {
		return fmt.Errorf("%w while creating NodeFacade", err)
	}

	ef.SetSyncer(node.CoreComponentsHolder.SyncTimer())

	node.facadeHandler = ef

	return nil
}

func (node *testOnlyProcessingNode) createHttpServer(configs config.Configs) error {
	httpServerArgs := gin.ArgsNewWebServer{
		Facade:          node.facadeHandler,
		ApiConfig:       *configs.ApiRoutesConfig,
		AntiFloodConfig: configs.GeneralConfig.WebServerAntiflood,
	}

	httpServerWrapper, err := gin.NewGinWebServerHandler(httpServerArgs)
	if err != nil {
		return err
	}

	err = httpServerWrapper.StartHttpServer()
	if err != nil {
		return err
	}

	node.httpServer = httpServerWrapper

	return nil
}

func (node *testOnlyProcessingNode) createMetrics(configs config.Configs) error {
	err := metrics.InitMetrics(
		node.StatusCoreComponents.AppStatusHandler(),
		node.CryptoComponentsHolder.PublicKeyString(),
		node.BootstrapComponentsHolder.NodeType(),
		node.BootstrapComponentsHolder.ShardCoordinator(),
		node.CoreComponentsHolder.GenesisNodesSetup(),
		configs.FlagsConfig.Version,
		configs.EconomicsConfig,
		configs.GeneralConfig.EpochStartConfig.RoundsPerEpoch,
		node.CoreComponentsHolder.MinTransactionVersion(),
	)

	if err != nil {
		return err
	}

	metrics.SaveStringMetric(node.StatusCoreComponents.AppStatusHandler(), common.MetricNodeDisplayName, configs.PreferencesConfig.Preferences.NodeDisplayName)
	metrics.SaveStringMetric(node.StatusCoreComponents.AppStatusHandler(), common.MetricRedundancyLevel, fmt.Sprintf("%d", configs.PreferencesConfig.Preferences.RedundancyLevel))
	metrics.SaveStringMetric(node.StatusCoreComponents.AppStatusHandler(), common.MetricRedundancyIsMainActive, common.MetricValueNA)
	metrics.SaveStringMetric(node.StatusCoreComponents.AppStatusHandler(), common.MetricChainId, node.CoreComponentsHolder.ChainID())
	metrics.SaveUint64Metric(node.StatusCoreComponents.AppStatusHandler(), common.MetricGasPerDataByte, node.CoreComponentsHolder.EconomicsData().GasPerDataByte())
	metrics.SaveUint64Metric(node.StatusCoreComponents.AppStatusHandler(), common.MetricMinGasPrice, node.CoreComponentsHolder.EconomicsData().MinGasPrice())
	metrics.SaveUint64Metric(node.StatusCoreComponents.AppStatusHandler(), common.MetricMinGasLimit, node.CoreComponentsHolder.EconomicsData().MinGasLimit())
	metrics.SaveUint64Metric(node.StatusCoreComponents.AppStatusHandler(), common.MetricExtraGasLimitGuardedTx, node.CoreComponentsHolder.EconomicsData().ExtraGasLimitGuardedTx())
	metrics.SaveUint64Metric(node.StatusCoreComponents.AppStatusHandler(), common.MetricExtraGasLimitRelayedTx, node.CoreComponentsHolder.EconomicsData().MinGasLimit())
	metrics.SaveStringMetric(node.StatusCoreComponents.AppStatusHandler(), common.MetricRewardsTopUpGradientPoint, node.CoreComponentsHolder.EconomicsData().RewardsTopUpGradientPoint().String())
	metrics.SaveStringMetric(node.StatusCoreComponents.AppStatusHandler(), common.MetricTopUpFactor, fmt.Sprintf("%g", node.CoreComponentsHolder.EconomicsData().RewardsTopUpFactor()))
	metrics.SaveStringMetric(node.StatusCoreComponents.AppStatusHandler(), common.MetricGasPriceModifier, fmt.Sprintf("%g", node.CoreComponentsHolder.EconomicsData().GasPriceModifier()))
	metrics.SaveUint64Metric(node.StatusCoreComponents.AppStatusHandler(), common.MetricMaxGasPerTransaction, node.CoreComponentsHolder.EconomicsData().MaxGasLimitPerTx())
	if configs.PreferencesConfig.Preferences.FullArchive {
		metrics.SaveStringMetric(node.StatusCoreComponents.AppStatusHandler(), common.MetricPeerType, core.ObserverPeer.String())
		metrics.SaveStringMetric(node.StatusCoreComponents.AppStatusHandler(), common.MetricPeerSubType, core.FullHistoryObserver.String())
	}

	return nil
}
