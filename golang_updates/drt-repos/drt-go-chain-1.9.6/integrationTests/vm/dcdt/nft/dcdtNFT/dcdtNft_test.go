package dcdtNFT

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"

	"github.com/TerraDharitri/drt-go-chain/integrationTests"
	"github.com/TerraDharitri/drt-go-chain/integrationTests/vm/dcdt"
	"github.com/TerraDharitri/drt-go-chain/integrationTests/vm/dcdt/nft"
	"github.com/TerraDharitri/drt-go-chain/vm"
)

func TestDCDTNonFungibleTokenCreateAndBurn(t *testing.T) {
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

	round := uint64(0)
	nonce := uint64(0)
	round = integrationTests.IncrementAndPrintRound(round)
	nonce++

	roles := [][]byte{
		[]byte(core.DCDTRoleNFTCreate),
		[]byte(core.DCDTRoleNFTBurn),
	}

	tokenIdentifier, nftMetaData := nft.PrepareNFTWithRoles(
		t,
		nodes,
		leaders,
		nodes[1],
		&round,
		&nonce,
		core.NonFungibleDCDT,
		1,
		roles,
	)

	// decrease quantity
	nonceArg := hex.EncodeToString(big.NewInt(0).SetUint64(1).Bytes())
	quantityToBurn := int64(1)
	quantityToBurnArg := hex.EncodeToString(big.NewInt(quantityToBurn).Bytes())
	txData := []byte(core.BuiltInFunctionDCDTNFTBurn + "@" + hex.EncodeToString([]byte(tokenIdentifier)) +
		"@" + nonceArg + "@" + quantityToBurnArg)
	integrationTests.CreateAndSendTransaction(
		nodes[1],
		nodes,
		big.NewInt(0),
		nodes[1].OwnAccount.Address,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)

	time.Sleep(time.Second)
	nrRoundsToPropagateMultiShard := 5
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	// the token data is removed from trie if the quantity is 0, so we should not find it
	nftMetaData.Quantity = 0
	nft.CheckNftData(
		t,
		nodes[1].OwnAccount.Address,
		nodes[1].OwnAccount.Address,
		nodes,
		[]byte(tokenIdentifier),
		nftMetaData,
		1,
	)
}

func TestDCDTSemiFungibleTokenCreateAddAndBurn(t *testing.T) {
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

	round := uint64(0)
	nonce := uint64(0)
	round = integrationTests.IncrementAndPrintRound(round)
	nonce++

	roles := [][]byte{
		[]byte(core.DCDTRoleNFTCreate),
		[]byte(core.DCDTRoleNFTAddQuantity),
		[]byte(core.DCDTRoleNFTBurn),
	}

	initialQuantity := int64(5)
	tokenIdentifier, nftMetaData := nft.PrepareNFTWithRoles(
		t,
		nodes,
		leaders,
		nodes[1],
		&round,
		&nonce,
		core.SemiFungibleDCDT,
		initialQuantity,
		roles,
	)

	// increase quantity
	nonceArg := hex.EncodeToString(big.NewInt(0).SetUint64(1).Bytes())
	quantityToAdd := int64(4)
	quantityToAddArg := hex.EncodeToString(big.NewInt(quantityToAdd).Bytes())
	txData := []byte(core.BuiltInFunctionDCDTNFTAddQuantity + "@" + hex.EncodeToString([]byte(tokenIdentifier)) +
		"@" + nonceArg + "@" + quantityToAddArg)
	integrationTests.CreateAndSendTransaction(
		nodes[1],
		nodes,
		big.NewInt(0),
		nodes[1].OwnAccount.Address,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)

	time.Sleep(time.Second)
	nrRoundsToPropagateMultiShard := 5
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	nftMetaData.Quantity += quantityToAdd
	nft.CheckNftData(
		t,
		nodes[1].OwnAccount.Address,
		nodes[1].OwnAccount.Address,
		nodes,
		[]byte(tokenIdentifier),
		nftMetaData,
		1,
	)

	time.Sleep(time.Second)
	nrRoundsToPropagateMultiShard = 5
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	nft.CheckNftData(
		t,
		nodes[1].OwnAccount.Address,
		nodes[1].OwnAccount.Address,
		nodes,
		[]byte(tokenIdentifier),
		nftMetaData,
		1,
	)

	// burn quantity
	quantityToBurn := int64(4)
	quantityToBurnArg := hex.EncodeToString(big.NewInt(quantityToBurn).Bytes())
	txData = []byte(core.BuiltInFunctionDCDTNFTBurn + "@" + hex.EncodeToString([]byte(tokenIdentifier)) +
		"@" + nonceArg + "@" + quantityToBurnArg)
	integrationTests.CreateAndSendTransaction(
		nodes[1],
		nodes,
		big.NewInt(0),
		nodes[1].OwnAccount.Address,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)

	time.Sleep(time.Second)
	nrRoundsToPropagateMultiShard = 5
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	nftMetaData.Quantity -= quantityToBurn
	nft.CheckNftData(
		t,
		nodes[1].OwnAccount.Address,
		nodes[1].OwnAccount.Address,
		nodes,
		[]byte(tokenIdentifier),
		nftMetaData,
		1,
	)
}

