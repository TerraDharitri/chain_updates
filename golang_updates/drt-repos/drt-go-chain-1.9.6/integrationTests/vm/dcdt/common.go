package dcdt

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data/dcdt"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/integrationTests"
	testVm "github.com/TerraDharitri/drt-go-chain/integrationTests/vm"
	"github.com/TerraDharitri/drt-go-chain/integrationTests/vm/wasm"
	"github.com/TerraDharitri/drt-go-chain/process"
	vmFactory "github.com/TerraDharitri/drt-go-chain/process/factory"
	"github.com/TerraDharitri/drt-go-chain/state"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/txDataBuilder"
	"github.com/TerraDharitri/drt-go-chain/vm"
)

// GetDCDTTokenData -
func GetDCDTTokenData(
	t *testing.T,
	address []byte,
	nodes []*integrationTests.TestProcessorNode,
	tickerID []byte,
	nonce uint64,
) *dcdt.DCDigitalToken {
	accShardID := nodes[0].ShardCoordinator.ComputeId(address)
	for _, node := range nodes {
		if node.ShardCoordinator.SelfId() != accShardID {
			continue
		}

		dcdtData, err := node.BlockchainHook.GetDCDTToken(address, tickerID, nonce)
		require.Nil(t, err)
		return dcdtData
	}

	return &dcdt.DCDigitalToken{Value: big.NewInt(0)}
}

// GetUserAccountWithAddress -
func GetUserAccountWithAddress(
	t *testing.T,
	address []byte,
	nodes []*integrationTests.TestProcessorNode,
) state.UserAccountHandler {
	for _, node := range nodes {
		accShardId := node.ShardCoordinator.ComputeId(address)

		for _, helperNode := range nodes {
			if helperNode.ShardCoordinator.SelfId() == accShardId {
				acc, err := helperNode.AccntState.LoadAccount(address)
				require.Nil(t, err)
				return acc.(state.UserAccountHandler)
			}
		}
	}

	return nil
}

// SetRoles -
func SetRoles(nodes []*integrationTests.TestProcessorNode, addrForRole []byte, tokenIdentifier []byte, roles [][]byte) {
	tokenIssuer := nodes[0]
	SetRolesWithSenderAccount(nodes, tokenIssuer.OwnAccount, addrForRole, tokenIdentifier, roles)
}

// SetRolesWithSenderAccount -
func SetRolesWithSenderAccount(nodes []*integrationTests.TestProcessorNode, issuerAccount *integrationTests.TestWalletAccount, addrForRole []byte, tokenIdentifier []byte, roles [][]byte) {
	tokenIssuer := nodes[0]

	txData := "setSpecialRole" +
		"@" + hex.EncodeToString(tokenIdentifier) +
		"@" + hex.EncodeToString(addrForRole)

	for _, role := range roles {
		txData += "@" + hex.EncodeToString(role)
	}

	integrationTests.CreateAndSendTransactionWithSenderAccount(tokenIssuer, nodes, big.NewInt(0), issuerAccount, vm.DCDTSCAddress, txData, core.MinMetaTxExtraGasCost)
}

// DeployNonPayableSmartContract -
func DeployNonPayableSmartContract(
	t *testing.T,
	nodes []*integrationTests.TestProcessorNode,
	leaders []*integrationTests.TestProcessorNode,
	nonce *uint64,
	round *uint64,
	fileName string,
) []byte {
	return DeployNonPayableSmartContractFromNode(t, nodes, 0, leaders, nonce, round, fileName)
}

// DeployNonPayableSmartContractFromNode -
func DeployNonPayableSmartContractFromNode(
	t *testing.T,
	nodes []*integrationTests.TestProcessorNode,
	idDeployer int,
	leaders []*integrationTests.TestProcessorNode,
	nonce *uint64,
	round *uint64,
	fileName string,
) []byte {
	scCode := wasm.GetSCCode(fileName)
	scAddress, _ := nodes[idDeployer].BlockchainHook.NewAddress(nodes[idDeployer].OwnAccount.Address, nodes[idDeployer].OwnAccount.Nonce, vmFactory.WasmVirtualMachine)

	integrationTests.CreateAndSendTransaction(
		nodes[idDeployer],
		nodes,
		big.NewInt(0),
		testVm.CreateEmptyAddress(),
		wasm.CreateDeployTxDataNonPayable(scCode),
		integrationTests.AdditionalGasLimit,
	)

	*nonce, *round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, 4, *nonce, *round)

	scShardID := nodes[0].ShardCoordinator.ComputeId(scAddress)
	for _, node := range nodes {
		if node.ShardCoordinator.SelfId() != scShardID {
			continue
		}
		_, err := node.AccntState.GetExistingAccount(scAddress)
		require.Nil(t, err)
	}

	return scAddress
}

