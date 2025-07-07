#![no_std]

dharitri_sc::imports!();

pub mod client_logic;
pub mod local_host_types;
pub mod views;

#[dharitri_sc::contract]
pub trait LocalHost:
    client_common::CommonClientLogicModule
    + client_logic::ClientLogicModule
    + views::ViewsModule
    + common_modules::utils::UtilsModule
{
    #[init]
    fn init(&self, ibc_handler: ManagedAddress) {
        self.set_ibc_handler(&ibc_handler);
    }

    #[upgrade]
    fn upgrade(&self) {}
}
