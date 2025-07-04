package metachain

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain-core/data/rewardTx"
	"github.com/TerraDharitri/drt-go-chain-core/hashing/sha256"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	"github.com/TerraDharitri/drt-go-chain/epochStart"
	"github.com/TerraDharitri/drt-go-chain/epochStart/mock"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/state/factory"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	txExecOrderStub "github.com/TerraDharitri/drt-go-chain/testscommon/common"
	dataRetrieverMock "github.com/TerraDharitri/drt-go-chain/testscommon/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/testscommon/enableEpochsHandlerMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/hashingMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/shardingMocks"
	stateMock "github.com/TerraDharitri/drt-go-chain/testscommon/state"
	"github.com/TerraDharitri/drt-go-chain/testscommon/storage"
	"github.com/TerraDharitri/drt-go-chain/trie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBaseRewardsCreator_NilShardCoordinator(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	args.ShardCoordinator = nil

	rwd, err := NewBaseRewardsCreator(args)

	assert.True(t, check.IfNil(rwd))
	assert.Equal(t, epochStart.ErrNilShardCoordinator, err)
}

func TestBaseRewardsCreator_NilPubkeyConverter(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	args.PubkeyConverter = nil

	rwd, err := NewBaseRewardsCreator(args)

	assert.True(t, check.IfNil(rwd))
	assert.Equal(t, epochStart.ErrNilPubkeyConverter, err)
}

func TestBaseRewardsCreator_NilRewardsStorage(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	args.RewardsStorage = nil

	rwd, err := NewBaseRewardsCreator(args)

	assert.True(t, check.IfNil(rwd))
	assert.Equal(t, epochStart.ErrNilStorage, err)
}

func TestBaseRewardsCreator_NilMiniBlockStorage(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	args.MiniBlockStorage = nil

	rwd, err := NewBaseRewardsCreator(args)

	assert.True(t, check.IfNil(rwd))
	assert.Equal(t, epochStart.ErrNilStorage, err)
}

func TestBaseRewardsCreator_NilHasher(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	args.Hasher = nil

	rwd, err := NewBaseRewardsCreator(args)

	assert.True(t, check.IfNil(rwd))
	assert.Equal(t, epochStart.ErrNilHasher, err)
}

func TestBaseRewardsCreator_NilMarshalizer(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	args.Marshalizer = nil

	rwd, err := NewBaseRewardsCreator(args)

	assert.True(t, check.IfNil(rwd))
	assert.Equal(t, epochStart.ErrNilMarshalizer, err)
}

func TestBaseRewardsCreator_NilDataPoolHolder(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	args.DataPool = nil

	rwd, err := NewBaseRewardsCreator(args)

	assert.True(t, check.IfNil(rwd))
	assert.Equal(t, epochStart.ErrNilDataPoolsHolder, err)
}

func TestBaseRewardsCreator_NilNodesConfigProvider(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	args.NodesConfigProvider = nil

	rwd, err := NewBaseRewardsCreator(args)

	assert.True(t, check.IfNil(rwd))
	assert.Equal(t, epochStart.ErrNilNodesConfigProvider, err)
}

func TestBaseRewardsCreator_NilUserAccountsDB(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	args.UserAccountsDB = nil

	rwd, err := NewBaseRewardsCreator(args)

	assert.True(t, check.IfNil(rwd))
	assert.Equal(t, epochStart.ErrNilAccountsDB, err)
}

func TestBaseRewardsCreator_NilEnableEpochsHandler(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	args.EnableEpochsHandler = nil

	rwd, err := NewBaseRewardsCreator(args)

	assert.True(t, check.IfNil(rwd))
	assert.Equal(t, epochStart.ErrNilEnableEpochsHandler, err)
}

func TestBaseRewardsCreator_NilRewardsHandler(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	args.RewardsHandler = nil

	rwd, err := NewBaseRewardsCreator(args)

	assert.True(t, check.IfNil(rwd))
	assert.Equal(t, epochStart.ErrNilRewardsHandler, err)
}

func TestBaseRewardsCreator_InvalidEnableEpochsHandler(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	args.EnableEpochsHandler = enableEpochsHandlerMock.NewEnableEpochsHandlerStubWithNoFlagsDefined()

	rwd, err := NewBaseRewardsCreator(args)

	assert.True(t, check.IfNil(rwd))
	assert.True(t, errors.Is(err, core.ErrInvalidEnableEpochsHandler))
}

func TestBaseRewardsCreator_clean(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)

	rwd.accumulatedRewards = big.NewInt(1000)
	rwd.protocolSustainabilityValue = big.NewInt(100)
	rwd.mapBaseRewardsPerBlockPerValidator[0] = big.NewInt(10)
	txHash := []byte("txHash")
	rwd.currTxs.AddTx(txHash, &rewardTx.RewardTx{})

	rwd.clean()
	require.Equal(t, big.NewInt(0), rwd.accumulatedRewards)
	require.Equal(t, big.NewInt(0), rwd.protocolSustainabilityValue)
	require.Equal(t, 0, len(rwd.mapBaseRewardsPerBlockPerValidator))
	tx, err := rwd.currTxs.GetTx(txHash)
	require.Nil(t, tx)
	require.NotNil(t, err)
}

