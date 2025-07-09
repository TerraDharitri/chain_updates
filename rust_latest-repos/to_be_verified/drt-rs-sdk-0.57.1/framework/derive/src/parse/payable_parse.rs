use super::{attributes::PayableAttribute, MethodAttributesPass1};
use crate::model::MethodPayableMetadata;

pub fn process_payable_attribute(
    attr: &syn::Attribute,
    pass_1_data: &mut MethodAttributesPass1,
) -> bool {
    PayableAttribute::parse(attr)
        .map(|payable_attr| {
            pass_1_data.payable = parse_payable_identifier(&payable_attr.identifier);
        })
        .is_some()
}

fn parse_payable_identifier(identifier: &str) -> MethodPayableMetadata {
    match identifier {
        "REWA" => MethodPayableMetadata::Rewa,
        "*" => MethodPayableMetadata::AnyToken,
        "" => panic!("empty token name not allowed in #[payable] attribute"),
        _ => MethodPayableMetadata::SingleDcdtToken(identifier.to_string()),
    }
}
