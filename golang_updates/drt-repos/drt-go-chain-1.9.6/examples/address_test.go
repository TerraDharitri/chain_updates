package examples

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/display"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain/vm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHexAddressToBech32Address(t *testing.T) {
	t.Parallel()

	hexEncodedAddress := "af006ece83473104ea91f7ff5605c4c1742f7214a1f46be299e30ee2e8707169"

	hexEncodedAddressBytes, err := hex.DecodeString(hexEncodedAddress)
	require.NoError(t, err)

	bech32Address, err := addressEncoder.Encode(hexEncodedAddressBytes)
	require.NoError(t, err)
	require.Equal(t, "drt14uqxan5rgucsf6537ll4vpwyc96z7us5586xhc5euv8w96rsw95s5rd7km", bech32Address)
}

func TestBech32AddressToHexAddress(t *testing.T) {
	t.Parallel()

	bech32Address := "drt14uqxan5rgucsf6537ll4vpwyc96z7us5586xhc5euv8w96rsw95s5rd7km"

	bech32AddressBytes, err := addressEncoder.Decode(bech32Address)
	require.NoError(t, err)

	hexEncodedAddress := hex.EncodeToString(bech32AddressBytes)
	require.Equal(t, "af006ece83473104ea91f7ff5605c4c1742f7214a1f46be299e30ee2e8707169", hexEncodedAddress)
}

func TestShardOfAddress(t *testing.T) {
	t.Parallel()

	// the shard of an address depends on the number of shards in the chain. The same address does not necessarily
	// belong to the same shard in a chain with a different number of shards.

	numberOfShards := uint32(3)
	shardCoordinator, err := sharding.NewMultiShardCoordinator(numberOfShards, 0)
	require.NoError(t, err)

	require.Equal(t, uint32(0), computeShardID(t, "drt1gn0y4l4rgkf2e7dg74u3nnugr7uycw5jwa44tlnqg2kxa37dr2kq8n3d0k", shardCoordinator))
	require.Equal(t, uint32(1), computeShardID(t, "drt1x23lzn8483xs2su4fak0r0dqx6w38enpmmqf2yrkylwq7mfnvyhsmueha6", shardCoordinator))
	require.Equal(t, uint32(2), computeShardID(t, "drt1zwkdd3k023llluhkd0963kdtfjh0xfgh8ngfwt2qj9da0l79qgpq3uglrn", shardCoordinator))
	require.Equal(t, core.MetachainShardId, computeShardID(t, "drt1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzllls6prdez", shardCoordinator))
}

func computeShardID(t *testing.T, addressBech32 string, shardCoordinator sharding.Coordinator) uint32 {
	addressBytes, err := addressEncoder.Decode(addressBech32)
	require.NoError(t, err)

	return shardCoordinator.ComputeId(addressBytes)
}

