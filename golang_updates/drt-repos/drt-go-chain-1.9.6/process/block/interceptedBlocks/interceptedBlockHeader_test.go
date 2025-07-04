package interceptedBlocks_test

import (
	"errors"
	"math/big"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	dataBlock "github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/common/graceperiod"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/block/interceptedBlocks"
	"github.com/TerraDharitri/drt-go-chain/process/mock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/consensus"
	"github.com/TerraDharitri/drt-go-chain/testscommon/enableEpochsHandlerMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/hashingMocks"
)

var testMarshalizer = &mock.MarshalizerMock{}
var testHasher = &hashingMocks.HasherMock{}
var hdrNonce = uint64(56)
var hdrShardId = uint32(1)
var hdrRound = uint64(67)
var hdrEpoch = uint32(78)

func createDefaultShardArgument() *interceptedBlocks.ArgInterceptedBlockHeader {
	gracePeriod, _ := graceperiod.NewEpochChangeGracePeriod([]config.EpochChangeGracePeriodByEpoch{{EnableEpoch: 0, GracePeriodInRounds: 1}})
	arg := &interceptedBlocks.ArgInterceptedBlockHeader{
		ShardCoordinator:              mock.NewOneShardCoordinatorMock(),
		Hasher:                        testHasher,
		Marshalizer:                   testMarshalizer,
		HeaderSigVerifier:             &consensus.HeaderSigVerifierMock{},
		HeaderIntegrityVerifier:       &mock.HeaderIntegrityVerifierStub{},
		ValidityAttester:              &mock.ValidityAttesterStub{},
		EpochStartTrigger:             &mock.EpochStartTriggerStub{},
		EnableEpochsHandler:           &enableEpochsHandlerMock.EnableEpochsHandlerStub{},
		EpochChangeGracePeriodHandler: gracePeriod,
	}

	hdr := createMockShardHeader()
	arg.HdrBuff, _ = arg.Marshalizer.Marshal(hdr)

	return arg
}

func createDefaultShardArgumentWithV2Support() *interceptedBlocks.ArgInterceptedBlockHeader {
	gracePeriod, _ := graceperiod.NewEpochChangeGracePeriod([]config.EpochChangeGracePeriodByEpoch{{EnableEpoch: 0, GracePeriodInRounds: 1}})
	arg := &interceptedBlocks.ArgInterceptedBlockHeader{
		ShardCoordinator:              mock.NewOneShardCoordinatorMock(),
		Hasher:                        testHasher,
		Marshalizer:                   &marshal.GogoProtoMarshalizer{},
		HeaderSigVerifier:             &consensus.HeaderSigVerifierMock{},
		HeaderIntegrityVerifier:       &mock.HeaderIntegrityVerifierStub{},
		ValidityAttester:              &mock.ValidityAttesterStub{},
		EpochStartTrigger:             &mock.EpochStartTriggerStub{},
		EnableEpochsHandler:           &enableEpochsHandlerMock.EnableEpochsHandlerStub{},
		EpochChangeGracePeriodHandler: gracePeriod,
	}
	hdr := createMockShardHeader()
	arg.HdrBuff, _ = arg.Marshalizer.Marshal(hdr)

	return arg
}

func createMockShardHeader() *dataBlock.Header {
	return &dataBlock.Header{
		Nonce:            hdrNonce,
		PrevHash:         []byte("prev hash"),
		PrevRandSeed:     []byte("prev rand seed"),
		RandSeed:         []byte("rand seed"),
		PubKeysBitmap:    []byte{1},
		ShardID:          hdrShardId,
		TimeStamp:        0,
		Round:            hdrRound,
		Epoch:            hdrEpoch,
		BlockBodyType:    dataBlock.TxBlock,
		Signature:        []byte("signature"),
		MiniBlockHeaders: nil,
		PeerChanges:      nil,
		RootHash:         []byte("root hash"),
		MetaBlockHashes:  nil,
		TxCount:          0,
		ChainID:          []byte("chain ID"),
		SoftwareVersion:  []byte("version"),
		AccumulatedFees:  big.NewInt(0),
		DeveloperFees:    big.NewInt(0),
	}
}

// ------- TestNewInterceptedHeader

