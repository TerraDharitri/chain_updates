package consensus

import (
	"fmt"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/core/throttler"
	"github.com/TerraDharitri/drt-go-chain-core/core/watchdog"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	logger "github.com/TerraDharitri/drt-go-chain-logger"
	"github.com/TerraDharitri/drt-go-chain-storage/timecache"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/common/disabled"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/consensus"
	"github.com/TerraDharitri/drt-go-chain/consensus/blacklist"
	"github.com/TerraDharitri/drt-go-chain/consensus/chronology"
	"github.com/TerraDharitri/drt-go-chain/consensus/spos"
	"github.com/TerraDharitri/drt-go-chain/consensus/spos/bls/proxy"
	"github.com/TerraDharitri/drt-go-chain/consensus/spos/sposFactory"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/errors"
	"github.com/TerraDharitri/drt-go-chain/factory"
	p2pFactory "github.com/TerraDharitri/drt-go-chain/p2p/factory"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/sync"
	"github.com/TerraDharitri/drt-go-chain/process/sync/storageBootstrap"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	nodesCoord "github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
	"github.com/TerraDharitri/drt-go-chain/state/syncer"
	"github.com/TerraDharitri/drt-go-chain/trie/statistics"
	"github.com/TerraDharitri/drt-go-chain/update"
)

var log = logger.GetOrCreate("factory")

const defaultSpan = 300 * time.Second

const numSignatureGoRoutinesThrottler = 30

// ConsensusComponentsFactoryArgs holds the arguments needed to create a consensus components factory
type ConsensusComponentsFactoryArgs struct {
	Config                config.Config
	FlagsConfig           config.ContextFlagsConfig
	BootstrapRoundIndex   uint64
	CoreComponents        factory.CoreComponentsHolder
	NetworkComponents     factory.NetworkComponentsHolder
	CryptoComponents      factory.CryptoComponentsHolder
	DataComponents        factory.DataComponentsHolder
	ProcessComponents     factory.ProcessComponentsHolder
	StateComponents       factory.StateComponentsHolder
	StatusComponents      factory.StatusComponentsHolder
	StatusCoreComponents  factory.StatusCoreComponentsHolder
	ScheduledProcessor    consensus.ScheduledProcessor
	IsInImportMode        bool
	ShouldDisableWatchdog bool
}

type consensusComponentsFactory struct {
	config                config.Config
	flagsConfig           config.ContextFlagsConfig
	bootstrapRoundIndex   uint64
	coreComponents        factory.CoreComponentsHolder
	networkComponents     factory.NetworkComponentsHolder
	cryptoComponents      factory.CryptoComponentsHolder
	dataComponents        factory.DataComponentsHolder
	processComponents     factory.ProcessComponentsHolder
	stateComponents       factory.StateComponentsHolder
	statusComponents      factory.StatusComponentsHolder
	statusCoreComponents  factory.StatusCoreComponentsHolder
	scheduledProcessor    consensus.ScheduledProcessor
	isInImportMode        bool
	shouldDisableWatchdog bool
}

type consensusComponents struct {
	chronology           consensus.ChronologyHandler
	bootstrapper         process.Bootstrapper
	broadcastMessenger   consensus.BroadcastMessenger
	worker               factory.ConsensusWorker
	peerBlacklistHandler consensus.PeerBlacklistHandler
	consensusTopic       string
}

// NewConsensusComponentsFactory creates an instance of consensusComponentsFactory
func NewConsensusComponentsFactory(args ConsensusComponentsFactoryArgs) (*consensusComponentsFactory, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	return &consensusComponentsFactory{
		config:                args.Config,
		flagsConfig:           args.FlagsConfig,
		bootstrapRoundIndex:   args.BootstrapRoundIndex,
		coreComponents:        args.CoreComponents,
		networkComponents:     args.NetworkComponents,
		cryptoComponents:      args.CryptoComponents,
		dataComponents:        args.DataComponents,
		processComponents:     args.ProcessComponents,
		stateComponents:       args.StateComponents,
		statusComponents:      args.StatusComponents,
		statusCoreComponents:  args.StatusCoreComponents,
		scheduledProcessor:    args.ScheduledProcessor,
		isInImportMode:        args.IsInImportMode,
		shouldDisableWatchdog: args.ShouldDisableWatchdog,
	}, nil
}

