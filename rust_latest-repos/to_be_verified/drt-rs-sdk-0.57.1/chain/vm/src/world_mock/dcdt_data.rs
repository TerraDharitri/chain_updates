use num_bigint::BigUint;
use num_traits::Zero;

use crate::{display_util::key_hex, types::VMTokenType};
use std::{
    collections::{hash_map::Iter, HashMap},
    fmt::{self, Write},
};

use super::{DcdtInstanceMetadata, DcdtInstances, DcdtRoles};

#[derive(Clone, Default, Debug)]
pub struct DcdtData {
    pub instances: DcdtInstances,
    pub last_nonce: u64,
    pub roles: DcdtRoles,
    pub frozen: bool,
}

impl DcdtData {
    pub fn is_empty(&self) -> bool {
        self.instances.is_empty_dcdt()
            && self.last_nonce == 0
            && self.roles.is_empty()
            && !self.frozen
    }

    pub fn get_roles(&self) -> Vec<Vec<u8>> {
        self.roles.get()
    }
}

#[derive(Clone, Default, Debug)]
pub struct AccountDcdt(HashMap<Vec<u8>, DcdtData>);

impl AccountDcdt {
    pub fn is_empty(&self) -> bool {
        self.0.is_empty()
    }

    pub fn get_by_identifier(&self, identifier: &[u8]) -> Option<&DcdtData> {
        self.0.get(identifier)
    }

    /// Will provide a clone.
    pub fn get_roles(&self, identifier: &[u8]) -> Vec<Vec<u8>> {
        self.get_by_identifier_or_default(identifier).get_roles()
    }

    pub fn set_roles(&mut self, token_identifier: Vec<u8>, roles: Vec<Vec<u8>>) {
        let dcdt_data = self.0.entry(token_identifier).or_insert_with(|| DcdtData {
            instances: DcdtInstances::new(),
            last_nonce: 0,
            roles: DcdtRoles::default(),
            frozen: false,
        });
        dcdt_data.roles = DcdtRoles::new(roles);
    }

    /// Will provide a clone.
    pub fn get_by_identifier_or_default(&self, identifier: &[u8]) -> DcdtData {
        if let Some(value) = self.0.get(identifier) {
            value.clone()
        } else {
            DcdtData::default()
        }
    }

    pub fn get_mut_by_identifier(&mut self, identifier: &[u8]) -> Option<&mut DcdtData> {
        self.0.get_mut(identifier)
    }

    pub fn new_from_raw_map(hash: HashMap<Vec<u8>, DcdtData>) -> Self {
        AccountDcdt(hash)
    }

    pub fn increase_balance(
        &mut self,
        token_identifier: Vec<u8>,
        nonce: u64,
        value: &BigUint,
        metadata: DcdtInstanceMetadata,
    ) {
        let dcdt_data = self.0.entry(token_identifier).or_insert_with(|| DcdtData {
            instances: DcdtInstances::new(),
            last_nonce: nonce,
            roles: DcdtRoles::default(),
            frozen: false,
        });
        dcdt_data.instances.increase_balance(nonce, value, metadata);
    }

    pub fn set_dcdt_balance(
        &mut self,
        token_identifier: Vec<u8>,
        nonce: u64,
        value: &BigUint,
        metadata: DcdtInstanceMetadata,
    ) {
        let dcdt_data = self.0.entry(token_identifier).or_insert_with(|| DcdtData {
            instances: DcdtInstances::new(),
            last_nonce: nonce,
            roles: DcdtRoles::default(),
            frozen: false,
        });
        dcdt_data.instances.set_balance(nonce, value, metadata);
    }

    pub fn get_dcdt_balance(&self, token_identifier: &[u8], nonce: u64) -> BigUint {
        if let Some(dcdt_data) = self.get_by_identifier(token_identifier) {
            if let Some(instance) = dcdt_data.instances.get_by_nonce(nonce) {
                instance.balance.clone()
            } else {
                BigUint::zero()
            }
        } else {
            BigUint::zero()
        }
    }

