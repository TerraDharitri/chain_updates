use dharitri_sc_scenario::imports::*;

fn world() -> ScenarioWorld {
    let mut blockchain = ScenarioWorld::new();
    blockchain.set_current_dir_from_workspace("contracts/feature-tests/composability");

    blockchain.register_contract(
        "drtsc:builtin-func-features/output/builtin-func-features.drtsc.json",
        builtin_func_features::ContractBuilder,
    );
    blockchain.register_contract(
        "drtsc:forwarder-queue/output/forwarder-queue.drtsc.json",
        forwarder_queue::ContractBuilder,
    );
    blockchain.register_contract(
        "drtsc:forwarder/output/forwarder.drtsc.json",
        forwarder::ContractBuilder,
    );
    blockchain.register_contract(
        "drtsc:forwarder-raw/output/forwarder-raw.drtsc.json",
        forwarder_raw::ContractBuilder,
    );
    blockchain.register_contract(
        "drtsc:promises-features/output/promises-features.drtsc.json",
        promises_features::ContractBuilder,
    );
    blockchain.register_contract(
        "drtsc:proxy-test-first/output/proxy-test-first.drtsc.json",
        proxy_test_first::ContractBuilder,
    );
    blockchain.register_contract(
        "drtsc:proxy-test-second/output/proxy-test-second.drtsc.json",
        proxy_test_second::ContractBuilder,
    );
    blockchain.register_contract(
        "drtsc:recursive-caller/output/recursive-caller.drtsc.json",
        recursive_caller::ContractBuilder,
    );

    let vault_sc_config =
        meta::multi_contract_config::<vault::AbiProvider>(&blockchain.current_dir().join("vault"));
    blockchain.register_contract_variant(
        "drtsc:vault/output/vault.drtsc.json",
        vault::ContractBuilder,
        vault_sc_config.find_contract("vault"),
    );
    blockchain.register_contract_variant(
        "drtsc:vault/output/vault-upgrade.drtsc.json",
        vault::ContractBuilder,
        vault_sc_config.find_contract("vault-upgrade"),
    );
    blockchain
}

#[test]
fn builtin_func_delete_user_name_rs() {
    world().run("scenarios/builtin_func_delete_user_name.scen.json");
}

#[test]
fn builtin_func_set_user_name_rs() {
    world().run("scenarios/builtin_func_set_user_name.scen.json");
}

#[test]
fn forw_queue_async_rs() {
    world().run("scenarios/forw_queue_async.scen.json");
}

#[test]
fn forw_raw_async_accept_rewa_rs() {
    world().run("scenarios/forw_raw_async_accept_rewa.scen.json");
}

#[test]
fn forw_raw_async_accept_dcdt_rs() {
    world().run("scenarios/forw_raw_async_accept_dcdt.scen.json");
}

#[test]
fn forw_raw_async_echo_rs() {
    world().run("scenarios/forw_raw_async_echo.scen.json");
}

#[test]
fn forw_raw_async_send_and_retrieve_multi_transfer_funds_rs() {
    world().run("scenarios/forw_raw_async_send_and_retrieve_multi_transfer_funds.scen.json");
}

#[test]
fn forw_raw_builtin_nft_local_mint_via_async_call_rs() {
    world().run("scenarios/forw_raw_builtin_nft_local_mint_via_async_call.scen.json");
}

#[test]
fn forw_raw_builtin_nft_local_mint_via_sync_call_rs() {
    world().run("scenarios/forw_raw_builtin_nft_local_mint_via_sync_call.scen.json");
}

#[test]
fn forw_raw_call_async_retrieve_multi_transfer_rs() {
    world().run("scenarios/forw_raw_call_async_retrieve_multi_transfer.scen.json");
}

#[test]
fn forw_raw_contract_deploy_rs() {
    world().run("scenarios/forw_raw_contract_deploy.scen.json");
}

#[test]
fn forw_raw_contract_upgrade_rs() {
    world().run("scenarios/forw_raw_contract_upgrade.scen.json");
}

