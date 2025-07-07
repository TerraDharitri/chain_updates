use core::marker::PhantomData;

use dharitri_chain_core::REWA_000000_TOKEN_IDENTIFIER;

use crate::{
    api::{
        const_handles, use_raw_handle, CallValueApi, CallValueApiImpl, ErrorApi, ErrorApiImpl,
        ManagedBufferApiImpl, ManagedTypeApi, RawHandle, StaticVarApiFlags, StaticVarApiImpl,
    },
    err_msg,
    types::{
        big_num_cmp::bi_gt_zero, BigUint, RewaDecimals, RewaOrDcdtTokenIdentifier,
        RewaOrDcdtTokenPayment, RewaOrMultiDcdtPayment, DcdtTokenPayment, ManagedDecimal,
        ManagedRef, ManagedType, ManagedVec, ManagedVecItem, ManagedVecItemPayload,
        ManagedVecPayloadIterator, ManagedVecRef, TokenIdentifier,
    },
};

#[derive(Default)]
pub struct CallValueWrapper<A>
where
    A: CallValueApi + ErrorApi + ManagedTypeApi,
{
    _phantom: PhantomData<A>,
}

impl<A> CallValueWrapper<A>
where
    A: CallValueApi + ErrorApi + ManagedTypeApi,
{
    pub fn new() -> Self {
        CallValueWrapper {
            _phantom: PhantomData,
        }
    }

    /// Cached transfers from the VM.
    fn all_dcdt_transfers_unchecked(&self) -> A::ManagedBufferHandle {
        let all_transfers_unchecked_handle: A::ManagedBufferHandle =
            use_raw_handle(const_handles::CALL_VALUE_MULTI_DCDT);
        if !A::static_var_api_impl()
            .flag_is_set_or_update(StaticVarApiFlags::CALL_VALUE_DCDT_UNCHECKED_INITIALIZED)
        {
            A::call_value_api_impl()
                .load_all_dcdt_transfers(all_transfers_unchecked_handle.clone());
        }
        all_transfers_unchecked_handle
    }

    /// Retrieves the REWA call value from the VM.
    ///
    /// Will return 0 in case of an DCDT transfer, even though REWA and DCDT transfers are now posible.
    pub fn rewa_direct_non_strict(&self) -> ManagedRef<'static, A, BigUint<A>> {
        let call_value_handle: A::BigIntHandle = use_raw_handle(const_handles::CALL_VALUE_REWA);
        if !A::static_var_api_impl()
            .flag_is_set_or_update(StaticVarApiFlags::CALL_VALUE_REWA_DIRECT_INITIALIZED)
        {
            A::call_value_api_impl().load_rewa_value(call_value_handle.clone());
        }
        unsafe { ManagedRef::wrap_handle(call_value_handle) }
    }

    /// Retrieves the REWA call value and crashes if anything else was transferred.
    ///
    /// Accepts both REWA sent directly, as well as REWA sent alone in a multi-transfer.
    ///
    /// Does not accept a multi-transfer with 2 or more transfers, not even 2 or more REWA transfers.
    pub fn rewa(&self) -> ManagedRef<'static, A, BigUint<A>> {
        let dcdt_transfers_handle = self.all_dcdt_transfers_unchecked();
        let dcdt_transfers: ManagedRef<'static, A, ManagedVec<A, RewaOrDcdtTokenPayment<A>>> =
            unsafe { ManagedRef::wrap_handle(dcdt_transfers_handle) };
        match dcdt_transfers.len() {
            0 => self.rewa_direct_non_strict(),
            1 => {
                let first = dcdt_transfers.get(0);
                if !first.token_identifier.is_rewa() {
                    A::error_api_impl().signal_error(err_msg::NON_PAYABLE_FUNC_DCDT.as_bytes());
                }
                unsafe { ManagedRef::wrap_handle(first.amount.get_handle()) }
            },
            _ => A::error_api_impl().signal_error(err_msg::INCORRECT_NUM_TRANSFERS.as_bytes()),
        }
    }

    /// Retrieves the REWA call value from the VM.
    ///
    /// Will return 0 in case of an DCDT transfer, even though REWA and DCDT transfers are now posible.
    ///
    /// ## Important!
    ///
    /// Does not cover multi-transfer scenarios properly, but left for backwards compatibility.
    ///
    /// Please use `.rewa()` instead!
    ///
    /// For raw handling, `.rewa_direct_non_strict()` is also acceptable.
    #[deprecated(
        since = "0.55.0",
        note = "Does not cover multi-transfer scenarios properly, but left for backwards compatibility. Please use .rewa() instead!"
    )]
    pub fn rewa_value(&self) -> ManagedRef<'static, A, BigUint<A>> {
        self.rewa_direct_non_strict()
    }

    /// Returns the REWA call value from the VM as ManagedDecimal
    pub fn rewa_decimal(&self) -> ManagedDecimal<A, RewaDecimals> {
        ManagedDecimal::<A, RewaDecimals>::const_decimals_from_raw(self.rewa_value().clone())
    }

    /// Returns all DCDT transfers that accompany this SC call.
    /// Will return 0 results if nothing was transfered, or just REWA.
    ///
    /// Will crash for REWA + DCDT multi transfers.
    pub fn all_dcdt_transfers(&self) -> ManagedRef<'static, A, ManagedVec<A, DcdtTokenPayment<A>>> {
        let multi_dcdt_handle: A::ManagedBufferHandle = self.all_dcdt_transfers_unchecked();
        let checked = A::static_var_api_impl()
            .flag_is_set_or_update(StaticVarApiFlags::CALL_VALUE_DCDT_INITIALIZED);
        if !checked && rewa_000000_transfer_exists::<A>(multi_dcdt_handle.clone()) {
            A::error_api_impl().signal_error(err_msg::DCDT_UNEXPECTED_REWA.as_bytes())
        }

        unsafe { ManagedRef::wrap_handle(multi_dcdt_handle) }
    }

    /// Will return all transfers in the form of a list of RewaOrDcdtTokenPayment.
    ///
    /// Both REWA and DCDT can be returned.
    ///
    /// In case of a single REWA transfer, only one item will be returned,
    /// the REWA payment represented as an DCDT transfer (REWA-000000).
    pub fn all_transfers(
        &self,
    ) -> ManagedRef<'static, A, ManagedVec<A, RewaOrDcdtTokenPayment<A>>> {
        let all_transfers_handle: A::ManagedBufferHandle =
            use_raw_handle(const_handles::CALL_VALUE_ALL);
        if !A::static_var_api_impl()
            .flag_is_set_or_update(StaticVarApiFlags::CALL_VALUE_ALL_INITIALIZED)
        {
            let rewa_single = self.rewa_direct_non_strict();
            if bi_gt_zero::<A>(rewa_single.get_handle()) {
                A::managed_type_impl().mb_overwrite(
                    use_raw_handle(const_handles::MBUF_REWA_000000),
                    REWA_000000_TOKEN_IDENTIFIER.as_bytes(),
                );
                A::managed_type_impl().mb_overwrite(
                    all_transfers_handle.clone(),
                    &const_handles::REWA_PAYMENT_PAYLOAD[..],
                );
            } else {
                // clone all_dcdt_transfers_unchecked -> all_transfers
                let all_transfers_unchecked_handle = self.all_dcdt_transfers_unchecked();
                A::managed_type_impl().mb_overwrite(all_transfers_handle.clone(), &[]);
                A::managed_type_impl()
                    .mb_append(all_transfers_handle.clone(), all_transfers_unchecked_handle);
            }
        }
        unsafe { ManagedRef::wrap_handle(all_transfers_handle) }
    }

    /// Verify and casts the received multi DCDT transfer in to an array.
    ///
    /// Can be used to extract all payments in one line like this:
    ///
    /// `let [payment_a, payment_b, payment_c] = self.call_value().multi_dcdt();`.
    ///
    /// Rejects REWA transfers. Switch to `multi_rewa_or_dcdt` to accept mixed transfers.
    pub fn multi_dcdt<const N: usize>(&self) -> [ManagedVecRef<'static, DcdtTokenPayment<A>>; N] {
        let dcdt_transfers = self.all_dcdt_transfers();
        let array = dcdt_transfers.to_array_of_refs::<N>().unwrap_or_else(|| {
            A::error_api_impl().signal_error(err_msg::INCORRECT_NUM_DCDT_TRANSFERS.as_bytes())
        });
        unsafe { core::mem::transmute(array) }
    }

    /// Verify and casts the received multi DCDT transfer in to an array.
    ///
    /// Can be used to extract all payments in one line like this:
    ///
    /// `let [payment_a, payment_b, payment_c] = self.call_value().multi_rewa_or_dcdt();`.
    pub fn multi_rewa_or_dcdt<const N: usize>(
        &self,
    ) -> [ManagedVecRef<'static, RewaOrDcdtTokenPayment<A>>; N] {
        let dcdt_transfers = self.all_transfers();
        let array = dcdt_transfers.to_array_of_refs::<N>().unwrap_or_else(|| {
            A::error_api_impl().signal_error(err_msg::INCORRECT_NUM_TRANSFERS.as_bytes())
        });
        unsafe { core::mem::transmute(array) }
    }

    /// Expects precisely one DCDT token transfer, fungible or not.
    ///
    /// Will return the received DCDT payment.
    ///
    /// The amount cannot be 0, since that would not qualify as an DCDT transfer.
    pub fn single_dcdt(&self) -> ManagedVecRef<'static, DcdtTokenPayment<A>> {
        let dcdt_transfers = self.all_dcdt_transfers();
        if dcdt_transfers.len() != 1 {
            A::error_api_impl().signal_error(err_msg::INCORRECT_NUM_DCDT_TRANSFERS.as_bytes())
        }
        let value = dcdt_transfers.get(0);
        unsafe { core::mem::transmute(value) }
    }

    /// Expects precisely one fungible DCDT token transfer.
    ///
    /// Returns the token ID and the amount for fungible DCDT transfers.
    ///
    /// The amount cannot be 0, since that would not qualify as an DCDT transfer.
    pub fn single_fungible_dcdt(
        &self,
    ) -> (
        ManagedRef<'static, A, TokenIdentifier<A>>,
        ManagedRef<'static, A, BigUint<A>>,
    ) {
        let payment = self.single_dcdt();
        if payment.token_nonce != 0 {
            A::error_api_impl().signal_error(err_msg::FUNGIBLE_TOKEN_EXPECTED_ERR_MSG.as_bytes());
        }

        unsafe {
            (
                ManagedRef::wrap_handle(payment.token_identifier.get_handle()),
                ManagedRef::wrap_handle(payment.amount.get_handle()),
            )
        }
    }

    /// Accepts and returns either an REWA payment, or a single DCDT token.
    ///
    /// Will halt execution if more than one DCDT transfer was received.
    ///
    /// In case no transfer of value happen, it will return a payment of 0 REWA.
    pub fn rewa_or_single_dcdt(&self) -> RewaOrDcdtTokenPayment<A> {
        let dcdt_transfers_handle = self.all_dcdt_transfers_unchecked();
        let dcdt_transfers: ManagedRef<'static, A, ManagedVec<A, RewaOrDcdtTokenPayment<A>>> =
            unsafe { ManagedRef::wrap_handle(dcdt_transfers_handle) };
        match dcdt_transfers.len() {
            0 => RewaOrDcdtTokenPayment {
                token_identifier: RewaOrDcdtTokenIdentifier::rewa(),
                token_nonce: 0,
                amount: self.rewa_direct_non_strict().clone(),
            },
            1 => dcdt_transfers.get(0).clone(),
            _ => A::error_api_impl().signal_error(err_msg::INCORRECT_NUM_DCDT_TRANSFERS.as_bytes()),
        }
    }

    /// Accepts and returns either an REWA payment, or a single fungible DCDT token.
    ///
    /// Will halt execution if more than one DCDT transfer was received, or if the received DCDT is non- or semi-fungible.
    ///
    /// Works similar to `rewa_or_single_dcdt`,
    /// but checks the nonce to be 0 and returns a tuple of just token identifier and amount, for convenience.
    ///
    /// In case no transfer of value happen, it will return a payment of 0 REWA.
    pub fn rewa_or_single_fungible_dcdt(&self) -> (RewaOrDcdtTokenIdentifier<A>, BigUint<A>) {
        let payment = self.rewa_or_single_dcdt();
        if payment.token_nonce != 0 {
            A::error_api_impl().signal_error(err_msg::FUNGIBLE_TOKEN_EXPECTED_ERR_MSG.as_bytes());
        }

        (payment.token_identifier, payment.amount)
    }

    /// Accepts any sort of patyment, which is either:
    /// - REWA (can be zero in case of no payment whatsoever);
    /// - Multi-DCDT (one or more DCDT transfers).
    pub fn any_payment(&self) -> RewaOrMultiDcdtPayment<A> {
        let dcdt_transfers = self.all_dcdt_transfers();
        if dcdt_transfers.is_empty() {
            RewaOrMultiDcdtPayment::Rewa(self.rewa_direct_non_strict().clone())
        } else {
            RewaOrMultiDcdtPayment::MultiDcdt(dcdt_transfers.clone())
        }
    }
}

fn rewa_000000_transfer_exists<A>(transfers_vec_handle: A::ManagedBufferHandle) -> bool
where
    A: CallValueApi + ErrorApi + ManagedTypeApi,
{
    A::managed_type_impl().mb_overwrite(
        use_raw_handle(const_handles::MBUF_REWA_000000),
        REWA_000000_TOKEN_IDENTIFIER.as_bytes(),
    );
    unsafe {
        let mut iter: ManagedVecPayloadIterator<
            A,
            <DcdtTokenPayment<A> as ManagedVecItem>::PAYLOAD,
        > = ManagedVecPayloadIterator::new(transfers_vec_handle);

        iter.any(|payload| {
            let token_identifier_handle = RawHandle::read_from_payload(payload.slice_unchecked(0));
            A::managed_type_impl().mb_eq(
                use_raw_handle(const_handles::MBUF_REWA_000000),
                use_raw_handle(token_identifier_handle),
            )
        })
    }
}
