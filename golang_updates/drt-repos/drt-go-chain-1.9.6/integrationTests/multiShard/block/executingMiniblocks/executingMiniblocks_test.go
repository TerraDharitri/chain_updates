package executingMiniblocks

import (
	"crypto"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/integrationTests"
	"github.com/TerraDharitri/drt-go-chain/process/factory"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain/state"
)

func TestShouldProcessBlocksInMultiShardArchitecture(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	fmt.Println("Setup nodes...")
	numOfShards := 6
	nodesPerShard := 3
	numMetachainNodes := 1

	senderShard := uint32(0)
	recvShards := []uint32{1, 2}
	round := uint64(0)
	nonce := uint64(0)

	valMinting := big.NewInt(10000000)
	valToTransferPerTx := big.NewInt(2)

	nodes := integrationTests.CreateNodes(
		numOfShards,
		nodesPerShard,
		numMetachainNodes,
	)
	leaders := []*integrationTests.TestProcessorNode{nodes[0], nodes[3], nodes[6], nodes[9], nodes[12], nodes[15], nodes[18]}
	integrationTests.DisplayAndStartNodes(nodes)

	defer func() {
		for _, n := range nodes {
			n.Close()
		}
	}()

	fmt.Println("Generating private keys for senders and receivers...")
	generateCoordinator, _ := sharding.NewMultiShardCoordinator(uint32(numOfShards), 0)
	txToGenerateInEachMiniBlock := 3

	proposerNode := nodes[0]

	// sender shard keys, receivers  keys
	sendersPrivateKeys := make([]crypto.PrivateKey, 3)
	receiversPublicKeys := make(map[uint32][]crypto.PublicKey)
	for i := 0; i < txToGenerateInEachMiniBlock; i++ {
		sendersPrivateKeys[i], _, _ = integrationTests.GenerateSkAndPkInShard(generateCoordinator, senderShard)
		// receivers in same shard with the sender
		_, pk, _ := integrationTests.GenerateSkAndPkInShard(generateCoordinator, senderShard)
		receiversPublicKeys[senderShard] = append(receiversPublicKeys[senderShard], pk)
		// receivers in other shards
		for _, shardId := range recvShards {
			_, pk, _ = integrationTests.GenerateSkAndPkInShard(generateCoordinator, shardId)
			receiversPublicKeys[shardId] = append(receiversPublicKeys[shardId], pk)
		}
	}

	fmt.Println("Minting sender addresses...")
	integrationTests.CreateMintingForSenders(nodes, senderShard, sendersPrivateKeys, valMinting)

	fmt.Println("Generating transactions...")
	integrationTests.GenerateAndDisseminateTxs(
		proposerNode,
		sendersPrivateKeys,
		receiversPublicKeys,
		valToTransferPerTx,
		integrationTests.MinTxGasPrice,
		integrationTests.MinTxGasLimit,
		integrationTests.ChainID,
		integrationTests.MinTransactionVersion,
	)
	fmt.Println("Delaying for disseminating transactions...")
	time.Sleep(time.Second * 5)

	round = integrationTests.IncrementAndPrintRound(round)
	nonce++
	roundsToWait := 6
	for i := 0; i < roundsToWait; i++ {
		round, nonce = integrationTests.ProposeAndSyncOneBlock(t, nodes, leaders, round, nonce)
	}

	gasPricePerTxBigInt := big.NewInt(0).SetUint64(integrationTests.MinTxGasPrice)
	gasLimitPerTxBigInt := big.NewInt(0).SetUint64(integrationTests.MinTxGasLimit)
	gasValue := big.NewInt(0).Mul(gasPricePerTxBigInt, gasLimitPerTxBigInt)
	totalValuePerTx := big.NewInt(0).Add(gasValue, valToTransferPerTx)
	fmt.Println("Test nodes from proposer shard to have the correct balances...")
	for _, n := range nodes {
		isNodeInSenderShard := n.ShardCoordinator.SelfId() == senderShard
		if !isNodeInSenderShard {
			continue
		}

		// test sender balances
		for _, sk := range sendersPrivateKeys {
			valTransferred := big.NewInt(0).Mul(totalValuePerTx, big.NewInt(int64(len(receiversPublicKeys))))
			valRemaining := big.NewInt(0).Sub(valMinting, valTransferred)
			integrationTests.TestPrivateKeyHasBalance(t, n, sk, valRemaining)
		}
		// test receiver balances from same shard
		for _, pk := range receiversPublicKeys[proposerNode.ShardCoordinator.SelfId()] {
			integrationTests.TestPublicKeyHasBalance(t, n, pk, valToTransferPerTx)
		}
	}

	fmt.Println("Test nodes from receiver shards to have the correct balances...")
	for _, n := range nodes {
		isNodeInReceiverShardAndNotProposer := false
		for _, shardId := range recvShards {
			if n.ShardCoordinator.SelfId() == shardId {
				isNodeInReceiverShardAndNotProposer = true
				break
			}
		}
		if !isNodeInReceiverShardAndNotProposer {
			continue
		}

		// test receiver balances from same shard
		for _, pk := range receiversPublicKeys[n.ShardCoordinator.SelfId()] {
			integrationTests.TestPublicKeyHasBalance(t, n, pk, valToTransferPerTx)
		}
	}
}

