use client_common::ClientStatus;
use common_types::{channel_types::height, ClientId};

use crate::mock_types::{client_state, consensus_state, header};

dharitri_sc::imports!();

#[dharitri_sc::module]
pub trait ClientLogicModule: client_common::CommonClientLogicModule {
    /// creates a new client with the given state
    #[endpoint(initializeClient)]
    fn initialize_client(
        &self,
        client_id: ClientId<Self::Api>,
        client_state: client_state::Data,
        consensus_state: consensus_state::Data,
    ) -> height::Data {
        self.require_ibc_handler_caller();
        require!(
            client_state.latest_height.revision_number == 0
                && client_state.latest_height.revision_height != 0,
            "Invalid client state"
        );
        require!(consensus_state.timestamp != 0, "Invalid consensus state");

        let mapper = self.client_states(&client_id);
        require!(mapper.is_empty(), "Client already known");

        mapper.set(&client_state);
        self.consensus_states(&client_id, &client_state.latest_height.to_biguint_concat())
            .set(consensus_state);
        self.statuses(&client_id).set(ClientStatus::Active);

        client_state.latest_height
    }

    /// sets the status of the client corresponding to `clientId`
    #[only_owner]
    #[endpoint(setStatus)]
    fn set_status(&self, client_id: ClientId<Self::Api>, status: ClientStatus) {
        self.require_known_client(&client_id);

        self.statuses(&client_id).set(status);
    }

    /// updates the client state and returns the updated heights
    #[endpoint(updateClient)]
    fn update_client(
        &self,
        client_id: ClientId<Self::Api>,
        header: header::Data,
    ) -> ManagedVec<height::Data> {
        self.require_known_client(&client_id);

        let client_status = self.statuses(&client_id).get();
        require!(
            client_status == ClientStatus::Active,
            "Client is not active"
        );
        require!(
            header.height.revision_number == 0
                && header.height.revision_height != 0
                && header.timestamp != 0,
            "Invalid header"
        );

        let mapper = self.client_states(&client_id);
        let latest_client_state = mapper.get();
        if header.height > latest_client_state.latest_height {
            mapper.set(client_state::Data::new(header.height));
        }

        self.consensus_states(&client_id, &header.height.to_biguint_concat())
            .set(consensus_state::Data::new(header.timestamp));

        ManagedVec::from_single_item(header.height)
    }

    fn require_known_client(&self, client_id: &ClientId<Self::Api>) {
        require!(!self.client_states(client_id).is_empty(), "Unknown client");
    }

    #[storage_mapper("clientStates")]
    fn client_states(
        &self,
        client_id: &ClientId<Self::Api>,
    ) -> SingleValueMapper<client_state::Data>;

    #[storage_mapper("consensusStates")]
    fn consensus_states(
        &self,
        client_id: &ClientId<Self::Api>,
        height: &BigUint,
    ) -> SingleValueMapper<consensus_state::Data>;

    #[storage_mapper("statuses")]
    fn statuses(&self, client_id: &ClientId<Self::Api>) -> SingleValueMapper<ClientStatus>;
}
