package multitransfer

import (
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain/integrationTests"
	testVm "github.com/TerraDharitri/drt-go-chain/integrationTests/vm"
	"github.com/TerraDharitri/drt-go-chain/integrationTests/vm/dcdt"
	"github.com/TerraDharitri/drt-go-chain/integrationTests/vm/wasm"
	vmFactory "github.com/TerraDharitri/drt-go-chain/process/factory"
	"github.com/TerraDharitri/drt-go-chain/testscommon/txDataBuilder"
	"github.com/TerraDharitri/drt-go-chain/vm"
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/stretchr/testify/require"
)

const numRoundsCrossShard = 15
const numRoundsSameShard = 1

// DcdtTransfer -
type DcdtTransfer struct {
	TokenIdentifier string
	Nonce           int64
	Amount          int64
}

// IssueFungibleToken -
func IssueFungibleToken(
	t *testing.T,
	net *integrationTests.TestNetwork,
	issuerNode *integrationTests.TestProcessorNode,
	ticker string,
	initialSupply int64,
) string {
	return IssueFungibleTokenWithIssuerAddress(t, net, issuerNode, issuerNode.OwnAccount, ticker, initialSupply)
}

// IssueFungibleTokenWithIssuerAddress -
func IssueFungibleTokenWithIssuerAddress(
	t *testing.T,
	net *integrationTests.TestNetwork,
	issuerNode *integrationTests.TestProcessorNode,
	issuerAccount *integrationTests.TestWalletAccount,
	ticker string,
	initialSupply int64,
) string {
	tokenName := "token"
	issuePrice := big.NewInt(1000)
	txData := txDataBuilder.NewBuilder()
	txData.IssueDCDT(tokenName, ticker, initialSupply, 6)
	txData.CanFreeze(true).CanWipe(true).CanPause(true).CanMint(true).CanBurn(true)

	integrationTests.CreateAndSendTransactionWithSenderAccount(
		issuerNode,
		net.Nodes,
		issuePrice,
		issuerAccount,
		vm.DCDTSCAddress,
		txData.ToString(), core.MinMetaTxExtraGasCost)
	WaitForOperationCompletion(net, numRoundsCrossShard)

	tokenIdentifier := integrationTests.GetTokenIdentifier(net.Nodes, []byte(ticker))

	dcdt.CheckAddressHasTokens(t, issuerAccount.Address, net.Nodes,
		tokenIdentifier, 0, initialSupply)

	return string(tokenIdentifier)
}

// IssueNft -
func IssueNft(
	net *integrationTests.TestNetwork,
	issuerNode *integrationTests.TestProcessorNode,
	ticker string,
	semiFungible bool,
) string {
	return IssueNftWithIssuerAddress(net, issuerNode, issuerNode.OwnAccount, ticker, semiFungible)
}

// IssueNftWithIssuerAddress -
func IssueNftWithIssuerAddress(
	net *integrationTests.TestNetwork,
	issuerNode *integrationTests.TestProcessorNode,
	issuerAccount *integrationTests.TestWalletAccount,
	ticker string,
	semiFungible bool,
) string {

	tokenName := "token"
	issuePrice := big.NewInt(1000)

	txData := txDataBuilder.NewBuilder()

	issueFunc := "issueNonFungible"
	if semiFungible {
		issueFunc = "issueSemiFungible"
	}

	txData.Func(issueFunc).Str(tokenName).Str(ticker)
	txData.CanFreeze(false).CanWipe(false).CanPause(false).CanTransferNFTCreateRole(true)

	integrationTests.CreateAndSendTransactionWithSenderAccount(
		issuerNode,
		net.Nodes,
		issuePrice,
		issuerAccount,
		vm.DCDTSCAddress,
		txData.ToString(),
		core.MinMetaTxExtraGasCost)
	WaitForOperationCompletion(net, numRoundsCrossShard)

	tokenIdentifier := string(integrationTests.GetTokenIdentifier(net.Nodes, []byte(ticker)))

	roles := [][]byte{
		[]byte(core.DCDTRoleNFTCreate),
	}
	if semiFungible {
		roles = append(roles, []byte(core.DCDTRoleNFTAddQuantity))
	}

	SetLocalRoles(net, issuerNode, issuerAccount, tokenIdentifier, roles)

	return tokenIdentifier
}