func TestDCDTNonFungibleTokenTransferSelfShard(t *testing.T) {
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

	round := uint64(0)
	nonce := uint64(0)
	round = integrationTests.IncrementAndPrintRound(round)
	nonce++

	roles := [][]byte{
		[]byte(core.DCDTRoleNFTCreate),
		[]byte(core.DCDTRoleNFTBurn),
	}
	tokenIdentifier, nftMetaData := nft.PrepareNFTWithRoles(
		t,
		nodes,
		leaders,
		nodes[1],
		&round,
		&nonce,
		core.NonFungibleDCDT,
		1,
		roles,
	)

	// transfer

	// get a node from the shard
	var nodeInSameShard = nodes[0]
	for _, node := range nodes {
		if node.ShardCoordinator.SelfId() == nodes[1].ShardCoordinator.SelfId() &&
			!bytes.Equal(node.OwnAccount.Address, node.OwnAccount.Address) {
			nodeInSameShard = node
			break
		}
	}

	nonceArg := hex.EncodeToString(big.NewInt(0).SetUint64(1).Bytes())
	quantityToTransfer := int64(1)
	quantityToTransferArg := hex.EncodeToString(big.NewInt(quantityToTransfer).Bytes())
	txData := []byte(core.BuiltInFunctionDCDTNFTTransfer + "@" + hex.EncodeToString([]byte(tokenIdentifier)) +
		"@" + nonceArg + "@" + quantityToTransferArg + "@" + hex.EncodeToString(nodeInSameShard.OwnAccount.Address))
	integrationTests.CreateAndSendTransaction(
		nodes[1],
		nodes,
		big.NewInt(0),
		nodes[1].OwnAccount.Address,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)

	time.Sleep(time.Second)
	nrRoundsToPropagateMultiShard := 5
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	// check that the new address owns the NFT
	nft.CheckNftData(
		t,
		nodes[1].OwnAccount.Address,
		nodeInSameShard.OwnAccount.Address,
		nodes,
		[]byte(tokenIdentifier),
		nftMetaData,
		1,
	)

	// check that the creator doesn't have the token data in trie anymore
	nftMetaData.Quantity = 0
	nft.CheckNftData(
		t,
		nodes[1].OwnAccount.Address,
		nodes[1].OwnAccount.Address,
		nodes,
		[]byte(tokenIdentifier),
		nftMetaData,
		1,
	)
}

