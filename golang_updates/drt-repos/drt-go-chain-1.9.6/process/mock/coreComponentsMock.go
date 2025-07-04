package mock

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data/endProcess"
	"github.com/TerraDharitri/drt-go-chain-core/data/typeConverters"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/consensus"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain/storage"
	"github.com/TerraDharitri/drt-go-chain/testscommon/economicsmocks"
)

// CoreComponentsMock -
type CoreComponentsMock struct {
	IntMarsh                           marshal.Marshalizer
	TxMarsh                            marshal.Marshalizer
	Hash                               hashing.Hasher
	TxSignHasherField                  hashing.Hasher
	UInt64ByteSliceConv                typeConverters.Uint64ByteSliceConverter
	AddrPubKeyConv                     core.PubkeyConverter
	ValPubKeyConv                      core.PubkeyConverter
	PathHdl                            storage.PathManagerHandler
	ChainIdCalled                      func() string
	MinTransactionVersionCalled        func() uint32
	GenesisNodesSetupCalled            func() sharding.GenesisNodesSetupHandler
	TxVersionCheckField                process.TxVersionCheckerHandler
	EpochNotifierField                 process.EpochNotifier
	EnableEpochsHandlerField           common.EnableEpochsHandler
	RoundNotifierField                 process.RoundNotifier
	EnableRoundsHandlerField           process.EnableRoundsHandler
	RoundField                         consensus.RoundHandler
	StatusField                        core.AppStatusHandler
	ChanStopNode                       chan endProcess.ArgEndProcess
	NodeTypeProviderField              core.NodeTypeProviderHandler
	EconomicsDataField                 process.EconomicsDataHandler
	ProcessStatusHandlerField          common.ProcessStatusHandler
	ChainParametersHandlerField        process.ChainParametersHandler
	HardforkTriggerPubKeyField         []byte
	ChainParametersSubscriberField     process.ChainParametersSubscriber
	FieldsSizeCheckerField             common.FieldsSizeChecker
	EpochChangeGracePeriodHandlerField common.EpochChangeGracePeriodHandler
}

// ChanStopNodeProcess -
func (ccm *CoreComponentsMock) ChanStopNodeProcess() chan endProcess.ArgEndProcess {
	if ccm.ChanStopNode != nil {
		return ccm.ChanStopNode
	}

	return endProcess.GetDummyEndProcessChannel()
}

// InternalMarshalizer -
func (ccm *CoreComponentsMock) InternalMarshalizer() marshal.Marshalizer {
	return ccm.IntMarsh
}

// SetInternalMarshalizer -
func (ccm *CoreComponentsMock) SetInternalMarshalizer(m marshal.Marshalizer) error {
	ccm.IntMarsh = m
	return nil
}

// TxMarshalizer -
func (ccm *CoreComponentsMock) TxMarshalizer() marshal.Marshalizer {
	return ccm.TxMarsh
}

// Hasher -
func (ccm *CoreComponentsMock) Hasher() hashing.Hasher {
	return ccm.Hash
}

// TxSignHasher -
func (ccm *CoreComponentsMock) TxSignHasher() hashing.Hasher {
	return ccm.TxSignHasherField
}

// Uint64ByteSliceConverter -
func (ccm *CoreComponentsMock) Uint64ByteSliceConverter() typeConverters.Uint64ByteSliceConverter {
	return ccm.UInt64ByteSliceConv
}

// ChainParametersHandler -
func (ccm *CoreComponentsMock) ChainParametersHandler() process.ChainParametersHandler {
	return ccm.ChainParametersHandlerField
}

// AddressPubKeyConverter -
func (ccm *CoreComponentsMock) AddressPubKeyConverter() core.PubkeyConverter {
	return ccm.AddrPubKeyConv
}

// ValidatorPubKeyConverter -
func (ccm *CoreComponentsMock) ValidatorPubKeyConverter() core.PubkeyConverter {
	return ccm.ValPubKeyConv
}

// PathHandler -
func (ccm *CoreComponentsMock) PathHandler() storage.PathManagerHandler {
	return ccm.PathHdl
}

// ChainID -
func (ccm *CoreComponentsMock) ChainID() string {
	if ccm.ChainIdCalled != nil {
		return ccm.ChainIdCalled()
	}
	return "undefined"
}

// MinTransactionVersion -
func (ccm *CoreComponentsMock) MinTransactionVersion() uint32 {
	if ccm.MinTransactionVersionCalled != nil {
		return ccm.MinTransactionVersionCalled()
	}
	return 1
}

// TxVersionChecker -
func (ccm *CoreComponentsMock) TxVersionChecker() process.TxVersionCheckerHandler {
	return ccm.TxVersionCheckField
}

// GenesisNodesSetup -
func (ccm *CoreComponentsMock) GenesisNodesSetup() sharding.GenesisNodesSetupHandler {
	if ccm.GenesisNodesSetupCalled != nil {
		return ccm.GenesisNodesSetupCalled()
	}
	return nil
}

// EpochNotifier -
func (ccm *CoreComponentsMock) EpochNotifier() process.EpochNotifier {
	return ccm.EpochNotifierField
}

// EnableEpochsHandler -
func (ccm *CoreComponentsMock) EnableEpochsHandler() common.EnableEpochsHandler {
	return ccm.EnableEpochsHandlerField
}

// RoundNotifier -
func (ccm *CoreComponentsMock) RoundNotifier() process.RoundNotifier {
	return ccm.RoundNotifierField
}

// EnableEpochsHandler -
func (ccm *CoreComponentsMock) EnableRoundsHandler() process.EnableRoundsHandler {
	return ccm.EnableRoundsHandlerField
}

// RoundHandler -
func (ccm *CoreComponentsMock) RoundHandler() consensus.RoundHandler {
	return ccm.RoundField
}

// NodeTypeProvider -
func (ccm *CoreComponentsMock) NodeTypeProvider() core.NodeTypeProviderHandler {
	return ccm.NodeTypeProviderField
}

// EconomicsData -
func (ccm *CoreComponentsMock) EconomicsData() process.EconomicsDataHandler {
	if !check.IfNil(ccm.EconomicsDataField) {
		return ccm.EconomicsDataField
	}

	return &economicsmocks.EconomicsHandlerMock{}
}

// ProcessStatusHandler -
func (ccm *CoreComponentsMock) ProcessStatusHandler() common.ProcessStatusHandler {
	return ccm.ProcessStatusHandlerField
}

// HardforkTriggerPubKey -
func (ccm *CoreComponentsMock) HardforkTriggerPubKey() []byte {
	return ccm.HardforkTriggerPubKeyField
}

// ChainParametersSubscriber -
func (ccm *CoreComponentsMock) ChainParametersSubscriber() process.ChainParametersSubscriber {
	return ccm.ChainParametersSubscriberField
}

// FieldsSizeChecker -
func (ccm *CoreComponentsMock) FieldsSizeChecker() common.FieldsSizeChecker {
	return ccm.FieldsSizeCheckerField
}

// EpochChangeGracePeriodHandler -
func (ccm *CoreComponentsMock) EpochChangeGracePeriodHandler() common.EpochChangeGracePeriodHandler {
	return ccm.EpochChangeGracePeriodHandlerField
}

// IsInterfaceNil -
func (ccm *CoreComponentsMock) IsInterfaceNil() bool {
	return ccm == nil
}
