package abi

import (
	"testing"
)

func TestStructValue(t *testing.T) {
	codec := &codec{}

	t.Run("should encode nested", func(t *testing.T) {
		testEncodeNested(t, codec,
			&StructValue{
				Fields: []Field{
					{
						Value: &U8Value{Value: 0x01},
					},
					{
						Value: &U16Value{Value: 0x4142},
					},
				},
			},
			"014142",
		)
	})

	t.Run("should encode top-level", func(t *testing.T) {
		testEncodeTopLevel(t, codec,
			&StructValue{
				Fields: []Field{
					{
						Value: &U8Value{Value: 0x01},
					},
					{
						Value: &U16Value{Value: 0x4142},
					},
				},
			},
			"014142",
		)
	})

	t.Run("should decode nested", func(t *testing.T) {
		testDecodeNested(t, codec,
			"014142",
			&StructValue{
				Fields: []Field{
					{
						Value: &U8Value{},
					},
					{
						Value: &U16Value{},
					},
				},
			},
			&StructValue{
				Fields: []Field{
					{
						Value: &U8Value{Value: 0x01},
					},
					{
						Value: &U16Value{Value: 0x4142},
					},
				},
			},
		)
	})

	t.Run("should decode top-level", func(t *testing.T) {
		testDecodeTopLevel(t, codec,
			"014142",
			&StructValue{
				Fields: []Field{
					{
						Value: &U8Value{},
					},
					{
						Value: &U16Value{},
					},
				},
			},
			&StructValue{
				Fields: []Field{
					{
						Value: &U8Value{Value: 0x01},
					},
					{
						Value: &U16Value{Value: 0x4142},
					},
				},
			},
		)
	})
}
