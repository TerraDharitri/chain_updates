#![no_std]

dharitri_sc::imports!();

#[dharitri_sc::contract]
pub trait SecondContract {
    #[init]
    fn init(&self, dcdt_token_identifier: RewaOrDcdtTokenIdentifier) {
        self.set_contract_dcdt_token_identifier(&dcdt_token_identifier);
    }

    #[payable("*")]
    #[endpoint(acceptDcdtPayment)]
    fn accept_dcdt_payment(&self) {
        let actual_token_identifier = self.call_value().rewa_or_single_dcdt().token_identifier;
        let expected_token_identifier = self.get_contract_dcdt_token_identifier();
        require!(
            actual_token_identifier == expected_token_identifier,
            "Wrong dcdt token"
        );
    }

    #[payable("*")]
    #[endpoint(rejectDcdtPayment)]
    fn reject_dcdt_payment(&self) {
        sc_panic!("Rejected")
    }

    // storage

    #[storage_set("dcdtTokenName")]
    fn set_contract_dcdt_token_identifier(&self, dcdt_token_identifier: &RewaOrDcdtTokenIdentifier);

    #[view(getdcdtTokenName)]
    #[storage_get("dcdtTokenName")]
    fn get_contract_dcdt_token_identifier(&self) -> RewaOrDcdtTokenIdentifier;
}
