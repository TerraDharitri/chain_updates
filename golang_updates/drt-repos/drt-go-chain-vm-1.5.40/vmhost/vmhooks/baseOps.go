package vmhooks

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data/dcdt"
	"github.com/TerraDharitri/drt-go-chain-core/data/vm"
	logger "github.com/TerraDharitri/drt-go-chain-logger"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	"github.com/TerraDharitri/drt-go-chain-vm-common/parsers"
	"github.com/TerraDharitri/drt-go-chain-vm/executor"
	"github.com/TerraDharitri/drt-go-chain-vm/math"
	"github.com/TerraDharitri/drt-go-chain-vm/vmhost"
)

const (
	getSCAddressName                 = "getSCAddress"
	getOwnerAddressName              = "getOwnerAddress"
	getShardOfAddressName            = "getShardOfAddress"
	isSmartContractName              = "isSmartContract"
	getExternalBalanceName           = "getExternalBalance"
	blockHashName                    = "blockHash"
	transferValueName                = "transferValue"
	transferDCDTExecuteName          = "transferDCDTExecute"
	transferDCDTNFTExecuteName       = "transferDCDTNFTExecute"
	multiTransferDCDTNFTExecuteName  = "multiTransferDCDTNFTExecute"
	transferValueExecuteName         = "transferValueExecute"
	createAsyncCallName              = "createAsyncCall"
	setAsyncGroupCallbackName        = "setAsyncGroupCallback"
	setAsyncContextCallbackName      = "setAsyncContextCallback"
	getArgumentLengthName            = "getArgumentLength"
	getArgumentName                  = "getArgument"
	getFunctionName                  = "getFunction"
	getNumArgumentsName              = "getNumArguments"
	storageStoreName                 = "storageStore"
	storageLoadLengthName            = "storageLoadLength"
	storageLoadName                  = "storageLoad"
	storageLoadFromAddressName       = "storageLoadFromAddress"
	getCallerName                    = "getCaller"
	checkNoPaymentName               = "checkNoPayment"
	callValueName                    = "callValue"
	getDCDTValueName                 = "getDCDTValue"
	getDCDTTokenNameName             = "getDCDTTokenName"
	getDCDTTokenNonceName            = "getDCDTTokenNonce"
	getDCDTTokenTypeName             = "getDCDTTokenType"
	getCallValueTokenNameName        = "getCallValueTokenName"
	getDCDTValueByIndexName          = "getDCDTValueByIndex"
	getDCDTTokenNameByIndexName      = "getDCDTTokenNameByIndex"
	getDCDTTokenNonceByIndexName     = "getDCDTTokenNonceByIndex"
	getDCDTTokenTypeByIndexName      = "getDCDTTokenTypeByIndex"
	getCallValueTokenNameByIndexName = "getCallValueTokenNameByIndex"
	getNumDCDTTransfersName          = "getNumDCDTTransfers"
	getCurrentDCDTNFTNonceName       = "getCurrentDCDTNFTNonce"
	writeLogName                     = "writeLog"
	writeEventLogName                = "writeEventLog"
	returnDataName                   = "returnData"
	signalErrorName                  = "signalError"
	getGasLeftName                   = "getGasLeft"
	getDCDTBalanceName               = "getDCDTBalance"
	getDCDTNFTNameLengthName         = "getDCDTNFTNameLength"
	getDCDTNFTAttributeLengthName    = "getDCDTNFTAttributeLength"
	getDCDTNFTURILengthName          = "getDCDTNFTURILength"
	getDCDTTokenDataName             = "getDCDTTokenData"
	getDCDTLocalRolesName            = "getDCDTLocalRoles"
	validateTokenIdentifierName      = "validateTokenIdentifier"
	executeOnDestContextName         = "executeOnDestContext"
	executeOnSameContextName         = "executeOnSameContext"
	executeReadOnlyName              = "executeReadOnly"
	createContractName               = "createContract"
	deployFromSourceContractName     = "deployFromSourceContract"
	upgradeContractName              = "upgradeContract"
	upgradeFromSourceContractName    = "upgradeFromSourceContract"
	deleteContractName               = "deleteContract"
	asyncCallName                    = "asyncCall"
	getNumReturnDataName             = "getNumReturnData"
	getReturnDataSizeName            = "getReturnDataSize"
	getReturnDataName                = "getReturnData"
	cleanReturnDataName              = "cleanReturnData"
	deleteFromReturnDataName         = "deleteFromReturnData"
	setStorageLockName               = "setStorageLock"
	getStorageLockName               = "getStorageLock"
	isStorageLockedName              = "isStorageLocked"
	clearStorageLockName             = "clearStorageLock"
	getBlockTimestampName            = "getBlockTimestamp"
	getBlockNonceName                = "getBlockNonce"
	getBlockRoundName                = "getBlockRound"
	getBlockEpochName                = "getBlockEpoch"
	getBlockRandomSeedName           = "getBlockRandomSeed"
	getStateRootHashName             = "getStateRootHash"
	getPrevBlockTimestampName        = "getPrevBlockTimestamp"
	getPrevBlockNonceName            = "getPrevBlockNonce"
	getPrevBlockRoundName            = "getPrevBlockRound"
	getPrevBlockEpochName            = "getPrevBlockEpoch"
	getPrevBlockRandomSeedName       = "getPrevBlockRandomSeed"
	getOriginalTxHashName            = "getOriginalTxHash"
	getCurrentTxHashName             = "getCurrentTxHash"
	getPrevTxHashName                = "getPrevTxHash"
)

type CreateContractCallType int

const (
	CreateContract = iota
	DeployContract
)

var logEEI = logger.GetOrCreate("vm/eei")

func getDCDTTransferFromInputFailIfWrongIndex(host vmhost.VMHost, index int32) *vmcommon.DCDTTransfer {
	dcdtTransfers := host.Runtime().GetVMInput().DCDTTransfers
	if int32(len(dcdtTransfers))-1 < index || index < 0 {
		WithFaultAndHost(host, vmhost.ErrInvalidTokenIndex, host.Runtime().BaseOpsErrorShouldFailExecution())
		return nil
	}
	return dcdtTransfers[index]
}

func failIfMoreThanOneDCDTTransfer(context *VMHooksImpl) bool {
	runtime := context.GetRuntimeContext()
	if len(runtime.GetVMInput().DCDTTransfers) > 1 {
		return context.WithFault(vmhost.ErrTooManyDCDTTransfers, true)
	}
	return false
}

// GetGasLeft VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetGasLeft() int64 {
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetGasLeft
	err := metering.UseGasBoundedAndAddTracedGas(getGasLeftName, gasToUse)
	if context.WithFault(err, context.GetRuntimeContext().BaseOpsErrorShouldFailExecution()) {
		return 0
	}

	return int64(metering.GasLeft())
}

// GetSCAddress VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetSCAddress(resultOffset executor.MemPtr) {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetSCAddress
	err := metering.UseGasBoundedAndAddTracedGas(getSCAddressName, gasToUse)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	owner := runtime.GetContextAddress()
	err = context.MemStore(resultOffset, owner)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}
}

// GetOwnerAddress VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetOwnerAddress(resultOffset executor.MemPtr) {
	blockchain := context.GetBlockchainContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetOwnerAddress
	err := metering.UseGasBoundedAndAddTracedGas(getOwnerAddressName, gasToUse)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	owner, err := blockchain.GetOwnerAddress()
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	err = context.MemStore(resultOffset, owner)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}
}

// GetShardOfAddress VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetShardOfAddress(addressOffset executor.MemPtr) int32 {
	blockchain := context.GetBlockchainContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetShardOfAddress
	err := metering.UseGasBoundedAndAddTracedGas(getShardOfAddressName, gasToUse)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	address, err := context.MemLoad(addressOffset, vmhost.AddressLen)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int32(blockchain.GetShardOfAddress(address))
}

// IsSmartContract VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) IsSmartContract(addressOffset executor.MemPtr) int32 {
	blockchain := context.GetBlockchainContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.IsSmartContract
	err := metering.UseGasBoundedAndAddTracedGas(isSmartContractName, gasToUse)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	address, err := context.MemLoad(addressOffset, vmhost.AddressLen)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	isSmartContract := blockchain.IsSmartContract(address)

	return int32(vmhost.BooleanToInt(isSmartContract))
}

// SignalError VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) SignalError(messageOffset executor.MemPtr, messageLength executor.MemLength) {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(signalErrorName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.SignalError
	gasToUse += metering.GasSchedule().BaseOperationCost.PersistPerByte * uint64(messageLength)

	err := metering.UseGasBounded(gasToUse)
	if err != nil && context.WithFault(err, runtime.UseGasBoundedShouldFailExecution()) {
		return
	}

	message, err := context.MemLoad(messageOffset, messageLength)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}
	runtime.SignalUserError(string(message))
}

// GetExternalBalance VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetExternalBalance(addressOffset executor.MemPtr, resultOffset executor.MemPtr) {
	blockchain := context.GetBlockchainContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetExternalBalance
	err := metering.UseGasBoundedAndAddTracedGas(getExternalBalanceName, gasToUse)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	address, err := context.MemLoad(addressOffset, vmhost.AddressLen)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	balance := blockchain.GetBalance(address)

	err = context.MemStore(resultOffset, balance)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}
}

// GetBlockHash VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetBlockHash(nonce int64, resultOffset executor.MemPtr) int32 {
	blockchain := context.GetBlockchainContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockHash
	err := metering.UseGasBoundedAndAddTracedGas(blockHashName, gasToUse)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	hash := blockchain.BlockHash(uint64(nonce))
	err = context.MemStore(resultOffset, hash)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

func getDCDTDataFromBlockchainHook(
	context *VMHooksImpl,
	addressOffset executor.MemPtr,
	tokenIDOffset executor.MemPtr,
	tokenIDLen executor.MemLength,
	nonce int64,
) (*dcdt.DCDigitalToken, error) {
	metering := context.GetMeteringContext()
	blockchain := context.GetBlockchainContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetExternalBalance
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return nil, err
	}

	address, err := context.MemLoad(addressOffset, vmhost.AddressLen)
	if err != nil {
		return nil, err
	}

	tokenID, err := context.MemLoad(tokenIDOffset, tokenIDLen)
	if err != nil {
		return nil, err
	}

	dcdtToken, err := blockchain.GetDCDTToken(address, tokenID, uint64(nonce))
	if err != nil {
		return nil, err
	}

	return dcdtToken, nil
}

// GetDCDTBalance VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetDCDTBalance(
	addressOffset executor.MemPtr,
	tokenIDOffset executor.MemPtr,
	tokenIDLen executor.MemLength,
	nonce int64,
	resultOffset executor.MemPtr,
) int32 {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(getDCDTBalanceName)

	dcdtData, err := getDCDTDataFromBlockchainHook(context, addressOffset, tokenIDOffset, tokenIDLen, nonce)

	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}
	err = context.MemStore(resultOffset, dcdtData.Value.Bytes())
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(dcdtData.Value.Bytes()))
}

