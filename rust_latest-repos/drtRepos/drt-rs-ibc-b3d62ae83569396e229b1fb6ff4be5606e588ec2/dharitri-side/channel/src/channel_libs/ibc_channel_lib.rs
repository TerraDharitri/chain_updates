use common_types::Hash;

use super::packet_types::PacketReceipt;

dharitri_sc::imports!();
dharitri_sc::derive_imports!();

#[dharitri_sc::module]
pub trait IbcChannelLibModule: common_modules::utils::UtilsModule {
    fn receipt_commitment_to_receipt(&self, commitment: &Hash<Self::Api>) -> PacketReceipt {
        if self.is_empty_hash(commitment) {
            return PacketReceipt::None;
        }

        let encoded_success = self.encode_to_buffer(&PacketReceipt::Successful);
        let successful_hash = self.crypto().keccak256(&encoded_success);
        if commitment == &successful_hash {
            return PacketReceipt::Successful;
        }

        sc_panic!("Unknown channel packet receipt commitment");
    }
}
