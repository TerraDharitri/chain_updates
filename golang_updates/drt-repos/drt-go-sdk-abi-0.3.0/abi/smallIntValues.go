package abi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/big"

	twos "github.com/multiversx/mx-components-big-int/twos-complement"
)

// U8Value is a wrapper for uint8
type U8Value struct {
	Value uint8
}

// EncodeNested encodes the value in the nested form
func (value *U8Value) EncodeNested(writer io.Writer) error {
	return encodeNestedSmallInt(writer, value.Value, 1)
}

// EncodeTopLevel encodes the value in the top-level form
func (value *U8Value) EncodeTopLevel(writer io.Writer) error {
	return encodeTopLevelUnsignedSmallInt(writer, uint64(value.Value))
}

// DecodeNested decodes the value from the nested form
func (value *U8Value) DecodeNested(reader io.Reader) error {
	return decodeNestedSmallInt(reader, &value.Value, 1)
}

// DecodeTopLevel decodes the value from the top-level form
func (value *U8Value) DecodeTopLevel(data []byte) error {
	decoded, err := decodeTopLevelUnsignedSmallInt(data, math.MaxUint8)
	if err != nil {
		return err
	}

	value.Value = uint8(decoded)
	return nil
}

// U16Value is a wrapper for uint16
type U16Value struct {
	Value uint16
}

// EncodeNested encodes the value in the nested form
func (value *U16Value) EncodeNested(writer io.Writer) error {
	return encodeNestedSmallInt(writer, value.Value, 2)
}

// EncodeTopLevel encodes the value in the top-level form
func (value *U16Value) EncodeTopLevel(writer io.Writer) error {
	return encodeTopLevelUnsignedSmallInt(writer, uint64(value.Value))
}

// DecodeNested decodes the value from the nested form
func (value *U16Value) DecodeNested(reader io.Reader) error {
	return decodeNestedSmallInt(reader, &value.Value, 2)
}

// DecodeTopLevel decodes the value from the top-level form
func (value *U16Value) DecodeTopLevel(data []byte) error {
	decoded, err := decodeTopLevelUnsignedSmallInt(data, math.MaxUint16)
	if err != nil {
		return err
	}

	value.Value = uint16(decoded)
	return nil
}

// U32Value is a wrapper for uint16
type U32Value struct {
	Value uint32
}

// EncodeNested encodes the value in the nested form
func (value *U32Value) EncodeNested(writer io.Writer) error {
	return encodeNestedSmallInt(writer, value.Value, 4)
}

// EncodeTopLevel encodes the value in the top-level form
func (value *U32Value) EncodeTopLevel(writer io.Writer) error {
	return encodeTopLevelUnsignedSmallInt(writer, uint64(value.Value))
}

// DecodeNested decodes the value from the nested form
func (value *U32Value) DecodeNested(reader io.Reader) error {
	return decodeNestedSmallInt(reader, &value.Value, 4)
}

// DecodeTopLevel decodes the value from the top-level form
func (value *U32Value) DecodeTopLevel(data []byte) error {
	decoded, err := decodeTopLevelUnsignedSmallInt(data, math.MaxUint32)
	if err != nil {
		return err
	}

	value.Value = uint32(decoded)
	return nil
}

// U64Value is a wrapper for uint16
type U64Value struct {
	Value uint64
}

// EncodeNested encodes the value in the nested form
func (value *U64Value) EncodeNested(writer io.Writer) error {
	return encodeNestedSmallInt(writer, value.Value, 8)
}

// EncodeTopLevel encodes the value in the top-level form
func (value *U64Value) EncodeTopLevel(writer io.Writer) error {
	return encodeTopLevelUnsignedSmallInt(writer, uint64(value.Value))
}

// DecodeNested decodes the value from the nested form
func (value *U64Value) DecodeNested(reader io.Reader) error {
	return decodeNestedSmallInt(reader, &value.Value, 8)
}

// DecodeTopLevel decodes the value from the top-level form
func (value *U64Value) DecodeTopLevel(data []byte) error {
	decoded, err := decodeTopLevelUnsignedSmallInt(data, math.MaxUint64)
	if err != nil {
		return err
	}

	value.Value = uint64(decoded)
	return nil
}

// I8Value is a wrapper for uint8
type I8Value struct {
	Value int8
}

// EncodeNested encodes the value in the nested form
func (value *I8Value) EncodeNested(writer io.Writer) error {
	return encodeNestedSmallInt(writer, value.Value, 1)
}

// EncodeTopLevel encodes the value in the top-level form
func (value *I8Value) EncodeTopLevel(writer io.Writer) error {
	return encodeTopLevelSignedSmallInt(writer, int64(value.Value))
}

// DecodeNested decodes the value from the nested form
func (value *I8Value) DecodeNested(reader io.Reader) error {
	return decodeNestedSmallInt(reader, &value.Value, 1)
}

// DecodeTopLevel decodes the value from the top-level form
func (value *I8Value) DecodeTopLevel(data []byte) error {
	decoded, err := decodeTopLevelSignedSmallInt(data, math.MaxInt8)
	if err != nil {
		return err
	}

	value.Value = int8(decoded)
	return nil
}

