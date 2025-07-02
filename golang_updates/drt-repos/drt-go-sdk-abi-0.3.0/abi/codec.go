package abi

import (
	"bytes"
	"fmt"
)

// codec is a component which follows the rules of the MultiversX Serialization format:
// https://docs.multiversx.com/developers/data/serialization-overview
type codec struct {
}

// EncodeNested encodes the given value following the nested encoding rules
func (c *codec) EncodeNested(value SingleValue) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	err := value.EncodeNested(buffer)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// EncodeTopLevel encodes the given value following the top-level encoding rules
func (c *codec) EncodeTopLevel(value SingleValue) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	err := value.EncodeTopLevel(buffer)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// DecodeNested decodes the given data into the provided object following the nested decoding rules
func (c *codec) DecodeNested(data []byte, value SingleValue) error {
	reader := bytes.NewReader(data)
	err := value.DecodeNested(reader)
	if err != nil {
		return fmt.Errorf("cannot decode (nested) %T, because of: %w", value, err)
	}

	return nil
}

// DecodeTopLevel decodes the given data into the provided object following the top-level decoding rules
func (c *codec) DecodeTopLevel(data []byte, value SingleValue) error {
	err := value.DecodeTopLevel(data)
	if err != nil {
		return fmt.Errorf("cannot decode (top-level) %T, because of: %w", value, err)
	}

	return nil
}
