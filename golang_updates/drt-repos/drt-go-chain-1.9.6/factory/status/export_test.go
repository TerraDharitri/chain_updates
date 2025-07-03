package status

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain/epochStart"
	outportDriverFactory "github.com/TerraDharitri/drt-go-chain/outport/factory"
	"github.com/TerraDharitri/drt-go-chain/p2p"
)

// EpochStartEventHandler -
func (pc *statusComponents) EpochStartEventHandler() epochStart.ActionHandler {
	return pc.epochStartEventHandler()
}

// ComputeNumConnectedPeers -
func ComputeNumConnectedPeers(
	appStatusHandler core.AppStatusHandler,
	netMessenger p2p.Messenger,
	suffix string,
) {
	computeNumConnectedPeers(appStatusHandler, netMessenger, suffix)
}

// ComputeConnectedPeers -
func ComputeConnectedPeers(
	appStatusHandler core.AppStatusHandler,
	netMessenger p2p.Messenger,
	suffix string,
) {
	computeConnectedPeers(appStatusHandler, netMessenger, suffix)
}

// MakeHostDriversArgs -
func (scf *statusComponentsFactory) MakeHostDriversArgs() ([]outportDriverFactory.ArgsHostDriverFactory, error) {
	return scf.makeHostDriversArgs()
}
