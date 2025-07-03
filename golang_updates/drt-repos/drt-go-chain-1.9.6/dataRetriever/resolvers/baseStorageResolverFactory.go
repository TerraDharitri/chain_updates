package resolvers

import "github.com/TerraDharitri/drt-go-chain/storage"

func createBaseStorageResolver(
	storer storage.Storer,
	isFullHistoryNode bool,
) baseStorageResolver {
	if isFullHistoryNode {
		return &baseFullHistoryResolver{
			storer: storer,
		}
	}

	return &baseSimpleResolver{
		storer: storer,
	}
}
