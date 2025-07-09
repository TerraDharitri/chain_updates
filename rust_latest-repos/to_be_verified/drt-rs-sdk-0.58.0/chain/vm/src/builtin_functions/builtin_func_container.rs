use super::{
    builtin_func_trait::BuiltinFunction,
    dcdt_nft::{
        DCDTLocalBurn, DCDTLocalMint, DCDTNftAddQuantity, DCDTNftAddUri, DCDTNftBurn,
        DCDTNftCreate, DCDTNftUpdateAttributes,
    },
    general::{ChangeOwner, ClaimDeveloperRewards, DeleteUsername, SetUsername, UpgradeContract},
    transfer::{DCDTMultiTransfer, DCDTNftTransfer, DCDTTransfer},
    BuiltinFunctionDcdtTransferInfo,
};

use crate::{
    host::context::{BlockchainUpdate, TxCache, TxInput, TxResult},
    host::runtime::{RuntimeInstanceCallLambda, RuntimeRef},
    types::DcdtLocalRole,
};

use crate::chain_core::builtin_func_names::*;

/// Container for builtin function logic.
///
/// Currently has no data, but could conceivably be configurable in the future.
#[derive(Default)]
pub struct BuiltinFunctionContainer;

impl BuiltinFunctionContainer {
    /// If the call points to a builtin function, it executes it, otherwise calls the `or_else` closure.
    ///
    /// It also checks that the appropriate roles are set, where applicable.
    pub fn execute_builtin_function_or_else<F, Else>(
        &self,
        runtime: &RuntimeRef,
        tx_input: TxInput,
        tx_cache: TxCache,
        f: F,
        or_else: Else,
    ) -> (TxResult, BlockchainUpdate)
    where
        F: RuntimeInstanceCallLambda,
        Else: FnOnce(TxInput, TxCache, F) -> (TxResult, BlockchainUpdate),
    {
        BuiltinFunctionCall::new(runtime, tx_input, tx_cache).execute_or_else(f, or_else)
    }

    /// Provides data on the builtin functions that perform DCDT token transfers.
    pub fn extract_token_transfers(&self, tx_input: &TxInput) -> BuiltinFunctionDcdtTransferInfo {
        match tx_input.func_name.as_str() {
            DCDT_MULTI_TRANSFER_FUNC_NAME => bf_extract_transfers(DCDTMultiTransfer, tx_input),
            DCDT_NFT_TRANSFER_FUNC_NAME => bf_extract_transfers(DCDTNftTransfer, tx_input),
            DCDT_TRANSFER_FUNC_NAME => bf_extract_transfers(DCDTTransfer, tx_input),
            _ => BuiltinFunctionDcdtTransferInfo::empty(tx_input),
        }
    }
}

fn bf_extract_transfers<B>(builtin_func: B, tx_input: &TxInput) -> BuiltinFunctionDcdtTransferInfo
where
    B: BuiltinFunction,
{
    builtin_func.extract_dcdt_transfers(tx_input)
}

/// Syntax helper for the big builtin function match in `execute_or_else`.
/// Thanks to it we do not need to write out the arguments for each match arm.
struct BuiltinFunctionCall<'a> {
    runtime: &'a RuntimeRef,
    tx_input: TxInput,
    tx_cache: TxCache,
}

impl<'a> BuiltinFunctionCall<'a> {
    pub fn new(runtime: &'a RuntimeRef, tx_input: TxInput, tx_cache: TxCache) -> Self {
        BuiltinFunctionCall {
            runtime,
            tx_input,
            tx_cache,
        }
    }

    pub fn execute_or_else<F, Else>(self, f: F, or_else: Else) -> (TxResult, BlockchainUpdate)
    where
        F: RuntimeInstanceCallLambda,
        Else: FnOnce(TxInput, TxCache, F) -> (TxResult, BlockchainUpdate),
    {
        match self.tx_input.func_name.as_str() {
            DCDT_LOCAL_MINT_FUNC_NAME => {
                self.check_role_and_execute(DcdtLocalRole::Mint, DCDTLocalMint, f)
            },
            DCDT_LOCAL_BURN_FUNC_NAME => {
                self.check_role_and_execute(DcdtLocalRole::Burn, DCDTLocalBurn, f)
            },
            DCDT_NFT_CREATE_FUNC_NAME => {
                self.check_role_and_execute(DcdtLocalRole::NftCreate, DCDTNftCreate, f)
            },
            DCDT_NFT_BURN_FUNC_NAME => {
                self.check_role_and_execute(DcdtLocalRole::NftBurn, DCDTNftBurn, f)
            },
            DCDT_NFT_ADD_QUANTITY_FUNC_NAME => {
                self.check_role_and_execute(DcdtLocalRole::NftAddQuantity, DCDTNftAddQuantity, f)
            },
            DCDT_NFT_ADD_URI_FUNC_NAME => {
                self.check_role_and_execute(DcdtLocalRole::NftAddUri, DCDTNftAddUri, f)
            },
            DCDT_NFT_UPDATE_ATTRIBUTES_FUNC_NAME => self.check_role_and_execute(
                DcdtLocalRole::NftUpdateAttributes,
                DCDTNftUpdateAttributes,
                f,
            ),

            DCDT_MULTI_TRANSFER_FUNC_NAME => self.execute_bf(DCDTMultiTransfer, f),
            DCDT_NFT_TRANSFER_FUNC_NAME => self.execute_bf(DCDTNftTransfer, f),
            DCDT_TRANSFER_FUNC_NAME => self.execute_bf(DCDTTransfer, f),
            CHANGE_OWNER_BUILTIN_FUNC_NAME => self.execute_bf(ChangeOwner, f),
            CLAIM_DEVELOPER_REWARDS_FUNC_NAME => self.execute_bf(ClaimDeveloperRewards, f),
            SET_USERNAME_FUNC_NAME => self.execute_bf(SetUsername, f),
            DELETE_USERNAME_FUNC_NAME => self.execute_bf(DeleteUsername, f),
            UPGRADE_CONTRACT_FUNC_NAME => self.execute_bf(UpgradeContract, f),
            MIGRATE_USERNAME_FUNC_NAME => {
                panic!("builtin function {MIGRATE_USERNAME_FUNC_NAME} was dropped")
            },
            _ => or_else(self.tx_input, self.tx_cache, f),
        }
    }

    fn execute_bf<B, F>(self, builtin_func: B, f: F) -> (TxResult, BlockchainUpdate)
    where
        B: BuiltinFunction,
        F: RuntimeInstanceCallLambda,
    {
        builtin_func.execute(self.tx_input, self.tx_cache, self.runtime, f)
    }

    fn check_role_and_execute<B, F>(
        self,
        role: DcdtLocalRole,
        builtin_func: B,
        f: F,
    ) -> (TxResult, BlockchainUpdate)
    where
        B: BuiltinFunction,
        F: RuntimeInstanceCallLambda,
    {
        if check_allowed_to_execute(role, &self.tx_input, &self.tx_cache) {
            self.execute_bf(builtin_func, f)
        } else {
            (
                TxResult::from_vm_error("action is not allowed"),
                BlockchainUpdate::empty(),
            )
        }
    }
}

fn check_allowed_to_execute(role: DcdtLocalRole, tx_input: &TxInput, tx_cache: &TxCache) -> bool {
    let token_identifier = tx_input.args[0].clone();
    let available_roles = tx_cache.with_account_mut(&tx_input.to, |account| {
        account.dcdt.get_roles(&token_identifier)
    });
    available_roles
        .iter()
        .any(|available_role| available_role.as_slice() == role.name().as_bytes())
}
