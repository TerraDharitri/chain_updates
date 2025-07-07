use client_common::{
    ClientStatus, GetLatestInfoResultType, VerifyMembershipArgs, VerifyNonMembershipArgs,
};
use common_types::{channel_types::height, ClientId, Hash, UnixTimestamp};
use host::host_views::ProxyTrait as _;

use crate::mock_types::{client_state, consensus_state};

dharitri_sc::imports!();

#[dharitri_sc::module]
pub trait ViewsModule:
    client_common::CommonClientLogicModule + crate::client_logic::ClientLogicModule
{
    /// returns the timestamp of the consensus state at the given height
    ///
    /// The timestamp is nanoseconds since unix epoch
    #[view(getTimestampAtHeight)]
    fn get_timestamp_at_height(
        &self,
        client_id: &ClientId<Self::Api>,
        height: &height::Data,
    ) -> UnixTimestamp {
        let consensus_state = self.get_consensus_state(client_id, height);

        consensus_state.timestamp
    }

    /// returns the latest height of the client state corresponding to `clientId`
    #[view(getLatestHeight)]
    fn get_latest_height(&self, client_id: &ClientId<Self::Api>) -> height::Data {
        let client_state = self.get_client_state(client_id);

        client_state.latest_height
    }

    /// returns the status of the client corresponding to `clientId`
    ///
    /// A client status of "None" means the client is unknown
    #[view(getStatus)]
    fn get_status(&self, client_id: &ClientId<Self::Api>) -> ClientStatus {
        self.statuses(client_id).get()
    }

    /// returns the latest height, the latest timestamp, and the status of the client corresponding to `clientId`
    #[view(getLatestInfo)]
    fn get_latest_info(&self, client_id: ClientId<Self::Api>) -> GetLatestInfoResultType {
        let latest_height = self.get_latest_height(&client_id);
        let latest_timestamp = self.get_timestamp_at_height(&client_id, &latest_height);
        let client_status = self.get_status(&client_id);

        GetLatestInfoResultType {
            latest_height,
            latest_timestamp,
            client_status,
        }
    }

    /// A generic proof verification method which verifies a proof of the existence of a value at a given CommitmentPath at the specified height
    ///
    /// The caller is expected to construct the full CommitmentPath from a CommitmentPrefix and a standardized path (as defined in ICS 24)
    #[view(verifyMembership)]
    fn verify_membership(&self, args: VerifyMembershipArgs<Self::Api>) -> bool {
        let _ = self.get_timestamp_at_height(&args.client_id, &args.height);
        self.require_ibc_prefix(&args.prefix);

        let local_proof = self.encode_and_hash(&args.height, &args.prefix, &args.path, &args.value);
        local_proof == args.proof
    }

    /// A generic proof verification method which verifies the absence of a given CommitmentPath at a specified height
    ///
    /// The caller is expected to construct the full CommitmentPath from a CommitmentPrefix and a standardized path (as defined in ICS 24)
    #[view(verifyNonMembership)]
    fn verify_non_membership(&self, args: VerifyNonMembershipArgs<Self::Api>) -> bool {
        let _ = self.get_timestamp_at_height(&args.client_id, &args.height);
        self.require_ibc_prefix(&args.prefix);

        let local_proof = self.encode_and_hash(
            &args.height,
            &args.prefix,
            &args.path,
            &ManagedBuffer::new(),
        );
        local_proof == args.proof
    }

    /// returns the clientState corresponding to `clientId`
    #[view(getClientState)]
    fn get_client_state(&self, client_id: &ClientId<Self::Api>) -> client_state::Data {
        let mapper = self.client_states(client_id);
        require!(!mapper.is_empty(), "Client state not found");

        mapper.get()
    }

    #[view(getConsensusState)]
    fn get_consensus_state(
        &self,
        client_id: &ClientId<Self::Api>,
        height: &height::Data,
    ) -> consensus_state::Data {
        let mapper = self.consensus_states(client_id, &height.to_biguint_concat());
        require!(!mapper.is_empty(), "Consensus state not found");

        mapper.get()
    }

    fn require_ibc_prefix(&self, prefix: &ManagedBuffer) {
        let ibc_handler = self.ibc_handler().get();
        let ibc_prefix: ManagedBuffer = self
            .host_proxy(ibc_handler)
            .get_commitment_prefix()
            .execute_on_dest_context();
        require!(prefix == &ibc_prefix, "Invalid prefix");
    }

    fn encode_and_hash(
        &self,
        height: &height::Data,
        prefix: &ManagedBuffer,
        path: &ManagedBuffer,
        value: &ManagedBuffer,
    ) -> Hash<Self::Api> {
        let prefix_hash = self.crypto().sha256(prefix);
        let path_hash = self.crypto().sha256(path);
        let value_hash = self.crypto().sha256(value);

        // abi.encodePacked(height.toUint128(), sha256(prefix), sha256(path), sha256(value))
        let mut buffer = ManagedBuffer::new();
        let _ = height
            .to_biguint_concat::<Self::Api>()
            .top_encode(&mut buffer);
        let _ = prefix_hash.dep_encode(&mut buffer);
        let _ = path_hash.dep_encode(&mut buffer);
        let _ = value_hash.dep_encode(&mut buffer);

        self.crypto().sha256(buffer)
    }

    #[proxy]
    fn host_proxy(&self, sc_address: ManagedAddress) -> host::Proxy<Self::Api>;
}
