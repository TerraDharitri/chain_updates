use core::ops::Deref;

use crate::{
    contract_base::SendRawWrapper,
    types::{BigUint, ManagedAddress, ManagedRef, MultiDcdtPayment, TxFrom, TxToSpecified},
};

use super::{FullPaymentData, FunctionCall, TxEnv, TxPayment};

/// Indicates that a payment object contains a multi-DCDT payment.
pub trait TxPaymentMultiDcdt<Env>: TxPayment<Env>
where
    Env: TxEnv,
{
}

impl<Env> TxPaymentMultiDcdt<Env> for MultiDcdtPayment<Env::Api> where Env: TxEnv {}
impl<Env> TxPaymentMultiDcdt<Env> for &MultiDcdtPayment<Env::Api> where Env: TxEnv {}
impl<Env> TxPaymentMultiDcdt<Env> for ManagedRef<'_, Env::Api, MultiDcdtPayment<Env::Api>> where
    Env: TxEnv
{
}

impl<Env> TxPayment<Env> for &MultiDcdtPayment<Env::Api>
where
    Env: TxEnv,
{
    fn is_no_payment(&self, _env: &Env) -> bool {
        self.is_empty()
    }

    fn perform_transfer_execute(
        self,
        _env: &Env,
        to: &ManagedAddress<Env::Api>,
        gas_limit: u64,
        fc: FunctionCall<Env::Api>,
    ) {
        let _ = SendRawWrapper::<Env::Api>::new().multi_dcdt_transfer_execute(
            to,
            self,
            gas_limit,
            &fc.function_name,
            &fc.arg_buffer,
        );
    }

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
        F: FnOnce(&ManagedAddress<Env::Api>, &BigUint<Env::Api>, FunctionCall<Env::Api>) -> R,
    {
        match self.len() {
            0 => ().with_normalized(env, from, to, fc, f),
            1 => self.get(0).as_refs().with_normalized(env, from, to, fc, f),
            _ => to.with_address_ref(env, |to_addr| {
                let fc_conv = fc.convert_to_multi_transfer_dcdt_call(
                    to_addr,
                    self.as_multi_rewa_or_dcdt_payment(),
                );
                f(&from.resolve_address(env), &*BigUint::zero_ref(), fc_conv)
            }),
        }
    }

    fn into_full_payment_data(self, _env: &Env) -> FullPaymentData<Env::Api> {
        FullPaymentData {
            rewa: None,
            multi_dcdt: self.as_multi_rewa_or_dcdt_payment().clone(),
        }
    }
}

impl<Env> TxPayment<Env> for ManagedRef<'_, Env::Api, MultiDcdtPayment<Env::Api>>
where
    Env: TxEnv,
{
    #[inline]
    fn is_no_payment(&self, _env: &Env) -> bool {
        self.deref().is_empty()
    }

    #[inline]
    fn perform_transfer_execute(
        self,
        env: &Env,
        to: &ManagedAddress<Env::Api>,
        gas_limit: u64,
        fc: FunctionCall<Env::Api>,
    ) {
        self.deref()
            .perform_transfer_execute(env, to, gas_limit, fc)
    }

    #[inline]
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
        F: FnOnce(&ManagedAddress<Env::Api>, &BigUint<Env::Api>, FunctionCall<Env::Api>) -> R,
    {
        self.deref().with_normalized(env, from, to, fc, f)
    }

    fn into_full_payment_data(self, env: &Env) -> FullPaymentData<Env::Api> {
        self.deref().into_full_payment_data(env)
    }
}

impl<Env> TxPayment<Env> for MultiDcdtPayment<Env::Api>
where
    Env: TxEnv,
{
    #[inline]
    fn is_no_payment(&self, _env: &Env) -> bool {
        self.is_empty()
    }

    #[inline]
    fn perform_transfer_execute(
        self,
        env: &Env,
        to: &ManagedAddress<Env::Api>,
        gas_limit: u64,
        fc: FunctionCall<Env::Api>,
    ) {
        (&self).perform_transfer_execute(env, to, gas_limit, fc);
    }

    #[inline]
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
        F: FnOnce(&ManagedAddress<Env::Api>, &BigUint<Env::Api>, FunctionCall<Env::Api>) -> R,
    {
        (&self).with_normalized(env, from, to, fc, f)
    }

    fn into_full_payment_data(self, _env: &Env) -> FullPaymentData<Env::Api> {
        FullPaymentData {
            rewa: None,
            multi_dcdt: self.into_multi_rewa_or_dcdt_payment(),
        }
    }
}
