package abi

import (
	"io"
	"math/big"
)

// BigUIntValue is a wrapper for a big integer (unsigned)
type BigUIntValue struct {
	Value *big.Int
}

// EncodeNested encodes the value in the nested form
func (value *BigUIntValue) EncodeNested(writer io.Writer) error {
	data := value.Value.Bytes()
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
func (value *BigUIntValue) EncodeTopLevel(writer io.Writer) error {
	data := value.Value.Bytes()
	_, err := writer.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// DecodeNested decodes the value from the nested form
func (value *BigUIntValue) DecodeNested(reader io.Reader) error {
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

	value.Value = big.NewInt(0).SetBytes(data)
	return nil
}

// DecodeTopLevel decodes the value from the top-level form
func (value *BigUIntValue) DecodeTopLevel(data []byte) error {
	value.Value = big.NewInt(0).SetBytes(data)
	return nil
}
