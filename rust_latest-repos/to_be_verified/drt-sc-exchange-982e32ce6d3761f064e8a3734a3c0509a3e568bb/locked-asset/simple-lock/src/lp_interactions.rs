dharitri_sc::imports!();
dharitri_sc::derive_imports!();

use crate::error_messages::*;

type AddLiquidityResultType<M> =
    MultiValue3<DcdtTokenPayment<M>, DcdtTokenPayment<M>, DcdtTokenPayment<M>>;
type RemoveLiquidityResultType<BigUint> =
    MultiValue2<DcdtTokenPayment<BigUint>, DcdtTokenPayment<BigUint>>;

pub struct AddLiquidityResultWrapper<M: ManagedTypeApi> {
    pub lp_tokens: DcdtTokenPayment<M>,
    pub first_token_refund: DcdtTokenPayment<M>,
    pub second_token_refund: DcdtTokenPayment<M>,
}

pub struct RemoveLiquidityResultWrapper<M: ManagedTypeApi> {
    pub first_token_payment_out: DcdtTokenPayment<M>,
    pub second_token_payment_out: DcdtTokenPayment<M>,
}

// Must manually declare, as Pair SC already depends on simple-lock
// This avoids circular dependency
mod lp_proxy {
    dharitri_sc::imports!();
    use super::{AddLiquidityResultType, RemoveLiquidityResultType};

    #[dharitri_sc::proxy]
    pub trait LpProxy {
        #[payable("*")]
        #[endpoint(addLiquidity)]
        fn add_liquidity(
            &self,
            first_token_amount_min: BigUint,
            second_token_amount_min: BigUint,
        ) -> AddLiquidityResultType<Self::Api>;

        #[payable("*")]
        #[endpoint(removeLiquidity)]
        fn remove_liquidity(
            &self,
            first_token_amount_min: BigUint,
            second_token_amount_min: BigUint,
        ) -> RemoveLiquidityResultType<Self::Api>;
    }
}

#[dharitri_sc::module]
pub trait LpInteractionsModule {
    fn call_pair_add_liquidity(
        &self,
        lp_address: ManagedAddress,
        first_payment: &DcdtTokenPayment<Self::Api>,
        second_payment: &DcdtTokenPayment<Self::Api>,
        first_token_amount_min: BigUint,
        second_token_amount_min: BigUint,
    ) -> AddLiquidityResultWrapper<Self::Api> {
        let mut lp_payments_in = ManagedVec::new();
        lp_payments_in.push(first_payment.clone());
        lp_payments_in.push(second_payment.clone());

        let lp_payments_out: AddLiquidityResultType<Self::Api> = self
            .lp_proxy(lp_address)
            .add_liquidity(first_token_amount_min, second_token_amount_min)
            .with_multi_token_transfer(lp_payments_in)
            .execute_on_dest_context();
        let (lp_tokens, first_token_optimal_payment, second_token_optimal_payment) =
            lp_payments_out.into_tuple();

        require!(
            first_payment.token_identifier == first_token_optimal_payment.token_identifier
                && second_payment.token_identifier == second_token_optimal_payment.token_identifier,
            INVALID_PAYMENTS_RECEIVED_FROM_LP_ERR_MSG
        );

        let first_refund_amount = &first_payment.amount - &first_token_optimal_payment.amount;
        let first_token_refund = DcdtTokenPayment::new(
            first_token_optimal_payment.token_identifier,
            0,
            first_refund_amount,
        );

        let second_refund_amount = &second_payment.amount - &second_token_optimal_payment.amount;
        let second_token_refund = DcdtTokenPayment::new(
            second_token_optimal_payment.token_identifier,
            0,
            second_refund_amount,
        );

        AddLiquidityResultWrapper {
            lp_tokens,
            first_token_refund,
            second_token_refund,
        }
    }

    fn call_pair_remove_liquidity(
        &self,
        lp_address: ManagedAddress,
        lp_token_id: TokenIdentifier,
        lp_token_amount: BigUint,
        first_token_amount_min: BigUint,
        second_token_amount_min: BigUint,
        expected_first_token_id_out: &TokenIdentifier,
        expected_second_token_id_out: &TokenIdentifier,
    ) -> RemoveLiquidityResultWrapper<Self::Api> {
        let lp_payments_out: RemoveLiquidityResultType<Self::Api> = self
            .lp_proxy(lp_address)
            .remove_liquidity(first_token_amount_min, second_token_amount_min)
            .with_dcdt_transfer((lp_token_id, 0, lp_token_amount))
            .execute_on_dest_context();

        let (first_token_payment_out, second_token_payment_out) = lp_payments_out.into_tuple();
        require!(
            &first_token_payment_out.token_identifier == expected_first_token_id_out
                && &second_token_payment_out.token_identifier == expected_second_token_id_out,
            INVALID_PAYMENTS_RECEIVED_FROM_LP_ERR_MSG
        );

        RemoveLiquidityResultWrapper {
            first_token_payment_out,
            second_token_payment_out,
        }
    }

    #[proxy]
    fn lp_proxy(&self, sc_address: ManagedAddress) -> lp_proxy::Proxy<Self::Api>;
}