func TestBaseRewardsCreator_OkValsShouldWork(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)

	assert.False(t, check.IfNil(rwd))
	assert.Nil(t, err)
}

func TestBaseRewardsCreator_GetLocalTxCache(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)

	txCache := rwd.GetLocalTxCache()
	require.False(t, check.IfNil(txCache))
}

func TestBaseRewardsCreator_GetProtocolSustainabilityRewards(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)

	// should return 0 as just initialized
	rewards := rwd.GetProtocolSustainabilityRewards()
	require.Zero(t, big.NewInt(0).Cmp(rewards))
}

func TestBaseRewardsCreator_addProtocolRewardToMiniblocks(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)

	initialProtRewardValue := big.NewInt(-100)
	protRwAddr, _ := args.PubkeyConverter.Decode(args.RewardsHandler.ProtocolSustainabilityAddressInEpoch(0))
	protRwTx := &rewardTx.RewardTx{
		Round:   100,
		Value:   big.NewInt(0).Set(initialProtRewardValue),
		RcvAddr: protRwAddr,
		Epoch:   1,
	}

	marshalled, err := args.Marshalizer.Marshal(protRwTx)
	require.Nil(t, err)

	protRwTxHash := args.Hasher.Compute(string(marshalled))

	protRwShard := args.ShardCoordinator.ComputeId(protRwAddr)
	mbSlice := createDefaultMiniBlocksSlice()
	err = rwd.addProtocolRewardToMiniBlocks(protRwTx, mbSlice, protRwShard)
	require.Nil(t, err)

	found := false
	for _, mb := range mbSlice {
		for _, txHash := range mb.TxHashes {
			if bytes.Equal(txHash, protRwTxHash) {
				found = true
			}
		}
	}
	require.True(t, found)
}

func TestBaseRewardsCreator_CreateMarshalizedDataNilMiniblocksEmptyMap(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)

	result := rwd.CreateMarshalledData(nil)
	require.Equal(t, 0, len(result))
}

func TestBaseRewardsCreator_CreateMarshalizedDataEmptyMiniblocksEmptyMap(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)

	result := rwd.CreateMarshalledData(&block.Body{})
	require.Equal(t, 0, len(result))
}

func TestBaseRewardsCreator_CreateMarshalizedDataOnlyRewardsMiniblocksGetMarshalized(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)

	dummyMiniBlock := createDummyRewardTxMiniblock(rwd)

	miniBlockTypes := []block.Type{
		block.TxBlock,
		block.StateBlock,
		block.PeerBlock,
		block.SmartContractResultBlock,
		block.InvalidBlock,
		block.ReceiptBlock,
	}

	for _, mbType := range miniBlockTypes {
		dummyMiniBlock.Type = mbType
		result := rwd.CreateMarshalledData(&block.Body{
			MiniBlocks: block.MiniBlockSlice{
				dummyMiniBlock,
			},
		})
		require.Equal(t, 0, len(result))
	}

	dummyMiniBlock.Type = block.RewardsBlock
	result := rwd.CreateMarshalledData(&block.Body{
		MiniBlocks: block.MiniBlockSlice{
			dummyMiniBlock,
		},
	})
	require.Greater(t, len(result), 0)

	readRwTx := &rewardTx.RewardTx{}
	var expectedTx data.TransactionHandler
	for _, resData := range result {
		for _, tx := range resData {
			err = args.Marshalizer.Unmarshal(readRwTx, tx)
			require.Nil(t, err)
			expectedTx, err = rwd.currTxs.GetTx(dummyMiniBlock.TxHashes[0])
			require.Nil(t, err)
			require.Equal(t, expectedTx, readRwTx)
		}
	}
}

func TestBaseRewardsCreator_CreateMarshalizedDataWrongSenderNotIncluded(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)

	dummyMiniBlock := createDummyRewardTxMiniblock(rwd)
	dummyMiniBlock.Type = block.RewardsBlock
	dummyMiniBlock.SenderShardID = args.ShardCoordinator.SelfId() + 1
	result := rwd.CreateMarshalledData(&block.Body{
		MiniBlocks: block.MiniBlockSlice{
			dummyMiniBlock,
		},
	})
	require.Equal(t, 0, len(result))
}

func TestBaseRewardsCreator_CreateMarshalizedDataNotFoundTxHashIgnored(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)

	dummyMiniBlock := createDummyRewardTxMiniblock(rwd)
	dummyMiniBlock.Type = block.RewardsBlock
	dummyMiniBlock.TxHashes = [][]byte{[]byte("not found txHash")}
	result := rwd.CreateMarshalledData(&block.Body{
		MiniBlocks: block.MiniBlockSlice{
			dummyMiniBlock,
		},
	})
	require.Equal(t, 0, len(result))
}

