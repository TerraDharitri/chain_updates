package txsFee

import (
	"context"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/common/errChan"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/integrationTests"
	"github.com/TerraDharitri/drt-go-chain/integrationTests/vm"
	"github.com/TerraDharitri/drt-go-chain/integrationTests/vm/txsFee/utils"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/state"
	"github.com/TerraDharitri/drt-go-chain/state/parsers"
	"github.com/stretchr/testify/require"
)

func TestAsyncDCDTCallShouldWork(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	testContext, err := vm.CreatePreparedTxProcessorWithVMs(config.EnableEpochs{
		DynamicGasCostForDataTrieStorageLoadEnableEpoch: integrationTests.UnreachableEpoch,
	}, 1)
	require.Nil(t, err)
	defer testContext.Close()

	localRewaBalance := big.NewInt(100000000)
	ownerAddr := []byte("12345678901234567890123456789010")
	_, _ = vm.CreateAccount(testContext.Accounts, ownerAddr, 0, localRewaBalance)

	// create an address with DCDT token
	sndAddr := []byte("12345678901234567890123456789012")

	localDcdtBalance := big.NewInt(100000000)
	token := []byte("miiutoken")
	utils.CreateAccountWithDCDTBalance(t, testContext.Accounts, sndAddr, localRewaBalance, token, 0, localDcdtBalance, uint32(core.Fungible))

	// deploy 2 contracts
	ownerAccount, _ := testContext.Accounts.LoadAccount(ownerAddr)
	deployGasLimit := uint64(50000)

	argsSecond := [][]byte{[]byte(hex.EncodeToString(token))}
	secondSCAddress := utils.DoDeploySecond(t, testContext, "../dcdt/testdata/second-contract.wasm", ownerAccount, gasPrice, deployGasLimit, argsSecond, big.NewInt(0))

	args := [][]byte{[]byte(hex.EncodeToString(token)), []byte(hex.EncodeToString(secondSCAddress))}
	ownerAccount, _ = testContext.Accounts.LoadAccount(ownerAddr)
	firstSCAddress := utils.DoDeploySecond(t, testContext, "../dcdt/testdata/first-contract.wasm", ownerAccount, gasPrice, deployGasLimit, args, big.NewInt(0))

	testContext.TxFeeHandler.CreateBlockStarted(getZeroGasAndFees())
	utils.CleanAccumulatedIntermediateTransactions(t, testContext)

	gasLimit := uint64(500000)
	tx := utils.CreateDCDTTransferTx(0, sndAddr, firstSCAddress, token, big.NewInt(5000), gasPrice, gasLimit)
	tx.Data = []byte(string(tx.Data) + "@" + hex.EncodeToString([]byte("transferToSecondContractHalf")))

	retCode, err := testContext.TxProcessor.ProcessTransaction(tx)
	require.Equal(t, vmcommon.Ok, retCode)
	require.Nil(t, err)

	_, err = testContext.Accounts.Commit()
	require.Nil(t, err)

	utils.CheckDCDTBalance(t, testContext, firstSCAddress, token, big.NewInt(2500))
	utils.CheckDCDTBalance(t, testContext, secondSCAddress, token, big.NewInt(2500))

	expectedSenderBalance := big.NewInt(98223470)
	utils.TestAccount(t, testContext.Accounts, sndAddr, 1, expectedSenderBalance)

	expectedAccumulatedFees := big.NewInt(1776530)
	accumulatedFees := testContext.TxFeeHandler.GetAccumulatedFees()
	require.Equal(t, expectedAccumulatedFees, accumulatedFees)
}

