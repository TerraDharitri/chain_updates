package hostCoretest

import (
	"encoding/hex"
	"errors"
	"math/big"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data/dcdt"
	"github.com/TerraDharitri/drt-go-chain-core/data/vm"
	"github.com/TerraDharitri/drt-go-chain-scenario/worldmock"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	"github.com/TerraDharitri/drt-go-chain-vm-common/builtInFunctions"
	"github.com/TerraDharitri/drt-go-chain-vm-common/txDataBuilder"
	"github.com/TerraDharitri/drt-go-chain-vm/config"
	"github.com/TerraDharitri/drt-go-chain-vm/executor"
	contextmock "github.com/TerraDharitri/drt-go-chain-vm/mock/context"
	"github.com/TerraDharitri/drt-go-chain-vm/testcommon"
	test "github.com/TerraDharitri/drt-go-chain-vm/testcommon"
	"github.com/TerraDharitri/drt-go-chain-vm/vmhost"
	"github.com/stretchr/testify/require"
)

//TODO package contains snake case named files, rename those.

func TestExecution_ExecuteOnDestContext_DCDTTransferWithoutExecute(t *testing.T) {
	gasSchedule := config.MakeGasMapForTests()

	world := worldmock.NewMockWorld()
	host := test.NewTestHostBuilder(t).
		WithBlockchainHook(world).
		WithBuiltinFunctions().
		WithGasSchedule(gasSchedule).
		Build()

	code := test.GetTestSCCodeModule("exec-dest-ctx-dcdt/basic", "basic", "../../")
	parentAccount := world.AcctMap.CreateSmartContractAccount(test.UserAddress, test.ParentAddress, code, world)
	parentAccount.Balance = big.NewInt(1000)

	defer func() {
		host.Reset()
	}()

	err := world.BuiltinFuncs.SetTokenData(
		test.ParentAddress,
		test.DCDTTestTokenName,
		0,
		&dcdt.DCDigitalToken{
			Value: big.NewInt(100),
			Type:  uint32(core.Fungible),
		})
	require.Nil(t, err)

	dcdtTokensToTransfer := int64(16)

	input := test.DefaultTestContractCallInput()
	input.Function = "basic_transfer"
	input.GasProvided = 100000
	input.DCDTTransfers = make([]*vmcommon.DCDTTransfer, 1)
	input.DCDTTransfers[0] = &vmcommon.DCDTTransfer{}
	input.DCDTTransfers[0].DCDTValue = big.NewInt(dcdtTokensToTransfer)
	input.DCDTTransfers[0].DCDTTokenName = test.DCDTTestTokenName

	expectedTransfer := txDataBuilder.NewBuilder()
	expectedTransfer.TransferDCDT(string(test.DCDTTestTokenName), dcdtTokensToTransfer)

	expectedTransfers := make([]testcommon.TransferEntry, 0)
	expectedTransfers = append(expectedTransfers,
		test.CreateTransferEntry(test.ParentAddress, test.ParentAddress, 1).
			WithData(expectedTransfer.ToBytes()).
			WithGasLimit(0).
			WithGasLocked(0).
			WithCallType(vm.DirectCall).
			WithValue(big.NewInt(0)))

	vmOutput, err := host.RunSmartContractCall(input)

	verify := test.NewVMOutputVerifier(t, vmOutput, err)
	verify.Ok().
		Transfers(expectedTransfers...)
}

func TestExecution_ExecuteOnDestContext_MockBuiltinFunctions_Claim(t *testing.T) {
	parentGasUsed := uint64(1973)
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("exec-dest-ctx-builtin", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(test.GasProvided).
			WithFunction("callBuiltinClaim").
			Build()).
		WithSetup(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.ProcessBuiltInFunctionCalled = dummyProcessBuiltInFunction
			host.SetBuiltInFunctionsContainer(getDummyBuiltinFunctionsContainer())
		}).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.Ok().
				BalanceDelta(test.ParentAddress, 42).
				GasUsed(test.ParentAddress, parentGasUsed).
				GasRemaining(test.GasProvided - parentGasUsed).
				ReturnData([]byte("succ"))
		})
}