// GetDCDTNFTNameLength VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetDCDTNFTNameLength(
	addressOffset executor.MemPtr,
	tokenIDOffset executor.MemPtr,
	tokenIDLen executor.MemLength,
	nonce int64,
) int32 {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(getDCDTNFTNameLengthName)

	dcdtData, err := getDCDTDataFromBlockchainHook(context, addressOffset, tokenIDOffset, tokenIDLen, nonce)

	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}
	if dcdtData == nil || dcdtData.TokenMetaData == nil {
		context.WithFault(vmhost.ErrNilDCDTData, runtime.BaseOpsErrorShouldFailExecution())
		return 0
	}

	return int32(len(dcdtData.TokenMetaData.Name))
}

// GetDCDTNFTAttributeLength VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetDCDTNFTAttributeLength(
	addressOffset executor.MemPtr,
	tokenIDOffset executor.MemPtr,
	tokenIDLen executor.MemLength,
	nonce int64,
) int32 {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(getDCDTNFTAttributeLengthName)

	dcdtData, err := getDCDTDataFromBlockchainHook(context, addressOffset, tokenIDOffset, tokenIDLen, nonce)

	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}
	if dcdtData == nil || dcdtData.TokenMetaData == nil {
		context.WithFault(vmhost.ErrNilDCDTData, runtime.BaseOpsErrorShouldFailExecution())
		return 0
	}

	return int32(len(dcdtData.TokenMetaData.Attributes))
}

// GetDCDTNFTURILength VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetDCDTNFTURILength(
	addressOffset executor.MemPtr,
	tokenIDOffset executor.MemPtr,
	tokenIDLen executor.MemLength,
	nonce int64,
) int32 {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(getDCDTNFTURILengthName)

	dcdtData, err := getDCDTDataFromBlockchainHook(context, addressOffset, tokenIDOffset, tokenIDLen, nonce)

	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}
	if dcdtData == nil || dcdtData.TokenMetaData == nil {
		context.WithFault(vmhost.ErrNilDCDTData, runtime.BaseOpsErrorShouldFailExecution())
		return 0
	}
	if len(dcdtData.TokenMetaData.URIs) == 0 {
		return 0
	}

	return int32(len(dcdtData.TokenMetaData.URIs[0]))
}

// GetDCDTTokenData VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetDCDTTokenData(
	addressOffset executor.MemPtr,
	tokenIDOffset executor.MemPtr,
	tokenIDLen executor.MemLength,
	nonce int64,
	valueHandle int32,
	propertiesOffset executor.MemPtr,
	hashOffset executor.MemPtr,
	nameOffset executor.MemPtr,
	attributesOffset executor.MemPtr,
	creatorOffset executor.MemPtr,
	royaltiesHandle int32,
	urisOffset executor.MemPtr,
) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(getDCDTTokenDataName)

	dcdtData, err := getDCDTDataFromBlockchainHook(context, addressOffset, tokenIDOffset, tokenIDLen, nonce)

	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	value := managedType.GetBigIntOrCreate(valueHandle)
	value.Set(dcdtData.Value)

	err = context.MemStore(propertiesOffset, dcdtData.Properties)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	if dcdtData.TokenMetaData != nil {
		err = context.MemStore(hashOffset, dcdtData.TokenMetaData.Hash)
		if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
			return -1
		}
		err = context.MemStore(nameOffset, dcdtData.TokenMetaData.Name)
		if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
			return -1
		}
		err = context.MemStore(attributesOffset, dcdtData.TokenMetaData.Attributes)
		if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
			return -1
		}
		err = context.MemStore(creatorOffset, dcdtData.TokenMetaData.Creator)
		if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
			return -1
		}

		royalties := managedType.GetBigIntOrCreate(royaltiesHandle)
		royalties.SetUint64(uint64(dcdtData.TokenMetaData.Royalties))

		if len(dcdtData.TokenMetaData.URIs) > 0 {
			err = context.MemStore(urisOffset, dcdtData.TokenMetaData.URIs[0])
			if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
				return -1
			}
		}
	}
	return int32(len(dcdtData.Value.Bytes()))
}

// GetDCDTLocalRoles VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetDCDTLocalRoles(tokenIdHandle int32) int64 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	storage := context.GetStorageContext()
	metering := context.GetMeteringContext()

	tokenID, err := managedType.GetBytes(tokenIdHandle)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	dcdtRoleKeyPrefix := []byte(core.ProtectedKeyPrefix + core.DCDTRoleIdentifier + core.DCDTKeyIdentifier)
	key := []byte(string(dcdtRoleKeyPrefix) + string(tokenID))

	data, trieDepth, usedCache, err := storage.GetStorage(key)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	err = storage.UseGasForStorageLoad(
		storageLoadName,
		int64(trieDepth),
		metering.GasSchedule().BaseOpsAPICost.StorageLoad,
		usedCache)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	enableEpochsHandler := context.host.EnableEpochsHandler()
	return getDCDTRoles(data, enableEpochsHandler.IsFlagEnabled(vmhost.CryptoOpcodesV2Flag))
}

// ValidateTokenIdentifier VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ValidateTokenIdentifier(
	tokenIdHandle int32,
) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetArgument
	err := metering.UseGasBoundedAndAddTracedGas(validateTokenIdentifierName, gasToUse)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	tokenID, err := managedType.GetBytes(tokenIdHandle)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	if ValidateToken(tokenID) {
		return 1
	} else {
		return 0
	}

}

// TransferValue VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) TransferValue(
	destOffset executor.MemPtr,
	valueOffset executor.MemPtr,
	dataOffset executor.MemPtr,
	length executor.MemLength) int32 {

	host := context.GetVMHost()
	runtime := host.Runtime()
	metering := host.Metering()
	output := host.Output()
	metering.StartGasTracing(transferValueName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.TransferValue
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.UseGasBoundedShouldFailExecution()) {
		return 1
	}

	sender := runtime.GetContextAddress()
	dest, err := context.MemLoad(destOffset, vmhost.AddressLen)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	valueBytes, err := context.MemLoad(valueOffset, vmhost.BalanceLen)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(length))
	err = metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.UseGasBoundedShouldFailExecution()) {
		return 1
	}

	data, err := context.MemLoad(dataOffset, length)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	if host.IsBuiltinFunctionCall(data) {
		context.WithFault(vmhost.ErrTransferValueOnDCDTCall, runtime.BaseOpsErrorShouldFailExecution())
		return 1
	}

	err = output.Transfer(dest, sender, 0, 0, big.NewInt(0).SetBytes(valueBytes), nil, data, vm.DirectCall)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

type indirectContractCallArguments struct {
	dest      []byte
	value     *big.Int
	function  []byte
	args      [][]byte
	actualLen int32
}

