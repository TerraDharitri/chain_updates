use crate::{
    api::{ErrorApi, ErrorApiImpl},
    contract_base::BlockchainWrapper,
    types::{DcdtTokenPayment, RHListItem, RHListItemExec, TxEnv},
};

/// Indicates that back-transfers will be returned.
pub struct ReturnsBackTransfersSingleDCDT;

impl<Env, Original> RHListItem<Env, Original> for ReturnsBackTransfersSingleDCDT
where
    Env: TxEnv,
{
    type Returns = DcdtTokenPayment<Env::Api>;
}

impl<RawResult, Env, Original> RHListItemExec<RawResult, Env, Original>
    for ReturnsBackTransfersSingleDCDT
where
    Env: TxEnv,
{
    fn item_process_result(self, _raw_result: &RawResult) -> Self::Returns {
        let dcdt_payments = BlockchainWrapper::<Env::Api>::new()
            .get_back_transfers()
            .dcdt_payments;

        if dcdt_payments.len() != 1 {
            Env::Api::error_api_impl().signal_error(b"Back transfers expected to be a single DCDT")
        }

        let x = dcdt_payments.get(0).clone();
        x
    }
}
