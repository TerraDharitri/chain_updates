package components

import (
	"bytes"
	"sync"
	"time"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/common/chainparametersnotifier"
	"github.com/TerraDharitri/drt-go-chain/common/enablers"
	factoryPubKey "github.com/TerraDharitri/drt-go-chain/common/factory"
	"github.com/TerraDharitri/drt-go-chain/common/fieldsChecker"
	"github.com/TerraDharitri/drt-go-chain/common/forking"
	"github.com/TerraDharitri/drt-go-chain/common/graceperiod"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/consensus"
	"github.com/TerraDharitri/drt-go-chain/epochStart/notifier"
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
	"github.com/TerraDharitri/drt-go-chain/testscommon"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/nodetype"
	"github.com/TerraDharitri/drt-go-chain-core/core/versioning"
	"github.com/TerraDharitri/drt-go-chain-core/core/watchdog"
	"github.com/TerraDharitri/drt-go-chain-core/data/endProcess"
	"github.com/TerraDharitri/drt-go-chain-core/data/typeConverters"
	"github.com/TerraDharitri/drt-go-chain-core/data/typeConverters/uint64ByteSlice"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	hashingFactory "github.com/TerraDharitri/drt-go-chain-core/hashing/factory"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	marshalFactory "github.com/TerraDharitri/drt-go-chain-core/marshal/factory"
)

type coreComponentsHolder struct {
	closeHandler                  *closeHandler
	internalMarshaller            marshal.Marshalizer
	txMarshaller                  marshal.Marshalizer
	vmMarshaller                  marshal.Marshalizer
	hasher                        hashing.Hasher
	txSignHasher                  hashing.Hasher
	uint64SliceConverter          typeConverters.Uint64ByteSliceConverter
	addressPubKeyConverter        core.PubkeyConverter
	validatorPubKeyConverter      core.PubkeyConverter
	pathHandler                   storage.PathManagerHandler
	watchdog                      core.WatchdogTimer
	alarmScheduler                core.TimersScheduler
	syncTimer                     ntp.SyncTimer
	roundHandler                  consensus.RoundHandler
	economicsData                 process.EconomicsDataHandler
	apiEconomicsData              process.EconomicsDataHandler
	ratingsData                   process.RatingsInfoHandler
	rater                         sharding.PeerAccountListAndRatingHandler
	genesisNodesSetup             sharding.GenesisNodesSetupHandler
	nodesShuffler                 nodesCoordinator.NodesShuffler
	epochNotifier                 process.EpochNotifier
	enableRoundsHandler           process.EnableRoundsHandler
	roundNotifier                 process.RoundNotifier
	epochStartNotifierWithConfirm factory.EpochStartNotifierWithConfirm
	chanStopNodeProcess           chan endProcess.ArgEndProcess
	genesisTime                   time.Time
	chainID                       string
	minTransactionVersion         uint32
	txVersionChecker              process.TxVersionCheckerHandler
	encodedAddressLen             uint32
	nodeTypeProvider              core.NodeTypeProviderHandler
	wasmVMChangeLocker            common.Locker
	processStatusHandler          common.ProcessStatusHandler
	hardforkTriggerPubKey         []byte
	enableEpochsHandler           common.EnableEpochsHandler
	chainParametersSubscriber     process.ChainParametersSubscriber
	chainParametersHandler        process.ChainParametersHandler
	fieldsSizeChecker             common.FieldsSizeChecker
	epochChangeGracePeriodHandler common.EpochChangeGracePeriodHandler
}

// ArgsCoreComponentsHolder will hold arguments needed for the core components holder
type ArgsCoreComponentsHolder struct {
	Config              config.Config
	EnableEpochsConfig  config.EnableEpochs
	RoundsConfig        config.RoundConfig
	EconomicsConfig     config.EconomicsConfig
	RatingConfig        config.RatingsConfig
	ChanStopNodeProcess chan endProcess.ArgEndProcess
	InitialRound        int64
	NodesSetupPath      string
	GasScheduleFilename string
	NumShards           uint32
	WorkingDir          string

	MinNodesPerShard            uint32
	ConsensusGroupSize          uint32
	MinNodesMeta                uint32
	MetaChainConsensusGroupSize uint32
	RoundDurationInMs           uint64
}

