use super::{DcdtTokenPayment, ManagedVec};
use crate::{
    api::ManagedTypeApi,
    codec::{
        self,
        derive::{NestedDecode, NestedEncode, TopDecode, TopEncode},
    },
    types::BigUint,
};

use crate as dharitri_sc; // needed by the TypeAbi generated code
use crate::derive::type_abi;

/// Encodes any type of payment, which either:
/// - REWA (can be zero in case of no payment whatsoever);
/// - Multi-DCDT (one or more DCDT transfers).
#[type_abi]
#[derive(TopDecode, TopEncode, NestedDecode, NestedEncode, Clone, PartialEq, Eq, Debug)]
pub enum RewaOrMultiDcdtPayment<M: ManagedTypeApi> {
    Rewa(BigUint<M>),
    MultiDcdt(ManagedVec<M, DcdtTokenPayment<M>>),
}

impl<M: ManagedTypeApi> RewaOrMultiDcdtPayment<M> {
    pub fn is_empty(&self) -> bool {
        match self {
            RewaOrMultiDcdtPayment::Rewa(rewa_value) => rewa_value == &0u32,
            RewaOrMultiDcdtPayment::MultiDcdt(dcdt_payments) => dcdt_payments.is_empty(),
        }
    }
}

/// The version of `RewaOrMultiDcdtPayment` that contains references instead of owned fields.
pub enum RewaOrMultiDcdtPaymentRefs<'a, M: ManagedTypeApi> {
    Rewa(&'a BigUint<M>),
    MultiDcdt(&'a ManagedVec<M, DcdtTokenPayment<M>>),
}

impl<M: ManagedTypeApi> RewaOrMultiDcdtPayment<M> {
    pub fn as_refs(&self) -> RewaOrMultiDcdtPaymentRefs<'_, M> {
        match self {
            RewaOrMultiDcdtPayment::Rewa(rewa_value) => {
                RewaOrMultiDcdtPaymentRefs::Rewa(rewa_value)
            },
            RewaOrMultiDcdtPayment::MultiDcdt(dcdt_payments) => {
                RewaOrMultiDcdtPaymentRefs::MultiDcdt(dcdt_payments)
            },
        }
    }
}

impl<M: ManagedTypeApi> RewaOrMultiDcdtPaymentRefs<'_, M> {
    pub fn to_owned_payment(&self) -> RewaOrMultiDcdtPayment<M> {
        match self {
            RewaOrMultiDcdtPaymentRefs::Rewa(rewa_value) => {
                RewaOrMultiDcdtPayment::Rewa((*rewa_value).clone())
            },
            RewaOrMultiDcdtPaymentRefs::MultiDcdt(dcdt_payments) => {
                RewaOrMultiDcdtPayment::MultiDcdt((*dcdt_payments).clone())
            },
        }
    }

    pub fn is_empty(&self) -> bool {
        match self {
            RewaOrMultiDcdtPaymentRefs::Rewa(rewa_value) => *rewa_value == &0u32,
            RewaOrMultiDcdtPaymentRefs::MultiDcdt(dcdt_payments) => dcdt_payments.is_empty(),
        }
    }
}
