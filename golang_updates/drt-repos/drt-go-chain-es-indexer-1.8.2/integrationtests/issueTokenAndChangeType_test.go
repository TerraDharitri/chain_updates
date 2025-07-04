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

func TestIssueTokenAndChangeToDynamic(t *testing.T) {
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

	address1 := "drt1k04pxr6c0gvlcx4rd5fje0a4uy33axqxwz0fpcrgtfdy3nrqauqq4s3zqj"
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
							Topics:     [][]byte{[]byte("CON-abcd"), []byte("semi-token"), []byte("SEMI"), []byte(core.SemiFungibleDCDT)},
						},
						{
							Address:    decodeAddress(address1),
							Identifier: []byte("upgradeProperties"),
							Topics:     [][]byte{[]byte("CON-abcd"), big.NewInt(0).Bytes(), []byte("canUpgrade"), []byte("true")},
						},
						nil,
					},
				},
			},
		},
	}

	err = esProc.SaveTransactions(createOutportBlockWithHeader(body, header, pool, map[string]*alteredAccount.AlteredAccount{}, testNumOfShards))
	require.Nil(t, err)

	ids := []string{"CON-abcd"}
	genericResponse := &GenericResponse{}
	err = esClient.DoMultiGet(context.Background(), ids, indexdrtata.TokensIndex, true, genericResponse)
	require.Nil(t, err)
	require.JSONEq(t, readExpectedResult("./testdata/issueTokenAndChangeType/token.json"), string(genericResponse.Docs[0].Source))

	header.TimeStamp = 10_000
	pool = &outport.TransactionPool{
		Logs: []*outport.LogData{
			{
				TxHash: hex.EncodeToString([]byte("h1")),
				Log: &transaction.Log{
					Address: decodeAddress(address1),
					Events: []*transaction.Event{
						{
							Address:    decodeAddress(address1),
							Identifier: []byte("changeToDynamic"),
							Topics:     [][]byte{[]byte("CON-abcd"), []byte("semi-token"), []byte("SEMI"), []byte(core.DynamicSFTDCDT)},
						},
						nil,
					},
				},
			},
		},
	}

	err = esProc.SaveTransactions(createOutportBlockWithHeader(body, header, pool, map[string]*alteredAccount.AlteredAccount{}, testNumOfShards))
	require.Nil(t, err)

	err = esClient.DoMultiGet(context.Background(), ids, indexdrtata.TokensIndex, true, genericResponse)
	require.Nil(t, err)
	require.JSONEq(t, readExpectedResult("./testdata/issueTokenAndChangeType/token-with-new-type.json"), string(genericResponse.Docs[0].Source))
}
