dharitri_sc::imports!();

use common_structs::Epoch;
use simple_lock::locked_token::LockedTokenAttributes;

use crate::energy::Energy;

static INVALID_EXTEND_PERIOD_ARG_ERR_MSG: &[u8] =
    b"New lock period must be longer than the current one";

#[dharitri_sc::module]
pub trait ExtendLockModule:
    simple_lock::basic_lock_unlock::BasicLockUnlock
    + simple_lock::locked_token::LockedTokenModule
    + simple_lock::token_attributes::TokenAttributesModule
    + dharitri_sc_modules::default_issue_callbacks::DefaultIssueCallbacksModule
    + crate::token_whitelist::TokenWhitelistModule
    + crate::energy::EnergyModule
    + crate::lock_options::LockOptionsModule
    + crate::events::EventsModule
    + crate::migration::SimpleLockMigrationModule
    + dharitri_sc_modules::pause::PauseModule
    + utils::UtilsModule
    + legacy_token_decode_module::LegacyTokenDecodeModule
{
    fn lock_by_token_type(
        &self,
        dest_address: &ManagedAddress,
        payment: DcdtTokenPayment,
        unlock_epoch: Epoch,
        current_epoch: Epoch,
    ) -> DcdtTokenPayment {
        let output_payment = self.update_energy(dest_address, |energy: &mut Energy<Self::Api>| {
            let payment_clone = payment.clone();
            if self.is_base_asset_token(&payment.token_identifier) {
                self.lock_base_asset(payment_clone, unlock_epoch, current_epoch, energy)
            } else {
                self.require_address_is_caller(dest_address);
                self.locked_token()
                    .require_same_token(&payment.token_identifier);

                self.extend_new_token_period(payment_clone, unlock_epoch, current_epoch, energy)
            }
        });

        self.send().dcdt_local_burn(
            &payment.token_identifier,
            payment.token_nonce,
            &payment.amount,
        );

        output_payment
    }

    fn lock_base_asset(
        &self,
        payment: DcdtTokenPayment,
        unlock_epoch: Epoch,
        current_epoch: Epoch,
        energy: &mut Energy<Self::Api>,
    ) -> DcdtTokenPayment {
        let output_tokens = self.lock_tokens(payment.into(), unlock_epoch);
        energy.add_after_token_lock(&output_tokens.amount, unlock_epoch, current_epoch);

        self.to_dcdt_payment(output_tokens)
    }

    fn extend_new_token_period(
        &self,
        payment: DcdtTokenPayment,
        new_unlock_epoch: Epoch,
        current_epoch: Epoch,
        energy: &mut Energy<Self::Api>,
    ) -> DcdtTokenPayment {
        let attributes: LockedTokenAttributes<Self::Api> = self
            .locked_token()
            .get_token_attributes(payment.token_nonce);

        require!(
            new_unlock_epoch > attributes.unlock_epoch,
            INVALID_EXTEND_PERIOD_ARG_ERR_MSG
        );

        energy.update_after_unlock_epoch_change(
            &payment.amount,
            attributes.unlock_epoch,
            new_unlock_epoch,
            current_epoch,
        );

        let unlocked_tokens = self.unlock_tokens_unchecked(payment, &attributes);
        let output_tokens = self.lock_tokens(unlocked_tokens, new_unlock_epoch);

        self.to_dcdt_payment(output_tokens)
    }

    fn require_address_is_caller(&self, address: &ManagedAddress) {
        let caller = self.blockchain().get_caller();
        require!(
            address == &caller,
            "May not use the optional destination argument here"
        );
    }
}
