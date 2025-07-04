package roles

import (
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/integrationTests"
	"github.com/TerraDharitri/drt-go-chain/integrationTests/vm/dcdt"
	"github.com/TerraDharitri/drt-go-chain/testscommon/txDataBuilder"
	"github.com/TerraDharitri/drt-go-chain/vm"
)

// Test scenario
// 1 - issue an DCDT token
// 2 - set special roles for the owner of the token (local burn and local mint)
// 3 - do a local burn - should work
// 4 - do a local mint - should work
func TestDCDTRolesIssueAndTransactionsOnMultiShardEnvironment(t *testing.T) {
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

	// ------- send token issue

	initialSupply := big.NewInt(10000000000)
	dcdt.IssueTestToken(nodes, initialSupply.Int64(), "FTT")
	tokenIssuer := nodes[0]

	time.Sleep(time.Second)
	nrRoundsToPropagateMultiShard := 6
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	tokenIdentifier := string(integrationTests.GetTokenIdentifier(nodes, []byte("FTT")))

	// ----- set special role
	setRole(nodes, nodes[0].OwnAccount.Address, []byte(tokenIdentifier), []byte(core.DCDTRoleLocalMint))
	setRole(nodes, nodes[0].OwnAccount.Address, []byte(tokenIdentifier), []byte(core.DCDTRoleLocalBurn))

	time.Sleep(time.Second)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	dcdt.CheckAddressHasTokens(t, tokenIssuer.OwnAccount.Address, nodes, []byte(tokenIdentifier), 0, initialSupply.Int64())

	// mint local new tokens
	txData := []byte(core.BuiltInFunctionDCDTLocalMint + "@" + hex.EncodeToString([]byte(tokenIdentifier)) +
		"@" + hex.EncodeToString(big.NewInt(500).Bytes()))
	integrationTests.CreateAndSendTransaction(
		nodes[0],
		nodes,
		big.NewInt(0),
		nodes[0].OwnAccount.Address,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)

	time.Sleep(time.Second)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	// check balance ofter local mint
	dcdt.CheckAddressHasTokens(t, tokenIssuer.OwnAccount.Address, nodes, []byte(tokenIdentifier), 0, big.NewInt(10000000500).Int64())

	// burn local  tokens
	txData = []byte(core.BuiltInFunctionDCDTLocalBurn + "@" + hex.EncodeToString([]byte(tokenIdentifier)) +
		"@" + hex.EncodeToString(big.NewInt(200).Bytes()))
	integrationTests.CreateAndSendTransaction(
		nodes[0],
		nodes,
		big.NewInt(0),
		nodes[0].OwnAccount.Address,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)

	time.Sleep(time.Second)
	_, _ = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	// check balance ofter local mint
	dcdt.CheckAddressHasTokens(t, tokenIssuer.OwnAccount.Address, nodes, []byte(tokenIdentifier), 0, big.NewInt(10000000300).Int64())
}

