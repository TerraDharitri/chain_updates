package process

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/state"
	disabledState "github.com/TerraDharitri/drt-go-chain/state/disabled"
	"github.com/TerraDharitri/drt-go-chain/state/storagePruningManager/disabled"
	"github.com/TerraDharitri/drt-go-chain/trie"
)

const maxTrieLevelInMemory = uint(5)

func createAccountAdapter(
	marshaller marshal.Marshalizer,
	hasher hashing.Hasher,
	accountFactory state.AccountFactory,
	trieStorage common.StorageManager,
	addressConverter core.PubkeyConverter,
	enableEpochsHandler common.EnableEpochsHandler,
) (state.AccountsAdapter, error) {
	tr, err := trie.NewTrie(trieStorage, marshaller, hasher, enableEpochsHandler, maxTrieLevelInMemory)
	if err != nil {
		return nil, err
	}

	args := state.ArgsAccountsDB{
		Trie:                  tr,
		Hasher:                hasher,
		Marshaller:            marshaller,
		AccountFactory:        accountFactory,
		StoragePruningManager: disabled.NewDisabledStoragePruningManager(),
		AddressConverter:      addressConverter,
		SnapshotsManager:      disabledState.NewDisabledSnapshotsManager(),
	}

	adb, err := state.NewAccountsDB(args)
	if err != nil {
		return nil, err
	}

	return adb, nil
}
