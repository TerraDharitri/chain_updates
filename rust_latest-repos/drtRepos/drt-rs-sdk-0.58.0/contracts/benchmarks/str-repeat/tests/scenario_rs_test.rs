use dharitri_sc_scenario::*;

fn world() -> ScenarioWorld {
    let mut blockchain = ScenarioWorld::new();
    blockchain.set_current_dir_from_workspace("contracts/benchmarks/str-repeat");
    blockchain.register_contract(
        "drtsc:output/str-repeat.drtsc.json",
        str_repeat::ContractBuilder,
    );
    blockchain.register_contract(
        "drtsc:output/str-repeat-mb-builder-basic.drtsc.json",
        str_repeat::ContractBuilder,
    );
    blockchain.register_contract(
        "drtsc:output/str-repeat-mb-builder-cached.drtsc.json",
        str_repeat::ContractBuilder,
    );
    blockchain
}

#[test]
fn mb_builder_basic_rs() {
    world().run("scenarios/mb_builder_basic.scen.json");
}

#[test]
fn mb_builder_cached_rs() {
    world().run("scenarios/mb_builder_cached.scen.json");
}

#[test]
fn str_repeat_rs() {
    world().run("scenarios/str_repeat.scen.json");
}