#[test]
fn forw_raw_contract_upgrade_self_rs() {
    world().run("scenarios/forw_raw_contract_upgrade_self.scen.json");
}

#[test]
fn forw_raw_direct_rewa_rs() {
    world().run("scenarios/forw_raw_direct_rewa.scen.json");
}

#[test]
fn forw_raw_direct_dcdt_rs() {
    world().run("scenarios/forw_raw_direct_dcdt.scen.json");
}

#[test]
fn forw_raw_direct_multi_dcdt_rs() {
    world().run("scenarios/forw_raw_direct_multi_dcdt.scen.json");
}

#[test]
#[ignore = "not yet supported"]
fn forw_raw_init_async_rs() {
    world().run("scenarios/forw_raw_init_async.scen.json");
}

#[test]
#[ignore = "not yet supported"]
fn forw_raw_init_sync_accept_rewa_rs() {
    world().run("scenarios/forw_raw_init_sync_accept_rewa.scen.json");
}

#[test]
#[ignore = "not yet supported"]
fn forw_raw_init_sync_echo_rs() {
    world().run("scenarios/forw_raw_init_sync_echo.scen.json");
}

#[test]
fn forw_raw_sync_echo_rs() {
    world().run("scenarios/forw_raw_sync_echo.scen.json");
}

#[test]
fn forw_raw_sync_echo_caller_rs() {
    world().run("scenarios/forw_raw_sync_echo_caller.scen.json");
}

#[test]
fn forw_raw_sync_rewa_rs() {
    world().run("scenarios/forw_raw_sync_rewa.scen.json");
}

#[test]
fn forw_raw_sync_readonly_rs() {
    world().run("scenarios/forw_raw_sync_readonly.scen.json");
}

#[test]
#[ignore = "different error code/message in Rust VM"]
fn forw_raw_sync_reject_rs() {
    world().run("scenarios/forw_raw_sync_reject.scen.json");
}

#[test]
#[ignore]
fn forw_raw_sync_same_context_rs() {
    world().run("scenarios/forw_raw_sync_same_context.scen.json");
}

#[test]
#[ignore]
fn forw_raw_sync_same_context_rewa_rs() {
    world().run("scenarios/forw_raw_sync_same_context_rewa.scen.json");
}

#[test]
fn forw_raw_transf_exec_accept_rewa_rs() {
    world().run("scenarios/forw_raw_transf_exec_accept_rewa.scen.json");
}

#[test]
fn forw_raw_transf_exec_reject_rewa_rs() {
    world().run("scenarios/forw_raw_transf_exec_reject_rewa.scen.json");
}

#[test]
fn forwarder_builtin_nft_add_quantity_rs() {
    world().run("scenarios/forwarder_builtin_nft_add_quantity.scen.json");
}

#[test]
fn forwarder_builtin_nft_burn_rs() {
    world().run("scenarios/forwarder_builtin_nft_burn.scen.json");
}

#[test]
fn forwarder_builtin_nft_create_rs() {
    world().run("scenarios/forwarder_builtin_nft_create.scen.json");
}

#[test]
fn forwarder_builtin_nft_local_burn_rs() {
    world().run("scenarios/forwarder_builtin_nft_local_burn.scen.json");
}

#[test]
fn forwarder_builtin_nft_local_mint_rs() {
    world().run("scenarios/forwarder_builtin_nft_local_mint.scen.json");
}

#[test]
fn forwarder_call_async_accept_rewa_rs() {
    world().run("scenarios/forwarder_call_async_accept_rewa.scen.json");
}

#[test]
fn forwarder_call_async_accept_dcdt_rs() {
    world().run("scenarios/forwarder_call_async_accept_dcdt.scen.json");
}

#[test]
fn forwarder_call_async_accept_nft_rs() {
    world().run("scenarios/forwarder_call_async_accept_nft.scen.json");
}

#[test]
fn forwarder_call_async_multi_transfer_rs() {
    world().run("scenarios/forwarder_call_async_multi_transfer.scen.json");
}

