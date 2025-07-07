use dharitri_sc_scenario::*;

// These tests don't really test any contract, but the testing framework itself.

fn world() -> ScenarioWorld {
    ScenarioWorld::new()
}

/// Checks that externalSteps work fine.
#[test]
fn external_steps_rs() {
    world().run("tests/scenarios-self/external_steps/external_steps.scen.json");
}

#[test]
#[should_panic]
fn set_account_addr_len_err1_rs() {
    world().run("tests/scenarios-self/set-check/set-account-addr-len.err1.json");
}

#[test]
#[should_panic]
fn set_account_addr_len_err2_rs() {
    world().run("tests/scenarios-self/set-check/set-account-addr-len.err2.json");
}

#[test]
#[should_panic]
fn set_account_sc_addr_err1_rs() {
    world().run("tests/scenarios-self/set-check/set-account-sc-addr.err1.json");
}

#[test]
#[should_panic]
fn set_account_sc_addr_err2_rs() {
    world().run("tests/scenarios-self/set-check/set-account-sc-addr.err2.json");
}

#[test]
#[should_panic]
fn set_account_sc_addr_err3_rs() {
    world().run("tests/scenarios-self/set-check/set-account-sc-addr.err3.json");
}

#[test]
#[should_panic]
fn set_check_balance_err_rs() {
    world().run("tests/scenarios-self/set-check/set-check-balance.err.json");
}

#[test]
fn set_check_balance_rs() {
    world().run("tests/scenarios-self/set-check/set-check-balance.scen.json");
}

#[test]
#[should_panic]
fn set_check_code_err_rs() {
    world().run("tests/scenarios-self/set-check/set-check-code.err.json");
}

#[test]
fn set_check_code() {
    world().run("tests/scenarios-self/set-check/set-check-code.scen.json");
}

#[test]
#[should_panic]
fn set_check_codemetadata_err_rs() {
    world().run("tests/scenarios-self/set-check/set-check-codemetadata.err.json");
}

#[test]
fn set_check_codemetadata() {
    world().run("tests/scenarios-self/set-check/set-check-codemetadata.scen.json");
}

#[test]
#[should_panic]
fn set_check_dcdt_err_rs() {
    world().run("tests/scenarios-self/set-check/set-check-dcdt.err1.json");
}

#[test]
fn set_check_dcdt_rs() {
    world().run("tests/scenarios-self/set-check/set-check-dcdt.scen.json");
}

#[test]
#[should_panic]
fn set_check_nonce_err_rs() {
    world().run("tests/scenarios-self/set-check/set-check-nonce.err.json");
}

#[test]
fn set_check_nonce_rs() {
    world().run("tests/scenarios-self/set-check/set-check-nonce.scen.json");
}

#[test]
#[should_panic]
fn set_check_storage_err1_rs() {
    world().run("tests/scenarios-self/set-check/set-check-storage.err1.json");
}

#[test]
#[should_panic]
fn set_check_storage_err2_rs() {
    world().run("tests/scenarios-self/set-check/set-check-storage.err2.json");
}

#[test]
#[should_panic]
fn set_check_storage_err3_rs() {
    world().run("tests/scenarios-self/set-check/set-check-storage.err3.json");
}

#[test]
#[should_panic]
fn set_check_storage_err4_rs() {
    world().run("tests/scenarios-self/set-check/set-check-storage.err4.json");
}

#[test]
#[should_panic]
fn set_check_storage_err5_rs() {
    world().run("tests/scenarios-self/set-check/set-check-storage.err5.json");
}

#[test]
fn set_check_storage_rs() {
    world().run("tests/scenarios-self/set-check/set-check-storage.scen.json");
}

#[test]
#[should_panic]
fn set_check_username_err_rs() {
    world().run("tests/scenarios-self/set-check/set-check-username.err.json");
}

#[test]
fn set_check_username_rs() {
    world().run("tests/scenarios-self/set-check/set-check-username.scen.json");
}

#[test]
fn builtin_func_dcdt_transfer() {
    world().run("tests/scenarios-self/builtin-func-dcdt-transfer.scen.json");
}

#[test]
#[should_panic]
fn dcdt_non_zero_balance_check_err_rs() {
    world().run("tests/scenarios-self/dcdt-non-zero-balance-check-err.scen.json");
}

#[test]
#[should_panic]
fn dcdt_zero_balance_check_err_rs() {
    world().run("tests/scenarios-self/dcdt-zero-balance-check-err.scen.json");
}

#[test]
#[ignore = "TODO: not yet implemented in Rust VM"]
fn multi_transfer_rewa_rs() {
    world().run("tests/scenarios-self/multi-transfer-dcdt.scen.json");
}

#[test]
fn multi_transfer_dcdt_rs() {
    world().run("tests/scenarios-self/multi-transfer-dcdt.scen.json");
}

#[test]
fn transfer_rewa_rs() {
    world().run("tests/scenarios-self/transfer-rewa.scen.json");
}

#[test]
fn transfer_dcdt_rs() {
    world().run("tests/scenarios-self/transfer-dcdt.scen.json");
}

#[test]
fn validator_reward_rs() {
    world().run("tests/scenarios-self/validatorReward.scen.json");
}
