package executingMiniblocks

import (
	"bytes"
	"crypto"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	logger "github.com/TerraDharitri/drt-go-chain-logger"
	"github.com/stretchr/testify/assert"

	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/integrationTests"
	testBlock "github.com/TerraDharitri/drt-go-chain/integrationTests/singleShard/block"
	"github.com/TerraDharitri/drt-go-chain/process"
)

// TestShardShouldNotProposeAndExecuteTwoBlocksInSameRound tests that a shard can not continue building on a
// chain with 2 blocks in the same round
func TestShardShouldNotProposeAndExecuteTwoBlocksInSameRound(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	maxShards := uint32(1)
	numOfNodes := 4

	nodes := make([]*integrationTests.TestProcessorNode, numOfNodes)
	connectableNodes := make([]integrationTests.Connectable, numOfNodes)
	for i := 0; i < numOfNodes; i++ {
		nodes[i] = integrationTests.NewTestProcessorNode(integrationTests.ArgTestProcessorNode{
			MaxShards:            maxShards,
			NodeShardId:          0,
			TxSignPrivKeyShardId: 0,
		})
		connectableNodes[i] = nodes[i]
	}

	integrationTests.ConnectNodes(connectableNodes)

	idxProposer := 0
	leader := nodes[idxProposer]

	defer func() {
		for _, n := range nodes {
			n.Close()
		}
	}()

	fmt.Println("Delaying for nodes p2p bootstrap...")
	time.Sleep(integrationTests.P2pBootstrapDelay)

	round := uint64(0)
	nonce := uint64(1)
	round = integrationTests.IncrementAndPrintRound(round)

	err := proposeAndCommitBlock(leader, round, nonce)
	assert.Nil(t, err)

	integrationTests.SyncBlock(t, nodes, []*integrationTests.TestProcessorNode{leader}, nonce)

	time.Sleep(testBlock.StepDelay)

	checkCurrentBlockHeight(t, nodes, nonce)

	// only nonce increases, round stays the same
	nonce++

	err = proposeAndCommitBlock(nodes[idxProposer], round, nonce)
	assert.Equal(t, process.ErrLowerRoundInBlock, err)

	// mockTestingT is used as in normal case SyncBlock would fail as it doesn't find the header with nonce 2
	mockTestingT := &testing.T{}
	integrationTests.SyncBlock(mockTestingT, nodes, []*integrationTests.TestProcessorNode{leader}, nonce)

	time.Sleep(testBlock.StepDelay)

	checkCurrentBlockHeight(t, nodes, nonce-1)
}

// TestShardShouldProposeBlockContainingInvalidTransactions tests the following scenario:
//  1. generate 3 move balance transactions: one that can be executed, one to be processed as invalid, and one that isn't executable (no balance left for fee).
//  2. proposer will have those 3 transactions in its pools and will propose a block
//  3. another node will be able to sync the proposed block (and request - receive) the 2 transactions that
//     will end up in the block (one valid and one invalid)
//  4. the non-executable transaction will not be immediately removed from the proposer's pool. See DRT-16200.
func TestShardShouldProposeBlockContainingInvalidTransactions(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	maxShards := uint32(1)
	numOfNodes := 2

	nodes := make([]*integrationTests.TestProcessorNode, numOfNodes)
	connectableNodes := make([]integrationTests.Connectable, numOfNodes)
	for i := 0; i < numOfNodes; i++ {
		nodes[i] = integrationTests.NewTestProcessorNode(integrationTests.ArgTestProcessorNode{
			MaxShards:            maxShards,
			NodeShardId:          0,
			TxSignPrivKeyShardId: 0,
		})
		connectableNodes[i] = nodes[i]
	}

	integrationTests.ConnectNodes(connectableNodes)

	idxProposer := 0
	leader := nodes[idxProposer]

	defer func() {
		for _, n := range nodes {
			n.Close()
		}
	}()

	fmt.Println("Delaying for nodes p2p bootstrap...")
	time.Sleep(integrationTests.P2pBootstrapDelay)

	round := uint64(0)
	nonce := uint64(1)
	round = integrationTests.IncrementAndPrintRound(round)

	transferValue := uint64(1000000)
	mintAllNodes(nodes, transferValue)

	txs, hashes := generateTransferTxs(transferValue, leader.OwnAccount.SkTxSign, nodes[1].OwnAccount.PkTxSign)
	addTxsInDataPool(leader, txs, hashes)

	_, _ = integrationTests.ProposeAndSyncOneBlock(t, nodes, []*integrationTests.TestProcessorNode{leader}, round, nonce)

	fmt.Println(integrationTests.MakeDisplayTable(nodes))

	testStateOnNodes(t, nodes, idxProposer, hashes)
}

func mintAllNodes(nodes []*integrationTests.TestProcessorNode, transferValue uint64) {
	balanceFirstTransaction := transferValue + integrationTests.MinTxGasLimit*integrationTests.MinTxGasPrice
	balanceSecondTransaction := integrationTests.MinTxGasLimit * integrationTests.MinTxGasPrice
	totalBalance := balanceFirstTransaction + balanceSecondTransaction

	integrationTests.MintAllNodes(nodes, big.NewInt(0).SetUint64(totalBalance))
}

