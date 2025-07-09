use super::*;
use crate::serde_raw::ValueSubTree;
use serde::{
    de::{self, Deserializer, MapAccess, Visitor},
    ser::Serializer,
    Deserialize, Serialize,
};
use std::fmt;

pub enum CheckDcdtRaw {
    Short(ValueSubTree),
    Full(CheckDcdtDataRaw),
}

impl Serialize for CheckDcdtRaw {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        match self {
            CheckDcdtRaw::Short(m) => m.serialize(serializer),
            CheckDcdtRaw::Full(m) => m.serialize(serializer),
        }
    }
}

impl<'de> Deserialize<'de> for CheckDcdtRaw {
    fn deserialize<D>(deserializer: D) -> Result<Self, D::Error>
    where
        D: Deserializer<'de>,
    {
        deserializer.deserialize_any(CheckDcdtRawVisitor)
    }
}
struct CheckDcdtRawVisitor;

impl<'de> Visitor<'de> for CheckDcdtRawVisitor {
    type Value = CheckDcdtRaw;

    // Format a message stating what data this Visitor expects to receive.
    fn expecting(&self, formatter: &mut fmt::Formatter) -> fmt::Result {
        formatter.write_str("serialized object JSON representation of dcdt check")
    }

    fn visit_str<E>(self, value: &str) -> Result<Self::Value, E>
    where
        E: de::Error,
    {
        Ok(CheckDcdtRaw::Short(ValueSubTree::Str(value.to_string())))
    }

    fn visit_map<M>(self, map: M) -> Result<Self::Value, M::Error>
    where
        M: MapAccess<'de>,
    {
        Ok(CheckDcdtRaw::Full(Deserialize::deserialize(
            de::value::MapAccessDeserializer::new(map),
        )?))
    }
}
