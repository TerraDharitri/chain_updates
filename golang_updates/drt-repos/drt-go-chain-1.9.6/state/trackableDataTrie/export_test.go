package trackableDataTrie

import "github.com/TerraDharitri/drt-go-chain-core/core"

// DirtyData -
type DirtyData struct {
	Value      []byte
	NewVersion core.TrieNodeVersion
}

// DirtyData -
func (tdt *trackableDataTrie) DirtyData() map[string]DirtyData {
	dd := make(map[string]DirtyData, len(tdt.dirtyData))

	for key, value := range tdt.dirtyData {
		dd[key] = DirtyData{
			Value:      value.value,
			NewVersion: value.newVersion,
		}
	}

	return dd
}
