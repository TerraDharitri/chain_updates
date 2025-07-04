package factory

import (
	"bytes"
	"errors"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/core/versioning"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	crypto "github.com/TerraDharitri/drt-go-chain-crypto"
	"github.com/stretchr/testify/assert"

	"github.com/TerraDharitri/drt-go-chain/common/graceperiod"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/block/interceptedBlocks"
	"github.com/TerraDharitri/drt-go-chain/process/mock"
	processMocks "github.com/TerraDharitri/drt-go-chain/process/mock"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/consensus"
	"github.com/TerraDharitri/drt-go-chain/testscommon/cryptoMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/economicsmocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/enableEpochsHandlerMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/epochNotifier"
	"github.com/TerraDharitri/drt-go-chain/testscommon/hashingMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/shardingMocks"
)

var errSingleSignKeyGenMock = errors.New("errSingleSignKeyGenMock")
var errSignerMockVerifySigFails = errors.New("errSignerMockVerifySigFails")
var sigOk = []byte("signature")

func createMockKeyGen() crypto.KeyGenerator {
	return &mock.SingleSignKeyGenMock{
		PublicKeyFromByteArrayCalled: func(b []byte) (key crypto.PublicKey, e error) {
			if string(b) == "" {
				return nil, errSingleSignKeyGenMock
			}

			return &mock.SingleSignPublicKey{}, nil
		},
	}
}

func createMockSigner() crypto.SingleSigner {
	return &mock.SignerMock{
		VerifyStub: func(public crypto.PublicKey, msg []byte, sig []byte) error {
			if !bytes.Equal(sig, sigOk) {
				return errSignerMockVerifySigFails
			}
			return nil
		},
	}
}

func createMockPubkeyConverter() core.PubkeyConverter {
	return testscommon.NewPubkeyConverterMock(32)
}

func createMockFeeHandler() process.FeeHandler {
	return &economicsmocks.EconomicsHandlerMock{}
}

func createMockComponentHolders() (*mock.CoreComponentsMock, *mock.CryptoComponentsMock) {
	gracePeriod, _ := graceperiod.NewEpochChangeGracePeriod([]config.EpochChangeGracePeriodByEpoch{{EnableEpoch: 0, GracePeriodInRounds: 1}})
	coreComponents := &mock.CoreComponentsMock{
		IntMarsh:            &mock.MarshalizerMock{},
		TxMarsh:             &mock.MarshalizerMock{},
		Hash:                &hashingMocks.HasherMock{},
		TxSignHasherField:   &hashingMocks.HasherMock{},
		UInt64ByteSliceConv: mock.NewNonceHashConverterMock(),
		AddrPubKeyConv:      createMockPubkeyConverter(),
		ChainIdCalled: func() string {
			return "chainID"
		},
		TxVersionCheckField:                versioning.NewTxVersionChecker(1),
		EpochNotifierField:                 &epochNotifier.EpochNotifierStub{},
		HardforkTriggerPubKeyField:         []byte("provided hardfork pub key"),
		EnableEpochsHandlerField:           &enableEpochsHandlerMock.EnableEpochsHandlerStub{},
		FieldsSizeCheckerField:             &testscommon.FieldsSizeCheckerMock{},
		EpochChangeGracePeriodHandlerField: gracePeriod,
	}
	cryptoComponents := &mock.CryptoComponentsMock{
		BlockSig:          createMockSigner(),
		TxSig:             createMockSigner(),
		MultiSigContainer: cryptoMocks.NewMultiSignerContainerMock(cryptoMocks.NewMultiSigner()),
		BlKeyGen:          createMockKeyGen(),
		TxKeyGen:          createMockKeyGen(),
	}

	return coreComponents, cryptoComponents
}