func TestSimpleTransactionsWithMoreGasWhichYieldInReceiptsInMultiShardedEnvironment(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	numOfShards := 2
	nodesPerShard := 3
	numMetachainNodes := 2

	nodes := integrationTests.CreateNodes(
		numOfShards,
		nodesPerShard,
		numMetachainNodes,
	)

	minGasLimit := uint64(10000)
	for _, node := range nodes {
		node.EconomicsData.SetMinGasLimit(minGasLimit, 0)
	}

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

	initialVal := big.NewInt(1000000000)
	sendValue := big.NewInt(5)
	integrationTests.MintAllNodes(nodes, initialVal)
	receiverAddress := []byte("12345678901234567890123456789012")

	round := uint64(0)
	nonce := uint64(0)
	round = integrationTests.IncrementAndPrintRound(round)
	nonce++

	gasLimit := minGasLimit * 2
	time.Sleep(time.Second)
	nrRoundsToTest := 10
	for i := 0; i <= nrRoundsToTest; i++ {
		integrationTests.UpdateRound(nodes, round)
		integrationTests.ProposeBlock(nodes, leaders, round, nonce)
		integrationTests.SyncBlock(t, nodes, leaders, round)
		round = integrationTests.IncrementAndPrintRound(round)
		nonce++

		for _, node := range nodes {
			integrationTests.PlayerSendsTransaction(
				nodes,
				node.OwnAccount,
				receiverAddress,
				sendValue,
				"",
				gasLimit,
			)
		}

		time.Sleep(integrationTests.StepDelay)
	}

	time.Sleep(time.Second)

	txGasNeed := nodes[0].EconomicsData.GetMinGasLimit(0)
	txGasPrice := nodes[0].EconomicsData.GetMinGasPrice()

	oneTxCost := big.NewInt(0).Add(sendValue, big.NewInt(0).SetUint64(txGasNeed*txGasPrice))
	txTotalCost := big.NewInt(0).Mul(oneTxCost, big.NewInt(int64(nrRoundsToTest)))

	expectedBalance := big.NewInt(0).Sub(initialVal, txTotalCost)
	for _, verifierNode := range nodes {
		for _, node := range nodes {
			accWrp, err := verifierNode.AccntState.GetExistingAccount(node.OwnAccount.Address)
			if err != nil {
				continue
			}

			account, _ := accWrp.(state.UserAccountHandler)
			assert.Equal(t, expectedBalance, account.GetBalance())
		}
	}
}

