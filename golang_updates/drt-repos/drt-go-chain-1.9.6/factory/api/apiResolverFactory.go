package api

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	logger "github.com/TerraDharitri/drt-go-chain-logger"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	"github.com/TerraDharitri/drt-go-chain-vm-common/parsers"
	datafield "github.com/TerraDharitri/drt-go-chain-vm-common/parsers/dataField"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/common/disabled"
	"github.com/TerraDharitri/drt-go-chain/common/operationmodes"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever/blockchain"
	"github.com/TerraDharitri/drt-go-chain/facade"
	"github.com/TerraDharitri/drt-go-chain/factory"
	"github.com/TerraDharitri/drt-go-chain/node/external"
	"github.com/TerraDharitri/drt-go-chain/node/external/blockAPI"
	"github.com/TerraDharitri/drt-go-chain/node/external/logs"
	"github.com/TerraDharitri/drt-go-chain/node/external/timemachine/fee"
	"github.com/TerraDharitri/drt-go-chain/node/external/transactionAPI"
	"github.com/TerraDharitri/drt-go-chain/node/trieIterators"
	trieIteratorsFactory "github.com/TerraDharitri/drt-go-chain/node/trieIterators/factory"
	"github.com/TerraDharitri/drt-go-chain/outport/process/alteredaccounts"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/coordinator"
	"github.com/TerraDharitri/drt-go-chain/process/factory/metachain"
	"github.com/TerraDharitri/drt-go-chain/process/factory/shard"
	"github.com/TerraDharitri/drt-go-chain/process/smartContract"
	"github.com/TerraDharitri/drt-go-chain/process/smartContract/builtInFunctions"
	"github.com/TerraDharitri/drt-go-chain/process/smartContract/hooks"
	"github.com/TerraDharitri/drt-go-chain/process/smartContract/hooks/counters"
	"github.com/TerraDharitri/drt-go-chain/process/txstatus"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain/state"
	"github.com/TerraDharitri/drt-go-chain/state/blockInfoProviders"
	disabledState "github.com/TerraDharitri/drt-go-chain/state/disabled"
	factoryState "github.com/TerraDharitri/drt-go-chain/state/factory"
	"github.com/TerraDharitri/drt-go-chain/state/storagePruningManager"
	"github.com/TerraDharitri/drt-go-chain/state/storagePruningManager/evictionWaitingList"
	"github.com/TerraDharitri/drt-go-chain/state/syncer"
	storageFactory "github.com/TerraDharitri/drt-go-chain/storage/factory"
	"github.com/TerraDharitri/drt-go-chain/storage/storageunit"
	trieFactory "github.com/TerraDharitri/drt-go-chain/trie/factory"
	"github.com/TerraDharitri/drt-go-chain/vm"
)

var log = logger.GetOrCreate("factory")

// ApiResolverArgs holds the argument needed to create an API resolver
type ApiResolverArgs struct {
	Configs              *config.Configs
	CoreComponents       factory.CoreComponentsHolder
	DataComponents       factory.DataComponentsHolder
	StateComponents      factory.StateComponentsHolder
	BootstrapComponents  factory.BootstrapComponentsHolder
	CryptoComponents     factory.CryptoComponentsHolder
	ProcessComponents    factory.ProcessComponentsHolder
	StatusCoreComponents factory.StatusCoreComponentsHolder
	StatusComponents     factory.StatusComponentsHolder
	GasScheduleNotifier  common.GasScheduleNotifierAPI
	Bootstrapper         process.Bootstrapper
	AllowVMQueriesChan   chan struct{}
	ProcessingMode       common.NodeProcessingMode
}

type scQueryServiceArgs struct {
	generalConfig              *config.Config
	epochConfig                *config.EpochConfig
	coreComponents             factory.CoreComponentsHolder
	stateComponents            factory.StateComponentsHolder
	dataComponents             factory.DataComponentsHolder
	processComponents          factory.ProcessComponentsHolder
	statusCoreComponents       factory.StatusCoreComponentsHolder
	gasScheduleNotifier        core.GasScheduleNotifier
	messageSigVerifier         vm.MessageSignVerifier
	systemSCConfig             *config.SystemSmartContractsConfig
	bootstrapper               process.Bootstrapper
	guardedAccountHandler      process.GuardedAccountHandler
	allowVMQueriesChan         chan struct{}
	workingDir                 string
	processingMode             common.NodeProcessingMode
	isInHistoricalBalancesMode bool
}

