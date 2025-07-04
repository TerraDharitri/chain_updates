package spos

import (
	"context"
	"fmt"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/consensus"
	"github.com/TerraDharitri/drt-go-chain/process"
)

// RedundancySingleKeySteppedIn exposes the redundancySingleKeySteppedIn constant
const RedundancySingleKeySteppedIn = redundancySingleKeySteppedIn

// LeaderSingleKeyStartMsg -
const LeaderSingleKeyStartMsg = singleKeyStartMsg

// LeaderMultiKeyStartMsg -
const LeaderMultiKeyStartMsg = multiKeyStartMsg

type RoundConsensus struct {
	*roundConsensus
}

// NewRoundConsensusWrapper -
func NewRoundConsensusWrapper(rcns *roundConsensus) *RoundConsensus {
	return &RoundConsensus{
		roundConsensus: rcns,
	}
}

// BlockProcessor -
func (wrk *Worker) BlockProcessor() process.BlockProcessor {
	return wrk.blockProcessor
}

// SetBlockProcessor -
func (wrk *Worker) SetBlockProcessor(blockProcessor process.BlockProcessor) {
	wrk.blockProcessor = blockProcessor
}

// Bootstrapper -
func (wrk *Worker) Bootstrapper() process.Bootstrapper {
	return wrk.bootstrapper
}

// SetBootstrapper -
func (wrk *Worker) SetBootstrapper(bootstrapper process.Bootstrapper) {
	wrk.bootstrapper = bootstrapper
}

// BroadcastMessenger -
func (wrk *Worker) BroadcastMessenger() consensus.BroadcastMessenger {
	return wrk.broadcastMessenger
}

// SetBroadcastMessenger -
func (wrk *Worker) SetBroadcastMessenger(broadcastMessenger consensus.BroadcastMessenger) {
	wrk.broadcastMessenger = broadcastMessenger
}

// ConsensusState -
func (wrk *Worker) ConsensusState() *ConsensusState {
	return wrk.consensusState
}

// SetConsensusState -
func (wrk *Worker) SetConsensusState(consensusState *ConsensusState) {
	wrk.consensusState = consensusState
}

// ForkDetector -
func (wrk *Worker) ForkDetector() process.ForkDetector {
	return wrk.forkDetector
}

// SetForkDetector -
func (wrk *Worker) SetForkDetector(forkDetector process.ForkDetector) {
	wrk.forkDetector = forkDetector
}

// Marshalizer -
func (wrk *Worker) Marshalizer() marshal.Marshalizer {
	return wrk.marshalizer
}

// SetMarshalizer -
func (wrk *Worker) SetMarshalizer(marshalizer marshal.Marshalizer) {
	wrk.marshalizer = marshalizer
}

// RoundHandler -
func (wrk *Worker) RoundHandler() consensus.RoundHandler {
	return wrk.roundHandler
}

// SetNodeRedundancyHandler -
func (wrk *Worker) SetNodeRedundancyHandler(nodeRedundancyHandler consensus.NodeRedundancyHandler) {
	wrk.nodeRedundancyHandler = nodeRedundancyHandler
}

// SetRoundHandler -
func (wrk *Worker) SetRoundHandler(roundHandler consensus.RoundHandler) {
	wrk.roundHandler = roundHandler
}

// CheckSignature -
func (wrk *Worker) CheckSignature(cnsData *consensus.Message) error {
	return wrk.peerSignatureHandler.VerifyPeerSignature(cnsData.PubKey, core.PeerID(cnsData.OriginatorPid), cnsData.Signature)
}

// ExecuteMessage -
func (wrk *Worker) ExecuteMessage(cnsDtaList []*consensus.Message) {
	wrk.executeMessage(cnsDtaList)
}

// InitReceivedMessages -
func (wrk *Worker) InitReceivedMessages() {
	wrk.initReceivedMessages()
}

// ReceivedSyncState -
func (wrk *Worker) ReceivedSyncState(isNodeSynchronized bool) {
	wrk.receivedSyncState(isNodeSynchronized)
}

// ReceivedMessages -
func (wrk *Worker) ReceivedMessages() map[consensus.MessageType][]*consensus.Message {
	wrk.mutReceivedMessages.RLock()
	defer wrk.mutReceivedMessages.RUnlock()

	return wrk.receivedMessages
}

