package dcdtMultiTransferToVaultSameShard

import (
	"testing"

	multitransfer "github.com/TerraDharitri/drt-go-chain/integrationTests/vm/dcdt/multi-transfer"
)

func TestDCDTMultiTransferToVaultSameShard(t *testing.T) {
	multitransfer.DcdtMultiTransferToVault(t, false, "../../testdata/vaultV2.wasm")
}
