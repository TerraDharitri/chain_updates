use dharitri_sc::types::{
    BoxedBytes, RewaOrDcdtTokenIdentifier, RewaOrDcdtTokenPayment, DcdtTokenPayment, ManagedBuffer,
    TokenIdentifier,
};
use dharitri_sc_scenario::{
    api::StaticApi, managed_rewa_token_id, managed_test_util::check_managed_top_encode_decode,
    managed_token_id, managed_token_id_wrapped, dharitri_sc,
};

#[test]
fn test_rewa() {
    assert!(RewaOrDcdtTokenIdentifier::<StaticApi>::rewa().is_rewa());
}

#[test]
fn test_codec() {
    check_managed_top_encode_decode(
        RewaOrDcdtTokenIdentifier::<StaticApi>::rewa(),
        RewaOrDcdtTokenIdentifier::<StaticApi>::REWA_REPRESENTATION,
    );

    let expected = BoxedBytes::from_concat(&[
        &[0, 0, 0, 4],
        &RewaOrDcdtTokenIdentifier::<StaticApi>::REWA_REPRESENTATION[..],
    ]);
    check_managed_top_encode_decode(
        vec![RewaOrDcdtTokenIdentifier::<StaticApi>::rewa()],
        expected.as_slice(),
    );
}

#[test]
#[rustfmt::skip]
fn test_is_valid_dcdt_identifier() {
    // valid identifier
    assert!(TokenIdentifier::<StaticApi>::from("ALC-6258d2").is_valid_dcdt_identifier());

    // valid identifier with numbers in ticker
    assert!(TokenIdentifier::<StaticApi>::from("ALC123-6258d2").is_valid_dcdt_identifier());

    // valid ticker only numbers
    assert!(TokenIdentifier::<StaticApi>::from("12345-6258d2").is_valid_dcdt_identifier());

    // missing dash
    assert!(!TokenIdentifier::<StaticApi>::from("ALC6258d2").is_valid_dcdt_identifier());

    // wrong dash position
    assert!(!TokenIdentifier::<StaticApi>::from("AL-C6258d2").is_valid_dcdt_identifier());

    // lowercase ticker
    assert!(!TokenIdentifier::<StaticApi>::from("alc-6258d2").is_valid_dcdt_identifier());

    // uppercase random chars
    assert!(!TokenIdentifier::<StaticApi>::from("ALC-6258D2").is_valid_dcdt_identifier());

    // too many random chars
    assert!(!TokenIdentifier::<StaticApi>::from("ALC-6258d2ff").is_valid_dcdt_identifier());

    // ticker too short
    assert!(!TokenIdentifier::<StaticApi>::from("AL-6258d2").is_valid_dcdt_identifier());

    // ticker too long
    assert!(!TokenIdentifier::<StaticApi>::from("ALCCCCCCCCC-6258d2").is_valid_dcdt_identifier());
}

#[test]
#[rustfmt::skip]
fn test_ticker() {
    // valid identifier
    assert_eq!(
        TokenIdentifier::<StaticApi>::from("ALC-6258d2").ticker(),
        ManagedBuffer::<StaticApi>::from("ALC"),
    );

    // valid identifier with numbers in ticker
    assert_eq!(
        TokenIdentifier::<StaticApi>::from("ALC123-6258d2").ticker(),
        ManagedBuffer::<StaticApi>::from("ALC123"),
    );

    // valid ticker only numbers
    assert_eq!(
        TokenIdentifier::<StaticApi>::from("12345-6258d2").ticker(),
        ManagedBuffer::<StaticApi>::from("12345"),
    );

    // missing dash
    assert_eq!(
        TokenIdentifier::<StaticApi>::from("ALC6258d2").ticker(),
        ManagedBuffer::<StaticApi>::from("AL"),
    );

    // wrong dash position
    assert_eq!(
        TokenIdentifier::<StaticApi>::from("AL-C6258d2").ticker(),
        ManagedBuffer::<StaticApi>::from("AL-"),
    );

    // lowercase ticker
    assert_eq!(
        TokenIdentifier::<StaticApi>::from("alc-6258d2").ticker(),
        ManagedBuffer::<StaticApi>::from("alc"),
    );

    // uppercase random chars
    assert_eq!(
        TokenIdentifier::<StaticApi>::from("ALC-6258D2").ticker(),
        ManagedBuffer::<StaticApi>::from("ALC"),
    );

    // too many random chars
    assert_eq!(
        TokenIdentifier::<StaticApi>::from("ALC-6258d2ff").ticker(),
        ManagedBuffer::<StaticApi>::from("ALC-6"),
    );

    // ticker too short
    assert_eq!(
        TokenIdentifier::<StaticApi>::from("AL-6258d2").ticker(),
        ManagedBuffer::<StaticApi>::from("AL"),
    );

    // ticker too long
    assert_eq!(
        TokenIdentifier::<StaticApi>::from("ALCCCCCCCCC-6258d2").ticker(),
        ManagedBuffer::<StaticApi>::from("ALCCCCCCCCC"),
    );
}

