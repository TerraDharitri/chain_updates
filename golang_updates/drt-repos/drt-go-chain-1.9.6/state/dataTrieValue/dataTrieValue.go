//go:generate protoc -I=. -I=$GOPATH/src -I=$GOPATH/src/github.com/TerraDharitri/protobuf/protobuf --gogoslick_out=. dataTrieValue.proto
package dataTrieValue
