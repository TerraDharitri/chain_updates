package transactionEvaluator

import (
	"context"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/common/errChan"
	"github.com/TerraDharitri/drt-go-chain/state"
	"github.com/TerraDharitri/drt-go-chain/state/parsers"
	stateMock "github.com/TerraDharitri/drt-go-chain/testscommon/state"
	"github.com/stretchr/testify/require"
)

func TestNewReadOnlyAccountsDB_NilOriginalAccountsDBShouldErr(t *testing.T) {
	t.Parallel()

	simAccountsDB, err := NewSimulationAccountsDB(nil)
	require.True(t, check.IfNil(simAccountsDB))
	require.Equal(t, ErrNilAccountsAdapter, err)
}

func TestNewReadOnlyAccountsDB(t *testing.T) {
	t.Parallel()

	simAccountsDB, err := NewSimulationAccountsDB(&stateMock.AccountsStub{})
	require.False(t, check.IfNil(simAccountsDB))
	require.NoError(t, err)
}

func TestReadOnlyAccountsDB_WriteOperationsShouldNotCalled(t *testing.T) {
	t.Parallel()

	failErrMsg := "this function should have not be called"
	accDb := &stateMock.AccountsStub{
		SaveAccountCalled: func(account vmcommon.AccountHandler) error {
			t.Errorf("%s", failErrMsg)
			return nil
		},
		RemoveAccountCalled: func(_ []byte) error {
			t.Errorf("%s", failErrMsg)
			return nil
		},
		CommitCalled: func() ([]byte, error) {
			t.Errorf("%s", failErrMsg)
			return nil, nil
		},
		RevertToSnapshotCalled: func(_ int) error {
			t.Errorf("%s", failErrMsg)
			return nil
		},
		RecreateTrieCalled: func(_ common.RootHashHolder) error {
			t.Errorf("%s", failErrMsg)
			return nil
		},
		PruneTrieCalled: func(_ []byte, _ state.TriePruningIdentifier, _ state.PruningHandler) {
			t.Errorf("%s", failErrMsg)
		},
		CancelPruneCalled: func(_ []byte, _ state.TriePruningIdentifier) {
			t.Errorf("%s", failErrMsg)
		},
		SnapshotStateCalled: func(_ []byte, _ uint32) {
			t.Errorf("%s", failErrMsg)
		},
		RecreateAllTriesCalled: func(_ []byte) (map[string]common.Trie, error) {
			t.Errorf("%s", failErrMsg)
			return nil, nil
		},
	}

	simAccountsDB, _ := NewSimulationAccountsDB(accDb)
	require.NotNil(t, simAccountsDB)

	err := simAccountsDB.SaveAccount(nil)
	require.NoError(t, err)

	err = simAccountsDB.RemoveAccount(nil)
	require.NoError(t, err)

	_, err = simAccountsDB.Commit()
	require.NoError(t, err)

	err = simAccountsDB.RevertToSnapshot(0)
	require.NoError(t, err)

	err = simAccountsDB.RecreateTrie(nil)
	require.NoError(t, err)

	simAccountsDB.PruneTrie(nil, state.NewRoot, state.NewPruningHandler(state.EnableDataRemoval))

	simAccountsDB.CancelPrune(nil, state.NewRoot)

	simAccountsDB.SnapshotState(nil, 0)

	_, err = simAccountsDB.RecreateAllTries(nil)
	require.NoError(t, err)
}

func TestReadOnlyAccountsDB_ReadOperationsShouldWork(t *testing.T) {
	t.Parallel()

	expectedAcc := &stateMock.AccountWrapMock{}
	expectedJournalLen := 37
	expectedRootHash := []byte("root")

	accDb := &stateMock.AccountsStub{
		GetExistingAccountCalled: func(_ []byte) (vmcommon.AccountHandler, error) {
			return expectedAcc, nil
		},
		LoadAccountCalled: func(_ []byte) (vmcommon.AccountHandler, error) {
			return expectedAcc, nil
		},
		JournalLenCalled: func() int {
			return expectedJournalLen
		},
		RootHashCalled: func() ([]byte, error) {
			return expectedRootHash, nil
		},
		IsPruningEnabledCalled: func() bool {
			return true
		},
	}

	simAccountsDB, _ := NewSimulationAccountsDB(accDb)
	require.NotNil(t, simAccountsDB)

	actualAcc, err := simAccountsDB.GetExistingAccount(nil)
	require.NoError(t, err)
	require.Equal(t, expectedAcc, actualAcc)

	actualAcc, err = simAccountsDB.LoadAccount(nil)
	require.NoError(t, err)
	require.Equal(t, expectedAcc, actualAcc)

	actualJournalLen := simAccountsDB.JournalLen()
	require.Equal(t, expectedJournalLen, actualJournalLen)

	actualRootHash, err := simAccountsDB.RootHash()
	require.NoError(t, err)
	require.Equal(t, expectedRootHash, actualRootHash)

	actualIsPruningEnabled := simAccountsDB.IsPruningEnabled()
	require.Equal(t, true, actualIsPruningEnabled)

	allLeaves := &common.TrieIteratorChannels{
		LeavesChan: make(chan core.KeyValueHolder),
		ErrChan:    errChan.NewErrChanWrapper(),
	}
	err = simAccountsDB.GetAllLeaves(allLeaves, context.Background(), nil, parsers.NewMainTrieLeafParser())
	require.NoError(t, err)

	err = allLeaves.ErrChan.ReadFromChanNonBlocking()
	require.NoError(t, err)
}