type scQueryElementArgs struct {
	generalConfig              *config.Config
	epochConfig                *config.EpochConfig
	coreComponents             factory.CoreComponentsHolder
	stateComponents            factory.StateComponentsHolder
	dataComponents             factory.DataComponentsHolder
	processComponents          factory.ProcessComponentsHolder
	statusCoreComponents       factory.StatusCoreComponentsHolder
	gasScheduleNotifier        core.GasScheduleNotifier
	messageSigVerifier         vm.MessageSignVerifier
	systemSCConfig             *config.SystemSmartContractsConfig
	bootstrapper               process.Bootstrapper
	guardedAccountHandler      process.GuardedAccountHandler
	allowVMQueriesChan         chan struct{}
	workingDir                 string
	index                      int
	processingMode             common.NodeProcessingMode
	isInHistoricalBalancesMode bool
}

// CreateApiResolver is able to create an ApiResolver instance that will solve the REST API requests through the node facade
// TODO: refactor to further decrease node's codebase
func CreateApiResolver(args *ApiResolverArgs) (facade.ApiResolver, error) {
	apiWorkingDir := filepath.Join(args.Configs.FlagsConfig.WorkingDir, common.TemporaryPath)
	argsSCQuery := &scQueryServiceArgs{
		generalConfig:              args.Configs.GeneralConfig,
		epochConfig:                args.Configs.EpochConfig,
		coreComponents:             args.CoreComponents,
		dataComponents:             args.DataComponents,
		stateComponents:            args.StateComponents,
		processComponents:          args.ProcessComponents,
		statusCoreComponents:       args.StatusCoreComponents,
		gasScheduleNotifier:        args.GasScheduleNotifier,
		messageSigVerifier:         args.CryptoComponents.MessageSignVerifier(),
		systemSCConfig:             args.Configs.SystemSCConfig,
		bootstrapper:               args.Bootstrapper,
		guardedAccountHandler:      args.BootstrapComponents.GuardedAccountHandler(),
		allowVMQueriesChan:         args.AllowVMQueriesChan,
		workingDir:                 apiWorkingDir,
		processingMode:             args.ProcessingMode,
		isInHistoricalBalancesMode: operationmodes.IsInHistoricalBalancesMode(args.Configs),
	}

	scQueryService, storageManagers, err := createScQueryService(argsSCQuery)
	if err != nil {
		return nil, err
	}

	pkConverter := args.CoreComponents.AddressPubKeyConverter()
	automaticCrawlerAddressesStrings := args.Configs.GeneralConfig.BuiltInFunctions.AutomaticCrawlerAddresses
	convertedAddresses, errDecode := factory.DecodeAddresses(pkConverter, automaticCrawlerAddressesStrings)
	if errDecode != nil {
		return nil, errDecode
	}

	dnsV2AddressesStrings := args.Configs.GeneralConfig.BuiltInFunctions.DNSV2Addresses
	convertedDNSV2Addresses, errDecode := factory.DecodeAddresses(pkConverter, dnsV2AddressesStrings)
	if errDecode != nil {
		return nil, errDecode
	}

	builtInFuncFactory, err := createBuiltinFuncs(
		args.GasScheduleNotifier,
		args.CoreComponents.InternalMarshalizer(),
		args.StateComponents.AccountsAdapterAPI(),
		args.BootstrapComponents.ShardCoordinator(),
		args.CoreComponents.EpochNotifier(),
		args.CoreComponents.EnableEpochsHandler(),
		args.BootstrapComponents.GuardedAccountHandler(),
		convertedAddresses,
		args.Configs.GeneralConfig.BuiltInFunctions.MaxNumAddressesInTransferRole,
		convertedDNSV2Addresses,
	)
	if err != nil {
		return nil, err
	}

	dcdtTransferParser, err := parsers.NewDCDTTransferParser(args.CoreComponents.InternalMarshalizer())
	if err != nil {
		return nil, err
	}

	argsTxTypeHandler := coordinator.ArgNewTxTypeHandler{
		PubkeyConverter:     args.CoreComponents.AddressPubKeyConverter(),
		ShardCoordinator:    args.ProcessComponents.ShardCoordinator(),
		BuiltInFunctions:    builtInFuncFactory.BuiltInFunctionContainer(),
		ArgumentParser:      parsers.NewCallArgsParser(),
		DCDTTransferParser:  dcdtTransferParser,
		EnableEpochsHandler: args.CoreComponents.EnableEpochsHandler(),
	}
	txTypeHandler, err := coordinator.NewTxTypeHandler(argsTxTypeHandler)
	if err != nil {
		return nil, err
	}

	accountsWrapper := &trieIterators.AccountsWrapper{
		Mutex:           &sync.Mutex{},
		AccountsAdapter: args.StateComponents.AccountsAdapterAPI(),
	}

	argsProcessors := trieIterators.ArgTrieIteratorProcessor{
		ShardID:            args.BootstrapComponents.ShardCoordinator().SelfId(),
		Accounts:           accountsWrapper,
		PublicKeyConverter: args.CoreComponents.AddressPubKeyConverter(),
		QueryService:       scQueryService,
	}
	totalStakedValueHandler, err := trieIteratorsFactory.CreateTotalStakedValueHandler(argsProcessors)
	if err != nil {
		return nil, err
	}

	directStakedListHandler, err := trieIteratorsFactory.CreateDirectStakedListHandler(argsProcessors)
	if err != nil {
		return nil, err
	}

	delegatedListHandler, err := trieIteratorsFactory.CreateDelegatedListHandler(argsProcessors)
	if err != nil {
		return nil, err
	}

	feeComputer, err := fee.NewFeeComputer(args.CoreComponents.EconomicsData())
	if err != nil {
		return nil, err
	}

	logsFacade, err := createLogsFacade(args)
	if err != nil {
		return nil, err
	}

	argsDataFieldParser := &datafield.ArgsOperationDataFieldParser{
		AddressLength: args.CoreComponents.AddressPubKeyConverter().Len(),
		Marshalizer:   args.CoreComponents.InternalMarshalizer(),
	}
	dataFieldParser, err := datafield.NewOperationDataFieldParser(argsDataFieldParser)
	if err != nil {
		return nil, err
	}

	argsAPITransactionProc := &transactionAPI.ArgAPITransactionProcessor{
		RoundDuration:            args.CoreComponents.GenesisNodesSetup().GetRoundDuration(),
		GenesisTime:              args.CoreComponents.GenesisTime(),
		Marshalizer:              args.CoreComponents.InternalMarshalizer(),
		AddressPubKeyConverter:   args.CoreComponents.AddressPubKeyConverter(),
		ShardCoordinator:         args.ProcessComponents.ShardCoordinator(),
		HistoryRepository:        args.ProcessComponents.HistoryRepository(),
		StorageService:           args.DataComponents.StorageService(),
		DataPool:                 args.DataComponents.Datapool(),
		Uint64ByteSliceConverter: args.CoreComponents.Uint64ByteSliceConverter(),
		FeeComputer:              feeComputer,
		TxTypeHandler:            txTypeHandler,
		LogsFacade:               logsFacade,
		DataFieldParser:          dataFieldParser,
		TxMarshaller:             args.CoreComponents.TxMarshalizer(),
		EnableEpochsHandler:      args.CoreComponents.EnableEpochsHandler(),
	}
	apiTransactionProcessor, err := transactionAPI.NewAPITransactionProcessor(argsAPITransactionProc)
	if err != nil {
		return nil, err
	}

	apiBlockProcessor, err := createAPIBlockProcessor(args, apiTransactionProcessor)
	if err != nil {
		return nil, err
	}

	apiInternalBlockProcessor, err := createAPIInternalBlockProcessor(args, apiTransactionProcessor)
	if err != nil {
		return nil, err
	}

	argsApiResolver := external.ArgNodeApiResolver{
		SCQueryService:           scQueryService,
		StatusMetricsHandler:     args.StatusCoreComponents.StatusMetrics(),
		APITransactionEvaluator:  args.ProcessComponents.APITransactionEvaluator(),
		TotalStakedValueHandler:  totalStakedValueHandler,
		DirectStakedListHandler:  directStakedListHandler,
		DelegatedListHandler:     delegatedListHandler,
		APITransactionHandler:    apiTransactionProcessor,
		APIBlockHandler:          apiBlockProcessor,
		APIInternalBlockHandler:  apiInternalBlockProcessor,
		GenesisNodesSetupHandler: args.CoreComponents.GenesisNodesSetup(),
		ValidatorPubKeyConverter: args.CoreComponents.ValidatorPubKeyConverter(),
		AccountsParser:           args.ProcessComponents.AccountsParser(),
		GasScheduleNotifier:      args.GasScheduleNotifier,
		ManagedPeersMonitor:      args.StatusComponents.ManagedPeersMonitor(),
		PublicKey:                args.CryptoComponents.PublicKeyString(),
		NodesCoordinator:         args.ProcessComponents.NodesCoordinator(),
		StorageManagers:          storageManagers,
	}

	return external.NewNodeApiResolver(argsApiResolver)
}

