//go:build integrationtests

package integrationtests

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data/alteredAccount"
	dataBlock "github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain-core/data/dcdt"
	"github.com/TerraDharitri/drt-go-chain-core/data/outport"
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	indexdrtata "github.com/TerraDharitri/drt-go-chain-es-indexer/process/dataindexer"
	"github.com/stretchr/testify/require"
)

func TestIndexAccountsBalance(t *testing.T) {
	setLogLevelDebug()

	esClient, err := createESClient(esURL)
	require.Nil(t, err)

	// ################ UPDATE ACCOUNT-DCDT BALANCE ##########################
	body := &dataBlock.Body{}

	dcdtToken := &dcdt.DCDigitalToken{
		Value: big.NewInt(1000),
	}

	addr := "drt17umc0uvel62ng30k5uprqcxh3ue33hq608njejaqljuqzqlxtzuqyq4umj"
	addr2 := "drt1m2pyjudsqt8gn0tnsstht35gfqcfx8ku5utz07mf2r6pq3sfxjzslt09es"

	account := &alteredAccount.AlteredAccount{
		Address: addr,
		Balance: "0",
		Tokens: []*alteredAccount.AccountTokenData{
			{
				Identifier: "TTTT-abcd",
				Balance:    "1000",
				Nonce:      0,
			},
		},
	}

	coreAlteredAccounts := map[string]*alteredAccount.AlteredAccount{
		addr:  account,
		addr2: account,
	}

	esProc, err := CreateElasticProcessor(esClient)
	require.Nil(t, err)

	header := &dataBlock.Header{
		Round:     51,
		TimeStamp: 5600,
		ShardID:   2,
	}

	pool := &outport.TransactionPool{
		Logs: []*outport.LogData{
			{
				TxHash: hex.EncodeToString([]byte("h1")),
				Log: &transaction.Log{
					Events: []*transaction.Event{
						{
							Address:    []byte("eeeebbbb"),
							Identifier: []byte(core.BuiltInFunctionDCDTTransfer),
							Topics:     [][]byte{[]byte("TTTT-abcd"), nil, big.NewInt(1).Bytes()},
						},
						nil,
					},
				},
			},
		},
	}

	err = esProc.SaveTransactions(createOutportBlockWithHeader(body, header, pool, coreAlteredAccounts, testNumOfShards))
	require.Nil(t, err)

	ids := []string{addr}
	genericResponse := &GenericResponse{}
	err = esClient.DoMultiGet(context.Background(), ids, indexdrtata.AccountsIndex, true, genericResponse)
	require.Nil(t, err)
	require.JSONEq(t, readExpectedResult("./testdata/accountsBalanceWithLowerTimestamp/account-balance-first-update.json"), string(genericResponse.Docs[0].Source))

	ids = []string{fmt.Sprintf("%s-TTTT-abcd-00", addr)}
	genericResponse = &GenericResponse{}
	err = esClient.DoMultiGet(context.Background(), ids, indexdrtata.AccountsDCDTIndex, true, genericResponse)
	require.Nil(t, err)
	require.JSONEq(t, readExpectedResult("./testdata/accountsBalanceWithLowerTimestamp/account-balance-dcdt-first-update.json"), string(genericResponse.Docs[0].Source))

	//////////////////// INDEX BALANCE LOWER TIMESTAMP ///////////////////////////////////

	header = &dataBlock.Header{
		Round:     51,
		TimeStamp: 5000,
		ShardID:   2,
	}

	err = esProc.SaveTransactions(createOutportBlockWithHeader(body, header, pool, map[string]*alteredAccount.AlteredAccount{}, testNumOfShards))
	require.Nil(t, err)

	ids = []string{addr}
	genericResponse = &GenericResponse{}
	err = esClient.DoMultiGet(context.Background(), ids, indexdrtata.AccountsIndex, true, genericResponse)
	require.Nil(t, err)
	require.JSONEq(t, readExpectedResult("./testdata/accountsBalanceWithLowerTimestamp/account-balance-first-update.json"), string(genericResponse.Docs[0].Source))

	ids = []string{fmt.Sprintf("%s-TTTT-abcd-00", addr)}
	genericResponse = &GenericResponse{}
	err = esClient.DoMultiGet(context.Background(), ids, indexdrtata.AccountsDCDTIndex, true, genericResponse)
	require.Nil(t, err)
	require.JSONEq(t, readExpectedResult("./testdata/accountsBalanceWithLowerTimestamp/account-balance-dcdt-first-update.json"), string(genericResponse.Docs[0].Source))

	//////////////////// INDEX BALANCE HIGHER TIMESTAMP ///////////////////////////////////
	header = &dataBlock.Header{
		Round:     51,
		TimeStamp: 6000,
		ShardID:   2,
	}

	coreAlteredAccounts[addr].Balance = "2000"
	coreAlteredAccounts[addr].AdditionalData = &alteredAccount.AdditionalAccountData{
		IsSender:       true,
		BalanceChanged: true,
	}
	pool = &outport.TransactionPool{
		Transactions: map[string]*outport.TxInfo{
			hex.EncodeToString([]byte("h1")): {
				Transaction: &transaction.Transaction{
					SndAddr: []byte(addr),
				},
				FeeInfo: &outport.FeeInfo{},
			},
		},
		Logs: []*outport.LogData{
			{
				TxHash: hex.EncodeToString([]byte("h1")),
				Log: &transaction.Log{
					Events: []*transaction.Event{
						{
							Address:    decodeAddress(addr2),
							Identifier: []byte(core.BuiltInFunctionDCDTTransfer),
							Topics:     [][]byte{[]byte("TTTT-abcd"), nil, big.NewInt(1).Bytes()},
						},
						nil,
					},
				},
			},
		},
	}
	body = &dataBlock.Body{
		MiniBlocks: []*dataBlock.MiniBlock{
			{
				Type:          dataBlock.TxBlock,
				TxHashes:      [][]byte{[]byte("h1")},
				SenderShardID: 2,
			},
		},
	}

	err = esProc.SaveTransactions(createOutportBlockWithHeader(body, header, pool, coreAlteredAccounts, testNumOfShards))
	require.Nil(t, err)

	ids = []string{addr}
	genericResponse = &GenericResponse{}
	err = esClient.DoMultiGet(context.Background(), ids, indexdrtata.AccountsIndex, true, genericResponse)
	require.Nil(t, err)
	require.JSONEq(t, readExpectedResult("./testdata/accountsBalanceWithLowerTimestamp/account-balance-second-update.json"), string(genericResponse.Docs[0].Source))

	ids = []string{fmt.Sprintf("%s-TTTT-abcd-00", addr)}
	genericResponse = &GenericResponse{}
	err = esClient.DoMultiGet(context.Background(), ids, indexdrtata.AccountsDCDTIndex, true, genericResponse)
	require.Nil(t, err)
	require.JSONEq(t, readExpectedResult("./testdata/accountsBalanceWithLowerTimestamp/account-balance-dcdt-second-update.json"), string(genericResponse.Docs[0].Source))

	//////////////////////// DELETE DCDT BALANCE LOWER TIMESTAMP ////////////////

	dcdtToken.Value = big.NewInt(0)
	esProc, err = CreateElasticProcessor(esClient)
	require.Nil(t, err)

	header = &dataBlock.Header{
		Round:     51,
		TimeStamp: 6001,
		ShardID:   2,
	}

	coreAlteredAccounts[addr].Balance = "2000"
	coreAlteredAccounts[addr].Tokens[0].Balance = "0"
	coreAlteredAccounts[addr].AdditionalData = &alteredAccount.AdditionalAccountData{
		IsSender:       false,
		BalanceChanged: false,
	}

	pool.Transactions = make(map[string]*outport.TxInfo)
	err = esProc.SaveTransactions(createOutportBlockWithHeader(body, header, pool, coreAlteredAccounts, testNumOfShards))
	require.Nil(t, err)

	ids = []string{addr}
	genericResponse = &GenericResponse{}
	err = esClient.DoMultiGet(context.Background(), ids, indexdrtata.AccountsIndex, true, genericResponse)
	require.Nil(t, err)
	require.JSONEq(t, readExpectedResult("./testdata/accountsBalanceWithLowerTimestamp/account-balance-second-update.json"), string(genericResponse.Docs[0].Source))

	ids = []string{fmt.Sprintf("%s-TTTT-abcd-00", addr)}
	genericResponse = &GenericResponse{}
	err = esClient.DoMultiGet(context.Background(), ids, indexdrtata.AccountsDCDTIndex, true, genericResponse)
	require.Nil(t, err)
	require.False(t, genericResponse.Docs[0].Found)
}
