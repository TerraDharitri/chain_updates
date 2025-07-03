package groups_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	apiErrors "github.com/TerraDharitri/drt-go-chain-proxy/api/errors"
	"github.com/TerraDharitri/drt-go-chain-proxy/api/groups"
	"github.com/TerraDharitri/drt-go-chain-proxy/api/mock"
	"github.com/TerraDharitri/drt-go-chain-proxy/common"
	"github.com/TerraDharitri/drt-go-chain-proxy/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const addressPath = "/address"

// General response structure
type GeneralResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

type accountResponseData struct {
	Account data.Account `json:"account"`
}

// accountResponse contains the account data and GeneralResponse fields
type accountResponse struct {
	GeneralResponse
	Data accountResponseData
}

type balanceResponseData struct {
	Balance string `json:"balance"`
}

// balanceResponse contains the balance and GeneralResponse fields
type balanceResponse struct {
	GeneralResponse
	Data balanceResponseData
}

type accountsResponseData struct {
	Accounts map[string]*data.Account `json:"accounts"`
}

type accountsResponse struct {
	GeneralResponse
	Data accountsResponseData `json:"data"`
}

type usernameResponseData struct {
	Username string `json:"username"`
}

// usernameResponse contains the username and GeneralResponse fields
type usernameResponse struct {
	GeneralResponse
	Data usernameResponseData
}

type getShardResponseData struct {
	ShardID uint32 `json:"shardID"`
}

type getShardResponse struct {
	GeneralResponse
	Data getShardResponseData
}

type guardianDataApiResponse struct {
	GeneralResponse
	Data guardianDataApiResponseData `json:"data"`
}

type guardianDataApiResponseData struct {
	GuardianData guardianData `json:"guardianData"`
}

type guardianData struct {
	ActiveGuardian  guardian `json:"activeGuardian,omitempty"`
	PendingGuardian guardian `json:"pendingGuardian,omitempty"`
	Guarded         bool     `json:"guarded,omitempty"`
}

type guardian struct {
	Address         string `json:"address"`
	ActivationEpoch uint32 `json:"activationEpoch"`
	ServiceUID      string `json:"serviceUID"`
}

type getDcdtTokensResponseData struct {
	Tokens []string `json:"tokens"`
}

type getDcdtTokensResponse struct {
	GeneralResponse
	Data getDcdtTokensResponseData
}

type dcdtTokenData struct {
	TokenIdentifier string `json:"tokenIdentifier"`
	Balance         string `json:"balance"`
	Properties      string `json:"properties"`
}

type dcdtNftData struct {
	TokenIdentifier string   `json:"tokenIdentifier"`
	Balance         string   `json:"balance"`
	Properties      string   `json:"properties"`
	Name            string   `json:"name"`
	Creator         string   `json:"creator"`
	Royalties       string   `json:"royalties"`
	Hash            []byte   `json:"hash"`
	URIs            [][]byte `json:"uris"`
	Attributes      []byte   `json:"attributes"`
}

type getDcdtTokenDataResponseData struct {
	TokenData dcdtTokenData `json:"tokenData"`
}

type getDcdtTokenDataResponse struct {
	GeneralResponse
	Data getDcdtTokenDataResponseData
}

type getDcdtNftTokenDataResponseData struct {
	TokenData dcdtNftData `json:"tokenData"`
}

type getDcdtNftTokenDataResponse struct {
	GeneralResponse
	Data getDcdtNftTokenDataResponseData
}

type getDCDTsRolesResponseData struct {
	Roles map[string][]string `json:"roles"`
}

type getDCDTsRolesResponse struct {
	GeneralResponse
	Data getDCDTsRolesResponseData
}

type getDcdtsWithRoleResponseData struct {
	Tokens []string `json:"tokenData"`
}

type getDcdtsWithRoleResponse struct {
	GeneralResponse
	Data getDcdtsWithRoleResponseData
}

