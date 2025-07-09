use crate::{
    scenario::model::{CheckValue, U64Value},
    scenario_format::{
        interpret_trait::{InterpretableFrom, InterpreterContext, IntoRaw},
        serde_raw::CheckDcdtDataRaw,
    },
};

use super::CheckDcdtInstances;

#[derive(Debug, Default, Clone)]
pub struct CheckDcdtData {
    pub instances: CheckDcdtInstances,
    pub last_nonce: CheckValue<U64Value>,
    pub frozen: CheckValue<U64Value>,
}

impl InterpretableFrom<CheckDcdtDataRaw> for CheckDcdtData {
    fn interpret_from(from: CheckDcdtDataRaw, context: &InterpreterContext) -> Self {
        CheckDcdtData {
            instances: CheckDcdtInstances::interpret_from(from.instances, context),
            last_nonce: CheckValue::<U64Value>::interpret_from(from.last_nonce, context),
            frozen: CheckValue::<U64Value>::interpret_from(from.frozen, context),
        }
    }
}

impl IntoRaw<CheckDcdtDataRaw> for CheckDcdtData {
    fn into_raw(self) -> CheckDcdtDataRaw {
        CheckDcdtDataRaw {
            instances: self.instances.into_raw(),
            last_nonce: self.last_nonce.into_raw(),
            roles: Vec::new(),
            frozen: self.frozen.into_raw(),
        }
    }
}
