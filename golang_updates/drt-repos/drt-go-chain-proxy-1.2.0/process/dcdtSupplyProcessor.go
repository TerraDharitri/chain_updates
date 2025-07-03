package process

import (
	"math/big"
	"strings"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-proxy/data"
)

const (
	dcdtContractAddress   = "drt1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzllls6prdez"
	initialDCDTSupplyFunc = "getTokenProperties"

	networkDCDTSupplyPath = "/network/dcdt/supply/"
	zeroBigIntStr         = "0"
)

type dcdtSupplyProcessor struct {
	baseProc    Processor
	scQueryProc SCQueryService
}

// NewDCDTSupplyProcessor will create a new instance of the DCDT supply processor
func NewDCDTSupplyProcessor(baseProc Processor, scQueryProc SCQueryService) (*dcdtSupplyProcessor, error) {
	if check.IfNil(baseProc) {
		return nil, ErrNilCoreProcessor
	}
	if check.IfNil(scQueryProc) {
		return nil, ErrNilSCQueryService
	}

	return &dcdtSupplyProcessor{
		baseProc:    baseProc,
		scQueryProc: scQueryProc,
	}, nil
}

// GetDCDTSupply will return the total supply for the provided token
func (esp *dcdtSupplyProcessor) GetDCDTSupply(tokenIdentifier string) (*data.DCDTSupplyResponse, error) {
	totalSupply, err := esp.getSupplyFromShards(tokenIdentifier)
	if err != nil {
		return nil, err
	}

	res := &data.DCDTSupplyResponse{
		Code: data.ReturnCodeSuccess,
	}
	if !isFungibleDCDT(tokenIdentifier) {
		res.Data = *totalSupply
		makeInitialMintedNotEmpty(res)
		return res, nil
	}

	initialSupply, err := esp.getInitialSupplyFromMeta(tokenIdentifier)
	if err != nil {
		return nil, err
	}

	res.Data.InitialMinted = initialSupply.String()
	if totalSupply.RecomputedSupply {
		res.Data.Supply = totalSupply.Supply
		res.Data.Burned = zeroBigIntStr
		res.Data.Minted = zeroBigIntStr
		res.Data.RecomputedSupply = true
	} else {
		res.Data.Supply = sumStr(totalSupply.Supply, initialSupply.String())
		res.Data.Burned = totalSupply.Burned
		res.Data.Minted = totalSupply.Minted
	}

	makeInitialMintedNotEmpty(res)
	return res, nil
}

func makeInitialMintedNotEmpty(resp *data.DCDTSupplyResponse) {
	if resp.Data.InitialMinted == "" {
		resp.Data.InitialMinted = zeroBigIntStr
	}
}

func (esp *dcdtSupplyProcessor) getSupplyFromShards(tokenIdentifier string) (*data.DCDTSupply, error) {
	totalSupply := &data.DCDTSupply{}
	shardIDs := esp.baseProc.GetShardIDs()
	numNodesQueried := 0
	numNodesWithRecomputedSupply := 0
	for _, shardID := range shardIDs {
		if shardID == core.MetachainShardId {
			continue
		}

		supply, err := esp.getShardSupply(tokenIdentifier, shardID)
		if err != nil {
			return nil, err
		}

		addToSupply(totalSupply, supply)
		if supply.RecomputedSupply {
			numNodesWithRecomputedSupply++
		}
		numNodesQueried++
	}

	if numNodesWithRecomputedSupply > 0 {
		totalSupply.RecomputedSupply = true
	}

	return totalSupply, nil
}

func addToSupply(dstSupply, sourceSupply *data.DCDTSupply) {
	dstSupply.Supply = sumStr(dstSupply.Supply, sourceSupply.Supply)
	dstSupply.Burned = sumStr(dstSupply.Burned, sourceSupply.Burned)
	dstSupply.Minted = sumStr(dstSupply.Minted, sourceSupply.Minted)
}

func sumStr(s1, s2 string) string {
	s1Big, ok := big.NewInt(0).SetString(s1, 10)
	if !ok {
		return s2
	}
	s2Big, ok := big.NewInt(0).SetString(s2, 10)
	if !ok {
		return s1
	}

	return big.NewInt(0).Add(s1Big, s2Big).String()
}

func (esp *dcdtSupplyProcessor) getInitialSupplyFromMeta(token string) (*big.Int, error) {
	scQuery := &data.SCQuery{
		ScAddress: dcdtContractAddress,
		FuncName:  initialDCDTSupplyFunc,
		Arguments: [][]byte{[]byte(token)},
	}

	res, _, err := esp.scQueryProc.ExecuteQuery(scQuery)
	if err != nil {
		return nil, err
	}

	if len(res.ReturnData) < 4 {
		return big.NewInt(0), nil
	}

	supplyBytes := res.ReturnData[3]
	supplyBig, ok := big.NewInt(0).SetString(string(supplyBytes), 10)
	if !ok {
		return big.NewInt(0), nil
	}

	return supplyBig, nil
}

func (esp *dcdtSupplyProcessor) getShardSupply(token string, shardID uint32) (*data.DCDTSupply, error) {
	shardObservers, errObs := esp.baseProc.GetObservers(shardID, data.AvailabilityAll)
	if errObs != nil {
		return nil, errObs
	}

	responseDcdtSupply := data.DCDTSupplyResponse{}
	apiPath := networkDCDTSupplyPath + token
	for _, observer := range shardObservers {

		_, errGet := esp.baseProc.CallGetRestEndPoint(observer.Address, apiPath, &responseDcdtSupply)
		if errGet != nil {
			log.Error("dcdt supply request", "shard ID", observer.ShardId, "observer", observer.Address, "error", errGet.Error())
			continue
		}

		log.Info("dcdt supply request", "shard ID", observer.ShardId, "observer", observer.Address)

		return &responseDcdtSupply.Data, nil

	}

	return nil, WrapObserversError(responseDcdtSupply.Error)
}

func isFungibleDCDT(tokenIdentifier string) bool {
	splitToken := strings.Split(tokenIdentifier, "-")

	return len(splitToken) < 3
}
