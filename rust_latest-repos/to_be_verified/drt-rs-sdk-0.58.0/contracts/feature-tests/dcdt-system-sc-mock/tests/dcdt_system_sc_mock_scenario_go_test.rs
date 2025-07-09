use dharitri_sc_scenario::*;

fn world() -> ScenarioWorld {
    ScenarioWorld::vm_go()
}

#[test]
fn dcdt_system_sc_go() {
    world().run("scenarios/dcdt_system_sc.scen.json");
}
