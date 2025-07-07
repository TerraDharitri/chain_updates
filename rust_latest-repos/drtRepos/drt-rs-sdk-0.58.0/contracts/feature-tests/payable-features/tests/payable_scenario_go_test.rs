use dharitri_sc_scenario::*;

fn world() -> ScenarioWorld {
    ScenarioWorld::vm_go()
}

#[test]
fn call_value_check_go() {
    world().run("scenarios/call-value-check.scen.json");
}

#[test]
fn call_value_check_multi_rewa_go() {
    world().run("scenarios/call-value-check-multi-rewa.scen.json");
}

#[test]
fn payable_all_transfers_1_go() {
    world().run("scenarios/payable_all_transfers_1.scen.json");
}

#[test]
fn payable_all_transfers_2_go() {
    world().run("scenarios/payable_all_transfers_2.scen.json");
}

#[test]
fn payable_any_1_go() {
    world().run("scenarios/payable_any_1.scen.json");
}

#[test]
fn payable_any_2_go() {
    world().run("scenarios/payable_any_2.scen.json");
}

#[test]
fn payable_any_3_go() {
    world().run("scenarios/payable_any_3.scen.json");
}

#[test]
fn payable_any_4_go() {
    world().run("scenarios/payable_any_4.scen.json");
}

#[test]
fn payable_rewa_1_go() {
    world().run("scenarios/payable_rewa_1.scen.json");
}

#[test]
fn payable_rewa_2_go() {
    world().run("scenarios/payable_rewa_2.scen.json");
}

#[test]
fn payable_rewa_3_go() {
    world().run("scenarios/payable_rewa_3.scen.json");
}

#[test]
fn payable_rewa_4_go() {
    world().run("scenarios/payable_rewa_4.scen.json");
}

#[test]
fn payable_multi_array_go() {
    world().run("scenarios/payable_multi_array.scen.json");
}

#[test]
fn payable_multi_array_rewa_go() {
    world().run("scenarios/payable_multi_array_rewa.scen.json");
}

#[test]
fn payable_multiple_go() {
    world().run("scenarios/payable_multiple.scen.json");
}

#[test]
fn payable_multiple_rewa_go() {
    world().run("scenarios/payable_multiple_rewa.scen.json");
}

#[test]
fn payable_token_1_go() {
    world().run("scenarios/payable_token_1.scen.json");
}

#[test]
fn payable_token_2_go() {
    world().run("scenarios/payable_token_2.scen.json");
}

#[test]
fn payable_token_3_go() {
    world().run("scenarios/payable_token_3.scen.json");
}

#[test]
fn payable_token_4_go() {
    world().run("scenarios/payable_token_4.scen.json");
}

#[test]
fn payable_token_5_go() {
    world().run("scenarios/payable_token_5.scen.json");
}