// SetReceivedMessages -
func (wrk *Worker) SetReceivedMessages(messageType consensus.MessageType, cnsDta []*consensus.Message) {
	wrk.mutReceivedMessages.Lock()
	wrk.receivedMessages[messageType] = cnsDta
	wrk.mutReceivedMessages.Unlock()
}

// NilReceivedMessages -
func (wrk *Worker) NilReceivedMessages() {
	wrk.mutReceivedMessages.Lock()
	wrk.receivedMessages = nil
	wrk.mutReceivedMessages.Unlock()
}

// ReceivedMessagesCalls -
func (wrk *Worker) ReceivedMessagesCalls() map[consensus.MessageType][]func(context.Context, *consensus.Message) bool {
	wrk.mutReceivedMessagesCalls.RLock()
	defer wrk.mutReceivedMessagesCalls.RUnlock()

	return wrk.receivedMessagesCalls
}

// AppendReceivedMessagesCalls -
func (wrk *Worker) AppendReceivedMessagesCalls(messageType consensus.MessageType, f func(context.Context, *consensus.Message) bool) {
	wrk.mutReceivedMessagesCalls.Lock()
	wrk.receivedMessagesCalls[messageType] = append(wrk.receivedMessagesCalls[messageType], f)
	wrk.mutReceivedMessagesCalls.Unlock()
}

// ExecuteMessageChannel -
func (wrk *Worker) ExecuteMessageChannel() chan *consensus.Message {
	return wrk.executeMessageChannel
}

// ConvertHeaderToConsensusMessage -
func (wrk *Worker) ConvertHeaderToConsensusMessage(header data.HeaderHandler) (*consensus.Message, error) {
	return wrk.convertHeaderToConsensusMessage(header)
}

// Hasher -
func (wrk *Worker) Hasher() data.Hasher {
	return wrk.hasher
}

// SetEnableEpochsHandler
func (wrk *Worker) SetEnableEpochsHandler(enableEpochsHandler common.EnableEpochsHandler) {
	wrk.enableEpochsHandler = enableEpochsHandler
}

// AddFutureHeaderToProcessIfNeeded -
func (wrk *Worker) AddFutureHeaderToProcessIfNeeded(header data.HeaderHandler) {
	wrk.addFutureHeaderToProcessIfNeeded(header)
}

// ConsensusStateChangedChannel -
func (wrk *Worker) ConsensusStateChangedChannel() chan bool {
	return wrk.consensusStateChangedChannel
}

// SetConsensusStateChangedChannel -
func (wrk *Worker) SetConsensusStateChangedChannel(consensusStateChangedChannel chan bool) {
	wrk.consensusStateChangedChannel = consensusStateChangedChannel
}

// CheckSelfState -
func (wrk *Worker) CheckSelfState(cnsDta *consensus.Message) error {
	return wrk.checkSelfState(cnsDta)
}

// SetRedundancyHandler -
func (wrk *Worker) SetRedundancyHandler(redundancyHandler consensus.NodeRedundancyHandler) {
	wrk.nodeRedundancyHandler = redundancyHandler
}

// SetKeysHandler -
func (wrk *Worker) SetKeysHandler(keysHandler consensus.KeysHandler) {
	wrk.consensusState.keysHandler = keysHandler
}

// EligibleList -
func (rcns *RoundConsensus) EligibleList() map[string]struct{} {
	return rcns.eligibleNodes
}

// AppStatusHandler -
func (wrk *Worker) AppStatusHandler() core.AppStatusHandler {
	return wrk.appStatusHandler
}

// CheckConsensusMessageValidity -
func (cmv *consensusMessageValidator) CheckConsensusMessageValidity(cnsMsg *consensus.Message, originator core.PeerID) error {
	return cmv.checkConsensusMessageValidity(cnsMsg, originator)
}

// CheckMessageWithFinalInfoValidity -
func (cmv *consensusMessageValidator) CheckMessageWithFinalInfoValidity(cnsMsg *consensus.Message) error {
	return cmv.checkMessageWithFinalInfoValidity(cnsMsg)
}

