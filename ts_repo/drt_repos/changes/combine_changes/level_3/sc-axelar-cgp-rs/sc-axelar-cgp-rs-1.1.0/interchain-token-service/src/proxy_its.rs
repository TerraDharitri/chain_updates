use core::ops::Deref;
use operatable::ProxyTrait as _;
use token_manager::flow_limit::ProxyTrait as _;
use token_manager::mintership::ProxyTrait as _;
use token_manager::ProxyTrait as _;

use crate::abi::AbiEncodeDecode;
use crate::abi_types::RegisterTokenMetadataPayload;
use crate::constants::{ManagedBufferAscii, TokenId, MESSAGE_TYPE_REGISTER_TOKEN_METADATA};
use crate::{address_tracker, events, proxy_gmp};

dharitri_sc::imports!();

pub mod executable_contract_proxy {
    use crate::constants::TokenId;

    dharitri_sc::imports!();

    // Contracts having these functions should check that the InterchainTokenService contract called them
    #[dharitri_sc::proxy]
    pub trait ExecutableContractProxy {
        #[payable("*")]
        #[endpoint(executeWithInterchainToken)]
        fn execute_with_interchain_token(
            &self,
            source_chain: &ManagedBuffer,
            message_id: &ManagedBuffer,
            source_address: ManagedBuffer,
            data: ManagedBuffer,
            token_id: TokenId<Self::Api>,
        );
    }
}

pub const DCDT_PROPERTIES_TOKEN_NAME_INDEX: usize = 0;
pub const DCDT_PROPERTIES_TOKEN_TYPE_INDEX: usize = 1;
pub const DCDT_PROPERTIES_DECIMALS_BUFFER_INDEX: usize = 5;

