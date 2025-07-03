package consensus

import (
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	crypto "github.com/TerraDharitri/drt-go-chain-crypto"

	"github.com/TerraDharitri/drt-go-chain/consensus"
	"github.com/TerraDharitri/drt-go-chain/consensus/mock"
	"github.com/TerraDharitri/drt-go-chain/consensus/spos"
	"github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/bootstrapperStubs"
	"github.com/TerraDharitri/drt-go-chain/testscommon/cryptoMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/testscommon/enableEpochsHandlerMock"
	epochNotifierMock "github.com/TerraDharitri/drt-go-chain/testscommon/epochNotifier"
	epochstartmock "github.com/TerraDharitri/drt-go-chain/testscommon/epochstartmock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/hashingMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/shardingMocks"
)

// InitChronologyHandlerMock -
func InitChronologyHandlerMock() consensus.ChronologyHandler {
	chr := &ChronologyHandlerMock{}
	return chr
}

// InitBlockProcessorMock -
func InitBlockProcessorMock(marshaller marshal.Marshalizer) *testscommon.BlockProcessorStub {
	blockProcessorMock := &testscommon.BlockProcessorStub{}
	blockProcessorMock.CreateBlockCalled = func(header data.HeaderHandler, haveTime func() bool) (data.HeaderHandler, data.BodyHandler, error) {
		emptyBlock := &block.Body{}
		_ = header.SetRootHash([]byte{})
		return header, emptyBlock, nil
	}
	blockProcessorMock.CommitBlockCalled = func(header data.HeaderHandler, body data.BodyHandler) error {
		return nil
	}
	blockProcessorMock.RevertCurrentBlockCalled = func() {}
	blockProcessorMock.ProcessBlockCalled = func(header data.HeaderHandler, body data.BodyHandler, haveTime func() time.Duration) error {
		return nil
	}
	blockProcessorMock.DecodeBlockBodyCalled = func(dta []byte) data.BodyHandler {
		return &block.Body{}
	}
	blockProcessorMock.DecodeBlockHeaderCalled = func(dta []byte) data.HeaderHandler {
		header := &block.Header{}
		_ = marshaller.Unmarshal(header, dta)

		return header
	}
	blockProcessorMock.MarshalizedDataToBroadcastCalled = func(header data.HeaderHandler, body data.BodyHandler) (map[uint32][]byte, map[string][][]byte, error) {
		return make(map[uint32][]byte), make(map[string][][]byte), nil
	}
	blockProcessorMock.CreateNewHeaderCalled = func(round uint64, nonce uint64) (data.HeaderHandler, error) {
		return &block.Header{
			Round: round,
			Nonce: nonce,
		}, nil
	}

	return blockProcessorMock
}

// InitBlockProcessorHeaderV2Mock -
func InitBlockProcessorHeaderV2Mock() *testscommon.BlockProcessorStub {
	blockProcessorMock := &testscommon.BlockProcessorStub{}
	blockProcessorMock.CreateBlockCalled = func(header data.HeaderHandler, haveTime func() bool) (data.HeaderHandler, data.BodyHandler, error) {
		emptyBlock := &block.Body{}
		_ = header.SetRootHash([]byte{})
		return header, emptyBlock, nil
	}
	blockProcessorMock.CommitBlockCalled = func(header data.HeaderHandler, body data.BodyHandler) error {
		return nil
	}
	blockProcessorMock.RevertCurrentBlockCalled = func() {}
	blockProcessorMock.ProcessBlockCalled = func(header data.HeaderHandler, body data.BodyHandler, haveTime func() time.Duration) error {
		return nil
	}
	blockProcessorMock.DecodeBlockBodyCalled = func(dta []byte) data.BodyHandler {
		return &block.Body{}
	}
	blockProcessorMock.DecodeBlockHeaderCalled = func(dta []byte) data.HeaderHandler {
		return &block.HeaderV2{
			Header:            &block.Header{},
			ScheduledRootHash: []byte{},
		}
	}
	blockProcessorMock.MarshalizedDataToBroadcastCalled = func(header data.HeaderHandler, body data.BodyHandler) (map[uint32][]byte, map[string][][]byte, error) {
		return make(map[uint32][]byte), make(map[string][][]byte), nil
	}
	blockProcessorMock.CreateNewHeaderCalled = func(round uint64, nonce uint64) (data.HeaderHandler, error) {
		return &block.HeaderV2{
			Header: &block.Header{

				Round: round,
				Nonce: nonce,
			},
			ScheduledRootHash: []byte{},
		}, nil
	}

	return blockProcessorMock
}

