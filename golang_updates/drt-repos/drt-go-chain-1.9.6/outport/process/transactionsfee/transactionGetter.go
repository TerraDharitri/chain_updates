package transactionsfee

import (
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/TerraDharitri/drt-go-chain/storage"
)

type txGetter struct {
	storer     storage.Storer
	marshaller marshal.Marshalizer
}

func newTxGetter(storer storage.Storer, marshaller marshal.Marshalizer) *txGetter {
	return &txGetter{
		storer:     storer,
		marshaller: marshaller,
	}
}

// GetTxByHash will return from storage transaction with the provided hash
func (tg *txGetter) GetTxByHash(txHash []byte) (*transaction.Transaction, error) {
	txBytes, err := tg.storer.Get(txHash)
	if err != nil {
		return nil, err
	}

	tx := &transaction.Transaction{}
	err = tg.marshaller.Unmarshal(tx, txBytes)
	return tx, err
}
