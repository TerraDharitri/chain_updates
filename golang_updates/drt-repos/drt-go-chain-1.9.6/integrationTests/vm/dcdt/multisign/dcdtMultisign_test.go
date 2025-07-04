package multisign

import (
	"encoding/hex"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	logger "github.com/TerraDharitri/drt-go-chain-logger"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/integrationTests"
	"github.com/TerraDharitri/drt-go-chain/integrationTests/vm/dcdt"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/vm"
)

var vmType = []byte{5, 0}
var emptyAddress = make([]byte, 32)
var log = logger.GetOrCreate("integrationtests/vm/dcdt")

func TestDCDTTransferWithMultisig(t *testing.T) {
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

	initialVal := big.NewInt(10000000000000)
	integrationTests.MintAllNodes(nodes, initialVal)

	round := uint64(0)
	nonce := uint64(0)
	round = integrationTests.IncrementAndPrintRound(round)
	nonce++

	multisignContractAddress := deployMultisig(t, nodes, 0, 1, 2)

	time.Sleep(time.Second)
	numRoundsToPropagateIntraShard := 2
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, numRoundsToPropagateIntraShard, nonce, round)
	time.Sleep(time.Second)

	// ----- issue DCDT token
	initalSupply := big.NewInt(10000000000)
	ticker := "TCK"
	proposeIssueTokenAndTransferFunds(nodes, multisignContractAddress, initalSupply, 0, ticker)

	time.Sleep(time.Second)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, numRoundsToPropagateIntraShard, nonce, round)
	time.Sleep(time.Second)

	actionID := getActionID(t, nodes, multisignContractAddress)
	log.Info("got action ID", "action ID", actionID)

	boardMembersSignActionID(nodes, multisignContractAddress, actionID, 1, 2)

	time.Sleep(time.Second)
	numRoundsToPropagateCrossShard := 10
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, numRoundsToPropagateCrossShard, nonce, round)
	time.Sleep(time.Second)

	performActionID(nodes, multisignContractAddress, actionID, 0)

	time.Sleep(time.Second)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, numRoundsToPropagateCrossShard, nonce, round)
	time.Sleep(time.Second)

	tokenIdentifier := integrationTests.GetTokenIdentifier(nodes, []byte(ticker))
	dcdt.CheckAddressHasTokens(t, multisignContractAddress, nodes, tokenIdentifier, 0, initalSupply.Int64())

	checkCallBackWasSaved(t, nodes, multisignContractAddress)

	// ----- transfer DCDT token
	destinationAddress, _ := integrationTests.TestAddressPubkeyConverter.Decode("drt1j25xk97yf820rgdp3mj5scavhjkn6tjyn0t63pmv5qyjj7wxlcfqa9aqd2")
	transferValue := big.NewInt(10)
	proposeTransferToken(nodes, multisignContractAddress, transferValue, 0, destinationAddress, tokenIdentifier)

	time.Sleep(time.Second)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, numRoundsToPropagateIntraShard, nonce, round)
	time.Sleep(time.Second)

	actionID = getActionID(t, nodes, multisignContractAddress)
	log.Info("got action ID", "action ID", actionID)

	boardMembersSignActionID(nodes, multisignContractAddress, actionID, 1, 2)

	time.Sleep(time.Second)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, numRoundsToPropagateCrossShard, nonce, round)
	time.Sleep(time.Second)

	performActionID(nodes, multisignContractAddress, actionID, 0)

	time.Sleep(time.Second)
	_, _ = integrationTests.WaitOperationToBeDone(t, leaders, nodes, numRoundsToPropagateCrossShard, nonce, round)
	time.Sleep(time.Second)

	expectedBalance := big.NewInt(0).Set(initalSupply)
	expectedBalance.Sub(expectedBalance, transferValue)
	dcdt.CheckAddressHasTokens(t, multisignContractAddress, nodes, tokenIdentifier, 0, expectedBalance.Int64())
	dcdt.CheckAddressHasTokens(t, destinationAddress, nodes, tokenIdentifier, 0, transferValue.Int64())
}

func checkCallBackWasSaved(t *testing.T, nodes []*integrationTests.TestProcessorNode, contract []byte) {
	contractID := nodes[0].ShardCoordinator.ComputeId(contract)
	for _, node := range nodes {
		if node.ShardCoordinator.SelfId() != contractID {
			continue
		}

		scQuery := &process.SCQuery{
			ScAddress:  contract,
			FuncName:   "callback_log",
			CallerAddr: contract,
			CallValue:  big.NewInt(0),
			Arguments:  [][]byte{},
		}
		vmOutput, _, err := node.SCQueryService.ExecuteQuery(scQuery)
		assert.Nil(t, err)
		assert.Equal(t, vmOutput.ReturnCode, vmcommon.Ok)
		assert.Equal(t, 1, len(vmOutput.ReturnData))
	}
}

