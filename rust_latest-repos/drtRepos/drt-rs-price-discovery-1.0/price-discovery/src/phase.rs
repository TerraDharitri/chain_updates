use crate::Timestamp;

dharitri_sc::imports!();
dharitri_sc::derive_imports!();

#[derive(TypeAbi, TopEncode, TopDecode, NestedEncode, PartialEq, PartialOrd)]
pub enum Phase {
    Idle,
    UserDepositWithdraw,
    OwnerDepositWithdraw,
    OwnerRedeem,
    UserRedeem,
}

pub const MAX_PHASE_DURATION: Timestamp = 60 * 60 * 24 * 30; // ~ 1 month
pub const MAX_LONG_PHASE_DURATION: Timestamp = MAX_PHASE_DURATION * 6; // ~ 6 months

#[dharitri_sc::module]
pub trait PhaseModule:
    crate::common_storage::CommonStorageModule + crate::events::EventsModule
{
    #[view(getCurrentPhase)]
    fn get_current_phase(&self) -> Phase {
        let current_time = self.blockchain().get_block_timestamp();
        let start_time = self.start_time().get();
        if current_time < start_time {
            return Phase::Idle;
        }

        let user_deposit_time = self.user_deposit_withdraw_time().get();
        let user_deposit_phase_end = start_time + user_deposit_time;
        if current_time < user_deposit_phase_end {
            return Phase::UserDepositWithdraw;
        }

        let owner_deposit_time = self.owner_deposit_withdraw_time().get();
        let owner_deposit_phase_end = user_deposit_phase_end + owner_deposit_time;
        if current_time < owner_deposit_phase_end {
            return Phase::OwnerDepositWithdraw;
        }

        let owner_redeem_time = self.owner_redeem_time().get();
        let owner_redeem_phase_end = owner_deposit_phase_end + owner_redeem_time;
        if current_time < owner_redeem_phase_end {
            return Phase::OwnerRedeem;
        }

        Phase::UserRedeem
    }

    fn require_user_deposit_withdraw_allowed(&self, phase: &Phase) {
        require!(
            phase == &Phase::UserDepositWithdraw,
            "User deposit/withdraw not allowed in this phase"
        );
    }

    fn require_owner_deposit_withdraw_allowed(&self, phase: &Phase) {
        require!(
            phase == &Phase::OwnerDepositWithdraw,
            "Owner deposit/withdraw not allowed in this phase"
        );
    }

    fn require_owner_redeem_allowed(&self, phase: &Phase) {
        require!(
            phase == &Phase::OwnerRedeem,
            "Owner redeem not allowed in this phase"
        );
    }

    fn require_user_redeem_allowed(&self, phase: &Phase) {
        require!(
            phase == &Phase::UserRedeem,
            "User redeem not allowed in this phase"
        );
    }

    fn require_before_redeem(&self, phase: &Phase) {
        require!(
            phase < &Phase::OwnerRedeem,
            "May only call this endpoint before redeem phase"
        );
    }

    #[view(getUserDepositWithdrawTime)]
    #[storage_mapper("userDepositWithdrawTime")]
    fn user_deposit_withdraw_time(&self) -> SingleValueMapper<Timestamp>;

    #[view(getOwnerDepositWithdrawTime)]
    #[storage_mapper("ownerDepositWithdrawTime")]
    fn owner_deposit_withdraw_time(&self) -> SingleValueMapper<Timestamp>;

    #[view(getOwnerRedeemTime)]
    #[storage_mapper("ownerRedeemTime")]
    fn owner_redeem_time(&self) -> SingleValueMapper<Timestamp>;
}
