#![no_std]

dharitri_sc::imports!();

pub mod builtin_func_features_proxy;
pub mod dcdt_features;

/// Test contract for investigating async calls.
#[dharitri_sc::contract]
pub trait BuiltinFuncFeatures: dcdt_features::DcdtFeaturesModule {
    #[init]
    fn init(&self, fungible_token_id: TokenIdentifier, non_fungible_token_id: TokenIdentifier) {
        self.fungible_dcdt_token_id()
            .set_token_id(fungible_token_id);
        self.non_fungible_dcdt_token_id()
            .set_token_id(non_fungible_token_id);
    }

    #[endpoint]
    fn call_set_user_name(&self, address: ManagedAddress, name: ManagedBuffer) {
        self.tx()
            .to(&address)
            .typed(system_proxy::UserBuiltinProxy)
            .set_user_name(name)
            .async_call_and_exit()
    }

    #[endpoint]
    fn call_delete_user_name(&self, address: ManagedAddress) {
        self.tx()
            .to(&address)
            .typed(system_proxy::UserBuiltinProxy)
            .delete_user_name()
            .async_call_and_exit()
    }
}