func deployMultisig(t *testing.T, nodes []*integrationTests.TestProcessorNode, ownerIdx int, proposersIndexes ...int) []byte {
	codeMetaData := &vmcommon.CodeMetadata{
		Payable:     true,
		Upgradeable: false,
		Readable:    true,
	}

	contractBytes, err := os.ReadFile("../testdata/multisig-callback.wasm")
	require.Nil(t, err)
	proposers := make([]string, 0, len(proposersIndexes)+1)
	proposers = append(proposers, hex.EncodeToString(nodes[ownerIdx].OwnAccount.Address))
	for _, proposerIdx := range proposersIndexes {
		walletAddressAsHex := hex.EncodeToString(nodes[proposerIdx].OwnAccount.Address)
		proposers = append(proposers, walletAddressAsHex)
	}

	parameters := []string{
		hex.EncodeToString(contractBytes),
		hex.EncodeToString(vmType),
		hex.EncodeToString(codeMetaData.ToBytes()),
		hex.EncodeToString(big.NewInt(int64(len(proposersIndexes))).Bytes()),
	}
	parameters = append(parameters, proposers...)

	txData := strings.Join(
		parameters,
		"@",
	)

	multisigContractAddress, err := nodes[ownerIdx].BlockchainHook.NewAddress(
		nodes[ownerIdx].OwnAccount.Address,
		nodes[ownerIdx].OwnAccount.Nonce,
		vmType,
	)
	require.Nil(t, err)

	encodedMultisigContractAddress, err := integrationTests.TestAddressPubkeyConverter.Encode(multisigContractAddress)
	require.Nil(t, err)

	log.Info("multisign contract", "address", encodedMultisigContractAddress)
	integrationTests.CreateAndSendTransaction(nodes[ownerIdx], nodes, big.NewInt(0), emptyAddress, txData, 1000000)

	return multisigContractAddress
}

func proposeIssueTokenAndTransferFunds(
	nodes []*integrationTests.TestProcessorNode,
	multisignContractAddress []byte,
	initalSupply *big.Int,
	ownerIdx int,
	ticker string,
) {
	tokenName := []byte("token")
	issuePrice := big.NewInt(1000)
	multisigParams := []string{
		"proposeSCCall",
		hex.EncodeToString(vm.DCDTSCAddress),
		hex.EncodeToString(issuePrice.Bytes()),
	}

	dcdtParams := []string{
		hex.EncodeToString([]byte("issue")),
		hex.EncodeToString(tokenName),
		hex.EncodeToString([]byte(ticker)),
		hex.EncodeToString(initalSupply.Bytes()),
		hex.EncodeToString([]byte{6}),
	}

	hexEncodedTrue := hex.EncodeToString([]byte("true"))
	tokenPropertiesParams := []string{
		hex.EncodeToString([]byte("canFreeze")),
		hexEncodedTrue,
		hex.EncodeToString([]byte("canWipe")),
		hexEncodedTrue,
		hex.EncodeToString([]byte("canPause")),
		hexEncodedTrue,
		hex.EncodeToString([]byte("canMint")),
		hexEncodedTrue,
		hex.EncodeToString([]byte("canBurn")),
		hexEncodedTrue,
	}

	params := append(multisigParams, dcdtParams...)
	params = append(params, tokenPropertiesParams...)
	txData := strings.Join(params, "@")

	integrationTests.CreateAndSendTransaction(nodes[ownerIdx], nodes, big.NewInt(1000000), multisignContractAddress, "deposit", 1000000)
	integrationTests.CreateAndSendTransaction(nodes[ownerIdx], nodes, big.NewInt(0), multisignContractAddress, txData, 1000000)
}

func getActionID(t *testing.T, nodes []*integrationTests.TestProcessorNode, multisignContractAddress []byte) []byte {
	node := getSameShardNode(nodes, multisignContractAddress)
	accnt, _ := node.AccntState.LoadAccount(multisignContractAddress)
	_ = accnt

	query := &process.SCQuery{
		ScAddress:  multisignContractAddress,
		FuncName:   "getPendingActionFullInfo",
		CallerAddr: make([]byte, 0),
		CallValue:  big.NewInt(0),
		Arguments:  make([][]byte, 0),
	}

	vmOutput, _, err := node.SCQueryService.ExecuteQuery(query)
	require.Nil(t, err)
	require.Equal(t, 1, len(vmOutput.ReturnData))

	actionFullInfo := vmOutput.ReturnData[0]

	return actionFullInfo[:4]
}

func getSameShardNode(nodes []*integrationTests.TestProcessorNode, address []byte) *integrationTests.TestProcessorNode {
	shId := nodes[0].ShardCoordinator.ComputeId(address)
	for _, n := range nodes {
		if n.ShardCoordinator.SelfId() == shId {
			return n
		}
	}

	return nil
}

func boardMembersSignActionID(
	nodes []*integrationTests.TestProcessorNode,
	multisignContractAddress []byte,
	actionID []byte,
	signersIndexes ...int,
) {
	for _, index := range signersIndexes {
		node := nodes[index]
		params := []string{
			"sign",
			hex.EncodeToString(actionID),
		}

		txData := strings.Join(params, "@")
		integrationTests.CreateAndSendTransaction(node, nodes, big.NewInt(0), multisignContractAddress, txData, 1000000)
	}
}

func performActionID(
	nodes []*integrationTests.TestProcessorNode,
	multisignContractAddress []byte,
	actionID []byte,
	nodeIndex int,
) {
	node := nodes[nodeIndex]
	params := []string{
		"performAction",
		hex.EncodeToString(actionID),
	}

	txData := strings.Join(params, "@")
	integrationTests.CreateAndSendTransaction(node, nodes, big.NewInt(0), multisignContractAddress, txData, 1500000)
}

func proposeTransferToken(
	nodes []*integrationTests.TestProcessorNode,
	multisignContractAddress []byte,
	transferValue *big.Int,
	ownerIdx int,
	destinationAddress []byte,
	tokenID []byte,
) {
	multisigParams := []string{
		"proposeSCCall",
		hex.EncodeToString(destinationAddress),
		"00",
	}

	dcdtParams := []string{
		hex.EncodeToString([]byte("DCDTTransfer")),
		hex.EncodeToString(tokenID),
		hex.EncodeToString(transferValue.Bytes()),
	}

	params := append(multisigParams, dcdtParams...)
	txData := strings.Join(params, "@")

	integrationTests.CreateAndSendTransaction(nodes[ownerIdx], nodes, big.NewInt(0), multisignContractAddress, txData, 1000000)
}
