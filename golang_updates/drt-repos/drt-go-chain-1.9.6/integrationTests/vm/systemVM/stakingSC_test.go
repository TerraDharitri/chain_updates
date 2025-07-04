package systemVM

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/integrationTests"
	integrationTestsVm "github.com/TerraDharitri/drt-go-chain/integrationTests/vm"
	"github.com/TerraDharitri/drt-go-chain/state"
	"github.com/TerraDharitri/drt-go-chain/state/accounts"
	"github.com/TerraDharitri/drt-go-chain/vm"
)

func TestStakingUnstakingAndUnbondingOnMultiShardEnvironment(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	numOfShards := 2
	nodesPerShard := 2
	numMetachainNodes := 2

	enableEpochsConfig := config.EnableEpochs{
		StakingV2EnableEpoch:                 integrationTests.UnreachableEpoch,
		ScheduledMiniBlocksEnableEpoch:       integrationTests.UnreachableEpoch,
		MiniBlockPartialExecutionEnableEpoch: integrationTests.UnreachableEpoch,
		StakingV4Step1EnableEpoch:            integrationTests.UnreachableEpoch,
		StakingV4Step2EnableEpoch:            integrationTests.UnreachableEpoch,
		StakingV4Step3EnableEpoch:            integrationTests.UnreachableEpoch,
		AndromedaEnableEpoch:                 integrationTests.UnreachableEpoch,
	}

	nodes := integrationTests.CreateNodesWithEnableEpochs(
		numOfShards,
		nodesPerShard,
		numMetachainNodes,
		enableEpochsConfig,
	)

	leaders := make([]*integrationTests.TestProcessorNode, numOfShards+1)
	for i := 0; i < numOfShards; i++ {
		leaders[i] = nodes[i*nodesPerShard]
	}
	leaders[numOfShards] = nodes[numOfShards*nodesPerShard]

	integrationTests.DisplayAndStartNodes(nodes)

	defer func() {
		for _, n := range nodes {
			n.Close()
		}
	}()

	initialVal := big.NewInt(10000000000)
	integrationTests.MintAllNodes(nodes, initialVal)
	verifyInitialBalance(t, nodes, initialVal)
	integrationTests.SaveDelegationManagerConfig(nodes)

	round := uint64(0)
	nonce := uint64(0)
	round = integrationTests.IncrementAndPrintRound(round)
	nonce++

	// ------- send stake tx and check sender's balance
	var txData string
	genesisBlock := nodes[0].GenesisBlocks[core.MetachainShardId]
	metaBlock := genesisBlock.(*block.MetaBlock)
	nodePrice := big.NewInt(0).Set(metaBlock.EpochStart.Economics.NodePrice)
	oneEncoded := hex.EncodeToString(big.NewInt(1).Bytes())
	for index, node := range nodes {
		pubKey := generateUniqueKey(index)
		txData = "stake" + "@" + oneEncoded + "@" + pubKey + "@" + hex.EncodeToString([]byte("msg"))
		integrationTests.CreateAndSendTransaction(node, nodes, nodePrice, vm.ValidatorSCAddress, txData, core.MinMetaTxExtraGasCost)
	}

	time.Sleep(time.Second)

	nrRoundsToPropagateMultiShard := 10
	integrationTests.AddSelfNotarizedHeaderByMetachain(nodes)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)

	time.Sleep(time.Second)

	checkAccountsAfterStaking(t, nodes)

	// ------ send unStake tx
	for index, node := range nodes {
		pubKey := generateUniqueKey(index)
		txData = "unStake" + "@" + pubKey
		integrationTests.CreateAndSendTransaction(node, nodes, big.NewInt(0), vm.ValidatorSCAddress, txData, core.MinMetaTxExtraGasCost)
	}

	roundsPerEpoch := uint64(10)
	for _, node := range nodes {
		node.EpochStartTrigger.SetRoundsPerEpoch(roundsPerEpoch)
		node.WaitTime = 100 * time.Millisecond
	}

	time.Sleep(time.Second)

	integrationTests.AddSelfNotarizedHeaderByMetachain(nodes)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)

	// ----- wait for unbond period
	integrationTests.AddSelfNotarizedHeaderByMetachain(nodes)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, 10, nonce, round)

	manualSetToInactiveStateStakedPeers(t, nodes)

	// ----- send unBond
	for index, node := range nodes {
		pubKey := generateUniqueKey(index)
		txData = "unBond" + "@" + pubKey
		integrationTests.CreateAndSendTransaction(node, nodes, big.NewInt(0), vm.ValidatorSCAddress, txData, core.MinMetaTxExtraGasCost)
	}

	time.Sleep(time.Second)

	integrationTests.AddSelfNotarizedHeaderByMetachain(nodes)
	_, _ = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)

	verifyUnbound(t, nodes)
}