func createScQueryService(
	args *scQueryServiceArgs,
) (process.SCQueryService, []common.StorageManager, error) {
	numConcurrentVms := args.generalConfig.VirtualMachine.Querying.NumConcurrentVMs
	if numConcurrentVms < 1 {
		return nil, nil, fmt.Errorf("VirtualMachine.Querying.NumConcurrentVms should be a positive number more than 1")
	}

	argsQueryElem := &scQueryElementArgs{
		generalConfig:              args.generalConfig,
		epochConfig:                args.epochConfig,
		coreComponents:             args.coreComponents,
		stateComponents:            args.stateComponents,
		dataComponents:             args.dataComponents,
		processComponents:          args.processComponents,
		statusCoreComponents:       args.statusCoreComponents,
		gasScheduleNotifier:        args.gasScheduleNotifier,
		messageSigVerifier:         args.messageSigVerifier,
		systemSCConfig:             args.systemSCConfig,
		bootstrapper:               args.bootstrapper,
		guardedAccountHandler:      args.guardedAccountHandler,
		allowVMQueriesChan:         args.allowVMQueriesChan,
		workingDir:                 args.workingDir,
		index:                      0,
		processingMode:             args.processingMode,
		isInHistoricalBalancesMode: args.isInHistoricalBalancesMode,
	}

	var err error
	var scQueryService process.SCQueryService
	var storageManager common.StorageManager
	storageManagers := make([]common.StorageManager, 0, numConcurrentVms)

	list := make([]process.SCQueryService, 0, numConcurrentVms)
	for i := 0; i < numConcurrentVms; i++ {
		argsQueryElem.index = i
		scQueryService, storageManager, err = createScQueryElement(*argsQueryElem)
		if err != nil {
			return nil, nil, err
		}

		list = append(list, scQueryService)
		storageManagers = append(storageManagers, storageManager)
	}

	sqQueryDispatcher, err := smartContract.NewScQueryServiceDispatcher(list)
	if err != nil {
		return nil, nil, err
	}

	return sqQueryDispatcher, storageManagers, nil
}

