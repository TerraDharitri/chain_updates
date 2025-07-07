pub mod merkle_prefix {
    dharitri_sc::imports!();
    dharitri_sc::derive_imports!();

    #[derive(TypeAbi, TopEncode, TopDecode, NestedEncode, NestedDecode, Clone)]
    pub struct Data<M: ManagedTypeApi> {
        pub key_prefix: ManagedBuffer<M>,
    }
}

pub mod connection_end {
    use crate::{ClientId, UnixTimestamp, VersionVec};

    use super::counterparty;

    dharitri_sc::imports!();
    dharitri_sc::derive_imports!();

    #[derive(TypeAbi, TopEncode, TopDecode, NestedEncode, NestedDecode, Clone, Copy)]
    pub enum State {
        UninitializedUnspecified,
        Init,
        TryOpen,
        Open,
    }

    #[derive(TypeAbi, TopEncode, TopDecode, Clone)]
    pub struct Data<M: ManagedTypeApi> {
        pub client_id: ClientId<M>,
        pub versions: VersionVec<M>,
        pub state: State,
        pub counterparty: counterparty::Data<M>,
        pub delay_period: UnixTimestamp,
    }
}

pub mod counterparty {
    use crate::{ClientId, ConnectionId};

    use super::merkle_prefix;

    dharitri_sc::imports!();
    dharitri_sc::derive_imports!();

    #[derive(TypeAbi, TopEncode, TopDecode, NestedEncode, NestedDecode, Clone)]
    pub struct Data<M: ManagedTypeApi> {
        pub client_id: ClientId<M>,
        pub connection_id: ConnectionId<M>,
        pub prefix: merkle_prefix::Data<M>,
    }
}

pub mod version {
    use crate::{FeatureId, FeatureVec};

    dharitri_sc::imports!();
    dharitri_sc::derive_imports!();

    #[derive(TypeAbi, TopEncode, TopDecode, NestedEncode, NestedDecode, ManagedVecItem, Clone)]
    pub struct Data<M: ManagedTypeApi> {
        pub identifier: FeatureId<M>,
        pub features: FeatureVec<M>,
    }
}
