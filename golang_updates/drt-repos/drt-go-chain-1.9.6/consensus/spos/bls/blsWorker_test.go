package bls_test

import (
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/stretchr/testify/assert"

	"github.com/TerraDharitri/drt-go-chain/consensus"
	"github.com/TerraDharitri/drt-go-chain/consensus/spos"
	"github.com/TerraDharitri/drt-go-chain/consensus/spos/bls"
	"github.com/TerraDharitri/drt-go-chain/testscommon/consensus/initializers"
)

func TestWorker_NewConsensusServiceShouldWork(t *testing.T) {
	t.Parallel()

	service, err := bls.NewConsensusService()
	assert.Nil(t, err)
	assert.False(t, check.IfNil(service))
}

func TestWorker_InitReceivedMessagesShouldWork(t *testing.T) {
	t.Parallel()

	bnService, _ := bls.NewConsensusService()
	messages := bnService.InitReceivedMessages()

	receivedMessages := make(map[consensus.MessageType][]*consensus.Message)
	receivedMessages[bls.MtBlockBodyAndHeader] = make([]*consensus.Message, 0)
	receivedMessages[bls.MtBlockBody] = make([]*consensus.Message, 0)
	receivedMessages[bls.MtBlockHeader] = make([]*consensus.Message, 0)
	receivedMessages[bls.MtSignature] = make([]*consensus.Message, 0)
	receivedMessages[bls.MtBlockHeaderFinalInfo] = make([]*consensus.Message, 0)
	receivedMessages[bls.MtInvalidSigners] = make([]*consensus.Message, 0)

	assert.Equal(t, len(receivedMessages), len(messages))
	assert.NotNil(t, messages[bls.MtBlockBodyAndHeader])
	assert.NotNil(t, messages[bls.MtBlockBody])
	assert.NotNil(t, messages[bls.MtBlockHeader])
	assert.NotNil(t, messages[bls.MtSignature])
	assert.NotNil(t, messages[bls.MtBlockHeaderFinalInfo])
	assert.NotNil(t, messages[bls.MtInvalidSigners])
}

func TestWorker_GetMessageRangeShouldWork(t *testing.T) {
	t.Parallel()

	v := make([]consensus.MessageType, 0)
	blsService, _ := bls.NewConsensusService()

	messagesRange := blsService.GetMessageRange()
	assert.NotNil(t, messagesRange)

	for i := bls.MtBlockBodyAndHeader; i <= bls.MtInvalidSigners; i++ {
		v = append(v, i)
	}
	assert.NotNil(t, v)

	for i, val := range messagesRange {
		assert.Equal(t, v[i], val)
	}
}

func TestWorker_CanProceedWithSrStartRoundFinishedForMtBlockBodyAndHeaderShouldWork(t *testing.T) {
	t.Parallel()

	blsService, _ := bls.NewConsensusService()

	consensusState := initializers.InitConsensusState()
	consensusState.SetStatus(bls.SrStartRound, spos.SsFinished)

	canProceed := blsService.CanProceed(consensusState, bls.MtBlockBodyAndHeader)
	assert.True(t, canProceed)
}

func TestWorker_CanProceedWithSrStartRoundNotFinishedForMtBlockBodyAndHeaderShouldNotWork(t *testing.T) {
	t.Parallel()

	blsService, _ := bls.NewConsensusService()

	consensusState := initializers.InitConsensusState()
	consensusState.SetStatus(bls.SrStartRound, spos.SsNotFinished)

	canProceed := blsService.CanProceed(consensusState, bls.MtBlockBodyAndHeader)
	assert.False(t, canProceed)
}

func TestWorker_CanProceedWithSrStartRoundFinishedForMtBlockBodyShouldWork(t *testing.T) {
	t.Parallel()

	blsService, _ := bls.NewConsensusService()

	consensusState := initializers.InitConsensusState()
	consensusState.SetStatus(bls.SrStartRound, spos.SsFinished)

	canProceed := blsService.CanProceed(consensusState, bls.MtBlockBody)
	assert.True(t, canProceed)
}

func TestWorker_CanProceedWithSrStartRoundNotFinishedForMtBlockBodyShouldNotWork(t *testing.T) {
	t.Parallel()

	blsService, _ := bls.NewConsensusService()

	consensusState := initializers.InitConsensusState()
	consensusState.SetStatus(bls.SrStartRound, spos.SsNotFinished)

	canProceed := blsService.CanProceed(consensusState, bls.MtBlockBody)
	assert.False(t, canProceed)
}