// Test scenario
// 1 - issue an DCDT token
// 2 - set special role for the owner of the token (local mint)
// 3 - unset special role (local mint)
// 3 - do a local mint - DCDT balance should not change
// 4 - set special role (local burn)
// 5 - do a local burn - should work
func TestDCDTRolesSetRolesAndUnsetRolesIssueAndTransactionsOnMultiShardEnvironment(t *testing.T) {
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

	// ------- send token issue

	initialSupply := big.NewInt(10000000000)
	dcdt.IssueTestToken(nodes, initialSupply.Int64(), "FTT")
	tokenIssuer := nodes[0]

	time.Sleep(time.Second)
	nrRoundsToPropagateMultiShard := 12
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	tokenIdentifier := string(integrationTests.GetTokenIdentifier(nodes, []byte("FTT")))

	// ----- set special role
	setRole(nodes, nodes[0].OwnAccount.Address, []byte(tokenIdentifier), []byte(core.DCDTRoleLocalMint))

	time.Sleep(time.Second)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	// unset special role
	unsetRole(nodes, nodes[0].OwnAccount.Address, []byte(tokenIdentifier), []byte(core.DCDTRoleLocalMint))

	time.Sleep(time.Second)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	dcdt.CheckAddressHasTokens(t, tokenIssuer.OwnAccount.Address, nodes, []byte(tokenIdentifier), 0, initialSupply.Int64())
	dcdt.CheckAddressHasTokens(t, tokenIssuer.OwnAccount.Address, nodes, []byte(tokenIdentifier), 0, initialSupply.Int64())

	// mint local new tokens
	txData := []byte(core.BuiltInFunctionDCDTLocalMint + "@" + hex.EncodeToString([]byte(tokenIdentifier)) +
		"@" + hex.EncodeToString(big.NewInt(500).Bytes()))
	integrationTests.CreateAndSendTransaction(
		nodes[0],
		nodes,
		big.NewInt(0),
		nodes[0].OwnAccount.Address,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)

	time.Sleep(time.Second)
	nrRoundsToPropagateMultiShard = 7
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	// check balance ofter local mint
	dcdt.CheckAddressHasTokens(t, tokenIssuer.OwnAccount.Address, nodes, []byte(tokenIdentifier), 0, big.NewInt(10000000000).Int64())

	setRole(nodes, nodes[0].OwnAccount.Address, []byte(tokenIdentifier), []byte(core.DCDTRoleLocalBurn))
	time.Sleep(time.Second)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	// burn local  tokens
	txData = []byte(core.BuiltInFunctionDCDTLocalBurn + "@" + hex.EncodeToString([]byte(tokenIdentifier)) +
		"@" + hex.EncodeToString(big.NewInt(200).Bytes()))
	integrationTests.CreateAndSendTransaction(
		nodes[0],
		nodes,
		big.NewInt(0),
		nodes[0].OwnAccount.Address,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)

	time.Sleep(time.Second)
	_, _ = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	// check balance ofter local mint
	dcdt.CheckAddressHasTokens(t, tokenIssuer.OwnAccount.Address, nodes, []byte(tokenIdentifier), 0, big.NewInt(9999999800).Int64())
}

func setRole(nodes []*integrationTests.TestProcessorNode, addrForRole []byte, tokenIdentifier []byte, roles []byte) {
	tokenIssuer := nodes[0]

	txData := "setSpecialRole" +
		"@" + hex.EncodeToString(tokenIdentifier) +
		"@" + hex.EncodeToString(addrForRole) +
		"@" + hex.EncodeToString(roles)
	integrationTests.CreateAndSendTransaction(tokenIssuer, nodes, big.NewInt(0), vm.DCDTSCAddress, txData, core.MinMetaTxExtraGasCost)
}

func unsetRole(nodes []*integrationTests.TestProcessorNode, addrForRole []byte, tokenIdentifier []byte, roles []byte) {
	tokenIssuer := nodes[0]

	txData := "unSetSpecialRole" +
		"@" + hex.EncodeToString(tokenIdentifier) +
		"@" + hex.EncodeToString(addrForRole) +
		"@" + hex.EncodeToString(roles)
	integrationTests.CreateAndSendTransaction(tokenIssuer, nodes, big.NewInt(0), vm.DCDTSCAddress, txData, core.MinMetaTxExtraGasCost)
}

