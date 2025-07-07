pub mod channel {
    use crate::{channel_types::channel_counterparty, ConnectionHops, Sequence, Version};

    dharitri_sc::imports!();
    dharitri_sc::derive_imports!();

    pub static ORDERED: &[u8] = b"ORDER_ORDERED";
    pub static UNORDERED: &[u8] = b"ORDER_UNORDERED";

    #[derive(TypeAbi, TopEncode, TopDecode, NestedEncode, NestedDecode, Clone, Copy)]
    pub enum State {
        UninitializedUnspecified,
        Init,
        TryOpen,
        Open,
        Closed,
        Flushing,
        FlushComplete,
    }

    #[derive(TypeAbi, TopEncode, TopDecode, NestedEncode, NestedDecode, Clone, Copy)]
    pub enum Order {
        NoneUnspecified,
        Unordered,
        Ordered,
    }

    impl Order {
        pub fn to_byte_slice<M: ManagedTypeApi>(&self) -> &[u8] {
            match *self {
                Order::NoneUnspecified => {
                    M::error_api_impl().signal_error(b"Unknown channel order")
                }
                Order::Unordered => UNORDERED,
                Order::Ordered => ORDERED,
            }
        }
    }

    #[derive(TypeAbi, TopEncode, TopDecode, NestedEncode, NestedDecode)]
    pub struct Data<M: ManagedTypeApi> {
        pub state: State,
        pub ordering: Order,
        pub counterparty: channel_counterparty::Data<M>,
        pub connection_hops: ConnectionHops<M>,
        pub version: Version<M>,
        pub upgrade_sequence: Sequence,
    }
}

pub mod channel_counterparty {
    dharitri_sc::imports!();
    dharitri_sc::derive_imports!();

    use crate::{ChannelId, PortId};

    #[derive(TypeAbi, TopEncode, TopDecode, NestedEncode, NestedDecode)]
    pub struct Data<M: ManagedTypeApi> {
        pub port_id: PortId<M>,
        pub channel_id: ChannelId<M>,
    }
}

pub mod height {
    dharitri_sc::imports!();
    dharitri_sc::derive_imports!();

    #[derive(
        TypeAbi,
        TopEncode,
        TopDecode,
        NestedEncode,
        NestedDecode,
        ManagedVecItem,
        PartialEq,
        PartialOrd,
        Clone,
        Copy,
    )]
    pub struct Data {
        pub revision_number: u64,
        pub revision_height: u64,
    }

    impl Data {
        pub fn is_zero(&self) -> bool {
            self.revision_number == 0 && self.revision_height == 0
        }

        pub fn to_biguint_concat<M: ManagedTypeApi>(&self) -> BigUint<M> {
            let mut buffer = ManagedBuffer::new();
            let _ = self.revision_number.dep_encode(&mut buffer);
            let _ = self.revision_height.dep_encode(&mut buffer);

            BigUint::from_bytes_be_buffer(&buffer)
        }
    }

    #[cfg(test)]
    mod tests {
        use core::cmp::Ordering;

        use super::*;

        impl Data {
            fn new(revision_number: u64, revision_height: u64) -> Self {
                Self {
                    revision_number,
                    revision_height,
                }
            }
        }

        #[test]
        fn partial_ord_test() {
            assert_eq!(
                Data::new(0, 100).partial_cmp(&Data::new(1, 50)),
                Some(Ordering::Less)
            );
            assert_eq!(
                Data::new(0, 100).partial_cmp(&Data::new(0, 100)),
                Some(Ordering::Equal)
            );
            assert_eq!(
                Data::new(0, 100).partial_cmp(&Data::new(0, 50)),
                Some(Ordering::Greater)
            );
            assert_eq!(
                Data::new(0, 50).partial_cmp(&Data::new(1, 50)),
                Some(Ordering::Less)
            );
        }
    }
}

pub mod timeout {
    use crate::channel_types::height;
    use crate::UnixTimestamp;

    dharitri_sc::imports!();
    dharitri_sc::derive_imports!();

    #[derive(TypeAbi, TopEncode, TopDecode, NestedEncode, NestedDecode)]
    pub struct Data {
        pub height: height::Data,
        pub timestamp: UnixTimestamp,
    }
}

pub mod upgrade {
    use crate::Sequence;

    use super::{timeout, upgrade_fields};

    dharitri_sc::imports!();
    dharitri_sc::derive_imports!();

    #[derive(TypeAbi, TopEncode, TopDecode, NestedEncode, NestedDecode)]
    pub struct Data<M: ManagedTypeApi> {
        pub fields: upgrade_fields::Data<M>,
        pub timeout: timeout::Data,
        pub next_sequence_send: Sequence,
    }
}

pub mod upgrade_fields {
    use crate::{ConnectionHops, Version};

    use super::channel;

    dharitri_sc::imports!();
    dharitri_sc::derive_imports!();

    #[derive(TypeAbi, TopEncode, TopDecode, NestedEncode, NestedDecode)]
    pub struct Data<M: ManagedTypeApi> {
        pub ordering: channel::Order,
        pub connection_hops: ConnectionHops<M>,
        pub version: Version<M>,
    }
}

pub mod error_receipt {
    use crate::Sequence;

    dharitri_sc::imports!();
    dharitri_sc::derive_imports!();

    #[derive(TypeAbi, TopEncode, TopDecode)]
    pub struct Data<M: ManagedTypeApi> {
        pub sequence: Sequence,
        pub message: ManagedBuffer<M>,
    }
}
