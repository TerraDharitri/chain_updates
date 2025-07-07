use client_common::{ClientStatus, GetLatestInfoResultType};
use common_types::{channel_types::height, ClientId, UnixTimestamp};

dharitri_sc::imports!();
dharitri_sc::derive_imports!();

#[dharitri_sc::module]
pub trait ViewsModule:
    client_common::CommonClientLogicModule
    + crate::client_logic::ClientLogicModule
    + host::host_views::HostViewsModule
    + host::storage::StorageModule
    + common_modules::utils::UtilsModule
{
    /// Timestamp is nanoseconds since unix epoch
    #[view(getTimestampAtHeight)]
    fn get_timestamp_at_height(
        &self,
        client_id: ClientId<Self::Api>,
        height: height::Data,
    ) -> UnixTimestamp {
        let mapper = self.consensus_states(&client_id, &height.to_biguint_concat());
        require!(!mapper.is_empty(), "Consensus state not found");

        let consensus_state = mapper.get();
        self.checked_timestamp_to_unix_mul(consensus_state.timestamp)
    }

    /// returns the latest height of the client state corresponding to `clientId`
    #[view(getLatestHeight)]
    fn get_latest_height(&self, client_id: &ClientId<Self::Api>) -> height::Data {
        let mapper = self.client_states(client_id);
        require!(!mapper.is_empty(), "Client state not found");

        let client_state = mapper.get();
        client_state.latest_height
    }

    /// returns the status of the client corresponding to `clientId`
    #[view(getStatus)]
    fn get_status(&self, _client_id: &ClientId<Self::Api>) -> ClientStatus {
        // TODO: Unsure why it's always considered active?
        ClientStatus::Active
    }

    /// returns the latest height, timestamp and status of the client corresponding to `clientId`
    #[view(getLatestInfo)]
    fn get_latest_info(&self, client_id: ClientId<Self::Api>) -> GetLatestInfoResultType {
        let latest_height = self.get_latest_height(&client_id);
        let consensus_state = self
            .consensus_states(&client_id, &latest_height.to_biguint_concat())
            .get();
        let client_status = self.get_status(&client_id);

        GetLatestInfoResultType {
            latest_height,
            latest_timestamp: consensus_state.timestamp,
            client_status,
        }
    }
}
