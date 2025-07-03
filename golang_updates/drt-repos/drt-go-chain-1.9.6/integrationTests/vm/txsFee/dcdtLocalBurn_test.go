package txsFee

import (
	"math/big"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/integrationTests/vm"
	"github.com/TerraDharitri/drt-go-chain/integrationTests/vm/txsFee/utils"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/stretchr/testify/require"
)

func TestDCDTLocalBurnShouldWork(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	testContext, err := vm.CreatePreparedTxProcessorWithVMs(config.EnableEpochs{}, 1)
	require.Nil(t, err)
	defer testContext.Close()

	sndAddr := []byte("12345678901234567890123456789012")

	rewaBalance := big.NewInt(100000000)
	dcdtBalance := big.NewInt(100000000)
	token := []byte("miiutoken")
	roles := [][]byte{[]byte(core.DCDTRoleLocalMint), []byte(core.DCDTRoleLocalBurn)}
	utils.CreateAccountWithDCDTBalanceAndRoles(t, testContext.Accounts, sndAddr, rewaBalance, token, 0, dcdtBalance, roles)

	gasLimit := uint64(40)
	tx := utils.CreateDCDTLocalBurnTx(0, sndAddr, sndAddr, token, big.NewInt(100), gasPrice, gasLimit)
	retCode, err := testContext.TxProcessor.ProcessTransaction(tx)
	require.Equal(t, vmcommon.Ok, retCode)
	require.Nil(t, err)

	_, err = testContext.Accounts.Commit()
	require.Nil(t, err)

	expectedBalanceSnd := big.NewInt(99999900)
	utils.CheckDCDTBalance(t, testContext, sndAddr, token, expectedBalanceSnd)

	// check accumulated fees
	accumulatedFees := testContext.TxFeeHandler.GetAccumulatedFees()
	require.Equal(t, big.NewInt(370), accumulatedFees)
}

func TestDCDTLocalBurnMoreThanTotalBalanceShouldErr(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	testContext, err := vm.CreatePreparedTxProcessorWithVMs(config.EnableEpochs{}, 1)
	require.Nil(t, err)
	defer testContext.Close()

	sndAddr := []byte("12345678901234567890123456789012")

	rewaBalance := big.NewInt(100000000)
	dcdtBalance := big.NewInt(100000000)
	token := []byte("miiutoken")
	roles := [][]byte{[]byte(core.DCDTRoleLocalMint), []byte(core.DCDTRoleLocalBurn)}
	utils.CreateAccountWithDCDTBalanceAndRoles(t, testContext.Accounts, sndAddr, rewaBalance, token, 0, dcdtBalance, roles)

	gasLimit := uint64(60)
	tx := utils.CreateDCDTLocalBurnTx(0, sndAddr, sndAddr, token, big.NewInt(100000001), gasPrice, gasLimit)
	retCode, err := testContext.TxProcessor.ProcessTransaction(tx)
	require.Equal(t, vmcommon.UserError, retCode)
	require.Equal(t, process.ErrFailedTransaction, err)

	_, err = testContext.Accounts.Commit()
	require.Nil(t, err)

	expectedBalanceSnd := big.NewInt(100000000)
	utils.CheckDCDTBalance(t, testContext, sndAddr, token, expectedBalanceSnd)

	// check accumulated fees
	accumulatedFees := testContext.TxFeeHandler.GetAccumulatedFees()
	require.Equal(t, big.NewInt(600), accumulatedFees)
}

func TestDCDTLocalBurnNotAllowedShouldErr(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	testContext, err := vm.CreatePreparedTxProcessorWithVMs(config.EnableEpochs{}, 1)
	require.Nil(t, err)
	defer testContext.Close()

	sndAddr := []byte("12345678901234567890123456789012")

	rewaBalance := big.NewInt(100000000)
	dcdtBalance := big.NewInt(100000000)
	token := []byte("miiutoken")
	utils.CreateAccountWithDCDTBalance(t, testContext.Accounts, sndAddr, rewaBalance, token, 0, dcdtBalance, uint32(core.Fungible))

	gasLimit := uint64(40)
	tx := utils.CreateDCDTLocalBurnTx(0, sndAddr, sndAddr, token, big.NewInt(100), gasPrice, gasLimit)
	retCode, err := testContext.TxProcessor.ProcessTransaction(tx)
	require.Equal(t, vmcommon.UserError, retCode)
	require.Equal(t, process.ErrFailedTransaction, err)

	_, err = testContext.Accounts.Commit()
	require.Nil(t, err)

	expectedBalanceSnd := big.NewInt(100000000)
	utils.CheckDCDTBalance(t, testContext, sndAddr, token, expectedBalanceSnd)

	// check accumulated fees
	accumulatedFees := testContext.TxFeeHandler.GetAccumulatedFees()
	require.Equal(t, big.NewInt(400), accumulatedFees)
}
