use std::collections::BTreeMap;

use dharitri_sc::abi::DcdtAttributeAbi;
use serde::{Deserialize, Serialize};

use super::{convert_type_descriptions_to_json, DcdtAttributeJson, TypeDescriptionJson};

/// Represents an entire DCDT attribute ABI file. The type descriptions only show up here.
#[derive(Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct DcdtAttributeAbiJson {
    pub dcdt_attribute: DcdtAttributeJson,

    #[serde(default)]
    #[serde(skip_serializing_if = "BTreeMap::is_empty")]
    pub types: BTreeMap<String, TypeDescriptionJson>,
}

impl DcdtAttributeAbiJson {
    pub fn new(attr: &DcdtAttributeAbi) -> Self {
        DcdtAttributeAbiJson {
            dcdt_attribute: DcdtAttributeJson::from(attr),
            types: convert_type_descriptions_to_json(&attr.type_descriptions),
        }
    }
}
