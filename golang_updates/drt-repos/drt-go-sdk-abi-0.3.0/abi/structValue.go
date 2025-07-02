package abi

import (
	"bytes"
	"fmt"
	"io"
)

// StructValue is a struct (collection of fields)
type StructValue struct {
	Fields []Field
}

// EncodeNested encodes the value in the nested form
func (value *StructValue) EncodeNested(writer io.Writer) error {
	for _, field := range value.Fields {
		err := field.Value.EncodeNested(writer)
		if err != nil {
			return fmt.Errorf("cannot encode field '%s' of struct, because of: %w", field.Name, err)
		}
	}

	return nil
}

// EncodeTopLevel encodes the value in the top-level form
func (value *StructValue) EncodeTopLevel(writer io.Writer) error {
	return value.EncodeNested(writer)
}

// DecodeNested decodes the value from the nested form
func (value *StructValue) DecodeNested(reader io.Reader) error {
	for _, field := range value.Fields {
		err := field.Value.DecodeNested(reader)
		if err != nil {
			return fmt.Errorf("cannot decode field '%s' of struct, because of: %w", field.Name, err)
		}
	}

	return nil
}

// DecodeTopLevel decodes the value from the top-level form
func (value *StructValue) DecodeTopLevel(data []byte) error {
	reader := bytes.NewReader(data)
	return value.DecodeNested(reader)
}
