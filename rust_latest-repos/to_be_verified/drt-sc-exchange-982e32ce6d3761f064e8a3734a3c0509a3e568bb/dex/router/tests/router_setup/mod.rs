use dharitri_sc::codec::multi_types::{MultiValue4, OptionalValue};
use dharitri_sc::types::{Address, DcdtLocalRole, ManagedAddress, MultiValueEncoded};
use dharitri_sc_scenario::whitebox_legacy::TxTokenTransfer;
use dharitri_sc_scenario::{
    managed_address, managed_biguint, managed_buffer, managed_token_id, rust_biguint,
    whitebox_legacy::*, DebugApi,
};

pub const PAIR_WASM_PATH: &str = "pair/output/pair.wasm";
pub const ROUTER_WASM_PATH: &str = "router/output/router.wasm";
pub const MOA_TOKEN_ID: &[u8] = b"MOA-abcdef";
pub const WREWA_TOKEN_ID: &[u8] = b"WREWA-abcdef";
pub const USDC_TOKEN_ID: &[u8] = b"USDC-abcdef";
pub const LPMOA_TOKEN_ID: &[u8] = b"LPMOA-abcdef";
pub const LPUSDC_TOKEN_ID: &[u8] = b"LPUSDC-abcdef";

pub const USER_TOTAL_MOA_TOKENS: u64 = 5_001_001_000;
pub const USER_TOTAL_WREWA_TOKENS: u64 = 5_002_002_000;
pub const USER_TOTAL_USDC_TOKENS: u64 = 5_001_001_000;

pub const ADD_LIQUIDITY_TOKENS: u64 = 1_001_000;

pub static CUSTOM_TOKEN_ID: &[u8] = b"CUSTOM-abcdef";
pub static LOCKED_TOKEN_ID: &[u8] = b"LOCKED-abcdef";
pub const MIN_LOCKED_TOKEN_VALUE: u64 = 500_000;
pub const MIN_LOCKED_PERIOD_EPOCHS: u64 = 100;
pub const USER_CUSTOM_TOKEN_BALANCE: u64 = 1_000_000_000;
pub const USER_USDC_BALANCE: u64 = 1_000_000;

use pair::config::ConfigModule as PairConfigModule;
use pair::pair_actions::add_liq::AddLiquidityModule;
use pair::*;
use pausable::{PausableModule, State};
use router::config::ConfigModule;
use router::factory::*;
use router::multi_pair_swap::*;
use router::*;

#[allow(dead_code)]
pub struct RouterSetup<RouterObjBuilder, PairObjBuilder>
where
    RouterObjBuilder: 'static + Copy + Fn() -> router::ContractObj<DebugApi>,
    PairObjBuilder: 'static + Copy + Fn() -> pair::ContractObj<DebugApi>,
{
    pub blockchain_wrapper: BlockchainStateWrapper,
    pub owner_address: Address,
    pub user_address: Address,
    pub router_wrapper: ContractObjWrapper<router::ContractObj<DebugApi>, RouterObjBuilder>,
    pub moa_pair_wrapper: ContractObjWrapper<pair::ContractObj<DebugApi>, PairObjBuilder>,
    pub usdc_pair_wrapper: ContractObjWrapper<pair::ContractObj<DebugApi>, PairObjBuilder>,
}

