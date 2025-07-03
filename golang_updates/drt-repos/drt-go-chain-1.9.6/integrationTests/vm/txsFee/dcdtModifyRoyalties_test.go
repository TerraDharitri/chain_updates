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

func TestDCDTModifyRoyalties(t *testing.T) {
	tokenTypes := getDynamicTokenTypes()
	for _, tokenType := range tokenTypes {
		testName := "dcdtModifyRoyalties for " + tokenType
		t.Run(testName, func(t *testing.T) {
			runDcdtModifyRoyaltiesTest(t, tokenType)
		})
	}
}

func runDcdtModifyRoyaltiesTest(t *testing.T, tokenType string) {
	creatorAddr := []byte("12345678901234567890123456789013")
	token := []byte("tokenId")
	baseDcdtKeyPrefix := core.ProtectedKeyPrefix + core.DCDTKeyIdentifier
	key := append([]byte(baseDcdtKeyPrefix), token...)

	testContext, err := vm.CreatePreparedTxProcessorWithVMs(config.EnableEpochs{}, 1)
	require.Nil(t, err)
	defer testContext.Close()
	testContext.BlockchainHook.(process.BlockChainHookHandler).SetCurrentHeader(&dataBlock.Header{Round: 7})

	createAccWithBalance(t, testContext.Accounts, creatorAddr, big.NewInt(100000000))
	createAccWithBalance(t, testContext.Accounts, core.DCDTSCAddress, big.NewInt(100000000))
	utils.SetDCDTRoles(t, testContext.Accounts, creatorAddr, token, [][]byte{[]byte(core.DCDTRoleModifyRoyalties), []byte(core.DCDTRoleNFTCreate)})

	tx := setTokenTypeTx(core.DCDTSCAddress, 100000, token, tokenType)
	retCode, err := testContext.TxProcessor.ProcessTransaction(tx)
	require.Equal(t, vmcommon.Ok, retCode)
	require.Nil(t, err)

	defaultMetaData := GetDefaultMetaData()
	tx = createTokenTx(creatorAddr, creatorAddr, 100000, 1, defaultMetaData)
	retCode, err = testContext.TxProcessor.ProcessTransaction(tx)
	require.Equal(t, vmcommon.Ok, retCode)
	require.Nil(t, err)

	defaultMetaData.Nonce = []byte(hex.EncodeToString(big.NewInt(1).Bytes()))
	defaultMetaData.Royalties = []byte(hex.EncodeToString(big.NewInt(20).Bytes()))
	tx = dcdtModifyRoyaltiesTx(creatorAddr, creatorAddr, 100000, defaultMetaData)
	retCode, err = testContext.TxProcessor.ProcessTransaction(tx)
	require.Equal(t, vmcommon.Ok, retCode)
	require.Nil(t, err)

	_, err = testContext.Accounts.Commit()
	require.Nil(t, err)

	retrievedMetaData := &dcdt.MetaData{}
	if tokenType == core.DynamicNFTDCDT {
		retrievedMetaData = getMetaDataFromAcc(t, testContext, creatorAddr, key)
	} else {
		retrievedMetaData = getMetaDataFromAcc(t, testContext, core.SystemAccountAddress, key)
	}
	require.Equal(t, uint32(big.NewInt(20).Uint64()), retrievedMetaData.Royalties)
}

func dcdtModifyRoyaltiesTx(
	sndAddr []byte,
	rcvAddr []byte,
	gasLimit uint64,
	metaData *MetaData,
) *transaction.Transaction {
	txDataField := bytes.Join(
		[][]byte{
			[]byte(core.DCDTModifyRoyalties),
			metaData.TokenId,
			metaData.Nonce,
			metaData.Royalties,
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