func TestSimpleTransactionsWithMoreValueThanBalanceYieldReceiptsInMultiShardedEnvironment(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	numOfShards := 2
	nodesPerShard := 2
	numMetachainNodes := 2

	nodes := integrationTests.CreateNodes(
		numOfShards,
		nodesPerShard,
		numMetachainNodes,
	)

	minGasLimit := uint64(10000)
	for _, node := range nodes {
		node.EconomicsData.SetMinGasLimit(minGasLimit, 0)
	}

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

	nrTxsToSend := uint64(10)
	initialVal := big.NewInt(0).SetUint64(nrTxsToSend * minGasLimit * integrationTests.MinTxGasPrice)
	halfInitVal := big.NewInt(0).Div(initialVal, big.NewInt(2))
	integrationTests.MintAllNodes(nodes, initialVal)
	receiverAddress := []byte("12345678901234567890123456789012")

	round := uint64(0)
	nonce := uint64(0)
	round = integrationTests.IncrementAndPrintRound(round)
	nonce++

	for _, node := range nodes {
		for j := uint64(0); j < nrTxsToSend; j++ {
			integrationTests.PlayerSendsTransaction(
				nodes,
				node.OwnAccount,
				receiverAddress,
				halfInitVal,
				"",
				minGasLimit,
			)
		}
	}

	time.Sleep(2 * time.Second)

	integrationTests.UpdateRound(nodes, round)
	integrationTests.ProposeBlock(nodes, leaders, round, nonce)
	integrationTests.SyncBlock(t, nodes, leaders, round)
	round = integrationTests.IncrementAndPrintRound(round)
	nonce++

	for _, node := range nodes {
		if node.ShardCoordinator.SelfId() == core.MetachainShardId {
			continue
		}

		header := node.BlockChain.GetCurrentBlockHeader()
		shardHdr, ok := header.(*block.Header)
		numInvalid := 0
		require.True(t, ok)
		for _, mb := range shardHdr.MiniBlockHeaders {
			if mb.Type == block.InvalidBlock {
				numInvalid++
			}
		}
		assert.Equal(t, 1, numInvalid)
	}

	time.Sleep(time.Second)
	numRoundsToTest := 6
	for i := 0; i < numRoundsToTest; i++ {
		integrationTests.UpdateRound(nodes, round)
		integrationTests.ProposeBlock(nodes, leaders, round, nonce)
		integrationTests.SyncBlock(t, nodes, leaders, round)
		round = integrationTests.IncrementAndPrintRound(round)
		nonce++

		time.Sleep(integrationTests.StepDelay)
	}

	time.Sleep(time.Second)

	expectedReceiverValue := big.NewInt(0).Mul(big.NewInt(int64(len(nodes))), halfInitVal)
	for _, verifierNode := range nodes {
		for _, node := range nodes {
			accWrp, err := verifierNode.AccntState.GetExistingAccount(node.OwnAccount.Address)
			if err != nil {
				continue
			}

			account, _ := accWrp.(state.UserAccountHandler)
			assert.Equal(t, big.NewInt(0), account.GetBalance())
		}

		accWrp, err := verifierNode.AccntState.GetExistingAccount(receiverAddress)
		if err != nil {
			continue
		}

		account, _ := accWrp.(state.UserAccountHandler)
		assert.Equal(t, expectedReceiverValue, account.GetBalance())
	}
}

