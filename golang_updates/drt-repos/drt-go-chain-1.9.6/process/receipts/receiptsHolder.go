package receipts

import (
	"fmt"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/common/holders"
	"github.com/TerraDharitri/drt-go-chain-core/data/batch"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
)

func marshalReceiptsHolder(holder common.ReceiptsHolder, marshaller marshal.Marshalizer) ([]byte, error) {
	receiptsBatch := &batch.Batch{}

	for _, miniBlock := range holder.GetMiniblocks() {
		miniblockBytes, err := marshaller.Marshal(miniBlock)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", errCannotMarshalReceipts, err)
		}

		receiptsBatch.Data = append(receiptsBatch.Data, miniblockBytes)
	}

	// No miniblocks, no other (to be defined) content
	if len(receiptsBatch.Data) == 0 {
		return make([]byte, 0), nil
	}

	receiptsBytes, err := marshaller.Marshal(receiptsBatch)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errCannotMarshalReceipts, err)
	}

	return receiptsBytes, nil
}

func unmarshalReceiptsHolder(receiptsBytes []byte, marshaller marshal.Marshalizer) (common.ReceiptsHolder, error) {
	if len(receiptsBytes) == 0 {
		return createEmptyReceiptsHolder(), nil
	}

	receiptsBatch := &batch.Batch{}
	err := marshaller.Unmarshal(receiptsBatch, receiptsBytes)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errCannotUnmarshalReceipts, err)
	}

	miniblocks := make([]*block.MiniBlock, 0, len(receiptsBatch.Data))

	for _, miniblockBytes := range receiptsBatch.Data {
		miniBlock := &block.MiniBlock{}
		err := marshaller.Unmarshal(miniBlock, miniblockBytes)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", errCannotUnmarshalReceipts, err)
		}

		miniblocks = append(miniblocks, miniBlock)
	}

	return holders.NewReceiptsHolder(miniblocks), nil
}

func createEmptyReceiptsHolder() common.ReceiptsHolder {
	return holders.NewReceiptsHolder([]*block.MiniBlock{})
}
