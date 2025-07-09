dharitri_sc::imports!();

use crate::config;

mod price_provider_proxy {
    dharitri_sc::imports!();

    #[dharitri_sc::proxy]
    pub trait PriceProvider {
        #[endpoint(getTokensForGivenPositionWithSafePrice)]
        fn get_tokens_for_given_position_with_safe_price(
            &self,
            liquidity: BigUint,
        ) -> MultiValue2<DcdtTokenPayment<Self::Api>, DcdtTokenPayment<Self::Api>>;
    }
}

#[dharitri_sc::module]
pub trait Lib: config::Config {
    fn get_vote_weight(&self, payment: &DcdtTokenPayment<Self::Api>) -> BigUint {
        let moa_token_id = self.moa_token_id().get();

        if payment.token_identifier == self.moa_token_id().get() {
            return payment.amount.clone();
        }

        if let Some(provider) = self.price_providers().get(&payment.token_identifier) {
            let call_result: MultiValue2<DcdtTokenPayment<Self::Api>, DcdtTokenPayment<Self::Api>> =
                self.price_provider_proxy(provider)
                    .get_tokens_for_given_position_with_safe_price(payment.amount.clone())
                    .execute_on_dest_context();
            let (token1, token2) = call_result.into_tuple();

            if token1.token_identifier == moa_token_id {
                return token1.amount;
            }

            if token2.token_identifier == moa_token_id {
                return token2.amount;
            }
        }

        BigUint::zero()
    }

    fn send_back(&self, payment: DcdtTokenPayment<Self::Api>) {
        self.send().direct_dcdt(
            &self.blockchain().get_caller(),
            &payment.token_identifier,
            payment.token_nonce,
            &payment.amount,
        );
    }

    #[proxy]
    fn price_provider_proxy(&self, to: ManagedAddress) -> price_provider_proxy::Proxy<Self::Api>;
}
