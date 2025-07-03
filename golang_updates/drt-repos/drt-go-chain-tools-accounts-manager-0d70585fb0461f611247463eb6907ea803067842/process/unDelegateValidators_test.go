package process

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core/pubkeyConverter"
	"github.com/TerraDharitri/drt-go-chain-core/data/vm"
	"github.com/TerraDharitri/drt-go-chain-tools-accounts-manager/config"
	"github.com/TerraDharitri/drt-go-chain-tools-accounts-manager/data"
	"github.com/TerraDharitri/drt-go-chain-tools-accounts-manager/mocks"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	"github.com/stretchr/testify/require"
)

func TestAccountsGetter_ValidatorsAccountsPutUnDelegatedValues(t *testing.T) {
	t.Parallel()

	pubKeyConverter, _ := pubkeyConverter.NewBech32PubkeyConverter(32, log)

	ag, err := NewAccountsGetter(&mocks.RestClientStub{
		CallPostRestEndPointCalled: func(path string, dataD interface{}, response interface{}, authenticationData data.RestApiAuthenticationData) error {
			responseVmValue := response.(*data.ResponseVmValue)
			responseVmValue.Data = data.VmValuesResponseData{
				Data: &vm.VMOutputApi{
					ReturnData: [][]byte{big.NewInt(1000000000000000000).Bytes(), []byte("")},
					ReturnCode: vmcommon.Ok.String(),
				},
			}

			return nil
		},
	}, pubKeyConverter, data.RestApiAuthenticationData{}, config.GeneralConfig{
		ValidatorsContract: "drt1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqplllskzf8kp",
	}, &mocks.ElasticClientStub{})
	require.Nil(t, err)

	accountsWithStakeJson := readJson("./testdata/accounts-with-stake.json")
	accountsWithStake := make(map[string]*data.AccountInfoWithStakeValues)
	err = json.Unmarshal([]byte(accountsWithStakeJson), &accountsWithStake)
	require.Nil(t, err)

	err = ag.putUndelegatedValuesFromValidatorsContract(accountsWithStake)
	require.Nil(t, err)

	for _, account := range accountsWithStake {
		require.Equal(t, account.UnDelegateValidator, "1000000000000000000")
		require.Equal(t, account.UnDelegateValidatorNum, float64(1))
	}
}
