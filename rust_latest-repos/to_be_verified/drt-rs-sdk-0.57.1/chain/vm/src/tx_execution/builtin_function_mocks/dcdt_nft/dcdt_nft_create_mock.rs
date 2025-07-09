use dharitri_chain_core::types::ReturnCode;
use num_bigint::BigUint;

use crate::{
    chain_core::builtin_func_names::DCDT_NFT_CREATE_FUNC_NAME,
    tx_execution::BlockchainVMRef,
    tx_mock::{BlockchainUpdate, TxCache, TxInput, TxLog, TxResult},
    types::{top_decode_u64, top_encode_u64},
    world_mock::{DcdtInstance, DcdtInstanceMetadata},
};

use super::super::builtin_func_trait::BuiltinFunction;

pub struct DCDTNftCreate;

impl BuiltinFunction for DCDTNftCreate {
    fn name(&self) -> &str {
        DCDT_NFT_CREATE_FUNC_NAME
    }

    fn execute<F>(
        &self,
        tx_input: TxInput,
        tx_cache: TxCache,
        _vm: &BlockchainVMRef,
        _f: F,
    ) -> (TxResult, BlockchainUpdate)
    where
        F: FnOnce(),
    {
        if tx_input.args.len() < 7 {
            let err_result = TxResult::from_vm_error("DCDTNFTCreate too few arguments");
            return (err_result, BlockchainUpdate::empty());
        }
        assert!(
            tx_input.to == tx_input.from,
            "DCDTNFTCreate expects that to == from"
        );

        let token_identifier = tx_input.args[0].as_slice();
        let amount = BigUint::from_bytes_be(tx_input.args[1].as_slice());
        let name = tx_input.args[2].clone();
        let royalties = top_decode_u64(tx_input.args[3].as_slice());
        let hash = tx_input.args[4].clone();
        let attributes = tx_input.args[5].clone();
        let uris = tx_input.args[6..].to_vec();

        let new_nonce = tx_cache.with_account_mut(&tx_input.to, |account| {
            let dcdt_data = account
                .dcdt
                .get_mut_by_identifier(token_identifier)
                .unwrap_or_else(|| panic!("DCDTNFTCreate unexpected token identifier"));
            dcdt_data.last_nonce += 1;
            dcdt_data.instances.push_instance(DcdtInstance {
                nonce: dcdt_data.last_nonce,
                balance: amount.clone(),
                metadata: DcdtInstanceMetadata {
                    name,
                    creator: Some(tx_input.from.clone()),
                    royalties,
                    hash: Some(hash),
                    uri: uris,
                    attributes,
                },
            });

            dcdt_data.last_nonce
        });

        let dcdt_nft_create_log = TxLog {
            address: tx_input.from.clone(),
            endpoint: DCDT_NFT_CREATE_FUNC_NAME.into(),
            topics: vec![
                token_identifier.to_vec(),
                top_encode_u64(new_nonce),
                amount.to_bytes_be(),
                Vec::new(), // in actuality this should contain the fully marshalled DCDT data
            ],
            data: vec![],
        };

        let tx_result = TxResult {
            result_status: ReturnCode::Success,
            result_values: vec![top_encode_u64(new_nonce)],
            result_logs: vec![dcdt_nft_create_log],
            ..Default::default()
        };

        (tx_result, tx_cache.into_blockchain_updates())
    }
}
