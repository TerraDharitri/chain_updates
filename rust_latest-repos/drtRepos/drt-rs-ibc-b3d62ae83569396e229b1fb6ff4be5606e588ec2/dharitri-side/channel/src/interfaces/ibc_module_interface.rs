use common_types::{
    channel_types::{channel, channel_counterparty},
    ChannelId, ConnectionHops, PortId, Version,
};

dharitri_sc::imports!();
dharitri_sc::derive_imports!();

#[derive(TypeAbi, TopEncode, TopDecode)]
pub struct MsgOnChanOpenInit<M: ManagedTypeApi> {
    pub order: channel::Order,
    pub connection_hops: ConnectionHops<M>,
    pub port_id: PortId<M>,
    pub channel_id: ChannelId<M>,
    pub counterparty: channel_counterparty::Data<M>,
    pub version: Version<M>,
}

#[derive(TypeAbi, TopEncode, TopDecode)]
pub struct MsgOnChanOpenTry<M: ManagedTypeApi> {
    pub order: channel::Order,
    pub connection_hops: ConnectionHops<M>,
    pub port_id: PortId<M>,
    pub channel_id: ChannelId<M>,
    pub counterparty: channel_counterparty::Data<M>,
    pub counterparty_version: Version<M>,
}

#[derive(TypeAbi, TopEncode, TopDecode)]
pub struct MsgOnChanOpenAck<M: ManagedTypeApi> {
    pub port_id: PortId<M>,
    pub channel_id: ChannelId<M>,
    pub counterparty_version: Version<M>,
}

#[derive(TypeAbi, TopEncode, TopDecode)]
pub struct MsgOnChanOpenConfirm<M: ManagedTypeApi> {
    pub port_id: PortId<M>,
    pub channel_id: ChannelId<M>,
}

pub type MsgOnChanCloseInit<M> = MsgOnChanOpenConfirm<M>;
pub type MsgOnChanCloseConfirm<M> = MsgOnChanOpenConfirm<M>;

pub mod ibc_module_proxy {
    use common_types::Version;

    use crate::channel_libs::packet_types::Packet;

    use super::{
        MsgOnChanCloseConfirm, MsgOnChanCloseInit, MsgOnChanOpenAck, MsgOnChanOpenConfirm,
        MsgOnChanOpenInit, MsgOnChanOpenTry,
    };

    dharitri_sc::imports!();

    #[dharitri_sc::proxy]
    pub trait IbcModuleProxy {
        /// Will verify that the relayer-chosen parameters are valid and perform any custom INIT logic.
        ///
        /// It may return an error if the chosen parameters are invalid, in which case the handshake is aborted.
        ///
        /// If the provided version string is non-empty, OnChanOpenInit should return the version string if valid,
        /// or an error if the provided version is invalid.
        ///
        /// If the version string is empty, OnChanOpenInit is expected to return a default version string representing the version(s) it supports.
        ///
        /// If there is no default version string for the application, it should return an error if provided version is empty string.
        #[endpoint(onChanOpenInit)]
        fn on_chan_open_init(&self, args: MsgOnChanOpenInit<Self::Api>) -> Version<Self::Api>;

        /// Will verify the relayer-chosen parameters along with the counterparty-chosen version string and perform custom TRY logic.
        ///
        /// If the relayer-chosen parameters are invalid, the callback must return an error to abort the handshake.
        ///
        /// If the counterparty-chosen version is not compatible with this modules supported versions,
        /// the callback must return an error to abort the handshake.
        ///
        /// If the versions are compatible, the try callback must select the final version string and return it to core IBC.
        ///
        /// OnChanOpenTry may also perform custom initialization logic
        #[endpoint(onChanOpenTry)]
        fn on_chan_open_try(&self, args: MsgOnChanOpenTry<Self::Api>) -> Version<Self::Api>;

        /// Will error if the counterparty selected version string is invalid to abort the handshake.
        ///
        /// It may also perform custom ACK logic.
        #[endpoint(onChanOpenAck)]
        fn on_chan_open_ack(&self, args: MsgOnChanOpenAck<Self::Api>);

        /// Will perform custom CONFIRM logic and may error to abort the handshake.
        #[endpoint(onChanOpenConfirm)]
        fn on_chan_open_confirm(&self, args: MsgOnChanOpenConfirm<Self::Api>);

        /// Will perform custom CLOSE_INIT logic and may error to abort the handshake.
        ///
        /// If the application does not allow the channel to be closed, this function must revert.
        #[endpoint(onChanCloseInit)]
        fn on_chan_close_init(&self, args: MsgOnChanCloseInit<Self::Api>);

        /// Will perform custom CLOSE_CONFIRM logic and may error to abort the handshake.
        #[endpoint(onChanCloseConfirm)]
        fn on_chan_close_confirm(&self, args: MsgOnChanCloseConfirm<Self::Api>);

        /// Must return an acknowledgement that implements the Acknowledgement interface.
        ///
        /// In the case of an asynchronous acknowledgement, nil should be returned.
        ///
        /// If the acknowledgement returned is successful, the state changes on callback are written,
        /// otherwise the application state changes are discarded.
        ///
        /// In either case the packet is received and the acknowledgement is written (in synchronous cases).
        #[endpoint(onRecvPacket)]
        fn on_recv_packet(
            &self,
            packet: Packet<Self::Api>,
            relayer: ManagedAddress,
        ) -> ManagedBuffer;

        /// Is called when a packet sent by this module has been acknowledged.
        #[endpoint(onAcknowledgementPacket)]
        fn on_ack_packet(
            &self,
            packet: Packet<Self::Api>,
            ack: ManagedBuffer,
            relayer: ManagedAddress,
        );

        /// Is called when a packet sent by this module has timed-out (such that it will not be received on the destination chain).
        #[endpoint(onTimeoutPacket)]
        fn on_timeout_packet(&self, packet: Packet<Self::Api>, relayer: ManagedAddress);
    }
}
