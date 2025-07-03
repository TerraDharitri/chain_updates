package abi

import (
	"testing"
)

func TestCodecForOption(t *testing.T) {
	codec := &codec{}

	t.Run("should encode nested", func(t *testing.T) {
		testEncodeNested(t, codec, &OptionValue{
			Value: nil,
		}, "00")

		testEncodeNested(t, codec, &OptionValue{
			Value: &U16Value{Value: 0x08},
		}, "010008")
	})

	t.Run("should encode top-level", func(t *testing.T) {
		testEncodeTopLevel(t, codec, &OptionValue{
			Value: nil,
		}, "")

		testEncodeTopLevel(t, codec, &OptionValue{
			Value: &U16Value{Value: 0x08},
		}, "010008")
	})

	t.Run("should decode nested", func(t *testing.T) {
		testDecodeNested(t, codec, "00",
			&OptionValue{
				Value: &U8Value{},
			},
			&OptionValue{
				Value: nil,
			},
		)

		testDecodeNested(t, codec, "010008",
			&OptionValue{
				Value: &U16Value{},
			},
			&OptionValue{
				Value: &U16Value{Value: 0x08},
			},
		)
	})

	t.Run("should err on decode nested (nil placeholder)", func(t *testing.T) {
		testDecodeNestedWithError(t, codec, "072a", &OptionValue{}, "placeholder value of option should be set before decoding")
	})

	t.Run("should err on decode nested (bad marker for value presence)", func(t *testing.T) {
		testDecodeNestedWithError(t, codec, "072a", &OptionValue{Value: &BytesValue{}}, "invalid first byte for nested encoded option: 7")
	})

	t.Run("should decode top-level", func(t *testing.T) {
		testDecodeTopLevel(t, codec, "",
			&OptionValue{
				Value: &U8Value{},
			},
			&OptionValue{
				Value: nil,
			},
		)

		testDecodeTopLevel(t, codec, "010008",
			&OptionValue{
				Value: &U16Value{},
			},
			&OptionValue{
				Value: &U16Value{Value: 0x08},
			},
		)
	})

	t.Run("should err on decode top-level (nil placeholder)", func(t *testing.T) {
		testDecodeTopLevelWithError(t, codec, "072a", &OptionValue{}, "placeholder value of option should be set before decoding")
	})

	t.Run("should err on decode top-level (bad marker for value presence)", func(t *testing.T) {
		testDecodeTopLevelWithError(t, codec, "002a", &OptionValue{Value: &BytesValue{}}, "invalid first byte for top-level encoded option: 0")
		testDecodeTopLevelWithError(t, codec, "072a", &OptionValue{Value: &BytesValue{}}, "invalid first byte for top-level encoded option: 7")
	})
}
