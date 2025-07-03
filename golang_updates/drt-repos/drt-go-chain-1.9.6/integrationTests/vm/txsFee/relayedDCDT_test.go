package txsFee

import (
	"math/big"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/integrationTests"
	"github.com/TerraDharitri/drt-go-chain/integrationTests/vm"
	"github.com/TerraDharitri/drt-go-chain/integrationTests/vm/txsFee/utils"
	"github.com/stretchr/testify/require"
)

func TestRelayedDCDTTransferShouldWork(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	t.Run("before relayed base cost fix", testRelayedDCDTTransferShouldWork(integrationTests.UnreachableEpoch, big.NewInt(9997614), big.NewInt(2386)))
	t.Run("after relayed base cost fix", testRelayedDCDTTransferShouldWork(0, big.NewInt(9997299), big.NewInt(2701)))
}

func testRelayedDCDTTransferShouldWork(
	relayedFixActivationEpoch uint32,
	expectedRelayerBalance *big.Int,
	expectedAccFees *big.Int,
) func(t *testing.T) {
	return func(t *testing.T) {
		testContext, err := vm.CreatePreparedTxProcessorWithVMs(config.EnableEpochs{
			FixRelayedBaseCostEnableEpoch: relayedFixActivationEpoch,
		}, gasPriceModifier)
		require.Nil(t, err)
		defer testContext.Close()

		relayerAddr := []byte("12345678901234567890123456789033")
		sndAddr := []byte("12345678901234567890123456789012")
		rcvAddr := []byte("12345678901234567890123456789022")

		relayerBalance := big.NewInt(10000000)
		localDcdtBalance := big.NewInt(100000000)
		token := []byte("miiutoken")
		utils.CreateAccountWithDCDTBalance(t, testContext.Accounts, sndAddr, big.NewInt(0), token, 0, localDcdtBalance, uint32(core.Fungible))
		_, _ = vm.CreateAccount(testContext.Accounts, relayerAddr, 0, relayerBalance)

		gasLimit := uint64(40)
		innerTx := utils.CreateDCDTTransferTx(0, sndAddr, rcvAddr, token, big.NewInt(100), gasPrice, gasLimit)

		rtxData := integrationTests.PrepareRelayedTxDataV1(innerTx)
		rTxGasLimit := minGasLimit + gasLimit + uint64(len(rtxData))
		rtx := vm.CreateTransaction(0, innerTx.Value, relayerAddr, sndAddr, gasPrice, rTxGasLimit, rtxData)

		retCode, err := testContext.TxProcessor.ProcessTransaction(rtx)
		require.Equal(t, vmcommon.Ok, retCode)
		require.Nil(t, err)

		_, err = testContext.Accounts.Commit()
		require.Nil(t, err)

		expectedBalanceSnd := big.NewInt(99999900)
		utils.CheckDCDTBalance(t, testContext, sndAddr, token, expectedBalanceSnd)

		expectedReceiverBalance := big.NewInt(100)
		utils.CheckDCDTBalance(t, testContext, rcvAddr, token, expectedReceiverBalance)

		expectedREWABalance := big.NewInt(0)
		utils.TestAccount(t, testContext.Accounts, sndAddr, 1, expectedREWABalance)

		utils.TestAccount(t, testContext.Accounts, relayerAddr, 1, expectedRelayerBalance)

		// check accumulated fees
		accumulatedFees := testContext.TxFeeHandler.GetAccumulatedFees()
		require.Equal(t, expectedAccFees, accumulatedFees)
	}
}

func TestRelayedESTTransferNotEnoughESTValueShouldConsumeGas(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	t.Run("before relayed base cost fix", testRelayedESTTransferNotEnoughESTValueShouldConsumeGas(integrationTests.UnreachableEpoch, big.NewInt(9997488), big.NewInt(2512)))
	t.Run("after relayed base cost fix", testRelayedESTTransferNotEnoughESTValueShouldConsumeGas(0, big.NewInt(9997119), big.NewInt(2881)))
}

func testRelayedESTTransferNotEnoughESTValueShouldConsumeGas(
	relayedFixActivationEpoch uint32,
	expectedRelayerBalance *big.Int,
	expectedAccFees *big.Int,
) func(t *testing.T) {
	return func(t *testing.T) {
		testContext, err := vm.CreatePreparedTxProcessorWithVMs(config.EnableEpochs{
			FixRelayedBaseCostEnableEpoch: relayedFixActivationEpoch,
		}, gasPriceModifier)
		require.Nil(t, err)
		defer testContext.Close()

		relayerAddr := []byte("12345678901234567890123456789033")
		sndAddr := []byte("12345678901234567890123456789012")
		rcvAddr := []byte("12345678901234567890123456789022")

		relayerBalance := big.NewInt(10000000)
		localDcdtBalance := big.NewInt(100000000)
		token := []byte("miiutoken")
		utils.CreateAccountWithDCDTBalance(t, testContext.Accounts, sndAddr, big.NewInt(0), token, 0, localDcdtBalance, uint32(core.Fungible))
		_, _ = vm.CreateAccount(testContext.Accounts, relayerAddr, 0, relayerBalance)

		gasLimit := uint64(42)
		innerTx := utils.CreateDCDTTransferTx(0, sndAddr, rcvAddr, token, big.NewInt(100000001), gasPrice, gasLimit)

		rtxData := integrationTests.PrepareRelayedTxDataV1(innerTx)
		rTxGasLimit := minGasLimit + gasLimit + uint64(len(rtxData))
		rtx := vm.CreateTransaction(0, innerTx.Value, relayerAddr, sndAddr, gasPrice, rTxGasLimit, rtxData)

		retCode, err := testContext.TxProcessor.ProcessTransaction(rtx)
		require.Equal(t, vmcommon.ExecutionFailed, retCode)
		require.Nil(t, err)

		_, err = testContext.Accounts.Commit()
		require.Nil(t, err)

		expectedBalanceSnd := big.NewInt(100000000)
		utils.CheckDCDTBalance(t, testContext, sndAddr, token, expectedBalanceSnd)

		expectedReceiverBalance := big.NewInt(0)
		utils.CheckDCDTBalance(t, testContext, rcvAddr, token, expectedReceiverBalance)

		expectedREWABalance := big.NewInt(0)
		utils.TestAccount(t, testContext.Accounts, sndAddr, 1, expectedREWABalance)

		utils.TestAccount(t, testContext.Accounts, relayerAddr, 1, expectedRelayerBalance)

		// check accumulated fees
		accumulatedFees := testContext.TxFeeHandler.GetAccumulatedFees()
		require.Equal(t, expectedAccFees, accumulatedFees)
	}
}
