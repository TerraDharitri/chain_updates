package abi

import (
	"testing"
)

func TestStringValue(t *testing.T) {
	codec := &codec{}

	t.Run("should encode nested", func(t *testing.T) {
		testEncodeNested(t, codec, &StringValue{Value: ""}, "00000000")
		testEncodeNested(t, codec, &StringValue{Value: "abc"}, "00000003616263")
	})

	t.Run("should encode top-level", func(t *testing.T) {
		testEncodeTopLevel(t, codec, &StringValue{Value: ""}, "")
		testEncodeTopLevel(t, codec, &StringValue{Value: "abc"}, "616263")
	})

	t.Run("should decode nested", func(t *testing.T) {
		testDecodeNested(t, codec, "00000000", &StringValue{}, &StringValue{})
		testDecodeNested(t, codec, "00000003616263", &StringValue{}, &StringValue{Value: "abc"})
	})

	t.Run("should decode top-level", func(t *testing.T) {
		testDecodeTopLevel(t, codec, "", &StringValue{}, &StringValue{})
		testDecodeTopLevel(t, codec, "616263", &StringValue{}, &StringValue{Value: "abc"})
	})
}
