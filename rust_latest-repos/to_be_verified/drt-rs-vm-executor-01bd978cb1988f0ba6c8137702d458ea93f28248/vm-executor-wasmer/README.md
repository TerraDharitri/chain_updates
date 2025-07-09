# dharitri-vm-executor-wasmer

Wasmer adapter for the Dharitri VM, used in production.

Currently used in 2 ways:
- Imported in the CAPI project that is called from dharitri-chain-vm-go, the Go implementation of the SpaceVM.
- Imported in the Rust VM, via a new set of traits. 

Wasmer 2.2 is currently used. The Wasmer implementation is forked from the official one, since it contains several adaptations and fixes.
