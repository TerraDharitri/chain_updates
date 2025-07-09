#![no_std]

dharitri_sc::imports!();

pub mod client_logic;
pub mod qbft_types;
pub mod views;

// Ignore this whole crate, still work in progress. Not even entirely sure we need something like this yet.

#[dharitri_sc::contract]
pub trait Qbft:
    client_common::CommonClientLogicModule
    + client_logic::ClientLogicModule
    + views::ViewsModule
    + host::host_views::HostViewsModule
    + host::storage::StorageModule
    + common_modules::utils::UtilsModule
{
    #[init]
    fn init(&self, ibc_handler: ManagedAddress) {
        self.set_ibc_handler(&ibc_handler);
    }

    #[upgrade]
    fn upgrade(&self) {}
}
