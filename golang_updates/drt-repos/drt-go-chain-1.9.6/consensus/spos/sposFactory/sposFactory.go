package sposFactory

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	crypto "github.com/TerraDharitri/mx-chain-crypto-go"

	"github.com/TerraDharitri/drt-go-chain/consensus"
	"github.com/TerraDharitri/drt-go-chain/consensus/broadcast"
	"github.com/TerraDharitri/drt-go-chain/consensus/spos"
	"github.com/TerraDharitri/drt-go-chain/consensus/spos/bls"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/sharding"
)

// GetConsensusCoreFactory returns a consensus service depending on the given parameter
func GetConsensusCoreFactory(consensusType string) (spos.ConsensusService, error) {
	switch consensusType {
	case blsConsensusType:
		return bls.NewConsensusService()
	default:
		return nil, ErrInvalidConsensusType
	}
}

// GetBroadcastMessenger returns a consensus service depending on the given parameter
func GetBroadcastMessenger(
	marshalizer marshal.Marshalizer,
	hasher hashing.Hasher,
	messenger consensus.P2PMessenger,
	shardCoordinator sharding.Coordinator,
	peerSignatureHandler crypto.PeerSignatureHandler,
	headersSubscriber consensus.HeadersPoolSubscriber,
	interceptorsContainer process.InterceptorsContainer,
	alarmScheduler core.TimersScheduler,
	keysHandler consensus.KeysHandler,
) (consensus.BroadcastMessenger, error) {

	if check.IfNil(shardCoordinator) {
		return nil, spos.ErrNilShardCoordinator
	}

	dbbArgs := &broadcast.ArgsDelayedBlockBroadcaster{
		InterceptorsContainer: interceptorsContainer,
		HeadersSubscriber:     headersSubscriber,
		ShardCoordinator:      shardCoordinator,
		LeaderCacheSize:       maxDelayCacheSize,
		ValidatorCacheSize:    maxDelayCacheSize,
		AlarmScheduler:        alarmScheduler,
	}

	delayedBroadcaster, err := broadcast.NewDelayedBlockBroadcaster(dbbArgs)
	if err != nil {
		return nil, err
	}

	commonMessengerArgs := broadcast.CommonMessengerArgs{
		Marshalizer:                marshalizer,
		Hasher:                     hasher,
		Messenger:                  messenger,
		ShardCoordinator:           shardCoordinator,
		PeerSignatureHandler:       peerSignatureHandler,
		HeadersSubscriber:          headersSubscriber,
		MaxDelayCacheSize:          maxDelayCacheSize,
		MaxValidatorDelayCacheSize: maxDelayCacheSize,
		InterceptorsContainer:      interceptorsContainer,
		AlarmScheduler:             alarmScheduler,
		KeysHandler:                keysHandler,
		DelayedBroadcaster:         delayedBroadcaster,
	}

	if shardCoordinator.SelfId() < shardCoordinator.NumberOfShards() {
		shardMessengerArgs := broadcast.ShardChainMessengerArgs{
			CommonMessengerArgs: commonMessengerArgs,
		}

		return broadcast.NewShardChainMessenger(shardMessengerArgs)
	}

	if shardCoordinator.SelfId() == core.MetachainShardId {
		metaMessengerArgs := broadcast.MetaChainMessengerArgs{
			CommonMessengerArgs: commonMessengerArgs,
		}

		return broadcast.NewMetaChainMessenger(metaMessengerArgs)
	}

	return nil, ErrInvalidShardId
}
