use crate::qbft_types::{client_state, consensus_state, header};
use common_types::{channel_types::height, ClientId, Hash, UnixTimestamp};

dharitri_sc::imports!();
dharitri_sc::derive_imports!();

// Likely don't need all this

// static HEADER_TYPE_URL: &[u8] = b"/ibc.lightclients.qbft.v1.Header";
// static CLIENT_STATE_TYPE_URL: &[u8] = b"/ibc.lightclients.qbft.v1.ClientState";
// static CONSENSUS_STATE_TYPE_URL: &[u8] = b"/ibc.lightclients.qbft.v1.ConsensusState";

pub struct ParsedBesuHeader<M: ManagedTypeApi> {
    pub base: header::Data<M>,
    pub height: height::Data,
    pub state_root: Hash<M>,
    pub time: UnixTimestamp,
    pub validators: ManagedVec<M, ManagedAddress<M>>, // TODO: Was RLPReader.RLPItem[]. Why?
}

#[dharitri_sc::module]
pub trait ClientLogicModule:
    client_common::CommonClientLogicModule
    + host::host_views::HostViewsModule
    + host::storage::StorageModule
    + common_modules::utils::UtilsModule
{
    #[endpoint(initializeClient)]
    fn initialize_client(
        &self,
        client_id: ClientId<Self::Api>,
        client_state: client_state::Data<Self::Api>,
        consensus_state: consensus_state::Data<Self::Api>,
    ) -> height::Data {
        self.require_ibc_handler_caller();
        require!(
            client_state.latest_height.revision_height != 0,
            "Invalid client state height"
        );
        require!(!consensus_state.validators.is_empty(), "Empty validators");

        self.client_states(&client_id).set(&client_state);
        self.consensus_states(&client_id, &client_state.latest_height.to_biguint_concat())
            .set(consensus_state);

        client_state.latest_height
    }

    /// updates the client with the given header
    #[endpoint(updateClient)]
    fn update_client(
        &self,
        client_id: ClientId<Self::Api>,
        _header: header::Data<Self::Api>,
    ) -> ManagedVec<height::Data> {
        let client_state_mapper = self.client_states(&client_id);
        require!(client_state_mapper.is_empty(), "Unknown client");

        // TODO

        ManagedVec::new()
    }

    /*
    function parseBesuHeader(Header.Data memory header) internal pure returns (ParsedBesuHeader memory) {
        ParsedBesuHeader memory parsedHeader;

        parsedHeader.base = header;
        RLPReader.RLPItem[] memory items = header.besu_header_rlp.toRlpItem().toList();
        if (items.length < 15) {
            revert UnexpectedEthereumHeaderFormat(items.length);
        }
        parsedHeader.stateRoot = bytes32(items[3].toUint());
        parsedHeader.height = Height.Data({revision_number: 0, revision_height: uint64(items[8].toUint())});
        parsedHeader.time = uint64(items[11].toUint());
        items = items[12].toBytes().toRlpItem().toList();
        // IBFT2: {Vanity, Validators, Vote, Round}
        // QBFT:  {Vanity, Validators, Vote, Round, Empty-Seals}
        if (items.length != 4 && items.length != 5) {
            revert UnexpectedExtraDataFormat(items.length);
        }
        parsedHeader.validators = items[1].toList();
        return parsedHeader;
    }
     */

    // fn parse_besu_header(&self, header: header::Data<Self::Api>) -> ParsedBesuHeader<Self::Api> {
    //     let items = header.besu_header

    //     let parsed_header = ParsedBesuHeader {
    //         base: header,
    //         height: todo!(),
    //         state_root: todo!(),
    //         time: todo!(),
    //         validators: todo!(),
    //     };
    // }

    // TODO: Replace the generic "ManagedBuffer" and "BigUint" with something specific once I figure out what it means
    // Was this:
    /*
        mapping(string => ClientState.Data) internal clientStates;
        mapping(string => mapping(uint128 => ConsensusState.Data)) internal consensusStates;
        mapping(string => mapping(uint128 => uint256)) internal processedTimes;
        mapping(string => mapping(uint128 => uint256)) internal processedHeights;
    */

    #[storage_mapper("clientStates")]
    fn client_states(
        &self,
        client_id: &ClientId<Self::Api>,
    ) -> SingleValueMapper<client_state::Data<Self::Api>>;

    #[storage_mapper("consensusStates")]
    fn consensus_states(
        &self,
        client_id: &ClientId<Self::Api>,
        height: &BigUint,
    ) -> SingleValueMapper<consensus_state::Data<Self::Api>>;

    #[storage_mapper("processedTimes")]
    fn processed_times(
        &self,
        buffer: &ManagedBuffer,
        biguint: &BigUint,
    ) -> SingleValueMapper<BigUint>;

    #[storage_mapper("processedHeights")]
    fn processed_heights(
        &self,
        buffer: &ManagedBuffer,
        biguint: &BigUint,
    ) -> SingleValueMapper<BigUint>;
}