func (context *VMHooksImpl) extractIndirectContractCallArgumentsWithValue(
	host vmhost.VMHost,
	destOffset executor.MemPtr,
	valueOffset executor.MemPtr,
	functionOffset executor.MemPtr,
	functionLength executor.MemLength,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) (*indirectContractCallArguments, error) {
	return context.extractIndirectContractCallArguments(
		host,
		destOffset,
		valueOffset,
		true,
		functionOffset,
		functionLength,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
}

func (context *VMHooksImpl) extractIndirectContractCallArgumentsWithoutValue(
	host vmhost.VMHost,
	destOffset executor.MemPtr,
	functionOffset executor.MemPtr,
	functionLength executor.MemLength,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) (*indirectContractCallArguments, error) {
	return context.extractIndirectContractCallArguments(
		host,
		destOffset,
		0,
		false,
		functionOffset,
		functionLength,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
}

func (context *VMHooksImpl) extractIndirectContractCallArguments(
	host vmhost.VMHost,
	destOffset executor.MemPtr,
	valueOffset executor.MemPtr,
	hasValueOffset bool,
	functionOffset executor.MemPtr,
	functionLength executor.MemLength,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) (*indirectContractCallArguments, error) {
	metering := host.Metering()

	dest, err := context.MemLoad(destOffset, vmhost.AddressLen)
	if err != nil {
		return nil, err
	}

	var value *big.Int

	if hasValueOffset {
		valueBytes, err := context.MemLoad(valueOffset, vmhost.BalanceLen)
		if err != nil {
			return nil, err
		}
		value = big.NewInt(0).SetBytes(valueBytes)
	}

	function, err := context.MemLoad(functionOffset, functionLength)
	if err != nil {
		return nil, err
	}

	args, actualLen, err := context.getArgumentsFromMemory(
		host,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
	if err != nil {
		return nil, err
	}

	gasToUse := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(actualLen))
	err = metering.UseGasBounded(gasToUse)
	if err != nil && host.Runtime().UseGasBoundedShouldFailExecution() {
		return nil, err
	}

	return &indirectContractCallArguments{
		dest:      dest,
		value:     value,
		function:  function,
		args:      args,
		actualLen: actualLen,
	}, nil
}

// TransferValueExecute VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) TransferValueExecute(
	destOffset executor.MemPtr,
	valueOffset executor.MemPtr,
	gasLimit int64,
	functionOffset executor.MemPtr,
	functionLength executor.MemLength,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) int32 {
	host := context.GetVMHost()
	return context.TransferValueExecuteWithHost(
		host,
		destOffset,
		valueOffset,
		gasLimit,
		functionOffset,
		functionLength,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
}

// TransferValueExecuteWithHost - transferValueExecute with host instead of pointer context
func (context *VMHooksImpl) TransferValueExecuteWithHost(
	host vmhost.VMHost,
	destOffset executor.MemPtr,
	valueOffset executor.MemPtr,
	gasLimit int64,
	functionOffset executor.MemPtr,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) int32 {
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(transferValueExecuteName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.TransferValue
	err := metering.UseGasBounded(gasToUse)
	if WithFaultAndHost(host, err, runtime.UseGasBoundedShouldFailExecution()) {
		return 1
	}

	callArgs, err := context.extractIndirectContractCallArgumentsWithValue(
		host, destOffset, valueOffset, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	return TransferValueExecuteWithTypedArgs(
		host,
		callArgs.dest,
		callArgs.value,
		gasLimit,
		callArgs.function,
		callArgs.args,
	)
}

// TransferValueExecuteWithTypedArgs - transferValueExecute with args already read from memory
func TransferValueExecuteWithTypedArgs(
	host vmhost.VMHost,
	dest []byte,
	value *big.Int,
	gasLimit int64,
	function []byte,
	args [][]byte,
) int32 {
	runtime := host.Runtime()
	metering := host.Metering()
	output := host.Output()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.TransferValue
	err := metering.UseGasBounded(gasToUse)
	if WithFaultAndHost(host, err, runtime.UseGasBoundedShouldFailExecution()) {
		return 1
	}

	sender := runtime.GetContextAddress()

	var contractCallInput *vmcommon.ContractCallInput
	if len(function) > 0 {
		contractCallInput, err = prepareIndirectContractCallInput(
			host,
			sender,
			value,
			gasLimit,
			dest,
			function,
			args,
			gasToUse,
			false,
		)
		if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
			return 1
		}
	}

	if contractCallInput != nil {
		if host.IsBuiltinFunctionName(contractCallInput.Function) {
			WithFaultAndHost(host, vmhost.ErrNilDCDTData, runtime.BaseOpsErrorShouldFailExecution())
			return 1
		}
	}

	data := ""
	if contractCallInput != nil {
		data = makeCrossShardCallFromInput(contractCallInput.Function, contractCallInput.Arguments)
	}

	lastRound := host.Blockchain().LastRound()
	if host.IsBuiltinFunctionCall([]byte(data)) &&
		lastRound >= uint64(host.EnableEpochsHandler().GetActivationEpoch(vmhost.CheckBuiltInCallOnTransferValueAndFailExecutionFlag)) {
		WithFaultAndHost(host, vmhost.ErrTransferValueOnDCDTCall, runtime.BaseOpsErrorShouldFailExecution())
		return 1
	}

	if host.AreInSameShard(sender, dest) && contractCallInput != nil && host.Blockchain().IsSmartContract(dest) {
		logEEI.Trace("rEWA pre-transfer execution begin")
		vmOutput, err := executeOnDestContextFromAPI(host, contractCallInput)
		if err != nil {
			logEEI.Trace("rEWA pre-transfer execution failed", "error", err)
			WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution())
			return 1
		}
		host.CompleteLogEntriesWithCallType(vmOutput, vmhost.TransferAndExecuteString)

		return 0
	}

	err = metering.UseGasBounded(uint64(gasLimit))
	if WithFaultAndHost(host, err, runtime.UseGasBoundedShouldFailExecution()) {
		return 1
	}

	err = output.Transfer(dest, sender, uint64(gasLimit), 0, value, nil, []byte(data), vm.DirectCall)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

func makeCrossShardCallFromInput(function string, arguments [][]byte) string {
	txData := function
	for _, arg := range arguments {
		txData += "@" + hex.EncodeToString(arg)
	}

	return txData
}

// TransferDCDTExecute VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) TransferDCDTExecute(
	destOffset executor.MemPtr,
	tokenIDOffset executor.MemPtr,
	tokenIDLen executor.MemLength,
	valueOffset executor.MemPtr,
	gasLimit int64,
	functionOffset executor.MemPtr,
	functionLength executor.MemLength,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) int32 {

	return context.TransferDCDTNFTExecute(destOffset, tokenIDOffset, tokenIDLen, valueOffset, 0,
		gasLimit, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)
}

// TransferDCDTNFTExecute VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) TransferDCDTNFTExecute(
	destOffset executor.MemPtr,
	tokenIDOffset executor.MemPtr,
	tokenIDLen executor.MemLength,
	valueOffset executor.MemPtr,
	nonce int64,
	gasLimit int64,
	functionOffset executor.MemPtr,
	functionLength executor.MemLength,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) int32 {
	host := context.GetVMHost()
	metering := host.Metering()
	metering.StartGasTracing(transferDCDTNFTExecuteName)
	return context.TransferDCDTNFTExecuteWithHost(
		host,
		destOffset,
		tokenIDOffset,
		tokenIDLen,
		valueOffset,
		nonce,
		gasLimit,
		functionOffset,
		functionLength,
		numArguments,
		argumentsLengthOffset,
		dataOffset)
}

// MultiTransferDCDTNFTExecute VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MultiTransferDCDTNFTExecute(
	destOffset executor.MemPtr,
	numTokenTransfers int32,
	tokenTransfersArgsLengthOffset executor.MemPtr,
	tokenTransferDataOffset executor.MemPtr,
	gasLimit int64,
	functionOffset executor.MemPtr,
	functionLength executor.MemLength,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) int32 {
	host := context.GetVMHost()
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(multiTransferDCDTNFTExecuteName)

	if numTokenTransfers == 0 {
		_ = WithFaultAndHost(host, vmhost.ErrFailedTransfer, runtime.BaseOpsErrorShouldFailExecution())
		return 1
	}

	callArgs, err := context.extractIndirectContractCallArgumentsWithoutValue(
		host, destOffset, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	gasToUse := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(callArgs.actualLen))
	err = metering.UseGasBounded(gasToUse)
	if WithFaultAndHost(host, err, runtime.UseGasBoundedShouldFailExecution()) {
		return 1
	}

	transferArgs, actualLen, err := context.getArgumentsFromMemory(
		host,
		numTokenTransfers*parsers.ArgsPerTransfer,
		tokenTransfersArgsLengthOffset,
		tokenTransferDataOffset,
	)

	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(actualLen))
	err = metering.UseGasBounded(gasToUse)
	if WithFaultAndHost(host, err, runtime.UseGasBoundedShouldFailExecution()) {
		return 1
	}

	transfers := make([]*vmcommon.DCDTTransfer, numTokenTransfers)
	for i := int32(0); i < numTokenTransfers; i++ {
		tokenStartIndex := i * parsers.ArgsPerTransfer
		transfer := &vmcommon.DCDTTransfer{
			DCDTTokenName:  transferArgs[tokenStartIndex],
			DCDTTokenNonce: big.NewInt(0).SetBytes(transferArgs[tokenStartIndex+1]).Uint64(),
			DCDTValue:      big.NewInt(0).SetBytes(transferArgs[tokenStartIndex+2]),
			DCDTTokenType:  uint32(core.Fungible),
		}
		if transfer.DCDTTokenNonce > 0 {
			transfer.DCDTTokenType = uint32(core.NonFungible)
		}
		transfers[i] = transfer
	}

	return TransferDCDTNFTExecuteWithTypedArgs(
		host,
		callArgs.dest,
		transfers,
		gasLimit,
		callArgs.function,
		callArgs.args,
	)
}

