use dharitri_sc_scenario::imports::*;

fn world() -> ScenarioWorld {
    let mut blockchain = ScenarioWorld::new().executor_config(ExecutorConfig::full_suite());

    blockchain.set_current_dir_from_workspace("contracts/feature-tests/managed-map-features");
    blockchain.register_contract(
        "drtsc:output/managed-map-features.drtsc.json",
        managed_map_features::ContractBuilder,
    );
    blockchain
}

#[test]
fn mmap_get_rs() {
    world().run("scenarios/mmap_get.scen.json");
}

#[test]
fn mmap_remove_rs() {
    world().run("scenarios/mmap_remove.scen.json");
}
