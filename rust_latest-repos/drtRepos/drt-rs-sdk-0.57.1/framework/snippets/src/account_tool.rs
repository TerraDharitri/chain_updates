use crate::sdk_core::data::dcdt::DcdtBalance;
use dharitri_chain_scenario_format::interpret_trait::IntoRaw;
use dharitri_sc_scenario::{
    imports::Bech32Address,
    scenario_model::{Account, BytesKey, BytesValue, Scenario, SetStateStep, Step},
};
use dharitri_sdk::gateway::{GatewayAsyncService, SetStateAccount};
use dharitri_sdk::gateway::{
    GetAccountDcdtRolesRequest, GetAccountDcdtTokensRequest, GetAccountRequest,
    GetAccountStorageRequest,
};
use std::collections::{BTreeMap, HashMap};

/// Called directly from CLI, from `sc-meta`.
///
/// Retrieves an account data via the API,
/// then formats it as a scenario set state step.
pub async fn print_account_as_scenario_set_state<GatewayProxy: GatewayAsyncService>(
    gateway_proxy: GatewayProxy,
    address_bech32_string: String,
) {
    let address = Bech32Address::from_bech32_string(address_bech32_string);
    let (_, set_state) = retrieve_account_as_scenario_set_state(&gateway_proxy, &address).await;
    let scenario = build_scenario(set_state);
    println!("{}", scenario.into_raw().to_json_string());
}

fn build_scenario(set_state: SetStateStep) -> Scenario {
    Scenario {
        name: None,
        comment: None,
        check_gas: None,
        steps: vec![Step::SetState(set_state)],
    }
}

pub async fn retrieve_account_as_scenario_set_state<GatewayProxy: GatewayAsyncService>(
    api: &GatewayProxy,
    bech32_address: &Bech32Address,
) -> (SetStateAccount, SetStateStep) {
    let address = bech32_address.as_address();
    let sdk_account = api.request(GetAccountRequest::new(address)).await.unwrap();

    let account_dcdt = api
        .request(GetAccountDcdtTokensRequest::new(address))
        .await
        .unwrap_or_else(|err| {
            eprintln!("failed to retrieve DCDT tokens for address {bech32_address}: {err}");
            HashMap::new()
        });
    let account_dcdt_roles = api
        .request(GetAccountDcdtRolesRequest::new(address))
        .await
        .unwrap_or_else(|err| {
            eprintln!("failed to retrieve DCDT roles for address {bech32_address}: {err}");
            HashMap::new()
        });
    let account_storage = api
        .request(GetAccountStorageRequest::new(address))
        .await
        .unwrap_or_else(|err| {
            panic!("failed to retrieve storage for address {bech32_address}: {err}")
        });

    let account_state = set_account(
        sdk_account.clone(),
        account_storage.clone(),
        account_dcdt,
        account_dcdt_roles,
    );

    let set_state_account = SetStateAccount::from(sdk_account).with_storage(account_storage);
    let set_state_step = SetStateStep::new();
    (
        set_state_account,
        set_state_step.put_account(bech32_address, account_state),
    )
}

fn set_account(
    account: crate::sdk::data::account::Account,
    account_storage: HashMap<String, String>,
    account_dcdt: HashMap<String, DcdtBalance>,
    account_dcdt_roles: HashMap<String, Vec<String>>,
) -> Account {
    let mut account_state = Account::new()
        .nonce(account.nonce)
        .balance(account.balance.as_str())
        .code(BytesValue::from_hex(&account.code));
    account_state.username = Some(format!("str:{}", account.username.as_str()).into());
    account_state.storage = convert_storage(account_storage);

    for (_, dcdt_balance) in account_dcdt.iter() {
        let token_id_expr = format!("str:{}", dcdt_balance.token_identifier);
        account_state =
            account_state.dcdt_balance(token_id_expr.as_str(), dcdt_balance.balance.as_str());
    }

    for (token_id, dcdt_roles) in account_dcdt_roles {
        let token_id_expr = format!("str:{token_id}");
        account_state = account_state.dcdt_roles(token_id_expr.as_str(), dcdt_roles);
    }

    account_state
}

fn convert_storage(account_storage: HashMap<String, String>) -> BTreeMap<BytesKey, BytesValue> {
    account_storage
        .into_iter()
        .filter(|(k, _)| !k.starts_with("4e554d424154"))
        .map(|(k, v)| (BytesKey::from_hex(&k), BytesValue::from_hex(&v)))
        .collect()
}
