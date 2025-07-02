package v2_test

import (
	"context"
	"sync"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	crypto "github.com/TerraDharitri/drt-go-chain-crypto"
	"github.com/TerraDharitri/drt-go-chain-crypto/signing"
	"github.com/TerraDharitri/drt-go-chain-crypto/signing/mcl"
	mclMultiSig "github.com/TerraDharitri/drt-go-chain-crypto/signing/mcl/multisig"
	"github.com/TerraDharitri/drt-go-chain-crypto/signing/multisig"
	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/consensus/spos"
	"github.com/TerraDharitri/drt-go-chain/consensus/spos/bls"
	v2 "github.com/TerraDharitri/drt-go-chain/consensus/spos/bls/v2"
	cryptoFactory "github.com/TerraDharitri/drt-go-chain/factory/crypto"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	nodeMock "github.com/TerraDharitri/drt-go-chain/testscommon/common"
	"github.com/TerraDharitri/drt-go-chain/testscommon/consensus"
	"github.com/TerraDharitri/drt-go-chain/testscommon/consensus/initializers"
	"github.com/TerraDharitri/drt-go-chain/testscommon/cryptoMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/enableEpochsHandlerMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/statusHandler"
)

func BenchmarkSubroundSignature_doSignatureJobForManagedKeys63(b *testing.B) {
	benchmarkSubroundSignatureDoSignatureJobForManagedKeys(b, 63)
}

func BenchmarkSubroundSignature_doSignatureJobForManagedKeys400(b *testing.B) {
	benchmarkSubroundSignatureDoSignatureJobForManagedKeys(b, 400)
}

func createMultiSignerSetup(grSize uint16, suite crypto.Suite) (crypto.KeyGenerator, map[string]crypto.PrivateKey) {
	kg := signing.NewKeyGenerator(suite)
	mapKeys := make(map[string]crypto.PrivateKey)

	for i := uint16(0); i < grSize; i++ {
		sk, pk := kg.GeneratePair()

		pubKey, _ := pk.ToByteArray()
		mapKeys[string(pubKey)] = sk
	}
	return kg, mapKeys
}

func benchmarkSubroundSignatureDoSignatureJobForManagedKeys(b *testing.B, numberOfKeys int) {
	container := consensus.InitConsensusCore()
	enableEpochsHandler := &enableEpochsHandlerMock.EnableEpochsHandlerStub{
		IsFlagEnabledInEpochCalled: func(flag core.EnableEpochFlag, epoch uint32) bool {
			return flag == common.AndromedaFlag
		},
	}
	container.SetEnableEpochsHandler(enableEpochsHandler)
	llSigner := &mclMultiSig.BlsMultiSignerKOSK{}

	suite := mcl.NewSuiteBLS12()
	kg, mapKeys := createMultiSignerSetup(uint16(numberOfKeys), suite)

	multiSigHandler, _ := multisig.NewBLSMultisig(llSigner, kg)

	keysHandlerMock := &testscommon.KeysHandlerStub{
		IsKeyManagedByCurrentNodeCalled: func(pkBytes []byte) bool {
			return true
		},
		GetHandledPrivateKeyCalled: func(pkBytes []byte) crypto.PrivateKey {
			return mapKeys[string(pkBytes)]
		},
	}

	args := cryptoFactory.ArgsSigningHandler{
		PubKeys: initializers.CreateEligibleListFromMap(mapKeys),
		MultiSignerContainer: &cryptoMocks.MultiSignerContainerStub{
			GetMultiSignerCalled: func(epoch uint32) (crypto.MultiSigner, error) {
				return multiSigHandler, nil
			}},
		SingleSigner: &cryptoMocks.SingleSignerStub{},
		KeyGenerator: kg,
		KeysHandler:  keysHandlerMock,
	}
	signingHandler, err := cryptoFactory.NewSigningHandler(args)
	require.Nil(b, err)

	container.SetSigningHandler(signingHandler)
	consensusState := initializers.InitConsensusStateWithArgs(keysHandlerMock, mapKeys)
	ch := make(chan bool, 1)

	sr, _ := spos.NewSubround(
		bls.SrBlock,
		bls.SrSignature,
		bls.SrEndRound,
		int64(70*roundTimeDuration/100),
		int64(85*roundTimeDuration/100),
		"(SIGNATURE)",
		consensusState,
		ch,
		executeStoredMessages,
		container,
		chainID,
		currentPid,
		&statusHandler.AppStatusHandlerStub{},
	)

	signatureSentForPks := make(map[string]struct{})
	mutex := sync.Mutex{}
	srSignature, _ := v2.NewSubroundSignature(
		sr,
		&statusHandler.AppStatusHandlerStub{},
		&testscommon.SentSignatureTrackerStub{
			SignatureSentCalled: func(pkBytes []byte) {
				mutex.Lock()
				signatureSentForPks[string(pkBytes)] = struct{}{}
				mutex.Unlock()
			},
		},
		&consensus.SposWorkerMock{},
		&nodeMock.ThrottlerStub{},
	)

	sr.SetHeader(&block.Header{})
	sr.SetSelfPubKey("OTHER")

	b.ResetTimer()
	b.StopTimer()

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		r := srSignature.DoSignatureJobForManagedKeys(context.TODO())
		b.StopTimer()

		require.True(b, r)
	}
}
