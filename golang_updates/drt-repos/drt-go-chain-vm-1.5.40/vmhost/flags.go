package vmhost

import "github.com/TerraDharitri/drt-go-chain-core/core"

const (
	// CryptoOpcodesV2Flag defines the flag that activates the new crypto APIs for RC1.7
	CryptoOpcodesV2Flag core.EnableEpochFlag = "CryptoOpcodesV2Flag"

	// MultiDCDTNFTTransferAndExecuteByUserFlag defines the flag that activates the enshrined sovereign functions
	MultiDCDTNFTTransferAndExecuteByUserFlag core.EnableEpochFlag = "MultiDCDTNFTTransferAndExecuteByUserFlag"

	// UseGasBoundedShouldFailExecutionFlag defines the flag that activates failing of execution if gas bounded check fails
	UseGasBoundedShouldFailExecutionFlag core.EnableEpochFlag = "UseGasBoundedShouldFailExecutionFlag"

	// CheckBuiltInCallOnTransferValueAndFailExecutionFlag defines the flag that activates failing of execution if gas bounded check fails
	CheckBuiltInCallOnTransferValueAndFailExecutionFlag core.EnableEpochFlag = "CheckBuiltInCallOnTransferValueAndFailExecutionFlag"
)