// CheckAddressHasTokens - Works for both fungible and non-fungible, according to nonce
func CheckAddressHasTokens(
	t *testing.T,
	address []byte,
	nodes []*integrationTests.TestProcessorNode,
	tickerID []byte,
	nonce int64,
	value int64,
) {
	nonceAsBigInt := big.NewInt(nonce)
	valueAsBigInt := big.NewInt(value)

	dcdtData := GetDCDTTokenData(t, address, nodes, tickerID, uint64(nonce))

	if dcdtData == nil {
		dcdtData = &dcdt.DCDigitalToken{
			Value: big.NewInt(0),
		}
	}
	if dcdtData.Value == nil {
		dcdtData.Value = big.NewInt(0)
	}

	if valueAsBigInt.Cmp(dcdtData.Value) != 0 {
		require.Fail(t, fmt.Sprintf("dcdt NFT balance difference. Token %s, nonce %s, expected %s, but got %s",
			tickerID, nonceAsBigInt.String(), valueAsBigInt.String(), dcdtData.Value.String()))
	}
}

// CreateNodesAndPrepareBalances -
func CreateNodesAndPrepareBalances(numOfShards int) ([]*integrationTests.TestProcessorNode, []*integrationTests.TestProcessorNode) {
	enableEpochs := config.EnableEpochs{
		OptimizeGasUsedInCrossMiniBlocksEnableEpoch: integrationTests.UnreachableEpoch,
		ScheduledMiniBlocksEnableEpoch:              integrationTests.UnreachableEpoch,
		MiniBlockPartialExecutionEnableEpoch:        integrationTests.UnreachableEpoch,
		AndromedaEnableEpoch:                        integrationTests.UnreachableEpoch,
	}
	roundsConfig := testscommon.GetDefaultRoundsConfig()
	return CreateNodesAndPrepareBalancesWithEpochsAndRoundsConfig(
		numOfShards,
		enableEpochs,
		roundsConfig,
	)
}

// CreateNodesAndPrepareBalancesWithEpochsAndRoundsConfig -
func CreateNodesAndPrepareBalancesWithEpochsAndRoundsConfig(
	numOfShards int,
	enableEpochs config.EnableEpochs,
	roundsConfig config.RoundConfig,
) ([]*integrationTests.TestProcessorNode, []*integrationTests.TestProcessorNode) {
	nodesPerShard := 1
	numMetachainNodes := 1

	nodes := integrationTests.CreateNodesWithEnableEpochsAndVmConfigWithRoundsConfig(
		numOfShards,
		nodesPerShard,
		numMetachainNodes,
		enableEpochs,
		roundsConfig,
		&config.VirtualMachineConfig{
			WasmVMVersions: []config.WasmVMVersionByEpoch{
				{StartEpoch: 0, Version: "*"},
			},
			TransferAndExecuteByUserAddresses: []string{"drt1fpkcgel4gcmh8zqqdt043yfcn5tyx8373kg6q2qmkxzu4dqamc0snh8ehx"},
		},
	)

	leaders := make([]*integrationTests.TestProcessorNode, numOfShards+1)
	for i := 0; i < numOfShards; i++ {
		leaders[i] = nodes[i*nodesPerShard]
	}
	leaders[numOfShards] = nodes[numOfShards*nodesPerShard]
	integrationTests.DisplayAndStartNodes(nodes)

	return nodes, leaders
}

// IssueNFT -
func IssueNFT(nodes []*integrationTests.TestProcessorNode, dcdtType string, ticker string) {
	tokenName := "token"
	issuePrice := big.NewInt(1000)

	tokenIssuer := nodes[0]

	txData := txDataBuilder.NewBuilder()

	issueFunc := "issueNonFungible"
	if dcdtType == core.SemiFungibleDCDT {
		issueFunc = "issueSemiFungible"
	}
	txData.Clear().Func(issueFunc).Str(tokenName).Str(ticker)
	txData.CanFreeze(false).CanWipe(false).CanPause(false).CanTransferNFTCreateRole(true)

	integrationTests.CreateAndSendTransaction(tokenIssuer, nodes, issuePrice, vm.DCDTSCAddress, txData.ToString(), core.MinMetaTxExtraGasCost)
}