// SetLocalRoles -
func SetLocalRoles(
	net *integrationTests.TestNetwork,
	issuerNode *integrationTests.TestProcessorNode,
	issuerAccount *integrationTests.TestWalletAccount,
	tokenIdentifier string,
	roles [][]byte,
) {
	txData := "setSpecialRole" +
		"@" + hex.EncodeToString([]byte(tokenIdentifier)) +
		"@" + hex.EncodeToString(issuerAccount.Address)

	for _, role := range roles {
		txData += "@" + hex.EncodeToString(role)
	}

	integrationTests.CreateAndSendTransactionWithSenderAccount(
		issuerNode,
		net.Nodes,
		big.NewInt(0),
		issuerAccount,
		vm.DCDTSCAddress,
		txData,
		core.MinMetaTxExtraGasCost)
	WaitForOperationCompletion(net, numRoundsCrossShard)
}

// CreateSFT -
func CreateSFT(
	t *testing.T,
	net *integrationTests.TestNetwork,
	issuerNode *integrationTests.TestProcessorNode,
	issuerAccount *integrationTests.TestWalletAccount,
	tokenIdentifier string,
	createdTokenNonce int64,
	initialSupply int64,
) {
	tokenName := "token"
	royalties := big.NewInt(0)
	hash := "someHash"
	attributes := "cool nft"
	uri := "www.my-cool-nfts.com"

	txData := txDataBuilder.NewBuilder()
	txData.Func(core.BuiltInFunctionDCDTNFTCreate)
	txData.Str(tokenIdentifier)
	txData.Int64(initialSupply)
	txData.Str(tokenName)
	txData.BigInt(royalties)
	txData.Str(hash)
	txData.Str(attributes)
	txData.Str(uri)

	integrationTests.CreateAndSendTransactionWithSenderAccount(
		issuerNode,
		net.Nodes,
		big.NewInt(0),
		issuerAccount,
		issuerAccount.Address,
		txData.ToString(),
		integrationTests.AdditionalGasLimit)
	WaitForOperationCompletion(net, numRoundsSameShard)

	dcdt.CheckAddressHasTokens(t, issuerAccount.Address, net.Nodes,
		[]byte(tokenIdentifier), createdTokenNonce, initialSupply)
}

// CreateNFT -
func CreateNFT(
	t *testing.T,
	net *integrationTests.TestNetwork,
	issuerNode *integrationTests.TestProcessorNode,
	issuerAccount *integrationTests.TestWalletAccount,
	tokenIdentifier string,
	createdTokenNonce int64,
) {

	CreateSFT(t, net, issuerNode, issuerAccount, tokenIdentifier, createdTokenNonce, 1)
}

// BuildDcdtMultiTransferTxData -
func BuildDcdtMultiTransferTxData(
	receiverAddress []byte,
	transfers []*DcdtTransfer,
	endpointName string,
	arguments ...[]byte,
) string {

	nrTransfers := len(transfers)

	txData := txDataBuilder.NewBuilder()
	txData.Func(core.BuiltInFunctionMultiDCDTNFTTransfer)
	txData.Bytes(receiverAddress)
	txData.Int(nrTransfers)

	for _, transfer := range transfers {
		txData.Str(transfer.TokenIdentifier)
		txData.Int64(transfer.Nonce)
		txData.Int64(transfer.Amount)
	}

	if len(endpointName) > 0 {
		txData.Str(endpointName)

		for _, arg := range arguments {
			txData.Bytes(arg)
		}
	}

	return txData.ToString()
}

// WaitForOperationCompletion -
func WaitForOperationCompletion(net *integrationTests.TestNetwork, roundsToWait int) {
	time.Sleep(time.Second)
	net.Steps(roundsToWait)
}