// TransferDCDTNFTExecuteWithHost contains only memory reading of arguments
func (context *VMHooksImpl) TransferDCDTNFTExecuteWithHost(
	host vmhost.VMHost,
	destOffset executor.MemPtr,
	tokenIDOffset executor.MemPtr,
	tokenIDLen executor.MemLength,
	valueOffset executor.MemPtr,
	nonce int64,
	gasLimit int64,
	functionOffset executor.MemPtr,
	functionLength executor.MemLength,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) int32 {
	runtime := host.Runtime()
	metering := host.Metering()

	tokenIdentifier, executeErr := context.MemLoad(tokenIDOffset, tokenIDLen)
	if WithFaultAndHost(host, executeErr, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	callArgs, err := context.extractIndirectContractCallArgumentsWithValue(
		host, destOffset, valueOffset, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	gasToUse := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(callArgs.actualLen))
	err = metering.UseGasBounded(gasToUse)
	if WithFaultAndHost(host, err, runtime.UseGasBoundedShouldFailExecution()) {
		return 1
	}

	transfer := &vmcommon.DCDTTransfer{
		DCDTValue:      callArgs.value,
		DCDTTokenName:  tokenIdentifier,
		DCDTTokenNonce: uint64(nonce),
		DCDTTokenType:  uint32(core.Fungible),
	}
	if nonce > 0 {
		transfer.DCDTTokenType = uint32(core.NonFungible)
	}
	return TransferDCDTNFTExecuteWithTypedArgs(
		host,
		callArgs.dest,
		[]*vmcommon.DCDTTransfer{transfer},
		gasLimit,
		callArgs.function,
		callArgs.args,
	)
}

// TransferDCDTNFTExecuteWithTypedArgs defines the actual transfer DCDT execute logic
func TransferDCDTNFTExecuteWithTypedArgs(
	host vmhost.VMHost,
	dest []byte,
	transfers []*vmcommon.DCDTTransfer,
	gasLimit int64,
	function []byte,
	data [][]byte,
) int32 {
	var executeErr error

	runtime := host.Runtime()
	metering := host.Metering()

	output := host.Output()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.TransferValue * uint64(len(transfers))
	err := metering.UseGasBounded(gasToUse)
	if WithFaultAndHost(host, err, runtime.UseGasBoundedShouldFailExecution()) {
		return 1
	}

	sender := runtime.GetContextAddress()

	var contractCallInput *vmcommon.ContractCallInput
	if len(function) > 0 {
		contractCallInput, executeErr = prepareIndirectContractCallInput(
			host,
			sender,
			big.NewInt(0),
			gasLimit,
			dest,
			function,
			data,
			gasToUse,
			false,
		)
		if WithFaultAndHost(host, executeErr, runtime.SyncExecAPIErrorShouldFailExecution()) {
			return 1
		}

		contractCallInput.DCDTTransfers = transfers
	}

	snapshotBeforeTransfer := host.Blockchain().GetSnapshot()

	originalCaller := host.Runtime().GetOriginalCallerAddress()
	transfersArgs := &vmhost.DCDTTransfersArgs{
		Destination:    dest,
		OriginalCaller: originalCaller,
		Sender:         sender,
		Transfers:      transfers,
		SenderForExec:  sender,
	}
	gasLimitForExec, executeErr := output.TransferDCDT(transfersArgs, contractCallInput)
	if WithFaultAndHost(host, executeErr, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	if host.AreInSameShard(sender, dest) && contractCallInput != nil && host.Blockchain().IsSmartContract(dest) {
		contractCallInput.GasProvided = gasLimitForExec
		contractCallInput.CallerAddr = sender
		logEEI.Trace("DCDT post-transfer execution begin")
		_, executeErr := executeOnDestContextFromAPI(host, contractCallInput)
		if executeErr != nil {
			logEEI.Trace("DCDT post-transfer execution failed", "error", executeErr)
			host.Blockchain().RevertToSnapshot(snapshotBeforeTransfer)
			WithFaultAndHost(host, executeErr, runtime.BaseOpsErrorShouldFailExecution())
			return 1
		}

		return 0
	}

	return 0

}

// TransferDCDTNFTExecuteByUserWithTypedArgs defines the actual transfer DCDT execute logic and execution
func TransferDCDTNFTExecuteByUserWithTypedArgs(
	host vmhost.VMHost,
	callerForExecution []byte,
	dest []byte,
	transfers []*vmcommon.DCDTTransfer,
	gasLimit int64,
	function []byte,
	data [][]byte,
) int32 {
	var executeErr error

	runtime := host.Runtime()
	metering := host.Metering()

	output := host.Output()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.TransferValue * uint64(len(transfers))
	err := metering.UseGasBounded(gasToUse)
	if WithFaultAndHost(host, err, runtime.SyncExecAPIErrorShouldFailExecution()) {
		return 1
	}

	sender := runtime.GetContextAddress()

	var contractCallInput *vmcommon.ContractCallInput
	if len(function) > 0 {
		contractCallInput, executeErr = prepareIndirectContractCallInput(
			host,
			sender,
			big.NewInt(0),
			gasLimit,
			dest,
			function,
			data,
			gasToUse,
			false,
		)
		if WithFaultAndHost(host, executeErr, runtime.SyncExecAPIErrorShouldFailExecution()) {
			return 1
		}

		contractCallInput.DCDTTransfers = transfers
	}

	originalCaller := host.Runtime().GetOriginalCallerAddress()
	transfersArgs := &vmhost.DCDTTransfersArgs{
		Destination:    dest,
		OriginalCaller: originalCaller,
		Sender:         sender,
		Transfers:      transfers,
		SenderForExec:  callerForExecution,
	}
	gasLimitForExec, executeErr := output.TransferDCDT(transfersArgs, contractCallInput)
	if WithFaultAndHost(host, executeErr, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	if host.AreInSameShard(sender, dest) && contractCallInput != nil && host.Blockchain().IsSmartContract(dest) {
		contractCallInput.GasProvided = gasLimitForExec
		contractCallInput.CallerAddr = callerForExecution
		logEEI.Trace("DCDT post-transfer execution begin")
		_, executeErr = executeOnDestContextFromAPI(host, contractCallInput)
		if executeErr != nil {
			logEEI.Trace("DCDT post-transfer execution failed, started transfer to user", "error", executeErr)

			// in case of failed execution, the funds have to be moved to the user
			returnTransferArgs := &vmhost.DCDTTransfersArgs{
				Destination:      callerForExecution,
				OriginalCaller:   originalCaller,
				Sender:           dest,
				Transfers:        transfers,
				SenderForExec:    dest,
				ReturnAfterError: true,
			}
			_, executeErr = output.TransferDCDT(returnTransferArgs, nil)
			if WithFaultAndHost(host, executeErr, runtime.BaseOpsErrorShouldFailExecution()) {
				return 1
			}

			return 1
		}

		return 0
	}

	return 0
}

// CreateAsyncCall VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) CreateAsyncCall(
	destOffset executor.MemPtr,
	valueOffset executor.MemPtr,
	dataOffset executor.MemPtr,
	dataLength executor.MemLength,
	successOffset executor.MemPtr,
	successLength executor.MemLength,
	errorOffset executor.MemPtr,
	errorLength executor.MemLength,
	gas int64,
	extraGasForCallback int64,
) int32 {
	host := context.GetVMHost()
	return context.CreateAsyncCallWithHost(
		host,
		destOffset,
		valueOffset,
		dataOffset,
		dataLength,
		successOffset,
		successLength,
		errorOffset,
		errorLength,
		gas,
		extraGasForCallback)
}

// CreateAsyncCallWithHost - createAsyncCall with host instead of pointer
func (context *VMHooksImpl) CreateAsyncCallWithHost(host vmhost.VMHost,
	destOffset executor.MemPtr,
	valueOffset executor.MemPtr,
	dataOffset executor.MemPtr,
	dataLength executor.MemLength,
	successOffset executor.MemPtr,
	successLength executor.MemLength,
	errorOffset executor.MemPtr,
	errorLength executor.MemLength,
	gas int64,
	extraGasForCallback int64,
) int32 {
	runtime := host.Runtime()

	calledSCAddress, err := context.MemLoad(destOffset, vmhost.AddressLen)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	value, err := context.MemLoad(valueOffset, vmhost.BalanceLen)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	data, err := context.MemLoad(dataOffset, dataLength)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	successFunc, err := context.MemLoad(successOffset, successLength)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	errorFunc, err := context.MemLoad(errorOffset, errorLength)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	return CreateAsyncCallWithTypedArgs(host,
		calledSCAddress,
		value,
		data,
		successFunc,
		errorFunc,
		gas,
		extraGasForCallback,
		nil)
}

// CreateAsyncCallWithTypedArgs - createAsyncCall with arguments already read from memory
func CreateAsyncCallWithTypedArgs(host vmhost.VMHost,
	calledSCAddress []byte,
	value []byte,
	data []byte,
	successFunc []byte,
	errorFunc []byte,
	gas int64,
	extraGasForCallback int64,
	callbackClosure []byte) int32 {

	metering := host.Metering()
	runtime := host.Runtime()
	async := host.Async()

	metering.StartGasTracing(createAsyncCallName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateAsyncCall
	err := metering.UseGasBounded(gasToUse)
	if WithFaultAndHost(host, err, runtime.UseGasBoundedShouldFailExecution()) {
		return 1
	}

	asyncCall := &vmhost.AsyncCall{
		Status:          vmhost.AsyncCallPending,
		Destination:     calledSCAddress,
		Data:            data,
		ValueBytes:      value,
		GasLimit:        uint64(gas),
		SuccessCallback: string(successFunc),
		ErrorCallback:   string(errorFunc),
		GasLocked:       uint64(extraGasForCallback),
		CallbackClosure: callbackClosure,
	}

	if asyncCall.HasDefinedAnyCallback() {
		gasToUse = metering.GasSchedule().BaseOpsAPICost.SetAsyncCallback
		err = metering.UseGasBounded(gasToUse)
		if WithFaultAndHost(host, err, runtime.UseGasBoundedShouldFailExecution()) {
			return 1
		}
	}

	err = async.RegisterAsyncCall("", asyncCall)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

// SetAsyncContextCallback VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) SetAsyncContextCallback(
	callback executor.MemPtr,
	callbackLength executor.MemLength,
	data executor.MemPtr,
	dataLength executor.MemLength,
	gas int64,
) int32 {
	host := context.GetVMHost()
	runtime := host.Runtime()
	metering := host.Metering()
	async := host.Async()
	metering.StartGasTracing(setAsyncContextCallbackName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.SetAsyncContextCallback
	err := metering.UseGasBounded(gasToUse)
	if WithFaultAndHost(host, err, runtime.UseGasBoundedShouldFailExecution()) {
		return 1
	}

	callbackNameBytes, err := context.MemLoad(callback, callbackLength)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	dataBytes, err := context.MemLoad(data, dataLength)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	err = async.SetContextCallback(
		string(callbackNameBytes),
		dataBytes,
		uint64(gas))
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

// UpgradeContract VMHooks implementation.
// @autogenerate(VMHooks)
// @autogenerate(VMHooks)
func (context *VMHooksImpl) UpgradeContract(
	destOffset executor.MemPtr,
	gasLimit int64,
	valueOffset executor.MemPtr,
	codeOffset executor.MemPtr,
	codeMetadataOffset executor.MemPtr,
	length executor.MemLength,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) {
	host := context.GetVMHost()
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(upgradeContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	err := metering.UseGasBounded(gasToUse)
	if WithFaultAndHost(host, err, runtime.UseGasBoundedShouldFailExecution()) {
		return
	}

	value, err := context.MemLoad(valueOffset, vmhost.BalanceLen)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	code, err := context.MemLoad(codeOffset, length)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	codeMetadata, err := context.MemLoad(codeMetadataOffset, vmhost.CodeMetadataLen)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	data, actualLen, err := context.getArgumentsFromMemory(
		host,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
	if WithFaultAndHost(host, err, runtime.UseGasBoundedShouldFailExecution()) {
		return
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(actualLen))
	err = metering.UseGasBounded(gasToUse)
	if WithFaultAndHost(host, err, runtime.UseGasBoundedShouldFailExecution()) {
		return
	}

	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	calledSCAddress, err := context.MemLoad(destOffset, vmhost.AddressLen)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	gasSchedule := metering.GasSchedule()
	gasToUse = math.MulUint64(gasSchedule.BaseOperationCost.DataCopyPerByte, uint64(length))
	err = metering.UseGasBounded(gasToUse)
	if WithFaultAndHost(host, err, runtime.UseGasBoundedShouldFailExecution()) {
		return
	}

	upgradeContract(host, calledSCAddress, code, codeMetadata, value, data, gasLimit)
}

// UpgradeFromSourceContract VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) UpgradeFromSourceContract(
	destOffset executor.MemPtr,
	gasLimit int64,
	valueOffset executor.MemPtr,
	sourceContractAddressOffset executor.MemPtr,
	codeMetadataOffset executor.MemPtr,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) {
	host := context.GetVMHost()
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(upgradeFromSourceContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	err := metering.UseGasBounded(gasToUse)
	if WithFaultAndHost(host, err, runtime.UseGasBoundedShouldFailExecution()) {
		return
	}

	value, err := context.MemLoad(valueOffset, vmhost.BalanceLen)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	sourceContractAddress, err := context.MemLoad(sourceContractAddressOffset, vmhost.AddressLen)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	codeMetadata, err := context.MemLoad(codeMetadataOffset, vmhost.CodeMetadataLen)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	data, actualLen, err := context.getArgumentsFromMemory(
		host,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
	if WithFaultAndHost(host, err, runtime.UseGasBoundedShouldFailExecution()) {
		return
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(actualLen))
	err = metering.UseGasBounded(gasToUse)
	if WithFaultAndHost(host, err, runtime.UseGasBoundedShouldFailExecution()) {
		return
	}

	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	calledSCAddress, err := context.MemLoad(destOffset, vmhost.AddressLen)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	UpgradeFromSourceContractWithTypedArgs(
		host,
		sourceContractAddress,
		calledSCAddress,
		value,
		data,
		gasLimit,
		codeMetadata,
	)
}

// UpgradeFromSourceContractWithTypedArgs - upgradeFromSourceContract with args already read from memory
func UpgradeFromSourceContractWithTypedArgs(
	host vmhost.VMHost,
	sourceContractAddress []byte,
	destContractAddress []byte,
	value []byte,
	data [][]byte,
	gasLimit int64,
	codeMetadata []byte,
) {
	runtime := host.Runtime()
	blockchain := host.Blockchain()

	code, err := blockchain.GetCode(sourceContractAddress)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	upgradeContract(host, destContractAddress, code, codeMetadata, value, data, gasLimit)
}

func upgradeContract(
	host vmhost.VMHost,
	destContractAddress []byte,
	code []byte,
	codeMetadata []byte,
	value []byte,
	data [][]byte,
	gasLimit int64,
) {
	runtime := host.Runtime()
	metering := host.Metering()
	gasSchedule := metering.GasSchedule()
	minAsyncCallCost := math.AddUint64(
		math.MulUint64(2, gasSchedule.BaseOpsAPICost.AsyncCallStep),
		gasSchedule.BaseOpsAPICost.AsyncCallbackGasLock)
	if uint64(gasLimit) < minAsyncCallCost {
		runtime.SetRuntimeBreakpointValue(vmhost.BreakpointOutOfGas)
		return
	}

	// Set up the async call as if it is not known whether the called SC
	// is in the same shard with the caller or not. This will be later resolved
	// by runtime.ExecuteAsyncCall().
	callData := vmhost.UpgradeFunctionName + "@" + hex.EncodeToString(code) + "@" + hex.EncodeToString(codeMetadata)
	for _, arg := range data {
		callData += "@" + hex.EncodeToString(arg)
	}

	async := host.Async()
	err := async.RegisterLegacyAsyncCall(
		destContractAddress,
		[]byte(callData),
		value,
	)
	logEEI.Trace("upgradeContract", "error", err)

	if errors.Is(err, vmhost.ErrNotEnoughGas) {
		runtime.SetRuntimeBreakpointValue(vmhost.BreakpointOutOfGas)
		return
	}
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}
}

// DeleteContract VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) DeleteContract(
	destOffset executor.MemPtr,
	gasLimit int64,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) {
	host := context.GetVMHost()
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(deleteContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	err := metering.UseGasBounded(gasToUse)
	if WithFaultAndHost(host, err, runtime.UseGasBoundedShouldFailExecution()) {
		return
	}

	data, actualLen, err := context.getArgumentsFromMemory(
		host,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
	if WithFaultAndHost(host, err, runtime.UseGasBoundedShouldFailExecution()) {
		return
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(actualLen))
	err = metering.UseGasBounded(gasToUse)
	if WithFaultAndHost(host, err, runtime.UseGasBoundedShouldFailExecution()) {
		return
	}

	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	calledSCAddress, err := context.MemLoad(destOffset, vmhost.AddressLen)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	deleteContract(
		host,
		calledSCAddress,
		data,
		gasLimit,
	)
}

func deleteContract(
	host vmhost.VMHost,
	dest []byte,
	data [][]byte,
	gasLimit int64,
) {
	runtime := host.Runtime()
	metering := host.Metering()
	gasSchedule := metering.GasSchedule()
	minAsyncCallCost := math.AddUint64(
		math.MulUint64(2, gasSchedule.BaseOpsAPICost.AsyncCallStep),
		gasSchedule.BaseOpsAPICost.AsyncCallbackGasLock)
	if uint64(gasLimit) < minAsyncCallCost {
		runtime.SetRuntimeBreakpointValue(vmhost.BreakpointOutOfGas)
		return
	}

	callData := vmhost.DeleteFunctionName
	for _, arg := range data {
		callData += "@" + hex.EncodeToString(arg)
	}

	async := host.Async()
	err := async.RegisterLegacyAsyncCall(
		dest,
		[]byte(callData),
		big.NewInt(0).Bytes(),
	)
	logEEI.Trace("deleteContract", "error", err)

	if errors.Is(err, vmhost.ErrNotEnoughGas) {
		runtime.SetRuntimeBreakpointValue(vmhost.BreakpointOutOfGas)
		return
	}
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}
}

// AsyncCall VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) AsyncCall(
	destOffset executor.MemPtr,
	valueOffset executor.MemPtr,
	dataOffset executor.MemPtr,
	length executor.MemLength) {

	host := context.GetVMHost()
	runtime := host.Runtime()
	async := host.Async()
	metering := host.Metering()
	metering.StartGasTracing(asyncCallName)

	gasSchedule := metering.GasSchedule()
	gasToUse := gasSchedule.BaseOpsAPICost.AsyncCallStep
	err := metering.UseGasBounded(gasToUse)
	if WithFaultAndHost(host, err, runtime.UseGasBoundedShouldFailExecution()) {
		return
	}

	calledSCAddress, err := context.MemLoad(destOffset, vmhost.AddressLen)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	value, err := context.MemLoad(valueOffset, vmhost.BalanceLen)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	gasToUse = math.MulUint64(gasSchedule.BaseOperationCost.DataCopyPerByte, uint64(length))
	err = metering.UseGasBounded(gasToUse)
	if WithFaultAndHost(host, err, runtime.UseGasBoundedShouldFailExecution()) {
		return
	}

	data, err := context.MemLoad(dataOffset, length)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	err = async.RegisterLegacyAsyncCall(calledSCAddress, data, value)
	if errors.Is(err, vmhost.ErrNotEnoughGas) {
		runtime.SetRuntimeBreakpointValue(vmhost.BreakpointOutOfGas)
		return
	}
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}
}

// GetArgumentLength VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetArgumentLength(id int32) int32 {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetArgument
	err := metering.UseGasBoundedAndAddTracedGas(getArgumentLengthName, gasToUse)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	args := runtime.Arguments()
	if id < 0 || int32(len(args)) <= id {
		context.WithFault(vmhost.ErrInvalidArgument, runtime.BaseOpsErrorShouldFailExecution())
		return -1
	}

	return int32(len(args[id]))
}

// GetArgument VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetArgument(id int32, argOffset executor.MemPtr) int32 {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetArgument
	err := metering.UseGasBoundedAndAddTracedGas(getArgumentName, gasToUse)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	args := runtime.Arguments()
	if id < 0 || int32(len(args)) <= id {
		context.WithFault(vmhost.ErrInvalidArgument, runtime.BaseOpsErrorShouldFailExecution())
		return -1
	}

	err = context.MemStore(argOffset, args[id])
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(args[id]))
}

// GetFunction VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetFunction(functionOffset executor.MemPtr) int32 {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetFunction
	err := metering.UseGasBoundedAndAddTracedGas(getFunctionName, gasToUse)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	function := runtime.FunctionName()
	err = context.MemStore(functionOffset, []byte(function))
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(function))
}

