package broadcast

import (
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
)

// SetMarshalizerMeta sets the unexported marshaller
func (mcm *metaChainMessenger) SetMarshalizerMeta(
	m marshal.Marshalizer,
) {
	mcm.marshalizer = m
}
