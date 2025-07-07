pub mod client_state {
    use common_types::channel_types::height;

    dharitri_sc::imports!();
    dharitri_sc::derive_imports!();

    #[derive(TypeAbi, TopEncode, TopDecode)]
    pub struct Data {
        pub latest_height: height::Data,
    }

    impl Data {
        #[inline]
        pub fn new(latest_height: height::Data) -> Self {
            Self { latest_height }
        }
    }
}

pub mod consensus_state {
    use common_types::UnixTimestamp;

    dharitri_sc::imports!();
    dharitri_sc::derive_imports!();

    #[derive(TypeAbi, TopEncode, TopDecode)]
    pub struct Data {
        pub timestamp: UnixTimestamp,
    }

    impl Data {
        #[inline]
        pub fn new(timestamp: UnixTimestamp) -> Self {
            Self { timestamp }
        }
    }
}

pub mod header {
    use common_types::{channel_types::height, UnixTimestamp};

    dharitri_sc::imports!();
    dharitri_sc::derive_imports!();

    #[derive(TypeAbi, TopEncode, TopDecode)]
    pub struct Data {
        pub height: height::Data,
        pub timestamp: UnixTimestamp,
    }
}
