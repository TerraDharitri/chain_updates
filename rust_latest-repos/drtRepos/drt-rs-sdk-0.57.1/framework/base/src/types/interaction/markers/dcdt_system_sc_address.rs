use hex_literal::hex;
use dharitri_chain_core::types::Address;
use dharitri_sc_codec::{EncodeErrorHandler, TopEncode, TopEncodeOutput};

use crate::{
    abi::TypeAbiFrom,
    api::ManagedTypeApi,
    types::{AnnotatedValue, ManagedAddress, ManagedBuffer, TxEnv, TxTo, TxToSpecified},
};

/// Address of the system smart contract that manages DCDT.
const DCDT_SYSTEM_SC_ADDRESS_BYTES: [u8; 32] =
    hex!("233300000000000000000000000000000002333000000000000000000002ffff");
const DCDT_SYSTEM_SC_ADDRESS_BECH32: &str =
    "drt1yvesqqqqqqqqqqqqqqqqqqqqqqqqyvesqqqqqqqqqqqqqqqzlllsd5j0s2";
const DCDT_SYSTEM_SC_ADDRESS_ANNOTATION: &str =
    "bech32:drt1yvesqqqqqqqqqqqqqqqqqqqqqqqqyvesqqqqqqqqqqqqqqqzlllsd5j0s2";

/// Indicates the system SC address, which is the same on any Dharitri blockchain.
pub struct DCDTSystemSCAddress;

impl DCDTSystemSCAddress {
    pub fn to_managed_address<Api>(self) -> ManagedAddress<Api>
    where
        Api: ManagedTypeApi,
    {
        ManagedAddress::from(DCDT_SYSTEM_SC_ADDRESS_BYTES)
    }

    pub fn to_address(&self) -> Address {
        DCDT_SYSTEM_SC_ADDRESS_BYTES.into()
    }

    pub fn to_bech32_str(&self) -> &str {
        DCDT_SYSTEM_SC_ADDRESS_BECH32
    }

    pub fn to_bech32_string(&self) -> alloc::string::String {
        DCDT_SYSTEM_SC_ADDRESS_BECH32.into()
    }
}

impl<Env> AnnotatedValue<Env, ManagedAddress<Env::Api>> for DCDTSystemSCAddress
where
    Env: TxEnv,
{
    fn annotation(&self, _env: &Env) -> ManagedBuffer<Env::Api> {
        ManagedBuffer::from(DCDT_SYSTEM_SC_ADDRESS_ANNOTATION)
    }

    fn to_value(&self, _env: &Env) -> ManagedAddress<Env::Api> {
        DCDTSystemSCAddress.to_managed_address()
    }
}

impl<Env> TxTo<Env> for DCDTSystemSCAddress where Env: TxEnv {}
impl<Env> TxToSpecified<Env> for DCDTSystemSCAddress where Env: TxEnv {}

impl TopEncode for DCDTSystemSCAddress {
    fn top_encode_or_handle_err<O, H>(&self, output: O, h: H) -> Result<(), H::HandledErr>
    where
        O: TopEncodeOutput,
        H: EncodeErrorHandler,
    {
        DCDT_SYSTEM_SC_ADDRESS_BYTES.top_encode_or_handle_err(output, h)
    }
}

impl<M> TypeAbiFrom<DCDTSystemSCAddress> for ManagedAddress<M> where M: ManagedTypeApi {}

impl core::fmt::Display for DCDTSystemSCAddress {
    fn fmt(&self, f: &mut core::fmt::Formatter<'_>) -> core::fmt::Result {
        f.write_str(DCDT_SYSTEM_SC_ADDRESS_BECH32)
    }
}
