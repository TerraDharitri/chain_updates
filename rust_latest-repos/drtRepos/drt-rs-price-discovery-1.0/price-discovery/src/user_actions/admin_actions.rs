use crate::{
    phase::{Phase, MAX_LONG_PHASE_DURATION, MAX_PHASE_DURATION},
    Timestamp,
};

dharitri_sc::imports!();

pub static INVALID_CURRENT_PHASE_ERR_MSG: &[u8] = b"Invalid current phase";
pub static INVALID_TIMESTAMP_CHANGE_ERR_MSG: &[u8] = b"Invalid timestamp change";
pub static INVALID_TIMESTAMP_DURATION_ERR_MGS: &[u8] = b"Invalid timestamp";

#[dharitri_sc::module]
pub trait AdminActionsModule:
    super::user_deposit_withdraw::UserDepositWithdrawModule
    + crate::common_storage::CommonStorageModule
    + crate::events::EventsModule
    + crate::phase::PhaseModule
    + dharitri_sc_modules::default_issue_callbacks::DefaultIssueCallbacksModule
{
    #[endpoint(setUserDepositWithdrawTime)]
    fn set_user_deposit_withdraw_time(&self, user_deposit_withdraw_time: Timestamp) {
        self.require_caller_admin();
        self.require_valid_timestamp(user_deposit_withdraw_time);

        let current_phase = self.get_current_phase();
        if current_phase == Phase::UserDepositWithdraw {
            let current_time = self.user_deposit_withdraw_time().get();
            require!(
                current_time <= user_deposit_withdraw_time,
                "May only extend phase at this point, not reduce"
            );
        }

        self.set_timestamp(
            user_deposit_withdraw_time,
            &Phase::UserDepositWithdraw,
            &self.user_deposit_withdraw_time(),
        );

        self.set_user_deposit_withdraw_time_event(user_deposit_withdraw_time);
    }

    #[endpoint(setOwnerDepositWithdrawTime)]
    fn set_owner_deposit_withdraw_time(&self, owner_deposit_withdraw_time: Timestamp) {
        self.require_caller_admin();
        self.require_valid_timestamp(owner_deposit_withdraw_time);

        self.set_timestamp(
            owner_deposit_withdraw_time,
            &Phase::OwnerDepositWithdraw,
            &self.owner_deposit_withdraw_time(),
        );

        self.set_owner_deposit_withdraw_time_event(owner_deposit_withdraw_time);
    }

    #[endpoint(setOwnerRedeemTime)]
    fn set_owner_redeem_time(&self, owner_redeem_time: Timestamp) {
        self.require_caller_admin();
        self.require_valid_long_timestamp(owner_redeem_time);

        self.set_timestamp(
            owner_redeem_time,
            &Phase::OwnerRedeem,
            &self.owner_redeem_time(),
        );

        self.set_owner_redeem_time_event(owner_redeem_time);
    }

    #[endpoint(setMinLaunchedTokens)]
    fn set_min_launched_tokens(&self, min_launched_tokens: BigUint) {
        self.require_caller_admin();
        require!(min_launched_tokens > 0, "Invalid min launched tokens");

        let phase = self.get_current_phase();
        self.require_before_redeem(&phase);

        self.min_launched_tokens().set(min_launched_tokens);
    }

    /// Pass `0` for `limit` if there is no limit
    #[endpoint(setUserLimit)]
    fn set_user_limit(&self, user: ManagedAddress, limit: BigUint) {
        self.require_caller_admin();

        let user_id = self.user_id_mapper().get_id_non_zero(&user);
        let user_current_deposit = self.total_deposit_by_user(user_id).get();
        if user_current_deposit == 0 || limit == 0 {
            self.set_user_deposit_limit(&user, user_id, &limit);

            return;
        }

        require!(
            user_current_deposit <= limit,
            "May not set user limit below current deposit value"
        );

        self.set_user_deposit_limit(&user, user_id, &limit);
    }

    /// `whitelist` arguments are pairs of (address, max_total_deposit). Pass `0` for `max_total_deposit` if there is no limit
    #[endpoint(addUsersToWhitelist)]
    fn add_users_to_whitelist(
        &self,
        whitelist: MultiValueEncoded<MultiValue2<ManagedAddress, BigUint>>,
    ) {
        self.require_caller_admin();

        let phase = self.get_current_phase();
        self.require_before_redeem(&phase);

        let id_mapper = self.user_id_mapper();
        let whitelist_mapper = self.user_whitelist();
        for pair in whitelist {
            let (user, limit) = pair.into_tuple();
            let user_id = id_mapper.insert_new(&user);
            whitelist_mapper.add(&user_id);

            if limit > 0 {
                self.user_deposit_limit(user_id).set(limit);
            }
        }
    }

    #[endpoint(refundUsers)]
    fn refund_users(&self, users: MultiValueEncoded<ManagedAddress>) {
        self.require_caller_admin();

        let phase = self.get_current_phase();
        self.require_before_redeem(&phase);

        let id_mapper = self.user_id_mapper();
        let whitelist_mapper = self.user_whitelist();
        let owner_address = self.blockchain().get_owner_address();
        for user in users {
            self.refund_single_user(&owner_address, &user, &id_mapper, &whitelist_mapper);
        }
    }

    fn refund_single_user(
        &self,
        owner_address: &ManagedAddress,
        user_addr: &ManagedAddress,
        id_mapper: &AddressToIdMapper,
        whitelist_mapper: &WhitelistMapper<AddressId>,
    ) {
        require!(user_addr != owner_address, "May not refund owner");

        let user_id = id_mapper.get_id_non_zero(user_addr);
        whitelist_mapper.require_whitelisted(&user_id);
        whitelist_mapper.remove(&user_id);

        let user_deposit = self.total_deposit_by_user(user_id).get();
        self.user_deposit_limit(user_id).clear();
        self.user_withdraw(user_addr, user_id, &user_deposit);

        self.emit_refund_user_event(user_addr);
    }

    fn require_caller_admin(&self) {
        let caller = self.blockchain().get_caller();
        let admin = self.admin().get();
        require!(caller == admin, "Only admin may call this function");
    }

    fn set_user_deposit_limit(
        &self,
        user_addr: &ManagedAddress,
        user_id: AddressId,
        limit: &BigUint,
    ) {
        if limit > &0 {
            let user_min_deposit = self.user_min_deposit().get();
            require!(
                limit >= &user_min_deposit,
                "May not set limit under min user deposit"
            );
        }

        self.user_deposit_limit(user_id).set(limit);
        self.set_user_limit_event(user_addr, limit);
    }

    fn set_timestamp(
        &self,
        new_timestamp: Timestamp,
        required_phase_limit: &Phase,
        mapper: &SingleValueMapper<Timestamp>,
    ) {
        let phase_before = self.get_current_phase();
        require!(
            &phase_before <= required_phase_limit,
            INVALID_CURRENT_PHASE_ERR_MSG
        );

        mapper.set(new_timestamp);

        let phase_after = self.get_current_phase();
        require!(
            phase_before == phase_after,
            INVALID_TIMESTAMP_CHANGE_ERR_MSG
        );
    }

    fn require_valid_timestamp(&self, timestamp: Timestamp) {
        require!(
            timestamp > 0 && timestamp <= MAX_PHASE_DURATION,
            INVALID_TIMESTAMP_DURATION_ERR_MGS
        );
    }

    fn require_valid_long_timestamp(&self, timestamp: Timestamp) {
        require!(
            timestamp > 0 && timestamp <= MAX_LONG_PHASE_DURATION,
            INVALID_TIMESTAMP_DURATION_ERR_MGS
        );
    }

    #[storage_mapper("admin")]
    fn admin(&self) -> SingleValueMapper<ManagedAddress>;
}
