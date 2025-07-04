package spos

import (
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"

	"github.com/TerraDharitri/drt-go-chain/common"
	cryptoCommon "github.com/TerraDharitri/drt-go-chain/common/crypto"
	"github.com/TerraDharitri/drt-go-chain/consensus"
	"github.com/TerraDharitri/drt-go-chain/epochStart"
	"github.com/TerraDharitri/drt-go-chain/ntp"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
)

// ConsensusCore implements ConsensusCoreHandler and provides access to common functionality
// for the rest of the consensus structures
type ConsensusCore struct {
	blockChain                    data.ChainHandler
	blockProcessor                process.BlockProcessor
	bootstrapper                  process.Bootstrapper
	broadcastMessenger            consensus.BroadcastMessenger
	chronologyHandler             consensus.ChronologyHandler
	hasher                        hashing.Hasher
	marshalizer                   marshal.Marshalizer
	multiSignerContainer          cryptoCommon.MultiSignerContainer
	roundHandler                  consensus.RoundHandler
	shardCoordinator              sharding.Coordinator
	nodesCoordinator              nodesCoordinator.NodesCoordinator
	syncTimer                     ntp.SyncTimer
	epochStartRegistrationHandler epochStart.RegistrationHandler
	antifloodHandler              consensus.P2PAntifloodHandler
	peerHonestyHandler            consensus.PeerHonestyHandler
	headerSigVerifier             consensus.HeaderSigVerifier
	fallbackHeaderValidator       consensus.FallbackHeaderValidator
	nodeRedundancyHandler         consensus.NodeRedundancyHandler
	scheduledProcessor            consensus.ScheduledProcessor
	messageSigningHandler         consensus.P2PSigningHandler
	peerBlacklistHandler          consensus.PeerBlacklistHandler
	signingHandler                consensus.SigningHandler
	enableEpochsHandler           common.EnableEpochsHandler
	equivalentProofsPool          consensus.EquivalentProofsPool
	epochNotifier                 process.EpochNotifier
	invalidSignersCache           InvalidSignersCache
}

// ConsensusCoreArgs store all arguments that are needed to create a ConsensusCore object
type ConsensusCoreArgs struct {
	BlockChain                    data.ChainHandler
	BlockProcessor                process.BlockProcessor
	Bootstrapper                  process.Bootstrapper
	BroadcastMessenger            consensus.BroadcastMessenger
	ChronologyHandler             consensus.ChronologyHandler
	Hasher                        hashing.Hasher
	Marshalizer                   marshal.Marshalizer
	MultiSignerContainer          cryptoCommon.MultiSignerContainer
	RoundHandler                  consensus.RoundHandler
	ShardCoordinator              sharding.Coordinator
	NodesCoordinator              nodesCoordinator.NodesCoordinator
	SyncTimer                     ntp.SyncTimer
	EpochStartRegistrationHandler epochStart.RegistrationHandler
	AntifloodHandler              consensus.P2PAntifloodHandler
	PeerHonestyHandler            consensus.PeerHonestyHandler
	HeaderSigVerifier             consensus.HeaderSigVerifier
	FallbackHeaderValidator       consensus.FallbackHeaderValidator
	NodeRedundancyHandler         consensus.NodeRedundancyHandler
	ScheduledProcessor            consensus.ScheduledProcessor
	MessageSigningHandler         consensus.P2PSigningHandler
	PeerBlacklistHandler          consensus.PeerBlacklistHandler
	SigningHandler                consensus.SigningHandler
	EnableEpochsHandler           common.EnableEpochsHandler
	EquivalentProofsPool          consensus.EquivalentProofsPool
	EpochNotifier                 process.EpochNotifier
	InvalidSignersCache           InvalidSignersCache
}