// Create will init all the components needed for a new instance of consensusComponents
func (ccf *consensusComponentsFactory) Create() (*consensusComponents, error) {
	var err error

	cc := &consensusComponents{}

	blockchain := ccf.dataComponents.Blockchain()
	notInitializedGenesisBlock := len(blockchain.GetGenesisHeaderHash()) == 0 ||
		check.IfNil(blockchain.GetGenesisHeader())
	if notInitializedGenesisBlock {
		return nil, errors.ErrGenesisBlockNotInitialized
	}

	cc.chronology, err = ccf.createChronology()
	if err != nil {
		return nil, err
	}

	cc.bootstrapper, err = ccf.createBootstrapper()
	if err != nil {
		return nil, err
	}

	err = cc.bootstrapper.StartSyncingBlocks()
	if err != nil {
		return nil, err
	}

	epoch := ccf.getEpoch()

	consensusGroupSize, err := getConsensusGroupSize(ccf.coreComponents.GenesisNodesSetup(), ccf.processComponents.ShardCoordinator(), ccf.processComponents.NodesCoordinator(), epoch)
	if err != nil {
		return nil, err
	}
	consensusState, err := ccf.createConsensusState(epoch, consensusGroupSize)
	if err != nil {
		return nil, err
	}

	consensusService, err := sposFactory.GetConsensusCoreFactory(ccf.config.Consensus.Type)
	if err != nil {
		return nil, err
	}

	cc.broadcastMessenger, err = sposFactory.GetBroadcastMessenger(
		ccf.coreComponents.InternalMarshalizer(),
		ccf.coreComponents.Hasher(),
		ccf.networkComponents.NetworkMessenger(),
		ccf.processComponents.ShardCoordinator(),
		ccf.cryptoComponents.PeerSignatureHandler(),
		ccf.dataComponents.Datapool().Headers(),
		ccf.processComponents.InterceptorsContainer(),
		ccf.coreComponents.AlarmScheduler(),
		ccf.cryptoComponents.KeysHandler(),
	)
	if err != nil {
		return nil, err
	}

	marshalizer := ccf.coreComponents.InternalMarshalizer()
	sizeCheckDelta := ccf.config.Marshalizer.SizeCheckDelta
	if sizeCheckDelta > 0 {
		marshalizer = marshal.NewSizeCheckUnmarshalizer(marshalizer, sizeCheckDelta)
	}

	cc.peerBlacklistHandler, err = ccf.createPeerBlacklistHandler()
	if err != nil {
		return nil, err
	}

	p2pSigningHandler, err := ccf.createP2pSigningHandler()
	if err != nil {
		return nil, err
	}

	argsInvalidSignersCacher := spos.ArgInvalidSignersCache{
		Hasher:         ccf.coreComponents.Hasher(),
		SigningHandler: p2pSigningHandler,
		Marshaller:     ccf.coreComponents.InternalMarshalizer(),
	}
	invalidSignersCache, err := spos.NewInvalidSignersCache(argsInvalidSignersCacher)
	if err != nil {
		return nil, err
	}

	workerArgs := &spos.WorkerArgs{
		ConsensusService:         consensusService,
		BlockChain:               ccf.dataComponents.Blockchain(),
		BlockProcessor:           ccf.processComponents.BlockProcessor(),
		ScheduledProcessor:       ccf.scheduledProcessor,
		Bootstrapper:             cc.bootstrapper,
		BroadcastMessenger:       cc.broadcastMessenger,
		ConsensusState:           consensusState,
		ForkDetector:             ccf.processComponents.ForkDetector(),
		PeerSignatureHandler:     ccf.cryptoComponents.PeerSignatureHandler(),
		Marshalizer:              marshalizer,
		Hasher:                   ccf.coreComponents.Hasher(),
		RoundHandler:             ccf.processComponents.RoundHandler(),
		ShardCoordinator:         ccf.processComponents.ShardCoordinator(),
		SyncTimer:                ccf.coreComponents.SyncTimer(),
		HeaderSigVerifier:        ccf.processComponents.HeaderSigVerifier(),
		HeaderIntegrityVerifier:  ccf.processComponents.HeaderIntegrityVerifier(),
		ChainID:                  []byte(ccf.coreComponents.ChainID()),
		NetworkShardingCollector: ccf.processComponents.PeerShardMapper(),
		AntifloodHandler:         ccf.networkComponents.InputAntiFloodHandler(),
		PoolAdder:                ccf.dataComponents.Datapool().MiniBlocks(),
		SignatureSize:            ccf.config.ValidatorPubkeyConverter.SignatureLength,
		PublicKeySize:            ccf.config.ValidatorPubkeyConverter.Length,
		AppStatusHandler:         ccf.statusCoreComponents.AppStatusHandler(),
		NodeRedundancyHandler:    ccf.processComponents.NodeRedundancyHandler(),
		PeerBlacklistHandler:     cc.peerBlacklistHandler,
		EnableEpochsHandler:      ccf.coreComponents.EnableEpochsHandler(),
		InvalidSignersCache:      invalidSignersCache,
	}

	cc.worker, err = spos.NewWorker(workerArgs)
	if err != nil {
		return nil, err
	}

	cc.worker.StartWorking()
	ccf.dataComponents.Datapool().Headers().RegisterHandler(cc.worker.ReceivedHeader)

	ccf.networkComponents.InputAntiFloodHandler().SetConsensusSizeNotifier(
		ccf.coreComponents.ChainParametersSubscriber(),
		ccf.processComponents.ShardCoordinator().SelfId(),
	)
	err = ccf.createConsensusTopic(cc)
	if err != nil {
		return nil, err
	}

	consensusArgs := &spos.ConsensusCoreArgs{
		BlockChain:                    ccf.dataComponents.Blockchain(),
		BlockProcessor:                ccf.processComponents.BlockProcessor(),
		Bootstrapper:                  cc.bootstrapper,
		BroadcastMessenger:            cc.broadcastMessenger,
		ChronologyHandler:             cc.chronology,
		Hasher:                        ccf.coreComponents.Hasher(),
		Marshalizer:                   ccf.coreComponents.InternalMarshalizer(),
		MultiSignerContainer:          ccf.cryptoComponents.MultiSignerContainer(),
		RoundHandler:                  ccf.processComponents.RoundHandler(),
		ShardCoordinator:              ccf.processComponents.ShardCoordinator(),
		NodesCoordinator:              ccf.processComponents.NodesCoordinator(),
		SyncTimer:                     ccf.coreComponents.SyncTimer(),
		EpochStartRegistrationHandler: ccf.processComponents.EpochStartNotifier(),
		AntifloodHandler:              ccf.networkComponents.InputAntiFloodHandler(),
		PeerHonestyHandler:            ccf.networkComponents.PeerHonestyHandler(),
		HeaderSigVerifier:             ccf.processComponents.HeaderSigVerifier(),
		FallbackHeaderValidator:       ccf.processComponents.FallbackHeaderValidator(),
		NodeRedundancyHandler:         ccf.processComponents.NodeRedundancyHandler(),
		ScheduledProcessor:            ccf.scheduledProcessor,
		MessageSigningHandler:         p2pSigningHandler,
		PeerBlacklistHandler:          cc.peerBlacklistHandler,
		SigningHandler:                ccf.cryptoComponents.ConsensusSigningHandler(),
		EnableEpochsHandler:           ccf.coreComponents.EnableEpochsHandler(),
		EquivalentProofsPool:          ccf.dataComponents.Datapool().Proofs(),
		EpochNotifier:                 ccf.coreComponents.EpochNotifier(),
		InvalidSignersCache:           invalidSignersCache,
	}

	consensusDataContainer, err := spos.NewConsensusCore(
		consensusArgs,
	)
	if err != nil {
		return nil, err
	}
	signatureThrottler, err := throttler.NewNumGoRoutinesThrottler(numSignatureGoRoutinesThrottler)
	if err != nil {
		return nil, err
	}

	subroundsHandlerArgs := &proxy.SubroundsHandlerArgs{
		Chronology:           cc.chronology,
		ConsensusCoreHandler: consensusDataContainer,
		ConsensusState:       consensusState,
		Worker:               cc.worker,
		SignatureThrottler:   signatureThrottler,
		AppStatusHandler:     ccf.statusCoreComponents.AppStatusHandler(),
		OutportHandler:       ccf.statusComponents.OutportHandler(),
		SentSignatureTracker: ccf.processComponents.SentSignaturesTracker(),
		EnableEpochsHandler:  ccf.coreComponents.EnableEpochsHandler(),
		ChainID:              []byte(ccf.coreComponents.ChainID()),
		CurrentPid:           ccf.networkComponents.NetworkMessenger().ID(),
	}

	subroundsHandler, err := proxy.NewSubroundsHandler(subroundsHandlerArgs)
	if err != nil {
		return nil, err
	}

	err = subroundsHandler.Start(epoch)
	if err != nil {
		return nil, err
	}

	err = ccf.addCloserInstances(cc.chronology, cc.bootstrapper, cc.worker, ccf.coreComponents.SyncTimer())
	if err != nil {
		return nil, err
	}

	return cc, nil
}