func TestDCDTSemiFungibleTokenTransferCrossShard(t *testing.T) {
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

	round := uint64(0)
	nonce := uint64(0)
	round = integrationTests.IncrementAndPrintRound(round)
	nonce++

	// get a node from a different shard
	var nodeInDifferentShard = nodes[0]
	for _, node := range nodes {
		if node.ShardCoordinator.SelfId() != nodes[0].ShardCoordinator.SelfId() {
			nodeInDifferentShard = node
			break
		}
	}

	roles := [][]byte{
		[]byte(core.DCDTRoleNFTCreate),
		[]byte(core.DCDTRoleNFTAddQuantity),
		[]byte(core.DCDTRoleNFTBurn),
	}

	initialQuantity := int64(5)
	tokenIdentifier, nftMetaData := nft.PrepareNFTWithRoles(
		t,
		nodes,
		leaders,
		nodeInDifferentShard,
		&round,
		&nonce,
		core.SemiFungibleDCDT,
		initialQuantity,
		roles,
	)

	// increase quantity
	nonceArg := hex.EncodeToString(big.NewInt(0).SetUint64(1).Bytes())
	quantityToAdd := int64(4)
	quantityToAddArg := hex.EncodeToString(big.NewInt(quantityToAdd).Bytes())
	txData := []byte(core.BuiltInFunctionDCDTNFTAddQuantity + "@" + hex.EncodeToString([]byte(tokenIdentifier)) +
		"@" + nonceArg + "@" + quantityToAddArg)
	integrationTests.CreateAndSendTransaction(
		nodeInDifferentShard,
		nodes,
		big.NewInt(0),
		nodeInDifferentShard.OwnAccount.Address,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)

	time.Sleep(time.Second)
	nrRoundsToPropagateMultiShard := 5
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	nftMetaData.Quantity += quantityToAdd
	nft.CheckNftData(
		t,
		nodeInDifferentShard.OwnAccount.Address,
		nodeInDifferentShard.OwnAccount.Address,
		nodes,
		[]byte(tokenIdentifier),
		nftMetaData,
		1,
	)

	time.Sleep(time.Second)
	nrRoundsToPropagateMultiShard = 5
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	nft.CheckNftData(
		t,
		nodeInDifferentShard.OwnAccount.Address,
		nodeInDifferentShard.OwnAccount.Address,
		nodes,
		[]byte(tokenIdentifier),
		nftMetaData,
		1,
	)

	// transfer
	quantityToTransfer := int64(4)
	quantityToTransferArg := hex.EncodeToString(big.NewInt(quantityToTransfer).Bytes())
	txData = []byte(core.BuiltInFunctionDCDTNFTTransfer + "@" + hex.EncodeToString([]byte(tokenIdentifier)) +
		"@" + nonceArg + "@" + quantityToTransferArg + "@" + hex.EncodeToString(nodes[0].OwnAccount.Address))
	integrationTests.CreateAndSendTransaction(
		nodeInDifferentShard,
		nodes,
		big.NewInt(0),
		nodeInDifferentShard.OwnAccount.Address,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)

	time.Sleep(time.Second)
	nrRoundsToPropagateMultiShard = 11
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	nftMetaData.Quantity = initialQuantity + quantityToAdd - quantityToTransfer
	nft.CheckNftData(
		t,
		nodeInDifferentShard.OwnAccount.Address,
		nodeInDifferentShard.OwnAccount.Address,
		nodes,
		[]byte(tokenIdentifier),
		nftMetaData,
		1,
	)

	nftMetaData.Quantity = quantityToTransfer
	nft.CheckNftData(
		t,
		nodeInDifferentShard.OwnAccount.Address,
		nodes[0].OwnAccount.Address,
		nodes,
		[]byte(tokenIdentifier),
		nftMetaData,
		1,
	)
}

