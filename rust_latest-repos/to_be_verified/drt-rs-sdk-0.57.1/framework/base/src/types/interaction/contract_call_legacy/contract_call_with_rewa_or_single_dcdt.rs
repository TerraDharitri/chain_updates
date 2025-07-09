use crate::codec::TopEncodeMulti;

use crate::{
    api::CallTypeApi,
    types::{
        BigUint, RewaOrDcdtTokenIdentifier, RewaOrDcdtTokenPayment, ManagedAddress, ManagedBuffer,
    },
};

use super::{
    contract_call_no_payment::ContractCallNoPayment, contract_call_trait::ContractCallBase,
    ContractCall, ContractCallWithRewa,
};

/// Holds data for calling another contract, with a single payment, either REWA or a single DCDT token.
///
/// Gets created when chaining method `with_rewa_or_single_dcdt_transfer`.
#[deprecated(
    since = "0.49.0",
    note = "Please use the unified transaction syntax instead."
)]
#[must_use]
pub struct ContractCallWithRewaOrSingleDcdt<SA, OriginalResult>
where
    SA: CallTypeApi + 'static,
{
    pub(super) basic: ContractCallNoPayment<SA, OriginalResult>,
    pub payment: RewaOrDcdtTokenPayment<SA>,
}

impl<SA, OriginalResult> ContractCallWithRewaOrSingleDcdt<SA, OriginalResult>
where
    SA: CallTypeApi + 'static,
    OriginalResult: TopEncodeMulti,
{
    fn into_normalized_rewa(self) -> ContractCallWithRewa<SA, OriginalResult> {
        ContractCallWithRewa {
            basic: self.basic,
            rewa_payment: self.payment.amount,
        }
    }

    fn into_normalized_dcdt(self) -> ContractCallWithRewa<SA, OriginalResult> {
        self.basic
            .into_normalized()
            .convert_to_single_transfer_dcdt_call(self.payment.unwrap_dcdt())
    }
}

impl<SA, OriginalResult> ContractCallBase<SA>
    for ContractCallWithRewaOrSingleDcdt<SA, OriginalResult>
where
    SA: CallTypeApi + 'static,
    OriginalResult: TopEncodeMulti,
{
    type OriginalResult = OriginalResult;

    fn into_normalized(self) -> ContractCallWithRewa<SA, Self::OriginalResult> {
        if self.payment.token_identifier.is_rewa() {
            self.into_normalized_rewa()
        } else {
            // Because we know that there can be at most one DCDT payment,
            // there is no need to call the full `convert_to_dcdt_transfer_call`.
            self.into_normalized_dcdt()
        }
    }
}

impl<SA, OriginalResult> ContractCall<SA> for ContractCallWithRewaOrSingleDcdt<SA, OriginalResult>
where
    SA: CallTypeApi + 'static,
    OriginalResult: TopEncodeMulti,
{
    #[inline]
    fn get_mut_basic(&mut self) -> &mut ContractCallNoPayment<SA, OriginalResult> {
        &mut self.basic
    }

    fn transfer_execute(self) {
        if self.payment.token_identifier.is_rewa() {
            self.basic.transfer_execute_rewa(self.payment.amount);
        } else {
            self.basic
                .transfer_execute_single_dcdt(self.payment.unwrap_dcdt());
        }
    }
}

impl<SA, OriginalResult> ContractCallWithRewaOrSingleDcdt<SA, OriginalResult>
where
    SA: CallTypeApi + 'static,
    OriginalResult: TopEncodeMulti,
{
    /// Creates a new instance directly.
    ///
    /// The constructor is mostly for hand-written proxies,
    /// the usual way of constructing this object is via the builder methods of other contract call types,
    /// especially `with_rewa_or_single_dcdt_transfer`.
    pub fn new<N: Into<ManagedBuffer<SA>>>(
        to: ManagedAddress<SA>,
        endpoint_name: N,
        token_identifier: RewaOrDcdtTokenIdentifier<SA>,
        token_nonce: u64,
        amount: BigUint<SA>,
    ) -> Self {
        ContractCallWithRewaOrSingleDcdt {
            basic: ContractCallNoPayment::new(to, endpoint_name),
            payment: RewaOrDcdtTokenPayment::new(token_identifier, token_nonce, amount),
        }
    }
}