func TestAsyncDCDTCallSecondScRefusesPayment(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	testContext, err := vm.CreatePreparedTxProcessorWithVMs(config.EnableEpochs{}, 1)
	require.Nil(t, err)
	defer testContext.Close()

	localRewaBalance := big.NewInt(100000000)
	ownerAddr := []byte("12345678901234567890123456789010")
	_, _ = vm.CreateAccount(testContext.Accounts, ownerAddr, 0, localRewaBalance)

	// create an address with DCDT token
	sndAddr := []byte("12345678901234567890123456789012")

	localDcdtBalance := big.NewInt(100000000)
	token := []byte("miiutoken")
	utils.CreateAccountWithDCDTBalance(t, testContext.Accounts, sndAddr, localRewaBalance, token, 0, localDcdtBalance, uint32(core.Fungible))

	// deploy 2 contracts
	ownerAccount, _ := testContext.Accounts.LoadAccount(ownerAddr)
	deployGasLimit := uint64(50000)

	argsSecond := [][]byte{[]byte(hex.EncodeToString(token))}
	secondSCAddress := utils.DoDeploySecond(t, testContext, "../dcdt/testdata/second-contract.wasm", ownerAccount, gasPrice, deployGasLimit, argsSecond, big.NewInt(0))

	args := [][]byte{[]byte(hex.EncodeToString(token)), []byte(hex.EncodeToString(secondSCAddress))}
	ownerAccount, _ = testContext.Accounts.LoadAccount(ownerAddr)
	firstSCAddress := utils.DoDeploySecond(t, testContext, "../dcdt/testdata/first-contract.wasm", ownerAccount, gasPrice, deployGasLimit, args, big.NewInt(0))

	testContext.TxFeeHandler.CreateBlockStarted(getZeroGasAndFees())
	utils.CleanAccumulatedIntermediateTransactions(t, testContext)
	require.Equal(t, big.NewInt(0), testContext.TxFeeHandler.GetAccumulatedFees())

	gasLimit := uint64(500000)
	tx := utils.CreateDCDTTransferTx(0, sndAddr, firstSCAddress, token, big.NewInt(5000), gasPrice, gasLimit)
	tx.Data = []byte(string(tx.Data) + "@" + hex.EncodeToString([]byte("transferToSecondContractRejected")))

	retCode, err := testContext.TxProcessor.ProcessTransaction(tx)
	require.Equal(t, vmcommon.Ok, retCode)
	require.Nil(t, err)

	_, err = testContext.Accounts.Commit()
	require.Nil(t, err)

	utils.CheckDCDTBalance(t, testContext, firstSCAddress, token, big.NewInt(5000))
	utils.CheckDCDTBalance(t, testContext, secondSCAddress, token, big.NewInt(0))

	expectedSenderBalance := big.NewInt(95999990)
	utils.TestAccount(t, testContext.Accounts, sndAddr, 1, expectedSenderBalance)

	expectedAccumulatedFees := big.NewInt(4000010)
	accumulatedFees := testContext.TxFeeHandler.GetAccumulatedFees()
	require.Equal(t, expectedAccumulatedFees, accumulatedFees)
}

func TestAsyncDCDTCallsOutOfGas(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	testContext, err := vm.CreatePreparedTxProcessorWithVMs(config.EnableEpochs{}, 1)
	require.Nil(t, err)
	defer testContext.Close()

	localRewaBalance := big.NewInt(100000000)
	ownerAddr := []byte("12345678901234567890123456789010")
	_, _ = vm.CreateAccount(testContext.Accounts, ownerAddr, 0, localRewaBalance)

	// create an address with DCDT token
	sndAddr := []byte("12345678901234567890123456789012")

	localDcdtBalance := big.NewInt(100000000)
	token := []byte("miiutoken")
	utils.CreateAccountWithDCDTBalance(t, testContext.Accounts, sndAddr, localRewaBalance, token, 0, localDcdtBalance, uint32(core.Fungible))

	// deploy 2 contracts
	ownerAccount, _ := testContext.Accounts.LoadAccount(ownerAddr)
	deployGasLimit := uint64(50000)

	argsSecond := [][]byte{[]byte(hex.EncodeToString(token))}
	secondSCAddress := utils.DoDeploySecond(t, testContext, "../dcdt/testdata/second-contract.wasm", ownerAccount, gasPrice, deployGasLimit, argsSecond, big.NewInt(0))

	args := [][]byte{[]byte(hex.EncodeToString(token)), []byte(hex.EncodeToString(secondSCAddress))}
	ownerAccount, _ = testContext.Accounts.LoadAccount(ownerAddr)
	firstSCAddress := utils.DoDeploySecond(t, testContext, "../dcdt/testdata/first-contract.wasm", ownerAccount, gasPrice, deployGasLimit, args, big.NewInt(0))

	testContext.TxFeeHandler.CreateBlockStarted(getZeroGasAndFees())
	utils.CleanAccumulatedIntermediateTransactions(t, testContext)

	gasLimit := uint64(2000)
	tx := utils.CreateDCDTTransferTx(0, sndAddr, firstSCAddress, token, big.NewInt(5000), gasPrice, gasLimit)
	tx.Data = []byte(string(tx.Data) + "@" + hex.EncodeToString([]byte("transferToSecondContractRejected")))

	retCode, err := testContext.TxProcessor.ProcessTransaction(tx)
	require.Equal(t, vmcommon.UserError, retCode)
	require.Nil(t, err)

	_, err = testContext.Accounts.Commit()
	require.Nil(t, err)

	utils.CheckDCDTBalance(t, testContext, firstSCAddress, token, big.NewInt(0))
	utils.CheckDCDTBalance(t, testContext, secondSCAddress, token, big.NewInt(0))

	expectedSenderBalance := big.NewInt(99980000)
	utils.TestAccount(t, testContext.Accounts, sndAddr, 1, expectedSenderBalance)

	expectedAccumulatedFees := big.NewInt(20000)
	accumulatedFees := testContext.TxFeeHandler.GetAccumulatedFees()
	require.Equal(t, expectedAccumulatedFees, accumulatedFees)
}

