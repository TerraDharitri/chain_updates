package abi

import (
	"testing"
)

func TestSmallIntValues(t *testing.T) {
	codec := &codec{}

	t.Run("should encode nested", func(t *testing.T) {
		testEncodeNested(t, codec, &U8Value{Value: 0x00}, "00")
		testEncodeNested(t, codec, &U8Value{Value: 0x01}, "01")
		testEncodeNested(t, codec, &U8Value{Value: 0x42}, "42")
		testEncodeNested(t, codec, &U8Value{Value: 0xff}, "ff")

		testEncodeNested(t, codec, &I8Value{Value: 0x00}, "00")
		testEncodeNested(t, codec, &I8Value{Value: 0x01}, "01")
		testEncodeNested(t, codec, &I8Value{Value: -1}, "ff")
		testEncodeNested(t, codec, &I8Value{Value: -128}, "80")
		testEncodeNested(t, codec, &I8Value{Value: 127}, "7f")

		testEncodeNested(t, codec, &U16Value{Value: 0x0000}, "0000")
		testEncodeNested(t, codec, &U16Value{Value: 0x0011}, "0011")
		testEncodeNested(t, codec, &U16Value{Value: 0x1234}, "1234")
		testEncodeNested(t, codec, &U16Value{Value: 0xffff}, "ffff")

		testEncodeNested(t, codec, &I16Value{Value: 0x0000}, "0000")
		testEncodeNested(t, codec, &I16Value{Value: 0x0011}, "0011")
		testEncodeNested(t, codec, &I16Value{Value: -1}, "ffff")
		testEncodeNested(t, codec, &I16Value{Value: -32768}, "8000")

		testEncodeNested(t, codec, &U32Value{Value: 0x00000000}, "00000000")
		testEncodeNested(t, codec, &U32Value{Value: 0x00000011}, "00000011")
		testEncodeNested(t, codec, &U32Value{Value: 0x00001122}, "00001122")
		testEncodeNested(t, codec, &U32Value{Value: 0x00112233}, "00112233")
		testEncodeNested(t, codec, &U32Value{Value: 0x11223344}, "11223344")
		testEncodeNested(t, codec, &U32Value{Value: 0xffffffff}, "ffffffff")

		testEncodeNested(t, codec, &I32Value{Value: 0x00000000}, "00000000")
		testEncodeNested(t, codec, &I32Value{Value: 0x00000011}, "00000011")
		testEncodeNested(t, codec, &I32Value{Value: -1}, "ffffffff")
		testEncodeNested(t, codec, &I32Value{Value: -2147483648}, "80000000")

		testEncodeNested(t, codec, &U64Value{Value: 0x0000000000000000}, "0000000000000000")
		testEncodeNested(t, codec, &U64Value{Value: 0x0000000000000011}, "0000000000000011")
		testEncodeNested(t, codec, &U64Value{Value: 0x0011223344556677}, "0011223344556677")
		testEncodeNested(t, codec, &U64Value{Value: 0xffffffffffffffff}, "ffffffffffffffff")

		testEncodeNested(t, codec, &I64Value{Value: 0x0000000000000000}, "0000000000000000")
		testEncodeNested(t, codec, &I64Value{Value: 0x0000000000000011}, "0000000000000011")
		testEncodeNested(t, codec, &I64Value{Value: -1}, "ffffffffffffffff")
		testEncodeNested(t, codec, &I64Value{Value: -9223372036854775808}, "8000000000000000")
	})

	t.Run("should encode top-level", func(t *testing.T) {
		testEncodeTopLevel(t, codec, &U8Value{Value: 0x00}, "")
		testEncodeTopLevel(t, codec, &U8Value{Value: 0x01}, "01")

		testEncodeTopLevel(t, codec, &I8Value{Value: 0x00}, "")
		testEncodeTopLevel(t, codec, &I8Value{Value: 0x01}, "01")
		testEncodeTopLevel(t, codec, &I8Value{Value: -1}, "ff")

		testEncodeTopLevel(t, codec, &U16Value{Value: 0x0000}, "")
		testEncodeTopLevel(t, codec, &U16Value{Value: 0x0011}, "11")

		testEncodeTopLevel(t, codec, &I16Value{Value: 0x0000}, "")
		testEncodeTopLevel(t, codec, &I16Value{Value: 0x0011}, "11")
		testEncodeTopLevel(t, codec, &I16Value{Value: -1}, "ff")

		testEncodeTopLevel(t, codec, &U32Value{Value: 0x00004242}, "4242")

		testEncodeTopLevel(t, codec, &I32Value{Value: 0x00000000}, "")
		testEncodeTopLevel(t, codec, &I32Value{Value: 0x00000011}, "11")
		testEncodeTopLevel(t, codec, &I32Value{Value: -1}, "ff")

		testEncodeTopLevel(t, codec, &U64Value{Value: 0x0042434445464748}, "42434445464748")

		testEncodeTopLevel(t, codec, &I64Value{Value: 0x0000000000000000}, "")
		testEncodeTopLevel(t, codec, &I64Value{Value: 0x0000000000000011}, "11")
		testEncodeTopLevel(t, codec, &I64Value{Value: -1}, "ff")
	})

	t.Run("should decode nested", func(t *testing.T) {
		testDecodeNested(t, codec, "00", &U8Value{}, &U8Value{Value: 0})
		testDecodeNested(t, codec, "01", &U8Value{}, &U8Value{Value: 1})
		testDecodeNested(t, codec, "ff", &U8Value{}, &U8Value{Value: 255})

		testDecodeNested(t, codec, "ff", &I8Value{}, &I8Value{Value: -1})

		testDecodeNested(t, codec, "4142", &U16Value{}, &U16Value{Value: 0x4142})
		testDecodeNested(t, codec, "ffff", &U16Value{}, &U16Value{Value: 65535})

		testDecodeNested(t, codec, "ffff", &I16Value{}, &I16Value{Value: -1})
		testDecodeNested(t, codec, "8000", &I16Value{}, &I16Value{Value: -32768})

		testDecodeNested(t, codec, "41424344", &U32Value{}, &U32Value{Value: 0x41424344})
		testDecodeNested(t, codec, "ffffffff", &U32Value{}, &U32Value{Value: 4294967295})

		testDecodeNested(t, codec, "ffffffff", &I32Value{}, &I32Value{Value: -1})
		testDecodeNested(t, codec, "80000000", &I32Value{}, &I32Value{Value: -2147483648})

		testDecodeNested(t, codec, "4142434445464748", &U64Value{}, &U64Value{Value: 0x4142434445464748})
		testDecodeNested(t, codec, "ffffffffffffffff", &U64Value{}, &U64Value{Value: 18446744073709551615})

		testDecodeNested(t, codec, "ffffffffffffffff", &I64Value{}, &I64Value{Value: -1})
		testDecodeNested(t, codec, "8000000000000000", &I64Value{}, &I64Value{Value: -9223372036854775808})
	})

	t.Run("should err on decode nested", func(t *testing.T) {
		testDecodeNestedWithError(t, codec, "01", &U16Value{}, "cannot read exactly 2 bytes")
		testDecodeNestedWithError(t, codec, "4142", &U32Value{}, "cannot read exactly 4 bytes")
		testDecodeNestedWithError(t, codec, "41424344", &U64Value{}, "cannot read exactly 8 bytes")
	})

	t.Run("should decode top-level", func(t *testing.T) {
		testDecodeTopLevel(t, codec, "", &U8Value{}, &U8Value{Value: 0})
		testDecodeNested(t, codec, "01", &U8Value{}, &U8Value{Value: 1})

		testDecodeTopLevel(t, codec, "ff", &I8Value{}, &I8Value{Value: -1})
		testDecodeTopLevel(t, codec, "80", &I8Value{}, &I8Value{Value: -128})

		testDecodeTopLevel(t, codec, "4242", &U16Value{}, &U16Value{Value: 0x4242})
		testDecodeTopLevel(t, codec, "ffff", &U16Value{}, &U16Value{Value: 65535})

		testDecodeTopLevel(t, codec, "ffff", &I16Value{}, &I16Value{Value: -1})
		testDecodeTopLevel(t, codec, "8000", &I16Value{}, &I16Value{Value: -32768})

		testDecodeTopLevel(t, codec, "41424344", &U32Value{}, &U32Value{Value: 0x41424344})
		testDecodeTopLevel(t, codec, "ffffffff", &U32Value{}, &U32Value{Value: 4294967295})

		testDecodeTopLevel(t, codec, "ffffffff", &I32Value{}, &I32Value{Value: -1})
		testDecodeTopLevel(t, codec, "80000000", &I32Value{}, &I32Value{Value: -2147483648})

		testDecodeTopLevel(t, codec, "4142434445464748", &U64Value{}, &U64Value{Value: 0x4142434445464748})
		testDecodeTopLevel(t, codec, "ffffffffffffffff", &U64Value{}, &U64Value{Value: 18446744073709551615})

		testDecodeTopLevel(t, codec, "ffffffffffffffff", &I64Value{}, &I64Value{Value: -1})
		testDecodeTopLevel(t, codec, "8000000000000000", &I64Value{}, &I64Value{Value: -9223372036854775808})
	})
}
