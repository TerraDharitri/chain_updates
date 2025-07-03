package accounts

import (
	"github.com/TerraDharitri/drt-go-chain/state"
)

type dataTrieInteractor interface {
	state.DataTrieTracker
}
