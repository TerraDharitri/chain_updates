package stateTrie

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/pubkeyConverter"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	dataTx "github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	"github.com/TerraDharitri/drt-go-chain-core/hashing/sha256"
	crypto "github.com/TerraDharitri/drt-go-chain-crypto"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/common/errChan"
	"github.com/TerraDharitri/drt-go-chain/common/holders"
	"github.com/TerraDharitri/drt-go-chain/common/statistics"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/epochStart"
	"github.com/TerraDharitri/drt-go-chain/integrationTests"
	"github.com/TerraDharitri/drt-go-chain/integrationTests/mock"
	dcdtCommon "github.com/TerraDharitri/drt-go-chain/integrationTests/vm/dcdt"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain/state"
	"github.com/TerraDharitri/drt-go-chain/state/factory"
	"github.com/TerraDharitri/drt-go-chain/state/iteratorChannelsProvider"
	"github.com/TerraDharitri/drt-go-chain/state/lastSnapshotMarker"
	"github.com/TerraDharitri/drt-go-chain/state/storagePruningManager"
	"github.com/TerraDharitri/drt-go-chain/state/storagePruningManager/evictionWaitingList"
	"github.com/TerraDharitri/drt-go-chain/storage"
	"github.com/TerraDharitri/drt-go-chain/storage/storageunit"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/enableEpochsHandlerMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/marshallerMock"
	stateMock "github.com/TerraDharitri/drt-go-chain/testscommon/state"
	testStorage "github.com/TerraDharitri/drt-go-chain/testscommon/storage"
	"github.com/TerraDharitri/drt-go-chain/trie"
)

const denomination = "000000000000000000"

func TestAccountsDB_RetrieveDataWithSomeValuesShouldWork(t *testing.T) {
	// test simulates creation of a new account, data trie retrieval,
	// adding a (key, value) pair in that data trie, committing changes
	// and then reloading the data trie based on the root hash generated before
	t.Parallel()

	key1 := []byte("ABC")
	val1 := []byte("123")
	key2 := []byte("DEF")
	val2 := []byte("456")
	_, account, adb := integrationTests.GenerateAddressJournalAccountAccountsDB()

	_ = account.SaveKeyValue(key1, val1)
	_ = account.SaveKeyValue(key2, val2)

	err := adb.SaveAccount(account)
	require.Nil(t, err)

	_, err = adb.Commit()
	require.Nil(t, err)

	acc, err := adb.LoadAccount(account.AddressBytes())
	require.Nil(t, err)
	recoveredAccount := acc.(state.UserAccountHandler)

	// verify data
	dataRecovered, _, err := recoveredAccount.RetrieveValue(key1)
	require.Nil(t, err)
	require.Equal(t, val1, dataRecovered)

	dataRecovered, _, err = recoveredAccount.RetrieveValue(key2)
	require.Nil(t, err)
	require.Equal(t, val2, dataRecovered)
}

func TestAccountsDB_PutCodeWithSomeValuesShouldWork(t *testing.T) {
	t.Parallel()

	_, account, adb := integrationTests.GenerateAddressJournalAccountAccountsDB()
	account.SetCode([]byte("Smart contract code"))
	err := adb.SaveAccount(account)
	require.Nil(t, err)
	require.NotNil(t, account.GetCodeHash())
	require.Equal(t, []byte("Smart contract code"), adb.GetCode(account.GetCodeHash()))

	fmt.Printf("SC code is at address: %v\n", account.GetCodeHash())

	acc, err := adb.LoadAccount(account.AddressBytes())
	require.Nil(t, err)
	recoveredAccount := acc.(state.UserAccountHandler)

	require.Equal(t, adb.GetCode(account.GetCodeHash()), adb.GetCode(recoveredAccount.GetCodeHash()))
	require.Equal(t, account.GetCodeHash(), recoveredAccount.GetCodeHash())
}

func TestAccountsDB_SaveAccountStateWithSomeValues_ShouldWork(t *testing.T) {
	t.Parallel()

	_, account, adb := integrationTests.GenerateAddressJournalAccountAccountsDB()

	err := adb.SaveAccount(account)
	require.Nil(t, err)
}

func TestAccountsDB_GetJournalizedAccountReturnExistingAccntShouldWork(t *testing.T) {
	t.Parallel()

	balance := big.NewInt(40)
	adr, account, adb := integrationTests.GenerateAddressJournalAccountAccountsDB()
	_ = account.AddToBalance(balance)

	err := adb.SaveAccount(account)
	require.Nil(t, err)

	accountHandlerRecovered, err := adb.LoadAccount(adr)
	require.Nil(t, err)
	accountRecovered := accountHandlerRecovered.(state.UserAccountHandler)
	require.NotNil(t, accountRecovered)
	require.Equal(t, balance, accountRecovered.GetBalance())
}

func TestAccountsDB_GetJournalizedAccountReturnNotFoundAccntShouldWork(t *testing.T) {
	// test when the account does not exist
	t.Parallel()

	adr, _, adb := integrationTests.GenerateAddressJournalAccountAccountsDB()

	// same address of the unsaved account
	accountHandlerRecovered, err := adb.LoadAccount(adr)
	require.Nil(t, err)
	accountRecovered := accountHandlerRecovered.(state.UserAccountHandler)
	require.NotNil(t, accountRecovered)
	require.Equal(t, big.NewInt(0), accountRecovered.GetBalance())
}

func TestAccountsDB_GetExistingAccountConcurrentlyShouldWork(t *testing.T) {
	t.Parallel()

	trieStorage, _ := integrationTests.CreateTrieStorageManager(integrationTests.CreateMemUnit())
	adb, _ := integrationTests.CreateAccountsDB(0, trieStorage)

	wg := sync.WaitGroup{}
	wg.Add(100)

	addresses := make([][]byte, 0)

	// generating 100 different addresses
	for len(addresses) < 100 {
		addr := integrationTests.CreateRandomAddress()

		found := false
		for i := 0; i < len(addresses); i++ {
			if bytes.Equal(addresses[i], addr) {
				found = true
				break
			}
		}

		if !found {
			addresses = append(addresses, addr)
		}
	}

	for i := 0; i < 50; i++ {
		go func(idx int) {
			accnt, err := adb.GetExistingAccount(addresses[idx*2])

			require.Equal(t, state.ErrAccNotFound, err)
			require.Nil(t, accnt)

			wg.Done()
		}(i)

		go func(idx int) {
			accnt, err := adb.LoadAccount(addresses[idx*2+1])
			require.Nil(t, err)
			require.NotNil(t, accnt)

			err = adb.SaveAccount(accnt)
			require.Nil(t, err)

			wg.Done()
		}(i)
	}

	wg.Wait()
}

func TestAccountsDB_CommitTwoOkAccountsShouldWork(t *testing.T) {
	// test creates 2 accounts (one with a data root)
	// verifies that commit saves the new tries and that can be loaded back
	t.Parallel()

	adr1, _, adb := integrationTests.GenerateAddressJournalAccountAccountsDB()
	adr2 := integrationTests.CreateRandomAddress()

	// first account has the balance of 40
	balance1 := big.NewInt(40)
	state1, err := adb.LoadAccount(adr1)
	require.Nil(t, err)
	_ = state1.(state.UserAccountHandler).AddToBalance(balance1)

	// second account has the balance of 50 and some data
	balance2 := big.NewInt(50)
	acc, err := adb.LoadAccount(adr2)
	require.Nil(t, err)

	userAccount := acc.(state.UserAccountHandler)
	_ = userAccount.AddToBalance(balance2)

	key := []byte("ABC")
	val := []byte("123")
	_ = userAccount.SaveKeyValue(key, val)

	_ = adb.SaveAccount(state1)
	_ = adb.SaveAccount(userAccount)

	// states are now prepared, committing

	h, err := adb.Commit()
	require.Nil(t, err)
	fmt.Printf("Result hash: %v\n", base64.StdEncoding.EncodeToString(h))

	rootHash, err := adb.RootHash()
	require.Nil(t, err)
	fmt.Printf("data committed! Root: %v\n", base64.StdEncoding.EncodeToString(rootHash))

	// reloading a new trie to test if data is inside
	rootHash, err = adb.RootHash()
	require.Nil(t, err)
	err = adb.RecreateTrie(holders.NewDefaultRootHashesHolder(rootHash))
	require.Nil(t, err)

	// checking state1
	newState1, err := adb.LoadAccount(adr1)
	require.Nil(t, err)
	require.Equal(t, balance1, newState1.(state.UserAccountHandler).GetBalance())

	// checking stateMock
	newState2, err := adb.LoadAccount(adr2)
	require.Nil(t, err)
	require.Equal(t, balance2, newState2.(state.UserAccountHandler).GetBalance())
	require.NotNil(t, newState2.(state.UserAccountHandler).GetRootHash())
	valRecovered, _, err := newState2.(state.UserAccountHandler).RetrieveValue(key)
	require.Nil(t, err)
	require.Equal(t, val, valRecovered)
}

func TestTrieDB_RecreateFromStorageShouldWork(t *testing.T) {
	hasher := integrationTests.TestHasher
	store := integrationTests.CreateMemUnit()
	args := testStorage.GetStorageManagerArgs()
	args.MainStorer = store
	args.Hasher = hasher
	trieStorage, _ := trie.NewTrieStorageManager(args)

	maxTrieLevelInMemory := uint(5)
	tr1, _ := trie.NewTrie(trieStorage, integrationTests.TestMarshalizer, hasher, &enableEpochsHandlerMock.EnableEpochsHandlerStub{}, maxTrieLevelInMemory)

	key := hasher.Compute("key")
	value := hasher.Compute("value")

	_ = tr1.Update(key, value)
	h1, _ := tr1.RootHash()
	err := tr1.Commit()
	require.Nil(t, err)

	tr2, err := tr1.Recreate(holders.NewDefaultRootHashesHolder(h1))
	require.Nil(t, err)

	valRecov, _, err := tr2.Get(key)
	require.Nil(t, err)
	require.Equal(t, value, valRecov)
}

func TestAccountsDB_CommitTwoOkAccountsWithRecreationFromStorageShouldWork(t *testing.T) {
	// test creates 2 accounts (one with a data root)
	// verifies that commit saves the new tries and that can be loaded back
	t.Parallel()

	trieStore, _ := integrationTests.CreateTrieStorageManager(integrationTests.CreateMemUnit())
	adb, _ := integrationTests.CreateAccountsDB(0, trieStore)
	adr1 := integrationTests.CreateRandomAddress()
	adr2 := integrationTests.CreateRandomAddress()

	// first account has the balance of 40
	balance1 := big.NewInt(40)
	state1, err := adb.LoadAccount(adr1)
	require.Nil(t, err)
	_ = state1.(state.UserAccountHandler).AddToBalance(balance1)

	// second account has the balance of 50 and some data
	balance2 := big.NewInt(50)
	acc, err := adb.LoadAccount(adr2)
	require.Nil(t, err)

	userAccount := acc.(state.UserAccountHandler)
	_ = userAccount.AddToBalance(balance2)

	key := []byte("ABC")
	val := []byte("123")
	_ = userAccount.SaveKeyValue(key, val)

	_ = adb.SaveAccount(state1)
	_ = adb.SaveAccount(userAccount)

	// states are now prepared, committing

	h, err := adb.Commit()
	require.Nil(t, err)
	fmt.Printf("Result hash: %v\n", base64.StdEncoding.EncodeToString(h))

	rootHash, err := adb.RootHash()
	require.Nil(t, err)
	fmt.Printf("data committed! Root: %v\n", base64.StdEncoding.EncodeToString(rootHash))

	// reloading a new trie to test if data is inside
	err = adb.RecreateTrie(holders.NewDefaultRootHashesHolder(h))
	require.Nil(t, err)

	// checking state1
	newState1, err := adb.LoadAccount(adr1)
	require.Nil(t, err)
	require.Equal(t, balance1, newState1.(state.UserAccountHandler).GetBalance())

	// checking stateMock
	acc2, err := adb.LoadAccount(adr2)
	require.Nil(t, err)
	newState2 := acc2.(state.UserAccountHandler)
	require.Equal(t, balance2, newState2.GetBalance())
	require.NotNil(t, newState2.GetRootHash())
	valRecovered, _, err := newState2.RetrieveValue(key)
	require.Nil(t, err)
	require.Equal(t, val, valRecovered)
}

func TestAccountsDB_CommitAnEmptyStateShouldWork(t *testing.T) {
	t.Parallel()

	defer func() {
		r := recover()
		if r != nil {
			require.Fail(t, "this test should not have panicked")
		}
	}()

	trieStorage, _ := integrationTests.CreateTrieStorageManager(integrationTests.CreateMemUnit())
	adb, _ := integrationTests.CreateAccountsDB(0, trieStorage)

	hash, err := adb.Commit()

	require.Nil(t, err)
	require.Equal(t, make([]byte, state.HashLength), hash)
}

