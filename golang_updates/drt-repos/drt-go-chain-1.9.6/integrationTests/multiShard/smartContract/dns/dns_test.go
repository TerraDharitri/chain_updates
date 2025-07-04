package dns

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/data/api"
	"github.com/TerraDharitri/drt-go-chain-core/hashing/keccak"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/genesis"
	"github.com/TerraDharitri/drt-go-chain/integrationTests"
	"github.com/TerraDharitri/drt-go-chain/integrationTests/multiShard/relayedTx"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/state"
)

func TestSCCallingDNSUserNames(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	nodes, players, leaders := prepareNodesAndPlayers()
	defer func() {
		for _, n := range nodes {
			n.Close()
		}
	}()

	dnsRegisterValue, sortedDNSAddresses := getDNSContractsData(nodes[0])

	round := uint64(0)
	nonce := uint64(0)
	round = integrationTests.IncrementAndPrintRound(round)
	nonce++

	userNames := sendRegisterUserNameTxForPlayers(players, nodes, sortedDNSAddresses, dnsRegisterValue)

	time.Sleep(time.Second)
	nrRoundsToPropagateMultiShard := 25
	_, _ = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)

	checkUserNamesAreSetCorrectly(t, players, nodes, userNames, sortedDNSAddresses)
}

func TestSCCallingDNSUserNamesTwice(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	nodes, players, leaders := prepareNodesAndPlayers()
	defer func() {
		for _, n := range nodes {
			n.Close()
		}
	}()

	dnsRegisterValue, sortedDNSAddresses := getDNSContractsData(nodes[0])

	round := uint64(0)
	nonce := uint64(0)
	round = integrationTests.IncrementAndPrintRound(round)
	nonce++

	userNames := sendRegisterUserNameTxForPlayers(players, nodes, sortedDNSAddresses, dnsRegisterValue)

	time.Sleep(time.Second)
	nrRoundsToPropagateMultiShard := 15
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)

	newUserNames := sendRegisterUserNameTxForPlayers(players, nodes, sortedDNSAddresses, dnsRegisterValue)

	time.Sleep(time.Second)
	_, _ = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)

	checkUserNamesAreSetCorrectly(t, players, nodes, userNames, sortedDNSAddresses)
	checkUserNamesAreDeleted(t, nodes, newUserNames, sortedDNSAddresses)
}

func TestDNSandRelayedTxNormal(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	nodes, players, leaders := prepareNodesAndPlayers()
	defer func() {
		for _, n := range nodes {
			n.Close()
		}
	}()

	relayer := createAndMintRelayer(nodes)
	dnsRegisterValue, sortedDNSAddresses := getDNSContractsData(nodes[0])

	round := uint64(0)
	nonce := uint64(0)
	round = integrationTests.IncrementAndPrintRound(round)
	nonce++

	userNames := sendRegisterUserNameAsRelayedTx(relayer, players, nodes, sortedDNSAddresses, dnsRegisterValue)

	time.Sleep(time.Second)
	nrRoundsToPropagateMultiShard := 30
	_, _ = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, nonce, round)

	checkUserNamesAreSetCorrectly(t, players, nodes, userNames, sortedDNSAddresses)
}

func createAndMintRelayer(nodes []*integrationTests.TestProcessorNode) *integrationTests.TestWalletAccount {
	relayer := integrationTests.CreateTestWalletAccount(nodes[0].ShardCoordinator, 0)

	initialVal := big.NewInt(10000000000000)
	initialVal.Mul(initialVal, initialVal)
	integrationTests.MintAllPlayers(nodes, []*integrationTests.TestWalletAccount{relayer}, initialVal)
	return relayer
}

func prepareNodesAndPlayers() ([]*integrationTests.TestProcessorNode, []*integrationTests.TestWalletAccount, []*integrationTests.TestProcessorNode) {
	numOfShards := 2
	nodesPerShard := 1
	numMetachainNodes := 1

	genesisFile := "smartcontracts.json"
	enableEpochsConfig := integrationTests.GetDefaultEnableEpochsConfig()
	enableEpochsConfig.StakingV2EnableEpoch = integrationTests.UnreachableEpoch
	enableEpochsConfig.ChangeUsernameEnableEpoch = integrationTests.UnreachableEpoch
	nodes, _ := integrationTests.CreateNodesWithFullGenesisCustomEnableEpochs(
		numOfShards,
		nodesPerShard,
		numMetachainNodes,
		genesisFile,
		enableEpochsConfig,
	)

	for _, node := range nodes {
		node.EconomicsData.SetMaxGasLimitPerBlock(1500000000, 0)
	}

	leaders := make([]*integrationTests.TestProcessorNode, numOfShards+1)
	for i := 0; i < numOfShards; i++ {
		leaders[i] = nodes[i*nodesPerShard]
	}
	leaders[numOfShards] = nodes[numOfShards*nodesPerShard]

	integrationTests.DisplayAndStartNodes(nodes)

	numPlayers := 6
	players := make([]*integrationTests.TestWalletAccount, numPlayers)
	for i := 0; i < numPlayers; i++ {
		players[i] = integrationTests.CreateTestWalletAccount(nodes[0].ShardCoordinator, 0)
	}

	initialVal := big.NewInt(10000000000000)
	initialVal.Mul(initialVal, initialVal)
	fmt.Printf("Initial minted sum: %s\n", initialVal.String())
	integrationTests.MintAllNodes(nodes, initialVal)
	integrationTests.MintAllPlayers(nodes, players, initialVal)

	return nodes, players, leaders
}

