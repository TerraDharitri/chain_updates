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

func TestDCDTModifyCreator(t *testing.T) {
	tokenTypes := getDynamicTokenTypes()
	for _, tokenType := range tokenTypes {
		dcdtType, _ := core.ConvertDCDTTypeToUint32(tokenType)
		if !core.IsDynamicDCDT(dcdtType) {
			continue
		}
		testName := "dcdtModifyCreator for " + tokenType
		t.Run(testName, func(t *testing.T) {
			runDcdtModifyCreatorTest(t, tokenType)
		})
	}
}

func runDcdtModifyCreatorTest(t *testing.T, tokenType string) {
	newCreator := []byte("12345678901234567890123456789012")
	creatorAddr := []byte("12345678901234567890123456789013")
	token := []byte("tokenId")
	baseDcdtKeyPrefix := core.ProtectedKeyPrefix + core.DCDTKeyIdentifier
	key := append([]byte(baseDcdtKeyPrefix), token...)

	testContext, err := vm.CreatePreparedTxProcessorWithVMs(config.EnableEpochs{}, 1)
	require.Nil(t, err)
	defer testContext.Close()
	testContext.BlockchainHook.(process.BlockChainHookHandler).SetCurrentHeader(&dataBlock.Header{Round: 7})

	createAccWithBalance(t, testContext.Accounts, newCreator, big.NewInt(100000000))
	createAccWithBalance(t, testContext.Accounts, creatorAddr, big.NewInt(100000000))
	createAccWithBalance(t, testContext.Accounts, core.DCDTSCAddress, big.NewInt(100000000))
	utils.SetDCDTRoles(t, testContext.Accounts, creatorAddr, token, [][]byte{[]byte(core.DCDTRoleNFTCreate)})
	utils.SetDCDTRoles(t, testContext.Accounts, newCreator, token, [][]byte{[]byte(core.DCDTRoleModifyCreator)})

	tx := setTokenTypeTx(core.DCDTSCAddress, 100000, token, tokenType)
	retCode, err := testContext.TxProcessor.ProcessTransaction(tx)
	require.Equal(t, vmcommon.Ok, retCode)
	require.Nil(t, err)

	defaultMetaData := GetDefaultMetaData()
	defaultMetaData.Nonce = []byte(hex.EncodeToString(big.NewInt(1).Bytes()))
	tx = createTokenTx(creatorAddr, creatorAddr, 100000, 1, defaultMetaData)
	retCode, err = testContext.TxProcessor.ProcessTransaction(tx)
	require.Equal(t, vmcommon.Ok, retCode)
	require.Nil(t, err)

	tx = dcdtModifyCreatorTx(newCreator, newCreator, 100000, defaultMetaData)
	retCode, err = testContext.TxProcessor.ProcessTransaction(tx)
	require.Equal(t, vmcommon.Ok, retCode)
	require.Nil(t, err)

	_, err = testContext.Accounts.Commit()
	require.Nil(t, err)

	retrievedMetaData := &dcdt.MetaData{}
	if tokenType == core.DynamicNFTDCDT {
		retrievedMetaData = getMetaDataFromAcc(t, testContext, newCreator, key)
	} else {
		retrievedMetaData = getMetaDataFromAcc(t, testContext, core.SystemAccountAddress, key)
	}
	require.Equal(t, newCreator, retrievedMetaData.Creator)
}

func dcdtModifyCreatorTx(
	sndAddr []byte,
	rcvAddr []byte,
	gasLimit uint64,
	metaData *MetaData,
) *transaction.Transaction {
	txDataField := bytes.Join(
		[][]byte{
			[]byte(core.DCDTModifyCreator),
			metaData.TokenId,
			metaData.Nonce,
		},
		[]byte("@"),
	)
	return &transaction.Transaction{
		Nonce:    0,
		SndAddr:  sndAddr,
		RcvAddr:  rcvAddr,
		GasLimit: gasLimit,
		GasPrice: gasPrice,

		Data:  txDataField,
		Value: big.NewInt(0),
	}
}