func TestAccountsDB_CommitAccountDataShouldWork(t *testing.T) {
	t.Parallel()

	adr1, _, adb := integrationTests.GenerateAddressJournalAccountAccountsDB()

	rootHash, err := adb.RootHash()
	require.Nil(t, err)
	hrEmpty := base64.StdEncoding.EncodeToString(rootHash)
	fmt.Printf("State root - empty: %v\n", hrEmpty)

	state1, err := adb.LoadAccount(adr1)
	require.Nil(t, err)
	_ = adb.SaveAccount(state1)

	rootHash, err = adb.RootHash()
	require.Nil(t, err)
	hrCreated := base64.StdEncoding.EncodeToString(rootHash)
	fmt.Printf("State root - created account: %v\n", hrCreated)

	_ = state1.(state.UserAccountHandler).AddToBalance(big.NewInt(40))
	_ = adb.SaveAccount(state1)

	rootHash, err = adb.RootHash()
	require.Nil(t, err)
	hrWithBalance := base64.StdEncoding.EncodeToString(rootHash)
	fmt.Printf("State root - account with balance 40: %v\n", hrWithBalance)

	_, err = adb.Commit()
	require.Nil(t, err)
	rootHash, err = adb.RootHash()
	require.Nil(t, err)
	hrCommit := base64.StdEncoding.EncodeToString(rootHash)
	fmt.Printf("State root - committed: %v\n", hrCommit)

	// commit hash == account with balance
	require.Equal(t, hrCommit, hrWithBalance)

	_ = state1.(state.UserAccountHandler).SubFromBalance(big.NewInt(40))
	_ = adb.SaveAccount(state1)

	// root hash == hrCreated
	rootHash, err = adb.RootHash()
	require.Nil(t, err)
	require.Equal(t, hrCreated, base64.StdEncoding.EncodeToString(rootHash))
	fmt.Printf("State root - account with balance 0: %v\n", base64.StdEncoding.EncodeToString(rootHash))

	err = adb.RemoveAccount(adr1)
	require.Nil(t, err)

	// root hash == hrEmpty
	rootHash, err = adb.RootHash()
	require.Nil(t, err)
	require.Equal(t, hrEmpty, base64.StdEncoding.EncodeToString(rootHash))
	fmt.Printf("State root - empty: %v\n", base64.StdEncoding.EncodeToString(rootHash))
}

func TestAccountsDB_RevertNonceStepByStepAccountDataShouldWork(t *testing.T) {
	t.Parallel()

	adr1 := integrationTests.CreateRandomAddress()
	adr2 := integrationTests.CreateRandomAddress()

	// Step 1. create accounts objects
	trieStorage, _ := integrationTests.CreateTrieStorageManager(integrationTests.CreateMemUnit())
	adb, _ := integrationTests.CreateAccountsDB(0, trieStorage)
	rootHash, err := adb.RootHash()
	require.Nil(t, err)
	hrEmpty := base64.StdEncoding.EncodeToString(rootHash)
	fmt.Printf("State root - empty: %v\n", hrEmpty)

	// Step 2. create 2 new accounts
	state1, err := adb.LoadAccount(adr1)
	require.Nil(t, err)
	_ = adb.SaveAccount(state1)

	snapshotCreated1 := adb.JournalLen()
	rootHash, err = adb.RootHash()
	require.Nil(t, err)
	hrCreated1 := base64.StdEncoding.EncodeToString(rootHash)

	fmt.Printf("State root - created 1-st account: %v\n", hrCreated1)

	userAccount, err := adb.LoadAccount(adr2)
	require.Nil(t, err)
	_ = adb.SaveAccount(userAccount)
	snapshotCreated2 := adb.JournalLen()
	rootHash, err = adb.RootHash()
	require.Nil(t, err)
	hrCreated2 := base64.StdEncoding.EncodeToString(rootHash)

	fmt.Printf("State root - created 2-nd account: %v\n", hrCreated2)

	// Test 2.1. test that hashes and snapshots ID are different
	require.NotEqual(t, snapshotCreated2, snapshotCreated1)
	require.NotEqual(t, hrCreated1, hrCreated2)

	// Save the preset snapshot id
	snapshotPreSet := adb.JournalLen()

	// Step 3. Set Nonces and save data
	state1.(state.UserAccountHandler).IncreaseNonce(40)
	_ = adb.SaveAccount(state1)

	rootHash, err = adb.RootHash()
	require.Nil(t, err)
	hrWithNonce1 := base64.StdEncoding.EncodeToString(rootHash)
	fmt.Printf("State root - account with nonce 40: %v\n", hrWithNonce1)

	userAccount.(state.UserAccountHandler).IncreaseNonce(50)
	_ = adb.SaveAccount(userAccount)

	rootHash, err = adb.RootHash()
	require.Nil(t, err)
	hrWithNonce2 := base64.StdEncoding.EncodeToString(rootHash)
	fmt.Printf("State root - account with nonce 50: %v\n", hrWithNonce2)

	// Test 3.1. current root hash shall not match created root hash hrCreated2
	rootHash, err = adb.RootHash()
	require.Nil(t, err)
	require.NotEqual(t, hrCreated2, rootHash)

	// Step 4. Revert account nonce and test
	err = adb.RevertToSnapshot(snapshotPreSet)
	require.Nil(t, err)

	// Test 4.1. current root hash shall match created root hash hrCreated
	rootHash, err = adb.RootHash()
	require.Nil(t, err)
	hrFinal := base64.StdEncoding.EncodeToString(rootHash)
	require.Equal(t, hrCreated2, hrFinal)
	fmt.Printf("State root - reverted last 2 nonces set: %v\n", hrFinal)
}

func TestAccountsDB_RevertBalanceStepByStepAccountDataShouldWork(t *testing.T) {
	t.Parallel()

	adr1 := integrationTests.CreateRandomAddress()
	adr2 := integrationTests.CreateRandomAddress()

	// Step 1. create accounts objects
	trieStorage, _ := integrationTests.CreateTrieStorageManager(integrationTests.CreateMemUnit())
	adb, _ := integrationTests.CreateAccountsDB(0, trieStorage)
	rootHash, err := adb.RootHash()
	require.Nil(t, err)
	hrEmpty := base64.StdEncoding.EncodeToString(rootHash)
	fmt.Printf("State root - empty: %v\n", hrEmpty)

	// Step 2. create 2 new accounts
	state1, err := adb.LoadAccount(adr1)
	require.Nil(t, err)
	_ = adb.SaveAccount(state1)

	snapshotCreated1 := adb.JournalLen()
	rootHash, err = adb.RootHash()
	require.Nil(t, err)
	hrCreated1 := base64.StdEncoding.EncodeToString(rootHash)

	fmt.Printf("State root - created 1-st account: %v\n", hrCreated1)

	userAccount, err := adb.LoadAccount(adr2)
	require.Nil(t, err)
	_ = adb.SaveAccount(userAccount)

	snapshotCreated2 := adb.JournalLen()
	rootHash, err = adb.RootHash()
	require.Nil(t, err)
	hrCreated2 := base64.StdEncoding.EncodeToString(rootHash)

	fmt.Printf("State root - created 2-nd account: %v\n", hrCreated2)

	// Test 2.1. test that hashes and snapshots ID are different
	require.NotEqual(t, snapshotCreated2, snapshotCreated1)
	require.NotEqual(t, hrCreated1, hrCreated2)

	// Save the preset snapshot id
	snapshotPreSet := adb.JournalLen()

	// Step 3. Set balances and save data
	_ = state1.(state.UserAccountHandler).AddToBalance(big.NewInt(40))
	_ = adb.SaveAccount(state1)

	rootHash, err = adb.RootHash()
	require.Nil(t, err)
	hrWithBalance1 := base64.StdEncoding.EncodeToString(rootHash)
	fmt.Printf("State root - account with balance 40: %v\n", hrWithBalance1)

	_ = userAccount.(state.UserAccountHandler).AddToBalance(big.NewInt(50))
	_ = adb.SaveAccount(userAccount)

	rootHash, err = adb.RootHash()
	require.Nil(t, err)
	hrWithBalance2 := base64.StdEncoding.EncodeToString(rootHash)
	fmt.Printf("State root - account with balance 50: %v\n", hrWithBalance2)

	// Test 3.1. current root hash shall not match created root hash hrCreated2
	require.NotEqual(t, hrCreated2, rootHash)

	// Step 4. Revert account balances and test
	err = adb.RevertToSnapshot(snapshotPreSet)
	require.Nil(t, err)

	// Test 4.1. current root hash shall match created root hash hrCreated
	rootHash, err = adb.RootHash()
	require.Nil(t, err)
	hrFinal := base64.StdEncoding.EncodeToString(rootHash)
	require.Equal(t, hrCreated2, hrFinal)
	fmt.Printf("State root - reverted last 2 balance set: %v\n", hrFinal)
}

func TestAccountsDB_RevertCodeStepByStepAccountDataShouldWork(t *testing.T) {
	t.Parallel()

	// adr1 puts code hash + code inside trie. adr2 has the same code hash
	// revert should work

	code := []byte("ABC")
	adr1 := integrationTests.CreateRandomAddress()
	adr2 := integrationTests.CreateRandomAddress()

	// Step 1. create accounts objects
	trieStorage, _ := integrationTests.CreateTrieStorageManager(integrationTests.CreateMemUnit())
	adb, _ := integrationTests.CreateAccountsDB(0, trieStorage)
	rootHash, err := adb.RootHash()
	require.Nil(t, err)
	hrEmpty := base64.StdEncoding.EncodeToString(rootHash)
	fmt.Printf("State root - empty: %v\n", hrEmpty)

	// Step 2. create 2 new accounts
	state1, err := adb.LoadAccount(adr1)
	require.Nil(t, err)
	state1.(state.UserAccountHandler).SetCode(code)
	_ = adb.SaveAccount(state1)

	snapshotCreated1 := adb.JournalLen()
	rootHash, err = adb.RootHash()
	require.Nil(t, err)
	hrCreated1 := base64.StdEncoding.EncodeToString(rootHash)

	fmt.Printf("State root - created 1-st account: %v\n", hrCreated1)

	userAccount, err := adb.LoadAccount(adr2)
	require.Nil(t, err)
	userAccount.(state.UserAccountHandler).SetCode(code)
	_ = adb.SaveAccount(userAccount)

	snapshotCreated2 := adb.JournalLen()
	rootHash, err = adb.RootHash()
	require.Nil(t, err)
	hrCreated2 := base64.StdEncoding.EncodeToString(rootHash)

	fmt.Printf("State root - created 2-nd account: %v\n", hrCreated2)

	// Test 2.1. test that hashes and snapshots ID are different
	require.NotEqual(t, snapshotCreated2, snapshotCreated1)
	require.NotEqual(t, hrCreated1, hrCreated2)

	// Step 3. Revert second account
	err = adb.RevertToSnapshot(snapshotCreated1)
	require.Nil(t, err)

	// Test 3.1. current root hash shall match created root hash hrCreated1
	rootHash, err = adb.RootHash()
	require.Nil(t, err)
	hrCrt := base64.StdEncoding.EncodeToString(rootHash)
	require.Equal(t, hrCreated1, hrCrt)
	fmt.Printf("State root - reverted last account: %v\n", hrCrt)

	// Step 4. Revert first account
	err = adb.RevertToSnapshot(0)
	require.Nil(t, err)

	// Test 4.1. current root hash shall match empty root hash
	rootHash, err = adb.RootHash()
	require.Nil(t, err)
	hrCrt = base64.StdEncoding.EncodeToString(rootHash)
	require.Equal(t, hrEmpty, hrCrt)
	fmt.Printf("State root - reverted first account: %v\n", hrCrt)
}

func TestAccountsDB_RevertDataStepByStepAccountDataShouldWork(t *testing.T) {
	t.Parallel()

	// adr1 puts data inside trie. adr2 puts the same data
	// revert should work

	key := []byte("ABC")
	val := []byte("123")
	adr1 := integrationTests.CreateRandomAddress()
	adr2 := integrationTests.CreateRandomAddress()

	// Step 1. create accounts objects
	trieStorage, _ := integrationTests.CreateTrieStorageManager(integrationTests.CreateMemUnit())
	adb, _ := integrationTests.CreateAccountsDB(0, trieStorage)
	rootHash, err := adb.RootHash()
	require.Nil(t, err)
	hrEmpty := base64.StdEncoding.EncodeToString(rootHash)
	fmt.Printf("State root - empty: %v\n", hrEmpty)

	// Step 2. create 2 new accounts
	state1, err := adb.LoadAccount(adr1)
	require.Nil(t, err)
	_ = state1.(state.UserAccountHandler).SaveKeyValue(key, val)
	err = adb.SaveAccount(state1)
	require.Nil(t, err)
	snapshotCreated1 := adb.JournalLen()
	rootHash, err = adb.RootHash()
	require.Nil(t, err)
	hrCreated1 := base64.StdEncoding.EncodeToString(rootHash)
	rootHash, err = state1.(state.UserAccountHandler).DataTrie().RootHash()
	require.Nil(t, err)
	hrRoot1 := base64.StdEncoding.EncodeToString(rootHash)

	fmt.Printf("State root - created 1-st account: %v\n", hrCreated1)
	fmt.Printf("data root - 1-st account: %v\n", hrRoot1)

	userAccount, err := adb.LoadAccount(adr2)
	require.Nil(t, err)
	_ = userAccount.(state.UserAccountHandler).SaveKeyValue(key, val)
	err = adb.SaveAccount(userAccount)
	require.Nil(t, err)
	snapshotCreated2 := adb.JournalLen()
	rootHash, err = adb.RootHash()
	require.Nil(t, err)
	hrCreated2 := base64.StdEncoding.EncodeToString(rootHash)
	rootHash, err = state1.(state.UserAccountHandler).DataTrie().RootHash()
	require.Nil(t, err)
	hrRoot2 := base64.StdEncoding.EncodeToString(rootHash)

	fmt.Printf("State root - created 2-nd account: %v\n", hrCreated2)
	fmt.Printf("data root - 2-nd account: %v\n", hrRoot2)

	// Test 2.1. test that hashes and snapshots ID are different
	require.NotEqual(t, snapshotCreated2, snapshotCreated1)
	require.NotEqual(t, hrCreated1, hrCreated2)

	// Test 2.2 test whether the datatrie roots match
	require.Equal(t, hrRoot1, hrRoot2)

	// Step 3. Revert 2-nd account ant test roots
	err = adb.RevertToSnapshot(snapshotCreated1)
	require.Nil(t, err)
	rootHash, err = adb.RootHash()
	require.Nil(t, err)
	hrCreated2Rev := base64.StdEncoding.EncodeToString(rootHash)

	require.Equal(t, hrCreated1, hrCreated2Rev)

	// Step 4. Revert 1-st account ant test roots
	err = adb.RevertToSnapshot(0)
	require.Nil(t, err)
	rootHash, err = adb.RootHash()
	require.Nil(t, err)
	hrCreated1Rev := base64.StdEncoding.EncodeToString(rootHash)

	require.Equal(t, hrEmpty, hrCreated1Rev)
}

