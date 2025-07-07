use dharitri_sc_scenario::imports::*;

use basic_features::token_identifier_features::TokenIdentifierFeatures;

#[test]
fn test_token_identifier_rewa() {
    let bf = basic_features::contract_obj::<StaticApi>();
    let result = bf.token_identifier_rewa();
    assert_eq!(RewaOrDcdtTokenIdentifier::rewa(), result);
}

/// This just tests the contract syntax.
/// For a complete suite of test cases, see `dharitri-sc-scenario/tests/managed_token_identifier_test.rs`.
#[test]
fn test_token_identifier_is_valid() {
    let bf = basic_features::contract_obj::<StaticApi>();
    let result = bf.token_identifier_is_valid_1(RewaOrDcdtTokenIdentifier::dcdt(
        TokenIdentifier::from(&b"ALC-6258d2"[..]),
    ));
    assert!(result);
    let result = bf.token_identifier_is_valid_1(RewaOrDcdtTokenIdentifier::dcdt(
        TokenIdentifier::from(&b"AL-C6258d2"[..]),
    ));
    assert!(!result);
    let result = bf.token_identifier_is_valid_2(ManagedBuffer::from(&b"12345-6258d2"[..]));
    assert!(result);
    let result = bf.token_identifier_is_valid_2(ManagedBuffer::from(&b"ALCCCCCCCCC-6258d2"[..]));
    assert!(!result);
}

#[test]
fn test_token_identifier_compare() {
    let token_id = TokenIdentifier::<StaticApi>::from(&b"ALC-6258d2"[..]);
    let dcdt_token_id = RewaOrDcdtTokenIdentifier::dcdt(token_id.clone());
    let wrong_dcdt_token_id =
        RewaOrDcdtTokenIdentifier::dcdt(TokenIdentifier::from(&b"AAA-111111"[..]));
    let rewa_token_id = RewaOrDcdtTokenIdentifier::rewa();

    assert!(token_id == dcdt_token_id);
    assert!(dcdt_token_id == token_id);

    assert!(token_id != rewa_token_id);
    assert!(rewa_token_id != token_id);

    assert!(token_id != wrong_dcdt_token_id);
    assert!(wrong_dcdt_token_id != token_id);
}