// IssueTestToken -
func IssueTestToken(nodes []*integrationTests.TestProcessorNode, initialSupply int64, ticker string) {
	issueTestToken(nodes, initialSupply, ticker, core.MinMetaTxExtraGasCost)
}

// IssueTestTokenWithIssuerAccount -
func IssueTestTokenWithIssuerAccount(nodes []*integrationTests.TestProcessorNode, issuerAccount *integrationTests.TestWalletAccount, initialSupply int64, ticker string) {
	issueTestTokenWithIssuerAccount(nodes, issuerAccount, initialSupply, ticker, core.MinMetaTxExtraGasCost)
}

// IssueTestTokenWithCustomGas -
func IssueTestTokenWithCustomGas(nodes []*integrationTests.TestProcessorNode, initialSupply int64, ticker string, gas uint64) {
	issueTestToken(nodes, initialSupply, ticker, gas)
}

// IssueTestTokenWithSpecialRoles -
func IssueTestTokenWithSpecialRoles(nodes []*integrationTests.TestProcessorNode, initialSupply int64, ticker string) {
	issueTestTokenWithSpecialRoles(nodes, initialSupply, ticker, core.MinMetaTxExtraGasCost)
}

func issueTestToken(nodes []*integrationTests.TestProcessorNode, initialSupply int64, ticker string, gas uint64) {
	tokenIssuer := nodes[0]
	issueTestTokenWithIssuerAccount(nodes, tokenIssuer.OwnAccount, initialSupply, ticker, gas)
}

func issueTestTokenWithIssuerAccount(nodes []*integrationTests.TestProcessorNode, issuerAccount *integrationTests.TestWalletAccount, initialSupply int64, ticker string, gas uint64) {
	tokenName := "token"
	issuePrice := big.NewInt(1000)

	tokenIssuer := nodes[0]
	txData := txDataBuilder.NewBuilder()
	txData.Clear().IssueDCDT(tokenName, ticker, initialSupply, 6)
	txData.CanFreeze(true).CanWipe(true).CanPause(true).CanMint(true).CanBurn(true)

	integrationTests.CreateAndSendTransactionWithSenderAccount(tokenIssuer, nodes, issuePrice, issuerAccount, vm.DCDTSCAddress, txData.ToString(), gas)
}

func issueTestTokenWithSpecialRoles(nodes []*integrationTests.TestProcessorNode, initialSupply int64, ticker string, gas uint64) {
	tokenName := "token"
	issuePrice := big.NewInt(1000)

	tokenIssuer := nodes[0]
	txData := txDataBuilder.NewBuilder()
	txData.Clear().IssueDCDT(tokenName, ticker, initialSupply, 6)
	txData.CanFreeze(true).CanWipe(true).CanPause(true).CanMint(true).CanBurn(true).CanAddSpecialRoles(true)

	integrationTests.CreateAndSendTransaction(tokenIssuer, nodes, issuePrice, vm.DCDTSCAddress, txData.ToString(), gas)
}

// CheckNumCallBacks -
func CheckNumCallBacks(
	t *testing.T,
	address []byte,
	nodes []*integrationTests.TestProcessorNode,
	expectedNumCallbacks int,
) {

	contractID := nodes[0].ShardCoordinator.ComputeId(address)
	for _, node := range nodes {
		if node.ShardCoordinator.SelfId() != contractID {
			continue
		}

		scQuery := &process.SCQuery{
			ScAddress:  address,
			FuncName:   "callback_args",
			CallerAddr: address,
			CallValue:  big.NewInt(0),
			Arguments:  [][]byte{},
		}
		vmOutput, _, err := node.SCQueryService.ExecuteQuery(scQuery)
		require.Nil(t, err)
		require.NotNil(t, vmOutput)
		require.Equal(t, vmOutput.ReturnCode, vmcommon.Ok)
		require.Equal(t, expectedNumCallbacks, len(vmOutput.ReturnData))
	}
}

