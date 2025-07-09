use dharitri_sc::{
    codec::{self, DefaultErrorHandler, TopEncode},
    types::{BigUint, DcdtTokenPayment, TokenIdentifier},
};
use dharitri_sc_scenario::api::StaticApi;

/// Helper top-decode that doesn't rely on the `dcdt-token-payment-legacy-decode` feature flag.
fn dcdt_token_payment_backwards_compatible_top_decode_or_handle_err<I, H>(
    top_input: I,
    h: H,
) -> Result<DcdtTokenPayment<StaticApi>, H::HandledErr>
where
    I: codec::TopDecodeInput,
    H: codec::DecodeErrorHandler,
{
    let mut nested_buffer = top_input.into_nested_buffer();
    let result =
        DcdtTokenPayment::backwards_compatible_dep_decode_or_handle_err(&mut nested_buffer, h)?;
    if !codec::NestedDecodeInput::is_depleted(&nested_buffer) {
        return Err(h.handle_error(codec::DecodeError::INPUT_TOO_LONG));
    }
    Ok(result)
}

/// Helper top-decode that doesn't rely on the `dcdt-token-payment-legacy-decode` feature flag.
fn dcdt_token_payment_regular_top_decode_or_handle_err<I, H>(
    top_input: I,
    h: H,
) -> Result<DcdtTokenPayment<StaticApi>, H::HandledErr>
where
    I: codec::TopDecodeInput,
    H: codec::DecodeErrorHandler,
{
    let mut nested_buffer = top_input.into_nested_buffer();
    let result = DcdtTokenPayment::regular_dep_decode_or_handle_err(&mut nested_buffer, h)?;
    if !codec::NestedDecodeInput::is_depleted(&nested_buffer) {
        return Err(h.handle_error(codec::DecodeError::INPUT_TOO_LONG));
    }
    Ok(result)
}

#[test]
fn dcdt_token_payment_backwards_compatibility_decode() {
    let token_payment = DcdtTokenPayment::<StaticApi>::new(
        TokenIdentifier::from("MYTOKEN-12345"),
        0u64,
        BigUint::from(42u64),
    );

    let mut bytes = Vec::<u8>::new();
    token_payment.top_encode(&mut bytes).unwrap();

    // 1. decode as-is
    let decoded1_regular =
        dcdt_token_payment_regular_top_decode_or_handle_err(bytes.as_slice(), DefaultErrorHandler)
            .unwrap();
    assert_eq!(token_payment, decoded1_regular);

    let decoded1_bc = dcdt_token_payment_backwards_compatible_top_decode_or_handle_err(
        bytes.as_slice(),
        DefaultErrorHandler,
    )
    .unwrap();
    assert_eq!(token_payment, decoded1_bc);

    // 2. legacy token type = 0
    bytes.insert(0, 0u8);

    let decoded2_regular_result =
        dcdt_token_payment_regular_top_decode_or_handle_err(bytes.as_slice(), DefaultErrorHandler);
    assert!(decoded2_regular_result.is_err());

    let decoded2_bc = dcdt_token_payment_backwards_compatible_top_decode_or_handle_err(
        bytes.as_slice(),
        DefaultErrorHandler,
    )
    .unwrap();
    assert_eq!(token_payment, decoded2_bc);

    // 3. legacy token type = 1
    bytes[0] = 1u8;

    let decoded3_regular_result =
        dcdt_token_payment_regular_top_decode_or_handle_err(bytes.as_slice(), DefaultErrorHandler);
    assert!(decoded3_regular_result.is_err());

    let decoded3_bc = dcdt_token_payment_backwards_compatible_top_decode_or_handle_err(
        bytes.as_slice(),
        DefaultErrorHandler,
    )
    .unwrap();
    assert_eq!(token_payment, decoded3_bc);
}

#[test]
fn dcdt_token_payment_backwards_compatibility_decode_real_data() {
    let bytes = dharitri_sc::hex_literal::hex!(
        "020000000f41534852455741462d3236356334350000000000000001000000065af3107a4000"
    );
    let decoded = dcdt_token_payment_backwards_compatible_top_decode_or_handle_err(
        &bytes[..],
        DefaultErrorHandler,
    )
    .unwrap();
    assert_eq!(decoded.token_identifier.to_string(), "ASHREWAF-265c45");
    assert_eq!(decoded.token_nonce, 1);
    assert_eq!(decoded.amount, BigUint::from(0x5af3107a4000u64));
}