func TestAccountsDB_RevertDataStepByStepWithCommitsAccountDataShouldWork(t *testing.T) {
	t.Parallel()

	// adr1 puts data inside trie. adr2 puts the same data
	// revert should work

	key := []byte("ABC")
	val := []byte("123")
	newVal := []byte("124")
	adr1 := integrationTests.CreateRandomAddress()
	adr2 := integrationTests.CreateRandomAddress()

	// Step 1. create accounts objects
	trieStorage, _ := integrationTests.CreateTrieStorageManager(integrationTests.CreateMemUnit())
	adb, _ := integrationTests.CreateAccountsDB(0, trieStorage)
	rootHash, err := adb.RootHash()
	require.Nil(t, err)
	hrEmpty := base64.StdEncoding.EncodeToString(rootHash)
	fmt.Printf("State root - empty: %v\n", hrEmpty)

	// Step 2. create 2 new accounts
	state1, err := adb.LoadAccount(adr1)
	require.Nil(t, err)
	_ = state1.(state.UserAccountHandler).SaveKeyValue(key, val)
	err = adb.SaveAccount(state1)
	require.Nil(t, err)
	snapshotCreated1 := adb.JournalLen()
	rootHash, err = adb.RootHash()
	require.Nil(t, err)
	hrCreated1 := base64.StdEncoding.EncodeToString(rootHash)
	rootHash, err = state1.(state.UserAccountHandler).DataTrie().RootHash()
	require.Nil(t, err)
	hrRoot1 := base64.StdEncoding.EncodeToString(rootHash)

	fmt.Printf("State root - created 1-st account: %v\n", hrCreated1)
	fmt.Printf("data root - 1-st account: %v\n", hrRoot1)

	userAccount, err := adb.LoadAccount(adr2)
	require.Nil(t, err)
	_ = userAccount.(state.UserAccountHandler).SaveKeyValue(key, val)
	err = adb.SaveAccount(userAccount)
	require.Nil(t, err)
	snapshotCreated2 := adb.JournalLen()
	rootHash, err = adb.RootHash()
	require.Nil(t, err)
	hrCreated2 := base64.StdEncoding.EncodeToString(rootHash)
	rootHash, err = userAccount.(state.UserAccountHandler).DataTrie().RootHash()
	require.Nil(t, err)
	hrRoot2 := base64.StdEncoding.EncodeToString(rootHash)

	fmt.Printf("State root - created 2-nd account: %v\n", hrCreated2)
	fmt.Printf("data root - 2-nd account: %v\n", hrRoot2)

	// Test 2.1. test that hashes and snapshots ID are different
	require.NotEqual(t, snapshotCreated2, snapshotCreated1)
	require.NotEqual(t, hrCreated1, hrCreated2)

	// Test 2.2 test that the datatrie roots are different
	require.NotEqual(t, hrRoot1, hrRoot2)

	// Step 3. Commit
	rootCommit, _ := adb.Commit()
	hrCommit := base64.StdEncoding.EncodeToString(rootCommit)
	fmt.Printf("State root - committed: %v\n", hrCommit)

	// Step 4. 2-nd account changes its data
	snapshotMod := adb.JournalLen()

	userAccount, err = adb.LoadAccount(adr2)
	require.Nil(t, err)
	_ = userAccount.(state.UserAccountHandler).SaveKeyValue(key, newVal)
	err = adb.SaveAccount(userAccount)
	require.Nil(t, err)
	rootHash, err = adb.RootHash()
	require.Nil(t, err)
	hrCreated2p1 := base64.StdEncoding.EncodeToString(rootHash)
	rootHash, err = userAccount.(state.UserAccountHandler).DataTrie().RootHash()
	require.Nil(t, err)
	hrRoot2p1 := base64.StdEncoding.EncodeToString(rootHash)

	fmt.Printf("State root - modified 2-nd account: %v\n", hrCreated2p1)
	fmt.Printf("data root - 2-nd account: %v\n", hrRoot2p1)

	// Test 4.1 test that hashes are different
	require.NotEqual(t, hrCreated2p1, hrCreated2)

	// Test 4.2 test whether the datatrie roots match/mismatch
	require.NotEqual(t, hrRoot2, hrRoot2p1)

	// Step 5. Revert 2-nd account modification
	err = adb.RevertToSnapshot(snapshotMod)
	require.Nil(t, err)
	rootHash, err = adb.RootHash()
	require.Nil(t, err)
	hrCreated2Rev := base64.StdEncoding.EncodeToString(rootHash)

	userAccount, err = adb.LoadAccount(adr2)
	require.Nil(t, err)
	rootHash, err = userAccount.(state.UserAccountHandler).DataTrie().RootHash()
	require.Nil(t, err)
	hrRoot2Rev := base64.StdEncoding.EncodeToString(rootHash)
	fmt.Printf("State root - reverted 2-nd account: %v\n", hrCreated2Rev)
	fmt.Printf("data root - 2-nd account: %v\n", hrRoot2Rev)
	require.Equal(t, hrCommit, hrCreated2Rev)
	require.Equal(t, hrRoot2, hrRoot2Rev)
}

func TestAccountsDB_ExecBalanceTxExecution(t *testing.T) {
	t.Parallel()

	adrSrc := integrationTests.CreateRandomAddress()
	adrDest := integrationTests.CreateRandomAddress()

	// Step 1. create accounts objects
	trieStorage, _ := integrationTests.CreateTrieStorageManager(integrationTests.CreateMemUnit())
	adb, _ := integrationTests.CreateAccountsDB(0, trieStorage)

	acntSrc, err := adb.LoadAccount(adrSrc)
	require.Nil(t, err)
	acntDest, err := adb.LoadAccount(adrDest)
	require.Nil(t, err)

	// Set a high balance to src's account
	_ = acntSrc.(state.UserAccountHandler).AddToBalance(big.NewInt(1000))
	_ = adb.SaveAccount(acntSrc)

	rootHash, err := adb.RootHash()
	require.Nil(t, err)
	hrOriginal := base64.StdEncoding.EncodeToString(rootHash)
	fmt.Printf("Original root hash: %s\n", hrOriginal)

	integrationTests.PrintShardAccount(acntSrc.(state.UserAccountHandler), "Source")
	integrationTests.PrintShardAccount(acntDest.(state.UserAccountHandler), "Destination")

	fmt.Println("Executing OK transaction...")
	integrationTests.AdbEmulateBalanceTxSafeExecution(acntSrc.(state.UserAccountHandler), acntDest.(state.UserAccountHandler), adb, big.NewInt(64))

	rootHash, err = adb.RootHash()
	require.Nil(t, err)
	hrOK := base64.StdEncoding.EncodeToString(rootHash)
	fmt.Printf("After executing an OK tx root hash: %s\n", hrOK)

	integrationTests.PrintShardAccount(acntSrc.(state.UserAccountHandler), "Source")
	integrationTests.PrintShardAccount(acntDest.(state.UserAccountHandler), "Destination")

	fmt.Println("Executing NOK transaction...")
	integrationTests.AdbEmulateBalanceTxSafeExecution(acntSrc.(state.UserAccountHandler), acntDest.(state.UserAccountHandler), adb, big.NewInt(10000))

	rootHash, err = adb.RootHash()
	require.Nil(t, err)
	hrNok := base64.StdEncoding.EncodeToString(rootHash)
	fmt.Printf("After executing a NOK tx root hash: %s\n", hrNok)

	integrationTests.PrintShardAccount(acntSrc.(state.UserAccountHandler), "Source")
	integrationTests.PrintShardAccount(acntDest.(state.UserAccountHandler), "Destination")

	require.NotEqual(t, hrOriginal, hrOK)
	require.Equal(t, hrOK, hrNok)

}

func TestAccountsDB_ExecALotOfBalanceTxOK(t *testing.T) {
	t.Parallel()

	adrSrc := integrationTests.CreateRandomAddress()
	adrDest := integrationTests.CreateRandomAddress()

	// Step 1. create accounts objects
	trieStorage, _ := integrationTests.CreateTrieStorageManager(integrationTests.CreateMemUnit())
	adb, _ := integrationTests.CreateAccountsDB(0, trieStorage)

	acntSrc, err := adb.LoadAccount(adrSrc)
	require.Nil(t, err)
	acntDest, err := adb.LoadAccount(adrDest)
	require.Nil(t, err)

	// Set a high balance to src's account
	_ = acntSrc.(state.UserAccountHandler).AddToBalance(big.NewInt(10000000))
	_ = adb.SaveAccount(acntSrc)

	rootHash, err := adb.RootHash()
	require.Nil(t, err)
	hrOriginal := base64.StdEncoding.EncodeToString(rootHash)
	fmt.Printf("Original root hash: %s\n", hrOriginal)

	for i := 1; i <= 1000; i++ {
		err = integrationTests.AdbEmulateBalanceTxExecution(adb, acntSrc.(state.UserAccountHandler), acntDest.(state.UserAccountHandler), big.NewInt(int64(i)))

		require.Nil(t, err)
	}

	integrationTests.PrintShardAccount(acntSrc.(state.UserAccountHandler), "Source")
	integrationTests.PrintShardAccount(acntDest.(state.UserAccountHandler), "Destination")
}

func TestAccountsDB_ExecALotOfBalanceTxOKorNOK(t *testing.T) {
	t.Parallel()

	adrSrc := integrationTests.CreateRandomAddress()
	adrDest := integrationTests.CreateRandomAddress()

	// Step 1. create accounts objects
	trieStorage, _ := integrationTests.CreateTrieStorageManager(integrationTests.CreateMemUnit())
	adb, _ := integrationTests.CreateAccountsDB(0, trieStorage)

	acntSrc, err := adb.LoadAccount(adrSrc)
	require.Nil(t, err)
	acntDest, err := adb.LoadAccount(adrDest)
	require.Nil(t, err)

	// Set a high balance to src's account
	_ = acntSrc.(state.UserAccountHandler).AddToBalance(big.NewInt(10000000))
	_ = adb.SaveAccount(acntSrc)

	rootHash, err := adb.RootHash()
	require.Nil(t, err)
	hrOriginal := base64.StdEncoding.EncodeToString(rootHash)
	fmt.Printf("Original root hash: %s\n", hrOriginal)

	st := time.Now()
	for i := 1; i <= 1000; i++ {
		err = integrationTests.AdbEmulateBalanceTxExecution(adb, acntSrc.(state.UserAccountHandler), acntDest.(state.UserAccountHandler), big.NewInt(int64(i)))
		require.Nil(t, err)

		err = integrationTests.AdbEmulateBalanceTxExecution(adb, acntDest.(state.UserAccountHandler), acntSrc.(state.UserAccountHandler), big.NewInt(int64(1000000)))
		require.NotNil(t, err)
	}

	fmt.Printf("Done in %v\n", time.Since(st))

	integrationTests.PrintShardAccount(acntSrc.(state.UserAccountHandler), "Source")
	integrationTests.PrintShardAccount(acntDest.(state.UserAccountHandler), "Destination")
}

func BenchmarkCreateOneMillionAccountsWithMockDB(b *testing.B) {
	nrOfAccounts := 1000000
	balance := 1500000
	persist := mock.MockDB{}

	adb, _, tr := createAccounts(nrOfAccounts, balance, persist)
	var rtm runtime.MemStats
	runtime.GC()

	runtime.ReadMemStats(&rtm)
	fmt.Printf("Mem before committing %v accounts with mockDB: go mem - %s, sys mem - %s \n",
		nrOfAccounts,
		core.ConvertBytes(rtm.Alloc),
		core.ConvertBytes(rtm.Sys),
	)

	_, _ = adb.Commit()

	runtime.ReadMemStats(&rtm)
	fmt.Printf("Mem after committing %v accounts with mockDB: go mem - %s, sys mem - %s \n",
		nrOfAccounts,
		core.ConvertBytes(rtm.Alloc),
		core.ConvertBytes(rtm.Sys),
	)

	runtime.GC()
	runtime.ReadMemStats(&rtm)
	fmt.Printf("Mem after committing %v accounts with mockDB and garbage collection : go mem - %s, sys mem - %s \n",
		nrOfAccounts,
		core.ConvertBytes(rtm.Alloc),
		core.ConvertBytes(rtm.Sys),
	)

	_ = tr.String()
}

func BenchmarkCreateOneMillionAccounts(b *testing.B) {
	nrOfAccounts := 1000000
	nrTxs := 15000
	txVal := 100
	balance := nrTxs * txVal
	persist := mock.NewCountingDB()

	adb, addr, _ := createAccounts(nrOfAccounts, balance, persist)
	var rtm runtime.MemStats

	_, _ = adb.Commit()
	runtime.GC()

	runtime.ReadMemStats(&rtm)
	goMemCollapsedTrie := rtm.Alloc

	fmt.Println("Partially collapsed trie")
	createAndExecTxs(b, addr, nrTxs, nrOfAccounts, txVal, adb)

	runtime.ReadMemStats(&rtm)
	goMemAfterTxExec := rtm.Alloc

	fmt.Printf("Extra memory after the exec of %v txs: %s \n",
		nrTxs,
		core.ConvertBytes(goMemAfterTxExec-goMemCollapsedTrie),
	)

	fmt.Println("Total nr. of nodes in trie: ", persist.GetCounter())
	persist.Reset()
	_, _ = adb.Commit()
	fmt.Printf("Nr. of modified nodes after %v txs: %v \n", nrTxs, persist.GetCounter())

	rootHash, err := adb.RootHash()
	require.Nil(b, err)

	_ = adb.RecreateTrie(holders.NewDefaultRootHashesHolder(rootHash))
	fmt.Println("Completely collapsed trie")
	createAndExecTxs(b, addr, nrTxs, nrOfAccounts, txVal, adb)
}