func createScQueryElement(
	args scQueryElementArgs,
) (process.SCQueryService, common.StorageManager, error) {
	var err error

	selfShardID := args.processComponents.ShardCoordinator().SelfId()

	pkConverter := args.coreComponents.AddressPubKeyConverter()
	automaticCrawlerAddressesStrings := args.generalConfig.BuiltInFunctions.AutomaticCrawlerAddresses
	convertedAddresses, errDecode := factory.DecodeAddresses(pkConverter, automaticCrawlerAddressesStrings)
	if errDecode != nil {
		return nil, nil, errDecode
	}

	dnsV2AddressesStrings := args.generalConfig.BuiltInFunctions.DNSV2Addresses
	convertedDNSV2Addresses, errDecode := factory.DecodeAddresses(pkConverter, dnsV2AddressesStrings)
	if errDecode != nil {
		return nil, nil, errDecode
	}

	apiBlockchain, err := createBlockchainForScQuery(selfShardID)
	if err != nil {
		return nil, nil, err
	}

	accountsAdapterApi, storageManager, err := createNewAccountsAdapterApi(args, apiBlockchain)
	if err != nil {
		return nil, nil, err
	}

	builtInFuncFactory, err := createBuiltinFuncs(
		args.gasScheduleNotifier,
		args.coreComponents.InternalMarshalizer(),
		accountsAdapterApi,
		args.processComponents.ShardCoordinator(),
		args.coreComponents.EpochNotifier(),
		args.coreComponents.EnableEpochsHandler(),
		args.guardedAccountHandler,
		convertedAddresses,
		args.generalConfig.BuiltInFunctions.MaxNumAddressesInTransferRole,
		convertedDNSV2Addresses,
	)
	if err != nil {
		return nil, nil, err
	}

	cacherCfg := storageFactory.GetCacherFromConfig(args.generalConfig.SmartContractDataPool)
	smartContractsCache, err := storageunit.NewCache(cacherCfg)
	if err != nil {
		return nil, nil, err
	}

	scStorage := args.generalConfig.SmartContractsStorageForSCQuery
	scStorage.DB.FilePath += fmt.Sprintf("%d", args.index)
	argsHook := hooks.ArgBlockChainHook{
		PubkeyConv:               args.coreComponents.AddressPubKeyConverter(),
		StorageService:           args.dataComponents.StorageService(),
		ShardCoordinator:         args.processComponents.ShardCoordinator(),
		Marshalizer:              args.coreComponents.InternalMarshalizer(),
		Uint64Converter:          args.coreComponents.Uint64ByteSliceConverter(),
		BuiltInFunctions:         builtInFuncFactory.BuiltInFunctionContainer(),
		NFTStorageHandler:        builtInFuncFactory.NFTStorageHandler(),
		GlobalSettingsHandler:    builtInFuncFactory.DCDTGlobalSettingsHandler(),
		DataPool:                 args.dataComponents.Datapool(),
		ConfigSCStorage:          scStorage,
		CompiledSCPool:           smartContractsCache,
		WorkingDir:               args.workingDir,
		EpochNotifier:            args.coreComponents.EpochNotifier(),
		EnableEpochsHandler:      args.coreComponents.EnableEpochsHandler(),
		NilCompiledSCStore:       true,
		GasSchedule:              args.gasScheduleNotifier,
		Counter:                  counters.NewDisabledCounter(),
		MissingTrieNodesNotifier: syncer.NewMissingTrieNodesNotifier(),
		Accounts:                 accountsAdapterApi,
		BlockChain:               apiBlockchain,
	}

	var vmFactory process.VirtualMachinesContainerFactory
	maxGasForVmQueries := args.generalConfig.VirtualMachine.GasConfig.ShardMaxGasPerVmQuery
	if selfShardID == core.MetachainShardId {
		maxGasForVmQueries = args.generalConfig.VirtualMachine.GasConfig.MetaMaxGasPerVmQuery
		vmFactory, err = createMetaVmContainerFactory(args, argsHook)
	} else {
		vmFactory, err = createShardVmContainerFactory(args, argsHook)
	}
	if err != nil {
		return nil, nil, err
	}

	log.Debug("maximum gas per VM Query", "value", maxGasForVmQueries)

	vmContainer, err := vmFactory.Create()
	if err != nil {
		return nil, nil, err
	}

	err = vmFactory.BlockChainHookImpl().SetVMContainer(vmContainer)
	if err != nil {
		return nil, nil, err
	}

	err = builtInFuncFactory.SetPayableHandler(vmFactory.BlockChainHookImpl())
	if err != nil {
		return nil, nil, err
	}

	argsNewSCQueryService := smartContract.ArgsNewSCQueryService{
		VmContainer:                vmContainer,
		EconomicsFee:               args.coreComponents.EconomicsData(),
		BlockChainHook:             vmFactory.BlockChainHookImpl(),
		MainBlockChain:             args.dataComponents.Blockchain(),
		APIBlockChain:              apiBlockchain,
		WasmVMChangeLocker:         args.coreComponents.WasmVMChangeLocker(),
		Bootstrapper:               args.bootstrapper,
		AllowExternalQueriesChan:   args.allowVMQueriesChan,
		MaxGasLimitPerQuery:        maxGasForVmQueries,
		HistoryRepository:          args.processComponents.HistoryRepository(),
		ShardCoordinator:           args.processComponents.ShardCoordinator(),
		StorageService:             args.dataComponents.StorageService(),
		Marshaller:                 args.coreComponents.InternalMarshalizer(),
		Hasher:                     args.coreComponents.Hasher(),
		Uint64ByteSliceConverter:   args.coreComponents.Uint64ByteSliceConverter(),
		IsInHistoricalBalancesMode: args.isInHistoricalBalancesMode,
	}

	scQueryService, err := smartContract.NewSCQueryService(argsNewSCQueryService)

	return scQueryService, storageManager, err
}

