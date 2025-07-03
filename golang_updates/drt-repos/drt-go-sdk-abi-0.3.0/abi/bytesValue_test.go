package abi

import (
	"testing"
)

func TestBytesValue(t *testing.T) {
	codec := &codec{}

	t.Run("should encode nested", func(t *testing.T) {
		testEncodeNested(t, codec, &BytesValue{Value: []byte{}}, "00000000")
		testEncodeNested(t, codec, &BytesValue{Value: []byte{'a', 'b', 'c'}}, "00000003616263")
	})

	t.Run("should encode top-level", func(t *testing.T) {
		testEncodeTopLevel(t, codec, &BytesValue{Value: []byte{}}, "")
		testEncodeTopLevel(t, codec, &BytesValue{Value: []byte{'a', 'b', 'c'}}, "616263")
	})

	t.Run("should decode nested", func(t *testing.T) {
		testDecodeNested(t, codec, "00000000", &BytesValue{}, &BytesValue{Value: []byte{}})
		testDecodeNested(t, codec, "00000003616263", &BytesValue{}, &BytesValue{Value: []byte{'a', 'b', 'c'}})
	})

	t.Run("should decode top-level", func(t *testing.T) {
		testDecodeTopLevel(t, codec, "", &BytesValue{}, &BytesValue{Value: []byte{}})
		testDecodeTopLevel(t, codec, "616263", &BytesValue{}, &BytesValue{Value: []byte{'a', 'b', 'c'}})
	})
}
