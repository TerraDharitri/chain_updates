package broadcast

import (
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain/consensus"
	"github.com/TerraDharitri/drt-go-chain/consensus/broadcast/shared"
)

// DelayedBroadcaster exposes functionality for handling the consensus members broadcasting of delay data
type DelayedBroadcaster interface {
	SetLeaderData(data *shared.DelayedBroadcastData) error
	SetValidatorData(data *shared.DelayedBroadcastData) error
	SetHeaderForValidator(vData *shared.ValidatorHeaderBroadcastData) error
	SetBroadcastHandlers(
		mbBroadcast func(mbData map[uint32][]byte, pkBytes []byte) error,
		txBroadcast func(txData map[string][][]byte, pkBytes []byte) error,
		headerBroadcast func(header data.HeaderHandler, pkBytes []byte) error,
		consensusMessageBroadcast func(message *consensus.Message) error,
	) error
	Close()
	IsInterfaceNil() bool
}
