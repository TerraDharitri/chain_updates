package core

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/alarm"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/core/nodetype"
	"github.com/TerraDharitri/drt-go-chain-core/core/versioning"
	"github.com/TerraDharitri/drt-go-chain-core/core/watchdog"
	"github.com/TerraDharitri/drt-go-chain-core/data/endProcess"
	"github.com/TerraDharitri/drt-go-chain-core/data/typeConverters"
	"github.com/TerraDharitri/drt-go-chain-core/data/typeConverters/uint64ByteSlice"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	hasherFactory "github.com/TerraDharitri/drt-go-chain-core/hashing/factory"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	marshalizerFactory "github.com/TerraDharitri/drt-go-chain-core/marshal/factory"
	logger "github.com/TerraDharitri/drt-go-chain-logger"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/common/chainparametersnotifier"
	"github.com/TerraDharitri/drt-go-chain/common/enablers"
	commonFactory "github.com/TerraDharitri/drt-go-chain/common/factory"
	"github.com/TerraDharitri/drt-go-chain/common/fieldsChecker"
	"github.com/TerraDharitri/drt-go-chain/common/forking"
	"github.com/TerraDharitri/drt-go-chain/common/graceperiod"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/consensus"
	"github.com/TerraDharitri/drt-go-chain/consensus/round"
	"github.com/TerraDharitri/drt-go-chain/epochStart/notifier"
	"github.com/TerraDharitri/drt-go-chain/errors"
	"github.com/TerraDharitri/drt-go-chain/factory"
	"github.com/TerraDharitri/drt-go-chain/ntp"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/economics"
	"github.com/TerraDharitri/drt-go-chain/process/rating"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
	"github.com/TerraDharitri/drt-go-chain/statusHandler"
	"github.com/TerraDharitri/drt-go-chain/storage"
	storageFactory "github.com/TerraDharitri/drt-go-chain/storage/factory"
)

var log = logger.GetOrCreate("factory")

// CoreComponentsFactoryArgs holds the arguments needed for creating a core components factory
type CoreComponentsFactoryArgs struct {
	Config              config.Config
	ConfigPathsHolder   config.ConfigurationPathsHolder
	EpochConfig         config.EpochConfig
	RoundConfig         config.RoundConfig
	RatingsConfig       config.RatingsConfig
	EconomicsConfig     config.EconomicsConfig
	ImportDbConfig      config.ImportDbConfig
	NodesConfig         config.NodesConfig
	WorkingDirectory    string
	ChanStopNodeProcess chan endProcess.ArgEndProcess
}

// coreComponentsFactory is responsible for creating the core components
type coreComponentsFactory struct {
	config              config.Config
	configPathsHolder   config.ConfigurationPathsHolder
	epochConfig         config.EpochConfig
	roundConfig         config.RoundConfig
	ratingsConfig       config.RatingsConfig
	economicsConfig     config.EconomicsConfig
	importDbConfig      config.ImportDbConfig
	nodesSetupConfig    config.NodesConfig
	workingDir          string
	chanStopNodeProcess chan endProcess.ArgEndProcess
}

// coreComponents is the DTO used for core components
type coreComponents struct {
	hasher                        hashing.Hasher
	txSignHasher                  hashing.Hasher
	internalMarshalizer           marshal.Marshalizer
	vmMarshalizer                 marshal.Marshalizer
	txSignMarshalizer             marshal.Marshalizer
	uint64ByteSliceConverter      typeConverters.Uint64ByteSliceConverter
	addressPubKeyConverter        core.PubkeyConverter
	validatorPubKeyConverter      core.PubkeyConverter
	pathHandler                   storage.PathManagerHandler
	syncTimer                     ntp.SyncTimer
	roundHandler                  consensus.RoundHandler
	alarmScheduler                core.TimersScheduler
	watchdog                      core.WatchdogTimer
	nodesSetupHandler             sharding.GenesisNodesSetupHandler
	economicsData                 process.EconomicsDataHandler
	apiEconomicsData              process.EconomicsDataHandler
	ratingsData                   process.RatingsInfoHandler
	rater                         sharding.PeerAccountListAndRatingHandler
	nodesShuffler                 nodesCoordinator.NodesShuffler
	txVersionChecker              process.TxVersionCheckerHandler
	genesisTime                   time.Time
	chainID                       string
	minTransactionVersion         uint32
	epochNotifier                 process.EpochNotifier
	roundNotifier                 process.RoundNotifier
	chainParametersSubscriber     process.ChainParametersSubscriber
	enableRoundsHandler           process.EnableRoundsHandler
	epochStartNotifierWithConfirm factory.EpochStartNotifierWithConfirm
	chanStopNodeProcess           chan endProcess.ArgEndProcess
	nodeTypeProvider              core.NodeTypeProviderHandler
	encodedAddressLen             uint32
	wasmVMChangeLocker            common.Locker
	processStatusHandler          common.ProcessStatusHandler
	hardforkTriggerPubKey         []byte
	enableEpochsHandler           common.EnableEpochsHandler
	chainParametersHandler        process.ChainParametersHandler
	fieldsSizeChecker             common.FieldsSizeChecker
	epochChangeGracePeriodHandler common.EpochChangeGracePeriodHandler
}

