package bootstrap

import (
	"context"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/epochStart"
	"github.com/TerraDharitri/drt-go-chain/epochStart/bootstrap/disabled"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/factory"
	"github.com/TerraDharitri/drt-go-chain/process/interceptors"
	interceptorsFactory "github.com/TerraDharitri/drt-go-chain/process/interceptors/factory"
	"github.com/TerraDharitri/drt-go-chain/sharding"
)

var _ epochStart.StartOfEpochMetaSyncer = (*epochStartMetaSyncer)(nil)

type epochStartMetaSyncer struct {
	requestHandler                 RequestHandler
	messenger                      Messenger
	marshalizer                    marshal.Marshalizer
	hasher                         hashing.Hasher
	singleDataInterceptor          process.Interceptor
	proofsInterceptor              process.Interceptor
	metaBlockProcessor             EpochStartMetaBlockInterceptorProcessor
	interceptedDataVerifierFactory process.InterceptedDataVerifierFactory
}

// ArgsNewEpochStartMetaSyncer -
type ArgsNewEpochStartMetaSyncer struct {
	CoreComponentsHolder           process.CoreComponentsHolder
	CryptoComponentsHolder         process.CryptoComponentsHolder
	RequestHandler                 RequestHandler
	Messenger                      Messenger
	ShardCoordinator               sharding.Coordinator
	EconomicsData                  process.EconomicsDataHandler
	WhitelistHandler               process.WhiteListHandler
	StartInEpochConfig             config.EpochStartConfig
	ArgsParser                     process.ArgumentsParser
	HeaderIntegrityVerifier        process.HeaderIntegrityVerifier
	MetaBlockProcessor             EpochStartMetaBlockInterceptorProcessor
	InterceptedDataVerifierFactory process.InterceptedDataVerifierFactory
	ProofsPool                     dataRetriever.ProofsPool
	ProofsInterceptorProcessor     process.InterceptorProcessor
}

// NewEpochStartMetaSyncer will return a new instance of epochStartMetaSyncer
func NewEpochStartMetaSyncer(args ArgsNewEpochStartMetaSyncer) (*epochStartMetaSyncer, error) {
	if check.IfNil(args.CoreComponentsHolder) {
		return nil, epochStart.ErrNilCoreComponentsHolder
	}
	if check.IfNil(args.CryptoComponentsHolder) {
		return nil, epochStart.ErrNilCryptoComponentsHolder
	}
	if check.IfNil(args.CoreComponentsHolder.AddressPubKeyConverter()) {
		return nil, epochStart.ErrNilPubkeyConverter
	}
	if check.IfNil(args.HeaderIntegrityVerifier) {
		return nil, epochStart.ErrNilHeaderIntegrityVerifier
	}
	if check.IfNil(args.MetaBlockProcessor) {
		return nil, epochStart.ErrNilMetablockProcessor
	}
	if check.IfNil(args.InterceptedDataVerifierFactory) {
		return nil, epochStart.ErrNilInterceptedDataVerifierFactory
	}
	if check.IfNil(args.ProofsInterceptorProcessor) {
		return nil, epochStart.ErrNilEquivalentProofsProcessor
	}

	e := &epochStartMetaSyncer{
		requestHandler:                 args.RequestHandler,
		messenger:                      args.Messenger,
		marshalizer:                    args.CoreComponentsHolder.InternalMarshalizer(),
		hasher:                         args.CoreComponentsHolder.Hasher(),
		metaBlockProcessor:             args.MetaBlockProcessor,
		interceptedDataVerifierFactory: args.InterceptedDataVerifierFactory,
	}

	argsInterceptedDataFactory := interceptorsFactory.ArgInterceptedDataFactory{
		CoreComponents:          args.CoreComponentsHolder,
		CryptoComponents:        args.CryptoComponentsHolder,
		ShardCoordinator:        args.ShardCoordinator,
		NodesCoordinator:        disabled.NewNodesCoordinator(),
		FeeHandler:              args.EconomicsData,
		HeaderSigVerifier:       disabled.NewHeaderSigVerifier(),
		HeaderIntegrityVerifier: args.HeaderIntegrityVerifier,
		ValidityAttester:        disabled.NewValidityAttester(),
		EpochStartTrigger:       disabled.NewEpochStartTrigger(),
		ArgsParser:              args.ArgsParser,
	}
	argsInterceptedMetaHeaderFactory := interceptorsFactory.ArgInterceptedMetaHeaderFactory{
		ArgInterceptedDataFactory: argsInterceptedDataFactory,
	}

	interceptedMetaHdrDataFactory, err := interceptorsFactory.NewInterceptedMetaHeaderDataFactory(&argsInterceptedMetaHeaderFactory)
	if err != nil {
		return nil, err
	}

	interceptedDataVerifier, err := e.interceptedDataVerifierFactory.Create(factory.MetachainBlocksTopic)
	if err != nil {
		return nil, err
	}

	e.singleDataInterceptor, err = interceptors.NewSingleDataInterceptor(
		interceptors.ArgSingleDataInterceptor{
			Topic:                   factory.MetachainBlocksTopic,
			DataFactory:             interceptedMetaHdrDataFactory,
			Processor:               args.MetaBlockProcessor,
			Throttler:               disabled.NewThrottler(),
			AntifloodHandler:        disabled.NewAntiFloodHandler(),
			WhiteListRequest:        args.WhitelistHandler,
			CurrentPeerId:           args.Messenger.ID(),
			PreferredPeersHolder:    disabled.NewPreferredPeersHolder(),
			InterceptedDataVerifier: interceptedDataVerifier,
		},
	)
	if err != nil {
		return nil, err
	}

	argsInterceptedEquivalentProofsFactory := interceptorsFactory.ArgInterceptedEquivalentProofsFactory{
		ArgInterceptedDataFactory: argsInterceptedDataFactory,
		ProofsPool:                args.ProofsPool,
	}
	interceptedEquivalentProofsFactory := interceptorsFactory.NewInterceptedEquivalentProofsFactory(argsInterceptedEquivalentProofsFactory)
	if err != nil {
		return nil, err
	}

	proofsTopic := common.EquivalentProofsTopic + core.CommunicationIdentifierBetweenShards(core.MetachainShardId, core.AllShardId)
	e.proofsInterceptor, err = interceptors.NewSingleDataInterceptor(
		interceptors.ArgSingleDataInterceptor{
			Topic:                   proofsTopic,
			DataFactory:             interceptedEquivalentProofsFactory,
			Processor:               args.ProofsInterceptorProcessor,
			Throttler:               disabled.NewThrottler(),
			AntifloodHandler:        disabled.NewAntiFloodHandler(),
			WhiteListRequest:        args.WhitelistHandler,
			CurrentPeerId:           args.Messenger.ID(),
			PreferredPeersHolder:    disabled.NewPreferredPeersHolder(),
			InterceptedDataVerifier: interceptedDataVerifier,
		},
	)
	if err != nil {
		return nil, err
	}

	return e, nil
}

