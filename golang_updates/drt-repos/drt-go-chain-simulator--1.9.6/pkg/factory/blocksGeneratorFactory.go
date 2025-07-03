package factory

import (
	"github.com/TerraDharitri/drt-go-chain-simulator/config"
	"github.com/TerraDharitri/drt-go-chain-simulator/pkg/process"
	"github.com/TerraDharitri/drt-go-chain-simulator/pkg/process/disabled"
)

// CreateBlocksGenerator creates a new instance of block generator
func CreateBlocksGenerator(simulator process.SimulatorHandler, config config.BlocksGeneratorConfig) (process.BlocksGenerator, error) {
	if config.AutoGenerateBlocks {
		return process.NewBlocksGenerator(process.ArgBlocksGenerator{
			Simulator:     simulator,
			BlockTimeInMs: config.BlockTimeInMs,
		})
	}

	return disabled.NewBlocksGenerator(), nil
}
