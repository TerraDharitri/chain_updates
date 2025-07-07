use crate::{
    chain_core::builtin_func_names::*,
    types::{DcdtLocalRole, DcdtLocalRoleFlags, RawHandle, VMAddress},
    vm_hooks::VMHooksHandlerSource,
    world_mock::{DcdtData, DcdtInstance},
};
use num_bigint::BigInt;
use num_traits::Zero;

// The Go VM doesn't do it, but if we change that, we can enable it easily here too via this constant.
const DCDT_TOKEN_DATA_FUNC_RESETS_VALUES: bool = false;
const VM_BUILTIN_FUNCTION_NAMES: [&str; 16] = [
    DCDT_LOCAL_MINT_FUNC_NAME,
    DCDT_LOCAL_BURN_FUNC_NAME,
    DCDT_MULTI_TRANSFER_FUNC_NAME,
    DCDT_NFT_TRANSFER_FUNC_NAME,
    DCDT_NFT_CREATE_FUNC_NAME,
    DCDT_NFT_ADD_QUANTITY_FUNC_NAME,
    DCDT_NFT_ADD_URI_FUNC_NAME,
    DCDT_NFT_UPDATE_ATTRIBUTES_FUNC_NAME,
    DCDT_NFT_BURN_FUNC_NAME,
    DCDT_TRANSFER_FUNC_NAME,
    CHANGE_OWNER_BUILTIN_FUNC_NAME,
    CLAIM_DEVELOPER_REWARDS_FUNC_NAME,
    SET_USERNAME_FUNC_NAME,
    MIGRATE_USERNAME_FUNC_NAME,
    DELETE_USERNAME_FUNC_NAME,
    UPGRADE_CONTRACT_FUNC_NAME,
];

pub trait VMHooksBlockchain: VMHooksHandlerSource {
    fn is_contract_address(&self, address_bytes: &[u8]) -> bool {
        let address = VMAddress::from_slice(address_bytes);
        &address == self.current_address()
    }

    fn managed_caller(&self, dest_handle: RawHandle) {
        self.m_types_lock()
            .mb_set(dest_handle, self.input_ref().from.to_vec());
    }

    fn managed_sc_address(&self, dest_handle: RawHandle) {
        self.m_types_lock()
            .mb_set(dest_handle, self.current_address().to_vec());
    }

    fn managed_owner_address(&self, dest_handle: RawHandle) {
        self.m_types_lock().mb_set(
            dest_handle,
            self.current_account_data()
                .contract_owner
                .unwrap_or_else(|| panic!("contract owner address not set"))
                .to_vec(),
        );
    }

    fn get_shard_of_address(&self, address_bytes: &[u8]) -> i32 {
        (address_bytes[address_bytes.len() - 1] % 3).into()
    }

    fn is_smart_contract(&self, address_bytes: &[u8]) -> bool {
        VMAddress::from_slice(address_bytes).is_smart_contract_address()
    }

    fn load_balance(&self, address_bytes: &[u8], dest: RawHandle) {
        assert!(
            self.is_contract_address(address_bytes),
            "get balance not yet implemented for accounts other than the contract itself"
        );
        self.m_types_lock()
            .bi_overwrite(dest, self.current_account_data().rewa_balance.into());
    }

    fn get_tx_hash(&self, dest: RawHandle) {
        self.m_types_lock()
            .mb_set(dest, self.input_ref().tx_hash.to_vec());
    }

    fn get_gas_left(&self) -> u64 {
        self.input_ref().gas_limit
    }

    fn get_block_timestamp(&self) -> u64 {
        self.get_current_block_info().block_timestamp
    }

    fn get_block_nonce(&self) -> u64 {
        self.get_current_block_info().block_nonce
    }

    fn get_block_round(&self) -> u64 {
        self.get_current_block_info().block_round
    }

    fn get_block_epoch(&self) -> u64 {
        self.get_current_block_info().block_epoch
    }

    fn get_block_random_seed(&self, dest: RawHandle) {
        self.m_types_lock().mb_set(
            dest,
            self.get_current_block_info().block_random_seed.to_vec(),
        );
    }

