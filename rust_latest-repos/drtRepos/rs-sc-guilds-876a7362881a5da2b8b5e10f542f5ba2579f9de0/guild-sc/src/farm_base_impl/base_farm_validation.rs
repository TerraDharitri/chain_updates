dharitri_sc::imports!();

use common_errors::ERROR_NO_FARM_TOKEN;

#[dharitri_sc::module]
pub trait BaseFarmValidationModule {
    fn require_valid_farm_token_id(&self, farm_token_id: &TokenIdentifier) {
        require!(
            farm_token_id.is_valid_dcdt_identifier(),
            ERROR_NO_FARM_TOKEN
        );
    }
}
