package abi

import (
	"bytes"
	"fmt"
	"io"
)

// OptionValue is a wrapper for an option value
type OptionValue struct {
	Value SingleValue
}

// EncodeNested encodes the value in the nested form
func (value *OptionValue) EncodeNested(writer io.Writer) error {
	if value.Value == nil {
		_, err := writer.Write([]byte{optionMarkerForAbsentValue})
		return err
	}

	_, err := writer.Write([]byte{optionMarkerForPresentValue})
	if err != nil {
		return err
	}

	return value.Value.EncodeNested(writer)
}

// EncodeTopLevel encodes the value in the top-level form
func (value *OptionValue) EncodeTopLevel(writer io.Writer) error {
	if value.Value == nil {
		return nil
	}

	_, err := writer.Write([]byte{optionMarkerForPresentValue})
	if err != nil {
		return err
	}

	return value.Value.EncodeNested(writer)
}

// DecodeNested decodes the value from the nested form
func (value *OptionValue) DecodeNested(reader io.Reader) error {
	if value.Value == nil {
		return fmt.Errorf("placeholder value of option should be set before decoding")
	}

	data, err := readBytesExactly(reader, 1)
	if err != nil {
		return err
	}

	firstByte := data[0]

	if firstByte == optionMarkerForAbsentValue {
		value.Value = nil
		return nil
	}

	if firstByte == optionMarkerForPresentValue {
		return value.Value.DecodeNested(reader)
	}

	return fmt.Errorf("invalid first byte for nested encoded option: %d", firstByte)
}

// DecodeTopLevel decodes the value from the top-level form
func (value *OptionValue) DecodeTopLevel(data []byte) error {
	if value.Value == nil {
		return fmt.Errorf("placeholder value of option should be set before decoding")
	}

	if len(data) == 0 {
		value.Value = nil
		return nil
	}

	firstByte := data[0]
	dataAfterFirstByte := data[1:]

	if firstByte != optionMarkerForPresentValue {
		return fmt.Errorf("invalid first byte for top-level encoded option: %d", firstByte)
	}

	reader := bytes.NewReader(dataAfterFirstByte)
	return value.Value.DecodeNested(reader)
}
