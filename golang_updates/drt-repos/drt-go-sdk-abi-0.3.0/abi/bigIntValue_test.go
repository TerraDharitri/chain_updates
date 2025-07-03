package abi

import (
	"math/big"
	"testing"
)

func TestBigIntValue(t *testing.T) {
	codec := &codec{}

	t.Run("should encode nested", func(t *testing.T) {
		testEncodeNested(t, codec, &BigIntValue{Value: big.NewInt(0)}, "00000000")
		testEncodeNested(t, codec, &BigIntValue{Value: big.NewInt(1)}, "0000000101")
		testEncodeNested(t, codec, &BigIntValue{Value: big.NewInt(-1)}, "00000001ff")
		testEncodeNested(t, codec, &BigIntValue{Value: big.NewInt(127)}, "000000017f")
		testEncodeNested(t, codec, &BigIntValue{Value: big.NewInt(128)}, "000000020080")
		testEncodeNested(t, codec, &BigIntValue{Value: big.NewInt(255)}, "0000000200ff")
		testEncodeNested(t, codec, &BigIntValue{Value: big.NewInt(256)}, "000000020100")
	})

	t.Run("should encode top-level", func(t *testing.T) {
		testEncodeTopLevel(t, codec, &BigIntValue{Value: big.NewInt(0)}, "")
		testEncodeTopLevel(t, codec, &BigIntValue{Value: big.NewInt(1)}, "01")
		testEncodeTopLevel(t, codec, &BigIntValue{Value: big.NewInt(-1)}, "ff")
		testEncodeTopLevel(t, codec, &BigIntValue{Value: big.NewInt(127)}, "7f")
		testEncodeTopLevel(t, codec, &BigIntValue{Value: big.NewInt(128)}, "0080")
		testEncodeTopLevel(t, codec, &BigIntValue{Value: big.NewInt(255)}, "00ff")
		testEncodeTopLevel(t, codec, &BigIntValue{Value: big.NewInt(256)}, "0100")
	})

	t.Run("should decode nested", func(t *testing.T) {
		testDecodeNested(t, codec, "00000000", &BigIntValue{}, &BigIntValue{Value: big.NewInt(0)})
		testDecodeNested(t, codec, "0000000101", &BigIntValue{}, &BigIntValue{Value: big.NewInt(1)})
		testDecodeNested(t, codec, "00000001ff", &BigIntValue{}, &BigIntValue{Value: big.NewInt(-1)})
		testDecodeNested(t, codec, "000000017f", &BigIntValue{}, &BigIntValue{Value: big.NewInt(127)})
		testDecodeNested(t, codec, "000000020080", &BigIntValue{}, &BigIntValue{Value: big.NewInt(128)})
		testDecodeNested(t, codec, "0000000200ff", &BigIntValue{}, &BigIntValue{Value: big.NewInt(255)})
		testDecodeNested(t, codec, "000000020100", &BigIntValue{}, &BigIntValue{Value: big.NewInt(256)})
	})

	t.Run("should err on decode nested", func(t *testing.T) {
		testDecodeNestedWithError(t, codec, "0000000301", &BigIntValue{}, "cannot decode (nested) *abi.BigIntValue, because of: cannot read exactly 3 bytes")
	})

	t.Run("should decode top-level", func(t *testing.T) {
		testDecodeTopLevel(t, codec, "", &BigIntValue{}, &BigIntValue{Value: big.NewInt(0)})
		testDecodeTopLevel(t, codec, "01", &BigIntValue{}, &BigIntValue{Value: big.NewInt(1)})
		testDecodeTopLevel(t, codec, "ff", &BigIntValue{}, &BigIntValue{Value: big.NewInt(-1)})
		testDecodeTopLevel(t, codec, "7f", &BigIntValue{}, &BigIntValue{Value: big.NewInt(127)})
		testDecodeTopLevel(t, codec, "0080", &BigIntValue{}, &BigIntValue{Value: big.NewInt(128)})
		testDecodeTopLevel(t, codec, "00ff", &BigIntValue{}, &BigIntValue{Value: big.NewInt(255)})
		testDecodeTopLevel(t, codec, "0100", &BigIntValue{}, &BigIntValue{Value: big.NewInt(256)})
	})
}