func TestDCDTSemiFungibleTokenTransferToSystemScAddressShouldReceiveBack(t *testing.T) {
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

	round := uint64(0)
	nonce := uint64(0)
	round = integrationTests.IncrementAndPrintRound(round)
	nonce++

	roles := [][]byte{
		[]byte(core.DCDTRoleNFTCreate),
		[]byte(core.DCDTRoleNFTAddQuantity),
		[]byte(core.DCDTRoleNFTBurn),
	}

	initialQuantity := int64(5)
	tokenIdentifier, nftMetaData := nft.PrepareNFTWithRoles(
		t,
		nodes,
		leaders,
		nodes[0],
		&round,
		&nonce,
		core.SemiFungibleDCDT,
		initialQuantity,
		roles,
	)

	// increase quantity
	nonceArg := hex.EncodeToString(big.NewInt(0).SetUint64(1).Bytes())
	quantityToAdd := int64(4)
	quantityToAddArg := hex.EncodeToString(big.NewInt(quantityToAdd).Bytes())
	txData := []byte(core.BuiltInFunctionDCDTNFTAddQuantity + "@" + hex.EncodeToString([]byte(tokenIdentifier)) +
		"@" + nonceArg + "@" + quantityToAddArg)
	integrationTests.CreateAndSendTransaction(
		nodes[0],
		nodes,
		big.NewInt(0),
		nodes[0].OwnAccount.Address,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)

	time.Sleep(time.Second)
	nrRoundsToPropagateMultiShard := 5
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	nftMetaData.Quantity += quantityToAdd
	nft.CheckNftData(
		t,
		nodes[0].OwnAccount.Address,
		nodes[0].OwnAccount.Address,
		nodes,
		[]byte(tokenIdentifier),
		nftMetaData,
		1,
	)

	time.Sleep(time.Second)
	nrRoundsToPropagateMultiShard = 5
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	nft.CheckNftData(
		t,
		nodes[0].OwnAccount.Address,
		nodes[0].OwnAccount.Address,
		nodes,
		[]byte(tokenIdentifier),
		nftMetaData,
		1,
	)

	// transfer
	quantityToTransfer := int64(4)
	quantityToTransferArg := hex.EncodeToString(big.NewInt(quantityToTransfer).Bytes())
	txData = []byte(core.BuiltInFunctionDCDTNFTTransfer + "@" + hex.EncodeToString([]byte(tokenIdentifier)) +
		"@" + nonceArg + "@" + quantityToTransferArg + "@" + hex.EncodeToString(vm.DCDTSCAddress))
	integrationTests.CreateAndSendTransaction(
		nodes[0],
		nodes,
		big.NewInt(0),
		nodes[0].OwnAccount.Address,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)

	time.Sleep(time.Second)
	nrRoundsToPropagateMultiShard = 15
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	nftMetaData.Quantity = 0 // make sure that the DCDT SC address didn't receive the token
	nft.CheckNftData(
		t,
		nodes[0].OwnAccount.Address,
		vm.DCDTSCAddress,
		nodes,
		[]byte(tokenIdentifier),
		nftMetaData,
		1,
	)

	nftMetaData.Quantity = initialQuantity + quantityToAdd // should have the same quantity as before transferring
	nft.CheckNftData(
		t,
		nodes[0].OwnAccount.Address,
		nodes[0].OwnAccount.Address,
		nodes,
		[]byte(tokenIdentifier),
		nftMetaData,
		1,
	)
}

