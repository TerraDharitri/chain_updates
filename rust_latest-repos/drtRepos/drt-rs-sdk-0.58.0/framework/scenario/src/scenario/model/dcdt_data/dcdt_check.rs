use super::{CheckDcdtData, CheckDcdtInstance, CheckDcdtInstances};
use crate::{
    scenario::model::{BigUintValue, CheckValue, U64Value},
    scenario_format::{
        interpret_trait::{InterpretableFrom, InterpreterContext, IntoRaw},
        serde_raw::{CheckDcdtRaw, ValueSubTree},
    },
    scenario_model::BytesValue,
};
use num_bigint::BigUint;

#[derive(Debug, Clone)]
pub enum CheckDcdt {
    Short(BigUintValue),
    Full(CheckDcdtData),
}

impl CheckDcdt {
    pub fn convert_to_short_if_possible(&mut self) {
        if let CheckDcdt::Full(dcdt_check) = self {
            let has_single_fungible_instance =
                if let CheckDcdtInstances::Equal(check_instance) = &dcdt_check.instances {
                    check_instance.len() == 1 && check_instance[0].is_simple_fungible()
                } else {
                    false
                };

            if has_single_fungible_instance
                && dcdt_check.frozen.is_star()
                && dcdt_check.last_nonce.is_star()
            {
                let balance =
                    if let CheckDcdtInstances::Equal(check_instances) = &dcdt_check.instances {
                        match &check_instances[0].balance {
                            CheckValue::Star => BigUintValue {
                                original: ValueSubTree::Str("*".to_string()),
                                value: BigUint::from(0u32),
                            },
                            CheckValue::Equal(val) => val.clone(),
                        }
                    } else {
                        unreachable!();
                    };

                *self = CheckDcdt::Short(balance);
            }
        }
    }

    pub fn convert_to_full(&mut self) {
        if let CheckDcdt::Short(prev_balance_check) = self {
            let new_instances_check = vec![CheckDcdtInstance {
                balance: CheckValue::Equal(prev_balance_check.clone()),
                ..Default::default()
            }];

            let new_dcdt_check = CheckDcdtData {
                instances: CheckDcdtInstances::Equal(new_instances_check),
                ..Default::default()
            };
            *self = CheckDcdt::Full(new_dcdt_check);
        }
    }

    pub fn add_balance_check<N, V>(&mut self, nonce_expr: N, balance_expr: V)
    where
        U64Value: InterpretableFrom<N>,
        BigUintValue: InterpretableFrom<V>,
    {
        let ctx = InterpreterContext::default();
        let nonce = U64Value::interpret_from(nonce_expr, &ctx);
        let balance = BigUintValue::interpret_from(balance_expr, &ctx);

        self.convert_to_full();

        if let CheckDcdt::Full(prev_dcdt_check) = self {
            match &mut prev_dcdt_check.instances {
                CheckDcdtInstances::Star => {
                    let new_instances_check = vec![CheckDcdtInstance {
                        nonce,
                        balance: CheckValue::Equal(balance),
                        ..Default::default()
                    }];

                    prev_dcdt_check.instances = CheckDcdtInstances::Equal(new_instances_check);
                },
                CheckDcdtInstances::Equal(dcdt_instance_check) => {
                    if let Some(i) = dcdt_instance_check
                        .iter()
                        .position(|item| item.nonce.value == nonce.value)
                    {
                        dcdt_instance_check[i].balance = CheckValue::Equal(balance);
                    } else {
                        dcdt_instance_check.push(CheckDcdtInstance {
                            nonce,
                            balance: CheckValue::Equal(balance),
                            ..Default::default()
                        });
                    }
                },
            }
        }
    }

    pub fn add_balance_and_attributes_check<N, V, T>(
        &mut self,
        nonce_expr: N,
        balance_expr: V,
        attributes_expr: T,
    ) where
        U64Value: From<N>,
        BigUintValue: From<V>,
        BytesValue: From<T>,
    {
        let nonce = U64Value::from(nonce_expr);
        let balance = BigUintValue::from(balance_expr);
        let attributes = BytesValue::from(attributes_expr);

        self.convert_to_full();

        if let CheckDcdt::Full(prev_dcdt_check) = self {
            match &mut prev_dcdt_check.instances {
                CheckDcdtInstances::Star => {
                    let new_instances_check = vec![CheckDcdtInstance {
                        nonce,
                        balance: CheckValue::Equal(balance),
                        attributes: CheckValue::Equal(attributes),
                        ..Default::default()
                    }];

                    prev_dcdt_check.instances = CheckDcdtInstances::Equal(new_instances_check);
                },
                CheckDcdtInstances::Equal(dcdt_instance_check) => {
                    if let Some(i) = dcdt_instance_check
                        .iter()
                        .position(|item| item.nonce.value == nonce.value)
                    {
                        dcdt_instance_check[i].balance = CheckValue::Equal(balance);
                        dcdt_instance_check[i].attributes = CheckValue::Equal(attributes);
                    } else {
                        dcdt_instance_check.push(CheckDcdtInstance {
                            nonce,
                            balance: CheckValue::Equal(balance),
                            attributes: CheckValue::Equal(attributes),
                            ..Default::default()
                        });
                    }
                },
            }
        }
    }
}

impl InterpretableFrom<CheckDcdtRaw> for CheckDcdt {
    fn interpret_from(from: CheckDcdtRaw, context: &InterpreterContext) -> Self {
        match from {
            CheckDcdtRaw::Full(m) => CheckDcdt::Full(CheckDcdtData::interpret_from(m, context)),
            CheckDcdtRaw::Short(v) => CheckDcdt::Short(BigUintValue::interpret_from(v, context)),
        }
    }
}

impl IntoRaw<CheckDcdtRaw> for CheckDcdt {
    fn into_raw(mut self) -> CheckDcdtRaw {
        self.convert_to_short_if_possible();

        match self {
            CheckDcdt::Full(m) => CheckDcdtRaw::Full(m.into_raw()),
            CheckDcdt::Short(v) => CheckDcdtRaw::Short(v.into_raw()),
        }
    }
}
