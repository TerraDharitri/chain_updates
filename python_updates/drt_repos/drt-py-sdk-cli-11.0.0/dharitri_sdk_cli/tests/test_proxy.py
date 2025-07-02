from dharitri_sdk import Address, ProxyNetworkProvider

from dharitri_sdk_cli.cli import main
from dharitri_sdk_cli.config import get_config_for_network_providers


def test_sync_nonce():
    account = Address.new_from_bech32("drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf")
    config = get_config_for_network_providers()
    proxy = ProxyNetworkProvider("https://testnet-api.dharitri.org", config=config)
    nonce = proxy.get_account(account).nonce
    assert True if nonce else False


def test_query_contract():
    result = main(
        [
            "contract",
            "query",
            "drt1qqqqqqqqqqqqqpgq6qr0w0zzyysklfneh32eqp2cf383zc89d8ssk0pue3",
            "--function",
            "getSum",
            "--proxy",
            "https://devnet-api.dharitri.org",
        ]
    )
    assert False if result else True
