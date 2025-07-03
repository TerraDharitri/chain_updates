package toml

import "github.com/TerraDharitri/drt-go-chain/config"

// OverrideConfig holds an array of configs to be overridden
type OverrideConfig struct {
	OverridableConfigTomlValues []config.OverridableConfig
}
