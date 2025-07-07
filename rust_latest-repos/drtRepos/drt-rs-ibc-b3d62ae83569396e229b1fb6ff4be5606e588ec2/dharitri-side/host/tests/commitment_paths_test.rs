use std::{cell::RefCell, rc::Rc};

use host::commitment::CommitmentModule;
use host_setup::HostSetup;
use dharitri_sc_scenario::{
    imports::BlockchainStateWrapper, managed_buffer, rust_biguint, DebugApi,
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
fn client_state_path_test() {
    let host_setup = get_host_setup(host::contract_obj);
    host_setup
        .b_mock
        .borrow_mut()
        .execute_tx(
            &host_setup.host_owner,
            &host_setup.host_wrapper,
            &rust_biguint!(0),
            |sc| {
                let client_state_path = sc.get_client_state_path(&managed_buffer!(b"10"));
                assert_eq!(client_state_path, b"clients/10/clientState");
            },
        )
        .assert_ok();
}

#[test]
fn consensus_state_path_test() {
    let host_setup = get_host_setup(host::contract_obj);
    host_setup
        .b_mock
        .borrow_mut()
        .execute_tx(
            &host_setup.host_owner,
            &host_setup.host_wrapper,
            &rust_biguint!(0),
            |sc| {
                let consensus_state_path =
                    sc.get_consensus_state_path(&managed_buffer!(b"10"), 20, 30);
                assert_eq!(consensus_state_path, b"clients/10/consensusStates/20-30");
            },
        )
        .assert_ok();
}

#[test]
fn connection_path_test() {
    let host_setup = get_host_setup(host::contract_obj);
    host_setup
        .b_mock
        .borrow_mut()
        .execute_tx(
            &host_setup.host_owner,
            &host_setup.host_wrapper,
            &rust_biguint!(0),
            |sc| {
                let connection_path = sc.get_connection_path(&managed_buffer!(b"10"));
                assert_eq!(connection_path, b"connections/10");
            },
        )
        .assert_ok();
}

#[test]
fn channel_path_test() {
    let host_setup = get_host_setup(host::contract_obj);
    host_setup
        .b_mock
        .borrow_mut()
        .execute_tx(
            &host_setup.host_owner,
            &host_setup.host_wrapper,
            &rust_biguint!(0),
            |sc| {
                let channel_path =
                    sc.get_channel_path(&managed_buffer!(b"10"), &managed_buffer!(b"20"));
                assert_eq!(channel_path, b"channelEnds/ports/10/channels/20");
            },
        )
        .assert_ok();
}

#[test]
fn packet_commitment_path_test() {
    let host_setup = get_host_setup(host::contract_obj);
    host_setup
        .b_mock
        .borrow_mut()
        .execute_tx(
            &host_setup.host_owner,
            &host_setup.host_wrapper,
            &rust_biguint!(0),
            |sc| {
                let packet_commitment_path = sc.get_packet_commitment_path(
                    &managed_buffer!(b"10"),
                    &managed_buffer!(b"20"),
                    30,
                );
                assert_eq!(
                    packet_commitment_path,
                    b"commitments/ports/10/channels/20/sequences/30"
                );
            },
        )
        .assert_ok();
}

#[test]
fn packet_acknowledgement_commitment_path_test() {
    let host_setup = get_host_setup(host::contract_obj);
    host_setup
        .b_mock
        .borrow_mut()
        .execute_tx(
            &host_setup.host_owner,
            &host_setup.host_wrapper,
            &rust_biguint!(0),
            |sc| {
                let packet_acknowledgement_commitment_path = sc
                    .get_packet_acknowledgement_commitment_path(
                        &managed_buffer!(b"10"),
                        &managed_buffer!(b"20"),
                        30,
                    );
                assert_eq!(
                    packet_acknowledgement_commitment_path,
                    b"acks/ports/10/channels/20/sequences/30"
                );
            },
        )
        .assert_ok();
}

#[test]
fn packet_receipt_commitment_path_test() {
    let host_setup = get_host_setup(host::contract_obj);
    host_setup
        .b_mock
        .borrow_mut()
        .execute_tx(
            &host_setup.host_owner,
            &host_setup.host_wrapper,
            &rust_biguint!(0),
            |sc| {
                let packet_receipt_commitment_path = sc.get_packet_receipt_commitment_path(
                    &managed_buffer!(b"10"),
                    &managed_buffer!(b"20"),
                    30,
                );
                assert_eq!(
                    packet_receipt_commitment_path,
                    b"receipts/ports/10/channels/20/sequences/30"
                );
            },
        )
        .assert_ok();
}

#[test]
fn next_seq_send_commitment_path_test() {
    let host_setup = get_host_setup(host::contract_obj);
    host_setup
        .b_mock
        .borrow_mut()
        .execute_tx(
            &host_setup.host_owner,
            &host_setup.host_wrapper,
            &rust_biguint!(0),
            |sc| {
                let next_seq_send_commitment_path = sc.get_next_seq_send_commitment_path(
                    &managed_buffer!(b"10"),
                    &managed_buffer!(b"20"),
                );
                assert_eq!(
                    next_seq_send_commitment_path,
                    b"nextSequenceSend/ports/10/channels/20"
                );
            },
        )
        .assert_ok();
}

#[test]
fn next_seq_recv_commitment_path_test() {
    let host_setup = get_host_setup(host::contract_obj);
    host_setup
        .b_mock
        .borrow_mut()
        .execute_tx(
            &host_setup.host_owner,
            &host_setup.host_wrapper,
            &rust_biguint!(0),
            |sc| {
                let next_seq_recv_commitment_path = sc.get_next_seq_recv_commitment_path(
                    &managed_buffer!(b"10"),
                    &managed_buffer!(b"20"),
                );
                assert_eq!(
                    next_seq_recv_commitment_path,
                    b"nextSequenceRecv/ports/10/channels/20"
                );
            },
        )
        .assert_ok();
}

#[test]
fn next_seq_ack_commitment_path_test() {
    let host_setup = get_host_setup(host::contract_obj);
    host_setup
        .b_mock
        .borrow_mut()
        .execute_tx(
            &host_setup.host_owner,
            &host_setup.host_wrapper,
            &rust_biguint!(0),
            |sc| {
                let next_seq_ack_commitment_path = sc.get_next_seq_ack_commitment_path(
                    &managed_buffer!(b"10"),
                    &managed_buffer!(b"20"),
                );
                assert_eq!(
                    next_seq_ack_commitment_path,
                    b"nextSequenceAck/ports/10/channels/20"
                );
            },
        )
        .assert_ok();
}

#[test]
fn channel_upgrade_path_test() {
    let host_setup = get_host_setup(host::contract_obj);
    host_setup
        .b_mock
        .borrow_mut()
        .execute_tx(
            &host_setup.host_owner,
            &host_setup.host_wrapper,
            &rust_biguint!(0),
            |sc| {
                let channel_upgrade_path =
                    sc.get_channel_upgrade_path(&managed_buffer!(b"10"), &managed_buffer!(b"20"));
                assert_eq!(
                    channel_upgrade_path,
                    b"channelUpgrades/upgrades/ports/10/channels/20"
                );
            },
        )
        .assert_ok();
}

#[test]
fn channel_upgrade_error_path_test() {
    let host_setup = get_host_setup(host::contract_obj);
    host_setup
        .b_mock
        .borrow_mut()
        .execute_tx(
            &host_setup.host_owner,
            &host_setup.host_wrapper,
            &rust_biguint!(0),
            |sc| {
                let channel_upgrade_error_path = sc.get_channel_upgrade_error_path(
                    &managed_buffer!(b"10"),
                    &managed_buffer!(b"20"),
                );
                assert_eq!(
                    channel_upgrade_error_path,
                    b"channelUpgrades/upgradeError/ports/10/channels/20"
                );
            },
        )
        .assert_ok();
}