// CheckForwarderRawSavedCallbackArgs -
func CheckForwarderRawSavedCallbackArgs(
	t *testing.T,
	address []byte,
	nodes []*integrationTests.TestProcessorNode,
	callbackIndex int,
	expectedResultCode vmcommon.ReturnCode,
	expectedArguments [][]byte) {

	contractID := nodes[0].ShardCoordinator.ComputeId(address)
	for _, node := range nodes {
		if node.ShardCoordinator.SelfId() != contractID {
			continue
		}

		scQueryArgs := &process.SCQuery{
			ScAddress:  address,
			FuncName:   "callback_args_at_index",
			CallerAddr: address,
			CallValue:  big.NewInt(0),
			Arguments: [][]byte{
				{byte(callbackIndex)},
			},
		}
		vmOutputArgs, _, err := node.SCQueryService.ExecuteQuery(scQueryArgs)
		require.Nil(t, err)
		require.Equal(t, vmcommon.Ok, vmOutputArgs.ReturnCode)
		require.GreaterOrEqual(t, len(vmOutputArgs.ReturnData), 1)
		if expectedResultCode == vmcommon.Ok {
			require.Equal(t, []byte{0x0}, vmOutputArgs.ReturnData[0])
			require.Equal(t, expectedArguments, vmOutputArgs.ReturnData[1:])
		} else {
			require.Equal(t, []byte{byte(expectedResultCode)}, vmOutputArgs.ReturnData[0])
		}
	}
}

// ForwarderRawSavedPaymentInfo contains token data to be checked in the forwarder-raw contract.
type ForwarderRawSavedPaymentInfo struct {
	TokenId string
	Nonce   uint64
	Payment *big.Int
}

// CheckForwarderRawSavedCallbackPayments -
func CheckForwarderRawSavedCallbackPayments(
	t *testing.T,
	address []byte,
	nodes []*integrationTests.TestProcessorNode,
	expectedPayments []*ForwarderRawSavedPaymentInfo) {

	scQueryPayment := &process.SCQuery{
		ScAddress:  address,
		FuncName:   "callback_payments_triples",
		CallerAddr: address,
		CallValue:  big.NewInt(0),
		Arguments:  [][]byte{},
	}

	contractID := nodes[0].ShardCoordinator.ComputeId(address)
	for _, node := range nodes {
		if node.ShardCoordinator.SelfId() != contractID {
			continue
		}
		vmOutputPayment, _, err := node.SCQueryService.ExecuteQuery(scQueryPayment)
		require.Nil(t, err)
		require.Equal(t, vmcommon.Ok, vmOutputPayment.ReturnCode)

		require.Equal(t, len(expectedPayments)*3, len(vmOutputPayment.ReturnData))
		for i, expectedPayment := range expectedPayments {
			require.Equal(t, []byte(expectedPayment.TokenId), vmOutputPayment.ReturnData[3*i])
			require.Equal(t, big.NewInt(0).SetUint64(expectedPayment.Nonce).Bytes(), vmOutputPayment.ReturnData[3*i+1])
			require.Equal(t, expectedPayment.Payment.Bytes(), vmOutputPayment.ReturnData[3*i+2])
		}
	}
}

// PrepareFungibleTokensWithLocalBurnAndMint -
func PrepareFungibleTokensWithLocalBurnAndMint(
	t *testing.T,
	nodes []*integrationTests.TestProcessorNode,
	addressWithRoles []byte,
	leaders []*integrationTests.TestProcessorNode,
	round *uint64,
	nonce *uint64,
) string {
	return PrepareFungibleTokensWithLocalBurnAndMintWithIssuerAccount(
		t,
		nodes,
		nodes[0].OwnAccount,
		addressWithRoles,
		leaders,
		round,
		nonce)
}

// PrepareFungibleTokensWithLocalBurnAndMintWithIssuerAccount -
func PrepareFungibleTokensWithLocalBurnAndMintWithIssuerAccount(
	t *testing.T,
	nodes []*integrationTests.TestProcessorNode,
	issuerAccount *integrationTests.TestWalletAccount,
	addressWithRoles []byte,
	leaders []*integrationTests.TestProcessorNode,
	round *uint64,
	nonce *uint64,
) string {
	IssueTestTokenWithIssuerAccount(nodes, issuerAccount, 100, "TKN")

	time.Sleep(time.Second)
	nrRoundsToPropagateMultiShard := 5
	*nonce, *round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, *nonce, *round)
	time.Sleep(time.Second)

	tokenIdentifier := string(integrationTests.GetTokenIdentifier(nodes, []byte("TKN")))

	SetRolesWithSenderAccount(nodes, issuerAccount, addressWithRoles, []byte(tokenIdentifier), [][]byte{[]byte(core.DCDTRoleLocalMint), []byte(core.DCDTRoleLocalBurn)})

	time.Sleep(time.Second)
	nrRoundsToPropagateMultiShard = 5
	*nonce, *round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, *nonce, *round)
	time.Sleep(time.Second)

	return tokenIdentifier
}
