package enablers

import (
	"github.com/TerraDharitri/drt-go-chain-core/core/atomic"
)

type roundFlag struct {
	*atomic.Flag
	round   uint64
	options []string
}
