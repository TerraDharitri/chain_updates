use crate::{
    scenario::model::BigUintValue,
    scenario_format::{
        interpret_trait::{InterpretableFrom, InterpreterContext},
        serde_raw::ValueSubTree,
    },
};

pub fn interpret_rewa_value(
    opt_legacy_value: Option<ValueSubTree>,
    opt_rewa_value: Option<ValueSubTree>,
    context: &InterpreterContext,
) -> BigUintValue {
    let mut rewa_value = BigUintValue::default();
    if let Some(parsed_legacy_value) = opt_legacy_value {
        rewa_value = BigUintValue::interpret_from(parsed_legacy_value, context);
    }
    if let Some(parsed_rewa_value) = opt_rewa_value {
        rewa_value = BigUintValue::interpret_from(parsed_rewa_value, context);
    }
    rewa_value
}
