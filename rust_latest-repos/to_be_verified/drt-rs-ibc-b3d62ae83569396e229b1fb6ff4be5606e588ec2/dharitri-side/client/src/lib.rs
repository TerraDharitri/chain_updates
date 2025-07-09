#![no_std]

dharitri_sc::imports!();

pub mod create_and_update_clients;
pub mod events;
pub mod merkle_proof;

#[dharitri_sc::contract]
pub trait Client:
    create_and_update_clients::CreateAndUpdateClientsModule
    + merkle_proof::MerkleProofModule
    + events::EventsModule
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
