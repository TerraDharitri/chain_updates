dharitri_sc::imports!();
dharitri_sc::derive_imports!();

use crate::liquidity_pool::State;

pub const MAX_PERCENTAGE: u64 = 100_000;
pub const UNBOND_PERIOD: u64 = 10;

#[type_abi]
#[derive(
    TopEncode, TopDecode, NestedEncode, NestedDecode, Clone, PartialEq, Eq, Debug, ManagedVecItem,
)]
pub struct UnstakeTokenAttributes<M: ManagedTypeApi> {
    pub delegation_contract: ManagedAddress<M>,
    pub unstake_epoch: u64,
    pub unstake_amount: BigUint<M>,
    pub unbond_epoch: u64,
}

#[dharitri_sc::module]
pub trait ConfigModule:
    dharitri_sc_modules::default_issue_callbacks::DefaultIssueCallbacksModule
{
    #[only_owner]
    #[payable("REWA")]
    #[endpoint(registerLsToken)]
    fn register_ls_token(
        &self,
        token_display_name: ManagedBuffer,
        token_ticker: ManagedBuffer,
        num_decimals: usize,
    ) {
        let payment_amount = self.call_value().rewa().clone_value();
        self.ls_token().issue_and_set_all_roles(
            payment_amount,
            token_display_name,
            token_ticker,
            num_decimals,
            None,
        );
    }

    #[only_owner]
    #[payable("REWA")]
    #[endpoint(registerUnstakeToken)]
    fn register_unstake_token(
        &self,
        token_display_name: ManagedBuffer,
        token_ticker: ManagedBuffer,
        num_decimals: usize,
    ) {
        let payment_amount = self.call_value().rewa().clone_value();
        self.unstake_token().issue_and_set_all_roles(
            DcdtTokenType::NonFungible,
            payment_amount,
            token_display_name,
            token_ticker,
            num_decimals,
            None,
        );
    }

    #[view(getState)]
    #[storage_mapper("state")]
    fn state(&self) -> SingleValueMapper<State>;

    #[view(getLsTokenId)]
    #[storage_mapper("lsTokenId")]
    fn ls_token(&self) -> FungibleTokenMapper<Self::Api>;

    #[view(getLsSupply)]
    #[storage_mapper("lsTokenSupply")]
    fn ls_token_supply(&self) -> SingleValueMapper<BigUint>;

    #[view(getVirtualRewaReserve)]
    #[storage_mapper("virtualRewaReserve")]
    fn virtual_rewa_reserve(&self) -> SingleValueMapper<BigUint>;

    #[view(getRewardsReserve)]
    #[storage_mapper("rewardsReserve")]
    fn rewards_reserve(&self) -> SingleValueMapper<BigUint>;

    #[view(getUnstakeTokenId)]
    #[storage_mapper("unstakeTokenId")]
    fn unstake_token(&self) -> NonFungibleTokenMapper<Self::Api>;
}