    fn get_prev_block_timestamp(&self) -> u64 {
        self.get_previous_block_info().block_timestamp
    }

    fn get_prev_block_nonce(&self) -> u64 {
        self.get_previous_block_info().block_nonce
    }

    fn get_prev_block_round(&self) -> u64 {
        self.get_previous_block_info().block_round
    }

    fn get_prev_block_epoch(&self) -> u64 {
        self.get_previous_block_info().block_epoch
    }

    fn get_prev_block_random_seed(&self, dest: RawHandle) {
        self.m_types_lock().mb_set(
            dest,
            self.get_previous_block_info().block_random_seed.to_vec(),
        );
    }

    fn get_current_dcdt_nft_nonce(&self, address_bytes: &[u8], token_id_bytes: &[u8]) -> u64 {
        assert!(
            self.is_contract_address(address_bytes),
            "get_current_dcdt_nft_nonce not yet implemented for accounts other than the contract itself"
        );

        self.current_account_data()
            .dcdt
            .get_by_identifier_or_default(token_id_bytes)
            .last_nonce
    }

    fn big_int_get_dcdt_external_balance(
        &self,
        address_bytes: &[u8],
        token_id_bytes: &[u8],
        nonce: u64,
        dest: RawHandle,
    ) {
        assert!(
            self.is_contract_address(address_bytes),
            "get_dcdt_balance not yet implemented for accounts other than the contract itself"
        );

        let dcdt_balance = self
            .current_account_data()
            .dcdt
            .get_dcdt_balance(token_id_bytes, nonce);
        self.m_types_lock().bi_overwrite(dest, dcdt_balance.into());
    }

    fn managed_get_code_metadata(&self, address_handle: i32, response_handle: i32) {
        let address = VMAddress::from_slice(self.m_types_lock().mb_get(address_handle));
        let Some(data) = self.account_data(&address) else {
            self.vm_error(&format!(
                "account not found: {}",
                hex::encode(address.as_bytes())
            ))
        };
        let code_metadata_bytes = data.code_metadata.to_byte_array();
        self.m_types_lock()
            .mb_set(response_handle, code_metadata_bytes.to_vec())
    }

    fn managed_is_builtin_function(&self, function_name_handle: i32) -> bool {
        VM_BUILTIN_FUNCTION_NAMES.contains(
            &self
                .m_types_lock()
                .mb_to_function_name(function_name_handle)
                .as_str(),
        )
    }

    #[allow(clippy::too_many_arguments)]
    fn managed_get_dcdt_token_data(
        &self,
        address_handle: RawHandle,
        token_id_handle: RawHandle,
        nonce: u64,
        value_handle: RawHandle,
        properties_handle: RawHandle,
        hash_handle: RawHandle,
        name_handle: RawHandle,
        attributes_handle: RawHandle,
        creator_handle: RawHandle,
        royalties_handle: RawHandle,
        uris_handle: RawHandle,
    ) {
        let address = VMAddress::from_slice(self.m_types_lock().mb_get(address_handle));
        let token_id_bytes = self.m_types_lock().mb_get(token_id_handle).to_vec();

        if let Some(account) = self.account_data(&address) {
            if let Some(dcdt_data) = account.dcdt.get_by_identifier(token_id_bytes.as_slice()) {
                if let Some(instance) = dcdt_data.instances.get_by_nonce(nonce) {
                    self.set_dcdt_data_values(
                        dcdt_data,
                        instance,
                        value_handle,
                        properties_handle,
                        hash_handle,
                        name_handle,
                        attributes_handle,
                        creator_handle,
                        royalties_handle,
                        uris_handle,
                    );
                    return;
                }
            }
        }

        // missing account/token identifier/nonce
        self.reset_dcdt_data_values(
            value_handle,
            properties_handle,
            hash_handle,
            name_handle,
            attributes_handle,
            creator_handle,
            royalties_handle,
            uris_handle,
        );
    }

