use dharitri_sc_scenario::*;

fn world() -> ScenarioWorld {
    let mut blockchain = ScenarioWorld::new();
    blockchain.set_current_dir_from_workspace("dcdt-system-sc-mock");
    blockchain.register_contract(
        "drtsc:output/dcdt-system-sc-mock.drtsc.json",
        dcdt_system_sc_mock::ContractBuilder,
    );
    blockchain
}

#[test]
fn issue_rs() {
    world().run("scenarios/dcdt_system_sc.scen.json");
}
