use serde::{Deserialize, Serialize};
use std::collections::HashMap;

// DcdtBalance  holds information about the dcdt balance
#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct DcdtBalance {
    pub token_identifier: String,
    pub balance: String,
    #[serde(default)]
    #[serde(skip_serializing_if = "Vec::is_empty")]
    pub uris: Vec<String>,
}

// DcdtBalanceData holds the dcdt balance data
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DcdtBalanceData {
    pub dcdts: HashMap<String, DcdtBalance>,
}

// DcdtBalanceResponse holds the dcdt balance endpoint response
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DcdtBalanceResponse {
    pub data: Option<DcdtBalanceData>,
    pub error: String,
    pub code: String,
}

// DcdtRolesData holds the dcdt roles data
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DcdtRolesData {
    pub roles: HashMap<String, Vec<String>>,
}

// DcdtBalanceResponse holds the dcdt roles endpoint response
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DcdtRolesResponse {
    pub data: Option<DcdtRolesData>,
    pub error: String,
    pub code: String,
}
