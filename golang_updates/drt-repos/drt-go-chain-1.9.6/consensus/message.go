//go:generate protoc -I=. -I=$GOPATH/src -I=$GOPATH/src/github.com/TerraDharitri/protobuf/protobuf  --gogoslick_out=. message.proto
package consensus

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
)

// MessageType specifies what type of message was received
type MessageType int

// NewConsensusMessage creates a new Message object
func NewConsensusMessage(
	blHeaderHash []byte,
	signatureShare []byte,
	body []byte,
	header []byte,
	pubKey []byte,
	sig []byte,
	msg int,
	roundIndex int64,
	chainID []byte,
	pubKeysBitmap []byte,
	aggregateSignature []byte,
	leaderSignature []byte,
	currentPid core.PeerID,
	invalidSigners []byte,
) *Message {
	return &Message{
		BlockHeaderHash:    blHeaderHash,
		SignatureShare:     signatureShare,
		Body:               body,
		Header:             header,
		PubKey:             pubKey,
		Signature:          sig,
		MsgType:            int64(msg),
		RoundIndex:         roundIndex,
		ChainID:            chainID,
		PubKeysBitmap:      pubKeysBitmap,
		AggregateSignature: aggregateSignature,
		LeaderSignature:    leaderSignature,
		OriginatorPid:      currentPid.Bytes(),
		InvalidSigners:     invalidSigners,
	}
}
