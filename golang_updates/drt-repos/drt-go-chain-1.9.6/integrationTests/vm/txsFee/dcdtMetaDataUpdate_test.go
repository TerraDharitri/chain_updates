package txsFee

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	dataBlock "github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/integrationTests/vm"
	"github.com/TerraDharitri/drt-go-chain/integrationTests/vm/txsFee/utils"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/stretchr/testify/require"
)

func TestDCDTMetaDataUpdate(t *testing.T) {
	tokenTypes := getDynamicTokenTypes()
	for _, tokenType := range tokenTypes {
		testName := "metaDataUpdate for " + tokenType
		t.Run(testName, func(t *testing.T) {
			runDcdtMetaDataUpdateTest(t, tokenType)
		})
	}
}

func runDcdtMetaDataUpdateTest(t *testing.T, tokenType string) {
	sndAddr := []byte("12345678901234567890123456789012")
	token := []byte("tokenId")
	roles := [][]byte{[]byte(core.DCDTRoleNFTUpdate), []byte(core.DCDTRoleNFTCreate)}
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

	// TODO change default metadata
	defaultMetaData.Nonce = []byte(hex.EncodeToString(big.NewInt(1).Bytes()))
	defaultMetaData.Name = []byte(hex.EncodeToString([]byte("newName")))
	defaultMetaData.Hash = []byte(hex.EncodeToString([]byte("newHash")))
	defaultMetaData.Uris = [][]byte{defaultMetaData.Uris[1]}
	tx = dcdtMetaDataUpdateTx(sndAddr, sndAddr, 100000, defaultMetaData)
	retCode, err = testContext.TxProcessor.ProcessTransaction(tx)
	require.Equal(t, vmcommon.Ok, retCode)
	require.Nil(t, err)

	_, err = testContext.Accounts.Commit()
	require.Nil(t, err)

	if tokenType == core.DynamicNFTDCDT {
		checkMetaData(t, testContext, sndAddr, key, defaultMetaData)
	} else {
		checkMetaData(t, testContext, core.SystemAccountAddress, key, defaultMetaData)
	}
}

func dcdtMetaDataUpdateTx(
	sndAddr []byte,
	rcvAddr []byte,
	gasLimit uint64,
	metaData *MetaData,
) *transaction.Transaction {
	txDataField := bytes.Join(
		[][]byte{
			[]byte(core.DCDTMetaDataUpdate),
			metaData.TokenId,
			metaData.Nonce,
			metaData.Name,
			metaData.Royalties,
			metaData.Hash,
			metaData.Attributes,
			metaData.Uris[0],
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
