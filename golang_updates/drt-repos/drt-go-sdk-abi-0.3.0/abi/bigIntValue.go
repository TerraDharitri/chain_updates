package abi

import (
	"io"
	"math/big"

	twos "github.com/multiversx/mx-components-big-int/twos-complement"
)

// BigIntValue is a wrapper for a big integer (signed)
type BigIntValue struct {
	Value *big.Int
}

// EncodeNested encodes the value in the nested form
func (value *BigIntValue) EncodeNested(writer io.Writer) error {
	data := twos.ToBytes(value.Value)
	dataLength := len(data)

	// Write the length of the payload
	err := encodeLength(writer, uint32(dataLength))
	if err != nil {
		return err
	}

	// Write the payload
	_, err = writer.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// EncodeTopLevel encodes the value in the top-level form
func (value *BigIntValue) EncodeTopLevel(writer io.Writer) error {
	data := twos.ToBytes(value.Value)
	_, err := writer.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// DecodeNested decodes the value from the nested form
func (value *BigIntValue) DecodeNested(reader io.Reader) error {
	// Read the length of the payload
	length, err := decodeLength(reader)
	if err != nil {
		return err
	}

	// Read the payload
	data, err := readBytesExactly(reader, int(length))
	if err != nil {
		return err
	}

	value.Value = twos.FromBytes(data)
	return nil
}

// DecodeTopLevel decodes the value from the top-level form
func (value *BigIntValue) DecodeTopLevel(data []byte) error {
	value.Value = twos.FromBytes(data)
	return nil
}
