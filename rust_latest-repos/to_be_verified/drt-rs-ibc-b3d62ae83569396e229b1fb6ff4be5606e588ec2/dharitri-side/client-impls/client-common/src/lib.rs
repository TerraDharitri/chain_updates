#![no_std]

use common_types::{channel_types::height, ClientId, Hash, Path, UnixTimestamp};

dharitri_sc::imports!();
dharitri_sc::derive_imports!();

pub struct ConsensusStateUpdate<M: ManagedTypeApi> {
    pub consensus_state_commitment: Hash<M>,
    pub height: height::Data,
}

#[derive(TypeAbi, TopEncode, TopDecode, NestedEncode, NestedDecode, PartialEq)]
pub enum ClientStatus {
    None,
    Active,
    Expired,
    Frozen,
}

#[derive(TypeAbi, TopEncode, TopDecode)]
pub struct GetLatestInfoResultType {
    pub latest_height: height::Data,
    pub latest_timestamp: UnixTimestamp,
    pub client_status: ClientStatus,
}

#[derive(TypeAbi, TopEncode, TopDecode)]
pub struct VerifyMembershipArgs<M: ManagedTypeApi> {
    pub client_id: ClientId<M>,
    pub height: height::Data,
    pub delay_time_period: UnixTimestamp,
    pub delay_block_period: u64,
    pub proof: Hash<M>,
    pub prefix: ManagedBuffer<M>,
    pub path: Path<M>,
    pub value: ManagedBuffer<M>,
}

#[derive(TypeAbi, TopEncode, TopDecode)]
pub struct VerifyNonMembershipArgs<M: ManagedTypeApi> {
    pub client_id: ClientId<M>,
    pub height: height::Data,
    pub delay_time_period: UnixTimestamp,
    pub delay_block_period: u64,
    pub proof: Hash<M>,
    pub prefix: ManagedBuffer<M>,
    pub path: Path<M>,
}

#[dharitri_sc::module]
pub trait CommonClientLogicModule {
    fn set_ibc_handler(&self, ibc_handler: &ManagedAddress) {
        require!(
            self.blockchain().is_smart_contract(ibc_handler) && !ibc_handler.is_zero(),
            "Invalid SC address"
        );

        self.ibc_handler().set(ibc_handler);
    }

    fn require_ibc_handler_caller(&self) {
        let caller = self.blockchain().get_caller();
        let ibc_handler = self.ibc_handler().get();
        require!(
            caller == ibc_handler,
            "Only the IBC handler may call this endpoint"
        );
    }

    #[view(getIbcHandler)]
    #[storage_mapper("ibcHandler")]
    fn ibc_handler(&self) -> SingleValueMapper<ManagedAddress>;
}