#[test]
fn test_is_valid_rewa_or_dcdt() {
    // rewa is always valid
    assert!(RewaOrDcdtTokenIdentifier::<StaticApi>::rewa().is_valid());

    // valid dcdt
    assert!(
        RewaOrDcdtTokenIdentifier::<StaticApi>::dcdt(TokenIdentifier::from("ALC-6258d2"))
            .is_valid()
    );

    // invalid dcdt, see above
    assert!(
        !RewaOrDcdtTokenIdentifier::<StaticApi>::dcdt(TokenIdentifier::from("ALCCCCCCCCC-6258d2"))
            .is_valid()
    );
}

#[test]
fn test_token_identifier_eq() {
    assert_eq!(
        TokenIdentifier::<StaticApi>::from("DCDT-00000"),
        TokenIdentifier::<StaticApi>::from("DCDT-00000")
    );
    assert_ne!(
        TokenIdentifier::<StaticApi>::from("DCDT-00001"),
        TokenIdentifier::<StaticApi>::from("DCDT-00002")
    );

    assert_eq!(
        RewaOrDcdtTokenIdentifier::<StaticApi>::dcdt(TokenIdentifier::from("DCDT-00003")),
        TokenIdentifier::<StaticApi>::from("DCDT-00003")
    );
    assert_ne!(
        RewaOrDcdtTokenIdentifier::<StaticApi>::rewa(),
        TokenIdentifier::<StaticApi>::from("ANYTHING-1234")
    );
    assert_ne!(
        RewaOrDcdtTokenIdentifier::<StaticApi>::rewa(),
        TokenIdentifier::<StaticApi>::from("REWA")
    );
}

#[test]
fn test_payment_eq() {
    assert_eq!(
        DcdtTokenPayment::<StaticApi>::new("PAY-00000".into(), 0, 1000u32.into()),
        DcdtTokenPayment::<StaticApi>::new("PAY-00000".into(), 0, 1000u32.into()),
    );
    assert_ne!(
        DcdtTokenPayment::<StaticApi>::new("PAY-00001".into(), 0, 1000u32.into()),
        DcdtTokenPayment::<StaticApi>::new("PAY-00002".into(), 0, 1000u32.into()),
    );
    assert_eq!(
        RewaOrDcdtTokenPayment::<StaticApi>::no_payment(),
        RewaOrDcdtTokenPayment::<StaticApi>::no_payment(),
    );
    assert_eq!(
        RewaOrDcdtTokenPayment::<StaticApi>::new(
            RewaOrDcdtTokenIdentifier::dcdt("DCDTPAY-00000"),
            0,
            1000u32.into()
        ),
        RewaOrDcdtTokenPayment::<StaticApi>::new(
            RewaOrDcdtTokenIdentifier::dcdt("DCDTPAY-00000"),
            0,
            1000u32.into()
        ),
    );
    assert_ne!(
        RewaOrDcdtTokenPayment::<StaticApi>::new(
            RewaOrDcdtTokenIdentifier::dcdt("DCDTPAY-00001"),
            0,
            1000u32.into()
        ),
        RewaOrDcdtTokenPayment::<StaticApi>::new(
            RewaOrDcdtTokenIdentifier::dcdt("DCDTPAY-00002"),
            0,
            1000u32.into()
        ),
    );
    assert_ne!(
        RewaOrDcdtTokenPayment::<StaticApi>::new(
            RewaOrDcdtTokenIdentifier::dcdt("DCDTPAY-00001"),
            0,
            1000u32.into()
        ),
        RewaOrDcdtTokenPayment::<StaticApi>::no_payment(),
    );
}

#[test]
fn test_managed_token_id_macro() {
    assert_eq!(
        managed_rewa_token_id!(),
        RewaOrDcdtTokenIdentifier::<StaticApi>::rewa()
    );
    assert_eq!(
        managed_token_id!(b"ALC-6258d2"),
        TokenIdentifier::<StaticApi>::from("ALC-6258d2")
    );
    assert_eq!(
        managed_token_id_wrapped!(b"ALC-6258d2").unwrap_dcdt(),
        TokenIdentifier::<StaticApi>::from("ALC-6258d2")
    )
}
