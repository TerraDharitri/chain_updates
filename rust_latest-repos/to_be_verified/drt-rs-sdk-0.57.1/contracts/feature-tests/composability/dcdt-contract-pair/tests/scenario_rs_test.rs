use dharitri_sc_scenario::*;

fn world() -> ScenarioWorld {
    let mut blockchain = ScenarioWorld::new();
    blockchain
        .set_current_dir_from_workspace("contracts/feature-tests/composability/dcdt-contract-pair");
    blockchain.register_contract(
        "drtsc:first-contract/output/first-contract.drtsc.json",
        first_contract::ContractBuilder,
    );

    blockchain.register_contract(
        "drtsc:second-contract/output/second-contract.drtsc.json",
        second_contract::ContractBuilder,
    );
    blockchain
}

#[test]
fn init_rs() {
    world().run("scenarios/init.scen.json");
}

// TODO: implement DCDTTransfer + async call
#[test]
#[ignore]
fn reject_transfer_rs() {
    world().run("scenarios/reject_transfer.scen.json");
}

#[test]
fn simple_transfer_full_rs() {
    world().run("scenarios/simple_transfer_full.scen.json");
}

#[test]
fn simple_transfer_full_wrong_token_rs() {
    world().run("scenarios/simple_transfer_full_wrong_token.scen.json");
}

#[test]
fn simple_transfer_half_rs() {
    world().run("scenarios/simple_transfer_half.scen.json");
}
