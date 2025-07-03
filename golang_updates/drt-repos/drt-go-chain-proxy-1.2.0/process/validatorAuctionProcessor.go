package process

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-proxy/data"
)

// GetAuctionList returns the auction list from a metachain observer node
func (vsp *ValidatorStatisticsProcessor) GetAuctionList() (*data.AuctionListResponse, error) {
	observers, errFetchObs := vsp.proc.GetObservers(core.MetachainShardId, data.AvailabilityRecent)
	if errFetchObs != nil {
		return nil, errFetchObs
	}

	var valStatsResponse data.AuctionListAPIResponse
	for _, observer := range observers {
		_, err := vsp.proc.CallGetRestEndPoint(observer.Address, auctionListPath, &valStatsResponse)
		if err == nil {
			log.Info("auction list fetched from API", "observer", observer.Address)
			return &valStatsResponse.Data, nil
		}

		log.Error("getAuctionListFromApi", "observer", observer.Address, "error", err)
	}

	return nil, ErrAuctionListNotAvailable
}
