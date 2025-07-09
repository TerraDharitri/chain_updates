# Smart contract `meta` crate support

[![crates.io](https://img.shields.io/crates/v/dharitri-sc-meta.svg)](https://crates.io/crates/dharitri-sc-meta)

The library that provides all the functionality of the individual contracts `meta` crates.

The purpose of the contract `meta` crates is to produce the contract ABI. Because of their access to the ABI, they have other ABI-based responsibilities, such as:
- generating the `wasm` crates,
- building the contracts,
- performing validations not possible otherwise,
- generating snippets,
- etc.

For more about the build process, see https://docs.dharitri.org/developers/developer-reference/sc-build-reference/
