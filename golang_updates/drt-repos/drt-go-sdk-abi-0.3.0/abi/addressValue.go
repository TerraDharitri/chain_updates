package abi

import (
	"fmt"
	"io"
)

// AddressValue is a wrapper for an address
type AddressValue struct {
	Value []byte
}

// EncodeNested encodes the value in the nested form
func (value *AddressValue) EncodeNested(writer io.Writer) error {
	err := value.checkPubKeyLength(value.Value)
	if err != nil {
		return err
	}

	_, err = writer.Write(value.Value)
	return err
}

// EncodeTopLevel encodes the value in the top-level form
func (value *AddressValue) EncodeTopLevel(writer io.Writer) error {
	return value.EncodeNested(writer)
}

// DecodeNested decodes the value from the nested form
func (value *AddressValue) DecodeNested(reader io.Reader) error {
	data, err := readBytesExactly(reader, pubKeyLength)
	if err != nil {
		return err
	}

	value.Value = data
	return nil
}

// DecodeTopLevel decodes the value from the top-level form
func (value *AddressValue) DecodeTopLevel(data []byte) error {
	err := value.checkPubKeyLength(data)
	if err != nil {
		return err
	}

	value.Value = data
	return nil
}

func (value *AddressValue) checkPubKeyLength(pubkey []byte) error {
	if len(pubkey) != pubKeyLength {
		return fmt.Errorf("public key (address) has invalid length: %d", len(pubkey))
	}

	return nil
}
