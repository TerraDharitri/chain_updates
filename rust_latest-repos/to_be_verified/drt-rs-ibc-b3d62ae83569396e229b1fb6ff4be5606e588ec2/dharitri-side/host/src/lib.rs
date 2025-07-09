#![no_std]

dharitri_sc::imports!();

pub mod commitment;
pub mod host_config;
pub mod host_views;
pub mod module_manager;
pub mod storage;

#[dharitri_sc::contract]
pub trait Host:
    commitment::CommitmentModule
    + host_config::HostConfigModule
    + host_views::HostViewsModule
    + module_manager::ModuleManagerModule
    + storage::StorageModule
    + common_modules::client_lib::ClientLibModule
    + common_modules::host_lib::HostLibModule
    + common_modules::utils::UtilsModule
{
    #[init]
    fn init(&self) {}

    #[upgrade]
    fn upgrade(&self) {}
}