// NewConsensusCore creates a new ConsensusCore instance
func NewConsensusCore(
	args *ConsensusCoreArgs,
) (*ConsensusCore, error) {
	consensusCore := &ConsensusCore{
		blockChain:                    args.BlockChain,
		blockProcessor:                args.BlockProcessor,
		bootstrapper:                  args.Bootstrapper,
		broadcastMessenger:            args.BroadcastMessenger,
		chronologyHandler:             args.ChronologyHandler,
		hasher:                        args.Hasher,
		marshalizer:                   args.Marshalizer,
		multiSignerContainer:          args.MultiSignerContainer,
		roundHandler:                  args.RoundHandler,
		shardCoordinator:              args.ShardCoordinator,
		nodesCoordinator:              args.NodesCoordinator,
		syncTimer:                     args.SyncTimer,
		epochStartRegistrationHandler: args.EpochStartRegistrationHandler,
		antifloodHandler:              args.AntifloodHandler,
		peerHonestyHandler:            args.PeerHonestyHandler,
		headerSigVerifier:             args.HeaderSigVerifier,
		fallbackHeaderValidator:       args.FallbackHeaderValidator,
		nodeRedundancyHandler:         args.NodeRedundancyHandler,
		scheduledProcessor:            args.ScheduledProcessor,
		messageSigningHandler:         args.MessageSigningHandler,
		peerBlacklistHandler:          args.PeerBlacklistHandler,
		signingHandler:                args.SigningHandler,
		enableEpochsHandler:           args.EnableEpochsHandler,
		equivalentProofsPool:          args.EquivalentProofsPool,
		epochNotifier:                 args.EpochNotifier,
		invalidSignersCache:           args.InvalidSignersCache,
	}

	err := ValidateConsensusCore(consensusCore)
	if err != nil {
		return nil, err
	}

	return consensusCore, nil
}

// Blockchain gets the ChainHandler stored in the ConsensusCore
func (cc *ConsensusCore) Blockchain() data.ChainHandler {
	return cc.blockChain
}

// GetAntiFloodHandler will return the antiflood handler which will be used in subrounds
func (cc *ConsensusCore) GetAntiFloodHandler() consensus.P2PAntifloodHandler {
	return cc.antifloodHandler
}

// BlockProcessor gets the BlockProcessor stored in the ConsensusCore
func (cc *ConsensusCore) BlockProcessor() process.BlockProcessor {
	return cc.blockProcessor
}

// BootStrapper gets the Bootstrapper stored in the ConsensusCore
func (cc *ConsensusCore) BootStrapper() process.Bootstrapper {
	return cc.bootstrapper
}

// BroadcastMessenger gets the BroadcastMessenger stored in the ConsensusCore
func (cc *ConsensusCore) BroadcastMessenger() consensus.BroadcastMessenger {
	return cc.broadcastMessenger
}

// Chronology gets the ChronologyHandler stored in the ConsensusCore
func (cc *ConsensusCore) Chronology() consensus.ChronologyHandler {
	return cc.chronologyHandler
}

// Hasher gets the Hasher stored in the ConsensusCore
func (cc *ConsensusCore) Hasher() hashing.Hasher {
	return cc.hasher
}

// Marshalizer gets the Marshalizer stored in the ConsensusCore
func (cc *ConsensusCore) Marshalizer() marshal.Marshalizer {
	return cc.marshalizer
}

// MultiSignerContainer gets the MultiSignerContainer stored in the ConsensusCore
func (cc *ConsensusCore) MultiSignerContainer() cryptoCommon.MultiSignerContainer {
	return cc.multiSignerContainer
}

// RoundHandler gets the RoundHandler stored in the ConsensusCore
func (cc *ConsensusCore) RoundHandler() consensus.RoundHandler {
	return cc.roundHandler
}

// ShardCoordinator gets the ShardCoordinator stored in the ConsensusCore
func (cc *ConsensusCore) ShardCoordinator() sharding.Coordinator {
	return cc.shardCoordinator
}

// SyncTimer gets the SyncTimer stored in the ConsensusCore
func (cc *ConsensusCore) SyncTimer() ntp.SyncTimer {
	return cc.syncTimer
}

// NodesCoordinator gets the NodesCoordinator stored in the ConsensusCore
func (cc *ConsensusCore) NodesCoordinator() nodesCoordinator.NodesCoordinator {
	return cc.nodesCoordinator
}

// EpochStartRegistrationHandler returns the epoch start registration handler
func (cc *ConsensusCore) EpochStartRegistrationHandler() epochStart.RegistrationHandler {
	return cc.epochStartRegistrationHandler
}