func TestDCDTMintTransferAndExecute(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	numOfShards := 2
	nodesPerShard := 1
	numMetachainNodes := 1

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

	scAddress := dcdt.DeployNonPayableSmartContract(t, nodes, leaders, &nonce, &round, "../testdata/rewa-dcdt-swap.wasm")

	// issue DCDT by calling exec on dest context on child contract
	ticker := "DSN"
	name := "DisplayName"
	issueCost := big.NewInt(1000)
	txIssueData := txDataBuilder.NewBuilder()
	txIssueData.Func("issueWrappedRewa").
		Str(name).
		Str(ticker)
	integrationTests.CreateAndSendTransaction(
		nodes[0],
		nodes,
		issueCost,
		scAddress,
		txIssueData.ToString(),
		integrationTests.AdditionalGasLimit,
	)

	time.Sleep(time.Second)
	nrRoundsToPropagateMultiShard := 15
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	tokenIdentifier := integrationTests.GetTokenIdentifier(nodes, []byte(ticker))
	integrationTests.CreateAndSendTransaction(
		nodes[0],
		nodes,
		big.NewInt(0),
		scAddress,
		"setLocalRoles",
		integrationTests.AdditionalGasLimit,
	)
	time.Sleep(time.Second)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	valueToWrap := big.NewInt(1000)
	for _, n := range nodes {
		txData := []byte("wrapRewa")
		integrationTests.CreateAndSendTransaction(
			n,
			nodes,
			valueToWrap,
			scAddress,
			string(txData),
			integrationTests.AdditionalGasLimit,
		)
	}

	time.Sleep(time.Second)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	for i, n := range nodes {
		if i == 0 {
			continue
		}
		dcdt.CheckAddressHasTokens(t, n.OwnAccount.Address, nodes, tokenIdentifier, 0, valueToWrap.Int64())
	}

	for _, n := range nodes {
		txUnWrap := txDataBuilder.NewBuilder()
		txUnWrap.Func(core.BuiltInFunctionDCDTTransfer).Str(string(tokenIdentifier)).BigInt(valueToWrap).Str("unwrapRewa")
		integrationTests.CreateAndSendTransaction(
			n,
			nodes,
			big.NewInt(0),
			scAddress,
			txUnWrap.ToString(),
			integrationTests.AdditionalGasLimit,
		)
	}
	time.Sleep(time.Second)

	_, _ = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	userAccount := dcdt.GetUserAccountWithAddress(t, scAddress, nodes)
	require.Equal(t, userAccount.GetBalance(), big.NewInt(0))
}