func TestStakingUnstakingAndUnbondingOnMultiShardEnvironmentWithValidatorStatistics(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	numOfShards := 2
	nodesPerShard := 2
	numMetachainNodes := 2
	shardConsensusGroupSize := 1
	metaConsensusGroupSize := 1

	nodesMap := integrationTests.CreateNodesWithNodesCoordinator(
		nodesPerShard,
		numMetachainNodes,
		numOfShards,
		shardConsensusGroupSize,
		metaConsensusGroupSize,
	)

	nodes := make([]*integrationTests.TestProcessorNode, 0)
	leaders := make([]*integrationTests.TestProcessorNode, numOfShards+1)

	for _, nds := range nodesMap {
		nodes = append(nodes, nds...)
	}

	for i := 0; i < numOfShards; i++ {
		leaders[i] = nodesMap[uint32(i)][0]
	}
	leaders[numOfShards] = nodesMap[core.MetachainShardId][0]

	integrationTests.DisplayAndStartNodes(nodes)

	defer func() {
		for _, n := range nodes {
			n.Close()
		}
	}()

	for _, nds := range nodesMap {
		fmt.Println(integrationTests.MakeDisplayTable(nds))
	}

	initialVal := big.NewInt(10000000000)
	integrationTests.MintAllNodes(nodes, initialVal)
	verifyInitialBalance(t, nodes, initialVal)
	integrationTests.SaveDelegationManagerConfig(nodes)

	round := uint64(0)
	nonce := uint64(0)
	round = integrationTests.IncrementAndPrintRound(round)
	nonce++

	// ------- send stake tx and check sender's balance
	genesisBlock := nodes[0].GenesisBlocks[core.MetachainShardId]
	metaBlock := genesisBlock.(*block.MetaBlock)
	nodePrice := big.NewInt(0).Set(metaBlock.EpochStart.Economics.NodePrice)
	oneEncoded := hex.EncodeToString(big.NewInt(1).Bytes())
	var txData string
	for index, node := range nodes {
		pubKey := generateUniqueKey(index)
		txData = "stake" + "@" + oneEncoded + "@" + pubKey + "@" + hex.EncodeToString([]byte("msg"))
		integrationTests.CreateAndSendTransaction(node, nodes, nodePrice, vm.ValidatorSCAddress, txData, core.MinMetaTxExtraGasCost)
	}

	time.Sleep(time.Second)

	nrRoundsToPropagateMultiShard := 10
	integrationTests.AddSelfNotarizedHeaderByMetachain(nodes)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)

	time.Sleep(time.Second)

	consumedBalance := big.NewInt(0).Add(big.NewInt(int64(len(txData))), big.NewInt(0).SetUint64(integrationTests.MinTxGasLimit))
	consumedBalance.Mul(consumedBalance, big.NewInt(0).SetUint64(integrationTests.MinTxGasPrice))

	checkAccountsAfterStaking(t, nodes)

	manualSetToInactiveStateStakedPeers(t, nodes)

	// ------ send unStake tx
	for index, node := range nodes {
		pubKey := generateUniqueKey(index)
		txData = "unStake" + "@" + pubKey
		integrationTests.CreateAndSendTransaction(node, nodes, big.NewInt(0), vm.ValidatorSCAddress, txData, core.MinMetaTxExtraGasCost)
	}
	consumed := big.NewInt(0).Add(big.NewInt(0).SetUint64(integrationTests.MinTxGasLimit), big.NewInt(int64(len(txData))))
	consumed.Mul(consumed, big.NewInt(0).SetUint64(integrationTests.MinTxGasPrice))
	consumedBalance.Add(consumedBalance, consumed)

	time.Sleep(time.Second)

	integrationTests.AddSelfNotarizedHeaderByMetachain(nodes)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)

	roundsPerEpoch := uint64(10)
	for _, node := range nodes {
		node.EpochStartTrigger.SetRoundsPerEpoch(roundsPerEpoch)
		node.WaitTime = 100 * time.Millisecond
	}

	// ----- wait for unbound period
	integrationTests.AddSelfNotarizedHeaderByMetachain(nodes)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, 10, nonce, round)

	// ----- send unBound
	for index, node := range nodes {
		pubKey := generateUniqueKey(index)
		txData = "unBond" + "@" + pubKey
		integrationTests.CreateAndSendTransaction(node, nodes, big.NewInt(0), vm.ValidatorSCAddress, txData, core.MinMetaTxExtraGasCost)
	}
	consumed = big.NewInt(0).Add(big.NewInt(0).SetUint64(integrationTests.MinTxGasLimit), big.NewInt(int64(len(txData))))
	consumed.Mul(consumed, big.NewInt(0).SetUint64(integrationTests.MinTxGasPrice))
	consumedBalance.Add(consumedBalance, consumed)

	time.Sleep(time.Second)

	integrationTests.AddSelfNotarizedHeaderByMetachain(nodes)
	_, _ = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)

	verifyUnbound(t, nodes)
}

