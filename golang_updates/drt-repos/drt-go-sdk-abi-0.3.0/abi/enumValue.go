package abi

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

// EnumValue is an enum (discriminant and fields)
type EnumValue struct {
	Discriminant   uint8
	Fields         []Field
	FieldsProvider func(uint8) []Field
}

// EncodeNested encodes the value in the nested form
func (value *EnumValue) EncodeNested(writer io.Writer) error {
	discriminant := U8Value{Value: value.Discriminant}
	err := discriminant.EncodeNested(writer)
	if err != nil {
		return err
	}

	for _, field := range value.Fields {
		err := field.Value.EncodeNested(writer)
		if err != nil {
			return fmt.Errorf("cannot encode field '%s' of enum, because of: %w", field.Name, err)
		}
	}

	return nil
}

// EncodeTopLevel encodes the value in the top-level form
func (value *EnumValue) EncodeTopLevel(writer io.Writer) error {
	if value.Discriminant == 0 && len(value.Fields) == 0 {
		// Write nothing
		return nil
	}

	return value.EncodeNested(writer)
}

// DecodeNested decodes the value from the nested form
func (value *EnumValue) DecodeNested(reader io.Reader) error {
	if value.FieldsProvider == nil {
		return errors.New("cannot decode enum: fields provider is nil")
	}

	discriminant := &U8Value{}
	err := discriminant.DecodeNested(reader)
	if err != nil {
		return err
	}

	value.Discriminant = discriminant.Value
	value.Fields = value.FieldsProvider(value.Discriminant)

	for _, field := range value.Fields {
		err := field.Value.DecodeNested(reader)
		if err != nil {
			return fmt.Errorf("cannot decode field '%s' of enum, because of: %w", field.Name, err)
		}
	}

	return nil
}

// DecodeTopLevel decodes the value from the top-level form
func (value *EnumValue) DecodeTopLevel(data []byte) error {
	if len(data) == 0 {
		value.Discriminant = 0
		return nil
	}

	reader := bytes.NewReader(data)
	return value.DecodeNested(reader)
}
