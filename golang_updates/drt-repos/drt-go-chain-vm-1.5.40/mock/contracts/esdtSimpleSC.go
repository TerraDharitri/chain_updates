package contracts

import (
	"fmt"
	"math/big"

	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	"github.com/TerraDharitri/drt-go-chain-vm-common/txDataBuilder"
	"github.com/TerraDharitri/drt-go-chain-vm/executor"
	mock "github.com/TerraDharitri/drt-go-chain-vm/mock/context"
	test "github.com/TerraDharitri/drt-go-chain-vm/testcommon"
	"github.com/TerraDharitri/drt-go-chain-vm/vmhost"
	"github.com/TerraDharitri/drt-go-chain-vm/vmhost/vmhooks"
)

// ExecDCDTTransferAndCallChild is an exposed mock contract method
func ExecDCDTTransferAndCallChild(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("execDCDTTransferAndCall", func() *mock.InstanceMock {
		testConfig := config.(*test.TestConfig)
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		err := host.Metering().UseGasBounded(testConfig.GasUsedByParent)
		if err != nil {
			host.Runtime().SetRuntimeBreakpointValue(vmhost.BreakpointOutOfGas)
			return instance
		}

		arguments := host.Runtime().Arguments()
		if len(arguments) < 3 {
			host.Runtime().SignalUserError("need 3 arguments")
			return instance
		}

		input := test.DefaultTestContractCallInput()
		input.CallerAddr = host.Runtime().GetContextAddress()
		input.GasProvided = testConfig.GasProvidedToChild
		input.Arguments = [][]byte{
			test.DCDTTestTokenName,
			big.NewInt(int64(testConfig.DCDTTokensToTransfer)).Bytes(),
		}
		input.Arguments = append(input.Arguments, arguments[2:]...)
		input.RecipientAddr = arguments[0]
		input.Function = string(arguments[1])

		returnValue := ExecuteOnDestContextInMockContracts(host, input)
		if returnValue != 0 {
			host.Runtime().FailExecution(fmt.Errorf("return value %d", returnValue))
		}

		return instance
	})
}

// ExecDCDTTransferWithAPICall is an exposed mock contract method
func ExecDCDTTransferWithAPICall(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("execDCDTTransferWithAPICall", func() *mock.InstanceMock {
		testConfig := config.(*test.TestConfig)
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		err := host.Metering().UseGasBounded(testConfig.GasUsedByParent)
		if err != nil {
			host.Runtime().SetRuntimeBreakpointValue(vmhost.BreakpointOutOfGas)
			return instance
		}

		arguments := host.Runtime().Arguments()
		if len(arguments) != 3 {
			host.Runtime().SignalUserError("need 3 arguments")
			return instance
		}

		input := test.DefaultTestContractCallInput()
		input.CallerAddr = host.Runtime().GetContextAddress()
		input.GasProvided = testConfig.GasProvidedToChild
		input.Arguments = [][]byte{
			test.DCDTTestTokenName,
			big.NewInt(int64(testConfig.DCDTTokensToTransfer)).Bytes(),
			arguments[2],
		}
		input.RecipientAddr = arguments[0]

		functionName := arguments[1]
		args := [][]byte{arguments[2]}

		transfer := &vmcommon.DCDTTransfer{
			DCDTValue:      big.NewInt(int64(testConfig.DCDTTokensToTransfer)),
			DCDTTokenName:  test.DCDTTestTokenName,
			DCDTTokenType:  0,
			DCDTTokenNonce: 0,
		}

		vmhooks.TransferDCDTNFTExecuteWithTypedArgs(
			host,
			input.RecipientAddr,
			[]*vmcommon.DCDTTransfer{transfer},
			int64(testConfig.GasProvidedToChild),
			functionName,
			args)

		return instance
	})
}

// ExecDCDTTransferAndAsyncCallChild is an exposed mock contract method
func ExecDCDTTransferAndAsyncCallChild(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("execDCDTTransferAndAsyncCall", func() *mock.InstanceMock {
		testConfig := config.(*test.TestConfig)
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		err := host.Metering().UseGasBounded(testConfig.GasUsedByParent)
		if err != nil {
			host.Runtime().SetRuntimeBreakpointValue(vmhost.BreakpointOutOfGas)
			return instance
		}

		arguments := host.Runtime().Arguments()
		if len(arguments) < 4 {
			host.Runtime().SignalUserError("need at least 4 arguments")
			return instance
		}

		receiver := arguments[0]
		builtInFunction := arguments[1]
		functionToCallOnChild := arguments[2]
		asyncCallType := arguments[3]
		numberOfBackTransfers := big.NewInt(int64(1)).Bytes()
		if len(arguments) > 4 {
			numberOfBackTransfers = arguments[4]
		}

		callData := txDataBuilder.NewBuilder()
		// function to be called on child
		callData.Func(string(builtInFunction))
		callData.Bytes(test.DCDTTestTokenName)
		callData.Bytes(big.NewInt(int64(testConfig.DCDTTokensToTransfer)).Bytes())
		callData.Bytes(functionToCallOnChild)
		callData.Bytes(numberOfBackTransfers)

		value := big.NewInt(0).Bytes()

		if asyncCallType[0] == 0 {
			err = host.Async().RegisterLegacyAsyncCall(receiver, callData.ToBytes(), value)
		} else {
			callbackName := "callBack"
			if host.Runtime().ValidateCallbackName(callbackName) == executor.ErrFuncNotFound {
				callbackName = ""
			}
			err = host.Async().RegisterAsyncCall("testGroup", &vmhost.AsyncCall{
				Status:          vmhost.AsyncCallPending,
				Destination:     receiver,
				Data:            callData.ToBytes(),
				ValueBytes:      value,
				SuccessCallback: callbackName,
				ErrorCallback:   callbackName,
				GasLimit:        testConfig.GasProvidedToChild,
				GasLocked:       testConfig.GasToLock,
			})
		}

		if err != nil {
			host.Runtime().FailExecution(err)
		}

		return instance
	})
}