func TestStakeWithRewardsAddressAndValidatorStatistics(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	numOfShards := 2
	nodesPerShard := 2
	numMetachainNodes := 2
	shardConsensusGroupSize := 1
	metaConsensusGroupSize := 1

	nodesMap := integrationTests.CreateNodesWithNodesCoordinatorAndTxKeys(
		nodesPerShard,
		numMetachainNodes,
		numOfShards,
		shardConsensusGroupSize,
		metaConsensusGroupSize,
	)

	nodes := make([]*integrationTests.TestProcessorNode, 0)

	for _, nds := range nodesMap {
		nodes = append(nodes, nds...)
	}

	for _, nds := range nodesMap {
		_, err := integrationTestsVm.GetNodeIndex(nodes, nds[0])
		assert.Nil(t, err)
	}
	integrationTests.DisplayAndStartNodes(nodes)

	roundsPerEpoch := uint64(5)
	for _, node := range nodes {
		node.EpochStartTrigger.SetRoundsPerEpoch(roundsPerEpoch)
	}

	defer func() {
		for _, n := range nodes {
			n.Close()
		}
	}()

	for _, node := range nodesMap {
		fmt.Println(integrationTests.MakeDisplayTable(node))
	}

	initialVal := big.NewInt(10000000000)
	integrationTests.MintAllNodes(nodes, initialVal)

	rewardAccount := integrationTests.CreateTestWalletAccount(nodes[0].ShardCoordinator, 0)

	verifyInitialBalance(t, nodes, initialVal)

	round := uint64(0)
	nonce := uint64(0)
	round = integrationTests.IncrementAndPrintRound(round)
	nonce++

	var txData string
	for _, node := range nodes {
		txData = "changeRewardAddress" + "@" + hex.EncodeToString(rewardAccount.Address)
		integrationTests.CreateAndSendTransaction(node, nodes, big.NewInt(0), vm.ValidatorSCAddress, txData, core.MinMetaTxExtraGasCost)
	}

	nbBlocksToProduce := roundsPerEpoch * 3

	for i := uint64(0); i < nbBlocksToProduce; i++ {
		for _, nodesSlice := range nodesMap {
			integrationTests.UpdateRound(nodesSlice, round)
			integrationTests.AddSelfNotarizedHeaderByMetachain(nodesSlice)
		}

		proposeData := integrationTests.AllShardsProposeBlock(round, nonce, nodesMap)
		integrationTests.SyncAllShardsWithRoundBlock(t, proposeData, nodesMap, round)
		round++
		nonce++

		time.Sleep(integrationTests.StepDelay)
	}

	rewardShardID := nodes[0].ShardCoordinator.ComputeId(rewardAccount.Address)
	for _, node := range nodes {
		if node.ShardCoordinator.SelfId() != rewardShardID {
			continue
		}

		rwdAccount := getAccountFromAddrBytes(node.AccntState, rewardAccount.Address)
		assert.True(t, rwdAccount.GetBalance().Cmp(big.NewInt(0)) > 0)
	}
}