func TestExecution_ExecuteOnDestContext_MockBuiltinFunctions_DoSomething(t *testing.T) {
	parentGasUsed := uint64(1977)
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("exec-dest-ctx-builtin", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(test.GasProvided).
			WithFunction("callBuiltinDoSomething").
			Build()).
		WithSetup(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.ProcessBuiltInFunctionCalled = dummyProcessBuiltInFunction
			host.SetBuiltInFunctionsContainer(getDummyBuiltinFunctionsContainer())
		}).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.Ok().
				BalanceDelta(test.ParentAddress, big.NewInt(0).Sub(vmhost.One, vmhost.One).Int64()).
				GasUsed(test.ParentAddress, parentGasUsed).
				GasRemaining(test.GasProvided - parentGasUsed).
				ReturnData([]byte("succ"))
		})
}

func TestExecution_ExecuteOnDestContext_MockBuiltinFunctions_Nonexistent(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("exec-dest-ctx-builtin", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(test.GasProvided).
			WithFunction("callNonexistingBuiltin").
			Build()).
		WithSetup(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.ProcessBuiltInFunctionCalled = dummyProcessBuiltInFunction
			host.SetBuiltInFunctionsContainer(getDummyBuiltinFunctionsContainer())
		}).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.ExecutionFailed().
				ReturnMessage(executor.ErrFuncNotFound.Error()).
				GasRemaining(0)
		})
}

func TestExecution_ExecuteOnDestContext_MockBuiltinFunctions_Fail(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("exec-dest-ctx-builtin", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(test.GasProvided).
			WithFunction("callBuiltinFail").
			Build()).
		WithSetup(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.ProcessBuiltInFunctionCalled = dummyProcessBuiltInFunction
			host.SetBuiltInFunctionsContainer(getDummyBuiltinFunctionsContainer())
		}).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.ExecutionFailed().
				ReturnMessage("whatdidyoudo").
				GasRemaining(0)
		})
}

func TestExecution_AsyncCall_MockBuiltinFails(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("async-call-builtin", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(test.GasProvided).
			WithFunction("performAsyncCallToBuiltin").
			WithArguments([]byte{1}).
			Build()).
		WithSetup(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.ProcessBuiltInFunctionCalled = dummyProcessBuiltInFunction
			host.SetBuiltInFunctionsContainer(getDummyBuiltinFunctionsContainer())
		}).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.Ok().
				ReturnData([]byte("hello"), []byte{10})
		})
}

func TestDCDT_GettersAPI(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("exchange", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(test.GasProvided).
			WithFunction("validateGetters").
			WithDCDTValue(big.NewInt(5)).
			WithDCDTTokenName([]byte("TT")).
			Build()).
		WithSetup(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.ProcessBuiltInFunctionCalled = dummyProcessBuiltInFunction
			host.SetBuiltInFunctionsContainer(getDummyBuiltinFunctionsContainer())
		}).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.Ok()
		})
}

func TestDCDT_GettersAPI_ExecuteAfterBuiltinCall(t *testing.T) {
	world := worldmock.NewMockWorld()
	host := test.NewTestHostBuilder(t).
		WithBlockchainHook(world).
		WithBuiltinFunctions().
		Build()
	defer func() {
		host.Reset()
	}()

	initialDCDTTokenBalance := uint64(1000)

	// Deploy the "parent" contract, which will call the exchange; the actual
	// code of the contract is not important, because the exchange will be called
	// by the "parent" using a manual call to host.ExecuteOnDestContext().
	dummyCode := test.GetTestSCCode("init-simple", "../../")
	testToken := []byte("TT")
	parentAccount := world.AcctMap.CreateSmartContractAccount(test.UserAddress, test.ParentAddress, dummyCode, world)
	_ = parentAccount.SetTokenBalanceUint64(testToken, 0, initialDCDTTokenBalance)

	// Deploy the exchange contract, which will receive DCDT and verify that it
	// can see the received token amount and token name.
	exchangeAddress := test.MakeTestSCAddress("exchange")
	exchangeCode := test.GetTestSCCode("exchange", "../../")
	exchange := world.AcctMap.CreateSmartContractAccount(test.UserAddress, exchangeAddress, exchangeCode, world)
	exchange.Balance = big.NewInt(1000)

	// Prepare VM to appear as if the parent contract is being executed
	input := test.DefaultTestContractCallInput()
	host.Runtime().InitStateFromContractCallInput(input)
	_ = host.Runtime().StartWasmerInstance(dummyCode, input.GasProvided, true)
	err := host.Output().TransferValueOnly(input.RecipientAddr, input.CallerAddr, input.CallValue, false)
	require.Nil(t, err)

	// Transfer DCDT to the exchange and call its "validateGetters" method
	dcdtValue := int64(5)
	input.CallerAddr = test.ParentAddress
	input.RecipientAddr = exchangeAddress
	input.Function = core.BuiltInFunctionDCDTTransfer
	input.GasProvided = 10000
	input.AsyncArguments = &vmcommon.AsyncArguments{
		CallID:       []byte{},
		CallerCallID: []byte{},
	}
	input.Arguments = [][]byte{
		testToken,
		big.NewInt(dcdtValue).Bytes(),
		[]byte("validateGetters"),
	}

	expectedTransfer := txDataBuilder.NewBuilder()
	expectedTransfer.TransferDCDT(string(test.DCDTTestTokenName), dcdtValue)

	expectedTransfers := make([]testcommon.TransferEntry, 0)
	expectedTransfers = append(expectedTransfers,
		test.CreateTransferEntry(test.ParentAddress, exchangeAddress, 1).
			WithData(expectedTransfer.ToBytes()).
			WithGasLimit(0).
			WithGasLocked(0).
			WithCallType(vm.DirectCall).
			WithValue(big.NewInt(0)))

	vmOutput, _, err := host.ExecuteOnDestContext(input)

	verify := test.NewVMOutputVerifier(t, vmOutput, err)
	verify.Ok().
		Transfers(expectedTransfers...)

	parentDCDTBalance, _ := parentAccount.GetTokenBalanceUint64(testToken, 0)
	require.Equal(t, initialDCDTTokenBalance-uint64(dcdtValue), parentDCDTBalance)
}

