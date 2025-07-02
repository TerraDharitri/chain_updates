package testscommon

import (
	"math/big"

	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/dcdt"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
)

// SimpleNFTStorageHandlerStub -
type SimpleNFTStorageHandlerStub struct {
	GetDCDTNFTTokenOnDestinationCalled   func(accnt vmcommon.UserAccountHandler, dcdtTokenKey []byte, nonce uint64) (*dcdt.DCDigitalToken, bool, error)
	SaveNFTMetaDataToSystemAccountCalled func(tx data.TransactionHandler) error
}

// GetDCDTNFTTokenOnDestination -
func (s *SimpleNFTStorageHandlerStub) GetDCDTNFTTokenOnDestination(accnt vmcommon.UserAccountHandler, dcdtTokenKey []byte, nonce uint64) (*dcdt.DCDigitalToken, bool, error) {
	if s.GetDCDTNFTTokenOnDestinationCalled != nil {
		return s.GetDCDTNFTTokenOnDestinationCalled(accnt, dcdtTokenKey, nonce)
	}
	return &dcdt.DCDigitalToken{Value: big.NewInt(0)}, true, nil
}

// SaveNFTMetaData -
func (s *SimpleNFTStorageHandlerStub) SaveNFTMetaData(tx data.TransactionHandler) error {
	if s.SaveNFTMetaDataToSystemAccountCalled != nil {
		return s.SaveNFTMetaDataToSystemAccountCalled(tx)
	}
	return nil
}

// IsInterfaceNil -
func (s *SimpleNFTStorageHandlerStub) IsInterfaceNil() bool {
	return s == nil
}
