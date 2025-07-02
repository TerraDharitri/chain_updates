package abi

import (
	"math/big"
	"testing"
)

func TestBigUIntValue(t *testing.T) {
	codec := &codec{}

	t.Run("should encode nested", func(t *testing.T) {
		testEncodeNested(t, codec, &BigUIntValue{Value: big.NewInt(0)}, "00000000")
		testEncodeNested(t, codec, &BigUIntValue{Value: big.NewInt(1)}, "0000000101")
		testEncodeNested(t, codec, &BigUIntValue{Value: big.NewInt(127)}, "000000017f")
		testEncodeNested(t, codec, &BigUIntValue{Value: big.NewInt(128)}, "0000000180")
		testEncodeNested(t, codec, &BigUIntValue{Value: big.NewInt(255)}, "00000001ff")
		testEncodeNested(t, codec, &BigUIntValue{Value: big.NewInt(256)}, "000000020100")
	})

	t.Run("should encode top-level", func(t *testing.T) {
		testEncodeTopLevel(t, codec, &BigUIntValue{Value: big.NewInt(0)}, "")
		testEncodeTopLevel(t, codec, &BigUIntValue{Value: big.NewInt(1)}, "01")
		testEncodeTopLevel(t, codec, &BigUIntValue{Value: big.NewInt(127)}, "7f")
		testEncodeTopLevel(t, codec, &BigUIntValue{Value: big.NewInt(128)}, "80")
		testEncodeTopLevel(t, codec, &BigUIntValue{Value: big.NewInt(256)}, "0100")
	})

	t.Run("should decode nested", func(t *testing.T) {
		testDecodeNested(t, codec, "00000000", &BigUIntValue{}, &BigUIntValue{Value: big.NewInt(0)})
		testDecodeNested(t, codec, "0000000101", &BigUIntValue{}, &BigUIntValue{Value: big.NewInt(1)})
		testDecodeNested(t, codec, "000000017f", &BigUIntValue{}, &BigUIntValue{Value: big.NewInt(127)})
		testDecodeNested(t, codec, "0000000180", &BigUIntValue{}, &BigUIntValue{Value: big.NewInt(128)})
		testDecodeNested(t, codec, "00000001ff", &BigUIntValue{}, &BigUIntValue{Value: big.NewInt(255)})
		testDecodeNested(t, codec, "000000020100", &BigUIntValue{}, &BigUIntValue{Value: big.NewInt(256)})
	})

	t.Run("should err on decode nested", func(t *testing.T) {
		testDecodeNestedWithError(t, codec, "0000000301", &BigUIntValue{}, "cannot decode (nested) *abi.BigUIntValue, because of: cannot read exactly 3 bytes")
	})

	t.Run("should decode top-level", func(t *testing.T) {
		testDecodeTopLevel(t, codec, "", &BigUIntValue{}, &BigUIntValue{Value: big.NewInt(0)})
		testDecodeTopLevel(t, codec, "01", &BigUIntValue{}, &BigUIntValue{Value: big.NewInt(1)})
		testDecodeTopLevel(t, codec, "7f", &BigUIntValue{}, &BigUIntValue{Value: big.NewInt(127)})
		testDecodeTopLevel(t, codec, "80", &BigUIntValue{}, &BigUIntValue{Value: big.NewInt(128)})
		testDecodeTopLevel(t, codec, "0100", &BigUIntValue{}, &BigUIntValue{Value: big.NewInt(256)})
	})
}
