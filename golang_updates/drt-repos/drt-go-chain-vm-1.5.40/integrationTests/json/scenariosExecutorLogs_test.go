package vmjsonintegrationtest

import (
	"testing"

	"github.com/TerraDharitri/drt-go-chain-scenario/worldmock"
	"github.com/TerraDharitri/drt-go-chain-vm/testcommon/testexecutor"
	"github.com/TerraDharitri/drt-go-chain-vm/wasmer"
	"github.com/TerraDharitri/drt-go-chain-vm/wasmer2"
)

func TestRustCompareAdderLog(t *testing.T) {
	if !testexecutor.IsWasmer1Allowed() {
		t.Skip("run exclusively with wasmer1")
	}

	expected := ScenariosTest(t).
		Folder("adder/scenarios").
		WithExecutorFactory(wasmer.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("adder/scenarios").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestRustFactorialLog(t *testing.T) {
	if !testexecutor.IsWasmer1Allowed() {
		t.Skip("run exclusively with wasmer1")
	}

	expected := ScenariosTest(t).
		Folder("factorial/scenarios").
		WithExecutorFactory(wasmer.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("factorial/scenarios").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestRustErc20Log(t *testing.T) {
	if !testexecutor.IsWasmer1Allowed() {
		t.Skip("run exclusively with wasmer1")
	}

	expected := ScenariosTest(t).
		Folder("erc20-rust/scenarios").
		WithExecutorFactory(wasmer.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("erc20-rust/scenarios").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestCErc20Log(t *testing.T) {
	if !testexecutor.IsWasmer1Allowed() {
		t.Skip("run exclusively with wasmer1")
	}

	expected := ScenariosTest(t).
		Folder("erc20-c").
		WithExecutorFactory(wasmer.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("erc20-c").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestDigitalCashLog(t *testing.T) {
	if !testexecutor.IsWasmer1Allowed() {
		t.Skip("run exclusively with wasmer1")
	}

	expected := ScenariosTest(t).
		Folder("digital-cash").
		WithExecutorFactory(wasmer.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("digital-cash").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestDCDTMultiTransferOnCallbackLog(t *testing.T) {
	if !testexecutor.IsWasmer1Allowed() {
		t.Skip("run exclusively with wasmer1")
	}

	expected := ScenariosTest(t).
		Folder("features/composability/scenarios").
		File("forw_raw_call_async_retrieve_multi_transfer.scen.json").
		WithExecutorFactory(wasmer.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("features/composability/scenarios").
		File("forw_raw_call_async_retrieve_multi_transfer.scen.json").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestCreateAsyncCallLog(t *testing.T) {
	if !testexecutor.IsWasmer1Allowed() {
		t.Skip("run exclusively with wasmer1")
	}

	expected := ScenariosTest(t).
		Folder("features/composability/scenarios").
		File("promises_single_transfer.scen.json").
		WithExecutorFactory(wasmer.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("features/composability/scenarios").
		File("promises_single_transfer.scen.json").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestDCDTMultiTransferOnCallAndCallbackLog(t *testing.T) {
	if !testexecutor.IsWasmer1Allowed() {
		t.Skip("run exclusively with wasmer1")
	}

	expected := ScenariosTest(t).
		Folder("features/composability/scenarios").
		File("forw_raw_async_send_and_retrieve_multi_transfer_funds.scen.json").
		WithExecutorFactory(wasmer.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("features/composability/scenarios").
		File("forw_raw_async_send_and_retrieve_multi_transfer_funds.scen.json").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestMultisigLog(t *testing.T) {
	if !testexecutor.IsWasmer1Allowed() {
		t.Skip("run exclusively with wasmer1")
	}

	expected := ScenariosTest(t).
		Folder("multisig/scenarios").
		Exclude("multisig/scenarios/interactor*").
		WithExecutorFactory(wasmer.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("multisig/scenarios").
		Exclude("multisig/scenarios/interactor*").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestDnsContractLog(t *testing.T) {
	if !testexecutor.IsWasmer1Allowed() {
		t.Skip("run exclusively with wasmer1")
	}

	if testing.Short() {
		t.Skip("not a short test")
	}

	expected := ScenariosTest(t).
		Folder("dns").
		WithEnableEpochsHandler(worldmock.EnableEpochsHandlerStubNoFlags()).
		WithExecutorFactory(wasmer.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("dns").
		WithEnableEpochsHandler(worldmock.EnableEpochsHandlerStubNoFlags()).
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestCrowdfundingDcdtLog(t *testing.T) {
	if !testexecutor.IsWasmer1Allowed() {
		t.Skip("run exclusively with wasmer1")
	}

	expected := ScenariosTest(t).
		Folder("crowdfunding-dcdt").
		WithExecutorFactory(wasmer.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("crowdfunding-dcdt").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestRewaDcdtSwapLog(t *testing.T) {
	if !testexecutor.IsWasmer1Allowed() {
		t.Skip("run exclusively with wasmer1")
	}

	expected := ScenariosTest(t).
		Folder("rewa-dcdt-swap").
		WithExecutorFactory(wasmer.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("rewa-dcdt-swap").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestPingPongRewaLog(t *testing.T) {
	if !testexecutor.IsWasmer1Allowed() {
		t.Skip("run exclusively with wasmer1")
	}

	expected := ScenariosTest(t).
		Folder("ping-pong-rewa").
		WithExecutorFactory(wasmer.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("ping-pong-rewa").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestRustAttestationLog(t *testing.T) {
	if !testexecutor.IsWasmer1Allowed() {
		t.Skip("run exclusively with wasmer1")
	}

	if testing.Short() {
		t.Skip("not a short test")
	}

	expected := ScenariosTest(t).
		Folder("attestation-rust").
		WithExecutorFactory(wasmer.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("attestation-rust").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestCAttestationLog(t *testing.T) {
	if !testexecutor.IsWasmer1Allowed() {
		t.Skip("run exclusively with wasmer1")
	}

	if testing.Short() {
		t.Skip("not a short test")
	}

	expected := ScenariosTest(t).
		Folder("attestation-c").
		WithExecutorFactory(wasmer.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("attestation-c").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expected)
}
