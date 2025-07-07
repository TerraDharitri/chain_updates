# Dharitri Smart Contract Framework

The crates in this folder form the Dharitri smart contract framework.

They are as follows:
    - `dharitri-sc` - the base crate for smart contract libraries, it is the only dependency the smart contract code sees.
    - `dharitri-sc-derive` - procedural macros for friendlier SC code
    - `dharitri-sc-meta` - smart contract meta-programming: build system and other tools
    - `dharitri-sc-scenario` - the main testing tool, contracts are tested by via scenarios
    - `dharitri-sc-snippets` - base crate for tools that interact with the blockchain
    - `dharitri-sc-wasm-adapter` - the API that connects contracts to the WASM backend
