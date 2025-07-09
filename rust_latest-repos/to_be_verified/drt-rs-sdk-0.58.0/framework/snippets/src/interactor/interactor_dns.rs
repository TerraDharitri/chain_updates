use dharitri_sc_scenario::{
    dharitri_sc::types::Address, scenario_format::value_interpreter::keccak256,
};

#[cfg(test)]
use dharitri_sc_scenario::bech32;

fn get_initial_dns_address() -> Address {
    Address::from_slice(&[1u8; 32])
}

fn compute_smart_contract_address(owner_address: Address, owner_nonce: u64) -> Address {
    // 8 bytes of zero + 2 bytes for VM type + 20 bytes of hash(owner) + 2 bytes of shard(owner)
    let owner_bytes = owner_address.as_bytes();
    let nonce_bytes = owner_nonce.to_le_bytes();
    let bytes_to_hash = [owner_bytes, &nonce_bytes].concat();
    let initial_padding = [0u8; 8];
    let vm_type: [u8; 2] = [5, 0];
    let address = keccak256(&bytes_to_hash);
    let address = [
        initial_padding.as_slice(),
        vm_type.as_slice(),
        &address[10..30],
        &owner_bytes[30..],
    ]
    .concat();
    Address::from_slice(&address)
}

fn compute_dns_address_for_shard_id(shard_id: u8) -> Address {
    let initial_dns_address = get_initial_dns_address();
    let initial_dns_address_slice = initial_dns_address.as_array();
    let shard_identifier = &[0u8, shard_id];
    let deployer_pubkey_prefix =
        &initial_dns_address_slice[0..initial_dns_address_slice.len() - shard_identifier.len()];

    let deployer_pubkey = [deployer_pubkey_prefix, shard_identifier].concat();
    let deployer_address = Address::from_slice(&deployer_pubkey);
    let deployer_nonce = 0;
    compute_smart_contract_address(deployer_address, deployer_nonce)
}

fn shard_id_from_name(name: &str) -> u8 {
    let name_hash = keccak256(name.as_bytes());
    name_hash[31]
}

pub fn dns_address_for_name(name: &str) -> Address {
    let shard_id = shard_id_from_name(name);
    compute_dns_address_for_shard_id(shard_id)
}

#[test]
fn test_compute_dns_address() {
    assert_eq!(
        bech32::encode(&compute_dns_address_for_shard_id(0)),
        "drt1qqqqqqqqqqqqqpgqnhvsujzd95jz6fyv3ldmynlf97tscs9nqqqqgewses"
    );
    assert_eq!(
        bech32::encode(&compute_dns_address_for_shard_id(1)),
        "drt1qqqqqqqqqqqqqpgqysmcsfkqed279x6jvs694th4e4v50p4pqqqsmjndc3"
    );
    assert_eq!(
        bech32::encode(&compute_dns_address_for_shard_id(2)),
        "drt1qqqqqqqqqqqqqpgqnk5fq8sgg4vc63ffzf7qez550xe2l5jgqqpqygxwm7"
    );
}

#[test]
fn test_dns_for_name() {
    assert_eq!(
        bech32::encode(&dns_address_for_name("test.numbat")),
        "drt1qqqqqqqqqqqqqpgqydsrr9rw6wgl83dm7mgyqzr7khqp982zqp4q36cj29"
    );
    assert_eq!(
        bech32::encode(&dns_address_for_name("helloworld.numbat")),
        "drt1qqqqqqqqqqqqqpgq37fmv57uqkayxctplh4swkw9vfz2jawvqqyqsqlm3z"
    );
}
