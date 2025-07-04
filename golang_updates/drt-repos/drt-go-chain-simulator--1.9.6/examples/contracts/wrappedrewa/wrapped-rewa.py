import json
import sys
import time
from pathlib import Path

from dharitri_sdk import (Address, ProxyNetworkProvider,
                            SmartContractTransactionsFactory, Token,
                            TransactionsFactoryConfig, UserSecretKey)

SIMULATOR_URL = "http://localhost:8085"
GENERATE_BLOCKS_UNTIL_EPOCH_REACHED_URL = "simulator/generate-blocks-until-epoch-reached"
GENERATE_BLOCKS_UNTIL_TX_PROCESSED = "simulator/generate-blocks-until-transaction-processed"
WRAPPED_REWA_TOKEN = "WREWA-bd4d79"

parent_directory = Path(__file__).parent

def main():
    # create a network provider
    provider = ProxyNetworkProvider(SIMULATOR_URL)

    key = UserSecretKey.generate()
    address = key.generate_public_key().to_address("drt")
    print(f"working with the generated address: {address.to_bech32()}")

    # call proxy faucet
    data = {"receiver": f"{address.to_bech32()}"}
    provider.do_post_generic("transaction/send-user-funds", data)

    # generate blocks until DCDTs are enabled
    provider.do_post_generic(f"{GENERATE_BLOCKS_UNTIL_EPOCH_REACHED_URL}/1", {})

    # set state for wrapped rewa contract and system account on shard 1
    # load JSON data from the file
    with open(parent_directory / "accounts-state.json", "r") as file:
        json_data = json.load(file)

    provider.do_post_generic("simulator/set-state", json_data)

    wrapp_contract_address = Address.new_from_bech32(
        "drt1qqqqqqqqqqqqqpgqhe8t5jewej70zupmh44jurgn29psua5l2jpsv0u330"
    )

    config = TransactionsFactoryConfig(provider.get_network_config().chain_id)
    sc_factory = SmartContractTransactionsFactory(config)
    amount_rewa = 5000000000000000000
    call_transaction = sc_factory.create_transaction_for_execute(
        sender=address,
        contract=wrapp_contract_address,
        function="wrapRewa",
        gas_limit=10000000,
        native_transfer_amount=amount_rewa,
        arguments=[]
    )

    call_transaction.nonce = provider.get_account(address).nonce
    call_transaction.signature = b"dummy"

    # send transaction
    tx_hash = provider.send_transaction(call_transaction)
    print(f"wrapp rewa tx hash: {tx_hash.hex()}")

    time.sleep(0.05)
    # generate enough blocks until the transaction is completed
    provider.do_post_generic(f"{GENERATE_BLOCKS_UNTIL_TX_PROCESSED}/{tx_hash.hex()}", {})

    token = provider.get_token_of_account(address, Token(WRAPPED_REWA_TOKEN))
    if token.amount != amount_rewa:
        sys.exit(f"amount of token from balance is no equal with the initial supply: "
                 f"actual-{token.amount}, expected-{amount_rewa}")

    print("transaction was executed, initial address received the wrapped rewa token")


if __name__ == "__main__":
    main()
