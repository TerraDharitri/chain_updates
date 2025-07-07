use client_common::{ClientStatus, GetLatestInfoResultType};
use common_types::{channel_types::height, ChannelId, PortId, Sequence, UnixTimestamp};

use crate::{channel_libs::events::SendPacketEventData, interfaces::client_interface};

dharitri_sc::imports!();

#[dharitri_sc::module]
pub trait SendModule:
    host::module_manager::ModuleManagerModule
    + host::storage::StorageModule
    + common_modules::utils::UtilsModule
    + host::commitment::CommitmentModule
    + crate::channel_libs::events::EventsModule
    + super::encoding::EncodingModule
{
    /// Is called by a module in order to send an IBC packet on a channel.
    ///
    /// The packet sequence generated for the packet to be sent is returned.
    ///
    /// An error is returned if one occurs.
    ///
    /// Also, `timeout_timestamp` is given in nanoseconds since unix epoch.
    #[endpoint(sendPacket)]
    fn send_packet(
        &self,
        src_port: PortId<Self::Api>,
        src_channel: ChannelId<Self::Api>,
        timeout_height: height::Data,
        timeout_timestamp: UnixTimestamp,
        data: ManagedBuffer,
    ) -> Sequence {
        let caller = self.blockchain().get_caller();
        self.authenticate_channel_capability(&src_port, &src_channel, &caller);

        let mut channel_info = self.try_get_channel_info(&src_port, &src_channel);
        let channel = &channel_info.channel;
        self.require_state_open(channel.state);
        require!(
            !timeout_height.is_zero() || timeout_timestamp != 0,
            "Zero packet timeout"
        );

        self.check_latest_info(
            &channel.connection_hops.get(0),
            timeout_height,
            timeout_timestamp,
        );

        let packet_seq = channel_info.next_seq_send;
        channel_info.next_seq_send += 1;
        self.channel_info(&src_port, &src_channel).set(channel_info);

        let commitment_hash = self.get_packet_commitment_key(&src_port, &src_channel, packet_seq);
        let encoded_data = self.encode_and_hash_twice(timeout_height, timeout_timestamp, &data);
        self.commitments(&commitment_hash).set(encoded_data);

        self.send_packet_event(SendPacketEventData {
            seq: packet_seq,
            source_port: &src_port,
            source_channel: &src_channel,
            timeout_height,
            timeout_timestamp,
            data: &data,
        });

        packet_seq
    }

    fn check_latest_info(
        &self,
        channel_id: &ChannelId<Self::Api>,
        timeout_height: height::Data,
        timeout_timestamp: UnixTimestamp,
    ) {
        let connection_info = self.try_get_connection_info(channel_id);
        let client_info = self.try_get_client_info(&connection_info.client_id);
        let latest_info: GetLatestInfoResultType = self
            .generic_client_proxy_impl_send(client_info.client_impl)
            .get_latest_info(connection_info.client_id)
            .execute_on_dest_context();

        require!(
            matches!(latest_info.client_status, ClientStatus::Active),
            "Client not active"
        );
        require!(
            timeout_height.is_zero() || latest_info.latest_height < timeout_height,
            "Past packet timeout height"
        );
        require!(
            timeout_timestamp == 0 || latest_info.latest_timestamp < timeout_timestamp,
            "Past packet timeout timestamp"
        );
    }

    #[proxy]
    fn generic_client_proxy_impl_send(
        &self,
        sc_address: ManagedAddress,
    ) -> client_interface::generic_client_proxy::GenericClientProxy<Self::Api>;
}
