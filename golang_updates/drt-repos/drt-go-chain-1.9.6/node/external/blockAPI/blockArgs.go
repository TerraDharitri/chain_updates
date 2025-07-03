package blockAPI

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	"github.com/TerraDharitri/drt-go-chain-core/data/typeConverters"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/dblookupext"
	outportProcess "github.com/TerraDharitri/drt-go-chain/outport/process"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/state"
)

// ArgAPIBlockProcessor is structure that store components that are needed to create an api block processor
type ArgAPIBlockProcessor struct {
	SelfShardID                  uint32
	Store                        dataRetriever.StorageService
	Marshalizer                  marshal.Marshalizer
	Uint64ByteSliceConverter     typeConverters.Uint64ByteSliceConverter
	HistoryRepo                  dblookupext.HistoryRepository
	APITransactionHandler        APITransactionHandler
	StatusComputer               transaction.StatusComputerHandler
	Hasher                       hashing.Hasher
	AddressPubkeyConverter       core.PubkeyConverter
	LogsFacade                   logsFacade
	ReceiptsRepository           receiptsRepository
	AlteredAccountsProvider      outportProcess.AlteredAccountsProviderHandler
	AccountsRepository           state.AccountsRepository
	ScheduledTxsExecutionHandler process.ScheduledTxsExecutionHandler
	EnableEpochsHandler          common.EnableEpochsHandler
	ProofsPool                   dataRetriever.ProofsPool
	BlockChain                   data.ChainHandler
}