func testNFTSendCreateRole(t *testing.T, numOfShards int) {
	nodes, leaders := dcdt.CreateNodesAndPrepareBalances(numOfShards)

	defer func() {
		for _, n := range nodes {
			n.Close()
		}
	}()

	initialVal := big.NewInt(10000000000)
	integrationTests.MintAllNodes(nodes, initialVal)

	round := uint64(0)
	nonce := uint64(0)
	round = integrationTests.IncrementAndPrintRound(round)
	nonce++

	roles := [][]byte{
		[]byte(core.DCDTRoleNFTCreate),
	}

	nftCreator := nodes[0]
	initialQuantity := int64(1)
	tokenIdentifier, nftMetaData := nft.PrepareNFTWithRoles(
		t,
		nodes,
		leaders,
		nftCreator,
		&round,
		&nonce,
		core.SemiFungibleDCDT,
		initialQuantity,
		roles,
	)

	nftCreatorShId := nftCreator.ShardCoordinator.ComputeId(nftCreator.OwnAccount.Address)
	nextNftCreator := nodes[1]
	for _, node := range nodes {
		if node.ShardCoordinator.ComputeId(node.OwnAccount.Address) != nftCreatorShId {
			nextNftCreator = node
			break
		}
	}

	// transferNFTCreateRole
	txData := []byte("transferNFTCreateRole" + "@" + hex.EncodeToString([]byte(tokenIdentifier)) +
		"@" + hex.EncodeToString(nftCreator.OwnAccount.Address) +
		"@" + hex.EncodeToString(nextNftCreator.OwnAccount.Address))
	integrationTests.CreateAndSendTransaction(
		nftCreator,
		nodes,
		big.NewInt(0),
		vm.DCDTSCAddress,
		string(txData),
		integrationTests.AdditionalGasLimit+core.MinMetaTxExtraGasCost,
	)

	time.Sleep(time.Second)
	nrRoundsToPropagateMultiShard := 20
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	nft.CreateNFT(
		[]byte(tokenIdentifier),
		nextNftCreator,
		nodes,
		nftMetaData,
	)

	time.Sleep(time.Second)
	nrRoundsToPropagateMultiShard = 2
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	nft.CheckNftData(
		t,
		nextNftCreator.OwnAccount.Address,
		nextNftCreator.OwnAccount.Address,
		nodes,
		[]byte(tokenIdentifier),
		nftMetaData,
		2,
	)
}

func TestDCDTNFTSendCreateRoleInShard(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	testNFTSendCreateRole(t, 1)
}

func TestDCDTNFTSendCreateRoleInCrossShard(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	testNFTSendCreateRole(t, 2)
}

func TestDCDTSemiFungibleWithTransferRoleIntraShard(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	testDCDTSemiFungibleTokenTransferRole(t, 1)
}

func TestDCDTSemiFungibleWithTransferRoleCrossShard(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	testDCDTSemiFungibleTokenTransferRole(t, 2)
}

