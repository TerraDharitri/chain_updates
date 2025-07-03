package groups

import (
	"fmt"
	"net/http"

	"github.com/TerraDharitri/drt-go-chain-proxy/api/errors"
	"github.com/TerraDharitri/drt-go-chain-proxy/api/shared"
	"github.com/TerraDharitri/drt-go-chain-proxy/data"
	"github.com/gin-gonic/gin"
)

type accountsGroup struct {
	facade AccountsFacadeHandler
	*baseGroup
}

// NewAccountsGroup returns a new instance of accountsGroup
func NewAccountsGroup(facadeHandler data.FacadeHandler) (*accountsGroup, error) {
	facade, ok := facadeHandler.(AccountsFacadeHandler)
	if !ok {
		return nil, ErrWrongTypeAssertion
	}

	ag := &accountsGroup{
		facade:    facade,
		baseGroup: &baseGroup{},
	}

	baseRoutesHandlers := []*data.EndpointHandlerData{
		{Path: "/:address", Handler: ag.getAccount, Method: http.MethodGet},
		{Path: "/:address/balance", Handler: ag.getBalance, Method: http.MethodGet},
		{Path: "/:address/username", Handler: ag.getUsername, Method: http.MethodGet},
		{Path: "/:address/nonce", Handler: ag.getNonce, Method: http.MethodGet},
		{Path: "/:address/shard", Handler: ag.getShard, Method: http.MethodGet},
		{Path: "/:address/code-hash", Handler: ag.getCodeHash, Method: http.MethodGet},
		{Path: "/:address/keys", Handler: ag.getKeyValuePairs, Method: http.MethodGet},
		{Path: "/:address/key/:key", Handler: ag.getValueForKey, Method: http.MethodGet},
		{Path: "/:address/dcdt", Handler: ag.getDCDTTokens, Method: http.MethodGet},
		{Path: "/:address/dcdt/:tokenIdentifier", Handler: ag.getDCDTTokenData, Method: http.MethodGet},
		{Path: "/:address/dcdts-with-role/:role", Handler: ag.getDCDTsWithRole, Method: http.MethodGet},
		{Path: "/:address/dcdts/roles", Handler: ag.getDCDTsRoles, Method: http.MethodGet},
		{Path: "/:address/registered-nfts", Handler: ag.getRegisteredNFTs, Method: http.MethodGet},
		{Path: "/:address/nft/:tokenIdentifier/nonce/:nonce", Handler: ag.getDCDTNftTokenData, Method: http.MethodGet},
		{Path: "/:address/guardian-data", Handler: ag.getGuardianData, Method: http.MethodGet},
		{Path: "/:address/is-data-trie-migrated", Handler: ag.isDataTrieMigrated, Method: http.MethodGet},
		{Path: "/bulk", Handler: ag.getAccounts, Method: http.MethodPost},
	}
	ag.baseGroup.endpoints = baseRoutesHandlers

	return ag, nil
}

func (group *accountsGroup) respondWithAccount(c *gin.Context, transform func(*data.AccountModel) gin.H) {
	address := c.Param("address")

	options, err := parseAccountQueryOptions(c, address)
	if err != nil {
		shared.RespondWithValidationError(c, errors.ErrBadUrlParams, err)
		return
	}

	model, err := group.facade.GetAccount(address, options)
	if err != nil {
		shared.RespondWithInternalError(c, errors.ErrGetAccount, err)
		return
	}

	response := transform(model)
	shared.RespondWith(c, http.StatusOK, response, "", data.ReturnCodeSuccess)
}

// getAccount returns an accountResponse containing information
// about the account correlated with provided address
func (group *accountsGroup) getAccount(c *gin.Context) {
	group.respondWithAccount(c, func(model *data.AccountModel) gin.H {
		return gin.H{"account": model.Account, "blockInfo": model.BlockInfo}
	})
}

// getBalance returns the balance for the address parameter
func (group *accountsGroup) getBalance(c *gin.Context) {
	group.respondWithAccount(c, func(model *data.AccountModel) gin.H {
		return gin.H{"balance": model.Account.Balance, "blockInfo": model.BlockInfo}
	})
}

// getUsername returns the username for the address parameter
func (group *accountsGroup) getUsername(c *gin.Context) {
	group.respondWithAccount(c, func(model *data.AccountModel) gin.H {
		return gin.H{"username": model.Account.Username, "blockInfo": model.BlockInfo}
	})
}

// getNonce returns the nonce for the address parameter
func (group *accountsGroup) getNonce(c *gin.Context) {
	group.respondWithAccount(c, func(model *data.AccountModel) gin.H {
		return gin.H{"nonce": model.Account.Nonce, "blockInfo": model.BlockInfo}
	})
}

// getCodeHash returns the code hash for the address parameter
func (group *accountsGroup) getCodeHash(c *gin.Context) {
	address := c.Param("address")
	options, err := parseAccountQueryOptions(c, address)
	if err != nil {
		shared.RespondWithValidationError(c, errors.ErrBadUrlParams, err)
		return
	}

	codeHashResponse, err := group.facade.GetCodeHash(address, options)
	if err != nil {
		shared.RespondWithInternalError(c, errors.ErrGetCodeHash, err)
		return
	}

	c.JSON(http.StatusOK, codeHashResponse)
}

