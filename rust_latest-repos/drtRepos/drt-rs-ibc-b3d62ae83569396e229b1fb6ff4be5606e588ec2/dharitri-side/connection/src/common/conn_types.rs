use common_types::{
    channel_types::height,
    connection_types::{counterparty, version},
    ClientId, ConnectionId, Hash, UnixTimestamp, VersionVec,
};

dharitri_sc::imports!();
dharitri_sc::derive_imports!();

#[derive(TypeAbi, TopDecode)]
pub struct MsgConnectionOpenInit<M: ManagedTypeApi> {
    pub client_id: ClientId<M>,
    pub counterparty: counterparty::Data<M>,
    pub version: version::Data<M>,
    pub delay_period: UnixTimestamp,
}

#[derive(TypeAbi, TopDecode)]
pub struct MsgConnectionOpenTry<M: ManagedTypeApi> {
    pub counterparty: counterparty::Data<M>, // counterpartyConnectionIdentifier, counterpartyPrefix and counterpartyClientIdentifier
    pub delay_period: UnixTimestamp,
    pub client_id: ClientId<M>,               // clientID of chainA
    pub client_state_bytes: ManagedBuffer<M>, // clientState that chainA has for chainB
    pub counterparty_versions: VersionVec<M>, // supported versions of chain A
    pub proof_init: Hash<M>, // proof that chainA stored connectionEnd in state (on ConnOpenInit)
    pub proof_client: Hash<M>, // proof that chainA stored a light client of chainB
    pub proof_consensus: Hash<M>, // proof that chainA stored chainB's consensus state at consensus height
    pub proof_height: height::Data, // height at which relayer constructs proof of A storing connectionEnd in state
    pub consensus_height: height::Data, // latest height of chain B which chain A has stored in its chain B client
    pub host_consensus_state_proof: Hash<M>, // optional proof data for host state machines that are unable to introspect their own consensus state
}

#[derive(TypeAbi, TopDecode, Clone)]
pub struct MsgConnectionOpenAck<M: ManagedTypeApi> {
    pub connection_id: ConnectionId<M>,
    pub client_state_bytes: ManagedBuffer<M>, // client state for chainA on chainB
    pub version: version::Data<M>,            // version that ChainB chose in ConnOpenTry
    pub counterparty_connection_id: ConnectionId<M>,
    pub proof_try: Hash<M>, // proof that connectionEnd was added to ChainB state in ConnOpenTry
    pub proof_client: Hash<M>, // proof of client state on chainB for chainA
    pub proof_consensus: Hash<M>, // proof that chainB has stored ConsensusState of chainA on its client
    pub proof_height: height::Data, // height that relayer constructed proofTry
    pub consensus_height: height::Data, // latest height of chainA that chainB has stored on its chainA client
    pub host_consensus_state_proof: Hash<M>, // optional proof data for host state machines that are unable to introspect their own consensus state
}

#[derive(TypeAbi, TopDecode)]
pub struct MsgConnectionOpenConfirm<M: ManagedTypeApi> {
    pub connection_id: ConnectionId<M>,
    pub proof_ack: Hash<M>,
    pub proof_height: height::Data,
}
