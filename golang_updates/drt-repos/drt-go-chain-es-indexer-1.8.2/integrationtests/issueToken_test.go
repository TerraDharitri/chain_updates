//go:build integrationtests

package integrationtests

import (
	"context"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	dataBlock "github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain-core/data/outport"
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	indexdrtata "github.com/TerraDharitri/drt-go-chain-es-indexer/process/dataindexer"
	"github.com/stretchr/testify/require"
)

func TestIssueTokenAndTransferOwnership(t *testing.T) {
	setLogLevelDebug()

	esClient, err := createESClient(esURL)
	require.Nil(t, err)

	esProc, err := CreateElasticProcessor(esClient)
	require.Nil(t, err)

	body := &dataBlock.Body{}
	header := &dataBlock.Header{
		Round:     50,
		TimeStamp: 5040,
		ShardID:   core.MetachainShardId,
	}

	address1 := "drt1v7e552pz9py4hv6raan0c4jflez3e6csdmzcgrncg0qrnk4tywvsa6c5hv"
	address2 := "drt1acjlnuhkd8773sqhmw85r0ur4lcyuqgm0n69h9ttxh0gwxtuuzxqgr045y"
	pool := &outport.TransactionPool{
		Logs: []*outport.LogData{
			{
				TxHash: hex.EncodeToString([]byte("h1")),
				Log: &transaction.Log{
					Address: decodeAddress(address1),
					Events: []*transaction.Event{
						{
							Address:    decodeAddress(address1),
							Identifier: []byte("issueSemiFungible"),
							Topics:     [][]byte{[]byte("SSSS-abcd"), []byte("semi-token"), []byte("SSSS"), []byte(core.SemiFungibleDCDT), big.NewInt(18).Bytes()},
						},
						nil,
					},
				},
			},
		},
	}

	err = esProc.SaveTransactions(createOutportBlockWithHeader(body, header, pool, nil, testNumOfShards))
	require.Nil(t, err)

	ids := []string{"SSSS-abcd"}
	genericResponse := &GenericResponse{}
	err = esClient.DoMultiGet(context.Background(), ids, indexdrtata.TokensIndex, true, genericResponse)
	require.Nil(t, err)
	require.JSONEq(t, readExpectedResult("./testdata/issueToken/token-semi.json"), string(genericResponse.Docs[0].Source))

	err = esClient.DoMultiGet(context.Background(), ids, indexdrtata.DCDTsIndex, true, genericResponse)
	require.Nil(t, err)
	require.JSONEq(t, readExpectedResult("./testdata/issueToken/token-semi.json"), string(genericResponse.Docs[0].Source))

	// transfer ownership
	pool = &outport.TransactionPool{
		Logs: []*outport.LogData{
			{
				TxHash: hex.EncodeToString([]byte("h1")),
				Log: &transaction.Log{
					Address: decodeAddress(address1),
					Events: []*transaction.Event{
						{
							Address:    decodeAddress(address1),
							Identifier: []byte("transferOwnership"),
							Topics:     [][]byte{[]byte("SSSS-abcd"), []byte("semi-token"), []byte("SSSS"), []byte(core.SemiFungibleDCDT), decodeAddress(address2)},
						},
						nil,
					},
				},
			},
		},
	}

	header.TimeStamp = 10000
	err = esProc.SaveTransactions(createOutportBlockWithHeader(body, header, pool, nil, testNumOfShards))
	require.Nil(t, err)

	err = esClient.DoMultiGet(context.Background(), ids, indexdrtata.TokensIndex, true, genericResponse)
	require.Nil(t, err)
	require.JSONEq(t, readExpectedResult("./testdata/issueToken/token-semi-after-transfer-ownership.json"), string(genericResponse.Docs[0].Source))

	err = esClient.DoMultiGet(context.Background(), ids, indexdrtata.DCDTsIndex, true, genericResponse)
	require.Nil(t, err)
	require.JSONEq(t, readExpectedResult("./testdata/issueToken/token-semi-after-transfer-ownership.json"), string(genericResponse.Docs[0].Source))

	// do pause
	pool = &outport.TransactionPool{
		Logs: []*outport.LogData{
			{
				TxHash: hex.EncodeToString([]byte("h1")),
				Log: &transaction.Log{
					Address: decodeAddress(address1),
					Events: []*transaction.Event{
						{
							Address:    decodeAddress(address1),
							Identifier: []byte("DCDTPause"),
							Topics:     [][]byte{[]byte("SSSS-abcd")},
						},
						nil,
					},
				},
			},
		},
	}

	header.TimeStamp = 10000
	err = esProc.SaveTransactions(createOutportBlockWithHeader(body, header, pool, nil, testNumOfShards))
	require.Nil(t, err)

	err = esClient.DoMultiGet(context.Background(), ids, indexdrtata.TokensIndex, true, genericResponse)
	require.Nil(t, err)
	require.JSONEq(t, readExpectedResult("./testdata/issueToken/token-semi-after-pause.json"), string(genericResponse.Docs[0].Source))

	// do unPause
	pool = &outport.TransactionPool{
		Logs: []*outport.LogData{
			{
				TxHash: hex.EncodeToString([]byte("h1")),
				Log: &transaction.Log{
					Address: decodeAddress(address1),
					Events: []*transaction.Event{
						{
							Address:    decodeAddress(address1),
							Identifier: []byte("DCDTUnPause"),
							Topics:     [][]byte{[]byte("SSSS-abcd")},
						},
						nil,
					},
				},
			},
		},
	}

	header.TimeStamp = 10000
	err = esProc.SaveTransactions(createOutportBlockWithHeader(body, header, pool, nil, testNumOfShards))
	require.Nil(t, err)

	err = esClient.DoMultiGet(context.Background(), ids, indexdrtata.TokensIndex, true, genericResponse)
	require.Nil(t, err)
	require.JSONEq(t, readExpectedResult("./testdata/issueToken/token-semi-after-un-pause.json"), string(genericResponse.Docs[0].Source))
}
