use crate::data::dcdt::{DcdtBalance, DcdtBalanceResponse};
use anyhow::anyhow;
use dharitri_chain_core::types::Address;
use std::collections::HashMap;

use super::{GatewayRequest, GatewayRequestType, ACCOUNT_ENDPOINT};

/// Retrieves an all dcdt tokens of an account from the network.
pub struct GetAccountDcdtTokensRequest<'a> {
    pub address: &'a Address,
}

impl<'a> GetAccountDcdtTokensRequest<'a> {
    pub fn new(address: &'a Address) -> Self {
        Self { address }
    }
}

impl GatewayRequest for GetAccountDcdtTokensRequest<'_> {
    type Payload = ();
    type DecodedJson = DcdtBalanceResponse;
    type Result = HashMap<String, DcdtBalance>;

    fn request_type(&self) -> GatewayRequestType {
        GatewayRequestType::Get
    }

    fn get_endpoint(&self) -> String {
        format!(
            "{ACCOUNT_ENDPOINT}/{}/dcdt",
            crate::bech32::encode(self.address)
        )
    }

    fn process_json(&self, decoded: Self::DecodedJson) -> anyhow::Result<Self::Result> {
        match decoded.data {
            None => Err(anyhow!("{}", decoded.error)),
            Some(b) => Ok(b.dcdts),
        }
    }
}