// EpochNotifier returns the epoch notifier
func (cc *ConsensusCore) EpochNotifier() process.EpochNotifier {
	return cc.epochNotifier
}

// PeerHonestyHandler will return the peer honesty handler which will be used in subrounds
func (cc *ConsensusCore) PeerHonestyHandler() consensus.PeerHonestyHandler {
	return cc.peerHonestyHandler
}

// HeaderSigVerifier returns the sig verifier handler which will be used in subrounds
func (cc *ConsensusCore) HeaderSigVerifier() consensus.HeaderSigVerifier {
	return cc.headerSigVerifier
}

// FallbackHeaderValidator will return the fallback header validator which will be used in subrounds
func (cc *ConsensusCore) FallbackHeaderValidator() consensus.FallbackHeaderValidator {
	return cc.fallbackHeaderValidator
}

// NodeRedundancyHandler will return the node redundancy handler which will be used in subrounds
func (cc *ConsensusCore) NodeRedundancyHandler() consensus.NodeRedundancyHandler {
	return cc.nodeRedundancyHandler
}

// ScheduledProcessor will return the scheduled processor
func (cc *ConsensusCore) ScheduledProcessor() consensus.ScheduledProcessor {
	return cc.scheduledProcessor
}

// MessageSigningHandler will return the message signing handler
func (cc *ConsensusCore) MessageSigningHandler() consensus.P2PSigningHandler {
	return cc.messageSigningHandler
}

// PeerBlacklistHandler will return the peer blacklist handler
func (cc *ConsensusCore) PeerBlacklistHandler() consensus.PeerBlacklistHandler {
	return cc.peerBlacklistHandler
}

// SigningHandler will return the signing handler component
func (cc *ConsensusCore) SigningHandler() consensus.SigningHandler {
	return cc.signingHandler
}

// EnableEpochsHandler returns the enable epochs handler component
func (cc *ConsensusCore) EnableEpochsHandler() common.EnableEpochsHandler {
	return cc.enableEpochsHandler
}

// EquivalentProofsPool returns the equivalent proofs component
func (cc *ConsensusCore) EquivalentProofsPool() consensus.EquivalentProofsPool {
	return cc.equivalentProofsPool
}

// InvalidSignersCache returns the invalid signers cache component
func (cc *ConsensusCore) InvalidSignersCache() InvalidSignersCache {
	return cc.invalidSignersCache
}

// SetBlockchain sets blockchain handler
func (cc *ConsensusCore) SetBlockchain(blockChain data.ChainHandler) {
	cc.blockChain = blockChain
}

// SetBlockProcessor sets block processor
func (cc *ConsensusCore) SetBlockProcessor(blockProcessor process.BlockProcessor) {
	cc.blockProcessor = blockProcessor
}

// SetBootStrapper sets process bootstrapper
func (cc *ConsensusCore) SetBootStrapper(bootstrapper process.Bootstrapper) {
	cc.bootstrapper = bootstrapper
}

// SetBroadcastMessenger sets broadcast messenger
func (cc *ConsensusCore) SetBroadcastMessenger(broadcastMessenger consensus.BroadcastMessenger) {
	cc.broadcastMessenger = broadcastMessenger
}

// SetChronology sets chronology
func (cc *ConsensusCore) SetChronology(chronologyHandler consensus.ChronologyHandler) {
	cc.chronologyHandler = chronologyHandler
}

// SetHasher sets hasher component
func (cc *ConsensusCore) SetHasher(hasher hashing.Hasher) {
	cc.hasher = hasher
}

// SetMarshalizer sets marshaller component
func (cc *ConsensusCore) SetMarshalizer(marshalizer marshal.Marshalizer) {
	cc.marshalizer = marshalizer
}

// SetMultiSignerContainer sets multi signer container
func (cc *ConsensusCore) SetMultiSignerContainer(multiSignerContainer cryptoCommon.MultiSignerContainer) {
	cc.multiSignerContainer = multiSignerContainer
}

// SetRoundHandler sets round handler
func (cc *ConsensusCore) SetRoundHandler(roundHandler consensus.RoundHandler) {
	cc.roundHandler = roundHandler
}

