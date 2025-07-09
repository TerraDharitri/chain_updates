use crate::{
    api::ManagedTypeApi,
    types::{
        BigUint, DcdtTokenPayment, FullPaymentData, FunctionCall, ManagedAddress,
        TestTokenIdentifier, TxEnv, TxFrom, TxPayment, TxToSpecified,
    },
};

/// Syntactic sugar for quickly writing DCDT transfers in tests.
///
/// The fields are:
/// 1. token identifier, as TestTokenIdentifier
/// 2. nonce
/// 3. amount
///
/// The amount is represented as u64, since for most tests it is enough.
#[derive(Clone, Copy, Debug, PartialEq, Eq)]
pub struct TestDcdtTransfer<'a>(pub TestTokenIdentifier<'a>, pub u64, pub u64);

impl<'a, Api> From<TestDcdtTransfer<'a>> for DcdtTokenPayment<Api>
where
    Api: ManagedTypeApi,
{
    fn from(value: TestDcdtTransfer<'a>) -> Self {
        DcdtTokenPayment::new(value.0.to_token_identifier(), value.1, value.2.into())
    }
}

impl<Env> TxPayment<Env> for TestDcdtTransfer<'_>
where
    Env: TxEnv,
{
    #[inline]
    fn is_no_payment(&self, _env: &Env) -> bool {
        self.2 == 0
    }

    #[inline]
    fn perform_transfer_execute(
        self,
        env: &Env,
        to: &ManagedAddress<Env::Api>,
        gas_limit: u64,
        fc: FunctionCall<Env::Api>,
    ) {
        DcdtTokenPayment::from(self).perform_transfer_execute(env, to, gas_limit, fc);
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
        DcdtTokenPayment::from(self).with_normalized(env, from, to, fc, f)
    }

    fn into_full_payment_data(self, env: &Env) -> FullPaymentData<Env::Api> {
        DcdtTokenPayment::from(self).into_full_payment_data(env)
    }
}
