use serde::{Deserialize, Serialize};
use std::{fs::File, io::Write, path::Path};

use crate::{
    abi_json::{BuildInfoAbiJson, ContractAbiJson},
    report_info_json::ReportInfoJson,
};

#[derive(Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct DrtscFileJson {
    pub build_info: BuildInfoAbiJson,
    pub abi: ContractAbiJson,
    pub code: String,
    pub report: ReportInfoJson,
}

pub fn serialize_drtsc_file_json(drtsc_file_json: &DrtscFileJson) -> String {
    let buf = Vec::new();
    let formatter = serde_json::ser::PrettyFormatter::with_indent(b"    ");
    let mut ser = serde_json::Serializer::with_formatter(buf, formatter);
    drtsc_file_json.serialize(&mut ser).unwrap();
    let mut serialized = String::from_utf8(ser.into_inner()).unwrap();
    serialized.push('\n');
    serialized
}

pub fn save_drtsc_file_json(drtsc_file_json: &DrtscFileJson, path: impl AsRef<Path>) {
    let drtsc_file_string = serialize_drtsc_file_json(drtsc_file_json);
    let mut drtsc_file = File::create(path).unwrap();
    write!(drtsc_file, "{drtsc_file_string}").unwrap();
}