#[test]
fn forwarder_call_async_multi_transfer_rewa_rs() {
    world().run("scenarios/forwarder_call_async_multi_transfer_rewa.scen.json");
}

#[test]
fn forwarder_call_async_retrieve_rewa_rs() {
    world().run("scenarios/forwarder_call_async_retrieve_rewa.scen.json");
}

#[test]
fn forwarder_call_async_retrieve_dcdt_rs() {
    world().run("scenarios/forwarder_call_async_retrieve_dcdt.scen.json");
}

#[test]
fn forwarder_call_async_retrieve_nft_rs() {
    world().run("scenarios/forwarder_call_async_retrieve_nft.scen.json");
}

#[test]
fn forwarder_call_sync_accept_rewa_rs() {
    world().run("scenarios/forwarder_call_sync_accept_rewa.scen.json");
}

#[test]
fn forwarder_call_sync_accept_dcdt_rs() {
    world().run("scenarios/forwarder_call_sync_accept_dcdt.scen.json");
}

#[test]
fn forwarder_call_sync_accept_multi_transfer_rs() {
    world().run("scenarios/forwarder_call_sync_accept_multi_transfer.scen.json");
}

#[test]
fn forwarder_call_sync_accept_multi_transfer_rewa_rs() {
    world().run("scenarios/forwarder_call_sync_accept_multi_transfer_rewa.scen.json");
}

#[test]
fn forwarder_call_sync_accept_nft_rs() {
    world().run("scenarios/forwarder_call_sync_accept_nft.scen.json");
}

#[test]
fn forwarder_call_sync_accept_then_read_rewa_rs() {
    world().run("scenarios/forwarder_call_sync_accept_then_read_rewa.scen.json");
}

#[test]
fn forwarder_call_sync_accept_then_read_dcdt_rs() {
    world().run("scenarios/forwarder_call_sync_accept_then_read_dcdt.scen.json");
}

#[test]
fn forwarder_call_sync_accept_then_read_nft_rs() {
    world().run("scenarios/forwarder_call_sync_accept_then_read_nft.scen.json");
}

#[test]
fn forwarder_call_sync_retrieve_rewa_rs() {
    world().run("scenarios/forwarder_call_sync_retrieve_rewa.scen.json");
}

#[test]
fn forwarder_call_sync_retrieve_rewa_bt_rs() {
    world().run("scenarios/forwarder_call_sync_retrieve_rewa_bt.scen.json");
}

#[test]
fn forwarder_call_sync_retrieve_dcdt_rs() {
    world().run("scenarios/forwarder_call_sync_retrieve_dcdt.scen.json");
}

#[test]
fn forwarder_call_sync_retrieve_dcdt_bt_rs() {
    world().run("scenarios/forwarder_call_sync_retrieve_dcdt_bt.scen.json");
}

#[test]
fn forwarder_call_sync_retrieve_nft_rs() {
    world().run("scenarios/forwarder_call_sync_retrieve_nft.scen.json");
}

#[test]
fn forwarder_call_sync_retrieve_nft_bt_rs() {
    world().run("scenarios/forwarder_call_sync_retrieve_nft_bt.scen.json");
}

#[test]
fn forwarder_call_transf_exec_accept_rewa_rs() {
    world().run("scenarios/forwarder_call_transf_exec_accept_rewa.scen.json");
}

#[test]
fn forwarder_call_transf_exec_accept_rewa_twice_rs() {
    world().run("scenarios/forwarder_call_transf_exec_accept_rewa_twice.scen.json");
}

#[test]
fn forwarder_call_transf_exec_accept_dcdt_rs() {
    world().run("scenarios/forwarder_call_transf_exec_accept_dcdt.scen.json");
}

#[test]
fn forwarder_call_transf_exec_accept_dcdt_twice_rs() {
    world().run("scenarios/forwarder_call_transf_exec_accept_dcdt_twice.scen.json");
}

