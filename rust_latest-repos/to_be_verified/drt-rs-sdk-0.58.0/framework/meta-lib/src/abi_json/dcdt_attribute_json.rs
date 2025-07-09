use dharitri_sc::abi::{DcdtAttributeAbi, TypeName};
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize)]
pub struct DcdtAttributeJson {
    pub ticker: String,
    #[serde(rename = "type")]
    pub ty: TypeName,
}

impl From<&DcdtAttributeAbi> for DcdtAttributeJson {
    fn from(attr: &DcdtAttributeAbi) -> Self {
        DcdtAttributeJson {
            ticker: attr.ticker.to_owned(),
            ty: attr.ty.clone(),
        }
    }
}