func testDCDTSemiFungibleTokenTransferRole(t *testing.T, numOfShards int) {
	nodesPerShard := 2
	numMetachainNodes := 2

	nodes := integrationTests.CreateNodes(
		numOfShards,
		nodesPerShard,
		numMetachainNodes,
	)

	leaders := make([]*integrationTests.TestProcessorNode, numOfShards+1)
	for i := 0; i < numOfShards; i++ {
		leaders[i] = nodes[i*nodesPerShard]
	}
	leaders[numOfShards] = nodes[numOfShards*nodesPerShard]

	integrationTests.DisplayAndStartNodes(nodes)

	defer func() {
		for _, n := range nodes {
			_ = n.MainMessenger.Close()
		}
	}()

	initialVal := big.NewInt(10000000000)
	integrationTests.MintAllNodes(nodes, initialVal)

	round := uint64(0)
	nonce := uint64(0)
	round = integrationTests.IncrementAndPrintRound(round)
	nonce++

	// get a node from a different shard
	var nodeInDifferentShard = nodes[0]
	for _, node := range nodes {
		if node.ShardCoordinator.SelfId() != nodes[0].ShardCoordinator.SelfId() {
			nodeInDifferentShard = node
			break
		}
	}

	roles := [][]byte{
		[]byte(core.DCDTRoleNFTCreate),
		[]byte(core.DCDTRoleNFTAddQuantity),
		[]byte(core.DCDTRoleNFTBurn),
		[]byte(core.DCDTRoleTransfer),
	}

	initialQuantity := int64(5)
	tokenIdentifier, nftMetaData := nft.PrepareNFTWithRoles(
		t,
		nodes,
		leaders,
		nodeInDifferentShard,
		&round,
		&nonce,
		core.SemiFungibleDCDT,
		initialQuantity,
		roles,
	)

	// increase quantity
	nonceArg := hex.EncodeToString(big.NewInt(0).SetUint64(1).Bytes())
	quantityToAdd := int64(4)
	quantityToAddArg := hex.EncodeToString(big.NewInt(quantityToAdd).Bytes())
	txData := []byte(core.BuiltInFunctionDCDTNFTAddQuantity + "@" + hex.EncodeToString([]byte(tokenIdentifier)) +
		"@" + nonceArg + "@" + quantityToAddArg)
	integrationTests.CreateAndSendTransaction(
		nodeInDifferentShard,
		nodes,
		big.NewInt(0),
		nodeInDifferentShard.OwnAccount.Address,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)

	time.Sleep(time.Second)
	nrRoundsToPropagateMultiShard := 5
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	nftMetaData.Quantity += quantityToAdd
	nft.CheckNftData(
		t,
		nodeInDifferentShard.OwnAccount.Address,
		nodeInDifferentShard.OwnAccount.Address,
		nodes,
		[]byte(tokenIdentifier),
		nftMetaData,
		1,
	)

	time.Sleep(time.Second)
	nrRoundsToPropagateMultiShard = 5
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	nft.CheckNftData(
		t,
		nodeInDifferentShard.OwnAccount.Address,
		nodeInDifferentShard.OwnAccount.Address,
		nodes,
		[]byte(tokenIdentifier),
		nftMetaData,
		1,
	)

	// transfer
	quantityToTransfer := int64(4)
	quantityToTransferArg := hex.EncodeToString(big.NewInt(quantityToTransfer).Bytes())
	txData = []byte(core.BuiltInFunctionDCDTNFTTransfer + "@" + hex.EncodeToString([]byte(tokenIdentifier)) +
		"@" + nonceArg + "@" + quantityToTransferArg + "@" + hex.EncodeToString(nodes[0].OwnAccount.Address))
	integrationTests.CreateAndSendTransaction(
		nodeInDifferentShard,
		nodes,
		big.NewInt(0),
		nodeInDifferentShard.OwnAccount.Address,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)

	time.Sleep(time.Second)
	nrRoundsToPropagateMultiShard = 11
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	nftMetaData.Quantity = initialQuantity + quantityToAdd - quantityToTransfer
	nft.CheckNftData(
		t,
		nodeInDifferentShard.OwnAccount.Address,
		nodeInDifferentShard.OwnAccount.Address,
		nodes,
		[]byte(tokenIdentifier),
		nftMetaData,
		1,
	)

	nftMetaData.Quantity = quantityToTransfer
	nft.CheckNftData(
		t,
		nodeInDifferentShard.OwnAccount.Address,
		nodes[0].OwnAccount.Address,
		nodes,
		[]byte(tokenIdentifier),
		nftMetaData,
		1,
	)
}

