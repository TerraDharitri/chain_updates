package process

import (
	"encoding/json"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core/pubkeyConverter"
	"github.com/TerraDharitri/drt-go-chain-tools-accounts-manager/config"
	"github.com/TerraDharitri/drt-go-chain-tools-accounts-manager/data"
	"github.com/TerraDharitri/drt-go-chain-tools-accounts-manager/mocks"
	"github.com/stretchr/testify/require"
)

func TestAccountsGetter_DelegationMetaPutUnDelegatedValues(t *testing.T) {
	t.Parallel()
	pubKeyConverter, _ := pubkeyConverter.NewBech32PubkeyConverter(32, log)

	ag, err := NewAccountsGetter(&mocks.RestClientStub{}, pubKeyConverter, data.RestApiAuthenticationData{}, config.GeneralConfig{}, &mocks.ElasticClientStub{
		DoScrollRequestAllDocumentsCalled: func(index string, body []byte, handlerFunc func(responseBytes []byte) error) error {
			delegatorsJson := readJson("./testdata/delegators-es.json")
			return handlerFunc([]byte(delegatorsJson))
		},
	})
	require.Nil(t, err)

	accountsWithStakeJson := readJson("./testdata/account-with-stake.json")
	accountsWithStake := make(map[string]*data.AccountInfoWithStakeValues)
	err = json.Unmarshal([]byte(accountsWithStakeJson), &accountsWithStake)
	require.Nil(t, err)

	err = ag.unDelegatedInfoProc.putUnDelegateInfoFromStakingProviders(accountsWithStake)
	require.Nil(t, err)

	accounts1 := accountsWithStake["drt102hpxzdawtka2usnmkqsk58v3k70jprhy50u4kdgc44j5azd6q5qr0ya25"]
	require.Equal(t, accounts1.UnDelegateDelegation, "2000000000000000000")
	require.Equal(t, accounts1.UnDelegateDelegationNum, float64(2))

	accounts2 := accountsWithStake["drt1063s32hkyj55dpvhtsadacpt268angz2rh2wu4zwqe54awxz5q5ss5r6yu"]
	require.Equal(t, accounts2.UnDelegateDelegation, "10000000000000000000")
	require.Equal(t, accounts2.UnDelegateDelegationNum, float64(10))
}