// CreateCoreComponents will create a new instance of factory.CoreComponentsHolder
func CreateCoreComponents(args ArgsCoreComponentsHolder) (*coreComponentsHolder, error) {
	var err error
	instance := &coreComponentsHolder{
		closeHandler: NewCloseHandler(),
	}

	instance.internalMarshaller, err = marshalFactory.NewMarshalizer(args.Config.Marshalizer.Type)
	if err != nil {
		return nil, err
	}
	instance.txMarshaller, err = marshalFactory.NewMarshalizer(args.Config.TxSignMarshalizer.Type)
	if err != nil {
		return nil, err
	}
	instance.vmMarshaller, err = marshalFactory.NewMarshalizer(args.Config.VmMarshalizer.Type)
	if err != nil {
		return nil, err
	}
	instance.hasher, err = hashingFactory.NewHasher(args.Config.Hasher.Type)
	if err != nil {
		return nil, err
	}
	instance.txSignHasher, err = hashingFactory.NewHasher(args.Config.TxSignHasher.Type)
	if err != nil {
		return nil, err
	}
	instance.uint64SliceConverter = uint64ByteSlice.NewBigEndianConverter()
	instance.addressPubKeyConverter, err = factoryPubKey.NewPubkeyConverter(args.Config.AddressPubkeyConverter)
	if err != nil {
		return nil, err
	}
	instance.validatorPubKeyConverter, err = factoryPubKey.NewPubkeyConverter(args.Config.ValidatorPubkeyConverter)
	if err != nil {
		return nil, err
	}

	instance.pathHandler, err = storageFactory.CreatePathManager(
		storageFactory.ArgCreatePathManager{
			WorkingDir: args.WorkingDir,
			ChainID:    args.Config.GeneralSettings.ChainID,
		},
	)
	if err != nil {
		return nil, err
	}

	instance.watchdog = &watchdog.DisabledWatchdog{}
	instance.alarmScheduler = &testscommon.AlarmSchedulerStub{}
	instance.syncTimer = &testscommon.SyncTimerStub{}

	instance.epochStartNotifierWithConfirm = notifier.NewEpochStartSubscriptionHandler()
	instance.chainParametersSubscriber = chainparametersnotifier.NewChainParametersNotifier()
	chainParametersNotifier := chainparametersnotifier.NewChainParametersNotifier()
	argsChainParametersHandler := sharding.ArgsChainParametersHolder{
		EpochStartEventNotifier: instance.epochStartNotifierWithConfirm,
		ChainParameters:         args.Config.GeneralSettings.ChainParametersByEpoch,
		ChainParametersNotifier: chainParametersNotifier,
	}
	instance.chainParametersHandler, err = sharding.NewChainParametersHolder(argsChainParametersHandler)
	if err != nil {
		return nil, err
	}

	instance.epochChangeGracePeriodHandler, err = graceperiod.NewEpochChangeGracePeriod(args.Config.GeneralSettings.EpochChangeGracePeriodByEpoch)
	if err != nil {
		return nil, err
	}

	var nodesSetup config.NodesConfig
	err = core.LoadJsonFile(&nodesSetup, args.NodesSetupPath)
	if err != nil {
		return nil, err
	}
	instance.genesisNodesSetup, err = sharding.NewNodesSetup(nodesSetup, instance.chainParametersHandler, instance.addressPubKeyConverter, instance.validatorPubKeyConverter, args.NumShards)
	if err != nil {
		return nil, err
	}

	roundDuration := time.Millisecond * time.Duration(instance.genesisNodesSetup.GetRoundDuration())
	instance.roundHandler = NewManualRoundHandler(instance.genesisNodesSetup.GetStartTime(), roundDuration, args.InitialRound)

	instance.wasmVMChangeLocker = &sync.RWMutex{}
	instance.txVersionChecker = versioning.NewTxVersionChecker(args.Config.GeneralSettings.MinTransactionVersion)
	instance.epochNotifier = forking.NewGenericEpochNotifier()
	instance.enableEpochsHandler, err = enablers.NewEnableEpochsHandler(args.EnableEpochsConfig, instance.epochNotifier)
	if err != nil {
		return nil, err
	}

	argsEconomicsHandler := economics.ArgsNewEconomicsData{
		TxVersionChecker:    instance.txVersionChecker,
		Economics:           &args.EconomicsConfig,
		EpochNotifier:       instance.epochNotifier,
		EnableEpochsHandler: instance.enableEpochsHandler,
		PubkeyConverter:     instance.addressPubKeyConverter,
		ShardCoordinator:    testscommon.NewMultiShardsCoordinatorMock(instance.genesisNodesSetup.NumberOfShards()),
	}

	instance.economicsData, err = economics.NewEconomicsData(argsEconomicsHandler)
	if err != nil {
		return nil, err
	}
	instance.apiEconomicsData = instance.economicsData

	instance.ratingsData, err = rating.NewRatingsData(rating.RatingsDataArg{
		EpochNotifier:             instance.epochNotifier,
		Config:                    args.RatingConfig,
		ChainParametersHolder:     instance.chainParametersHandler,
		RoundDurationMilliseconds: args.RoundDurationInMs,
	})
	if err != nil {
		return nil, err
	}

	instance.rater, err = rating.NewBlockSigningRater(instance.ratingsData)
	if err != nil {
		return nil, err
	}

	instance.nodesShuffler, err = nodesCoordinator.NewHashValidatorsShuffler(&nodesCoordinator.NodesShufflerArgs{
		ShuffleBetweenShards: true,
		MaxNodesEnableConfig: args.EnableEpochsConfig.MaxNodesChangeEnableEpoch,
		EnableEpochsHandler:  instance.enableEpochsHandler,
		EnableEpochs:         args.EnableEpochsConfig,
	})
	if err != nil {
		return nil, err
	}

	instance.roundNotifier = forking.NewGenericRoundNotifier()
	instance.enableRoundsHandler, err = enablers.NewEnableRoundsHandler(args.RoundsConfig, instance.roundNotifier)
	if err != nil {
		return nil, err
	}

	instance.chanStopNodeProcess = args.ChanStopNodeProcess
	instance.genesisTime = time.Unix(instance.genesisNodesSetup.GetStartTime(), 0)
	instance.chainID = args.Config.GeneralSettings.ChainID
	instance.minTransactionVersion = args.Config.GeneralSettings.MinTransactionVersion
	instance.encodedAddressLen, err = computeEncodedAddressLen(instance.addressPubKeyConverter)
	if err != nil {
		return nil, err
	}

	instance.nodeTypeProvider = nodetype.NewNodeTypeProvider(core.NodeTypeObserver)
	instance.processStatusHandler = statusHandler.NewProcessStatusHandler()

	pubKeyBytes, err := instance.validatorPubKeyConverter.Decode(args.Config.Hardfork.PublicKeyToListenFrom)
	if err != nil {
		return nil, err
	}
	instance.hardforkTriggerPubKey = pubKeyBytes

	fchecker, err := fieldsChecker.NewFieldsSizeChecker(instance.chainParametersHandler, hasher)
	if err != nil {
		return nil, err
	}
	instance.fieldsSizeChecker = fchecker

	instance.collectClosableComponents()

	return instance, nil
}

