use crate::{
    contract_base::SendRawWrapper,
    types::{
        AnnotatedValue, BigUint, ManagedAddress, ManagedBuffer, ManagedVec, TxFrom, TxToSpecified,
    },
};

use super::{
    AnnotatedRewaPayment, FullPaymentData, FunctionCall, TxRewaValue, TxEnv, TxPayment,
    TxPaymentRewaOnly,
};

/// Indicates the REWA payment in a transaction.
pub struct Rewa<RewaValue>(pub RewaValue);

pub type RewaPayment<Api> = Rewa<BigUint<Api>>;

impl<RewaValue: Clone> Clone for Rewa<RewaValue> {
    fn clone(&self) -> Self {
        Rewa(self.0.clone())
    }
}

impl<Env, RewaValue> TxPayment<Env> for Rewa<RewaValue>
where
    Env: TxEnv,
    RewaValue: TxRewaValue<Env>,
{
    fn is_no_payment(&self, env: &Env) -> bool {
        self.0.with_value_ref(env, |rewa_value| rewa_value == &0u32)
    }

    fn perform_transfer_execute(
        self,
        env: &Env,
        to: &ManagedAddress<Env::Api>,
        gas_limit: u64,
        fc: FunctionCall<Env::Api>,
    ) {
        self.0.with_value_ref(env, |rewa_value| {
            let _ = SendRawWrapper::<Env::Api>::new().direct_rewa_execute(
                to,
                rewa_value,
                gas_limit,
                &fc.function_name,
                &fc.arg_buffer,
            );
        })
    }

    #[inline]
    fn with_normalized<From, To, F, R>(
        self,
        env: &Env,
        _from: &From,
        to: To,
        fc: FunctionCall<Env::Api>,
        f: F,
    ) -> R
    where
        From: TxFrom<Env>,
        To: TxToSpecified<Env>,
        F: FnOnce(&ManagedAddress<Env::Api>, &BigUint<Env::Api>, FunctionCall<Env::Api>) -> R,
    {
        to.with_address_ref(env, |to_addr| {
            self.0
                .with_value_ref(env, |rewa_value| f(to_addr, rewa_value, fc))
        })
    }

    fn into_full_payment_data(self, env: &Env) -> FullPaymentData<Env::Api> {
        FullPaymentData {
            rewa: Some(AnnotatedRewaPayment::new_rewa(self.0.into_value(env))),
            multi_dcdt: ManagedVec::new(),
        }
    }
}

impl<Env, RewaValue> AnnotatedValue<Env, BigUint<Env::Api>> for Rewa<RewaValue>
where
    Env: TxEnv,
    RewaValue: TxRewaValue<Env>,
{
    fn annotation(&self, env: &Env) -> ManagedBuffer<Env::Api> {
        self.0.annotation(env)
    }

    #[inline]
    fn to_value(&self, env: &Env) -> BigUint<Env::Api> {
        self.0.to_value(env)
    }

    #[inline]
    fn into_value(self, env: &Env) -> BigUint<Env::Api> {
        self.0.into_value(env)
    }

    #[inline]
    fn with_value_ref<F, R>(&self, env: &Env, f: F) -> R
    where
        F: FnOnce(&BigUint<Env::Api>) -> R,
    {
        self.0.with_value_ref(env, f)
    }
}

impl<Env, RewaValue> TxPaymentRewaOnly<Env> for Rewa<RewaValue>
where
    Env: TxEnv,
    RewaValue: TxRewaValue<Env>,
{
}
