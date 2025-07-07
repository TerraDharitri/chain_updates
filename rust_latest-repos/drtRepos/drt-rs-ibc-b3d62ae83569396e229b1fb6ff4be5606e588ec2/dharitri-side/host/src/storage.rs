use common_types::{
    channel_types::{channel, upgrade},
    connection_types::connection_end,
    ChannelId, ClientId, ClientType, ConnectionId, Hash, PortId, Sequence, UnixTimestamp,
};

dharitri_sc::imports!();
dharitri_sc::derive_imports!();

#[derive(TypeAbi, TopEncode, TopDecode)]
pub struct ClientInfo<M: ManagedTypeApi> {
    pub client_type: ClientType<M>,
    pub client_impl: ManagedAddress<M>,
}

#[derive(TypeAbi, TopEncode, TopDecode, Default, PartialEq, Debug)]
pub struct HostInfo {
    pub next_client_seq: Sequence,
    pub next_connection_seq: Sequence,
    pub next_channel_seq: Sequence,
    pub expected_time_per_block: UnixTimestamp,
}

#[derive(TypeAbi, TopEncode, TopDecode, NestedEncode, NestedDecode, Clone, Copy)]
pub struct RecvStartSequence {
    pub seq: Sequence,
    pub prev_seq: Sequence,
}

#[derive(TypeAbi, TopEncode, TopDecode)]
pub struct ChannelInfo<M: ManagedTypeApi> {
    pub channel: channel::Data<M>,
    pub next_seq_send: Sequence,
    pub next_seq_recv: Sequence,
    pub next_seq_ack: Sequence,
    pub upgrade: upgrade::Data<M>,
    pub latest_error_rec_seq: Sequence,
    pub recv_start_seq: RecvStartSequence,
    pub ack_start_seq: Sequence,
}

#[dharitri_sc::module]
pub trait StorageModule {
    // Needs to be its own function for proxy
    #[view(getCommitment)]
    fn get_commitment(&self, commitment_hash: &Hash<Self::Api>) -> Hash<Self::Api> {
        self.commitments(commitment_hash).get()
    }

    /// calculates the block delay based on the expected time per block
    fn calculate_block_delay(&self, time_delay: UnixTimestamp) -> UnixTimestamp {
        if time_delay == 0 {
            return 0;
        }

        let host_info = self.host_info().get();
        if host_info.expected_time_per_block == 0 {
            return 0;
        }

        (time_delay + host_info.expected_time_per_block - 1) / host_info.expected_time_per_block
    }

    fn get_next_client_seq(&self) -> Sequence {
        self.host_info().update(|host_info| {
            let ret_val = host_info.next_client_seq;
            host_info.next_client_seq += 1;

            ret_val
        })
    }

    fn get_next_connection_seq(&self) -> Sequence {
        self.host_info().update(|host_info| {
            let ret_val = host_info.next_connection_seq;
            host_info.next_connection_seq += 1;

            ret_val
        })
    }

    fn get_next_channel_seq(&self) -> Sequence {
        self.host_info().update(|host_info| {
            let ret_val = host_info.next_channel_seq;
            host_info.next_channel_seq += 1;

            ret_val
        })
    }

    fn try_get_client_info(&self, client_id: &ClientId<Self::Api>) -> ClientInfo<Self::Api> {
        let mapper = self.client_info(client_id);
        require!(!mapper.is_empty(), "Client not found");

        mapper.get()
    }

    fn try_get_connection_info(
        &self,
        connection_id: &ConnectionId<Self::Api>,
    ) -> connection_end::Data<Self::Api> {
        let mapper = self.connection_info(connection_id);
        require!(!mapper.is_empty(), "Connection not found");

        mapper.get()
    }

    fn try_get_channel_info(
        &self,
        port_id: &PortId<Self::Api>,
        channel_id: &ChannelId<Self::Api>,
    ) -> ChannelInfo<Self::Api> {
        let mapper = self.channel_info(port_id, channel_id);
        require!(!mapper.is_empty(), "Channel not found");

        mapper.get()
    }

    #[storage_mapper("commitments")]
    fn commitments(&self, comm_key: &Hash<Self::Api>) -> SingleValueMapper<Hash<Self::Api>>;

    #[storage_mapper("clientReg")]
    fn client_registry(
        &self,
        client_type: &ClientType<Self::Api>,
    ) -> SingleValueMapper<ManagedAddress>;

    #[storage_mapper("clientInfo")]
    fn client_info(
        &self,
        client_id: &ClientId<Self::Api>,
    ) -> SingleValueMapper<ClientInfo<Self::Api>>;

    #[storage_mapper("portCap")]
    fn port_capabilities(&self, port_id: &PortId<Self::Api>) -> SingleValueMapper<ManagedAddress>;

    #[storage_mapper("channelCap")]
    fn channel_capabilities(
        &self,
        port_id: &PortId<Self::Api>,
        channel_id: &ChannelId<Self::Api>,
    ) -> SingleValueMapper<ManagedAddress>;

    #[storage_mapper("hostInfo")]
    fn host_info(&self) -> SingleValueMapper<HostInfo>;

    #[storage_mapper("connInfo")]
    fn connection_info(
        &self,
        connection_id: &ConnectionId<Self::Api>,
    ) -> SingleValueMapper<connection_end::Data<Self::Api>>;

    #[storage_mapper("channelInfo")]
    fn channel_info(
        &self,
        port_id: &PortId<Self::Api>,
        channel_id: &ChannelId<Self::Api>,
    ) -> SingleValueMapper<ChannelInfo<Self::Api>>;
}
