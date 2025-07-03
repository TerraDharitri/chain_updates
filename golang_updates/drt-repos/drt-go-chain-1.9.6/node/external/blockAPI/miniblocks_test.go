package blockAPI

import (
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/data/api"
	"github.com/stretchr/testify/require"
)

func TestFilterOutDuplicatedMiniblocks(t *testing.T) {
	miniblocks := []*api.MiniBlock{
		{Hash: "abba"},
		{Hash: "aabb"},
		{Hash: "aaaa"},
		{Hash: "abba"},
	}

	filteredMiniblocks := filterOutDuplicatedMiniblocks(miniblocks)

	expectedFilteredMiniblocks := []*api.MiniBlock{
		{Hash: "abba"},
		{Hash: "aabb"},
		{Hash: "aaaa"},
	}

	require.Equal(t, expectedFilteredMiniblocks, filteredMiniblocks)
}
