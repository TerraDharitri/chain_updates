use super::*;
use serde::{
    de::{self, Deserializer, MapAccess, Visitor},
    ser::Serializer,
    Deserialize, Serialize,
};
use std::fmt;

#[derive(Default)]
pub enum CheckDcdtMapRaw {
    #[default]
    Unspecified,
    Star,
    Equal(CheckDcdtMapContentsRaw),
}

impl CheckDcdtMapRaw {
    pub fn is_unspecified(&self) -> bool {
        matches!(self, CheckDcdtMapRaw::Unspecified)
    }

    pub fn is_star(&self) -> bool {
        matches!(self, CheckDcdtMapRaw::Star)
    }
}

impl Serialize for CheckDcdtMapRaw {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        match self {
            CheckDcdtMapRaw::Unspecified => serializer.serialize_str(""),
            CheckDcdtMapRaw::Star => serializer.serialize_str("*"),
            CheckDcdtMapRaw::Equal(m) => m.serialize(serializer),
        }
    }
}

impl<'de> Deserialize<'de> for CheckDcdtMapRaw {
    fn deserialize<D>(deserializer: D) -> Result<Self, D::Error>
    where
        D: Deserializer<'de>,
    {
        deserializer.deserialize_any(CheckDcdtMapRawVisitor)
    }
}
struct CheckDcdtMapRawVisitor;

impl<'de> Visitor<'de> for CheckDcdtMapRawVisitor {
    type Value = CheckDcdtMapRaw;

    // Format a message stating what data this Visitor expects to receive.
    fn expecting(&self, formatter: &mut fmt::Formatter) -> fmt::Result {
        formatter.write_str("serialized object JSON representation of log check")
    }

    fn visit_str<E>(self, value: &str) -> Result<Self::Value, E>
    where
        E: de::Error,
    {
        if value == "*" {
            Ok(CheckDcdtMapRaw::Star)
        } else {
            Err(de::Error::custom("only '*' allowed as logs string value"))
        }
    }

    fn visit_map<M>(self, map: M) -> Result<Self::Value, M::Error>
    where
        M: MapAccess<'de>,
    {
        Ok(CheckDcdtMapRaw::Equal(Deserialize::deserialize(
            de::value::MapAccessDeserializer::new(map),
        )?))
    }
}
