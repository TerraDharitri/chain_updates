use client_common::VerifyMembershipArgs;
use common_modules::utils::UNEXPECTED_CHANNEL_STATE_ERR_MSG;
use common_types::{
    channel_types::{channel, height},
    connection_types::connection_end,
    Hash, Path, Sequence,
};
use host::storage::ChannelInfo;

use crate::{
    channel_libs::packet_types::{MsgPacketRecv, Packet, PacketReceipt},
    interfaces::{client_interface, ibc_module_interface},
    packet_handling::errors::{PACKET_ALREADY_PROCESSED_ERR_MSG, UNEXPECTED_PACKET_SOURCE_ERR_MSG},
};

dharitri_sc::imports!();

pub struct VerifyPacketCommitmentReceiveArgs<'a, M: ManagedTypeApi> {
    pub connection_info: &'a connection_end::Data<M>,
    pub height: height::Data,
    pub proof: Hash<M>,
    pub path: Path<M>,
    pub commitment: Hash<M>,
}

#[dharitri_sc::module]
pub trait ReceiveModule:
    host::module_manager::ModuleManagerModule
    + host::storage::StorageModule
    + common_modules::utils::UtilsModule
    + host::commitment::CommitmentModule
    + host::host_views::HostViewsModule
    + crate::channel_libs::events::EventsModule
    + super::ack::AckModule
    + super::encoding::EncodingModule
{
    /// Is called by a module in order to receive & process an IBC packet sent on the corresponding channel end on the counterparty chain.
    #[endpoint(recvPacket)]
    fn recieve_packet(&self, args: MsgPacketRecv<Self::Api>) {
        let mut channel_info =
            self.try_get_channel_info(&args.packet.dest_port, &args.packet.dest_channel);
        let channel = &channel_info.channel;
        match channel.state {
            channel::State::Open => {}
            channel::State::Flushing | channel::State::FlushComplete => {
                self.flush_channel(&channel_info, args.packet.seq)
            }
            _ => sc_panic!(UNEXPECTED_CHANNEL_STATE_ERR_MSG),
        }

        self.verify_valid_packet(&args.packet, channel);

        let connection_info = self.try_get_connection_info(&channel.connection_hops.get(0));
        self.verify_packet_commitment_receive(VerifyPacketCommitmentReceiveArgs {
            connection_info: &connection_info,
            height: args.proof_height,
            proof: args.proof,
            path: self.get_packet_commitment_path(
                &args.packet.src_port,
                &args.packet.src_channel,
                args.packet.seq,
            ),
            commitment: self.encode_and_hash(
                args.packet.timeout_height,
                args.packet.timeout_timestamp,
                &args.packet.data,
            ),
        });

        self.receive_packet_by_channel_order(&mut channel_info, &args.packet);

        let module =
            self.lookup_module_by_channel(&args.packet.dest_port, &args.packet.dest_channel);
        let caller = self.blockchain().get_caller();
        let ack: ManagedBuffer = self
            .ibc_module_proxy_impl_receive(module)
            .on_recv_packet(args.packet.clone(), caller)
            .execute_on_dest_context();

        if !ack.is_empty() {
            self.write_ack(
                &args.packet.dest_port,
                &args.packet.dest_channel,
                args.packet.seq,
                &ack,
            );
        }

        self.receive_packet_event(&args.packet);
    }

    fn flush_channel(&self, channel_info: &ChannelInfo<Self::Api>, packet_seq: Sequence) {
        let rseq = channel_info.recv_start_seq;

        // prevSequence=0 means the channel is not in the process of being upgraded or counterparty has not been upgraded yet
        if rseq.prev_seq == 0 {
            return;
        }

        require!(
            packet_seq < rseq.seq,
            "Channel cannot receive next upgrade packet"
        );
        require!(
            packet_seq >= rseq.prev_seq,
            PACKET_ALREADY_PROCESSED_ERR_MSG
        );
    }

    fn verify_valid_packet(&self, packet: &Packet<Self::Api>, channel: &channel::Data<Self::Api>) {
        require!(
            packet.src_port == channel.counterparty.port_id,
            UNEXPECTED_PACKET_SOURCE_ERR_MSG
        );
        require!(
            packet.src_channel == channel.counterparty.channel_id,
            UNEXPECTED_PACKET_SOURCE_ERR_MSG
        );

        let block_number = self.blockchain().get_block_nonce();
        let revision_height = packet.timeout_height.revision_height;
        require!(
            revision_height == 0 || block_number < revision_height,
            "Packet timeout height"
        );

        let host_timestamp = self.get_host_timestamp();
        let timeout_timestamp = packet.timeout_timestamp;
        require!(
            timeout_timestamp == 0 || host_timestamp < timeout_timestamp,
            "Packet timeout timestamp"
        );
    }

    fn verify_packet_commitment_receive(&self, args: VerifyPacketCommitmentReceiveArgs<Self::Api>) {
        let client = self.check_and_get_client(&args.connection_info.client_id);
        let membership_args = VerifyMembershipArgs {
            client_id: args.connection_info.client_id.clone(),
            height: args.height,
            delay_time_period: args.connection_info.delay_period,
            delay_block_period: self.calculate_block_delay(args.connection_info.delay_period),
            proof: args.proof,
            prefix: args.connection_info.counterparty.prefix.key_prefix.clone(),
            path: args.path,
            value: args.commitment.as_managed_buffer().clone(),
        };
        let membership_result: bool = self
            .generic_client_proxy_impl_receive(client)
            .verify_membership(membership_args)
            .execute_on_dest_context();
        require!(
            membership_result,
            "Failed to verify received packet commitment"
        );
    }

    fn receive_packet_by_channel_order(
        &self,
        channel_info: &mut ChannelInfo<Self::Api>,
        packet: &Packet<Self::Api>,
    ) {
        match channel_info.channel.ordering {
            channel::Order::Ordered => self.receive_packet_ordered(channel_info, packet),
            channel::Order::Unordered => self.receive_packet_unordered(channel_info, packet),
            channel::Order::NoneUnspecified => sc_panic!("Unknown order"),
        }
    }

    fn receive_packet_ordered(
        &self,
        channel_info: &mut ChannelInfo<Self::Api>,
        packet: &Packet<Self::Api>,
    ) {
        require!(
            channel_info.next_seq_recv == packet.seq,
            "Unexpected sequence receive"
        );

        channel_info.next_seq_recv += 1;

        let commitment_key =
            self.get_next_seq_recv_commitment_key(&packet.dest_port, &packet.dest_channel);
        let encoded_val = self.encode_to_buffer(&channel_info.next_seq_recv);
        let hashed_val = self.crypto().keccak256(&encoded_val);
        self.commitments(&commitment_key).set(hashed_val);

        self.channel_info(&packet.dest_port, &packet.dest_channel)
            .set(channel_info);
    }

    fn receive_packet_unordered(
        &self,
        channel_info: &ChannelInfo<Self::Api>,
        packet: &Packet<Self::Api>,
    ) {
        if matches!(channel_info.channel.state, channel::State::Open) {
            let rseq = channel_info.recv_start_seq;
            require!(packet.seq >= rseq.seq, PACKET_ALREADY_PROCESSED_ERR_MSG);
        }

        let commitment_key = self.get_packet_receipt_commitment_key(
            &packet.dest_port,
            &packet.dest_channel,
            packet.seq,
        );
        let comm_mapper = self.commitments(&commitment_key);
        require!(
            comm_mapper.is_empty(),
            "Channel packet receipt already exists"
        );

        let encoded_success = self.encode_to_buffer(&PacketReceipt::Successful);
        let successful_hash = self.crypto().keccak256(&encoded_success);
        comm_mapper.set(successful_hash);
    }

    #[proxy]
    fn generic_client_proxy_impl_receive(
        &self,
        sc_address: ManagedAddress,
    ) -> client_interface::generic_client_proxy::GenericClientProxy<Self::Api>;

    #[proxy]
    fn ibc_module_proxy_impl_receive(
        &self,
        sc_address: ManagedAddress,
    ) -> ibc_module_interface::ibc_module_proxy::IbcModuleProxy<Self::Api>;
}