// GetNumArguments VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetNumArguments() int32 {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetNumArguments
	err := metering.UseGasBoundedAndAddTracedGas(getNumArgumentsName, gasToUse)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	args := runtime.Arguments()
	return int32(len(args))
}

// StorageStore VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) StorageStore(
	keyOffset executor.MemPtr,
	keyLength executor.MemLength,
	dataOffset executor.MemPtr,
	dataLength executor.MemLength) int32 {

	host := context.GetVMHost()
	return context.StorageStoreWithHost(
		host,
		keyOffset,
		keyLength,
		dataOffset,
		dataLength,
	)
}

// StorageStoreWithHost - storageStore with host instead of pointer context
func (context *VMHooksImpl) StorageStoreWithHost(
	host vmhost.VMHost,
	keyOffset executor.MemPtr,
	keyLength executor.MemLength,
	dataOffset executor.MemPtr,
	dataLength executor.MemLength) int32 {

	runtime := host.Runtime()

	key, err := context.MemLoad(keyOffset, keyLength)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	data, err := context.MemLoad(dataOffset, dataLength)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return StorageStoreWithTypedArgs(host, key, data)
}

// StorageStoreWithTypedArgs - storageStore with args already read from memory
func StorageStoreWithTypedArgs(host vmhost.VMHost, key []byte, data []byte) int32 {
	runtime := host.Runtime()
	storage := host.Storage()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.StorageStore
	err := metering.UseGasBoundedAndAddTracedGas(storageStoreName, gasToUse)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	storageStatus, err := storage.SetStorage(key, data)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int32(storageStatus)
}

// StorageLoadLength VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) StorageLoadLength(keyOffset executor.MemPtr, keyLength executor.MemLength) int32 {
	runtime := context.GetRuntimeContext()
	storage := context.GetStorageContext()
	metering := context.GetMeteringContext()

	key, err := context.MemLoad(keyOffset, keyLength)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	data, trieDepth, usedCache, err := storage.GetStorageUnmetered(key)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	err = storage.UseGasForStorageLoad(
		storageLoadLengthName,
		int64(trieDepth),
		metering.GasSchedule().BaseOpsAPICost.StorageLoad,
		usedCache)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(data))
}

// StorageLoadFromAddress VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) StorageLoadFromAddress(
	addressOffset executor.MemPtr,
	keyOffset executor.MemPtr,
	keyLength executor.MemLength,
	dataOffset executor.MemPtr) int32 {

	host := context.GetVMHost()
	return context.StorageLoadFromAddressWithHost(
		host,
		addressOffset,
		keyOffset,
		keyLength,
		dataOffset,
	)
}

// StorageLoadFromAddressWithHost - storageLoadFromAddress with host instead of pointer context
func (context *VMHooksImpl) StorageLoadFromAddressWithHost(
	host vmhost.VMHost,
	addressOffset executor.MemPtr,
	keyOffset executor.MemPtr,
	keyLength executor.MemLength,
	dataOffset executor.MemPtr) int32 {

	runtime := host.Runtime()

	key, err := context.MemLoad(keyOffset, keyLength)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	address, err := context.MemLoad(addressOffset, vmhost.AddressLen)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	data, err := StorageLoadFromAddressWithTypedArgs(host, address, key)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	err = context.MemStore(dataOffset, data)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(data))
}

// StorageLoadFromAddressWithTypedArgs - storageLoadFromAddress with args already read from memory
func StorageLoadFromAddressWithTypedArgs(host vmhost.VMHost, address []byte, key []byte) ([]byte, error) {
	storage := host.Storage()
	metering := host.Metering()
	data, trieDepth, usedCache, err := storage.GetStorageFromAddress(address, key)
	if err != nil {
		return nil, err
	}
	err = storage.UseGasForStorageLoad(
		storageLoadFromAddressName,
		int64(trieDepth),
		metering.GasSchedule().BaseOpsAPICost.StorageLoad,
		usedCache)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// StorageLoad VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) StorageLoad(keyOffset executor.MemPtr, keyLength executor.MemLength, dataOffset executor.MemPtr) int32 {
	host := context.GetVMHost()
	return context.StorageLoadWithHost(
		host,
		keyOffset,
		keyLength,
		dataOffset,
	)
}

// StorageLoadWithHost - storageLoad with host instead of pointer context
func (context *VMHooksImpl) StorageLoadWithHost(host vmhost.VMHost, keyOffset executor.MemPtr, keyLength executor.MemLength, dataOffset executor.MemPtr) int32 {
	runtime := host.Runtime()

	key, err := context.MemLoad(keyOffset, keyLength)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	data, err := StorageLoadWithWithTypedArgs(host, key)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	err = context.MemStore(dataOffset, data)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(data))
}

// StorageLoadWithWithTypedArgs - storageLoad with args already read from memory
func StorageLoadWithWithTypedArgs(host vmhost.VMHost, key []byte) ([]byte, error) {
	storage := host.Storage()
	metering := host.Metering()
	data, trieDepth, usedCache, err := storage.GetStorage(key)
	if err != nil {
		return nil, err
	}

	err = storage.UseGasForStorageLoad(
		storageLoadName,
		int64(trieDepth),
		metering.GasSchedule().BaseOpsAPICost.StorageLoad,
		usedCache)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// SetStorageLock VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) SetStorageLock(keyOffset executor.MemPtr, keyLength executor.MemLength, lockTimestamp int64) int32 {
	host := context.GetVMHost()
	return context.SetStorageLockWithHost(
		host,
		keyOffset,
		keyLength,
		lockTimestamp,
	)
}

// SetStorageLockWithHost - setStorageLock with host instead of pointer context
func (context *VMHooksImpl) SetStorageLockWithHost(host vmhost.VMHost, keyOffset executor.MemPtr, keyLength executor.MemLength, lockTimestamp int64) int32 {
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.Int64StorageStore
	err := metering.UseGasBoundedAndAddTracedGas(setStorageLockName, gasToUse)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	key, err := context.MemLoad(keyOffset, keyLength)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return SetStorageLockWithTypedArgs(host, key, lockTimestamp)
}

// SetStorageLockWithTypedArgs - setStorageLock with args already read from memory
func SetStorageLockWithTypedArgs(host vmhost.VMHost, key []byte, lockTimestamp int64) int32 {
	runtime := host.Runtime()
	storage := host.Storage()
	timeLockKeyPrefix := string(storage.GetVmProtectedPrefix(vmhost.TimeLockKeyPrefix))
	timeLockKey := vmhost.CustomStorageKey(timeLockKeyPrefix, key)
	bigTimestamp := big.NewInt(0).SetInt64(lockTimestamp)
	storageStatus, err := storage.SetProtectedStorage(timeLockKey, bigTimestamp.Bytes())
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}
	return int32(storageStatus)
}

