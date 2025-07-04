package integrationTests

import (
	"fmt"
	"math/big"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/core/pubkeyConverter"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	dataBlock "github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain-core/data/endProcess"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/hashing/blake2b"
	crypto "github.com/TerraDharitri/drt-go-chain-crypto"
	mclMultiSig "github.com/TerraDharitri/drt-go-chain-crypto/signing/mcl/multisig"
	"github.com/TerraDharitri/drt-go-chain-crypto/signing/multisig"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/common/enablers"
	"github.com/TerraDharitri/drt-go-chain/common/forking"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/consensus"
	"github.com/TerraDharitri/drt-go-chain/consensus/round"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	epochStartDisabled "github.com/TerraDharitri/drt-go-chain/epochStart/bootstrap/disabled"
	"github.com/TerraDharitri/drt-go-chain/epochStart/metachain"
	"github.com/TerraDharitri/drt-go-chain/epochStart/notifier"
	"github.com/TerraDharitri/drt-go-chain/epochStart/shardchain"
	cryptoFactory "github.com/TerraDharitri/drt-go-chain/factory/crypto"
	"github.com/TerraDharitri/drt-go-chain/factory/peerSignatureHandler"
	"github.com/TerraDharitri/drt-go-chain/integrationTests/mock"
	"github.com/TerraDharitri/drt-go-chain/keysManagement"
	"github.com/TerraDharitri/drt-go-chain/node"
	"github.com/TerraDharitri/drt-go-chain/ntp"
	"github.com/TerraDharitri/drt-go-chain/p2p"
	p2pFactory "github.com/TerraDharitri/drt-go-chain/p2p/factory"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/factory"
	"github.com/TerraDharitri/drt-go-chain/process/factory/interceptorscontainer"
	"github.com/TerraDharitri/drt-go-chain/process/interceptors"
	disabledInterceptors "github.com/TerraDharitri/drt-go-chain/process/interceptors/disabled"
	interceptorsFactory "github.com/TerraDharitri/drt-go-chain/process/interceptors/factory"
	processMock "github.com/TerraDharitri/drt-go-chain/process/mock"
	"github.com/TerraDharitri/drt-go-chain/process/smartContract"
	syncFork "github.com/TerraDharitri/drt-go-chain/process/sync"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	chainShardingMocks "github.com/TerraDharitri/drt-go-chain/sharding/mock"
	"github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
	"github.com/TerraDharitri/drt-go-chain/state"
	"github.com/TerraDharitri/drt-go-chain/storage"
	"github.com/TerraDharitri/drt-go-chain/storage/cache"
	"github.com/TerraDharitri/drt-go-chain/storage/storageunit"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/chainParameters"
	consensusMocks "github.com/TerraDharitri/drt-go-chain/testscommon/consensus"
	"github.com/TerraDharitri/drt-go-chain/testscommon/cryptoMocks"
	dataRetrieverMock "github.com/TerraDharitri/drt-go-chain/testscommon/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/testscommon/economicsmocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/enableEpochsHandlerMock"
	testFactory "github.com/TerraDharitri/drt-go-chain/testscommon/factory"
	"github.com/TerraDharitri/drt-go-chain/testscommon/genesisMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/nodeTypeProviderMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/p2pmocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/shardingMocks"
	stateMock "github.com/TerraDharitri/drt-go-chain/testscommon/state"
	statusHandlerMock "github.com/TerraDharitri/drt-go-chain/testscommon/statusHandler"
	vic "github.com/TerraDharitri/drt-go-chain/testscommon/validatorInfoCacher"
)

const (
	blsConsensusType = "bls"
	signatureSize    = 48
	publicKeySize    = 96
	maxShards        = 1
)

var testPubkeyConverter, _ = pubkeyConverter.NewHexPubkeyConverter(32)

// ArgsTestConsensusNode represents the arguments for the test consensus node constructor(s)
type ArgsTestConsensusNode struct {
	ShardID            uint32
	ConsensusSize      int
	RoundTime          uint64
	ConsensusType      string
	NodeKeys           *TestNodeKeys
	EligibleMap        map[uint32][]nodesCoordinator.Validator
	WaitingMap         map[uint32][]nodesCoordinator.Validator
	KeyGen             crypto.KeyGenerator
	P2PKeyGen          crypto.KeyGenerator
	MultiSigner        *cryptoMocks.MultisignerMock
	StartTime          int64
	EnableEpochsConfig config.EnableEpochs
}

