package broadcast

import (
	"strings"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/core/partitioning"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	crypto "github.com/TerraDharitri/drt-go-chain-crypto"
	logger "github.com/TerraDharitri/drt-go-chain-logger"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/consensus"
	"github.com/TerraDharitri/drt-go-chain/consensus/spos"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/factory"
	"github.com/TerraDharitri/drt-go-chain/sharding"
)

var log = logger.GetOrCreate("consensus/broadcast")

type commonMessenger struct {
	marshalizer             marshal.Marshalizer
	hasher                  hashing.Hasher
	messenger               consensus.P2PMessenger
	shardCoordinator        sharding.Coordinator
	peerSignatureHandler    crypto.PeerSignatureHandler
	delayedBlockBroadcaster DelayedBroadcaster
	keysHandler             consensus.KeysHandler
}

// CommonMessengerArgs holds the arguments for creating commonMessenger instance
type CommonMessengerArgs struct {
	Marshalizer                marshal.Marshalizer
	Hasher                     hashing.Hasher
	Messenger                  consensus.P2PMessenger
	ShardCoordinator           sharding.Coordinator
	PeerSignatureHandler       crypto.PeerSignatureHandler
	HeadersSubscriber          consensus.HeadersPoolSubscriber
	InterceptorsContainer      process.InterceptorsContainer
	MaxDelayCacheSize          uint32
	MaxValidatorDelayCacheSize uint32
	AlarmScheduler             core.TimersScheduler
	KeysHandler                consensus.KeysHandler
	DelayedBroadcaster         DelayedBroadcaster
}

func checkCommonMessengerNilParameters(
	args CommonMessengerArgs,
) error {
	if check.IfNil(args.Marshalizer) {
		return spos.ErrNilMarshalizer
	}
	if check.IfNil(args.Hasher) {
		return spos.ErrNilHasher
	}
	if check.IfNil(args.Messenger) {
		return spos.ErrNilMessenger
	}
	if check.IfNil(args.ShardCoordinator) {
		return spos.ErrNilShardCoordinator
	}
	if check.IfNil(args.PeerSignatureHandler) {
		return spos.ErrNilPeerSignatureHandler
	}
	if check.IfNil(args.InterceptorsContainer) {
		return spos.ErrNilInterceptorsContainer
	}
	if check.IfNil(args.HeadersSubscriber) {
		return spos.ErrNilHeadersSubscriber
	}
	if check.IfNil(args.AlarmScheduler) {
		return spos.ErrNilAlarmScheduler
	}
	if args.MaxDelayCacheSize == 0 || args.MaxValidatorDelayCacheSize == 0 {
		return spos.ErrInvalidCacheSize
	}
	if check.IfNil(args.KeysHandler) {
		return ErrNilKeysHandler
	}
	if check.IfNil(args.DelayedBroadcaster) {
		return ErrNilDelayedBroadcaster
	}

	return nil
}

// BroadcastConsensusMessage will send on consensus topic the consensus message
func (cm *commonMessenger) BroadcastConsensusMessage(message *consensus.Message) error {
	privateKey := cm.keysHandler.GetHandledPrivateKey(message.PubKey)
	signature, err := cm.peerSignatureHandler.GetPeerSignature(privateKey, message.OriginatorPid)
	if err != nil {
		return err
	}

	message.Signature = signature

	buff, err := cm.marshalizer.Marshal(message)
	if err != nil {
		return err
	}

	consensusTopic := common.ConsensusTopic +
		cm.shardCoordinator.CommunicationIdentifier(cm.shardCoordinator.SelfId())

	cm.broadcast(consensusTopic, buff, message.PubKey)

	return nil
}

// BroadcastMiniBlocks will send on miniblocks topic the cross-shard miniblocks
func (cm *commonMessenger) BroadcastMiniBlocks(miniBlocks map[uint32][]byte, pkBytes []byte) error {
	for k, v := range miniBlocks {
		miniBlocksTopic := factory.MiniBlocksTopic +
			cm.shardCoordinator.CommunicationIdentifier(k)

		cm.broadcast(miniBlocksTopic, v, pkBytes)
	}

	if len(miniBlocks) > 0 {
		log.Debug("commonMessenger.BroadcastMiniBlocks",
			"num minblocks", len(miniBlocks),
		)
	}

	return nil
}

