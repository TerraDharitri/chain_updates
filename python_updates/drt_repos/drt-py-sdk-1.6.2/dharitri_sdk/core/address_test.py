import pytest

from dharitri_sdk.core.address import (
    Address,
    AddressComputer,
    AddressFactory,
    is_valid_bech32,
)
from dharitri_sdk.core.config import LibraryConfig
from dharitri_sdk.core.errors import BadAddressError, BadPubkeyLengthError


def test_address():
    address = Address.new_from_bech32("drt1l453hd0gt5gzdp7czpuall8ggt2dcv5zwmfdf3sd3lguxseux2fsxvluwu")
    assert str(address) == "drt1l453hd0gt5gzdp7czpuall8ggt2dcv5zwmfdf3sd3lguxseux2fsxvluwu"
    assert address == Address.new_from_bech32("drt1l453hd0gt5gzdp7czpuall8ggt2dcv5zwmfdf3sd3lguxseux2fsxvluwu")
    assert not address == Address(
        bytes.fromhex("fd691bb5e85d102687d81079dffce842d4dc328276d2d4c60d8fd1c3433c3293"),
        "test",
    )
    assert "fd691bb5e85d102687d81079dffce842d4dc328276d2d4c60d8fd1c3433c3293" == address.to_hex()
    assert "drt1l453hd0gt5gzdp7czpuall8ggt2dcv5zwmfdf3sd3lguxseux2fsxvluwu" == address.to_bech32()

    address = Address.new_from_hex("fd691bb5e85d102687d81079dffce842d4dc328276d2d4c60d8fd1c3433c3293", "drt")
    assert "fd691bb5e85d102687d81079dffce842d4dc328276d2d4c60d8fd1c3433c3293" == address.to_hex()
    assert "drt1l453hd0gt5gzdp7czpuall8ggt2dcv5zwmfdf3sd3lguxseux2fsxvluwu" == address.to_bech32()

    with pytest.raises(BadPubkeyLengthError):
        address = Address(bytes.fromhex("fd691bb5e85d102687d8"), "drt")

    with pytest.raises(BadAddressError):
        address = Address.new_from_bech32("bad")

    address = Address.empty()
    assert address.pubkey == bytes()
    assert address.hrp == "drt"


def test_address_with_custom_hrp():
    address = Address.new_from_hex("0139472eff6886771a982f3083da5d421f24c29181e63888228dc81ca60d69e1", "test")
    assert address.hrp == "test"
    assert address.to_bech32() == "test1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ss5hqhtr"

    address = Address.new_from_bech32("test1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ss5hqhtr")
    assert address.hrp == "test"
    assert address.to_hex() == "0139472eff6886771a982f3083da5d421f24c29181e63888228dc81ca60d69e1"


def test_address_factory():
    factory_foo = AddressFactory("foo")
    factory_drt = AddressFactory("drt")
    pubkey = bytes.fromhex("0139472eff6886771a982f3083da5d421f24c29181e63888228dc81ca60d69e1")

    assert (
        factory_foo.create_from_public_key(pubkey).to_bech32()
        == "foo1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssunhpj4"
    )
    assert (
        factory_drt.create_from_public_key(pubkey).to_bech32()
        == "drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf"
    )


def test_is_valid_bech32():
    assert is_valid_bech32("drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf", "drt")
    assert is_valid_bech32("foo1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssunhpj4", "foo")
    assert not is_valid_bech32("foobar", "foo")
    assert not is_valid_bech32("drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf", "foo")


def test_get_address_shard():
    address_computer = AddressComputer()
    address = Address.new_from_bech32("drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf")
    assert address_computer.get_shard_of_address(address) == 1

    address = Address.new_from_bech32("drt1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqlqde3c")
    assert address_computer.get_shard_of_address(address) == 0

    address = Address.new_from_bech32("drt1k2s324ww2g0yj38qn2ch2jwctdy8mnfxep94q9arncc6xecg3xaq889n6e")
    assert address_computer.get_shard_of_address(address) == 2


def test_compute_contract_address():
    deployer = Address.new_from_bech32("drt1j0hxzs7dcyxw08c4k2nv9tfcaxmqy8rj59meq505w92064x0h40q96qj7l")
    address_computer = AddressComputer()

    contract_address = address_computer.compute_contract_address(deployer, deployment_nonce=0)
    assert contract_address.to_hex() == "00000000000000000500bb652200ed1f994200ab6699462cab4b1af7b11ebd5e"
    assert contract_address.to_bech32() == "drt1qqqqqqqqqqqqqpgqhdjjyq8dr7v5yq9tv6v5vt9tfvd00vg7h40q8zfxpd"

    contract_address = address_computer.compute_contract_address(deployer, deployment_nonce=1)
    assert contract_address.to_hex() == "000000000000000005006e4f90488e27342f9a46e1809452c85ee7186566bd5e"
    assert contract_address.to_bech32() == "drt1qqqqqqqqqqqqqpgqde8eqjywyu6zlxjxuxqfg5kgtmn3setxh40qy0s6t6"


def test_address_with_library_config_hrp():
    address = Address(bytes.fromhex("0139472eff6886771a982f3083da5d421f24c29181e63888228dc81ca60d69e1"))
    assert address.to_bech32() == "drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf"

    LibraryConfig.default_address_hrp = "test"
    address = Address(bytes.fromhex("0139472eff6886771a982f3083da5d421f24c29181e63888228dc81ca60d69e1"))
    assert address.to_bech32() == "test1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ss5hqhtr"
    LibraryConfig.default_address_hrp = "drt"