    fn managed_get_back_transfers(
        &self,
        dcdt_transfer_value_handle: RawHandle,
        call_value_handle: RawHandle,
    ) {
        let back_transfers = self.back_transfers_lock();
        let mut m_types = self.m_types_lock();
        m_types.bi_overwrite(call_value_handle, back_transfers.call_value.clone().into());
        m_types.mb_set_vec_of_dcdt_payments(
            dcdt_transfer_value_handle,
            &back_transfers.dcdt_transfers,
        );
    }

    fn check_dcdt_frozen(
        &self,
        address_handle: RawHandle,
        token_id_handle: RawHandle,
        _nonce: u64,
    ) -> bool {
        let address = VMAddress::from_slice(self.m_types_lock().mb_get(address_handle));
        let token_id_bytes = self.m_types_lock().mb_get(token_id_handle).to_vec();
        if let Some(account) = self.account_data(&address) {
            if let Some(dcdt_data) = account.dcdt.get_by_identifier(token_id_bytes.as_slice()) {
                return dcdt_data.frozen;
            }
        }

        false
    }

    fn get_dcdt_local_roles_bits(&self, token_id_handle: RawHandle) -> u64 {
        let token_id_bytes = self.m_types_lock().mb_get(token_id_handle).to_vec();
        let account = self.current_account_data();
        let mut result = DcdtLocalRoleFlags::NONE;
        if let Some(dcdt_data) = account.dcdt.get_by_identifier(token_id_bytes.as_slice()) {
            for role_name in dcdt_data.roles.get() {
                result |= DcdtLocalRole::from(role_name.as_slice()).to_flag();
            }
        }
        result.bits()
    }

    #[allow(clippy::too_many_arguments)]
    fn set_dcdt_data_values(
        &self,
        dcdt_data: &DcdtData,
        instance: &DcdtInstance,
        value_handle: RawHandle,
        properties_handle: RawHandle,
        hash_handle: RawHandle,
        name_handle: RawHandle,
        attributes_handle: RawHandle,
        creator_handle: RawHandle,
        royalties_handle: RawHandle,
        uris_handle: RawHandle,
    ) {
        let mut m_types = self.m_types_lock();
        m_types.bi_overwrite(value_handle, instance.balance.clone().into());
        if dcdt_data.frozen {
            m_types.mb_set(properties_handle, vec![1, 0]);
        } else {
            m_types.mb_set(properties_handle, vec![0, 0]);
        }
        m_types.mb_set(
            hash_handle,
            instance.metadata.hash.clone().unwrap_or_default(),
        );
        m_types.mb_set(name_handle, instance.metadata.name.clone());
        m_types.mb_set(attributes_handle, instance.metadata.attributes.clone());
        if let Some(creator) = &instance.metadata.creator {
            m_types.mb_set(creator_handle, creator.to_vec());
        } else {
            m_types.mb_set(creator_handle, vec![0u8; 32]);
        };
        m_types.bi_overwrite(
            royalties_handle,
            num_bigint::BigInt::from(instance.metadata.royalties),
        );
        m_types.mb_set_vec_of_bytes(uris_handle, instance.metadata.uri.clone());
    }

    #[allow(clippy::too_many_arguments)]
    fn reset_dcdt_data_values(
        &self,
        value_handle: RawHandle,
        properties_handle: RawHandle,
        hash_handle: RawHandle,
        name_handle: RawHandle,
        attributes_handle: RawHandle,
        creator_handle: RawHandle,
        royalties_handle: RawHandle,
        uris_handle: RawHandle,
    ) {
        if DCDT_TOKEN_DATA_FUNC_RESETS_VALUES {
            let mut m_types = self.m_types_lock();
            m_types.bi_overwrite(value_handle, BigInt::zero());
            m_types.mb_set(properties_handle, vec![0, 0]);
            m_types.mb_set(hash_handle, vec![]);
            m_types.mb_set(name_handle, vec![]);
            m_types.mb_set(attributes_handle, vec![]);
            m_types.mb_set(creator_handle, vec![0u8; 32]);
            m_types.bi_overwrite(royalties_handle, BigInt::zero());
            m_types.bi_overwrite(uris_handle, BigInt::zero());
        }
    }
}
