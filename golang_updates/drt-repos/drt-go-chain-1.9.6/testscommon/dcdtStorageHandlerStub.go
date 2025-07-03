package testscommon

import (
	"math/big"

	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/dcdt"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
)

// DcdtStorageHandlerStub -
type DcdtStorageHandlerStub struct {
	SaveDCDTNFTTokenCalled                                    func(senderAddress []byte, acnt vmcommon.UserAccountHandler, dcdtTokenKey []byte, nonce uint64, dcdtData *dcdt.DCDigitalToken, saveArgs vmcommon.NftSaveArgs) ([]byte, error)
	GetDCDTNFTTokenOnSenderCalled                             func(acnt vmcommon.UserAccountHandler, dcdtTokenKey []byte, nonce uint64) (*dcdt.DCDigitalToken, error)
	GetDCDTNFTTokenOnDestinationCalled                        func(acnt vmcommon.UserAccountHandler, dcdtTokenKey []byte, nonce uint64) (*dcdt.DCDigitalToken, bool, error)
	GetDCDTNFTTokenOnDestinationWithCustomSystemAccountCalled func(accnt vmcommon.UserAccountHandler, dcdtTokenKey []byte, nonce uint64, systemAccount vmcommon.UserAccountHandler) (*dcdt.DCDigitalToken, bool, error)
	WasAlreadySentToDestinationShardAndUpdateStateCalled      func(tickerID []byte, nonce uint64, dstAddress []byte) (bool, error)
	SaveNFTMetaDataCalled                                     func(tx data.TransactionHandler) error
	AddToLiquiditySystemAccCalled                             func(dcdtTokenKey []byte, tokenType uint32, nonce uint64, transferValue *big.Int, keepMetadataOnZeroLiquidity bool) error
	SaveMetaDataToSystemAccountCalled                         func(tokenKey []byte, nonce uint64, dcdtData *dcdt.DCDigitalToken) error
	GetMetaDataFromSystemAccountCalled                        func(bytes []byte, u uint64) (*dcdt.DCDigitalToken, error)
}

// SaveMetaDataToSystemAccount -
func (e *DcdtStorageHandlerStub) SaveMetaDataToSystemAccount(tokenKey []byte, nonce uint64, dcdtData *dcdt.DCDigitalToken) error {
	if e.SaveMetaDataToSystemAccountCalled != nil {
		return e.SaveMetaDataToSystemAccountCalled(tokenKey, nonce, dcdtData)
	}

	return nil
}

// GetMetaDataFromSystemAccount -
func (e *DcdtStorageHandlerStub) GetMetaDataFromSystemAccount(bytes []byte, u uint64) (*dcdt.DCDigitalToken, error) {
	if e.GetMetaDataFromSystemAccountCalled != nil {
		return e.GetMetaDataFromSystemAccountCalled(bytes, u)
	}

	return nil, nil
}

// SaveDCDTNFTToken -
func (e *DcdtStorageHandlerStub) SaveDCDTNFTToken(senderAddress []byte, acnt vmcommon.UserAccountHandler, dcdtTokenKey []byte, nonce uint64, dcdtData *dcdt.DCDigitalToken, saveArgs vmcommon.NftSaveArgs) ([]byte, error) {
	if e.SaveDCDTNFTTokenCalled != nil {
		return e.SaveDCDTNFTTokenCalled(senderAddress, acnt, dcdtTokenKey, nonce, dcdtData, saveArgs)
	}

	return nil, nil
}

// GetDCDTNFTTokenOnSender -
func (e *DcdtStorageHandlerStub) GetDCDTNFTTokenOnSender(acnt vmcommon.UserAccountHandler, dcdtTokenKey []byte, nonce uint64) (*dcdt.DCDigitalToken, error) {
	if e.GetDCDTNFTTokenOnSenderCalled != nil {
		return e.GetDCDTNFTTokenOnSenderCalled(acnt, dcdtTokenKey, nonce)
	}

	return nil, nil
}

// GetDCDTNFTTokenOnDestination -
func (e *DcdtStorageHandlerStub) GetDCDTNFTTokenOnDestination(acnt vmcommon.UserAccountHandler, dcdtTokenKey []byte, nonce uint64) (*dcdt.DCDigitalToken, bool, error) {
	if e.GetDCDTNFTTokenOnDestinationCalled != nil {
		return e.GetDCDTNFTTokenOnDestinationCalled(acnt, dcdtTokenKey, nonce)
	}

	return nil, false, nil
}

// GetDCDTNFTTokenOnDestinationWithCustomSystemAccount -
func (e *DcdtStorageHandlerStub) GetDCDTNFTTokenOnDestinationWithCustomSystemAccount(accnt vmcommon.UserAccountHandler, dcdtTokenKey []byte, nonce uint64, systemAccount vmcommon.UserAccountHandler) (*dcdt.DCDigitalToken, bool, error) {
	if e.GetDCDTNFTTokenOnDestinationWithCustomSystemAccountCalled != nil {
		return e.GetDCDTNFTTokenOnDestinationWithCustomSystemAccountCalled(accnt, dcdtTokenKey, nonce, systemAccount)
	}

	return nil, false, nil
}

// WasAlreadySentToDestinationShardAndUpdateState -
func (e *DcdtStorageHandlerStub) WasAlreadySentToDestinationShardAndUpdateState(tickerID []byte, nonce uint64, dstAddress []byte) (bool, error) {
	if e.WasAlreadySentToDestinationShardAndUpdateStateCalled != nil {
		return e.WasAlreadySentToDestinationShardAndUpdateStateCalled(tickerID, nonce, dstAddress)
	}

	return false, nil
}

// SaveNFTMetaData -
func (e *DcdtStorageHandlerStub) SaveNFTMetaData(tx data.TransactionHandler) error {
	if e.SaveNFTMetaDataCalled != nil {
		return e.SaveNFTMetaDataCalled(tx)
	}

	return nil
}

// AddToLiquiditySystemAcc -
func (e *DcdtStorageHandlerStub) AddToLiquiditySystemAcc(dcdtTokenKey []byte, tokenType uint32, nonce uint64, transferValue *big.Int, keepMetadataOnZeroLiquidity bool) error {
	if e.AddToLiquiditySystemAccCalled != nil {
		return e.AddToLiquiditySystemAccCalled(dcdtTokenKey, tokenType, nonce, transferValue, keepMetadataOnZeroLiquidity)
	}

	return nil
}

// IsInterfaceNil -
func (e *DcdtStorageHandlerStub) IsInterfaceNil() bool {
	return e == nil
}