func TestWorker_CanProceedWithSrStartRoundFinishedForMtBlockHeaderShouldWork(t *testing.T) {
	t.Parallel()

	blsService, _ := bls.NewConsensusService()

	consensusState := initializers.InitConsensusState()
	consensusState.SetStatus(bls.SrStartRound, spos.SsFinished)

	canProceed := blsService.CanProceed(consensusState, bls.MtBlockHeader)
	assert.True(t, canProceed)
}

func TestWorker_CanProceedWithSrStartRoundNotFinishedForMtBlockHeaderShouldNotWork(t *testing.T) {
	t.Parallel()

	blsService, _ := bls.NewConsensusService()

	consensusState := initializers.InitConsensusState()
	consensusState.SetStatus(bls.SrStartRound, spos.SsNotFinished)

	canProceed := blsService.CanProceed(consensusState, bls.MtBlockHeader)
	assert.False(t, canProceed)
}

func TestWorker_CanProceedWithSrBlockFinishedForMtBlockHeaderShouldWork(t *testing.T) {
	t.Parallel()

	blsService, _ := bls.NewConsensusService()

	consensusState := initializers.InitConsensusState()
	consensusState.SetStatus(bls.SrBlock, spos.SsFinished)

	canProceed := blsService.CanProceed(consensusState, bls.MtSignature)
	assert.True(t, canProceed)
}

func TestWorker_CanProceedWithSrBlockRoundNotFinishedForMtBlockHeaderShouldNotWork(t *testing.T) {
	t.Parallel()

	blsService, _ := bls.NewConsensusService()

	consensusState := initializers.InitConsensusState()
	consensusState.SetStatus(bls.SrBlock, spos.SsNotFinished)

	canProceed := blsService.CanProceed(consensusState, bls.MtSignature)
	assert.False(t, canProceed)
}

func TestWorker_CanProceedWithSrSignatureFinishedForMtBlockHeaderFinalInfoShouldWork(t *testing.T) {
	t.Parallel()

	blsService, _ := bls.NewConsensusService()

	consensusState := initializers.InitConsensusState()
	consensusState.SetStatus(bls.SrSignature, spos.SsFinished)

	canProceed := blsService.CanProceed(consensusState, bls.MtBlockHeaderFinalInfo)
	assert.True(t, canProceed)
}

func TestWorker_CanProceedWithSrSignatureRoundNotFinishedForMtBlockHeaderFinalInfoShouldNotWork(t *testing.T) {
	t.Parallel()

	blsService, _ := bls.NewConsensusService()

	consensusState := initializers.InitConsensusState()
	consensusState.SetStatus(bls.SrSignature, spos.SsNotFinished)

	canProceed := blsService.CanProceed(consensusState, bls.MtBlockHeaderFinalInfo)
	assert.False(t, canProceed)
}

func TestWorker_CanProceedWitUnkownMessageTypeShouldNotWork(t *testing.T) {
	t.Parallel()

	blsService, _ := bls.NewConsensusService()
	consensusState := initializers.InitConsensusState()

	canProceed := blsService.CanProceed(consensusState, -1)
	assert.False(t, canProceed)
}

func TestWorker_GetSubroundName(t *testing.T) {
	t.Parallel()

	service, _ := bls.NewConsensusService()

	r := service.GetSubroundName(bls.SrStartRound)
	assert.Equal(t, "(START_ROUND)", r)
	r = service.GetSubroundName(bls.SrBlock)
	assert.Equal(t, "(BLOCK)", r)
	r = service.GetSubroundName(bls.SrSignature)
	assert.Equal(t, "(SIGNATURE)", r)
	r = service.GetSubroundName(bls.SrEndRound)
	assert.Equal(t, "(END_ROUND)", r)
	r = service.GetSubroundName(-1)
	assert.Equal(t, "Undefined subround", r)
}

func TestWorker_GetStringValue(t *testing.T) {
	t.Parallel()

	service, _ := bls.NewConsensusService()

	r := service.GetStringValue(bls.MtBlockBodyAndHeader)
	assert.Equal(t, bls.BlockBodyAndHeaderStringValue, r)
	r = service.GetStringValue(bls.MtBlockBody)
	assert.Equal(t, bls.BlockBodyStringValue, r)
	r = service.GetStringValue(bls.MtBlockHeader)
	assert.Equal(t, bls.BlockHeaderStringValue, r)
	r = service.GetStringValue(bls.MtSignature)
	assert.Equal(t, bls.BlockSignatureStringValue, r)
	r = service.GetStringValue(bls.MtBlockHeaderFinalInfo)
	assert.Equal(t, bls.BlockHeaderFinalInfoStringValue, r)
	r = service.GetStringValue(bls.MtUnknown)
	assert.Equal(t, bls.BlockUnknownStringValue, r)
	r = service.GetStringValue(-1)
	assert.Equal(t, bls.BlockDefaultStringValue, r)
}

