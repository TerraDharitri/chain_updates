use common_types::{ClientType, PortId, UnixTimestamp};

use crate::storage::HostInfo;

dharitri_sc::imports!();

#[dharitri_sc::module]
pub trait HostConfigModule:
    crate::module_manager::ModuleManagerModule
    + crate::storage::StorageModule
    + common_modules::client_lib::ClientLibModule
    + common_modules::host_lib::HostLibModule
    + common_modules::utils::UtilsModule
{
    #[only_owner]
    #[endpoint(setExpectedTimePerBlock)]
    fn set_expected_time_per_block(&self, exp_time_per_block: UnixTimestamp) {
        let mapper = self.host_info();
        if !mapper.is_empty() {
            mapper.update(|host_info| host_info.expected_time_per_block = exp_time_per_block);

            return;
        }

        let default_host_value = HostInfo {
            expected_time_per_block: exp_time_per_block,
            ..Default::default()
        };
        mapper.set(default_host_value);
    }

    #[only_owner]
    #[endpoint(registerClient)]
    fn register_client(&self, client_type: ClientType<Self::Api>, client: ManagedAddress) {
        require!(self.is_valid_client_type(&client_type), "Invalid client ID");

        let mapper = self.client_registry(&client_type);
        require!(mapper.is_empty(), "Client already exists");
        self.require_valid_address(&client);

        mapper.set(client);
    }

    #[only_owner]
    #[endpoint(bindPort)]
    fn bind_port(&self, port_id: PortId<Self::Api>, module: ManagedAddress) {
        require!(self.is_valid_port_id(&port_id), "Invalid Port ID");
        self.require_valid_address(&module);

        self.claim_port_capability(&port_id, &module);
    }
}