func computeEncodedAddressLen(converter core.PubkeyConverter) (uint32, error) {
	emptyAddress := bytes.Repeat([]byte{0}, converter.Len())
	encodedEmptyAddress, err := converter.Encode(emptyAddress)
	if err != nil {
		return 0, err
	}

	return uint32(len(encodedEmptyAddress)), nil
}

// InternalMarshalizer will return the internal marshaller
func (c *coreComponentsHolder) InternalMarshalizer() marshal.Marshalizer {
	return c.internalMarshaller
}

// SetInternalMarshalizer will set the internal marshaller
func (c *coreComponentsHolder) SetInternalMarshalizer(marshalizer marshal.Marshalizer) error {
	c.internalMarshaller = marshalizer
	return nil
}

// TxMarshalizer will return the transaction marshaller
func (c *coreComponentsHolder) TxMarshalizer() marshal.Marshalizer {
	return c.txMarshaller
}

// VmMarshalizer will return the vm marshaller
func (c *coreComponentsHolder) VmMarshalizer() marshal.Marshalizer {
	return c.vmMarshaller
}

// Hasher will return the hasher
func (c *coreComponentsHolder) Hasher() hashing.Hasher {
	return c.hasher
}

// TxSignHasher will return the transaction sign hasher
func (c *coreComponentsHolder) TxSignHasher() hashing.Hasher {
	return c.txSignHasher
}

