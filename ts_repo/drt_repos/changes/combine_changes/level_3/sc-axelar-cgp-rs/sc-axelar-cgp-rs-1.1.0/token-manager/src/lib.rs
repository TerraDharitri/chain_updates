#![no_std]

use core::ops::Deref;
use dharitri_sc::api::KECCAK256_RESULT_LEN;

use constants::{DeployTokenManagerParams, TokenManagerType};
use operatable::roles::Roles;

use crate::constants::{ManagedBufferAscii, DEFAULT_DCDT_ISSUE_COST};

dharitri_sc::imports!();

pub mod constants;
pub mod flow_limit;
pub mod mintership;

#[dharitri_sc::contract]
pub trait TokenManagerLockUnlockContract:
    flow_limit::FlowLimit
    + operatable::Operatable
    + operatable::roles::AccountRoles
    + mintership::Mintership
{
    #[init]
    fn init(
        &self,
        interchain_token_service: ManagedAddress,
        implementation_type: TokenManagerType,
        interchain_token_id: ManagedByteArray<KECCAK256_RESULT_LEN>,
        params: DeployTokenManagerParams<Self::Api>,
    ) {
        require!(!interchain_token_service.is_zero(), "Zero address");

        self.interchain_token_service()
            .set_if_empty(interchain_token_service.clone());
        self.implementation_type().set_if_empty(implementation_type);
        self.interchain_token_id().set_if_empty(interchain_token_id);

        // From setup of TokenManager
        let operator = if params.operator.is_none() {
            ManagedAddress::zero()
        } else {
            params.operator.unwrap()
        };

        // If an operator is not provided, set zero address as the operator.
        // This allows anyone to easily check if a custom operator was set on the token manager.
        self.add_role(operator, Roles::FLOW_LIMITER | Roles::OPERATOR);
        // Add operator and flow limiter role to the service. The operator can remove the flow limiter role if they so chose and the service has no way to use the operator role for now.
        self.add_role(
            interchain_token_service,
            Roles::FLOW_LIMITER | Roles::OPERATOR,
        );

        match implementation_type {
            TokenManagerType::NativeInterchainToken => {
                require!(params.token_identifier.is_none(), "Invalid token address");
            }
            TokenManagerType::LockUnlock | TokenManagerType::LockUnlockFee => {
                require!(params.token_identifier.is_some(), "Invalid token address");
            }
            TokenManagerType::MintBurn | TokenManagerType::MintBurnFrom => {
                require!(
                    params.token_identifier.is_some()
                        && params.token_identifier.clone().unwrap().is_dcdt(),
                    "Invalid token address"
                );
            }
        }

        if params.token_identifier.is_some() {
            self.token_identifier()
                .set_if_empty(params.token_identifier.unwrap());
        }
    }

    #[upgrade]
    fn upgrade(&self) {}

    #[endpoint(addFlowLimiter)]
    fn add_flow_limiter(&self, flow_limiter: ManagedAddress) {
        self.only_operator();

        self.add_role(flow_limiter, Roles::FLOW_LIMITER);
    }

    #[endpoint(removeFlowLimiter)]
    fn remove_flow_limiter(&self, flow_limiter: ManagedAddress) {
        self.only_operator();

        self.remove_role(flow_limiter, Roles::FLOW_LIMITER);
    }

    #[endpoint(transferFlowLimiter)]
    fn transfer_flow_limiter(&self, from: ManagedAddress, to: ManagedAddress) {
        self.only_operator();

        self.transfer_role(from, to, Roles::FLOW_LIMITER);
    }

    #[endpoint(setFlowLimit)]
    fn set_flow_limit(&self, flow_limit: Option<BigUint>) {
        self.only_flow_limiter();

        self.set_flow_limit_raw(flow_limit, self.interchain_token_id().get());
    }

    #[endpoint(giveToken)]
    fn give_token(
        &self,
        destination_address: &ManagedAddress,
        amount: BigUint,
    ) -> MultiValue2<RewaOrDcdtTokenIdentifier, BigUint> {
        self.only_service();

        let token_identifier = self.token_identifier().get();

        self.add_flow_in_raw(&amount);

        let implementation_type = self.implementation_type().get();
        match implementation_type {
            TokenManagerType::NativeInterchainToken
            | TokenManagerType::MintBurn
            | TokenManagerType::MintBurnFrom => {
                self.give_token_mint_burn(&token_identifier, destination_address, &amount);
            }
            TokenManagerType::LockUnlock | TokenManagerType::LockUnlockFee => {
                self.give_token_lock_unlock(&token_identifier, destination_address, &amount);
            }
        }

        (token_identifier, amount).into()
    }

    #[payable("*")]
    #[endpoint(takeToken)]
    fn take_token(&self) -> BigUint {
        self.only_service();

        let (token_identifier, amount) = self.require_correct_token();

        self.add_flow_out_raw(&amount);

        let implementation_type = self.implementation_type().get();
        match implementation_type {
            TokenManagerType::NativeInterchainToken
            | TokenManagerType::MintBurn
            | TokenManagerType::MintBurnFrom => {
                self.take_token_mint_burn(token_identifier, &amount);
            }
            // nothing to do for lock/unlock, tokens remain in contract
            TokenManagerType::LockUnlock | TokenManagerType::LockUnlockFee => {}
        }

        amount
    }

    /// Lock/unlock type only functions

    #[payable("*")]
    #[endpoint(donateTokens)]
    fn donate_tokens(&self) {
        require!(
            self.implementation_type().get() == TokenManagerType::LockUnlock,
            "Not lock/unlock token manager"
        );

        let _ = self.require_correct_token();

        // Nothing else to do, tokens will be left in the contract
    }

    /// NativeInterchainToken type only functions

    // Somewhat equivalent to Axelar InterchainToken init method
    #[payable("REWA")]
    #[endpoint(deployInterchainToken)]
    fn deploy_interchain_token(
        &self,
        minter: Option<ManagedAddress>,
        name: ManagedBuffer,
        symbol: ManagedBuffer,
        decimals: u8,
        initial_caller: OptionalValue<ManagedAddress>,
    ) {
        require!(
            self.implementation_type().get() == TokenManagerType::NativeInterchainToken,
            "Not native interchain token manager"
        );

        require!(
            self.token_identifier().is_empty(),
            "Token address already exists"
        );

        let caller = self.blockchain().get_caller();
        let interchain_token_service = self.interchain_token_service().get();

        // Also allow minter to call this (if set) in case issue dcdt fails
        require!(
            caller == interchain_token_service || self.is_minter(&caller),
            "Not service or minter"
        );

        require!(!name.is_empty(), "Empty token name");
        require!(!symbol.is_empty(), "Empty token symbol");

        if self.minter_address().is_empty() {
            if minter.is_some() {
                self.add_minter(minter.unwrap());
            } else {
                self.add_minter(ManagedAddress::zero());
            }
        }

        let issue_cost = BigUint::from(DEFAULT_DCDT_ISSUE_COST);

        require!(
            self.call_value().rewa_value().deref() == &issue_cost,
            "Invalid dcdt issue cost"
        );

        self.send()
            .dcdt_system_sc_proxy()
            .issue_and_set_all_roles(
                issue_cost,
                name.to_normalized_token_name(),
                symbol.to_normalized_token_ticker(),
                DcdtTokenType::Fungible,
                decimals as usize,
            )
            .with_callback(
                self.callbacks().deploy_token_callback(
                    initial_caller
                        .into_option()
                        .unwrap_or(self.blockchain().get_caller()),
                ),
            )
            .async_call_and_exit();
    }

    #[endpoint]
    fn mint(&self, address: ManagedAddress, amount: &BigUint) {
        require!(
            self.implementation_type().get() == TokenManagerType::NativeInterchainToken,
            "Not native interchain token manager"
        );

        self.only_minter();

        require!(
            !self.token_identifier().is_empty(),
            "Token address not yet set"
        );

        let token_identifier = self.token_identifier().get().into_dcdt_option().unwrap();

        self.send().dcdt_local_mint(&token_identifier, 0, amount);
        self.send()
            .direct_dcdt(&address, &token_identifier, 0, amount)
    }

    #[payable("*")]
    #[endpoint]
    fn burn(&self) {
        require!(
            self.implementation_type().get() == TokenManagerType::NativeInterchainToken,
            "Not native interchain token manager"
        );

        self.only_minter();

        require!(
            !self.token_identifier().is_empty(),
            "Token address not yet set"
        );

        let (token_identifier, amount) = self.require_correct_token();

        self.send()
            .dcdt_local_burn(&token_identifier.unwrap_dcdt(), 0, &amount);
    }

    fn only_service(&self) {
        require!(
            self.blockchain().get_caller() == self.interchain_token_service().get(),
            "Not service"
        );
    }

    fn only_flow_limiter(&self) {
        self.only_role(Roles::FLOW_LIMITER);
    }

    fn require_correct_token(&self) -> (RewaOrDcdtTokenIdentifier, BigUint) {
        let (token_identifier, amount) = self.call_value().rewa_or_single_fungible_dcdt();

        let required_token_identifier = self.token_identifier().get();

        require!(
            token_identifier == required_token_identifier,
            "Wrong token sent"
        );

        (token_identifier, amount)
    }

    fn give_token_lock_unlock(
        &self,
        token_identifier: &RewaOrDcdtTokenIdentifier,
        destination_address: &ManagedAddress,
        amount: &BigUint,
    ) {
        self.send()
            .direct(destination_address, token_identifier, 0, amount);
    }

    fn give_token_mint_burn(
        &self,
        token_identifier: &RewaOrDcdtTokenIdentifier,
        destination_address: &ManagedAddress,
        amount: &BigUint,
    ) {
        self.send()
            .dcdt_local_mint(&token_identifier.clone().unwrap_dcdt(), 0, amount);

        self.send()
            .direct(destination_address, token_identifier, 0, amount);
    }

    fn take_token_mint_burn(&self, token_identifier: RewaOrDcdtTokenIdentifier, amount: &BigUint) {
        self.send()
            .dcdt_local_burn(&token_identifier.unwrap_dcdt(), 0, amount);
    }

    #[view(getImplementationTypeAndTokenIdentifier)]
    fn get_implementation_type_and_token_identifier(
        &self,
    ) -> MultiValue2<TokenManagerType, RewaOrDcdtTokenIdentifier> {
        MultiValue2::from((
            self.implementation_type().get(),
            self.token_identifier().get(),
        ))
    }

    #[view(isFlowLimiter)]
    fn is_flow_limiter(&self, address: &ManagedAddress) -> bool {
        self.has_role(address, Roles::FLOW_LIMITER)
    }

    // Mainly be used by frontends
    #[view(params)]
    fn params(
        &self,
        operator: Option<ManagedAddress>,
        token_identifier: Option<RewaOrDcdtTokenIdentifier>,
    ) -> DeployTokenManagerParams<Self::Api> {
        DeployTokenManagerParams {
            operator,
            token_identifier,
        }
    }

    #[view(optTokenIdentifier)]
    fn get_opt_token_identifier(&self) -> Option<RewaOrDcdtTokenIdentifier> {
        let token_identifier_mapper = self.token_identifier();

        if token_identifier_mapper.is_empty() {
            return None;
        }

        Some(token_identifier_mapper.get())
    }

    #[view(interchainTokenService)]
    #[storage_mapper("interchain_token_service")]
    fn interchain_token_service(&self) -> SingleValueMapper<ManagedAddress>;

    #[view(implementationType)]
    #[storage_mapper("implementation_type")]
    fn implementation_type(&self) -> SingleValueMapper<TokenManagerType>;

    #[view(interchainTokenId)]
    #[storage_mapper("interchain_token_id")]
    fn interchain_token_id(&self) -> SingleValueMapper<ManagedByteArray<KECCAK256_RESULT_LEN>>;

    #[view(tokenIdentifier)]
    #[storage_mapper("token_identifier")]
    fn token_identifier(&self) -> SingleValueMapper<RewaOrDcdtTokenIdentifier>;

    #[callback]
    fn deploy_token_callback(
        &self,
        user: ManagedAddress,
        #[call_result] result: ManagedAsyncCallResult<TokenIdentifier>,
    ) {
        match result {
            ManagedAsyncCallResult::Ok(token_id_raw) => {
                let token_identifier = RewaOrDcdtTokenIdentifier::dcdt(token_id_raw);

                self.interchain_token_deployed_event(
                    self.interchain_token_id().get(),
                    &token_identifier,
                );

                self.token_identifier().set_if_empty(token_identifier);
            }
            ManagedAsyncCallResult::Err(_) => {
                self.interchain_token_deployment_failed();

                self.send().direct_rewa(&user, self.call_value().rewa_value().deref());
            }
        }
    }

    #[event("interchain_token_deployment_failed")]
    fn interchain_token_deployment_failed(&self);

    #[event("interchain_token_deployed_event")]
    fn interchain_token_deployed_event(
        &self,
        #[indexed] token_id: ManagedByteArray<KECCAK256_RESULT_LEN>,
        #[indexed] token_identifier: &RewaOrDcdtTokenIdentifier,
    );
}
