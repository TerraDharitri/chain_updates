package parsers

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/core/keyValStorage"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/errors"
	"github.com/TerraDharitri/drt-go-chain/state/dataTrieValue"
)

type dataTrieLeafParser struct {
	address             []byte
	marshaller          marshal.Marshalizer
	enableEpochsHandler common.EnableEpochsHandler
}

// NewDataTrieLeafParser returns a new instance of dataTrieLeafParser
func NewDataTrieLeafParser(address []byte, marshaller marshal.Marshalizer, enableEpochsHandler common.EnableEpochsHandler) (*dataTrieLeafParser, error) {
	if check.IfNil(marshaller) {
		return nil, errors.ErrNilMarshalizer
	}
	if check.IfNil(enableEpochsHandler) {
		return nil, errors.ErrNilEnableEpochsHandler
	}
	err := core.CheckHandlerCompatibility(enableEpochsHandler, []core.EnableEpochFlag{
		common.AutoBalanceDataTriesFlag,
	})
	if err != nil {
		return nil, err
	}

	return &dataTrieLeafParser{
		address:             address,
		marshaller:          marshaller,
		enableEpochsHandler: enableEpochsHandler,
	}, nil
}

// ParseLeaf returns a new KeyValStorage with the actual key and value
func (tlp *dataTrieLeafParser) ParseLeaf(trieKey []byte, trieVal []byte, version core.TrieNodeVersion) (core.KeyValueHolder, error) {
	isAutoBalanceDataTriesFlagEnabled := tlp.enableEpochsHandler.IsFlagEnabled(common.AutoBalanceDataTriesFlag)
	if isAutoBalanceDataTriesFlagEnabled && version == core.AutoBalanceEnabled {
		data := &dataTrieValue.TrieLeafData{}
		err := tlp.marshaller.Unmarshal(data, trieVal)
		if err != nil {
			return nil, err
		}

		return keyValStorage.NewKeyValStorage(data.Key, data.Value), nil
	}

	suffix := append(trieKey, tlp.address...)
	value, err := common.TrimSuffixFromValue(trieVal, len(suffix))
	if err != nil {
		return nil, err
	}

	return keyValStorage.NewKeyValStorage(trieKey, value), nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (tlp *dataTrieLeafParser) IsInterfaceNil() bool {
	return tlp == nil
}
