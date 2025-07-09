use common_types::{channel_types::height, ClientId};

use crate::local_host_types::{client_state, consensus_state};

dharitri_sc::imports!();

// static CLIENT_TYPE: &[u8] = b"09-localhost";
static CLIENT_ID: &[u8] = b"09-localhost-0";

mod client_proxy {
    use common_types::ClientId;

    dharitri_sc::imports!();

    #[dharitri_sc::proxy]
    pub trait ClientProxy {
        #[endpoint(updateClientCommitments)]
        fn update_client_commitments(
            &self,
            client_id: ClientId<Self::Api>,
            encoded_heights: ManagedBuffer,
        );
    }
}

#[dharitri_sc::module]
pub trait ClientLogicModule: client_common::CommonClientLogicModule {
    /// initializes a new localhost client with the given client identifier, client state, and consensus state.
    ///
    /// `client_id`` the client identifier must be match with `CLIENT_ID`
    ///
    /// `client_state` the client state's latest height must be match with the current block number
    ///
    /// `consensus_state` the consensus state must be match with the sentinel consensus state (i.e. 0)
    #[endpoint(initializeClient)]
    fn initialize_client(
        &self,
        client_id: ClientId<Self::Api>,
        client_state: client_state::Data,
        consensus_state: consensus_state::Data,
    ) -> height::Data {
        self.require_ibc_handler_caller();
        self.require_valid_client_id(&client_id);
        require!(consensus_state.timestamp == 0, "Invalid consensus state");
        require!(
            client_state.latest_height.revision_number == 0,
            "Invalid revision number"
        );

        let current_block = self.blockchain().get_block_nonce();
        require!(
            client_state.latest_height.revision_height == current_block,
            "Invalid revision height"
        );

        height::Data {
            revision_number: 0,
            revision_height: current_block,
        }
    }

    /// updates the client state commitment with the current block number
    ///
    /// `client_id`` the client identifier must be match with `CLIENT_ID`
    #[endpoint(updateClient)]
    fn update_client(&self, client_id: ClientId<Self::Api>) -> ManagedVec<height::Data> {
        self.require_valid_client_id(&client_id);

        let ibc_handler = self.ibc_handler().get();
        let _: () = self
            .client_proxy(ibc_handler)
            .update_client_commitments(client_id, ManagedBuffer::new())
            .execute_on_dest_context();

        ManagedVec::new()
    }

    fn require_valid_client_id(&self, client_id: &ClientId<Self::Api>) {
        require!(client_id == CLIENT_ID, "Invalid client ID");
    }

    #[proxy]
    fn client_proxy(&self, sc_address: ManagedAddress) -> client_proxy::ClientProxy<Self::Api>;
}