// SetShardCoordinator set shard coordinator
func (cc *ConsensusCore) SetShardCoordinator(shardCoordinator sharding.Coordinator) {
	cc.shardCoordinator = shardCoordinator
}

// SetSyncTimer sets sync timer
func (cc *ConsensusCore) SetSyncTimer(syncTimer ntp.SyncTimer) {
	cc.syncTimer = syncTimer
}

// SetNodesCoordinator sets nodes coordinaotr
func (cc *ConsensusCore) SetNodesCoordinator(nodesCoordinator nodesCoordinator.NodesCoordinator) {
	cc.nodesCoordinator = nodesCoordinator
}

// SetEpochStartNotifier sets epoch start notifier
func (cc *ConsensusCore) SetEpochStartNotifier(epochStartNotifier epochStart.RegistrationHandler) {
	cc.epochStartRegistrationHandler = epochStartNotifier
}

// SetAntifloodHandler sets antiflood handler
func (cc *ConsensusCore) SetAntifloodHandler(antifloodHandler consensus.P2PAntifloodHandler) {
	cc.antifloodHandler = antifloodHandler
}

// SetPeerHonestyHandler sets peer honesty handler
func (cc *ConsensusCore) SetPeerHonestyHandler(peerHonestyHandler consensus.PeerHonestyHandler) {
	cc.peerHonestyHandler = peerHonestyHandler
}

// SetScheduledProcessor set scheduled processor
func (cc *ConsensusCore) SetScheduledProcessor(scheduledProcessor consensus.ScheduledProcessor) {
	cc.scheduledProcessor = scheduledProcessor
}

// SetPeerBlacklistHandler sets peer blacklist handlerc
func (cc *ConsensusCore) SetPeerBlacklistHandler(peerBlacklistHandler consensus.PeerBlacklistHandler) {
	cc.peerBlacklistHandler = peerBlacklistHandler
}

// SetHeaderSigVerifier sets header sig verifier
func (cc *ConsensusCore) SetHeaderSigVerifier(headerSigVerifier consensus.HeaderSigVerifier) {
	cc.headerSigVerifier = headerSigVerifier
}

// SetFallbackHeaderValidator sets fallback header validaor
func (cc *ConsensusCore) SetFallbackHeaderValidator(fallbackHeaderValidator consensus.FallbackHeaderValidator) {
	cc.fallbackHeaderValidator = fallbackHeaderValidator
}

// SetNodeRedundancyHandler set nodes redundancy handler
func (cc *ConsensusCore) SetNodeRedundancyHandler(nodeRedundancyHandler consensus.NodeRedundancyHandler) {
	cc.nodeRedundancyHandler = nodeRedundancyHandler
}

// SetMessageSigningHandler sets message signing handler
func (cc *ConsensusCore) SetMessageSigningHandler(messageSigningHandler consensus.P2PSigningHandler) {
	cc.messageSigningHandler = messageSigningHandler
}

// SetSigningHandler sets signing handler
func (cc *ConsensusCore) SetSigningHandler(signingHandler consensus.SigningHandler) {
	cc.signingHandler = signingHandler
}

// SetEnableEpochsHandler sets enable eopchs handler
func (cc *ConsensusCore) SetEnableEpochsHandler(enableEpochsHandler common.EnableEpochsHandler) {
	cc.enableEpochsHandler = enableEpochsHandler
}

// SetEquivalentProofsPool sets equivalent proofs pool
func (cc *ConsensusCore) SetEquivalentProofsPool(proofPool consensus.EquivalentProofsPool) {
	cc.equivalentProofsPool = proofPool
}

// SetEpochNotifier sets epoch notifier
func (cc *ConsensusCore) SetEpochNotifier(epochNotifier process.EpochNotifier) {
	cc.epochNotifier = epochNotifier
}

// SetInvalidSignersCache sets the invalid signers cache
func (cc *ConsensusCore) SetInvalidSignersCache(cache InvalidSignersCache) {
	cc.invalidSignersCache = cache
}

// IsInterfaceNil returns true if there is no value under the interface
func (cc *ConsensusCore) IsInterfaceNil() bool {
	return cc == nil
}
