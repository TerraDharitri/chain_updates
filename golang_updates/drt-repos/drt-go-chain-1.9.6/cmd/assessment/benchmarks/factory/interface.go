package factory

import "github.com/TerraDharitri/drt-go-chain/cmd/assessment/benchmarks"

type benchmarkCoordinator interface {
	RunAllTests() *benchmarks.TestResults
	IsInterfaceNil() bool
}