// InitMultiSignerMock -
func InitMultiSignerMock() *cryptoMocks.MultisignerMock {
	multiSigner := cryptoMocks.NewMultiSigner()
	multiSigner.VerifySignatureShareCalled = func(publicKey []byte, message []byte, sig []byte) error {
		return nil
	}
	multiSigner.VerifyAggregatedSigCalled = func(pubKeysSigners [][]byte, message []byte, aggSig []byte) error {
		return nil
	}
	multiSigner.AggregateSigsCalled = func(pubKeysSigners [][]byte, signatures [][]byte) ([]byte, error) {
		return []byte("aggregatedSig"), nil
	}
	multiSigner.CreateSignatureShareCalled = func(privateKeyBytes []byte, message []byte) ([]byte, error) {
		return []byte("partialSign"), nil
	}
	return multiSigner
}

// InitKeys -
func InitKeys() (*mock.KeyGenMock, *mock.PrivateKeyMock, *mock.PublicKeyMock) {
	toByteArrayMock := func() ([]byte, error) {
		return []byte("byteArray"), nil
	}
	privKeyMock := &mock.PrivateKeyMock{
		ToByteArrayMock: toByteArrayMock,
	}
	pubKeyMock := &mock.PublicKeyMock{
		ToByteArrayMock: toByteArrayMock,
	}
	privKeyFromByteArr := func(b []byte) (crypto.PrivateKey, error) {
		return privKeyMock, nil
	}
	pubKeyFromByteArr := func(b []byte) (crypto.PublicKey, error) {
		return pubKeyMock, nil
	}
	keyGenMock := &mock.KeyGenMock{
		PrivateKeyFromByteArrayMock: privKeyFromByteArr,
		PublicKeyFromByteArrayMock:  pubKeyFromByteArr,
	}
	return keyGenMock, privKeyMock, pubKeyMock
}

// InitConsensusCoreHeaderV2 -
func InitConsensusCoreHeaderV2() *spos.ConsensusCore {
	consensusCoreMock := InitConsensusCore()
	consensusCoreMock.SetBlockProcessor(InitBlockProcessorHeaderV2Mock())

	return consensusCoreMock
}

// InitConsensusCore -
func InitConsensusCore() *spos.ConsensusCore {
	multiSignerMock := InitMultiSignerMock()

	return InitConsensusCoreWithMultiSigner(multiSignerMock)
}

