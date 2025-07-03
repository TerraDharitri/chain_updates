package factory

import (
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/debug/handler"
)

// NewInterceptorDebuggerFactory will instantiate an InterceptorDebugHandler based on the provided config
func NewInterceptorDebuggerFactory(config config.InterceptorResolverDebugConfig) (InterceptorDebugHandler, error) {
	if !config.Enabled {
		return handler.NewDisabledInterceptorDebugHandler(), nil
	}

	return handler.NewInterceptorDebugHandler(config)
}
