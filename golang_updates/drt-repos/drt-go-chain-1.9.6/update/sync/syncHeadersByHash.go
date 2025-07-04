package sync

import (
	"context"
	"sync"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/storage"
	"github.com/TerraDharitri/drt-go-chain/update"
)

var _ update.MissingHeadersByHashSyncer = (*syncHeadersByHash)(nil)

type syncHeadersByHash struct {
	mutMissingHdrs          sync.Mutex
	mapHeaders              map[string]data.HeaderHandler
	mapHashes               map[string]struct{}
	pool                    dataRetriever.HeadersPool
	proofsPool              dataRetriever.ProofsPool
	storage                 update.HistoryStorer
	chReceivedAll           chan bool
	marshalizer             marshal.Marshalizer
	stopSyncing             bool
	syncedAll               bool
	requestHandler          process.RequestHandler
	waitTimeBetweenRequests time.Duration
	enableEpochsHandler     common.EnableEpochsHandler
}

// ArgsNewMissingHeadersByHashSyncer defines the arguments needed for the sycner
type ArgsNewMissingHeadersByHashSyncer struct {
	Storage             storage.Storer
	Cache               dataRetriever.HeadersPool
	ProofsPool          dataRetriever.ProofsPool
	Marshalizer         marshal.Marshalizer
	RequestHandler      process.RequestHandler
	EnableEpochsHandler common.EnableEpochsHandler
}

// NewMissingheadersByHashSyncer creates a syncer for all missing headers
func NewMissingheadersByHashSyncer(args ArgsNewMissingHeadersByHashSyncer) (*syncHeadersByHash, error) {
	if check.IfNil(args.Storage) {
		return nil, dataRetriever.ErrNilHeadersStorage
	}
	if check.IfNil(args.Cache) {
		return nil, update.ErrNilCacher
	}
	if check.IfNil(args.ProofsPool) {
		return nil, dataRetriever.ErrNilProofsPool
	}
	if check.IfNil(args.Marshalizer) {
		return nil, dataRetriever.ErrNilMarshalizer
	}
	if check.IfNil(args.RequestHandler) {
		return nil, process.ErrNilRequestHandler
	}
	if check.IfNil(args.EnableEpochsHandler) {
		return nil, process.ErrNilEnableEpochsHandler
	}

	p := &syncHeadersByHash{
		mutMissingHdrs:          sync.Mutex{},
		mapHeaders:              make(map[string]data.HeaderHandler),
		mapHashes:               make(map[string]struct{}),
		pool:                    args.Cache,
		proofsPool:              args.ProofsPool,
		storage:                 args.Storage,
		chReceivedAll:           make(chan bool),
		requestHandler:          args.RequestHandler,
		stopSyncing:             true,
		syncedAll:               false,
		marshalizer:             args.Marshalizer,
		waitTimeBetweenRequests: args.RequestHandler.RequestInterval(),
		enableEpochsHandler:     args.EnableEpochsHandler,
	}

	p.pool.RegisterHandler(p.receivedHeader)
	p.proofsPool.RegisterHandler(p.receivedProof)

	return p, nil
}

// SyncMissingHeadersByHash syncs the missing headers
func (m *syncHeadersByHash) SyncMissingHeadersByHash(shardIDs []uint32, headersHashes [][]byte, ctx context.Context) error {
	_ = core.EmptyChannel(m.chReceivedAll)

	mapHashesToRequest := make(map[string]uint32)
	for index, hash := range headersHashes {
		mapHashesToRequest[string(hash)] = shardIDs[index]
	}

	for {
		requestedHdrs := 0
		requestedProofs := 0

		m.mutMissingHdrs.Lock()
		m.stopSyncing = false
		for hash, shardId := range mapHashesToRequest {
			requestedHeader, requestedProof := m.updateMapsAndRequestIfNeeded(shardId, hash, mapHashesToRequest)
			if requestedHeader {
				requestedHdrs++
			}
			if requestedProof {
				requestedProofs++
			}
		}

		if requestedHdrs == 0 && requestedProofs == 0 {
			m.stopSyncing = true
			m.syncedAll = true
			m.mutMissingHdrs.Unlock()
			return nil
		}

		m.mutMissingHdrs.Unlock()

		select {
		case <-m.chReceivedAll:
			m.mutMissingHdrs.Lock()
			m.stopSyncing = true
			m.syncedAll = true
			m.mutMissingHdrs.Unlock()
			return nil
		case <-time.After(m.waitTimeBetweenRequests):
			continue
		case <-ctx.Done():
			m.mutMissingHdrs.Lock()
			m.stopSyncing = true
			m.mutMissingHdrs.Unlock()
			return update.ErrTimeIsOut
		}
	}
}

