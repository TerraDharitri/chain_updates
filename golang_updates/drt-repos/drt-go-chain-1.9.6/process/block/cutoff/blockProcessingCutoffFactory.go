package cutoff

import "github.com/TerraDharitri/drt-go-chain/config"

// CreateBlockProcessingCutoffHandler will create the desired block processing cutoff handler based on configuration
func CreateBlockProcessingCutoffHandler(cfg config.BlockProcessingCutoffConfig) (BlockProcessingCutoffHandler, error) {
	if !cfg.Enabled {
		return NewDisabledBlockProcessingCutoff(), nil
	}

	return NewBlockProcessingCutoffHandler(cfg)
}