// GetStorageLock VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetStorageLock(keyOffset executor.MemPtr, keyLength executor.MemLength) int64 {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	storage := context.GetStorageContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.StorageLoad
	err := metering.UseGasBoundedAndAddTracedGas(getStorageLockName, gasToUse)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	key, err := context.MemLoad(keyOffset, keyLength)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	timeLockKeyPrefix := string(storage.GetVmProtectedPrefix(vmhost.TimeLockKeyPrefix))
	timeLockKey := vmhost.CustomStorageKey(timeLockKeyPrefix, key)

	data, trieDepth, usedCache, err := storage.GetStorage(timeLockKey)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	err = storage.UseGasForStorageLoad(
		getStorageLockName,
		int64(trieDepth),
		metering.GasSchedule().BaseOpsAPICost.StorageLoad,
		usedCache)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	timeLock := big.NewInt(0).SetBytes(data).Int64()

	// TODO if timelock <= currentTimeStamp { fail somehow }

	return timeLock
}

// IsStorageLocked VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) IsStorageLocked(keyOffset executor.MemPtr, keyLength executor.MemLength) int32 {
	timeLock := context.GetStorageLock(keyOffset, keyLength)
	if timeLock < 0 {
		return -1
	}

	currentTimestamp := context.GetBlockTimestamp()
	if timeLock <= currentTimestamp {
		return 0
	}

	return 1
}

// ClearStorageLock VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ClearStorageLock(keyOffset executor.MemPtr, keyLength executor.MemLength) int32 {
	return context.SetStorageLock(keyOffset, keyLength, 0)
}

// GetCaller VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetCaller(resultOffset executor.MemPtr) {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCaller
	err := metering.UseGasBoundedAndAddTracedGas(getCallerName, gasToUse)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	caller := runtime.GetVMInput().CallerAddr

	err = context.MemStore(resultOffset, caller)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}
}

// CheckNoPayment VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) CheckNoPayment() {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	err := metering.UseGasBoundedAndAddTracedGas(checkNoPaymentName, gasToUse)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	vmInput := runtime.GetVMInput()
	if vmInput.CallValue.Sign() > 0 {
		_ = context.WithFault(vmhost.ErrNonPayableFunctionRewa, runtime.BaseOpsErrorShouldFailExecution())
		return
	}
	if len(vmInput.DCDTTransfers) > 0 {
		_ = context.WithFault(vmhost.ErrNonPayableFunctionDcdt, runtime.BaseOpsErrorShouldFailExecution())
		return
	}
}

// GetCallValue VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetCallValue(resultOffset executor.MemPtr) int32 {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	err := metering.UseGasBoundedAndAddTracedGas(callValueName, gasToUse)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	value := runtime.GetVMInput().CallValue.Bytes()
	value = vmhost.PadBytesLeft(value, vmhost.BalanceLen)

	err = context.MemStore(resultOffset, value)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(value))
}

// GetDCDTValue VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetDCDTValue(resultOffset executor.MemPtr) int32 {
	isFail := failIfMoreThanOneDCDTTransfer(context)
	if isFail {
		return -1
	}
	return context.GetDCDTValueByIndex(resultOffset, 0)
}

// GetDCDTValueByIndex VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetDCDTValueByIndex(resultOffset executor.MemPtr, index int32) int32 {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	err := metering.UseGasBoundedAndAddTracedGas(getDCDTValueByIndexName, gasToUse)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	var value []byte

	dcdtTransfer := getDCDTTransferFromInputFailIfWrongIndex(context.GetVMHost(), index)
	if dcdtTransfer != nil && dcdtTransfer.DCDTValue.Cmp(vmhost.Zero) > 0 {
		value = dcdtTransfer.DCDTValue.Bytes()
		value = vmhost.PadBytesLeft(value, vmhost.BalanceLen)
	}

	err = context.MemStore(resultOffset, value)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(value))
}

// GetDCDTTokenName VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetDCDTTokenName(resultOffset executor.MemPtr) int32 {
	isFail := failIfMoreThanOneDCDTTransfer(context)
	if isFail {
		return -1
	}
	return context.GetDCDTTokenNameByIndex(resultOffset, 0)
}

// GetDCDTTokenNameByIndex VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetDCDTTokenNameByIndex(resultOffset executor.MemPtr, index int32) int32 {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	err := metering.UseGasBoundedAndAddTracedGas(getDCDTTokenNameByIndexName, gasToUse)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	dcdtTransfer := getDCDTTransferFromInputFailIfWrongIndex(context.GetVMHost(), index)
	var tokenName []byte
	if dcdtTransfer != nil {
		tokenName = dcdtTransfer.DCDTTokenName
	}

	err = context.MemStore(resultOffset, tokenName)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(tokenName))
}

// GetDCDTTokenNonce VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetDCDTTokenNonce() int64 {
	isFail := failIfMoreThanOneDCDTTransfer(context)
	if isFail {
		return -1
	}
	return context.GetDCDTTokenNonceByIndex(0)
}

// GetDCDTTokenNonceByIndex VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetDCDTTokenNonceByIndex(index int32) int64 {
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	err := metering.UseGasBoundedAndAddTracedGas(getDCDTTokenNonceByIndexName, gasToUse)
	if context.WithFault(err, context.GetRuntimeContext().BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	dcdtTransfer := getDCDTTransferFromInputFailIfWrongIndex(context.GetVMHost(), index)
	nonce := uint64(0)
	if dcdtTransfer != nil {
		nonce = dcdtTransfer.DCDTTokenNonce
	}
	return int64(nonce)
}

// GetCurrentDCDTNFTNonce VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetCurrentDCDTNFTNonce(
	addressOffset executor.MemPtr,
	tokenIDOffset executor.MemPtr,
	tokenIDLen executor.MemLength) int64 {

	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	storage := context.GetStorageContext()

	destination, err := context.MemLoad(addressOffset, vmhost.AddressLen)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 0
	}

	tokenID, err := context.MemLoad(tokenIDOffset, tokenIDLen)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 0
	}

	key := []byte(core.ProtectedKeyPrefix + core.DCDTNFTLatestNonceIdentifier + string(tokenID))
	data, trieDepth, _, err := storage.GetStorageFromAddress(destination, key)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 0
	}

	err = storage.UseGasForStorageLoad(
		getCurrentDCDTNFTNonceName,
		int64(trieDepth),
		metering.GasSchedule().BaseOpsAPICost.StorageLoad,
		false)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 0
	}

	nonce := big.NewInt(0).SetBytes(data).Uint64()
	return int64(nonce)
}

// GetDCDTTokenType VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetDCDTTokenType() int32 {
	isFail := failIfMoreThanOneDCDTTransfer(context)
	if isFail {
		return -1
	}
	return context.GetDCDTTokenTypeByIndex(0)
}

// GetDCDTTokenTypeByIndex VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetDCDTTokenTypeByIndex(index int32) int32 {
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	err := metering.UseGasBoundedAndAddTracedGas(getDCDTTokenTypeByIndexName, gasToUse)
	if context.WithFault(err, context.GetRuntimeContext().BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	dcdtTransfer := getDCDTTransferFromInputFailIfWrongIndex(context.GetVMHost(), index)
	if dcdtTransfer != nil {
		return int32(dcdtTransfer.DCDTTokenType)
	}
	return 0
}

// GetNumDCDTTransfers VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetNumDCDTTransfers() int32 {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	err := metering.UseGasBoundedAndAddTracedGas(getNumDCDTTransfersName, gasToUse)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(runtime.GetVMInput().DCDTTransfers))
}

// GetCallValueTokenName VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetCallValueTokenName(callValueOffset executor.MemPtr, tokenNameOffset executor.MemPtr) int32 {
	isFail := failIfMoreThanOneDCDTTransfer(context)
	if isFail {
		return -1
	}
	return context.GetCallValueTokenNameByIndex(callValueOffset, tokenNameOffset, 0)
}

// GetCallValueTokenNameByIndex VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetCallValueTokenNameByIndex(
	callValueOffset executor.MemPtr,
	tokenNameOffset executor.MemPtr,
	index int32) int32 {

	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	err := metering.UseGasBoundedAndAddTracedGas(getCallValueTokenNameByIndexName, gasToUse)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	callValue := runtime.GetVMInput().CallValue.Bytes()
	tokenName := make([]byte, 0)
	dcdtTransfer := getDCDTTransferFromInputFailIfWrongIndex(context.GetVMHost(), index)

	if dcdtTransfer != nil {
		tokenName = make([]byte, len(dcdtTransfer.DCDTTokenName))
		copy(tokenName, dcdtTransfer.DCDTTokenName)
		callValue = dcdtTransfer.DCDTValue.Bytes()
	}
	callValue = vmhost.PadBytesLeft(callValue, vmhost.BalanceLen)

	err = context.MemStore(tokenNameOffset, tokenName)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	err = context.MemStore(callValueOffset, callValue)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(tokenName))
}

// IsReservedFunctionName VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) IsReservedFunctionName(nameHandle int32) int32 {
	host := context.host
	managedTypes := context.GetManagedTypesContext()
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.IsReservedFunctionName
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	name, err := managedTypes.GetBytes(nameHandle)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	if runtime.IsReservedFunctionName(string(name)) {
		return 1
	}

	return 0
}

// WriteLog VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) WriteLog(
	dataPointer executor.MemPtr,
	dataLength executor.MemLength,
	topicPtr executor.MemPtr,
	numTopics int32) {

	// note: deprecated
	runtime := context.GetRuntimeContext()
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.Log
	gas := math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(numTopics*vmhost.HashLen+dataLength))
	gasToUse = math.AddUint64(gasToUse, gas)

	if numTopics < 0 || dataLength < 0 {
		err := vmhost.ErrNegativeLength
		context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution())
		return
	}

	err := metering.UseGasBoundedAndAddTracedGas(writeLogName, gasToUse)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	log, err := context.MemLoad(dataPointer, dataLength)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	topics := make([][]byte, numTopics)
	for i := int32(0); i < numTopics; i++ {
		topics[i], err = context.MemLoad(topicPtr.Offset(i*vmhost.HashLen), vmhost.HashLen)
		if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
			return
		}
	}

	output.WriteLog(runtime.GetContextAddress(), topics, [][]byte{log})
}

// WriteEventLog VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) WriteEventLog(
	numTopics int32,
	topicLengthsOffset executor.MemPtr,
	topicOffset executor.MemPtr,
	dataOffset executor.MemPtr,
	dataLength executor.MemLength,
) {

	host := context.GetVMHost()
	runtime := context.GetRuntimeContext()
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()

	topics, topicDataTotalLen, err := context.getArgumentsFromMemory(
		host,
		numTopics,
		topicLengthsOffset,
		topicOffset,
	)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	data, err := context.MemLoad(dataOffset, dataLength)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	gasToUse := metering.GasSchedule().BaseOpsAPICost.Log
	gasForData := math.MulUint64(
		metering.GasSchedule().BaseOperationCost.DataCopyPerByte,
		uint64(topicDataTotalLen+dataLength))
	gasToUse = math.AddUint64(gasToUse, gasForData)
	err = metering.UseGasBoundedAndAddTracedGas(writeEventLogName, gasToUse)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	output.WriteLog(runtime.GetContextAddress(), topics, [][]byte{data})
}

// GetBlockTimestamp VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetBlockTimestamp() int64 {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockTimeStamp
	err := metering.UseGasBoundedAndAddTracedGas(getBlockTimestampName, gasToUse)
	if context.WithFault(err, context.GetRuntimeContext().BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int64(blockchain.CurrentTimeStamp())
}

