use crate::types::{AnnotatedValue, BigUint, ManagedRef};

use super::TxEnv;

pub trait TxRewaValue<Env>: AnnotatedValue<Env, BigUint<Env::Api>>
where
    Env: TxEnv,
{
}

impl<Env> TxRewaValue<Env> for BigUint<Env::Api> where Env: TxEnv {}
impl<Env> TxRewaValue<Env> for &BigUint<Env::Api> where Env: TxEnv {}
impl<Env> TxRewaValue<Env> for ManagedRef<'_, Env::Api, BigUint<Env::Api>> where Env: TxEnv {}
impl<Env> TxRewaValue<Env> for u64 where Env: TxEnv {}
