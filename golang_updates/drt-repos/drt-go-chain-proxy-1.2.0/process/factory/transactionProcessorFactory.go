package factory

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/TerraDharitri/drt-go-chain-proxy/facade"
	"github.com/TerraDharitri/drt-go-chain-proxy/process"
	"github.com/TerraDharitri/drt-go-chain-proxy/process/logsevents"
	"github.com/TerraDharitri/drt-go-chain-proxy/process/txcost"
)

// CreateTransactionProcessor will return the transaction processor needed for current settings
func CreateTransactionProcessor(
	proc process.Processor,
	pubKeyConverter core.PubkeyConverter,
	hasher hashing.Hasher,
	marshalizer marshal.Marshalizer,
	allowEntireTxPoolFetch bool,
) (facade.TransactionProcessor, error) {
	newTxCostProcessor := func() (process.TransactionCostHandler, error) {
		return txcost.NewTransactionCostProcessor(
			proc,
			pubKeyConverter,
		)
	}

	logsMerger, err := logsevents.NewLogsMerger(hasher, &marshal.JsonMarshalizer{})
	if err != nil {
		return nil, err
	}

	return process.NewTransactionProcessor(
		proc,
		pubKeyConverter,
		hasher,
		marshalizer,
		newTxCostProcessor,
		logsMerger,
		allowEntireTxPoolFetch,
	)
}
