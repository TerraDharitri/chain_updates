package endOfEpoch

import (
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/stretchr/testify/assert"

	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/integrationTests"
)

// CreateAndPropagateBlocks -
func CreateAndPropagateBlocks(
	t *testing.T,
	nbRounds uint64,
	currentRound uint64,
	currentNonce uint64,
	nodes []*integrationTests.TestProcessorNode,
	leaders []*integrationTests.TestProcessorNode,
) (uint64, uint64) {
	for i := uint64(0); i <= nbRounds; i++ {
		integrationTests.UpdateRound(nodes, currentRound)
		integrationTests.ProposeBlock(nodes, leaders, currentRound, currentNonce)
		integrationTests.SyncBlock(t, nodes, leaders, currentRound)
		currentRound = integrationTests.IncrementAndPrintRound(currentRound)
		currentNonce++
	}
	time.Sleep(time.Second)

	return currentRound, currentNonce
}

// VerifyThatNodesHaveCorrectEpoch -
func VerifyThatNodesHaveCorrectEpoch(
	t *testing.T,
	epoch uint32,
	nodes []*integrationTests.TestProcessorNode,
) {
	for _, node := range nodes {
		currentShId := node.ShardCoordinator.SelfId()
		currentHeader := node.BlockChain.GetCurrentBlockHeader()
		assert.Equal(t, epoch, currentHeader.GetEpoch())

		for _, testNode := range nodes {
			if testNode.ShardCoordinator.SelfId() == currentShId {
				testHeader := testNode.BlockChain.GetCurrentBlockHeader()
				assert.Equal(t, testHeader.GetNonce(), currentHeader.GetNonce())
			}
		}
	}
}

// VerifyIfAddedShardHeadersAreWithNewEpoch -
func VerifyIfAddedShardHeadersAreWithNewEpoch(
	t *testing.T,
	nodes []*integrationTests.TestProcessorNode,
) {
	for _, node := range nodes {
		if node.ShardCoordinator.SelfId() != core.MetachainShardId {
			continue
		}

		currentMetaHdr, ok := node.BlockChain.GetCurrentBlockHeader().(*block.MetaBlock)
		if !ok {
			assert.Fail(t, "metablock should have been in current block header")
		}

		shardHDrStorage, err := node.Storage.GetStorer(dataRetriever.BlockHeaderUnit)
		assert.Nil(t, err)
		for _, shardInfo := range currentMetaHdr.ShardInfo {
			header, err := node.DataPool.Headers().GetHeaderByHash(shardInfo.HeaderHash)
			if err == nil {
				assert.Equal(t, header.GetEpoch(), currentMetaHdr.GetEpoch())
				continue
			}

			buff, err := shardHDrStorage.Get(shardInfo.HeaderHash)
			assert.Nil(t, err)

			shardHeader := block.Header{}
			err = integrationTests.TestMarshalizer.Unmarshal(&shardHeader, buff)
			assert.Nil(t, err)
			assert.Equal(t, shardHeader.Epoch, currentMetaHdr.Epoch)
		}
	}
}

// GetBlockProposersIndexes -
func GetBlockProposersIndexes(
	consensusMap map[uint32][]*integrationTests.TestProcessorNode,
	nodesMap map[uint32][]*integrationTests.TestProcessorNode,
) map[uint32]int {

	indexProposer := make(map[uint32]int)

	for sh, testNodeList := range nodesMap {
		for k, testNode := range testNodeList {
			if consensusMap[sh][0] == testNode {
				indexProposer[sh] = k
			}
		}
	}

	return indexProposer
}