func getDNSContractsData(node *integrationTests.TestProcessorNode) (*big.Int, []string) {
	dnsRegisterValue := big.NewInt(100)
	genesisSCs := node.SmartContractParser.InitialSmartContracts()
	for _, genesisSC := range genesisSCs {
		if genesisSC.GetType() == genesis.DNSType {
			decodedValue, _ := hex.DecodeString(genesisSC.GetInitParameters())
			dnsRegisterValue.SetBytes(decodedValue)
			break
		}
	}

	mapDNSAddresses, _ := node.SmartContractParser.GetDeployedSCAddresses(genesis.DNSType)
	sortedDNSAddresses := make([]string, 0, len(mapDNSAddresses))
	for address := range mapDNSAddresses {
		sortedDNSAddresses = append(sortedDNSAddresses, address)
	}
	sort.Slice(sortedDNSAddresses, func(i, j int) bool {
		return sortedDNSAddresses[i][31] < sortedDNSAddresses[j][31]
	})

	return dnsRegisterValue, sortedDNSAddresses
}

func sendRegisterUserNameAsRelayedTx(
	relayer *integrationTests.TestWalletAccount,
	players []*integrationTests.TestWalletAccount,
	nodes []*integrationTests.TestProcessorNode,
	sortedDNSAddresses []string,
	dnsRegisterValue *big.Int,
) []string {

	userNames := make([]string, len(players))
	gasLimit := uint64(2000000)
	for i, player := range players {
		userName := generateNewUserName()
		scAddress := selectDNSAddressFromUserName(sortedDNSAddresses, userName)
		_, _ = relayedTx.CreateAndSendRelayedAndUserTx(nodes, relayer, player, []byte(scAddress), dnsRegisterValue,
			gasLimit, []byte("register@"+hex.EncodeToString([]byte(userName))))
		userNames[i] = userName
	}

	return userNames
}

func sendRegisterUserNameTxForPlayers(
	players []*integrationTests.TestWalletAccount,
	nodes []*integrationTests.TestProcessorNode,
	sortedDNSAddresses []string,
	dnsRegisterValue *big.Int,
) []string {

	gasLimit := uint64(2000000)
	userNames := make([]string, len(players))
	for i, player := range players {
		userName := generateNewUserName()
		scAddress := selectDNSAddressFromUserName(sortedDNSAddresses, userName)
		integrationTests.PlayerSendsTransaction(
			nodes,
			player,
			[]byte(scAddress),
			dnsRegisterValue,
			"register@"+hex.EncodeToString([]byte(userName)),
			gasLimit,
		)
		userNames[i] = userName
	}
	return userNames
}

func checkUserNamesAreSetCorrectly(
	t *testing.T,
	players []*integrationTests.TestWalletAccount,
	nodes []*integrationTests.TestProcessorNode,
	userNames []string,
	sortedDNSAddresses []string,
) {
	for i, player := range players {
		playerShID := nodes[0].ShardCoordinator.ComputeId(player.Address)
		for _, node := range nodes {
			if node.ShardCoordinator.SelfId() != playerShID {
				continue
			}

			acnt, _ := node.AccntState.GetExistingAccount(player.Address)
			userAcc, _ := acnt.(state.UserAccountHandler)

			assert.Equal(t, userNames[i], string(userAcc.GetUserName()))

			bech32c := integrationTests.TestAddressPubkeyConverter
			playerAddress, err := bech32c.Encode(player.Address)
			require.NoError(t, err)
			usernameReportedByNode, _, err := node.Node.GetUsername(playerAddress, api.AccountQueryOptions{})
			require.NoError(t, err)
			require.Equal(t, userNames[i], usernameReportedByNode)
		}

		dnsAddress := selectDNSAddressFromUserName(sortedDNSAddresses, userNames[i])
		scQuery := &process.SCQuery{
			CallerAddr: player.Address,
			ScAddress:  []byte(dnsAddress),
			FuncName:   "resolve",
			Arguments:  [][]byte{[]byte(userNames[i])},
		}

		dnsSHId := nodes[0].ShardCoordinator.ComputeId([]byte(dnsAddress))
		for _, node := range nodes {
			if node.ShardCoordinator.SelfId() != dnsSHId {
				continue
			}

			vmOutput, _, _ := node.SCQueryService.ExecuteQuery(scQuery)

			require.NotNil(t, vmOutput)
			require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
			require.Equal(t, 1, len(vmOutput.ReturnData))
			assert.True(t, bytes.Equal(player.Address, vmOutput.ReturnData[0]))
		}
	}
}

func checkUserNamesAreDeleted(
	t *testing.T,
	nodes []*integrationTests.TestProcessorNode,
	userNames []string,
	sortedDNSAddresses []string,
) {
	for _, userName := range userNames {
		dnsAddress := selectDNSAddressFromUserName(sortedDNSAddresses, userName)

		dnsSHId := nodes[0].ShardCoordinator.ComputeId([]byte(dnsAddress))
		for _, node := range nodes {
			if node.ShardCoordinator.SelfId() != dnsSHId {
				continue
			}

			acnt, _ := node.AccntState.GetExistingAccount([]byte(dnsAddress))
			dnsAcc, _ := acnt.(state.UserAccountHandler)

			keyFromTrie := "value_state" + string(keccak.NewKeccak().Compute(userName))
			value, _, err := dnsAcc.RetrieveValue([]byte(keyFromTrie))
			assert.Nil(t, err)
			assert.Nil(t, value)
		}
	}
}

func selectDNSAddressFromUserName(sortedDNSAddresses []string, userName string) string {
	hashedAddr := keccak.NewKeccak().Compute(userName)
	return sortedDNSAddresses[hashedAddr[31]]
}

func generateNewUserName() string {
	return RandStringBytes(10) + ".numbat"
}

const letterBytes = "abcdefghijklmnopqrstuvwxyz"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
