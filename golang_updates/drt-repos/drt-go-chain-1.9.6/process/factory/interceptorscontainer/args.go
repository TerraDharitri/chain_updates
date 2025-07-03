package interceptorscontainer

import (
	crypto "github.com/TerraDharitri/drt-go-chain-crypto"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/heartbeat"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
	"github.com/TerraDharitri/drt-go-chain/state"
)

// CommonInterceptorsContainerFactoryArgs holds the arguments needed for the metachain/shard interceptors factories
type CommonInterceptorsContainerFactoryArgs struct {
	CoreComponents                 process.CoreComponentsHolder
	CryptoComponents               process.CryptoComponentsHolder
	Accounts                       state.AccountsAdapter
	ShardCoordinator               sharding.Coordinator
	NodesCoordinator               nodesCoordinator.NodesCoordinator
	MainMessenger                  process.TopicHandler
	FullArchiveMessenger           process.TopicHandler
	Store                          dataRetriever.StorageService
	DataPool                       dataRetriever.PoolsHolder
	MaxTxNonceDeltaAllowed         int
	TxFeeHandler                   process.FeeHandler
	BlockBlackList                 process.TimeCacher
	HeaderSigVerifier              process.InterceptedHeaderSigVerifier
	HeaderIntegrityVerifier        process.HeaderIntegrityVerifier
	ValidityAttester               process.ValidityAttester
	EpochStartTrigger              process.EpochStartTriggerHandler
	WhiteListHandler               process.WhiteListHandler
	WhiteListerVerifiedTxs         process.WhiteListHandler
	AntifloodHandler               process.P2PAntifloodHandler
	ArgumentsParser                process.ArgumentsParser
	PreferredPeersHolder           process.PreferredPeersHolderHandler
	SizeCheckDelta                 uint32
	RequestHandler                 process.RequestHandler
	PeerSignatureHandler           crypto.PeerSignatureHandler
	SignaturesHandler              process.SignaturesHandler
	HeartbeatExpiryTimespanInSec   int64
	MainPeerShardMapper            process.PeerShardMapper
	FullArchivePeerShardMapper     process.PeerShardMapper
	HardforkTrigger                heartbeat.HardforkTrigger
	NodeOperationMode              common.NodeOperation
	InterceptedDataVerifierFactory process.InterceptedDataVerifierFactory
}