// getAccounts will handle the request for a bulk of addresses data
func (group *accountsGroup) getAccounts(c *gin.Context) {
	var addresses []string
	err := c.ShouldBindJSON(&addresses)
	if err != nil {
		shared.RespondWithBadRequest(c, errors.ErrInvalidAddressesArray.Error())
		return
	}

	addr := ""
	if len(addresses) > 0 {
		addr = addresses[0]
	}

	options, err := parseAccountQueryOptions(c, addr)
	if err != nil {
		shared.RespondWithValidationError(c, errors.ErrInvalidFields, err)
		return
	}

	response, err := group.facade.GetAccounts(addresses, options)
	if err != nil {
		shared.RespondWithInternalError(c, errors.ErrCannotGetAddresses, err)
		return
	}

	shared.RespondWith(c, http.StatusOK, response, "", data.ReturnCodeSuccess)
}

// getKeyValuePairs returns the key-value pairs for the address parameter
func (group *accountsGroup) getKeyValuePairs(c *gin.Context) {
	addr := c.Param("address")
	if addr == "" {
		shared.RespondWithValidationError(c, errors.ErrGetKeyValuePairs, errors.ErrEmptyAddress)
		return
	}

	options, err := parseAccountQueryOptions(c, addr)
	if err != nil {
		shared.RespondWithValidationError(c, errors.ErrGetKeyValuePairs, err)
		return
	}

	keyValuePairs, err := group.facade.GetKeyValuePairs(addr, options)
	if err != nil {
		shared.RespondWithInternalError(c, errors.ErrGetKeyValuePairs, err)
		return
	}

	c.JSON(http.StatusOK, keyValuePairs)
}

// getValueForKey returns the value for the given address and key
func (group *accountsGroup) getValueForKey(c *gin.Context) {
	addr := c.Param("address")
	if addr == "" {
		shared.RespondWithValidationError(c, errors.ErrGetValueForKey, errors.ErrEmptyAddress)
		return
	}

	options, err := parseAccountQueryOptions(c, addr)
	if err != nil {
		shared.RespondWithValidationError(c, errors.ErrGetValueForKey, err)
		return
	}

	key := c.Param("key")
	if key == "" {
		shared.RespondWithValidationError(c, errors.ErrGetValueForKey, errors.ErrEmptyKey)
		return
	}

	value, err := group.facade.GetValueForKey(addr, key, options)
	if err != nil {
		shared.RespondWithInternalError(c, errors.ErrGetValueForKey, err)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"value": value}, "", data.ReturnCodeSuccess)
}

// getShard returns the shard for the given address based on the current proxy's configuration
func (group *accountsGroup) getShard(c *gin.Context) {
	addr := c.Param("address")
	if addr == "" {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%v: %v", errors.ErrComputeShardForAddress, errors.ErrEmptyAddress),
			data.ReturnCodeRequestError,
		)
		return
	}

	shardID, err := group.facade.GetShardIDForAddress(addr)
	if err != nil {
		shared.RespondWithInternalError(c, errors.ErrComputeShardForAddress, err)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"shardID": shardID}, "", data.ReturnCodeSuccess)
}

// getDCDTTokenData returns the balance for the given address and dcdt token
func (group *accountsGroup) getDCDTTokenData(c *gin.Context) {
	addr := c.Param("address")
	if addr == "" {
		shared.RespondWithValidationError(c, errors.ErrGetDCDTTokenData, errors.ErrEmptyAddress)
		return
	}

	options, err := parseAccountQueryOptions(c, addr)
	if err != nil {
		shared.RespondWithValidationError(c, errors.ErrGetDCDTTokenData, err)
		return
	}

	tokenIdentifier := c.Param("tokenIdentifier")
	if tokenIdentifier == "" {
		shared.RespondWithValidationError(c, errors.ErrEmptyTokenIdentifier, err)
		return
	}

	dcdtTokenResponse, err := group.facade.GetDCDTTokenData(addr, tokenIdentifier, options)
	if err != nil {
		shared.RespondWithInternalError(c, errors.ErrGetDCDTTokenData, err)
		return
	}

	c.JSON(http.StatusOK, dcdtTokenResponse)
}

func (group *accountsGroup) getDCDTsRoles(c *gin.Context) {
	addr := c.Param("address")
	if addr == "" {
		shared.RespondWithValidationError(c, errors.ErrGetRolesForAccount, errors.ErrEmptyAddress)
		return
	}

	options, err := parseAccountQueryOptions(c, addr)
	if err != nil {
		shared.RespondWithValidationError(c, errors.ErrGetRolesForAccount, err)
		return
	}

	tokensRoles, err := group.facade.GetDCDTsRoles(addr, options)
	if err != nil {
		shared.RespondWithInternalError(c, errors.ErrEmptyTokenIdentifier, err)
		return
	}

	c.JSON(http.StatusOK, tokensRoles)
}