// Close will close all the inner components
func (cc *consensusComponents) Close() error {
	err := cc.chronology.Close()
	if err != nil {
		// todo: maybe just log error and try to close as much as possible
		return err
	}
	err = cc.worker.Close()
	if err != nil {
		return err
	}
	err = cc.bootstrapper.Close()
	if err != nil {
		return err
	}
	err = cc.peerBlacklistHandler.Close()
	if err != nil {
		return err
	}

	return nil
}

func (ccf *consensusComponentsFactory) createChronology() (consensus.ChronologyHandler, error) {
	wd := ccf.coreComponents.Watchdog()
	if ccf.statusComponents.OutportHandler().HasDrivers() {
		log.Warn("node is running with an outport with attached drivers. Chronology watchdog will be turned off as " +
			"it is incompatible with the indexing process.")
		wd = &watchdog.DisabledWatchdog{}
	}
	if ccf.isInImportMode {
		log.Warn("node is running in import mode. Chronology watchdog will be turned off as " +
			"it is incompatible with the import-db process.")
		wd = &watchdog.DisabledWatchdog{}
	}
	if ccf.shouldDisableWatchdog {
		log.Warn("Chronology watchdog will be turned off (explicitly).")
		wd = &watchdog.DisabledWatchdog{}
	}

	chronologyArg := chronology.ArgChronology{
		GenesisTime:      ccf.coreComponents.GenesisTime(),
		RoundHandler:     ccf.processComponents.RoundHandler(),
		SyncTimer:        ccf.coreComponents.SyncTimer(),
		Watchdog:         wd,
		AppStatusHandler: ccf.statusCoreComponents.AppStatusHandler(),
	}
	return chronology.NewChronology(chronologyArg)
}

