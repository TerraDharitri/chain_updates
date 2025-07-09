mod gateway_http_proxy;

pub use gateway_http_proxy::GatewayHttpProxy;

pub use dharitri_sdk as core;

pub const MAINNET_GATEWAY: &str = "https://gateway.dharitri.org";
pub const TESTNET_GATEWAY: &str = "https://testnet-gateway.dharitri.org";
pub const DEVNET_GATEWAY: &str = "https://devnet-gateway.dharitri.org";
pub const CHAIN_SIMULATOR_GATEWAY: &str = "http://localhost:8085";

// MetachainShardId will be used to identify a shard ID as metachain
pub const METACHAIN_SHARD_ID: u32 = 0xFFFFFFFF;

pub const DEFAULT_USE_CHAIN_SIMULATOR: bool = false;