func TestBaseRewardsCreator_GetRewardsTxsNonRewardsMiniBlocksGetIgnored(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)

	dummyMb := createDummyRewardTxMiniblock(rwd)
	dummyMb.Type = block.StateBlock
	result := rwd.GetRewardsTxs(&block.Body{MiniBlocks: block.MiniBlockSlice{dummyMb}})
	require.Equal(t, 0, len(result))
}

func TestBaseRewardsCreator_GetRewardsTxsRewardMiniBlockNotFoundTxIgnored(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)
	rwTxHash := []byte("not found")
	dummyMb := createDummyRewardTxMiniblock(rwd)
	dummyMb.TxHashes = [][]byte{rwTxHash}

	result := rwd.GetRewardsTxs(&block.Body{MiniBlocks: block.MiniBlockSlice{dummyMb}})
	require.Equal(t, 0, len(result))
}

func TestBaseRewardsCreator_GetRewardsTxsRewardMiniBlockOK(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)

	rwTx := &rewardTx.RewardTx{
		Value: big.NewInt(1000),
	}
	rwTxHash := []byte("rwTxHash")
	rwd.currTxs.AddTx(rwTxHash, rwTx)

	dummyMb := createDummyRewardTxMiniblock(rwd)
	dummyMb.TxHashes = [][]byte{rwTxHash}

	result := rwd.GetRewardsTxs(&block.Body{MiniBlocks: block.MiniBlockSlice{dummyMb}})
	require.Equal(t, 1, len(result))
	for _, tx := range result {
		require.Equal(t, rwTx, tx)
	}
}

func TestBaseRewardsCreator_SaveTxBlockToStorageNilBodyNoPanic(t *testing.T) {
	t.Parallel()

	defer func() {
		r := recover()
		require.Nil(t, r)
	}()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)

	rwd.SaveBlockDataToStorage(nil, nil)
}

func TestBaseRewardsCreator_SaveTxBlockToStorageNonRewardsMiniBlocksAreIgnored(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)

	dummyMiniBlock := createDummyRewardTxMiniblock(rwd)

	miniBlockTypes := []block.Type{
		block.TxBlock,
		block.StateBlock,
		block.PeerBlock,
		block.SmartContractResultBlock,
		block.InvalidBlock,
		block.ReceiptBlock,
	}

	var mb, mmb []byte
	for _, mbType := range miniBlockTypes {
		dummyMiniBlock.Type = mbType

		rwd.SaveBlockDataToStorage(nil, &block.Body{
			MiniBlocks: block.MiniBlockSlice{
				dummyMiniBlock,
			},
		})

		mmb, err = args.Marshalizer.Marshal(dummyMiniBlock)
		require.Nil(t, err)
		mbHash := args.Hasher.Compute(string(mmb))
		mb, err = args.MiniBlockStorage.Get(mbHash)
		require.Nil(t, mb)
		require.NotNil(t, err)
	}

	dummyMiniBlock.Type = block.RewardsBlock
	rwd.SaveBlockDataToStorage(nil, &block.Body{
		MiniBlocks: block.MiniBlockSlice{
			dummyMiniBlock,
		},
	})

	mmb, err = args.Marshalizer.Marshal(dummyMiniBlock)
	require.Nil(t, err)
	mbHash := args.Hasher.Compute(string(mmb))
	mb, err = rwd.miniBlockStorage.Get(mbHash)
	require.Equal(t, mmb, mb)
	require.Nil(t, err)

	rwTx, err := rwd.rewardsStorage.Get(dummyMiniBlock.TxHashes[0])
	require.NotNil(t, rwTx)
	require.Nil(t, err)
}

func TestBaseRewardsCreator_SaveTxBlockToStorageNotFoundTxIgnored(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)
	rwTxHash := []byte("not found")
	dummyMb := createDummyRewardTxMiniblock(rwd)
	dummyMb.TxHashes = [][]byte{rwTxHash}

	rwd.SaveBlockDataToStorage(nil, &block.Body{MiniBlocks: block.MiniBlockSlice{dummyMb}})

	mmb, err := args.Marshalizer.Marshal(dummyMb)
	require.Nil(t, err)
	mbHash := args.Hasher.Compute(string(mmb))
	mb, err := rwd.miniBlockStorage.Get(mbHash)
	require.Equal(t, mmb, mb)
	require.Nil(t, err)

	rwTx, err := rwd.rewardsStorage.Get(rwTxHash)
	require.Nil(t, rwTx)
	require.NotNil(t, err)
}

func TestBaseRewardsCreator_DeleteTxsFromStorageNilMetablockNoPanic(t *testing.T) {
	t.Parallel()

	defer func() {
		r := recover()
		require.Nil(t, r)
	}()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)

	dummyMb := createDummyRewardTxMiniblock(rwd)
	rwd.DeleteBlockDataFromStorage(nil, &block.Body{MiniBlocks: block.MiniBlockSlice{dummyMb}})
}

