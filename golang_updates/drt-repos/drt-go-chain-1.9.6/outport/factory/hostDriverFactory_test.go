package factory

import (
	"fmt"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-communication/websocket/data"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/testscommon/marshallerMock"
	"github.com/stretchr/testify/require"
)

func TestCreateHostDriver(t *testing.T) {
	t.Parallel()

	args := ArgsHostDriverFactory{
		HostConfig: config.HostDriversConfig{
			URL:                "ws://localhost",
			RetryDurationInSec: 1,
			MarshallerType:     "json",
			Mode:               data.ModeClient,
		},
		Marshaller: &marshallerMock.MarshalizerStub{},
	}

	driver, err := CreateHostDriver(args)
	require.Nil(t, err)
	require.NotNil(t, driver)
	require.Equal(t, "*host.hostDriver", fmt.Sprintf("%T", driver))
}
