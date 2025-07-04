package upgrades

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/integrationTests"
	"github.com/TerraDharitri/drt-go-chain/integrationTests/vm/wasm"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/factory"
)

func TestUpgrades_Hello(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	context := wasm.SetupTestContext(t)
	defer context.Close()

	fmt.Println("Deploy v1")

	context.ScCodeMetadata.Upgradeable = true
	err := context.DeploySC("../testdata/hello-v1/output/answer.wasm", "")
	require.Nil(t, err)
	require.Equal(t, uint64(24), context.QuerySCInt("getUltimateAnswer", [][]byte{}))

	fmt.Println("Upgrade to v2")

	err = context.UpgradeSC("../testdata/hello-v2/output/answer.wasm", "")
	require.Nil(t, err)
	require.Equal(t, uint64(42), context.QuerySCInt("getUltimateAnswer", [][]byte{}))

	fmt.Println("Upgrade to v3")

	err = context.UpgradeSC("../testdata/hello-v3/output/answer.wasm", "")
	require.Nil(t, err)
	require.Equal(t, "forty-two", context.QuerySCString("getUltimateAnswer", [][]byte{}))
}

func TestUpgrades_HelloDoesNotUpgradeWhenNotUpgradeable(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	context := wasm.SetupTestContext(t)
	defer context.Close()

	fmt.Println("Deploy v1")

	context.ScCodeMetadata.Upgradeable = false
	err := context.DeploySC("../testdata/hello-v1/output/answer.wasm", "")
	require.Nil(t, err)
	require.Equal(t, uint64(24), context.QuerySCInt("getUltimateAnswer", [][]byte{}))

	fmt.Println("Upgrade to v2 will not be performed")

	err = context.UpgradeSC("../testdata/hello-v2/output/answer.wasm", "")
	require.Equal(t, process.ErrUpgradeNotAllowed, err)
	require.Equal(t, uint64(24), context.QuerySCInt("getUltimateAnswer", [][]byte{}))
}

func TestUpgrades_HelloUpgradesToNotUpgradeable(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	context := wasm.SetupTestContext(t)
	defer context.Close()

	fmt.Println("Deploy v1")

	context.ScCodeMetadata.Upgradeable = true
	err := context.DeploySC("../testdata/hello-v1/output/answer.wasm", "")
	require.Nil(t, err)
	require.Equal(t, uint64(24), context.QuerySCInt("getUltimateAnswer", [][]byte{}))

	fmt.Println("Upgrade to v2, becomes not upgradeable")

	context.ScCodeMetadata.Upgradeable = false
	err = context.UpgradeSC("../testdata/hello-v2/output/answer.wasm", "")
	require.Nil(t, err)
	require.Equal(t, uint64(42), context.QuerySCInt("getUltimateAnswer", [][]byte{}))

	fmt.Println("Upgrade to v3, should not be possible")

	err = context.UpgradeSC("../testdata/hello-v3/output/answer.wasm", "")
	require.Equal(t, process.ErrUpgradeNotAllowed, err)
	require.Equal(t, uint64(42), context.QuerySCInt("getUltimateAnswer", [][]byte{}))
}

func TestUpgrades_ParentAndChildContracts(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	context := wasm.SetupTestContext(t)
	defer context.Close()

	var parentAddress []byte
	var childAddress []byte
	owner := &context.Owner

	fmt.Println("Deploy parent")

	err := context.DeploySC("../testdata/upgrades-parent/output/parent.wasm", "")
	require.Nil(t, err)
	require.Equal(t, uint64(45), context.QuerySCInt("getUltimateAnswer", [][]byte{}))
	parentAddress = context.ScAddress

	fmt.Println("Deploy child v1")

	childInitialCode := wasm.GetSCCode("../testdata/hello-v1/output/answer.wasm")
	err = context.ExecuteSC(owner, "createChild@"+childInitialCode)
	require.Nil(t, err)

	fmt.Println("Aquire child address, do query")

	childAddress = context.QuerySCBytes("getChildAddress", [][]byte{})
	context.ScAddress = childAddress
	require.Equal(t, uint64(24), context.QuerySCInt("getUltimateAnswer", [][]byte{}))

	fmt.Println("Deploy child v2")
	context.ScAddress = parentAddress

	childUpgradedCode := wasm.GetSCCode("../testdata/hello-v2/output/answer.wasm")
	context.GasLimit = 21700000000
	err = context.ExecuteSC(owner, "upgradeChild@"+childUpgradedCode)
	require.Nil(t, err)

	context.ScAddress = childAddress
	require.Equal(t, uint64(42), context.QuerySCInt("getUltimateAnswer", [][]byte{}))
}

func TestUpgrades_HelloCannotBeUpgradedByNonOwner(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	context := wasm.SetupTestContext(t)
	defer context.Close()

	fmt.Println("Deploy v1")

	context.ScCodeMetadata.Upgradeable = true
	err := context.DeploySC("../testdata/hello-v1/output/answer.wasm", "")
	require.Nil(t, err)
	require.Equal(t, uint64(24), context.QuerySCInt("getUltimateAnswer", [][]byte{}))

	fmt.Println("Upgrade to v2 will not be performed")

	// Alice states that she is the owner of the contract (though she is not)
	context.Owner = context.Alice
	err = context.UpgradeSC("../testdata/hello-v2/output/answer.wasm", "")
	require.Equal(t, process.ErrUpgradeNotAllowed, err)
	require.Equal(t, uint64(24), context.QuerySCInt("getUltimateAnswer", [][]byte{}))
}