func createAccounts(
	nrOfAccounts int,
	balance int,
	persist storage.Persister,
) (*state.AccountsDB, [][]byte, common.Trie) {
	cache, _ := storageunit.NewCache(storageunit.CacheConfig{Type: storageunit.LRUCache, Capacity: 10, Shards: 1, SizeInBytes: 0})
	store, _ := storageunit.NewStorageUnit(cache, persist)
	evictionWaitListSize := uint(100)

	ewlArgs := evictionWaitingList.MemoryEvictionWaitingListArgs{
		RootHashesSize: evictionWaitListSize,
		HashesSize:     evictionWaitListSize * 100,
	}
	ewl, _ := evictionWaitingList.NewMemoryEvictionWaitingList(ewlArgs)
	args := testStorage.GetStorageManagerArgs()
	args.MainStorer = store
	trieStorage, _ := trie.NewTrieStorageManager(args)
	maxTrieLevelInMemory := uint(5)
	tr, _ := trie.NewTrie(trieStorage, integrationTests.TestMarshalizer, integrationTests.TestHasher, &enableEpochsHandlerMock.EnableEpochsHandlerStub{}, maxTrieLevelInMemory)
	spm, _ := storagePruningManager.NewStoragePruningManager(ewl, 10)
	argsAccCreator := factory.ArgsAccountCreator{
		Hasher:              integrationTests.TestHasher,
		Marshaller:          integrationTests.TestMarshalizer,
		EnableEpochsHandler: &enableEpochsHandlerMock.EnableEpochsHandlerStub{},
	}
	accCreator, _ := factory.NewAccountCreator(argsAccCreator)
	snapshotsManager, _ := state.NewSnapshotsManager(state.ArgsNewSnapshotsManager{
		ProcessingMode:       common.Normal,
		Marshaller:           &marshallerMock.MarshalizerMock{},
		AddressConverter:     &testscommon.PubkeyConverterMock{},
		ProcessStatusHandler: &testscommon.ProcessStatusHandlerStub{},
		StateMetrics:         &stateMock.StateMetricsStub{},
		AccountFactory:       accCreator,
		ChannelsProvider:     iteratorChannelsProvider.NewUserStateIteratorChannelsProvider(),
		LastSnapshotMarker:   lastSnapshotMarker.NewLastSnapshotMarker(),
		StateStatsHandler:    statistics.NewStateStatistics(),
	})
	argsAccountsDB := state.ArgsAccountsDB{
		Trie:                  tr,
		Hasher:                integrationTests.TestHasher,
		Marshaller:            integrationTests.TestMarshalizer,
		AccountFactory:        accCreator,
		StoragePruningManager: spm,
		AddressConverter:      &testscommon.PubkeyConverterMock{},
		SnapshotsManager:      snapshotsManager,
	}
	adb, _ := state.NewAccountsDB(argsAccountsDB)

	addr := make([][]byte, nrOfAccounts)
	for i := 0; i < nrOfAccounts; i++ {
		addr[i] = integrationTests.CreateAccount(adb, 0, big.NewInt(int64(balance)))
	}

	return adb, addr, tr
}

func createAndExecTxs(
	b *testing.B,
	addr [][]byte,
	nrTxs int,
	nrOfAccounts int,
	txVal int,
	adb *state.AccountsDB,
) {

	txProcessor := integrationTests.CreateSimpleTxProcessor(adb)
	var totalTime int64 = 0

	for i := 0; i < nrTxs; i++ {
		sender := rand.Intn(nrOfAccounts)
		receiver := rand.Intn(nrOfAccounts)
		for sender == receiver {
			receiver = rand.Intn(nrOfAccounts)
		}

		tx := &dataTx.Transaction{
			Nonce:   1,
			Value:   big.NewInt(int64(txVal)),
			SndAddr: addr[sender],
			RcvAddr: addr[receiver],
		}

		startTime := time.Now()
		_, err := txProcessor.ProcessTransaction(tx)
		duration := time.Since(startTime)
		totalTime += int64(duration)
		require.Nil(b, err)
	}
	fmt.Printf("Time needed for executing %v transactions: %v \n", nrTxs, time.Duration(totalTime))
}

func BenchmarkTxExecution(b *testing.B) {
	adrSrc := integrationTests.CreateRandomAddress()
	adrDest := integrationTests.CreateRandomAddress()

	// Step 1. create accounts objects
	trieStorage, _ := integrationTests.CreateTrieStorageManager(integrationTests.CreateMemUnit())
	adb, _ := integrationTests.CreateAccountsDB(0, trieStorage)

	acntSrc, err := adb.LoadAccount(adrSrc)
	require.Nil(b, err)
	acntDest, err := adb.LoadAccount(adrDest)
	require.Nil(b, err)

	// Set a high balance to src's account
	_ = acntSrc.(state.UserAccountHandler).AddToBalance(big.NewInt(10000000))
	_ = adb.SaveAccount(acntSrc)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		integrationTests.AdbEmulateBalanceTxSafeExecution(acntSrc.(state.UserAccountHandler), acntDest.(state.UserAccountHandler), adb, big.NewInt(1))
	}
}

func TestTrieDbPruning_GetAccountAfterPruning(t *testing.T) {
	t.Parallel()

	adb := createAccountsDBTestSetup()

	hexPubkeyConverter, _ := pubkeyConverter.NewHexPubkeyConverter(32)
	address1, _ := hexPubkeyConverter.Decode("0000000000000000000000000000000000000000000000000000000000000000")
	address2, _ := hexPubkeyConverter.Decode("0000000000000000000000000000000000000000000000000000000000000001")
	address3, _ := hexPubkeyConverter.Decode("0000000000000000000000000000000000000000000000000000000000000002")

	newDefaultAccount(adb, address1)
	newDefaultAccount(adb, address2)
	account := newDefaultAccount(adb, address3)

	rootHash1, _ := adb.Commit()
	_ = account.(state.UserAccountHandler).AddToBalance(big.NewInt(1))
	_ = adb.SaveAccount(account)
	rootHash2, _ := adb.Commit()
	adb.PruneTrie(rootHash1, state.OldRoot, state.NewPruningHandler(state.EnableDataRemoval))

	err := adb.RecreateTrie(holders.NewDefaultRootHashesHolder(rootHash2))
	require.Nil(t, err)
	acc, err := adb.GetExistingAccount(address1)
	require.NotNil(t, acc)
	require.Nil(t, err)
}

func newDefaultAccount(adb *state.AccountsDB, address []byte) vmcommon.AccountHandler {
	account, _ := adb.LoadAccount(address)
	_ = adb.SaveAccount(account)

	return account
}

func TestAccountsDB_RecreateTrieInvalidatesDataTriesCache(t *testing.T) {
	adb := createAccountsDBTestSetup()

	hexAddressPubkeyConverter, _ := pubkeyConverter.NewHexPubkeyConverter(32)
	address1, _ := hexAddressPubkeyConverter.Decode("0000000000000000000000000000000000000000000000000000000000000000")

	key1 := []byte("ABC")
	key2 := []byte("ABD")
	value1 := []byte("dog")
	value2 := []byte("puppy")

	acc1, _ := adb.LoadAccount(address1)
	state1 := acc1.(state.UserAccountHandler)
	_ = state1.SaveKeyValue(key1, value1)
	_ = state1.SaveKeyValue(key2, value1)
	_ = adb.SaveAccount(state1)
	rootHash, err := adb.Commit()
	require.Nil(t, err)

	acc1, _ = adb.LoadAccount(address1)
	state1 = acc1.(state.UserAccountHandler)
	_ = state1.SaveKeyValue(key1, value2)
	_ = adb.SaveAccount(state1)
	_, err = adb.Commit()
	require.Nil(t, err)

	acc1, _ = adb.LoadAccount(address1)
	state1 = acc1.(state.UserAccountHandler)
	_ = state1.SaveKeyValue(key2, value2)
	_ = adb.SaveAccount(state1)
	err = adb.RevertToSnapshot(0)
	require.Nil(t, err)

	err = adb.RecreateTrie(holders.NewDefaultRootHashesHolder(rootHash))
	require.Nil(t, err)
	acc1, _ = adb.LoadAccount(address1)
	state1 = acc1.(state.UserAccountHandler)

	retrievedVal, _, _ := state1.RetrieveValue(key1)
	require.Equal(t, value1, retrievedVal)
}

func TestTrieDbPruning_GetDataTrieTrackerAfterPruning(t *testing.T) {
	t.Parallel()

	adb := createAccountsDBTestSetup()

	hexAddressPubkeyConverter, _ := pubkeyConverter.NewHexPubkeyConverter(32)
	address1, _ := hexAddressPubkeyConverter.Decode("0000000000000000000000000000000000000000000000000000000000000000")
	address2, _ := hexAddressPubkeyConverter.Decode("0000000000000000000000000000000000000000000000000000000000000001")

	key1 := []byte("ABC")
	key2 := []byte("ABD")
	value1 := []byte("dog")
	value2 := []byte("puppy")

	acc1, _ := adb.LoadAccount(address1)
	state1 := acc1.(state.UserAccountHandler)
	_ = state1.SaveKeyValue(key1, value1)
	_ = state1.SaveKeyValue(key2, value1)
	_ = adb.SaveAccount(state1)

	acc2, _ := adb.LoadAccount(address2)
	userAccount := acc2.(state.UserAccountHandler)
	_ = userAccount.SaveKeyValue(key1, value1)
	_ = userAccount.SaveKeyValue(key2, value1)
	_ = adb.SaveAccount(userAccount)

	oldRootHash, _ := adb.Commit()

	acc2, _ = adb.LoadAccount(address2)
	userAccount = acc2.(state.UserAccountHandler)
	_ = userAccount.SaveKeyValue(key1, value2)
	_ = adb.SaveAccount(userAccount)

	newRootHash, _ := adb.Commit()
	adb.PruneTrie(oldRootHash, state.OldRoot, state.NewPruningHandler(state.EnableDataRemoval))

	err := adb.RecreateTrie(holders.NewDefaultRootHashesHolder(newRootHash))
	require.Nil(t, err)
	acc, err := adb.GetExistingAccount(address1)
	require.NotNil(t, acc)
	require.Nil(t, err)

	collapseTrie(state1, t)
	collapseTrie(userAccount, t)

	val, _, err := state1.RetrieveValue(key1)
	require.Nil(t, err)
	require.Equal(t, value1, val)

	val, _, err = userAccount.RetrieveValue(key2)
	require.Nil(t, err)
	require.Equal(t, value1, val)
}

func collapseTrie(state state.UserAccountHandler, t *testing.T) {
	stateRootHash := state.GetRootHash()
	stateTrie := state.DataTrie().(common.Trie)
	stateNewTrie, _ := stateTrie.Recreate(holders.NewDefaultRootHashesHolder(stateRootHash))
	require.NotNil(t, stateNewTrie)

	state.SetDataTrie(stateNewTrie)
}