func generateTransferTxs(
	transferValue uint64,
	sk crypto.PrivateKey,
	pkReceiver crypto.PublicKey,
) ([]data.TransactionHandler, [][]byte) {

	numTxs := 3
	txs := make([]data.TransactionHandler, numTxs)
	hashes := make([][]byte, numTxs)
	for i := 0; i < numTxs; i++ {
		txs[i] = integrationTests.GenerateTransferTx(
			uint64(i),
			sk,
			pkReceiver,
			big.NewInt(0).SetUint64(transferValue),
			integrationTests.MinTxGasPrice,
			integrationTests.MinTxGasLimit,
			integrationTests.ChainID,
			integrationTests.MinTransactionVersion,
		)

		hashes[i], _ = core.CalculateHash(integrationTests.TestMarshalizer, integrationTests.TestHasher, txs[i])
	}

	return txs, hashes
}

func addTxsInDataPool(proposer *integrationTests.TestProcessorNode, txs []data.TransactionHandler, hashes [][]byte) {
	shardId := proposer.ShardCoordinator.SelfId()
	cacherIdentifier := process.ShardCacherIdentifier(shardId, shardId)
	txCache := proposer.DataPool.Transactions()

	for i := 0; i < len(txs); i++ {
		txCache.AddData(hashes[i], txs[i], txs[i].Size(), cacherIdentifier)
	}
}

func testStateOnNodes(t *testing.T, nodes []*integrationTests.TestProcessorNode, idxProposer int, hashes [][]byte) {
	proposer := nodes[idxProposer]

	expectedHeaderNonce := uint64(1)
	txValidIdx := 0
	txInvalidIdx := 1
	txDeletedIdx := 2

	testSameBlockHeight(t, nodes, idxProposer, expectedHeaderNonce)
	testTxIsInMiniblock(t, proposer, hashes[txValidIdx], block.TxBlock)
	testTxIsInMiniblock(t, proposer, hashes[txInvalidIdx], block.InvalidBlock)
	testTxIsInNotInBody(t, proposer, hashes[txDeletedIdx])

	// Removed from mempool.
	_, ok := proposer.DataPool.Transactions().SearchFirstData(hashes[txValidIdx])
	assert.False(t, ok)

	// Removed from mempool.
	_, ok = proposer.DataPool.Transactions().SearchFirstData(hashes[txInvalidIdx])
	assert.False(t, ok)

	// Not removed from mempool (see DRT-16200).
	_, ok = proposer.DataPool.Transactions().SearchFirstData(hashes[txDeletedIdx])
	assert.True(t, ok)
}

func testSameBlockHeight(t *testing.T, nodes []*integrationTests.TestProcessorNode, idxProposer int, expectedHeight uint64) {
	proposer := nodes[idxProposer]

	for _, n := range nodes {
		assert.NotNil(t, n.BlockChain.GetCurrentBlockHeader())
		assert.Equal(t, expectedHeight, n.BlockChain.GetCurrentBlockHeader().GetNonce())
		assert.Equal(t, proposer.BlockChain.GetCurrentBlockHeaderHash(), n.BlockChain.GetCurrentBlockHeaderHash())
	}
}

func testTxIsInMiniblock(t *testing.T, proposer *integrationTests.TestProcessorNode, hash []byte, bt block.Type) {
	hdrHandler := proposer.BlockChain.GetCurrentBlockHeader()
	hdr := hdrHandler.(*block.Header)

	for _, mbh := range hdr.MiniBlockHeaders {
		if mbh.Type != bt {
			continue
		}

		mbBuff, err := proposer.Storage.Get(dataRetriever.MiniBlockUnit, mbh.Hash)
		assert.Nil(t, err)

		miniblock := &block.MiniBlock{}
		_ = integrationTests.TestMarshalizer.Unmarshal(miniblock, mbBuff)

		for _, txHash := range miniblock.TxHashes {
			if bytes.Equal(hash, txHash) {
				return
			}
		}
	}

	assert.Fail(t, fmt.Sprintf("hash %s not found in miniblock type %s", logger.DisplayByteSlice(hash), bt.String()))
}

func testTxIsInNotInBody(t *testing.T, proposer *integrationTests.TestProcessorNode, hash []byte) {
	hdrHandler := proposer.BlockChain.GetCurrentBlockHeader()
	hdr := hdrHandler.(*block.Header)

	for _, mbh := range hdr.MiniBlockHeaders {
		mbBuff, err := proposer.Storage.Get(dataRetriever.MiniBlockUnit, mbh.Hash)
		assert.Nil(t, err)

		miniblock := &block.MiniBlock{}
		_ = integrationTests.TestMarshalizer.Unmarshal(miniblock, mbBuff)

		for _, txHash := range miniblock.TxHashes {
			if bytes.Equal(hash, txHash) {
				assert.Fail(t, fmt.Sprintf("hash %s should not have been not found in miniblock type %s",
					logger.DisplayByteSlice(hash), miniblock.Type.String()))
			}
		}
	}
}

func proposeAndCommitBlock(node *integrationTests.TestProcessorNode, round uint64, nonce uint64) error {
	body, hdr, _ := node.ProposeBlock(round, nonce)
	err := node.BlockProcessor.CommitBlock(hdr, body)
	if err != nil {
		return err
	}

	pk := node.NodeKeys.MainKey.Pk
	node.BroadcastBlock(body, hdr, pk)
	time.Sleep(testBlock.StepDelay)
	return nil
}

func checkCurrentBlockHeight(t *testing.T, nodes []*integrationTests.TestProcessorNode, nonce uint64) {
	for _, n := range nodes {
		assert.Equal(t, nonce, n.BlockChain.GetCurrentBlockHeader().GetNonce())
	}
}
