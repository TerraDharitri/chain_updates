package epochChangeWithNodesShufflingAndRater

import (
	"math/big"
	"testing"
	"time"

	logger "github.com/TerraDharitri/drt-go-chain-logger"

	"github.com/TerraDharitri/drt-go-chain/integrationTests"
	"github.com/TerraDharitri/drt-go-chain/integrationTests/multiShard/endOfEpoch"
	"github.com/TerraDharitri/drt-go-chain/process/rating"
)

func TestEpochChangeWithNodesShufflingAndRater(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	_ = logger.SetDisplayByteSlice(logger.ToHexShort)

	nodesPerShard := 1
	nbMetaNodes := 1
	nbShards := 1
	consensusGroupSize := 1
	maxGasLimitPerBlock := uint64(100000)

	rater, _ := rating.NewBlockSigningRater(integrationTests.CreateRatingsData())
	coordinatorFactory := &integrationTests.IndexHashedNodesCoordinatorWithRaterFactory{
		PeerAccountListAndRatingHandler: rater,
	}

	// create map of shard - testNodeProcessors for metachain and shard chain
	nodesMap := integrationTests.CreateNodesWithNodesCoordinatorFactory(
		nodesPerShard,
		nbMetaNodes,
		nbShards,
		consensusGroupSize,
		consensusGroupSize,
		coordinatorFactory,
	)

	gasPrice := uint64(10)
	gasLimit := uint64(100)
	valToTransfer := big.NewInt(100)
	nbTxsPerShard := uint32(100)
	mintValue := big.NewInt(1000000)

	defer func() {
		for _, nodes := range nodesMap {
			for _, n := range nodes {
				n.Close()
			}
		}
	}()

	roundsPerEpoch := uint64(7)
	for _, nodes := range nodesMap {
		integrationTests.SetEconomicsParameters(nodes, maxGasLimitPerBlock, gasPrice, gasLimit)
		integrationTests.DisplayAndStartNodes(nodes)
		for _, node := range nodes {
			node.EpochStartTrigger.SetRoundsPerEpoch(roundsPerEpoch)
		}
	}

	integrationTests.GenerateIntraShardTransactions(nodesMap, nbTxsPerShard, mintValue, valToTransfer, gasPrice, gasLimit)

	round := uint64(1)
	nonce := uint64(1)
	nbBlocksToProduce := uint64(20)
	expectedLastEpoch := uint32(nbBlocksToProduce / roundsPerEpoch)

	for i := uint64(0); i < nbBlocksToProduce; i++ {
		for _, nodes := range nodesMap {
			integrationTests.UpdateRound(nodes, round)
		}

		proposeData := integrationTests.AllShardsProposeBlock(round, nonce, nodesMap)
		integrationTests.SyncAllShardsWithRoundBlock(t, proposeData, nodesMap, round)
		round++
		nonce++

		time.Sleep(integrationTests.StepDelay)
	}

	for _, nodes := range nodesMap {
		endOfEpoch.VerifyThatNodesHaveCorrectEpoch(t, expectedLastEpoch, nodes)
	}
}
