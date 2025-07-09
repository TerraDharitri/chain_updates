use common_types::{channel_types::height, ClientId, ClientType};
use host::storage::ClientInfo;

dharitri_sc::imports!();
dharitri_sc::derive_imports!();

#[derive(TypeAbi, TopDecode)]
pub struct MsgCreateClient<M: ManagedTypeApi> {
    pub client_type: ClientType<M>,
    pub encoded_client_state: ManagedBuffer<M>,
    pub encoded_consensus_state: ManagedBuffer<M>,
}

#[derive(TypeAbi, TopDecode)]
pub struct MsgUpdateClient<M: ManagedTypeApi> {
    pub client_id: ClientId<M>,
    pub encoded_client_message: ManagedBuffer<M>,
}

mod generic_client_proxy {
    use common_types::{channel_types::height, ClientId};

    dharitri_sc::imports!();

    #[dharitri_sc::proxy]
    pub trait GenericClientProxy {
        #[endpoint(initializeClient)]
        fn initialize_client(
            &self,
            encoded_client_state: ManagedBuffer,
            encoded_consensus_state: ManagedBuffer,
        ) -> height::Data;

        #[endpoint(updateClient)]
        fn update_client(
            &self,
            client_id: ClientId<Self::Api>,
            encoded_client_message: ManagedBuffer,
        ) -> ManagedVec<height::Data>;

        #[view(getClientState)]
        fn get_client_state(&self, client_id: &ClientId<Self::Api>) -> ManagedBuffer;

        #[view(getConsensusState)]
        fn get_consensus_state(
            &self,
            client_id: &ClientId<Self::Api>,
            height: &height::Data,
        ) -> ManagedBuffer;
    }
}

#[dharitri_sc::module]
pub trait CreateAndUpdateClientsModule:
    crate::events::EventsModule
    + host::commitment::CommitmentModule
    + host::host_config::HostConfigModule
    + host::host_views::HostViewsModule
    + host::module_manager::ModuleManagerModule
    + host::storage::StorageModule
    + common_modules::client_lib::ClientLibModule
    + common_modules::host_lib::HostLibModule
    + common_modules::utils::UtilsModule
{
    /// creates a new client state and populates it with a given consensus state
    #[endpoint(createClient)]
    fn create_client(&self, args: MsgCreateClient<Self::Api>) -> ClientId<Self::Api> {
        let client_impl_mapper = self.client_registry(&args.client_type);
        require!(!client_impl_mapper.is_empty(), "Client not registered");

        let client_impl = client_impl_mapper.get();
        let client_id = self.generate_client_identifier(&args.client_type);
        self.client_info(&client_id).set(ClientInfo {
            client_type: args.client_type.clone(),
            client_impl: client_impl.clone(),
        });

        self.update_commitments_after_create(args, client_impl, &client_id);
        self.generated_client_id_event(&client_id);

        client_id
    }

    /// updates the consensus state and the state root from a provided header
    #[endpoint(updateClient)]
    fn update_client(&self, args: MsgUpdateClient<Self::Api>) {
        let client_impl = self.check_and_get_client(&args.client_id);
        let heights: ManagedVec<height::Data> = self
            .generic_client_proxy_impl(client_impl)
            .update_client(args.client_id.clone(), args.encoded_client_message)
            .execute_on_dest_context();

        if !heights.is_empty() {
            self.update_client_commitments(args.client_id, heights);
        }
    }

    #[endpoint(updateClientCommitments)]
    fn update_client_commitments(
        &self,
        client_id: ClientId<Self::Api>,
        heights: ManagedVec<height::Data>,
    ) {
        let client = self.check_and_get_client(&client_id);
        let encoded_client_state: ManagedBuffer = self
            .generic_client_proxy_impl(client.clone())
            .get_client_state(client_id.clone())
            .execute_on_dest_context();

        let client_state_comm_key = self.get_client_state_commitment_key(&client_id);
        let client_state_hash = self.crypto().keccak256(encoded_client_state);
        self.commitments(&client_state_comm_key)
            .set(client_state_hash);

        for height in &heights {
            self.update_single_commitment(client.clone(), &client_id, &height);
        }
    }

    fn generate_client_identifier(
        &self,
        client_type: &ClientType<Self::Api>,
    ) -> ClientId<Self::Api> {
        let next_client_seq = self.get_next_client_seq();

        sc_format!("{}-{}", client_type, next_client_seq)
    }

    fn update_commitments_after_create(
        &self,
        args: MsgCreateClient<Self::Api>,
        client_impl: ManagedAddress,
        client_id: &ClientId<Self::Api>,
    ) {
        let client_state_hash = self.crypto().keccak256(&args.encoded_client_state);
        let consensus_state_hash = self.crypto().keccak256(&args.encoded_consensus_state);
        let height: height::Data = self
            .generic_client_proxy_impl(client_impl)
            .initialize_client(args.encoded_client_state, args.encoded_consensus_state)
            .execute_on_dest_context();

        let client_comm_key = self.get_client_state_commitment_key(client_id);
        let consensus_comm_key = self.get_consensus_state_commitment_key(
            client_id,
            height.revision_number,
            height.revision_height,
        );
        self.commitments(&client_comm_key).set(&client_state_hash);
        self.commitments(&consensus_comm_key)
            .set(&consensus_state_hash);
    }

    fn update_single_commitment(
        &self,
        client: ManagedAddress,
        client_id: &ClientId<Self::Api>,
        height: &height::Data,
    ) {
        let encoded_consensus_state: ManagedBuffer = self
            .generic_client_proxy_impl(client)
            .get_consensus_state(client_id, height)
            .execute_on_dest_context();

        let consensus_state_comm_key = self.get_consensus_state_commitment_key(
            client_id,
            height.revision_number,
            height.revision_height,
        );
        let comm_mapper = self.commitments(&consensus_state_comm_key);
        if comm_mapper.is_empty() {
            let hashed_consensus_state = self.crypto().keccak256(encoded_consensus_state);
            comm_mapper.set(hashed_consensus_state);
        }
    }

    #[proxy]
    fn generic_client_proxy_impl(
        &self,
        sc_address: ManagedAddress,
    ) -> generic_client_proxy::GenericClientProxy<Self::Api>;
}
