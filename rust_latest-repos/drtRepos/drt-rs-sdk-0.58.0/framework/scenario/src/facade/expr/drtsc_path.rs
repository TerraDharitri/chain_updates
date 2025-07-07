use dharitri_chain_scenario_format::{
    interpret_trait::InterpreterContext, value_interpreter::interpret_string,
};
use dharitri_sc::types::{AnnotatedValue, ManagedBuffer, TxCodeValue};

use crate::{ScenarioTxEnv, ScenarioTxEnvData};

use super::RegisterCodeSource;

const DRTSC_PREFIX: &str = "drtsc:";

#[derive(Clone, Copy, Debug, PartialEq, Eq)]
pub struct DrtscPath<'a> {
    path: &'a str,
}

impl<'a> DrtscPath<'a> {
    pub const fn new(path: &'a str) -> Self {
        DrtscPath { path }
    }
}

impl DrtscPath<'_> {
    pub fn eval_to_expr(&self) -> String {
        format!("{DRTSC_PREFIX}{}", self.path)
    }

    pub fn resolve_contents(&self, context: &InterpreterContext) -> Vec<u8> {
        interpret_string(&format!("{DRTSC_PREFIX}{}", self.path), context)
    }
}

impl<Env> AnnotatedValue<Env, ManagedBuffer<Env::Api>> for DrtscPath<'_>
where
    Env: ScenarioTxEnv,
{
    fn annotation(&self, _env: &Env) -> ManagedBuffer<Env::Api> {
        self.eval_to_expr().into()
    }

    fn to_value(&self, env: &Env) -> ManagedBuffer<Env::Api> {
        self.resolve_contents(&env.env_data().interpreter_context())
            .into()
    }
}

impl<Env> TxCodeValue<Env> for DrtscPath<'_> where Env: ScenarioTxEnv {}

impl RegisterCodeSource for DrtscPath<'_> {
    fn into_code(self, env_data: ScenarioTxEnvData) -> Vec<u8> {
        self.resolve_contents(&env_data.interpreter_context())
    }
}

#[cfg(test)]
pub mod tests {
    use crate::imports::DrtscPath;

    fn assert_eq_eval(expr: &'static str, expected: &str) {
        assert_eq!(&DrtscPath::new(expr).eval_to_expr(), expected);
    }

    #[test]
    fn test_address_value() {
        assert_eq_eval("output/adder.drtsc.json", "drtsc:output/adder.drtsc.json");
    }
}
