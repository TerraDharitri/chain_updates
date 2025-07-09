use common_types::{
    channel_types::{channel, height},
    ChannelId, Hash, PortId, Version,
};

dharitri_sc::imports!();
dharitri_sc::derive_imports!();

#[derive(TypeAbi, TopEncode, TopDecode)]
pub struct MsgChannelOpenInit<M: ManagedTypeApi> {
    pub port_id: PortId<M>,
    pub channel: channel::Data<M>,
}

#[derive(TypeAbi, TopEncode, TopDecode)]
pub struct MsgChannelOpenTry<M: ManagedTypeApi> {
    pub port_id: PortId<M>,
    pub channel: channel::Data<M>,
    pub counterparty_version: Version<M>,
    pub proof_init: Hash<M>,
    pub proof_height: height::Data,
}

#[derive(TypeAbi, TopEncode, TopDecode)]
pub struct MsgChannelOpenAck<M: ManagedTypeApi> {
    pub port_id: PortId<M>,
    pub channel_id: ChannelId<M>,
    pub counterparty_version: Version<M>,
    pub counterparty_channel_id: ChannelId<M>,
    pub proof_try: Hash<M>,
    pub proof_height: height::Data,
}

#[derive(TypeAbi, TopEncode, TopDecode)]
pub struct MsgChannelOpenConfirm<M: ManagedTypeApi> {
    pub port_id: PortId<M>,
    pub channel_id: ChannelId<M>,
    pub proof_ack: Hash<M>,
    pub proof_height: height::Data,
}

#[derive(TypeAbi, TopEncode, TopDecode)]
pub struct MsgChannelCloseInit<M: ManagedTypeApi> {
    pub port_id: PortId<M>,
    pub channel_id: ChannelId<M>,
}

#[derive(TypeAbi, TopEncode, TopDecode)]
pub struct MsgChannelCloseConfirm<M: ManagedTypeApi> {
    pub port_id: PortId<M>,
    pub channel_id: ChannelId<M>,
    pub proof_init: Hash<M>,
    pub proof_height: height::Data,
}
