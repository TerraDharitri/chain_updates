import requests
from config import DEFAULT_PROXY
from helpers import base64_to_hex
from dharitri_sdk_core import Address


def get_delegation_contract_address_from_tx(tx_hash):
    response = requests.get(f"{DEFAULT_PROXY}/transaction/{tx_hash}?withResults=True")
    response.raise_for_status()
    parsed = response.json()

    general_data = parsed.get("data")
    transaction_data = general_data.get("transaction")
    logs_data = transaction_data.get("logs")
    events_data = logs_data.get("events")
    first_set_of_events = events_data[0]
    topics = first_set_of_events.get("topics")
    delegation_contract_address = topics[1]

    delegation_contract_address = base64_to_hex(delegation_contract_address)
    delegation_contract_address = Address.from_hex(delegation_contract_address, "drt").to_bech32()

    return delegation_contract_address
