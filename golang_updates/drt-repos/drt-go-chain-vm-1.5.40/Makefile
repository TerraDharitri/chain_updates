.PHONY: test test-short build vmserver clean

VM_VERSION := $(shell git describe --tags --long --dirty --always)

clean:
	go clean -cache -testcache

build:
	go build ./...

vmserver:
ifndef VMSERVER_PATH
	$(error VMSERVER_PATH is undefined)
endif
	go build -o ./cmd/vmserver/vmserver ./cmd/vmserver
	cp ./cmd/vmserver/vmserver ${VMSERVER_PATH}

test:
	go clean -cache -testcache
	VMOAECUTOR="wasmer1" go test ./...
	go clean -cache -testcache
	VMOAECUTOR="wasmer2" go test ./...

test-w1: clean
	VMOAECUTOR="wasmer1" go test ./...

test-w2: clean
	VMOAECUTOR="wasmer2" go test ./...

test-v: clean
	go test ./... -v

test-serial: clean
	go test ./... -failfast -p 1

test-short: clean
	go test ./... -short

test-short-v: clean
	go test ./... -short -v

test-short-serial:
	go test ./... -short -failfast -p 1

print-api-costs:
	@echo "bigIntOps.go:"
	@grep "func v1_5\|GasSchedule" vmhost/vmhooks/bigIntOps.go | sed -e "/func/ s:func v1_5_\(.*\)(.*:\1:" -e "/GasSchedule/ s:metering.GasSchedule()::"
	@echo "----------------"
	@echo "baseOps.go:"
	@grep "func v1_5\|GasSchedule" vmhost/vmhooks/baseOps.go | sed -e "/func/ s:func v1_5_\(.*\)(.*:\1:" -e "/GasSchedule/ s:metering.GasSchedule()::"
	@echo "----------------"
	@echo "managedei.go:"
	@grep "func v1_5\|GasSchedule" vmhost/vmhooks/managedei.go | sed -e "/func/ s:func v1_5_\(.*\)(.*:\1:" -e "/GasSchedule/ s:metering.GasSchedule()::"
	@echo "----------------"
	@echo "manBufOps.go:"
	@grep "func v1_5\|GasSchedule" vmhost/vmhooks/manBufOps.go | sed -e "/func/ s:func v1_5_\(.*\)(.*:\1:" -e "/GasSchedule/ s:metering.GasSchedule()::"
	@echo "----------------"
	@echo "smallIntOps.go:"
	@grep "func v1_5\|GasSchedule" vmhost/vmhooks/smallIntOps.go | sed -e "/func/ s:func v1_5_\(.*\)(.*:\1:" -e "/GasSchedule/ s:metering.GasSchedule()::"


build-test-contracts: build-test-contracts-drtpy build-test-contracts-wat

build-test-contracts-wat:
	cd test/contracts/init-simple-popcnt && wat2wasm *.wat
	cd test/contracts/forbidden-opcodes/data-drop/output && wat2wasm *.wat
	cd test/contracts/forbidden-opcodes/memory-copy/output && wat2wasm *.wat
	cd test/contracts/forbidden-opcodes/memory-fill/output && wat2wasm *.wat
	cd test/contracts/forbidden-opcodes/memory-init/output && wat2wasm *.wat
	cd test/contracts/forbidden-opcodes/simd/output && wat2wasm *.wat
	cd test/contracts/wasmbacking/imported-global/output && wat2wasm *.wat
	cd test/contracts/wasmbacking/mem-exceeded-max-pages/output && wat2wasm *.wat
	cd test/contracts/wasmbacking/mem-exceeded-pages/output && wat2wasm *.wat
	cd test/contracts/wasmbacking/mem-grow/output/ && wat2wasm *.wat
	cd test/contracts/wasmbacking/mem-min-pages-greater-than-max-pages/output && wat2wasm --no-check *.wat
	cd test/contracts/wasmbacking/mem-multiple-max-pages/output/ && wat2wasm *.wat
	cd test/contracts/wasmbacking/mem-multiple-pages/output/ && wat2wasm *.wat
	cd test/contracts/wasmbacking/mem-no-max-pages/output/ && wat2wasm *.wat
	cd test/contracts/wasmbacking/mem-no-pages/output/ && wat2wasm *.wat
	cd test/contracts/wasmbacking/mem-single-page/output/ && wat2wasm *.wat
	cd test/contracts/wasmbacking/memoryless/output/ && wat2wasm *.wat
	cd test/contracts/wasmbacking/multiple-memories/output/ && wat2wasm --no-check *.wat
	cd test/contracts/wasmbacking/multiple-mutable/output/ && wat2wasm *.wat
	cd test/contracts/wasmbacking/noglobals/output/ && wat2wasm *.wat
	cd test/contracts/wasmbacking/single-immutable/output/ && wat2wasm --no-check *.wat
	cd test/contracts/wasmbacking/single-mutable/output/ && wat2wasm *.wat

