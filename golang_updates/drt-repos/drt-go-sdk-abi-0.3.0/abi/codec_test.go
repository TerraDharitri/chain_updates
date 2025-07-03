package abi

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

// testEncodeNested is a helper function to test nested encoding.
func testEncodeNested(t *testing.T, codec *codec, value SingleValue, expected string) {
	encoded, err := codec.EncodeNested(value)

	require.NoError(t, err)
	require.Equal(t, expected, hex.EncodeToString(encoded))
}

// testEncodeTopLevel is a helper function to test top-level encoding.
func testEncodeTopLevel(t *testing.T, codec *codec, value SingleValue, expected string) {
	encoded, err := codec.EncodeTopLevel(value)

	require.NoError(t, err)
	require.Equal(t, expected, hex.EncodeToString(encoded))
}

// testDecodeNested is a helper function to test nested decoding.
func testDecodeNested(t *testing.T, codec *codec, encodedData string, destination SingleValue, expected SingleValue) {
	data, _ := hex.DecodeString(encodedData)
	err := codec.DecodeNested(data, destination)

	require.NoError(t, err)
	require.Equal(t, expected, destination)
}

// testDecodeNestedWithError is a helper function to test nested decoding.
func testDecodeNestedWithError(t *testing.T, codec *codec, encodedData string, destination SingleValue, expectedError string) {
	data, _ := hex.DecodeString(encodedData)
	err := codec.DecodeNested(data, destination)

	require.ErrorContains(t, err, expectedError)
}

// testDecodeTopLevel is a helper function to test top-level decoding.
func testDecodeTopLevel(t *testing.T, codec *codec, encodedData string, destination SingleValue, expected SingleValue) {
	data, _ := hex.DecodeString(encodedData)
	err := codec.DecodeTopLevel(data, destination)

	require.NoError(t, err)
	require.Equal(t, expected, destination)
}

// testDecodeTopLevelWithError is a helper function to test top-level decoding.
func testDecodeTopLevelWithError(t *testing.T, codec *codec, encodedData string, destination SingleValue, expectedError string) {
	data, _ := hex.DecodeString(encodedData)
	err := codec.DecodeTopLevel(data, destination)

	require.ErrorContains(t, err, expectedError)
}