func TestNewInterceptedHeader_NilArgumentShouldErr(t *testing.T) {
	t.Parallel()

	inHdr, err := interceptedBlocks.NewInterceptedHeader(nil)

	assert.Nil(t, inHdr)
	assert.Equal(t, process.ErrNilArgumentStruct, err)
}

func TestNewInterceptedHeader_MarshalizerFailShouldErr(t *testing.T) {
	t.Parallel()

	arg := createDefaultShardArgument()
	arg.HdrBuff = []byte("invalid buffer")

	inHdr, err := interceptedBlocks.NewInterceptedHeader(arg)

	assert.Nil(t, inHdr)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}

func TestNewInterceptedHeader_NotForThisShardShouldWork(t *testing.T) {
	t.Parallel()

	arg := createDefaultShardArgument()
	arg.ShardCoordinator = &mock.CoordinatorStub{
		NumberOfShardsCalled: func() uint32 {
			return hdrShardId + 2
		},
		SelfIdCalled: func() uint32 {
			return hdrShardId + 1
		},
	}

	inHdr, err := interceptedBlocks.NewInterceptedHeader(arg)

	assert.False(t, check.IfNil(inHdr))
	assert.Nil(t, err)
	assert.False(t, inHdr.IsForCurrentShard())
}

func TestNewInterceptedHeader_ForThisShardShouldWork(t *testing.T) {
	t.Parallel()

	arg := createDefaultShardArgumentWithV2Support()
	arg.ShardCoordinator = &mock.CoordinatorStub{
		NumberOfShardsCalled: func() uint32 {
			return hdrShardId + 2
		},
		SelfIdCalled: func() uint32 {
			return hdrShardId
		},
	}

	inHdr, err := interceptedBlocks.NewInterceptedHeader(arg)

	assert.False(t, check.IfNil(inHdr))
	assert.Nil(t, err)
	assert.True(t, inHdr.IsForCurrentShard())
}

func TestNewInterceptedHeader_MetachainForThisShardShouldWork(t *testing.T) {
	t.Parallel()

	arg := createDefaultShardArgument()
	arg.ShardCoordinator = &mock.CoordinatorStub{
		NumberOfShardsCalled: func() uint32 {
			return hdrShardId + 2
		},
		SelfIdCalled: func() uint32 {
			return core.MetachainShardId
		},
	}

	inHdr, err := interceptedBlocks.NewInterceptedHeader(arg)

	assert.False(t, check.IfNil(inHdr))
	assert.Nil(t, err)
	assert.True(t, inHdr.IsForCurrentShard())
}

// ------- Verify

func TestInterceptedHeader_CheckValidityNilPubKeyBitmapShouldErr(t *testing.T) {
	t.Parallel()

	hdr := createMockShardHeader()
	hdr.PubKeysBitmap = nil
	buff, _ := testMarshalizer.Marshal(hdr)

	arg := createDefaultShardArgument()
	arg.HdrBuff = buff
	inHdr, _ := interceptedBlocks.NewInterceptedHeader(arg)

	err := inHdr.CheckValidity()

	assert.Equal(t, process.ErrNilPubKeysBitmap, err)
}

func TestInterceptedHeader_CheckValidityLeaderSignatureNotCorrectShouldErr(t *testing.T) {
	t.Parallel()

	arg := createDefaultShardArgumentWithV2Support()
	marshaller := arg.Marshalizer
	hdr := createMockShardHeader()
	expectedErr := errors.New("expected err")
	buff, _ := marshaller.Marshal(hdr)

	arg.HeaderSigVerifier = &consensus.HeaderSigVerifierMock{
		VerifyRandSeedAndLeaderSignatureCalled: func(header data.HeaderHandler) error {
			return expectedErr
		},
	}
	arg.EpochStartTrigger = &mock.EpochStartTriggerStub{}
	arg.HdrBuff = buff
	inHdr, _ := interceptedBlocks.NewInterceptedHeader(arg)

	err := inHdr.CheckValidity()
	assert.Equal(t, expectedErr, err)
}

