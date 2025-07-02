import base64

from dharitri_sdk.core.address import Address
from dharitri_sdk.core.transaction_on_network import TransactionEvent
from dharitri_sdk.governance.governance_transactions_outcome_parser import (
    GovernanceTransactionsOutcomeParser,
)
from dharitri_sdk.testutils.mock_transaction_on_network import (
    get_empty_transaction_on_network,
)


class TestGovernanceOutcomeParser:
    parser = GovernanceTransactionsOutcomeParser()

    def test_parse_propose_proposal(self):
        transaction = get_empty_transaction_on_network()
        transaction.logs.events = [
            TransactionEvent(
                raw={},
                address=Address.empty(),
                identifier="proposal",
                topics=[b"\x01", b"1db734c0315f9ec422b88f679ccfe3e0197b9d67", b"5", b"7"],
                data=b"",
                additional_data=[],
            )
        ]

        outcome = self.parser.parse_new_proposal(transaction)
        assert len(outcome) == 1
        assert outcome[0].proposal_nonce == 1
        assert outcome[0].commit_hash == "1db734c0315f9ec422b88f679ccfe3e0197b9d67"
        assert outcome[0].start_vote_epoch == 53
        assert outcome[0].end_vote_epoch == 55

    def test_parse_vote(self):
        transaction = get_empty_transaction_on_network()
        transaction.logs.events = [
            TransactionEvent(
                raw={},
                address=Address.empty(),
                identifier="vote",
                topics=[
                    base64.b64decode("AQ=="),
                    base64.b64decode("eWVz"),
                    base64.b64decode("BlpNol0wFsAAAA=="),
                    base64.b64decode("BlpNol0wFsAAAA=="),
                ],
                data=b"",
                additional_data=[],
            )
        ]

        outcome = self.parser.parse_vote(transaction)
        assert len(outcome) == 1
        assert outcome[0].proposal_nonce == 1
        assert outcome[0].vote == "yes"
        assert outcome[0].total_stake == 30000_000000000000000000
        assert outcome[0].total_voting_power == 30000_000000000000000000

    def test_parse_delegate_vote(self):
        transaction = get_empty_transaction_on_network()
        transaction.logs.events = [
            TransactionEvent(
                raw={},
                address=Address.empty(),
                identifier="delegateVote",
                topics=[
                    base64.b64decode("AQ=="),
                    base64.b64decode("YWJzdGFpbg=="),
                    base64.b64decode("a3Qc0P1f8raaWzOVkcJbHHxHOx2+LI6S8CM9aV+W6KY="),
                    base64.b64decode("Ah4Z4Mm6skAAAA=="),
                    base64.b64decode("Ah4Z4Mm6skAAAA=="),
                ],
                data=b"",
                additional_data=[],
            )
        ]

        outcome = self.parser.parse_delegate_vote(transaction)
        assert len(outcome) == 1
        assert outcome[0].proposal_nonce == 1
        assert outcome[0].vote == "abstain"
        assert outcome[0].voter == Address.new_from_bech32(
            "drt1dd6pe58atletdxjmxw2ersjmr37ywwcahckgayhsyv7kjhukaznqmkvprh"
        )
        assert outcome[0].user_stake == 10000_000000000000000000
        assert outcome[0].voting_power == 10000_000000000000000000

    def test_parse_close_proposal(self):
        transaction = get_empty_transaction_on_network()
        transaction.logs.events = [
            TransactionEvent(
                raw={},
                address=Address.empty(),
                identifier="closeProposal",
                topics=[
                    base64.b64decode("ZDVkMjRhYTY1ZWY5OWM3NDcxMjkxMmZkOGJiMmE1MDVjY2RmMDYyYw=="),
                    base64.b64decode("dHJ1ZQ=="),
                ],
                data=b"",
                additional_data=[],
            )
        ]

        outcome = self.parser.parse_close_proposal(transaction)
        assert len(outcome) == 1
        assert outcome[0].commit_hash == "d5d24aa65ef99c74712912fd8bb2a505ccdf062c"
        assert outcome[0].passed
