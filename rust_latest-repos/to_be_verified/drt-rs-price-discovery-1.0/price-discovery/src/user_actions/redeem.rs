dharitri_sc::imports!();

pub type UserRedeemFlag = bool;
pub const USER_REDEEMED: UserRedeemFlag = true;

#[dharitri_sc::module]
pub trait RedeemModule:
    super::user_deposit_withdraw::UserDepositWithdrawModule
    + crate::common_storage::CommonStorageModule
    + crate::events::EventsModule
    + crate::phase::PhaseModule
    + dharitri_sc_modules::default_issue_callbacks::DefaultIssueCallbacksModule
{
    /// After all phases have ended,
    /// users can withdraw their fair share of launched tokens.
    #[endpoint(userRedeem)]
    fn user_redeem_endpoint(&self) -> RewaOrDcdtTokenPayment {
        let phase = self.get_current_phase();
        self.require_user_redeem_allowed(&phase);

        let caller = self.blockchain().get_caller();
        let bought_tokens = self.user_redeem(&caller);
        self.emit_redeem_event(&bought_tokens.token_identifier, &bought_tokens.amount);

        bought_tokens
    }

    /// After the OwnerDepositWithdraw phase has ended,
    /// the owner can withdraw the accepted tokens.
    #[only_owner]
    #[endpoint(ownerRedeem)]
    fn owner_redeem_endpoint(&self) -> RewaOrDcdtTokenPayment {
        let phase = self.get_current_phase();
        self.require_owner_redeem_allowed(&phase);

        let caller = self.blockchain().get_caller();
        let redeemed_tokens = self.owner_redeem(&caller);
        self.emit_redeem_event(&redeemed_tokens.token_identifier, &redeemed_tokens.amount);

        redeemed_tokens
    }

    /// Only to be used in the cases where the owner somehow missed the long owner redeem phase
    #[only_owner]
    #[endpoint(withdrawLaunchpadTokens)]
    fn withdraw_launchpad_tokens(&self) {
        let phase = self.get_current_phase();
        self.require_user_redeem_allowed(&phase);
        self.require_owner_didnt_redeem();

        let launched_token_id = self.launched_token_id().get();
        let launched_tokens_supply = self.launched_token_balance().take();
        let owner = self.blockchain().get_caller();
        self.send()
            .direct_dcdt(&owner, &launched_token_id, 0, &launched_tokens_supply);

        self.owner_redeemed().set(USER_REDEEMED);
    }

    fn owner_redeem(&self, owner: &ManagedAddress) -> RewaOrDcdtTokenPayment {
        self.require_owner_didnt_redeem();

        let launched_token_supply = self.launched_token_balance().get();
        require!(
            launched_token_supply > 0,
            "May not withdraw tokens as launched tokens were not deposited"
        );

        let accepted_token_id = self.accepted_token_id().get();
        let accepted_token_balance = self.accepted_token_balance().get();
        self.send()
            .direct(owner, &accepted_token_id, 0, &accepted_token_balance);

        self.owner_redeemed().set(USER_REDEEMED);

        RewaOrDcdtTokenPayment::new(accepted_token_id, 0, accepted_token_balance)
    }

    fn user_redeem(&self, user: &ManagedAddress) -> RewaOrDcdtTokenPayment {
        let user_id = self.require_user_whitelisted(user);
        let user_redeemed_mapper = self.user_redeemed(user_id);
        require!(
            user_redeemed_mapper.get() != USER_REDEEMED,
            "User already redeemed"
        );

        let total_user_deposit = self.total_deposit_by_user(user_id).take();

        let accepted_token_id = self.accepted_token_id().get();
        let accepted_token_sc_balance = self.blockchain().get_sc_balance(&accepted_token_id, 0);
        let launched_token_supply = self.launched_token_balance().get();

        // only allow users to withdraw if the launched tokens were deposited AND the owner withdrew his accepted tokens
        let output_tokens = if launched_token_supply != 0 && accepted_token_sc_balance == 0 {
            let bought_tokens = self.compute_user_bought_tokens(&total_user_deposit);
            self.send().direct_non_zero(
                user,
                &bought_tokens.token_identifier,
                0,
                &bought_tokens.amount,
            );

            bought_tokens
        } else {
            self.send()
                .direct_non_zero(user, &accepted_token_id, 0, &total_user_deposit);

            RewaOrDcdtTokenPayment::new(accepted_token_id, 0, total_user_deposit)
        };

        user_redeemed_mapper.set(USER_REDEEMED);

        output_tokens
    }

    fn compute_user_bought_tokens(&self, redeem_amount: &BigUint) -> RewaOrDcdtTokenPayment {
        let total_deposit_all_users = self.accepted_token_balance().get();
        let launched_token_id = RewaOrDcdtTokenIdentifier::dcdt(self.launched_token_id().get());
        let total_launched_token_supply = self.launched_token_balance().get();
        let reward_amount = total_launched_token_supply * redeem_amount / total_deposit_all_users;

        RewaOrDcdtTokenPayment::new(launched_token_id, 0, reward_amount)
    }

    fn require_owner_didnt_redeem(&self) {
        require!(
            self.owner_redeemed().get() != USER_REDEEMED,
            "Owner already redeemed"
        );
    }

    #[storage_mapper("userRedeemed")]
    fn user_redeemed(&self, user_id: AddressId) -> SingleValueMapper<UserRedeemFlag>;

    #[storage_mapper("ownerRedeemed")]
    fn owner_redeemed(&self) -> SingleValueMapper<UserRedeemFlag>;
}
