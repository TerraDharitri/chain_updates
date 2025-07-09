dharitri_sc::imports!();
dharitri_sc::derive_imports!();

use crate::basics::errors::{ERROR_INSUFFICIENT_LIQUIDITY, ERROR_INSUFFICIENT_LIQ_BURNED};
use crate::contexts::base::StorageCache;

use crate::config;

#[type_abi]
#[derive(TopEncode, TopDecode, PartialEq, Eq, Copy, Clone, Debug)]
pub enum State {
    Inactive,
    Active,
}

#[dharitri_sc::module]
pub trait LiquidityPoolModule:
    config::ConfigModule + dharitri_sc_modules::default_issue_callbacks::DefaultIssueCallbacksModule
{
    fn pool_add_liquidity(
        &self,
        token_amount: &BigUint,
        storage_cache: &mut StorageCache<Self>,
    ) -> BigUint {
        let ls_amount = if storage_cache.virtual_rewa_reserve > 0 {
            token_amount.clone() * &storage_cache.ls_token_supply
                / &storage_cache.virtual_rewa_reserve
        } else {
            token_amount.clone()
        };

        require!(ls_amount > 0, ERROR_INSUFFICIENT_LIQUIDITY);

        storage_cache.ls_token_supply += &ls_amount;
        storage_cache.virtual_rewa_reserve += token_amount;

        ls_amount
    }

    fn pool_remove_liquidity(
        &self,
        token_amount: &BigUint,
        storage_cache: &mut StorageCache<Self>,
    ) -> BigUint {
        let rewa_amount = self.get_rewa_amount(token_amount, storage_cache);
        storage_cache.ls_token_supply -= token_amount;
        storage_cache.virtual_rewa_reserve -= &rewa_amount;

        rewa_amount
    }

    fn get_rewa_amount(
        &self,
        ls_token_amount: &BigUint,
        storage_cache: &StorageCache<Self>,
    ) -> BigUint {
        let rewa_amount =
            ls_token_amount * &storage_cache.virtual_rewa_reserve / &storage_cache.ls_token_supply;
        require!(rewa_amount > 0u64, ERROR_INSUFFICIENT_LIQ_BURNED);

        rewa_amount
    }

    fn mint_ls_token(&self, amount: BigUint) -> DcdtTokenPayment<Self::Api> {
        self.ls_token().mint(amount)
    }

    fn burn_ls_token(&self, amount: &BigUint) {
        self.ls_token().burn(amount);
    }

    fn mint_unstake_tokens<T: TopEncode>(&self, attributes: &T) -> DcdtTokenPayment<Self::Api> {
        self.unstake_token()
            .nft_create(BigUint::from(1u64), attributes)
    }

    fn burn_unstake_tokens(&self, token_nonce: u64) {
        self.unstake_token()
            .nft_burn(token_nonce, &BigUint::from(1u64));
    }
}