func TestWorker_IsMessageWithBlockBodyAndHeader(t *testing.T) {
	t.Parallel()

	service, _ := bls.NewConsensusService()

	ret := service.IsMessageWithBlockBodyAndHeader(bls.MtBlockBody)
	assert.False(t, ret)

	ret = service.IsMessageWithBlockBodyAndHeader(bls.MtBlockHeader)
	assert.False(t, ret)

	ret = service.IsMessageWithBlockBodyAndHeader(bls.MtBlockBodyAndHeader)
	assert.True(t, ret)
}

func TestWorker_IsMessageWithBlockBody(t *testing.T) {
	t.Parallel()

	service, _ := bls.NewConsensusService()

	ret := service.IsMessageWithBlockBody(bls.MtBlockHeader)
	assert.False(t, ret)

	ret = service.IsMessageWithBlockBody(bls.MtBlockBody)
	assert.True(t, ret)
}

func TestWorker_IsMessageWithBlockHeader(t *testing.T) {
	t.Parallel()

	service, _ := bls.NewConsensusService()

	ret := service.IsMessageWithBlockHeader(bls.MtBlockBody)
	assert.False(t, ret)

	ret = service.IsMessageWithBlockHeader(bls.MtBlockHeader)
	assert.True(t, ret)
}

func TestWorker_IsMessageWithSignature(t *testing.T) {
	t.Parallel()

	service, _ := bls.NewConsensusService()

	ret := service.IsMessageWithSignature(bls.MtBlockBodyAndHeader)
	assert.False(t, ret)

	ret = service.IsMessageWithSignature(bls.MtSignature)
	assert.True(t, ret)
}

func TestWorker_IsMessageWithFinalInfo(t *testing.T) {
	t.Parallel()

	service, _ := bls.NewConsensusService()

	ret := service.IsMessageWithFinalInfo(bls.MtSignature)
	assert.False(t, ret)

	ret = service.IsMessageWithFinalInfo(bls.MtBlockHeaderFinalInfo)
	assert.True(t, ret)
}

func TestWorker_IsMessageWithInvalidSigners(t *testing.T) {
	t.Parallel()

	service, _ := bls.NewConsensusService()

	ret := service.IsMessageWithInvalidSigners(bls.MtBlockHeaderFinalInfo)
	assert.False(t, ret)

	ret = service.IsMessageWithInvalidSigners(bls.MtInvalidSigners)
	assert.True(t, ret)
}

func TestWorker_IsSubroundSignature(t *testing.T) {
	t.Parallel()

	service, _ := bls.NewConsensusService()

	ret := service.IsSubroundSignature(bls.SrEndRound)
	assert.False(t, ret)

	ret = service.IsSubroundSignature(bls.SrSignature)
	assert.True(t, ret)
}

func TestWorker_IsSubroundStartRound(t *testing.T) {
	t.Parallel()

	service, _ := bls.NewConsensusService()

	ret := service.IsSubroundStartRound(bls.SrSignature)
	assert.False(t, ret)

	ret = service.IsSubroundStartRound(bls.SrStartRound)
	assert.True(t, ret)
}

func TestWorker_IsMessageTypeValid(t *testing.T) {
	t.Parallel()

	service, _ := bls.NewConsensusService()

	ret := service.IsMessageTypeValid(bls.MtBlockBody)
	assert.True(t, ret)

	ret = service.IsMessageTypeValid(666)
	assert.False(t, ret)
}

func TestWorker_GetMaxNumOfMessageTypeAccepted(t *testing.T) {
	t.Parallel()

	service, _ := bls.NewConsensusService()
	t.Run("message type signature", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, bls.MaxNumOfMessageTypeSignatureAccepted, service.GetMaxNumOfMessageTypeAccepted(bls.MtSignature))
	})
	t.Run("other message types", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, bls.DefaultMaxNumOfMessageTypeAccepted, service.GetMaxNumOfMessageTypeAccepted(bls.MtUnknown))
		assert.Equal(t, bls.DefaultMaxNumOfMessageTypeAccepted, service.GetMaxNumOfMessageTypeAccepted(bls.MtBlockBody))
		assert.Equal(t, bls.DefaultMaxNumOfMessageTypeAccepted, service.GetMaxNumOfMessageTypeAccepted(bls.MtBlockHeader))
		assert.Equal(t, bls.DefaultMaxNumOfMessageTypeAccepted, service.GetMaxNumOfMessageTypeAccepted(bls.MtBlockBodyAndHeader))
		assert.Equal(t, bls.DefaultMaxNumOfMessageTypeAccepted, service.GetMaxNumOfMessageTypeAccepted(bls.MtBlockHeaderFinalInfo))
	})
}