func TestInterceptedHeader_CheckValidityLeaderSignatureOkShouldWork(t *testing.T) {
	t.Parallel()

	arg := createDefaultShardArgumentWithV2Support()
	marshaller := arg.Marshalizer
	hdr := createMockShardHeader()
	expectedSignature := []byte("ran")
	hdr.LeaderSignature = expectedSignature
	buff, _ := marshaller.Marshal(hdr)

	arg.HdrBuff = buff
	inHdr, err := interceptedBlocks.NewInterceptedHeader(arg)
	require.Nil(t, err)
	require.NotNil(t, inHdr)

	err = inHdr.CheckValidity()
	assert.Nil(t, err)
}

func TestInterceptedHeader_CheckValidityLeaderSignatureOkWithFlagActiveShouldWork(t *testing.T) {
	t.Parallel()

	arg := createDefaultShardArgumentWithV2Support()
	arg.EnableEpochsHandler = &enableEpochsHandlerMock.EnableEpochsHandlerStub{
		IsFlagEnabledInEpochCalled: func(flag core.EnableEpochFlag, epoch uint32) bool {
			return flag == common.AndromedaFlag
		},
	}
	wasVerifySignatureCalled := false
	arg.HeaderSigVerifier = &consensus.HeaderSigVerifierMock{
		VerifySignatureCalled: func(header data.HeaderHandler) error {
			wasVerifySignatureCalled = true

			return nil
		},
	}
	marshaller := arg.Marshalizer
	hdr := &dataBlock.HeaderV2{
		Header:                   createMockShardHeader(),
		ScheduledRootHash:        []byte("root hash"),
		ScheduledAccumulatedFees: big.NewInt(0),
		ScheduledDeveloperFees:   big.NewInt(0),
	}
	buff, _ := marshaller.Marshal(hdr)

	arg.HdrBuff = buff
	inHdr, err := interceptedBlocks.NewInterceptedHeader(arg)
	require.Nil(t, err)
	require.NotNil(t, inHdr)

	err = inHdr.CheckValidity()
	assert.Nil(t, err)
	assert.True(t, wasVerifySignatureCalled)
}

func TestInterceptedHeader_ErrorInMiniBlockShouldErr(t *testing.T) {
	t.Parallel()

	hdr := createMockShardHeader()
	badShardId := uint32(2)
	hdr.MiniBlockHeaders = []dataBlock.MiniBlockHeader{
		{
			Hash:            make([]byte, 0),
			SenderShardID:   badShardId,
			ReceiverShardID: 0,
			TxCount:         0,
			Type:            0,
		},
	}

	arg := createDefaultShardArgumentWithV2Support()
	marshaller := arg.Marshalizer
	buff, _ := marshaller.Marshal(hdr)
	arg.HdrBuff = buff
	inHdr, err := interceptedBlocks.NewInterceptedHeader(arg)
	require.Nil(t, err)
	require.NotNil(t, inHdr)

	err = inHdr.CheckValidity()

	assert.Equal(t, process.ErrInvalidShardId, err)
}

func TestInterceptedHeader_CheckValidityShouldWork(t *testing.T) {
	t.Parallel()

	arg := createDefaultShardArgumentWithV2Support()
	inHdr, err := interceptedBlocks.NewInterceptedHeader(arg)
	require.Nil(t, err)
	require.NotNil(t, inHdr)

	err = inHdr.CheckValidity()

	assert.Nil(t, err)
}

func TestInterceptedHeader_CheckAgainstRoundHandlerErrorsShouldErr(t *testing.T) {
	t.Parallel()

	arg := createDefaultShardArgumentWithV2Support()
	expectedErr := errors.New("expected error")
	arg.ValidityAttester = &mock.ValidityAttesterStub{
		CheckBlockAgainstRoundHandlerCalled: func(headerHandler data.HeaderHandler) error {
			return expectedErr
		},
	}
	inHdr, err := interceptedBlocks.NewInterceptedHeader(arg)
	require.Nil(t, err)
	require.NotNil(t, inHdr)

	err = inHdr.CheckValidity()

	assert.Equal(t, expectedErr, err)
}