func TestSystemSCsAddressesAndSpecialAddresses(t *testing.T) {
	contractDeployScAdress, err := addressEncoder.Encode(make([]byte, addressEncoder.Len()))
	require.NoError(t, err)
	stakingScAddress, err := addressEncoder.Encode(vm.StakingSCAddress)
	require.NoError(t, err)
	validatorScAddress, err := addressEncoder.Encode(vm.ValidatorSCAddress)
	require.NoError(t, err)
	dcdtScAddress, err := addressEncoder.Encode(vm.DCDTSCAddress)
	require.NoError(t, err)
	governanceScAddress, err := addressEncoder.Encode(vm.GovernanceSCAddress)
	require.NoError(t, err)
	jailingAddress, err := addressEncoder.Encode(vm.JailingAddress)
	require.NoError(t, err)
	endOfEpochAddress, err := addressEncoder.Encode(vm.EndOfEpochAddress)
	require.NoError(t, err)
	delegationManagerScAddress, err := addressEncoder.Encode(vm.DelegationManagerSCAddress)
	require.NoError(t, err)
	firstDelegationScAddress, err := addressEncoder.Encode(vm.FirstDelegationSCAddress)
	require.NoError(t, err)

	genesisMintingAddressBytes, err := hex.DecodeString("f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0")
	require.NoError(t, err)
	genesisMintingAddress, err := addressEncoder.Encode(genesisMintingAddressBytes)
	require.NoError(t, err)
	systemAccountAddress, err := addressEncoder.Encode(core.SystemAccountAddress)
	require.NoError(t, err)

	dcdtGlobalSettingsAddresses := getGlobalSettingsAddresses()

	header := []string{"Smart contract/Special address", "Address"}
	lines := []*display.LineData{
		display.NewLineData(false, []string{"Contract deploy", contractDeployScAdress}),
		display.NewLineData(false, []string{"Staking", stakingScAddress}),
		display.NewLineData(false, []string{"Validator", validatorScAddress}),
		display.NewLineData(false, []string{"DCDT", dcdtScAddress}),
		display.NewLineData(false, []string{"Governance", governanceScAddress}),
		display.NewLineData(false, []string{"Jailing address", jailingAddress}),
		display.NewLineData(false, []string{"End of epoch address", endOfEpochAddress}),
		display.NewLineData(false, []string{"Delegation manager", delegationManagerScAddress}),
		display.NewLineData(false, []string{"First delegation", firstDelegationScAddress}),
		display.NewLineData(false, []string{"Genesis Minting Address", genesisMintingAddress}),
		display.NewLineData(false, []string{"System Account Address", systemAccountAddress}),
		display.NewLineData(false, []string{"DCDT Global Settings Shard 0", dcdtGlobalSettingsAddresses[0]}),
		display.NewLineData(false, []string{"DCDT Global Settings Shard 1", dcdtGlobalSettingsAddresses[1]}),
		display.NewLineData(false, []string{"DCDT Global Settings Shard 2", dcdtGlobalSettingsAddresses[2]}),
	}

	table, _ := display.CreateTableString(header, lines)
	fmt.Println(table)

	assert.Equal(t, "drt1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqlllsjr0pnq", stakingScAddress)
	assert.Equal(t, "drt1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqplllskzf8kp", validatorScAddress)
	assert.Equal(t, "drt1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzllls6prdez", dcdtScAddress)
	assert.Equal(t, "drt1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqrllls7q9tur", governanceScAddress)
	assert.Equal(t, "drt1qqqqqqqqqqqqqqqpqqqqqqqqqrlllllllllllllllllllllllllswxc2vg", jailingAddress)
	assert.Equal(t, "drt1qqqqqqqqqqqqqqqpqqqqqqqqllllllllllllllllllllllllllls7el70e", endOfEpochAddress)
	assert.Equal(t, "drt1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqylllsz8he8y", delegationManagerScAddress)
	assert.Equal(t, "drt1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq0llllsa22qry", firstDelegationScAddress)
	assert.Equal(t, "drt1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq85hk5z", contractDeployScAdress)
	assert.Equal(t, "drt17rc0pu8s7rc0pu8s7rc0pu8s7rc0pu8s7rc0pu8s7rc0pu8s7rcqa2qg80", genesisMintingAddress)
	assert.Equal(t, "drt1llllllllllllllllllllllllllllllllllllllllllllllllllls9258a4", systemAccountAddress)
	assert.Equal(t, "drt1llllllllllllllllllllllllllllllllllllllllllllllllluqqh8x2vh", dcdtGlobalSettingsAddresses[0])
	assert.Equal(t, "drt1llllllllllllllllllllllllllllllllllllllllllllllllluqszxwjey", dcdtGlobalSettingsAddresses[1])
	assert.Equal(t, "drt1lllllllllllllllllllllllllllllllllllllllllllllllllupqc6fs83", dcdtGlobalSettingsAddresses[2])
}

func getGlobalSettingsAddresses() map[uint32]string {
	numShards := uint32(3)
	addressesMap := make(map[uint32]string, numShards)
	for i := uint32(0); i < numShards; i++ {
		addressesMap[i] = computeGlobalSettingsAddr(i)
	}

	return addressesMap
}

func computeGlobalSettingsAddr(shardID uint32) string {
	baseSystemAccountAddress := core.SystemAccountAddress
	globalSettingsAddress := baseSystemAccountAddress
	globalSettingsAddress[len(globalSettingsAddress)-1] = uint8(shardID)

	computedAddress, _ := addressEncoder.Encode(globalSettingsAddress)

	return computedAddress
}