func TestBaseRewardsCreator_DeleteTxsFromStorageNilBlockBodyNoPanic(t *testing.T) {
	t.Parallel()

	defer func() {
		r := recover()
		require.Nil(t, r)
	}()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)

	metaBlk := &block.MetaBlock{
		EpochStart:     getDefaultEpochStart(),
		DevFeesInEpoch: big.NewInt(0),
	}

	rwd.DeleteBlockDataFromStorage(metaBlk, nil)
}

func TestBaseRewardsCreator_DeleteTxsFromStorageNonRewardsMiniBlocksIgnored(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)

	miniBlockTypes := []block.Type{
		block.TxBlock,
		block.StateBlock,
		block.PeerBlock,
		block.SmartContractResultBlock,
		block.InvalidBlock,
		block.ReceiptBlock,
	}
	metaBlk := &block.MetaBlock{
		EpochStart:     getDefaultEpochStart(),
		DevFeesInEpoch: big.NewInt(0),
	}

	var tx, mb []byte
	for _, mbType := range miniBlockTypes {
		dummyMb := createDummyRewardTxMiniblock(rwd)
		dummyMb.Type = mbType
		rwTxHash := []byte("txHash")
		marshalledRwTx, _ := args.Marshalizer.Marshal(&rewardTx.RewardTx{})
		dummyMb.TxHashes = [][]byte{rwTxHash}

		_ = rwd.rewardsStorage.Put(rwTxHash, marshalledRwTx)

		mbHash := []byte("mb1")
		metaBlk.MiniBlockHeaders = []block.MiniBlockHeader{
			{
				Hash: mbHash,
				Type: mbType,
			},
		}
		dummyMbMarshalled, _ := args.Marshalizer.Marshal(dummyMb)
		_ = rwd.miniBlockStorage.Put(mbHash, dummyMbMarshalled)

		rwd.DeleteBlockDataFromStorage(metaBlk, &block.Body{MiniBlocks: block.MiniBlockSlice{dummyMb}})
		tx, err = rwd.rewardsStorage.Get(rwTxHash)
		require.Nil(t, err)
		require.NotNil(t, tx)

		mb, err = rwd.miniBlockStorage.Get(mbHash)
		require.Nil(t, err)
		require.NotNil(t, mb)
	}
}

func TestBaseRewardsCreator_DeleteTxsFromStorage(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)

	metaBlk := &block.MetaBlock{
		EpochStart:     getDefaultEpochStart(),
		DevFeesInEpoch: big.NewInt(0),
	}

	dummyMb := createDummyRewardTxMiniblock(rwd)
	dummyMb.Type = block.RewardsBlock
	rwTxHash := []byte("txHash")
	dummyMb.TxHashes = [][]byte{rwTxHash}
	marshalledRwTx, _ := args.Marshalizer.Marshal(&rewardTx.RewardTx{})

	_ = rwd.rewardsStorage.Put(rwTxHash, marshalledRwTx)

	mbHash := []byte("mb1")
	metaBlk.MiniBlockHeaders = []block.MiniBlockHeader{
		{
			Hash: mbHash,
			Type: block.RewardsBlock,
		},
	}
	dummyMbMarshalled, _ := args.Marshalizer.Marshal(dummyMb)
	_ = rwd.miniBlockStorage.Put(mbHash, dummyMbMarshalled)

	rwd.DeleteBlockDataFromStorage(metaBlk, &block.Body{MiniBlocks: block.MiniBlockSlice{dummyMb}})
	tx, err := rwd.rewardsStorage.Get(rwTxHash)
	require.NotNil(t, err)
	require.Nil(t, tx)

	mb, err := rwd.miniBlockStorage.Get(mbHash)
	require.NotNil(t, err)
	require.Nil(t, mb)
}

func TestBaseRewardsCreator_RemoveBlockDataFromPoolsNilMetablockNoPanic(t *testing.T) {
	t.Parallel()

	defer func() {
		r := recover()
		require.Nil(t, r)
	}()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)

	dummyMb := createDummyRewardTxMiniblock(rwd)
	rwd.RemoveBlockDataFromPools(nil, &block.Body{MiniBlocks: block.MiniBlockSlice{dummyMb}})
}

func TestBaseRewardsCreator_RemoveBlockDataFromPoolsNilBlockBodyNoPanic(t *testing.T) {
	t.Parallel()

	defer func() {
		r := recover()
		require.Nil(t, r)
	}()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)

	metaBlk := &block.MetaBlock{
		EpochStart:     getDefaultEpochStart(),
		DevFeesInEpoch: big.NewInt(0),
	}

	rwd.DeleteBlockDataFromStorage(metaBlk, nil)
}

