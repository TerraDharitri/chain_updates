package resolvers

import (
	"encoding/hex"
	"fmt"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data/batch"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	logger "github.com/TerraDharitri/drt-go-chain-logger"

	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/p2p"
	"github.com/TerraDharitri/drt-go-chain/storage"
)

// maxBuffToSendValidatorsInfo represents max buffer size to send in bytes
const maxBuffToSendValidatorsInfo = 1 << 18 // 256KB

// ArgValidatorInfoResolver is the argument structure used to create a new validator info resolver instance
type ArgValidatorInfoResolver struct {
	SenderResolver       dataRetriever.TopicResolverSender
	Marshaller           marshal.Marshalizer
	AntifloodHandler     dataRetriever.P2PAntifloodHandler
	Throttler            dataRetriever.ResolverThrottler
	ValidatorInfoPool    dataRetriever.ShardedDataCacherNotifier
	ValidatorInfoStorage storage.Storer
	DataPacker           dataRetriever.DataPacker
	IsFullHistoryNode    bool
}

// validatorInfoResolver is a wrapper over Resolver that is specialized in resolving validator info requests
type validatorInfoResolver struct {
	dataRetriever.TopicResolverSender
	messageProcessor
	baseStorageResolver
	validatorInfoPool    dataRetriever.ShardedDataCacherNotifier
	validatorInfoStorage storage.Storer
	dataPacker           dataRetriever.DataPacker
}

// NewValidatorInfoResolver creates a validator info resolver
func NewValidatorInfoResolver(args ArgValidatorInfoResolver) (*validatorInfoResolver, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	return &validatorInfoResolver{
		TopicResolverSender: args.SenderResolver,
		messageProcessor: messageProcessor{
			marshalizer:      args.Marshaller,
			antifloodHandler: args.AntifloodHandler,
			throttler:        args.Throttler,
			topic:            args.SenderResolver.RequestTopic(),
		},
		baseStorageResolver:  createBaseStorageResolver(args.ValidatorInfoStorage, args.IsFullHistoryNode),
		validatorInfoPool:    args.ValidatorInfoPool,
		validatorInfoStorage: args.ValidatorInfoStorage,
		dataPacker:           args.DataPacker,
	}, nil
}

func checkArgs(args ArgValidatorInfoResolver) error {
	if check.IfNil(args.SenderResolver) {
		return dataRetriever.ErrNilResolverSender
	}
	if check.IfNil(args.Marshaller) {
		return dataRetriever.ErrNilMarshalizer
	}
	if check.IfNil(args.AntifloodHandler) {
		return dataRetriever.ErrNilAntifloodHandler
	}
	if check.IfNil(args.Throttler) {
		return dataRetriever.ErrNilThrottler
	}
	if check.IfNil(args.ValidatorInfoPool) {
		return dataRetriever.ErrNilValidatorInfoPool
	}
	if check.IfNil(args.ValidatorInfoStorage) {
		return dataRetriever.ErrNilValidatorInfoStorage
	}
	if check.IfNil(args.DataPacker) {
		return dataRetriever.ErrNilDataPacker
	}

	return nil
}

// ProcessReceivedMessage represents the callback func from the p2p.Messenger that is called each time a new message is received
// (for the topic this validator was registered to, usually a request topic)
func (res *validatorInfoResolver) ProcessReceivedMessage(message p2p.MessageP2P, fromConnectedPeer core.PeerID, source p2p.MessageHandler) ([]byte, error) {
	err := res.canProcessMessage(message, fromConnectedPeer)
	if err != nil {
		return nil, err
	}

	res.throttler.StartProcessing()
	defer res.throttler.EndProcessing()

	rd, err := res.parseReceivedMessage(message, fromConnectedPeer)
	if err != nil {
		return nil, err
	}

	switch rd.Type {
	case dataRetriever.HashType:
		err = res.resolveHashRequest(rd.Value, rd.Epoch, fromConnectedPeer, source)
	case dataRetriever.HashArrayType:
		err = res.resolveMultipleHashesRequest(rd.Value, rd.Epoch, fromConnectedPeer, source)
	default:
		err = fmt.Errorf("%w for value %s", dataRetriever.ErrRequestTypeNotImplemented, logger.DisplayByteSlice(rd.Value))
	}

	if err != nil {
		return nil, err
	}
	return []byte{}, nil
}

