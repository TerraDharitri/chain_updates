dharitri_sc::imports!();

use crate::config;
use crate::errors::*;

#[dharitri_sc::module]
pub trait Validation: config::Config {
    fn require_is_accepted_payment(&self, payment: &DcdtTokenPayment<Self::Api>) {
        require!(
            self.governance_token_ids()
                .contains(&payment.token_identifier),
            UNREGISTERED_TOKEN_ID
        );
    }
}