func createBlockchainForScQuery(selfShardID uint32) (data.ChainHandler, error) {
	isMetachain := selfShardID == core.MetachainShardId
	if isMetachain {
		return blockchain.NewMetaChain(disabled.NewAppStatusHandler())
	}

	return blockchain.NewBlockChain(disabled.NewAppStatusHandler())
}

func createMetaVmContainerFactory(args scQueryElementArgs, argsHook hooks.ArgBlockChainHook) (process.VirtualMachinesContainerFactory, error) {
	blockChainHookImpl, errBlockChainHook := hooks.NewBlockChainHookImpl(argsHook)
	if errBlockChainHook != nil {
		return nil, errBlockChainHook
	}

	argsNewVmFactory := metachain.ArgsNewVMContainerFactory{
		BlockChainHook:      blockChainHookImpl,
		PubkeyConv:          argsHook.PubkeyConv,
		Economics:           args.coreComponents.EconomicsData(),
		MessageSignVerifier: args.messageSigVerifier,
		GasSchedule:         args.gasScheduleNotifier,
		NodesConfigProvider: args.coreComponents.GenesisNodesSetup(),
		Hasher:              args.coreComponents.Hasher(),
		Marshalizer:         args.coreComponents.InternalMarshalizer(),
		SystemSCConfig:      args.systemSCConfig,
		ValidatorAccountsDB: args.stateComponents.PeerAccounts(),
		UserAccountsDB:      args.stateComponents.AccountsAdapterAPI(),
		ChanceComputer:      args.coreComponents.Rater(),
		ShardCoordinator:    args.processComponents.ShardCoordinator(),
		EnableEpochsHandler: args.coreComponents.EnableEpochsHandler(),
		NodesCoordinator:    args.processComponents.NodesCoordinator(),
	}
	vmFactory, err := metachain.NewVMContainerFactory(argsNewVmFactory)
	if err != nil {
		return nil, err
	}

	return vmFactory, nil
}

