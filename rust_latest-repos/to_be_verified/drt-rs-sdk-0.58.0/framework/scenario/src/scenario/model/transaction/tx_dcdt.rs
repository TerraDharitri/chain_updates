use dharitri_sc::{
    api::ManagedTypeApi,
    chain_core::REWA_000000_TOKEN_IDENTIFIER,
    types::{RewaOrDcdtTokenPayment, DcdtTokenPayment},
};

use crate::{
    scenario::model::{BigUintValue, BytesValue, U64Value},
    scenario_format::{
        interpret_trait::{InterpretableFrom, InterpreterContext, IntoRaw},
        serde_raw::{TxDCDTRaw, ValueSubTree},
    },
};

#[derive(Debug, Clone)]
pub struct TxDCDT {
    pub dcdt_token_identifier: BytesValue,
    pub nonce: U64Value,
    pub dcdt_value: BigUintValue,
}

impl TxDCDT {
    pub fn is_rewa(&self) -> bool {
        self.dcdt_token_identifier.value == REWA_000000_TOKEN_IDENTIFIER.as_bytes()
    }
}

impl InterpretableFrom<TxDCDTRaw> for TxDCDT {
    fn interpret_from(from: TxDCDTRaw, context: &InterpreterContext) -> Self {
        TxDCDT {
            dcdt_token_identifier: interpret_dcdt_token_identifier(from.token_identifier, context),
            nonce: interpret_opt_u64(from.nonce, context),
            dcdt_value: BigUintValue::interpret_from(from.value, context),
        }
    }
}

impl IntoRaw<TxDCDTRaw> for TxDCDT {
    fn into_raw(self) -> TxDCDTRaw {
        TxDCDTRaw {
            token_identifier: Some(self.dcdt_token_identifier.into_raw()),
            nonce: self.nonce.into_raw_opt(),
            value: self.dcdt_value.into_raw(),
        }
    }
}

impl<M: ManagedTypeApi> From<DcdtTokenPayment<M>> for TxDCDT {
    fn from(value: DcdtTokenPayment<M>) -> Self {
        TxDCDT {
            dcdt_token_identifier: BytesValue::from(
                value.token_identifier.as_managed_buffer().to_vec(),
            ),
            nonce: U64Value::from(value.token_nonce),
            dcdt_value: BigUintValue::from(value.amount),
        }
    }
}

impl<M: ManagedTypeApi> From<RewaOrDcdtTokenPayment<M>> for TxDCDT {
    fn from(value: RewaOrDcdtTokenPayment<M>) -> Self {
        TxDCDT {
            dcdt_token_identifier: BytesValue::from(
                value.token_identifier.as_managed_buffer().to_vec(),
            ),
            nonce: U64Value::from(value.token_nonce),
            dcdt_value: BigUintValue::from(value.amount),
        }
    }
}

fn interpret_dcdt_token_identifier(
    dcdt_token_identifier: Option<ValueSubTree>,
    context: &InterpreterContext,
) -> BytesValue {
    if let Some(dcdt_token_identifier_raw) = dcdt_token_identifier {
        BytesValue::interpret_from(dcdt_token_identifier_raw, context)
    } else {
        BytesValue::empty()
    }
}

fn interpret_opt_u64(opt_u64: Option<ValueSubTree>, context: &InterpreterContext) -> U64Value {
    if let Some(u) = opt_u64 {
        U64Value::interpret_from(u, context)
    } else {
        U64Value::empty()
    }
}
