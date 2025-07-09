use common_types::{connection_types::connection_end, ConnectionId, VersionVec};

use crate::common::conn_types::{
    MsgConnectionOpenAck, MsgConnectionOpenConfirm, MsgConnectionOpenInit, MsgConnectionOpenTry,
};

dharitri_sc::imports!();

static CONNECTION_ALREADY_EXISTS_ERR_MSG: &[u8] = b"Connection already exists";
static CONNECTION_DOES_NOT_EXIST_ERR_MSG: &[u8] = b"Connection does not exist";
static INVALID_CONNECTION_STATE_ERR_MSG: &[u8] = b"Invalid connection state";

#[dharitri_sc::module]
pub trait ConnectionEndpointsModule:
    crate::common::conn_lib::ConnectionLibModule
    + crate::common::verify_states::VerifyStatesModule
    + crate::common::conn_internal::ConnectionInternalModule
    + crate::common::events::EventsModule
    + host::commitment::CommitmentModule
    + host::host_config::HostConfigModule
    + host::host_views::HostViewsModule
    + host::module_manager::ModuleManagerModule
    + host::storage::StorageModule
    + common_modules::client_lib::ClientLibModule
    + common_modules::host_lib::HostLibModule
    + common_modules::utils::UtilsModule
{
    // TODO: Check if those endpoints need special permissions

    /// Initialises a connection attempt on chain A.
    ///
    /// The generated connection identifier is returned.
    #[endpoint(connectionOpenInit)]
    fn connection_open_init(
        &self,
        args: MsgConnectionOpenInit<Self::Api>,
    ) -> ConnectionId<Self::Api> {
        let connection_id = self.generate_connection_id();
        let connection_mapper = self.connection_info(&connection_id);
        require!(
            connection_mapper.is_empty(),
            CONNECTION_ALREADY_EXISTS_ERR_MSG
        );

        // ensure the client exists
        let _ = self.check_and_get_client(&args.client_id);
        require!(
            !args.counterparty.connection_id.is_empty(),
            "Invalid counterparty connection ID"
        );

        let mut connection_info = connection_end::Data {
            client_id: args.client_id,
            counterparty: args.counterparty,
            delay_period: args.delay_period,
            state: connection_end::State::Init,
            versions: VersionVec::new(),
        };

        self.set_versions_after_init(args.version, &mut connection_info.versions);
        self.update_connection_commitment(&connection_id, &connection_info);
        connection_mapper.set(connection_info);

        self.generated_connection_id_event(&connection_id);

        connection_id
    }

    /// relays notice of a connection attempt on chain A to chain B (this code is executed on chain B)
    #[endpoint(connectionOpenTry)]
    fn connection_open_try(
        &self,
        args: MsgConnectionOpenTry<Self::Api>,
    ) -> ConnectionId<Self::Api> {
        require!(
            !args.counterparty_versions.is_empty(),
            "Empty counterparty versions"
        );

        let self_consensus_state = args.host_consensus_state_proof.clone();
        let connection_id = self.generate_connection_id();
        let connection_mapper = self.connection_info(&connection_id);
        require!(
            connection_mapper.is_empty(),
            CONNECTION_ALREADY_EXISTS_ERR_MSG
        );

        // ensure the client exists
        let _ = self.check_and_get_client(&args.client_id);

        let compatible_versions = self.get_compatible_versions();
        let picked_version = self.pick_version(&compatible_versions, &args.counterparty_versions);
        let connection_info = connection_end::Data {
            client_id: args.client_id.clone(),
            counterparty: args.counterparty.clone(),
            delay_period: args.delay_period,
            state: connection_end::State::TryOpen,
            versions: ManagedVec::from_single_item(picked_version),
        };
        connection_mapper.set(&connection_info);

        self.verify_all_states_open_try(connection_info.clone(), &self_consensus_state, args);
        self.update_connection_commitment(&connection_id, &connection_info);
        self.generated_connection_id_event(&connection_id);

        connection_id
    }

    /// relays acceptance of a connection open attempt from chain B back to chain A (this code is executed on chain A)
    #[endpoint(connectionOpenAck)]
    fn connection_open_ack(&self, args: MsgConnectionOpenAck<Self::Api>) {
        let connection_mapper = self.connection_info(&args.connection_id);
        require!(
            !connection_mapper.is_empty(),
            CONNECTION_DOES_NOT_EXIST_ERR_MSG
        );

        let mut connection_info = connection_mapper.get();
        require!(
            matches!(connection_info.state, connection_end::State::Init),
            INVALID_CONNECTION_STATE_ERR_MSG
        );

        let self_consensus_state = args.host_consensus_state_proof.clone();
        self.verify_all_states_open_ack(
            connection_info.clone(),
            &self_consensus_state,
            args.clone(),
        );

        connection_info.state = connection_end::State::Open;
        connection_info.counterparty.connection_id = args.counterparty_connection_id;
        connection_info.versions = ManagedVec::from_single_item(args.version);

        connection_mapper.set(&connection_info);
        self.update_connection_commitment(&args.connection_id, &connection_info);
    }

    /// confirms opening of a connection on chain A to chain B, after which the connection is open on both chains (this code is executed on chain B)
    #[endpoint(connectionOpenConfirm)]
    fn connection_open_confirm(&self, args: MsgConnectionOpenConfirm<Self::Api>) {
        let connection_mapper = self.connection_info(&args.connection_id);
        require!(
            !connection_mapper.is_empty(),
            CONNECTION_DOES_NOT_EXIST_ERR_MSG
        );

        let mut connection_info = connection_mapper.get();
        require!(
            matches!(connection_info.state, connection_end::State::TryOpen),
            INVALID_CONNECTION_STATE_ERR_MSG
        );

        let connection_id = args.connection_id.clone();
        self.verify_all_states_open_confirm(connection_info.clone(), args);

        connection_info.state = connection_end::State::Open;
        connection_mapper.set(&connection_info);
        self.update_connection_commitment(&connection_id, &connection_info);
    }
}
