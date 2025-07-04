package basicSync

import (
	"fmt"
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	logger "github.com/TerraDharitri/drt-go-chain-logger"
	"github.com/stretchr/testify/assert"

	"github.com/TerraDharitri/drt-go-chain/integrationTests"
)

var log = logger.GetOrCreate("basicSync")

func TestSyncWorksInShard_EmptyBlocksNoForks(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}
	maxShards := uint32(1)
	shardId := uint32(0)
	numNodesPerShard := 6

	nodes := make([]*integrationTests.TestProcessorNode, numNodesPerShard+1)
	connectableNodes := make([]integrationTests.Connectable, 0)
	for i := 0; i < numNodesPerShard; i++ {
		nodes[i] = integrationTests.NewTestProcessorNode(integrationTests.ArgTestProcessorNode{
			MaxShards:            maxShards,
			NodeShardId:          shardId,
			TxSignPrivKeyShardId: shardId,
			WithSync:             true,
		})
		connectableNodes = append(connectableNodes, nodes[i])
	}

	metachainNode := integrationTests.NewTestProcessorNode(integrationTests.ArgTestProcessorNode{
		MaxShards:            maxShards,
		NodeShardId:          core.MetachainShardId,
		TxSignPrivKeyShardId: shardId,
		WithSync:             true,
	})
	idxProposerMeta := numNodesPerShard
	nodes[idxProposerMeta] = metachainNode
	connectableNodes = append(connectableNodes, metachainNode)

	idxProposerShard0 := 0
	leaders := []*integrationTests.TestProcessorNode{nodes[idxProposerShard0], nodes[idxProposerMeta]}

	integrationTests.ConnectNodes(connectableNodes)

	defer func() {
		for _, n := range nodes {
			n.Close()
		}
	}()

	for _, n := range nodes {
		_ = n.StartSync()
	}

	fmt.Println("Delaying for nodes p2p bootstrap...")
	time.Sleep(integrationTests.P2pBootstrapDelay)

	round := uint64(0)
	nonce := uint64(0)
	round = integrationTests.IncrementAndPrintRound(round)
	integrationTests.UpdateRound(nodes, round)
	nonce++

	numRoundsToTest := 5
	for i := 0; i < numRoundsToTest; i++ {
		integrationTests.ProposeBlock(nodes, leaders, round, nonce)

		time.Sleep(integrationTests.SyncDelay)

		round = integrationTests.IncrementAndPrintRound(round)
		integrationTests.UpdateRound(nodes, round)
		nonce++
	}

	time.Sleep(integrationTests.SyncDelay)

	testAllNodesHaveTheSameBlockHeightInBlockchain(t, nodes)
}

func TestSyncWorksInShard_EmptyBlocksDoubleSign(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	maxShards := uint32(1)
	shardId := uint32(0)
	numNodesPerShard := 6

	nodes := make([]*integrationTests.TestProcessorNode, numNodesPerShard)
	connectableNodes := make([]integrationTests.Connectable, 0)
	for i := 0; i < numNodesPerShard; i++ {
		nodes[i] = integrationTests.NewTestProcessorNode(integrationTests.ArgTestProcessorNode{
			MaxShards:            maxShards,
			NodeShardId:          shardId,
			TxSignPrivKeyShardId: shardId,
			WithSync:             true,
		})
		connectableNodes = append(connectableNodes, nodes[i])
	}

	integrationTests.ConnectNodes(connectableNodes)

	idxProposerShard0 := 0
	leaders := []*integrationTests.TestProcessorNode{nodes[idxProposerShard0]}

	defer func() {
		for _, n := range nodes {
			n.Close()
		}
	}()

	for _, n := range nodes {
		_ = n.StartSync()
	}

	fmt.Println("Delaying for nodes p2p bootstrap...")
	time.Sleep(integrationTests.P2pBootstrapDelay)

	round := uint64(0)
	nonce := uint64(0)
	round = integrationTests.IncrementAndPrintRound(round)
	integrationTests.UpdateRound(nodes, round)
	nonce++

	numRoundsToTest := 2
	for i := 0; i < numRoundsToTest; i++ {
		integrationTests.ProposeBlock(nodes, leaders, round, nonce)

		time.Sleep(integrationTests.SyncDelay)

		round = integrationTests.IncrementAndPrintRound(round)
		integrationTests.UpdateRound(nodes, round)
		nonce++
	}

	time.Sleep(integrationTests.SyncDelay)

	pubKeysVariant1 := []byte{3}
	pubKeysVariant2 := []byte{1}

	proposeBlockWithPubKeyBitmap(nodes[idxProposerShard0], round, nonce, pubKeysVariant1)
	proposeBlockWithPubKeyBitmap(nodes[1], round, nonce, pubKeysVariant2)

	time.Sleep(integrationTests.StepDelay)

	round = integrationTests.IncrementAndPrintRound(round)
	integrationTests.UpdateRound(nodes, round)

	stepDelayForkResolving := 4 * integrationTests.StepDelay
	time.Sleep(stepDelayForkResolving)

	testAllNodesHaveTheSameBlockHeightInBlockchain(t, nodes)
	testAllNodesHaveSameLastBlock(t, nodes)
}