func TestUpgrades_CounterCannotBeUpgradedByNonOwner(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	context := wasm.SetupTestContext(t)
	defer context.Close()

	context.ScCodeMetadata.Upgradeable = true
	err := context.DeploySC("../testdata/counter/output/counter.wasm", "")
	require.Nil(t, err)
	require.Equal(t, uint64(1), context.QuerySCInt("get", [][]byte{}))

	err = context.ExecuteSC(&context.Alice, "increment")
	require.Nil(t, err)
	require.Equal(t, uint64(2), context.QuerySCInt("get", [][]byte{}))

	// Alice states that she is the owner of the contract (though she is not)
	// Neither code, nor storage get modified
	context.Owner = context.Alice
	err = context.UpgradeSC("../testdata/counter/output/counter.wasm", "")
	require.Equal(t, process.ErrUpgradeNotAllowed, err)
	require.Equal(t, uint64(2), context.QuerySCInt("get", [][]byte{}))
}

func TestUpgrades_HelloTrialAndError(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	network := integrationTests.NewMiniNetwork()
	defer network.Stop()

	alice := network.AddUser(big.NewInt(10000000000000))
	bob := network.AddUser(big.NewInt(10000000000000))

	network.Start()

	deployTxData := fmt.Sprintf("%s@%s@0100", wasm.GetSCCode("../testdata/hello-v1/output/answer.wasm"), hex.EncodeToString(factory.WasmVirtualMachine))
	upgradeTxData := fmt.Sprintf("upgradeContract@%s@0100", wasm.GetSCCode("../testdata/hello-v2/output/answer.wasm"))

	// Deploy the smart contract. Alice is the owner
	_, err := network.SendTransaction(
		alice.Address,
		make([]byte, 32),
		big.NewInt(0),
		deployTxData,
		100000,
	)
	require.Nil(t, err)

	scAddress, _ := network.ShardNode.BlockchainHook.NewAddress(alice.Address, 0, factory.WasmVirtualMachine)
	network.Continue(t, 2)
	require.Equal(t, []byte{24}, query(t, network.ShardNode, scAddress, "getUltimateAnswer"))

	// Upgrade as Bob - upgrade should fail, since Alice is the owner
	_, err = network.SendTransaction(
		bob.Address,
		scAddress,
		big.NewInt(0),
		upgradeTxData,
		100000,
	)
	require.Nil(t, err)

	network.Continue(t, 2)
	require.Equal(t, []byte{24}, query(t, network.ShardNode, scAddress, "getUltimateAnswer"))

	// Now upgrade as Alice, should work
	_, err = network.SendTransaction(
		alice.Address,
		scAddress,
		big.NewInt(0),
		upgradeTxData,
		100000,
	)
	require.Nil(t, err)

	network.Continue(t, 2)
	require.Equal(t, []byte{42}, query(t, network.ShardNode, scAddress, "getUltimateAnswer"))
}

func TestUpgrades_CounterTrialAndError(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	network := integrationTests.NewMiniNetwork()
	defer network.Stop()

	alice := network.AddUser(big.NewInt(10000000000000))
	bob := network.AddUser(big.NewInt(10000000000000))

	network.Start()

	deployTxData := fmt.Sprintf("%s@%s@0100", wasm.GetSCCode("../testdata/counter/output/counter.wasm"), hex.EncodeToString(factory.WasmVirtualMachine))
	upgradeTxData := fmt.Sprintf("upgradeContract@%s@0100", wasm.GetSCCode("../testdata/counter/output/counter.wasm"))

	// Deploy the smart contract. Alice is the owner
	_, err := network.SendTransaction(
		alice.Address,
		make([]byte, 32),
		big.NewInt(0),
		deployTxData,
		100000,
	)
	require.Nil(t, err)

	scAddress, _ := network.ShardNode.BlockchainHook.NewAddress(alice.Address, 0, factory.WasmVirtualMachine)
	network.Continue(t, 2)
	require.Equal(t, []byte{1}, query(t, network.ShardNode, scAddress, "get"))

	// Increment the counter (could be either Bob or Alice)
	_, err = network.SendTransaction(
		alice.Address,
		scAddress,
		big.NewInt(0),
		"increment",
		100000,
	)
	require.Nil(t, err)

	network.Continue(t, 2)
	require.Equal(t, []byte{2}, query(t, network.ShardNode, scAddress, "get"))

	// Upgrade as Bob - upgrade should fail, since Alice is the owner (counter.init() not executed, state not reset)
	_, err = network.SendTransaction(
		bob.Address,
		scAddress,
		big.NewInt(0),
		upgradeTxData,
		100000,
	)
	require.Nil(t, err)

	network.Continue(t, 2)
	require.Equal(t, []byte{2}, query(t, network.ShardNode, scAddress, "get"))

	// Now upgrade as Alice, should work (state is reset by counter.init())
	_, err = network.SendTransaction(
		alice.Address,
		scAddress,
		big.NewInt(0),
		upgradeTxData,
		100000,
	)
	require.Nil(t, err)

	network.Continue(t, 2)
	require.Equal(t, []byte{1}, query(t, network.ShardNode, scAddress, "get"))
}

func query(t *testing.T, node *integrationTests.TestProcessorNode, scAddress []byte, function string) []byte {
	scQuery := node.SCQueryService
	vmOutput, _, err := scQuery.ExecuteQuery(&process.SCQuery{
		ScAddress: scAddress,
		FuncName:  function,
		Arguments: [][]byte{},
	})

	require.Nil(t, err)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	return vmOutput.ReturnData[0]
}
