package mock

import (
	"github.com/TerraDharitri/drt-go-chain-proxy/common"
	"github.com/TerraDharitri/drt-go-chain-proxy/data"
)

// AccountProcessorStub -
type AccountProcessorStub struct {
	GetAccountCalled                        func(address string, options common.AccountQueryOptions) (*data.AccountModel, error)
	GetAccountsCalled                       func(addresses []string, options common.AccountQueryOptions) (*data.AccountsModel, error)
	GetValueForKeyCalled                    func(address string, key string, options common.AccountQueryOptions) (string, error)
	GetShardIDForAddressCalled              func(address string) (uint32, error)
	GetTransactionsCalled                   func(address string) ([]data.DatabaseTransaction, error)
	ValidatorStatisticsCalled               func() (map[string]*data.ValidatorApiResponse, error)
	GetAllDCDTTokensCalled                  func(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetDCDTTokenDataCalled                  func(address string, key string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetDCDTNftTokenDataCalled               func(address string, key string, nonce uint64, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetDCDTsWithRoleCalled                  func(address string, role string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetNFTTokenIDsRegisteredByAddressCalled func(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetKeyValuePairsCalled                  func(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetDCDTsRolesCalled                     func(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetCodeHashCalled                       func(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetGuardianDataCalled                   func(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	IsDataTrieMigratedCalled                func(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
}

// GetKeyValuePairs -
func (aps *AccountProcessorStub) GetKeyValuePairs(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return aps.GetKeyValuePairsCalled(address, options)
}

// GetAllDCDTTokens -
func (aps *AccountProcessorStub) GetAllDCDTTokens(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return aps.GetAllDCDTTokensCalled(address, options)
}

// GetDCDTTokenData -
func (aps *AccountProcessorStub) GetDCDTTokenData(address string, key string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return aps.GetDCDTTokenDataCalled(address, key, options)
}

// GetDCDTNftTokenData -
func (aps *AccountProcessorStub) GetDCDTNftTokenData(address string, key string, nonce uint64, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return aps.GetDCDTNftTokenDataCalled(address, key, nonce, options)
}

// GetDCDTsWithRole -
func (aps *AccountProcessorStub) GetDCDTsWithRole(address string, role string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return aps.GetDCDTsWithRoleCalled(address, role, options)
}

// GetDCDTsRoles -
func (aps *AccountProcessorStub) GetDCDTsRoles(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	if aps.GetDCDTsRolesCalled != nil {
		return aps.GetDCDTsRolesCalled(address, options)
	}

	return &data.GenericAPIResponse{}, nil
}

// GetNFTTokenIDsRegisteredByAddress -
func (aps *AccountProcessorStub) GetNFTTokenIDsRegisteredByAddress(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return aps.GetNFTTokenIDsRegisteredByAddressCalled(address, options)
}

// GetAccount -
func (aps *AccountProcessorStub) GetAccount(address string, options common.AccountQueryOptions) (*data.AccountModel, error) {
	return aps.GetAccountCalled(address, options)
}

// GetAccounts -
func (aps *AccountProcessorStub) GetAccounts(addresses []string, options common.AccountQueryOptions) (*data.AccountsModel, error) {
	return aps.GetAccountsCalled(addresses, options)
}

// GetValueForKey -
func (aps *AccountProcessorStub) GetValueForKey(address string, key string, options common.AccountQueryOptions) (string, error) {
	return aps.GetValueForKeyCalled(address, key, options)
}

// GetGuardianData -
func (aps *AccountProcessorStub) GetGuardianData(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return aps.GetGuardianDataCalled(address, options)
}

// GetShardIDForAddress -
func (aps *AccountProcessorStub) GetShardIDForAddress(address string) (uint32, error) {
	return aps.GetShardIDForAddressCalled(address)
}

// GetCodeHash -
func (aps *AccountProcessorStub) GetCodeHash(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return aps.GetCodeHashCalled(address, options)
}

// ValidatorStatistics -
func (aps *AccountProcessorStub) ValidatorStatistics() (map[string]*data.ValidatorApiResponse, error) {
	return aps.ValidatorStatisticsCalled()
}

// IsDataTrieMigrated --
func (aps *AccountProcessorStub) IsDataTrieMigrated(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	if aps.IsDataTrieMigratedCalled != nil {
		return aps.IsDataTrieMigratedCalled(address, options)
	}

	return &data.GenericAPIResponse{}, nil
}

// AuctionList -
func (aps *AccountProcessorStub) AuctionList() ([]*data.AuctionListValidatorAPIResponse, error) {
	return nil, nil
}
