use generic_array::typenum::U16;
use dharitri_sc_codec::IntoMultiValue;

use crate::{
    abi::TypeAbiFrom,
    api::ManagedTypeApi,
    types::{BigUint, RewaOrDcdtTokenIdentifier, RewaOrDcdtTokenPaymentMultiValue},
};

use crate::codec::{
    self,
    derive::{NestedDecode, NestedEncode, TopDecode, TopEncode},
};

use crate as dharitri_sc; // needed by the TypeAbi generated code
use crate::derive::type_abi;

use super::{
    managed_vec_item_read_from_payload_index, managed_vec_item_save_to_payload_index,
    DcdtTokenPayment, DcdtTokenPaymentRefs, ManagedVec, ManagedVecItem,
    ManagedVecItemPayloadBuffer, ManagedVecRef,
};

#[type_abi]
#[derive(TopDecode, TopEncode, NestedDecode, NestedEncode, Clone, PartialEq, Eq, Debug)]
pub struct RewaOrDcdtTokenPayment<M: ManagedTypeApi> {
    pub token_identifier: RewaOrDcdtTokenIdentifier<M>,
    pub token_nonce: u64,
    pub amount: BigUint<M>,
}

/// Alias for a list of payments of REWA or DCDT tokens.
pub type MultiRewaOrDcdtPayment<Api> = ManagedVec<Api, RewaOrDcdtTokenPayment<Api>>;

impl<M: ManagedTypeApi> RewaOrDcdtTokenPayment<M> {
    pub fn no_payment() -> Self {
        Self::rewa_payment(BigUint::zero())
    }

    pub fn new(
        token_identifier: RewaOrDcdtTokenIdentifier<M>,
        token_nonce: u64,
        amount: BigUint<M>,
    ) -> Self {
        RewaOrDcdtTokenPayment {
            token_identifier,
            token_nonce,
            amount,
        }
    }

    /// A payment of token REWA-000000.
    pub fn rewa_payment(amount: BigUint<M>) -> Self {
        Self::new(RewaOrDcdtTokenIdentifier::rewa(), 0, amount)
    }

    /// Will convert to just DCDT or terminate execution if the token is REWA.
    pub fn unwrap_dcdt(self) -> DcdtTokenPayment<M> {
        DcdtTokenPayment::new(
            self.token_identifier.unwrap_dcdt(),
            self.token_nonce,
            self.amount,
        )
    }

    /// Equivalent to a `match { Rewa | Dcdt }`.
    ///
    /// Context passed on from function to closures, to avoid ownership issues.
    /// More precisely, since only one of the two closures `for_rewa` and `for_dcdt` is called,
    /// it is ok for them to have fully owned access to anything from the environment.
    /// The compiler doesn't know that only one of them can ever be called,
    /// so if we pass context to both closures, it will complain that they are moved twice.
    pub fn map_rewa_or_dcdt<Context, D, F, U>(self, context: Context, for_rewa: D, for_dcdt: F) -> U
    where
        D: FnOnce(Context, BigUint<M>) -> U,
        F: FnOnce(Context, DcdtTokenPayment<M>) -> U,
    {
        self.token_identifier.map_or_else(
            (context, self.amount),
            |(context, amount)| for_rewa(context, amount),
            |(context, amount), token_identifier| {
                for_dcdt(
                    context,
                    DcdtTokenPayment::new(token_identifier, self.token_nonce, amount),
                )
            },
        )
    }

    /// Same as `map_rewa_or_dcdt`, but only takes a reference,
    /// and consequently, the closures also only get references.
    pub fn map_ref_rewa_or_dcdt<Context, D, F, U>(
        &self,
        context: Context,
        for_rewa: D,
        for_dcdt: F,
    ) -> U
    where
        D: FnOnce(Context, &BigUint<M>) -> U,
        F: FnOnce(Context, DcdtTokenPaymentRefs<'_, M>) -> U,
    {
        self.token_identifier.map_ref_or_else(
            context,
            |context| for_rewa(context, &self.amount),
            |context, token_identifier| {
                for_dcdt(
                    context,
                    DcdtTokenPaymentRefs::new(token_identifier, self.token_nonce, &self.amount),
                )
            },
        )
    }

    pub fn into_tuple(self) -> (RewaOrDcdtTokenIdentifier<M>, u64, BigUint<M>) {
        (self.token_identifier, self.token_nonce, self.amount)
    }
}

impl<M: ManagedTypeApi> From<(RewaOrDcdtTokenIdentifier<M>, u64, BigUint<M>)>
    for RewaOrDcdtTokenPayment<M>
{
    #[inline]
    fn from(value: (RewaOrDcdtTokenIdentifier<M>, u64, BigUint<M>)) -> Self {
        let (token_identifier, token_nonce, amount) = value;
        Self::new(token_identifier, token_nonce, amount)
    }
}

impl<M: ManagedTypeApi> From<DcdtTokenPayment<M>> for RewaOrDcdtTokenPayment<M> {
    fn from(dcdt_payment: DcdtTokenPayment<M>) -> Self {
        RewaOrDcdtTokenPayment {
            token_identifier: RewaOrDcdtTokenIdentifier::dcdt(dcdt_payment.token_identifier),
            token_nonce: dcdt_payment.token_nonce,
            amount: dcdt_payment.amount,
        }
    }
}

impl<M: ManagedTypeApi> IntoMultiValue for RewaOrDcdtTokenPayment<M> {
    type MultiValue = RewaOrDcdtTokenPaymentMultiValue<M>;

