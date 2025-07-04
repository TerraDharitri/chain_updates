package softfork

import (
	"encoding/hex"
	"math/big"
	"os"
	"testing"
	"time"

	crypto "github.com/TerraDharitri/drt-go-chain-crypto"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	logger "github.com/TerraDharitri/drt-go-chain-logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/integrationTests"
	"github.com/TerraDharitri/drt-go-chain/process/factory"
	"github.com/TerraDharitri/drt-go-chain/state"
)

var log = logger.GetOrCreate("integrationtests/singleshard/block/softfork")

func TestScDeploy(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	deployEnableEpoch := uint32(1)
	relayedTxEnableEpoch := uint32(0)
	penalizedTooMuchGasEnableEpoch := uint32(0)
	roundsPerEpoch := uint64(10)
	scProcessorV2EnableEpoch := integrationTests.UnreachableEpoch

	enableEpochs := integrationTests.CreateEnableEpochsConfig()
	enableEpochs.SCDeployEnableEpoch = deployEnableEpoch
	enableEpochs.RelayedTransactionsEnableEpoch = relayedTxEnableEpoch
	enableEpochs.PenalizedTooMuchGasEnableEpoch = penalizedTooMuchGasEnableEpoch
	enableEpochs.SCProcessorV2EnableEpoch = scProcessorV2EnableEpoch
	enableEpochs.StakingV4Step1EnableEpoch = integrationTests.StakingV4Step1EnableEpoch
	enableEpochs.StakingV4Step2EnableEpoch = integrationTests.StakingV4Step2EnableEpoch
	enableEpochs.StakingV4Step3EnableEpoch = integrationTests.StakingV4Step3EnableEpoch

	shardNode := integrationTests.NewTestProcessorNode(integrationTests.ArgTestProcessorNode{
		MaxShards:            1,
		NodeShardId:          0,
		TxSignPrivKeyShardId: 0,
		EpochsConfig:         &enableEpochs,
	})
	shardNode.EpochStartTrigger.SetRoundsPerEpoch(roundsPerEpoch)

	metaNode := integrationTests.NewTestProcessorNode(integrationTests.ArgTestProcessorNode{
		MaxShards:            1,
		NodeShardId:          core.MetachainShardId,
		TxSignPrivKeyShardId: 0,
		EpochsConfig:         &enableEpochs,
	})
	metaNode.EpochStartTrigger.SetRoundsPerEpoch(roundsPerEpoch)

	nodes := []*integrationTests.TestProcessorNode{
		shardNode,
		metaNode,
	}
	connectableNodes := make([]integrationTests.Connectable, 0)
	for _, n := range nodes {
		connectableNodes = append(connectableNodes, n)
	}
	integrationTests.ConnectNodes(connectableNodes)

	leaders := []*integrationTests.TestProcessorNode{nodes[0], nodes[1]}

	defer func() {
		for _, n := range nodes {
			n.Close()
		}
	}()

	log.Info("delaying for nodes p2p bootstrap...")
	time.Sleep(integrationTests.P2pBootstrapDelay)

	round := uint64(1)
	nonce := uint64(1)
	numRounds := roundsPerEpoch + 5

	integrationTests.CreateMintingForSenders(nodes, 0, []crypto.PrivateKey{shardNode.OwnAccount.SkTxSign}, big.NewInt(1000000000))

	accnt, _ := shardNode.AccntState.GetExistingAccount(shardNode.OwnAccount.Address)
	userAccnt := accnt.(state.UserAccountHandler)
	balance := userAccnt.GetBalance()
	log.Info("balance", "value", balance.String())

	deployedFailedAddress := deploySc(t, nodes)

	for i := uint64(0); i < numRounds; i++ {
		integrationTests.UpdateRound(nodes, round)
		integrationTests.ProposeBlock(nodes, leaders, round, nonce)
		round = integrationTests.IncrementAndPrintRound(round)
		nonce++

		time.Sleep(integrationTests.StepDelay)
	}

	encodedDeployFailedAddr, err := integrationTests.TestAddressPubkeyConverter.Encode(deployedFailedAddress)
	assert.Nil(t, err)
	log.Info("resulted sc address (failed)", "address", encodedDeployFailedAddr)
	assert.False(t, scAccountExists(shardNode, deployedFailedAddress))

	deploySucceeded := deploySc(t, nodes)
	for i := uint64(0); i < 5; i++ {
		integrationTests.UpdateRound(nodes, round)
		integrationTests.ProposeBlock(nodes, leaders, round, nonce)
		round = integrationTests.IncrementAndPrintRound(round)
		nonce++

		time.Sleep(integrationTests.StepDelay)
	}

	encodedDeploySucceededAddr, err := integrationTests.TestAddressPubkeyConverter.Encode(deploySucceeded)
	assert.Nil(t, err)
	log.Info("resulted sc address (success)", "address", encodedDeploySucceededAddr)
	assert.True(t, scAccountExists(shardNode, deploySucceeded))
}

func deploySc(t *testing.T, nodes []*integrationTests.TestProcessorNode) []byte {
	scCode, err := os.ReadFile("./testdata/answer.wasm")
	require.Nil(t, err)

	node := nodes[0]
	scAddress, err := node.BlockchainHook.NewAddress(node.OwnAccount.Address, node.OwnAccount.Nonce, factory.WasmVirtualMachine)
	require.Nil(t, err)

	integrationTests.DeployScTx(nodes, 0, hex.EncodeToString(scCode), factory.WasmVirtualMachine, "001000000000")

	return scAddress
}

func scAccountExists(node *integrationTests.TestProcessorNode, address []byte) bool {
	accnt, _ := node.AccntState.GetExistingAccount(address)

	return !check.IfNil(accnt)
}
