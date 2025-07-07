use dharitri_sdk::data::sdk_address::SdkAddress;
use dharitri_sdk_http::{GatewayHttpProxy, DEVNET_GATEWAY};

#[tokio::main]
async fn main() {
    let blockchain = GatewayHttpProxy::new(DEVNET_GATEWAY.to_string());
    let network_config = blockchain.get_network_config().await.unwrap();
    let addr = SdkAddress::from_bech32_string(
        "drt1qqqqqqqqqqqqqpgqfzydqmdw7m2vazsp6u5p95yxz76t2p9rd8ssj7kxgw",
    )
    .unwrap();

    let arg = blockchain
        .get_default_transaction_arguments(&addr.0, &network_config)
        .await
        .unwrap();

    assert!(arg.rcv_addr.to_bech32_string().is_ok());
    println!("default tx arg: {arg:#?}");
}