// SyncEpochStartMeta syncs the latest epoch start metablock
func (e *epochStartMetaSyncer) SyncEpochStartMeta(timeToWait time.Duration) (data.MetaHeaderHandler, error) {
	err := e.initTopicForEpochStartMetaBlockInterceptor()
	if err != nil {
		return nil, err
	}
	defer func() {
		e.resetTopicsAndInterceptors()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), timeToWait)
	mb, errConsensusNotReached := e.metaBlockProcessor.GetEpochStartMetaBlock(ctx)
	cancel()

	if errConsensusNotReached != nil {
		return nil, errConsensusNotReached
	}

	return mb, nil
}

func (e *epochStartMetaSyncer) resetTopicsAndInterceptors() {
	err := e.messenger.UnregisterMessageProcessor(factory.MetachainBlocksTopic, common.EpochStartInterceptorsIdentifier)
	if err != nil {
		log.Trace("error unregistering message processors", "error", err)
	}

	proofsTopic := common.EquivalentProofsTopic + core.CommunicationIdentifierBetweenShards(core.MetachainShardId, core.AllShardId)
	err = e.messenger.UnregisterMessageProcessor(proofsTopic, common.EpochStartInterceptorsIdentifier)
	if err != nil {
		log.Trace("error unregistering message processors", "error", err)
	}
}

func (e *epochStartMetaSyncer) initTopicForEpochStartMetaBlockInterceptor() error {
	err := e.messenger.CreateTopic(factory.MetachainBlocksTopic, true)
	if err != nil {
		log.Warn("error messenger create topic", "error", err)
		return err
	}

	proofsTopic := common.EquivalentProofsTopic + core.CommunicationIdentifierBetweenShards(core.MetachainShardId, core.AllShardId)
	err = e.messenger.CreateTopic(proofsTopic, true)
	if err != nil {
		log.Warn("error messenger create topic", "topic", proofsTopic, "error", err)
		return err
	}

	e.resetTopicsAndInterceptors()
	err = e.messenger.RegisterMessageProcessor(factory.MetachainBlocksTopic, common.EpochStartInterceptorsIdentifier, e.singleDataInterceptor)
	if err != nil {
		return err
	}

	return e.messenger.RegisterMessageProcessor(proofsTopic, common.EpochStartInterceptorsIdentifier, e.proofsInterceptor)
}

// IsInterfaceNil returns true if underlying object is nil
func (e *epochStartMetaSyncer) IsInterfaceNil() bool {
	return e == nil
}