func (m *syncHeadersByHash) updateMapsAndRequestIfNeeded(
	shardId uint32,
	hash string,
	mapHashesToRequest map[string]uint32,
) (bool, bool) {
	hasProof := false
	hasHeader := false
	hasRequestedProof := false
	if header, ok := m.mapHeaders[hash]; ok {
		hasHeader = ok
		hasProof = m.hasProof(shardId, []byte(hash), header.GetEpoch())
		if hasProof {
			delete(mapHashesToRequest, hash)
			return false, false
		}
	}

	m.mapHashes[hash] = struct{}{}
	header, ok := m.getHeaderFromPoolOrStorage([]byte(hash))
	if ok {
		hasHeader = ok
		hasProof = m.hasProof(shardId, []byte(hash), header.GetEpoch())
		if hasProof {
			m.mapHeaders[hash] = header
			delete(mapHashesToRequest, hash)
			return false, false
		}
	}

	// if header is missing, do not request the proof
	// if a proof is needed for the header, it will be requested when header is received
	if hasHeader {
		if !hasProof {
			hasRequestedProof = true
			m.requestHandler.RequestEquivalentProofByHash(shardId, []byte(hash))
		}

		return false, hasRequestedProof
	}

	if shardId == core.MetachainShardId {
		m.requestHandler.RequestMetaHeader([]byte(hash))
		return true, hasRequestedProof
	}

	m.requestHandler.RequestShardHeader(shardId, []byte(hash))

	return true, hasRequestedProof
}

func (m *syncHeadersByHash) hasProof(shardID uint32, hash []byte, epoch uint32) bool {
	if !m.enableEpochsHandler.IsFlagEnabledInEpoch(common.AndromedaFlag, epoch) {
		return true
	}

	return m.proofsPool.HasProof(shardID, hash)
}

// receivedHeader is a callback function when a new header was received
// it will further ask for missing transactions
func (m *syncHeadersByHash) receivedHeader(hdrHandler data.HeaderHandler, hdrHash []byte) {
	m.mutMissingHdrs.Lock()
	if m.stopSyncing {
		m.mutMissingHdrs.Unlock()
		return
	}

	if _, ok := m.mapHashes[string(hdrHash)]; !ok {
		m.mutMissingHdrs.Unlock()
		return
	}

	if !m.hasProof(hdrHandler.GetShardID(), hdrHash, hdrHandler.GetEpoch()) {
		go m.requestHandler.RequestEquivalentProofByHash(hdrHandler.GetShardID(), hdrHash)
		m.mutMissingHdrs.Unlock()
		return
	}

	if _, ok := m.mapHeaders[string(hdrHash)]; ok {
		m.mutMissingHdrs.Unlock()
		return
	}

	m.mapHeaders[string(hdrHash)] = hdrHandler
	receivedAll := len(m.mapHashes) == len(m.mapHeaders)
	m.mutMissingHdrs.Unlock()
	if receivedAll {
		m.chReceivedAll <- true
	}
}

func (m *syncHeadersByHash) receivedProof(proofHandler data.HeaderProofHandler) {
	m.mutMissingHdrs.Lock()
	if m.stopSyncing {
		m.mutMissingHdrs.Unlock()
		return
	}

	hdrHash := proofHandler.GetHeaderHash()
	if _, ok := m.mapHashes[string(hdrHash)]; !ok {
		m.mutMissingHdrs.Unlock()
		return
	}

	hdrHandler, ok := m.mapHeaders[string(hdrHash)]
	if !ok {
		hdrHandler, ok = m.getHeaderFromPoolOrStorage(hdrHash)
		if !ok {
			m.mutMissingHdrs.Unlock()
			return
		}
	}

	m.mapHeaders[string(hdrHash)] = hdrHandler
	receivedAll := len(m.mapHashes) == len(m.mapHeaders)
	m.mutMissingHdrs.Unlock()
	if receivedAll {
		m.chReceivedAll <- true
	}
}

func (m *syncHeadersByHash) getHeaderFromPoolOrStorage(hash []byte) (data.HeaderHandler, bool) {
	header, ok := m.getHeaderFromPool(hash)
	if ok {
		return header, true
	}

	hdrData, err := GetDataFromStorage(hash, m.storage)
	if err != nil {
		return nil, false
	}

	hdr, err := process.UnmarshalShardHeader(m.marshalizer, hdrData)
	if err != nil {
		return nil, false
	}

	return hdr, true
}

func (m *syncHeadersByHash) getHeaderFromPool(hash []byte) (data.HeaderHandler, bool) {
	val, err := m.pool.GetHeaderByHash(hash)
	if err != nil {
		return nil, false
	}

	return val, true
}

// GetHeaders returns the synced headers
func (m *syncHeadersByHash) GetHeaders() (map[string]data.HeaderHandler, error) {
	m.mutMissingHdrs.Lock()
	defer m.mutMissingHdrs.Unlock()
	if !m.syncedAll {
		return nil, update.ErrNotSynced
	}

	return m.mapHeaders, nil
}

// ClearFields will clear all the maps
func (m *syncHeadersByHash) ClearFields() {
	m.mutMissingHdrs.Lock()
	m.mapHashes = make(map[string]struct{})
	m.mapHeaders = make(map[string]data.HeaderHandler)
	m.mutMissingHdrs.Unlock()
}

// IsInterfaceNil returns nil if underlying object is nil
func (m *syncHeadersByHash) IsInterfaceNil() bool {
	return m == nil
}