func TestAsyncMultiTransferOnCallback(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	testContext, err := vm.CreatePreparedTxProcessorWithVMs(config.EnableEpochs{
		DynamicGasCostForDataTrieStorageLoadEnableEpoch: integrationTests.UnreachableEpoch,
	}, 1)
	require.Nil(t, err)
	defer testContext.Close()

	ownerAddr := []byte("12345678901234567890123456789010")
	sftTokenID := []byte("SFT-123456")
	sftNonce := uint64(1)
	sftBalance := big.NewInt(1000)
	halfBalance := big.NewInt(500)

	utils.CreateAccountWithDCDTBalance(t, testContext.Accounts, ownerAddr, big.NewInt(1000000000), sftTokenID, sftNonce, sftBalance, uint32(core.SemiFungible))
	utils.CheckDCDTNFTBalance(t, testContext, ownerAddr, sftTokenID, sftNonce, sftBalance)

	ownerAccount, _ := testContext.Accounts.LoadAccount(ownerAddr)
	deployGasLimit := uint64(1000000)
	txGasLimit := uint64(1000000)

	// deploy forwarder
	forwarderAddr := utils.DoDeploySecond(t,
		testContext,
		"../dcdt/testdata/forwarder-raw-0.34.0.wasm",
		ownerAccount,
		gasPrice,
		deployGasLimit,
		nil,
		big.NewInt(0),
	)

	// deploy vault
	ownerAccount, _ = testContext.Accounts.LoadAccount(ownerAddr)
	vaultAddr := utils.DoDeploySecond(t,
		testContext,
		"../dcdt/testdata/vault-0.34.2.wasm",
		ownerAccount,
		gasPrice,
		deployGasLimit,
		nil,
		big.NewInt(0),
	)

	// send the tokens to vault
	ownerAccount, _ = testContext.Accounts.LoadAccount(ownerAddr)
	tx := utils.CreateDCDTNFTTransferTx(
		ownerAccount.GetNonce(),
		ownerAddr,
		vaultAddr,
		sftTokenID,
		sftNonce,
		sftBalance,
		gasPrice,
		txGasLimit,
		"accept_funds",
	)
	retCode, err := testContext.TxProcessor.ProcessTransaction(tx)
	require.Equal(t, vmcommon.Ok, retCode)
	require.Nil(t, err)

	_, err = testContext.Accounts.Commit()
	require.Nil(t, err)

	utils.CheckDCDTNFTBalance(t, testContext, vaultAddr, sftTokenID, sftNonce, sftBalance)

	lenSCRs := len(testContext.GetIntermediateTransactions(t))
	// receive tokens from vault to forwarder on callback
	// receive 500 + 500 of the SFT through multi-transfer
	ownerAccount, _ = testContext.Accounts.LoadAccount(ownerAddr)
	tx = utils.CreateSmartContractCall(
		ownerAccount.GetNonce(),
		ownerAddr,
		forwarderAddr,
		gasPrice,
		txGasLimit,
		"forward_async_retrieve_multi_transfer_funds",
		vaultAddr,
		sftTokenID,
		big.NewInt(int64(sftNonce)).Bytes(),
		halfBalance.Bytes(),
		sftTokenID,
		big.NewInt(int64(sftNonce)).Bytes(),
		halfBalance.Bytes(),
	)
	retCode, err = testContext.TxProcessor.ProcessTransaction(tx)
	require.Equal(t, vmcommon.Ok, retCode)
	require.Nil(t, err)
	require.Equal(t, 1, len(testContext.GetIntermediateTransactions(t))-lenSCRs)

	_, err = testContext.Accounts.Commit()
	require.Nil(t, err)

	utils.CheckDCDTNFTBalance(t, testContext, forwarderAddr, sftTokenID, sftNonce, sftBalance)
}

