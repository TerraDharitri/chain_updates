package abi

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnumValue(t *testing.T) {
	codec := &codec{}

	t.Run("should encode nested", func(t *testing.T) {
		testEncodeNested(t, codec,
			&EnumValue{
				Discriminant: 0,
			},
			"00",
		)

		testEncodeNested(t, codec,
			&EnumValue{
				Discriminant: 42,
			},
			"2a",
		)

		testEncodeNested(t, codec,
			&EnumValue{
				Discriminant: 42,
				Fields: []Field{
					{
						Value: &U8Value{Value: 0x01},
					},
					{
						Value: &U16Value{Value: 0x4142},
					},
				},
			},
			"2a014142",
		)
	})

	t.Run("should encode top-level", func(t *testing.T) {
		testEncodeTopLevel(t, codec,
			&EnumValue{
				Discriminant: 0,
			},
			"",
		)

		testEncodeTopLevel(t, codec,
			&EnumValue{
				Discriminant: 42,
			},
			"2a",
		)

		testEncodeTopLevel(t, codec,
			&EnumValue{
				Discriminant: 42,
				Fields: []Field{
					{
						Value: &U8Value{Value: 0x01},
					},
					{
						Value: &U16Value{Value: 0x4142},
					},
				},
			},
			"2a014142",
		)
	})

	t.Run("should decode nested (simple)", func(t *testing.T) {
		data, _ := hex.DecodeString("2a")

		destination := &EnumValue{
			FieldsProvider: func(discriminant uint8) []Field {
				return nil
			},
		}

		err := codec.DecodeNested(data, destination)
		require.NoError(t, err)
		require.Equal(t, uint8(42), destination.Discriminant)
		require.Empty(t, destination.Fields)
	})

	t.Run("should decode nested (simple, zero)", func(t *testing.T) {
		data, _ := hex.DecodeString("00")

		destination := &EnumValue{
			FieldsProvider: func(discriminant uint8) []Field {
				return nil
			},
		}

		err := codec.DecodeNested(data, destination)
		require.NoError(t, err)
		require.Equal(t, uint8(0), destination.Discriminant)
		require.Empty(t, destination.Fields)
	})

	t.Run("should decode nested (heterogeneous)", func(t *testing.T) {
		data, _ := hex.DecodeString("01014142")

		destination := &EnumValue{
			FieldsProvider: func(discriminant uint8) []Field {
				return []Field{
					{
						Value: &U8Value{},
					},
					{
						Value: &U16Value{},
					},
				}
			},
		}

		err := codec.DecodeNested(data, destination)
		require.NoError(t, err)
		require.Equal(t, uint8(1), destination.Discriminant)
		require.Equal(t,
			[]Field{
				{
					Value: &U8Value{Value: 0x01},
				},
				{
					Value: &U16Value{Value: 0x4142},
				},
			},
			destination.Fields,
		)
	})

	t.Run("should decode top-level (simple)", func(t *testing.T) {
		data, _ := hex.DecodeString("2a")

		destination := &EnumValue{
			FieldsProvider: func(discriminant uint8) []Field {
				return nil
			},
		}

		err := codec.DecodeTopLevel(data, destination)
		require.NoError(t, err)
		require.Equal(t, uint8(42), destination.Discriminant)
		require.Empty(t, destination.Fields)
	})

	t.Run("should decode top-level (simple, zero)", func(t *testing.T) {
		data, _ := hex.DecodeString("")

		destination := &EnumValue{
			FieldsProvider: func(discriminant uint8) []Field {
				return nil
			},
		}

		err := codec.DecodeTopLevel(data, destination)
		require.NoError(t, err)
		require.Equal(t, uint8(0), destination.Discriminant)
		require.Empty(t, destination.Fields)
	})

	t.Run("should decode top-level (heterogeneous)", func(t *testing.T) {
		data, _ := hex.DecodeString("01014142")

		destination := &EnumValue{
			FieldsProvider: func(discriminant uint8) []Field {
				return []Field{
					{
						Value: &U8Value{},
					},
					{
						Value: &U16Value{},
					},
				}
			},
		}

		err := codec.DecodeNested(data, destination)
		require.NoError(t, err)
		require.Equal(t, uint8(1), destination.Discriminant)
		require.Equal(t,
			[]Field{
				{
					Value: &U8Value{Value: 0x01},
				},
				{
					Value: &U16Value{Value: 0x4142},
				},
			},
			destination.Fields,
		)
	})
}
