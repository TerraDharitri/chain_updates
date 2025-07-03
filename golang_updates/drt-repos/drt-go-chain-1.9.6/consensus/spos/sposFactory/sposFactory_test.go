package sposFactory_test

import (
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/stretchr/testify/assert"

	"github.com/TerraDharitri/drt-go-chain/consensus"
	"github.com/TerraDharitri/drt-go-chain/consensus/mock"
	"github.com/TerraDharitri/drt-go-chain/consensus/spos"
	"github.com/TerraDharitri/drt-go-chain/consensus/spos/sposFactory"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/hashingMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/p2pmocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/pool"
)

func TestGetConsensusCoreFactory_InvalidTypeShouldErr(t *testing.T) {
	t.Parallel()

	csf, err := sposFactory.GetConsensusCoreFactory("invalid")

	assert.Nil(t, csf)
	assert.Equal(t, sposFactory.ErrInvalidConsensusType, err)
}

func TestGetConsensusCoreFactory_BlsShouldWork(t *testing.T) {
	t.Parallel()

	csf, err := sposFactory.GetConsensusCoreFactory(consensus.BlsConsensusType)

	assert.Nil(t, err)
	assert.False(t, check.IfNil(csf))
}

func TestGetBroadcastMessenger_ShardShouldWork(t *testing.T) {
	t.Parallel()

	marshalizer := &mock.MarshalizerMock{}
	hasher := &hashingMocks.HasherMock{}
	messenger := &p2pmocks.MessengerStub{}
	shardCoord := mock.NewMultiShardsCoordinatorMock(3)
	shardCoord.SelfIDCalled = func() uint32 {
		return 0
	}
	peerSigHandler := &mock.PeerSignatureHandler{}
	headersSubscriber := &pool.HeadersPoolStub{}
	interceptosContainer := &testscommon.InterceptorsContainerStub{}
	alarmSchedulerStub := &testscommon.AlarmSchedulerStub{}

	bm, err := sposFactory.GetBroadcastMessenger(
		marshalizer,
		hasher,
		messenger,
		shardCoord,
		peerSigHandler,
		headersSubscriber,
		interceptosContainer,
		alarmSchedulerStub,
		&testscommon.KeysHandlerStub{},
	)

	assert.Nil(t, err)
	assert.NotNil(t, bm)
}

func TestGetBroadcastMessenger_MetachainShouldWork(t *testing.T) {
	t.Parallel()

	marshalizer := &mock.MarshalizerMock{}
	hasher := &hashingMocks.HasherMock{}
	messenger := &p2pmocks.MessengerStub{}
	shardCoord := mock.NewMultiShardsCoordinatorMock(3)
	shardCoord.SelfIDCalled = func() uint32 {
		return core.MetachainShardId
	}
	peerSigHandler := &mock.PeerSignatureHandler{}
	headersSubscriber := &pool.HeadersPoolStub{}
	interceptosContainer := &testscommon.InterceptorsContainerStub{}
	alarmSchedulerStub := &testscommon.AlarmSchedulerStub{}

	bm, err := sposFactory.GetBroadcastMessenger(
		marshalizer,
		hasher,
		messenger,
		shardCoord,
		peerSigHandler,
		headersSubscriber,
		interceptosContainer,
		alarmSchedulerStub,
		&testscommon.KeysHandlerStub{},
	)

	assert.Nil(t, err)
	assert.NotNil(t, bm)
}

func TestGetBroadcastMessenger_NilShardCoordinatorShouldErr(t *testing.T) {
	t.Parallel()

	headersSubscriber := &pool.HeadersPoolStub{}
	interceptosContainer := &testscommon.InterceptorsContainerStub{}
	alarmSchedulerStub := &testscommon.AlarmSchedulerStub{}

	bm, err := sposFactory.GetBroadcastMessenger(
		nil,
		nil,
		nil,
		nil,
		nil,
		headersSubscriber,
		interceptosContainer,
		alarmSchedulerStub,
		&testscommon.KeysHandlerStub{},
	)

	assert.Nil(t, bm)
	assert.Equal(t, spos.ErrNilShardCoordinator, err)
}

func TestGetBroadcastMessenger_InvalidShardIdShouldErr(t *testing.T) {
	t.Parallel()

	shardCoord := mock.NewMultiShardsCoordinatorMock(3)
	shardCoord.SelfIDCalled = func() uint32 {
		return 37
	}
	headersSubscriber := &pool.HeadersPoolStub{}
	interceptosContainer := &testscommon.InterceptorsContainerStub{}
	alarmSchedulerStub := &testscommon.AlarmSchedulerStub{}

	bm, err := sposFactory.GetBroadcastMessenger(
		nil,
		nil,
		nil,
		shardCoord,
		nil,
		headersSubscriber,
		interceptosContainer,
		alarmSchedulerStub,
		&testscommon.KeysHandlerStub{},
	)

	assert.Nil(t, bm)
	assert.Equal(t, sposFactory.ErrInvalidShardId, err)
}