func (ccf *consensusComponentsFactory) getEpoch() uint32 {
	blockchain := ccf.dataComponents.Blockchain()
	epoch := blockchain.GetGenesisHeader().GetEpoch()
	crtBlockHeader := blockchain.GetCurrentBlockHeader()
	if !check.IfNil(crtBlockHeader) {
		epoch = crtBlockHeader.GetEpoch()
	}
	log.Info("starting consensus", "epoch", epoch)

	return epoch
}

// createConsensusState method creates a consensusState object
func (ccf *consensusComponentsFactory) createConsensusState(epoch uint32, consensusGroupSize int) (*spos.ConsensusState, error) {
	if ccf.cryptoComponents.PublicKey() == nil {
		return nil, errors.ErrNilPublicKey
	}
	selfId, err := ccf.cryptoComponents.PublicKey().ToByteArray()
	if err != nil {
		return nil, err
	}

	if check.IfNil(ccf.processComponents.NodesCoordinator()) {
		return nil, errors.ErrNilNodesCoordinator
	}
	eligibleNodesPubKeys, err := ccf.processComponents.NodesCoordinator().GetConsensusWhitelistedNodes(epoch)
	if err != nil {
		return nil, err
	}

	roundConsensus, err := spos.NewRoundConsensus(
		eligibleNodesPubKeys,
		// TODO: move the consensus data from nodesSetup json to config
		consensusGroupSize,
		string(selfId),
		ccf.cryptoComponents.KeysHandler(),
	)
	if err != nil {
		return nil, err
	}

	roundConsensus.ResetRoundState()

	roundThreshold := spos.NewRoundThreshold()

	roundStatus := spos.NewRoundStatus()
	roundStatus.ResetRoundStatus()

	consensusState := spos.NewConsensusState(
		roundConsensus,
		roundThreshold,
		roundStatus)

	return consensusState, nil
}

func (ccf *consensusComponentsFactory) createBootstrapper() (process.Bootstrapper, error) {
	shardCoordinator := ccf.processComponents.ShardCoordinator()
	if check.IfNil(shardCoordinator) {
		return nil, errors.ErrNilShardCoordinator
	}

	if shardCoordinator.SelfId() < shardCoordinator.NumberOfShards() {
		return ccf.createShardBootstrapper()
	}

	if shardCoordinator.SelfId() == core.MetachainShardId {
		return ccf.createMetaChainBootstrapper()
	}

	return nil, sharding.ErrShardIdOutOfRange
}

