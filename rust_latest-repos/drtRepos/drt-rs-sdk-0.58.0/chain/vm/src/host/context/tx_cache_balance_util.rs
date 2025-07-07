use dharitri_chain_core::REWA_000000_TOKEN_IDENTIFIER;
use num_bigint::BigUint;

use crate::{
    blockchain::state::DcdtInstanceMetadata, host::context::TxPanic,
    system_sc::is_system_sc_address, types::VMAddress,
};

use super::TxCache;

impl TxCache {
    pub fn subtract_rewa_balance(
        &self,
        address: &VMAddress,
        call_value: &BigUint,
    ) -> Result<(), TxPanic> {
        self.with_account_mut(address, |account| {
            if call_value > &account.rewa_balance {
                return Err(TxPanic::vm_error("failed transfer (insufficient funds)"));
            }
            account.rewa_balance -= call_value;
            Ok(())
        })
    }

    pub fn subtract_tx_gas(&self, address: &VMAddress, gas_limit: u64, gas_price: u64) {
        self.with_account_mut(address, |account| {
            let gas_cost = BigUint::from(gas_limit) * BigUint::from(gas_price);
            assert!(
                account.rewa_balance >= gas_cost,
                "Not enough balance to pay gas upfront"
            );
            account.rewa_balance -= &gas_cost;
        });
    }

    pub fn increase_rewa_balance(&self, address: &VMAddress, amount: &BigUint) {
        self.with_account_mut(address, |account| {
            account.rewa_balance += amount;
        });
    }

    pub fn subtract_dcdt_balance(
        &self,
        address: &VMAddress,
        dcdt_token_identifier: &[u8],
        nonce: u64,
        value: &BigUint,
    ) -> Result<DcdtInstanceMetadata, TxPanic> {
        self.with_account_mut(address, |account| {
            let dcdt_data_map = &mut account.dcdt;
            let dcdt_data = dcdt_data_map
                .get_mut_by_identifier(dcdt_token_identifier)
                .ok_or_else(err_insufficient_funds)?;

            let dcdt_instances = &mut dcdt_data.instances;
            let dcdt_instance = dcdt_instances
                .get_mut_by_nonce(nonce)
                .ok_or_else(err_insufficient_funds)?;

            let dcdt_balance = &mut dcdt_instance.balance;
            if &*dcdt_balance < value {
                return Err(err_insufficient_funds());
            }

            *dcdt_balance -= value;

            Ok(dcdt_instance.metadata.clone())
        })
    }

    pub fn increase_dcdt_balance(
        &self,
        address: &VMAddress,
        dcdt_token_identifier: &[u8],
        nonce: u64,
        value: &BigUint,
        dcdt_metadata: DcdtInstanceMetadata,
    ) {
        self.with_account_mut(address, |account| {
            account.dcdt.increase_balance(
                dcdt_token_identifier.to_vec(),
                nonce,
                value,
                dcdt_metadata,
            );
        });
    }

    pub fn transfer_rewa_balance(
        &self,
        from: &VMAddress,
        to: &VMAddress,
        value: &BigUint,
    ) -> Result<(), TxPanic> {
        if !is_system_sc_address(from) {
            self.subtract_rewa_balance(from, value)?;
        }
        if !is_system_sc_address(to) {
            self.increase_rewa_balance(to, value);
        }
        Ok(())
    }

    pub fn transfer_dcdt_balance(
        &self,
        from: &VMAddress,
        to: &VMAddress,
        dcdt_token_identifier: &[u8],
        nonce: u64,
        value: &BigUint,
    ) -> Result<(), TxPanic> {
        if dcdt_token_identifier == REWA_000000_TOKEN_IDENTIFIER.as_bytes() {
            return self.transfer_rewa_balance(from, to, value);
        }

        if !is_system_sc_address(from) && !is_system_sc_address(to) {
            let metadata = self.subtract_dcdt_balance(from, dcdt_token_identifier, nonce, value)?;
            self.increase_dcdt_balance(to, dcdt_token_identifier, nonce, value, metadata);
        }
        Ok(())
    }
}

fn err_insufficient_funds() -> TxPanic {
    TxPanic::vm_error("insufficient funds")
}
