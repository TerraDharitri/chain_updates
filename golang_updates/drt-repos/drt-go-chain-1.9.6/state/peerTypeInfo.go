package state

import "github.com/TerraDharitri/drt-go-chain-core/core"

// PeerTypeInfo contains information related to the peertypes needed by the peerTypeProvider
type PeerTypeInfo struct {
	PublicKey   string
	PeerType    string
	PeerSubType core.P2PPeerSubType
	ShardId     uint32
}