func proposeBlockWithPubKeyBitmap(n *integrationTests.TestProcessorNode, round uint64, nonce uint64, pubKeys []byte) {
	body, header, _ := n.ProposeBlock(round, nonce)
	err := header.SetPubKeysBitmap(pubKeys)
	if err != nil {
		log.Error("header.SetPubKeysBitmap", "error", err.Error())
	}

	pk := n.NodeKeys.MainKey.Pk
	n.BroadcastBlock(body, header, pk)
	n.CommitBlock(body, header)
}

func testAllNodesHaveTheSameBlockHeightInBlockchain(t *testing.T, nodes []*integrationTests.TestProcessorNode) {
	expectedNonce := nodes[0].BlockChain.GetCurrentBlockHeader().GetNonce()
	for i := 1; i < len(nodes); i++ {
		if check.IfNil(nodes[i].BlockChain.GetCurrentBlockHeader()) {
			assert.Fail(t, fmt.Sprintf("Node with idx %d does not have a current block", i))
		} else {
			assert.Equal(t, expectedNonce, nodes[i].BlockChain.GetCurrentBlockHeader().GetNonce())
		}
	}
}

func testAllNodesHaveSameLastBlock(t *testing.T, nodes []*integrationTests.TestProcessorNode) {
	mapBlocksByHash := make(map[string]data.HeaderHandler)

	for _, n := range nodes {
		hdr := n.BlockChain.GetCurrentBlockHeader()
		buff, _ := core.CalculateHash(integrationTests.TestMarshalizer, integrationTests.TestHasher, hdr)

		mapBlocksByHash[string(buff)] = hdr
	}

	assert.Equal(t, 1, len(mapBlocksByHash))
}

func TestSyncWorksInShard_EmptyBlocksNoForks_With_EquivalentProofs(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	// 3 shard nodes and 1 metachain node
	maxShards := uint32(1)
	shardId := uint32(0)
	numNodesPerShard := 3

	enableEpochs := integrationTests.CreateEnableEpochsConfig()
	enableEpochs.AndromedaEnableEpoch = uint32(0)

	nodes := make([]*integrationTests.TestProcessorNode, numNodesPerShard+1)
	connectableNodes := make([]integrationTests.Connectable, 0)
	for i := 0; i < numNodesPerShard; i++ {
		nodes[i] = integrationTests.NewTestProcessorNode(integrationTests.ArgTestProcessorNode{
			MaxShards:            maxShards,
			NodeShardId:          shardId,
			TxSignPrivKeyShardId: shardId,
			WithSync:             true,
			EpochsConfig:         &enableEpochs,
		})
		connectableNodes = append(connectableNodes, nodes[i])
	}

	metachainNode := integrationTests.NewTestProcessorNode(integrationTests.ArgTestProcessorNode{
		MaxShards:            maxShards,
		NodeShardId:          core.MetachainShardId,
		TxSignPrivKeyShardId: shardId,
		WithSync:             true,
		EpochsConfig:         &enableEpochs,
	})
	idxProposerMeta := numNodesPerShard
	nodes[idxProposerMeta] = metachainNode
	connectableNodes = append(connectableNodes, metachainNode)

	idxProposerShard0 := 0
	leaders := []*integrationTests.TestProcessorNode{nodes[idxProposerShard0], nodes[idxProposerMeta]}

	integrationTests.ConnectNodes(connectableNodes)

	defer func() {
		for _, n := range nodes {
			n.Close()
		}
	}()

	for _, n := range nodes {
		_ = n.StartSync()
	}

	fmt.Println("Delaying for nodes p2p bootstrap...")
	time.Sleep(integrationTests.P2pBootstrapDelay)

	round := uint64(0)
	nonce := uint64(0)
	round = integrationTests.IncrementAndPrintRound(round)
	integrationTests.UpdateRound(nodes, round)
	nonce++

	numRoundsToTest := 5

	for i := 0; i < numRoundsToTest; i++ {
		integrationTests.ProposeBlockWithProof(nodes, leaders, round, nonce)

		time.Sleep(integrationTests.SyncDelay)

		round = integrationTests.IncrementAndPrintRound(round)
		integrationTests.UpdateRound(nodes, round)
		nonce++
	}

	time.Sleep(integrationTests.SyncDelay)

	expectedNonce := nodes[0].BlockChain.GetCurrentBlockHeader().GetNonce()
	for i := 1; i < len(nodes); i++ {
		if check.IfNil(nodes[i].BlockChain.GetCurrentBlockHeader()) {
			assert.Fail(t, fmt.Sprintf("Node with idx %d does not have a current block", i))
		} else {
			// all nodes must have proofs now
			assert.Equal(t, expectedNonce, nodes[i].BlockChain.GetCurrentBlockHeader().GetNonce())
		}
	}
}

