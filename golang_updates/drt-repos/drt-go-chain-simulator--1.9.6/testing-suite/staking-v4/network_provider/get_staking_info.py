import requests
import json
from dharitri_sdk_core import Address
from constants import VALIDATOR_CONTRACT
from config import DEFAULT_PROXY
from helpers import base64_to_string


def get_total_staked(owner: str):
    address_in_hex = Address.from_bech32(owner).to_hex()
    post_body = {
        "scAddress": VALIDATOR_CONTRACT,
        "funcName": "getTotalStaked",
        "args": [address_in_hex]
    }

    json_structure = json.dumps(post_body)
    response = requests.post(f"{DEFAULT_PROXY}/vm-values/query", data=json_structure)
    response.raise_for_status()
    parsed = response.json()

    general_data = parsed.get("data")
    tx_response_data = general_data.get("data")
    total_staked_list = tx_response_data.get("returnData")
    total_staked = total_staked_list[0]

    total_staked = base64_to_string(total_staked)
    return total_staked