func (ccf *consensusComponentsFactory) createShardBootstrapper() (process.Bootstrapper, error) {
	argsBaseStorageBootstrapper := storageBootstrap.ArgsBaseStorageBootstrapper{
		BootStorer:                   ccf.processComponents.BootStorer(),
		ForkDetector:                 ccf.processComponents.ForkDetector(),
		BlockProcessor:               ccf.processComponents.BlockProcessor(),
		ChainHandler:                 ccf.dataComponents.Blockchain(),
		Marshalizer:                  ccf.coreComponents.InternalMarshalizer(),
		Store:                        ccf.dataComponents.StorageService(),
		Uint64Converter:              ccf.coreComponents.Uint64ByteSliceConverter(),
		BootstrapRoundIndex:          ccf.bootstrapRoundIndex,
		ShardCoordinator:             ccf.processComponents.ShardCoordinator(),
		NodesCoordinator:             ccf.processComponents.NodesCoordinator(),
		EpochStartTrigger:            ccf.processComponents.EpochStartTrigger(),
		BlockTracker:                 ccf.processComponents.BlockTracker(),
		ChainID:                      ccf.coreComponents.ChainID(),
		ScheduledTxsExecutionHandler: ccf.processComponents.ScheduledTxsExecutionHandler(),
		MiniblocksProvider:           ccf.dataComponents.MiniBlocksProvider(),
		EpochNotifier:                ccf.coreComponents.EpochNotifier(),
		ProcessedMiniBlocksTracker:   ccf.processComponents.ProcessedMiniBlocksTracker(),
		AppStatusHandler:             ccf.statusCoreComponents.AppStatusHandler(),
		EnableEpochsHandler:          ccf.coreComponents.EnableEpochsHandler(),
		ProofsPool:                   ccf.dataComponents.Datapool().Proofs(),
	}

	argsShardStorageBootstrapper := storageBootstrap.ArgsShardStorageBootstrapper{
		ArgsBaseStorageBootstrapper: argsBaseStorageBootstrapper,
	}

	shardStorageBootstrapper, err := storageBootstrap.NewShardStorageBootstrapper(argsShardStorageBootstrapper)
	if err != nil {
		return nil, err
	}

	accountsDBSyncer, err := ccf.createUserAccountsSyncer()
	if err != nil {
		return nil, err
	}

	stateNodesNotifierSubscriber, ok := accountsDBSyncer.(common.StateSyncNotifierSubscriber)
	if !ok {
		return nil, fmt.Errorf("wrong type conversion for accountsDBSyncer, type: %T", accountsDBSyncer)
	}
	err = ccf.stateComponents.MissingTrieNodesNotifier().RegisterHandler(stateNodesNotifierSubscriber)
	if err != nil {
		return nil, err
	}

	argsBaseBootstrapper := sync.ArgBaseBootstrapper{
		PoolsHolder:                  ccf.dataComponents.Datapool(),
		Store:                        ccf.dataComponents.StorageService(),
		ChainHandler:                 ccf.dataComponents.Blockchain(),
		RoundHandler:                 ccf.processComponents.RoundHandler(),
		BlockProcessor:               ccf.processComponents.BlockProcessor(),
		WaitTime:                     ccf.processComponents.RoundHandler().TimeDuration(),
		Hasher:                       ccf.coreComponents.Hasher(),
		Marshalizer:                  ccf.coreComponents.InternalMarshalizer(),
		ForkDetector:                 ccf.processComponents.ForkDetector(),
		RequestHandler:               ccf.processComponents.RequestHandler(),
		ShardCoordinator:             ccf.processComponents.ShardCoordinator(),
		Accounts:                     ccf.stateComponents.AccountsAdapter(),
		BlackListHandler:             ccf.processComponents.BlackListHandler(),
		NetworkWatcher:               ccf.networkComponents.NetworkMessenger(),
		BootStorer:                   ccf.processComponents.BootStorer(),
		StorageBootstrapper:          shardStorageBootstrapper,
		EpochHandler:                 ccf.processComponents.EpochStartTrigger(),
		MiniblocksProvider:           ccf.dataComponents.MiniBlocksProvider(),
		Uint64Converter:              ccf.coreComponents.Uint64ByteSliceConverter(),
		AppStatusHandler:             ccf.statusCoreComponents.AppStatusHandler(),
		OutportHandler:               ccf.statusComponents.OutportHandler(),
		AccountsDBSyncer:             accountsDBSyncer,
		CurrentEpochProvider:         ccf.processComponents.CurrentEpochProvider(),
		IsInImportMode:               ccf.isInImportMode,
		HistoryRepo:                  ccf.processComponents.HistoryRepository(),
		ScheduledTxsExecutionHandler: ccf.processComponents.ScheduledTxsExecutionHandler(),
		ProcessWaitTime:              time.Duration(ccf.config.GeneralSettings.SyncProcessTimeInMillis) * time.Millisecond,
		RepopulateTokensSupplies:     ccf.flagsConfig.RepopulateTokensSupplies,
		EnableEpochsHandler:          ccf.coreComponents.EnableEpochsHandler(),
	}

	argsShardBootstrapper := sync.ArgShardBootstrapper{
		ArgBaseBootstrapper: argsBaseBootstrapper,
	}

	return sync.NewShardBootstrap(argsShardBootstrapper)
}

