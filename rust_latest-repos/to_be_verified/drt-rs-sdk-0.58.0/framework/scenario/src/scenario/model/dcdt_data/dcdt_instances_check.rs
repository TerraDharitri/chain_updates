use crate::scenario_format::{
    interpret_trait::{InterpretableFrom, InterpreterContext, IntoRaw},
    serde_raw::CheckDcdtInstancesRaw,
};

use super::CheckDcdtInstance;

#[derive(Debug, Clone)]
pub enum CheckDcdtInstances {
    Star,
    Equal(Vec<CheckDcdtInstance>),
}

impl CheckDcdtInstances {
    pub fn is_star(&self) -> bool {
        matches!(self, CheckDcdtInstances::Star)
    }

    pub fn contains_nonce(&self, nonce: u64) -> bool {
        match &self {
            CheckDcdtInstances::Equal(eq) => {
                for expected_value in eq.iter() {
                    if expected_value.nonce.value == nonce {
                        return true;
                    }
                }
            },
            CheckDcdtInstances::Star => {},
        }
        false
    }
}

impl Default for CheckDcdtInstances {
    fn default() -> Self {
        CheckDcdtInstances::Equal(Vec::new())
    }
}

impl InterpretableFrom<CheckDcdtInstancesRaw> for CheckDcdtInstances {
    fn interpret_from(from: CheckDcdtInstancesRaw, context: &InterpreterContext) -> Self {
        match from {
            CheckDcdtInstancesRaw::Unspecified => CheckDcdtInstances::Star,
            CheckDcdtInstancesRaw::Star => CheckDcdtInstances::Star,
            CheckDcdtInstancesRaw::Equal(m) => CheckDcdtInstances::Equal(
                m.into_iter()
                    .map(|v| CheckDcdtInstance::interpret_from(v, context))
                    .collect(),
            ),
        }
    }
}

impl IntoRaw<CheckDcdtInstancesRaw> for CheckDcdtInstances {
    fn into_raw(self) -> CheckDcdtInstancesRaw {
        match self {
            CheckDcdtInstances::Equal(eq) => {
                CheckDcdtInstancesRaw::Equal(eq.into_iter().map(|cei| cei.into_raw()).collect())
            },
            CheckDcdtInstances::Star => CheckDcdtInstancesRaw::Star,
        }
    }
}