// TestConsensusNode represents a structure used in integration tests used for consensus tests
type TestConsensusNode struct {
	Node                      *node.Node
	MainMessenger             p2p.Messenger
	FullArchiveMessenger      p2p.Messenger
	NodesCoordinator          nodesCoordinator.NodesCoordinator
	ShardCoordinator          sharding.Coordinator
	ChainHandler              data.ChainHandler
	BlockProcessor            *mock.BlockProcessorMock
	RequestersFinder          dataRetriever.RequestersFinder
	AccountsDB                *state.AccountsDB
	NodeKeys                  *TestKeyPair
	MultiSigner               *cryptoMocks.MultisignerMock
	MainInterceptorsContainer process.InterceptorsContainer
	DataPool                  dataRetriever.PoolsHolder
	RequestHandler            process.RequestHandler
}

// NewTestConsensusNode returns a new TestConsensusNode
func NewTestConsensusNode(args ArgsTestConsensusNode) *TestConsensusNode {

	shardCoordinator, _ := sharding.NewMultiShardCoordinator(maxShards, args.ShardID)

	tcn := &TestConsensusNode{
		NodeKeys:         args.NodeKeys.MainKey,
		ShardCoordinator: shardCoordinator,
		MultiSigner:      args.MultiSigner,
	}
	tcn.initNode(args)

	return tcn
}

// CreateNodesWithTestConsensusNode returns a map with nodes per shard each using TestConsensusNode
func CreateNodesWithTestConsensusNode(
	numMetaNodes int,
	nodesPerShard int,
	consensusSize int,
	roundTime uint64,
	consensusType string,
	numKeysOnEachNode int,
	enableEpochsConfig config.EnableEpochs,
) map[uint32][]*TestConsensusNode {

	nodes := make(map[uint32][]*TestConsensusNode, nodesPerShard)
	cp := CreateCryptoParams(nodesPerShard, numMetaNodes, maxShards, numKeysOnEachNode)
	keysMap := PubKeysMapFromNodesKeysMap(cp.NodesKeys)
	validatorsMap := GenValidatorsFromPubKeys(keysMap, maxShards)
	eligibleMap, _ := nodesCoordinator.NodesInfoToValidators(validatorsMap)
	waitingMap := make(map[uint32][]nodesCoordinator.Validator)
	connectableNodes := make(map[uint32][]Connectable, 0)

	startTime := time.Now().Unix()
	testHasher := createHasher(consensusType)

	for shardID := range cp.NodesKeys {
		for _, keysPair := range cp.NodesKeys[shardID] {
			multiSigner, _ := multisig.NewBLSMultisig(&mclMultiSig.BlsMultiSigner{Hasher: testHasher}, cp.KeyGen)
			multiSignerMock := createCustomMultiSignerMock(multiSigner)

			args := ArgsTestConsensusNode{
				ShardID:            shardID,
				ConsensusSize:      consensusSize,
				RoundTime:          roundTime,
				ConsensusType:      consensusType,
				NodeKeys:           keysPair,
				EligibleMap:        eligibleMap,
				WaitingMap:         waitingMap,
				KeyGen:             cp.KeyGen,
				P2PKeyGen:          cp.P2PKeyGen,
				MultiSigner:        multiSignerMock,
				StartTime:          startTime,
				EnableEpochsConfig: enableEpochsConfig,
			}

			tcn := NewTestConsensusNode(args)
			nodes[shardID] = append(nodes[shardID], tcn)
			connectableNodes[shardID] = append(connectableNodes[shardID], tcn)
		}
	}

	for shardID := range nodes {
		ConnectNodes(connectableNodes[shardID])
	}

	return nodes
}

func createCustomMultiSignerMock(multiSigner crypto.MultiSigner) *cryptoMocks.MultisignerMock {
	multiSignerMock := &cryptoMocks.MultisignerMock{}
	multiSignerMock.CreateSignatureShareCalled = func(privateKeyBytes, message []byte) ([]byte, error) {
		return multiSigner.CreateSignatureShare(privateKeyBytes, message)
	}
	multiSignerMock.VerifySignatureShareCalled = func(publicKey, message, sig []byte) error {
		return multiSigner.VerifySignatureShare(publicKey, message, sig)
	}
	multiSignerMock.AggregateSigsCalled = func(pubKeysSigners, signatures [][]byte) ([]byte, error) {
		return multiSigner.AggregateSigs(pubKeysSigners, signatures)
	}
	multiSignerMock.VerifyAggregatedSigCalled = func(pubKeysSigners [][]byte, message, aggSig []byte) error {
		return multiSigner.VerifyAggregatedSig(pubKeysSigners, message, aggSig)
	}

	return multiSignerMock
}