func TestBaseRewardsCreator_RemoveBlockDataFromPoolsNonRewardsMiniBlocksIgnored(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)

	miniBlockTypes := []block.Type{
		block.TxBlock,
		block.StateBlock,
		block.PeerBlock,
		block.SmartContractResultBlock,
		block.InvalidBlock,
		block.ReceiptBlock,
	}
	metaBlk := &block.MetaBlock{
		EpochStart:     getDefaultEpochStart(),
		DevFeesInEpoch: big.NewInt(0),
	}
	for _, mbType := range miniBlockTypes {
		dummyMb := createDummyRewardTxMiniblock(rwd)
		dummyMb.Type = mbType
		rwTxHash := []byte("txHash")
		marshalledRwTx, _ := args.Marshalizer.Marshal(&rewardTx.RewardTx{})
		dummyMb.TxHashes = [][]byte{rwTxHash}

		_ = rwd.rewardsStorage.Put(rwTxHash, marshalledRwTx)

		mbHash := []byte("mb1")
		metaBlk.MiniBlockHeaders = []block.MiniBlockHeader{
			{
				Hash: mbHash,
				Type: mbType,
			},
		}
		dummyMbMarshalled, _ := args.Marshalizer.Marshal(dummyMb)
		strCache := process.ShardCacherIdentifier(dummyMb.SenderShardID, dummyMb.ReceiverShardID)
		_ = rwd.dataPool.MiniBlocks().Put(mbHash, dummyMbMarshalled, len(dummyMbMarshalled))
		rwd.dataPool.RewardTransactions().AddData(rwTxHash, rwd, 100, strCache)

		rwd.RemoveBlockDataFromPools(metaBlk, &block.Body{MiniBlocks: block.MiniBlockSlice{dummyMb}})
		// non reward txs do not get removed
		tx, ok := rwd.dataPool.RewardTransactions().ShardDataStore(strCache).Get(rwTxHash)
		require.True(t, ok)
		require.NotNil(t, tx)

		// non reward miniBlocks do not get removed
		mb, ok := rwd.dataPool.MiniBlocks().Get(mbHash)
		require.True(t, ok)
		require.NotNil(t, mb)
	}
}

func TestBaseRewardsCreator_RemoveBlockDataFromPools(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)

	metaBlk := &block.MetaBlock{
		EpochStart:     getDefaultEpochStart(),
		DevFeesInEpoch: big.NewInt(0),
	}

	mbType := block.RewardsBlock
	dummyMb := createDummyRewardTxMiniblock(rwd)
	dummyMb.Type = mbType
	rwTxHash := []byte("txHash")
	marshalledRwTx, _ := args.Marshalizer.Marshal(&rewardTx.RewardTx{})
	dummyMb.TxHashes = [][]byte{rwTxHash}

	_ = rwd.rewardsStorage.Put(rwTxHash, marshalledRwTx)

	mbHash := []byte("mb1")
	metaBlk.MiniBlockHeaders = []block.MiniBlockHeader{
		{
			Hash: mbHash,
			Type: mbType,
		},
	}
	dummyMbMarshalled, _ := args.Marshalizer.Marshal(dummyMb)
	strCache := process.ShardCacherIdentifier(dummyMb.SenderShardID, dummyMb.ReceiverShardID)
	_ = rwd.dataPool.MiniBlocks().Put(mbHash, dummyMbMarshalled, len(dummyMbMarshalled))
	rwd.dataPool.Transactions().AddData(rwTxHash, rwd, 100, strCache)

	rwd.RemoveBlockDataFromPools(metaBlk, &block.Body{MiniBlocks: block.MiniBlockSlice{dummyMb}})
	// reward txs get removed
	tx, ok := rwd.dataPool.Transactions().ShardDataStore(strCache).Get(rwTxHash)
	require.False(t, ok)
	require.Nil(t, tx)

	// reward miniBlocks get removed
	mb, ok := rwd.dataPool.MiniBlocks().Get(mbHash)
	require.False(t, ok)
	require.Nil(t, mb)
}

func TestBaseRewardsCreator_isSystemDelegationSC(t *testing.T) {
	t.Parallel()

	nonExistentAccountAddress := []byte("address")
	peerAccountAddress := []byte("addressPeer")
	userAccountAddress := []byte("addressUser")

	args := getBaseRewardsArguments()
	args.UserAccountsDB = &stateMock.AccountsStub{
		GetExistingAccountCalled: func(addressContainer []byte) (vmcommon.AccountHandler, error) {
			if bytes.Equal(addressContainer, nonExistentAccountAddress) {
				return nil, fmt.Errorf("account does not exist")
			}

			if bytes.Equal(addressContainer, peerAccountAddress) {
				peerAccount := &stateMock.PeerAccountHandlerMock{
					AddressBytesCalled: func() []byte {
						return peerAccountAddress
					},
				}

				return peerAccount, nil
			}

			if bytes.Equal(addressContainer, userAccountAddress) {
				userAccount := &stateMock.UserAccountStub{
					RetrieveValueCalled: func(key []byte) ([]byte, uint32, error) {
						if bytes.Equal(key, []byte(core.DelegationSystemSCKey)) {
							return []byte("delegation"), 0, nil
						}
						return nil, 0, fmt.Errorf("not found")
					},
				}

				return userAccount, nil
			}

			return &stateMock.UserAccountStub{}, nil
		},
	}
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)

	// not existing account
	isDelegationSCAddress := rwd.isSystemDelegationSC(nonExistentAccountAddress)
	require.False(t, isDelegationSCAddress)

	// peer account
	isDelegationSCAddress = rwd.isSystemDelegationSC(peerAccountAddress)
	require.False(t, isDelegationSCAddress)

	// existing user account
	isDelegationSCAddress = rwd.isSystemDelegationSC(userAccountAddress)
	require.True(t, isDelegationSCAddress)
}

