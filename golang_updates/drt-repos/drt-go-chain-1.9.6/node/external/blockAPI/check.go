package blockAPI

import (
	"errors"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/process"
)

var (
	errNilArgAPIBlockProcessor         = errors.New("nil arg api block processor")
	errNilTransactionHandler           = errors.New("nil API transaction handler")
	errNilLogsFacade                   = errors.New("nil logs facade")
	errNilReceiptsRepository           = errors.New("nil receipts repository")
	errNilAlteredAccountsProvider      = errors.New("nil altered accounts provider")
	errNilAccountsRepository           = errors.New("nil accounts repository")
	errNilScheduledTxsExecutionHandler = errors.New("nil scheduled txs execution handler")
	errNilEnableEpochsHandler          = errors.New("nil enable epochs handler")
)

func checkNilArg(arg *ArgAPIBlockProcessor) error {
	if arg == nil {
		return errNilArgAPIBlockProcessor
	}
	if check.IfNil(arg.Uint64ByteSliceConverter) {
		return process.ErrNilUint64Converter
	}
	if check.IfNil(arg.Store) {
		return process.ErrNilStorage
	}
	if check.IfNil(arg.Marshalizer) {
		return process.ErrNilMarshalizer
	}
	if check.IfNil(arg.HistoryRepo) {
		return process.ErrNilHistoryRepository
	}
	if check.IfNil(arg.APITransactionHandler) {
		return errNilTransactionHandler
	}
	if check.IfNil(arg.Hasher) {
		return process.ErrNilHasher
	}
	if check.IfNil(arg.AddressPubkeyConverter) {
		return process.ErrNilPubkeyConverter
	}
	if check.IfNil(arg.LogsFacade) {
		return errNilLogsFacade
	}
	if check.IfNil(arg.ReceiptsRepository) {
		return errNilReceiptsRepository
	}
	if check.IfNil(arg.AlteredAccountsProvider) {
		return errNilAlteredAccountsProvider
	}
	if check.IfNil(arg.AccountsRepository) {
		return errNilAccountsRepository
	}
	if check.IfNil(arg.ScheduledTxsExecutionHandler) {
		return errNilScheduledTxsExecutionHandler
	}
	if check.IfNil(arg.EnableEpochsHandler) {
		return errNilEnableEpochsHandler
	}
	if check.IfNil(arg.ProofsPool) {
		return process.ErrNilProofsPool
	}
	if check.IfNil(arg.BlockChain) {
		return process.ErrNilBlockChain
	}

	return core.CheckHandlerCompatibility(arg.EnableEpochsHandler, []core.EnableEpochFlag{
		common.RefactorPeersMiniBlocksFlag,
	})
}