func (ccf *consensusComponentsFactory) createArgsBaseAccountsSyncer(trieStorageManager common.StorageManager) syncer.ArgsNewBaseAccountsSyncer {
	return syncer.ArgsNewBaseAccountsSyncer{
		Hasher:                            ccf.coreComponents.Hasher(),
		Marshalizer:                       ccf.coreComponents.InternalMarshalizer(),
		TrieStorageManager:                trieStorageManager,
		RequestHandler:                    ccf.processComponents.RequestHandler(),
		Timeout:                           common.TimeoutGettingTrieNodes,
		Cacher:                            ccf.dataComponents.Datapool().TrieNodes(),
		MaxTrieLevelInMemory:              ccf.config.StateTriesConfig.MaxStateTrieLevelInMemory,
		MaxHardCapForMissingNodes:         ccf.config.TrieSync.MaxHardCapForMissingNodes,
		TrieSyncerVersion:                 ccf.config.TrieSync.TrieSyncerVersion,
		CheckNodesOnDisk:                  ccf.config.TrieSync.CheckNodesOnDisk,
		UserAccountsSyncStatisticsHandler: statistics.NewTrieSyncStatistics(),
		AppStatusHandler:                  disabled.NewAppStatusHandler(),
		EnableEpochsHandler:               ccf.coreComponents.EnableEpochsHandler(),
	}
}

func (ccf *consensusComponentsFactory) createValidatorAccountsSyncer() (process.AccountsDBSyncer, error) {
	trieStorageManager, ok := ccf.stateComponents.TrieStorageManagers()[dataRetriever.PeerAccountsUnit.String()]
	if !ok {
		return nil, errors.ErrNilTrieStorageManager
	}

	args := syncer.ArgsNewValidatorAccountsSyncer{
		ArgsNewBaseAccountsSyncer: ccf.createArgsBaseAccountsSyncer(trieStorageManager),
	}
	return syncer.NewValidatorAccountsSyncer(args)
}

func (ccf *consensusComponentsFactory) createUserAccountsSyncer() (process.AccountsDBSyncer, error) {
	trieStorageManager, ok := ccf.stateComponents.TrieStorageManagers()[dataRetriever.UserAccountsUnit.String()]
	if !ok {
		return nil, errors.ErrNilTrieStorageManager
	}

	thr, err := throttler.NewNumGoRoutinesThrottler(int32(ccf.config.TrieSync.NumConcurrentTrieSyncers))
	if err != nil {
		return nil, err
	}

	argsUserAccountsSyncer := syncer.ArgsNewUserAccountsSyncer{
		ArgsNewBaseAccountsSyncer: ccf.createArgsBaseAccountsSyncer(trieStorageManager),
		ShardId:                   ccf.processComponents.ShardCoordinator().SelfId(),
		Throttler:                 thr,
		AddressPubKeyConverter:    ccf.coreComponents.AddressPubKeyConverter(),
	}
	return syncer.NewUserAccountsSyncer(argsUserAccountsSyncer)
}

