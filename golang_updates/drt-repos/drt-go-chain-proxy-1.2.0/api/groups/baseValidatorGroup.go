package groups

import (
	"net/http"

	"github.com/TerraDharitri/drt-go-chain-proxy/api/shared"
	"github.com/TerraDharitri/drt-go-chain-proxy/data"
	"github.com/gin-gonic/gin"
)

type validatorGroup struct {
	facade ValidatorFacadeHandler
	*baseGroup
}

// NewValidatorGroup returns a new instance of validatorGroup
func NewValidatorGroup(facadeHandler data.FacadeHandler) (*validatorGroup, error) {
	facade, ok := facadeHandler.(ValidatorFacadeHandler)
	if !ok {
		return nil, ErrWrongTypeAssertion
	}

	vg := &validatorGroup{
		facade:    facade,
		baseGroup: &baseGroup{},
	}

	baseRoutesHandlers := []*data.EndpointHandlerData{
		{Path: "/statistics", Handler: vg.statistics, Method: http.MethodGet},
		{Path: "/auction", Handler: vg.auctionList, Method: http.MethodGet},
	}
	vg.baseGroup.endpoints = baseRoutesHandlers

	return vg, nil
}

// statistics returns the validator statistics
func (group *validatorGroup) statistics(c *gin.Context) {
	validatorStatistics, err := group.facade.ValidatorStatistics()
	if err != nil {
		shared.RespondWith(c, http.StatusBadRequest, nil, err.Error(), data.ReturnCodeRequestError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"statistics": validatorStatistics}, "", data.ReturnCodeSuccess)
}

func (group *validatorGroup) auctionList(c *gin.Context) {
	auctionList, err := group.facade.AuctionList()
	if err != nil {
		shared.RespondWith(c, http.StatusBadRequest, nil, err.Error(), data.ReturnCodeRequestError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"auctionList": auctionList}, "", data.ReturnCodeSuccess)
}
