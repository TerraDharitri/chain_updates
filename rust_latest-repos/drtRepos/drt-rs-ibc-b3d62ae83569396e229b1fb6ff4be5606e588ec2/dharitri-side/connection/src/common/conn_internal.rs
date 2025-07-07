use client_common::VerifyMembershipArgs;
use common_types::{
    channel_types::height,
    connection_types::{connection_end, version},
    ConnectionId, Hash, Path, VersionVec,
};

dharitri_sc::imports!();

pub struct VerifyClientStateArgs<M: ManagedTypeApi> {
    pub connection_info: connection_end::Data<M>,
    pub height: height::Data,
    pub path: Path<M>,
    pub proof: Hash<M>,
    pub client_state_bytes: ManagedBuffer<M>,
}

pub struct VerifyConsensusStateArgs<M: ManagedTypeApi> {
    pub connection_info: connection_end::Data<M>,
    pub height: height::Data,
    pub consensus_height: height::Data,
    pub proof: Hash<M>,
    pub consensus_state_bytes: ManagedBuffer<M>,
}

pub struct VerifyConnectionStateArgs<M: ManagedTypeApi> {
    pub connection_info: connection_end::Data<M>,
    pub height: height::Data,
    pub proof: Hash<M>,
    pub counterparty_connection_id: ConnectionId<M>,
    pub counterparty_connection_info: connection_end::Data<M>,
}

mod client_proxy {
    use client_common::{VerifyMembershipArgs, VerifyNonMembershipArgs};

    dharitri_sc::imports!();

    #[dharitri_sc::proxy]
    pub trait ClientProxy {
        #[view(verifyMembership)]
        fn verify_membership(&self, args: VerifyMembershipArgs<Self::Api>) -> bool;

        #[view(verifyNonMembership)]
        fn verify_non_membership(&self, args: VerifyNonMembershipArgs<Self::Api>) -> bool;
    }
}

#[dharitri_sc::module]
pub trait ConnectionInternalModule:
    crate::common::conn_lib::ConnectionLibModule
    + host::commitment::CommitmentModule
    + host::host_config::HostConfigModule
    + host::host_views::HostViewsModule
    + host::module_manager::ModuleManagerModule
    + host::storage::StorageModule
    + common_modules::client_lib::ClientLibModule
    + common_modules::host_lib::HostLibModule
    + common_modules::utils::UtilsModule
{
    #[view(getCompatibleVersions)]
    fn get_compatible_versions(&self) -> VersionVec<Self::Api> {
        VersionVec::from_single_item(self.default_ibc_version())
    }

    fn set_versions_after_init(
        &self,
        args_version: version::Data<Self::Api>,
        output: &mut VersionVec<Self::Api>,
    ) {
        let compatible_versions = self.get_compatible_versions();
        if !args_version.features.is_empty() {
            require!(
                self.is_supported_version(&compatible_versions, &args_version),
                "Version not supported"
            );

            output.push(args_version);
        } else {
            self.set_supported_versions(compatible_versions, output);
        }
    }

    fn update_connection_commitment(
        &self,
        connection_id: &ConnectionId<Self::Api>,
        connection_info: &connection_end::Data<Self::Api>,
    ) {
        let connection_key = self.get_connection_commitment_key(connection_id);
        let encoded_connection = self.encode_to_buffer(&connection_info);
        let hashed_connection = self.crypto().keccak256(encoded_connection);

        self.commitments(&connection_key).set(&hashed_connection);
    }

    fn generate_connection_id(&self) -> ConnectionId<Self::Api> {
        let next_conn_seq = self.get_next_connection_seq();

        sc_format!("connection-{}", next_conn_seq)
    }

    fn verify_client_state(&self, args: VerifyClientStateArgs<Self::Api>) {
        let client = self.check_and_get_client(&args.connection_info.client_id);
        let args = VerifyMembershipArgs {
            client_id: args.connection_info.client_id,
            height: args.height,
            delay_time_period: 0,
            delay_block_period: 0,
            proof: args.proof,
            prefix: args.connection_info.counterparty.prefix.key_prefix,
            path: args.path,
            value: args.client_state_bytes,
        };
        let membership_result: bool = self
            .client_proxy_impl(client)
            .verify_membership(args)
            .execute_on_dest_context();
        require!(membership_result, "Failed to verify client state");
    }

    fn verify_consensus_state(&self, args: VerifyConsensusStateArgs<Self::Api>) {
        let client = self.check_and_get_client(&args.connection_info.client_id);
        let consensus_state_path = self.get_consensus_state_path(
            &args.connection_info.counterparty.client_id,
            args.consensus_height.revision_number,
            args.consensus_height.revision_height,
        );
        let args = VerifyMembershipArgs {
            client_id: args.connection_info.client_id,
            height: args.height,
            delay_time_period: 0,
            delay_block_period: 0,
            proof: args.proof,
            prefix: args.connection_info.counterparty.prefix.key_prefix,
            path: consensus_state_path,
            value: args.consensus_state_bytes,
        };
        let membership_result: bool = self
            .client_proxy_impl(client)
            .verify_membership(args)
            .execute_on_dest_context();
        require!(membership_result, "Failed to verify consensus state");
    }

    fn verify_connection_state(&self, args: VerifyConnectionStateArgs<Self::Api>) {
        let client = self.check_and_get_client(&args.connection_info.client_id);
        let connection_path = self.get_connection_path(&args.counterparty_connection_id);
        let encoded_connection = self.encode_to_buffer(&args.counterparty_connection_info);

        let args = VerifyMembershipArgs {
            client_id: args.connection_info.client_id,
            height: args.height,
            delay_time_period: 0,
            delay_block_period: 0,
            proof: args.proof,
            prefix: args.connection_info.counterparty.prefix.key_prefix,
            path: connection_path,
            value: encoded_connection,
        };
        let membership_result: bool = self
            .client_proxy_impl(client)
            .verify_membership(args)
            .execute_on_dest_context();
        require!(membership_result, "Failed to verify connection state");
    }

    #[proxy]
    fn client_proxy_impl(&self, sc_address: ManagedAddress)
        -> client_proxy::ClientProxy<Self::Api>;
}