func (ccf *consensusComponentsFactory) createMetaChainBootstrapper() (process.Bootstrapper, error) {
	argsBaseStorageBootstrapper := storageBootstrap.ArgsBaseStorageBootstrapper{
		BootStorer:                   ccf.processComponents.BootStorer(),
		ForkDetector:                 ccf.processComponents.ForkDetector(),
		BlockProcessor:               ccf.processComponents.BlockProcessor(),
		ChainHandler:                 ccf.dataComponents.Blockchain(),
		Marshalizer:                  ccf.coreComponents.InternalMarshalizer(),
		Store:                        ccf.dataComponents.StorageService(),
		Uint64Converter:              ccf.coreComponents.Uint64ByteSliceConverter(),
		BootstrapRoundIndex:          ccf.bootstrapRoundIndex,
		ShardCoordinator:             ccf.processComponents.ShardCoordinator(),
		NodesCoordinator:             ccf.processComponents.NodesCoordinator(),
		EpochStartTrigger:            ccf.processComponents.EpochStartTrigger(),
		BlockTracker:                 ccf.processComponents.BlockTracker(),
		ChainID:                      ccf.coreComponents.ChainID(),
		ScheduledTxsExecutionHandler: ccf.processComponents.ScheduledTxsExecutionHandler(),
		MiniblocksProvider:           ccf.dataComponents.MiniBlocksProvider(),
		EpochNotifier:                ccf.coreComponents.EpochNotifier(),
		ProcessedMiniBlocksTracker:   ccf.processComponents.ProcessedMiniBlocksTracker(),
		AppStatusHandler:             ccf.statusCoreComponents.AppStatusHandler(),
		EnableEpochsHandler:          ccf.coreComponents.EnableEpochsHandler(),
		ProofsPool:                   ccf.dataComponents.Datapool().Proofs(),
	}

	argsMetaStorageBootstrapper := storageBootstrap.ArgsMetaStorageBootstrapper{
		ArgsBaseStorageBootstrapper: argsBaseStorageBootstrapper,
		PendingMiniBlocksHandler:    ccf.processComponents.PendingMiniBlocksHandler(),
	}

	metaStorageBootstrapper, err := storageBootstrap.NewMetaStorageBootstrapper(argsMetaStorageBootstrapper)
	if err != nil {
		return nil, err
	}

	accountsDBSyncer, err := ccf.createUserAccountsSyncer()
	if err != nil {
		return nil, err
	}

	validatorAccountsDBSyncer, err := ccf.createValidatorAccountsSyncer()
	if err != nil {
		return nil, err
	}

	argsBaseBootstrapper := sync.ArgBaseBootstrapper{
		PoolsHolder:                  ccf.dataComponents.Datapool(),
		Store:                        ccf.dataComponents.StorageService(),
		ChainHandler:                 ccf.dataComponents.Blockchain(),
		RoundHandler:                 ccf.processComponents.RoundHandler(),
		BlockProcessor:               ccf.processComponents.BlockProcessor(),
		WaitTime:                     ccf.processComponents.RoundHandler().TimeDuration(),
		Hasher:                       ccf.coreComponents.Hasher(),
		Marshalizer:                  ccf.coreComponents.InternalMarshalizer(),
		ForkDetector:                 ccf.processComponents.ForkDetector(),
		RequestHandler:               ccf.processComponents.RequestHandler(),
		ShardCoordinator:             ccf.processComponents.ShardCoordinator(),
		Accounts:                     ccf.stateComponents.AccountsAdapter(),
		BlackListHandler:             ccf.processComponents.BlackListHandler(),
		NetworkWatcher:               ccf.networkComponents.NetworkMessenger(),
		BootStorer:                   ccf.processComponents.BootStorer(),
		StorageBootstrapper:          metaStorageBootstrapper,
		EpochHandler:                 ccf.processComponents.EpochStartTrigger(),
		MiniblocksProvider:           ccf.dataComponents.MiniBlocksProvider(),
		Uint64Converter:              ccf.coreComponents.Uint64ByteSliceConverter(),
		AppStatusHandler:             ccf.statusCoreComponents.AppStatusHandler(),
		OutportHandler:               ccf.statusComponents.OutportHandler(),
		AccountsDBSyncer:             accountsDBSyncer,
		CurrentEpochProvider:         ccf.processComponents.CurrentEpochProvider(),
		IsInImportMode:               ccf.isInImportMode,
		HistoryRepo:                  ccf.processComponents.HistoryRepository(),
		ScheduledTxsExecutionHandler: ccf.processComponents.ScheduledTxsExecutionHandler(),
		ProcessWaitTime:              time.Duration(ccf.config.GeneralSettings.SyncProcessTimeInMillis) * time.Millisecond,
		RepopulateTokensSupplies:     ccf.flagsConfig.RepopulateTokensSupplies,
		EnableEpochsHandler:          ccf.coreComponents.EnableEpochsHandler(),
	}

	argsMetaBootstrapper := sync.ArgMetaBootstrapper{
		ArgBaseBootstrapper:         argsBaseBootstrapper,
		EpochBootstrapper:           ccf.processComponents.EpochStartTrigger(),
		ValidatorAccountsDB:         ccf.stateComponents.PeerAccounts(),
		ValidatorStatisticsDBSyncer: validatorAccountsDBSyncer,
	}

	return sync.NewMetaBootstrap(argsMetaBootstrapper)
}

func (ccf *consensusComponentsFactory) createConsensusTopic(cc *consensusComponents) error {
	shardCoordinator := ccf.processComponents.ShardCoordinator()
	messenger := ccf.networkComponents.NetworkMessenger()

	if check.IfNil(shardCoordinator) {
		return errors.ErrNilShardCoordinator
	}
	if check.IfNil(messenger) {
		return errors.ErrNilMessenger
	}

	cc.consensusTopic = common.ConsensusTopic + shardCoordinator.CommunicationIdentifier(shardCoordinator.SelfId())
	if !ccf.networkComponents.NetworkMessenger().HasTopic(cc.consensusTopic) {
		err := ccf.networkComponents.NetworkMessenger().CreateTopic(cc.consensusTopic, true)
		if err != nil {
			return err
		}
	}

	return ccf.networkComponents.NetworkMessenger().RegisterMessageProcessor(cc.consensusTopic, common.DefaultInterceptorsIdentifier, cc.worker)
}

