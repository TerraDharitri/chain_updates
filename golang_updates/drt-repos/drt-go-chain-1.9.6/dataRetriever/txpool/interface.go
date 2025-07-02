package txpool

import (
	"math/big"

	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain/storage"
	"github.com/TerraDharitri/drt-go-chain/storage/txcache"
)

type txCache interface {
	storage.Cacher

	AddTx(tx *txcache.WrappedTransaction) (ok bool, added bool)
	GetByTxHash(txHash []byte) (*txcache.WrappedTransaction, bool)
	RemoveTxByHash(txHash []byte) bool
	ImmunizeTxsAgainstEviction(keys [][]byte)
	ForEachTransaction(function txcache.ForEachTransaction)
	NumBytes() int
	Diagnose(deep bool)
	GetTransactionsPoolForSender(sender string) []*txcache.WrappedTransaction
}

type txGasHandler interface {
	ComputeTxFee(tx data.TransactionWithFeeHandler) *big.Int
	IsInterfaceNil() bool
}
