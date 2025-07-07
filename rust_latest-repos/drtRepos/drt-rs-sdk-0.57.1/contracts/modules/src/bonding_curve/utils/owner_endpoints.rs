dharitri_sc::imports!();
dharitri_sc::derive_imports!();

use dharitri_sc::contract_base::ManagedSerializer;

use crate::bonding_curve::{
    curves::curve_function::CurveFunction,
    utils::{
        events, storage,
        structs::{BondingCurve, TokenOwnershipData},
    },
};

use super::structs::CurveArguments;

#[dharitri_sc::module]
pub trait OwnerEndpointsModule: storage::StorageModule + events::EventsModule {
    #[endpoint(setLocalRoles)]
    fn set_local_roles(
        &self,
        address: ManagedAddress,
        token_identifier: TokenIdentifier,
        roles: MultiValueEncoded<DcdtLocalRole>,
    ) {
        self.send()
            .dcdt_system_sc_proxy()
            .set_special_roles(&address, &token_identifier, roles.into_iter())
            .async_call_and_exit()
    }

    #[endpoint(unsetLocalRoles)]
    fn unset_local_roles(
        &self,
        address: ManagedAddress,
        token_identifier: TokenIdentifier,
        roles: MultiValueEncoded<DcdtLocalRole>,
    ) {
        self.send()
            .dcdt_system_sc_proxy()
            .unset_special_roles(&address, &token_identifier, roles.into_iter())
            .async_call_and_exit()
    }

    fn set_bonding_curve<T>(
        &self,
        identifier: TokenIdentifier,
        function: T,
        sell_availability: bool,
    ) where
        T: CurveFunction<Self::Api>
            + TopEncode
            + TopDecode
            + NestedEncode
            + NestedDecode
            + TypeAbi
            + PartialEq
            + Default,
    {
        require!(
            !self.token_details(&identifier).is_empty(),
            "Token is not issued yet!"
        );

        let caller = self.blockchain().get_caller();

        let details = self.token_details(&identifier).get();
        require!(
            details.owner == caller,
            "The price function can only be set by the seller."
        );
        self.bonding_curve(&identifier).update(|buffer| {
            let serializer = ManagedSerializer::new();

            let mut bonding_curve: BondingCurve<Self::Api, T> =
                serializer.top_decode_from_managed_buffer(buffer);
            bonding_curve.curve = function;
            bonding_curve.sell_availability = sell_availability;
            *buffer = serializer.top_encode_to_managed_buffer(&bonding_curve);
        });
    }

    fn deposit<T>(&self, payment_token: OptionalValue<TokenIdentifier>)
    where
        T: CurveFunction<Self::Api>
            + TopEncode
            + TopDecode
            + NestedEncode
            + NestedDecode
            + TypeAbi
            + PartialEq
            + Default,
    {
        let dcdt_payment = self.call_value().single_dcdt();
        let identifier = &dcdt_payment.token_identifier;
        let nonce = dcdt_payment.token_nonce;
        let amount = &dcdt_payment.amount;
        let caller = self.blockchain().get_caller();
        let mut set_payment = RewaOrDcdtTokenIdentifier::rewa();

        if self.bonding_curve(identifier).is_empty() {
            match payment_token {
                OptionalValue::Some(token) => set_payment = RewaOrDcdtTokenIdentifier::dcdt(token),
                OptionalValue::None => {
                    sc_panic!("Expected provided accepted_payment for the token");
                },
            };
        }
        if self.token_details(identifier).is_empty() {
            let nonces = ManagedVec::from_single_item(nonce);
            self.token_details(identifier).set(&TokenOwnershipData {
                token_nonces: nonces,
                owner: caller.clone(),
            });
        } else {
            let mut details = self.token_details(identifier).get();
            require!(
                details.owner == caller,
                "The token was already deposited by another address"
            );
            if !details.token_nonces.contains(&nonce) {
                details.token_nonces.push(nonce);
                self.token_details(identifier).set(&details);
            }
        }

        self.set_curve_storage::<T>(identifier, amount.clone(), set_payment);
        self.owned_tokens(&caller).insert(identifier.clone());
        self.nonce_amount(identifier, nonce)
            .update(|current_amount| *current_amount += amount);
    }

    fn claim<T>(&self)
    where
        T: CurveFunction<Self::Api>
            + TopEncode
            + TopDecode
            + NestedEncode
            + NestedDecode
            + TypeAbi
            + PartialEq
            + Default,
    {
        let caller = self.blockchain().get_caller();
        require!(
            !self.owned_tokens(&caller).is_empty(),
            "You have nothing to claim"
        );

        let mut tokens_to_claim = ManagedVec::<Self::Api, DcdtTokenPayment<Self::Api>>::new();
        let mut rewa_to_claim = BigUint::zero();
        let serializer = ManagedSerializer::new();
        for token in self.owned_tokens(&caller).iter() {
            let nonces = self.token_details(&token).get().token_nonces;
            for nonce in &nonces {
                tokens_to_claim.push(DcdtTokenPayment::new(
                    token.clone(),
                    nonce,
                    self.nonce_amount(&token, nonce).get(),
                ));

                self.nonce_amount(&token, nonce).clear();
            }

            let bonding_curve: BondingCurve<Self::Api, T> =
                serializer.top_decode_from_managed_buffer(&self.bonding_curve(&token).get());

            if let Some(dcdt_token_identifier) =
                bonding_curve.payment.token_identifier.into_dcdt_option()
            {
                tokens_to_claim.push(DcdtTokenPayment::new(
                    dcdt_token_identifier,
                    bonding_curve.payment.token_nonce,
                    bonding_curve.payment.amount,
                ));
            } else {
                rewa_to_claim += bonding_curve.payment.amount;
            }

            self.token_details(&token).clear();
            self.bonding_curve(&token).clear();
        }
        self.owned_tokens(&caller).clear();
        self.tx().to(&caller).multi_dcdt(tokens_to_claim).transfer();
        if rewa_to_claim > BigUint::zero() {
            self.tx().to(&caller).rewa(&rewa_to_claim).transfer();
        }
    }

    fn set_curve_storage<T>(
        &self,
        identifier: &TokenIdentifier,
        amount: BigUint,
        payment_token_identifier: RewaOrDcdtTokenIdentifier,
    ) where
        T: CurveFunction<Self::Api>
            + TopEncode
            + TopDecode
            + NestedEncode
            + NestedDecode
            + TypeAbi
            + PartialEq
            + Default,
    {
        let mut curve: T = T::default();
        let mut arguments;
        let payment;
        let sell_availability: bool;
        let serializer = ManagedSerializer::new();

        if self.bonding_curve(identifier).is_empty() {
            arguments = CurveArguments {
                available_supply: amount.clone(),
                balance: amount,
            };
            payment = RewaOrDcdtTokenPayment::new(payment_token_identifier, 0, BigUint::zero());
            sell_availability = false;
        } else {
            let bonding_curve: BondingCurve<Self::Api, T> =
                serializer.top_decode_from_managed_buffer(&self.bonding_curve(identifier).get());

            payment = bonding_curve.payment;
            curve = bonding_curve.curve;
            arguments = bonding_curve.arguments;
            arguments.balance += &amount;
            arguments.available_supply += amount;
            sell_availability = bonding_curve.sell_availability;
        }
        let encoded_curve = serializer.top_encode_to_managed_buffer(&BondingCurve {
            curve,
            arguments,
            sell_availability,
            payment,
        });
        self.bonding_curve(identifier).set(encoded_curve);
    }
}
