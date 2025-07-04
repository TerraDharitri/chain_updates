import requests
import json

from config import *
from network_provider.get_transaction_info import get_status_of_tx
from constants import *
import time
from core.validatorKey import ValidatorKey


def send_rewa_to_address(rewa_amount, drt_address):
    details = {
        'address': f'{drt_address}',
        'balance': f'{rewa_amount}'
    }

    details_list = [details]
    json_structure = json.dumps(details_list)
    response = requests.post(f"{DEFAULT_PROXY}/simulator/set-state", data=json_structure)
    response.raise_for_status()

    return response.text


def add_blocks(nr_of_blocks):
    response = requests.post(f"{DEFAULT_PROXY}/simulator/generate-blocks/{nr_of_blocks}")
    response.raise_for_status()
    return response.text


def get_block() -> int:
    response = requests.get(f"{DEFAULT_PROXY}/network/status/0")
    parsed = response.json()

    general_data = parsed.get("data")
    general_status = general_data.get("status")
    nonce = general_status.get("drt_nonce")
    return nonce


def add_blocks_until_epoch_reached(epoch_to_be_reached: int):
    req = requests.post(f"{DEFAULT_PROXY}/simulator/generate-blocks-until-epoch-reached/{str(epoch_to_be_reached)}")
    add_blocks(1)
    return req.text


def add_blocks_until_tx_fully_executed(tx_hash) -> str:
    print("Checking: ", tx_hash)
    counter = 0

    while counter < MAX_NUM_OF_BLOCKS_UNTIL_TX_SHOULD_BE_EXECUTED:
        add_blocks(1)

        time.sleep(WAIT_UNTIL_API_REQUEST_IN_SEC)
        if get_status_of_tx(tx_hash) == "pending":
            counter += 1
        else:
            print("Tx fully executed after", counter, " blocks.")
            return get_status_of_tx(tx_hash)


def is_chain_online() -> bool:
    flag = False

    while not flag:
        time.sleep(1)
        try:
            response = requests.get(f"{DEFAULT_PROXY}/network/status/0")
            print(response)
            flag = True
        except requests.exceptions.ConnectionError:
            print("Chain not started jet")

    return flag


def add_key(keys: list[ValidatorKey]) -> str:
    private_keys = []
    for key in keys:
        private_keys.append(key.get_private_key())

    post_body = {
        "privateKeysBase64": private_keys
    }

    json_structure = json.dumps(post_body)
    req = requests.post(f"{DEFAULT_PROXY}/simulator/add-keys", data=json_structure)

    return req.text


def add_blocks_until_key_eligible(keys: list[ValidatorKey]) -> ValidatorKey:
    flag = False
    while not flag:
        for key in keys:
            if key.get_state() == "eligible":
                eligible_key = key
                print("eligible key found")
                flag = True

            else:
                print("no eligible key found , moving to next epoch...")
                current_epoch = proxy_default.get_network_status().epoch_number
                add_blocks_until_epoch_reached(current_epoch+1)
                add_blocks(3)

    return eligible_key


def add_blocks_until_last_block_of_current_epoch() -> str:
    response = requests.get(f"{DEFAULT_PROXY}/network/status/4294967295")
    response.raise_for_status()
    parsed = response.json()

    general_data = parsed.get("data")
    status = general_data.get("status")
    passed_nonces = status.get("drt_nonces_passed_in_current_epoch")

    blocks_to_be_added = rounds_per_epoch - passed_nonces
    response_from_add_blocks = add_blocks(blocks_to_be_added)
    return response_from_add_blocks