func TestBaseRewardsCreator_isSystemDelegationSCTrue(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	args.UserAccountsDB = &stateMock.AccountsStub{
		GetExistingAccountCalled: func(address []byte) (vmcommon.AccountHandler, error) {
			return &stateMock.UserAccountStub{
				RetrieveValueCalled: func(key []byte) ([]byte, uint32, error) {
					if bytes.Equal(key, []byte("delegation")) {
						return []byte("value"), 0, nil
					}

					return nil, 0, fmt.Errorf("error")
				},
			}, nil
		},
	}
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)

	isDelegationSCAddress := rwd.isSystemDelegationSC([]byte("userAddress"))

	require.True(t, isDelegationSCAddress)
}

func TestBaseRewardsCreator_createProtocolSustainabilityRewardTransaction(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)

	metaBlk := &block.MetaBlock{
		EpochStart:     getDefaultEpochStart(),
		DevFeesInEpoch: big.NewInt(0),
	}

	rwTx, _, err := rwd.createProtocolSustainabilityRewardTransaction(metaBlk, &metaBlk.EpochStart.Economics)
	require.Nil(t, err)
	require.NotNil(t, rwTx)
	require.Equal(t, metaBlk.EpochStart.Economics.RewardsForProtocolSustainability, rwTx.Value)
}

func TestBaseRewardsCreator_createRewardFromRwdInfo(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)

	metaBlk := &block.MetaBlock{
		EpochStart:     getDefaultEpochStart(),
		DevFeesInEpoch: big.NewInt(0),
	}

	rwInfo := &rewardInfoData{
		accumulatedFees:     big.NewInt(100),
		address:             "addressRewards",
		rewardsFromProtocol: big.NewInt(1000),
	}

	rwTx, rwTxHash, err := rwd.createRewardFromRwdInfo(rwInfo, metaBlk)
	require.Nil(t, err)
	require.NotNil(t, rwTx)
	require.NotNil(t, rwTxHash)
	require.Equal(t, big.NewInt(0).Add(rwInfo.accumulatedFees, rwInfo.rewardsFromProtocol), rwTx.Value)
}

func TestBaseRewardsCreator_initializeRewardsMiniBlocks(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)

	mbSlice := rwd.initializeRewardsMiniBlocks()
	require.NotNil(t, mbSlice)
	require.Equal(t, int(args.ShardCoordinator.NumberOfShards()+1), len(mbSlice))
	for _, mb := range mbSlice {
		require.Equal(t, block.RewardsBlock, mb.Type)
		require.Equal(t, 0, len(mb.TxHashes))
	}
}

func TestBaseRewardsCreator_finalizeMiniBlocksOrdersTxsAscending(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)

	orderedTxHashes := [][]byte{[]byte("txHash0"), []byte("txHash1"), []byte("txHash2")}

	mbSlice := createDefaultMiniBlocksSlice()
	for _, mb := range mbSlice {
		nbTxs := len(orderedTxHashes)
		mb.TxHashes = make([][]byte, nbTxs)
		for i := range orderedTxHashes {
			// in descending order
			mb.TxHashes[i] = orderedTxHashes[nbTxs-i-1]
		}
	}

	resultedMbs := rwd.finalizeMiniBlocks(mbSlice)
	for _, mb := range resultedMbs {
		require.Equal(t, orderedTxHashes, mb.TxHashes)
	}
}

func TestBaseRewardsCreator_finalizeMiniBlocksEmptyMbsAreRemoved(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)

	orderedTxHashes := [][]byte{[]byte("txHash0"), []byte("txHash1"), []byte("txHash2")}

	mbSlice := createDefaultMiniBlocksSlice()
	for i := 0; i < len(mbSlice)-1; i++ {
		nbTxs := len(orderedTxHashes)
		mbSlice[i].TxHashes = make([][]byte, nbTxs)
		for j := range orderedTxHashes {
			// in descending order
			mbSlice[i].TxHashes[j] = orderedTxHashes[nbTxs-j-1]
		}
	}

	resultedMbs := rwd.finalizeMiniBlocks(mbSlice)

	// mb without txs is removed
	require.Equal(t, len(mbSlice)-1, len(resultedMbs))
}

func TestBaseRewardsCreator_fillBaseRewardsPerBlockPerNode(t *testing.T) {
	t.Parallel()

	// should work for epoch 0 even if this is a bad input
	testFillBaseRewardsPerBlockPerNode(t, 0)

	// should work for an epoch higher than 0
	testFillBaseRewardsPerBlockPerNode(t, 1)
}