func createMockArgMetaHeaderFactoryArgument(
	coreComponents *mock.CoreComponentsMock,
	cryptoComponents *mock.CryptoComponentsMock,
) *ArgInterceptedMetaHeaderFactory {
	return &ArgInterceptedMetaHeaderFactory{
		ArgInterceptedDataFactory: ArgInterceptedDataFactory{
			CoreComponents:               coreComponents,
			CryptoComponents:             cryptoComponents,
			ShardCoordinator:             mock.NewOneShardCoordinatorMock(),
			NodesCoordinator:             shardingMocks.NewNodesCoordinatorMock(),
			FeeHandler:                   createMockFeeHandler(),
			WhiteListerVerifiedTxs:       &testscommon.WhiteListHandlerStub{},
			HeaderSigVerifier:            &consensus.HeaderSigVerifierMock{},
			ValidityAttester:             &mock.ValidityAttesterStub{},
			HeaderIntegrityVerifier:      &mock.HeaderIntegrityVerifierStub{},
			EpochStartTrigger:            &mock.EpochStartTriggerStub{},
			ArgsParser:                   &testscommon.ArgumentParserMock{},
			PeerSignatureHandler:         &processMocks.PeerSignatureHandlerStub{},
			SignaturesHandler:            &processMocks.SignaturesHandlerStub{},
			HeartbeatExpiryTimespanInSec: 30,
			PeerID:                       "pid",
		},
	}
}

func createMockArgument(
	coreComponents *mock.CoreComponentsMock,
	cryptoComponents *mock.CryptoComponentsMock,
) *ArgInterceptedDataFactory {
	return &ArgInterceptedDataFactory{
		CoreComponents:               coreComponents,
		CryptoComponents:             cryptoComponents,
		ShardCoordinator:             mock.NewOneShardCoordinatorMock(),
		NodesCoordinator:             shardingMocks.NewNodesCoordinatorMock(),
		FeeHandler:                   createMockFeeHandler(),
		WhiteListerVerifiedTxs:       &testscommon.WhiteListHandlerStub{},
		HeaderSigVerifier:            &consensus.HeaderSigVerifierMock{},
		ValidityAttester:             &mock.ValidityAttesterStub{},
		HeaderIntegrityVerifier:      &mock.HeaderIntegrityVerifierStub{},
		EpochStartTrigger:            &mock.EpochStartTriggerStub{},
		ArgsParser:                   &testscommon.ArgumentParserMock{},
		PeerSignatureHandler:         &processMocks.PeerSignatureHandlerStub{},
		SignaturesHandler:            &processMocks.SignaturesHandlerStub{},
		HeartbeatExpiryTimespanInSec: 30,
		PeerID:                       "pid",
	}
}

func TestNewInterceptedMetaHeaderDataFactory_NilArgumentShouldErr(t *testing.T) {
	t.Parallel()

	imh, err := NewInterceptedMetaHeaderDataFactory(nil)

	assert.True(t, check.IfNil(imh))
	assert.Equal(t, process.ErrNilArgumentStruct, err)
}

func TestNewInterceptedMetaHeaderDataFactory_NilMarshalizerShouldErr(t *testing.T) {
	t.Parallel()

	coreComp, cryptoComp := createMockComponentHolders()
	coreComp.IntMarsh = nil
	arg := createMockArgMetaHeaderFactoryArgument(coreComp, cryptoComp)

	imh, err := NewInterceptedMetaHeaderDataFactory(arg)
	assert.Nil(t, imh)
	assert.Equal(t, process.ErrNilMarshalizer, err)
}

func TestNewInterceptedMetaHeaderDataFactory_NilSignMarshalizerShouldErr(t *testing.T) {
	t.Parallel()

	coreComp, cryptoComp := createMockComponentHolders()
	coreComp.TxMarsh = nil
	arg := createMockArgMetaHeaderFactoryArgument(coreComp, cryptoComp)

	imh, err := NewInterceptedMetaHeaderDataFactory(arg)
	assert.True(t, check.IfNil(imh))
	assert.Equal(t, process.ErrNilMarshalizer, err)
}

