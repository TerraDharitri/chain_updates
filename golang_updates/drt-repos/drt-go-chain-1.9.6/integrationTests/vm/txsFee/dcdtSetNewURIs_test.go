package txsFee

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	dataBlock "github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain-core/data/dcdt"
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/integrationTests/vm"
	"github.com/TerraDharitri/drt-go-chain/integrationTests/vm/txsFee/utils"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/stretchr/testify/require"
)

func TestDCDTSetNewURIs(t *testing.T) {
	tokenTypes := getDynamicTokenTypes()
	for _, tokenType := range tokenTypes {
		testName := "DCDTsetNewURIs for " + tokenType
		t.Run(testName, func(t *testing.T) {
			runDcdtSetNewURIsTest(t, tokenType)
		})
	}
}

func runDcdtSetNewURIsTest(t *testing.T, tokenType string) {
	sndAddr := []byte("12345678901234567890123456789012")
	token := []byte("tokenId")
	roles := [][]byte{[]byte(core.DCDTRoleSetNewURI), []byte(core.DCDTRoleNFTCreate)}
	baseDcdtKeyPrefix := core.ProtectedKeyPrefix + core.DCDTKeyIdentifier
	key := append([]byte(baseDcdtKeyPrefix), token...)

	testContext, err := vm.CreatePreparedTxProcessorWithVMs(config.EnableEpochs{}, 1)
	require.Nil(t, err)
	defer testContext.Close()
	testContext.BlockchainHook.(process.BlockChainHookHandler).SetCurrentHeader(&dataBlock.Header{Round: 7})

	createAccWithBalance(t, testContext.Accounts, sndAddr, big.NewInt(100000000))
	createAccWithBalance(t, testContext.Accounts, core.DCDTSCAddress, big.NewInt(100000000))
	utils.SetDCDTRoles(t, testContext.Accounts, sndAddr, token, roles)

	tx := setTokenTypeTx(core.DCDTSCAddress, 100000, token, tokenType)
	retCode, err := testContext.TxProcessor.ProcessTransaction(tx)
	require.Equal(t, vmcommon.Ok, retCode)
	require.Nil(t, err)

	defaultMetaData := GetDefaultMetaData()
	tx = createTokenTx(sndAddr, sndAddr, 100000, 1, defaultMetaData)
	retCode, err = testContext.TxProcessor.ProcessTransaction(tx)
	require.Equal(t, vmcommon.Ok, retCode)
	require.Nil(t, err)

	defaultMetaData.Uris = [][]byte{[]byte(hex.EncodeToString([]byte("newUri1"))), []byte(hex.EncodeToString([]byte("newUri2")))}
	defaultMetaData.Nonce = []byte(hex.EncodeToString(big.NewInt(1).Bytes()))
	tx = dcdtSetNewUrisTx(sndAddr, sndAddr, 100000, defaultMetaData)
	retCode, err = testContext.TxProcessor.ProcessTransaction(tx)
	require.Equal(t, vmcommon.Ok, retCode)
	require.Nil(t, err)

	_, err = testContext.Accounts.Commit()
	require.Nil(t, err)

	retrievedMetaData := &dcdt.MetaData{}
	if tokenType == core.DynamicNFTDCDT {
		retrievedMetaData = getMetaDataFromAcc(t, testContext, sndAddr, key)
	} else {
		retrievedMetaData = getMetaDataFromAcc(t, testContext, core.SystemAccountAddress, key)
	}
	require.Equal(t, [][]byte{[]byte("newUri1"), []byte("newUri2")}, retrievedMetaData.URIs)
}

func dcdtSetNewUrisTx(
	sndAddr []byte,
	rcvAddr []byte,
	gasLimit uint64,
	metaData *MetaData,
) *transaction.Transaction {
	txDataField := bytes.Join(
		[][]byte{
			[]byte(core.DCDTSetNewURIs),
			metaData.TokenId,
			metaData.Nonce,
			metaData.Uris[0],
			metaData.Uris[1],
		},
		[]byte("@"),
	)

	return &transaction.Transaction{
		Nonce:    1,
		SndAddr:  sndAddr,
		RcvAddr:  rcvAddr,
		GasLimit: gasLimit,
		GasPrice: gasPrice,

		Data:  txDataField,
		Value: big.NewInt(0),
	}
}