func TestAsyncMultiTransferOnCallAndOnCallback(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	testContext, err := vm.CreatePreparedTxProcessorWithVMs(config.EnableEpochs{}, 1)
	require.Nil(t, err)
	defer testContext.Close()

	ownerAddr := []byte("12345678901234567890123456789010")
	sftTokenID := []byte("SFT-123456")
	sftNonce := uint64(1)
	sftBalance := big.NewInt(1000)
	halfBalance := big.NewInt(500)

	utils.CreateAccountWithDCDTBalance(t, testContext.Accounts, ownerAddr, big.NewInt(1000000000), sftTokenID, sftNonce, sftBalance, uint32(core.SemiFungible))
	utils.CheckDCDTNFTBalance(t, testContext, ownerAddr, sftTokenID, sftNonce, sftBalance)

	ownerAccount, _ := testContext.Accounts.LoadAccount(ownerAddr)
	deployGasLimit := uint64(1000000)
	txGasLimit := uint64(1000000)

	// deploy forwarder
	forwarderAddr := utils.DoDeploySecond(t,
		testContext,
		"../dcdt/testdata/forwarder-raw-0.34.0.wasm",
		ownerAccount,
		gasPrice,
		deployGasLimit,
		nil,
		big.NewInt(0),
	)

	// deploy vault
	ownerAccount, _ = testContext.Accounts.LoadAccount(ownerAddr)
	vaultAddr := utils.DoDeploySecond(t,
		testContext,
		"../dcdt/testdata/vault-0.34.2.wasm",
		ownerAccount,
		gasPrice,
		deployGasLimit,
		nil,
		big.NewInt(0),
	)

	// set vault roles
	utils.SetDCDTRoles(t, testContext.Accounts, vaultAddr, sftTokenID, [][]byte{
		[]byte(core.DCDTRoleNFTAddQuantity),
		[]byte(core.DCDTRoleNFTCreate),
		[]byte(core.DCDTRoleNFTBurn),
	})
	// set lastNonce for vault
	utils.SetLastNFTNonce(t, testContext.Accounts, vaultAddr, sftTokenID, 1)

	// send the tokens to forwarder
	ownerAccount, _ = testContext.Accounts.LoadAccount(ownerAddr)
	tx := utils.CreateDCDTNFTTransferTx(
		ownerAccount.GetNonce(),
		ownerAddr,
		forwarderAddr,
		sftTokenID,
		sftNonce,
		sftBalance,
		gasPrice,
		txGasLimit,
		"",
	)
	retCode, err := testContext.TxProcessor.ProcessTransaction(tx)
	require.Equal(t, vmcommon.Ok, retCode)
	require.Nil(t, err)

	_, err = testContext.Accounts.Commit()
	require.Nil(t, err)

	utils.CheckDCDTNFTBalance(t, testContext, forwarderAddr, sftTokenID, sftNonce, sftBalance)

	// send tokens to vault, vault burns and creates new ones, sending them on forwarder's callback
	ownerAccount, _ = testContext.Accounts.LoadAccount(ownerAddr)
	tx = utils.CreateSmartContractCall(
		ownerAccount.GetNonce(),
		ownerAddr,
		forwarderAddr,
		gasPrice,
		txGasLimit,
		"forwarder_async_send_and_retrieve_multi_transfer_funds",
		vaultAddr,
		sftTokenID,
		big.NewInt(int64(sftNonce)).Bytes(),
		halfBalance.Bytes(),
		sftTokenID,
		big.NewInt(int64(sftNonce)).Bytes(),
		halfBalance.Bytes(),
	)
	retCode, err = testContext.TxProcessor.ProcessTransaction(tx)
	require.Equal(t, vmcommon.Ok, retCode)
	require.Nil(t, err)

	_, err = testContext.Accounts.Commit()
	require.Nil(t, err)

	utils.CheckDCDTNFTBalance(t, testContext, forwarderAddr, sftTokenID, 2, halfBalance)
	utils.CheckDCDTNFTBalance(t, testContext, forwarderAddr, sftTokenID, 3, halfBalance)
}