func TestRollbackBlockAndCheckThatPruningIsCancelledOnAccountsTrie(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	numNodesPerShard := 1
	numNodesMeta := 1

	nodes, leaders := integrationTests.SetupSyncNodesOneShardAndMeta(numNodesPerShard, numNodesMeta)
	defer integrationTests.CloseProcessorNodes(nodes)

	integrationTests.BootstrapDelay()
	integrationTests.StartSyncingBlocks(nodes)

	round := uint64(0)
	nonce := uint64(0)

	valMinting := big.NewInt(1000000000)
	valToTransferPerTx := big.NewInt(2)

	fmt.Println("Generating private keys for senders and receivers...")
	generateCoordinator, _ := sharding.NewMultiShardCoordinator(uint32(1), 0)
	nrTxs := 20

	// sender shard keys, receivers  keys
	sendersPrivateKeys := make([]crypto.PrivateKey, nrTxs)
	receiversPublicKeys := make(map[uint32][]crypto.PublicKey)
	for i := 0; i < nrTxs; i++ {
		sendersPrivateKeys[i], _, _ = integrationTests.GenerateSkAndPkInShard(generateCoordinator, 0)
		_, pk, _ := integrationTests.GenerateSkAndPkInShard(generateCoordinator, 0)
		receiversPublicKeys[0] = append(receiversPublicKeys[0], pk)
	}

	fmt.Println("Minting sender addresses...")
	integrationTests.CreateMintingForSenders(nodes, 0, sendersPrivateKeys, valMinting)

	shardNode := nodes[0]

	round = integrationTests.IncrementAndPrintRound(round)
	nonce++
	round, nonce = integrationTests.ProposeAndSyncOneBlock(t, nodes, leaders, round, nonce)

	rootHashOfFirstBlock, _ := shardNode.AccntState.RootHash()

	require.Equal(t, uint64(1), nodes[0].BlockChain.GetCurrentBlockHeader().GetNonce())
	require.Equal(t, uint64(1), nodes[1].BlockChain.GetCurrentBlockHeader().GetNonce())

	delayRounds := 10
	for i := 0; i < delayRounds; i++ {
		round, nonce = integrationTests.ProposeAndSyncOneBlock(t, nodes, leaders, round, nonce)
	}

	fmt.Println("Generating transactions...")
	integrationTests.GenerateAndDisseminateTxs(
		shardNode,
		sendersPrivateKeys,
		receiversPublicKeys,
		valToTransferPerTx,
		1000,
		1000,
		integrationTests.ChainID,
		integrationTests.MinTransactionVersion,
	)
	fmt.Println("Delaying for disseminating transactions...")
	time.Sleep(time.Second * 5)

	round, _ = integrationTests.ProposeAndSyncOneBlock(t, nodes, leaders, round, nonce)
	time.Sleep(time.Second * 5)

	rootHashOfRollbackedBlock, _ := shardNode.AccntState.RootHash()

	require.Equal(t, uint64(12), nodes[0].BlockChain.GetCurrentBlockHeader().GetNonce())
	require.Equal(t, uint64(12), nodes[1].BlockChain.GetCurrentBlockHeader().GetNonce())

	shardIdToRollbackLastBlock := uint32(0)
	integrationTests.ForkChoiceOneBlock(nodes, shardIdToRollbackLastBlock)
	integrationTests.ResetHighestProbableNonce(nodes, shardIdToRollbackLastBlock, 1)
	integrationTests.EmptyDataPools(nodes, shardIdToRollbackLastBlock)

	require.Equal(t, uint64(11), nodes[0].BlockChain.GetCurrentBlockHeader().GetNonce())
	require.Equal(t, uint64(12), nodes[1].BlockChain.GetCurrentBlockHeader().GetNonce())

	rootHash, err := shardNode.AccntState.RootHash()
	require.Nil(t, err)

	if !bytes.Equal(rootHash, rootHashOfRollbackedBlock) {
		time.Sleep(time.Second * 6)
		err = shardNode.AccntState.RecreateTrie(holders.NewDefaultRootHashesHolder(rootHashOfRollbackedBlock))
		require.True(t, strings.Contains(err.Error(), trie.ErrKeyNotFound.Error()))
	}

	nonces := []*uint64{new(uint64), new(uint64)}
	atomic.AddUint64(nonces[0], 2)
	atomic.AddUint64(nonces[1], 3)

	numOfRounds := 2
	integrationTests.ProposeBlocks(
		nodes,
		&round,
		leaders,
		nonces,
		numOfRounds,
	)
	time.Sleep(time.Second * 5)

	err = shardNode.AccntState.RecreateTrie(holders.NewDefaultRootHashesHolder(rootHashOfFirstBlock))
	require.Nil(t, err)
	require.Equal(t, uint64(11), nodes[0].BlockChain.GetCurrentBlockHeader().GetNonce())
	require.Equal(t, uint64(12), nodes[1].BlockChain.GetCurrentBlockHeader().GetNonce())
}

func TestRollbackBlockWithSameRootHashAsPreviousAndCheckThatPruningIsNotDone(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	numNodesPerShard := 1
	numNodesMeta := 1

	nodes, idxProposers := integrationTests.SetupSyncNodesOneShardAndMeta(numNodesPerShard, numNodesMeta)
	defer integrationTests.CloseProcessorNodes(nodes)

	integrationTests.BootstrapDelay()
	integrationTests.StartSyncingBlocks(nodes)

	round := uint64(0)
	nonce := uint64(0)

	valMinting := big.NewInt(1000000000)

	fmt.Println("Generating private keys for senders and receivers...")
	generateCoordinator, _ := sharding.NewMultiShardCoordinator(uint32(1), 0)
	nrTxs := 20

	// sender shard keys, receivers  keys
	sendersPrivateKeys := make([]crypto.PrivateKey, nrTxs)
	for i := 0; i < nrTxs; i++ {
		sendersPrivateKeys[i], _, _ = integrationTests.GenerateSkAndPkInShard(generateCoordinator, 0)
	}

	fmt.Println("Minting sender addresses...")
	integrationTests.CreateMintingForSenders(nodes, 0, sendersPrivateKeys, valMinting)

	shardNode := nodes[0]

	round = integrationTests.IncrementAndPrintRound(round)
	nonce++
	round, nonce = integrationTests.ProposeAndSyncOneBlock(t, nodes, idxProposers, round, nonce)

	rootHashOfFirstBlock, _ := shardNode.AccntState.RootHash()

	require.Equal(t, uint64(1), nodes[0].BlockChain.GetCurrentBlockHeader().GetNonce())
	require.Equal(t, uint64(1), nodes[1].BlockChain.GetCurrentBlockHeader().GetNonce())

	_, _ = integrationTests.ProposeAndSyncOneBlock(t, nodes, idxProposers, round, nonce)
	time.Sleep(time.Second * 5)

	require.Equal(t, uint64(2), nodes[0].BlockChain.GetCurrentBlockHeader().GetNonce())
	require.Equal(t, uint64(2), nodes[1].BlockChain.GetCurrentBlockHeader().GetNonce())

	shardIdToRollbackLastBlock := uint32(0)
	integrationTests.ForkChoiceOneBlock(nodes, shardIdToRollbackLastBlock)
	integrationTests.ResetHighestProbableNonce(nodes, shardIdToRollbackLastBlock, 1)
	integrationTests.EmptyDataPools(nodes, shardIdToRollbackLastBlock)

	require.Equal(t, uint64(1), nodes[0].BlockChain.GetCurrentBlockHeader().GetNonce())
	require.Equal(t, uint64(2), nodes[1].BlockChain.GetCurrentBlockHeader().GetNonce())

	err := shardNode.AccntState.RecreateTrie(holders.NewDefaultRootHashesHolder(rootHashOfFirstBlock))
	require.Nil(t, err)
}

func TestTriePruningWhenBlockIsFinal(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	fmt.Println("Setup nodes...")
	numOfShards := 1
	nodesPerShard := 1
	numMetachainNodes := 1

	senderShard := uint32(0)
	round := uint64(0)
	nonce := uint64(0)

	valMinting := big.NewInt(1000000000)
	valToTransferPerTx := big.NewInt(2)

	nodes, idxProposers := integrationTests.SetupSyncNodesOneShardAndMeta(nodesPerShard, numMetachainNodes)
	integrationTests.DisplayAndStartNodes(nodes)

	defer integrationTests.CloseProcessorNodes(nodes)

	fmt.Println("Generating private keys for senders and receivers...")
	generateCoordinator, _ := sharding.NewMultiShardCoordinator(uint32(numOfShards), 0)
	nrTxs := 20

	// sender shard keys, receivers  keys
	sendersPrivateKeys := make([]crypto.PrivateKey, nrTxs)
	receiversPublicKeys := make(map[uint32][]crypto.PublicKey)
	for i := 0; i < nrTxs; i++ {
		sendersPrivateKeys[i], _, _ = integrationTests.GenerateSkAndPkInShard(generateCoordinator, senderShard)
		_, pk, _ := integrationTests.GenerateSkAndPkInShard(generateCoordinator, senderShard)
		receiversPublicKeys[senderShard] = append(receiversPublicKeys[senderShard], pk)
	}

	fmt.Println("Minting sender addresses...")
	integrationTests.CreateMintingForSenders(nodes, senderShard, sendersPrivateKeys, valMinting)

	shardNode := nodes[0]

	round = integrationTests.IncrementAndPrintRound(round)
	nonce++
	round, nonce = integrationTests.ProposeAndSyncOneBlock(t, nodes, idxProposers, round, nonce)

	delayRounds := 10
	for i := 0; i < delayRounds; i++ {
		round, nonce = integrationTests.ProposeAndSyncOneBlock(t, nodes, idxProposers, round, nonce)
	}

	require.Equal(t, uint64(11), nodes[0].BlockChain.GetCurrentBlockHeader().GetNonce())
	require.Equal(t, uint64(11), nodes[1].BlockChain.GetCurrentBlockHeader().GetNonce())

	rootHashOfFirstBlock, _ := shardNode.AccntState.RootHash()

	fmt.Println("Generating transactions...")
	integrationTests.GenerateAndDisseminateTxs(
		shardNode,
		sendersPrivateKeys,
		receiversPublicKeys,
		valToTransferPerTx,
		1000,
		1000,
		integrationTests.ChainID,
		integrationTests.MinTransactionVersion,
	)
	fmt.Println("Delaying for disseminating transactions...")
	time.Sleep(time.Second * 5)

	roundsToWait := 6
	for i := 0; i < roundsToWait; i++ {
		round, nonce = integrationTests.ProposeAndSyncOneBlock(t, nodes, idxProposers, round, nonce)
	}

	require.Equal(t, uint64(17), nodes[0].BlockChain.GetCurrentBlockHeader().GetNonce())
	require.Equal(t, uint64(17), nodes[1].BlockChain.GetCurrentBlockHeader().GetNonce())

	err := shardNode.AccntState.RecreateTrie(holders.NewDefaultRootHashesHolder(rootHashOfFirstBlock))
	require.True(t, strings.Contains(err.Error(), trie.ErrKeyNotFound.Error()))
}

func TestStatePruningIsNotBuffered(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	numOfShards := 1
	nodesPerShard := 1
	numMetachainNodes := 1

	nodes := integrationTests.CreateNodes(
		numOfShards,
		nodesPerShard,
		numMetachainNodes,
	)

	shardNode := nodes[0]
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

	time.Sleep(integrationTests.StepDelay)

	round, nonce = integrationTests.ProposeAndSyncOneBlock(t, nodes, leaders, round, nonce)

	delayRounds := 5
	for j := 0; j < 8; j++ {
		// alter the shardNode's state by placing the value0 variable inside it's data trie
		alterState(t, shardNode, nodes, []byte("key"), []byte("value0"))
		for i := 0; i < delayRounds; i++ {
			round, nonce = integrationTests.ProposeAndSyncOneBlock(t, nodes, leaders, round, nonce)
		}
		checkTrieCanBeRecreated(t, shardNode)

		// alter the shardNode's state by placing the value1 variable inside it's data trie
		alterState(t, shardNode, nodes, []byte("key"), []byte("value1"))
		for i := 0; i < delayRounds; i++ {
			round, nonce = integrationTests.ProposeAndSyncOneBlock(t, nodes, leaders, round, nonce)
		}
		checkTrieCanBeRecreated(t, shardNode)
	}
}

func TestStatePruningIsNotBufferedOnConsecutiveBlocks(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	numOfShards := 1
	nodesPerShard := 1
	numMetachainNodes := 1

	nodes := integrationTests.CreateNodes(
		numOfShards,
		nodesPerShard,
		numMetachainNodes,
	)

	shardNode := nodes[0]
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

	time.Sleep(integrationTests.StepDelay)

	round, nonce = integrationTests.ProposeAndSyncOneBlock(t, nodes, leaders, round, nonce)

	for j := 0; j < 30; j++ {
		// alter the shardNode's state by placing the value0 variable inside it's data trie
		alterState(t, shardNode, nodes, []byte("key"), []byte("value0"))
		round, nonce = integrationTests.ProposeAndSyncOneBlock(t, nodes, leaders, round, nonce)
		checkTrieCanBeRecreated(t, shardNode)

		// alter the shardNode's state by placing the value1 variable inside it's data trie
		alterState(t, shardNode, nodes, []byte("key"), []byte("value1"))
		round, nonce = integrationTests.ProposeAndSyncOneBlock(t, nodes, leaders, round, nonce)
		checkTrieCanBeRecreated(t, shardNode)
	}
}

func alterState(tb testing.TB, node *integrationTests.TestProcessorNode, nodes []*integrationTests.TestProcessorNode, key []byte, value []byte) {
	shardID := node.ShardCoordinator.SelfId()
	for _, n := range nodes {
		if n.ShardCoordinator.SelfId() != shardID {
			continue
		}

		account, err := n.AccntState.LoadAccount(node.OwnAccount.Address)
		assert.Nil(tb, err)

		userAccount := account.(state.UserAccountHandler)
		err = userAccount.SaveKeyValue(key, value)
		assert.Nil(tb, err)

		err = n.AccntState.SaveAccount(userAccount)
		assert.Nil(tb, err)

		_, err = n.AccntState.Commit()
		assert.Nil(tb, err)
	}
}

func checkTrieCanBeRecreated(tb testing.TB, node *integrationTests.TestProcessorNode) {
	if node.ShardCoordinator.SelfId() == core.MetachainShardId {
		return
	}

	stateTrie := node.TrieContainer.Get([]byte(dataRetriever.UserAccountsUnit.String()))
	roothash := node.BlockChain.GetCurrentBlockRootHash()
	tr, err := stateTrie.Recreate(holders.NewDefaultRootHashesHolder(roothash))
	require.Nil(tb, err)
	require.NotNil(tb, tr)

	_, _, finalRoothash := node.BlockChain.GetFinalBlockInfo()
	tr, err = stateTrie.Recreate(holders.NewDefaultRootHashesHolder(finalRoothash))
	require.Nil(tb, err)
	require.NotNil(tb, tr)

	currentBlockHeader := node.BlockChain.GetCurrentBlockHeader()
	prevHeaderHash := currentBlockHeader.GetPrevHash()
	hdrBytes, err := node.Storage.Get(dataRetriever.BlockHeaderUnit, prevHeaderHash)
	require.Nil(tb, err)
	hdr := &block.Header{}
	err = integrationTests.TestMarshalizer.Unmarshal(hdr, hdrBytes)
	require.Nil(tb, err)

	tr, err = stateTrie.Recreate(holders.NewDefaultRootHashesHolder(hdr.GetRootHash()))
	require.Nil(tb, err)
	require.NotNil(tb, tr)
}

