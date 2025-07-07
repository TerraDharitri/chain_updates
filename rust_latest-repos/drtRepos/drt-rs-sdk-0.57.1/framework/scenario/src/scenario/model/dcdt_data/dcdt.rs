use crate::scenario_format::{
    interpret_trait::{InterpretableFrom, InterpreterContext, IntoRaw},
    serde_raw::{DcdtFullRaw, DcdtRaw},
};

use super::{DcdtInstance, DcdtObject};
use crate::scenario::model::{BigUintValue, BytesValue, U64Value};

#[derive(Debug, Clone)]
pub enum Dcdt {
    Short(BigUintValue),
    Full(DcdtObject),
}

impl Dcdt {
    pub fn convert_to_short_if_possible(&mut self) {
        if let Dcdt::Full(dcdt_obj) = self {
            if dcdt_obj.is_short_form() {
                *self = Self::Short(dcdt_obj.instances[0].balance.clone().unwrap())
            }
        }
    }

    pub fn convert_to_full(&mut self) {
        if let Dcdt::Short(balance) = self {
            let mut new_dcdt_obj = DcdtObject::default();
            new_dcdt_obj.set_balance(0u64, balance.clone());

            *self = Self::Full(new_dcdt_obj);
        }
    }

    pub fn set_balance<N, A>(&mut self, token_nonce_expr: N, amount_expr: A)
    where
        U64Value: From<N>,
        BigUintValue: From<A>,
    {
        self.convert_to_full();

        if let Dcdt::Full(dcdt_obj) = self {
            dcdt_obj.set_balance(token_nonce_expr, amount_expr);
        }
    }

    pub fn get_mut_dcdt_object(&mut self) -> &mut DcdtObject {
        self.convert_to_full();

        if let Dcdt::Full(dcdt_obj) = self {
            return dcdt_obj;
        }

        unreachable!()
    }
}

impl InterpretableFrom<DcdtRaw> for Dcdt {
    fn interpret_from(from: DcdtRaw, context: &InterpreterContext) -> Self {
        match from {
            DcdtRaw::Short(short_dcdt) => {
                Dcdt::Short(BigUintValue::interpret_from(short_dcdt, context))
            },
            DcdtRaw::Full(full_dcdt) => Dcdt::Full(DcdtObject {
                token_identifier: full_dcdt
                    .token_identifier
                    .map(|b| BytesValue::interpret_from(b, context)),
                instances: full_dcdt
                    .instances
                    .into_iter()
                    .map(|instance| DcdtInstance::interpret_from(instance, context))
                    .collect(),
                last_nonce: full_dcdt
                    .last_nonce
                    .map(|b| U64Value::interpret_from(b, context)),
                roles: full_dcdt.roles,
                frozen: full_dcdt
                    .frozen
                    .map(|b| U64Value::interpret_from(b, context)),
            }),
        }
    }
}

impl IntoRaw<DcdtRaw> for Dcdt {
    fn into_raw(mut self) -> DcdtRaw {
        self.convert_to_short_if_possible();

        match self {
            Dcdt::Short(short) => DcdtRaw::Short(short.original),
            Dcdt::Full(eo) => DcdtRaw::Full(DcdtFullRaw {
                token_identifier: eo.token_identifier.map(|ti| ti.original),
                instances: eo
                    .instances
                    .into_iter()
                    .map(|inst| inst.into_raw())
                    .collect(),
                last_nonce: eo.last_nonce.map(|ti| ti.original),
                roles: eo.roles,
                frozen: eo.frozen.map(|ti| ti.original),
            }),
        }
    }
}