func TestSyncMetaAndShard_With_EquivalentProofs(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	// 3 shard nodes and 3 metachain node
	maxShards := uint32(1)
	shardId := uint32(0)
	numNodesPerShard := 3

	enableEpochs := integrationTests.CreateEnableEpochsConfig()
	enableEpochs.AndromedaEnableEpoch = uint32(0)

	nodes := make([]*integrationTests.TestProcessorNode, 2*numNodesPerShard)
	leaders := make([]*integrationTests.TestProcessorNode, 0)
	connectableNodes := make([]integrationTests.Connectable, 0)

	for i := 0; i < numNodesPerShard; i++ {
		nodes[i] = integrationTests.NewTestProcessorNode(integrationTests.ArgTestProcessorNode{
			MaxShards:            maxShards,
			NodeShardId:          shardId,
			TxSignPrivKeyShardId: shardId,
			WithSync:             true,
			EpochsConfig:         &enableEpochs,
		})
		connectableNodes = append(connectableNodes, nodes[i])
	}

	idxProposerShard0 := 0
	leaders = append(leaders, nodes[idxProposerShard0])

	idxProposerMeta := numNodesPerShard
	for i := 0; i < numNodesPerShard; i++ {
		metachainNode := integrationTests.NewTestProcessorNode(integrationTests.ArgTestProcessorNode{
			MaxShards:            maxShards,
			NodeShardId:          core.MetachainShardId,
			TxSignPrivKeyShardId: shardId,
			WithSync:             true,
			EpochsConfig:         &enableEpochs,
		})
		nodes[idxProposerMeta+i] = metachainNode
		connectableNodes = append(connectableNodes, metachainNode)
	}
	leaders = append(leaders, nodes[idxProposerMeta])

	integrationTests.ConnectNodes(connectableNodes)

	defer func() {
		for _, n := range nodes {
			n.Close()
		}
	}()

	for _, n := range nodes {
		_ = n.StartSync()
	}

	fmt.Println("Delaying for nodes p2p bootstrap...")
	time.Sleep(integrationTests.P2pBootstrapDelay)

	round := uint64(0)
	nonce := uint64(0)
	round = integrationTests.IncrementAndPrintRound(round)
	integrationTests.UpdateRound(nodes, round)
	nonce++

	numRoundsToTest := 5
	for i := 0; i < numRoundsToTest; i++ {
		integrationTests.ProposeBlockWithProof(nodes, leaders, round, nonce)

		time.Sleep(integrationTests.SyncDelay)

		round = integrationTests.IncrementAndPrintRound(round)
		integrationTests.UpdateRound(nodes, round)
		nonce++
	}

	time.Sleep(integrationTests.SyncDelay)

	expectedNonce := nodes[0].BlockChain.GetCurrentBlockHeader().GetNonce()
	for i := 1; i < len(nodes); i++ {
		if check.IfNil(nodes[i].BlockChain.GetCurrentBlockHeader()) {
			assert.Fail(t, fmt.Sprintf("Node with idx %d does not have a current block", i))
		} else {
			// all nodes must have proofs now
			assert.Equal(t, expectedNonce, nodes[i].BlockChain.GetCurrentBlockHeader().GetNonce())
		}
	}
}
