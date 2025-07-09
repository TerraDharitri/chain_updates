#![no_std]

dharitri_sc::imports!();

pub mod common;
pub mod conn_endpoints;

#[dharitri_sc::contract]
pub trait Connection:
    common::conn_lib::ConnectionLibModule
    + common::verify_states::VerifyStatesModule
    + common::conn_internal::ConnectionInternalModule
    + conn_endpoints::ConnectionEndpointsModule
    + common::events::EventsModule
    + host::commitment::CommitmentModule
    + host::host_config::HostConfigModule
    + host::host_views::HostViewsModule
    + host::module_manager::ModuleManagerModule
    + host::storage::StorageModule
    + common_modules::client_lib::ClientLibModule
    + common_modules::host_lib::HostLibModule
    + common_modules::utils::UtilsModule
{
    #[init]
    fn init(&self) {}

    #[upgrade]
    fn upgrade(&self) {}
}