impl<RouterObjBuilder, PairObjBuilder> RouterSetup<RouterObjBuilder, PairObjBuilder>
where
    RouterObjBuilder: 'static + Copy + Fn() -> router::ContractObj<DebugApi>,
    PairObjBuilder: 'static + Copy + Fn() -> pair::ContractObj<DebugApi>,
{
    pub fn new(router_builder: RouterObjBuilder, pair_builder: PairObjBuilder) -> Self {
        let rust_zero = rust_biguint!(0u64);
        let mut blockchain_wrapper = BlockchainStateWrapper::new();
        let owner_addr = blockchain_wrapper.create_user_account(&rust_zero);

        let router_wrapper = blockchain_wrapper.create_sc_account(
            &rust_zero,
            Some(&owner_addr),
            router_builder,
            ROUTER_WASM_PATH,
        );

        let moa_pair_wrapper = blockchain_wrapper.create_sc_account(
            &rust_zero,
            Some(&owner_addr),
            pair_builder,
            PAIR_WASM_PATH,
        );

        let usdc_pair_wrapper = blockchain_wrapper.create_sc_account(
            &rust_zero,
            Some(&owner_addr),
            pair_builder,
            PAIR_WASM_PATH,
        );

        blockchain_wrapper
            .execute_tx(&owner_addr, &moa_pair_wrapper, &rust_zero, |sc| {
                let first_token_id = managed_token_id!(WREWA_TOKEN_ID);
                let second_token_id = managed_token_id!(MOA_TOKEN_ID);
                let router_address = managed_address!(&owner_addr);
                let router_owner_address = managed_address!(&owner_addr);
                let total_fee_percent = 300u64;
                let special_fee_percent = 50u64;

                sc.init(
                    first_token_id,
                    second_token_id,
                    router_address,
                    router_owner_address,
                    total_fee_percent,
                    special_fee_percent,
                    ManagedAddress::<DebugApi>::zero(),
                    MultiValueEncoded::<DebugApi, ManagedAddress<DebugApi>>::new(),
                );

                let lp_token_id = managed_token_id!(LPMOA_TOKEN_ID);
                sc.lp_token_identifier().set(&lp_token_id);

                sc.state().set(State::Active);
            })
            .assert_ok();

        blockchain_wrapper
            .execute_tx(&owner_addr, &usdc_pair_wrapper, &rust_zero, |sc| {
                let first_token_id = managed_token_id!(WREWA_TOKEN_ID);
                let second_token_id = managed_token_id!(USDC_TOKEN_ID);
                let router_address = managed_address!(&owner_addr);
                let router_owner_address = managed_address!(&owner_addr);
                let total_fee_percent = 300u64;
                let special_fee_percent = 50u64;

                sc.init(
                    first_token_id,
                    second_token_id,
                    router_address,
                    router_owner_address,
                    total_fee_percent,
                    special_fee_percent,
                    ManagedAddress::<DebugApi>::zero(),
                    MultiValueEncoded::<DebugApi, ManagedAddress<DebugApi>>::new(),
                );

                let lp_token_id = managed_token_id!(LPUSDC_TOKEN_ID);
                sc.lp_token_identifier().set(&lp_token_id);

                sc.state().set(State::Active);
            })
            .assert_ok();

        blockchain_wrapper
            .execute_tx(&owner_addr, &router_wrapper, &rust_zero, |sc| {
                sc.init(OptionalValue::None);

                sc.pair_map().insert(
                    PairTokens {
                        first_token_id: managed_token_id!(WREWA_TOKEN_ID),
                        second_token_id: managed_token_id!(MOA_TOKEN_ID),
                    },
                    managed_address!(moa_pair_wrapper.address_ref()),
                );
                sc.pair_map().insert(
                    PairTokens {
                        first_token_id: managed_token_id!(WREWA_TOKEN_ID),
                        second_token_id: managed_token_id!(USDC_TOKEN_ID),
                    },
                    managed_address!(usdc_pair_wrapper.address_ref()),
                );
            })
            .assert_ok();

        let lp_token_roles = [DcdtLocalRole::Mint, DcdtLocalRole::Burn];
        blockchain_wrapper.set_dcdt_local_roles(
            moa_pair_wrapper.address_ref(),
            LPMOA_TOKEN_ID,
            &lp_token_roles[..],
        );

        let lp_token_roles = [DcdtLocalRole::Mint, DcdtLocalRole::Burn];
        blockchain_wrapper.set_dcdt_local_roles(
            usdc_pair_wrapper.address_ref(),
            LPUSDC_TOKEN_ID,
            &lp_token_roles[..],
        );

        let user_addr = blockchain_wrapper.create_user_account(&rust_biguint!(100_000_000));
        blockchain_wrapper.set_dcdt_balance(
            &user_addr,
            WREWA_TOKEN_ID,
            &rust_biguint!(USER_TOTAL_WREWA_TOKENS),
        );
        blockchain_wrapper.set_dcdt_balance(
            &user_addr,
            MOA_TOKEN_ID,
            &rust_biguint!(USER_TOTAL_MOA_TOKENS),
        );
        blockchain_wrapper.set_dcdt_balance(
            &user_addr,
            USDC_TOKEN_ID,
            &rust_biguint!(USER_TOTAL_USDC_TOKENS),
        );

        RouterSetup {
            blockchain_wrapper,
            owner_address: owner_addr,
            user_address: user_addr,
            router_wrapper,
            moa_pair_wrapper,
            usdc_pair_wrapper,
        }
    }

    pub fn add_liquidity(&mut self) {
        let payments = vec![
            TxTokenTransfer {
                token_identifier: WREWA_TOKEN_ID.to_vec(),
                nonce: 0,
                value: rust_biguint!(ADD_LIQUIDITY_TOKENS),
            },
            TxTokenTransfer {
                token_identifier: MOA_TOKEN_ID.to_vec(),
                nonce: 0,
                value: rust_biguint!(ADD_LIQUIDITY_TOKENS),
            },
        ];

        self.blockchain_wrapper
            .execute_dcdt_multi_transfer(
                &self.user_address,
                &self.moa_pair_wrapper,
                &payments,
                |sc| {
                    sc.add_liquidity(
                        managed_biguint!(ADD_LIQUIDITY_TOKENS),
                        managed_biguint!(ADD_LIQUIDITY_TOKENS),
                    );
                },
            )
            .assert_ok();

        let payments = vec![
            TxTokenTransfer {
                token_identifier: WREWA_TOKEN_ID.to_vec(),
                nonce: 0,
                value: rust_biguint!(ADD_LIQUIDITY_TOKENS),
            },
            TxTokenTransfer {
                token_identifier: USDC_TOKEN_ID.to_vec(),
                nonce: 0,
                value: rust_biguint!(ADD_LIQUIDITY_TOKENS),
            },
        ];

        self.blockchain_wrapper
            .execute_dcdt_multi_transfer(
                &self.user_address,
                &self.usdc_pair_wrapper,
                &payments,
                |sc| {
                    sc.add_liquidity(
                        managed_biguint!(ADD_LIQUIDITY_TOKENS),
                        managed_biguint!(ADD_LIQUIDITY_TOKENS),
                    );
                },
            )
            .assert_ok();
    }

    pub fn multi_pair_swap(
        &mut self,
        payment_token: &[u8],
        payment_amount: u64,
        args: &[(Address, &[u8], &[u8], u64)],
    ) {
        let payment_amount_big = rust_biguint!(payment_amount);

        self.blockchain_wrapper
            .execute_dcdt_transfer(
                &self.user_address,
                &self.router_wrapper,
                payment_token,
                0,
                &payment_amount_big,
                |sc| {
                    let mut swap_operations = MultiValueEncoded::new();
                    for x in args.iter() {
                        swap_operations.push(MultiValue4::from((
                            managed_address!(&x.0),
                            managed_buffer!(x.1),
                            managed_token_id!(x.2.to_owned()),
                            managed_biguint!(x.3),
                        )));
                    }

                    sc.multi_pair_swap(swap_operations);
                },
            )
            .assert_ok();
    }
}