func TestDCDTSFTWithEnhancedTransferRole(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	nodesPerShard := 2
	numMetachainNodes := 2
	numOfShards := 3

	nodes := integrationTests.CreateNodes(
		numOfShards,
		nodesPerShard,
		numMetachainNodes,
	)

	leaders := make([]*integrationTests.TestProcessorNode, numOfShards+1)
	for i := 0; i < numOfShards; i++ {
		leaders[i] = nodes[i*nodesPerShard]
	}
	leaders[numOfShards] = nodes[numOfShards*nodesPerShard]

	integrationTests.DisplayAndStartNodes(nodes)

	defer func() {
		for _, n := range nodes {
			_ = n.MainMessenger.Close()
		}
	}()

	initialVal := big.NewInt(10000000000)
	integrationTests.MintAllNodes(nodes, initialVal)

	round := uint64(0)
	nonce := uint64(0)
	round = integrationTests.IncrementAndPrintRound(round)
	nonce++

	roles := [][]byte{
		[]byte(core.DCDTRoleNFTCreate),
		[]byte(core.DCDTRoleNFTAddQuantity),
		[]byte(core.DCDTRoleNFTBurn),
		[]byte(core.DCDTRoleTransfer),
	}

	tokenIssuer := nodes[0]
	initialQuantity := int64(5)
	tokenIdentifier, nftMetaData := nft.PrepareNFTWithRoles(
		t,
		nodes,
		leaders,
		tokenIssuer,
		&round,
		&nonce,
		core.SemiFungibleDCDT,
		initialQuantity,
		roles,
	)

	// increase quantity
	nonceArg := hex.EncodeToString(big.NewInt(0).SetUint64(1).Bytes())
	quantityToAdd := int64(100)
	quantityToAddArg := hex.EncodeToString(big.NewInt(quantityToAdd).Bytes())
	txData := []byte(core.BuiltInFunctionDCDTNFTAddQuantity + "@" + hex.EncodeToString([]byte(tokenIdentifier)) +
		"@" + nonceArg + "@" + quantityToAddArg)
	integrationTests.CreateAndSendTransaction(
		tokenIssuer,
		nodes,
		big.NewInt(0),
		tokenIssuer.OwnAccount.Address,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)

	time.Sleep(time.Second)
	nrRoundsToPropagateMultiShard := 2
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	nftMetaData.Quantity += quantityToAdd
	nft.CheckNftData(
		t,
		tokenIssuer.OwnAccount.Address,
		tokenIssuer.OwnAccount.Address,
		nodes,
		[]byte(tokenIdentifier),
		nftMetaData,
		1,
	)

	// transfer
	quantityToTransfer := int64(5)
	quantityToTransferArg := hex.EncodeToString(big.NewInt(quantityToTransfer).Bytes())

	for _, node := range nodes[1:] {
		txData = []byte(core.BuiltInFunctionDCDTNFTTransfer + "@" + hex.EncodeToString([]byte(tokenIdentifier)) +
			"@" + nonceArg + "@" + quantityToTransferArg + "@" + hex.EncodeToString(node.OwnAccount.Address))
		integrationTests.CreateAndSendTransaction(
			tokenIssuer,
			nodes,
			big.NewInt(0),
			tokenIssuer.OwnAccount.Address,
			string(txData),
			integrationTests.AdditionalGasLimit,
		)
	}

	time.Sleep(time.Second)
	nrRoundsToPropagateMultiShard = 12
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	nftMetaData.Quantity = initialQuantity + quantityToAdd - int64(len(nodes)-1)*quantityToTransfer
	nft.CheckNftData(
		t,
		tokenIssuer.OwnAccount.Address,
		tokenIssuer.OwnAccount.Address,
		nodes,
		[]byte(tokenIdentifier),
		nftMetaData,
		1,
	)

	nftMetaData.Quantity = quantityToTransfer
	for _, node := range nodes[1:] {
		nft.CheckNftData(
			t,
			tokenIssuer.OwnAccount.Address,
			node.OwnAccount.Address,
			nodes,
			[]byte(tokenIdentifier),
			nftMetaData,
			1,
		)
	}

	// every account will send back the tokens
	for _, node := range nodes[1:] {
		txData = []byte(core.BuiltInFunctionDCDTNFTTransfer + "@" + hex.EncodeToString([]byte(tokenIdentifier)) +
			"@" + nonceArg + "@" + quantityToTransferArg + "@" + hex.EncodeToString(tokenIssuer.OwnAccount.Address))
		integrationTests.CreateAndSendTransaction(
			node,
			nodes,
			big.NewInt(0),
			node.OwnAccount.Address,
			string(txData),
			integrationTests.AdditionalGasLimit,
		)
	}

	time.Sleep(time.Second)
	nrRoundsToPropagateMultiShard = 12
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	nftMetaData.Quantity = initialQuantity + quantityToAdd
	nft.CheckNftData(
		t,
		tokenIssuer.OwnAccount.Address,
		tokenIssuer.OwnAccount.Address,
		nodes,
		[]byte(tokenIdentifier),
		nftMetaData,
		1,
	)

	nftMetaData.Quantity = 0
	for _, node := range nodes[1:] {
		nft.CheckNftData(
			t,
			tokenIssuer.OwnAccount.Address,
			node.OwnAccount.Address,
			nodes,
			[]byte(tokenIdentifier),
			nftMetaData,
			1,
		)
	}
}

func TestNFTTransferCreateAndSetRolesInShard(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	testNFTTransferCreateRoleAndStop(t, 1)
}

func TestNFTTransferCreateAndSetRolesCrossShard(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	testNFTTransferCreateRoleAndStop(t, 2)
}

