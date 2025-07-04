package factory

import (
	"github.com/TerraDharitri/eth-chain-sovereign-notifier-go/config"
	"github.com/TerraDharitri/eth-chain-sovereign-notifier-go/process/client"
)

// CreateWSETHClientNotifier creates a ws eth client notifier
func CreateWSETHClientNotifier(cfg config.Config) (ETHClient, error) {
	return client.NewClient(cfg.ClientConfig.Url)
}