    #[inline]
    fn into_multi_value(self) -> Self::MultiValue {
        self.into()
    }
}

impl<M> TypeAbiFrom<&[u8]> for RewaOrDcdtTokenPayment<M> where M: ManagedTypeApi {}

impl<M: ManagedTypeApi> RewaOrDcdtTokenPayment<M> {
    pub fn as_refs(&self) -> RewaOrDcdtTokenPaymentRefs<'_, M> {
        RewaOrDcdtTokenPaymentRefs::new(&self.token_identifier, self.token_nonce, &self.amount)
    }
}

/// Similar to `RewaOrDcdtTokenPayment`, but only contains references.
pub struct RewaOrDcdtTokenPaymentRefs<'a, M: ManagedTypeApi> {
    pub token_identifier: &'a RewaOrDcdtTokenIdentifier<M>,
    pub token_nonce: u64,
    pub amount: &'a BigUint<M>,
}

impl<'a, M: ManagedTypeApi> RewaOrDcdtTokenPaymentRefs<'a, M> {
    pub fn new(
        token_identifier: &'a RewaOrDcdtTokenIdentifier<M>,
        token_nonce: u64,
        amount: &'a BigUint<M>,
    ) -> RewaOrDcdtTokenPaymentRefs<'a, M> {
        RewaOrDcdtTokenPaymentRefs {
            token_identifier,
            token_nonce,
            amount,
        }
    }

    pub fn to_owned_payment(&self) -> RewaOrDcdtTokenPayment<M> {
        RewaOrDcdtTokenPayment {
            token_identifier: self.token_identifier.clone(),
            token_nonce: self.token_nonce,
            amount: self.amount.clone(),
        }
    }

    pub fn is_empty(&self) -> bool {
        self.amount == &0u32
    }

    pub fn map_rewa_or_dcdt<Context, D, F, U>(self, context: Context, for_rewa: D, for_dcdt: F) -> U
    where
        D: FnOnce(Context, &BigUint<M>) -> U,
        F: FnOnce(Context, DcdtTokenPaymentRefs<M>) -> U,
    {
        self.token_identifier.map_ref_or_else(
            context,
            |context| for_rewa(context, self.amount),
            |context, token_identifier| {
                for_dcdt(
                    context,
                    DcdtTokenPaymentRefs::new(token_identifier, self.token_nonce, self.amount),
                )
            },
        )
    }
}

impl<M: ManagedTypeApi> ManagedVecItem for RewaOrDcdtTokenPayment<M> {
    type PAYLOAD = ManagedVecItemPayloadBuffer<U16>;
    const SKIPS_RESERIALIZATION: bool = false;
    type Ref<'a> = ManagedVecRef<'a, Self>;

    fn read_from_payload(payload: &Self::PAYLOAD) -> Self {
        let mut index = 0;
        unsafe {
            RewaOrDcdtTokenPayment {
                token_identifier: managed_vec_item_read_from_payload_index(payload, &mut index),
                token_nonce: managed_vec_item_read_from_payload_index(payload, &mut index),
                amount: managed_vec_item_read_from_payload_index(payload, &mut index),
            }
        }
    }

    unsafe fn borrow_from_payload<'a>(payload: &Self::PAYLOAD) -> Self::Ref<'a> {
        ManagedVecRef::new(Self::read_from_payload(payload))
    }

    fn save_to_payload(self, payload: &mut Self::PAYLOAD) {
        let mut index = 0;

        unsafe {
            managed_vec_item_save_to_payload_index(self.token_identifier, payload, &mut index);
            managed_vec_item_save_to_payload_index(self.token_nonce, payload, &mut index);
            managed_vec_item_save_to_payload_index(self.amount, payload, &mut index);
        }
    }
}
