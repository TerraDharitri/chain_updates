package hostCoretest

import (
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data/vm"
	"github.com/TerraDharitri/drt-go-chain-scenario/worldmock"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	test "github.com/TerraDharitri/drt-go-chain-vm/testcommon"
	"github.com/TerraDharitri/drt-go-chain-vm/vmhost"

	"github.com/stretchr/testify/require"
)

func Test_RunDEXPairBenchmark(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	owner := vmhost.MakeTestSCAddress("owner")
	user := vmhost.MakeTestSCAddress("user")

	_, host, mex := setupMEXPair(t, owner, user)

	mex.AddLiquidity(user, 1_000_000, 1, 1_000_000, 1)

	numBatches := 10
	numSwapsPerBatch := 2000

	wrewaToMexDCDT, wrewaToMexDCDTSwap := mex.CreateSwapVMInputs(mex.WREWAToken, 100, mex.MEXToken, 1)
	mexToWrewaDCDT, mexToWrewaSwap := mex.CreateSwapVMInputs(mex.MEXToken, 100, mex.WREWAToken, 1)

	for batch := 0; batch < numBatches; batch++ {
		start := time.Now()
		for i := 0; i < numSwapsPerBatch/2; i++ {
			mex.ExecuteSwap(wrewaToMexDCDT, wrewaToMexDCDTSwap)
			mex.ExecuteSwap(mexToWrewaDCDT, mexToWrewaSwap)
		}
		elapsedTime := time.Since(start)
		logBenchmark.Trace(
			"swap batch finished",
			"numSwapsPerBatch",
			numSwapsPerBatch,
			"duration",
			elapsedTime,
		)
	}

	defer func() {
		host.Reset()
	}()
}

func setupMEXPair(t *testing.T, owner Address, user Address) (*worldmock.MockWorld, vmhost.VMHost, *MEXSetup) {
	world, ownerAccount, host, err := prepare(t, owner)
	require.Nil(t, err)

	userAccount := world.AcctMap.CreateAccount(user, world)
	userAccount.Balance = big.NewInt(100)
	mex := NewMEXSetup(t, host, world, ownerAccount, userAccount)
	mex.Deploy()

	mex.ApplyInitialSetup()

	return world, host, mex
}

type MEXSetup struct {
	WREWAToken               []byte
	MEXToken                 []byte
	LPToken                  []byte
	OwnerAccount             *worldmock.Account
	OwnerAddress             Address
	RouterAddress            Address
	PairAddress              Address
	TotalFeePercent          uint64
	SpecialFeePercent        uint64
	MaxObservationsPerRecord int
	Code                     []byte
	UserAccount              *worldmock.Account
	UserWREWABalance         uint64
	UserMEXBalance           uint64

	T     *testing.T
	Host  vmhost.VMHost
	World *worldmock.MockWorld
}

func NewMEXSetup(
	t *testing.T,
	host vmhost.VMHost,
	world *worldmock.MockWorld,
	ownerAccount *worldmock.Account,
	userAccount *worldmock.Account,
) *MEXSetup {
	return &MEXSetup{
		WREWAToken:               []byte("WREWA-abcdef"),
		MEXToken:                 []byte("MEX-abcdef"),
		LPToken:                  []byte("LPTOK-abcdef"),
		OwnerAccount:             ownerAccount,
		OwnerAddress:             ownerAccount.Address,
		RouterAddress:            ownerAccount.Address,
		PairAddress:              test.MakeTestSCAddress("pairSC"),
		TotalFeePercent:          300,
		SpecialFeePercent:        50,
		MaxObservationsPerRecord: 10,
		Code:                     test.GetTestSCCode("pair", "../../"),
		UserAccount:              userAccount,
		UserWREWABalance:         5_000_000_000,
		UserMEXBalance:           5_000_000_000,

		T:     t,
		Host:  host,
		World: world,
	}
}

