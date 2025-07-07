use std::collections::{btree_map::Entry, BTreeMap};

use dharitri_chain_scenario_format::interpret_trait::{InterpretableFrom, InterpreterContext};
use dharitri_sc::{
    codec::{top_encode_to_vec_u8, TopEncode},
    types::{AnnotatedValue, BigUint, ManagedAddress, ManagedBuffer, TokenIdentifier},
};

use crate::{
    api::StaticApi,
    scenario::{
        tx_to_step::{
            address_annotated, big_uint_annotated, bytes_annotated, token_identifier_annotated,
            u64_annotated,
        },
        ScenarioRunner,
    },
    scenario_model::{
        AddressKey, BytesKey, BytesValue, CheckAccount, CheckDcdt, CheckDcdtData,
        CheckDcdtInstances, CheckDcdtMap, CheckDcdtMapContents, CheckStateStep, CheckStorage,
        CheckStorageDetails, CheckValue,
    },
    ScenarioTxEnvData, ScenarioWorld,
};

impl ScenarioWorld {
    pub fn check_account<A>(&mut self, address: A) -> CheckStateBuilder<'_>
    where
        A: AnnotatedValue<ScenarioTxEnvData, ManagedAddress<StaticApi>>,
    {
        let address_value = address_annotated(&self.new_env_data(), &address);
        CheckStateBuilder::new(self, address_value.into())
    }
}

pub struct CheckStateBuilder<'w> {
    world: &'w mut ScenarioWorld,
    check_state_step: CheckStateStep,
    current_account: CheckAccount,
    current_address: AddressKey,
}

