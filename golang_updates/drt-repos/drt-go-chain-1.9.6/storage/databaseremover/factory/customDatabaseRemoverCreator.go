package factory

import (
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/storage"
	"github.com/TerraDharitri/drt-go-chain/storage/databaseremover"
	"github.com/TerraDharitri/drt-go-chain/storage/databaseremover/disabled"
)

// CreateCustomDatabaseRemover will handle the creation of a custom database remover based on the configuration
func CreateCustomDatabaseRemover(storagePruningConfig config.StoragePruningConfig) (storage.CustomDatabaseRemoverHandler, error) {
	if storagePruningConfig.AccountsTrieCleanOldEpochsData {
		return databaseremover.NewCustomDatabaseRemover(storagePruningConfig)
	}

	return disabled.NewDisabledCustomDatabaseRemover(), nil
}
