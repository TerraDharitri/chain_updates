package abi

import (
	"io"
)

// StringValue is a wrapper for a string
type StringValue struct {
	Value string
}

// EncodeNested encodes the value in the nested form
func (value *StringValue) EncodeNested(writer io.Writer) error {
	data := []byte(value.Value)
	err := encodeLength(writer, uint32(len(data)))
	if err != nil {
		return err
	}

	_, err = writer.Write(data)
	return err
}

// EncodeTopLevel encodes the value in the top-level form
func (value *StringValue) EncodeTopLevel(writer io.Writer) error {
	_, err := writer.Write([]byte(value.Value))
	return err
}

// DecodeNested decodes the value from the nested form
func (value *StringValue) DecodeNested(reader io.Reader) error {
	length, err := decodeLength(reader)
	if err != nil {
		return err
	}

	data, err := readBytesExactly(reader, int(length))
	if err != nil {
		return err
	}

	value.Value = string(data)
	return nil
}

// DecodeTopLevel decodes the value from the top-level form
func (value *StringValue) DecodeTopLevel(data []byte) error {
	value.Value = string(data)
	return nil
}
