use crate::{
    api::{
        const_handles, use_raw_handle, CallValueApi, CallValueApiImpl, ErrorApi, ErrorApiImpl,
        ManagedBufferApiImpl, ManagedTypeApi,
    },
    contract_base::CallValueWrapper,
    err_msg,
    types::{
        BigUint, RewaOrDcdtTokenIdentifier, DcdtTokenPayment, ManagedRef, ManagedType, ManagedVec,
    },
};

/// Called initially in the generated code whenever no payable annotation is provided.
pub fn not_payable<A>()
where
    A: CallValueApi,
{
    A::call_value_api_impl().check_not_payable();
}

/// Called initially in the generated code whenever `#[payable("*")]` annotation is provided.
pub fn payable_any<A>()
where
    A: CallValueApi,
{
}

/// Called initially in the generated code whenever `#[payable("REWA")]` annotation is provided.
pub fn payable_rewa<A>()
where
    A: CallValueApi + ErrorApi + ManagedTypeApi,
{
    // will crash if anything other than (single) REWA was transferred
    let _ = CallValueWrapper::<A>::new().rewa();
}

/// Called initially in the generated code whenever `#[payable("<token identifier>")]` annotation is provided.
///
/// Was never really used, expected to be deprecated/removed.
pub fn payable_single_specific_token<A>(expected_tokend_identifier: &str)
where
    A: CallValueApi + ManagedTypeApi + ErrorApi,
{
    let transfers = CallValueWrapper::<A>::new().all_dcdt_transfers();
    if transfers.len() != 1 {
        A::error_api_impl().signal_error(err_msg::SINGLE_DCDT_EXPECTED.as_bytes());
    }
    let expected_token_handle: A::ManagedBufferHandle =
        use_raw_handle(const_handles::MBUF_TEMPORARY_1);
    A::managed_type_impl().mb_overwrite(
        expected_token_handle.clone(),
        expected_tokend_identifier.as_bytes(),
    );
    let transfer = transfers.get(0);
    if !A::managed_type_impl().mb_eq(
        transfer.token_identifier.get_handle(),
        expected_token_handle,
    ) {
        A::error_api_impl().signal_error(err_msg::BAD_TOKEN_PROVIDED.as_bytes());
    }
}

/// Initializes an argument annotated with `#[payment_amount]` or `#[payment]`.
pub fn arg_payment_amount<A>() -> BigUint<A>
where
    A: CallValueApi + ManagedTypeApi,
{
    CallValueWrapper::<A>::new().rewa_or_single_dcdt().amount
}

/// Initializes an argument annotated with `#[payment_token]`.
pub fn arg_payment_token<A>() -> RewaOrDcdtTokenIdentifier<A>
where
    A: CallValueApi + ManagedTypeApi,
{
    CallValueWrapper::<A>::new()
        .rewa_or_single_dcdt()
        .token_identifier
}

/// Initializes an argument annotated with `#[payment_nonce]`.
pub fn arg_payment_nonce<A>() -> u64
where
    A: CallValueApi + ManagedTypeApi,
{
    CallValueWrapper::<A>::new()
        .rewa_or_single_dcdt()
        .token_nonce
}

/// Initializes an argument annotated with `#[payment_multi]`.
pub fn arg_payment_multi<A>() -> ManagedRef<'static, A, ManagedVec<A, DcdtTokenPayment<A>>>
where
    A: CallValueApi + ManagedTypeApi,
{
    CallValueWrapper::<A>::new().all_dcdt_transfers()
}
