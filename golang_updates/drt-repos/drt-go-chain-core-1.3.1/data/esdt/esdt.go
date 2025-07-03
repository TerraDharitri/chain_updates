//go:generate protoc -I=proto -I=$GOPATH/src -I=$GOPATH/src/github.com/TerraDharitri/protobuf/protobuf  --gogoslick_out=. dcdt.proto
package dcdt

import "math/big"

// New returns a new batch from given buffers
func New() *DCDigitalToken {
	return &DCDigitalToken{
		Value: big.NewInt(0),
	}
}