func createShardVmContainerFactory(args scQueryElementArgs, argsHook hooks.ArgBlockChainHook) (process.VirtualMachinesContainerFactory, error) {
	queryVirtualMachineConfig := args.generalConfig.VirtualMachine.Querying.VirtualMachineConfig
	dcdtTransferParser, errParser := parsers.NewDCDTTransferParser(args.coreComponents.InternalMarshalizer())
	if errParser != nil {
		return nil, errParser
	}

	blockChainHookImpl, errBlockChainHook := hooks.NewBlockChainHookImpl(argsHook)
	if errBlockChainHook != nil {
		return nil, errBlockChainHook
	}

	argsNewVMFactory := shard.ArgVMContainerFactory{
		BlockChainHook:      blockChainHookImpl,
		BuiltInFunctions:    argsHook.BuiltInFunctions,
		Config:              queryVirtualMachineConfig,
		BlockGasLimit:       args.coreComponents.EconomicsData().MaxGasLimitPerBlock(args.processComponents.ShardCoordinator().SelfId()),
		GasSchedule:         args.gasScheduleNotifier,
		EpochNotifier:       args.coreComponents.EpochNotifier(),
		EnableEpochsHandler: args.coreComponents.EnableEpochsHandler(),
		WasmVMChangeLocker:  args.coreComponents.WasmVMChangeLocker(),
		DCDTTransferParser:  dcdtTransferParser,
		Hasher:              args.coreComponents.Hasher(),
		PubKeyConverter:     args.coreComponents.AddressPubKeyConverter(),
	}

	log.Debug("apiResolver: enable epoch for sc deploy", "epoch", args.epochConfig.EnableEpochs.SCDeployEnableEpoch)
	log.Debug("apiResolver: enable epoch for ahead of time gas usage", "epoch", args.epochConfig.EnableEpochs.AheadOfTimeGasUsageEnableEpoch)
	log.Debug("apiResolver: enable epoch for repair callback", "epoch", args.epochConfig.EnableEpochs.RepairCallbackEnableEpoch)

	vmFactory, err := shard.NewVMContainerFactory(argsNewVMFactory)
	if err != nil {
		return nil, err
	}

	return vmFactory, nil
}