func testFillBaseRewardsPerBlockPerNode(t *testing.T, epoch uint32) {
	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)

	baseRewardsPerNode := big.NewInt(1000000)
	rwd.fillBaseRewardsPerBlockPerNode(baseRewardsPerNode, epoch)
	consensusShard := args.NodesConfigProvider.ConsensusGroupSizeForShardAndEpoch(0, epoch)
	consensusMeta := args.NodesConfigProvider.ConsensusGroupSizeForShardAndEpoch(core.MetachainShardId, epoch)
	expectedRewardPerNodeInShard := big.NewInt(0).Div(baseRewardsPerNode, big.NewInt(int64(consensusShard)))
	expectedRewardPerNodeInMeta := big.NewInt(0).Div(baseRewardsPerNode, big.NewInt(int64(consensusMeta)))

	for shardID, rewardPerNode := range rwd.mapBaseRewardsPerBlockPerValidator {
		if shardID == core.MetachainShardId {
			require.Equal(t, expectedRewardPerNodeInMeta, rewardPerNode)
			continue
		}
		require.Equal(t, expectedRewardPerNodeInShard, rewardPerNode)
	}
}

func TestBaseRewardsCreator_verifyCreatedRewardMiniBlocksWithMetaBlockNonRewardsMbsHeadersAreIgnored(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)

	miniBlockTypes := []block.Type{
		block.TxBlock,
		block.StateBlock,
		block.PeerBlock,
		block.SmartContractResultBlock,
		block.InvalidBlock,
		block.ReceiptBlock,
	}

	metaBlk := &block.MetaBlock{
		EpochStart:       getDefaultEpochStart(),
		DevFeesInEpoch:   big.NewInt(0),
		MiniBlockHeaders: make([]block.MiniBlockHeader, 4),
	}
	for _, mbType := range miniBlockTypes {
		for _, mbHeader := range metaBlk.MiniBlockHeaders {
			mbHeader.Type = mbType
		}
		err = rwd.verifyCreatedRewardMiniBlocksWithMetaBlock(metaBlk, nil)
		require.Nil(t, err)
	}
}

func TestBaseRewardsCreator_verifyCreatedRewardMiniBlocksWithMetaBlockMiniBlockHashMismatchReceiverShardMB(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)

	metaBlk := &block.MetaBlock{
		EpochStart:       getDefaultEpochStart(),
		DevFeesInEpoch:   big.NewInt(0),
		MiniBlockHeaders: make([]block.MiniBlockHeader, 4),
	}

	mbs := createDefaultMiniBlocksSlice()
	mbs[0].ReceiverShardID = core.MetachainShardId
	mbs[1].ReceiverShardID = core.MetachainShardId

	for i := range metaBlk.MiniBlockHeaders {
		metaBlk.MiniBlockHeaders[i].Type = block.RewardsBlock
	}
	err = rwd.verifyCreatedRewardMiniBlocksWithMetaBlock(metaBlk, mbs)
	require.Equal(t, epochStart.ErrRewardMiniBlockHashDoesNotMatch, err)
}

func TestBaseRewardsCreator_verifyCreatedRewardMiniBlocksWithMetaBlockMiniBlockHashMismatch(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)

	metaBlk := &block.MetaBlock{
		EpochStart:       getDefaultEpochStart(),
		DevFeesInEpoch:   big.NewInt(0),
		MiniBlockHeaders: make([]block.MiniBlockHeader, 4),
	}

	mbs := createDefaultMiniBlocksSlice()
	mbs[0].TxHashes = [][]byte{[]byte("txHash")}

	for i := range metaBlk.MiniBlockHeaders {
		metaBlk.MiniBlockHeaders[i].Type = block.RewardsBlock
	}
	err = rwd.verifyCreatedRewardMiniBlocksWithMetaBlock(metaBlk, mbs)
	require.Equal(t, epochStart.ErrRewardMiniBlockHashDoesNotMatch, err)
}

func TestBaseRewardsCreator_verifyCreatedRewardMiniBlocksWithMetaBlockMiniBlockNumMismatch(t *testing.T) {
	t.Parallel()

	args := getBaseRewardsArguments()
	rwd, err := NewBaseRewardsCreator(args)
	require.Nil(t, err)
	require.NotNil(t, rwd)

	metaBlk := &block.MetaBlock{
		EpochStart:       getDefaultEpochStart(),
		DevFeesInEpoch:   big.NewInt(0),
		MiniBlockHeaders: make([]block.MiniBlockHeader, 4),
	}

	mbs := createDefaultMiniBlocksSlice()

	for i := range metaBlk.MiniBlockHeaders {
		metaBlk.MiniBlockHeaders[i].Type = block.TxBlock
	}
	err = rwd.verifyCreatedRewardMiniBlocksWithMetaBlock(metaBlk, mbs)
	require.Equal(t, epochStart.ErrRewardMiniBlocksNumDoesNotMatch, err)
}

