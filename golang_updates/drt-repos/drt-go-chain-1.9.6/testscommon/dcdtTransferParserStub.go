package testscommon

import (
	"errors"

	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
)

// DCDTTransferParserStub -
type DCDTTransferParserStub struct {
	ParseDCDTTransfersCalled func(sndAddr []byte, rcvAddr []byte, function string, args [][]byte) (*vmcommon.ParsedDCDTTransfers, error)
}

// ParseDCDTTransfers -
func (stub *DCDTTransferParserStub) ParseDCDTTransfers(sndAddr []byte, rcvAddr []byte, function string, args [][]byte) (*vmcommon.ParsedDCDTTransfers, error) {
	if stub.ParseDCDTTransfersCalled != nil {
		return stub.ParseDCDTTransfersCalled(sndAddr, rcvAddr, function, args)
	}

	return nil, errors.New("not implemented")
}

// IsInterfaceNil -
func (stub *DCDTTransferParserStub) IsInterfaceNil() bool {
	return stub == nil
}
