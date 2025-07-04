from typing import Any

from dharitri_sdk_cli.localnet import genesis
from dharitri_sdk_cli.localnet.config_root import ConfigRoot


def patch(data: Any, config: ConfigRoot):
    owner = genesis.get_owner_of_genesis_contracts()

    delegation_config = data[0]
    dns_config = data[1]

    delegation_config["owner"] = owner.address.bech32()
    dns_config["owner"] = owner.address.bech32()
    # registration price = 100 atto-REWA
    dns_config["init-parameters"] = "0064"
