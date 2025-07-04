package getAccount

import (
	"math/big"
	"testing"

	chainData "github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/api"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/common/enablers"
	"github.com/TerraDharitri/drt-go-chain/common/forking"
	"github.com/TerraDharitri/drt-go-chain/integrationTests"
	"github.com/TerraDharitri/drt-go-chain/node"
	"github.com/TerraDharitri/drt-go-chain/state"
	"github.com/TerraDharitri/drt-go-chain/state/blockInfoProviders"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
)

func createAccountsRepository(accDB state.AccountsAdapter, blockchain chainData.ChainHandler) state.AccountsRepository {
	provider, _ := blockInfoProviders.NewCurrentBlockInfo(blockchain)
	wrapper, _ := state.NewAccountsDBApi(accDB, provider)

	args := state.ArgsAccountsRepository{
		FinalStateAccountsWrapper:      wrapper,
		CurrentStateAccountsWrapper:    wrapper,
		HistoricalStateAccountsWrapper: wrapper,
	}
	accountsRepo, _ := state.NewAccountsRepository(args)

	return accountsRepo
}

func TestNode_GetAccountAccountDoesNotExistsShouldRetEmpty(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	trieStorage, _ := integrationTests.CreateTrieStorageManager(testscommon.CreateMemUnit())
	accDB, _ := integrationTests.CreateAccountsDB(0, trieStorage)
	rootHash, _ := accDB.Commit()

	genericEpochNotifier := forking.NewGenericEpochNotifier()
	enableEpochsHandler, _ := enablers.NewEnableEpochsHandler(integrationTests.CreateEnableEpochsConfig(), genericEpochNotifier)
	coreComponents := integrationTests.GetDefaultCoreComponents(enableEpochsHandler, genericEpochNotifier)
	coreComponents.AddressPubKeyConverterField = integrationTests.TestAddressPubkeyConverter

	dataComponents := integrationTests.GetDefaultDataComponents()
	_ = dataComponents.BlockChain.SetCurrentBlockHeaderAndRootHash(&block.Header{Nonce: 42}, rootHash)
	dataComponents.BlockChain.SetCurrentBlockHeaderHash([]byte("header hash"))

	stateComponents := integrationTests.GetDefaultStateComponents()
	stateComponents.AccountsRepo = createAccountsRepository(accDB, dataComponents.BlockChain)

	n, _ := node.NewNode(
		node.WithCoreComponents(coreComponents),
		node.WithDataComponents(dataComponents),
		node.WithStateComponents(stateComponents),
	)

	encodedAddress, err := integrationTests.TestAddressPubkeyConverter.Encode(integrationTests.CreateRandomBytes(32))
	require.Nil(t, err)
	recovAccnt, _, err := n.GetAccount(encodedAddress, api.AccountQueryOptions{})

	require.Nil(t, err)
	assert.Equal(t, uint64(0), recovAccnt.Nonce)
	assert.Equal(t, "0", recovAccnt.Balance)
	assert.Equal(t, "0", recovAccnt.DeveloperReward)
	assert.Empty(t, recovAccnt.OwnerAddress)
	assert.Nil(t, recovAccnt.CodeHash)
	assert.Nil(t, recovAccnt.RootHash)
}

func TestNode_GetAccountAccountExistsShouldReturn(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	testNonce := uint64(7)
	testBalance := big.NewInt(100)

	trieStorage, _ := integrationTests.CreateTrieStorageManager(testscommon.CreateMemUnit())
	accDB, _ := integrationTests.CreateAccountsDB(0, trieStorage)
	testPubkey := integrationTests.CreateAccount(accDB, testNonce, testBalance)
	rootHash, _ := accDB.Commit()

	genericEpochNotifier := forking.NewGenericEpochNotifier()
	enableEpochsHandler, _ := enablers.NewEnableEpochsHandler(integrationTests.CreateEnableEpochsConfig(), genericEpochNotifier)
	coreComponents := integrationTests.GetDefaultCoreComponents(enableEpochsHandler, genericEpochNotifier)
	coreComponents.AddressPubKeyConverterField = testscommon.RealWorldBech32PubkeyConverter

	dataComponents := integrationTests.GetDefaultDataComponents()
	_ = dataComponents.BlockChain.SetCurrentBlockHeaderAndRootHash(&block.Header{Nonce: 42}, rootHash)
	dataComponents.BlockChain.SetCurrentBlockHeaderHash([]byte("header hash"))

	stateComponents := integrationTests.GetDefaultStateComponents()
	stateComponents.AccountsRepo = createAccountsRepository(accDB, dataComponents.BlockChain)

	n, _ := node.NewNode(
		node.WithCoreComponents(coreComponents),
		node.WithDataComponents(dataComponents),
		node.WithStateComponents(stateComponents),
	)

	testAddress, err := coreComponents.AddressPubKeyConverter().Encode(testPubkey)
	require.Nil(t, err)
	recovAccnt, _, err := n.GetAccount(testAddress, api.AccountQueryOptions{})

	require.Nil(t, err)
	require.Equal(t, testNonce, recovAccnt.Nonce)
	require.Equal(t, testBalance.String(), recovAccnt.Balance)
}
