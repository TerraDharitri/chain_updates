import pytest

from dharitri_sdk.core.address import Address
from dharitri_sdk.core.transaction import Transaction
from dharitri_sdk.core.transaction_computer import TransactionComputer
from dharitri_sdk.network_providers.account_awaiter import AccountAwaiter
from dharitri_sdk.network_providers.api_network_provider import ApiNetworkProvider
from dharitri_sdk.network_providers.errors import (
    ExpectedAccountConditionNotReachedError,
)
from dharitri_sdk.network_providers.resources import AccountOnNetwork
from dharitri_sdk.testutils.mock_network_provider import (
    MockNetworkProvider,
    TimelinePointMarkCompleted,
    TimelinePointWait,
)
from dharitri_sdk.testutils.utils import create_account_rewa_balance
from dharitri_sdk.testutils.wallets import load_wallets


class TestAccountAwaiter:
    provider = MockNetworkProvider()
    watcher = AccountAwaiter(
        fetcher=provider,
        polling_interval_in_milliseconds=42,
        timeout_interval_in_milliseconds=42 * 42,
        patience_time_in_milliseconds=0,
    )

    def test_await_on_balance_increase(self):
        alice = Address.new_from_bech32("drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf")
        # alice account is created with 1000 REWA
        initial_balance = self.provider.get_account(alice).balance

        # adds 7 REWA to the account balance
        self.provider.mock_account_balance_timeline_by_address(
            alice,
            [
                TimelinePointWait(40),
                TimelinePointWait(40),
                TimelinePointWait(45),
                TimelinePointMarkCompleted(),
            ],
        )

        def condition(account: AccountOnNetwork):
            return account.balance == initial_balance + create_account_rewa_balance(7)

        account = self.watcher.await_on_condition(alice, condition)
        assert account.balance == create_account_rewa_balance(1007)

    @pytest.mark.networkInteraction
    def test_await_for_account_balance_increase_on_network(self):
        alice = load_wallets()["alice"]
        alice_address = Address.new_from_bech32(alice.label)
        frank = Address.new_from_bech32("drt1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asacg58j")

        api = ApiNetworkProvider("https://devnet-api.dharitri.org")
        watcher = AccountAwaiter(fetcher=api)
        tx_computer = TransactionComputer()
        value = 100_000

        transaction = Transaction(
            sender=alice_address,
            receiver=frank,
            gas_limit=50000,
            chain_id="D",
            value=value,
        )
        transaction.nonce = api.get_account(alice_address).nonce
        transaction.signature = alice.secret_key.sign(tx_computer.compute_bytes_for_signing(transaction))

        initial_balance = api.get_account(frank).balance

        def condition(account: AccountOnNetwork):
            return account.balance == initial_balance + value

        api.send_transaction(transaction)

        account_on_network = watcher.await_on_condition(frank, condition)
        assert account_on_network.balance == initial_balance + value

    @pytest.mark.networkInteraction
    def test_ensure_error_if_timeout(self):
        alice = load_wallets()["alice"]
        alice_address = Address.new_from_bech32(alice.label)
        bob = Address.new_from_bech32("drt1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqlqde3c")

        api = ApiNetworkProvider("https://devnet-api.dharitri.org")
        watcher = AccountAwaiter(
            fetcher=api,
            polling_interval_in_milliseconds=1000,
            timeout_interval_in_milliseconds=10000,
        )

        value = 100_000
        transaction = Transaction(
            sender=alice_address,
            receiver=bob,
            gas_limit=50000,
            chain_id="D",
            value=value,
        )
        transaction.nonce = api.get_account(alice_address).nonce

        tx_computer = TransactionComputer()
        transaction.signature = alice.secret_key.sign(tx_computer.compute_bytes_for_signing(transaction))

        initial_balance = api.get_account(bob).balance

        def condition(account: AccountOnNetwork):
            return account.balance == initial_balance + value

        api.send_transaction(transaction)

        with pytest.raises(ExpectedAccountConditionNotReachedError):
            watcher.await_on_condition(bob, condition)
