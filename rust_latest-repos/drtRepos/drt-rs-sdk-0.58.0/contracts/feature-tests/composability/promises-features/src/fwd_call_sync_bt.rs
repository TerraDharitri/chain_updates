use crate::vault_proxy;

dharitri_sc::imports!();

/// Not directly related to promises, but this contract already has the setup for VM 1.5.
#[dharitri_sc::module]
pub trait BackTransfersFeatureModule {
    #[endpoint]
    fn forward_sync_retrieve_funds_bt(
        &self,
        to: ManagedAddress,
        token: RewaOrDcdtTokenIdentifier,
        token_nonce: u64,
        amount: BigUint,
    ) {
        let back_transfers = self
            .tx()
            .to(&to)
            .typed(vault_proxy::VaultProxy)
            .retrieve_funds(token, token_nonce, amount)
            .returns(ReturnsBackTransfers)
            .sync_call();

        require!(
            back_transfers.dcdt_payments.len() == 1 || back_transfers.total_rewa_amount != 0,
            "Only one DCDT payment expected"
        );

        self.back_transfers_event(
            &back_transfers.total_rewa_amount,
            &back_transfers.dcdt_payments.into_multi_value(),
        );
    }

    #[endpoint]
    fn forward_sync_retrieve_funds_bt_reset_twice(
        &self,
        to: ManagedAddress,
        token: RewaOrDcdtTokenIdentifier,
        token_nonce: u64,
        amount: BigUint,
    ) {
        let back_transfers = self
            .tx()
            .to(&to)
            .typed(vault_proxy::VaultProxy)
            .retrieve_funds(token.clone(), token_nonce, amount.clone())
            .returns(ReturnsBackTransfersReset)
            .sync_call();

        require!(
            back_transfers.dcdt_payments.len() == 1 || back_transfers.total_rewa_amount != 0,
            "Only one DCDT payment expected"
        );

        self.back_transfers_event(
            &back_transfers.total_rewa_amount,
            &back_transfers.dcdt_payments.into_multi_value(),
        );

        let back_transfers = self
            .tx()
            .to(&to)
            .typed(vault_proxy::VaultProxy)
            .retrieve_funds(token, token_nonce, amount)
            .returns(ReturnsBackTransfersReset)
            .sync_call();

        require!(
            back_transfers.dcdt_payments.len() == 1 || back_transfers.total_rewa_amount != 0,
            "Only one DCDT payment expected"
        );

        self.back_transfers_event(
            &back_transfers.total_rewa_amount,
            &back_transfers.dcdt_payments.into_multi_value(),
        );
    }

    #[endpoint]
    fn forward_sync_retrieve_funds_bt_twice(
        &self,
        to: ManagedAddress,
        token: RewaOrDcdtTokenIdentifier,
        token_nonce: u64,
        amount: BigUint,
    ) {
        let back_transfers = self
            .tx()
            .to(&to)
            .typed(vault_proxy::VaultProxy)
            .retrieve_funds(token.clone(), token_nonce, amount.clone())
            .returns(ReturnsBackTransfers)
            .sync_call();

        require!(
            back_transfers.dcdt_payments.len() == 1 || back_transfers.total_rewa_amount != 0,
            "Only one DCDT payment expected"
        );

        self.back_transfers_event(
            &back_transfers.total_rewa_amount,
            &back_transfers.dcdt_payments.into_multi_value(),
        );

        let back_transfers = self
            .tx()
            .to(&to)
            .typed(vault_proxy::VaultProxy)
            .retrieve_funds(token, token_nonce, amount)
            .returns(ReturnsBackTransfers)
            .sync_call();

        require!(
            back_transfers.dcdt_payments.len() == 1 || back_transfers.total_rewa_amount != 0,
            "Only one DCDT payment expected"
        );

        self.back_transfers_event(
            &back_transfers.total_rewa_amount,
            &back_transfers.dcdt_payments.into_multi_value(),
        );
    }

    #[event("back_transfers")]
    fn back_transfers_event(
        &self,
        #[indexed] rewa_value: &BigUint,
        #[indexed] multi_dcdt: &MultiValueEncoded<DcdtTokenPaymentMultiValue>,
    );
}