func verifyUnbound(t *testing.T, nodes []*integrationTests.TestProcessorNode) {
	expectedValue := big.NewInt(0).SetUint64(9996639300)
	for _, node := range nodes {
		accShardId := node.ShardCoordinator.ComputeId(node.OwnAccount.Address)

		for _, helperNode := range nodes {
			if helperNode.ShardCoordinator.SelfId() == accShardId {
				sndAcc := getAccountFromAddrBytes(helperNode.AccntState, node.OwnAccount.Address)
				assert.Equal(t, expectedValue.String(), sndAcc.GetBalance().String())
			}
		}
	}
}

func checkAccountsAfterStaking(t *testing.T, nodes []*integrationTests.TestProcessorNode) {
	expectedValue := big.NewInt(0).SetUint64(9998878200)
	for _, node := range nodes {
		accShardId := node.ShardCoordinator.ComputeId(node.OwnAccount.Address)

		for _, helperNode := range nodes {
			if helperNode.ShardCoordinator.SelfId() == accShardId {

				sndAcc := getAccountFromAddrBytes(helperNode.AccntState, node.OwnAccount.Address)
				assert.Equal(t, sndAcc.GetBalance().String(), expectedValue.String())
				break
			}
		}
	}
}

func verifyInitialBalance(t *testing.T, nodes []*integrationTests.TestProcessorNode, initialVal *big.Int) {
	for _, node := range nodes {
		accShardId := node.ShardCoordinator.ComputeId(node.OwnAccount.Address)

		for _, helperNode := range nodes {
			if helperNode.ShardCoordinator.SelfId() == accShardId {
				sndAcc := getAccountFromAddrBytes(helperNode.AccntState, node.OwnAccount.Address)
				require.Equal(t, initialVal, sndAcc.GetBalance())
				break
			}
		}
	}
}

func getAccountFromAddrBytes(accState state.AccountsAdapter, address []byte) state.UserAccountHandler {
	sndrAcc, _ := accState.GetExistingAccount(address)

	sndAccSt, _ := sndrAcc.(state.UserAccountHandler)

	return sndAccSt
}

func generateUniqueKey(identifier int) string {
	neededLength := 192
	uniqueIdentifier := fmt.Sprintf("%d", identifier)
	return strings.Repeat("0", neededLength-len(uniqueIdentifier)) + uniqueIdentifier
}

func manualSetToInactiveStateStakedPeers(t *testing.T, nodes []*integrationTests.TestProcessorNode) {
	for _, node := range nodes {
		if node.ShardCoordinator.SelfId() != core.MetachainShardId {
			continue
		}

		for index := range nodes {
			pubKey, _ := hex.DecodeString(generateUniqueKey(index))
			peerAccount, _ := accounts.NewPeerAccount(pubKey)
			peerAccount.List = string(common.InactiveList)
			peerAccount.BLSPublicKey = pubKey
			err := node.PeerState.SaveAccount(peerAccount)
			require.Nil(t, err)
		}

		_, err := node.PeerState.Commit()
		require.Nil(t, err)
	}
}
