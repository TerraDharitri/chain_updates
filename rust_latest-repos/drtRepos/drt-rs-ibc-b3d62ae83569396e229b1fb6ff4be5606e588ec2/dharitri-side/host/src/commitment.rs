use common_types::{ChannelId, ClientId, ConnectionId, Hash, Path, PortId};

dharitri_sc::imports!();

// https://github.com/hyperledger-labs/yui-ibc-solidity/blob/main/contracts/core/24-host/IBCCommitment.sol
// https://github.com/cosmos/ibc/tree/main/spec/core/ics-024-host-requirements#path-space

#[dharitri_sc::module]
pub trait CommitmentModule {
    /// "clients/{identifier}/clientState"
    fn get_client_state_path(&self, client_id: &ClientId<Self::Api>) -> Path<Self::Api> {
        sc_format!("clients/{}/clientState", client_id)
    }

    /// "clients/{identifier}/consensusStates/{revision_number}-{height}"
    fn get_consensus_state_path(
        &self,
        client_id: &ClientId<Self::Api>,
        revision_number: u64,
        revision_height: u64,
    ) -> Path<Self::Api> {
        sc_format!(
            "clients/{}/consensusStates/{}-{}",
            client_id,
            revision_number,
            revision_height
        )
    }

    /// "connections/{identifier}"
    fn get_connection_path(&self, connection_id: &ConnectionId<Self::Api>) -> Path<Self::Api> {
        sc_format!("connections/{}", connection_id)
    }

    /// "channelEnds/ports/{identifier}/channels/{identifier}"
    fn get_channel_path(
        &self,
        port_id: &PortId<Self::Api>,
        channel_id: &ChannelId<Self::Api>,
    ) -> Path<Self::Api> {
        sc_format!("channelEnds/ports/{}/channels/{}", port_id, channel_id)
    }

    /// "commitments/ports/{identifier}/channels/{identifier}/sequences/{sequence}"
    fn get_packet_commitment_path(
        &self,
        port_id: &PortId<Self::Api>,
        channel_id: &ChannelId<Self::Api>,
        sequence: u64,
    ) -> Path<Self::Api> {
        sc_format!(
            "commitments/ports/{}/channels/{}/sequences/{}",
            port_id,
            channel_id,
            sequence
        )
    }

    /// "acks/ports/{identifier}/channels/{identifier}/sequences/{sequence}"
    fn get_packet_acknowledgement_commitment_path(
        &self,
        port_id: &PortId<Self::Api>,
        channel_id: &ChannelId<Self::Api>,
        sequence: u64,
    ) -> Path<Self::Api> {
        sc_format!(
            "acks/ports/{}/channels/{}/sequences/{}",
            port_id,
            channel_id,
            sequence
        )
    }

    /// "receipts/ports/{identifier}/channels/{identifier}/sequences/{sequence}"
    fn get_packet_receipt_commitment_path(
        &self,
        port_id: &PortId<Self::Api>,
        channel_id: &ChannelId<Self::Api>,
        sequence: u64,
    ) -> Path<Self::Api> {
        sc_format!(
            "receipts/ports/{}/channels/{}/sequences/{}",
            port_id,
            channel_id,
            sequence
        )
    }

    /// "nextSequenceSend/ports/{identifier}/channels/{identifier}"
    fn get_next_seq_send_commitment_path(
        &self,
        port_id: &PortId<Self::Api>,
        channel_id: &ChannelId<Self::Api>,
    ) -> Path<Self::Api> {
        sc_format!("nextSequenceSend/ports/{}/channels/{}", port_id, channel_id)
    }

    /// "nextSequenceRecv/ports/{identifier}/channels/{identifier}"
    fn get_next_seq_recv_commitment_path(
        &self,
        port_id: &PortId<Self::Api>,
        channel_id: &ChannelId<Self::Api>,
    ) -> Path<Self::Api> {
        sc_format!("nextSequenceRecv/ports/{}/channels/{}", port_id, channel_id)
    }

