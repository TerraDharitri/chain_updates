use std::{cell::RefCell, rc::Rc};

use host_setup::HostSetup;
use dharitri_sc_scenario::{imports::BlockchainStateWrapper, rust_biguint};

pub mod host_setup;

#[test]
fn setup_test() {
    let mut b_mock = BlockchainStateWrapper::new();
    let owner = b_mock.create_user_account(&rust_biguint!(0));
    let _ = HostSetup::new(Rc::new(RefCell::new(b_mock)), &owner, host::contract_obj);
}