// NewCoreComponentsFactory initializes the factory which is responsible to creating core components
func NewCoreComponentsFactory(args CoreComponentsFactoryArgs) (*coreComponentsFactory, error) {
	return &coreComponentsFactory{
		config:              args.Config,
		configPathsHolder:   args.ConfigPathsHolder,
		epochConfig:         args.EpochConfig,
		roundConfig:         args.RoundConfig,
		ratingsConfig:       args.RatingsConfig,
		importDbConfig:      args.ImportDbConfig,
		economicsConfig:     args.EconomicsConfig,
		workingDir:          args.WorkingDirectory,
		chanStopNodeProcess: args.ChanStopNodeProcess,
		nodesSetupConfig:    args.NodesConfig,
	}, nil
}

// Create creates the core components
func (ccf *coreComponentsFactory) Create() (*coreComponents, error) {
	hasher, err := hasherFactory.NewHasher(ccf.config.Hasher.Type)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errors.ErrHasherCreation, err.Error())
	}

	internalMarshalizer, err := marshalizerFactory.NewMarshalizer(ccf.config.Marshalizer.Type)
	if err != nil {
		return nil, fmt.Errorf("%w (internal): %s", errors.ErrMarshalizerCreation, err.Error())
	}

	vmMarshalizer, err := marshalizerFactory.NewMarshalizer(ccf.config.VmMarshalizer.Type)
	if err != nil {
		return nil, fmt.Errorf("%w (vm): %s", errors.ErrMarshalizerCreation, err.Error())
	}

	txSignMarshalizer, err := marshalizerFactory.NewMarshalizer(ccf.config.TxSignMarshalizer.Type)
	if err != nil {
		return nil, fmt.Errorf("%w (tx sign): %s", errors.ErrMarshalizerCreation, err.Error())
	}

	txSignHasher, err := hasherFactory.NewHasher(ccf.config.TxSignHasher.Type)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errors.ErrHasherCreation, err.Error())
	}

	uint64ByteSliceConverter := uint64ByteSlice.NewBigEndianConverter()

	addressPubkeyConverter, err := commonFactory.NewPubkeyConverter(ccf.config.AddressPubkeyConverter)
	if err != nil {
		return nil, fmt.Errorf("%w for AddressPubkeyConverter", err)
	}

	validatorPubkeyConverter, err := commonFactory.NewPubkeyConverter(ccf.config.ValidatorPubkeyConverter)
	if err != nil {
		return nil, fmt.Errorf("%w for AddressPubkeyConverter", err)
	}

	epochChangeGracePeriodHandler, err := graceperiod.NewEpochChangeGracePeriod(ccf.config.GeneralSettings.EpochChangeGracePeriodByEpoch)
	if err != nil {
		return nil, fmt.Errorf("%w for epochChangeGracePeriod", err)
	}

	pathHandler, err := storageFactory.CreatePathManager(
		storageFactory.ArgCreatePathManager{
			WorkingDir: ccf.workingDir,
			ChainID:    ccf.config.GeneralSettings.ChainID,
		},
	)
	if err != nil {
		return nil, err
	}

	syncer := ntp.NewSyncTime(ccf.config.NTPConfig, nil)
	syncer.StartSyncingTime()
	log.Debug("NTP average clock offset", "value", syncer.ClockOffset())

	epochNotifier := forking.NewGenericEpochNotifier()
	epochStartHandlerWithConfirm := notifier.NewEpochStartSubscriptionHandler()

	chainParametersNotifier := chainparametersnotifier.NewChainParametersNotifier()
	argsChainParametersHandler := sharding.ArgsChainParametersHolder{
		EpochStartEventNotifier: epochStartHandlerWithConfirm,
		ChainParameters:         ccf.config.GeneralSettings.ChainParametersByEpoch,
		ChainParametersNotifier: chainParametersNotifier,
	}
	chainParametersHandler, err := sharding.NewChainParametersHolder(argsChainParametersHandler)
	if err != nil {
		return nil, err
	}

	genesisNodesConfig, err := sharding.NewNodesSetup(
		ccf.nodesSetupConfig,
		chainParametersHandler,
		addressPubkeyConverter,
		validatorPubkeyConverter,
		ccf.config.GeneralSettings.GenesisMaxNumberOfShards,
	)
	if err != nil {
		return nil, err
	}

	startRound := int64(0)
	if ccf.config.Hardfork.AfterHardFork {
		log.Debug("changed genesis time after hardfork",
			"old genesis time", genesisNodesConfig.StartTime,
			"new genesis time", ccf.config.Hardfork.GenesisTime)
		genesisNodesConfig.StartTime = ccf.config.Hardfork.GenesisTime
		startRound = int64(ccf.config.Hardfork.StartRound)
	}

	if genesisNodesConfig.StartTime == 0 {
		time.Sleep(1000 * time.Millisecond)
		ntpTime := syncer.CurrentTime()
		genesisNodesConfig.StartTime = (ntpTime.Unix()/60 + 1) * 60
	}

	startTime := time.Unix(genesisNodesConfig.StartTime, 0)

	log.Info("start time",
		"formatted", startTime.Format("Mon Jan 2 15:04:05 MST 2006"),
		"seconds", startTime.Unix())

	genesisTime := time.Unix(genesisNodesConfig.StartTime, 0)
	roundHandler, err := round.NewRound(
		genesisTime,
		syncer.CurrentTime(),
		time.Millisecond*time.Duration(genesisNodesConfig.RoundDuration),
		syncer,
		startRound,
	)
	if err != nil {
		return nil, err
	}

	alarmScheduler := alarm.NewAlarmScheduler()
	// TODO: disable watchdog if block processing cutoff is enabled
	watchdogTimer, err := watchdog.NewWatchdog(alarmScheduler, ccf.chanStopNodeProcess, log)
	if err != nil {
		return nil, err
	}

	roundNotifier := forking.NewGenericRoundNotifier()
	enableRoundsHandler, err := enablers.NewEnableRoundsHandler(ccf.roundConfig, roundNotifier)
	if err != nil {
		return nil, err
	}

	enableEpochsHandler, err := enablers.NewEnableEpochsHandler(ccf.epochConfig.EnableEpochs, epochNotifier)
	if err != nil {
		return nil, err
	}

	wasmVMChangeLocker := &sync.RWMutex{}

	txVersionChecker := versioning.NewTxVersionChecker(ccf.config.GeneralSettings.MinTransactionVersion)

	// This shard coordinator uses a hardcoded selfId of 0 as it does not know its selfId.
	// Its main purpose is to validate the rewards config (protocol sustainability address shard against meta),
	// inside economics data and should not be used for another scope.
	// The real component will be created later on, as part of bootstrap components.
	shardCoordinator, err := sharding.NewMultiShardCoordinator(genesisNodesConfig.NumberOfShards(), 0)
	if err != nil {
		return nil, err
	}

	log.Trace("creating economics data components")
	argsNewEconomicsData := economics.ArgsNewEconomicsData{
		Economics:           &ccf.economicsConfig,
		EpochNotifier:       epochNotifier,
		EnableEpochsHandler: enableEpochsHandler,
		TxVersionChecker:    txVersionChecker,
		PubkeyConverter:     addressPubkeyConverter,
		ShardCoordinator:    shardCoordinator,
	}
	economicsData, err := economics.NewEconomicsData(argsNewEconomicsData)
	if err != nil {
		return nil, err
	}

	apiEconomicsData, err := economics.NewAPIEconomicsData(economicsData)
	if err != nil {
		return nil, err
	}

	log.Trace("creating ratings data")
	ratingDataArgs := rating.RatingsDataArg{
		Config:                    ccf.ratingsConfig,
		ChainParametersHolder:     chainParametersHandler,
		RoundDurationMilliseconds: genesisNodesConfig.RoundDuration,
		EpochNotifier:             epochNotifier,
	}
	ratingsData, err := rating.NewRatingsData(ratingDataArgs)
	if err != nil {
		return nil, err
	}

	rater, err := rating.NewBlockSigningRater(ratingsData)
	if err != nil {
		return nil, err
	}

	argsNodesShuffler := &nodesCoordinator.NodesShufflerArgs{
		ShuffleBetweenShards: true,
		MaxNodesEnableConfig: ccf.epochConfig.EnableEpochs.MaxNodesChangeEnableEpoch,
		EnableEpochsHandler:  enableEpochsHandler,
		EnableEpochs:         ccf.epochConfig.EnableEpochs,
	}

	nodesShuffler, err := nodesCoordinator.NewHashValidatorsShuffler(argsNodesShuffler)
	if err != nil {
		return nil, err
	}

	// set as observer at first - it will be updated when creating the nodes coordinator
	nodeTypeProvider := nodetype.NewNodeTypeProvider(core.NodeTypeObserver)

	pubKeyStr := ccf.config.Hardfork.PublicKeyToListenFrom
	pubKeyBytes, err := validatorPubkeyConverter.Decode(pubKeyStr)
	if err != nil {
		return nil, err
	}

	encodedAddressLen, err := computeEncodedAddressLen(addressPubkeyConverter)
	if err != nil {
		return nil, err
	}

	fieldsSizeChecker, err := fieldsChecker.NewFieldsSizeChecker(chainParametersHandler, hasher)
	if err != nil {
		return nil, err
	}

	return &coreComponents{
		hasher:                        hasher,
		txSignHasher:                  txSignHasher,
		internalMarshalizer:           internalMarshalizer,
		vmMarshalizer:                 vmMarshalizer,
		txSignMarshalizer:             txSignMarshalizer,
		uint64ByteSliceConverter:      uint64ByteSliceConverter,
		addressPubKeyConverter:        addressPubkeyConverter,
		validatorPubKeyConverter:      validatorPubkeyConverter,
		pathHandler:                   pathHandler,
		syncTimer:                     syncer,
		roundHandler:                  roundHandler,
		alarmScheduler:                alarmScheduler,
		watchdog:                      watchdogTimer,
		nodesSetupHandler:             genesisNodesConfig,
		economicsData:                 economicsData,
		apiEconomicsData:              apiEconomicsData,
		ratingsData:                   ratingsData,
		rater:                         rater,
		nodesShuffler:                 nodesShuffler,
		txVersionChecker:              txVersionChecker,
		genesisTime:                   genesisTime,
		chainID:                       ccf.config.GeneralSettings.ChainID,
		minTransactionVersion:         ccf.config.GeneralSettings.MinTransactionVersion,
		epochNotifier:                 epochNotifier,
		roundNotifier:                 roundNotifier,
		chainParametersSubscriber:     chainParametersNotifier,
		enableRoundsHandler:           enableRoundsHandler,
		epochStartNotifierWithConfirm: epochStartHandlerWithConfirm,
		chanStopNodeProcess:           ccf.chanStopNodeProcess,
		encodedAddressLen:             encodedAddressLen,
		nodeTypeProvider:              nodeTypeProvider,
		wasmVMChangeLocker:            wasmVMChangeLocker,
		processStatusHandler:          statusHandler.NewProcessStatusHandler(),
		hardforkTriggerPubKey:         pubKeyBytes,
		enableEpochsHandler:           enableEpochsHandler,
		chainParametersHandler:        chainParametersHandler,
		fieldsSizeChecker:             fieldsSizeChecker,
		epochChangeGracePeriodHandler: epochChangeGracePeriodHandler,
	}, nil
}

// Close closes all underlying components
func (cc *coreComponents) Close() error {
	if !check.IfNil(cc.alarmScheduler) {
		cc.alarmScheduler.Close()
	}
	if !check.IfNil(cc.syncTimer) {
		err := cc.syncTimer.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func computeEncodedAddressLen(converter core.PubkeyConverter) (uint32, error) {
	emptyAddress := bytes.Repeat([]byte{0}, converter.Len())
	encodedEmptyAddress, err := converter.Encode(emptyAddress)
	if err != nil {
		return 0, err
	}

	return uint32(len(encodedEmptyAddress)), nil
}