type nonceResponseData struct {
	Nonce uint64 `json:"nonce"`
}

// nonceResponse contains the nonce and GeneralResponse fields
type nonceResponse struct {
	GeneralResponse
	Data nonceResponseData
}

func TestNewAccountGroup_WrongFacadeShouldErr(t *testing.T) {
	wrongFacade := &mock.WrongFacade{}
	group, err := groups.NewAccountsGroup(wrongFacade)
	require.Nil(t, group)
	require.Equal(t, groups.ErrWrongTypeAssertion, err)
}

func TestAddressRoute_EmptyTrailReturns404(t *testing.T) {
	t.Parallel()

	facade := &mock.FacadeStub{}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	req, _ := http.NewRequest("GET", "/address", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

//------- GetAccount

func TestGetAccount_FailWhenFacadeGetAccountFails(t *testing.T) {
	t.Parallel()

	returnedError := "i am an error"
	facade := &mock.FacadeStub{
		GetAccountHandler: func(address string, _ common.AccountQueryOptions) (*data.AccountModel, error) {
			return nil, errors.New(returnedError)
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	req, _ := http.NewRequest("GET", "/address/test", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	accountResponse := accountResponse{}
	loadResponse(resp.Body, &accountResponse)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Empty(t, accountResponse.Data)
	assert.Contains(t, accountResponse.Error, returnedError)
}

func TestGetAccount_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	facade := &mock.FacadeStub{
		GetAccountHandler: func(address string, _ common.AccountQueryOptions) (*data.AccountModel, error) {
			return &data.AccountModel{
				Account: data.Account{
					Address: address,
					Nonce:   1,
					Balance: "100",
				},
			}, nil
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	accountResponse := accountResponse{}
	loadResponse(resp.Body, &accountResponse)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, accountResponse.Data.Account.Address, reqAddress)
	assert.Equal(t, accountResponse.Data.Account.Nonce, uint64(1))
	assert.Equal(t, accountResponse.Data.Account.Balance, "100")
	assert.Empty(t, accountResponse.Error)
}

//------- GetAccounts

func TestGetAccount_FailsWhenInvalidRequest(t *testing.T) {
	t.Parallel()

	facade := &mock.FacadeStub{}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	req, _ := http.NewRequest("POST", "/address/bulk", bytes.NewBuffer([]byte(`invalid request`)))
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	accountsResponse := accountsResponse{}
	loadResponse(resp.Body, &accountsResponse)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Empty(t, accountsResponse.Data)
	assert.Equal(t, accountsResponse.Error, apiErrors.ErrInvalidAddressesArray.Error())
}

func TestGetAccount_FailWhenFacadeGetAccountsFails(t *testing.T) {
	t.Parallel()

	returnedError := "i am an error"
	facade := &mock.FacadeStub{
		GetAccountsHandler: func(addresses []string, _ common.AccountQueryOptions) (*data.AccountsModel, error) {
			return nil, errors.New(returnedError)
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	req, _ := http.NewRequest("POST", "/address/bulk", bytes.NewBuffer([]byte(`["test", "test"]`)))
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	accountsResponse := accountsResponse{}
	loadResponse(resp.Body, &accountsResponse)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Empty(t, accountsResponse.Data)
	assert.Contains(t, accountsResponse.Error, returnedError)
}

func TestGetAccounts_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	accounts := map[string]*data.Account{
		"drt1alice": {
			Address: "drt1alice",
			Nonce:   1,
			Balance: "100",
		},
		"drt1bob": {
			Address: "drt1bob",
			Nonce:   1,
			Balance: "101",
		},
	}
	facade := &mock.FacadeStub{
		GetAccountsHandler: func(addresses []string, _ common.AccountQueryOptions) (*data.AccountsModel, error) {
			return &data.AccountsModel{
				Accounts: accounts,
			}, nil
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddresses := []string{"drt1alice", "drt1bob"}
	addressBytes, _ := json.Marshal(reqAddresses)
	fmt.Println(string(addressBytes))
	req, _ := http.NewRequest("POST", "/address/bulk", bytes.NewBuffer(addressBytes))
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	accountsResponse := accountsResponse{}
	loadResponse(resp.Body, &accountsResponse)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, accountsResponse.Data.Accounts, accounts)
	assert.Empty(t, accountsResponse.Error)
}

//------- GetBalance

func TestGetBalance_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	facade := &mock.FacadeStub{
		GetAccountHandler: func(address string, _ common.AccountQueryOptions) (*data.AccountModel, error) {
			return &data.AccountModel{
				Account: data.Account{
					Address: address,
					Nonce:   1,
					Balance: "100",
				},
			}, nil
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/balance", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	balanceResponse := balanceResponse{}
	loadResponse(resp.Body, &balanceResponse)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, balanceResponse.Data.Balance, "100")
	assert.Empty(t, balanceResponse.Error)
}

//------- GetUsername

func TestGetUsername_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	expectedUsername := "testUser"
	facade := &mock.FacadeStub{
		GetAccountHandler: func(address string, _ common.AccountQueryOptions) (*data.AccountModel, error) {
			return &data.AccountModel{
				Account: data.Account{
					Address:  address,
					Nonce:    1,
					Balance:  "100",
					Username: expectedUsername,
				},
			}, nil
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/username", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	usernameResponse := usernameResponse{}
	loadResponse(resp.Body, &usernameResponse)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, expectedUsername, usernameResponse.Data.Username)
	assert.Empty(t, usernameResponse.Error)
}

//------- GetNonce

func TestGetNonce_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	facade := &mock.FacadeStub{
		GetAccountHandler: func(address string, _ common.AccountQueryOptions) (*data.AccountModel, error) {
			return &data.AccountModel{
				Account: data.Account{
					Address: address,
					Nonce:   1,
					Balance: "100",
				},
			}, nil
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/nonce", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	nonceResponse := nonceResponse{}
	loadResponse(resp.Body, &nonceResponse)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, uint64(1), nonceResponse.Data.Nonce)
	assert.Empty(t, nonceResponse.Error)
}

// ---- GetShard

func TestGetShard_FailWhenFacadeErrors(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("cannot compute shard ID")
	facade := &mock.FacadeStub{
		GetShardIDForAddressHandler: func(_ string) (uint32, error) {
			return 0, expectedErr
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/shard", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	shardResponse := getShardResponse{}
	loadResponse(resp.Body, &shardResponse)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.True(t, strings.Contains(shardResponse.Error, expectedErr.Error()))
}

func TestGetShard_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	expectedShardID := uint32(37)
	facade := &mock.FacadeStub{
		GetShardIDForAddressHandler: func(_ string) (uint32, error) {
			return expectedShardID, nil
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/shard", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	shardResponse := getShardResponse{}
	loadResponse(resp.Body, &shardResponse)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, shardResponse.Data.ShardID, expectedShardID)
	assert.Empty(t, shardResponse.Error)
}

// ---- GetDCDTTokens

func TestGetDCDTTokens_FailsWhenFacadeErrors(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("internal err")
	facade := &mock.FacadeStub{
		GetAllDCDTTokensCalled: func(_ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return nil, expectedErr
		},
	}

	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/dcdt", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	shardResponse := getDcdtTokensResponse{}
	loadResponse(resp.Body, &shardResponse)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.True(t, strings.Contains(shardResponse.Error, expectedErr.Error()))
}

func TestGetDCDTTokens_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	expectedTokens := []string{"abc", "def"}
	facade := &mock.FacadeStub{
		GetAllDCDTTokensCalled: func(_ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return &data.GenericAPIResponse{Data: getDcdtTokensResponseData{Tokens: expectedTokens}}, nil
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/dcdt", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	shardResponse := getDcdtTokensResponse{}
	loadResponse(resp.Body, &shardResponse)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, expectedTokens, shardResponse.Data.Tokens)
	assert.Empty(t, shardResponse.Error)
}

// ---- GetGuardianData

func TestGetGuardianData(t *testing.T) {
	t.Parallel()

	expectedGuardianData := guardianDataApiResponseData{
		GuardianData: guardianData{
			ActiveGuardian:  guardian{Address: "address1", ActivationEpoch: 0, ServiceUID: "serviceUID"},
			PendingGuardian: guardian{Address: "address2", ActivationEpoch: 1, ServiceUID: "serviceUID2"},
			Guarded:         false,
		}}

	expectedErr := errors.New("expected error")

	t.Run("internal error", func(t *testing.T) {
		t.Parallel()

		facade := &mock.FacadeStub{
			GetGuardianDataCalled: func(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
				return nil, expectedErr
			},
		}
		addressGroup, err := groups.NewAccountsGroup(facade)
		require.NoError(t, err)
		ws := startProxyServer(addressGroup, addressPath)
		reqAddress := "test"
		req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/guardian-data", reqAddress), nil)
		resp := httptest.NewRecorder()
		ws.ServeHTTP(resp, req)
		shardResponse := data.GenericAPIResponse{}
		loadResponse(resp.Body, &shardResponse)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.True(t, strings.Contains(shardResponse.Error, expectedErr.Error()))
	})
	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		facade := &mock.FacadeStub{
			GetGuardianDataCalled: func(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
				return &data.GenericAPIResponse{
					Data: expectedGuardianData,
				}, nil
			},
		}

		addressGroup, err := groups.NewAccountsGroup(facade)
		require.NoError(t, err)
		ws := startProxyServer(addressGroup, addressPath)
		reqAddress := "test"
		req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/guardian-data", reqAddress), nil)
		resp := httptest.NewRecorder()
		ws.ServeHTTP(resp, req)
		shardResponse := guardianDataApiResponse{}
		loadResponse(resp.Body, &shardResponse)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, expectedGuardianData, shardResponse.Data)
		assert.Empty(t, shardResponse.Error)
	})
}

// ---- GetDCDTsRoles

func TestGetDCDTsRoles_FailsWhenFacadeErrors(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("internal err")
	facade := &mock.FacadeStub{
		GetDCDTsRolesCalled: func(_ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return nil, expectedErr
		},
	}

	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/dcdts/roles", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	shardResponse := getDCDTsRolesResponse{}
	loadResponse(resp.Body, &shardResponse)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.True(t, strings.Contains(shardResponse.Error, expectedErr.Error()))
}

func TestGetDCDTsRoles_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	expectedRoles := map[string][]string{
		"tkn0": {"role0", "role1"},
		"tkn1": {"role1"},
	}
	facade := &mock.FacadeStub{
		GetDCDTsRolesCalled: func(_ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return &data.GenericAPIResponse{Data: getDCDTsRolesResponseData{Roles: expectedRoles}}, nil
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/dcdts/roles", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	shardResponse := getDCDTsRolesResponse{}
	loadResponse(resp.Body, &shardResponse)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, shardResponse.Data.Roles, expectedRoles)
	assert.Empty(t, shardResponse.Error)
}

// ---- GetDCDTTokenData

func TestGetDCDTTokenData_FailWhenFacadeErrors(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("internal err")
	facade := &mock.FacadeStub{
		GetDCDTTokenDataCalled: func(_ string, _ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return nil, expectedErr
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/dcdt/tkn", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	shardResponse := getDcdtTokenDataResponse{}
	loadResponse(resp.Body, &shardResponse)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.True(t, strings.Contains(shardResponse.Error, expectedErr.Error()))
}

func TestGetDCDTTokenData_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	expectedTokenData := dcdtTokenData{
		TokenIdentifier: "name",
		Balance:         "123",
		Properties:      "1",
	}
	facade := &mock.FacadeStub{
		GetDCDTTokenDataCalled: func(_ string, _ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return &data.GenericAPIResponse{Data: getDcdtTokenDataResponseData{TokenData: expectedTokenData}}, nil
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/dcdt/tkn", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	shardResponse := getDcdtTokenDataResponse{}
	loadResponse(resp.Body, &shardResponse)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, shardResponse.Data.TokenData, expectedTokenData)
	assert.Empty(t, shardResponse.Error)
}

// ---- GetDCDTNftTokenData

func TestGetDCDTNftTokenData_FailWhenFacadeErrors(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("internal err")
	facade := &mock.FacadeStub{
		GetDCDTNftTokenDataCalled: func(_ string, _ string, _ uint64, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return nil, expectedErr
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/nft/tkn/nonce/0", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	shardResponse := getDcdtNftTokenDataResponse{}
	loadResponse(resp.Body, &shardResponse)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.True(t, strings.Contains(shardResponse.Error, expectedErr.Error()))
}

func TestGetDCDTNftTokenData_FailWhenNonceParamIsInvalid(t *testing.T) {
	t.Parallel()

	facade := &mock.FacadeStub{}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/nft/tkn/nonce/qq", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := getDcdtNftTokenDataResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.True(t, strings.Contains(response.Error, apiErrors.ErrCannotParseNonce.Error()))
}

func TestGetDCDTNftTokenData_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	expectedTokenData := dcdtNftData{
		TokenIdentifier: "name",
		Balance:         "123",
		Properties:      "1",
		Royalties:       "10000",
	}
	facade := &mock.FacadeStub{
		GetDCDTNftTokenDataCalled: func(_ string, _ string, _ uint64, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return &data.GenericAPIResponse{Data: getDcdtNftTokenDataResponseData{TokenData: expectedTokenData}}, nil
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/nft/tkn/nonce/0", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := getDcdtNftTokenDataResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, response.Data.TokenData, expectedTokenData)
	assert.Empty(t, response.Error)
}

// ---- GetDCDTsWithRole

func TestGetDCDTsWithRole_FailWhenFacadeErrors(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("internal err")
	facade := &mock.FacadeStub{
		GetDCDTsWithRoleCalled: func(_ string, _ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return nil, expectedErr
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/dcdts-with-role/DCDTRoleNFTBurn", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	dcdtsWithRoleResponse := getDcdtsWithRoleResponse{}
	loadResponse(resp.Body, &dcdtsWithRoleResponse)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.True(t, strings.Contains(dcdtsWithRoleResponse.Error, expectedErr.Error()))
}

func TestGetDCDTsWithRole_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	expectedTokens := []string{"FDF-00rr44", "CVC-2598v7"}
	facade := &mock.FacadeStub{
		GetDCDTsWithRoleCalled: func(_ string, _ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return &data.GenericAPIResponse{Data: getDcdtsWithRoleResponseData{Tokens: expectedTokens}}, nil
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/dcdts-with-role/DCDTRoleNFTBurn", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := getDcdtsWithRoleResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, response.Data.Tokens, expectedTokens)
	assert.Empty(t, response.Error)
}

// ---- GetNFTTokenIDsRegisteredByAddress

func TestGetNFTTokenIDsRegisteredByAddress_FailWhenFacadeErrors(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("internal err")
	facade := &mock.FacadeStub{
		GetNFTTokenIDsRegisteredByAddressCalled: func(_ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return nil, expectedErr
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/registered-nfts", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	tokensResponse := getDcdtsWithRoleResponse{}
	loadResponse(resp.Body, &tokensResponse)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.True(t, strings.Contains(tokensResponse.Error, expectedErr.Error()))
}

func TestGetNFTTokenIDsRegisteredByAddress_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	expectedTokens := []string{"FDF-00rr44", "CVC-2598v7"}
	facade := &mock.FacadeStub{
		GetNFTTokenIDsRegisteredByAddressCalled: func(_ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return &data.GenericAPIResponse{Data: getDcdtsWithRoleResponseData{Tokens: expectedTokens}}, nil
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/registered-nfts", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := getDcdtsWithRoleResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, response.Data.Tokens, expectedTokens)
	assert.Empty(t, response.Error)
}

// ---- GetKeyValuePairs

func TestGetKeyValuePairs_FailWhenFacadeErrors(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("internal err")
	facade := &mock.FacadeStub{
		GetKeyValuePairsHandler: func(_ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return nil, expectedErr
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/keys", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := &data.GenericAPIResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.True(t, strings.Contains(response.Error, expectedErr.Error()))
}

func TestGetKeyValuePairs_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	expectedResponse := &data.GenericAPIResponse{
		Data: map[string]interface{}{"pairs": map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		}},
		Error: "",
		Code:  data.ReturnCodeSuccess,
	}
	facade := &mock.FacadeStub{
		GetKeyValuePairsHandler: func(_ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return expectedResponse, nil
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/keys", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	actualResponse := &data.GenericAPIResponse{}
	loadResponse(resp.Body, &actualResponse)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, expectedResponse, actualResponse)
	assert.Empty(t, actualResponse.Error)
}

// ---- get code hash

func TestGetCodeHash_FailWhenFacadeErrors(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("internal err")
	facade := &mock.FacadeStub{
		GetCodeHashCalled: func(_ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return nil, expectedErr
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/code-hash", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := &data.GenericAPIResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.True(t, strings.Contains(response.Error, expectedErr.Error()))
}

func TestGetCodeHash_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	expectedResponse := &data.GenericAPIResponse{
		Data:  "code hash",
		Error: "",
		Code:  data.ReturnCodeSuccess,
	}
	facade := &mock.FacadeStub{
		GetCodeHashCalled: func(_ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return expectedResponse, nil
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/code-hash", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	actualResponse := &data.GenericAPIResponse{}
	loadResponse(resp.Body, &actualResponse)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, expectedResponse, actualResponse)
	assert.Empty(t, actualResponse.Error)
}

func TestAccountsGroup_IsDataTrieMigrated(t *testing.T) {
	t.Parallel()

	t.Run("should return error when facade returns error", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("internal err")
		facade := &mock.FacadeStub{
			IsDataTrieMigratedCalled: func(_ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
				return nil, expectedErr
			},
		}
		addressGroup, err := groups.NewAccountsGroup(facade)
		require.NoError(t, err)
		ws := startProxyServer(addressGroup, addressPath)

		reqAddress := "test"
		req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/is-data-trie-migrated", reqAddress), nil)
		resp := httptest.NewRecorder()
		ws.ServeHTTP(resp, req)

		response := &data.GenericAPIResponse{}
		loadResponse(resp.Body, &response)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.True(t, strings.Contains(response.Error, expectedErr.Error()))
	})

	t.Run("should return successfully", func(t *testing.T) {
		t.Parallel()

		expectedResponse := &data.GenericAPIResponse{
			Data: map[string]interface{}{
				"isMigrated": "true",
			},
			Error: "",
			Code:  data.ReturnCodeSuccess,
		}
		facade := &mock.FacadeStub{
			IsDataTrieMigratedCalled: func(_ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
				return expectedResponse, nil
			},
		}
		addressGroup, err := groups.NewAccountsGroup(facade)
		require.NoError(t, err)
		ws := startProxyServer(addressGroup, addressPath)

		reqAddress := "test"
		req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/is-data-trie-migrated", reqAddress), nil)
		resp := httptest.NewRecorder()
		ws.ServeHTTP(resp, req)

		actualResponse := &data.GenericAPIResponse{}
		loadResponse(resp.Body, &actualResponse)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, expectedResponse, actualResponse)
		assert.Empty(t, actualResponse.Error)
	})
}
