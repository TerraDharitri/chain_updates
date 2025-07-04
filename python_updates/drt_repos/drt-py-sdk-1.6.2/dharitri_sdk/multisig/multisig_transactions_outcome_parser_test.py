import base64
from pathlib import Path

from dharitri_sdk.abi.abi import Abi
from dharitri_sdk.core.address import Address
from dharitri_sdk.core.transaction_on_network import TransactionEvent
from dharitri_sdk.multisig.multisig_transactions_outcome_parser import (
    MultisigTransactionsOutcomeParser,
)
from dharitri_sdk.network_providers.api_network_provider import ApiNetworkProvider
from dharitri_sdk.testutils.mock_transaction_on_network import (
    get_empty_transaction_on_network,
)
from dharitri_sdk.testutils.utils import create_network_providers_config


class TestMultisigTransactionsOutcomeParser:
    testdata = Path(__file__).parent.parent / "testutils" / "testdata"
    abi = Abi.load(testdata / "multisig-full.abi.json")
    parser = MultisigTransactionsOutcomeParser(abi)
    network_provider = ApiNetworkProvider(
        url="https://devnet-api.dharitri.org", config=create_network_providers_config()
    )

    def test_parse_deploy(self):
        tx_on_network = get_empty_transaction_on_network()
        event = TransactionEvent(
            raw={},
            identifier="SCDeploy",
            address=Address.new_from_bech32("drt1qqqqqqqqqqqqqpgqe832k3l6d02ww7l9cvqum25539nmmdxa9ncssqu3lh"),
            topics=[
                base64.b64decode("AAAAAAAAAAAFAMniq0f6a9Tne+XDAc2qlIlnvbTdLPE="),
                base64.b64decode("P7gfQwO+b3N3NQuKWV+UsT/Wy85MTH0sY+nh+PDULPE="),
                base64.b64decode("JPGpNEze5Sasa1UvkfIqvzapnqikyDefUsZcf44gxIY="),
            ],
            data=b"",
            additional_data=[],
        )
        tx_on_network.logs.events = [event]

        outcome = self.parser.parse_deploy(tx_on_network)
        assert len(outcome.contracts) == 1
        assert (
            outcome.contracts[0].address.to_bech32() == "drt1qqqqqqqqqqqqqpgqe832k3l6d02ww7l9cvqum25539nmmdxa9ncssqu3lh"
        )