func TestSnapshotOnEpochChange(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	numOfShards := 1
	nodesPerShard := 1
	numMetachainNodes := 1

	enableEpochsConfig := integrationTests.GetDefaultEnableEpochsConfig()
	enableEpochsConfig.StakingV2EnableEpoch = integrationTests.UnreachableEpoch

	nodes := integrationTests.CreateNodesWithEnableEpochsConfig(
		numOfShards,
		nodesPerShard,
		numMetachainNodes,
		enableEpochsConfig,
	)

	roundsPerEpoch := uint64(17)
	for _, node := range nodes {
		node.EpochStartTrigger.SetRoundsPerEpoch(roundsPerEpoch)
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

	sendValue := big.NewInt(5)
	receiverAddress := []byte("12345678901234567890123456789012")
	initialVal := big.NewInt(10000000000)

	integrationTests.MintAllNodes(nodes, initialVal)

	round := uint64(0)
	nonce := uint64(0)
	round = integrationTests.IncrementAndPrintRound(round)
	nonce++

	time.Sleep(integrationTests.StepDelay)

	snapshotsRootHashes := make(map[uint32][][]byte)
	prunedRootHashes := make(map[int][][]byte)

	numShardNodes := numOfShards * nodesPerShard
	numRounds := uint32(20)
	for i := uint64(0); i < uint64(numRounds); i++ {

		round, nonce = integrationTests.ProposeAndSyncOneBlock(t, nodes, leaders, round, nonce)

		for _, node := range nodes {
			integrationTests.CreateAndSendTransaction(node, nodes, sendValue, receiverAddress, "", integrationTests.AdditionalGasLimit)
		}
		time.Sleep(integrationTests.StepDelay)

		collectSnapshotHashes(
			nodes,
			numShardNodes,
			snapshotsRootHashes,
			prunedRootHashes,
			roundsPerEpoch,
		)
		time.Sleep(time.Second)
	}

	numDelayRounds := uint32(15)
	for i := uint64(0); i < uint64(numDelayRounds); i++ {
		round, nonce = integrationTests.ProposeAndSyncOneBlock(t, nodes, leaders, round, nonce)

		for _, node := range nodes {
			integrationTests.CreateAndSendTransaction(node, nodes, sendValue, receiverAddress, "", integrationTests.AdditionalGasLimit)
		}
		time.Sleep(integrationTests.StepDelay)
	}

	for i := 0; i < numOfShards*nodesPerShard; i++ {
		shId := nodes[i].ShardCoordinator.SelfId()
		testNodeStateSnapshotAndPruning(t, nodes[i], snapshotsRootHashes[shId], prunedRootHashes[i])
	}
}

func collectSnapshotHashes(
	nodes []*integrationTests.TestProcessorNode,
	numShardNodes int,
	snapshotsRootHashes map[uint32][][]byte,
	prunedRootHashes map[int][][]byte,
	roundsPerEpoch uint64,
) {
	pruningQueueSize := uint64(5)
	finality := uint64(2)
	pruningDelayMultiplier := uint64(2)

	for j := 0; j < numShardNodes; j++ {
		currentBlockHeader := nodes[j].BlockChain.GetCurrentBlockHeader()
		if currentBlockHeader.IsStartOfEpochBlock() {
			continue
		}

		if currentBlockHeader.GetNonce() > roundsPerEpoch-pruningQueueSize-finality {
			continue
		}

		if currentBlockHeader.GetNonce() < pruningQueueSize*pruningDelayMultiplier {
			continue
		}

		prunedRootHashes[j] = append(prunedRootHashes[j], currentBlockHeader.GetRootHash())
	}

	for _, node := range nodes {
		if node.ShardCoordinator.SelfId() != core.MetachainShardId {
			continue
		}

		currentBlockHeader := node.BlockChain.GetCurrentBlockHeader()
		if !currentBlockHeader.IsStartOfEpochBlock() {
			continue
		}

		metaHdr := currentBlockHeader.(*block.MetaBlock)
		snapshotsRootHashes[core.MetachainShardId] = append(snapshotsRootHashes[core.MetachainShardId], metaHdr.GetRootHash())
		for _, epochStartData := range metaHdr.EpochStart.LastFinalizedHeaders {
			snapshotsRootHashes[epochStartData.ShardID] = append(snapshotsRootHashes[epochStartData.ShardID], epochStartData.RootHash)
		}
	}
}

func testNodeStateSnapshotAndPruning(
	t *testing.T,
	node *integrationTests.TestProcessorNode,
	snapshotsRootHashes [][]byte,
	prunedRootHashes [][]byte,
) {
	stateTrie := node.TrieContainer.Get([]byte(dataRetriever.UserAccountsUnit.String()))
	assert.Equal(t, 1, len(snapshotsRootHashes))
	for i := range snapshotsRootHashes {
		tr, err := stateTrie.Recreate(holders.NewDefaultRootHashesHolder(snapshotsRootHashes[i]))
		require.Nil(t, err)
		require.NotNil(t, tr)
	}

	assert.Equal(t, 1, len(prunedRootHashes))
	// if pruning is called for a root hash in a different epoch than the commit, then recreate trie should work
	for i := 0; i < len(prunedRootHashes)-1; i++ {
		tr, err := stateTrie.Recreate(holders.NewDefaultRootHashesHolder(prunedRootHashes[i]))
		require.Nil(t, tr)
		require.NotNil(t, err)
	}
}

func TestContinuouslyAccountCodeChanges(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	numOfShards := 1
	nodesPerShard := 1
	numMetachainNodes := 1
	senderShard := uint32(0)
	round := uint64(0)
	nonce := uint64(0)
	valMinting := big.NewInt(1000000000)

	nodes, idxProposers := integrationTests.SetupSyncNodesOneShardAndMeta(nodesPerShard, numMetachainNodes)
	integrationTests.DisplayAndStartNodes(nodes)

	defer integrationTests.CloseProcessorNodes(nodes)

	fmt.Println("Generating private keys for senders...")
	generateCoordinator, _ := sharding.NewMultiShardCoordinator(uint32(numOfShards), 0)
	nrAccounts := 20

	sendersPrivateKeys := make([]crypto.PrivateKey, nrAccounts)
	for i := 0; i < nrAccounts; i++ {
		sendersPrivateKeys[i], _, _ = integrationTests.GenerateSkAndPkInShard(generateCoordinator, senderShard)
	}

	accounts := make([][]byte, len(sendersPrivateKeys))
	for i := range sendersPrivateKeys {
		accounts[i], _ = sendersPrivateKeys[i].GeneratePublic().ToByteArray()
	}

	fmt.Println("Minting sender addresses...")
	integrationTests.CreateMintingForSenders(nodes, senderShard, sendersPrivateKeys, valMinting)

	round = integrationTests.IncrementAndPrintRound(round)
	nonce++
	round, nonce = integrationTests.ProposeAndSyncOneBlock(t, nodes, idxProposers, round, nonce)

	numCodes := 10
	codeMap := getCodeMap(numCodes)
	codeArray := make([][]byte, 0)
	for code := range codeMap {
		codeArray = append(codeArray, []byte(code))
	}

	maxCodeUpdates := 10
	maxCodeDeletes := 10

	shardNode := nodes[0]
	roundsToWait := 50
	for i := 0; i < roundsToWait; i++ {
		numCodeUpdates := rand.Intn(maxCodeUpdates)
		numCodeDeletes := rand.Intn(maxCodeDeletes)

		for j := 0; j < numCodeUpdates; j++ {
			accountIndex := rand.Intn(nrAccounts)
			account, _ := shardNode.AccntState.LoadAccount(accounts[accountIndex])

			updateCode(t, shardNode.AccntState, codeArray, codeMap, account, numCodes)
		}

		for j := 0; j < numCodeDeletes; j++ {
			accountIndex := rand.Intn(nrAccounts)
			account, _ := shardNode.AccntState.LoadAccount(accounts[accountIndex])

			removeCode(t, shardNode.AccntState, codeMap, account)
		}
		_, _ = shardNode.AccntState.Commit()

		round, nonce = integrationTests.ProposeAndSyncOneBlock(t, nodes, idxProposers, round, nonce)

		for j := range codeArray {
			fmt.Printf("%v - %v \n", codeArray[j], codeMap[string(codeArray[j])])
		}

		checkCodeConsistency(t, shardNode, codeMap)
	}
}

func shouldRevert() bool {
	return (rand.Intn(10) % 3) == 0
}

func getCodeMap(numCodes int) map[string]int {
	codeMap := make(map[string]int)
	for i := 0; i < numCodes; i++ {
		code := "code" + strconv.Itoa(i)
		codeMap[code] = 0
	}

	return codeMap
}

func updateCode(
	t *testing.T,
	AccntState state.AccountsAdapter,
	codeArray [][]byte,
	codeMap map[string]int,
	account vmcommon.AccountHandler,
	numCodes int,
) {
	snapshot := AccntState.JournalLen()
	codeIndex := rand.Intn(numCodes)
	code := codeArray[codeIndex]

	oldCode := AccntState.GetCode(account.(state.UserAccountHandler).GetCodeHash())
	account.(state.UserAccountHandler).SetCode(code)
	_ = AccntState.SaveAccount(account)

	if shouldRevert() && snapshot != 0 {
		err := AccntState.RevertToSnapshot(snapshot)
		require.Nil(t, err)
		fmt.Printf("updated code %v to account %v and reverted\n", code, hex.EncodeToString(account.AddressBytes()))
		return
	}

	codeMap[string(code)]++
	if len(oldCode) != 0 {
		codeMap[string(oldCode)]--
	}

	fmt.Printf("updated code %v to account %v \n", code, hex.EncodeToString(account.AddressBytes()))
}

func removeCode(
	t *testing.T,
	AccntState state.AccountsAdapter,
	codeMap map[string]int,
	account vmcommon.AccountHandler,
) {
	snapshot := AccntState.JournalLen()
	code := AccntState.GetCode(account.(state.UserAccountHandler).GetCodeHash())
	account.(state.UserAccountHandler).SetCode(nil)
	_ = AccntState.SaveAccount(account)

	if shouldRevert() && snapshot != 0 {
		err := AccntState.RevertToSnapshot(snapshot)
		require.Nil(t, err)
		fmt.Printf("removed old code %v from account %v and reverted\n", code, hex.EncodeToString(account.AddressBytes()))
		return
	}

	if len(code) != 0 {
		codeMap[string(code)]--
	}

	fmt.Printf("removed old code %v from account %v \n", code, hex.EncodeToString(account.AddressBytes()))
}

func checkCodeConsistency(
	t *testing.T,
	shardNode *integrationTests.TestProcessorNode,
	codeMap map[string]int,
) {
	for code := range codeMap {
		codeHash := integrationTests.TestHasher.Compute(code)
		tr := shardNode.TrieContainer.Get([]byte(dataRetriever.UserAccountsUnit.String()))

		if codeMap[code] != 0 {
			val, _, err := tr.Get(codeHash)
			require.Nil(t, err)
			require.NotNil(t, val)

			var codeEntry state.CodeEntry
			err = integrationTests.TestMarshalizer.Unmarshal(&codeEntry, val)
			require.Nil(t, err)

			require.Equal(t, uint32(codeMap[code]), codeEntry.NumReferences)
		}
	}
}

func TestAccountRemoval(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	numOfShards := 1
	nodesPerShard := 1
	numMetachainNodes := 1
	senderShard := uint32(0)
	round := uint64(0)
	nonce := uint64(0)
	valMinting := big.NewInt(1000000000)

	nodes, idxProposers := integrationTests.SetupSyncNodesOneShardAndMeta(nodesPerShard, numMetachainNodes)
	integrationTests.DisplayAndStartNodes(nodes)

	defer integrationTests.CloseProcessorNodes(nodes)

	fmt.Println("Generating private keys for senders...")
	generateCoordinator, _ := sharding.NewMultiShardCoordinator(uint32(numOfShards), 0)
	nrAccounts := 10000

	sendersPrivateKeys := make([]crypto.PrivateKey, nrAccounts)
	for i := 0; i < nrAccounts; i++ {
		sendersPrivateKeys[i], _, _ = integrationTests.GenerateSkAndPkInShard(generateCoordinator, senderShard)
	}

	accounts := make([][]byte, len(sendersPrivateKeys))
	for i := range sendersPrivateKeys {
		accounts[i], _ = sendersPrivateKeys[i].GeneratePublic().ToByteArray()
	}

	fmt.Println("Minting sender addresses...")
	integrationTests.CreateMintingForSenders(nodes, senderShard, sendersPrivateKeys, valMinting)

	round = integrationTests.IncrementAndPrintRound(round)
	nonce++
	round, nonce = integrationTests.ProposeAndSyncOneBlock(t, nodes, idxProposers, round, nonce)

	shardNode := nodes[0]

	dataTriesRootHashes, codeMap := generateAccounts(shardNode, accounts)

	_, _ = shardNode.AccntState.Commit()
	round, nonce = integrationTests.ProposeAndSyncOneBlock(t, nodes, idxProposers, round, nonce)

	numAccountsToRemove := 2
	roundsToWait := 50

	delayRounds := 10
	for i := 0; i < delayRounds; i++ {
		round, nonce = integrationTests.ProposeAndSyncOneBlock(t, nodes, idxProposers, round, nonce)
	}

	removedAccounts := make(map[int]struct{})
	for i := 0; i < roundsToWait; i++ {
		for j := 0; j < numAccountsToRemove; j++ {
			accountIndex := rand.Intn(nrAccounts)
			removedAccounts[accountIndex] = struct{}{}
			account, err := shardNode.AccntState.GetExistingAccount(accounts[accountIndex])
			if err != nil {
				continue
			}
			code := shardNode.AccntState.GetCode(account.(state.UserAccountHandler).GetCodeHash())

			_ = shardNode.AccntState.RemoveAccount(account.AddressBytes())

			codeMap[string(code)]--
		}

		_, _ = shardNode.AccntState.Commit()
		round, nonce = integrationTests.ProposeAndSyncOneBlock(t, nodes, idxProposers, round, nonce)
		checkCodeConsistency(t, shardNode, codeMap)
	}

	delayRounds = 5
	for i := 0; i < delayRounds; i++ {
		round, nonce = integrationTests.ProposeAndSyncOneBlock(t, nodes, idxProposers, round, nonce)
	}

	checkDataTrieConsistency(t, shardNode.AccntState, removedAccounts, dataTriesRootHashes)
}

func generateAccounts(
	shardNode *integrationTests.TestProcessorNode,
	accounts [][]byte,
) ([][]byte, map[string]int) {
	numCodes := 100
	codeMap := getCodeMap(numCodes)
	codeArray := make([][]byte, 0)
	for code := range codeMap {
		codeArray = append(codeArray, []byte(code))
	}

	dataTriesRootHashes := make([][]byte, 0)
	dataTrieSize := 5
	for i := 0; i < len(accounts); i++ {
		account, _ := shardNode.AccntState.LoadAccount(accounts[i])

		code := codeArray[rand.Intn(numCodes)]
		account.(state.UserAccountHandler).SetCode(code)
		codeMap[string(code)]++

		for j := 0; j < dataTrieSize; j++ {
			_ = account.(state.UserAccountHandler).SaveKeyValue(getDataTrieEntry())
		}

		_ = shardNode.AccntState.SaveAccount(account)

		rootHash := account.(state.UserAccountHandler).GetRootHash()
		dataTriesRootHashes = append(dataTriesRootHashes, rootHash)
	}

	return dataTriesRootHashes, codeMap
}

func getDataTrieEntry() ([]byte, []byte) {
	index := strconv.Itoa(rand.Intn(math.MaxInt32))
	key := []byte("key" + index)
	value := []byte("value" + index)

	return key, value
}

func checkDataTrieConsistency(
	t *testing.T,
	adb state.AccountsAdapter,
	removedAccounts map[int]struct{},
	dataTriesRootHashes [][]byte,
) {
	for i, rootHash := range dataTriesRootHashes {
		_, ok := removedAccounts[i]
		if ok {
			err := adb.RecreateTrie(holders.NewDefaultRootHashesHolder(rootHash))
			assert.NotNil(t, err)
		} else {
			err := adb.RecreateTrie(holders.NewDefaultRootHashesHolder(rootHash))
			require.Nil(t, err)
		}
	}
}

func TestProofAndVerifyProofDataTrie(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	numOfShards := 1
	nodesPerShard := 1
	numMetachainNodes := 1
	senderShard := uint32(0)
	round := uint64(0)
	nonce := uint64(0)

	nodes, idxProposers := integrationTests.SetupSyncNodesOneShardAndMeta(nodesPerShard, numMetachainNodes)
	integrationTests.DisplayAndStartNodes(nodes)

	defer integrationTests.CloseProcessorNodes(nodes)

	generateCoordinator, _ := sharding.NewMultiShardCoordinator(uint32(numOfShards), 0)
	senderPrivateKey, _, _ := integrationTests.GenerateSkAndPkInShard(generateCoordinator, senderShard)
	address, _ := senderPrivateKey.GeneratePublic().ToByteArray()

	shardNode := nodes[0]

	account, _ := shardNode.AccntState.LoadAccount(address)
	numValsInDataTrie := 500
	for i := 0; i < numValsInDataTrie; i++ {
		index := strconv.Itoa(i)
		key := []byte("key" + index)
		value := []byte("value" + index)

		err := account.(state.UserAccountHandler).SaveKeyValue(key, value)
		assert.Nil(t, err)
	}

	err := shardNode.AccntState.SaveAccount(account)
	assert.Nil(t, err)

	_, _ = shardNode.AccntState.Commit()
	_, _ = integrationTests.ProposeAndSyncOneBlock(t, nodes, idxProposers, round, nonce)

	rootHash, _ := shardNode.AccntState.RootHash()
	rootHashHex := hex.EncodeToString(rootHash)
	encodedAddr, _ := shardNode.Node.EncodeAddressPubkey(address)
	account, err = shardNode.AccntState.GetExistingAccount(address)
	assert.Nil(t, err)
	dataTrieRootHashBytes := account.(state.UserAccountHandler).GetRootHash()
	mainTrie, _ := shardNode.AccntState.GetTrie(rootHash)

	for i := 0; i < numValsInDataTrie; i++ {
		index := strconv.Itoa(i)
		keyBytes := []byte("key" + index)
		key := hex.EncodeToString(keyBytes)
		value := []byte("value" + index)

		mainTrieProof, dataTrieProof, errProof := shardNode.Node.GetProofDataTrie(rootHashHex, encodedAddr, key)
		assert.Nil(t, errProof)

		response, errVerify := mainTrie.VerifyProof(rootHash, address, mainTrieProof.Proof)
		assert.Nil(t, errVerify)
		assert.True(t, response)

		response, errVerify = mainTrie.VerifyProof(dataTrieRootHashBytes, keyBytes, dataTrieProof.Proof)
		assert.Nil(t, errVerify)
		assert.True(t, response)
		assert.Equal(t, value, dataTrieProof.Value)
	}
}

func TestTrieDBPruning_PruningOldData(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}
	adb := createAccountsDBTestSetup()
	numAccounts := uint32(10000)
	numIncreaseDecreaseIterations := 100

	rootHashes := make([][]byte, 0)
	rootHash, err := createDummyAccountsWith100REWA(numAccounts, adb)
	require.Nil(t, err)
	rootHashes = append(rootHashes, rootHash)

	for i := 0; i < numIncreaseDecreaseIterations; i++ {
		// change some accounts
		rootHash, err = increaseBalanceForAccountsStartingWithIndex(100, 1000, 10, adb)
		require.Nil(t, err)
		rootHashes = append(rootHashes, rootHash)
		checkAccountsBalances(t, 100, 1000, 110, adb)

		// change same accounts state back
		rootHash, err = decreaseBalanceForAccountsStartingWithIndex(100, 1000, 10, adb)
		require.Nil(t, err)
		rootHashes = append(rootHashes, rootHash)
		adb.PruneTrie(rootHashes[len(rootHashes)-2], state.OldRoot, state.NewPruningHandler(state.EnableDataRemoval))
		checkAccountsBalances(t, 100, 1000, 100, adb)
	}
}

