use forwarder_interact::forwarder_cli;
use dharitri_sc_snippets::imports::*;

#[tokio::main]
async fn main() {
    forwarder_cli().await;
}
