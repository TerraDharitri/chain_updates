use dharitri_sc_scenario::*;

fn world() -> ScenarioWorld {
    ScenarioWorld::vm_go()
}

#[test]
fn multi_transfer_with_rewa_test() {
    world().run("tests/scenarios-self/multi_transfer_with_rewa.scen.json");
}
