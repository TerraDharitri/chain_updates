package factory

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain/errors"
)

// DecodeAddresses will decode the provided string addresses
func DecodeAddresses(pkConverter core.PubkeyConverter, stringAddresses []string) ([][]byte, error) {
	if check.IfNil(pkConverter) {
		return nil, errors.ErrNilPubKeyConverter
	}
	decodedAddresses := make([][]byte, len(stringAddresses))
	for i, stringAddress := range stringAddresses {
		decodedAddress, errDecode := pkConverter.Decode(stringAddress)
		if errDecode != nil {
			return nil, errDecode
		}
		decodedAddresses[i] = decodedAddress
	}
	return decodedAddresses, nil
}