#[test]
fn forwarder_call_transf_exec_accept_multi_transfer_rs() {
    world().run("scenarios/forwarder_call_transf_exec_accept_multi_transfer.scen.json");
}

#[test]
fn forwarder_call_transf_exec_accept_multi_transfer_rewa_rs() {
    world().run("scenarios/forwarder_call_transf_exec_accept_multi_transfer_rewa.scen.json");
}

#[test]
fn forwarder_call_transf_exec_accept_nft_rs() {
    world().run("scenarios/forwarder_call_transf_exec_accept_nft.scen.json");
}

#[test]
fn forwarder_call_transf_exec_accept_return_values_rs() {
    world().run("scenarios/forwarder_call_transf_exec_accept_return_values.scen.json");
}

#[test]
fn forwarder_call_transf_exec_accept_sft_twice_rs() {
    world().run("scenarios/forwarder_call_transf_exec_accept_sft_twice.scen.json");
}

#[test]
fn forwarder_call_transf_exec_reject_multi_transfer_rs() {
    world().run("scenarios/forwarder_call_transf_exec_reject_multi_transfer.scen.json");
}

#[test]
fn forwarder_call_transf_exec_reject_nft_rs() {
    world().run("scenarios/forwarder_call_transf_exec_reject_nft.scen.json");
}

#[test]
fn forwarder_contract_change_owner_rs() {
    world().run("scenarios/forwarder_contract_change_owner.scen.json");
}

#[test]
fn forwarder_contract_deploy_rs() {
    world().run("scenarios/forwarder_contract_deploy.scen.json");
}

#[test]
fn forwarder_contract_upgrade_rs() {
    world().run("scenarios/forwarder_contract_upgrade.scen.json");
}

#[test]
fn forwarder_get_dcdt_local_roles_rs() {
    world().run("scenarios/forwarder_get_dcdt_local_roles.scen.json");
}

#[test]
fn forwarder_get_dcdt_token_data_rs() {
    world().run("scenarios/forwarder_get_dcdt_token_data.scen.json");
}

#[test]
fn forwarder_nft_add_uri_rs() {
    world().run("scenarios/forwarder_nft_add_uri.scen.json");
}

#[test]
fn forwarder_nft_create_rs() {
    world().run("scenarios/forwarder_nft_create.scen.json");
}

#[test]
fn forwarder_nft_create_and_send_rs() {
    world().run("scenarios/forwarder_nft_create_and_send.scen.json");
}

#[test]
fn forwarder_nft_current_nonce_rs() {
    world().run("scenarios/forwarder_nft_current_nonce.scen.json");
}

#[test]
fn forwarder_nft_decode_complex_attributes_rs() {
    world().run("scenarios/forwarder_nft_decode_complex_attributes.scen.json");
}

#[test]
fn forwarder_nft_transfer_async_rs() {
    world().run("scenarios/forwarder_nft_transfer_async.scen.json");
}

#[test]
fn forwarder_nft_transfer_exec_rs() {
    world().run("scenarios/forwarder_nft_transfer_exec.scen.json");
}

#[test]
fn forwarder_nft_update_attributes_rs() {
    world().run("scenarios/forwarder_nft_update_attributes.scen.json");
}

#[test]
fn forwarder_no_endpoint_rs() {
    world().run("scenarios/forwarder_no_endpoint.scen.json");
}

#[test]
fn forwarder_retrieve_funds_with_accept_func_rs() {
    world().run("scenarios/forwarder_retrieve_funds_with_accept_func.scen.json");
}

#[test]
fn forwarder_send_dcdt_multi_transfer_rs() {
    world().run("scenarios/forwarder_send_dcdt_multi_transfer.scen.json");
}

#[test]
fn forwarder_sync_echo_rs() {
    world().run("scenarios/forwarder_sync_echo.scen.json");
}

#[test]
fn forwarder_transfer_dcdt_with_fees_rs() {
    world().run("scenarios/forwarder_transfer_dcdt_with_fees.scen.json");
}

