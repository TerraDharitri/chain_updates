# dharitri-vm-executor-experimental

Experimental Wasmer adapter for the Dharitri VM, not used in production.

This adapter is currently only usable from the Rust VM, via a new set of traits.

The plan is to stabilize this implementation enough, so that it can be used in production. This implementation is intended to eventually replace the one in `dharitri-vm-executor-wasmer`.

Wasmer 6.0 is used in this implementation, and is not forked. However, several of the fixes used in production were not yet applied to this version. It remains to be seen if the official Wasmer implementation will remain in use, or if another fork will be necessary.
