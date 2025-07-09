use dharitri_chain_vm::types::{top_encode_big_uint, top_encode_u64};
use dharitri_sc::api::{
    DCDT_MULTI_TRANSFER_FUNC_NAME, DCDT_NFT_TRANSFER_FUNC_NAME, DCDT_TRANSFER_FUNC_NAME,
};

use crate::{
    api::StaticApi,
    dharitri_sc::types::{ContractCall, DcdtTokenPayment},
    scenario::model::{AddressValue, BigUintValue, BytesValue, U64Value},
    scenario_format::{
        interpret_trait::{InterpretableFrom, InterpreterContext, IntoRaw},
        serde_raw::TxCallRaw,
    },
};

use super::{tx_interpret_util::interpret_rewa_value, TxDCDT};

pub const DEFAULT_GAS_EXPR: &str = "5,000,000";

#[derive(Debug, Clone)]
pub struct TxCall {
    pub from: AddressValue,
    pub to: AddressValue,
    pub rewa_value: BigUintValue,
    pub dcdt_value: Vec<TxDCDT>,
    pub function: String,
    pub arguments: Vec<BytesValue>,
    pub gas_limit: U64Value,
    pub gas_price: U64Value,
}

impl Default for TxCall {
    fn default() -> Self {
        Self {
            from: Default::default(),
            to: Default::default(),
            rewa_value: Default::default(),
            dcdt_value: Default::default(),
            function: Default::default(),
            arguments: Default::default(),
            gas_limit: U64Value::from(DEFAULT_GAS_EXPR),
            gas_price: Default::default(),
        }
    }
}

impl InterpretableFrom<TxCallRaw> for TxCall {
    fn interpret_from(from: TxCallRaw, context: &InterpreterContext) -> Self {
        TxCall {
            from: AddressValue::interpret_from(from.from, context),
            to: AddressValue::interpret_from(from.to, context),
            rewa_value: interpret_rewa_value(from.value, from.rewa_value, context),
            dcdt_value: from
                .dcdt_value
                .into_iter()
                .map(|dcdt_value| TxDCDT::interpret_from(dcdt_value, context))
                .collect(),
            function: from.function,
            arguments: from
                .arguments
                .into_iter()
                .map(|t| BytesValue::interpret_from(t, context))
                .collect(),
            gas_limit: U64Value::interpret_from(from.gas_limit, context),
            gas_price: U64Value::interpret_from(from.gas_price.unwrap_or_default(), context),
        }
    }
}

impl IntoRaw<TxCallRaw> for TxCall {
    fn into_raw(self) -> TxCallRaw {
        TxCallRaw {
            from: self.from.into_raw(),
            to: self.to.into_raw(),
            value: None,
            rewa_value: self.rewa_value.into_raw_opt(),
            dcdt_value: self
                .dcdt_value
                .into_iter()
                .map(|dcdt_value| dcdt_value.into_raw())
                .collect(),
            function: self.function,
            arguments: self
                .arguments
                .into_iter()
                .map(|arg| arg.into_raw())
                .collect(),
            gas_limit: self.gas_limit.into_raw(),
            gas_price: self.gas_price.into_raw_opt(),
        }
    }
}

