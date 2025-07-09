#![no_std]

use dharitri_sc::imports::*;

use dharitri_sc_modules::{default_issue_callbacks, subscription};

#[dharitri_sc::contract]
pub trait NftSubscription:
    default_issue_callbacks::DefaultIssueCallbacksModule + subscription::SubscriptionModule
{
    #[init]
    fn init(&self) {}

    #[endpoint]
    #[payable("REWA")]
    fn issue(&self) {
        self.token_id().issue_and_set_all_roles(
            DcdtTokenType::NonFungible,
            self.call_value().rewa().clone(),
            ManagedBuffer::from(b"Subscription"),
            ManagedBuffer::from(b"SUB"),
            0,
            None,
        )
    }

    #[endpoint]
    fn mint(&self) {
        let nonce = self.create_subscription_nft(
            self.token_id().get_token_id_ref(),
            &BigUint::from(1u8),
            &ManagedBuffer::new(),
            &BigUint::from(0u8),
            &ManagedBuffer::new(),
            0,
            ManagedBuffer::from(b"common"),
            &ManagedVec::new(),
        );

        self.tx()
            .to(ToCaller)
            .single_dcdt(
                self.token_id().get_token_id_ref(),
                nonce,
                &BigUint::from(1u8),
            )
            .transfer();
    }

    #[payable]
    #[endpoint]
    fn update_attributes(&self, attributes: ManagedBuffer) {
        let payment = self.call_value().single_dcdt();
        self.update_subscription_attributes::<ManagedBuffer>(
            &payment.token_identifier,
            payment.token_nonce,
            attributes,
        );
        self.tx()
            .to(ToCaller)
            .single_dcdt(
                &payment.token_identifier,
                payment.token_nonce,
                &BigUint::from(1u8),
            )
            .transfer();
    }

    #[payable]
    #[endpoint]
    fn renew(&self, duration: u64) {
        let payment = self.call_value().single_dcdt();
        self.renew_subscription::<ManagedBuffer>(
            &payment.token_identifier,
            payment.token_nonce,
            duration,
        );
        self.tx()
            .to(ToCaller)
            .single_dcdt(
                &payment.token_identifier,
                payment.token_nonce,
                &BigUint::from(1u8),
            )
            .transfer();
    }

    #[payable]
    #[endpoint]
    fn cancel(&self) {
        let payment = self.call_value().single_dcdt();
        self.cancel_subscription::<ManagedBuffer>(&payment.token_identifier, payment.token_nonce);

        self.tx()
            .to(ToCaller)
            .single_dcdt(
                &payment.token_identifier,
                payment.token_nonce,
                &BigUint::from(1u8),
            )
            .transfer();
    }

    #[storage_mapper("tokenId")]
    fn token_id(&self) -> NonFungibleTokenMapper;
}
