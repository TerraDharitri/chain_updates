use crate::{
    api::ManagedTypeApi,
    types::{BigUint, MultiDcdtPayment},
};

/// Holding back-transfer data, as retrieved from the VM.
#[derive(Clone)]
pub struct BackTransfers<A>
where
    A: ManagedTypeApi,
{
    pub total_rewa_amount: BigUint<A>,
    pub dcdt_payments: MultiDcdtPayment<A>,
}
