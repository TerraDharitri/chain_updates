use common_types::{channel_types::height, ChannelId, PortId, Sequence, UnixTimestamp};

use super::packet_types::Packet;

dharitri_sc::imports!();
dharitri_sc::derive_imports!();

#[derive(TypeAbi, TopEncode)]
pub struct SendPacketEventData<'a, M: ManagedTypeApi> {
    pub seq: Sequence,
    pub source_port: &'a PortId<M>,
    pub source_channel: &'a ChannelId<M>,
    pub timeout_height: height::Data,
    pub timeout_timestamp: UnixTimestamp,
    pub data: &'a ManagedBuffer<M>,
}

#[dharitri_sc::module]
pub trait EventsModule {
    #[event("generatedChannelIdEvent")]
    fn generated_channel_id_event(&self, #[indexed] channel_id: &ChannelId<Self::Api>);

    #[event("writeAckEvent")]
    fn write_ack_event(
        &self,
        #[indexed] dest_port_id: &PortId<Self::Api>,
        #[indexed] dest_channel: &ChannelId<Self::Api>,
        #[indexed] seq: Sequence,
        ack: &ManagedBuffer,
    );

    #[event("sendPacketEvent")]
    fn send_packet_event(&self, data: SendPacketEventData<Self::Api>);

    #[event("receivePacketEvent")]
    fn receive_packet_event(&self, #[indexed] packet: &Packet<Self::Api>);

    #[event("ackPacketEvent")]
    fn ack_packet_event(&self, #[indexed] packet: &Packet<Self::Api>, ack: &ManagedBuffer);

    #[event("timeoutPacketEvent")]
    fn timeout_packet_event(&self, #[indexed] packet: &Packet<Self::Api>);
}
