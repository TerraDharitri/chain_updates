dharitri_sc::imports!();
dharitri_sc::derive_imports!();

static DCDT_TRANSFER_FUNC_NAME: &[u8] = b"DCDTTransfer";

const GAS_LIMIT_DCDT_TRANSFER: u64 = 500_000;
const CALLBACK_DCDT_TRANSFER_GAS_LIMIT: u64 = 100_000;

#[derive(TopEncode, TopDecode)]
pub enum TransferResult {
    None,
    Success,
    Fail,
}

#[dharitri_sc::module]
pub trait DcdtFeaturesModule {
    #[endpoint(transferFungiblePromiseNoCallback)]
    fn transfer_fungible_promise_no_callback(&self, to: ManagedAddress, amount: BigUint) {
        let token_id = self.fungible_dcdt_token_id().get_token_id();
        self.tx()
            .to(to)
            .raw_call(DCDT_TRANSFER_FUNC_NAME)
            .argument(&token_id)
            .argument(&amount)
            .gas(GAS_LIMIT_DCDT_TRANSFER)
            .register_promise();
    }

    #[endpoint(transferFungiblePromiseWithCallback)]
    fn transfer_fungible_promise_with_callback(&self, to: ManagedAddress, amount: BigUint) {
        let token_id = self.fungible_dcdt_token_id().get_token_id();
        self.tx()
            .to(to)
            .raw_call(DCDT_TRANSFER_FUNC_NAME)
            .argument(&token_id)
            .argument(&amount)
            .gas(GAS_LIMIT_DCDT_TRANSFER)
            .callback(self.callbacks().transfer_callback())
            .gas_for_callback(CALLBACK_DCDT_TRANSFER_GAS_LIMIT)
            .register_promise();
    }

    #[promises_callback]
    fn transfer_callback(&self, #[call_result] result: ManagedAsyncCallResult<()>) {
        match result {
            ManagedAsyncCallResult::Ok(()) => {
                self.latest_transfer_result().set(TransferResult::Success);
            },
            ManagedAsyncCallResult::Err(_) => {
                self.latest_transfer_result().set(TransferResult::Fail);
            },
        }
    }

    #[storage_mapper("fungibleDcdtTokenId")]
    fn fungible_dcdt_token_id(&self) -> FungibleTokenMapper;

    #[storage_mapper("nonFungibleDcdtTokenId")]
    fn non_fungible_dcdt_token_id(&self) -> NonFungibleTokenMapper;

    #[storage_mapper("latestTransferResult")]
    fn latest_transfer_result(&self) -> SingleValueMapper<TransferResult>;
}
