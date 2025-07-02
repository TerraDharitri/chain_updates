import pytest

from dharitri_sdk.network_providers.config import NetworkProviderConfig
from dharitri_sdk.network_providers.proxy_network_provider import ProxyNetworkProvider
from dharitri_sdk.smart_contracts.smart_contract_transactions_outcome_parser import (
    SmartContractTransactionsOutcomeParser,
)


@pytest.mark.networkInteraction
class TestSmartContractDeployDevnet:
    parser = SmartContractTransactionsOutcomeParser()
    network_config = NetworkProviderConfig(client_name="drt-py-sdk/tests")
    provider = ProxyNetworkProvider(url="https://devnet-gateway.dharitri.org", config=network_config)

    def test_parse_deploy_transaction_1(self):
        tx_hash = "5d2ff2af8eb3fe7f2acb7e29c0436854b4c6c44de02878b6afff582888024a55"
        transaction = self.provider.get_transaction(tx_hash)
        parsed = self.parser.parse_deploy(transaction)

        assert parsed.return_code == "ok"
        assert len(parsed.contracts) == 1

        contract = parsed.contracts[0]
        assert contract.address.to_bech32() == "drt1qqqqqqqqqqqqqpgqpayq2es08gq8798xhnpr0kzgn7495qt5q6uqszgzdh"
        assert contract.owner_address.to_bech32() == "drt1tn62hjp72rznp8vq0lplva5csav6rccpqqdungpxtqz0g2hcq6uqc2zmmy"
        assert contract.code_hash == bytes.fromhex("c876625ec34a04445cfd99067777ebe488afdbc6899cd958f4c1d36107ca02d9")

    def test_parse_deploy_transaction_2(self):
        tx_hash = "76683e926dad142fc9651afca208487f2a80d327fc87e5c876eec9d028196352"
        transaction = self.provider.get_transaction(tx_hash)
        parsed = self.parser.parse_deploy(transaction)

        assert parsed.return_code == "execution failed"
        assert len(parsed.contracts) == 0
