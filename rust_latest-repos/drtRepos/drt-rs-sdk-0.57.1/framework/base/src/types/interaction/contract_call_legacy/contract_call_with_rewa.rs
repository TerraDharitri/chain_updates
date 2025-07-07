use crate::codec::TopEncodeMulti;

use crate::{
    api::CallTypeApi,
    types::{BigUint, ManagedAddress, ManagedBuffer},
};

use super::{
    contract_call_no_payment::ContractCallNoPayment, contract_call_trait::ContractCallBase,
    ContractCall,
};

/// Holds data for calling another contract, with REWA payment only.
///
/// Gets created when chaining method `with_rewa_transfer`.
///
/// If the payment is zero, it bevahes exactly like `ContractCallNoPayment`.
///
/// It also represents the normalized form of any contract call, since DCDT transfers
/// (the only payment not available here) get converted to builtin function calls in normalized form.
#[deprecated(
    since = "0.49.0",
    note = "Please use the unified transaction syntax instead."
)]
#[must_use]
pub struct ContractCallWithRewa<SA, OriginalResult>
where
    SA: CallTypeApi + 'static,
{
    pub basic: ContractCallNoPayment<SA, OriginalResult>,
    pub rewa_payment: BigUint<SA>,
}

impl<SA, OriginalResult> ContractCallBase<SA> for ContractCallWithRewa<SA, OriginalResult>
where
    SA: CallTypeApi + 'static,
    OriginalResult: TopEncodeMulti,
{
    type OriginalResult = OriginalResult;

    #[inline]
    fn into_normalized(self) -> ContractCallWithRewa<SA, Self::OriginalResult> {
        // no DCDT, no conversion needed
        self
    }
}

impl<SA, OriginalResult> ContractCall<SA> for ContractCallWithRewa<SA, OriginalResult>
where
    SA: CallTypeApi + 'static,
    OriginalResult: TopEncodeMulti,
{
    #[inline]
    fn get_mut_basic(&mut self) -> &mut ContractCallNoPayment<SA, OriginalResult> {
        &mut self.basic
    }

    fn transfer_execute(self) {
        self.basic.transfer_execute_rewa(self.rewa_payment);
    }
}

impl<SA, OriginalResult> ContractCallWithRewa<SA, OriginalResult>
where
    SA: CallTypeApi + 'static,
{
    /// Creates a new instance directly.
    ///
    /// The constructor is mostly for hand-written proxies,
    /// the usual way of constructing this object is via the builder methods of other contract call types,
    /// especially `with_rewa_transfer`.
    pub fn new<N: Into<ManagedBuffer<SA>>>(
        to: ManagedAddress<SA>,
        endpoint_name: N,
        rewa_payment: BigUint<SA>,
    ) -> Self {
        ContractCallWithRewa {
            basic: ContractCallNoPayment::new(to, endpoint_name),
            rewa_payment,
        }
    }
}