func (mex *MEXSetup) Deploy() {
	t := mex.T
	host := mex.Host
	world := mex.World

	vmInput := test.CreateTestContractCreateInputBuilder().
		WithCallerAddr(mex.OwnerAddress).
		WithContractCode(mex.Code).
		WithContractCodeMetadata([]byte{5, 6}).
		WithArguments(
			mex.WREWAToken,
			mex.MEXToken,
			mex.OwnerAddress,
			mex.RouterAddress,
			big.NewInt(int64(mex.TotalFeePercent)).Bytes(),
			big.NewInt(int64(mex.SpecialFeePercent)).Bytes(),
		).
		WithGasProvided(math.MaxInt64).
		Build()

	world.NewAddressMocks = append(world.NewAddressMocks, &worldmock.NewAddressMock{
		CreatorAddress: mex.OwnerAddress,
		CreatorNonce:   mex.OwnerAccount.Nonce,
		NewAddress:     mex.PairAddress,
	})

	mex.OwnerAccount.Nonce++ // nonce increases before deploy
	vmOutput, err := host.RunSmartContractCreate(vmInput)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, "", vmOutput.ReturnMessage)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	_ = world.UpdateAccounts(vmOutput.OutputAccounts, nil)
}

func (mex *MEXSetup) ApplyInitialSetup() {
	mex.setLPToken()
	mex.setActiveState()
	mex.setMaxObservationsPerRecord()
	mex.setRequiredTokenRoles()
	mex.setDCDTBalances()
}

func (mex *MEXSetup) setLPToken() {
	t := mex.T
	host := mex.Host
	world := mex.World

	vmInput := test.CreateTestContractCallInputBuilder().
		WithCallerAddr(mex.OwnerAddress).
		WithRecipientAddr(mex.PairAddress).
		WithFunction("setLpTokenIdentifier").
		WithArguments(mex.LPToken).
		WithGasProvided(math.MaxInt64).
		Build()

	vmOutput, err := host.RunSmartContractCall(vmInput)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, "", vmOutput.ReturnMessage)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	_ = world.UpdateAccounts(vmOutput.OutputAccounts, nil)
}

func (mex *MEXSetup) setActiveState() {
	t := mex.T
	host := mex.Host
	world := mex.World

	vmInput := test.CreateTestContractCallInputBuilder().
		WithCallerAddr(mex.OwnerAddress).
		WithRecipientAddr(mex.PairAddress).
		WithFunction("resume").
		WithGasProvided(math.MaxInt64).
		Build()

	vmOutput, err := host.RunSmartContractCall(vmInput)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, "", vmOutput.ReturnMessage)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	_ = world.UpdateAccounts(vmOutput.OutputAccounts, nil)
}

func (mex *MEXSetup) setMaxObservationsPerRecord() {
	t := mex.T
	host := mex.Host
	world := mex.World

	vmInput := test.CreateTestContractCallInputBuilder().
		WithCallerAddr(mex.OwnerAddress).
		WithRecipientAddr(mex.PairAddress).
		WithFunction("setMaxObservationsPerRecord").
		WithArguments(big.NewInt(int64(mex.MaxObservationsPerRecord)).Bytes()).
		WithGasProvided(math.MaxInt64).
		Build()

	vmOutput, err := host.RunSmartContractCall(vmInput)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, "", vmOutput.ReturnMessage)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	_ = world.UpdateAccounts(vmOutput.OutputAccounts, nil)
}

func (mex *MEXSetup) setRequiredTokenRoles() {
	world := mex.World
	pairAccount := world.AcctMap.GetAccount(mex.PairAddress)

	roles := []string{core.DCDTRoleLocalMint, core.DCDTRoleLocalBurn}
	_ = pairAccount.SetTokenRolesAsStrings(mex.LPToken, roles)
}

func (mex *MEXSetup) setDCDTBalances() {
	_ = mex.UserAccount.SetTokenBalanceUint64(mex.WREWAToken, 0, mex.UserWREWABalance)
	_ = mex.UserAccount.SetTokenBalanceUint64(mex.MEXToken, 0, mex.UserMEXBalance)
}