func TestBaseRewardsCreator_getMiniBlockWithReceiverShardIDNotFound(t *testing.T) {
	mbSlice := createDefaultMiniBlocksSlice()
	for i := range mbSlice {
		mbSlice[i].ReceiverShardID = 0
	}

	mb := getMiniBlockWithReceiverShardID(core.MetachainShardId, mbSlice)
	require.Nil(t, mb)
}

func TestBaseRewardsCreator_getMiniBlockWithReceiverShardIDFound(t *testing.T) {
	mbSlice := createDefaultMiniBlocksSlice()
	mbSlice[0].ReceiverShardID = 0
	mb := getMiniBlockWithReceiverShardID(0, mbSlice)
	require.Equal(t, mbSlice[0], mb)
}

func getBaseRewardsArguments() BaseRewardsCreatorArgs {
	hasher := sha256.NewSha256()
	marshalizer := &marshal.GogoProtoMarshalizer{}

	storageManagerArgs := storage.GetStorageManagerArgs()
	storageManagerArgs.Marshalizer = marshalizer
	storageManagerArgs.Hasher = hasher

	trieFactoryManager, _ := trie.CreateTrieStorageManager(storageManagerArgs, storage.GetStorageManagerOptions())
	argsAccCreator := factory.ArgsAccountCreator{
		Hasher:              hasher,
		Marshaller:          marshalizer,
		EnableEpochsHandler: &enableEpochsHandlerMock.EnableEpochsHandlerStub{},
	}
	accCreator, _ := factory.NewAccountCreator(argsAccCreator)
	enableEpochsHandler := &enableEpochsHandlerMock.EnableEpochsHandlerStub{}
	userAccountsDB := createAccountsDB(hasher, marshalizer, accCreator, trieFactoryManager, enableEpochsHandler)
	shardCoordinator := mock.NewMultiShardsCoordinatorMock(2)
	shardCoordinator.CurrentShard = core.MetachainShardId
	shardCoordinator.ComputeIdCalled = func(address []byte) uint32 {
		return 0
	}

	rewardsTopUpGradientPoint, _ := big.NewInt(0).SetString("3000000000000000000000000", 10)

	rewardsHandler := &mock.RewardsHandlerStub{
		RewardsTopUpGradientPointInEpochCalled: func(_ uint32) *big.Int {
			return rewardsTopUpGradientPoint
		},
		RewardsTopUpFactorInEpochCalled: func(_ uint32) float64 {
			return 0.25
		},
		ProtocolSustainabilityAddressInEpochCalled: func(_ uint32) string {
			return "11"
		},
	}

	return BaseRewardsCreatorArgs{
		ShardCoordinator: shardCoordinator,
		PubkeyConverter:  testscommon.NewPubkeyConverterMock(32),
		RewardsStorage:   mock.NewStorerMock(),
		MiniBlockStorage: mock.NewStorerMock(),
		Hasher:           &hashingMocks.HasherMock{},
		Marshalizer:      &mock.MarshalizerMock{},
		DataPool:         dataRetrieverMock.NewPoolsHolderMock(),
		NodesConfigProvider: &shardingMocks.NodesCoordinatorStub{
			ConsensusGroupSizeCalled: func(shardID uint32, _ uint32) int {
				if shardID == core.MetachainShardId {
					return 400
				}
				return 63
			},
		},
		UserAccountsDB:        userAccountsDB,
		EnableEpochsHandler:   enableEpochsHandler,
		ExecutionOrderHandler: &txExecOrderStub.TxExecutionOrderHandlerStub{},
		RewardsHandler:        rewardsHandler,
	}
}

func createDummyRewardTxMiniblock(rwd *baseRewardsCreator) *block.MiniBlock {
	dummyTx := &rewardTx.RewardTx{}
	dummyTxHash := []byte("rwdTxHash")
	rwd.currTxs.AddTx(dummyTxHash, dummyTx)

	return &block.MiniBlock{
		Type:            block.RewardsBlock,
		SenderShardID:   rwd.shardCoordinator.SelfId(),
		ReceiverShardID: rwd.shardCoordinator.SelfId() + 1,
		TxHashes:        [][]byte{dummyTxHash},
	}
}

func createDefaultMiniBlocksSlice() block.MiniBlockSlice {
	return block.MiniBlockSlice{
		&block.MiniBlock{
			TxHashes:        nil,
			ReceiverShardID: 0,
			SenderShardID:   core.MetachainShardId,
			Type:            block.RewardsBlock,
		},
		&block.MiniBlock{
			TxHashes:        nil,
			ReceiverShardID: 1,
			SenderShardID:   core.MetachainShardId,
			Type:            block.RewardsBlock,
		},
		&block.MiniBlock{
			TxHashes:        nil,
			ReceiverShardID: 2,
			SenderShardID:   core.MetachainShardId,
			Type:            block.RewardsBlock,
		},
	}
}
