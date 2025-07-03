package testscommon

import (
	"encoding/hex"

	"github.com/TerraDharitri/drt-go-chain-core/core/pubkeyConverter"
)

// RealWorldBech32PubkeyConverter is a bech32 converter, to be used in tests
var RealWorldBech32PubkeyConverter, _ = pubkeyConverter.NewBech32PubkeyConverter(32, "drt")

var (
	// TestAddressAlice is a test address
	TestAddressAlice = "drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf"
	// TestPubKeyAlice is a test pubkey
	TestPubKeyAlice, _ = RealWorldBech32PubkeyConverter.Decode(TestAddressAlice)
	// TestPubKeyHexAlice is a test pubkey
	TestPubKeyHexAlice = hex.EncodeToString(TestPubKeyAlice)

	// TestAddressBob is a test address
	TestAddressBob = "drt1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqlqde3c"
	// TestPubKeyBob is a test pubkey
	TestPubKeyBob, _ = RealWorldBech32PubkeyConverter.Decode(TestAddressBob)
	// TestPubKeyHexBob is a test pubkey
	TestPubKeyHexBob = hex.EncodeToString(TestPubKeyBob)
)
