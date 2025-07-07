use common_types::{ClientId, UnixTimestamp};

dharitri_sc::imports!();

static DEFAULT_COMMITMENT_PREFIX: &[u8] = b"ibc";

#[dharitri_sc::module]
pub trait HostViewsModule:
    crate::storage::StorageModule + common_modules::utils::UtilsModule
{
    /// Returns the current timestamp (Unix time in nanoseconds) of the host chain.
    #[view(getHostTimestamp)]
    fn get_host_timestamp(&self) -> UnixTimestamp {
        let block_timestamp = self.blockchain().get_block_timestamp();

        self.checked_timestamp_to_unix_mul(block_timestamp)
    }

    #[view(getCommitmentPrefix)]
    fn get_commitment_prefix(&self) -> ManagedBuffer {
        ManagedBuffer::from(DEFAULT_COMMITMENT_PREFIX)
    }

    #[view(checkAndGetClient)]
    fn check_and_get_client(&self, client_id: &ClientId<Self::Api>) -> ManagedAddress {
        let client_info = self.try_get_client_info(client_id);

        client_info.client_impl
    }
}