    pub fn add_uris(&mut self, token_identifier: &[u8], nonce: u64, mut new_uris: Vec<Vec<u8>>) {
        self.0
            .get_mut(token_identifier)
            .unwrap_or_else(|| panic!("invalid token"))
            .instances
            .get_mut_by_nonce(nonce)
            .unwrap_or_else(|| panic!("invalid token nonce"))
            .metadata
            .uri
            .append(&mut new_uris);
    }

    pub fn update_attributes(
        &mut self,
        token_identifier: &[u8],
        nonce: u64,
        new_attribute_bytes: Vec<u8>,
    ) {
        self.0
            .get_mut(token_identifier)
            .unwrap_or_else(|| panic!("invalid token"))
            .instances
            .get_mut_by_nonce(nonce)
            .unwrap_or_else(|| panic!("invalid token nonce"))
            .metadata
            .attributes = new_attribute_bytes;
    }

    pub fn iter(&self) -> Iter<Vec<u8>, DcdtData> {
        self.0.iter()
    }

    pub fn set_special_role(&mut self, token_identifier: &[u8], role: &[u8]) {
        if let Some(dcdt_data) = self.get_mut_by_identifier(token_identifier) {
            let roles = dcdt_data.roles.get();
            if !roles.contains(role.to_vec().as_ref()) {
                let mut new_roles = roles;
                new_roles.push(role.to_vec());
                dcdt_data.roles = DcdtRoles::new(new_roles);
            }
        }
    }

    pub fn register_and_set_roles(&mut self, token_identifier: &[u8], token_type: VMTokenType) {
        self.issue_token(token_identifier);
        self.set_roles(
            token_identifier.to_vec(),
            Self::get_all_roles_for_token_type(token_type),
        );
    }

    fn issue_token(&mut self, token_identifier: &[u8]) {
        self.0.insert(
            token_identifier.to_vec(),
            DcdtData {
                instances: DcdtInstances::new(),
                last_nonce: 0,
                roles: DcdtRoles::default(),
                frozen: false,
            },
        );
    }

    fn get_all_roles_for_token_type(token_type: VMTokenType) -> Vec<Vec<u8>> {
        match token_type {
            VMTokenType::NonFungible => vec![
                "DCDTRoleNFTCreate".as_bytes().to_vec(),
                "DCDTRoleNFTBurn".as_bytes().to_vec(),
                "DCDTRoleNFTUpdateAttributes".as_bytes().to_vec(),
                "DCDTRoleNFTAddURI".as_bytes().to_vec(),
            ],
            VMTokenType::SemiFungible | VMTokenType::Meta => vec![
                "DCDTRoleNFTCreate".as_bytes().to_vec(),
                "DCDTRoleNFTBurn".as_bytes().to_vec(),
                "DCDTRoleNFTAddQuantity".as_bytes().to_vec(),
            ],
            VMTokenType::Fungible => vec![
                "DCDTRoleLocalMint".as_bytes().to_vec(),
                "DCDTRoleLocalBurn".as_bytes().to_vec(),
                "DCDTRoleLocalTransfer".as_bytes().to_vec(),
            ],
        }
    }
}

impl fmt::Display for DcdtData {
    fn fmt(&self, _f: &mut fmt::Formatter<'_>) -> fmt::Result {
        let mut dcdt_buf = String::new();
        write!(
            dcdt_buf,
            "{{
                instances: [{}],
                last_nonce: {},
                roles: [{}],
                frozen: {},
            }}",
            self.instances, self.last_nonce, self.roles, self.frozen
        )?;
        Ok(())
    }
}

impl fmt::Display for AccountDcdt {
    fn fmt(&self, _f: &mut fmt::Formatter<'_>) -> fmt::Result {
        let mut dcdt_buf = String::new();
        let dcdt_keys: Vec<Vec<u8>> = self.0.keys().cloned().collect();

        for key in &dcdt_keys {
            let value = self.0.get(key).unwrap();
            write!(dcdt_buf, "\n\t\t\t{} -> {value}", key_hex(key.as_slice()))?;
        }
        Ok(())
    }
}
