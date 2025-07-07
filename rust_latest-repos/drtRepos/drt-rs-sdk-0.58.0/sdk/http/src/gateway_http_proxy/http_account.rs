use anyhow::Result;
use dharitri_sdk::{
    chain_core::types::Address,
    data::{account::Account, dcdt::DcdtBalance},
    gateway::{
        GetAccountDcdtRolesRequest, GetAccountDcdtTokensRequest, GetAccountRequest,
        GetAccountStorageRequest,
    },
};
use std::collections::HashMap;

use super::GatewayHttpProxy;

impl GatewayHttpProxy {
    // get_account retrieves an account info from the network (nonce, balance)
    pub async fn get_account(&self, address: &Address) -> Result<Account> {
        self.http_request(GetAccountRequest::new(address)).await
    }

    // get_account_dcdt_roles retrieves an all dcdt roles of an account from the network
    pub async fn get_account_dcdt_roles(
        &self,
        address: &Address,
    ) -> Result<HashMap<String, Vec<String>>> {
        self.http_request(GetAccountDcdtRolesRequest::new(address))
            .await
    }

    // get_account_dcdt_tokens retrieves an all dcdt token of an account from the network
    pub async fn get_account_dcdt_tokens(
        &self,
        address: &Address,
    ) -> Result<HashMap<String, DcdtBalance>> {
        self.http_request(GetAccountDcdtTokensRequest::new(address))
            .await
    }

    // get_account_dcdt_tokens retrieves an all dcdt token of an account from the network
    pub async fn get_account_storage_keys(
        &self,
        address: &Address,
    ) -> Result<HashMap<String, String>> {
        self.http_request(GetAccountStorageRequest::new(address))
            .await
    }
}
