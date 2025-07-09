use std::{cell::RefCell, rc::Rc};

use host::Host;
use dharitri_sc::types::Address;
use dharitri_sc_scenario::{
    imports::{BlockchainStateWrapper, ContractObjWrapper},
    rust_biguint, DebugApi,
};

pub struct HostSetup<HostObjBuilder>
where
    HostObjBuilder: 'static + Copy + Fn() -> host::ContractObj<DebugApi>,
{
    pub b_mock: Rc<RefCell<BlockchainStateWrapper>>,
    pub host_owner: Address,
    pub host_wrapper: ContractObjWrapper<host::ContractObj<DebugApi>, HostObjBuilder>,
}

impl<HostObjBuilder> HostSetup<HostObjBuilder>
where
    HostObjBuilder: 'static + Copy + Fn() -> host::ContractObj<DebugApi>,
{
    pub fn new(
        b_mock: Rc<RefCell<BlockchainStateWrapper>>,
        owner: &Address,
        host_builder: HostObjBuilder,
    ) -> Self {
        let rust_zero = rust_biguint!(0u64);
        let host_wrapper =
            b_mock
                .borrow_mut()
                .create_sc_account(&rust_zero, Some(owner), host_builder, "host");
        b_mock
            .borrow_mut()
            .execute_tx(owner, &host_wrapper, &rust_zero, |sc| {
                sc.init();
            })
            .assert_ok();

        Self {
            b_mock,
            host_owner: owner.clone(),
            host_wrapper,
        }
    }
}
