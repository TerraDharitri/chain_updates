dharitri_sc::imports!();

pub type PaymentsVec<M> = ManagedVec<M, DcdtTokenPayment<M>>;

// lp farm

pub struct LpFarmClaimRewardsResult<M: ManagedTypeApi> {
    pub new_lp_farm_tokens: DcdtTokenPayment<M>,
    pub lp_farm_rewards: DcdtTokenPayment<M>,
}

pub struct LpFarmExitResult<M: ManagedTypeApi> {
    pub lp_tokens: DcdtTokenPayment<M>,
    pub lp_farm_rewards: DcdtTokenPayment<M>,
}

// staking farm

pub struct StakingFarmEnterResult<M: ManagedTypeApi> {
    pub received_staking_farm_token: DcdtTokenPayment<M>,
}

pub struct StakingFarmClaimRewardsResult<M: ManagedTypeApi> {
    pub new_staking_farm_tokens: DcdtTokenPayment<M>,
    pub staking_farm_rewards: DcdtTokenPayment<M>,
}

pub struct StakingFarmExitResult<M: ManagedTypeApi> {
    pub unbond_staking_farm_token: DcdtTokenPayment<M>,
    pub staking_rewards: DcdtTokenPayment<M>,
}

// pair

pub struct PairRemoveLiquidityResult<M: ManagedTypeApi> {
    pub staking_token_payment: DcdtTokenPayment<M>,
    pub other_token_payment: DcdtTokenPayment<M>,
}
