package mocks

import "github.com/TerraDharitri/drt-go-chain-tools-accounts-manager/data"

// RestClientStub -
type RestClientStub struct {
	CallPostRestEndPointCalled func(path string, data interface{}, response interface{}, authenticationData data.RestApiAuthenticationData) error
}

// CallGetRestEndPoint -
func (r RestClientStub) CallGetRestEndPoint(_ string, _ interface{}, _ data.RestApiAuthenticationData) error {
	panic("implement me")
}

// CallPostRestEndPoint -
func (r RestClientStub) CallPostRestEndPoint(
	path string,
	data interface{},
	response interface{},
	authenticationData data.RestApiAuthenticationData,
) error {
	if r.CallPostRestEndPointCalled != nil {
		return r.CallPostRestEndPointCalled(path, data, response, authenticationData)
	}

	return nil
}
