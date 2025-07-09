use client_common::VerifyMembershipArgs;
use common_modules::utils::UNEXPECTED_CHANNEL_STATE_ERR_MSG;
use common_types::{
    channel_types::{channel, height},
    connection_types::connection_end,
    ChannelId, Hash, Path, PortId, Sequence,
};
use host::storage::ChannelInfo;

use crate::{
    channel_libs::packet_types::{MsgPacketAcknowledgement, Packet},
    interfaces::{client_interface, ibc_module_interface},
    packet_handling::errors::UNEXPECTED_PACKET_DEST_ERR_MSG,
};

use super::errors::UNKNOW_CHANNEL_ORDER_ERR_MSG;

dharitri_sc::imports!();

pub struct VerifyPacketAckArgs<'a, M: ManagedTypeApi> {
    pub connection_info: &'a connection_end::Data<M>,
    pub height: height::Data,
    pub proof: Hash<M>,
    pub path: Path<M>,
    pub ack_comm: Hash<M>,
}

#[dharitri_sc::module]
pub trait AckModule:
    host::module_manager::ModuleManagerModule
    + host::storage::StorageModule
    + common_modules::utils::UtilsModule
    + host::commitment::CommitmentModule
    + host::host_views::HostViewsModule
    + crate::channel_libs::events::EventsModule
    + super::encoding::EncodingModule
{
    /// Writes the packet execution acknowledgement to the state, which will be verified by the counterparty chain using AcknowledgePacket
    #[endpoint(writeAcknowledgement)]
    fn write_ack_endpoint(
        &self,
        dest_port: PortId<Self::Api>,
        dest_channel: ChannelId<Self::Api>,
        seq: Sequence,
        ack: ManagedBuffer,
    ) {
        let caller = self.blockchain().get_caller();
        self.authenticate_channel_capability(&dest_port, &dest_channel, &caller);

        let channel_info = self.try_get_channel_info(&dest_port, &dest_channel);
        let channel = &channel_info.channel;
        self.require_state_open(channel.state);
        require!(!ack.is_empty(), "Empty ack");

        self.write_ack(&dest_port, &dest_channel, seq, &ack);
    }

    /// Is called by a module to process the acknowledgement of a packet previously sent by the calling module on a channel to a counterparty
    /// module on the counterparty chain.
    ///
    /// Its intended usage is within the ante handler.
    ///
    /// AcknowledgePacket will clean up the packet commitment, which is no longer necessary since the packet has been received and acted upon.
    ///
    /// It will also increment NextSequenceAck in case of ORDERED channels.
    #[endpoint(acknowledgePacket)]
    fn ack_packet(&self, args: MsgPacketAcknowledgement<Self::Api>) {
        let mut channel_info =
            self.try_get_channel_info(&args.packet.src_port, &args.packet.src_channel);
        let channel = &channel_info.channel;
        self.check_expected_ack_pack_data(channel, &args.packet);

        let comm_mapper = self.verify_packet_commitment_ack(&args.packet);
        let connection_info = self.try_get_connection_info(&channel.connection_hops.get(0));
        self.verify_packet_ack(VerifyPacketAckArgs {
            connection_info: &connection_info,
            height: args.proof_height,
            proof: args.proof,
            path: self.get_packet_acknowledgement_commitment_path(
                &args.packet.dest_port,
                &args.packet.dest_channel,
                args.packet.seq,
            ),
            ack_comm: self.crypto().sha256(&args.ack),
        });
        self.check_packet_seq_by_order(&mut channel_info, &args.packet);

        comm_mapper.clear();

        let caller = self.blockchain().get_caller();
        let module = self.lookup_module_by_channel(&args.packet.src_port, &args.packet.src_channel);
        let _: () = self
            .ibc_module_proxy_impl_ack(module)
            .on_ack_packet(args.packet.clone(), args.ack.clone(), caller)
            .execute_on_dest_context();

        self.ack_packet_event(&args.packet, &args.ack);
    }

    fn write_ack(
        &self,
        dest_port: &PortId<Self::Api>,
        dest_channel: &ChannelId<Self::Api>,
        seq: Sequence,
        ack: &ManagedBuffer,
    ) {
        let ack_comm_key =
            self.get_packet_acknowledgement_commitment_key(dest_port, dest_channel, seq);
        let comm_mapper = self.commitments(&ack_comm_key);
        require!(comm_mapper.is_empty(), "Ack already written");

        let first_hash = self.crypto().sha256(ack);
        let second_hash = self.crypto().keccak256(first_hash.as_managed_buffer());
        comm_mapper.set(second_hash);

        self.write_ack_event(dest_port, dest_channel, seq, ack);
    }

    fn verify_packet_ack(&self, args: VerifyPacketAckArgs<Self::Api>) {
        let client = self.check_and_get_client(&args.connection_info.client_id);
        let membership_args = VerifyMembershipArgs {
            client_id: args.connection_info.client_id.clone(),
            height: args.height,
            delay_time_period: args.connection_info.delay_period,
            delay_block_period: self.calculate_block_delay(args.connection_info.delay_period),
            proof: args.proof,
            prefix: args.connection_info.counterparty.prefix.key_prefix.clone(),
            path: args.path,
            value: args.ack_comm.as_managed_buffer().clone(),
        };
        let membership_result: bool = self
            .generic_client_proxy_impl_ack(client)
            .verify_membership(membership_args)
            .execute_on_dest_context();
        require!(membership_result, "Failed to verify received packet ack");
    }

    fn check_expected_ack_pack_data(
        &self,
        channel: &channel::Data<Self::Api>,
        packet: &Packet<Self::Api>,
    ) {
        if !matches!(channel.state, channel::State::Open) {
            require!(
                matches!(channel.state, channel::State::Flushing),
                UNEXPECTED_CHANNEL_STATE_ERR_MSG
            );
        }

        require!(
            packet.dest_port == channel.counterparty.port_id,
            UNEXPECTED_PACKET_DEST_ERR_MSG
        );
        require!(
            packet.dest_channel == channel.counterparty.channel_id,
            UNEXPECTED_PACKET_DEST_ERR_MSG
        );
    }

    fn verify_packet_commitment_ack(
        &self,
        packet: &Packet<Self::Api>,
    ) -> SingleValueMapper<Hash<Self::Api>> {
        let comm_key =
            self.get_packet_commitment_key(&packet.src_port, &packet.src_channel, packet.seq);
        let comm_mapper = self.commitments(&comm_key);
        require!(!comm_mapper.is_empty(), "Commitment not found");

        let packet_comm = comm_mapper.get();
        let expected_comm = self.encode_and_hash_twice(
            packet.timeout_height,
            packet.timeout_timestamp,
            &packet.data,
        );
        require!(packet_comm == expected_comm, "Packet commitment mismatch");

        return comm_mapper;
    }

    fn check_packet_seq_by_order(
        &self,
        channel_info: &mut ChannelInfo<Self::Api>,
        packet: &Packet<Self::Api>,
    ) {
        match channel_info.channel.ordering {
            channel::Order::Ordered => {
                require!(
                    packet.seq == channel_info.next_seq_ack,
                    "Unexpected next seq ack"
                );

                channel_info.next_seq_ack += 1;

                self.channel_info(&packet.src_port, &packet.src_channel)
                    .set(channel_info);
            }
            channel::Order::Unordered => require!(
                packet.seq >= channel_info.ack_start_seq,
                "Ack already processed in prev upgrade"
            ),
            channel::Order::NoneUnspecified => sc_panic!(UNKNOW_CHANNEL_ORDER_ERR_MSG),
        }
    }

    #[proxy]
    fn ibc_module_proxy_impl_ack(
        &self,
        sc_address: ManagedAddress,
    ) -> ibc_module_interface::ibc_module_proxy::IbcModuleProxy<Self::Api>;

    #[proxy]
    fn generic_client_proxy_impl_ack(
        &self,
        sc_address: ManagedAddress,
    ) -> client_interface::generic_client_proxy::GenericClientProxy<Self::Api>;
}
