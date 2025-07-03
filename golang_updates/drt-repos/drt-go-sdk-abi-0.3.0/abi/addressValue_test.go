package abi

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddressValue(t *testing.T) {
	codec := &codec{}

	alicePubKeyHex := "0139472eff6886771a982f3083da5d421f24c29181e63888228dc81ca60d69e1"
	alicePubKey, _ := hex.DecodeString(alicePubKeyHex)

	shortPubKeyHex := "0139472eff6886771a982f3083da5d42"
	shortPubKey, _ := hex.DecodeString(shortPubKeyHex)

	t.Run("should encode nested", func(t *testing.T) {
		testEncodeNested(t, codec, &AddressValue{Value: alicePubKey}, alicePubKeyHex)
	})

	t.Run("should err on encode nested (bad public key length)", func(t *testing.T) {
		_, err := codec.EncodeNested(&AddressValue{Value: shortPubKey})
		require.ErrorContains(t, err, "public key (address) has invalid length")
	})

	t.Run("should encode top-level", func(t *testing.T) {
		testEncodeTopLevel(t, codec, &AddressValue{Value: alicePubKey}, alicePubKeyHex)
	})

	t.Run("should err on encode top-level (bad public key length)", func(t *testing.T) {
		_, err := codec.EncodeTopLevel(&AddressValue{Value: shortPubKey})
		require.ErrorContains(t, err, "public key (address) has invalid length")
	})

	t.Run("should decode nested", func(t *testing.T) {
		testDecodeNested(t, codec, alicePubKeyHex, &AddressValue{}, &AddressValue{Value: alicePubKey})
	})

	t.Run("should err on decode nested (shorter public key)", func(t *testing.T) {
		err := codec.DecodeNested(shortPubKey, &AddressValue{})
		require.ErrorContains(t, err, "cannot read exactly 32 bytes")
	})

	t.Run("should decode top-level", func(t *testing.T) {
		testDecodeTopLevel(t, codec, alicePubKeyHex, &AddressValue{}, &AddressValue{Value: alicePubKey})
	})

	t.Run("should err on decode top-level (shorter public key)", func(t *testing.T) {
		err := codec.DecodeTopLevel(shortPubKey, &AddressValue{})
		require.ErrorContains(t, err, "public key (address) has invalid length")
	})
}