impl TxCall {
    #[deprecated(
        since = "0.49.0",
        note = "Please use the unified transaction syntax instead."
    )]
    #[allow(deprecated)]
    pub fn to_contract_call(&self) -> dharitri_sc::types::ContractCallWithRewa<StaticApi, ()> {
        let mut contract_call = dharitri_sc::types::ContractCallWithRewa::new(
            (&self.to.value).into(),
            self.function.as_bytes(),
            (&self.rewa_value.value).into(),
        );

        contract_call.basic.explicit_gas_limit = self.gas_limit.value;

        contract_call = contract_call.convert_to_dcdt_transfer_call(
            self.dcdt_value
                .iter()
                .map(|dcdt| {
                    DcdtTokenPayment::new(
                        dcdt.dcdt_token_identifier.value.as_slice().into(),
                        dcdt.nonce.value,
                        (&dcdt.dcdt_value.value).into(),
                    )
                })
                .collect(),
        );

        // For some contract calls from == to.
        // The contract call objects have no "from" field, since that is always part of the execution context.
        // On the static API there is no execution context, but a placeholder value is provided.
        // Here we already know the sender, so we can replace the placeholder with the actual value.
        if StaticApi::is_current_address_placeholder(&contract_call.basic.to.to_address()) {
            contract_call.basic.to = self.from.value.clone().into();
        }

        for argument in &self.arguments {
            contract_call.push_raw_argument(argument.value.as_slice());
        }
        contract_call
    }

    /// Converts call to builtin function DCDT transfer call, if necessary.
    pub fn normalize(&self) -> TxCall {
        let (function, arguments, to_self) = self.process_payments();
        TxCall {
            from: self.from.clone(),
            to: if to_self {
                self.from.clone()
            } else {
                self.to.clone()
            },
            rewa_value: self.rewa_value.clone(),
            dcdt_value: Vec::new(),
            function,
            arguments,
            gas_limit: self.gas_limit.clone(),
            gas_price: self.gas_price.clone(),
        }
    }

    fn process_payments(&self) -> (String, Vec<BytesValue>, bool) {
        assert!(
            self.rewa_value.is_zero() || self.dcdt_value.is_empty(),
            "Cannot have both REWA and DCDT fields filled. To transfer REWA and DCDT in the same transaction, represent REWA as REWA-000000 in the DCDTs.");

        match self.dcdt_value.len() {
            0 => (self.function.clone(), self.arguments.clone(), false),
            1 => {
                let payment = self.dcdt_value.first().unwrap();
                if payment.is_rewa() {
                    self.construct_multi_transfer_dcdt_call()
                } else if payment.nonce.value == 0 {
                    self.construct_single_transfer_fungible_call(payment)
                } else {
                    self.construct_single_transfer_nft_call(payment)
                }
            },
            _ => self.construct_multi_transfer_dcdt_call(),
        }
    }

    fn append_function_call_to_arguments(&self, arguments: &mut Vec<BytesValue>) {
        if !self.function.is_empty() {
            arguments.push(BytesValue::from_str_expr(&self.function));
        }
        for regular_arg in &self.arguments {
            arguments.push(regular_arg.clone())
        }
    }

    /// Constructs `DCDTTransfer` builtin function call.
    pub(crate) fn construct_single_transfer_fungible_call(
        &self,
        payment: &TxDCDT,
    ) -> (String, Vec<BytesValue>, bool) {
        let mut arguments = Vec::new();
        arguments.push(payment.dcdt_token_identifier.value.as_slice().into());
        arguments.push(top_encode_big_uint(&payment.dcdt_value.value).into());
        self.append_function_call_to_arguments(&mut arguments);

        (DCDT_TRANSFER_FUNC_NAME.to_owned(), arguments, false)
    }

    /// Constructs `DCDTNFTTransfer` builtin function call.
    ///
    /// `DCDTNFTTransfer` takes 4 arguments:
    /// arg0 - token identifier
    /// arg1 - nonce
    /// arg2 - quantity to transfer
    /// arg3 - destination address
    pub(crate) fn construct_single_transfer_nft_call(
        &self,
        payment: &TxDCDT,
    ) -> (String, Vec<BytesValue>, bool) {
        let mut arguments = vec![
            payment.dcdt_token_identifier.value.as_slice().into(),
            top_encode_u64(payment.nonce.value).into(),
            top_encode_big_uint(&payment.dcdt_value.value).into(),
            self.to.value.as_bytes().into(), // TODO: preserve representation
        ];

        self.append_function_call_to_arguments(&mut arguments);

        (DCDT_NFT_TRANSFER_FUNC_NAME.to_owned(), arguments, true)
    }

    /// Constructs `MultiDCDTNFTTransfer` builtin function call.
    pub(crate) fn construct_multi_transfer_dcdt_call(&self) -> (String, Vec<BytesValue>, bool) {
        let mut arguments = Vec::new();
        arguments.push(self.to.value.as_bytes().into()); // TODO: preserve representation
        arguments.push(top_encode_u64(self.dcdt_value.len() as u64).into());

        for payment in &self.dcdt_value {
            arguments.push(payment.dcdt_token_identifier.value.as_slice().into());
            arguments.push(top_encode_u64(payment.nonce.value).into());
            arguments.push(top_encode_big_uint(&payment.dcdt_value.value).into());
        }

        self.append_function_call_to_arguments(&mut arguments);

        (DCDT_MULTI_TRANSFER_FUNC_NAME.to_owned(), arguments, true)
    }

    /// Creates the data field of the transaction represented by object.
    pub fn compute_data_field(&self) -> String {
        let mut result = self.function.clone();
        for argument in &self.arguments {
            result.push('@');
            result.push_str(hex::encode(argument.value.as_slice()).as_str());
        }
        result
    }
}
