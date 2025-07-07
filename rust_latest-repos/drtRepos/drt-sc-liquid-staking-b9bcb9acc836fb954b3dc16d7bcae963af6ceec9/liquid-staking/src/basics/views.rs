dharitri_sc::imports!();

use crate::{basics::errors::ERROR_INSUFFICIENT_LIQ_BURNED, liquidity_pool};

#[dharitri_sc::module]
pub trait ViewsModule:
    crate::config::ConfigModule
    + dharitri_sc_modules::default_issue_callbacks::DefaultIssueCallbacksModule
    + liquidity_pool::LiquidityPoolModule
{
    // views
    #[view(getLsValueForPosition)]
    fn get_ls_value_for_position(&self, ls_token_amount: BigUint) -> BigUint {
        let ls_token_supply = self.ls_token_supply().get();
        let virtual_rewa_reserve = self.virtual_rewa_reserve().get();

        let rewa_amount = ls_token_amount * &virtual_rewa_reserve / ls_token_supply;
        require!(rewa_amount > 0u64, ERROR_INSUFFICIENT_LIQ_BURNED);

        rewa_amount
    }
}