impl<'w> CheckStateBuilder<'w> {
    pub(crate) fn new(world: &'w mut ScenarioWorld, address: AddressKey) -> CheckStateBuilder<'w> {
        let mut builder = CheckStateBuilder {
            world,
            check_state_step: CheckStateStep::new(),
            current_account: CheckAccount::new(),
            current_address: AddressKey::default(),
        };
        builder.reset_account(address);
        builder
    }

    fn new_env_data(&self) -> ScenarioTxEnvData {
        self.world.new_env_data()
    }

    /// Starts building of a new account.
    pub fn check_account<A>(mut self, address: A) -> Self
    where
        A: AnnotatedValue<ScenarioTxEnvData, ManagedAddress<StaticApi>>,
    {
        self.add_current_account();
        let env = self.new_env_data();
        let address_value = address_annotated(&env, &address);
        self.reset_account(address_value.into());
        self
    }

    fn add_current_account(&mut self) {
        if let Entry::Vacant(entry) = self
            .check_state_step
            .accounts
            .accounts
            .entry(core::mem::take(&mut self.current_address))
        {
            entry.insert(core::mem::take(&mut self.current_account));
        };
    }

    fn reset_account(&mut self, address: AddressKey) {
        self.current_address = address;
        self.current_account = CheckAccount::default();
    }

    /// Finished and sets all account in the blockchain mock.
    fn commit_accounts(&mut self) {
        self.add_current_account();
        self.world.run_check_state_step(&self.check_state_step);
    }

    /// Forces value drop and commit accounts.
    pub fn commit(self) {}

    pub fn nonce<V>(mut self, nonce: V) -> Self
    where
        V: AnnotatedValue<ScenarioTxEnvData, u64>,
    {
        let env = self.new_env_data();
        self.current_account.nonce = CheckValue::Equal(u64_annotated(&env, &nonce));
        self
    }

    pub fn balance<V>(mut self, balance: V) -> Self
    where
        V: AnnotatedValue<ScenarioTxEnvData, BigUint<StaticApi>>,
    {
        let env = self.new_env_data();
        self.current_account.balance = CheckValue::Equal(big_uint_annotated(&env, &balance));
        self
    }

    pub fn code<V>(mut self, code: V) -> Self
    where
        V: AnnotatedValue<ScenarioTxEnvData, ManagedBuffer<StaticApi>>,
    {
        let env = self.new_env_data();
        let code_value = bytes_annotated(&env, code);

        self.current_account.code = CheckValue::Equal(code_value);
        self
    }

    pub fn code_metadata<V>(mut self, code_metadata_expr: V) -> Self
    where
        BytesValue: InterpretableFrom<V>,
    {
        self.current_account.code_metadata = CheckValue::Equal(BytesValue::interpret_from(
            code_metadata_expr,
            &InterpreterContext::default(),
        ));
        self
    }

    pub fn dcdt_balance<K, V>(mut self, token_id: K, balance: V) -> Self
    where
        K: AnnotatedValue<ScenarioTxEnvData, TokenIdentifier<StaticApi>>,
        V: AnnotatedValue<ScenarioTxEnvData, BigUint<StaticApi>>,
    {
        let env = self.new_env_data();
        let token_id_key = token_identifier_annotated(&env, token_id);
        let balance_value = big_uint_annotated(&env, &balance);

        match &mut self.current_account.dcdt {
            CheckDcdtMap::Unspecified | CheckDcdtMap::Star => {
                let mut new_dcdt_map = BTreeMap::new();
                let _ = new_dcdt_map.insert(token_id_key, CheckDcdt::Short(balance_value));

                let new_check_dcdt_map = CheckDcdtMapContents {
                    contents: new_dcdt_map,
                    other_dcdts_allowed: true,
                };

                self.current_account.dcdt = CheckDcdtMap::Equal(new_check_dcdt_map);
            },
            CheckDcdtMap::Equal(check_dcdt_map) => {
                if check_dcdt_map.contents.contains_key(&token_id_key) {
                    let prev_entry = check_dcdt_map.contents.get_mut(&token_id_key).unwrap();
                    match prev_entry {
                        CheckDcdt::Short(prev_balance_check) => *prev_balance_check = balance_value,
                        CheckDcdt::Full(prev_dcdt_check) => match prev_dcdt_check.instances {
                            CheckDcdtInstances::Star => todo!(),
                            CheckDcdtInstances::Equal(_) => todo!(),
                        },
                    }
                }
            },
        }

        self
    }

    pub fn dcdt_nft_balance_and_attributes<K, N, V, T>(
        mut self,
        token_id: K,
        nonce: N,
        balance: V,
        attributes: T,
    ) -> Self
    where
        K: AnnotatedValue<ScenarioTxEnvData, TokenIdentifier<StaticApi>>,
        N: AnnotatedValue<ScenarioTxEnvData, u64>,
        V: AnnotatedValue<ScenarioTxEnvData, BigUint<StaticApi>>,
        T: TopEncode,
    {
        let env = self.new_env_data();
        let token_id_key = token_identifier_annotated(&env, token_id);
        let nonce_value = u64_annotated(&env, &nonce);
        let balance_value = big_uint_annotated(&env, &balance);
        let attributes_value = top_encode_to_vec_u8(&attributes).unwrap();

        if let CheckDcdtMap::Unspecified = &self.current_account.dcdt {
            let mut check_dcdt = CheckDcdt::Full(CheckDcdtData::default());

            check_dcdt.add_balance_and_attributes_check(
                nonce_value,
                balance_value,
                attributes_value,
            );

            let mut new_dcdt_map = BTreeMap::new();
            let _ = new_dcdt_map.insert(token_id_key, check_dcdt);

            let new_check_dcdt_map = CheckDcdtMapContents {
                contents: new_dcdt_map,
                other_dcdts_allowed: true,
            };

            self.current_account.dcdt = CheckDcdtMap::Equal(new_check_dcdt_map);
        }

        self
    }

    pub fn check_storage(mut self, key: &str, value: &str) -> Self {
        let mut details = match &self.current_account.storage {
            CheckStorage::Star => CheckStorageDetails::default(),
            CheckStorage::Equal(details) => details.clone(),
        };
        details.other_storages_allowed = true;
        details.storages.insert(
            BytesKey::interpret_from(key, &InterpreterContext::default()),
            CheckValue::Equal(BytesValue::interpret_from(
                value,
                &InterpreterContext::default(),
            )),
        );
        self.current_account.storage = CheckStorage::Equal(details);
        self
    }
}

impl Drop for CheckStateBuilder<'_> {
    fn drop(&mut self) {
        self.commit_accounts();
    }
}