func TestDCDTLocalBurnFromAnyoneOfThisToken(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	numOfShards := 2
	nodesPerShard := 2
	numMetachainNodes := 2

	enableEpochs := config.EnableEpochs{
		ScheduledMiniBlocksEnableEpoch: integrationTests.UnreachableEpoch,
		AndromedaEnableEpoch:           integrationTests.UnreachableEpoch,
	}
	nodes := integrationTests.CreateNodesWithEnableEpochs(
		numOfShards,
		nodesPerShard,
		numMetachainNodes,
		enableEpochs,
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

	// send token issue
	ticker := "ALC"
	issuePrice := big.NewInt(1000)
	initialSupply := int64(10000000000)
	tokenIssuer := nodes[0]
	txData := txDataBuilder.NewBuilder()

	txData.Clear().IssueDCDT("aliceToken", ticker, initialSupply, 6)
	txData.CanFreeze(true).CanWipe(true).CanPause(true).CanMint(true).CanBurn(false)
	integrationTests.CreateAndSendTransaction(tokenIssuer, nodes, issuePrice, vm.DCDTSCAddress, txData.ToString(), core.MinMetaTxExtraGasCost)

	time.Sleep(time.Second)
	nrRoundsToPropagateMultiShard := 12
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	tokenIdentifier := string(integrationTests.GetTokenIdentifier(nodes, []byte(ticker)))

	dcdt.CheckAddressHasTokens(t, tokenIssuer.OwnAccount.Address, nodes, []byte(tokenIdentifier), 0, initialSupply)

	// send tx to other nodes
	valueToSend := int64(100)
	for _, node := range nodes[1:] {
		txData.Clear().TransferDCDT(tokenIdentifier, valueToSend)
		integrationTests.CreateAndSendTransaction(tokenIssuer, nodes, big.NewInt(0), node.OwnAccount.Address, txData.ToString(), integrationTests.AdditionalGasLimit)
	}

	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	finalSupply := initialSupply
	for _, node := range nodes[1:] {
		dcdt.CheckAddressHasTokens(t, node.OwnAccount.Address, nodes, []byte(tokenIdentifier), 0, valueToSend)
		finalSupply = finalSupply - valueToSend
		txData.Clear().LocalBurnDCDT(tokenIdentifier, valueToSend)
		integrationTests.CreateAndSendTransaction(node, nodes, big.NewInt(0), node.OwnAccount.Address, txData.ToString(), integrationTests.AdditionalGasLimit)
	}

	dcdt.CheckAddressHasTokens(t, tokenIssuer.OwnAccount.Address, nodes, []byte(tokenIdentifier), 0, finalSupply)
	txData.Clear().LocalBurnDCDT(tokenIdentifier, finalSupply)
	integrationTests.CreateAndSendTransaction(tokenIssuer, nodes, big.NewInt(0), tokenIssuer.OwnAccount.Address, txData.ToString(), integrationTests.AdditionalGasLimit)

	_, _ = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	for _, node := range nodes {
		dcdt.CheckAddressHasTokens(t, node.OwnAccount.Address, nodes, []byte(tokenIdentifier), 0, 0)
	}
}

func TestDCDTWithTransferRoleCrossShardShouldWork(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	numOfShards := 2
	nodesPerShard := 2
	numMetachainNodes := 2

	enableEpochs := config.EnableEpochs{
		ScheduledMiniBlocksEnableEpoch: integrationTests.UnreachableEpoch,
		AndromedaEnableEpoch:           integrationTests.UnreachableEpoch,
	}
	nodes := integrationTests.CreateNodesWithEnableEpochs(
		numOfShards,
		nodesPerShard,
		numMetachainNodes,
		enableEpochs,
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

	// send token issue
	ticker := "ALC"
	issuePrice := big.NewInt(1000)
	initialSupply := int64(10000000000)
	tokenIssuer := nodes[0]
	txData := txDataBuilder.NewBuilder()

	txData.Clear().IssueDCDT("aliceToken", ticker, initialSupply, 6)
	txData.CanFreeze(true).CanWipe(true).CanPause(true).CanMint(true).CanBurn(false)
	integrationTests.CreateAndSendTransaction(tokenIssuer, nodes, issuePrice, vm.DCDTSCAddress, txData.ToString(), core.MinMetaTxExtraGasCost)

	time.Sleep(time.Second)
	nrRoundsToPropagateMultiShard := 12
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	tokenIdentifier := string(integrationTests.GetTokenIdentifier(nodes, []byte(ticker)))
	setRole(nodes, tokenIssuer.OwnAccount.Address, []byte(tokenIdentifier), []byte(core.DCDTRoleTransfer))
	dcdt.CheckAddressHasTokens(t, tokenIssuer.OwnAccount.Address, nodes, []byte(tokenIdentifier), 0, initialSupply)

	time.Sleep(time.Second)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, 2, nonce, round)
	time.Sleep(time.Second)

	// send tx to other nodes
	valueToSend := int64(100)
	for _, node := range nodes[1:] {
		txData.Clear().TransferDCDT(tokenIdentifier, valueToSend)
		integrationTests.CreateAndSendTransaction(tokenIssuer, nodes, big.NewInt(0), node.OwnAccount.Address, txData.ToString(), integrationTests.AdditionalGasLimit)
	}

	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	// send value back to the initial node
	for _, node := range nodes[1:] {
		dcdt.CheckAddressHasTokens(t, node.OwnAccount.Address, nodes, []byte(tokenIdentifier), 0, valueToSend)
		txData.Clear().TransferDCDT(tokenIdentifier, valueToSend)
		integrationTests.CreateAndSendTransaction(node, nodes, big.NewInt(0), tokenIssuer.OwnAccount.Address, txData.ToString(), integrationTests.AdditionalGasLimit)
	}

	_, _ = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)
	time.Sleep(time.Second)

	for _, node := range nodes[1:] {
		dcdt.CheckAddressHasTokens(t, node.OwnAccount.Address, nodes, []byte(tokenIdentifier), 0, 0)
	}
	dcdt.CheckAddressHasTokens(t, tokenIssuer.OwnAccount.Address, nodes, []byte(tokenIdentifier), 0, initialSupply)
}