// Uint64ByteSliceConverter will return the uint64 to slice converter
func (c *coreComponentsHolder) Uint64ByteSliceConverter() typeConverters.Uint64ByteSliceConverter {
	return c.uint64SliceConverter
}

// AddressPubKeyConverter will return the address pub key converter
func (c *coreComponentsHolder) AddressPubKeyConverter() core.PubkeyConverter {
	return c.addressPubKeyConverter
}

// ValidatorPubKeyConverter will return the validator pub key converter
func (c *coreComponentsHolder) ValidatorPubKeyConverter() core.PubkeyConverter {
	return c.validatorPubKeyConverter
}

// PathHandler will return the path handler
func (c *coreComponentsHolder) PathHandler() storage.PathManagerHandler {
	return c.pathHandler
}

// Watchdog will return the watch dog
func (c *coreComponentsHolder) Watchdog() core.WatchdogTimer {
	return c.watchdog
}

// AlarmScheduler will return the alarm scheduler
func (c *coreComponentsHolder) AlarmScheduler() core.TimersScheduler {
	return c.alarmScheduler
}

// SyncTimer will return the sync timer
func (c *coreComponentsHolder) SyncTimer() ntp.SyncTimer {
	return c.syncTimer
}

// RoundHandler will return the round handler
func (c *coreComponentsHolder) RoundHandler() consensus.RoundHandler {
	return c.roundHandler
}

// EconomicsData will return the economics data handler
func (c *coreComponentsHolder) EconomicsData() process.EconomicsDataHandler {
	return c.economicsData
}

// APIEconomicsData will return the api economics data handler
func (c *coreComponentsHolder) APIEconomicsData() process.EconomicsDataHandler {
	return c.apiEconomicsData
}

// RatingsData will return the ratings data handler
func (c *coreComponentsHolder) RatingsData() process.RatingsInfoHandler {
	return c.ratingsData
}

// Rater will return the rater handler
func (c *coreComponentsHolder) Rater() sharding.PeerAccountListAndRatingHandler {
	return c.rater
}

// GenesisNodesSetup will return the genesis nodes setup handler
func (c *coreComponentsHolder) GenesisNodesSetup() sharding.GenesisNodesSetupHandler {
	return c.genesisNodesSetup
}

// NodesShuffler will return the nodes shuffler
func (c *coreComponentsHolder) NodesShuffler() nodesCoordinator.NodesShuffler {
	return c.nodesShuffler
}

// EpochNotifier will return the epoch notifier
func (c *coreComponentsHolder) EpochNotifier() process.EpochNotifier {
	return c.epochNotifier
}

// EnableRoundsHandler will return the enable rounds handler
func (c *coreComponentsHolder) EnableRoundsHandler() process.EnableRoundsHandler {
	return c.enableRoundsHandler
}

