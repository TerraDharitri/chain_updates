dharitri_sc::imports!();
dharitri_sc::derive_imports!();

#[derive(
    ManagedVecItem,
    TopEncode,
    TopDecode,
    NestedEncode,
    NestedDecode,
    TypeAbi,
    Clone,
    PartialEq,
    Debug,
)]
pub struct UnstakePair<M: ManagedTypeApi> {
    pub unlock_epoch: u64,
    pub locked_tokens: DcdtTokenPayment<M>,
    pub unlocked_tokens: DcdtTokenPayment<M>,
}

#[dharitri_sc::module]
pub trait TokensPerUserModule {
    #[view(getUnbondEpochs)]
    #[storage_mapper("unbondEpochs")]
    fn unbond_epochs(&self) -> SingleValueMapper<u64>;

    #[view(getUnlockedTokensForUser)]
    #[storage_mapper("unlockedTokensForUser")]
    fn unlocked_tokens_for_user(
        &self,
        address: &ManagedAddress,
    ) -> SingleValueMapper<ManagedVec<UnstakePair<Self::Api>>>;
}
