use common_types::{channel_types::channel, Hash, UnixTimestamp, HASH_LENGTH};

dharitri_sc::imports!();

const NANO_SECONDS_MULT: u64 = 1_000_000_000;
const EMPTY_HASH: &[u8; HASH_LENGTH] = &[0u8; HASH_LENGTH];

pub static UNEXPECTED_CHANNEL_STATE_ERR_MSG: &[u8] = b"Unexpected channel state";

#[dharitri_sc::module]
pub trait UtilsModule {
    fn require_valid_address(&self, address: &ManagedAddress) {
        let own_sc_address = self.blockchain().get_sc_address();
        require!(
            address != &own_sc_address && !address.is_zero(),
            "Invalid address"
        );
    }

    fn require_state_open(&self, state: channel::State) {
        require!(
            matches!(state, channel::State::Open),
            UNEXPECTED_CHANNEL_STATE_ERR_MSG
        );
    }

    fn checked_timestamp_to_unix_mul(&self, timestamp: u64) -> UnixTimestamp {
        match timestamp.checked_mul(NANO_SECONDS_MULT) {
            Some(result) => result,
            None => sc_panic!("Overlow!!!"),
        }
    }

    fn encode_to_buffer<T: TopEncode>(&self, value: &T) -> ManagedBuffer {
        let mut encoded_value = ManagedBuffer::new();
        let _ = value.top_encode(&mut encoded_value);

        encoded_value
    }

    fn is_empty_hash(&self, hash: &Hash<Self::Api>) -> bool {
        hash == &Hash::new_from_bytes(EMPTY_HASH)
    }
}