func TestNewInterceptedMetaHeaderDataFactory_NilHasherShouldErr(t *testing.T) {
	t.Parallel()

	coreComp, cryptoComp := createMockComponentHolders()
	coreComp.Hash = nil
	arg := createMockArgMetaHeaderFactoryArgument(coreComp, cryptoComp)

	imh, err := NewInterceptedMetaHeaderDataFactory(arg)
	assert.True(t, check.IfNil(imh))
	assert.Equal(t, process.ErrNilHasher, err)
}

func TestNewInterceptedMetaHeaderDataFactory_NilHeaderSigVerifierShouldErr(t *testing.T) {
	t.Parallel()
	coreComp, cryptoComp := createMockComponentHolders()
	arg := createMockArgMetaHeaderFactoryArgument(coreComp, cryptoComp)
	arg.HeaderSigVerifier = nil

	imh, err := NewInterceptedMetaHeaderDataFactory(arg)
	assert.True(t, check.IfNil(imh))
	assert.Equal(t, process.ErrNilHeaderSigVerifier, err)
}

func TestNewInterceptedMetaHeaderDataFactory_NilHeaderIntegrityVerifierShouldErr(t *testing.T) {
	t.Parallel()

	coreComp, cryptoComp := createMockComponentHolders()
	arg := createMockArgMetaHeaderFactoryArgument(coreComp, cryptoComp)
	arg.HeaderIntegrityVerifier = nil

	imh, err := NewInterceptedMetaHeaderDataFactory(arg)
	assert.True(t, check.IfNil(imh))
	assert.Equal(t, process.ErrNilHeaderIntegrityVerifier, err)
}

func TestNewInterceptedMetaHeaderDataFactory_NilShardCoordinatorShouldErr(t *testing.T) {
	t.Parallel()

	coreComp, cryptoComp := createMockComponentHolders()
	arg := createMockArgMetaHeaderFactoryArgument(coreComp, cryptoComp)
	arg.ShardCoordinator = nil

	imh, err := NewInterceptedMetaHeaderDataFactory(arg)
	assert.True(t, check.IfNil(imh))
	assert.Equal(t, process.ErrNilShardCoordinator, err)
}

func TestNewInterceptedMetaHeaderDataFactory_NilChainIdShouldErr(t *testing.T) {
	t.Parallel()

	coreComp, cryptoComp := createMockComponentHolders()
	coreComp.ChainIdCalled = func() string {
		return ""
	}
	arg := createMockArgMetaHeaderFactoryArgument(coreComp, cryptoComp)

	imh, err := NewInterceptedMetaHeaderDataFactory(arg)
	assert.True(t, check.IfNil(imh))
	assert.Equal(t, process.ErrInvalidChainID, err)
}

func TestNewInterceptedMetaHeaderDataFactory_NilValidityAttesterShouldErr(t *testing.T) {
	t.Parallel()

	coreComp, cryptoComp := createMockComponentHolders()
	arg := createMockArgMetaHeaderFactoryArgument(coreComp, cryptoComp)
	arg.ValidityAttester = nil

	imh, err := NewInterceptedMetaHeaderDataFactory(arg)
	assert.True(t, check.IfNil(imh))
	assert.Equal(t, process.ErrNilValidityAttester, err)
}

func TestNewInterceptedMetaHeaderDataFactory_ShouldWorkAndCreate(t *testing.T) {
	t.Parallel()

	coreComp, cryptoComp := createMockComponentHolders()
	arg := createMockArgMetaHeaderFactoryArgument(coreComp, cryptoComp)

	imh, err := NewInterceptedMetaHeaderDataFactory(arg)
	assert.False(t, check.IfNil(imh))
	assert.Nil(t, err)
	assert.False(t, imh.IsInterfaceNil())

	marshalizer := &mock.MarshalizerMock{}
	emptyMetaHeader := &block.Header{}
	emptyMetaHeaderBuff, _ := marshalizer.Marshal(emptyMetaHeader)
	interceptedData, err := imh.Create(emptyMetaHeaderBuff, "")
	assert.Nil(t, err)

	_, ok := interceptedData.(*interceptedBlocks.InterceptedMetaHeader)
	assert.True(t, ok)
}