// GetBlockNonce VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetBlockNonce() int64 {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockNonce
	err := metering.UseGasBoundedAndAddTracedGas(getBlockNonceName, gasToUse)
	if context.WithFault(err, context.GetRuntimeContext().BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int64(blockchain.CurrentNonce())
}

// GetBlockRound VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetBlockRound() int64 {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockRound
	err := metering.UseGasBoundedAndAddTracedGas(getBlockRoundName, gasToUse)
	if context.WithFault(err, context.GetRuntimeContext().BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int64(blockchain.CurrentRound())
}

// GetBlockEpoch VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetBlockEpoch() int64 {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockEpoch
	err := metering.UseGasBoundedAndAddTracedGas(getBlockEpochName, gasToUse)
	if context.WithFault(err, context.GetRuntimeContext().BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int64(blockchain.CurrentEpoch())
}

// GetBlockRandomSeed VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetBlockRandomSeed(pointer executor.MemPtr) {
	runtime := context.GetRuntimeContext()
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockRandomSeed
	err := metering.UseGasBoundedAndAddTracedGas(getBlockRandomSeedName, gasToUse)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	randomSeed := blockchain.CurrentRandomSeed()
	err = context.MemStore(pointer, randomSeed)
	context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution())
}

// GetStateRootHash VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetStateRootHash(pointer executor.MemPtr) {
	runtime := context.GetRuntimeContext()
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetStateRootHash
	err := metering.UseGasBoundedAndAddTracedGas(getStateRootHashName, gasToUse)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	stateRootHash := blockchain.GetStateRootHash()
	err = context.MemStore(pointer, stateRootHash)
	context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution())
}

// GetPrevBlockTimestamp VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetPrevBlockTimestamp() int64 {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockTimeStamp
	err := metering.UseGasBoundedAndAddTracedGas(getPrevBlockTimestampName, gasToUse)
	if context.WithFault(err, context.GetRuntimeContext().BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int64(blockchain.LastTimeStamp())
}

// GetPrevBlockNonce VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetPrevBlockNonce() int64 {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockNonce
	err := metering.UseGasBoundedAndAddTracedGas(getPrevBlockNonceName, gasToUse)
	if context.WithFault(err, context.GetRuntimeContext().BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int64(blockchain.LastNonce())
}

// GetPrevBlockRound VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetPrevBlockRound() int64 {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockRound
	err := metering.UseGasBoundedAndAddTracedGas(getPrevBlockRoundName, gasToUse)
	if context.WithFault(err, context.GetRuntimeContext().BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int64(blockchain.LastRound())
}

// GetPrevBlockEpoch VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetPrevBlockEpoch() int64 {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockEpoch
	err := metering.UseGasBoundedAndAddTracedGas(getPrevBlockEpochName, gasToUse)
	if context.WithFault(err, context.GetRuntimeContext().BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int64(blockchain.LastEpoch())
}

// GetPrevBlockRandomSeed VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetPrevBlockRandomSeed(pointer executor.MemPtr) {
	runtime := context.GetRuntimeContext()
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockRandomSeed
	err := metering.UseGasBoundedAndAddTracedGas(getPrevBlockRandomSeedName, gasToUse)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	randomSeed := blockchain.LastRandomSeed()
	err = context.MemStore(pointer, randomSeed)
	context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution())
}

// Finish VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) Finish(pointer executor.MemPtr, length executor.MemLength) {
	runtime := context.GetRuntimeContext()
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(returnDataName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.Finish
	gas := math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(length))
	gasToUse = math.AddUint64(gasToUse, gas)
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		_ = context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution())
		return
	}

	data, err := context.MemLoad(pointer, length)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	output.Finish(data)
}

// ExecuteOnSameContext VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ExecuteOnSameContext(
	gasLimit int64,
	addressOffset executor.MemPtr,
	valueOffset executor.MemPtr,
	functionOffset executor.MemPtr,
	functionLength executor.MemLength,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) int32 {
	host := context.GetVMHost()
	metering := host.Metering()
	metering.StartGasTracing(executeOnSameContextName)

	return context.ExecuteOnSameContextWithHost(
		host,
		gasLimit,
		addressOffset,
		valueOffset,
		functionOffset,
		functionLength,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
}

// ExecuteOnSameContextWithHost - executeOnSameContext with host instead of pointer context
func (context *VMHooksImpl) ExecuteOnSameContextWithHost(
	host vmhost.VMHost,
	gasLimit int64,
	addressOffset executor.MemPtr,
	valueOffset executor.MemPtr,
	functionOffset executor.MemPtr,
	functionLength executor.MemLength,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) int32 {
	runtime := host.Runtime()

	callArgs, err := context.extractIndirectContractCallArgumentsWithValue(
		host, addressOffset, valueOffset, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	return ExecuteOnSameContextWithTypedArgs(
		host,
		gasLimit,
		callArgs.value,
		callArgs.function,
		callArgs.dest,
		callArgs.args,
	)
}

// ExecuteOnSameContextWithTypedArgs - executeOnSameContext with args already read from memory
func ExecuteOnSameContextWithTypedArgs(
	host vmhost.VMHost,
	gasLimit int64,
	value *big.Int,
	function []byte,
	dest []byte,
	args [][]byte,
) int32 {
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.ExecuteOnSameContext
	err := metering.UseGasBounded(gasToUse)
	if WithFaultAndHost(host, err, runtime.UseGasBoundedShouldFailExecution()) {
		return -1
	}

	sender := runtime.GetContextAddress()

	contractCallInput, err := prepareIndirectContractCallInput(
		host,
		sender,
		value,
		gasLimit,
		dest,
		function,
		args,
		gasToUse,
		true,
	)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	if host.IsBuiltinFunctionName(contractCallInput.Function) {
		WithFaultAndHost(host, vmhost.ErrInvalidBuiltInFunctionCall, runtime.BaseOpsErrorShouldFailExecution())
		return 1
	}

	err = host.ExecuteOnSameContext(contractCallInput)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return 0
}

// ExecuteOnDestContext VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ExecuteOnDestContext(
	gasLimit int64,
	addressOffset executor.MemPtr,
	valueOffset executor.MemPtr,
	functionOffset executor.MemPtr,
	functionLength executor.MemLength,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) int32 {
	host := context.GetVMHost()
	metering := host.Metering()
	metering.StartGasTracing(executeOnDestContextName)

	return context.ExecuteOnDestContextWithHost(
		host,
		gasLimit,
		addressOffset,
		valueOffset,
		functionOffset,
		functionLength,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
}

// ExecuteOnDestContextWithHost - executeOnDestContext with host instead of pointer context
func (context *VMHooksImpl) ExecuteOnDestContextWithHost(
	host vmhost.VMHost,
	gasLimit int64,
	addressOffset executor.MemPtr,
	valueOffset executor.MemPtr,
	functionOffset executor.MemPtr,
	functionLength executor.MemLength,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) int32 {
	runtime := host.Runtime()

	callArgs, err := context.extractIndirectContractCallArgumentsWithValue(
		host, addressOffset, valueOffset, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	return ExecuteOnDestContextWithTypedArgs(
		host,
		gasLimit,
		callArgs.value,
		callArgs.function,
		callArgs.dest,
		callArgs.args,
	)
}

// ExecuteOnDestContextWithTypedArgs - executeOnDestContext with args already read from memory
func ExecuteOnDestContextWithTypedArgs(
	host vmhost.VMHost,
	gasLimit int64,
	value *big.Int,
	function []byte,
	dest []byte,
	args [][]byte,
) int32 {
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.ExecuteOnDestContext
	err := metering.UseGasBounded(gasToUse)
	if WithFaultAndHost(host, err, runtime.UseGasBoundedShouldFailExecution()) {
		return -1
	}

	sender := runtime.GetContextAddress()

	contractCallInput, err := prepareIndirectContractCallInput(
		host,
		sender,
		value,
		gasLimit,
		dest,
		function,
		args,
		gasToUse,
		true,
	)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	vmOutput, err := executeOnDestContextFromAPI(host, contractCallInput)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	host.CompleteLogEntriesWithCallType(vmOutput, vmhost.ExecuteOnDestContextString)

	return 0
}

// ExecuteReadOnly VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ExecuteReadOnly(
	gasLimit int64,
	addressOffset executor.MemPtr,
	functionOffset executor.MemPtr,
	functionLength executor.MemLength,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) int32 {
	host := context.GetVMHost()
	metering := host.Metering()
	metering.StartGasTracing(executeReadOnlyName)

	return context.ExecuteReadOnlyWithHost(
		host,
		gasLimit,
		addressOffset,
		functionOffset,
		functionLength,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
}

// ExecuteReadOnlyWithHost - executeReadOnly with host instead of pointer context
func (context *VMHooksImpl) ExecuteReadOnlyWithHost(
	host vmhost.VMHost,
	gasLimit int64,
	addressOffset executor.MemPtr,
	functionOffset executor.MemPtr,
	functionLength executor.MemLength,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) int32 {
	runtime := host.Runtime()

	callArgs, err := context.extractIndirectContractCallArgumentsWithoutValue(
		host, addressOffset, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return ExecuteReadOnlyWithTypedArguments(
		host,
		gasLimit,
		callArgs.function,
		callArgs.dest,
		callArgs.args,
	)
}

// ExecuteReadOnlyWithTypedArguments - executeReadOnly with args already read from memory
func ExecuteReadOnlyWithTypedArguments(
	host vmhost.VMHost,
	gasLimit int64,
	function []byte,
	dest []byte,
	args [][]byte,
) int32 {
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.ExecuteReadOnly
	err := metering.UseGasBounded(gasToUse)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	sender := runtime.GetContextAddress()

	contractCallInput, err := prepareIndirectContractCallInput(
		host,
		sender,
		big.NewInt(0),
		gasLimit,
		dest,
		function,
		args,
		gasToUse,
		true,
	)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	if host.IsBuiltinFunctionName(contractCallInput.Function) {
		WithFaultAndHost(host, vmhost.ErrInvalidBuiltInFunctionCall, runtime.BaseOpsErrorShouldFailExecution())
		return 1
	}

	wasReadOnly := runtime.ReadOnly()
	runtime.SetReadOnly(true)
	_, err = executeOnDestContextFromAPI(host, contractCallInput)
	runtime.SetReadOnly(wasReadOnly)

	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return 0
}

// CreateContract VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) CreateContract(
	gasLimit int64,
	valueOffset executor.MemPtr,
	codeOffset executor.MemPtr,
	codeMetadataOffset executor.MemPtr,
	length executor.MemLength,
	resultOffset executor.MemPtr,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) int32 {
	host := context.GetVMHost()
	return context.createContractWithHost(
		host,
		gasLimit,
		valueOffset,
		codeOffset,
		codeMetadataOffset,
		length,
		resultOffset,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
}

func (context *VMHooksImpl) createContractWithHost(
	host vmhost.VMHost,
	gasLimit int64,
	valueOffset executor.MemPtr,
	codeOffset executor.MemPtr,
	codeMetadataOffset executor.MemPtr,
	length executor.MemLength,
	resultOffset executor.MemPtr,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) int32 {
	runtime := host.Runtime()

	metering := host.Metering()
	metering.StartGasTracing(createContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	err := metering.UseGasBounded(gasToUse)
	if WithFaultAndHost(host, err, runtime.UseGasBoundedShouldFailExecution()) {
		return -1
	}

	sender := runtime.GetContextAddress()
	value, err := context.MemLoad(valueOffset, vmhost.BalanceLen)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	code, err := context.MemLoad(codeOffset, length)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	codeMetadata, err := context.MemLoad(codeMetadataOffset, vmhost.CodeMetadataLen)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	data, actualLen, err := context.getArgumentsFromMemory(
		host,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(actualLen))
	err = metering.UseGasBounded(gasToUse)
	if WithFaultAndHost(host, err, runtime.UseGasBoundedShouldFailExecution()) {
		return 1
	}

	valueAsInt := big.NewInt(0).SetBytes(value)
	newAddress, err := createContract(sender, data, valueAsInt, gasLimit, code, codeMetadata, host, CreateContract)

	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	err = context.MemStore(resultOffset, newAddress)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

// DeployFromSourceContract VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) DeployFromSourceContract(
	gasLimit int64,
	valueOffset executor.MemPtr,
	sourceContractAddressOffset executor.MemPtr,
	codeMetadataOffset executor.MemPtr,
	resultAddressOffset executor.MemPtr,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) int32 {
	host := context.GetVMHost()
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(deployFromSourceContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	err := metering.UseGasBounded(gasToUse)
	if WithFaultAndHost(host, err, runtime.UseGasBoundedShouldFailExecution()) {
		return 1
	}

	value, err := context.MemLoad(valueOffset, vmhost.BalanceLen)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	sourceContractAddress, err := context.MemLoad(sourceContractAddressOffset, vmhost.AddressLen)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	codeMetadata, err := context.MemLoad(codeMetadataOffset, vmhost.CodeMetadataLen)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	data, actualLen, err := context.getArgumentsFromMemory(
		host,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(actualLen))
	err = metering.UseGasBounded(gasToUse)
	if WithFaultAndHost(host, err, runtime.UseGasBoundedShouldFailExecution()) {
		return 1
	}

	newAddress, err := DeployFromSourceContractWithTypedArgs(
		host,
		sourceContractAddress,
		codeMetadata,
		big.NewInt(0).SetBytes(value),
		data,
		gasLimit,
	)

	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	err = context.MemStore(resultAddressOffset, newAddress)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

// DeployFromSourceContractWithTypedArgs - deployFromSourceContract with args already read from memory
func DeployFromSourceContractWithTypedArgs(
	host vmhost.VMHost,
	sourceContractAddress []byte,
	codeMetadata []byte,
	value *big.Int,
	data [][]byte,
	gasLimit int64,
) ([]byte, error) {
	runtime := host.Runtime()
	sender := runtime.GetContextAddress()

	blockchain := host.Blockchain()
	code, err := blockchain.GetCode(sourceContractAddress)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return nil, err
	}

	return createContract(sender, data, value, gasLimit, code, codeMetadata, host, DeployContract)
}

func createContract(
	sender []byte,
	data [][]byte,
	value *big.Int,
	gasLimit int64,
	code []byte,
	codeMetadata []byte,
	host vmhost.VMHost,
	createContractCallType CreateContractCallType,
) ([]byte, error) {
	originalCaller := host.Runtime().GetOriginalCallerAddress()
	metering := host.Metering()
	contractCreate := &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			OriginalCallerAddr: originalCaller,
			CallerAddr:         sender,
			Arguments:          data,
			CallValue:          value,
			GasPrice:           0,
			GasProvided:        metering.BoundGasLimit(gasLimit),
		},
		ContractCode:         code,
		ContractCodeMetadata: codeMetadata,
	}

	currentVMInput := host.Runtime().GetVMInput()
	if len(currentVMInput.RelayerAddr) > 0 {
		contractCreate.RelayerAddr = make([]byte, len(currentVMInput.RelayerAddr))
		copy(contractCreate.RelayerAddr, currentVMInput.RelayerAddr)
	}

	return host.CreateNewContract(contractCreate, int(createContractCallType))
}

// GetNumReturnData VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetNumReturnData() int32 {
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetNumReturnData
	err := metering.UseGasBoundedAndAddTracedGas(getNumReturnDataName, gasToUse)
	if context.WithFault(err, context.GetRuntimeContext().BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	returnData := output.ReturnData()
	return int32(len(returnData))
}

// GetReturnDataSize VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetReturnDataSize(resultID int32) int32 {
	runtime := context.GetRuntimeContext()
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetReturnDataSize
	err := metering.UseGasBoundedAndAddTracedGas(getReturnDataSizeName, gasToUse)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	returnData := output.ReturnData()
	if resultID >= int32(len(returnData)) || resultID < 0 {
		context.WithFault(vmhost.ErrInvalidArgument, runtime.BaseOpsErrorShouldFailExecution())
		return 0
	}

	return int32(len(returnData[resultID]))
}

// GetReturnData VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetReturnData(resultID int32, dataOffset executor.MemPtr) int32 {
	host := context.GetVMHost()

	result := GetReturnDataWithHostAndTypedArgs(host, resultID)
	if result == nil {
		return 0
	}

	runtime := context.GetRuntimeContext()
	err := context.MemStore(dataOffset, result)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 0
	}

	return int32(len(result))
}

func GetReturnDataWithHostAndTypedArgs(host vmhost.VMHost, resultID int32) []byte {
	output := host.Output()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetReturnData
	err := metering.UseGasBoundedAndAddTracedGas(getReturnDataName, gasToUse)
	if WithFaultAndHost(host, err, host.Runtime().BaseOpsErrorShouldFailExecution()) {
		return nil
	}

	returnData := output.ReturnData()
	if resultID >= int32(len(returnData)) || resultID < 0 {
		WithFaultAndHost(host, vmhost.ErrInvalidArgument, host.Runtime().BaseOpsErrorShouldFailExecution())
		return nil
	}

	return returnData[resultID]
}

// CleanReturnData VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) CleanReturnData() {
	host := context.GetVMHost()
	CleanReturnDataWithHost(host)
}

