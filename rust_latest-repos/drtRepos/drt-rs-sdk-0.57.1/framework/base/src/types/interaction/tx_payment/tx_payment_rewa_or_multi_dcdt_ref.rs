use crate::types::{BigUint, RewaOrMultiDcdtPaymentRefs, ManagedAddress, TxFrom, TxToSpecified};

use super::{Rewa, FullPaymentData, FunctionCall, TxEnv, TxPayment};

impl<Env> TxPayment<Env> for RewaOrMultiDcdtPaymentRefs<'_, Env::Api>
where
    Env: TxEnv,
{
    fn is_no_payment(&self, _env: &Env) -> bool {
        self.is_empty()
    }

    fn perform_transfer_execute(
        self,
        env: &Env,
        to: &ManagedAddress<Env::Api>,
        gas_limit: u64,
        fc: FunctionCall<Env::Api>,
    ) {
        match self {
            RewaOrMultiDcdtPaymentRefs::Rewa(rewa_amount) => {
                Rewa(rewa_amount).perform_transfer_execute(env, to, gas_limit, fc);
            },
            RewaOrMultiDcdtPaymentRefs::MultiDcdt(multi_dcdt_payment) => {
                multi_dcdt_payment.perform_transfer_execute(env, to, gas_limit, fc);
            },
        }
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
        match self {
            RewaOrMultiDcdtPaymentRefs::Rewa(rewa_amount) => {
                Rewa(rewa_amount).with_normalized(env, from, to, fc, f)
            },
            RewaOrMultiDcdtPaymentRefs::MultiDcdt(multi_dcdt_payment) => {
                multi_dcdt_payment.with_normalized(env, from, to, fc, f)
            },
        }
    }

    fn into_full_payment_data(self, env: &Env) -> FullPaymentData<Env::Api> {
        match self {
            RewaOrMultiDcdtPaymentRefs::Rewa(rewa_amount) => {
                Rewa(rewa_amount).into_full_payment_data(env)
            },
            RewaOrMultiDcdtPaymentRefs::MultiDcdt(multi_dcdt_payment) => {
                multi_dcdt_payment.into_full_payment_data(env)
            },
        }
    }
}
