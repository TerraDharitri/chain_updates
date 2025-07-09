use common_types::{channel_types::height, ChannelId, Hash, PortId, Sequence, UnixTimestamp};

dharitri_sc::imports!();
dharitri_sc::derive_imports!();

#[derive(TopEncode)]
pub enum PacketReceipt {
    None,
    Successful,
}

/// `Packet` defines a type that carries data across different chains through IBC.
///
/// `seq` corresponds to the order of sends and receives, where a packet with an earlier sequence number must be sent and received before a packet with a later sequence number
///
/// `src_port` identifies the port on the sending chain
///
/// `src_channel` identifies the channel end on the sending chain
///
/// `dest_port` identifies the port on the receiving chain
///
/// `dest_channel` identifies the channel end on the receiving chain
///
/// `data` is an opaque value which can be defined by the application logic of the associated modules
///
/// `timeout_height` indicates a consensus height on the destination chain after which the packet will no longer be processed, and will instead count as having timed-out
///
/// `timeout_timestamp` indicates a timestamp on the destination chain after which the packet will no longer be processed, and will instead count as having timed-out
#[derive(TypeAbi, TopEncode, TopDecode, NestedEncode, NestedDecode, Clone)]
pub struct Packet<M: ManagedTypeApi> {
    pub seq: Sequence,
    pub src_port: PortId<M>,
    pub src_channel: ChannelId<M>,
    pub dest_port: PortId<M>,
    pub dest_channel: ChannelId<M>,
    pub data: ManagedBuffer<M>,
    pub timeout_height: height::Data,
    pub timeout_timestamp: UnixTimestamp,
}

#[derive(TypeAbi, TopEncode, TopDecode)]
pub struct MsgPacketRecv<M: ManagedTypeApi> {
    pub packet: Packet<M>,
    pub proof: Hash<M>,
    pub proof_height: height::Data,
}

#[derive(TypeAbi, TopEncode, TopDecode)]
pub struct MsgPacketAcknowledgement<M: ManagedTypeApi> {
    pub packet: Packet<M>,
    pub ack: ManagedBuffer<M>, // TODO: Or is it Hash<M>?
    pub proof: Hash<M>,
    pub proof_height: height::Data,
}

#[derive(TypeAbi, TopEncode, TopDecode)]
pub struct MsgTimeoutPacket<M: ManagedTypeApi> {
    pub packet: Packet<M>,
    pub proof: Hash<M>,
    pub proof_height: height::Data,
    pub next_seq_recv: Sequence,
}

#[derive(TypeAbi, TopEncode, TopDecode)]
pub struct MsgTimeoutOnClose<M: ManagedTypeApi> {
    pub packet: Packet<M>,
    pub proof_unreceived: Hash<M>,
    pub proof_close: Hash<M>,
    pub proof_height: height::Data,
    pub next_seq_recv: Sequence,
    pub counterparty_upgrade_seq: Sequence,
}

pub trait TimeoutArgs<M: ManagedTypeApi> {
    fn get_packet(&self) -> &Packet<M>;

    fn get_proof(&self) -> &Hash<M>;

    fn get_proof_height(&self) -> height::Data;

    fn get_next_seq_recv(&self) -> Sequence;
}

impl<M: ManagedTypeApi> TimeoutArgs<M> for MsgTimeoutPacket<M> {
    #[inline(always)]
    fn get_packet(&self) -> &Packet<M> {
        &self.packet
    }

    #[inline(always)]
    fn get_proof(&self) -> &Hash<M> {
        &self.proof
    }

    #[inline(always)]
    fn get_proof_height(&self) -> height::Data {
        self.proof_height
    }

    #[inline(always)]
    fn get_next_seq_recv(&self) -> Sequence {
        self.next_seq_recv
    }
}

impl<M: ManagedTypeApi> TimeoutArgs<M> for MsgTimeoutOnClose<M> {
    #[inline(always)]
    fn get_packet(&self) -> &Packet<M> {
        &self.packet
    }

    #[inline(always)]
    fn get_proof(&self) -> &Hash<M> {
        &self.proof_unreceived
    }

    #[inline(always)]
    fn get_proof_height(&self) -> height::Data {
        self.proof_height
    }

    #[inline(always)]
    fn get_next_seq_recv(&self) -> Sequence {
        self.next_seq_recv
    }
}
