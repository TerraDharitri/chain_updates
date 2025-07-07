use forwarder_legacy::fwd_nft_legacy::{Color, ForwarderNftModule};

use dharitri_sc_scenario::imports::*;

const USER_ADDRESS: TestAddress = TestAddress::new("user");
const FORWARDER_ADDRESS: TestSCAddress = TestSCAddress::new("forwarder_legacy");
const FORWARDER_PATH: DrtscPath = DrtscPath::new("output/forwarder_legacy.drtsc.json");
const NFT_TOKEN_ID: TestTokenIdentifier = TestTokenIdentifier::new("COOL-123456");

fn world() -> ScenarioWorld {
    let mut blockchain = ScenarioWorld::new();

    blockchain.set_current_dir_from_workspace("contracts/feature-tests/composability");
    blockchain.register_contract(FORWARDER_PATH, forwarder_legacy::ContractBuilder);
    blockchain
}

#[test]
fn test_nft_update_attributes_and_send() {
    let mut world = world();

    let roles = vec![
        "DCDTRoleNFTCreate".to_string(),
        "DCDTRoleNFTUpdateAttributes".to_string(),
    ];

    world.account(USER_ADDRESS).nonce(1);
    world
        .account(FORWARDER_ADDRESS)
        .nonce(1)
        .code(FORWARDER_PATH)
        .dcdt_roles(NFT_TOKEN_ID, roles);

    let original_attributes = Color { r: 0, g: 0, b: 0 };

    world
        .tx()
        .from(USER_ADDRESS)
        .to(FORWARDER_ADDRESS)
        .whitebox(forwarder_legacy::contract_obj, |sc| {
            sc.nft_create_compact(
                NFT_TOKEN_ID.to_token_identifier(),
                managed_biguint!(1),
                original_attributes,
            );

            sc.tx()
                .to(USER_ADDRESS)
                .dcdt((NFT_TOKEN_ID.to_token_identifier(), 1, 1u32.into()))
                .transfer();
        });

    world
        .check_account(USER_ADDRESS)
        .dcdt_nft_balance_and_attributes(NFT_TOKEN_ID, 1, 1, original_attributes);

    let new_attributes = Color {
        r: 255,
        g: 255,
        b: 255,
    };

    world
        .tx()
        .from(USER_ADDRESS)
        .to(FORWARDER_ADDRESS)
        .payment(TestDcdtTransfer(NFT_TOKEN_ID, 1, 1))
        .whitebox(forwarder_legacy::contract_obj, |sc| {
            sc.nft_update_attributes(NFT_TOKEN_ID.to_token_identifier(), 1, new_attributes);

            sc.tx()
                .to(USER_ADDRESS)
                .dcdt((NFT_TOKEN_ID.to_token_identifier(), 1, 1u32.into()))
                .transfer();
        });

    world
        .check_account(USER_ADDRESS)
        .dcdt_nft_balance_and_attributes(NFT_TOKEN_ID, 1, 1, new_attributes);
}
