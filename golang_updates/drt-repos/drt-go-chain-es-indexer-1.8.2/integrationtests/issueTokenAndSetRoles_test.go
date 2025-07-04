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

func TestIssueTokenAndSetRole(t *testing.T) {
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
	address2 := "drt1suhxyflu4w4pqdxmushpxzc6a3qszr89m8uswzqcvyh0mh9mzxwqsjpcvc"
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
							Topics:     [][]byte{[]byte("TOK-abcd"), []byte("semi-token"), []byte("SEMI"), []byte(core.SemiFungibleDCDT)},
						},
						{
							Address:    decodeAddress(address1),
							Identifier: []byte("upgradeProperties"),
							Topics:     [][]byte{[]byte("TOK-abcd"), big.NewInt(0).Bytes(), []byte("canUpgrade"), []byte("true")},
						},
						nil,
					},
				},
			},
		},
	}

	err = esProc.SaveTransactions(createOutportBlockWithHeader(body, header, pool, map[string]*alteredAccount.AlteredAccount{}, testNumOfShards))
	require.Nil(t, err)

	ids := []string{"TOK-abcd"}
	genericResponse := &GenericResponse{}
	err = esClient.DoMultiGet(context.Background(), ids, indexdrtata.TokensIndex, true, genericResponse)
	require.Nil(t, err)
	require.JSONEq(t, readExpectedResult("./testdata/issueTokenAndSetRoles/token-after-issue-ok.json"), string(genericResponse.Docs[0].Source))

	// SET ROLES
	pool = &outport.TransactionPool{
		Logs: []*outport.LogData{
			{
				TxHash: hex.EncodeToString([]byte("h1")),
				Log: &transaction.Log{
					Address: decodeAddress(address1),
					Events: []*transaction.Event{
						{
							Address:    decodeAddress(address1),
							Identifier: []byte(core.BuiltInFunctionSetDCDTRole),
							Topics:     [][]byte{[]byte("TOK-abcd"), big.NewInt(0).Bytes(), big.NewInt(0).Bytes(), []byte(core.DCDTRoleNFTCreate), []byte(core.DCDTRoleNFTBurn)},
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

	ids = []string{"TOK-abcd"}
	genericResponse = &GenericResponse{}
	err = esClient.DoMultiGet(context.Background(), ids, indexdrtata.TokensIndex, true, genericResponse)
	require.Nil(t, err)
	require.JSONEq(t, readExpectedResult("./testdata/issueTokenAndSetRoles/token-after-set-role.json"), string(genericResponse.Docs[0].Source))

	// TRANSFER ROLE
	pool = &outport.TransactionPool{
		Logs: []*outport.LogData{
			{
				TxHash: hex.EncodeToString([]byte("h1")),
				Log: &transaction.Log{
					Address: decodeAddress(address1),
					Events: []*transaction.Event{
						{
							Address:    decodeAddress(address1),
							Identifier: []byte(core.BuiltInFunctionDCDTNFTCreateRoleTransfer),
							Topics:     [][]byte{[]byte("TOK-abcd"), big.NewInt(0).Bytes(), big.NewInt(0).Bytes(), []byte("false")},
						},
						{
							Address:    decodeAddress(address2),
							Identifier: []byte(core.BuiltInFunctionDCDTNFTCreateRoleTransfer),
							Topics:     [][]byte{[]byte("TOK-abcd"), big.NewInt(0).Bytes(), big.NewInt(0).Bytes(), []byte("true")},
						},
					},
				},
			},
		},
	}

	header.TimeStamp = 10000
	err = esProc.SaveTransactions(createOutportBlockWithHeader(body, header, pool, nil, testNumOfShards))
	require.Nil(t, err)

	ids = []string{"TOK-abcd"}
	genericResponse = &GenericResponse{}
	err = esClient.DoMultiGet(context.Background(), ids, indexdrtata.TokensIndex, true, genericResponse)
	require.Nil(t, err)
	require.JSONEq(t, readExpectedResult("./testdata/issueTokenAndSetRoles/token-after-transfer-role.json"), string(genericResponse.Docs[0].Source))

	// UNSET ROLES
	pool = &outport.TransactionPool{
		Logs: []*outport.LogData{
			{
				TxHash: hex.EncodeToString([]byte("h1")),
				Log: &transaction.Log{
					Address: decodeAddress(address1),
					Events: []*transaction.Event{
						{
							Address:    decodeAddress(address1),
							Identifier: []byte(core.BuiltInFunctionUnSetDCDTRole),
							Topics:     [][]byte{[]byte("TOK-abcd"), big.NewInt(0).Bytes(), big.NewInt(0).Bytes(), []byte(core.DCDTRoleNFTBurn)},
						},
						nil,
					},
				},
			},
		},
	}

	header.TimeStamp = 10000
	err = esProc.SaveTransactions(createOutportBlockWithHeader(body, header, pool, map[string]*alteredAccount.AlteredAccount{}, testNumOfShards))
	require.Nil(t, err)

	ids = []string{"TOK-abcd"}
	genericResponse = &GenericResponse{}
	err = esClient.DoMultiGet(context.Background(), ids, indexdrtata.TokensIndex, true, genericResponse)
	require.Nil(t, err)
	require.JSONEq(t, readExpectedResult("./testdata/issueTokenAndSetRoles/token-after-unset-role.json"), string(genericResponse.Docs[0].Source))
}

func TestIssueSetRolesEventAndAfterTokenIssue(t *testing.T) {
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
	// SET ROLES
	pool := &outport.TransactionPool{
		Logs: []*outport.LogData{
			{
				TxHash: hex.EncodeToString([]byte("h1")),
				Log: &transaction.Log{
					Address: decodeAddress(address1),
					Events: []*transaction.Event{
						{
							Address:    decodeAddress(address1),
							Identifier: []byte(core.BuiltInFunctionSetDCDTRole),
							Topics:     [][]byte{[]byte("TTT-abcd"), big.NewInt(0).Bytes(), big.NewInt(0).Bytes(), []byte(core.DCDTRoleNFTCreate), []byte(core.DCDTRoleNFTBurn)},
						},
						nil,
					},
				},
			},
		},
	}

	header.TimeStamp = 10000
	err = esProc.SaveTransactions(createOutportBlockWithHeader(body, header, pool, map[string]*alteredAccount.AlteredAccount{}, testNumOfShards))
	require.Nil(t, err)

	ids := []string{"TTT-abcd"}
	genericResponse := &GenericResponse{}
	err = esClient.DoMultiGet(context.Background(), ids, indexdrtata.TokensIndex, true, genericResponse)
	require.Nil(t, err)
	require.JSONEq(t, readExpectedResult("./testdata/issueTokenAndSetRoles/token-after-set-roles-first.json"), string(genericResponse.Docs[0].Source))

	// ISSUE
	pool = &outport.TransactionPool{
		Logs: []*outport.LogData{
			{
				TxHash: hex.EncodeToString([]byte("h1")),
				Log: &transaction.Log{
					Address: decodeAddress(address1),
					Events: []*transaction.Event{
						{
							Address:    decodeAddress(address1),
							Identifier: []byte("issueSemiFungible"),
							Topics:     [][]byte{[]byte("TTT-abcd"), []byte("semi-token"), []byte("SEMI"), []byte(core.SemiFungibleDCDT)},
						},
						nil,
					},
				},
			},
		},
	}

	err = esProc.SaveTransactions(createOutportBlockWithHeader(body, header, pool, map[string]*alteredAccount.AlteredAccount{}, testNumOfShards))
	require.Nil(t, err)

	ids = []string{"TTT-abcd"}
	genericResponse = &GenericResponse{}
	err = esClient.DoMultiGet(context.Background(), ids, indexdrtata.TokensIndex, true, genericResponse)
	require.Nil(t, err)
	require.JSONEq(t, readExpectedResult("./testdata/issueTokenAndSetRoles/token-after-set-roles-and-issue.json"), string(genericResponse.Docs[0].Source))
}
