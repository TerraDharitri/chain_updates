use common_types::{
    channel_types::{channel, height},
    ChannelId, ClientId, Hash, PortId, Sequence, UnixTimestamp,
};

use crate::{
    channel_libs::packet_types::{MsgTimeoutOnClose, MsgTimeoutPacket, Packet},
    interfaces::{client_interface, ibc_module_interface},
    packet_handling::errors::{PACKET_COMM_MISMATCH_ERR_MSG, UNEXPECTED_PACKET_DEST_ERR_MSG},
};

dharitri_sc::imports!();

#[dharitri_sc::module]
pub trait TimeoutModule:
    host::module_manager::ModuleManagerModule
    + host::storage::StorageModule
    + host::commitment::CommitmentModule
    + common_modules::utils::UtilsModule
    + super::membership::MembershipModule
    + crate::channel_libs::events::EventsModule
{
    #[endpoint(timeoutPacket)]
    fn timeout_packet(&self, args: MsgTimeoutPacket<Self::Api>) {
        let channel_info =
            self.try_get_channel_info(&args.packet.src_port, &args.packet.src_channel);
        let channel = &channel_info.channel;
        self.check_expected_args(&args.packet, channel);

        let connection_info = self.try_get_connection_info(&channel.connection_hops.get(0));
        let client_info = self.try_get_client_info(&connection_info.client_id);
        self.check_timeout_reached(
            &args.packet,
            args.proof_height,
            client_info.client_impl.clone(),
            connection_info.client_id.clone(),
        );

        let commitment_mapper = self.check_and_get_commitment_mapper(
            &args.packet.src_port,
            &args.packet.src_channel,
            args.packet.seq,
        );
        let commitment = commitment_mapper.get();
        let packet_commitment = self.get_packet_commitment(&args.packet);
        require!(
            commitment == packet_commitment,
            PACKET_COMM_MISMATCH_ERR_MSG
        );

        self.check_channel_membership(
            channel.ordering,
            client_info.client_impl,
            &connection_info,
            &args,
        );

        commitment_mapper.clear();

        self.timeout_packet_final(args.packet);
    }

    #[endpoint(timeoutOnClose)]
    fn timeout_on_close(&self, args: MsgTimeoutOnClose<Self::Api>) {
        let channel_info =
            self.try_get_channel_info(&args.packet.src_port, &args.packet.src_channel);
        let channel = &channel_info.channel;
        self.check_expected_args(&args.packet, channel);

        let connection_info = self.try_get_connection_info(&channel.connection_hops.get(0));
        let client_info = self.try_get_client_info(&connection_info.client_id);

        let commitment_mapper = self.check_and_get_commitment_mapper(
            &args.packet.src_port,
            &args.packet.src_channel,
            args.packet.seq,
        );
        let commitment = commitment_mapper.get();
        let packet_commitment = self.get_packet_commitment(&args.packet);
        require!(
            commitment == packet_commitment,
            PACKET_COMM_MISMATCH_ERR_MSG
        );

        self.check_expected_channel_membership(
            client_info.client_impl.clone(),
            channel,
            &connection_info,
            &args,
        );
        self.check_channel_membership(
            channel.ordering,
            client_info.client_impl,
            &connection_info,
            &args,
        );

        self.timeout_packet_final(args.packet);
    }

    fn check_expected_args(&self, packet: &Packet<Self::Api>, channel: &channel::Data<Self::Api>) {
        self.require_state_open(channel.state);
        require!(
            packet.dest_port == channel.counterparty.port_id,
            UNEXPECTED_PACKET_DEST_ERR_MSG
        );
        require!(
            packet.dest_channel == channel.counterparty.channel_id,
            UNEXPECTED_PACKET_DEST_ERR_MSG
        );
    }

    fn check_timeout_reached(
        &self,
        packet: &Packet<Self::Api>,
        proof_height: height::Data,
        client_impl: ManagedAddress,
        client_id: ClientId<Self::Api>,
    ) {
        if !packet.timeout_height.is_zero() && proof_height >= packet.timeout_height {
            return;
        }

        let timestamp_at_height: UnixTimestamp = self
            .generic_client_proxy_impl_timeout(client_impl)
            .get_timestamp_at_height(client_id, proof_height)
            .execute_on_dest_context();
        require!(
            packet.timeout_timestamp != 0 && timestamp_at_height >= packet.timeout_timestamp,
            "Channel timeout not reached"
        );
    }

    fn check_and_get_commitment_mapper(
        &self,
        src_port: &PortId<Self::Api>,
        src_channel: &ChannelId<Self::Api>,
        seq: Sequence,
    ) -> SingleValueMapper<Hash<Self::Api>> {
        let comm_hash = self.get_packet_commitment_key(src_port, src_channel, seq);
        let comm_mapper = self.commitments(&comm_hash);
        require!(!comm_mapper.is_empty(), "Packet commitment not found");

        comm_mapper
    }

    fn get_packet_commitment(&self, packet: &Packet<Self::Api>) -> Hash<Self::Api> {
        let hashed_data = self.crypto().sha256(&packet.data);

        let encoded_timeout_timestamp = self.encode_to_buffer(&packet.timeout_timestamp);
        let encoded_revision_number = self.encode_to_buffer(&packet.timeout_height.revision_number);
        let encoded_revision_height = self.encode_to_buffer(&packet.timeout_height.revision_height);
        let encoded_hashed_data = self.encode_to_buffer(&hashed_data);

        let mut encoded_buffer = ManagedBuffer::new();
        encoded_buffer = encoded_buffer.concat(encoded_timeout_timestamp);
        encoded_buffer = encoded_buffer.concat(encoded_revision_number);
        encoded_buffer = encoded_buffer.concat(encoded_revision_height);
        encoded_buffer = encoded_buffer.concat(encoded_hashed_data);

        let first_hashing = self.crypto().sha256(&encoded_buffer);

        self.crypto().keccak256(first_hashing.as_managed_buffer())
    }

    fn timeout_packet_final(&self, packet: Packet<Self::Api>) {
        let caller = self.blockchain().get_caller();
        let ibc_module = self.lookup_module_by_channel(&packet.src_port, &packet.src_channel);
        let _: () = self
            .ibc_module_proxy_impl_timeout(ibc_module)
            .on_timeout_packet(packet.clone(), caller)
            .execute_on_dest_context();

        self.timeout_packet_event(&packet);
    }

    #[proxy]
    fn generic_client_proxy_impl_timeout(
        &self,
        sc_address: ManagedAddress,
    ) -> client_interface::generic_client_proxy::GenericClientProxy<Self::Api>;

    #[proxy]
    fn ibc_module_proxy_impl_timeout(
        &self,
        sc_address: ManagedAddress,
    ) -> ibc_module_interface::ibc_module_proxy::IbcModuleProxy<Self::Api>;
}