func (tcn *TestConsensusNode) initNode(args ArgsTestConsensusNode) {
	var err error

	testHasher := createHasher(args.ConsensusType)
	epochStartRegistrationHandler := notifier.NewEpochStartSubscriptionHandler()
	consensusCache, _ := cache.NewLRUCache(10000)
	pkBytes, _ := tcn.NodeKeys.Pk.ToByteArray()

	tcn.initNodesCoordinator(args.ConsensusSize, testHasher, epochStartRegistrationHandler, args.EligibleMap, args.WaitingMap, pkBytes, consensusCache)
	tcn.MainMessenger = CreateMessengerWithNoDiscovery()
	tcn.FullArchiveMessenger = &p2pmocks.MessengerStub{}
	tcn.initBlockChain(testHasher)
	tcn.initBlockProcessor(tcn.ShardCoordinator.SelfId())

	syncer := ntp.NewSyncTime(ntp.NewNTPGoogleConfig(), nil)
	syncer.StartSyncingTime()

	genericEpochNotifier := forking.NewGenericEpochNotifier()

	epochsConfig := GetDefaultEnableEpochsConfig()
	enableEpochsHandler, _ := enablers.NewEnableEpochsHandler(*epochsConfig, genericEpochNotifier)

	storage := CreateStore(tcn.ShardCoordinator.NumberOfShards())

	roundHandler, _ := round.NewRound(
		time.Unix(args.StartTime, 0),
		syncer.CurrentTime(),
		time.Millisecond*time.Duration(args.RoundTime),
		syncer,
		0)

	dataPool := dataRetrieverMock.CreatePoolsHolder(1, 0)
	tcn.DataPool = dataPool

	var epochTrigger TestEpochStartTrigger
	if tcn.ShardCoordinator.SelfId() == core.MetachainShardId {
		argsNewMetaEpochStart := &metachain.ArgsNewMetaEpochStartTrigger{
			GenesisTime:        time.Unix(args.StartTime, 0),
			EpochStartNotifier: notifier.NewEpochStartSubscriptionHandler(),
			Settings: &config.EpochStartConfig{
				MinRoundsBetweenEpochs: 1,
				RoundsPerEpoch:         1000,
			},
			Epoch:            0,
			Storage:          createTestStore(),
			Marshalizer:      TestMarshalizer,
			Hasher:           testHasher,
			AppStatusHandler: &statusHandlerMock.AppStatusHandlerStub{},
			DataPool:         dataPool,
		}
		epochStartTrigger, err := metachain.NewEpochStartTrigger(argsNewMetaEpochStart)
		if err != nil {
			fmt.Println(err.Error())
		}
		epochTrigger = &metachain.TestTrigger{}
		epochTrigger.SetTrigger(epochStartTrigger)
	} else {
		argsPeerMiniBlocksSyncer := shardchain.ArgPeerMiniBlockSyncer{
			MiniBlocksPool:     tcn.DataPool.MiniBlocks(),
			ValidatorsInfoPool: tcn.DataPool.ValidatorsInfo(),
			RequestHandler:     &testscommon.RequestHandlerStub{},
		}
		peerMiniBlockSyncer, _ := shardchain.NewPeerMiniBlockSyncer(argsPeerMiniBlocksSyncer)

		argsShardEpochStart := &shardchain.ArgsShardEpochStartTrigger{
			Marshalizer:          TestMarshalizer,
			Hasher:               TestHasher,
			HeaderValidator:      &mock.HeaderValidatorStub{},
			Uint64Converter:      TestUint64Converter,
			DataPool:             tcn.DataPool,
			Storage:              storage,
			RequestHandler:       &testscommon.RequestHandlerStub{},
			Epoch:                0,
			Validity:             1,
			Finality:             1,
			EpochStartNotifier:   notifier.NewEpochStartSubscriptionHandler(),
			PeerMiniBlocksSyncer: peerMiniBlockSyncer,
			RoundHandler:         roundHandler,
			AppStatusHandler:     &statusHandlerMock.AppStatusHandlerStub{},
			EnableEpochsHandler:  enableEpochsHandler,
		}
		epochStartTrigger, err := shardchain.NewEpochStartTrigger(argsShardEpochStart)
		if err != nil {
			fmt.Println("NewEpochStartTrigger shard")
			fmt.Println(err.Error())
		}
		epochTrigger = &shardchain.TestTrigger{}
		epochTrigger.SetTrigger(epochStartTrigger)
	}

	tcn.initRequestersFinder()

	peerSigCache, _ := storageunit.NewCache(storageunit.CacheConfig{Type: storageunit.LRUCache, Capacity: 1000})
	peerSigHandler, _ := peerSignatureHandler.NewPeerSignatureHandler(peerSigCache, TestSingleBlsSigner, args.KeyGen)

	multiSigContainer := cryptoMocks.NewMultiSignerContainerMock(tcn.MultiSigner)
	privKey := tcn.NodeKeys.Sk
	pubKey := tcn.NodeKeys.Sk.GeneratePublic()

	tcn.initAccountsDB()

	genericEpochNotifier = forking.NewGenericEpochNotifier()
	enableEpochsHandler, _ = enablers.NewEnableEpochsHandler(args.EnableEpochsConfig, genericEpochNotifier)
	coreComponents := GetDefaultCoreComponents(enableEpochsHandler, genericEpochNotifier)
	coreComponents.SyncTimerField = syncer
	coreComponents.RoundHandlerField = roundHandler
	coreComponents.InternalMarshalizerField = TestMarshalizer
	coreComponents.HasherField = testHasher
	coreComponents.AddressPubKeyConverterField = testPubkeyConverter
	coreComponents.ChainIdCalled = func() string {
		return string(ChainID)
	}
	coreComponents.GenesisTimeField = time.Unix(args.StartTime, 0)
	coreComponents.GenesisNodesSetupField = &genesisMocks.NodesSetupStub{
		GetShardConsensusGroupSizeCalled: func() uint32 {
			return uint32(args.ConsensusSize)
		},
		GetMetaConsensusGroupSizeCalled: func() uint32 {
			return uint32(args.ConsensusSize)
		},
	}
	coreComponents.HardforkTriggerPubKeyField = []byte("provided hardfork pub key")

	argsKeysHolder := keysManagement.ArgsManagedPeersHolder{
		KeyGenerator:          args.KeyGen,
		P2PKeyGenerator:       args.P2PKeyGen,
		MaxRoundsOfInactivity: 0, // 0 for main node, non-0 for backup node
		PrefsConfig:           config.Preferences{},
		P2PKeyConverter:       p2pFactory.NewP2PKeyConverter(),
	}
	keysHolder, _ := keysManagement.NewManagedPeersHolder(argsKeysHolder)

	// adding provided handled keys
	for _, key := range args.NodeKeys.HandledKeys {
		skBytes, _ := key.Sk.ToByteArray()
		_ = keysHolder.AddManagedPeer(skBytes)
	}

	pubKeyBytes, _ := pubKey.ToByteArray()
	pubKeyString := coreComponents.ValidatorPubKeyConverterField.SilentEncode(pubKeyBytes, log)
	argsKeysHandler := keysManagement.ArgsKeysHandler{
		ManagedPeersHolder: keysHolder,
		PrivateKey:         tcn.NodeKeys.Sk,
		Pid:                tcn.MainMessenger.ID(),
	}
	keysHandler, _ := keysManagement.NewKeysHandler(argsKeysHandler)

	signingHandlerArgs := cryptoFactory.ArgsSigningHandler{
		PubKeys:              []string{pubKeyString},
		MultiSignerContainer: multiSigContainer,
		KeyGenerator:         args.KeyGen,
		KeysHandler:          keysHandler,
		SingleSigner:         TestSingleBlsSigner,
	}
	sigHandler, _ := cryptoFactory.NewSigningHandler(signingHandlerArgs)

	networkComponents := GetDefaultNetworkComponents()
	networkComponents.Messenger = tcn.MainMessenger
	networkComponents.InputAntiFlood = &mock.NilAntifloodHandler{}
	networkComponents.PeerHonesty = &mock.PeerHonestyHandlerStub{}

	cryptoComponents := GetDefaultCryptoComponents()
	cryptoComponents.PrivKey = privKey
	cryptoComponents.PubKey = pubKey
	cryptoComponents.BlockSig = TestSingleBlsSigner
	cryptoComponents.TxSig = TestSingleSigner
	cryptoComponents.MultiSigContainer = multiSigContainer
	cryptoComponents.BlKeyGen = args.KeyGen
	cryptoComponents.PeerSignHandler = peerSigHandler
	cryptoComponents.SigHandler = sigHandler
	cryptoComponents.KeysHandlerField = keysHandler

	forkDetector, _ := syncFork.NewShardForkDetector(
		roundHandler,
		cache.NewTimeCache(time.Second),
		&mock.BlockTrackerStub{},
		args.StartTime,
		enableEpochsHandler,
		dataPool.Proofs(),
	)

	processComponents := GetDefaultProcessComponents()
	processComponents.ForkDetect = forkDetector
	processComponents.ShardCoord = tcn.ShardCoordinator
	processComponents.NodesCoord = tcn.NodesCoordinator
	processComponents.BlockProcess = tcn.BlockProcessor
	processComponents.ReqFinder = tcn.RequestersFinder
	processComponents.EpochTrigger = epochTrigger
	processComponents.EpochNotifier = epochStartRegistrationHandler
	processComponents.BlackListHdl = &testscommon.TimeCacheStub{}
	processComponents.BootSore = &mock.BoostrapStorerMock{}
	processComponents.HeaderSigVerif = &consensusMocks.HeaderSigVerifierMock{}
	processComponents.HeaderIntegrVerif = &mock.HeaderIntegrityVerifierStub{}
	processComponents.ReqHandler = &testscommon.RequestHandlerStub{}
	processComponents.MainPeerMapper = mock.NewNetworkShardingCollectorMock()
	processComponents.FullArchivePeerMapper = mock.NewNetworkShardingCollectorMock()
	processComponents.RoundHandlerField = roundHandler
	processComponents.ScheduledTxsExecutionHandlerInternal = &testscommon.ScheduledTxsExecutionStub{}
	processComponents.ProcessedMiniBlocksTrackerInternal = &testscommon.ProcessedMiniBlocksTrackerStub{}
	processComponents.SentSignaturesTrackerInternal = &testscommon.SentSignatureTrackerStub{}

	tcn.initInterceptors(coreComponents, cryptoComponents, roundHandler, enableEpochsHandler, storage, epochTrigger)
	processComponents.IntContainer = tcn.MainInterceptorsContainer

	dataComponents := GetDefaultDataComponents()
	dataComponents.BlockChain = tcn.ChainHandler
	dataComponents.DataPool = dataPool
	dataComponents.Store = createTestStore()

	stateComponents := GetDefaultStateComponents()
	stateComponents.Accounts = tcn.AccountsDB
	stateComponents.AccountsAPI = tcn.AccountsDB

	statusCoreComponents := &testFactory.StatusCoreComponentsStub{
		AppStatusHandlerField: &statusHandlerMock.AppStatusHandlerStub{},
	}

	tcn.Node, err = node.NewNode(
		node.WithCoreComponents(coreComponents),
		node.WithStatusCoreComponents(statusCoreComponents),
		node.WithCryptoComponents(cryptoComponents),
		node.WithProcessComponents(processComponents),
		node.WithDataComponents(dataComponents),
		node.WithStateComponents(stateComponents),
		node.WithNetworkComponents(networkComponents),
		node.WithRoundDuration(args.RoundTime),
		node.WithConsensusType(args.ConsensusType),
		node.WithGenesisTime(time.Unix(args.StartTime, 0)),
		node.WithValidatorSignatureSize(signatureSize),
		node.WithPublicKeySize(publicKeySize),
	)

	if err != nil {
		fmt.Println(err.Error())
	}
}