func TestInterceptedHeader_CheckAgainstFinalHeaderErrorsShouldErr(t *testing.T) {
	t.Parallel()

	arg := createDefaultShardArgumentWithV2Support()
	expectedErr := errors.New("expected error")
	arg.ValidityAttester = &mock.ValidityAttesterStub{
		CheckBlockAgainstFinalCalled: func(headerHandler data.HeaderHandler) error {
			return expectedErr
		},
	}
	inHdr, err := interceptedBlocks.NewInterceptedHeader(arg)
	require.Nil(t, err)
	require.NotNil(t, inHdr)

	err = inHdr.CheckValidity()

	assert.Equal(t, expectedErr, err)
}

// ------- getters

func TestInterceptedHeader_Getters(t *testing.T) {
	t.Parallel()

	arg := createDefaultShardArgument()
	inHdr, _ := interceptedBlocks.NewInterceptedHeader(arg)

	hash := testHasher.Compute(string(arg.HdrBuff))

	assert.Equal(t, hash, inHdr.Hash())
}

// ------- IsInterfaceNil

func TestInterceptedHeader_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	var inHdr *interceptedBlocks.InterceptedHeader

	assert.True(t, check.IfNil(inHdr))
}

func Test_UnmarshalHeaderV1FromHeaderV1(t *testing.T) {
	marshalizer := &marshal.GogoProtoMarshalizer{}

	// normal HeaderV1
	hdr := createMockShardHeader()
	hdrBuff, err := hdr.Marshal()
	require.Nil(t, err)
	require.NotNil(t, hdrBuff)

	recreatedHdr, err := process.UnmarshalShardHeaderV1(marshalizer, hdrBuff)
	require.Nil(t, err)
	require.NotNil(t, recreatedHdr)
	require.Equal(t, hdr, recreatedHdr)
}

func Test_UnmarshalHeaderV1FromHeaderV2(t *testing.T) {
	marshalizer := &marshal.GogoProtoMarshalizer{}

	// normal HeaderV1
	hdr := createMockShardHeader()

	hdrV2 := &dataBlock.HeaderV2{
		Header:            hdr,
		ScheduledRootHash: nil,
	}

	hdrBuffV2, err := hdrV2.Marshal()
	require.Nil(t, err)
	require.NotNil(t, hdrBuffV2)

	recreatedHdr, err := process.UnmarshalShardHeaderV1(marshalizer, hdrBuffV2)
	require.Nil(t, recreatedHdr)
	require.NotNil(t, err)
}

func Test_UnmarshalHeaderV2FromHeaderV1(t *testing.T) {
	marshalizer := &marshal.GogoProtoMarshalizer{}

	// normal HeaderV1
	hdr := createMockShardHeader()
	hdrBuff, err := hdr.Marshal()
	require.Nil(t, err)
	require.NotNil(t, hdrBuff)

	recreatedHdr, err := process.UnmarshalShardHeaderV2(marshalizer, hdrBuff)
	require.NotNil(t, err)
	require.Nil(t, recreatedHdr)
}

func Test_UnmarshalHeaderV2FromHeaderV2(t *testing.T) {
	marshalizer := &marshal.GogoProtoMarshalizer{}

	// normal HeaderV1
	hdr := createMockShardHeader()

	hdrV2 := &dataBlock.HeaderV2{
		Header:            hdr,
		ScheduledRootHash: nil,
	}

	hdrBuffV2, err := hdrV2.Marshal()
	require.Nil(t, err)
	require.NotNil(t, hdrBuffV2)

	recreatedHdrV2, err := process.UnmarshalShardHeaderV2(marshalizer, hdrBuffV2)
	require.Nil(t, err)
	require.NotNil(t, recreatedHdrV2)
	require.Equal(t, hdrV2, recreatedHdrV2)
}

func Test_UnmarshalHeaderV2FromHeaderV1WithJson(t *testing.T) {
	marshalizer := &marshal.JsonMarshalizer{}

	// normal HeaderV1
	hdr := createMockShardHeader()
	hdrBuff, err := marshalizer.Marshal(hdr)
	require.Nil(t, err)
	require.NotNil(t, hdrBuff)

	recreatedHdr, err := process.UnmarshalShardHeaderV2(marshalizer, hdrBuff)
	require.NotNil(t, err)
	require.Nil(t, recreatedHdr)
}
