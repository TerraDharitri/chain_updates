use crate::Timestamp;

dharitri_sc::imports!();

#[dharitri_sc::module]
pub trait CommonStorageModule {
    #[storage_mapper("launchedTokenId")]
    fn launched_token_id(&self) -> SingleValueMapper<TokenIdentifier>;

    #[storage_mapper("launchedTokenBalance")]
    fn launched_token_balance(&self) -> SingleValueMapper<BigUint>;

    #[storage_mapper("minLaunchedTokens")]
    fn min_launched_tokens(&self) -> SingleValueMapper<BigUint>;

    #[storage_mapper("acceptedTokenId")]
    fn accepted_token_id(&self) -> SingleValueMapper<RewaOrDcdtTokenIdentifier>;

    #[storage_mapper("acceptedTokenBalance")]
    fn accepted_token_balance(&self) -> SingleValueMapper<BigUint>;

    #[view(getStartTime)]
    #[storage_mapper("startTime")]
    fn start_time(&self) -> SingleValueMapper<Timestamp>;

    #[storage_mapper("pricePrecision")]
    fn price_precision(&self) -> SingleValueMapper<u64>;
}
