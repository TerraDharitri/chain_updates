package factory

import (
	"github.com/TerraDharitri/drt-go-chain-vm/crypto"
	"github.com/TerraDharitri/drt-go-chain-vm/crypto/hashing"
	"github.com/TerraDharitri/drt-go-chain-vm/crypto/signing/bls"
	"github.com/TerraDharitri/drt-go-chain-vm/crypto/signing/ed25519"
	"github.com/TerraDharitri/drt-go-chain-vm/crypto/signing/secp256"
)

// NewVMCrypto returns a composite struct containing VMCrypto functionality implementations
func NewVMCrypto() (crypto.VMCrypto, error) {
	blsVerifier, err := bls.NewBLS()
	if err != nil {
		return nil, err
	}

	secp, err := secp256.NewSecp256()
	if err != nil {
		return nil, err
	}

	return struct {
		crypto.Hasher
		crypto.Ed25519
		crypto.BLS
		crypto.Secp256
	}{
		Hasher:  hashing.NewHasher(),
		Ed25519: ed25519.NewEd25519Signer(),
		BLS:     blsVerifier,
		Secp256: secp,
	}, nil
}