build-test-contracts-drtpy:
	drtpy contract build --no-optimization ./test/contracts/answer
	drtpy contract build ./test/contracts/async-call-builtin
	drtpy contract build ./test/contracts/async-call-child
	drtpy contract build ./test/contracts/async-call-parent
	drtpy contract build ./test/contracts/breakpoint
	drtpy contract build ./test/contracts/big-floats
	drtpy contract build ./test/contracts/counter
	drtpy contract build ./test/contracts/deployer
	drtpy contract build ./test/contracts/deployer-child
	drtpy contract build ./test/contracts/deployer-fromanother-contract
	drtpy contract build ./test/contracts/deployer-parent
	drtpy contract build ./test/contracts/vmhooks
	drtpy contract build ./test/contracts/erc20
	drtpy contract build ./test/contracts/exchange
	
	drtpy contract build ./test/contracts/exec-dest-ctx-builtin
	drtpy contract build ./test/contracts/exec-dest-ctx-by-caller/child
	drtpy contract build ./test/contracts/exec-dest-ctx-by-caller/parent
	drtpy contract build ./test/contracts/exec-dest-ctx-child
	drtpy contract build ./test/contracts/exec-dest-ctx-dcdt/basic
	drtpy contract build ./test/contracts/exec-dest-ctx-parent
	drtpy contract build ./test/contracts/exec-dest-ctx-recursive
	drtpy contract build ./test/contracts/exec-dest-ctx-recursive-child
	drtpy contract build ./test/contracts/exec-dest-ctx-recursive-parent
	
	drtpy contract build ./test/contracts/exec-same-ctx-child
	drtpy contract build ./test/contracts/exec-same-ctx-parent
	drtpy contract build ./test/contracts/exec-same-ctx-recursive
	drtpy contract build ./test/contracts/exec-same-ctx-recursive-child
	drtpy contract build ./test/contracts/exec-same-ctx-recursive-parent
	drtpy contract build ./test/contracts/exec-same-ctx-simple-child
	drtpy contract build ./test/contracts/exec-same-ctx-simple-parent
	
	drtpy contract build ./test/contracts/exec-sync-ctx-multiple/alpha
	drtpy contract build ./test/contracts/exec-sync-ctx-multiple/beta
	drtpy contract build ./test/contracts/exec-sync-ctx-multiple/delta
	drtpy contract build ./test/contracts/exec-sync-ctx-multiple/gamma
	
	drtpy contract build ./test/contracts/init-correct
	drtpy contract build ./test/contracts/init-simple
	drtpy contract build ./test/contracts/init-wrong
	drtpy contract build ./test/contracts/managed-buffers
	drtpy contract build ./test/contracts/misc
	drtpy contract build --no-optimization ./test/contracts/num-with-fp
	drtpy contract build ./test/contracts/promises
	drtpy contract build ./test/contracts/promises-train
	drtpy contract build ./test/contracts/promises-tracking
	drtpy contract build ./test/contracts/signatures
	drtpy contract build ./test/contracts/timelocks
	drtpy contract build ./test/contracts/upgrader-fromanother-contract

build-delegation:
ifndef SANDBOX
	$(error SANDBOX variable is undefined)
endif
	rm -rf ${SANDBOX}/sc-delegation-rs
	git clone --depth=1 --branch=master https://github.com/TerraDharitri/sc-delegation-rs.git ${SANDBOX}/sc-delegation-rs
	rm -rf ${SANDBOX}/sc-delegation-rs/.git
	drtpy contract build ${SANDBOX}/sc-delegation-rs
	drtpy contract test --directory="tests" ${SANDBOX}/sc-delegation-rs
	cp ${SANDBOX}/sc-delegation-rs/output/delegation.wasm ./test/delegation/delegation.wasm


build-dns:
ifndef SANDBOX
	$(error SANDBOX variable is undefined)
endif
	rm -rf ${SANDBOX}/sc-dns-rs
	git clone --depth=1 --branch=master https://github.com/TerraDharitri/sc-dns-rs.git ${SANDBOX}/sc-dns-rs
	rm -rf ${SANDBOX}/sc-dns-rs/.git
	drtpy contract build ${SANDBOX}/sc-dns-rs
	drtpy contract test --directory="tests" ${SANDBOX}/sc-dns-rs
	cp ${SANDBOX}/sc-dns-rs/output/dns.wasm ./test/dns/dns.wasm


build-sc-examples:
ifndef SANDBOX
	$(error SANDBOX variable is undefined)
endif
	rm -rf ${SANDBOX}/sc-examples

	drtpy contract new --template=erc20-c --directory ${SANDBOX}/sc-examples erc20-c
	drtpy contract build ${SANDBOX}/sc-examples/erc20-c
	cp ${SANDBOX}/sc-examples/erc20-c/output/wrc20.wasm ./test/erc20/contracts/erc20-c.wasm


build-sc-examples-rs:
ifndef SANDBOX
	$(error SANDBOX variable is undefined)
endif
	rm -rf ${SANDBOX}/sc-examples-rs
	
	drtpy contract new --template=simple-coin --directory ${SANDBOX}/sc-examples-rs simple-coin
	drtpy contract new --template=adder --directory ${SANDBOX}/sc-examples-rs adder
	drtpy contract build ${SANDBOX}/sc-examples-rs/adder
	drtpy contract build ${SANDBOX}/sc-examples-rs/simple-coin
	drtpy contract test ${SANDBOX}/sc-examples-rs/adder
	drtpy contract test ${SANDBOX}/sc-examples-rs/simple-coin
	cp ${SANDBOX}/sc-examples-rs/adder/output/adder.wasm ./test/adder/adder.wasm
	cp ${SANDBOX}/sc-examples-rs/simple-coin/output/simple-coin.wasm ./test/erc20/contracts/simple-coin.wasm

lint-install:
ifeq (,$(wildcard test -f bin/golangci-lint))
	@echo "Installing golint"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s
endif

run-lint:
	@echo "Running golint"
	bin/golangci-lint run --max-issues-per-linter 0 --max-same-issues 0 --timeout=2m

lint: lint-install run-lint
