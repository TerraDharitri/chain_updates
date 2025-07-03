package process

import (
	"errors"
	"fmt"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data/vm"
	"github.com/TerraDharitri/drt-go-chain-proxy/data"
	"github.com/TerraDharitri/drt-go-chain-proxy/process/mock"
	"github.com/stretchr/testify/require"
)

func TestNewDCDTSupplyProcessor(t *testing.T) {
	t.Parallel()

	_, err := NewDCDTSupplyProcessor(nil, &mock.SCQueryServiceStub{})
	require.Equal(t, ErrNilCoreProcessor, err)

	_, err = NewDCDTSupplyProcessor(&mock.ProcessorStub{}, nil)
	require.Equal(t, ErrNilSCQueryService, err)
}

func TestDcdtSupplyProcessor_GetDCDTSupplyFungible(t *testing.T) {
	t.Parallel()

	baseProc := &mock.ProcessorStub{
		GetShardIDsCalled: func() []uint32 {
			return []uint32{0, 1, core.MetachainShardId}
		},
		GetObserversCalled: func(shardID uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return []*data.NodeData{
				{
					ShardId: shardID,
					Address: fmt.Sprintf("shard-%d", shardID),
				},
			}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			switch address {
			case "shard-0":
				valResp := value.(*data.DCDTSupplyResponse)
				valResp.Data.Supply = "1000"
				valResp.Data.Burned = "500"
				valResp.Data.Minted = "2000"
				return 200, nil
			case "shard-1":
				valResp := value.(*data.DCDTSupplyResponse)
				valResp.Data.Supply = "3000"
				valResp.Data.Burned = "100"
				valResp.Data.Minted = "300"
				return 200, nil
			}
			return 0, nil
		},
	}
	scQueryProc := &mock.SCQueryServiceStub{
		ExecuteQueryCalled: func(query *data.SCQuery) (*vm.VMOutputApi, data.BlockInfo, error) {
			return &vm.VMOutputApi{
				ReturnData: [][]byte{nil, nil, nil, []byte("500")},
			}, data.BlockInfo{}, nil
		},
	}
	dcdtProc, err := NewDCDTSupplyProcessor(baseProc, scQueryProc)
	require.Nil(t, err)

	supplyRes, err := dcdtProc.GetDCDTSupply("TOKEN-ABCD")
	require.Nil(t, err)
	require.Equal(t, "4500", supplyRes.Data.Supply)
	require.Equal(t, "600", supplyRes.Data.Burned)
	require.Equal(t, "2300", supplyRes.Data.Minted)
}

func TestDcdtSupplyProcessor_GetDCDTSupplyNonFungible(t *testing.T) {
	t.Parallel()

	called := false
	baseProc := &mock.ProcessorStub{
		GetShardIDsCalled: func() []uint32 {
			return []uint32{0, 1, core.MetachainShardId}
		},
		GetObserversCalled: func(shardID uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return []*data.NodeData{
				{
					ShardId: shardID,
					Address: fmt.Sprintf("shard-%d", shardID),
				},
				{
					ShardId: shardID,
					Address: fmt.Sprintf("shard-%d", shardID),
				},
			}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			switch address {
			case "shard-0":
				if !called {
					called = true
					return 400, errors.New("local err")
				}
				valResp := value.(*data.DCDTSupplyResponse)
				valResp.Data.Supply = "-1000"
				return 200, nil
			case "shard-1":
				valResp := value.(*data.DCDTSupplyResponse)
				valResp.Data.Supply = "3000"
				return 200, nil
			}
			return 0, nil
		},
	}
	scQueryProc := &mock.SCQueryServiceStub{}
	dcdtProc, err := NewDCDTSupplyProcessor(baseProc, scQueryProc)
	require.Nil(t, err)

	supplyRes, err := dcdtProc.GetDCDTSupply("SEMI-ABCD-0A")
	require.Nil(t, err)
	require.Equal(t, "2000", supplyRes.Data.Supply)
	require.Equal(t, "0", supplyRes.Data.InitialMinted)
}

func TestDcdtSupplyProcessor_GetDCDTSupplyShouldReturnRecomputed(t *testing.T) {
	t.Parallel()

	called := false
	baseProc := &mock.ProcessorStub{
		GetShardIDsCalled: func() []uint32 {
			return []uint32{0, 1, core.MetachainShardId}
		},
		GetObserversCalled: func(shardID uint32, _ data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return []*data.NodeData{
				{
					ShardId: shardID,
					Address: fmt.Sprintf("shard-%d", shardID),
				},
				{
					ShardId: shardID,
					Address: fmt.Sprintf("shard-%d", shardID),
				},
			}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			switch address {
			case "shard-0":
				if !called {
					called = true
					return 400, errors.New("local err")
				}
				valResp := value.(*data.DCDTSupplyResponse)
				valResp.Data.Supply = "300"
				valResp.Data.RecomputedSupply = true
				return 200, nil
			case "shard-1":
				valResp := value.(*data.DCDTSupplyResponse)
				valResp.Data.Supply = "600"
				valResp.Data.Minted = "50"
				valResp.Data.Burned = "100"
				valResp.Data.RecomputedSupply = true
				return 200, nil
			}
			return 0, nil
		},
	}
	scQueryProc := &mock.SCQueryServiceStub{
		ExecuteQueryCalled: func(query *data.SCQuery) (*vm.VMOutputApi, data.BlockInfo, error) {
			return &vm.VMOutputApi{
				ReturnData: [][]byte{nil, nil, nil, []byte("500")},
			}, data.BlockInfo{}, nil
		},
	}
	dcdtProc, err := NewDCDTSupplyProcessor(baseProc, scQueryProc)
	require.Nil(t, err)

	supplyRes, err := dcdtProc.GetDCDTSupply("SEMI-ABCDEF")
	require.Nil(t, err)
	require.Equal(t, "900", supplyRes.Data.Supply)
	require.Equal(t, "0", supplyRes.Data.Burned)
	require.Equal(t, "0", supplyRes.Data.Minted)
	require.True(t, supplyRes.Data.RecomputedSupply)
}
