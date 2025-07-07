use super::*;
use serde::{
    de::{self, Deserializer, SeqAccess, Visitor},
    ser::{SerializeMap, SerializeSeq, Serializer},
    Deserialize, Serialize,
};
use std::fmt;

#[derive(Default)]
pub enum CheckDcdtInstancesRaw {
    #[default]
    Unspecified,
    Star,
    Equal(Vec<CheckDcdtInstanceRaw>),
}

impl CheckDcdtInstancesRaw {
    pub fn is_star(&self) -> bool {
        matches!(self, CheckDcdtInstancesRaw::Star)
    }

    pub fn is_unspecified(&self) -> bool {
        matches!(self, CheckDcdtInstancesRaw::Unspecified)
    }
}

impl Serialize for CheckDcdtInstancesRaw {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        match self {
            CheckDcdtInstancesRaw::Unspecified => {
                let map = serializer.serialize_map(Some(0))?;
                map.end()
            },
            CheckDcdtInstancesRaw::Star => serializer.serialize_str("*"),
            CheckDcdtInstancesRaw::Equal(m) => {
                let mut map = serializer.serialize_seq(Some(m.len()))?;
                for v in m {
                    map.serialize_element(v)?;
                }
                map.end()
            },
        }
    }
}

struct CheckDcdtInstancesRawVisitor;

impl<'de> Visitor<'de> for CheckDcdtInstancesRawVisitor {
    type Value = CheckDcdtInstancesRaw;

    // Format a message stating what data this Visitor expects to receive.
    fn expecting(&self, formatter: &mut fmt::Formatter) -> fmt::Result {
        formatter.write_str("serialized object JSON representation of an DCDT instances list")
    }

    fn visit_str<E>(self, value: &str) -> Result<Self::Value, E>
    where
        E: de::Error,
    {
        if value == "*" {
            Ok(CheckDcdtInstancesRaw::Star)
        } else {
            Err(de::Error::custom(
                "only '*' allowed as DCDT instances value",
            ))
        }
    }

    fn visit_seq<A>(self, mut seq: A) -> Result<Self::Value, A::Error>
    where
        A: SeqAccess<'de>,
    {
        let mut list = Vec::<CheckDcdtInstanceRaw>::new();

        while let Some(item) = seq.next_element()? {
            list.push(item);
        }

        Ok(CheckDcdtInstancesRaw::Equal(list))
    }
}

impl<'de> Deserialize<'de> for CheckDcdtInstancesRaw {
    fn deserialize<D>(deserializer: D) -> Result<Self, D::Error>
    where
        D: Deserializer<'de>,
    {
        deserializer.deserialize_any(CheckDcdtInstancesRawVisitor)
    }
}
