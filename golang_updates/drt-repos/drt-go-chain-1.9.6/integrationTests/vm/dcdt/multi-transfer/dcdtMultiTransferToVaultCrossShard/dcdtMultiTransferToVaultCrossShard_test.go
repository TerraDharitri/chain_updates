package dcdtMultiTransferToVaultCrossShard

import (
	"testing"

	multitransfer "github.com/TerraDharitri/drt-go-chain/integrationTests/vm/dcdt/multi-transfer"
)

func TestDCDTMultiTransferToVaultCrossShard(t *testing.T) {
	multitransfer.DcdtMultiTransferToVault(t, true, "../../testdata/vaultV2.wasm")
}
