mod test_dcdt_transfer;
mod tx_payment_rewa;
mod tx_payment_rewa_or_dcdt;
mod tx_payment_rewa_or_dcdt_refs;
mod tx_payment_rewa_or_multi_dcdt;
mod tx_payment_rewa_or_multi_dcdt_ref;
mod tx_payment_rewa_value;
mod tx_payment_multi_rewa_or_dcdt;
mod tx_payment_multi_dcdt;
mod tx_payment_none;
mod tx_payment_not_payable;
mod tx_payment_single_dcdt;
mod tx_payment_single_dcdt_ref;
mod tx_payment_single_dcdt_triple;

pub use test_dcdt_transfer::TestDcdtTransfer;
pub use tx_payment_rewa::{Rewa, RewaPayment};
pub use tx_payment_rewa_value::TxRewaValue;
pub use tx_payment_multi_dcdt::TxPaymentMultiDcdt;
pub use tx_payment_not_payable::NotPayable;

use crate::{
    api::ManagedTypeApi,
    types::{BigUint, ManagedAddress, ManagedBuffer, MultiRewaOrDcdtPayment},
};

use super::{AnnotatedValue, FunctionCall, TxEnv, TxFrom, TxToSpecified};

/// Describes a payment that is part of a transaction.
#[diagnostic::on_unimplemented(
    message = "Type `{Self}` cannot be used as payment (does not implement `TxPayment<{Env}>`)",
    label = "not a valid payment type",
    note = "there are multiple ways to specify the transaction payment, but `{Self}` is not one of them"
)]
pub trait TxPayment<Env>
where
    Env: TxEnv,
{
    /// Returns true if payment indicates transfer of either non-zero REWA or DCDT amounts.
    fn is_no_payment(&self, env: &Env) -> bool;

    /// Transfer-execute calls have different APIs for different payments types.
    /// This method selects between them.
    fn perform_transfer_execute(
        self,
        env: &Env,
        to: &ManagedAddress<Env::Api>,
        gas_limit: u64,
        fc: FunctionCall<Env::Api>,
    );

    /// Converts an DCDT call to a built-in function call, if necessary.
    fn with_normalized<From, To, F, R>(
        self,
        env: &Env,
        from: &From,
        to: To,
        fc: FunctionCall<Env::Api>,
        f: F,
    ) -> R
    where
        From: TxFrom<Env>,
        To: TxToSpecified<Env>,
        F: FnOnce(&ManagedAddress<Env::Api>, &BigUint<Env::Api>, FunctionCall<Env::Api>) -> R;

    /// Payment data to be used by the testing framework. Will be refactored.
    fn into_full_payment_data(self, env: &Env) -> FullPaymentData<Env::Api>;
}

/// Marker trait that indicates that payment field contains no payment.
///
/// Implemented by `()` and `NotPayable`.
pub trait TxNoPayment<Env>: TxPayment<Env>
where
    Env: TxEnv,
{
}

/// Marks a payment object that only contains REWA or nothing at all.
pub trait TxPaymentRewaOnly<Env>: TxPayment<Env> + AnnotatedValue<Env, BigUint<Env::Api>>
where
    Env: TxEnv,
{
    #[inline]
    fn with_rewa_value<F, R>(&self, env: &Env, f: F) -> R
    where
        F: FnOnce(&BigUint<Env::Api>) -> R,
    {
        self.with_value_ref(env, f)
    }

    fn into_rewa_payment(self, env: &Env) -> BigUint<Env::Api> {
        self.into_value(env)
    }
}

#[derive(Clone)]
pub struct AnnotatedRewaPayment<Api>
where
    Api: ManagedTypeApi,
{
    pub value: BigUint<Api>,
    pub annotation: ManagedBuffer<Api>,
}

impl<Api> AnnotatedRewaPayment<Api>
where
    Api: ManagedTypeApi,
{
    pub fn new_rewa(value: BigUint<Api>) -> Self {
        let annotation = value.to_display();
        AnnotatedRewaPayment { value, annotation }
    }
}

#[derive(Clone)]
pub struct FullPaymentData<Api>
where
    Api: ManagedTypeApi,
{
    pub rewa: Option<AnnotatedRewaPayment<Api>>,
    pub multi_dcdt: MultiRewaOrDcdtPayment<Api>,
}

impl<Api> Default for FullPaymentData<Api>
where
    Api: ManagedTypeApi,
{
    fn default() -> Self {
        Self {
            rewa: None,
            multi_dcdt: Default::default(),
        }
    }
}
