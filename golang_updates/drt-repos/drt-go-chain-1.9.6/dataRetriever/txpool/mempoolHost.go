package txpool

import (
	"bytes"
	"math/big"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	"github.com/TerraDharitri/drt-go-chain-vm-common/parsers"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/process"
)

type argsMempoolHost struct {
	txGasHandler txGasHandler
	marshalizer  marshal.Marshalizer
}

type mempoolHost struct {
	txGasHandler        txGasHandler
	callArgumentsParser process.CallArgumentsParser
	dcdtTransferParser  vmcommon.DCDTTransferParser
}

func newMempoolHost(args argsMempoolHost) (*mempoolHost, error) {
	if check.IfNil(args.txGasHandler) {
		return nil, dataRetriever.ErrNilTxGasHandler
	}
	if check.IfNil(args.marshalizer) {
		return nil, dataRetriever.ErrNilMarshalizer
	}

	argsParser := parsers.NewCallArgsParser()

	dcdtTransferParser, err := parsers.NewDCDTTransferParser(args.marshalizer)
	if err != nil {
		return nil, err
	}

	return &mempoolHost{
		txGasHandler:        args.txGasHandler,
		callArgumentsParser: argsParser,
		dcdtTransferParser:  dcdtTransferParser,
	}, nil
}

// ComputeTxFee computes the fee for a transaction.
func (host *mempoolHost) ComputeTxFee(tx data.TransactionWithFeeHandler) *big.Int {
	return host.txGasHandler.ComputeTxFee(tx)
}

// GetTransferredValue returns the value transferred by a transaction.
func (host *mempoolHost) GetTransferredValue(tx data.TransactionHandler) *big.Int {
	value := tx.GetValue()
	hasValue := value != nil && value.Sign() != 0
	if hasValue {
		// Early exit (optimization): a transaction can either bear a regular value or be a "MultiDCDTNFTTransfer".
		return value
	}

	data := tx.GetData()
	hasData := len(data) > 0
	if !hasData {
		// Early exit (optimization): no "MultiDCDTNFTTransfer" to parse.
		return tx.GetValue()
	}

	maybeMultiTransfer := bytes.HasPrefix(data, []byte(core.BuiltInFunctionMultiDCDTNFTTransfer))
	if !maybeMultiTransfer {
		// Early exit (optimization).
		return big.NewInt(0)
	}

	function, args, err := host.callArgumentsParser.ParseData(string(data))
	if err != nil {
		return big.NewInt(0)
	}

	if function != core.BuiltInFunctionMultiDCDTNFTTransfer {
		// Early exit (optimization).
		return big.NewInt(0)
	}

	dcdtTransfers, err := host.dcdtTransferParser.ParseDCDTTransfers(tx.GetSndAddr(), tx.GetRcvAddr(), function, args)
	if err != nil {
		return big.NewInt(0)
	}

	accumulatedNativeValue := big.NewInt(0)

	for _, transfer := range dcdtTransfers.DCDTTransfers {
		if transfer.DCDTTokenNonce != 0 {
			continue
		}
		if string(transfer.DCDTTokenName) != vmcommon.REWAIdentifier {
			// We only care about native transfers.
			continue
		}

		_ = accumulatedNativeValue.Add(accumulatedNativeValue, transfer.DCDTValue)
	}

	return accumulatedNativeValue
}

// IsInterfaceNil returns true if there is no value under the interface
func (host *mempoolHost) IsInterfaceNil() bool {
	return host == nil
}
