package interceptedBlocks

import (
	"fmt"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	logger "github.com/TerraDharitri/drt-go-chain-logger"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/sharding"
)

var _ process.HdrValidatorHandler = (*InterceptedMetaHeader)(nil)
var _ process.InterceptedData = (*InterceptedMetaHeader)(nil)

var log = logger.GetOrCreate("process/block/interceptedBlocks")

// InterceptedMetaHeader represents the wrapper over the meta block header struct
type InterceptedMetaHeader struct {
	hdr                 data.MetaHeaderHandler
	sigVerifier         process.InterceptedHeaderSigVerifier
	integrityVerifier   process.HeaderIntegrityVerifier
	hasher              hashing.Hasher
	shardCoordinator    sharding.Coordinator
	hash                []byte
	validityAttester    process.ValidityAttester
	epochStartTrigger   process.EpochStartTriggerHandler
	enableEpochsHandler common.EnableEpochsHandler
}

// NewInterceptedMetaHeader creates a new instance of InterceptedMetaHeader struct
func NewInterceptedMetaHeader(arg *ArgInterceptedBlockHeader) (*InterceptedMetaHeader, error) {
	err := checkBlockHeaderArgument(arg)
	if err != nil {
		return nil, err
	}

	hdr, err := createMetaHdr(arg.Marshalizer, arg.HdrBuff)
	if err != nil {
		return nil, err
	}

	inHdr := &InterceptedMetaHeader{
		hdr:                 hdr,
		hasher:              arg.Hasher,
		sigVerifier:         arg.HeaderSigVerifier,
		integrityVerifier:   arg.HeaderIntegrityVerifier,
		shardCoordinator:    arg.ShardCoordinator,
		validityAttester:    arg.ValidityAttester,
		epochStartTrigger:   arg.EpochStartTrigger,
		enableEpochsHandler: arg.EnableEpochsHandler,
	}
	inHdr.processFields(arg.HdrBuff)

	return inHdr, nil
}

func createMetaHdr(marshalizer marshal.Marshalizer, hdrBuff []byte) (*block.MetaBlock, error) {
	hdr := &block.MetaBlock{
		ShardInfo: make([]block.ShardData, 0),
	}
	err := marshalizer.Unmarshal(hdr, hdrBuff)
	if err != nil {
		return nil, err
	}

	return hdr, nil
}

func (imh *InterceptedMetaHeader) processFields(txBuff []byte) {
	imh.hash = imh.hasher.Compute(string(txBuff))
}

// Hash gets the hash of this header
func (imh *InterceptedMetaHeader) Hash() []byte {
	return imh.hash
}

// HeaderHandler returns the MetaBlock pointer that holds the data
func (imh *InterceptedMetaHeader) HeaderHandler() data.HeaderHandler {
	return imh.hdr
}

// CheckValidity checks if the received meta header is valid (not nil fields, valid sig and so on)
func (imh *InterceptedMetaHeader) CheckValidity() error {
	log.Trace("CheckValidity for header with", "epoch", imh.hdr.GetEpoch(), "hash", logger.DisplayByteSlice(imh.hash))

	err := imh.integrity()
	if err != nil {
		return err
	}

	if !imh.validityAttester.CheckBlockAgainstWhitelist(imh) {
		err = imh.validityAttester.CheckBlockAgainstFinal(imh.HeaderHandler())
		if err != nil {
			return err
		}

		if imh.isMetaHeaderEpochOutOfRange() {
			log.Trace("InterceptedMetaHeader.CheckValidity",
				"trigger epoch", imh.epochStartTrigger.Epoch(),
				"metaBlock epoch", imh.hdr.GetEpoch(),
				"error", process.ErrMetaHeaderEpochOutOfRange)

			return process.ErrMetaHeaderEpochOutOfRange
		}
	}

	err = imh.validityAttester.CheckBlockAgainstRoundHandler(imh.HeaderHandler())
	if err != nil {
		return err
	}

	err = imh.sigVerifier.VerifyRandSeedAndLeaderSignature(imh.hdr)
	if err != nil {
		return err
	}

	err = imh.sigVerifier.VerifySignature(imh.hdr)
	if err != nil {
		return err
	}

	return imh.integrityVerifier.Verify(imh.hdr)
}

func (imh *InterceptedMetaHeader) isMetaHeaderEpochOutOfRange() bool {
	if imh.shardCoordinator.SelfId() == core.MetachainShardId {
		return false
	}

	if imh.hdr.GetEpoch() > imh.epochStartTrigger.Epoch()+1 {
		return true
	}

	return false
}

// integrity checks the integrity of the meta header block wrapper
func (imh *InterceptedMetaHeader) integrity() error {
	err := checkHeaderHandler(imh.HeaderHandler(), imh.enableEpochsHandler)
	if err != nil {
		return err
	}

	err = checkMetaShardInfo(imh.hdr.GetShardInfoHandlers(), imh.shardCoordinator)
	if err != nil {
		return err
	}

	return nil
}

// IsForCurrentShard always returns true
func (imh *InterceptedMetaHeader) IsForCurrentShard() bool {
	return true
}

// Type returns the type of this intercepted data
func (imh *InterceptedMetaHeader) Type() string {
	return "intercepted meta header"
}

// String returns the meta header's most important fields as string
func (imh *InterceptedMetaHeader) String() string {
	return fmt.Sprintf("epoch=%d, round=%d, nonce=%d",
		imh.hdr.GetEpoch(),
		imh.hdr.GetRound(),
		imh.hdr.GetNonce(),
	)
}

// Identifiers returns the identifiers used in requests
func (imh *InterceptedMetaHeader) Identifiers() [][]byte {
	keyNonce := []byte(fmt.Sprintf("%d-%d", core.MetachainShardId, imh.hdr.GetNonce()))
	keyEpoch := []byte(core.EpochStartIdentifier(imh.hdr.GetEpoch()))

	return [][]byte{imh.hash, keyNonce, keyEpoch}
}

// IsInterfaceNil returns true if there is no value under the interface
func (imh *InterceptedMetaHeader) IsInterfaceNil() bool {
	return imh == nil
}
