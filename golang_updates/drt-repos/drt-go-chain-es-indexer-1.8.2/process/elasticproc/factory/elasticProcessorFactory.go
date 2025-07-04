package factory

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/TerraDharitri/drt-go-chain-es-indexer/process/dataindexer"
	"github.com/TerraDharitri/drt-go-chain-es-indexer/process/elasticproc"
	"github.com/TerraDharitri/drt-go-chain-es-indexer/process/elasticproc/accounts"
	blockProc "github.com/TerraDharitri/drt-go-chain-es-indexer/process/elasticproc/block"
	"github.com/TerraDharitri/drt-go-chain-es-indexer/process/elasticproc/converters"
	"github.com/TerraDharitri/drt-go-chain-es-indexer/process/elasticproc/logsevents"
	"github.com/TerraDharitri/drt-go-chain-es-indexer/process/elasticproc/miniblocks"
	"github.com/TerraDharitri/drt-go-chain-es-indexer/process/elasticproc/operations"
	"github.com/TerraDharitri/drt-go-chain-es-indexer/process/elasticproc/statistics"
	"github.com/TerraDharitri/drt-go-chain-es-indexer/process/elasticproc/templatesAndPolicies"
	"github.com/TerraDharitri/drt-go-chain-es-indexer/process/elasticproc/transactions"
	"github.com/TerraDharitri/drt-go-chain-es-indexer/process/elasticproc/validators"
)

// ArgElasticProcessorFactory is struct that is used to store all components that are needed to create an elastic processor factory
type ArgElasticProcessorFactory struct {
	Marshalizer              marshal.Marshalizer
	Hasher                   hashing.Hasher
	AddressPubkeyConverter   core.PubkeyConverter
	ValidatorPubkeyConverter core.PubkeyConverter
	DBClient                 elasticproc.DatabaseClientHandler
	EnabledIndexes           []string
	Version                  string
	Denomination             int
	BulkRequestMaxSize       int
	UseKibana                bool
	ImportDB                 bool
}

// CreateElasticProcessor will create a new instance of ElasticProcessor
func CreateElasticProcessor(arguments ArgElasticProcessorFactory) (dataindexer.ElasticProcessor, error) {
	templatesAndPoliciesReader := templatesAndPolicies.CreateTemplatesAndPoliciesReader(arguments.UseKibana)
	indexTemplates, indexPolicies, err := templatesAndPoliciesReader.GetElasticTemplatesAndPolicies()
	if err != nil {
		return nil, err
	}
	extraMappings, err := templatesAndPoliciesReader.GetExtraMappings()
	if err != nil {
		return nil, err
	}

	enabledIndexesMap := make(map[string]struct{})
	for _, index := range arguments.EnabledIndexes {
		enabledIndexesMap[index] = struct{}{}
	}
	if len(enabledIndexesMap) == 0 {
		return nil, dataindexer.ErrEmptyEnabledIndexes
	}

	balanceConverter, err := converters.NewBalanceConverter(arguments.Denomination)
	if err != nil {
		return nil, err
	}

	accountsProc, err := accounts.NewAccountsProcessor(
		arguments.AddressPubkeyConverter,
		balanceConverter,
	)
	if err != nil {
		return nil, err
	}

	blockProcHandler, err := blockProc.NewBlockProcessor(arguments.Hasher, arguments.Marshalizer, arguments.ValidatorPubkeyConverter)
	if err != nil {
		return nil, err
	}

	miniblocksProc, err := miniblocks.NewMiniblocksProcessor(arguments.Hasher, arguments.Marshalizer)
	if err != nil {
		return nil, err
	}
	validatorsProc, err := validators.NewValidatorsProcessor(arguments.ValidatorPubkeyConverter, arguments.BulkRequestMaxSize)
	if err != nil {
		return nil, err
	}

	generalInfoProc := statistics.NewStatisticsProcessor()

	argsTxsProc := &transactions.ArgsTransactionProcessor{
		AddressPubkeyConverter: arguments.AddressPubkeyConverter,
		Hasher:                 arguments.Hasher,
		Marshalizer:            arguments.Marshalizer,
		BalanceConverter:       balanceConverter,
	}
	txsProc, err := transactions.NewTransactionsProcessor(argsTxsProc)
	if err != nil {
		return nil, err
	}

	argsLogsAndEventsProc := logsevents.ArgsLogsAndEventsProcessor{
		PubKeyConverter:  arguments.AddressPubkeyConverter,
		Marshalizer:      arguments.Marshalizer,
		BalanceConverter: balanceConverter,
		Hasher:           arguments.Hasher,
	}
	logsAndEventsProc, err := logsevents.NewLogsAndEventsProcessor(argsLogsAndEventsProc)
	if err != nil {
		return nil, err
	}

	operationsProc, err := operations.NewOperationsProcessor()
	if err != nil {
		return nil, err
	}

	args := &elasticproc.ArgElasticProcessor{
		BulkRequestMaxSize: arguments.BulkRequestMaxSize,
		TransactionsProc:   txsProc,
		AccountsProc:       accountsProc,
		BlockProc:          blockProcHandler,
		MiniblocksProc:     miniblocksProc,
		ValidatorsProc:     validatorsProc,
		StatisticsProc:     generalInfoProc,
		LogsAndEventsProc:  logsAndEventsProc,
		DBClient:           arguments.DBClient,
		EnabledIndexes:     enabledIndexesMap,
		UseKibana:          arguments.UseKibana,
		IndexTemplates:     indexTemplates,
		IndexPolicies:      indexPolicies,
		ExtraMappings:      extraMappings,
		OperationsProc:     operationsProc,
		ImportDB:           arguments.ImportDB,
		Version:            arguments.Version,
	}

	return elasticproc.NewElasticProcessor(args)
}
