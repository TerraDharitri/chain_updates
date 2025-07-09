use crate::{
    scenario::model::{AddressValue, BigUintValue},
    scenario_format::{
        interpret_trait::{InterpretableFrom, InterpreterContext, IntoRaw},
        serde_raw::TxValidatorRewardRaw,
    },
};

use super::tx_interpret_util::interpret_rewa_value;

#[derive(Debug, Clone)]
pub struct TxValidatorReward {
    pub to: AddressValue,
    pub rewa_value: BigUintValue,
}

impl InterpretableFrom<TxValidatorRewardRaw> for TxValidatorReward {
    fn interpret_from(from: TxValidatorRewardRaw, context: &InterpreterContext) -> Self {
        TxValidatorReward {
            to: AddressValue::interpret_from(from.to, context),
            rewa_value: interpret_rewa_value(from.value, from.rewa_value, context),
        }
    }
}

impl IntoRaw<TxValidatorRewardRaw> for TxValidatorReward {
    fn into_raw(self) -> TxValidatorRewardRaw {
        TxValidatorRewardRaw {
            to: self.to.into_raw(),
            value: None,
            rewa_value: Some(self.rewa_value.into_raw()),
        }
    }
}
