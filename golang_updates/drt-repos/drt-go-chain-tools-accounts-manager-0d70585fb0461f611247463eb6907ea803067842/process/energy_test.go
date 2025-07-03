package process

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core/pubkeyConverter"
	"github.com/TerraDharitri/drt-go-chain-tools-accounts-manager/config"
	"github.com/TerraDharitri/drt-go-chain-tools-accounts-manager/core"
	"github.com/TerraDharitri/drt-go-chain-tools-accounts-manager/data"
	"github.com/TerraDharitri/drt-go-chain-tools-accounts-manager/mocks"
	"github.com/stretchr/testify/require"
)

func TestExtractAddressesAndEnergy(t *testing.T) {
	t.Parallel()

	pubKey, _ := pubkeyConverter.NewBech32PubkeyConverter(32, log)
	auth := core.FetchAuthenticationData(config.APIConfig{})
	accountsWithEnergyGetter, err := NewAccountsGetter(&mocks.RestClientStub{}, pubKey, auth, config.GeneralConfig{}, &mocks.ElasticClientStub{})
	require.Nil(t, err)

	testData := readJson("./testdata/account-storage.json")
	res, err := accountsWithEnergyGetter.extractAddressesAndEnergy([]byte(testData), 2047)
	require.Nil(t, err)
	require.NotNil(t, res)

	resNegativeEnergy := res["drt1ytknlprw8lyfn9x5yn0e0c8wtttkzumzm5z7dp4ynadnhq26aczszvqxrq"]
	require.Equal(t, &data.AccountInfoWithStakeValues{
		StakeInfo: data.StakeInfo{
			Energy:    "-2613000000000000000000",
			EnergyNum: -2613,
			EnergyDetails: &data.EnergyDetails{
				LastUpdateEpoch:   1891,
				Amount:            "1599000000000000000000",
				TotalLockedTokens: "27000000000000000000",
			},
		},
	}, resNegativeEnergy)

	require.Equal(t, &data.AccountInfoWithStakeValues{
		StakeInfo: data.StakeInfo{
			Energy:    "336000000000000000000",
			EnergyNum: 336,
			EnergyDetails: &data.EnergyDetails{
				LastUpdateEpoch:   1891,
				Amount:            "5328000000000000000000",
				TotalLockedTokens: "32000000000000000000",
			},
		},
	}, res["drt10f7nnvqk8xvyd50f2sc5p4e0ru4alf99p3v7zfe4uvenra2esgesve2x9q"])

	require.Equal(t, &data.AccountInfoWithStakeValues{
		StakeInfo: data.StakeInfo{
			Energy:    "273000000000000000000",
			EnergyNum: 273,
			EnergyDetails: &data.EnergyDetails{
				LastUpdateEpoch:   1891,
				Amount:            "4173000000000000000000",
				TotalLockedTokens: "25000000000000000000",
			},
		},
	}, res["drt1ejjwyzrdj053vcs5nhupxn6kha8audf4mla6tth9339zmcx52w5qr39765"])

	require.Equal(t, &data.AccountInfoWithStakeValues{
		StakeInfo: data.StakeInfo{
			Energy:    "12625000000000000000000000",
			EnergyNum: 12625000,
			EnergyDetails: &data.EnergyDetails{
				LastUpdateEpoch:   1881,
				Amount:            "96455000000000000000000000",
				TotalLockedTokens: "505000000000000000000000",
			},
		},
	}, res["drt1yhhzgv5ql3h8gppy5286grre23vfgw68tnth7dmcl8ywpd9puluqzym0dm"])

	require.Equal(t, &data.AccountInfoWithStakeValues{
		StakeInfo: data.StakeInfo{
			Energy:    "63371454581200312235",
			EnergyNum: 63.3714545812,
			EnergyDetails: &data.EnergyDetails{
				LastUpdateEpoch:   1881,
				Amount:            "4544871637244977820221",
				TotalLockedTokens: "26996989052191430771",
			},
		},
	}, res["drt188lxgu4m889yht73t3svs4lxknfqtv2vgymgzz283x6wv4hw9nwqjyldvj"])

}

func readJson(path string) string {
	jsonFile, _ := os.Open(path)
	byteValue, _ := ioutil.ReadAll(jsonFile)

	return string(byteValue)
}

func TestExtractEnergyValue(t *testing.T) {
	input := "0000000bff6ccc6b2db50c5215ad8b000000000000035e0000000a0cd2f65ea43133838fec"

	energyDetails, ok := extractEnergyFromValue(input)
	require.True(t, ok)
	require.NotNil(t, energyDetails)

	require.Equal(t, &data.EnergyDetails{
		LastUpdateEpoch:   862,
		Amount:            "-695139380645670211441269",
		TotalLockedTokens: "60559966857227114090476",
	}, energyDetails)
}

func TestExtractBlockInfo(t *testing.T) {
	t.Parallel()

	response := `{"blockInfo":{"hash":"52a2e3c800d03b1499e3cbc57431ee5f122e1bf0e1065fa05578f2d58621f7a0","nonce":3576295,"rootHash":"1829f8c869318f1c5ddc8a887fc2bb206b42fa9a484b8beb94e7873b633cdc61"}}`

	blockInfo, err := extractBlockInfo([]byte(response))
	require.Nil(t, err)
	require.Equal(t, &data.BlockInfo{
		Hash:     "52a2e3c800d03b1499e3cbc57431ee5f122e1bf0e1065fa05578f2d58621f7a0",
		Nonce:    uint64(3576295),
		RootHash: "1829f8c869318f1c5ddc8a887fc2bb206b42fa9a484b8beb94e7873b633cdc61",
	}, blockInfo)
}