func TestSendNFTToContractWith0Function(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	testContext, err := vm.CreatePreparedTxProcessorWithVMs(config.EnableEpochs{}, 1)
	require.Nil(t, err)
	defer testContext.Close()

	ownerAddr := []byte("12345678901234567890123456789010")
	sftTokenID := []byte("SFT-123456")
	sftNonce := uint64(1)
	sftBalance := big.NewInt(1000)

	utils.CreateAccountWithDCDTBalance(t, testContext.Accounts, ownerAddr, big.NewInt(1000000000), sftTokenID, sftNonce, sftBalance, uint32(core.SemiFungible))
	utils.CheckDCDTNFTBalance(t, testContext, ownerAddr, sftTokenID, sftNonce, sftBalance)

	ownerAccount, _ := testContext.Accounts.LoadAccount(ownerAddr)
	deployGasLimit := uint64(1000000)
	txGasLimit := uint64(1000000)

	vaultAddr := utils.DoDeploySecond(t,
		testContext,
		"../dcdt/testdata/vault-0.34.0.wasm",
		ownerAccount,
		gasPrice,
		deployGasLimit,
		nil,
		big.NewInt(0),
	)

	// send the tokens to vault
	ownerAccount, _ = testContext.Accounts.LoadAccount(ownerAddr)
	tx := utils.CreateDCDTNFTTransferTx(
		ownerAccount.GetNonce(),
		ownerAddr,
		vaultAddr,
		sftTokenID,
		sftNonce,
		sftBalance,
		gasPrice,
		txGasLimit,
		"",
	)
	tx.Data = append(tx.Data, []byte("@")...)
	retCode, err := testContext.TxProcessor.ProcessTransaction(tx)
	require.Equal(t, vmcommon.Ok, retCode)
	require.Nil(t, err)

	_, err = testContext.Accounts.Commit()
	require.Nil(t, err)
}

func TestSendNFTToContractWith0FunctionNonPayable(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	testContext, err := vm.CreatePreparedTxProcessorWithVMs(config.EnableEpochs{}, 1)
	require.Nil(t, err)
	defer testContext.Close()

	ownerAddr := []byte("12345678901234567890123456789010")
	sftTokenID := []byte("SFT-123456")
	sftNonce := uint64(1)
	sftBalance := big.NewInt(1000)

	utils.CreateAccountWithDCDTBalance(t, testContext.Accounts, ownerAddr, big.NewInt(1000000000), sftTokenID, sftNonce, sftBalance, uint32(core.SemiFungible))
	utils.CheckDCDTNFTBalance(t, testContext, ownerAddr, sftTokenID, sftNonce, sftBalance)

	ownerAccount, _ := testContext.Accounts.LoadAccount(ownerAddr)
	deployGasLimit := uint64(1000000)
	txGasLimit := uint64(1000000)

	vaultAddr := utils.DoDeployWithMetadata(t,
		testContext,
		"../dcdt/testdata/vault-0.34.0.wasm",
		ownerAccount,
		gasPrice,
		deployGasLimit,
		[]byte("0000"),
		nil,
		big.NewInt(0),
	)

	// send the tokens to vault
	ownerAccount, _ = testContext.Accounts.LoadAccount(ownerAddr)
	tx := utils.CreateDCDTNFTTransferTx(
		ownerAccount.GetNonce(),
		ownerAddr,
		vaultAddr,
		sftTokenID,
		sftNonce,
		sftBalance,
		gasPrice,
		txGasLimit,
		"",
	)
	tx.Data = append(tx.Data, []byte("@")...)
	retCode, err := testContext.TxProcessor.ProcessTransaction(tx)
	require.Equal(t, vmcommon.UserError, retCode)
	require.Equal(t, process.ErrFailedTransaction, err)

	_, err = testContext.Accounts.Commit()
	require.Nil(t, err)
}

