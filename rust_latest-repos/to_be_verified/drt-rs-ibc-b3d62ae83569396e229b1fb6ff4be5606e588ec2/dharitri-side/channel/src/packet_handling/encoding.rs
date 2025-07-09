use common_types::{channel_types::height, Hash, UnixTimestamp};

dharitri_sc::imports!();

#[dharitri_sc::module]
pub trait EncodingModule: common_modules::utils::UtilsModule {
    fn encode_and_hash(
        &self,
        timeout_height: height::Data,
        timeout_timestamp: UnixTimestamp,
        data: &ManagedBuffer,
    ) -> Hash<Self::Api> {
        let hashed_data = self.crypto().sha256(data);

        let mut encoded_buffer = ManagedBuffer::new();
        let encoded_timestamp = self.encode_to_buffer(&timeout_timestamp);
        let encoded_rev_number = self.encode_to_buffer(&timeout_height.revision_number);
        let encoded_rev_height = self.encode_to_buffer(&timeout_height.revision_height);
        let encoded_hashed_data = self.encode_to_buffer(&hashed_data);

        encoded_buffer = encoded_buffer.concat(encoded_timestamp);
        encoded_buffer = encoded_buffer.concat(encoded_rev_number);
        encoded_buffer = encoded_buffer.concat(encoded_rev_height);
        encoded_buffer = encoded_buffer.concat(encoded_hashed_data);

        self.crypto().sha256(&encoded_buffer)
    }

    fn encode_and_hash_twice(
        &self,
        timeout_height: height::Data,
        timeout_timestamp: UnixTimestamp,
        data: &ManagedBuffer,
    ) -> Hash<Self::Api> {
        let encoded_data = self.encode_and_hash(timeout_height, timeout_timestamp, data);

        self.crypto().keccak256(encoded_data.as_managed_buffer())
    }
}