// TestShouldSubtractTheCorrectTxFee uses the mock VM as it's gas model is predictable
// The test checks the tx fee subtraction from the sender account when deploying a SC
// It also checks the fee obtained by the leader is correct
// Test parameters: 2 shards + meta, each with 2 nodes
func TestShouldSubtractTheCorrectTxFee(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	maxShards := 2
	consensusGroupSize := 2
	nodesPerShard := 2

	// create map of shards - testNodeProcessors for metachain and shard chain
	nodesMap := integrationTests.CreateNodesWithNodesCoordinator(
		nodesPerShard,
		nodesPerShard,
		maxShards,
		consensusGroupSize,
		consensusGroupSize,
	)

	for _, nodes := range nodesMap {
		integrationTests.DisplayAndStartNodes(nodes)
		integrationTests.SetEconomicsParameters(nodes, integrationTests.MaxGasLimitPerBlock, integrationTests.MinTxGasPrice, integrationTests.MinTxGasLimit)
	}

	defer func() {
		for _, nodes := range nodesMap {
			for _, n := range nodes {
				n.Close()
			}
		}
	}()

	fmt.Println("Delaying for nodes p2p bootstrap...")
	time.Sleep(integrationTests.StepDelay)

	round := uint64(0)
	nonce := uint64(0)
	round = integrationTests.IncrementAndPrintRound(round)
	nonce++

	initialVal := big.NewInt(10000000)
	senders := integrationTests.CreateSendersWithInitialBalances(nodesMap, initialVal)

	deployValue := big.NewInt(0)
	nodeShard0 := nodesMap[0][0]
	txData := "DEADBEEF@" + hex.EncodeToString(factory.InternalTestingVM) + "@00"
	dummyTx := &transaction.Transaction{
		Data: []byte(txData),
	}
	gasLimit := nodeShard0.EconomicsData.ComputeGasLimit(dummyTx)
	gasLimit += integrationTests.OpGasValueForMockVm
	gasPrice := integrationTests.MinTxGasPrice
	txNonce := uint64(0)
	owner := senders[0][0]
	ownerPk, _ := owner.GeneratePublic().ToByteArray()
	integrationTests.ScCallTxWithParams(
		nodeShard0,
		owner,
		txNonce,
		txData,
		deployValue,
		gasLimit,
		gasPrice,
	)

	proposeData := integrationTests.AllShardsProposeBlock(round, nonce, nodesMap)
	shardId0 := uint32(0)

	_ = integrationTests.IncrementAndPrintRound(round)

	// test sender account decreased its balance with gasPrice * gasLimit
	accnt, err := proposeData[shardId0].Leader.AccntState.GetExistingAccount(ownerPk)
	assert.Nil(t, err)
	ownerAccnt := accnt.(state.UserAccountHandler)
	expectedBalance := big.NewInt(0).Set(initialVal)
	tx := &transaction.Transaction{GasPrice: gasPrice, GasLimit: gasLimit, Data: []byte(txData)}
	txCost := proposeData[shardId0].Leader.EconomicsData.ComputeTxFee(tx)
	expectedBalance.Sub(expectedBalance, txCost)
	assert.Equal(t, expectedBalance, ownerAccnt.GetBalance())

	printContainingTxs(proposeData[shardId0].Leader, proposeData[shardId0].Leader.BlockChain.GetCurrentBlockHeader().(*block.Header))
}

func printContainingTxs(tpn *integrationTests.TestProcessorNode, hdr data.HeaderHandler) {
	for _, miniblockHdr := range hdr.GetMiniBlockHeaderHandlers() {
		miniblockBytes, err := tpn.Storage.Get(dataRetriever.MiniBlockUnit, miniblockHdr.GetHash())
		if err != nil {
			fmt.Println("miniblock " + base64.StdEncoding.EncodeToString(miniblockHdr.GetHash()) + "not found")
			continue
		}

		miniblock := &block.MiniBlock{}
		err = integrationTests.TestMarshalizer.Unmarshal(miniblock, miniblockBytes)
		if err != nil {
			fmt.Println("can not unmarshal miniblock " + base64.StdEncoding.EncodeToString(miniblockHdr.GetHash()))
			continue
		}

		for _, txHash := range miniblock.TxHashes {
			txBytes := []byte("not found")

			mbType := block.Type(miniblockHdr.GetTypeInt32())
			switch mbType {
			case block.TxBlock:
				txBytes, err = tpn.Storage.Get(dataRetriever.TransactionUnit, txHash)
				if err != nil {
					fmt.Println("tx hash " + base64.StdEncoding.EncodeToString(txHash) + " not found")
					continue
				}
			case block.SmartContractResultBlock:
				txBytes, err = tpn.Storage.Get(dataRetriever.UnsignedTransactionUnit, txHash)
				if err != nil {
					fmt.Println("scr hash " + base64.StdEncoding.EncodeToString(txHash) + " not found")
					continue
				}
			case block.RewardsBlock:
				txBytes, err = tpn.Storage.Get(dataRetriever.RewardTransactionUnit, txHash)
				if err != nil {
					fmt.Println("reward hash " + base64.StdEncoding.EncodeToString(txHash) + " not found")
					continue
				}
			}

			fmt.Println(string(txBytes))
		}
	}
}
