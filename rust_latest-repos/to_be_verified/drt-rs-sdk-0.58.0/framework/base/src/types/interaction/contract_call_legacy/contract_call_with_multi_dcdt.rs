use crate::codec::TopEncodeMulti;

use crate::{
    api::CallTypeApi,
    types::{
        BigUint, DcdtTokenPayment, ManagedAddress, ManagedBuffer, ManagedVec, TokenIdentifier,
    },
};

use super::{
    contract_call_no_payment::ContractCallNoPayment, contract_call_trait::ContractCallBase,
    ContractCall, ContractCallWithRewa,
};

#[deprecated(
    since = "0.49.0",
    note = "Please use the unified transaction syntax instead."
)]
#[must_use]
pub struct ContractCallWithMultiDcdt<SA, OriginalResult>
where
    SA: CallTypeApi + 'static,
{
    pub basic: ContractCallNoPayment<SA, OriginalResult>,
    pub dcdt_payments: ManagedVec<SA, DcdtTokenPayment<SA>>,
}

impl<SA, OriginalResult> ContractCallBase<SA> for ContractCallWithMultiDcdt<SA, OriginalResult>
where
    SA: CallTypeApi + 'static,
    OriginalResult: TopEncodeMulti,
{
    type OriginalResult = OriginalResult;

    fn into_normalized(self) -> ContractCallWithRewa<SA, Self::OriginalResult> {
        self.basic
            .into_normalized()
            .convert_to_dcdt_transfer_call(self.dcdt_payments)
    }
}

impl<SA, OriginalResult> ContractCall<SA> for ContractCallWithMultiDcdt<SA, OriginalResult>
where
    SA: CallTypeApi + 'static,
    OriginalResult: TopEncodeMulti,
{
    #[inline]
    fn get_mut_basic(&mut self) -> &mut ContractCallNoPayment<SA, OriginalResult> {
        &mut self.basic
    }

    fn transfer_execute(self) {
        self.basic.transfer_execute_dcdt(self.dcdt_payments);
    }
}

impl<SA, OriginalResult> ContractCallWithMultiDcdt<SA, OriginalResult>
where
    SA: CallTypeApi + 'static,
{
    /// Creates a new instance directly.
    ///
    /// The constructor is mostly for hand-written proxies,
    /// the usual way of constructing this object is via the builder methods of other contract call types,
    /// especially `with_dcdt_transfer` or `with_multi_token_transfer`.
    pub fn new<N: Into<ManagedBuffer<SA>>>(
        to: ManagedAddress<SA>,
        endpoint_name: N,
        payments: ManagedVec<SA, DcdtTokenPayment<SA>>,
    ) -> Self {
        ContractCallWithMultiDcdt {
            basic: ContractCallNoPayment::new(to, endpoint_name),
            dcdt_payments: payments,
        }
    }

    /// Adds a single DCDT token transfer to a contract call.
    ///
    /// Can be called multiple times on the same call.
    pub fn with_dcdt_transfer<P: Into<DcdtTokenPayment<SA>>>(mut self, payment: P) -> Self {
        self.dcdt_payments.push(payment.into());
        self
    }

    #[deprecated(
        since = "0.39.0",
        note = "Replace by `contract_call.with_dcdt_transfer((payment_token, payment_nonce, payment_amount))`. 
        The tuple argument will get automatically converted to DcdtTokenPayment."
    )]
    pub fn add_dcdt_token_transfer(
        self,
        payment_token: TokenIdentifier<SA>,
        payment_nonce: u64,
        payment_amount: BigUint<SA>,
    ) -> Self {
        self.with_dcdt_transfer((payment_token, payment_nonce, payment_amount))
    }
}
