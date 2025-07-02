package abi

import (
	"io"
)

// BytesValue is a wrapper for a byte slice
type BytesValue struct {
	Value []byte
}

// EncodeNested encodes the value in the nested form
func (value *BytesValue) EncodeNested(writer io.Writer) error {
	err := encodeLength(writer, uint32(len(value.Value)))
	if err != nil {
		return err
	}

	_, err = writer.Write(value.Value)
	return err
}

// EncodeTopLevel encodes the value in the top-level form
func (value *BytesValue) EncodeTopLevel(writer io.Writer) error {
	_, err := writer.Write(value.Value)
	return err
}

// DecodeNested decodes the value from the nested form
func (value *BytesValue) DecodeNested(reader io.Reader) error {
	length, err := decodeLength(reader)
	if err != nil {
		return err
	}

	data, err := readBytesExactly(reader, int(length))
	if err != nil {
		return err
	}

	value.Value = data
	return nil
}

// DecodeTopLevel decodes the value from the top-level form
func (value *BytesValue) DecodeTopLevel(data []byte) error {
	value.Value = data
	return nil
}
