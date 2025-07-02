package abi

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

type serializer struct {
	codec          *codec
	partsSeparator string
}

// ArgsNewSerializer defines the arguments needed for a new serializer
type ArgsNewSerializer struct {
	PartsSeparator string
}

// NewSerializer creates a new serializer.
// The serializer follows the rules of the MultiversX Serialization format:
// https://docs.multiversx.com/developers/data/serialization-overview
func NewSerializer(args ArgsNewSerializer) (*serializer, error) {
	if args.PartsSeparator == "" {
		return nil, errors.New("cannot create serializer: parts separator must not be empty")
	}

	codec := &codec{}

	return &serializer{
		codec:          codec,
		partsSeparator: args.PartsSeparator,
	}, nil
}

// Serialize serializes the given input values into a string
func (s *serializer) Serialize(inputValues []any) (string, error) {
	parts, err := s.serializeToParts(inputValues)
	if err != nil {
		return "", err
	}

	return s.encodeParts(parts), nil
}

func (s *serializer) serializeToParts(inputValues []any) ([][]byte, error) {
	partsHolder := newEmptyPartsHolder()

	err := s.doSerialize(partsHolder, inputValues)
	if err != nil {
		return nil, err
	}

	return partsHolder.getParts(), nil
}

func (s *serializer) doSerialize(partsHolder *partsHolder, inputValues []any) error {
	var err error

	for i, value := range inputValues {
		if value == nil {
			return errors.New("cannot serialize nil value")
		}

		switch value := value.(type) {
		case *OptionalValue:
			if i != len(inputValues)-1 {
				// Usage of multiple optional values is not recommended:
				// https://docs.multiversx.com/developers/data/multi-values
				// Thus, here, we disallow them.
				return errors.New("an optional value must be last among input values")
			}

			if value.Value != nil {
				err = s.doSerialize(partsHolder, []any{value.Value})
			}
		case *MultiValue:
			err = s.doSerialize(partsHolder, value.Items)
		case *VariadicValues:
			if i != len(inputValues)-1 {
				return errors.New("variadic values must be last among input values")
			}

			err = s.doSerialize(partsHolder, value.Items)
		case SingleValue:
			partsHolder.appendEmptyPart()
			err = s.serializeSingleValue(partsHolder, value)
		default:
			return fmt.Errorf("unsupported type for serialization: %T", value)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// Deserialize deserializes the given data into the output values
func (s *serializer) Deserialize(data string, outputValues []any) error {
	parts, err := s.decodeIntoParts(data)
	if err != nil {
		return err
	}

	return s.deserializeParts(parts, outputValues)
}

func (s *serializer) deserializeParts(parts [][]byte, outputValues []any) error {
	partsHolder := newPartsHolder(parts)

	err := s.doDeserialize(partsHolder, outputValues)
	if err != nil {
		return err
	}

	return nil
}

func (s *serializer) doDeserialize(partsHolder *partsHolder, outputValues []any) error {
	var err error

	for i, value := range outputValues {
		if value == nil {
			return errors.New("cannot deserialize into nil value")
		}

		switch value := value.(type) {
		case *OptionalValue:
			if i != len(outputValues)-1 {
				// Usage of multiple optional values is not recommended:
				// https://docs.multiversx.com/developers/data/multi-values
				// Thus, here, we disallow them.
				return errors.New("an optional value must be last among output values")
			}

			if partsHolder.isFocusedBeyondLastPart() {
				value.Value = nil
			} else {
				err = s.doDeserialize(partsHolder, []any{value.Value})
			}
		case *MultiValue:
			err = s.doDeserialize(partsHolder, value.Items)
		case *VariadicValues:
			if i != len(outputValues)-1 {
				return errors.New("variadic values must be last among output values")
			}

			err = s.deserializeVariadicValues(partsHolder, value)
		case SingleValue:
			err = s.deserializeSingleValue(partsHolder, value)
		default:
			return fmt.Errorf("unsupported type for deserialization: %T", value)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *serializer) serializeSingleValue(partsHolder *partsHolder, value SingleValue) error {
	data, err := s.codec.EncodeTopLevel(value)
	if err != nil {
		return err
	}

	return partsHolder.appendToLastPart(data)
}

func (s *serializer) deserializeVariadicValues(partsHolder *partsHolder, value *VariadicValues) error {
	if value.ItemCreator == nil {
		return errors.New("cannot deserialize variadic values: item creator is nil")
	}

	for !partsHolder.isFocusedBeyondLastPart() {
		newItem := value.ItemCreator()

		err := s.doDeserialize(partsHolder, []any{newItem})
		if err != nil {
			return err
		}

		value.Items = append(value.Items, newItem)
	}

	return nil
}

func (s *serializer) deserializeSingleValue(partsHolder *partsHolder, value SingleValue) error {
	part, err := partsHolder.readWholeFocusedPart()
	if err != nil {
		return err
	}

	err = s.codec.DecodeTopLevel(part, value)
	if err != nil {
		return err
	}

	err = partsHolder.focusOnNextPart()
	if err != nil {
		return err
	}

	return nil
}

func (s *serializer) encodeParts(parts [][]byte) string {
	partsHex := make([]string, len(parts))

	for i, part := range parts {
		partsHex[i] = hex.EncodeToString(part)
	}

	return strings.Join(partsHex, s.partsSeparator)
}

func (s *serializer) decodeIntoParts(encoded string) ([][]byte, error) {
	partsHex := strings.Split(encoded, s.partsSeparator)
	parts := make([][]byte, len(partsHex))

	for i, partHex := range partsHex {
		part, err := hex.DecodeString(partHex)
		if err != nil {
			return nil, err
		}

		parts[i] = part
	}

	return parts, nil
}