func (mex *MEXSetup) AddLiquidity(
	userAddress Address,
	WREWAAmount uint64,
	minWREWAAmount uint64,
	MEXAmount uint64,
	minMEXAmount uint64,
) {
	t := mex.T
	host := mex.Host
	world := mex.World

	vmInputBuiler := test.CreateTestContractCallInputBuilder().
		WithCallerAddr(mex.UserAccount.Address).
		WithRecipientAddr(mex.PairAddress).
		WithFunction("addLiquidity").
		WithArguments(
			big.NewInt(int64(minWREWAAmount)).Bytes(),
			big.NewInt(int64(minMEXAmount)).Bytes(),
		).
		WithGasProvided(math.MaxInt64)

	vmInputBuiler.
		WithDCDTTokenName(mex.WREWAToken).
		WithDCDTValue(big.NewInt(int64(WREWAAmount))).
		NextDCDTTransfer().
		WithDCDTTokenName(mex.MEXToken).
		WithDCDTValue(big.NewInt(int64(MEXAmount)))

	vmInput := vmInputBuiler.Build()

	addLiquidityDCDT := mex.createMultiDCDTTransferVMInput(
		vmInput.CallerAddr,
		vmInput.RecipientAddr,
		vmInput.DCDTTransfers,
	)

	mex.performMultiDCDTTransfer(addLiquidityDCDT)

	vmOutput, err := host.RunSmartContractCall(vmInput)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, "", vmOutput.ReturnMessage)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	_ = world.UpdateAccounts(vmOutput.OutputAccounts, nil)
}

func (mex *MEXSetup) CreateSwapVMInputs(
	leftToken []byte,
	leftAmount uint64,
	rightToken []byte,
	rightAmount uint64,
) (*vmcommon.ContractCallInput, *vmcommon.ContractCallInput) {
	vmInputBuiler := test.CreateTestContractCallInputBuilder().
		WithCallerAddr(mex.UserAccount.Address).
		WithRecipientAddr(mex.PairAddress)

	vmInputBuiler.
		WithDCDTTokenName(leftToken).
		WithDCDTValue(big.NewInt(int64(leftAmount))).
		WithFunction("swapTokensFixedInput").
		WithArguments(
			rightToken,
			big.NewInt(int64(rightAmount)).Bytes(),
		).
		WithGasProvided(math.MaxInt64)

	vmInput := vmInputBuiler.Build()

	multiTransferInput := mex.createMultiDCDTTransferVMInput(
		vmInput.CallerAddr,
		vmInput.RecipientAddr,
		vmInput.DCDTTransfers,
	)

	return multiTransferInput, vmInput

}

func (mex *MEXSetup) ExecuteSwap(
	multiTransferInput *vmcommon.ContractCallInput,
	vmInput *vmcommon.ContractCallInput,
) {
	t := mex.T
	host := mex.Host
	world := mex.World

	mex.performMultiDCDTTransfer(multiTransferInput)

	vmOutput, err := host.RunSmartContractCall(vmInput)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, "", vmOutput.ReturnMessage)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	_ = world.UpdateAccounts(vmOutput.OutputAccounts, nil)

	err = world.CommitChanges()
	require.Nil(t, err)
}

func (mex *MEXSetup) createMultiDCDTTransferVMInput(
	sender Address,
	receiver Address,
	dcdtTransfers []*vmcommon.DCDTTransfer,
) *vmcommon.ContractCallInput {
	nrTransfers := len(dcdtTransfers)
	nrTransfersAsBytes := big.NewInt(0).SetUint64(uint64(nrTransfers)).Bytes()

	multiTransferInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   make([][]byte, 0),
			CallValue:   big.NewInt(0),
			CallType:    vm.DirectCall,
			GasPrice:    1,
			GasProvided: math.MaxInt64,
			GasLocked:   0,
		},
		RecipientAddr:     sender,
		Function:          core.BuiltInFunctionMultiDCDTNFTTransfer,
		AllowInitFunction: false,
	}
	multiTransferInput.Arguments = append(multiTransferInput.Arguments, receiver, nrTransfersAsBytes)

	for i := 0; i < nrTransfers; i++ {
		token := dcdtTransfers[i].DCDTTokenName
		nonceAsBytes := big.NewInt(0).SetUint64(dcdtTransfers[i].DCDTTokenNonce).Bytes()
		value := dcdtTransfers[i].DCDTValue.Bytes()

		multiTransferInput.Arguments = append(multiTransferInput.Arguments, token, nonceAsBytes, value)
	}

	return multiTransferInput
}

func (mex *MEXSetup) performMultiDCDTTransfer(
	multiTransferInput *vmcommon.ContractCallInput,
) {
	t := mex.T
	world := mex.World

	vmOutput, err := world.BuiltinFuncs.ProcessBuiltInFunction(multiTransferInput)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, "", vmOutput.ReturnMessage)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	_ = world.UpdateAccounts(vmOutput.OutputAccounts, nil)
}