func TestTrieDBPruning_PruningOldDataWithDataTries(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	adb := createAccountsDBTestSetup()
	numAccounts := uint32(100)
	numIncreaseDecreaseIterations := 100
	numAccountsChances := uint32(10)

	rootHashes := make([][]byte, 0)
	rootHash, err := createDummyAccountsWith100REWA(numAccounts, adb)
	require.Nil(t, err)
	rootHashes = append(rootHashes, rootHash)
	rootHash, err = addDataTriesForAccountsStartingWithIndex(10, numAccountsChances, 100, adb)
	require.Nil(t, err)
	rootHashes = append(rootHashes, rootHash)

	for i := 0; i < numIncreaseDecreaseIterations; i++ {
		// change some accounts
		rootHash, err = addDataTriesForAccountsStartingWithIndex(10, numAccountsChances, 10, adb)
		require.Nil(t, err)
		rootHashes = append(rootHashes, rootHash)
		checkAccountsDataTries(t, 10, numAccountsChances, 0, adb)

		// change same accounts state back
		rootHash, err = removeKeysFromAccountsStartingWithIndex(10, numAccountsChances, 10, adb)
		require.Nil(t, err)
		rootHashes = append(rootHashes, rootHash)
		adb.PruneTrie(rootHashes[len(rootHashes)-2], state.OldRoot, state.NewPruningHandler(state.EnableDataRemoval))
		checkAccountsDataTries(t, 10, numAccountsChances, 10, adb)
	}
}

func Test_SnapshotStateRemovesLastSnapshotStartedAfterSnapshotFinished(t *testing.T) {
	t.Parallel()

	lastSnapshotStartedKey := "lastSnapshot"
	epochsNotifier := &mock.EpochStartNotifierStub{
		RegisterHandlerCalled: func(handler epochStart.ActionHandler) {
			header := &block.Header{Epoch: 1}
			handler.EpochStartAction(header)
		},
	}
	tsm := integrationTests.CreateTrieStorageManagerWithPruningStorer(testscommon.NewMultiShardsCoordinatorMock(1), epochsNotifier)
	adb, _ := integrationTests.CreateAccountsDB(0, tsm)
	_ = adb.SetSyncer(&mock.AccountsDBSyncerStub{})
	rootHash, err := addDataTriesForAccountsStartingWithIndex(0, 1, 1, adb)
	assert.Nil(t, err)

	adb.SnapshotState(rootHash, 1)
	for adb.IsSnapshotInProgress() {
		time.Sleep(10 * time.Millisecond)
	}

	val, err := tsm.Get([]byte(lastSnapshotStartedKey))
	assert.Nil(t, val)
	assert.NotNil(t, err)
}

func TestMigrateDataTrieBuiltinFunc(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	t.Run("migrate shard 0 system account", func(t *testing.T) {
		shardId := byte(0)
		nodes, idxProposers, nonce, round := startNodesAndIssueToken(t, 2, shardId)
		defer func() {
			for _, n := range nodes {
				n.Close()
			}
		}()

		valuesBeforeMigration := getValuesFromAccount(t, nodes[shardId].AccntState, core.SystemAccountAddress)
		migrateDataTrieBuiltInFunc(t, nodes, shardId, core.SystemAccountAddress, nonce, round, idxProposers)
		valuesAfterMigration := getValuesFromAccount(t, nodes[shardId].AccntState, core.SystemAccountAddress)

		require.Equal(t, len(valuesBeforeMigration), len(valuesAfterMigration))
		require.True(t, len(valuesAfterMigration) > 0)
		for i := range valuesBeforeMigration {
			require.Equal(t, valuesBeforeMigration[i], valuesAfterMigration[i])
		}
	})
	t.Run("migrate shard 0 user account", func(t *testing.T) {
		shardId := byte(0)
		nodes, idxProposers, nonce, round := startNodesAndIssueToken(t, 2, shardId)
		defer func() {
			for _, n := range nodes {
				n.Close()
			}
		}()

		migrationAddress := nodes[shardId].OwnAccount.Address
		valuesBeforeMigration := getValuesFromAccount(t, nodes[shardId].AccntState, migrationAddress)
		migrateDataTrieBuiltInFunc(t, nodes, shardId, migrationAddress, nonce, round, idxProposers)
		valuesAfterMigration := getValuesFromAccount(t, nodes[shardId].AccntState, migrationAddress)

		require.Equal(t, len(valuesBeforeMigration), len(valuesAfterMigration))
		require.True(t, len(valuesAfterMigration) > 0)
		for i := range valuesBeforeMigration {
			require.Equal(t, valuesBeforeMigration[i], valuesAfterMigration[i])
		}
	})
	t.Run("migrate shard 1 system account", func(t *testing.T) {
		shardId := byte(1)
		nodes, idxProposers, nonce, round := startNodesAndIssueToken(t, 2, shardId)
		defer func() {
			for _, n := range nodes {
				n.Close()
			}
		}()

		valuesBeforeMigration := getValuesFromAccount(t, nodes[shardId].AccntState, core.SystemAccountAddress)
		migrateDataTrieBuiltInFunc(t, nodes, shardId, core.SystemAccountAddress, nonce, round, idxProposers)
		valuesAfterMigration := getValuesFromAccount(t, nodes[shardId].AccntState, core.SystemAccountAddress)

		require.Equal(t, len(valuesBeforeMigration), len(valuesAfterMigration))
		require.True(t, len(valuesAfterMigration) > 0)
		for i := range valuesBeforeMigration {
			require.Equal(t, valuesBeforeMigration[i], valuesAfterMigration[i])
		}
	})
	t.Run("migrate shard 1 user account", func(t *testing.T) {
		shardId := byte(1)
		nodes, idxProposers, nonce, round := startNodesAndIssueToken(t, 2, shardId)
		defer func() {
			for _, n := range nodes {
				n.Close()
			}
		}()

		migrationAddress := nodes[shardId].OwnAccount.Address
		valuesBeforeMigration := getValuesFromAccount(t, nodes[shardId].AccntState, migrationAddress)
		migrateDataTrieBuiltInFunc(t, nodes, shardId, nodes[shardId].OwnAccount.Address, nonce, round, idxProposers)
		valuesAfterMigration := getValuesFromAccount(t, nodes[shardId].AccntState, migrationAddress)

		require.Equal(t, len(valuesBeforeMigration), len(valuesAfterMigration))
		require.True(t, len(valuesAfterMigration) > 0)
		for i := range valuesBeforeMigration {
			require.Equal(t, valuesBeforeMigration[i], valuesAfterMigration[i])
		}
	})
}

func getValuesFromAccount(t *testing.T, adb state.AccountsAdapter, address []byte) [][]byte {
	account, err := adb.GetExistingAccount(address)
	require.Nil(t, err)

	chLeaves := &common.TrieIteratorChannels{
		LeavesChan: make(chan core.KeyValueHolder, common.TrieLeavesChannelDefaultCapacity),
		ErrChan:    errChan.NewErrChanWrapper(),
	}
	err = account.(state.UserAccountHandler).GetAllLeaves(chLeaves, context.Background())
	require.Nil(t, err)

	values := make([][]byte, 0)
	for leaf := range chLeaves.LeavesChan {
		values = append(values, leaf.Value())
	}

	err = chLeaves.ErrChan.ReadFromChanNonBlocking()
	require.Nil(t, err)

	return values
}