#[dharitri_sc::module]
pub trait ProxyItsModule:
    events::EventsModule + proxy_gmp::ProxyGmpModule + address_tracker::AddressTracker
{
    fn token_manager_take_token(
        &self,
        token_id: &TokenId<Self::Api>,
        token_identifier: RewaOrDcdtTokenIdentifier,
        amount: BigUint,
    ) {
        self.token_manager_proxy(self.deployed_token_manager(token_id))
            .take_token()
            .with_rewa_or_single_dcdt_transfer(RewaOrDcdtTokenPayment::new(
                token_identifier,
                0,
                amount,
            ))
            .execute_on_dest_context::<()>();
    }

    fn token_manager_set_flow_limit(
        &self,
        token_id: &TokenId<Self::Api>,
        flow_limit: Option<BigUint>,
    ) {
        self.token_manager_proxy(self.deployed_token_manager(token_id))
            .set_flow_limit(flow_limit)
            .execute_on_dest_context::<()>();
    }

    fn token_manager_give_token(
        &self,
        token_id: &TokenId<Self::Api>,
        destination_address: &ManagedAddress,
        amount: &BigUint,
    ) -> (RewaOrDcdtTokenIdentifier, BigUint) {
        self.token_manager_proxy(self.deployed_token_manager(token_id))
            .give_token(destination_address, amount)
            .execute_on_dest_context::<MultiValue2<RewaOrDcdtTokenIdentifier, BigUint>>()
            .into_tuple()
    }

    fn token_manager_deploy_interchain_token(
        &self,
        token_id: &TokenId<Self::Api>,
        minter: Option<ManagedAddress>,
        name: ManagedBuffer,
        symbol: ManagedBuffer,
        decimals: u8,
        initial_caller: ManagedAddress,
    ) {
        self.token_manager_proxy(self.deployed_token_manager(token_id))
            .deploy_interchain_token(
                minter,
                name,
                symbol,
                decimals,
                OptionalValue::Some(initial_caller),
            )
            .with_rewa_transfer(self.call_value().rewa_value().clone_value())
            .with_gas_limit(100_000_000) // Need to specify gas manually here because the function does an async call. This should be plenty
            .execute_on_dest_context::<()>();
    }

    fn token_manager_get_opt_token_identifier(
        &self,
        sc_address: ManagedAddress,
    ) -> Option<RewaOrDcdtTokenIdentifier> {
        self.token_manager_proxy(sc_address)
            .get_opt_token_identifier()
            .execute_on_dest_context()
    }

    fn token_manager_mint(
        &self,
        sc_address: ManagedAddress,
        address: ManagedAddress,
        amount: BigUint,
    ) {
        self.token_manager_proxy(sc_address)
            .mint(address, amount)
            .execute_on_dest_context::<()>();
    }

    fn token_manager_transfer_mintership(
        &self,
        sc_address: ManagedAddress,
        minter: ManagedAddress,
    ) {
        self.token_manager_proxy(sc_address)
            .transfer_mintership(minter)
            .execute_on_dest_context::<()>();
    }

    fn token_manager_remove_flow_limiter(
        &self,
        sc_address: ManagedAddress,
        flow_limiter: ManagedAddress,
    ) {
        self.token_manager_proxy(sc_address)
            .remove_flow_limiter(flow_limiter)
            .execute_on_dest_context::<()>();
    }

    fn token_manager_add_flow_limiter(
        &self,
        sc_address: ManagedAddress,
        flow_limiter: ManagedAddress,
    ) {
        self.token_manager_proxy(sc_address)
            .add_flow_limiter(flow_limiter)
            .execute_on_dest_context::<()>();
    }

    fn token_manager_transfer_operatorship(
        &self,
        sc_address: ManagedAddress,
        operator: ManagedAddress,
    ) {
        self.token_manager_proxy(sc_address)
            .transfer_operatorship(operator)
            .execute_on_dest_context::<()>();
    }

    fn token_manager_is_minter(&self, sc_address: ManagedAddress, minter: &ManagedAddress) -> bool {
        self.token_manager_proxy(sc_address)
            .is_minter(minter)
            .execute_on_dest_context()
    }

    fn executable_contract_execute_with_interchain_token(
        &self,
        destination_address: ManagedAddress,
        original_source_chain: ManagedBuffer,
        message_id: ManagedBuffer,
        original_source_address: ManagedBuffer,
        data: ManagedBuffer,
        token_id: TokenId<Self::Api>,
        token_identifier: RewaOrDcdtTokenIdentifier,
        amount: BigUint,
    ) {
        self.executable_contract_proxy(destination_address)
            .execute_with_interchain_token(
                &original_source_chain,
                &message_id,
                original_source_address,
                data,
                token_id.clone(),
            )
            .with_rewa_or_single_dcdt_transfer((token_identifier.clone(), 0, amount.clone()))
            .execute_on_dest_context::<()>();
    }

    fn register_token_metadata_async_call(
        &self,
        token_identifier: TokenIdentifier,
        gas_value: BigUint,
    ) {
        self.dcdt_get_token_properties(
            token_identifier.clone(),
            self.callbacks().register_token_metadata_callback(
                token_identifier,
                gas_value,
                self.blockchain().get_caller(),
            ),
        );
    }

    fn dcdt_get_token_properties(
        &self,
        token_identifier: TokenIdentifier,
        callback: CallbackClosure<Self::Api>,
    ) {
        let mut contract_call = self.send().contract_call::<()>(
            DCDTSystemSCAddress.to_managed_address(),
            ManagedBuffer::from("getTokenProperties"),
        );
        contract_call.push_raw_argument(token_identifier.into_managed_buffer());

        contract_call
            .async_call()
            .with_callback(callback)
            .call_and_exit();
    }

    #[view(flowLimit)]
    fn flow_limit(&self, token_id: TokenId<Self::Api>) -> BigUint {
        self.token_manager_proxy(self.deployed_token_manager(&token_id))
            .flow_limit()
            .execute_on_dest_context()
    }

    #[view(flowOutAmount)]
    fn flow_out_amount(&self, token_id: TokenId<Self::Api>) -> BigUint {
        self.token_manager_proxy(self.deployed_token_manager(&token_id))
            .get_flow_out_amount()
            .execute_on_dest_context()
    }

    #[view(flowInAmount)]
    fn flow_in_amount(&self, token_id: TokenId<Self::Api>) -> BigUint {
        self.token_manager_proxy(self.deployed_token_manager(&token_id))
            .get_flow_in_amount()
            .execute_on_dest_context()
    }

    #[view(deployedTokenManager)]
    fn deployed_token_manager(&self, token_id: &TokenId<Self::Api>) -> ManagedAddress {
        let token_manager_address_mapper = self.token_manager_address(token_id);

        require!(
            !token_manager_address_mapper.is_empty(),
            "Token manager does not exist"
        );

        token_manager_address_mapper.get()
    }

    #[view(registeredTokenIdentifier)]
    fn registered_token_identifier(
        &self,
        token_id: &TokenId<Self::Api>,
    ) -> RewaOrDcdtTokenIdentifier {
        self.token_manager_proxy(self.deployed_token_manager(token_id))
            .token_identifier()
            .execute_on_dest_context()
    }

    #[view(invalidTokenManagerAddress)]
    fn get_opt_token_manager_address(
        &self,
        token_id: &TokenId<Self::Api>,
    ) -> Option<ManagedAddress> {
        let token_manager_address_mapper = self.token_manager_address(token_id);

        if token_manager_address_mapper.is_empty() {
            return None;
        }

        Some(token_manager_address_mapper.get())
    }

    #[view(tokenManagerAddress)]
    #[storage_mapper("token_manager_address")]
    fn token_manager_address(
        &self,
        token_id: &TokenId<Self::Api>,
    ) -> SingleValueMapper<ManagedAddress>;

    #[proxy]
    fn token_manager_proxy(&self, address: ManagedAddress) -> token_manager::Proxy<Self::Api>;

    #[proxy]
    fn executable_contract_proxy(
        &self,
        sc_address: ManagedAddress,
    ) -> executable_contract_proxy::Proxy<Self::Api>;

    #[callback]
    fn register_token_metadata_callback(
        &self,
        token_identifier: TokenIdentifier,
        gas_value: BigUint,
        caller: ManagedAddress,
        #[call_result] result: ManagedAsyncCallResult<MultiValueEncoded<ManagedBuffer>>,
    ) {
        match result {
            ManagedAsyncCallResult::Ok(values) => {
                let vec: ManagedVec<ManagedBuffer> = values.into_vec_of_buffers();

                let token_type = vec.get(DCDT_PROPERTIES_TOKEN_TYPE_INDEX);
                let decimals_buffer_ref = vec.get(DCDT_PROPERTIES_DECIMALS_BUFFER_INDEX);

                if token_type.deref() != DcdtTokenType::Fungible.as_type_name() {
                    // Send back paid cross chain gas value to initial caller if token is non fungible
                    self.send().direct_non_zero_rewa(&caller, &gas_value);

                    return;
                }

                let decimals_buffer = decimals_buffer_ref.deref();
                // num decimals is in format string NumDecimals-DECIMALS
                // skip `NumDecimals-` part and convert to number
                let token_decimals_buf: ManagedBuffer = decimals_buffer
                    .copy_slice(12, decimals_buffer.len() - 12)
                    .unwrap();
                let token_decimals = token_decimals_buf.ascii_to_u8();

                self.register_token_metadata_raw(
                    RewaOrDcdtTokenIdentifier::dcdt(token_identifier),
                    token_decimals,
                    gas_value,
                );
            }
            ManagedAsyncCallResult::Err(_) => {
                // Send back paid gas value to initial caller
                self.send().direct_non_zero_rewa(&caller, &gas_value);
            }
        }
    }

    fn register_token_metadata_raw(
        &self,
        token_identifier: RewaOrDcdtTokenIdentifier,
        decimals: u8,
        gas_value: BigUint,
    ) {
        self.token_metadata_registered_event(&token_identifier, decimals);

        let data = RegisterTokenMetadataPayload {
            message_type: BigUint::from(MESSAGE_TYPE_REGISTER_TOKEN_METADATA),
            token_identifier: token_identifier.into_name(),
            decimals,
        };

        let payload = data.abi_encode();

        self.call_contract_its_hub(payload, gas_value);
    }
}
