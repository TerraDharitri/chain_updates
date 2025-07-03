package transactionAPI

import (
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data/typeConverters"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/dblookupext"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/sharding"
)

// ArgAPITransactionProcessor is structure that store components that are needed to create an api transaction processor
type ArgAPITransactionProcessor struct {
	RoundDuration            uint64
	GenesisTime              time.Time
	Marshalizer              marshal.Marshalizer
	AddressPubKeyConverter   core.PubkeyConverter
	ShardCoordinator         sharding.Coordinator
	HistoryRepository        dblookupext.HistoryRepository
	StorageService           dataRetriever.StorageService
	DataPool                 dataRetriever.PoolsHolder
	Uint64ByteSliceConverter typeConverters.Uint64ByteSliceConverter
	FeeComputer              feeComputer
	TxTypeHandler            process.TxTypeHandler
	LogsFacade               LogsFacade
	DataFieldParser          DataFieldParser
	TxMarshaller             marshal.Marshalizer
	EnableEpochsHandler      common.EnableEpochsHandler
}