// getDCDTsWithRole returns the token identifiers of the tokens where  the given address has the given role
func (group *accountsGroup) getDCDTsWithRole(c *gin.Context) {
	addr := c.Param("address")
	if addr == "" {
		shared.RespondWithValidationError(c, errors.ErrGetDCDTsWithRole, errors.ErrEmptyAddress)
		return
	}

	options, err := parseAccountQueryOptions(c, addr)
	if err != nil {
		shared.RespondWithValidationError(c, errors.ErrGetDCDTsWithRole, err)
		return
	}

	role := c.Param("role")
	if role == "" {
		shared.RespondWithValidationError(c, errors.ErrGetDCDTsWithRole, err)
		return
	}

	dcdtsWithRole, err := group.facade.GetDCDTsWithRole(addr, role, options)
	if err != nil {
		shared.RespondWithInternalError(c, errors.ErrGetDCDTsWithRole, err)
		return
	}

	c.JSON(http.StatusOK, dcdtsWithRole)
}

// getRegisteredNFTs returns the token identifiers of the NFTs registered by the address
func (group *accountsGroup) getRegisteredNFTs(c *gin.Context) {
	addr := c.Param("address")
	if addr == "" {
		shared.RespondWithValidationError(c, errors.ErrGetNFTTokenIDsRegisteredByAddress, errors.ErrEmptyAddress)
		return
	}

	options, err := parseAccountQueryOptions(c, addr)
	if err != nil {
		shared.RespondWithValidationError(c, errors.ErrGetNFTTokenIDsRegisteredByAddress, err)
		return
	}

	tokens, err := group.facade.GetNFTTokenIDsRegisteredByAddress(addr, options)
	if err != nil {
		shared.RespondWithInternalError(c, errors.ErrGetNFTTokenIDsRegisteredByAddress, err)
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// getDCDTNftTokenData returns the dcdt nft data for the given address, dcdt token and nonce
func (group *accountsGroup) getDCDTNftTokenData(c *gin.Context) {
	addr := c.Param("address")
	if addr == "" {
		shared.RespondWithValidationError(c, errors.ErrGetDCDTTokenData, errors.ErrEmptyAddress)
		return
	}

	options, err := parseAccountQueryOptions(c, addr)
	if err != nil {
		shared.RespondWithValidationError(c, errors.ErrGetDCDTTokenData, err)
		return
	}

	tokenIdentifier := c.Param("tokenIdentifier")
	if tokenIdentifier == "" {
		shared.RespondWithValidationError(c, errors.ErrGetDCDTTokenData, errors.ErrEmptyTokenIdentifier)
		return
	}

	nonce, err := shared.FetchNonceFromRequest(c)
	if err != nil {
		shared.RespondWithValidationError(c, errors.ErrGetDCDTTokenData, errors.ErrCannotParseNonce)
		return
	}

	dcdtTokenResponse, err := group.facade.GetDCDTNftTokenData(addr, tokenIdentifier, nonce, options)
	if err != nil {
		shared.RespondWithInternalError(c, errors.ErrGetDCDTTokenData, err)
		return
	}

	c.JSON(http.StatusOK, dcdtTokenResponse)
}

func (group *accountsGroup) getGuardianData(c *gin.Context) {
	addr := c.Param("address")
	if addr == "" {
		shared.RespondWithValidationError(c, errors.ErrGetGuardianData, errors.ErrEmptyAddress)
		return
	}

	options, err := parseAccountQueryOptions(c, addr)
	if err != nil {
		shared.RespondWithValidationError(c, errors.ErrGetGuardianData, err)
		return
	}

	guardianData, err := group.facade.GetGuardianData(addr, options)
	if err != nil {
		shared.RespondWithInternalError(c, errors.ErrGetGuardianData, err)
		return
	}

	c.JSON(http.StatusOK, guardianData)
}

// getDCDTTokens returns the tokens list from this account
func (group *accountsGroup) getDCDTTokens(c *gin.Context) {
	addr := c.Param("address")
	if addr == "" {
		shared.RespondWithValidationError(c, errors.ErrGetDCDTTokenData, errors.ErrEmptyAddress)
		return
	}

	options, err := parseAccountQueryOptions(c, addr)
	if err != nil {
		shared.RespondWithValidationError(c, errors.ErrGetDCDTTokenData, err)
		return
	}
	tokens, err := group.facade.GetAllDCDTTokens(addr, options)
	if err != nil {
		shared.RespondWithInternalError(c, errors.ErrGetDCDTTokenData, err)
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (group *accountsGroup) isDataTrieMigrated(c *gin.Context) {
	addr := c.Param("address")
	if addr == "" {
		shared.RespondWithValidationError(c, errors.ErrIsDataTrieMigrated, errors.ErrEmptyAddress)
		return
	}

	options, err := parseAccountQueryOptions(c, addr)
	if err != nil {
		shared.RespondWithValidationError(c, errors.ErrIsDataTrieMigrated, err)
		return
	}

	isMigrated, err := group.facade.IsDataTrieMigrated(addr, options)
	if err != nil {
		shared.RespondWithInternalError(c, errors.ErrIsDataTrieMigrated, err)
		return
	}

	c.JSON(http.StatusOK, isMigrated)
}