// RoundNotifier will return the round notifier
func (c *coreComponentsHolder) RoundNotifier() process.RoundNotifier {
	return c.roundNotifier
}

// EpochStartNotifierWithConfirm will return the epoch start notifier with confirm
func (c *coreComponentsHolder) EpochStartNotifierWithConfirm() factory.EpochStartNotifierWithConfirm {
	return c.epochStartNotifierWithConfirm
}

// ChanStopNodeProcess will return the channel for stop node process
func (c *coreComponentsHolder) ChanStopNodeProcess() chan endProcess.ArgEndProcess {
	return c.chanStopNodeProcess
}

// GenesisTime will return the genesis time
func (c *coreComponentsHolder) GenesisTime() time.Time {
	return c.genesisTime
}

// ChainID will return the chain id
func (c *coreComponentsHolder) ChainID() string {
	return c.chainID
}

// MinTransactionVersion will return the min transaction version
func (c *coreComponentsHolder) MinTransactionVersion() uint32 {
	return c.minTransactionVersion
}

// TxVersionChecker will return the tx version checker
func (c *coreComponentsHolder) TxVersionChecker() process.TxVersionCheckerHandler {
	return c.txVersionChecker
}

// EncodedAddressLen will return the len of encoded address
func (c *coreComponentsHolder) EncodedAddressLen() uint32 {
	return c.encodedAddressLen
}

// NodeTypeProvider will return the node type provider
func (c *coreComponentsHolder) NodeTypeProvider() core.NodeTypeProviderHandler {
	return c.nodeTypeProvider
}

// WasmVMChangeLocker will return the wasm vm change locker
func (c *coreComponentsHolder) WasmVMChangeLocker() common.Locker {
	return c.wasmVMChangeLocker
}

// ProcessStatusHandler will return the process status handler
func (c *coreComponentsHolder) ProcessStatusHandler() common.ProcessStatusHandler {
	return c.processStatusHandler
}

// HardforkTriggerPubKey will return the pub key for the hard fork trigger
func (c *coreComponentsHolder) HardforkTriggerPubKey() []byte {
	return c.hardforkTriggerPubKey
}

// EnableEpochsHandler will return the enable epoch handler
func (c *coreComponentsHolder) EnableEpochsHandler() common.EnableEpochsHandler {
	return c.enableEpochsHandler
}

// ChainParametersSubscriber will return the chain parameters subscriber
func (c *coreComponentsHolder) ChainParametersSubscriber() process.ChainParametersSubscriber {
	return c.chainParametersSubscriber
}

// ChainParametersHandler will return the chain parameters handler
func (c *coreComponentsHolder) ChainParametersHandler() process.ChainParametersHandler {
	return c.chainParametersHandler
}

// FieldsSizeChecker will return the fields size checker component
func (c *coreComponentsHolder) FieldsSizeChecker() common.FieldsSizeChecker {
	return c.fieldsSizeChecker
}

// EpochChangeGracePeriodHandler will return the epoch change grace period handler
func (c *coreComponentsHolder) EpochChangeGracePeriodHandler() common.EpochChangeGracePeriodHandler {
	return c.epochChangeGracePeriodHandler
}

func (c *coreComponentsHolder) collectClosableComponents() {
	c.closeHandler.AddComponent(c.alarmScheduler)
	c.closeHandler.AddComponent(c.syncTimer)
}

// Close will call the Close methods on all inner components
func (c *coreComponentsHolder) Close() error {
	return c.closeHandler.Close()
}

// IsInterfaceNil returns true if there is no value under the interface
func (c *coreComponentsHolder) IsInterfaceNil() bool {
	return c == nil
}

// Create will do nothing
func (c *coreComponentsHolder) Create() error {
	return nil
}

// CheckSubcomponents will do nothing
func (c *coreComponentsHolder) CheckSubcomponents() error {
	return nil
}

// String will do nothing
func (c *coreComponentsHolder) String() string {
	return ""
}