// MultiTransferToVault -
func MultiTransferToVault(
	t *testing.T,
	net *integrationTests.TestNetwork,
	senderNode *integrationTests.TestProcessorNode,
	vaultScAddress []byte,
	transfers []*DcdtTransfer,
	nrRoundsToWait int,
	userBalances map[string]map[int64]int64,
	scBalances map[string]map[int64]int64,
) {

	acceptMultiTransferEndpointName := "accept_funds_multi_transfer"
	senderAddress := senderNode.OwnAccount.Address

	txData := BuildDcdtMultiTransferTxData(vaultScAddress,
		transfers,
		acceptMultiTransferEndpointName,
	)

	integrationTests.CreateAndSendTransaction(
		senderNode,
		net.Nodes,
		big.NewInt(0),
		senderAddress,
		txData,
		integrationTests.AdditionalGasLimit,
	)
	WaitForOperationCompletion(net, nrRoundsToWait)

	// update expected balances after transfers
	for _, transfer := range transfers {
		userBalances[transfer.TokenIdentifier][transfer.Nonce] -= transfer.Amount
		scBalances[transfer.TokenIdentifier][transfer.Nonce] += transfer.Amount
	}

	// check expected vs actual values
	for _, transfer := range transfers {
		expectedUserBalance := userBalances[transfer.TokenIdentifier][transfer.Nonce]
		expectedScBalance := scBalances[transfer.TokenIdentifier][transfer.Nonce]

		dcdt.CheckAddressHasTokens(t, senderAddress, net.Nodes,
			[]byte(transfer.TokenIdentifier), transfer.Nonce, expectedUserBalance)
		dcdt.CheckAddressHasTokens(t, vaultScAddress, net.Nodes,
			[]byte(transfer.TokenIdentifier), transfer.Nonce, expectedScBalance)
	}
}

// DeployNonPayableSmartContract -
func DeployNonPayableSmartContract(
	t *testing.T,
	net *integrationTests.TestNetwork,
	deployerNode *integrationTests.TestProcessorNode,
	fileName string,
) []byte {
	scCode := wasm.GetSCCode(fileName)
	scAddress, _ := deployerNode.BlockchainHook.NewAddress(
		deployerNode.OwnAccount.Address,
		deployerNode.OwnAccount.Nonce,
		vmFactory.WasmVirtualMachine)

	integrationTests.CreateAndSendTransaction(
		deployerNode,
		net.Nodes,
		big.NewInt(0),
		testVm.CreateEmptyAddress(),
		wasm.CreateDeployTxDataNonPayable(scCode),
		integrationTests.AdditionalGasLimit,
	)
	WaitForOperationCompletion(net, 4)

	_, err := deployerNode.AccntState.GetExistingAccount(scAddress)
	require.Nil(t, err)

	return scAddress
}

