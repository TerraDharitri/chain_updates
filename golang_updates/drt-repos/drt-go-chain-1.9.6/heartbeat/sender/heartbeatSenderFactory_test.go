package sender

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	crypto "github.com/TerraDharitri/drt-go-chain-crypto"
	"github.com/TerraDharitri/drt-go-chain/heartbeat"
	"github.com/TerraDharitri/drt-go-chain/heartbeat/mock"
	"github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/cryptoMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/shardingMocks"
	"github.com/stretchr/testify/assert"
)

func createMockHeartbeatSenderFactoryArgs() argHeartbeatSenderFactory {
	return argHeartbeatSenderFactory{
		argBaseSender:              createMockBaseArgs(),
		baseVersionNumber:          "base version number",
		versionNumber:              "version number",
		nodeDisplayName:            "node name",
		identity:                   "identity",
		peerSubType:                core.RegularPeer,
		currentBlockProvider:       &mock.CurrentBlockProviderStub{},
		peerTypeProvider:           &mock.PeerTypeProviderStub{},
		managedPeersHolder:         &testscommon.ManagedPeersHolderStub{},
		shardCoordinator:           createShardCoordinatorInShard(0),
		nodesCoordinator:           &shardingMocks.NodesCoordinatorStub{},
		trieSyncStatisticsProvider: &testscommon.SizeSyncStatisticsHandlerStub{},
	}
}

func TestHeartbeatSenderFactory_createHeartbeatSender(t *testing.T) {
	t.Parallel()

	t.Run("ToByteArray fails should error", func(t *testing.T) {
		t.Parallel()

		args := createMockHeartbeatSenderFactoryArgs()
		args.privKey = &cryptoMocks.PrivateKeyStub{
			GeneratePublicStub: func() crypto.PublicKey {
				return &cryptoMocks.PublicKeyStub{
					ToByteArrayStub: func() ([]byte, error) {
						return nil, expectedErr
					},
				}
			},
		}
		hbSender, err := createHeartbeatSender(args)
		assert.True(t, errors.Is(err, expectedErr))
		assert.True(t, check.IfNil(hbSender))
	})
	t.Run("validator with keys managed should error", func(t *testing.T) {
		t.Parallel()

		args := createMockHeartbeatSenderFactoryArgs()
		args.nodesCoordinator = &shardingMocks.NodesCoordinatorStub{
			GetValidatorWithPublicKeyCalled: func(publicKey []byte) (validator nodesCoordinator.Validator, shardId uint32, err error) {
				return nil, 0, nil
			},
		}
		args.managedPeersHolder = &testscommon.ManagedPeersHolderStub{
			IsMultiKeyModeCalled: func() bool {
				return true
			},
		}
		hbSender, err := createHeartbeatSender(args)
		assert.True(t, errors.Is(err, heartbeat.ErrInvalidConfiguration))
		assert.True(t, strings.Contains(err.Error(), "isValidator"))
		assert.True(t, check.IfNil(hbSender))
	})
	t.Run("validator should create regular sender", func(t *testing.T) {
		t.Parallel()

		args := createMockHeartbeatSenderFactoryArgs()
		args.nodesCoordinator = &shardingMocks.NodesCoordinatorStub{
			GetValidatorWithPublicKeyCalled: func(publicKey []byte) (validator nodesCoordinator.Validator, shardId uint32, err error) {
				return nil, 0, nil
			},
		}
		hbSender, err := createHeartbeatSender(args)
		assert.Nil(t, err)
		assert.False(t, check.IfNil(hbSender))
		assert.Equal(t, "*sender.heartbeatSender", fmt.Sprintf("%T", hbSender))
	})
	t.Run("regular observer should create regular sender", func(t *testing.T) {
		t.Parallel()

		args := createMockHeartbeatSenderFactoryArgs()
		args.nodesCoordinator = &shardingMocks.NodesCoordinatorStub{
			GetValidatorWithPublicKeyCalled: func(publicKey []byte) (validator nodesCoordinator.Validator, shardId uint32, err error) {
				return nil, 0, errors.New("not validator")
			},
		}
		hbSender, err := createHeartbeatSender(args)
		assert.Nil(t, err)
		assert.False(t, check.IfNil(hbSender))
		assert.Equal(t, "*sender.heartbeatSender", fmt.Sprintf("%T", hbSender))
	})
	t.Run("not validator with keys managed should create multikey sender", func(t *testing.T) {
		t.Parallel()

		args := createMockHeartbeatSenderFactoryArgs()
		args.nodesCoordinator = &shardingMocks.NodesCoordinatorStub{
			GetValidatorWithPublicKeyCalled: func(publicKey []byte) (validator nodesCoordinator.Validator, shardId uint32, err error) {
				return nil, 0, errors.New("not validator")
			},
		}
		args.managedPeersHolder = &testscommon.ManagedPeersHolderStub{
			IsMultiKeyModeCalled: func() bool {
				return true
			},
		}
		hbSender, err := createHeartbeatSender(args)
		assert.Nil(t, err)
		assert.False(t, check.IfNil(hbSender))
		assert.Equal(t, "*sender.multikeyHeartbeatSender", fmt.Sprintf("%T", hbSender))
	})
}
