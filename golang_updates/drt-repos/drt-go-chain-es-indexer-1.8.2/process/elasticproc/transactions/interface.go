package transactions

import (
	"github.com/TerraDharitri/drt-go-chain-core/data/outport"
	datafield "github.com/TerraDharitri/drt-go-chain-vm-common/parsers/dataField"
)

// DataFieldParser defines what a data field parser should be able to do
type DataFieldParser interface {
	Parse(dataField []byte, sender, receiver []byte, numOfShards uint32) *datafield.ResponseParseData
}

type feeInfoHandler interface {
	GetFeeInfo() *outport.FeeInfo
}