// I16Value is a wrapper for uint16
type I16Value struct {
	Value int16
}

// EncodeNested encodes the value in the nested form
func (value *I16Value) EncodeNested(writer io.Writer) error {
	return encodeNestedSmallInt(writer, value.Value, 2)
}

// EncodeTopLevel encodes the value in the top-level form
func (value *I16Value) EncodeTopLevel(writer io.Writer) error {
	return encodeTopLevelSignedSmallInt(writer, int64(value.Value))
}

// DecodeNested decodes the value from the nested form
func (value *I16Value) DecodeNested(reader io.Reader) error {
	return decodeNestedSmallInt(reader, &value.Value, 2)
}

// DecodeTopLevel decodes the value from the top-level form
func (value *I16Value) DecodeTopLevel(data []byte) error {
	decoded, err := decodeTopLevelSignedSmallInt(data, math.MaxInt16)
	if err != nil {
		return err
	}

	value.Value = int16(decoded)
	return nil
}

// I32Value is a wrapper for uint16
type I32Value struct {
	Value int32
}

// EncodeNested encodes the value in the nested form
func (value *I32Value) EncodeNested(writer io.Writer) error {
	return encodeNestedSmallInt(writer, value.Value, 4)
}

// EncodeTopLevel encodes the value in the top-level form
func (value *I32Value) EncodeTopLevel(writer io.Writer) error {
	return encodeTopLevelSignedSmallInt(writer, int64(value.Value))
}

// DecodeNested decodes the value from the nested form
func (value *I32Value) DecodeNested(reader io.Reader) error {
	return decodeNestedSmallInt(reader, &value.Value, 4)
}

// DecodeTopLevel decodes the value from the top-level form
func (value *I32Value) DecodeTopLevel(data []byte) error {
	decoded, err := decodeTopLevelSignedSmallInt(data, math.MaxInt32)
	if err != nil {
		return err
	}

	value.Value = int32(decoded)
	return nil
}

// I64Value is a wrapper for uint16
type I64Value struct {
	Value int64
}

// EncodeNested encodes the value in the nested form
func (value *I64Value) EncodeNested(writer io.Writer) error {
	return encodeNestedSmallInt(writer, value.Value, 8)
}

// EncodeTopLevel encodes the value in the top-level form
func (value *I64Value) EncodeTopLevel(writer io.Writer) error {
	return encodeTopLevelSignedSmallInt(writer, int64(value.Value))
}

// DecodeNested decodes the value from the nested form
func (value *I64Value) DecodeNested(reader io.Reader) error {
	return decodeNestedSmallInt(reader, &value.Value, 8)
}

// DecodeTopLevel decodes the value from the top-level form
func (value *I64Value) DecodeTopLevel(data []byte) error {
	decoded, err := decodeTopLevelSignedSmallInt(data, math.MaxInt64)
	if err != nil {
		return err
	}

	value.Value = int64(decoded)
	return nil
}

func encodeNestedSmallInt(writer io.Writer, value any, numBytes int) error {
	buffer := new(bytes.Buffer)

	err := binary.Write(buffer, binary.BigEndian, value)
	if err != nil {
		return err
	}

	data := buffer.Bytes()
	if len(data) != numBytes {
		return fmt.Errorf("unexpected number of bytes: %d != %d", len(data), numBytes)
	}

	_, err = writer.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func encodeTopLevelUnsignedSmallInt(writer io.Writer, value uint64) error {
	b := big.NewInt(0).SetUint64(value)
	data := b.Bytes()
	_, err := writer.Write(data)
	return err
}

func encodeTopLevelSignedSmallInt(writer io.Writer, value int64) error {
	data := twos.ToBytes(big.NewInt(value))
	_, err := writer.Write(data)
	return err
}

func decodeNestedSmallInt(reader io.Reader, value any, numBytes int) error {
	data, err := readBytesExactly(reader, numBytes)
	if err != nil {
		return err
	}

	buffer := bytes.NewReader(data)
	err = binary.Read(buffer, binary.BigEndian, value)
	if err != nil {
		return err
	}

	return nil
}

func decodeTopLevelUnsignedSmallInt(data []byte, maxValue uint64) (uint64, error) {
	b := big.NewInt(0).SetBytes(data)
	if !b.IsUint64() {
		return 0, fmt.Errorf("decoded value is too large or invalid: %s", b)
	}

	n := b.Uint64()
	if n > maxValue {
		return 0, fmt.Errorf("decoded value is too large: %d > %d", n, maxValue)
	}

	return n, nil
}

func decodeTopLevelSignedSmallInt(data []byte, maxValue int64) (int64, error) {
	b := twos.FromBytes(data)

	if !b.IsInt64() {
		return 0, fmt.Errorf("decoded value is too large or invalid: %s", b)
	}

	n := b.Int64()
	if n > maxValue {
		return 0, fmt.Errorf("decoded value is too large: %d > %d", n, maxValue)
	}

	return n, nil
}