func createNewAccountsAdapterApi(args scQueryElementArgs, chainHandler data.ChainHandler) (state.AccountsAdapterAPI, common.StorageManager, error) {
	argsAccCreator := factoryState.ArgsAccountCreator{
		Hasher:              args.coreComponents.Hasher(),
		Marshaller:          args.coreComponents.InternalMarshalizer(),
		EnableEpochsHandler: args.coreComponents.EnableEpochsHandler(),
	}
	accountFactory, err := factoryState.NewAccountCreator(argsAccCreator)
	if err != nil {
		return nil, nil, err
	}

	storagePruning, err := newStoragePruningManager(args)
	if err != nil {
		return nil, nil, err
	}
	storageService := args.dataComponents.StorageService()
	trieStorer, err := storageService.GetStorer(dataRetriever.UserAccountsUnit)
	if err != nil {
		return nil, nil, err
	}

	trieFactoryArgs := trieFactory.TrieFactoryArgs{
		Marshalizer:              args.coreComponents.InternalMarshalizer(),
		Hasher:                   args.coreComponents.Hasher(),
		PathManager:              args.coreComponents.PathHandler(),
		TrieStorageManagerConfig: args.generalConfig.TrieStorageManagerConfig,
	}
	trFactory, err := trieFactory.NewTrieFactory(trieFactoryArgs)
	if err != nil {
		return nil, nil, err
	}

	trieCreatorArgs := trieFactory.TrieCreateArgs{
		MainStorer:          trieStorer,
		PruningEnabled:      args.generalConfig.StateTriesConfig.AccountsStatePruningEnabled,
		MaxTrieLevelInMem:   args.generalConfig.StateTriesConfig.MaxStateTrieLevelInMemory,
		SnapshotsEnabled:    args.generalConfig.StateTriesConfig.SnapshotsEnabled,
		IdleProvider:        args.coreComponents.ProcessStatusHandler(),
		Identifier:          dataRetriever.UserAccountsUnit.String(),
		EnableEpochsHandler: args.coreComponents.EnableEpochsHandler(),
		StatsCollector:      args.statusCoreComponents.StateStatsHandler(),
	}
	trieStorageManager, merkleTrie, err := trFactory.Create(trieCreatorArgs)
	if err != nil {
		return nil, nil, err
	}

	argsAPIAccountsDB := state.ArgsAccountsDB{
		Trie:                  merkleTrie,
		Hasher:                args.coreComponents.Hasher(),
		Marshaller:            args.coreComponents.InternalMarshalizer(),
		AccountFactory:        accountFactory,
		StoragePruningManager: storagePruning,
		AddressConverter:      args.coreComponents.AddressPubKeyConverter(),
		SnapshotsManager:      disabledState.NewDisabledSnapshotsManager(),
	}

	provider, err := blockInfoProviders.NewCurrentBlockInfo(chainHandler)
	if err != nil {
		return nil, nil, err
	}

	accounts, err := state.NewAccountsDB(argsAPIAccountsDB)
	if err != nil {
		return nil, nil, err
	}

	accountsDB, err := state.NewAccountsDBApi(accounts, provider)

	return accountsDB, trieStorageManager, err
}

func newStoragePruningManager(args scQueryElementArgs) (state.StoragePruningManager, error) {
	argsMemEviction := evictionWaitingList.MemoryEvictionWaitingListArgs{
		RootHashesSize: args.generalConfig.EvictionWaitingList.RootHashesSize,
		HashesSize:     args.generalConfig.EvictionWaitingList.HashesSize,
	}
	trieEvictionWaitingList, err := evictionWaitingList.NewMemoryEvictionWaitingList(argsMemEviction)
	if err != nil {
		return nil, err
	}

	storagePruning, err := storagePruningManager.NewStoragePruningManager(
		trieEvictionWaitingList,
		args.generalConfig.TrieStorageManagerConfig.PruningBufferLen,
	)
	if err != nil {
		return nil, err
	}

	return storagePruning, nil
}

func createBuiltinFuncs(
	gasScheduleNotifier core.GasScheduleNotifier,
	marshalizer marshal.Marshalizer,
	accnts state.AccountsAdapter,
	shardCoordinator sharding.Coordinator,
	epochNotifier vmcommon.EpochNotifier,
	enableEpochsHandler vmcommon.EnableEpochsHandler,
	guardedAccountHandler vmcommon.GuardedAccountHandler,
	automaticCrawlerAddresses [][]byte,
	maxNumAddressesInTransferRole uint32,
	dnsV2Addresses [][]byte,
) (vmcommon.BuiltInFunctionFactory, error) {
	mapDNSV2Addresses := make(map[string]struct{})
	for _, address := range dnsV2Addresses {
		mapDNSV2Addresses[string(address)] = struct{}{}
	}

	argsBuiltIn := builtInFunctions.ArgsCreateBuiltInFunctionContainer{
		GasSchedule:               gasScheduleNotifier,
		MapDNSAddresses:           make(map[string]struct{}),
		MapDNSV2Addresses:         mapDNSV2Addresses,
		Marshalizer:               marshalizer,
		Accounts:                  accnts,
		ShardCoordinator:          shardCoordinator,
		EpochNotifier:             epochNotifier,
		EnableEpochsHandler:       enableEpochsHandler,
		GuardedAccountHandler:     guardedAccountHandler,
		AutomaticCrawlerAddresses: automaticCrawlerAddresses,
		MaxNumNodesInTransferRole: maxNumAddressesInTransferRole,
	}
	return builtInFunctions.CreateBuiltInFunctionsFactory(argsBuiltIn)
}

