package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain/api/errors"
	"github.com/TerraDharitri/drt-go-chain/api/shared"
)

type throttlerGetter interface {
	GetThrottlerForEndpoint(endpoint string) (core.Throttler, bool)
}

// CreateEndpointThrottlerFromFacade will create a middleware-type of handler to be used in conjunction with special
// REST API end points that need to be better protected
func CreateEndpointThrottlerFromFacade(throttlerName string, facade interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		tg, ok := facade.(throttlerGetter)
		if !ok {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				shared.GenericAPIResponse{
					Data:  nil,
					Error: errors.ErrInvalidAppContext.Error(),
					Code:  shared.ReturnCodeInternalError,
				},
			)
			return
		}

		endpointThrottler, ok := tg.GetThrottlerForEndpoint(throttlerName)
		if !ok {
			c.Next()
			return
		}

		if !endpointThrottler.CanProcess() {
			c.AbortWithStatusJSON(
				http.StatusTooManyRequests,
				shared.GenericAPIResponse{
					Data:  nil,
					Error: fmt.Sprintf("%s for endpoint %s", errors.ErrTooManyRequests.Error(), throttlerName),
					Code:  shared.ReturnCodeSystemBusy,
				},
			)
			return
		}

		endpointThrottler.StartProcessing()
		defer endpointThrottler.EndProcessing()

		c.Next()
	}
}
