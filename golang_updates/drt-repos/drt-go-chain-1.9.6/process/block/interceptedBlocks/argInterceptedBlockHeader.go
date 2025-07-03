package interceptedBlocks

import (
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/sharding"
)

// ArgInterceptedBlockHeader is the argument for the intercepted header
type ArgInterceptedBlockHeader struct {
	HdrBuff                       []byte
	Marshalizer                   marshal.Marshalizer
	Hasher                        hashing.Hasher
	ShardCoordinator              sharding.Coordinator
	HeaderSigVerifier             process.InterceptedHeaderSigVerifier
	HeaderIntegrityVerifier       process.HeaderIntegrityVerifier
	ValidityAttester              process.ValidityAttester
	EpochStartTrigger             process.EpochStartTriggerHandler
	EnableEpochsHandler           common.EnableEpochsHandler
	EpochChangeGracePeriodHandler common.EpochChangeGracePeriodHandler
}