func testNFTTransferCreateRoleAndStop(t *testing.T, numOfShards int) {
	nodes, leaders := dcdt.CreateNodesAndPrepareBalances(numOfShards)

	defer func() {
		for _, n := range nodes {
			n.Close()
		}
	}()

	initialVal := big.NewInt(10000000000)
	integrationTests.MintAllNodes(nodes, initialVal)

	round := uint64(0)
	nonce := uint64(0)
	round = integrationTests.IncrementAndPrintRound(round)
	nonce++

	roles := [][]byte{
		[]byte(core.DCDTRoleNFTCreate),
	}

	nftCreator := nodes[0]
	initialQuantity := int64(1)
	tokenIdentifier, nftMetaData := nft.PrepareNFTWithRoles(
		t,
		nodes,
		leaders,
		nftCreator,
		&round,
		&nonce,
		core.SemiFungibleDCDT,
		initialQuantity,
		roles,
	)

	nftCreatorShId := nftCreator.ShardCoordinator.ComputeId(nftCreator.OwnAccount.Address)
	nextNftCreator := nodes[1]
	for _, node := range nodes {
		if node.ShardCoordinator.ComputeId(node.OwnAccount.Address) != nftCreatorShId {
			nextNftCreator = node
			break
		}
	}

	// transferNFTCreateRole
	txData := []byte("transferNFTCreateRole" + "@" + hex.EncodeToString([]byte(tokenIdentifier)) +
		"@" + hex.EncodeToString(nftCreator.OwnAccount.Address) +
		"@" + hex.EncodeToString(nextNftCreator.OwnAccount.Address))
	integrationTests.CreateAndSendTransaction(
		nftCreator,
		nodes,
		big.NewInt(0),
		vm.DCDTSCAddress,
		string(txData),
		integrationTests.AdditionalGasLimit+core.MinMetaTxExtraGasCost,
	)

	time.Sleep(time.Second)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, 15, nonce, round)
	time.Sleep(time.Second)

	// stopNFTCreate
	txData = []byte("stopNFTCreate" + "@" + hex.EncodeToString([]byte(tokenIdentifier)))
	integrationTests.CreateAndSendTransaction(
		nftCreator,
		nodes,
		big.NewInt(0),
		vm.DCDTSCAddress,
		string(txData),
		integrationTests.AdditionalGasLimit+core.MinMetaTxExtraGasCost,
	)

	time.Sleep(time.Second)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, 2, nonce, round)
	time.Sleep(time.Second)

	// setCreateRole
	txData = []byte("setSpecialRole" + "@" + hex.EncodeToString([]byte(tokenIdentifier)) +
		"@" + hex.EncodeToString(nftCreator.OwnAccount.Address) +
		"@" + hex.EncodeToString([]byte(core.DCDTRoleNFTCreate)))
	integrationTests.CreateAndSendTransaction(
		nftCreator,
		nodes,
		big.NewInt(0),
		vm.DCDTSCAddress,
		string(txData),
		integrationTests.AdditionalGasLimit+core.MinMetaTxExtraGasCost,
	)

	time.Sleep(time.Second)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, 20, nonce, round)
	time.Sleep(time.Second)

	newNFTMetaData := nft.NftArguments{
		Name:       []byte("NEW"),
		Quantity:   1,
		Royalties:  9000,
		Hash:       []byte("NEW"),
		Attributes: []byte("NEW"),
		URI:        [][]byte{[]byte("NEW")},
	}

	nft.CreateNFT(
		[]byte(tokenIdentifier),
		nftCreator,
		nodes,
		&newNFTMetaData,
	)

	time.Sleep(time.Second)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, 2, nonce, round)
	time.Sleep(time.Second)

	// we check that old data remains on NONCE 1 - as creation must return failure
	nft.CheckNftData(
		t,
		nftCreator.OwnAccount.Address,
		nftCreator.OwnAccount.Address,
		nodes,
		[]byte(tokenIdentifier),
		nftMetaData,
		1,
	)
}