#[test]
fn forwarder_validate_token_identifier_rs() {
    world().run("scenarios/forwarder_validate_token_identifier.scen.json");
}

#[test]
fn promises_call_async_accept_rewa_rs() {
    world().run("scenarios/promises_call_async_accept_rewa.scen.json");
}

#[test]
fn promises_call_async_accept_dcdt_rs() {
    world().run("scenarios/promises_call_async_accept_dcdt.scen.json");
}

#[test]
fn promises_call_async_retrieve_rewa_rs() {
    world().run("scenarios/promises_call_async_retrieve_rewa.scen.json");
}

#[test]
fn promises_call_async_retrieve_dcdt_rs() {
    world().run("scenarios/promises_call_async_retrieve_dcdt.scen.json");
}

#[test]
fn promises_call_callback_directly_rs() {
    world().run("scenarios/promises_call_callback_directly.scen.json");
}

#[test]
fn promises_call_transfer_callback_rewa_rs() {
    world().run("scenarios/promises_call_transfer_callback_rewa.scen.json");
}

#[test]
fn promises_call_transfer_callback_dcdt_rs() {
    world().run("scenarios/promises_call_transfer_callback_dcdt.scen.json");
}

#[test]
fn promises_multi_transfer_rs() {
    world().run("scenarios/promises_multi_transfer.scen.json");
}

#[test]
#[ignore = "TODO - some log is missing"]
fn promises_multi_transfer_err_rs() {
    world().run("scenarios/promises_multi_transfer_err.scen.json");
}

#[test]
#[ignore = "gas"]
fn promises_single_transfer_rs() {
    world().run("scenarios/promises_single_transfer.scen.json");
}

#[test]
#[ignore = "gas"]
fn promises_single_transfer_gas_1_rs() {
    world().run("scenarios/promises_single_transfer_gas1.scen.json");
}

#[test]
#[ignore = "gas"]
fn promises_single_transfer_gas_2_rs() {
    world().run("scenarios/promises_single_transfer_gas2.scen.json");
}

#[test]
fn proxy_test_init_rs() {
    world().run("scenarios/proxy_test_init.scen.json");
}

#[test]
fn proxy_test_message_other_shard_rs() {
    world().run("scenarios/proxy_test_message_otherShard.scen.json");
}

#[test]
fn proxy_test_message_other_shard_callback_rs() {
    world().run("scenarios/proxy_test_message_otherShard_callback.scen.json");
}

#[test]
fn proxy_test_message_same_shard_rs() {
    world().run("scenarios/proxy_test_message_sameShard.scen.json");
}

#[test]
fn proxy_test_message_same_shard_callback_rs() {
    world().run("scenarios/proxy_test_message_sameShard_callback.scen.json");
}

#[test]
fn proxy_test_payment_other_shard_rs() {
    world().run("scenarios/proxy_test_payment_otherShard.scen.json");
}

#[test]
fn proxy_test_payment_other_shard_callback_rs() {
    world().run("scenarios/proxy_test_payment_otherShard_callback.scen.json");
}

#[test]
fn proxy_test_payment_same_shard_rs() {
    world().run("scenarios/proxy_test_payment_sameShard.scen.json");
}

#[test]
fn proxy_test_payment_same_shard_callback_rs() {
    world().run("scenarios/proxy_test_payment_sameShard_callback.scen.json");
}

#[test]
fn proxy_test_upgrade_rs() {
    world().run("scenarios/proxy_test_upgrade.scen.json");
}

#[test]
fn recursive_caller_rewa_1_rs() {
    world().run("scenarios/recursive_caller_rewa_1.scen.json");
}

#[test]
fn recursive_caller_dcdt_1_rs() {
    world().run("scenarios/recursive_caller_dcdt_1.scen.json");
}

#[test]
fn send_rewa_rs() {
    world().run("scenarios/send_rewa.scen.json");
}

#[test]
fn send_dcdt_rs() {
    world().run("scenarios/send_dcdt.scen.json");
}