func migrateDataTrieBuiltInFunc(
	t *testing.T,
	nodes []*integrationTests.TestProcessorNode,
	shardId byte,
	migrationAddress []byte,
	nonce uint64,
	round uint64,
	leaders []*integrationTests.TestProcessorNode,
) {
	require.True(t, nodes[shardId].EnableEpochsHandler.IsFlagEnabled(common.AutoBalanceDataTriesFlag))
	isMigrated := getAddressMigrationStatus(t, nodes[shardId].AccntState, migrationAddress)
	require.False(t, isMigrated)

	integrationTests.CreateAndSendTransactionWithSenderAccount(nodes[shardId], nodes, big.NewInt(0), nodes[shardId].OwnAccount, getDestAccountAddress(migrationAddress, shardId), core.BuiltInFunctionMigrateDataTrie, 1000000)

	time.Sleep(time.Second)
	nrRoundsToPropagate := 5
	_, _ = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagate, nonce, round)

	isMigrated = getAddressMigrationStatus(t, nodes[shardId].AccntState, migrationAddress)
	require.True(t, isMigrated)
}

func startNodesAndIssueToken(
	t *testing.T,
	numOfShards int,
	issuerShardId byte,
) (leaders []*integrationTests.TestProcessorNode, nodes []*integrationTests.TestProcessorNode, nonce uint64, round uint64) {
	nodesPerShard := 1
	numMetachainNodes := 1

	enableEpochs := config.EnableEpochs{
		GlobalMintBurnDisableEpoch:                  integrationTests.UnreachableEpoch,
		OptimizeGasUsedInCrossMiniBlocksEnableEpoch: integrationTests.UnreachableEpoch,
		ScheduledMiniBlocksEnableEpoch:              integrationTests.UnreachableEpoch,
		MiniBlockPartialExecutionEnableEpoch:        integrationTests.UnreachableEpoch,
		StakingV2EnableEpoch:                        integrationTests.UnreachableEpoch,
		StakeLimitsEnableEpoch:                      integrationTests.UnreachableEpoch,
		StakingV4Step1EnableEpoch:                   integrationTests.UnreachableEpoch,
		StakingV4Step2EnableEpoch:                   integrationTests.UnreachableEpoch,
		StakingV4Step3EnableEpoch:                   integrationTests.UnreachableEpoch,
		AndromedaEnableEpoch:                        integrationTests.UnreachableEpoch,
		AutoBalanceDataTriesEnableEpoch:             1,
	}
	nodes = integrationTests.CreateNodesWithEnableEpochs(
		numOfShards,
		nodesPerShard,
		numMetachainNodes,
		enableEpochs,
	)

	roundsPerEpoch := uint64(5)
	for _, node := range nodes {
		node.EpochStartTrigger.SetRoundsPerEpoch(roundsPerEpoch)
	}

	leaders = make([]*integrationTests.TestProcessorNode, numOfShards+1)
	for i := 0; i < numOfShards; i++ {
		leaders[i] = nodes[i*nodesPerShard]
	}
	leaders[numOfShards] = nodes[numOfShards*nodesPerShard]

	integrationTests.DisplayAndStartNodes(nodes)

	initialVal := int64(10000000000)
	integrationTests.MintAllNodes(nodes, big.NewInt(initialVal))

	round = uint64(0)
	nonce = uint64(0)
	round = integrationTests.IncrementAndPrintRound(round)
	nonce++

	// send token issue
	initialSupply := int64(10000000000)
	ticker := "TCK"
	dcdtCommon.IssueTestTokenWithIssuerAccount(nodes, nodes[issuerShardId].OwnAccount, initialSupply, ticker)

	time.Sleep(time.Second)
	nrRoundsToPropagate := 8
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagate, nonce, round)
	time.Sleep(time.Second)

	tokenIdentifier := string(integrationTests.GetTokenIdentifier(nodes, []byte(ticker)))

	dcdtCommon.CheckAddressHasTokens(t, nodes[issuerShardId].OwnAccount.Address, nodes, []byte(tokenIdentifier), 0, initialSupply)

	return nodes, leaders, nonce, round
}

func getDestAccountAddress(migrationAddress []byte, shardId byte) []byte {
	if bytes.Equal(migrationAddress, core.SystemAccountAddress) && shardId == 0 {
		systemAccountAddress := bytes.Repeat([]byte{255}, 30)
		systemAccountAddress = append(systemAccountAddress, []byte{0, 0}...)
		return systemAccountAddress
	}

	return migrationAddress
}

func getAddressMigrationStatus(t *testing.T, adb state.AccountsAdapter, address []byte) bool {
	account, err := adb.LoadAccount(address)
	require.Nil(t, err)

	userAccount, ok := account.(state.UserAccountHandler)
	require.True(t, ok)

	isMigrated, err := userAccount.DataTrie().IsMigratedToLatestVersion()
	require.Nil(t, err)

	return isMigrated
}

func addDataTriesForAccountsStartingWithIndex(
	startIndex uint32,
	nbAccounts uint32,
	numKeys uint32,
	adb *state.AccountsDB,
) ([]byte, error) {
	for i := startIndex; i < startIndex+nbAccounts; i++ {
		addValuesInAccountDataTrie(i, numKeys, adb)
	}
	return adb.Commit()
}

func removeKeysFromAccountsStartingWithIndex(startIndex uint32,
	nbAccounts uint32,
	numKeys uint32,
	adb *state.AccountsDB,
) ([]byte, error) {
	for i := startIndex; i < startIndex+nbAccounts; i++ {
		removeValuesFromAccountDataTrie(i, numKeys, adb)
	}
	return adb.Commit()
}

func increaseBalanceForAccountsStartingWithIndex(
	startIndex uint32,
	nbAccounts uint32,
	rewaValue uint32,
	adb *state.AccountsDB,
) ([]byte, error) {
	for i := startIndex; i < startIndex+nbAccounts; i++ {
		increaseBalanceForAccountWithIndex(i, rewaValue, adb)
	}
	return adb.Commit()
}

func decreaseBalanceForAccountsStartingWithIndex(
	startIndex uint32,
	nbAccounts uint32,
	rewaValue uint32,
	adb *state.AccountsDB,
) ([]byte, error) {
	for i := startIndex; i < startIndex+nbAccounts; i++ {
		decreaseBalanceForAccountWithIndex(i, rewaValue, adb)
	}
	return adb.Commit()
}

func checkAccountsBalances(t *testing.T, startIndex uint32, nbAccounts uint32, expectedREWAValue uint32, adb *state.AccountsDB) {
	for i := startIndex; i < startIndex+nbAccounts; i++ {
		checkAccountBalance(t, i, expectedREWAValue, adb)
	}
}

func checkAccountBalance(t *testing.T, index uint32, expectedREWAValue uint32, adb *state.AccountsDB) {
	expectedREWAValueDenominated, _ := big.NewInt(0).SetString(fmt.Sprintf("%d", expectedREWAValue)+denomination, 10)

	acc, err := adb.LoadAccount(getDummyAccountAddressFromIndex(index))
	require.Nil(t, err)

	accState := acc.(state.UserAccountHandler)
	actualValue := accState.GetBalance()
	require.Equal(t, expectedREWAValueDenominated, actualValue)
}

func decreaseBalanceForAccountWithIndex(index uint32, rewaValue uint32, adb *state.AccountsDB) {
	rewaValueDenominated, _ := big.NewInt(0).SetString(fmt.Sprintf("%d", rewaValue)+denomination, 10)

	acc, _ := adb.LoadAccount(getDummyAccountAddressFromIndex(index))
	accState := acc.(state.UserAccountHandler)
	_ = accState.SubFromBalance(rewaValueDenominated)
	_ = adb.SaveAccount(accState)
}

func increaseBalanceForAccountWithIndex(index uint32, rewaValue uint32, adb *state.AccountsDB) {
	rewaValueDenominated, _ := big.NewInt(0).SetString(fmt.Sprintf("%d", rewaValue)+denomination, 10)

	acc, _ := adb.LoadAccount(getDummyAccountAddressFromIndex(index))
	accState := acc.(state.UserAccountHandler)
	_ = accState.AddToBalance(rewaValueDenominated)
	_ = adb.SaveAccount(accState)
}

func getDummyAccountAddressFromIndex(index uint32) []byte {
	addrLen := 32
	indexLen := 4
	address := make([]byte, addrLen)
	lastBytes := address[addrLen-indexLen:]
	binary.LittleEndian.PutUint32(lastBytes, index)

	return address
}

func addValuesInAccountDataTrie(index uint32, numKeys uint32, adb *state.AccountsDB) {
	acc, _ := adb.LoadAccount(getDummyAccountAddressFromIndex(index))
	accState := acc.(state.UserAccountHandler)
	for i := 0; i < int(numKeys); i++ {
		k, v := createDummyKeyValue(i)
		_ = accState.SaveKeyValue(k, v)
	}
	_ = adb.SaveAccount(accState)
}

func removeValuesFromAccountDataTrie(index uint32, numKeys uint32, adb *state.AccountsDB) {
	acc, _ := adb.LoadAccount(getDummyAccountAddressFromIndex(index))
	accState := acc.(state.UserAccountHandler)
	for i := 0; i < int(numKeys); i++ {
		k, _ := createDummyKeyValue(i)
		_ = accState.SaveKeyValue(k, nil)
	}
	_ = adb.SaveAccount(accState)
}

func checkAccountsDataTries(t *testing.T, startIndex uint32, nbAccounts uint32, startingKey uint32, adb *state.AccountsDB) {
	for i := startIndex; i < startIndex+nbAccounts; i++ {
		checkAccountsDataTrie(t, i, startingKey, adb)
	}
}

func checkAccountsDataTrie(t *testing.T, index uint32, startingKey uint32, adb *state.AccountsDB) {
	acc, err := adb.LoadAccount(getDummyAccountAddressFromIndex(index))
	require.Nil(t, err)

	numKeys := 100
	accState := acc.(state.UserAccountHandler)
	for i := int(startingKey); i < numKeys; i++ {
		k, v := createDummyKeyValue(i)
		actualValue, _, errKey := accState.RetrieveValue(k)
		require.Nil(t, errKey)
		require.Equal(t, v, actualValue)
	}
}

func createDummyKeyValue(index int) ([]byte, []byte) {
	hasher := sha256.NewSha256()
	key := hasher.Compute(fmt.Sprintf("%d", index))
	value := hasher.Compute(string(key))
	return key, value
}

func createDummyAccountsWith100REWA(numAccounts uint32, adb *state.AccountsDB) ([]byte, error) {
	val100Denominated, _ := big.NewInt(0).SetString("100"+denomination, 10)

	for i := 0; i < int(numAccounts); i++ {
		addr := getDummyAccountAddressFromIndex(uint32(i))
		acc, _ := adb.LoadAccount(addr)
		accState := acc.(state.UserAccountHandler)
		_ = accState.AddToBalance(val100Denominated)
		_ = adb.SaveAccount(accState)
	}

	return adb.Commit()
}

func createAccountsDBTestSetup() *state.AccountsDB {
	generalCfg := config.TrieStorageManagerConfig{
		PruningBufferLen:      1000,
		SnapshotsBufferLen:    10,
		SnapshotsGoroutineNum: 1,
	}
	evictionWaitListSize := uint(100)
	ewlArgs := evictionWaitingList.MemoryEvictionWaitingListArgs{
		RootHashesSize: evictionWaitListSize,
		HashesSize:     evictionWaitListSize * 100,
	}
	ewl, _ := evictionWaitingList.NewMemoryEvictionWaitingList(ewlArgs)
	args := testStorage.GetStorageManagerArgs()
	args.GeneralConfig = generalCfg
	trieStorage, _ := trie.NewTrieStorageManager(args)
	maxTrieLevelInMemory := uint(5)
	tr, _ := trie.NewTrie(trieStorage, integrationTests.TestMarshalizer, integrationTests.TestHasher, &enableEpochsHandlerMock.EnableEpochsHandlerStub{}, maxTrieLevelInMemory)
	spm, _ := storagePruningManager.NewStoragePruningManager(ewl, 10)
	argsAccCreator := factory.ArgsAccountCreator{
		Hasher:              integrationTests.TestHasher,
		Marshaller:          integrationTests.TestMarshalizer,
		EnableEpochsHandler: &enableEpochsHandlerMock.EnableEpochsHandlerStub{},
	}
	accCreator, _ := factory.NewAccountCreator(argsAccCreator)

	snapshotsManager, _ := state.NewSnapshotsManager(state.ArgsNewSnapshotsManager{
		ProcessingMode:       common.Normal,
		Marshaller:           &marshallerMock.MarshalizerMock{},
		AddressConverter:     &testscommon.PubkeyConverterMock{},
		ProcessStatusHandler: &testscommon.ProcessStatusHandlerStub{},
		StateMetrics:         &stateMock.StateMetricsStub{},
		AccountFactory:       accCreator,
		ChannelsProvider:     iteratorChannelsProvider.NewUserStateIteratorChannelsProvider(),
		LastSnapshotMarker:   lastSnapshotMarker.NewLastSnapshotMarker(),
		StateStatsHandler:    statistics.NewStateStatistics(),
	})

	argsAccountsDB := state.ArgsAccountsDB{
		Trie:                  tr,
		Hasher:                integrationTests.TestHasher,
		Marshaller:            integrationTests.TestMarshalizer,
		AccountFactory:        accCreator,
		StoragePruningManager: spm,
		AddressConverter:      &testscommon.PubkeyConverterMock{},
		SnapshotsManager:      snapshotsManager,
	}
	adb, _ := state.NewAccountsDB(argsAccountsDB)

	return adb
}
