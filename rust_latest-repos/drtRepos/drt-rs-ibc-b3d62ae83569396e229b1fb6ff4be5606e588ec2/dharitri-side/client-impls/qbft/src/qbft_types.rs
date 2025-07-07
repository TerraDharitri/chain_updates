pub mod client_state {
    use common_types::{channel_types::height, ChainId};

    dharitri_sc::imports!();
    dharitri_sc::derive_imports!();

    #[derive(TypeAbi, TopEncode, TopDecode)]
    pub struct Data<M: ManagedTypeApi> {
        pub chain_id: ChainId<M>,
        pub ibc_store_address: ManagedAddress<M>,
        pub latest_height: height::Data,
        pub trusting_period: u64, // TODO: Likely timestamp
        pub max_clock_drift: u64, // TODO: No idea what this means, likely a timestamp too?
    }
}

pub mod consensus_state {
    use common_types::{Hash, UnixTimestamp};

    dharitri_sc::imports!();
    dharitri_sc::derive_imports!();

    #[derive(TypeAbi, TopEncode, TopDecode)]
    pub struct Data<M: ManagedTypeApi> {
        pub timestamp: UnixTimestamp,
        pub root: Hash<M>,
        pub validators: ManagedVec<M, ManagedAddress<M>>, // TODO: Why was this bytes[]? Also, might be able to use ArrayVec over ManagedVec
    }
}

pub mod header {
    use common_types::{channel_types::height, Hash};

    dharitri_sc::imports!();
    dharitri_sc::derive_imports!();

    #[derive(TypeAbi, TopEncode, TopDecode)]
    pub struct Data<M: ManagedTypeApi> {
        pub besu_header: ManagedBuffer<M>,
        pub seals: ManagedVec<M, ManagedBuffer<M>>,
        pub trusted_height: height::Data,
        pub account_state_proof: Hash<M>,
    }
}
