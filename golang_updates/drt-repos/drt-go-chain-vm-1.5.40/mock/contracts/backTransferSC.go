package contracts

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/TerraDharitri/drt-go-chain-core/data/vm"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	"github.com/TerraDharitri/drt-go-chain-vm-common/txDataBuilder"
	mock "github.com/TerraDharitri/drt-go-chain-vm/mock/context"
	test "github.com/TerraDharitri/drt-go-chain-vm/testcommon"
	"github.com/TerraDharitri/drt-go-chain-vm/vmhost"
	"github.com/TerraDharitri/drt-go-chain-vm/vmhost/vmhooks"
)

// BackTransfer_ParentCallsChild is an exposed mock contract method
func BackTransfer_ParentCallsChild(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("callChild", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)

		storedResult := []byte("ok")

		testConfig := config.(*test.TestConfig)
		input := test.DefaultTestContractCallInput()
		input.GasProvided = testConfig.GasProvidedToChild
		input.CallerAddr = testConfig.ParentAddress
		input.RecipientAddr = testConfig.ChildAddress
		input.Function = "childFunction"
		returnValue := ExecuteOnDestContextInMockContracts(host, input)
		if returnValue != 0 {
			host.Runtime().FailExecution(fmt.Errorf("return value %d", returnValue))
		}
		managedTypes := host.ManagedTypes()

		arguments := host.Runtime().Arguments()
		if len(arguments) > 0 {
			checkBackTransfers := arguments[0]
			if checkBackTransfers[0] == 1 {
				dcdtTransfers, rewa := managedTypes.GetBackTransfers()
				if len(dcdtTransfers) != 1 {
					host.Runtime().FailExecution(fmt.Errorf("found dcdt transfers %d", len(dcdtTransfers)))
					storedResult = []byte("err")
				}
				if !bytes.Equal(test.DCDTTestTokenName, dcdtTransfers[0].DCDTTokenName) {
					host.Runtime().FailExecution(fmt.Errorf("invalid token name %s", string(dcdtTransfers[0].DCDTTokenName)))
					storedResult = []byte("err")
				}
				if big.NewInt(0).SetUint64(testConfig.DCDTTokensToTransfer).Cmp(dcdtTransfers[0].DCDTValue) != 0 {
					host.Runtime().FailExecution(fmt.Errorf("invalid token value %d", dcdtTransfers[0].DCDTValue.Uint64()))
					storedResult = []byte("err")
				}
				if rewa.Cmp(big.NewInt(testConfig.TransferFromChildToParent)) != 0 {
					host.Runtime().FailExecution(fmt.Errorf("invalid rewa value %d", rewa))
					storedResult = []byte("err")
				}
			}
		}

		_, err := host.Storage().SetStorage(test.ParentKeyA, storedResult)
		if err != nil {
			host.Runtime().FailExecution(err)
		}

		return instance
	})
}

// BackTransfer_ChildMakesAsync is an exposed mock contract method
func BackTransfer_ChildMakesAsync(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("childFunction", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		testConfig := config.(*test.TestConfig)

		callData := txDataBuilder.NewBuilder()
		callData.Func("wasteGas")
		callData.Int64(0)

		err := host.Async().RegisterAsyncCall("testGroup", &vmhost.AsyncCall{
			Status:          vmhost.AsyncCallPending,
			Destination:     testConfig.NephewAddress,
			Data:            callData.ToBytes(),
			ValueBytes:      big.NewInt(0).Bytes(),
			SuccessCallback: testConfig.SuccessCallback,
			ErrorCallback:   testConfig.ErrorCallback,
			GasLimit:        uint64(300),
			GasLocked:       testConfig.GasToLock,
			CallbackClosure: nil,
		})
		if err != nil {
			host.Runtime().FailExecution(err)
		}
		return instance
	})
}

// BackTransfer_ChildCallback is an exposed mock contract method
func BackTransfer_ChildCallback(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("myCallback", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		testConfig := config.(*test.TestConfig)

		valueBytes := big.NewInt(testConfig.TransferFromChildToParent).Bytes()
		err := host.Output().Transfer(
			testConfig.ParentAddress,
			testConfig.ChildAddress, 0, 0, big.NewInt(0).SetBytes(valueBytes), nil, []byte{}, vm.DirectCall)
		if err != nil {
			host.Runtime().FailExecution(err)
		}

		transfer := &vmcommon.DCDTTransfer{
			DCDTValue:      big.NewInt(int64(testConfig.DCDTTokensToTransfer)),
			DCDTTokenName:  test.DCDTTestTokenName,
			DCDTTokenType:  0,
			DCDTTokenNonce: 0,
		}

		ret := vmhooks.TransferDCDTNFTExecuteWithTypedArgs(
			host,
			testConfig.ParentAddress,
			[]*vmcommon.DCDTTransfer{transfer},
			int64(testConfig.GasProvidedToChild),
			nil,
			nil)
		if ret != 0 {
			host.Runtime().FailExecution(fmt.Errorf("Transfer DCDT failed"))
		}

		return instance
	})
}
