package vmhooks

import (
	"encoding/binary"
	"errors"
	"math/big"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	"github.com/TerraDharitri/drt-go-chain-vm/math"
	"github.com/TerraDharitri/drt-go-chain-vm/vmhost"
)

const dcdtTransferLen = 16

// Deserializes a vmcommon.DCDTTransfer object.
func readDCDTTransfer(
	managedType vmhost.ManagedTypesContext,
	runtime vmhost.RuntimeContext,
	data []byte,
) (*vmcommon.DCDTTransfer, error) {
	if len(data) != dcdtTransferLen {
		return nil, errors.New("invalid DCDT transfer object encoding")
	}

	tokenIdentifierHandle := int32(binary.BigEndian.Uint32(data[0:4]))
	tokenIdentifier, err := managedType.GetBytes(tokenIdentifierHandle)
	if err != nil {
		return nil, err
	}

	err = managedType.ConsumeGasForBytes(tokenIdentifier)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return nil, err
	}

	nonce := binary.BigEndian.Uint64(data[4:12])
	valueHandle := int32(binary.BigEndian.Uint32(data[12:16]))
	value, err := managedType.GetBigInt(valueHandle)
	if err != nil {
		return nil, err
	}

	err = managedType.ConsumeGasForBigIntCopy(value)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return nil, err
	}

	tokenType := core.Fungible
	if nonce > 0 {
		tokenType = core.NonFungible
	}

	return &vmcommon.DCDTTransfer{
		DCDTTokenName:  tokenIdentifier,
		DCDTTokenType:  uint32(tokenType),
		DCDTTokenNonce: nonce,
		DCDTValue:      value,
	}, nil
}

// Converts a managed buffer of serialized data to a slice of DCDTTransfer.
// The format is:
// - token identifier handle - 4 bytes
// - nonce - 8 bytes
// - value handle - 4 bytes
// Total: 16 bytes.
func readDCDTTransfers(
	managedType vmhost.ManagedTypesContext,
	runtime vmhost.RuntimeContext,
	managedVecHandle int32,
) ([]*vmcommon.DCDTTransfer, error) {
	managedVecBytes, err := managedType.GetBytes(managedVecHandle)
	if err != nil {
		return nil, err
	}

	err = managedType.ConsumeGasForBytes(managedVecBytes)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return nil, err
	}

	if len(managedVecBytes)%dcdtTransferLen != 0 {
		return nil, errors.New("invalid managed vector of DCDT transfers")
	}

	numTransfers := len(managedVecBytes) / dcdtTransferLen
	result := make([]*vmcommon.DCDTTransfer, 0, numTransfers)
	for i := 0; i < len(managedVecBytes); i += dcdtTransferLen {
		dcdtTransfer, err := readDCDTTransfer(managedType, runtime, managedVecBytes[i:i+dcdtTransferLen])
		if err != nil {
			return nil, err
		}
		result = append(result, dcdtTransfer)
	}

	return result, nil
}

// Serializes a vmcommon.DCDTTransfer object.
func writeDCDTTransfer(
	managedType vmhost.ManagedTypesContext,
	dcdtTransfer *vmcommon.DCDTTransfer,
	destinationBytes []byte,
) {
	tokenIdentifierHandle := managedType.NewManagedBufferFromBytes(dcdtTransfer.DCDTTokenName)
	valueHandle := managedType.NewBigInt(dcdtTransfer.DCDTValue)

	binary.BigEndian.PutUint32(destinationBytes[0:4], uint32(tokenIdentifierHandle))
	binary.BigEndian.PutUint64(destinationBytes[4:12], dcdtTransfer.DCDTTokenNonce)
	binary.BigEndian.PutUint32(destinationBytes[12:16], uint32(valueHandle))
}

// Serializes a list of DCDTTransfer one after the other into a byte slice.
// The format is (for each DCDTTransfer):
// - token identifier handle - 4 bytes
// - nonce - 8 bytes
// - value handle - 4 bytes
// Total: 16 bytes.
func writeDCDTTransfersToBytes(
	managedType vmhost.ManagedTypesContext,
	dcdtTransfers []*vmcommon.DCDTTransfer,
) []byte {
	destinationBytes := make([]byte, dcdtTransferLen*len(dcdtTransfers))
	dataIndex := 0
	for _, dcdtTransfer := range dcdtTransfers {
		writeDCDTTransfer(managedType, dcdtTransfer, destinationBytes[dataIndex:dataIndex+dcdtTransferLen])
		dataIndex += dcdtTransferLen
	}

	return destinationBytes
}

type vmInputData struct {
	destination []byte
	function    string
	value       *big.Int
	arguments   [][]byte
}

func readDestinationValueFunctionArguments(
	host vmhost.VMHost,
	destHandle int32,
	valueHandle int32,
	functionHandle int32,
	argumentsHandle int32,
) (*vmInputData, error) {
	managedType := host.ManagedTypes()

	vmInput, err := readDestinationValueArguments(host, destHandle, valueHandle, argumentsHandle)
	if err != nil {
		return nil, err
	}

	function, err := managedType.GetBytes(functionHandle)
	if err != nil {
		return nil, err
	}
	vmInput.function = string(function)

	return vmInput, nil
}

func readDestinationValueArguments(
	host vmhost.VMHost,
	destHandle int32,
	valueHandle int32,
	argumentsHandle int32,
) (*vmInputData, error) {
	managedType := host.ManagedTypes()

	vmInput, err := readDestinationArguments(host, destHandle, argumentsHandle)
	if err != nil {
		return nil, err
	}

	vmInput.value, err = managedType.GetBigInt(valueHandle)
	if err != nil {
		return nil, err
	}

	return vmInput, nil
}

func readDestinationFunctionArguments(
	host vmhost.VMHost,
	destHandle int32,
	functionHandle int32,
	argumentsHandle int32,
) (*vmInputData, error) {
	managedType := host.ManagedTypes()

	vmInput, err := readDestinationArguments(host, destHandle, argumentsHandle)
	if err != nil {
		return nil, err
	}

	function, err := managedType.GetBytes(functionHandle)
	if err != nil {
		return nil, err
	}
	vmInput.function = string(function)

	return vmInput, nil
}

func readDestinationArguments(
	host vmhost.VMHost,
	destHandle int32,
	argumentsHandle int32,
) (*vmInputData, error) {
	managedType := host.ManagedTypes()
	metering := host.Metering()

	var err error
	vmInput := &vmInputData{}

	vmInput.destination, err = managedType.GetBytes(destHandle)
	if err != nil {
		return nil, err
	}

	vmInput.value = big.NewInt(0)
	data, actualLen, err := managedType.ReadManagedVecOfManagedBuffers(argumentsHandle)
	if err != nil {
		return nil, err
	}
	vmInput.arguments = data

	gasToUse := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, actualLen)
	err = metering.UseGasBounded(gasToUse)
	if err != nil && host.Runtime().UseGasBoundedShouldFailExecution() {
		return nil, err
	}

	return vmInput, nil
}