// CleanReturnDataWithHost - exposed version of v1_5_deleteFromReturnData for tests
func CleanReturnDataWithHost(host vmhost.VMHost) {
	output := host.Output()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CleanReturnData
	err := metering.UseGasBoundedAndAddTracedGas(cleanReturnDataName, gasToUse)
	if WithFaultAndHost(host, err, host.Runtime().BaseOpsErrorShouldFailExecution()) {
		return
	}

	output.ClearReturnData()
}

// DeleteFromReturnData VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) DeleteFromReturnData(resultID int32) {
	host := context.GetVMHost()
	DeleteFromReturnDataWithHost(host, resultID)
}

// DeleteFromReturnDataWithHost - exposed version of v1_5_deleteFromReturnData for tests
func DeleteFromReturnDataWithHost(host vmhost.VMHost, resultID int32) {
	output := host.Output()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.DeleteFromReturnData
	err := metering.UseGasBoundedAndAddTracedGas(deleteFromReturnDataName, gasToUse)
	if WithFaultAndHost(host, err, host.Runtime().BaseOpsErrorShouldFailExecution()) {
		return
	}

	returnData := output.ReturnData()
	if resultID < int32(len(returnData)) {
		output.RemoveReturnData(uint32(resultID))
	}
}

// GetOriginalTxHash VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetOriginalTxHash(dataOffset executor.MemPtr) {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetOriginalTxHash
	err := metering.UseGasBoundedAndAddTracedGas(getOriginalTxHashName, gasToUse)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	err = context.MemStore(dataOffset, runtime.GetOriginalTxHash())
	_ = context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution())
}

// GetCurrentTxHash VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetCurrentTxHash(dataOffset executor.MemPtr) {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCurrentTxHash
	err := metering.UseGasBoundedAndAddTracedGas(getCurrentTxHashName, gasToUse)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	err = context.MemStore(dataOffset, runtime.GetCurrentTxHash())
	_ = context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution())
}

// GetPrevTxHash VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetPrevTxHash(dataOffset executor.MemPtr) {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetPrevTxHash
	err := metering.UseGasBoundedAndAddTracedGas(getPrevTxHashName, gasToUse)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	err = context.MemStore(dataOffset, runtime.GetPrevTxHash())
	_ = context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution())
}

func prepareIndirectContractCallInput(
	host vmhost.VMHost,
	sender []byte,
	value *big.Int,
	gasLimit int64,
	destination []byte,
	function []byte,
	data [][]byte,
	_ uint64,
	syncExecutionRequired bool,
) (*vmcommon.ContractCallInput, error) {
	runtime := host.Runtime()
	metering := host.Metering()

	if syncExecutionRequired && !host.AreInSameShard(runtime.GetContextAddress(), destination) {
		return nil, vmhost.ErrSyncExecutionNotInSameShard
	}

	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			OriginalCallerAddr: host.Runtime().GetOriginalCallerAddress(),
			CallerAddr:         sender,
			Arguments:          data,
			CallValue:          value,
			GasPrice:           0,
			GasProvided:        metering.BoundGasLimit(gasLimit),
			CallType:           vm.DirectCall,
		},
		RecipientAddr: destination,
		Function:      string(function),
	}

	currentVMInput := runtime.GetVMInput()
	if len(currentVMInput.RelayerAddr) > 0 {
		contractCallInput.RelayerAddr = make([]byte, len(currentVMInput.RelayerAddr))
		copy(contractCallInput.RelayerAddr, currentVMInput.RelayerAddr)
	}

	return contractCallInput, nil
}

func (context *VMHooksImpl) getArgumentsFromMemory(
	_ vmhost.VMHost,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) ([][]byte, int32, error) {
	if numArguments < 0 {
		return nil, 0, fmt.Errorf("negative numArguments (%d)", numArguments)
	}

	argumentsLengthData, err := context.MemLoad(argumentsLengthOffset, numArguments*4)
	if err != nil {
		return nil, 0, err
	}

	argumentLengths := createInt32Array(argumentsLengthData, numArguments)
	data, err := context.MemLoadMultiple(dataOffset, argumentLengths)
	if err != nil {
		return nil, 0, err
	}

	totalArgumentBytes := int32(0)
	for _, length := range argumentLengths {
		totalArgumentBytes += length
	}

	return data, totalArgumentBytes, nil
}

func createInt32Array(rawData []byte, numIntegers int32) []int32 {
	integers := make([]int32, numIntegers)
	index := 0
	for cursor := 0; cursor < len(rawData); cursor += 4 {
		rawInt := rawData[cursor : cursor+4]
		actualInt := binary.LittleEndian.Uint32(rawInt)
		integers[index] = int32(actualInt)
		index++
	}
	return integers
}

func executeOnDestContextFromAPI(host vmhost.VMHost, input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	host.Async().SetAsyncArgumentsForCall(input)
	vmOutput, isChildComplete, err := host.ExecuteOnDestContext(input)
	if err != nil {
		return nil, err
	}

	err = host.Async().CompleteChildConditional(isChildComplete, nil, 0)
	if err != nil {
		return nil, err
	}
	return vmOutput, err
}