// InitConsensusCoreWithMultiSigner -
func InitConsensusCoreWithMultiSigner(multiSigner crypto.MultiSigner) *spos.ConsensusCore {
	blockChain := &testscommon.ChainHandlerStub{
		GetGenesisHeaderCalled: func() data.HeaderHandler {
			return &block.Header{
				RandSeed: []byte("randSeed"),
			}
		},
	}
	marshalizerMock := mock.MarshalizerMock{}
	blockProcessorMock := InitBlockProcessorMock(marshalizerMock)
	bootstrapperMock := &bootstrapperStubs.BootstrapperStub{}
	broadcastMessengerMock := &BroadcastMessengerMock{
		BroadcastConsensusMessageCalled: func(message *consensus.Message) error {
			return nil
		},
	}

	chronologyHandlerMock := InitChronologyHandlerMock()
	hasherMock := &hashingMocks.HasherMock{}
	roundHandlerMock := &RoundHandlerMock{}
	shardCoordinatorMock := mock.ShardCoordinatorMock{}
	syncTimerMock := &SyncTimerMock{}
	nodesCoordinator := &shardingMocks.NodesCoordinatorMock{
		ComputeValidatorsGroupCalled: func(randomness []byte, round uint64, shardId uint32, epoch uint32) (nodesCoordinator.Validator, []nodesCoordinator.Validator, error) {
			defaultSelectionChances := uint32(1)
			leader := shardingMocks.NewValidatorMock([]byte("A"), 1, defaultSelectionChances)
			return leader, []nodesCoordinator.Validator{
				leader,
				shardingMocks.NewValidatorMock([]byte("B"), 1, defaultSelectionChances),
				shardingMocks.NewValidatorMock([]byte("C"), 1, defaultSelectionChances),
				shardingMocks.NewValidatorMock([]byte("D"), 1, defaultSelectionChances),
				shardingMocks.NewValidatorMock([]byte("E"), 1, defaultSelectionChances),
				shardingMocks.NewValidatorMock([]byte("F"), 1, defaultSelectionChances),
				shardingMocks.NewValidatorMock([]byte("G"), 1, defaultSelectionChances),
				shardingMocks.NewValidatorMock([]byte("H"), 1, defaultSelectionChances),
				shardingMocks.NewValidatorMock([]byte("I"), 1, defaultSelectionChances),
			}, nil
		},
	}
	epochStartSubscriber := &epochstartmock.EpochStartNotifierStub{}
	antifloodHandler := &mock.P2PAntifloodHandlerStub{}
	peerHonestyHandler := &testscommon.PeerHonestyHandlerStub{}
	headerSigVerifier := &HeaderSigVerifierMock{}
	fallbackHeaderValidator := &testscommon.FallBackHeaderValidatorStub{}
	nodeRedundancyHandler := &mock.NodeRedundancyHandlerStub{}
	scheduledProcessor := &ScheduledProcessorStub{}
	messageSigningHandler := &mock.MessageSigningHandlerStub{}
	peerBlacklistHandler := &mock.PeerBlacklistHandlerStub{}
	multiSignerContainer := cryptoMocks.NewMultiSignerContainerMock(multiSigner)
	signingHandler := &SigningHandlerStub{}
	enableEpochsHandler := &enableEpochsHandlerMock.EnableEpochsHandlerStub{}
	equivalentProofsPool := &dataRetriever.ProofsPoolMock{}
	epochNotifier := &epochNotifierMock.EpochNotifierStub{}

	container, _ := spos.NewConsensusCore(&spos.ConsensusCoreArgs{
		BlockChain:                    blockChain,
		BlockProcessor:                blockProcessorMock,
		Bootstrapper:                  bootstrapperMock,
		BroadcastMessenger:            broadcastMessengerMock,
		ChronologyHandler:             chronologyHandlerMock,
		Hasher:                        hasherMock,
		Marshalizer:                   marshalizerMock,
		MultiSignerContainer:          multiSignerContainer,
		RoundHandler:                  roundHandlerMock,
		ShardCoordinator:              shardCoordinatorMock,
		SyncTimer:                     syncTimerMock,
		NodesCoordinator:              nodesCoordinator,
		EpochStartRegistrationHandler: epochStartSubscriber,
		AntifloodHandler:              antifloodHandler,
		PeerHonestyHandler:            peerHonestyHandler,
		HeaderSigVerifier:             headerSigVerifier,
		FallbackHeaderValidator:       fallbackHeaderValidator,
		NodeRedundancyHandler:         nodeRedundancyHandler,
		ScheduledProcessor:            scheduledProcessor,
		MessageSigningHandler:         messageSigningHandler,
		PeerBlacklistHandler:          peerBlacklistHandler,
		SigningHandler:                signingHandler,
		EnableEpochsHandler:           enableEpochsHandler,
		EquivalentProofsPool:          equivalentProofsPool,
		EpochNotifier:                 epochNotifier,
		InvalidSignersCache:           &InvalidSignersCacheMock{},
	})

	return container
}
