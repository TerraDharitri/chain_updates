package node

import (
	"bytes"
	"encoding/hex"
	"errors"
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/data/endProcess"
	"github.com/TerraDharitri/drt-go-chain-core/data/dcdt"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	"github.com/TerraDharitri/drt-go-chain/node/mock"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/stretchr/testify/assert"
)

func TestWithInitialNodesPubKeys(t *testing.T) {
	t.Parallel()

	node, _ := NewNode()

	pubKeys := make(map[uint32][]string, 1)
	pubKeys[0] = []string{"pk1", "pk2", "pk3"}

	opt := WithInitialNodesPubKeys(pubKeys)
	err := opt(node)

	assert.Equal(t, pubKeys, node.initialNodesPubkeys)
	assert.Nil(t, err)
}

func TestWithPublicKey(t *testing.T) {
	t.Parallel()

	node, _ := NewNode()

	pubKeys := make(map[uint32][]string, 1)
	pubKeys[0] = []string{"pk1", "pk2", "pk3"}

	opt := WithInitialNodesPubKeys(pubKeys)
	err := opt(node)

	assert.Equal(t, pubKeys, node.initialNodesPubkeys)
	assert.Nil(t, err)
}

func TestWithRoundDuration_ZeroDurationShouldErr(t *testing.T) {
	t.Parallel()

	node, _ := NewNode()

	opt := WithRoundDuration(0)
	err := opt(node)

	assert.Equal(t, uint64(0), node.roundDuration)
	assert.Equal(t, ErrZeroRoundDurationNotSupported, err)
}

func TestWithRoundDuration_ShouldWork(t *testing.T) {
	t.Parallel()

	node, _ := NewNode()

	duration := uint64(5664)

	opt := WithRoundDuration(duration)
	err := opt(node)

	assert.True(t, node.roundDuration == duration)
	assert.Nil(t, err)
}

func TestWithGenesisTime(t *testing.T) {
	t.Parallel()

	node, _ := NewNode()

	aTime := time.Time{}.Add(time.Duration(uint64(78)))

	opt := WithGenesisTime(aTime)
	err := opt(node)

	assert.Equal(t, node.genesisTime, aTime)
	assert.Nil(t, err)
}

func TestWithConsensusBls_ShouldWork(t *testing.T) {
	t.Parallel()

	node, _ := NewNode()

	consensusType := "bls"
	opt := WithConsensusType(consensusType)
	err := opt(node)

	assert.Equal(t, consensusType, node.consensusType)
	assert.Nil(t, err)
}

func TestWithRequestedItemsHandler_NilRequestedItemsHandlerShouldErr(t *testing.T) {
	t.Parallel()

	node, _ := NewNode()

	opt := WithRequestedItemsHandler(nil)
	err := opt(node)

	assert.Equal(t, ErrNilRequestedItemsHandler, err)
}

func TestWithRequestedItemsHandler_OkRequestedItemsHandlerShouldWork(t *testing.T) {
	t.Parallel()

	node, _ := NewNode()

	requestedItemsHeanlder := &testscommon.TimeCacheStub{}
	opt := WithRequestedItemsHandler(requestedItemsHeanlder)
	err := opt(node)

	assert.True(t, node.requestedItemsHandler == requestedItemsHeanlder)
	assert.Nil(t, err)
}

func TestWithBootstrapRoundIndex(t *testing.T) {
	t.Parallel()

	node, _ := NewNode()
	roundIndex := uint64(0)
	opt := WithBootstrapRoundIndex(roundIndex)

	err := opt(node)
	assert.Equal(t, roundIndex, node.bootstrapRoundIndex)
	assert.Nil(t, err)
}

func TestWithPeerDenialEvaluator_NilBlackListHandlerShouldErr(t *testing.T) {
	t.Parallel()

	node, _ := NewNode()

	opt := WithPeerDenialEvaluator(nil)
	err := opt(node)

	assert.True(t, errors.Is(err, ErrNilPeerDenialEvaluator))
}

func TestWithPeerDenialEvaluator_OkHandlerShouldWork(t *testing.T) {
	t.Parallel()

	node, _ := NewNode()

	blackListHandler := &mock.PeerDenialEvaluatorStub{}
	opt := WithPeerDenialEvaluator(blackListHandler)
	err := opt(node)

	assert.True(t, node.peerDenialEvaluator == blackListHandler)
	assert.Nil(t, err)
}

func TestWithAddressSignatureSize(t *testing.T) {
	t.Parallel()

	node, _ := NewNode()
	signatureSize := 32
	opt := WithAddressSignatureSize(signatureSize)

	err := opt(node)
	assert.Equal(t, signatureSize, node.addressSignatureSize)
	assert.Nil(t, err)

	expectedHexSize := len(hex.EncodeToString(bytes.Repeat([]byte{0}, signatureSize)))
	assert.Equal(t, expectedHexSize, node.addressSignatureHexSize)
}

func TestWithValidatorSignatureSize(t *testing.T) {
	t.Parallel()

	node, _ := NewNode()
	signatureSize := 48
	opt := WithValidatorSignatureSize(signatureSize)

	err := opt(node)
	assert.Equal(t, signatureSize, node.validatorSignatureSize)
	assert.Nil(t, err)
}

func TestWithPublicKeySize(t *testing.T) {
	t.Parallel()

	node, _ := NewNode()
	publicKeySize := 96
	opt := WithPublicKeySize(publicKeySize)

	err := opt(node)
	assert.Equal(t, publicKeySize, node.publicKeySize)
	assert.Nil(t, err)
}

func TestWithNodeStopChannel_NilNodeStopChannelShouldErr(t *testing.T) {
	t.Parallel()

	node, _ := NewNode()

	opt := WithNodeStopChannel(nil)
	err := opt(node)

	assert.Equal(t, ErrNilNodeStopChannel, err)
}

func TestWithNodeStopChannel_OkNodeStopChannelShouldWork(t *testing.T) {
	t.Parallel()

	node, _ := NewNode()

	ch := make(chan endProcess.ArgEndProcess, 1)
	opt := WithNodeStopChannel(ch)
	err := opt(node)

	assert.True(t, node.chanStopNodeProcess == ch)
	assert.Nil(t, err)
}

func TestWithSignTxWithHashEpoch_EnableSignTxWithHashEpochShouldWork(t *testing.T) {
	t.Parallel()

	node, _ := NewNode()

	epochEnable := uint32(10)
	opt := WithEnableSignTxWithHashEpoch(epochEnable)
	err := opt(node)

	assert.Equal(t, epochEnable, node.enableSignTxWithHashEpoch)
	assert.Nil(t, err)
}

func TestWithDCDTNFTStorageHandler(t *testing.T) {
	t.Parallel()

	t.Run("nil dcdt nft storage, should error", func(t *testing.T) {
		t.Parallel()

		node, _ := NewNode()
		opt := WithDCDTNFTStorageHandler(nil)
		err := opt(node)

		assert.Equal(t, ErrNilDCDTNFTStorageHandler, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		dcdtStorer := &testscommon.DcdtStorageHandlerStub{
			GetDCDTNFTTokenOnDestinationCalled: func(_ vmcommon.UserAccountHandler, _ []byte, _ uint64) (*dcdt.DCDigitalToken, bool, error) {
				return nil, true, nil
			},
		}

		node, _ := NewNode()
		opt := WithDCDTNFTStorageHandler(dcdtStorer)
		err := opt(node)

		assert.NoError(t, err)
		assert.Equal(t, dcdtStorer, node.dcdtStorageHandler)
	})
}
