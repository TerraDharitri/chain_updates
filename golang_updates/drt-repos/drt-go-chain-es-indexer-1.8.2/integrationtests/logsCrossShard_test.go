//go:build integrationtests

package integrationtests

import (
	"context"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data/alteredAccount"
	dataBlock "github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain-core/data/outport"
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	indexdrtata "github.com/TerraDharitri/drt-go-chain-es-indexer/process/dataindexer"
	"github.com/stretchr/testify/require"
)

func TestIndexLogSourceShardAndAfterDestinationAndAgainSource(t *testing.T) {
	setLogLevelDebug()

	esClient, err := createESClient(esURL)
	require.Nil(t, err)

	esProc, err := CreateElasticProcessor(esClient)
	require.Nil(t, err)

	header := &dataBlock.Header{
		Round:     50,
		TimeStamp: 5040,
	}

	txHash := []byte("cross-log")
	logID := hex.EncodeToString(txHash)

	body := &dataBlock.Body{
		MiniBlocks: []*dataBlock.MiniBlock{
			{
				TxHashes: [][]byte{txHash},
			},
		},
	}

	address1 := "drt1ju8pkvg57cwdmjsjx58jlmnuf4l9yspstrhr9tgsrt98n9edpm2qkrl8xm"
	address2 := "drt1w7jyzuj6cv4ngw8luhlkakatjpmjh3ql95lmxphd3vssc4vpymks82rg5q"

	// index on source
	pool := &outport.TransactionPool{
		Logs: []*outport.LogData{
			{
				TxHash: logID,
				Log: &transaction.Log{
					Address: decodeAddress(address1),
					Events: []*transaction.Event{
						{
							Address:    decodeAddress(address1),
							Identifier: []byte(core.BuiltInFunctionDCDTTransfer),
							Topics:     [][]byte{[]byte("DCDT-abcd"), big.NewInt(0).Bytes(), big.NewInt(1).Bytes()},
						},
						nil,
					},
				},
			},
		},
		Transactions: map[string]*outport.TxInfo{
			logID: {
				Transaction:    &transaction.Transaction{},
				ExecutionOrder: 0,
			},
		},
	}
	err = esProc.SaveTransactions(createOutportBlockWithHeader(body, header, pool, map[string]*alteredAccount.AlteredAccount{}, testNumOfShards))
	require.Nil(t, err)

	ids := []string{logID}
	genericResponse := &GenericResponse{}
	err = esClient.DoMultiGet(context.Background(), ids, indexdrtata.LogsIndex, true, genericResponse)
	require.Nil(t, err)
	require.JSONEq(t,
		readExpectedResult("./testdata/logsCrossShard/log-at-source.json"),
		string(genericResponse.Docs[0].Source),
	)

	event1ID := logID + "-0-0"
	ids = []string{event1ID}
	err = esClient.DoMultiGet(context.Background(), ids, indexdrtata.EventsIndex, true, genericResponse)
	require.Nil(t, err)
	require.JSONEq(t,
		readExpectedResult("./testdata/logsCrossShard/event-transfer-source-first.json"),
		string(genericResponse.Docs[0].Source),
	)

	// INDEX ON DESTINATION
	header = &dataBlock.Header{
		Round:     50,
		TimeStamp: 6040,
		ShardID:   1,
	}
	pool = &outport.TransactionPool{
		Logs: []*outport.LogData{
			{
				TxHash: logID,
				Log: &transaction.Log{
					Address: decodeAddress(address1),
					Events: []*transaction.Event{
						{
							Address:    decodeAddress(address1),
							Identifier: []byte(core.BuiltInFunctionDCDTTransfer),
							Topics:     [][]byte{[]byte("DCDT-abcd"), big.NewInt(0).Bytes(), big.NewInt(1).Bytes()},
						},
						{

							Address:    decodeAddress(address2),
							Identifier: []byte("do-something"),
							Topics:     [][]byte{[]byte("topic1"), []byte("topic2")},
						},
						nil,
					},
				},
			},
		},
		Transactions: map[string]*outport.TxInfo{
			logID: {
				Transaction:    &transaction.Transaction{},
				ExecutionOrder: 0,
			},
		},
	}
	err = esProc.SaveTransactions(createOutportBlockWithHeader(body, header, pool, map[string]*alteredAccount.AlteredAccount{}, testNumOfShards))
	require.Nil(t, err)

	ids = []string{logID}
	err = esClient.DoMultiGet(context.Background(), ids, indexdrtata.LogsIndex, true, genericResponse)
	require.Nil(t, err)
	require.JSONEq(t,
		readExpectedResult("./testdata/logsCrossShard/log-at-destination.json"),
		string(genericResponse.Docs[0].Source),
	)

	event2ID, event3ID := logID+"-1-0", logID+"-1-1"
	ids = []string{event2ID, event3ID}
	err = esClient.DoMultiGet(context.Background(), ids, indexdrtata.EventsIndex, true, genericResponse)
	require.Nil(t, err)
	require.JSONEq(t,
		readExpectedResult("./testdata/logsCrossShard/event-transfer-destination.json"),
		string(genericResponse.Docs[0].Source),
	)
	require.JSONEq(t,
		readExpectedResult("./testdata/logsCrossShard/event-do-something.json"),
		string(genericResponse.Docs[1].Source),
	)

	// index on source again should not change the log
	header = &dataBlock.Header{
		Round:     50,
		TimeStamp: 5000,
	}
	pool = &outport.TransactionPool{
		Logs: []*outport.LogData{
			{
				TxHash: logID,
				Log: &transaction.Log{
					Address: decodeAddress(address1),
					Events: []*transaction.Event{
						{
							Address:    decodeAddress(address1),
							Identifier: []byte(core.BuiltInFunctionDCDTTransfer),
							Topics:     [][]byte{[]byte("DCDT-abcd"), big.NewInt(0).Bytes(), big.NewInt(1).Bytes()},
						},
						nil,
					},
				},
			},
		},
		Transactions: map[string]*outport.TxInfo{
			logID: {
				Transaction:    &transaction.Transaction{},
				ExecutionOrder: 0,
			},
		},
	}
	err = esProc.SaveTransactions(createOutportBlockWithHeader(body, header, pool, map[string]*alteredAccount.AlteredAccount{}, testNumOfShards))
	require.Nil(t, err)

	ids = []string{logID}
	err = esClient.DoMultiGet(context.Background(), ids, indexdrtata.LogsIndex, true, genericResponse)
	require.Nil(t, err)
	require.JSONEq(t,
		readExpectedResult("./testdata/logsCrossShard/log-at-destination.json"),
		string(genericResponse.Docs[0].Source),
	)

	// do rollback
	header = &dataBlock.Header{
		Round:     50,
		TimeStamp: 6040,
		MiniBlockHeaders: []dataBlock.MiniBlockHeader{
			{},
		},
		ShardID: 1,
	}
	body = &dataBlock.Body{
		MiniBlocks: []*dataBlock.MiniBlock{
			{
				TxHashes: [][]byte{[]byte("cross-log")},
			},
		},
	}

	err = esProc.RemoveTransactions(header, body)
	require.Nil(t, err)

	err = esClient.DoMultiGet(context.Background(), ids, indexdrtata.LogsIndex, true, genericResponse)
	require.Nil(t, err)

	require.False(t, genericResponse.Docs[0].Found)

	ids = []string{event2ID, event3ID}
	err = esClient.DoMultiGet(context.Background(), ids, indexdrtata.EventsIndex, true, genericResponse)
	require.Nil(t, err)

	require.False(t, genericResponse.Docs[0].Found)
	require.False(t, genericResponse.Docs[1].Found)
}