    /// "nextSequenceAck/ports/{identifier}/channels/{identifier}"
    fn get_next_seq_ack_commitment_path(
        &self,
        port_id: &PortId<Self::Api>,
        channel_id: &ChannelId<Self::Api>,
    ) -> Path<Self::Api> {
        sc_format!("nextSequenceAck/ports/{}/channels/{}", port_id, channel_id)
    }

    /// "channelUpgrades/upgrades/ports/{identifier}/channels/{identifier}"
    fn get_channel_upgrade_path(
        &self,
        port_id: &PortId<Self::Api>,
        channel_id: &ChannelId<Self::Api>,
    ) -> Path<Self::Api> {
        sc_format!(
            "channelUpgrades/upgrades/ports/{}/channels/{}",
            port_id,
            channel_id
        )
    }

    /// "channelUpgrades/upgradeError/ports/{identifier}/channels/{identifier}"
    fn get_channel_upgrade_error_path(
        &self,
        port_id: &PortId<Self::Api>,
        channel_id: &ChannelId<Self::Api>,
    ) -> Path<Self::Api> {
        sc_format!(
            "channelUpgrades/upgradeError/ports/{}/channels/{}",
            port_id,
            channel_id
        )
    }

    // key gen

    fn get_client_state_commitment_key(&self, client_id: &ClientId<Self::Api>) -> Hash<Self::Api> {
        self.crypto()
            .keccak256(self.get_client_state_path(client_id))
    }

    fn get_consensus_state_commitment_key(
        &self,
        client_id: &ClientId<Self::Api>,
        revision_number: u64,
        revision_height: u64,
    ) -> Hash<Self::Api> {
        self.crypto().keccak256(self.get_consensus_state_path(
            client_id,
            revision_number,
            revision_height,
        ))
    }

    fn get_connection_commitment_key(
        &self,
        connection_id: &ConnectionId<Self::Api>,
    ) -> Hash<Self::Api> {
        self.crypto()
            .keccak256(self.get_connection_path(connection_id))
    }

    fn get_channel_commitment_key(
        &self,
        port_id: &PortId<Self::Api>,
        channel_id: &ChannelId<Self::Api>,
    ) -> Hash<Self::Api> {
        self.crypto()
            .keccak256(self.get_channel_path(port_id, channel_id))
    }

    fn get_next_seq_recv_commitment_key(
        &self,
        port_id: &PortId<Self::Api>,
        channel_id: &ChannelId<Self::Api>,
    ) -> Hash<Self::Api> {
        self.crypto()
            .keccak256(self.get_next_seq_recv_commitment_path(port_id, channel_id))
    }

    fn get_packet_commitment_key(
        &self,
        port_id: &PortId<Self::Api>,
        channel_id: &ChannelId<Self::Api>,
        sequence: u64,
    ) -> Hash<Self::Api> {
        self.crypto()
            .keccak256(self.get_packet_commitment_path(port_id, channel_id, sequence))
    }

    fn get_packet_acknowledgement_commitment_key(
        &self,
        port_id: &PortId<Self::Api>,
        channel_id: &ChannelId<Self::Api>,
        sequence: u64,
    ) -> Hash<Self::Api> {
        self.crypto().keccak256(
            self.get_packet_acknowledgement_commitment_path(port_id, channel_id, sequence),
        )
    }

    fn get_packet_receipt_commitment_key(
        &self,
        port_id: &PortId<Self::Api>,
        channel_id: &ChannelId<Self::Api>,
        sequence: u64,
    ) -> Hash<Self::Api> {
        self.crypto()
            .keccak256(self.get_packet_receipt_commitment_path(port_id, channel_id, sequence))
    }

    fn get_channel_upgrade_commitment_key(
        &self,
        port_id: &PortId<Self::Api>,
        channel_id: &ChannelId<Self::Api>,
    ) -> Hash<Self::Api> {
        self.crypto()
            .keccak256(self.get_channel_upgrade_path(port_id, channel_id))
    }

    fn get_channel_upgrade_error_commitment_key(
        &self,
        port_id: &PortId<Self::Api>,
        channel_id: &ChannelId<Self::Api>,
    ) -> Hash<Self::Api> {
        self.crypto()
            .keccak256(self.get_channel_upgrade_error_path(port_id, channel_id))
    }
}
