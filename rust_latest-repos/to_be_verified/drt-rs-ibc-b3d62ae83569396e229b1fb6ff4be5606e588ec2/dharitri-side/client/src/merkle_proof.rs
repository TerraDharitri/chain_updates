use common_types::Hash;

dharitri_sc::imports!();

#[dharitri_sc::module]
pub trait MerkleProofModule {
    fn verify_merkle_proof(
        &self,
        item_to_prove: &ManagedBuffer,
        proof: &ManagedVec<Hash<Self::Api>>,
        root_hash: &Hash<Self::Api>,
    ) -> bool {
        let caller = self.blockchain().get_caller();
        let mut leaf_bytes = caller.as_managed_buffer().clone();
        leaf_bytes.append(item_to_prove);

        let mut hash = self.crypto().sha256(&leaf_bytes);
        for proof_item in proof {
            let hash_buffer = hash.as_managed_buffer();
            let proof_buffer = proof_item.as_managed_buffer();

            if BigUint::from(hash_buffer) < BigUint::from(proof_buffer) {
                hash = self.append_and_hash(hash_buffer, proof_buffer);
            } else {
                hash = self.append_and_hash(proof_buffer, hash_buffer);
            }
        }

        &hash == root_hash
    }

    fn append_and_hash(
        &self,
        first_buffer: &ManagedBuffer,
        second_buffer: &ManagedBuffer,
    ) -> Hash<Self::Api> {
        let mut hash_input = first_buffer.clone();
        hash_input.append(second_buffer);

        self.crypto().sha256(hash_input)
    }
}