// DcdtMultiTransferToVault -
func DcdtMultiTransferToVault(t *testing.T, crossShard bool, scCodeFilename string) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	// For cross shard, we use 2 nodes, with node[1] being the SC deployer, and node[0] being the caller
	numShards := 1
	nrRoundsToWait := numRoundsSameShard

	if crossShard {
		numShards = 2
		nrRoundsToWait = numRoundsCrossShard
	}

	net := integrationTests.NewTestNetworkSized(t, numShards, 1, 1)
	net.Start()
	defer net.Close()

	net.MintNodeAccountsUint64(10000000000)
	net.Step()

	senderNode := net.NodesSharded[0][0]
	if crossShard {
		senderNode = net.NodesSharded[1][0]
	}

	expectedIssuerBalance := make(map[string]map[int64]int64)
	expectedVaultBalance := make(map[string]map[int64]int64)

	// deploy vault SC
	vaultScAddress := DeployNonPayableSmartContract(t, net, net.NodesSharded[0][0], scCodeFilename)

	// issue two fungible tokens
	fungibleTokenIdentifier1 := IssueFungibleToken(t, net, senderNode, "FUNG1", 1000)
	fungibleTokenIdentifier2 := IssueFungibleToken(t, net, senderNode, "FUNG2", 1000)

	expectedIssuerBalance[fungibleTokenIdentifier1] = make(map[int64]int64)
	expectedIssuerBalance[fungibleTokenIdentifier2] = make(map[int64]int64)
	expectedVaultBalance[fungibleTokenIdentifier1] = make(map[int64]int64)
	expectedVaultBalance[fungibleTokenIdentifier2] = make(map[int64]int64)

	expectedIssuerBalance[fungibleTokenIdentifier1][0] = 1000
	expectedIssuerBalance[fungibleTokenIdentifier2][0] = 1000

	// issue two NFT, with multiple NFTCreate
	nonFungibleTokenIdentifier1 := IssueNft(net, senderNode, "NFT1", false)
	nonFungibleTokenIdentifier2 := IssueNft(net, senderNode, "NFT2", false)

	expectedIssuerBalance[nonFungibleTokenIdentifier1] = make(map[int64]int64)
	expectedIssuerBalance[nonFungibleTokenIdentifier2] = make(map[int64]int64)

	expectedVaultBalance[nonFungibleTokenIdentifier1] = make(map[int64]int64)
	expectedVaultBalance[nonFungibleTokenIdentifier2] = make(map[int64]int64)

	for i := int64(1); i <= 10; i++ {
		CreateNFT(t, net, senderNode, senderNode.OwnAccount, nonFungibleTokenIdentifier1, i)
		CreateNFT(t, net, senderNode, senderNode.OwnAccount, nonFungibleTokenIdentifier2, i)

		expectedIssuerBalance[nonFungibleTokenIdentifier1][i] = 1
		expectedIssuerBalance[nonFungibleTokenIdentifier2][i] = 1
	}

	// issue two SFTs, with two NFTCreate for each
	semiFungibleTokenIdentifier1 := IssueNft(net, senderNode, "SFT1", true)
	semiFungibleTokenIdentifier2 := IssueNft(net, senderNode, "SFT2", true)

	expectedIssuerBalance[semiFungibleTokenIdentifier1] = make(map[int64]int64)
	expectedIssuerBalance[semiFungibleTokenIdentifier2] = make(map[int64]int64)

	expectedVaultBalance[semiFungibleTokenIdentifier1] = make(map[int64]int64)
	expectedVaultBalance[semiFungibleTokenIdentifier2] = make(map[int64]int64)

	for i := int64(1); i <= 2; i++ {
		CreateSFT(t, net, senderNode, senderNode.OwnAccount, semiFungibleTokenIdentifier1, i, 1000)
		CreateSFT(t, net, senderNode, senderNode.OwnAccount, semiFungibleTokenIdentifier2, i, 1000)

		expectedIssuerBalance[semiFungibleTokenIdentifier1][i] = 1000
		expectedIssuerBalance[semiFungibleTokenIdentifier2][i] = 1000
	}

	// send a single DCDT with multi-transfer
	transfers := []*DcdtTransfer{
		{
			TokenIdentifier: fungibleTokenIdentifier1,
			Nonce:           0,
			Amount:          100,
		}}
	MultiTransferToVault(t, net, senderNode,
		vaultScAddress, transfers, nrRoundsToWait,
		expectedIssuerBalance, expectedVaultBalance,
	)

	// send two identical transfers with multi-transfer
	transfers = []*DcdtTransfer{
		{
			TokenIdentifier: fungibleTokenIdentifier1,
			Nonce:           0,
			Amount:          50,
		},
		{
			TokenIdentifier: fungibleTokenIdentifier1,
			Nonce:           0,
			Amount:          50,
		}}
	MultiTransferToVault(t, net, senderNode,
		vaultScAddress, transfers, nrRoundsToWait,
		expectedIssuerBalance, expectedVaultBalance,
	)

	// send two different transfers amounts, same token
	transfers = []*DcdtTransfer{
		{
			TokenIdentifier: fungibleTokenIdentifier1,
			Nonce:           0,
			Amount:          50,
		},
		{
			TokenIdentifier: fungibleTokenIdentifier1,
			Nonce:           0,
			Amount:          100,
		}}
	MultiTransferToVault(t, net, senderNode,
		vaultScAddress, transfers, nrRoundsToWait,
		expectedIssuerBalance, expectedVaultBalance,
	)

	// send two different tokens, same amount
	transfers = []*DcdtTransfer{
		{
			TokenIdentifier: fungibleTokenIdentifier1,
			Nonce:           0,
			Amount:          100,
		},
		{
			TokenIdentifier: fungibleTokenIdentifier2,
			Nonce:           0,
			Amount:          100,
		}}
	MultiTransferToVault(t, net, senderNode,
		vaultScAddress, transfers, nrRoundsToWait,
		expectedIssuerBalance, expectedVaultBalance,
	)

	// send single NFT
	transfers = []*DcdtTransfer{
		{
			TokenIdentifier: nonFungibleTokenIdentifier1,
			Nonce:           1,
			Amount:          1,
		}}
	MultiTransferToVault(t, net, senderNode,
		vaultScAddress, transfers, nrRoundsToWait,
		expectedIssuerBalance, expectedVaultBalance,
	)

	// send two NFTs, same token ID
	transfers = []*DcdtTransfer{
		{
			TokenIdentifier: nonFungibleTokenIdentifier1,
			Nonce:           2,
			Amount:          1,
		},
		{
			TokenIdentifier: nonFungibleTokenIdentifier1,
			Nonce:           3,
			Amount:          1,
		}}
	MultiTransferToVault(t, net, senderNode,
		vaultScAddress, transfers, nrRoundsToWait,
		expectedIssuerBalance, expectedVaultBalance,
	)

	// send two NFTs, different token ID
	transfers = []*DcdtTransfer{
		{
			TokenIdentifier: nonFungibleTokenIdentifier1,
			Nonce:           4,
			Amount:          1,
		},
		{
			TokenIdentifier: nonFungibleTokenIdentifier2,
			Nonce:           1,
			Amount:          1,
		}}
	MultiTransferToVault(t, net, senderNode,
		vaultScAddress, transfers, nrRoundsToWait,
		expectedIssuerBalance, expectedVaultBalance,
	)

	// send fours NFTs, two of each different token ID
	transfers = []*DcdtTransfer{
		{
			TokenIdentifier: nonFungibleTokenIdentifier1,
			Nonce:           5,
			Amount:          1,
		},
		{
			TokenIdentifier: nonFungibleTokenIdentifier2,
			Nonce:           2,
			Amount:          1,
		},
		{
			TokenIdentifier: nonFungibleTokenIdentifier1,
			Nonce:           6,
			Amount:          1,
		},
		{
			TokenIdentifier: nonFungibleTokenIdentifier2,
			Nonce:           3,
			Amount:          1,
		}}
	MultiTransferToVault(t, net, senderNode,
		vaultScAddress, transfers, nrRoundsToWait,
		expectedIssuerBalance, expectedVaultBalance,
	)

	// send single SFT
	transfers = []*DcdtTransfer{
		{
			TokenIdentifier: semiFungibleTokenIdentifier1,
			Nonce:           1,
			Amount:          100,
		}}
	MultiTransferToVault(t, net, senderNode,
		vaultScAddress, transfers, nrRoundsToWait,
		expectedIssuerBalance, expectedVaultBalance,
	)

	// send two SFTs, same token ID
	transfers = []*DcdtTransfer{
		{
			TokenIdentifier: semiFungibleTokenIdentifier1,
			Nonce:           1,
			Amount:          100,
		},
		{
			TokenIdentifier: semiFungibleTokenIdentifier1,
			Nonce:           2,
			Amount:          100,
		}}
	MultiTransferToVault(t, net, senderNode,
		vaultScAddress, transfers, nrRoundsToWait,
		expectedIssuerBalance, expectedVaultBalance,
	)

	// send two SFTs, different token ID
	transfers = []*DcdtTransfer{
		{
			TokenIdentifier: semiFungibleTokenIdentifier1,
			Nonce:           1,
			Amount:          100,
		},
		{
			TokenIdentifier: semiFungibleTokenIdentifier2,
			Nonce:           1,
			Amount:          100,
		}}
	MultiTransferToVault(t, net, senderNode,
		vaultScAddress, transfers, nrRoundsToWait,
		expectedIssuerBalance, expectedVaultBalance,
	)

	// send fours SFTs, two of each different token ID
	transfers = []*DcdtTransfer{
		{
			TokenIdentifier: semiFungibleTokenIdentifier1,
			Nonce:           1,
			Amount:          100,
		},
		{
			TokenIdentifier: semiFungibleTokenIdentifier2,
			Nonce:           2,
			Amount:          100,
		},
		{
			TokenIdentifier: semiFungibleTokenIdentifier1,
			Nonce:           2,
			Amount:          50,
		},
		{
			TokenIdentifier: semiFungibleTokenIdentifier2,
			Nonce:           1,
			Amount:          200,
		}}
	MultiTransferToVault(t, net, senderNode,
		vaultScAddress, transfers, nrRoundsToWait,
		expectedIssuerBalance, expectedVaultBalance,
	)

	// transfer all 3 types
	transfers = []*DcdtTransfer{
		{
			TokenIdentifier: fungibleTokenIdentifier1,
			Nonce:           0,
			Amount:          100,
		},
		{
			TokenIdentifier: semiFungibleTokenIdentifier2,
			Nonce:           2,
			Amount:          100,
		},
		{
			TokenIdentifier: nonFungibleTokenIdentifier1,
			Nonce:           7,
			Amount:          1,
		}}
	MultiTransferToVault(t, net, senderNode,
		vaultScAddress, transfers, nrRoundsToWait,
		expectedIssuerBalance, expectedVaultBalance,
	)
}
