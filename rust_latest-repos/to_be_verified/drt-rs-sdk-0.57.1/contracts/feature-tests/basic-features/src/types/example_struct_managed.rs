use dharitri_sc::{
    api::ManagedTypeApi,
    types::{BigUint, ManagedBuffer},
};

dharitri_sc::derive_imports!();

#[type_abi]
#[derive(NestedEncode, NestedDecode, TopEncode, TopDecode, PartialEq, Eq, Debug, Clone)]
pub struct ExampleStructManaged<M: ManagedTypeApi> {
    pub big_uint: BigUint<M>,
    pub int: u32,
    pub bytes: ManagedBuffer<M>,
}
