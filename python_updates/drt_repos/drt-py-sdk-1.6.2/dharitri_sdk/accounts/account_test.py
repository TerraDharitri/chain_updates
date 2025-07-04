from pathlib import Path

from dharitri_sdk.accounts.account import Account
from dharitri_sdk.core.address import Address
from dharitri_sdk.core.message import Message
from dharitri_sdk.core.transaction import Transaction
from dharitri_sdk.wallet.keypair import KeyPair
from dharitri_sdk.wallet.user_keys import UserSecretKey

testwallets = Path(__file__).parent.parent / "testutils" / "testwallets"
DUMMY_MNEMONIC = "moral volcano peasant pass circle pen over picture flat shop clap goat never lyrics gather prepare woman film husband gravity behind test tiger improve"
alice = testwallets / "alice.pem"


def test_create_account_from_pem():
    account = Account.new_from_pem(alice)

    assert account.secret_key.get_bytes().hex() == "413f42575f7f26fad3317a778771212fdb80245850981e48b58a4f25e344e8f9"
    assert account.address.to_bech32() == "drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf"


def test_create_account_from_keystore():
    account = Account.new_from_keystore(testwallets / "withDummyMnemonic.json", "password")

    assert account.secret_key.get_bytes().hex() == "413f42575f7f26fad3317a778771212fdb80245850981e48b58a4f25e344e8f9"
    assert account.address.to_bech32() == "drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf"


def test_create_account_from_mnemonic():
    account = Account.new_from_mnemonic(DUMMY_MNEMONIC)

    assert account.secret_key.get_bytes().hex() == "413f42575f7f26fad3317a778771212fdb80245850981e48b58a4f25e344e8f9"
    assert account.address.to_bech32() == "drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf"


def test_create_account_from_keypair():
    secret_key = UserSecretKey.new_from_string("413f42575f7f26fad3317a778771212fdb80245850981e48b58a4f25e344e8f9")
    keypair = KeyPair(secret_key)
    account = Account.new_from_keypair(keypair)

    assert account.secret_key == secret_key
    assert account.address.to_bech32() == "drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf"


def test_account_nonce_holder():
    account = Account.new_from_pem(alice)
    account.nonce = 42
    assert account.get_nonce_then_increment() == 42
    assert account.get_nonce_then_increment() == 43

    account.get_nonce_then_increment()
    account.get_nonce_then_increment()
    account.get_nonce_then_increment()
    assert account.nonce == 47


def test_sign_transaction():
    """
    Also see: https://github.com/TerraDharitri/drt-go-chain/blob/main/examples/construction_test.go
    """

    tx = Transaction(
        nonce=89,
        value=0,
        receiver=Address.new_from_bech32("drt1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqlqde3c"),
        sender=Address.new_from_bech32("drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf"),
        data=None,
        gas_price=1000000000,
        gas_limit=50000,
        chain_id="local-testnet",
        version=1,
        options=0,
    )

    account = Account.new_from_pem(alice)
    tx.signature = account.sign_transaction(tx)
    assert (
        tx.signature.hex()
        == "b56769014f2bdc5cf9fc4a05356807d71fcf8775c819b0f1b0964625b679c918ffa64862313bfef86f99b38cb84fcdb16fa33ad6eb565276616723405cd8f109"
    )


def test_sign_message():
    message = Message(
        "hello".encode(),
        address=Address.new_from_bech32("drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf"),
    )

    account = Account.new_from_pem(alice)
    message.signature = account.sign_message(message)
    assert (
        message.signature.hex()
        == "561bc58f1dc6b10de208b2d2c22c9a474ea5e8cabb59c3d3ce06bbda21cc46454aa71a85d5a60442bd7784effa2e062fcb8fb421c521f898abf7f5ec165e5d0f"
    )


def test_sign_tx_by_hash():
    account = Account.new_from_pem(alice)

    tx = Transaction(
        sender=Address.new_from_bech32("drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf"),
        receiver=Address.new_from_bech32("drt1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqlqde3c"),
        value=0,
        gas_limit=50000,
        version=2,
        options=1,
        chain_id="integration tests chain ID",
        nonce=89,
    )

    tx.signature = account.sign_transaction(tx)

    assert (
        tx.signature.hex()
        == "f0c81f2393b1ec5972c813f817bae8daa00ade91c6f75ea604ab6a4d2797aca4378d783023ff98f1a02717fe4f24240cdfba0b674ee9abb18042203d713bc70a"
    )