// resolveHashRequest sends the response for a hash request
func (res *validatorInfoResolver) resolveHashRequest(hash []byte, epoch uint32, pid core.PeerID, source p2p.MessageHandler) error {
	data, err := res.fetchValidatorInfoByteSlice(hash, epoch)
	if err != nil {
		return err
	}

	return res.marshalAndSend(data, pid, source)
}

// resolveMultipleHashesRequest sends the response for a hash array type request
func (res *validatorInfoResolver) resolveMultipleHashesRequest(hashesBuff []byte, epoch uint32, pid core.PeerID, source p2p.MessageHandler) error {
	b := batch.Batch{}
	err := res.marshalizer.Unmarshal(&b, hashesBuff)
	if err != nil {
		return err
	}
	hashes := b.Data

	validatorInfoForHashes, err := res.fetchValidatorInfoForHashes(hashes, epoch)
	if err != nil {
		outputHashes := ""
		for _, hash := range hashes {
			outputHashes += hex.EncodeToString(hash) + " "
		}
		return fmt.Errorf("resolveMultipleHashesRequest error %w from buff %s", err, outputHashes)
	}

	return res.sendValidatorInfoForHashes(validatorInfoForHashes, pid, source)
}

func (res *validatorInfoResolver) sendValidatorInfoForHashes(validatorInfoForHashes [][]byte, pid core.PeerID, source p2p.MessageHandler) error {
	buffsToSend, err := res.dataPacker.PackDataInChunks(validatorInfoForHashes, maxBuffToSendValidatorsInfo)
	if err != nil {
		return err
	}

	for _, buff := range buffsToSend {
		err = res.Send(buff, pid, source)
		if err != nil {
			return err
		}
	}

	return nil
}

func (res *validatorInfoResolver) fetchValidatorInfoForHashes(hashes [][]byte, epoch uint32) ([][]byte, error) {
	validatorInfos := make([][]byte, 0)
	for _, hash := range hashes {
		validatorInfoForHash, _ := res.fetchValidatorInfoByteSlice(hash, epoch)
		if validatorInfoForHash != nil {
			validatorInfos = append(validatorInfos, validatorInfoForHash)
		}
	}

	if len(validatorInfos) == 0 {
		return nil, dataRetriever.ErrValidatorInfoNotFound
	}

	return validatorInfos, nil
}

func (res *validatorInfoResolver) fetchValidatorInfoByteSlice(hash []byte, epoch uint32) ([]byte, error) {
	data, ok := res.validatorInfoPool.SearchFirstData(hash)
	if ok {
		return res.marshalizer.Marshal(data)
	}

	buff, err := res.getFromStorage(hash, epoch)
	if err != nil {
		res.DebugHandler().LogFailedToResolveData(
			res.topic,
			hash,
			err,
		)
		return nil, err
	}

	res.DebugHandler().LogSucceededToResolveData(res.topic, hash)

	return buff, nil
}

func (res *validatorInfoResolver) marshalAndSend(data []byte, pid core.PeerID, source p2p.MessageHandler) error {
	b := &batch.Batch{
		Data: [][]byte{data},
	}
	buff, err := res.marshalizer.Marshal(b)
	if err != nil {
		return err
	}

	return res.Send(buff, pid, source)
}

// SetDebugHandler sets a debug handler
func (res *validatorInfoResolver) SetDebugHandler(handler dataRetriever.DebugHandler) error {
	return res.TopicResolverSender.SetDebugHandler(handler)
}

// Close returns nil
func (res *validatorInfoResolver) Close() error {
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (res *validatorInfoResolver) IsInterfaceNil() bool {
	return res == nil
}