func TestAsyncDCDTCallForThirdContractShouldWork(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	testContext, err := vm.CreatePreparedTxProcessorWithVMsMultiShard(0, config.EnableEpochs{}, 1)
	require.Nil(t, err)
	defer testContext.Close()

	function1 := []byte("add_queued_call")
	function2 := []byte("forward_queued_calls")

	localRewaBalance := big.NewInt(100000000)
	ownerAddr := []byte("owner-78901234567890123456789000")
	_, _ = vm.CreateAccount(testContext.Accounts, ownerAddr, 0, localRewaBalance)

	// create an address with DCDT token
	sndAddr := []byte("sender-8901234567890123456789000")

	localDcdtBalance := big.NewInt(100000000)
	dcdtTransferValue := big.NewInt(5000)
	token := []byte("miiutoken")
	utils.CreateAccountWithDCDTBalance(t, testContext.Accounts, sndAddr, localRewaBalance, token, 0, localDcdtBalance, uint32(core.Fungible))

	// deploy contract
	ownerAccount, _ := testContext.Accounts.LoadAccount(ownerAddr)
	deployGasLimit := uint64(50000)
	scAddress := utils.DoDeploySecond(t, testContext, "./testdata/third/third.wasm", ownerAccount, gasPrice, deployGasLimit, nil, big.NewInt(0))

	testContext.TxFeeHandler.CreateBlockStarted(getZeroGasAndFees())
	utils.CleanAccumulatedIntermediateTransactions(t, testContext)

	// execute first call
	gasLimit := uint64(500000)
	tx := utils.CreateDCDTTransferTx(0, sndAddr, scAddress, token, dcdtTransferValue, gasPrice, gasLimit)
	tx.Data = []byte(string(tx.Data) + "@" + hex.EncodeToString(function1) + "@01@" + hex.EncodeToString(scAddress) + "@" + hex.EncodeToString(function2))

	retCode, err := testContext.TxProcessor.ProcessTransaction(tx)
	require.Equal(t, vmcommon.UserError, retCode)
	require.Nil(t, err)

	utils.CheckDCDTBalance(t, testContext, sndAddr, token, localDcdtBalance)
	utils.CheckDCDTBalance(t, testContext, scAddress, token, big.NewInt(0))

	// execute second call
	tx = utils.CreateDCDTTransferTx(1, sndAddr, scAddress, token, big.NewInt(5000), gasPrice, gasLimit)
	tx.Data = []byte(string(tx.Data) + "@" + hex.EncodeToString(function2))

	retCode, err = testContext.TxProcessor.ProcessTransaction(tx)
	require.Equal(t, vmcommon.Ok, retCode)
	require.Nil(t, err)

	_, err = testContext.Accounts.Commit()
	require.Nil(t, err)

	utils.CheckDCDTBalance(t, testContext, sndAddr, token, big.NewInt(0).Sub(localDcdtBalance, dcdtTransferValue))
	utils.CheckDCDTBalance(t, testContext, scAddress, token, dcdtTransferValue)

	// try to recreate the data trie
	scAccount, err := testContext.Accounts.LoadAccount(scAddress)
	require.Nil(t, err)
	userScAccount := scAccount.(state.UserAccountHandler)
	roothash := userScAccount.GetRootHash()
	log.Info("recreating data trie", "roothash", roothash)

	leaves := &common.TrieIteratorChannels{
		LeavesChan: make(chan core.KeyValueHolder, 1),
		ErrChan:    errChan.NewErrChanWrapper(),
	}
	err = testContext.Accounts.GetAllLeaves(leaves, context.Background(), roothash, parsers.NewMainTrieLeafParser())
	require.Nil(t, err)

	for range leaves.LeavesChan {
		// do nothing, just iterate
	}

	err = leaves.ErrChan.ReadFromChanNonBlocking()
	require.Nil(t, err)
}
