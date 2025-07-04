package mock

import (
	"github.com/TerraDharitri/drt-go-chain-core/data"

	"github.com/TerraDharitri/drt-go-chain/epochStart"
)

// EpochStartNotifierStub -
type EpochStartNotifierStub struct {
	RegisterHandlerCalled                 func(handler epochStart.ActionHandler)
	UnregisterHandlerCalled               func(handler epochStart.ActionHandler)
	NotifyAllCalled                       func(hdr data.HeaderHandler)
	NotifyAllPrepareCalled                func(hdr data.HeaderHandler, body data.BodyHandler)
	epochStartHdls                        []epochStart.ActionHandler
	NotifyEpochChangeConfirmedCalled      func(epoch uint32)
	RegisterForEpochChangeConfirmedCalled func(handler func(epoch uint32))
}

// RegisterForEpochChangeConfirmed -
func (esnm *EpochStartNotifierStub) RegisterForEpochChangeConfirmed(handler func(epoch uint32)) {
	if esnm.RegisterForEpochChangeConfirmedCalled != nil {
		esnm.RegisterForEpochChangeConfirmedCalled(handler)
	}
}

// NotifyEpochChangeConfirmed -
func (esnm *EpochStartNotifierStub) NotifyEpochChangeConfirmed(epoch uint32) {
	if esnm.NotifyEpochChangeConfirmedCalled != nil {
		esnm.NotifyEpochChangeConfirmedCalled(epoch)
	}
}

// RegisterHandler -
func (esnm *EpochStartNotifierStub) RegisterHandler(handler epochStart.ActionHandler) {
	if esnm.RegisterHandlerCalled != nil {
		esnm.RegisterHandlerCalled(handler)
	}

	esnm.epochStartHdls = append(esnm.epochStartHdls, handler)
}

// UnregisterHandler -
func (esnm *EpochStartNotifierStub) UnregisterHandler(handler epochStart.ActionHandler) {
	if esnm.UnregisterHandlerCalled != nil {
		esnm.UnregisterHandlerCalled(handler)
	}

	for i, hdl := range esnm.epochStartHdls {
		if hdl == handler {
			esnm.epochStartHdls = append(esnm.epochStartHdls[:i], esnm.epochStartHdls[i+1:]...)
			break
		}
	}
}

// NotifyAllPrepare -
func (esnm *EpochStartNotifierStub) NotifyAllPrepare(metaHdr data.HeaderHandler, body data.BodyHandler) {
	if esnm.NotifyAllPrepareCalled != nil {
		esnm.NotifyAllPrepareCalled(metaHdr, body)
	}

	for _, hdl := range esnm.epochStartHdls {
		hdl.EpochStartPrepare(metaHdr, body)
	}
}

// NotifyAll -
func (esnm *EpochStartNotifierStub) NotifyAll(hdr data.HeaderHandler) {
	if esnm.NotifyAllCalled != nil {
		esnm.NotifyAllCalled(hdr)
	}

	for _, hdl := range esnm.epochStartHdls {
		hdl.EpochStartAction(hdr)
	}
}

// IsInterfaceNil -
func (esnm *EpochStartNotifierStub) IsInterfaceNil() bool {
	return esnm == nil
}
