use crate::{
    contract_base::BlockchainWrapper,
    types::{DcdtTokenPayment, ManagedVec, RHListItem, RHListItemExec, TxEnv},
};

/// Indicates that back-transfers will be returned.
pub struct ReturnsBackTransfersMultiDCDT;

impl<Env, Original> RHListItem<Env, Original> for ReturnsBackTransfersMultiDCDT
where
    Env: TxEnv,
{
    type Returns = ManagedVec<Env::Api, DcdtTokenPayment<Env::Api>>;
}

impl<RawResult, Env, Original> RHListItemExec<RawResult, Env, Original>
    for ReturnsBackTransfersMultiDCDT
where
    Env: TxEnv,
{
    fn item_process_result(self, _raw_result: &RawResult) -> Self::Returns {
        BlockchainWrapper::<Env::Api>::new()
            .get_back_transfers()
            .dcdt_payments
    }
}
