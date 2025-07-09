# Dharitri SDK for Rust

[![Crates.io](https://img.shields.io/crates/v/dharitri-sdk)](https://crates.io/crates/dharitri-sdk)

General purpose collection of tools & SDKs to interact with the Dharitri blockchain from Rust projects.

## Example

```rust
use dharitri_sdk::blockchain::rpc::{CommunicationProxy, DEVNET_GATEWAY};

#[tokio::test]
async fn get_network_config() {
    let blockchain = CommunicationProxy::new(DEVNET_GATEWAY.to_string());
    let network_config = blockchain.get_network_config().await.unwrap();

    println!("network_config: {:?}", network_config)
}
```

More examples in `./examples`.
