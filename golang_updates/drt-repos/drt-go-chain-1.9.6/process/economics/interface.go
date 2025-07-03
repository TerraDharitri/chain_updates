package economics

import (
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
)

// EpochNotifier raises epoch change events
type EpochNotifier interface {
	RegisterNotifyHandler(handler vmcommon.EpochSubscriberHandler)
	IsInterfaceNil() bool
}
