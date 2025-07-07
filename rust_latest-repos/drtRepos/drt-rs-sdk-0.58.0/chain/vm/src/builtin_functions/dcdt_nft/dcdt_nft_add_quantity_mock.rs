use dharitri_chain_core::types::ReturnCode;
use num_bigint::BigUint;

use crate::{
    blockchain::state::DcdtInstanceMetadata,
    chain_core::builtin_func_names::DCDT_NFT_ADD_QUANTITY_FUNC_NAME,
    host::context::{BlockchainUpdate, TxCache, TxInput, TxLog, TxResult},
    host::runtime::{RuntimeInstanceCallLambda, RuntimeRef},
    types::{top_decode_u64, top_encode_u64},
};

use super::super::builtin_func_trait::BuiltinFunction;

pub struct DCDTNftAddQuantity;

impl BuiltinFunction for DCDTNftAddQuantity {
    fn name(&self) -> &str {
        DCDT_NFT_ADD_QUANTITY_FUNC_NAME
    }

    fn execute<F>(
        &self,
        tx_input: TxInput,
        tx_cache: TxCache,
        _runtime: &RuntimeRef,
        _f: F,
    ) -> (TxResult, BlockchainUpdate)
    where
        F: RuntimeInstanceCallLambda,
    {
        if tx_input.args.len() != 3 {
            let err_result = TxResult::from_vm_error("DCDTNFTAddQuantity expects 3 arguments");
            return (err_result, BlockchainUpdate::empty());
        }

        let token_identifier = tx_input.args[0].clone();
        let nonce = top_decode_u64(tx_input.args[1].as_slice());
        let value = BigUint::from_bytes_be(tx_input.args[2].as_slice());

        tx_cache.increase_dcdt_balance(
            &tx_input.to,
            &token_identifier,
            nonce,
            &value,
            DcdtInstanceMetadata::default(),
        );

        let dcdt_nft_create_log = TxLog {
            address: tx_input.from,
            endpoint: DCDT_NFT_ADD_QUANTITY_FUNC_NAME.into(),
            topics: vec![
                token_identifier.to_vec(),
                top_encode_u64(nonce),
                value.to_bytes_be(),
            ],
            data: vec![],
        };

        let tx_result = TxResult {
            result_status: ReturnCode::Success,
            result_logs: vec![dcdt_nft_create_log],
            ..Default::default()
        };

        (tx_result, tx_cache.into_blockchain_updates())
    }
}
