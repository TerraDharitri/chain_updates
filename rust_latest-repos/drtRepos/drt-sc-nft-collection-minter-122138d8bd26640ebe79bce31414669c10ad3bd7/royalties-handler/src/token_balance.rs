use nft_minter::common_storage::RewaValuePaymentsVecPair;

dharitri_sc::imports!();

#[dharitri_sc::module]
pub trait TokenBalanceModule {
    fn add_balance(&self, token: RewaOrDcdtTokenIdentifier, amount: &BigUint) {
        self.balance_for_token(&token).update(|b| {
            *b += amount;
        });
        let _ = self.known_tokens().insert(token);
    }

    fn update_balance_from_results(&self, result: RewaValuePaymentsVecPair<Self::Api>) {
        let (rewa_value, other_payments) = result.into_tuple();

        if rewa_value > 0 {
            self.add_balance(RewaOrDcdtTokenIdentifier::rewa(), &rewa_value);
        }
        for p in &other_payments {
            self.add_balance(
                RewaOrDcdtTokenIdentifier::dcdt(p.token_identifier),
                &p.amount,
            );
        }
    }

    #[view(getTokenBalances)]
    fn get_token_balances(
        &self,
    ) -> MultiValueEncoded<MultiValue2<RewaOrDcdtTokenIdentifier, BigUint>> {
        let mut balances = MultiValueEncoded::new();

        for token_id in self.known_tokens().iter() {
            let balance_for_token = self.balance_for_token(&token_id).get();
            if balance_for_token > 0 {
                balances.push((token_id, balance_for_token).into());
            }
        }

        balances
    }

    #[storage_mapper("knownTokens")]
    fn known_tokens(&self) -> UnorderedSetMapper<RewaOrDcdtTokenIdentifier>;

    #[storage_mapper("balanceForToken")]
    fn balance_for_token(&self, token_id: &RewaOrDcdtTokenIdentifier)
        -> SingleValueMapper<BigUint>;
}