func createAPIBlockProcessor(args *ApiResolverArgs, apiTransactionHandler external.APITransactionHandler) (blockAPI.APIBlockHandler, error) {
	blockApiArgs, err := createAPIBlockProcessorArgs(args, apiTransactionHandler)
	if err != nil {
		return nil, err
	}

	return blockAPI.CreateAPIBlockProcessor(blockApiArgs)
}

func createAPIInternalBlockProcessor(args *ApiResolverArgs, apiTransactionHandler external.APITransactionHandler) (blockAPI.APIInternalBlockHandler, error) {
	blockApiArgs, err := createAPIBlockProcessorArgs(args, apiTransactionHandler)
	if err != nil {
		return nil, err
	}

	return blockAPI.CreateAPIInternalBlockProcessor(blockApiArgs)
}

func createAPIBlockProcessorArgs(args *ApiResolverArgs, apiTransactionHandler external.APITransactionHandler) (*blockAPI.ArgAPIBlockProcessor, error) {
	statusComputer, err := txstatus.NewStatusComputer(
		args.ProcessComponents.ShardCoordinator().SelfId(),
		args.CoreComponents.Uint64ByteSliceConverter(),
		args.DataComponents.StorageService(),
	)
	if err != nil {
		return nil, errors.New("error creating transaction status computer " + err.Error())
	}

	logsFacade, err := createLogsFacade(args)
	if err != nil {
		return nil, err
	}

	alteredAccountsProvider, err := alteredaccounts.NewAlteredAccountsProvider(alteredaccounts.ArgsAlteredAccountsProvider{
		ShardCoordinator:       args.ProcessComponents.ShardCoordinator(),
		AddressConverter:       args.CoreComponents.AddressPubKeyConverter(),
		AccountsDB:             args.StateComponents.AccountsAdapterAPI(),
		DcdtDataStorageHandler: args.ProcessComponents.DCDTDataStorageHandlerForAPI(),
	})
	if err != nil {
		return nil, err
	}

	blockApiArgs := &blockAPI.ArgAPIBlockProcessor{
		SelfShardID:                  args.ProcessComponents.ShardCoordinator().SelfId(),
		Store:                        args.DataComponents.StorageService(),
		Marshalizer:                  args.CoreComponents.InternalMarshalizer(),
		Uint64ByteSliceConverter:     args.CoreComponents.Uint64ByteSliceConverter(),
		HistoryRepo:                  args.ProcessComponents.HistoryRepository(),
		APITransactionHandler:        apiTransactionHandler,
		StatusComputer:               statusComputer,
		AddressPubkeyConverter:       args.CoreComponents.AddressPubKeyConverter(),
		Hasher:                       args.CoreComponents.Hasher(),
		LogsFacade:                   logsFacade,
		ReceiptsRepository:           args.ProcessComponents.ReceiptsRepository(),
		AlteredAccountsProvider:      alteredAccountsProvider,
		AccountsRepository:           args.StateComponents.AccountsRepository(),
		ScheduledTxsExecutionHandler: args.ProcessComponents.ScheduledTxsExecutionHandler(),
		EnableEpochsHandler:          args.CoreComponents.EnableEpochsHandler(),
		ProofsPool:                   args.DataComponents.Datapool().Proofs(),
		BlockChain:                   args.DataComponents.Blockchain(),
	}

	return blockApiArgs, nil
}

func createLogsFacade(args *ApiResolverArgs) (factory.LogsFacade, error) {
	return logs.NewLogsFacade(logs.ArgsNewLogsFacade{
		StorageService:  args.DataComponents.StorageService(),
		Marshaller:      args.CoreComponents.InternalMarshalizer(),
		PubKeyConverter: args.CoreComponents.AddressPubKeyConverter(),
	})
}
