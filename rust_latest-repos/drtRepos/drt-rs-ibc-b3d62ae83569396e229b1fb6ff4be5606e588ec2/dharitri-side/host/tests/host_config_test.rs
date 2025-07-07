use std::{cell::RefCell, rc::Rc};

use host::{
    host_config::HostConfigModule,
    storage::{HostInfo, StorageModule},
};
use host_setup::HostSetup;
use dharitri_sc_scenario::{
    imports::BlockchainStateWrapper, managed_address, managed_buffer, rust_biguint, DebugApi,
};

pub mod host_setup;

fn get_host_setup<HostObjBuilder: 'static + Copy + Fn() -> host::ContractObj<DebugApi>>(
    host_builder: HostObjBuilder,
) -> HostSetup<HostObjBuilder> {
    let mut b_mock = BlockchainStateWrapper::new();
    let owner = b_mock.create_user_account(&rust_biguint!(0));
    HostSetup::new(Rc::new(RefCell::new(b_mock)), &owner, host_builder)
}

#[test]
fn set_expected_time_per_block_new_info_test() {
    let host_setup = get_host_setup(host::contract_obj);
    host_setup
        .b_mock
        .borrow_mut()
        .execute_tx(
            &host_setup.host_owner,
            &host_setup.host_wrapper,
            &rust_biguint!(0),
            |sc| {
                sc.set_expected_time_per_block(6);

                let host_info = sc.host_info().get();
                assert_eq!(
                    host_info,
                    HostInfo {
                        expected_time_per_block: 6,
                        ..Default::default()
                    }
                );
            },
        )
        .assert_ok();
}

#[test]
fn set_expected_time_per_block_update_info_test() {
    let host_setup = get_host_setup(host::contract_obj);
    host_setup
        .b_mock
        .borrow_mut()
        .execute_tx(
            &host_setup.host_owner,
            &host_setup.host_wrapper,
            &rust_biguint!(0),
            |sc| {
                sc.host_info().set(HostInfo {
                    next_channel_seq: 1,
                    next_client_seq: 2,
                    next_connection_seq: 3,
                    expected_time_per_block: 42,
                });
                sc.set_expected_time_per_block(6);

                let host_info = sc.host_info().get();
                assert_eq!(
                    host_info,
                    HostInfo {
                        next_channel_seq: 1,
                        next_client_seq: 2,
                        next_connection_seq: 3,
                        expected_time_per_block: 6,
                    }
                );
            },
        )
        .assert_ok();
}

#[test]
fn register_new_client_ok_test() {
    let host_setup = get_host_setup(host::contract_obj);
    let my_cool_client = host_setup
        .b_mock
        .borrow_mut()
        .create_user_account(&rust_biguint!(0));
    host_setup
        .b_mock
        .borrow_mut()
        .execute_tx(
            &host_setup.host_owner,
            &host_setup.host_wrapper,
            &rust_biguint!(0),
            |sc| {
                sc.register_client(
                    managed_buffer!(b"my-cool-client"),
                    managed_address!(&my_cool_client),
                );

                assert_eq!(
                    sc.client_registry(&managed_buffer!(b"my-cool-client"))
                        .get(),
                    managed_address!(&my_cool_client),
                );
            },
        )
        .assert_ok();
}

#[test]
fn try_register_already_existing_client_test() {
    let host_setup = get_host_setup(host::contract_obj);
    let my_cool_client = host_setup
        .b_mock
        .borrow_mut()
        .create_user_account(&rust_biguint!(0));
    host_setup
        .b_mock
        .borrow_mut()
        .execute_tx(
            &host_setup.host_owner,
            &host_setup.host_wrapper,
            &rust_biguint!(0),
            |sc| {
                sc.register_client(
                    managed_buffer!(b"my-cool-client"),
                    managed_address!(&my_cool_client),
                );

                assert_eq!(
                    sc.client_registry(&managed_buffer!(b"my-cool-client"))
                        .get(),
                    managed_address!(&my_cool_client),
                );
            },
        )
        .assert_ok();

    host_setup
        .b_mock
        .borrow_mut()
        .execute_tx(
            &host_setup.host_owner,
            &host_setup.host_wrapper,
            &rust_biguint!(0),
            |sc| {
                sc.register_client(
                    managed_buffer!(b"my-cool-client"),
                    managed_address!(&my_cool_client),
                );
            },
        )
        .assert_user_error("Client already exists");
}

#[test]
fn try_register_invalid_client_id_test() {
    let host_setup = get_host_setup(host::contract_obj);
    let my_cool_client = host_setup
        .b_mock
        .borrow_mut()
        .create_user_account(&rust_biguint!(0));
    host_setup
        .b_mock
        .borrow_mut()
        .execute_tx(
            &host_setup.host_owner,
            &host_setup.host_wrapper,
            &rust_biguint!(0),
            |sc| {
                sc.register_client(
                    managed_buffer!(b"EVIL-CLIENT"),
                    managed_address!(&my_cool_client),
                );
            },
        )
        .assert_user_error("Invalid client ID");
}

#[test]
fn bind_port_ok_test() {
    let host_setup = get_host_setup(host::contract_obj);
    let my_cool_module = host_setup
        .b_mock
        .borrow_mut()
        .create_user_account(&rust_biguint!(0));
    host_setup
        .b_mock
        .borrow_mut()
        .execute_tx(
            &host_setup.host_owner,
            &host_setup.host_wrapper,
            &rust_biguint!(0),
            |sc| {
                sc.bind_port(
                    managed_buffer!(b"my-cool-module"),
                    managed_address!(&my_cool_module),
                );

                assert_eq!(
                    sc.port_capabilities(&managed_buffer!(b"my-cool-module"))
                        .get(),
                    managed_address!(&my_cool_module),
                );
            },
        )
        .assert_ok();
}

#[test]
fn try_bind_same_port_twice_test() {
    let host_setup = get_host_setup(host::contract_obj);
    let my_cool_module = host_setup
        .b_mock
        .borrow_mut()
        .create_user_account(&rust_biguint!(0));
    host_setup
        .b_mock
        .borrow_mut()
        .execute_tx(
            &host_setup.host_owner,
            &host_setup.host_wrapper,
            &rust_biguint!(0),
            |sc| {
                sc.bind_port(
                    managed_buffer!(b"my-cool-module"),
                    managed_address!(&my_cool_module),
                );

                assert_eq!(
                    sc.port_capabilities(&managed_buffer!(b"my-cool-module"))
                        .get(),
                    managed_address!(&my_cool_module),
                );
            },
        )
        .assert_ok();

    host_setup
        .b_mock
        .borrow_mut()
        .execute_tx(
            &host_setup.host_owner,
            &host_setup.host_wrapper,
            &rust_biguint!(0),
            |sc| {
                sc.bind_port(
                    managed_buffer!(b"my-cool-module"),
                    managed_address!(&my_cool_module),
                );
            },
        )
        .assert_user_error("Port already claimed");
}

#[test]
fn try_bind_invalid_port_id() {
    let host_setup = get_host_setup(host::contract_obj);
    let my_cool_module = host_setup
        .b_mock
        .borrow_mut()
        .create_user_account(&rust_biguint!(0));
    host_setup
        .b_mock
        .borrow_mut()
        .execute_tx(
            &host_setup.host_owner,
            &host_setup.host_wrapper,
            &rust_biguint!(0),
            |sc| {
                sc.bind_port(
                    managed_buffer!(b"I'm @wesome!"),
                    managed_address!(&my_cool_module),
                );
            },
        )
        .assert_user_error("Invalid Port ID");
}