// CheckMessageWithSignatureValidity -
func (cmv *consensusMessageValidator) CheckMessageWithSignatureValidity(cnsMsg *consensus.Message) error {
	return cmv.checkMessageWithSignatureValidity(cnsMsg)
}

// CheckMessageWithBlockHeaderValidity -
func (cmv *consensusMessageValidator) CheckMessageWithBlockHeaderValidity(cnsMsg *consensus.Message) error {
	return cmv.checkMessageWithBlockHeaderValidity(cnsMsg)
}

// CheckMessageWithBlockBodyValidity -
func (cmv *consensusMessageValidator) CheckMessageWithBlockBodyValidity(cnsMsg *consensus.Message) error {
	return cmv.checkMessageWithBlockBodyValidity(cnsMsg)
}

// CheckMessageWithBlockBodyAndHeaderValidity -
func (cmv *consensusMessageValidator) CheckMessageWithBlockBodyAndHeaderValidity(cnsMsg *consensus.Message) error {
	return cmv.checkMessageWithBlockBodyAndHeaderValidity(cnsMsg)
}

// CheckConsensusMessageValidityForMessageType -
func (cmv *consensusMessageValidator) CheckConsensusMessageValidityForMessageType(cnsMsg *consensus.Message) error {
	return cmv.checkConsensusMessageValidityForMessageType(cnsMsg)
}

// IsBlockHeaderHashSizeValid -
func (cmv *consensusMessageValidator) IsBlockHeaderHashSizeValid(cnsMsg *consensus.Message) bool {
	return cmv.isBlockHeaderHashSizeValid(cnsMsg)
}

// AddMessageTypeToPublicKey -
func (cmv *consensusMessageValidator) AddMessageTypeToPublicKey(pk []byte, round int64, msgType consensus.MessageType) {
	cmv.addMessageTypeToPublicKey(pk, round, msgType)
}

// IsMessageTypeLimitReached -
func (cmv *consensusMessageValidator) IsMessageTypeLimitReached(pk []byte, round int64, msgType consensus.MessageType) bool {
	return cmv.isMessageTypeLimitReached(pk, round, msgType)
}

// GetNumOfMessageTypeForPublicKey -
func (cmv *consensusMessageValidator) GetNumOfMessageTypeForPublicKey(pk []byte, round int64, msgType consensus.MessageType) uint32 {
	cmv.mutPkConsensusMessages.RLock()
	defer cmv.mutPkConsensusMessages.RUnlock()

	key := fmt.Sprintf("%s_%d", string(pk), round)

	mapMsgType, ok := cmv.mapPkConsensusMessages[key]
	if !ok {
		return uint32(0)
	}

	return mapMsgType[msgType]
}

// ResetConsensusMessages -
func (cmv *consensusMessageValidator) ResetConsensusMessages() {
	cmv.resetConsensusMessages()
}

// SetStatus -
func (sp *scheduledProcessorWrapper) SetStatus(status processingStatus) {
	sp.setStatus(status)
}

// GetStatus -
func (sp *scheduledProcessorWrapper) GetStatus() processingStatus {
	return sp.getStatus()
}

// SetStartTime -
func (sp *scheduledProcessorWrapper) SetStartTime(t time.Time) {
	sp.startTime = t
}

// GetStartTime -
func (sp *scheduledProcessorWrapper) GetStartTime() time.Time {
	return sp.startTime
}

// GetRoundTimeHandler -
func (sp *scheduledProcessorWrapper) GetRoundTimeHandler() process.RoundTimeDurationHandler {
	return sp.roundTimeDurationHandler
}

// ProcessingNotStarted -
var ProcessingNotStarted = processingNotStarted

// ProcessingError -
var ProcessingError = processingError

// InProgress -
var InProgress = inProgress

// ProcessingOK -
var ProcessingOK = processingOK

// Stopped -
var Stopped = stopped

// ProcessingNotStartedString -
var ProcessingNotStartedString = processingNotStartedString

// ProcessingErrorString -
var ProcessingErrorString = processingErrorString

// InProgressString -
var InProgressString = inProgressString

// ProcessingOKString -
var ProcessingOKString = processingOKString

// StoppedString -
var StoppedString = stoppedString

// UnexpectedString -
var UnexpectedString = unexpectedString
