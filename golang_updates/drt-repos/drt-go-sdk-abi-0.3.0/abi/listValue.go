package abi

import (
	"bytes"
	"errors"
	"io"
)

// ListValue is a list of values
type ListValue struct {
	Items       []SingleValue
	ItemCreator func() SingleValue
}

// EncodeNested encodes the value in the nested form
func (value *ListValue) EncodeNested(writer io.Writer) error {
	err := encodeLength(writer, uint32(len(value.Items)))
	if err != nil {
		return err
	}

	return value.encodeItems(writer)
}

// EncodeTopLevel encodes the value in the top-level form
func (value *ListValue) EncodeTopLevel(writer io.Writer) error {
	return value.encodeItems(writer)
}

func (value *ListValue) encodeItems(writer io.Writer) error {
	for _, item := range value.Items {
		err := item.EncodeNested(writer)
		if err != nil {
			return err
		}
	}

	return nil
}

// DecodeNested decodes the value from the nested form
func (value *ListValue) DecodeNested(reader io.Reader) error {
	length, err := decodeLength(reader)
	if err != nil {
		return err
	}

	value.Items = make([]SingleValue, 0, length)

	for i := uint32(0); i < length; i++ {
		err := value.decodeItem(reader)
		if err != nil {
			return err
		}
	}

	return nil
}

// DecodeTopLevel decodes the value from the top-level form
func (value *ListValue) DecodeTopLevel(data []byte) error {
	reader := bytes.NewReader(data)
	value.Items = make([]SingleValue, 0)

	for reader.Len() > 0 {
		err := value.decodeItem(reader)
		if err != nil {
			return err
		}
	}

	return nil
}

func (value *ListValue) decodeItem(reader io.Reader) error {
	if value.ItemCreator == nil {
		return errors.New("cannot decode list: item creator is nil")
	}

	newItem := value.ItemCreator()

	err := newItem.DecodeNested(reader)
	if err != nil {
		return err
	}

	value.Items = append(value.Items, newItem)
	return nil
}