// BroadcastTransactions will send on transaction topic the transactions
func (cm *commonMessenger) BroadcastTransactions(transactions map[string][][]byte, pkBytes []byte) error {
	dataPacker, err := partitioning.NewSimpleDataPacker(cm.marshalizer)
	if err != nil {
		return err
	}

	txs := 0
	var packets [][]byte
	for topic, v := range transactions {
		txs += len(v)
		// forward txs to the destination shards in packets
		packets, err = dataPacker.PackDataInChunks(v, common.MaxBulkTransactionSize)
		if err != nil {
			return err
		}

		for _, buff := range packets {
			cm.broadcast(topic, buff, pkBytes)
		}
	}

	if txs > 0 {
		log.Debug("commonMessenger.BroadcastTransactions",
			"num txs", txs,
		)
	}

	return nil
}

// BroadcastBlockData broadcasts the miniblocks and transactions
func (cm *commonMessenger) BroadcastBlockData(
	miniBlocks map[uint32][]byte,
	transactions map[string][][]byte,
	pkBytes []byte,
	extraDelayForBroadcast time.Duration,
) {
	time.Sleep(extraDelayForBroadcast)

	if len(miniBlocks) > 0 {
		err := cm.BroadcastMiniBlocks(miniBlocks, pkBytes)
		if err != nil {
			log.Warn("commonMessenger.BroadcastBlockData: broadcast miniblocks", "error", err.Error())
		}
	}

	time.Sleep(common.ExtraDelayBetweenBroadcastMbsAndTxs)

	if len(transactions) > 0 {
		err := cm.BroadcastTransactions(transactions, pkBytes)
		if err != nil {
			log.Warn("commonMessenger.BroadcastBlockData: broadcast transactions", "error", err.Error())
		}
	}
}

func (cm *commonMessenger) extractMetaMiniBlocksAndTransactions(
	miniBlocks map[uint32][]byte,
	transactions map[string][][]byte,
) (map[uint32][]byte, map[string][][]byte) {

	metaMiniBlocks := make(map[uint32][]byte)
	metaTransactions := make(map[string][][]byte)

	for shardID, mbsMarshalized := range miniBlocks {
		if shardID != core.MetachainShardId {
			continue
		}

		metaMiniBlocks[shardID] = mbsMarshalized
		delete(miniBlocks, shardID)
	}

	identifier := cm.shardCoordinator.CommunicationIdentifier(core.MetachainShardId)

	for broadcastTopic, txsMarshalized := range transactions {
		if !strings.Contains(broadcastTopic, identifier) {
			continue
		}

		metaTransactions[broadcastTopic] = txsMarshalized
		delete(transactions, broadcastTopic)
	}

	return metaMiniBlocks, metaTransactions
}

func (cm *commonMessenger) broadcast(topic string, data []byte, pkBytes []byte) {
	if cm.keysHandler.IsOriginalPublicKeyOfTheNode(pkBytes) {
		cm.messenger.Broadcast(topic, data)
		return
	}

	skBytes, pid, err := cm.keysHandler.GetP2PIdentity(pkBytes)
	if err != nil {
		log.Error("setup error in commonMessenger.broadcast - public key is managed but does not contain p2p sign info",
			"pk", pkBytes, "error", err)
		return
	}

	cm.messenger.BroadcastUsingPrivateKey(topic, data, pid, skBytes)
}

func (cm *commonMessenger) broadcastEquivalentProof(proof data.HeaderProofHandler, pkBytes []byte, topic string) error {
	if check.IfNil(proof) {
		return spos.ErrNilHeaderProof
	}

	msgProof, err := cm.marshalizer.Marshal(proof)
	if err != nil {
		return err
	}

	cm.broadcast(topic, msgProof, pkBytes)

	return nil
}
