package abi

import (
	"fmt"
	"io"
)

// BoolValue is a wrapper for a boolean
type BoolValue struct {
	Value bool
}

// EncodeNested encodes the value in the nested form
func (value *BoolValue) EncodeNested(writer io.Writer) error {
	if value.Value {
		_, err := writer.Write([]byte{trueAsByte})
		return err
	}

	_, err := writer.Write([]byte{falseAsByte})
	return err
}

// EncodeTopLevel encodes the value in the top-level form
func (value *BoolValue) EncodeTopLevel(writer io.Writer) error {
	if !value.Value {
		// For "false", write nothing.
		return nil
	}

	_, err := writer.Write([]byte{trueAsByte})
	return err
}

// DecodeNested decodes the value from the nested form
func (value *BoolValue) DecodeNested(reader io.Reader) error {
	data, err := readBytesExactly(reader, 1)
	if err != nil {
		return err
	}

	value.Value, err = value.byteToBool(data[0])
	if err != nil {
		return err
	}

	return nil
}

// DecodeTopLevel decodes the value from the top-level form
func (value *BoolValue) DecodeTopLevel(data []byte) error {
	if len(data) == 0 {
		value.Value = false
		return nil
	}

	if len(data) == 1 {
		boolValue, err := value.byteToBool(data[0])
		if err != nil {
			return err
		}

		value.Value = boolValue
		return nil
	}

	return fmt.Errorf("unexpected boolean value: %v", data)
}

func (value *BoolValue) byteToBool(data uint8) (bool, error) {
	switch data {
	case trueAsByte:
		return true, nil
	case falseAsByte:
		return false, nil
	default:
		return false, fmt.Errorf("unexpected boolean value: %d", data)
	}
}
