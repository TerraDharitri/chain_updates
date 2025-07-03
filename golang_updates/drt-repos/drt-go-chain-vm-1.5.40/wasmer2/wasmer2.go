// Package wasmer2 is a Go library to run WebAssembly binaries.
package wasmer2

import logger "github.com/TerraDharitri/drt-go-chain-logger"

// VM logger.
var logWasmer2 = logger.GetOrCreate("vm/executor")

// Executor logger.
type LogLevel uint64

const (
	LogLevelOff LogLevel = iota
	LogLevelError
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
	LogLevelTrace
)