func (ccf *consensusComponentsFactory) createPeerBlacklistHandler() (consensus.PeerBlacklistHandler, error) {
	cache := timecache.NewTimeCache(defaultSpan)
	peerCacher, err := timecache.NewPeerTimeCache(cache)
	if err != nil {
		return nil, err
	}
	blacklistArgs := blacklist.PeerBlackListArgs{
		PeerCacher: peerCacher,
	}

	return blacklist.NewPeerBlacklist(blacklistArgs)
}

func (ccf *consensusComponentsFactory) createP2pSigningHandler() (consensus.P2PSigningHandler, error) {
	p2pSignerArgs := p2pFactory.ArgsMessageVerifier{
		Marshaller: ccf.coreComponents.InternalMarshalizer(),
		P2PSigner:  ccf.networkComponents.NetworkMessenger(),
		Logger:     logger.GetOrCreate("main/p2p/messagecheck"),
	}

	return p2pFactory.NewMessageVerifier(p2pSignerArgs)
}

func (ccf *consensusComponentsFactory) addCloserInstances(closers ...update.Closer) error {
	hardforkTrigger := ccf.processComponents.HardforkTrigger()
	for _, c := range closers {
		err := hardforkTrigger.AddCloser(c)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkArgs(args ConsensusComponentsFactoryArgs) error {
	if check.IfNil(args.CoreComponents) {
		return errors.ErrNilCoreComponentsHolder
	}
	if check.IfNil(args.CoreComponents.GenesisNodesSetup()) {
		return errors.ErrNilGenesisNodesSetupHandler
	}
	if check.IfNil(args.DataComponents) {
		return errors.ErrNilDataComponentsHolder
	}
	if check.IfNil(args.DataComponents.Datapool()) {
		return errors.ErrNilDataPoolsHolder
	}
	if check.IfNil(args.DataComponents.Blockchain()) {
		return errors.ErrNilBlockChainHandler
	}
	if check.IfNil(args.CryptoComponents) {
		return errors.ErrNilCryptoComponentsHolder
	}
	if check.IfNil(args.CryptoComponents.PublicKey()) {
		return errors.ErrNilPublicKey
	}
	if check.IfNil(args.CryptoComponents.PrivateKey()) {
		return errors.ErrNilPrivateKey
	}
	if check.IfNil(args.NetworkComponents) {
		return errors.ErrNilNetworkComponentsHolder
	}
	if check.IfNil(args.NetworkComponents.NetworkMessenger()) {
		return errors.ErrNilMessenger
	}
	if check.IfNil(args.ProcessComponents) {
		return errors.ErrNilProcessComponentsHolder
	}
	if check.IfNil(args.ProcessComponents.NodesCoordinator()) {
		return errors.ErrNilNodesCoordinator
	}
	if check.IfNil(args.ProcessComponents.ShardCoordinator()) {
		return errors.ErrNilShardCoordinator
	}
	if check.IfNil(args.ProcessComponents.RoundHandler()) {
		return errors.ErrNilRoundHandler
	}
	if check.IfNil(args.ProcessComponents.HardforkTrigger()) {
		return errors.ErrNilHardforkTrigger
	}
	if check.IfNil(args.StateComponents) {
		return errors.ErrNilStateComponentsHolder
	}
	if check.IfNil(args.StatusComponents) {
		return errors.ErrNilStatusComponentsHolder
	}
	if check.IfNil(args.StatusComponents.OutportHandler()) {
		return errors.ErrNilOutportHandler
	}
	if check.IfNil(args.ScheduledProcessor) {
		return errors.ErrNilScheduledProcessor
	}
	if check.IfNil(args.StatusCoreComponents) {
		return errors.ErrNilStatusCoreComponents
	}

	return nil
}

func getConsensusGroupSize(nodesConfig sharding.GenesisNodesSetupHandler, shardCoordinator sharding.Coordinator, nodesCoordinator nodesCoord.NodesCoordinator, epoch uint32) (int, error) {
	consensusGroupSize := nodesCoordinator.ConsensusGroupSizeForShardAndEpoch(shardCoordinator.SelfId(), epoch)
	if consensusGroupSize > 0 {
		return consensusGroupSize, nil
	}

	if shardCoordinator.SelfId() == core.MetachainShardId {
		return int(nodesConfig.GetMetaConsensusGroupSize()), nil
	}
	if shardCoordinator.SelfId() < shardCoordinator.NumberOfShards() {
		return int(nodesConfig.GetShardConsensusGroupSize()), nil
	}

	return 0, sharding.ErrShardIdOutOfRange
}
