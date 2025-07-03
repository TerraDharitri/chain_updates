package mock

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	crypto "github.com/TerraDharitri/drt-go-chain-crypto"
)

// PeerSignatureHandler -
type PeerSignatureHandler struct{}

// VerifyPeerSignature -
func (p *PeerSignatureHandler) VerifyPeerSignature(_ []byte, _ core.PeerID, _ []byte) error {
	return nil
}

// GetPeerSignature -
func (p *PeerSignatureHandler) GetPeerSignature(_ crypto.PrivateKey, _ []byte) ([]byte, error) {
	return nil, nil
}

// IsInterfaceNil -
func (p *PeerSignatureHandler) IsInterfaceNil() bool {
	return p == nil
}