func (tcn *TestConsensusNode) initInterceptors(
	coreComponents process.CoreComponentsHolder,
	cryptoComponents process.CryptoComponentsHolder,
	roundHandler consensus.RoundHandler,
	enableEpochsHandler common.EnableEpochsHandler,
	storage dataRetriever.StorageService,
	epochStartTrigger TestEpochStartTrigger,
) {
	interceptorDataVerifierArgs := interceptorsFactory.InterceptedDataVerifierFactoryArgs{
		CacheSpan:   time.Second * 10,
		CacheExpiry: time.Second * 10,
	}

	accountsAdapter := epochStartDisabled.NewAccountsAdapter()

	blockBlackListHandler := cache.NewTimeCache(TimeSpanForBadHeaders)

	genesisBlocks := make(map[uint32]data.HeaderHandler)
	blockTracker := processMock.NewBlockTrackerMock(tcn.ShardCoordinator, genesisBlocks)

	whiteLstHandler, _ := disabledInterceptors.NewDisabledWhiteListDataVerifier()

	cacherVerifiedCfg := storageunit.CacheConfig{Capacity: 5000, Type: storageunit.LRUCache, Shards: 1}
	cacheVerified, _ := storageunit.NewCache(cacherVerifiedCfg)
	whiteListerVerifiedTxs, _ := interceptors.NewWhiteListDataVerifier(cacheVerified)

	interceptorContainerFactoryArgs := interceptorscontainer.CommonInterceptorsContainerFactoryArgs{
		CoreComponents:                 coreComponents,
		CryptoComponents:               cryptoComponents,
		Accounts:                       accountsAdapter,
		ShardCoordinator:               tcn.ShardCoordinator,
		NodesCoordinator:               tcn.NodesCoordinator,
		MainMessenger:                  tcn.MainMessenger,
		FullArchiveMessenger:           tcn.FullArchiveMessenger,
		Store:                          storage,
		DataPool:                       tcn.DataPool,
		MaxTxNonceDeltaAllowed:         common.MaxTxNonceDeltaAllowed,
		TxFeeHandler:                   &economicsmocks.EconomicsHandlerMock{},
		BlockBlackList:                 blockBlackListHandler,
		HeaderSigVerifier:              &consensusMocks.HeaderSigVerifierMock{},
		HeaderIntegrityVerifier:        CreateHeaderIntegrityVerifier(),
		ValidityAttester:               blockTracker,
		EpochStartTrigger:              epochStartTrigger,
		WhiteListHandler:               whiteLstHandler,
		WhiteListerVerifiedTxs:         whiteListerVerifiedTxs,
		AntifloodHandler:               &mock.NilAntifloodHandler{},
		ArgumentsParser:                smartContract.NewArgumentParser(),
		PreferredPeersHolder:           &p2pmocks.PeersHolderStub{},
		SizeCheckDelta:                 sizeCheckDelta,
		RequestHandler:                 &testscommon.RequestHandlerStub{},
		PeerSignatureHandler:           &processMock.PeerSignatureHandlerStub{},
		SignaturesHandler:              &processMock.SignaturesHandlerStub{},
		HeartbeatExpiryTimespanInSec:   30,
		MainPeerShardMapper:            mock.NewNetworkShardingCollectorMock(),
		FullArchivePeerShardMapper:     mock.NewNetworkShardingCollectorMock(),
		HardforkTrigger:                &testscommon.HardforkTriggerStub{},
		NodeOperationMode:              common.NormalOperation,
		InterceptedDataVerifierFactory: interceptorsFactory.NewInterceptedDataVerifierFactory(interceptorDataVerifierArgs),
	}
	if tcn.ShardCoordinator.SelfId() == core.MetachainShardId {
		interceptorContainerFactory, err := interceptorscontainer.NewMetaInterceptorsContainerFactory(interceptorContainerFactoryArgs)
		if err != nil {
			fmt.Println(err.Error())
		}

		tcn.MainInterceptorsContainer, _, err = interceptorContainerFactory.Create()
		if err != nil {
			log.Debug("interceptor container factory Create", "error", err.Error())
		}
	} else {
		argsPeerMiniBlocksSyncer := shardchain.ArgPeerMiniBlockSyncer{
			MiniBlocksPool:     tcn.DataPool.MiniBlocks(),
			ValidatorsInfoPool: tcn.DataPool.ValidatorsInfo(),
			RequestHandler:     &testscommon.RequestHandlerStub{},
		}
		peerMiniBlockSyncer, _ := shardchain.NewPeerMiniBlockSyncer(argsPeerMiniBlocksSyncer)
		argsShardEpochStart := &shardchain.ArgsShardEpochStartTrigger{
			Marshalizer:          TestMarshalizer,
			Hasher:               TestHasher,
			HeaderValidator:      &mock.HeaderValidatorStub{},
			Uint64Converter:      TestUint64Converter,
			DataPool:             tcn.DataPool,
			Storage:              storage,
			RequestHandler:       &testscommon.RequestHandlerStub{},
			Epoch:                0,
			Validity:             1,
			Finality:             1,
			EpochStartNotifier:   notifier.NewEpochStartSubscriptionHandler(),
			PeerMiniBlocksSyncer: peerMiniBlockSyncer,
			RoundHandler:         roundHandler,
			AppStatusHandler:     &statusHandlerMock.AppStatusHandlerStub{},
			EnableEpochsHandler:  enableEpochsHandler,
		}
		_, _ = shardchain.NewEpochStartTrigger(argsShardEpochStart)

		interceptorContainerFactory, err := interceptorscontainer.NewShardInterceptorsContainerFactory(interceptorContainerFactoryArgs)
		if err != nil {
			fmt.Println(err.Error())
		}

		tcn.MainInterceptorsContainer, _, err = interceptorContainerFactory.Create()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func (tcn *TestConsensusNode) initNodesCoordinator(
	consensusSize int,
	hasher hashing.Hasher,
	epochStartRegistrationHandler notifier.EpochStartNotifier,
	eligibleMap map[uint32][]nodesCoordinator.Validator,
	waitingMap map[uint32][]nodesCoordinator.Validator,
	pkBytes []byte,
	cache storage.Cacher,
) {
	argumentsNodesCoordinator := nodesCoordinator.ArgNodesCoordinator{
		ChainParametersHandler: &chainParameters.ChainParametersHandlerStub{
			ChainParametersForEpochCalled: func(_ uint32) (config.ChainParametersByEpochConfig, error) {
				return config.ChainParametersByEpochConfig{
					ShardConsensusGroupSize:     uint32(consensusSize),
					MetachainConsensusGroupSize: uint32(consensusSize),
				}, nil
			},
		},
		Marshalizer:                     TestMarshalizer,
		Hasher:                          hasher,
		Shuffler:                        &shardingMocks.NodeShufflerMock{},
		EpochStartNotifier:              epochStartRegistrationHandler,
		BootStorer:                      CreateMemUnit(),
		NbShards:                        maxShards,
		EligibleNodes:                   eligibleMap,
		WaitingNodes:                    waitingMap,
		SelfPublicKey:                   pkBytes,
		ConsensusGroupCache:             cache,
		ShuffledOutHandler:              &chainShardingMocks.ShuffledOutHandlerStub{},
		ChanStopNode:                    endProcess.GetDummyEndProcessChannel(),
		NodeTypeProvider:                &nodeTypeProviderMock.NodeTypeProviderStub{},
		IsFullArchive:                   false,
		EnableEpochsHandler:             &enableEpochsHandlerMock.EnableEpochsHandlerStub{},
		ValidatorInfoCacher:             &vic.ValidatorInfoCacherStub{},
		ShardIDAsObserver:               tcn.ShardCoordinator.SelfId(),
		GenesisNodesSetupHandler:        &genesisMocks.NodesSetupStub{},
		NodesCoordinatorRegistryFactory: &shardingMocks.NodesCoordinatorRegistryFactoryMock{},
	}

	tcn.NodesCoordinator, _ = nodesCoordinator.NewIndexHashedNodesCoordinator(argumentsNodesCoordinator)
}

func (tcn *TestConsensusNode) initBlockChain(hasher hashing.Hasher) {
	if tcn.ShardCoordinator.SelfId() == core.MetachainShardId {
		tcn.ChainHandler = CreateMetaChain()
	} else {
		tcn.ChainHandler = CreateShardChain()
	}

	rootHash := []byte("roothash")
	header := &dataBlock.Header{
		Nonce:         0,
		ShardID:       tcn.ShardCoordinator.SelfId(),
		BlockBodyType: dataBlock.StateBlock,
		Signature:     rootHash,
		RootHash:      rootHash,
		PrevRandSeed:  rootHash,
		RandSeed:      rootHash,
	}

	metaHeader := &dataBlock.MetaBlock{
		Nonce:        0,
		Signature:    rootHash,
		RootHash:     rootHash,
		PrevRandSeed: rootHash,
		RandSeed:     rootHash,
	}

	if tcn.ShardCoordinator.SelfId() == core.MetachainShardId {
		_ = tcn.ChainHandler.SetGenesisHeader(metaHeader)
	} else {
		_ = tcn.ChainHandler.SetGenesisHeader(header)
	}
	hdrMarshalized, _ := TestMarshalizer.Marshal(header)
	tcn.ChainHandler.SetGenesisHeaderHash(hasher.Compute(string(hdrMarshalized)))
}

func (tcn *TestConsensusNode) initBlockProcessor(shardId uint32) {
	tcn.BlockProcessor = &mock.BlockProcessorMock{
		Marshalizer: TestMarshalizer,
		CommitBlockCalled: func(header data.HeaderHandler, body data.BodyHandler) error {
			tcn.BlockProcessor.NumCommitBlockCalled++
			headerHash, _ := core.CalculateHash(TestMarshalizer, TestHasher, header)
			tcn.ChainHandler.SetCurrentBlockHeaderHash(headerHash)
			_ = tcn.ChainHandler.SetCurrentBlockHeaderAndRootHash(header, header.GetRootHash())

			return nil
		},
		CreateBlockCalled: func(header data.HeaderHandler, haveTime func() bool) (data.HeaderHandler, data.BodyHandler, error) {
			_ = header.SetAccumulatedFees(big.NewInt(0))
			_ = header.SetDeveloperFees(big.NewInt(0))
			_ = header.SetRootHash([]byte("roothash"))

			return header, &dataBlock.Body{}, nil
		},
		MarshalizedDataToBroadcastCalled: func(header data.HeaderHandler, body data.BodyHandler) (map[uint32][]byte, map[string][][]byte, error) {
			mrsData := make(map[uint32][]byte)
			mrsTxs := make(map[string][][]byte)
			return mrsData, mrsTxs, nil
		},
		CreateNewHeaderCalled: func(round uint64, nonce uint64) (data.HeaderHandler, error) {
			if shardId == common.MetachainShardId {
				return &dataBlock.MetaBlock{
					Round:                  round,
					Nonce:                  nonce,
					SoftwareVersion:        []byte("version"),
					ValidatorStatsRootHash: []byte("validator stats root hash"),
					AccumulatedFeesInEpoch: big.NewInt(0),
					DeveloperFees:          big.NewInt(0),
					DevFeesInEpoch:         big.NewInt(0),
				}, nil
			}

			return &dataBlock.HeaderV2{
				Header: &dataBlock.Header{
					Round:           round,
					Nonce:           nonce,
					SoftwareVersion: []byte("version"),
				},
				ScheduledDeveloperFees:   big.NewInt(0),
				ScheduledAccumulatedFees: big.NewInt(0),
			}, nil
		},
		DecodeBlockHeaderCalled: func(dta []byte) data.HeaderHandler {
			var header data.HeaderHandler
			header = &dataBlock.HeaderV2{}
			if shardId == common.MetachainShardId {
				header = &dataBlock.MetaBlock{}
			}

			_ = TestMarshalizer.Unmarshal(header, dta)
			return header
		},
	}
}

func (tcn *TestConsensusNode) initRequestersFinder() {
	hdrRequester := &dataRetrieverMock.HeaderRequesterStub{}
	mbRequester := &dataRetrieverMock.HashSliceRequesterStub{}
	tcn.RequestersFinder = &dataRetrieverMock.RequestersFinderStub{
		IntraShardRequesterCalled: func(baseTopic string) (resolver dataRetriever.Requester, e error) {
			if baseTopic == factory.MiniBlocksTopic {
				return mbRequester, nil
			}
			return nil, nil
		},
		CrossShardRequesterCalled: func(baseTopic string, crossShard uint32) (resolver dataRetriever.Requester, err error) {
			if baseTopic == factory.ShardBlocksTopic {
				return hdrRequester, nil
			}
			return nil, nil
		},
	}
}

func (tcn *TestConsensusNode) initAccountsDB() {
	storer, _, err := stateMock.CreateTestingTriePruningStorer(tcn.ShardCoordinator, notifier.NewEpochStartSubscriptionHandler())
	if err != nil {
		log.Error("initAccountsDB", "error", err.Error())
	}
	trieStorage, _ := CreateTrieStorageManager(storer)

	tcn.AccountsDB, _ = CreateAccountsDB(UserAccount, trieStorage)
}

func createHasher(consensusType string) hashing.Hasher {
	if consensusType == blsConsensusType {
		hasher, _ := blake2b.NewBlake2bWithSize(mclMultiSig.HasherOutputSize)
		return hasher
	}
	return blake2b.NewBlake2b()
}

func createTestStore() dataRetriever.StorageService {
	store := dataRetriever.NewChainStorer()
	store.AddStorer(dataRetriever.TransactionUnit, CreateMemUnit())
	store.AddStorer(dataRetriever.MiniBlockUnit, CreateMemUnit())
	store.AddStorer(dataRetriever.RewardTransactionUnit, CreateMemUnit())
	store.AddStorer(dataRetriever.MetaBlockUnit, CreateMemUnit())
	store.AddStorer(dataRetriever.PeerChangesUnit, CreateMemUnit())
	store.AddStorer(dataRetriever.BlockHeaderUnit, CreateMemUnit())
	store.AddStorer(dataRetriever.BootstrapUnit, CreateMemUnit())
	store.AddStorer(dataRetriever.ReceiptsUnit, CreateMemUnit())
	store.AddStorer(dataRetriever.ScheduledSCRsUnit, CreateMemUnit())
	store.AddStorer(dataRetriever.ShardHdrNonceHashDataUnit, CreateMemUnit())
	store.AddStorer(dataRetriever.MetaHdrNonceHashDataUnit, CreateMemUnit())

	return store
}

// ConnectOnMain will try to initiate a connection to the provided parameter on the main messenger
func (tcn *TestConsensusNode) ConnectOnMain(connectable Connectable) error {
	if check.IfNil(connectable) {
		return fmt.Errorf("trying to connect to a nil Connectable parameter")
	}

	return tcn.MainMessenger.ConnectToPeer(connectable.GetMainConnectableAddress())
}

// ConnectOnFullArchive will try to initiate a connection to the provided parameter on the full archive messenger
func (tcn *TestConsensusNode) ConnectOnFullArchive(connectable Connectable) error {
	if check.IfNil(connectable) {
		return fmt.Errorf("trying to connect to a nil Connectable parameter")
	}

	return tcn.FullArchiveMessenger.ConnectToPeer(connectable.GetMainConnectableAddress())
}

// GetMainConnectableAddress returns a non circuit, non windows default connectable p2p address
func (tcn *TestConsensusNode) GetMainConnectableAddress() string {
	if tcn == nil {
		return "nil"
	}

	return GetConnectableAddress(tcn.MainMessenger)
}

// GetFullArchiveConnectableAddress returns a non circuit, non windows default connectable p2p address of the full archive network
func (tcn *TestConsensusNode) GetFullArchiveConnectableAddress() string {
	if tcn == nil {
		return "nil"
	}

	return GetConnectableAddress(tcn.FullArchiveMessenger)
}

// IsInterfaceNil returns true if there is no value under the interface
func (tcn *TestConsensusNode) IsInterfaceNil() bool {
	return tcn == nil
}