func dummyProcessBuiltInFunction(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	outputAccounts := make(map[string]*vmcommon.OutputAccount)
	outputAccounts[string(test.ParentAddress)] = &vmcommon.OutputAccount{
		BalanceDelta: big.NewInt(0),
		Address:      test.ParentAddress}

	if input.Function == "builtinClaim" {
		outputAccounts[string(test.ParentAddress)].BalanceDelta = big.NewInt(42)
		return &vmcommon.VMOutput{
			GasRemaining:   400 + input.GasLocked,
			OutputAccounts: outputAccounts,
		}, nil
	}
	if input.Function == "builtinDoSomething" {
		return &vmcommon.VMOutput{
			GasRemaining:   400 + input.GasLocked,
			OutputAccounts: outputAccounts,
		}, nil
	}
	if input.Function == "builtinFail" {
		return nil, errors.New("whatdidyoudo")
	}
	if input.Function == core.BuiltInFunctionDCDTTransfer {
		vmOutput := &vmcommon.VMOutput{
			GasRemaining: 0,
		}
		function := string(input.Arguments[2])
		dcdtTransferTxData := function
		for _, arg := range input.Arguments[3:] {
			dcdtTransferTxData += "@" + hex.EncodeToString(arg)
		}
		outTransfer := vmcommon.OutputTransfer{
			Value:         big.NewInt(0),
			GasLimit:      input.GasProvided - test.DCDTTransferGasCost + input.GasLocked,
			Data:          []byte(dcdtTransferTxData),
			CallType:      vm.AsynchronousCall,
			SenderAddress: input.CallerAddr,
		}
		vmOutput.OutputAccounts = make(map[string]*vmcommon.OutputAccount)
		vmOutput.OutputAccounts[string(input.RecipientAddr)] = &vmcommon.OutputAccount{
			Address:         input.RecipientAddr,
			OutputTransfers: []vmcommon.OutputTransfer{outTransfer},
		}
		// TODO when DCDT token balance querying is implemented, ensure the
		// transfers that happen here are persisted in the mock accounts
		return vmOutput, nil
	}

	return nil, executor.ErrFuncNotFound
}

func getDummyBuiltinFunctionsContainer() vmcommon.BuiltInFunctionContainer {
	builtInContainer := builtInFunctions.NewBuiltInFunctionContainer()
	_ = builtInContainer.Add("builtinClaim", &test.MockBuiltin{})
	_ = builtInContainer.Add("builtinDoSomething", &test.MockBuiltin{})
	_ = builtInContainer.Add("builtinFail", &test.MockBuiltin{})
	_ = builtInContainer.Add(core.BuiltInFunctionDCDTTransfer, &test.MockBuiltin{})

	return builtInContainer
}
