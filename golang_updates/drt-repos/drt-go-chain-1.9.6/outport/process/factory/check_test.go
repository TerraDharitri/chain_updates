package factory

import (
	"testing"

	"github.com/TerraDharitri/drt-go-chain/outport/process"
	"github.com/TerraDharitri/drt-go-chain/outport/process/alteredaccounts"
	"github.com/TerraDharitri/drt-go-chain/outport/process/transactionsfee"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	commonMocks "github.com/TerraDharitri/drt-go-chain/testscommon/common"
	"github.com/TerraDharitri/drt-go-chain/testscommon/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/testscommon/economicsmocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/enableEpochsHandlerMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/genericMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/marshallerMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/shardingMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/state"
	"github.com/stretchr/testify/require"
)

func createArgOutportDataProviderFactory() ArgOutportDataProviderFactory {
	return ArgOutportDataProviderFactory{
		HasDrivers:             false,
		AddressConverter:       testscommon.NewPubkeyConverterMock(32),
		AccountsDB:             &state.AccountsStub{},
		Marshaller:             &marshallerMock.MarshalizerMock{},
		DcdtDataStorageHandler: &testscommon.DcdtStorageHandlerStub{},
		TransactionsStorer:     &genericMocks.StorerMock{},
		ShardCoordinator:       &testscommon.ShardsCoordinatorMock{},
		TxCoordinator:          &testscommon.TransactionCoordinatorMock{},
		NodesCoordinator:       &shardingMocks.NodesCoordinatorMock{},
		GasConsumedProvider:    &testscommon.GasHandlerStub{},
		EconomicsData:          &economicsmocks.EconomicsHandlerMock{},
		Hasher:                 &testscommon.KeccakMock{},
		MbsStorer:              &genericMocks.StorerMock{},
		EnableEpochsHandler:    &enableEpochsHandlerMock.EnableEpochsHandlerStub{},
		ExecutionOrderGetter:   &commonMocks.TxExecutionOrderHandlerStub{},
		ProofsPool:             &dataRetriever.ProofsPoolMock{},
	}
}

func TestCheckArgCreateOutportDataProvider(t *testing.T) {
	t.Parallel()

	arg := createArgOutportDataProviderFactory()
	arg.AddressConverter = nil
	require.Equal(t, alteredaccounts.ErrNilPubKeyConverter, checkArgOutportDataProviderFactory(arg))

	arg = createArgOutportDataProviderFactory()
	arg.AccountsDB = nil
	require.Equal(t, alteredaccounts.ErrNilAccountsDB, checkArgOutportDataProviderFactory(arg))

	arg = createArgOutportDataProviderFactory()
	arg.Marshaller = nil
	require.Equal(t, transactionsfee.ErrNilMarshaller, checkArgOutportDataProviderFactory(arg))

	arg = createArgOutportDataProviderFactory()
	arg.DcdtDataStorageHandler = nil
	require.Equal(t, alteredaccounts.ErrNilDCDTDataStorageHandler, checkArgOutportDataProviderFactory(arg))

	arg = createArgOutportDataProviderFactory()
	arg.TransactionsStorer = nil
	require.Equal(t, transactionsfee.ErrNilStorage, checkArgOutportDataProviderFactory(arg))

	arg = createArgOutportDataProviderFactory()
	arg.ShardCoordinator = nil
	require.Equal(t, transactionsfee.ErrNilShardCoordinator, checkArgOutportDataProviderFactory(arg))

	arg = createArgOutportDataProviderFactory()
	arg.TxCoordinator = nil
	require.Equal(t, process.ErrNilTransactionCoordinator, checkArgOutportDataProviderFactory(arg))

	arg = createArgOutportDataProviderFactory()
	arg.NodesCoordinator = nil
	require.Equal(t, process.ErrNilNodesCoordinator, checkArgOutportDataProviderFactory(arg))

	arg = createArgOutportDataProviderFactory()
	arg.GasConsumedProvider = nil
	require.Equal(t, process.ErrNilGasConsumedProvider, checkArgOutportDataProviderFactory(arg))

	arg = createArgOutportDataProviderFactory()
	arg.EconomicsData = nil
	require.Equal(t, transactionsfee.ErrNilTransactionFeeCalculator, checkArgOutportDataProviderFactory(arg))

	arg = createArgOutportDataProviderFactory()
	arg.Hasher = nil
	require.Equal(t, process.ErrNilHasher, checkArgOutportDataProviderFactory(arg))

	arg = createArgOutportDataProviderFactory()
	arg.ProofsPool = nil
	require.Equal(t, process.ErrNilProofsPool, checkArgOutportDataProviderFactory(arg))

	arg = createArgOutportDataProviderFactory()
	require.Nil(t, checkArgOutportDataProviderFactory(arg))
}
